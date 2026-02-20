# Repository Guidelines

## Project Structure & Module Organization
- `client.go`, `errors.go`: top-level SDK entry points and shared errors.
- `cmd/`: CLI entry (`cmd/main.go`).
- `modules/`: service clients (e.g., `modules/vps`, `modules/iam`, `modules/vrm`).
- `models/`: request/response types grouped by service (e.g., `models/vps/servers`).
- `internal/`: private utilities (`internal/http`, `internal/backoff`, `internal/types`, `internal/waiter`).
- `docs/`: service docs and examples.
- `swagger/`: OpenAPI specs.
- `specs/`: design docs and checklists (excluded from build/test/lint).
- `hack/`: REST examples and ad hoc scripts.
- Tests live alongside code as `*_test.go` (including module/model tests).

## Build, Test, and Development Commands
- `make help`: list all targets.
- `make fmt`: run `gofmt` and `goimports` (if installed).
- `make lint`: run `golangci-lint` or fallback to `go vet`.
- `make test`: run unit tests with coverage.
- `make test-race`: run tests with the race detector (set `CGO_ENABLED=1`).
- `make coverage`: generate `coverage.html`.
- `make build`: build all packages.
- `make check`: `fmt + lint + test`.

## Coding Style & Naming Conventions
- Go standard formatting: tabs for indentation via `gofmt`.
- Imports are grouped and normalized by `goimports` with `-local github.com/Zillaforge/cloud-sdk`.
- Package names are short and lowercase (`vps`, `iam`, `vrm`).
- Tests use Go’s `*_test.go` naming and live next to implementation.

## Testing Guidelines
- Framework: Go `testing` package.
- Run all tests: `make test` or `go test ./...`.
- Prefer service-level tests under `modules/<service>/...` and model tests under `models/<service>/...`.
- Add regression tests for API changes (update `docs/` and `swagger/` if behavior changes).

## Commit & Pull Request Guidelines
- Commit messages follow Conventional Commits in practice (e.g., `feat: ...`, `refactor: ...`, `chore: ...`, `build: ...`).
- PRs should include:
  - A concise summary and rationale.
  - Tests run (e.g., `make test`, `make lint`).
  - Links to relevant `specs/` docs or API contracts when applicable.

## Security & Configuration Tips
- Do not commit credentials or tokens; use environment variables in examples.
- When updating API contracts, keep `swagger/` and `docs/` in sync.
