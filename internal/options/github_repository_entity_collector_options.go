package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type GitHubRepositoryEntityCollectorOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubRepositoryEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://github/collectors/repository_entity_collector_options"
}

func (o *GitHubRepositoryEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Applications, spaces.ApplicationAccounts}
}

func (o *GitHubRepositoryEntityCollectorOptions) GetRequirements() []string {
	return githubRequirements()
}
