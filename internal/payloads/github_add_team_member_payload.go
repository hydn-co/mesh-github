package payloads

import (
	"fmt"
	"strings"

	githubroles "github.com/hydn-co/mesh-github/internal/roles"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
)

type GitHubAddTeamMemberPayload struct {
	TeamSlug string `json:"team_slug"      binding:"required" title:"Team"     description:"GitHub team to add the member to"      x-lookup:"{\"entity-type\": \"groups\", \"display-key\": \"name\", \"submit-key\": \"group_ref\", \"form-input-type\": \"select\" }"`
	Username string `json:"username"       binding:"required" title:"Username" description:"GitHub username of the member to add"`
	Role     string `json:"role,omitempty"                    title:"Role"     description:"GitHub team membership role to assign" x-lookup:"{\"entity-type\": \"roles\", \"display-key\": \"name\", \"submit-key\": \"role_ref\", \"form-input-type\": \"select\" }"`
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

	if p.Role != "" && !githubroles.IsTeamMembershipRole(p.Role) {
		return fmt.Errorf("role must be one of %s", quoteAndJoin(githubroles.TeamMembershipRoleRefs()))
	}

	return nil
}

func quoteAndJoin(values []string) string {
	quoted := make([]string, 0, len(values))
	for _, value := range values {
		quoted = append(quoted, fmt.Sprintf("'%s'", value))
	}

	if len(quoted) == 0 {
		return ""
	}

	if len(quoted) == 1 {
		return quoted[0]
	}

	return strings.Join(quoted[:len(quoted)-1], ", ") + " or " + quoted[len(quoted)-1]
}
