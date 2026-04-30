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
	if err := manifest.RegisterFeature(
		"github_member_entity_collector",
		"GitHub Member Entity Collector",
		"Collects organization members from GitHub and emits them as account entities.",
		true,
		runner.FeatureTypeCollector,
		new(options.GitHubMemberEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewGitHubMemberEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"github_team_entity_collector",
		"GitHub Team Entity Collector",
		"Collects teams and team memberships from GitHub organizations.",
		true,
		runner.FeatureTypeCollector,
		new(options.GitHubTeamEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewGitHubTeamEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"github_repository_entity_collector",
		"GitHub Repository Entity Collector",
		"Collects repositories and collaborators from GitHub organizations.",
		true,
		runner.FeatureTypeCollector,
		new(options.GitHubRepositoryEntityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewGitHubRepositoryEntityCollector),
	); err != nil {
		log.Fatal(err)
	}

	// Activity Collectors
	if err := manifest.RegisterFeature(
		"github_audit_log_activity_collector",
		"GitHub Audit Log Activity Collector",
		"Collects audit log events from GitHub organizations.",
		true,
		runner.FeatureTypeCollector,
		new(options.GitHubAuditLogActivityCollectorOptions),
		(*connector.NoPayload)(nil),
		runner.FeatureResumeBehaviorLastActivity,
		runner.APIKeyCredential,
		runner.Factory(collectors.NewGitHubAuditLogActivityCollector),
	); err != nil {
		log.Fatal(err)
	}

	// Actions
	if err := manifest.RegisterFeature(
		"github_add_team_member_action",
		"GitHub Add Team Member Action",
		"Adds a member to a GitHub team.",
		false,
		runner.FeatureTypeAction,
		new(options.GitHubAddTeamMemberActionOptions),
		new(payloads.GitHubAddTeamMemberPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewGitHubAddTeamMemberAction),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.RegisterFeature(
		"github_remove_team_member_action",
		"GitHub Remove Team Member Action",
		"Removes a member from a GitHub team.",
		false,
		runner.FeatureTypeAction,
		new(options.GitHubRemoveTeamMemberActionOptions),
		new(payloads.GitHubRemoveTeamMemberPayload),
		runner.FeatureResumeBehaviorNone,
		runner.APIKeyCredential,
		runner.Factory(actions.NewGitHubRemoveTeamMemberAction),
	); err != nil {
		log.Fatal(err)
	}

	if err := manifest.Validate(); err != nil {
		log.Fatal(err)
	}

	return manifest
}
