package collectors

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/fgrzl/enumerators"
	"github.com/hydn-co/mesh-github/internal/api"
	"github.com/hydn-co/mesh-github/internal/credentials"
	"github.com/hydn-co/mesh-github/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/events"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/types"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

const auditLogUnavailableHint = "GitHub returns 404 when the organization audit log is unavailable for the org plan or the token lacks organization-owner/read:audit_log access"

// GitHubAuditLogActivityCollector collects audit log events from GitHub.
type GitHubAuditLogActivityCollector struct {
	*connector.TypedFeatureContext[*options.GitHubAuditLogActivityCollectorOptions, *connector.NoPayload]
	client *api.Client
	state  connectorutil.FeatureState
}

func NewGitHubAuditLogActivityCollector(
	ctx *connector.TypedFeatureContext[*options.GitHubAuditLogActivityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &GitHubAuditLogActivityCollector{TypedFeatureContext: ctx}
}

func (c *GitHubAuditLogActivityCollector) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := connectorutil.Validate(c.GetOptions(), "github audit log activity collector options"); err != nil {
		return err
	}

	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract token: %w", err)
	}

	opts := c.GetOptions()
	c.client = api.NewClient(token, opts.Organization)
	c.state.MarkReady()
	return nil
}

func (c *GitHubAuditLogActivityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := c.state.RequireReady(); err != nil {
		return err
	}

	connectorutil.LogFeature(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting GitHub audit log activity collector")

	// Extract resume cursor from payload
	var since time.Time
	if c.Configuration != nil && c.Payload != nil && c.Payload.Content != nil {
		switch e := c.Payload.Content.(type) {
		case *events.GroupMemberAdded:
			since = e.Timestamp
		case *events.GroupMemberRemoved:
			since = e.Timestamp
		case *events.PermissionGranted:
			since = e.Timestamp
		case *events.PermissionRevoked:
			since = e.Timestamp
		case *events.ResourceAccessed:
			since = e.Timestamp
		}
	}

	if err := enumerators.ForEach(c.client.AuditLogEnumerator(ctx, "", since), func(entry api.AuditLogEntry) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		event := mapAuditLogEntry(entry)
		if event == nil {
			return nil
		}

		if err := c.Emit(ctx, event); err != nil {
			return fmt.Errorf("failed to emit audit event %s: %w", entry.DocumentID, err)
		}

		return nil
	}); err != nil {
		if api.IsAuditLogUnavailable(err) {
			slog.WarnContext(
				ctx,
				"GitHub audit log endpoint unavailable; skipping collector",
				"organization", c.GetOptions().Organization,
				"error", err,
				"hint", auditLogUnavailableHint,
			)
			return newAuditLogUnavailableError(err)
		}

		return fmt.Errorf("failed to enumerate audit log: %w", err)
	}

	return nil
}

func (c *GitHubAuditLogActivityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.state.Reset()
	c.client = nil
	return nil
}

func newAuditLogUnavailableError(err error) error {
	return fmt.Errorf("failed to enumerate audit log: %w; hint: %s", err, auditLogUnavailableHint)
}

func mapAuditLogEntry(entry api.AuditLogEntry) any {
	ts := time.UnixMilli(entry.Timestamp).UTC()
	if ts.IsZero() {
		ts = time.UnixMilli(entry.CreatedAt).UTC()
	}

	actor := types.Actor{
		Ref:  entry.Actor,
		Type: "account",
	}

	switch {
	case strings.HasPrefix(entry.Action, "team.add_member"):
		return &events.GroupMemberAdded{
			EventRef:  entry.DocumentID,
			Timestamp: ts,
			Actor:     actor,
			Target: types.Target{
				Ref:  entry.User,
				Type: "account",
			},
			Outcome: types.EventOutcome{
				Action: "add_member",
				Result: "success",
			},
			GroupRef:  entry.Team,
			GroupName: entry.Team,
		}

	case strings.HasPrefix(entry.Action, "team.remove_member"):
		return &events.GroupMemberRemoved{
			EventRef:  entry.DocumentID,
			Timestamp: ts,
			Actor:     actor,
			Target: types.Target{
				Ref:  entry.User,
				Type: "account",
			},
			Outcome: types.EventOutcome{
				Action: "remove_member",
				Result: "success",
			},
			GroupRef:  entry.Team,
			GroupName: entry.Team,
		}

	case strings.HasPrefix(entry.Action, "repo.access"):
		return &events.ResourceAccessed{
			EventRef:  entry.DocumentID,
			Timestamp: ts,
			Actor:     actor,
			Target: types.Target{
				Ref:  entry.Repo,
				Type: "application",
			},
			Outcome: types.EventOutcome{
				Action: "access",
				Result: "success",
			},
			ResourceType: "repository",
		}

	case strings.HasPrefix(entry.Action, "org.add_member"):
		return &events.GroupMemberAdded{
			EventRef:  entry.DocumentID,
			Timestamp: ts,
			Actor:     actor,
			Target: types.Target{
				Ref:  entry.User,
				Type: "account",
			},
			Outcome: types.EventOutcome{
				Action: "add_member",
				Result: "success",
			},
			GroupRef:  entry.Org,
			GroupName: entry.Org,
			GroupType: "organization",
		}

	case strings.HasPrefix(entry.Action, "org.remove_member"):
		return &events.GroupMemberRemoved{
			EventRef:  entry.DocumentID,
			Timestamp: ts,
			Actor:     actor,
			Target: types.Target{
				Ref:  entry.User,
				Type: "account",
			},
			Outcome: types.EventOutcome{
				Action: "remove_member",
				Result: "success",
			},
			GroupRef:  entry.Org,
			GroupName: entry.Org,
			GroupType: "organization",
		}

	default:
		// Emit generic resource accessed for other audit events
		return &events.ResourceAccessed{
			EventRef:  entry.DocumentID,
			Timestamp: ts,
			Actor:     actor,
			Target: types.Target{
				Ref:  entry.Repo,
				Type: "resource",
			},
			Outcome: types.EventOutcome{
				Action: entry.Action,
				Result: "success",
			},
			ResourceType: "audit_event",
			AccessType:   entry.Action,
		}
	}
}
