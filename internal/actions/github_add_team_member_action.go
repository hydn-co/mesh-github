package actions

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/hydn-co/mesh-github/internal/api"
	"github.com/hydn-co/mesh-github/internal/credentials"
	"github.com/hydn-co/mesh-github/internal/options"
	"github.com/hydn-co/mesh-github/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

// GitHubAddTeamMemberAction adds a member to a GitHub team.
type GitHubAddTeamMemberAction struct {
	*connector.TypedFeatureContext[*options.GitHubAddTeamMemberActionOptions, *payloads.GitHubAddTeamMemberPayload]
	client *api.Client
	state  connectorutil.FeatureState
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

	if err := connectorutil.Validate(a.GetOptions(), "github add team member action options"); err != nil {
		return err
	}

	if err := connectorutil.Validate(a.GetPayload(), "github add team member payload"); err != nil {
		return err
	}

	token, err := credentials.ExtractToken(a.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract token: %w", err)
	}

	opts := a.GetOptions()
	a.client = api.NewClient(token, opts.Organization)
	a.state.MarkReady()
	return nil
}

func (a *GitHubAddTeamMemberAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := a.state.RequireReady(); err != nil {
		return err
	}

	connectorutil.LogFeature(ctx, a.TypedFeatureContext, slog.LevelInfo, "Starting GitHub add team member action")

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

func (a *GitHubAddTeamMemberAction) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	a.state.Reset()
	a.client = nil
	return nil
}
