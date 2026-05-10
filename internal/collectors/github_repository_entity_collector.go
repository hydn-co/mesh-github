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

// GitHubRepositoryEntityCollector collects repositories and collaborator links from GitHub.
type GitHubRepositoryEntityCollector struct {
	*connector.TypedFeatureContext[*options.GitHubRepositoryEntityCollectorOptions, *connector.NoPayload]
	client *api.Client
	state  connectorutil.FeatureState
}

func NewGitHubRepositoryEntityCollector(
	ctx *connector.TypedFeatureContext[*options.GitHubRepositoryEntityCollectorOptions, *connector.NoPayload],
) runner.Feature {
	return &GitHubRepositoryEntityCollector{TypedFeatureContext: ctx}
}

func (c *GitHubRepositoryEntityCollector) Init(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := connectorutil.Validate(c.GetOptions(), "github repository entity collector options"); err != nil {
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

func (c *GitHubRepositoryEntityCollector) Start(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	if err := c.state.RequireReady(); err != nil {
		return err
	}

	connectorutil.LogFeature(ctx, c.TypedFeatureContext, slog.LevelInfo, "Starting GitHub repository entity collector")

	owner := c.GetOptions().Organization

	if err := enumerators.ForEach(c.client.RepositoryEnumerator(ctx), func(repository api.Repository) error {
		if err := ctx.Err(); err != nil {
			return err
		}

		applicationRef := repository.FullName
		if applicationRef == "" {
			applicationRef = repository.Name
		}
		if applicationRef == "" {
			return fmt.Errorf("repository missing reference")
		}

		application := &entities.Application{
			ApplicationRef: applicationRef,
			Name:           repository.Name,
			Description:    repository.Description,
		}

		if err := c.Emit(ctx, application); err != nil {
			return fmt.Errorf("failed to emit repository %s: %w", applicationRef, err)
		}

		if err := enumerators.ForEach(
			c.client.RepositoryCollaboratorEnumerator(ctx, owner, repository.Name),
			func(collaborator api.RepositoryCollaborator) error {
				if err := ctx.Err(); err != nil {
					return err
				}

				applicationAccount := &entities.ApplicationAccount{
					ApplicationRef: applicationRef,
					AccountRef:     collaborator.Login,
				}

				if err := c.Emit(ctx, applicationAccount); err != nil {
					return fmt.Errorf(
						"failed to emit repository collaborator %s/%s: %w",
						applicationRef,
						collaborator.Login,
						err,
					)
				}

				return nil
			},
		); err != nil {
			return fmt.Errorf("failed to enumerate collaborators for repository %s: %w", applicationRef, err)
		}

		return nil
	}); err != nil {
		return fmt.Errorf("failed to enumerate repositories: %w", err)
	}

	return nil
}

func (c *GitHubRepositoryEntityCollector) Stop(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	c.state.Reset()
	c.client = nil
	return nil
}
