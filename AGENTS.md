# Agent instructions — mesh-github

## Before any task
- Read this file first.
- Read the collector you are modifying, or the closest existing collector, before changing code.
- If you touch capability registration or option schemas, read `cmd/main.go` and `internal/options/options.go` first.

## Project summary
- **mesh-github**: A Mesh platform connector that collects identity and access data from GitHub.
- **Entry point**: `cmd/main.go`
- **Framework**: Built on [mesh-sdk](https://github.com/hydn-co/mesh-sdk)
- **Language**: Go (see `go.mod` for version)

## Non-negotiable rules
1. All collectors embed `*connector.TypedFeatureContext` and implement `Init`/`Start`/`Stop`.
2. Factory functions take `*connector.TypedFeatureContext` and return `runner.Feature`.
3. All option types live in `internal/options/options.go` with polymorphic registration.
4. Use only `net/http` for provider API calls — no provider SDKs or third-party HTTP clients.
5. Wrap errors with `fmt.Errorf("context: %w", err)`.
6. Check `ctx.Err()` inside loops and after enumeration.
7. Use `testkit.TestPolymorphicRegistrations` for option registration coverage.
8. Use behavioral test names: `TestShould{Expectation}When{Condition}`.

## Primary sources
| Need | Source |
|------|--------|
| Collector examples | `internal/collectors/*.go` |
| Options | `internal/options/options.go` |
| API client | `internal/api/` |
| Manifest | `cmd/main.go` |
| SDK framework | [mesh-sdk](https://github.com/hydn-co/mesh-sdk) |
