# mesh-github

Mesh platform connector that collects identity and access data from GitHub organizations.

## Features

### Entity Collectors
- `github_member_entity_collector` — Collects organization members as accounts
- `github_team_entity_collector` — Collects teams (groups) and team memberships
- `github_repository_entity_collector` — Collects repositories (applications) and collaborators

### Activity Collectors
- `github_audit_log_activity_collector` — Collects organization audit log events

### Actions
- `github_add_team_member_action` — Add a member to a GitHub team
- `github_remove_team_member_action` — Remove a member from a GitHub team

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
