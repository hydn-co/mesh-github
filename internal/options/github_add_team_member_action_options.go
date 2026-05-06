package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

type GitHubAddTeamMemberActionOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubAddTeamMemberActionOptions) GetDiscriminator() string {
	return "mesh://github/actions/add_team_member_action_options"
}

func (o *GitHubAddTeamMemberActionOptions) GetSpaces() []spaces.Space {
	return groupMembershipSpaces()
}

func (o *GitHubAddTeamMemberActionOptions) GetRequirements() []string {
	return githubRequirements()
}
