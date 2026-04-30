package collectors

import (
	"context"
	"fmt"

	"github.com/hydn-co/mesh-github/internal/api"
	"github.com/hydn-co/mesh-github/internal/credentials"
	"github.com/hydn-co/mesh-github/internal/options"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/entities"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

// GitHubTeamEntityCollector collects teams and team memberships from GitHub.
type GitHubTeamEntityCollector struct {
	*connector.TypedFeatureContext[*options.GitHubTeamEntityCollectorOptions, *connector.NoPayload]
	client *api.Client
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

	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract token: %w", err)
	}

	opts := c.GetOptions()
	c.client = api.NewClient(token, opts.Organization)
	return nil
}

func (c *GitHubTeamEntityCollector) Start(ctx context.Context) error {
	logCollector(ctx, "github_team_entity_collector")

	teams, err := c.client.ListTeams(ctx)
	if err != nil {
		return fmt.Errorf("failed to list teams: %w", err)
	}

	for _, team := range teams {
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

		// Collect team members
		members, err := c.client.ListTeamMembers(ctx, team.Slug)
		if err != nil {
			return fmt.Errorf("failed to list members for team %s: %w", team.Slug, err)
		}

		for _, member := range members {
			if err := ctx.Err(); err != nil {
				return err
			}

			groupMember := &entities.GroupMember{
				GroupRef:   team.Slug,
				AccountRef: member.Login,
				RoleRef:    member.Role,
			}

			if err := c.Emit(ctx, groupMember); err != nil {
				return fmt.Errorf("failed to emit team member %s/%s: %w", team.Slug, member.Login, err)
			}
		}
	}

	return nil
}

func (c *GitHubTeamEntityCollector) Stop(_ context.Context) error { return nil }
