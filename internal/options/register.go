package options

import "github.com/fgrzl/json/polymorphic"

func init() {
	polymorphic.RegisterType[GitHubMemberEntityCollectorOptions]()
	polymorphic.RegisterType[GitHubRepositoryEntityCollectorOptions]()
	polymorphic.RegisterType[GitHubTeamEntityCollectorOptions]()
	polymorphic.RegisterType[GitHubAuditLogActivityCollectorOptions]()
	polymorphic.RegisterType[GitHubAddTeamMemberActionOptions]()
	polymorphic.RegisterType[GitHubRemoveTeamMemberActionOptions]()
}
