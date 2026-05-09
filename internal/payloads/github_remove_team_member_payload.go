package payloads

import (
	"fmt"

	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
)

type GitHubRemoveTeamMemberPayload struct {
	TeamSlug string `json:"team_slug" binding:"required" x-lookup:"{\"entity-type\": \"groups\", \"display-key\": \"name\", \"submit-key\": \"group_ref\", \"form-input-type\": \"select\" }"`
	Username string `json:"username"  binding:"required"`
}

func (p *GitHubRemoveTeamMemberPayload) GetDiscriminator() string {
	return "mesh://github/actions/remove_team_member_payload"
}

func (p *GitHubRemoveTeamMemberPayload) Validate() error {
	if connectorutil.IsNil(p) {
		return fmt.Errorf("github remove team member payload is required")
	}

	return connectorutil.RequireStrings(
		"github remove team member payload",
		connectorutil.RequiredString{Name: "team_slug", Value: p.TeamSlug},
		connectorutil.RequiredString{Name: "username", Value: p.Username},
	)
}
