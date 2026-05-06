# Agent instructions — mesh-github

## Before any task
- Read this file first.
- Read the collector you are modifying, or the closest existing collector, before changing code.
- If you touch capability registration or option schemas, read `cmd/main.go`, `internal/options/`, and `internal/payloads/` first.

## Project summary
- **mesh-github**: A Mesh platform connector that collects identity and access data from GitHub.
- **Entry point**: `cmd/main.go`
- **Framework**: Built on [mesh-sdk](https://github.com/hydn-co/mesh-sdk)
- **Language**: Go (see `go.mod` for version)

## Non-negotiable rules
1. All collectors embed `*connector.TypedFeatureContext` and implement `Init`/`Start`/`Stop`.
2. Factory functions take `*connector.TypedFeatureContext` and return `runner.Feature`.
3. Keep option and payload types split by feature under `internal/options/` and `internal/payloads/`, with central registration files in each package.
4. Prefer `mesh-sdk/pkg/connectorutil` for feature logging, payload validation, lifecycle state, and shared credential parsing before adding local helpers.
5. Prefer `connectorutil.FeatureState` for collector and action readiness guards instead of ad hoc initialized flags.
6. Prefer `connectorutil.ExtractAPIKey(...)` for the standard `api_key` credential shape; only keep provider-specific fallbacks when required for backward compatibility.
7. Use only `net/http` for provider API calls — no provider SDKs or third-party HTTP clients.
8. Wrap errors with `fmt.Errorf("context: %w", err)`.
9. Check `ctx.Err()` inside loops and after enumeration.
10. Use `testkit.TestPolymorphicRegistrations` for option and payload registration coverage.
11. Use behavioral test names: `TestShould{Expectation}When{Condition}`.

## Primary sources
| Need | Source |
|------|--------|
| Collector examples | `internal/collectors/*.go` |
| Options | `internal/options/*.go` |
| Payloads | `internal/payloads/*.go` |
| API client | `internal/api/` |
| Manifest | `cmd/main.go` |
| SDK framework | [mesh-sdk](https://github.com/hydn-co/mesh-sdk) |
