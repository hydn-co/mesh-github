package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type GitHubRemoveTeamMemberActionOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubRemoveTeamMemberActionOptions) GetDiscriminator() string {
	return "mesh://github/actions/remove_team_member_action_options"
}

func (o *GitHubRemoveTeamMemberActionOptions) GetSpaces() []spaces.Space {
	return groupMembershipSpaces()
}

func (o *GitHubRemoveTeamMemberActionOptions) GetRequirements() []string {
	return githubRequirements()
}
