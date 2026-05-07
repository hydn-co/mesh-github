# mesh-github

Mesh platform connector that collects identity and access data from GitHub organizations.

## Features

### Entity Collectors
- `collect_members` — Collects organization members as accounts
- `collect_teams` — Collects teams (groups) and team memberships

### Activity Collectors
- `collect_audit_log` — Collects organization audit log events

### Actions
- `add_team_member` — Add a member to a GitHub team
- `remove_team_member` — Remove a member from a GitHub team

## Authentication

`mesh-github` uses a GitHub personal access token.

Prefer a fine-grained personal access token when your organization allows it. Use a classic personal access token only when you specifically need classic scopes or your organization has not enabled fine-grained token approval for the target org.

### Acquire a Fine-Grained Personal Access Token

1. In GitHub, open `Settings -> Developer settings -> Personal access tokens -> Fine-grained tokens`.
2. Select `Generate new token`.
3. Set the `Resource owner` to the target organization used by this connector.
4. Grant only the organization permissions needed for the features you will enable.
5. Generate the token and store it in your connector credentials.

### Acquire a Classic Personal Access Token

1. In GitHub, open `Settings -> Developer settings -> Personal access tokens -> Tokens (classic)`.
2. Select `Generate new token (classic)`.
3. Grant the scopes required for the enabled features.
4. Generate the token and store it in your connector credentials.

### Credential Format

Preferred credentials payload:

```json
{
	"api_key": "github_pat_..."
}
```

Legacy credentials payload is still accepted:

```json
{
	"token": "github_pat_..."
}
```

### Required Permissions By Feature

| Feature | Fine-grained PAT permissions | Classic PAT scopes | Additional GitHub requirements |
| --- | --- | --- | --- |
| `collect_members` | `Members` organization permission: `read` | `read:org` | Token holder must be a member of the target organization to retrieve concealed members. |
| `collect_teams` | `Members` organization permission: `read` | `read:org` | Token holder must be able to see the teams being enumerated. |
| `collect_audit_log` | `Administration` organization permission: `read` | `read:audit_log` | Token holder must be an organization owner. GitHub may return `404` when the audit log is unavailable for the org plan or the token lacks required access. |
| `add_team_member` | `Members` organization permission: `write` | `admin:org` | Token holder must be an organization owner or a team maintainer. Adding a non-member to a team requires an organization owner. Team synchronization can block API-based membership changes. |
| `remove_team_member` | `Members` organization permission: `write` | `admin:org` | Token holder must have team admin privileges or be an organization owner. Team synchronization can block API-based membership changes. |

### Least-Privilege Recommendations

- If you only run entity collectors, start with `Members: read` on a fine-grained token.
- If you also enable `collect_audit_log`, add `Administration: read`.
- If you enable `add_team_member` or `remove_team_member`, upgrade `Members` from `read` to `write`.
- For classic PATs, use `read:org` for read-only collection, add `read:audit_log` for audit collection, and use `admin:org` for team membership mutations.

## Development

### Prerequisites
- Go 1.25.6+

### Build
```bash
go build ./...
```

### Test
```bash
go test ./...
```

### Lint
```bash
golangci-lint run
```
