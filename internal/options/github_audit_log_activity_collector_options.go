package options

import "github.com/hydn-co/mesh-sdk/pkg/catalog/spaces"

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
	return githubRequirements()
}
