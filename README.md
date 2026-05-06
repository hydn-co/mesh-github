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
