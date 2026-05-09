package roles

import (
	"reflect"
	"testing"
)

func TestTeamMembershipRoleRefs(t *testing.T) {
	refs := TeamMembershipRoleRefs()
	want := []string{TeamMembershipRoleMember, TeamMembershipRoleMaintainer}

	if !reflect.DeepEqual(refs, want) {
		t.Fatalf("unexpected team membership role refs: got %v want %v", refs, want)
	}
}

func TestTeamMembershipCatalogRoles(t *testing.T) {
	roles := TeamMembershipCatalogRoles()

	if len(roles) != 2 {
		t.Fatalf("unexpected catalog role count: got %d want 2", len(roles))
	}

	if roles[0].RoleRef != TeamMembershipRoleMember || roles[0].Name != "Member" {
		t.Fatalf("unexpected member role: %+v", roles[0])
	}

	if roles[1].RoleRef != TeamMembershipRoleMaintainer || roles[1].Name != "Maintainer" {
		t.Fatalf("unexpected maintainer role: %+v", roles[1])
	}
}

func TestIsTeamMembershipRole(t *testing.T) {
	if !IsTeamMembershipRole(TeamMembershipRoleMember) {
		t.Fatal("expected member role to be valid")
	}

	if !IsTeamMembershipRole(TeamMembershipRoleMaintainer) {
		t.Fatal("expected maintainer role to be valid")
	}

	if IsTeamMembershipRole("owner") {
		t.Fatal("expected unsupported role to be invalid")
	}
}
