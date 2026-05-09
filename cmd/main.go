package main

import (
	"log"

	"github.com/hydn-co/mesh-github/internal/actions"
	"github.com/hydn-co/mesh-github/internal/collectors"
	"github.com/hydn-co/mesh-github/internal/options"
	"github.com/hydn-co/mesh-github/internal/payloads"
	"github.com/hydn-co/mesh-sdk/pkg/connector"
	"github.com/hydn-co/mesh-sdk/pkg/runner"
)

func main() {
	runner.Run(WithManifest())
}

func WithManifest() *runner.Manifest {
	manifest := runner.CreateManifest(
		"mesh-github",
		"",
		"GitHub",
		"Mesh integration with GitHub",
	)

	// Entity Collectors
	manifest.MustRegisterFeature(
		"collect_members",
		"GitHub Member Entity Collector",
		"Collects organization members from GitHub and emits them as account entities.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.GitHubMemberEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewGitHubMemberEntityCollector),
	)

	manifest.MustRegisterFeature(
		"collect_teams",
		"GitHub Team Entity Collector",
		"Collects teams and team memberships from GitHub organizations.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.GitHubTeamEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewGitHubTeamEntityCollector),
	)

	// Activity Collectors
	manifest.MustRegisterFeature(
		"collect_audit_log",
		"GitHub Audit Log Activity Collector",
		"Collects audit log events from GitHub organizations.",
		runner.FeatureSchedulable,
		runner.FeatureTypeCollector,
		new(options.GitHubAuditLogActivityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorLastActivity,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewGitHubAuditLogActivityCollector),
	)

	// Actions
	manifest.MustRegisterFeature(
		"add_team_member",
		"GitHub Add Team Member Action",
		"Adds a member to a GitHub team.",
		runner.FeatureUnschedulable,
		runner.FeatureTypeAction,
		new(options.GitHubAddTeamMemberActionOptions),
		new(payloads.GitHubAddTeamMemberPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewGitHubAddTeamMemberAction),
	)

	manifest.MustRegisterFeature(
		"remove_team_member",
		"GitHub Remove Team Member Action",
		"Removes a member from a GitHub team.",
		runner.FeatureUnschedulable,
		runner.FeatureTypeAction,
		new(options.GitHubRemoveTeamMemberActionOptions),
		new(payloads.GitHubRemoveTeamMemberPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewGitHubRemoveTeamMemberAction),
	)

	if err := manifest.Validate(); err != nil {
		log.Fatal(err)
	}

	return manifest
}
