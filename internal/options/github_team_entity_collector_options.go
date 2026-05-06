package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type GitHubTeamEntityCollectorOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubTeamEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://github/collectors/team_entity_collector_options"
}

func (o *GitHubTeamEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Groups, spaces.GroupMembers}
}

func (o *GitHubTeamEntityCollectorOptions) GetRequirements() []string {
	return githubRequirements()
}
