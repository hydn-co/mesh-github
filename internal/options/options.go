package options

import (
	"github.com/fgrzl/json/polymorphic"
	"github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"
)

func init() {
	polymorphic.RegisterType[GitHubMemberEntityCollectorOptions]()
	polymorphic.RegisterType[GitHubTeamEntityCollectorOptions]()
	polymorphic.RegisterType[GitHubRepositoryEntityCollectorOptions]()
	polymorphic.RegisterType[GitHubAuditLogActivityCollectorOptions]()
	polymorphic.RegisterType[GitHubAddTeamMemberActionOptions]()
	polymorphic.RegisterType[GitHubRemoveTeamMemberActionOptions]()
}

// GitHubOptionsCore contains shared options for GitHub collectors.
type GitHubOptionsCore struct {
	Organization string `json:"organization" title:"Organization" description:"GitHub organization login name" binding:"required"`
}

// GitHubMemberEntityCollectorOptions configures the member entity collector.
type GitHubMemberEntityCollectorOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubMemberEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://github/collectors/member_entity_collector_options"
}

func (o *GitHubMemberEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Accounts}
}

func (o *GitHubMemberEntityCollectorOptions) GetRequirements() []string {
	return []string{"github"}
}

// GitHubTeamEntityCollectorOptions configures the team entity collector.
type GitHubTeamEntityCollectorOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubTeamEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://github/collectors/team_entity_collector_options"
}

func (o *GitHubTeamEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Groups, spaces.GroupMembers}
}

func (o *GitHubTeamEntityCollectorOptions) GetRequirements() []string {
	return []string{"github"}
}

// GitHubRepositoryEntityCollectorOptions configures the repository entity collector.
type GitHubRepositoryEntityCollectorOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubRepositoryEntityCollectorOptions) GetDiscriminator() string {
	return "mesh://github/collectors/repository_entity_collector_options"
}

func (o *GitHubRepositoryEntityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Applications, spaces.ApplicationAccounts}
}

func (o *GitHubRepositoryEntityCollectorOptions) GetRequirements() []string {
	return []string{"github"}
}

// GitHubAuditLogActivityCollectorOptions configures the audit log activity collector.
type GitHubAuditLogActivityCollectorOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubAuditLogActivityCollectorOptions) GetDiscriminator() string {
	return "mesh://github/collectors/audit_log_activity_collector_options"
}

func (o *GitHubAuditLogActivityCollectorOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.Activity}
}

func (o *GitHubAuditLogActivityCollectorOptions) GetRequirements() []string {
	return []string{"github"}
}

// GitHubAddTeamMemberActionOptions configures the add team member action.
type GitHubAddTeamMemberActionOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubAddTeamMemberActionOptions) GetDiscriminator() string {
	return "mesh://github/actions/add_team_member_action_options"
}

func (o *GitHubAddTeamMemberActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.GroupMembers}
}

func (o *GitHubAddTeamMemberActionOptions) GetRequirements() []string {
	return []string{"github"}
}

// GitHubRemoveTeamMemberActionOptions configures the remove team member action.
type GitHubRemoveTeamMemberActionOptions struct {
	GitHubOptionsCore `json:",inline"`
}

func (o *GitHubRemoveTeamMemberActionOptions) GetDiscriminator() string {
	return "mesh://github/actions/remove_team_member_action_options"
}

func (o *GitHubRemoveTeamMemberActionOptions) GetSpaces() []spaces.Space {
	return []spaces.Space{spaces.GroupMembers}
}

func (o *GitHubRemoveTeamMemberActionOptions) GetRequirements() []string {
	return []string{"github"}
}
