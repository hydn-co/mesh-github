package options

func (o *GitHubMemberEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateGitHubOptionsCore(nil)
	}

	return validateGitHubOptionsCore(&o.GitHubOptionsCore)
}

func (o *GitHubRepositoryEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateGitHubOptionsCore(nil)
	}

	return validateGitHubOptionsCore(&o.GitHubOptionsCore)
}

func (o *GitHubTeamEntityCollectorOptions) Validate() error {
	if o == nil {
		return validateGitHubOptionsCore(nil)
	}

	return validateGitHubOptionsCore(&o.GitHubOptionsCore)
}

func (o *GitHubAuditLogActivityCollectorOptions) Validate() error {
	if o == nil {
		return validateGitHubOptionsCore(nil)
	}

	return validateGitHubOptionsCore(&o.GitHubOptionsCore)
}

func (o *GitHubAddTeamMemberActionOptions) Validate() error {
	if o == nil {
		return validateGitHubOptionsCore(nil)
	}

	return validateGitHubOptionsCore(&o.GitHubOptionsCore)
}

func (o *GitHubRemoveTeamMemberActionOptions) Validate() error {
	if o == nil {
		return validateGitHubOptionsCore(nil)
	}

	return validateGitHubOptionsCore(&o.GitHubOptionsCore)
}
