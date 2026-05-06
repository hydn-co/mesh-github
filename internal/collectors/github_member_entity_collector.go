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
	"github.com/hydn-co/mesh-sdk/pkg/catalog/types"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

// GitHubMemberEntityCollector collects organization members from GitHub.
type GitHubMemberEntityCollector struct {
	*connector.TypedFeatureContext[*options.GitHubMemberEntityCollectorOptions, *connector.NoPayload]
	client *api.Client
	state  connectorutil.FeatureState
}

func NewGitHubMemberEntityCollector(
	ctx *connector.TypedFeatureContext[*options.GitHubMemberEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &GitHubMemberEntityCollector{TypedFeatureContext: ctx}
}

func (c *GitHubMemberEntityCollector) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := connectorutil.Validate(c.GetOptions(), "github member entity collector options"); err != nil {
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

func (c *GitHubMemberEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := c.state.RequireReady(); err != nil {
		return err
	}

	connectorutil.LogFeature(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting GitHub member entity collector")

	if err := enumerators.ForEach(c.client.MemberEnumerator(ctx), func(member api.Member) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		// Fetch detailed user info for name/email
		user, err := c.client.GetUser(ctx, member.Login)
		if err != nil {
			return fmt.Errorf("failed to get user details for %s: %w", member.Login, err)
		}

		account := &entities.Account{
			AccountRef:  member.Login,
			AccountType: toAccountType(member.Type),
			Name:        member.Login,
			DisplayName: user.Name,
			Enabled:     true,
		}

		if user.Email != "" {
			account.PrimaryEmail = &types.Email{
				Address: user.Email,
			}
		}

		if err := c.Emit(ctx, account); err != nil {
			return fmt.Errorf("failed to emit account %s: %w", member.Login, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to enumerate org members: %w", err)
	}

	return nil
}

func (c *GitHubMemberEntityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.state.Reset()
	c.client = nil
	return nil
}

func toAccountType(ghType string) types.AccountType {
	switch ghType {
	case "Bot":
		return types.AccountTypeServicePrincipal
	default:
		return types.AccountTypeUser
	}
}
