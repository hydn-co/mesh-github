package payloads

import (
	"testing"

	"github.com/hydn-co/mesh-sdk/pkg/testkit"
)

func TestShouldRegisterPolymorphicPayloads(t *testing.T) {
	testkit.TestPolymorphicRegistrations(t, map[string]any{
		"mesh://github/actions/add_team_member_payload":    &GitHubAddTeamMemberPayload{},
		"mesh://github/actions/remove_team_member_payload": &GitHubRemoveTeamMemberPayload{},
	})
}
