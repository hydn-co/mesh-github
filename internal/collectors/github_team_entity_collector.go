package collectors

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/fgrzl/enumerators"
	"github.com/hydn-co/mesh-github/internal/api"
	"github.com/hydn-co/mesh-github/internal/credentials"
	"github.com/hydn-co/mesh-github/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/entities"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

// GitHubTeamEntityCollector collects teams and team memberships from GitHub.
type GitHubTeamEntityCollector struct {
	*connector.TypedFeatureContext[*options.GitHubTeamEntityCollectorOptions, *connector.NoPayload]
	client *api.Client
	state  connectorutil.FeatureState
}

func NewGitHubTeamEntityCollector(
	ctx *connector.TypedFeatureContext[*options.GitHubTeamEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &GitHubTeamEntityCollector{TypedFeatureContext: ctx}
}

func (c *GitHubTeamEntityCollector) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := connectorutil.Validate(c.GetOptions(), "github team entity collector options"); err != nil {
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

func (c *GitHubTeamEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := c.state.RequireReady(); err != nil {
		return err
	}

	connectorutil.LogFeature(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting GitHub team entity collector")

	if err := enumerators.ForEach(c.client.TeamEnumerator(ctx), func(team api.Team) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		group := &entities.Group{
			GroupRef:    team.Slug,
			Name:        team.Name,
			Description: team.Description,
		}

		if err := c.Emit(ctx, group); err != nil {
			return fmt.Errorf("failed to emit team %s: %w", team.Slug, err)
		}

		if err := enumerators.ForEach(c.client.TeamMemberEnumerator(ctx, team.Slug), func(member api.TeamMember) error {
			if err := ctx.Err(); err != nil {
				return err
			}

			role, err := c.client.GetTeamMembershipRole(ctx, team.Slug, member.Login)
			if err != nil {
				return fmt.Errorf("failed to get team membership for %s/%s: %w", team.Slug, member.Login, err)
			}

			groupMember := &entities.GroupMember{
				GroupRef:   team.Slug,
				AccountRef: member.Login,
				RoleRef:    role,
			}

			if err := c.Emit(ctx, groupMember); err != nil {
				return fmt.Errorf("failed to emit team member %s/%s: %w", team.Slug, member.Login, err)
			}

			return nil
		}); err != nil {
			return fmt.Errorf("failed to enumerate members for team %s: %w", team.Slug, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to enumerate teams: %w", err)
	}

	return nil
}

func (c *GitHubTeamEntityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.state.Reset()
	c.client = nil
	return nil
}
