package roles

import "github.com/hydn-co/mesh-sdk/pkg/catalog/entities"

const (
	TeamMembershipRoleMember     = "member"
	TeamMembershipRoleMaintainer = "maintainer"
)

type teamMembershipRoleDefinition struct {
	ref         string
	name        string
	description string
}

var teamMembershipRoleDefinitions = []teamMembershipRoleDefinition{
	{
		ref:         TeamMembershipRoleMember,
		name:        "Member",
		description: "GitHub team member role",
	},
	{
		ref:         TeamMembershipRoleMaintainer,
		name:        "Maintainer",
		description: "GitHub team maintainer role",
	},
}

func IsTeamMembershipRole(role string) bool {
	for _, definition := range teamMembershipRoleDefinitions {
		if definition.ref == role {
			return true
		}
	}

	return false
}

func TeamMembershipRoleRefs() []string {
	refs := make([]string, 0, len(teamMembershipRoleDefinitions))
	for _, definition := range teamMembershipRoleDefinitions {
		refs = append(refs, definition.ref)
	}

	return refs
}

func TeamMembershipCatalogRoles() []*entities.Role {
	roles := make([]*entities.Role, 0, len(teamMembershipRoleDefinitions))
	for _, definition := range teamMembershipRoleDefinitions {
		roles = append(roles, &entities.Role{
			RoleRef:     definition.ref,
			Name:        definition.name,
			Description: definition.description,
		})
	}

	return roles
}
