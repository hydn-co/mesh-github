package options

import (
	"fmt"

	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
	"github.com/hydn-co/mesh-sdk/pkg/connectorutil"
)

type GitHubOptionsCore struct {
	Organization string `json:"organization" title:"Organization" description:"GitHub organization login name" binding:"required"`
}

func validateGitHubOptionsCore(core *GitHubOptionsCore) error {
	if connectorutil.IsNil(core) {
		return fmt.Errorf("github options required but not provided")
	}

	return connectorutil.RequireStrings(
		"github options",
		connectorutil.RequiredString{Name: "organization", Value: core.Organization},
	)
}

func githubRequirements() []string {
	return []string{"github"}
}

func groupMembershipSpaces() []spaces.Space {
	return []spaces.Space{spaces.GroupMembers}
}
