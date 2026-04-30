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

// GitHubRepositoryEntityCollector collects repositories and collaborators from GitHub.
type GitHubRepositoryEntityCollector struct {
	*connector.TypedFeatureContext[*options.GitHubRepositoryEntityCollectorOptions, *connector.NoPayload]
	client *api.Client
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

	token, err := credentials.ExtractToken(c.GetCredentials())
	if err != nil {
		return fmt.Errorf("failed to extract token: %w", err)
	}

	opts := c.GetOptions()
	c.client = api.NewClient(token, opts.Organization)
	return nil
}

func (c *GitHubRepositoryEntityCollector) Start(ctx context.Context) error {
	logCollector(ctx, "github_repository_entity_collector")

	repos, err := c.client.ListRepositories(ctx)
	if err != nil {
		return fmt.Errorf("failed to list repositories: %w", err)
	}

	for _, repo := range repos {
		if err := ctx.Err(); err != nil {
			return err
		}

		app := &entities.Application{
			ApplicationRef: repo.FullName,
			Name:           repo.Name,
			Description:    repo.Description,
		}

		if err := c.Emit(ctx, app); err != nil {
			return fmt.Errorf("failed to emit repository %s: %w", repo.FullName, err)
		}

		// Collect collaborators
		collabs, err := c.client.ListCollaborators(ctx, repo.Name)
		if err != nil {
			return fmt.Errorf("failed to list collaborators for %s: %w", repo.Name, err)
		}

		for _, collab := range collabs {
			if err := ctx.Err(); err != nil {
				return err
			}

			appAccount := &entities.ApplicationAccount{
				ApplicationRef: repo.FullName,
				AccountRef:     collab.Login,
			}

			if err := c.Emit(ctx, appAccount); err != nil {
				return fmt.Errorf("failed to emit collaborator %s/%s: %w", repo.FullName, collab.Login, err)
			}
		}
	}

	return nil
}

func (c *GitHubRepositoryEntityCollector) Stop(_ context.Context) error { return nil }
