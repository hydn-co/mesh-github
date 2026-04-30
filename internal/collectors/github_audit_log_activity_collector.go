package collectors

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/hydn-co/mesh-github/internal/api"
	"github.com/hydn-co/mesh-github/internal/credentials"
	"github.com/hydn-co/mesh-github/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/events"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/types"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

// GitHubAuditLogActivityCollector collects audit log events from GitHub.
type GitHubAuditLogActivityCollector struct {
	*connector.TypedFeatureContext[*options.GitHubAuditLogActivityCollectorOptions, *connector.NoPayload]
	client *api.Client
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

	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract token: %w", err)
	}

	opts := c.GetOptions()
	c.client = api.NewClient(token, opts.Organization)
	return nil
}

func (c *GitHubAuditLogActivityCollector) Start(ctx context.Context) error {
	logCollector(ctx, "github_audit_log_activity_collector")

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

	entries, err := c.client.ListAuditLog(ctx, "", since)
	if err != nil {
		return fmt.Errorf("failed to list audit log: %w", err)
	}

	for _, entry := range entries {
		if err := ctx.Err(); err != nil {
			return err
		}

		event := mapAuditLogEntry(entry)
		if event == nil {
			continue
		}

		if err := c.Emit(ctx, event); err != nil {
			return fmt.Errorf("failed to emit audit event %s: %w", entry.DocumentID, err)
		}
	}

	return nil
}

func (c *GitHubAuditLogActivityCollector) Stop(_ context.Context) error { return nil }

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
