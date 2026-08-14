# Quickstart: Fix Network Model Definition

## Prerequisites
- Go 1.21+
- Dev container or local environment aligned with repo tooling

## Implementation Steps
1. Update `models/vps/networks/network.go` to mirror `pb.NetworkInfo` and `NetCreateInput` fields.
2. Update `models/vps/networks/ports.go` to mirror the `NetPort` definition, including the embedded server summary.
3. Adjust related test fixtures in `models/vps/networks` and `modules/vps/networks/test` so they cover the new network and port fields.
4. Refresh contract tests to validate Swagger parity (ensure generated payloads include new JSON keys).
5. Re-run `go fmt` on modified files.
6. Execute `make check` to run linting and all Go test suites.
7. Prepare to validate the six network module operations against the live environment using `curl` once base URL and token are supplied.

## Verification
- Run unit tests: `go test ./models/vps/networks`
- Run module-level integration + contract tests: `go test ./modules/vps/networks/...`
- Optional: Execute repository-wide tests if affected by shared types: `go test ./...`
- Run `make check` and ensure it completes without errors.
- After automated checks, execute the six network module operations via `curl` against the real system (base URL/token provided at test time) and compare responses with SDK results.

Expected Result: All suites pass with newly asserted fields present in responses and requests.
