package payloads

import (
	"fmt"

	"github.com/fgrzl/json/polymorphic"
)

func init() {
	polymorphic.RegisterType[GitHubAddTeamMemberPayload]()
	polymorphic.RegisterType[GitHubRemoveTeamMemberPayload]()
}

// GitHubAddTeamMemberPayload is the input payload for adding a member to a team.
type GitHubAddTeamMemberPayload struct {
	TeamSlug string `json:"team_slug"      binding:"required"`
	Username string `json:"username"       binding:"required"`
	Role     string `json:"role,omitempty"`
}

func (p *GitHubAddTeamMemberPayload) GetDiscriminator() string {
	return "mesh://github/actions/add_team_member_payload"
}

func (p *GitHubAddTeamMemberPayload) Validate() error {
	if p.TeamSlug == "" {
		return fmt.Errorf("team_slug is required")
	}
	if p.Username == "" {
		return fmt.Errorf("username is required")
	}
	if p.Role != "" && p.Role != "member" && p.Role != "maintainer" {
		return fmt.Errorf("role must be 'member' or 'maintainer'")
	}
	return nil
}

// GitHubRemoveTeamMemberPayload is the input payload for removing a member from a team.
type GitHubRemoveTeamMemberPayload struct {
	TeamSlug string `json:"team_slug" binding:"required"`
	Username string `json:"username"  binding:"required"`
}

func (p *GitHubRemoveTeamMemberPayload) GetDiscriminator() string {
	return "mesh://github/actions/remove_team_member_payload"
}

func (p *GitHubRemoveTeamMemberPayload) Validate() error {
	if p.TeamSlug == "" {
		return fmt.Errorf("team_slug is required")
	}
	if p.Username == "" {
		return fmt.Errorf("username is required")
	}
	return nil
}
