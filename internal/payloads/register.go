package payloads

import "github.com/fgrzl/json/polymorphic"

func init() {
	polymorphic.RegisterType[GitHubAddTeamMemberPayload]()
	polymorphic.RegisterType[GitHubRemoveTeamMemberPayload]()
}
