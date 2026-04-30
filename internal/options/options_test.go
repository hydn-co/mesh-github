package options

import (
	"testing"

	"github.com/hydn-co/mesh-sdk/pkg/testkit"
)

func TestShouldRegisterPolymorphicOptions(t *testing.T) {
	testkit.TestPolymorphicRegistrations(t, map[string]any{
		"mesh://github/collectors/member_entity_collector_options":      &GitHubMemberEntityCollectorOptions{},
		"mesh://github/collectors/team_entity_collector_options":        &GitHubTeamEntityCollectorOptions{},
		"mesh://github/collectors/repository_entity_collector_options":  &GitHubRepositoryEntityCollectorOptions{},
		"mesh://github/collectors/audit_log_activity_collector_options": &GitHubAuditLogActivityCollectorOptions{},
		"mesh://github/actions/add_team_member_action_options":          &GitHubAddTeamMemberActionOptions{},
		"mesh://github/actions/remove_team_member_action_options":       &GitHubRemoveTeamMemberActionOptions{},
	})
}
