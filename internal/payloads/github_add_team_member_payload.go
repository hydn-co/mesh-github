package payloads

import (
	"fmt"

	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
)

type GitHubAddTeamMemberPayload struct {
	TeamSlug string `json:"team_slug"      binding:"required"`
	Username string `json:"username"       binding:"required"`
	Role     string `json:"role,omitempty"`
}

func (p *GitHubAddTeamMemberPayload) GetDiscriminator() string {
	return "mesh://github/actions/add_team_member_payload"
}

func (p *GitHubAddTeamMemberPayload) Validate() error {
	if connectorutil.IsNil(p) {
		return fmt.Errorf("github add team member payload is required")
	}

	if err := connectorutil.RequireStrings(
		"github add team member payload",
		connectorutil.RequiredString{Name: "team_slug", Value: p.TeamSlug},
		connectorutil.RequiredString{Name: "username", Value: p.Username},
	); err != nil {
		return err
	}

	if p.Role != "" && p.Role != "member" && p.Role != "maintainer" {
		return fmt.Errorf("role must be 'member' or 'maintainer'")
	}

	return nil
}
