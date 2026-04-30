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

// GitHubAddTeamMemberAction adds a member to a GitHub team.
type GitHubAddTeamMemberAction struct {
	*connector.TypedFeatureContext[*options.GitHubAddTeamMemberActionOptions, *payloads.GitHubAddTeamMemberPayload]
	client *api.Client
}

func NewGitHubAddTeamMemberAction(
	ctx *connector.TypedFeatureContext[*options.GitHubAddTeamMemberActionOptions, *payloads.GitHubAddTeamMemberPayload],
) runner.Feature {
	return &GitHubAddTeamMemberAction{TypedFeatureContext: ctx}
}

func (a *GitHubAddTeamMemberAction) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	payload := a.GetPayload()
	if payload == nil {
		return fmt.Errorf("add team member payload is required")
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

func (a *GitHubAddTeamMemberAction) Start(ctx context.Context) error {
	logAction(ctx, "github_add_team_member_action")

	payload := a.GetPayload()
	role := payload.Role
	if role == "" {
		role = "member"
	}

	if err := a.client.AddTeamMember(ctx, payload.TeamSlug, payload.Username, role); err != nil {
		return fmt.Errorf("failed to add team member: %w", err)
	}

	return nil
}

func (a *GitHubAddTeamMemberAction) Stop(_ context.Context) error { return nil }
