package actions

import (
	"context"
	"fmt"

	"github.com/hydn-co/mesh-github/internal/api"
	"github.com/hydn-co/mesh-github/internal/credentials"
	"github.com/hydn-co/mesh-github/internal/options"
	"github.com/hydn-co/mesh-github/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

// GitHubRemoveTeamMemberAction removes a member from a GitHub team.
type GitHubRemoveTeamMemberAction struct {
	*connector.TypedFeatureContext[*options.GitHubRemoveTeamMemberActionOptions, *payloads.GitHubRemoveTeamMemberPayload]
	client *api.Client
}

func NewGitHubRemoveTeamMemberAction(
	ctx *connector.TypedFeatureContext[*options.GitHubRemoveTeamMemberActionOptions, *payloads.GitHubRemoveTeamMemberPayload],
) runner.Feature {
	return &GitHubRemoveTeamMemberAction{TypedFeatureContext: ctx}
}

func (a *GitHubRemoveTeamMemberAction) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	payload := a.GetPayload()
	if payload == nil {
		return fmt.Errorf("remove team member payload is required")
	}
	if err := payload.Validate(); err != nil {
		return err
	}

	token, err := credentials.ExtractToken(a.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract token: %w", err)
	}

	opts := a.GetOptions()
	a.client = api.NewClient(token, opts.Organization)
	return nil
}

func (a *GitHubRemoveTeamMemberAction) Start(ctx context.Context) error {
	logAction(ctx, "github_remove_team_member_action")

	payload := a.GetPayload()
	if err := a.client.RemoveTeamMember(ctx, payload.TeamSlug, payload.Username); err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}

	return nil
}

func (a *GitHubRemoveTeamMemberAction) Stop(_ context.Context) error { return nil }
