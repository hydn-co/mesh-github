package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type GitHubMemberEntityCollectorOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubMemberEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://github/collectors/member_entity_collector_options"
}

func (o *GitHubMemberEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Accounts}
}

func (o *GitHubMemberEntityCollectorOptions) GetRequirements() []string {
	return githubRequirements()
}
