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

// GitHubRemoveTeamMemberAction removes a member from a GitHub team.
type GitHubRemoveTeamMemberAction struct {
	*connector.TypedFeatureContext[*options.GitHubRemoveTeamMemberActionOptions, *payloads.GitHubRemoveTeamMemberPayload]
	client *api.Client
	state  connectorutil.FeatureState
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

	if err := connectorutil.Validate(a.GetOptions(), "github remove team member action options"); err != nil {
		return err
	}

	if err := connectorutil.Validate(a.GetPayload(), "github remove team member payload"); err != nil {
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

func (a *GitHubRemoveTeamMemberAction) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := a.state.RequireReady(); err != nil {
		return err
	}

	connectorutil.LogFeature(ctx, a.TypedFeatureContext, slog.LevelInfo, "Starting GitHub remove team member action")

	payload := a.GetPayload()
	if err := a.client.RemoveTeamMember(ctx, payload.TeamSlug, payload.Username); err != nil {
		return fmt.Errorf("failed to remove team member: %w", err)
	}

	return nil
}

func (a *GitHubRemoveTeamMemberAction) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	a.state.Reset()
	a.client = nil
	return nil
}
