# Implementation Plan: [FEATURE]

**Branch**: `[###-feature-name]` | **Date**: [DATE] | **Spec**: [link]
**Input**: Feature specification from `/specs/[###-feature-name]/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

 Implement a vps Snapshot resource following `swagger/vps.yaml`. The plan covers the data model (Snapshot entity), SDK model generation and manual adjustments, a `modules/vps/snapshots` client in Go following existing `modules/vps/volumes` patterns, and unit and contract tests that provide at least 85% coverage.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

**Language/Version**: Go 1.21+ (repository standard)  
**Primary Dependencies**: Go standard library; use repo internal packages (internal/http, models/vps/common); minimal external dependencies only if justified.  
**Storage**: N/A (cloud service-backed snapshot API; SDK does not change storage implementation)  
**Testing**: Go's builtin `testing` package; unit tests and contract tests under `hack/` (use rest test harness)  
**Target Platform**: Linux (CI), cross-compiled for Go as needed  
**Project Type**: API client SDK (Go package modules/vps/snapshots)  
**Performance Goals**: List queries must return 95% within 1s under 1k snapshots (test environment)  
**Constraints**: All public APIs use `context.Context` and typed results; prefer existing `models/vps` and `modules/vps` patterns  
**Scale/Scope**: Client usage within VPS SDK; do not implement lifecycle/TTL or retention in v1

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The feature plan MUST satisfy all applicable Cloud SDK Constitution principles:

- TDD mandatory: tests are written first and initially fail; include unit tests and
 - TDD mandatory: tests are written first and initially fail; include unit tests and contract tests derived from `swagger/vps.yaml`.
- Public API shape: expose idiomatic Go packages and a Client with methods per API;
  do not expose raw HTTP to callers; all public methods accept `context.Context`
  and return typed responses and wrapped errors.
- Dependencies: prefer standard library; any external dependency MUST be justified
  here with clear value (security, maintainability, or functionality).
- Versioning: document whether the change is MAJOR/MINOR/PATCH and include a
  migration note for any breaking change.
- Observability: provide hooks for logging/metrics without forcing a vendor.
- Security: do not log secrets; TLS verification enabled by default; configuration
  through env vars or explicit constructors.
 - Observability: provide hooks for logging/metrics without forcing a vendor.
 - Security: do not log secrets; TLS verification enabled by default; configuration through env vars or explicit constructors.

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)
<!--
  ACTION REQUIRED: Replace the placeholder tree below with the concrete layout
  for this feature. Delete unused options and expand the chosen structure with
  real paths (e.g., apps/admin, packages/something). The delivered plan must
  not include Option labels.
-->

```text
# [REMOVE IF UNUSED] Option 1: Single project (DEFAULT)
src/
├── models/
├── services/
├── cli/
└── lib/

tests/
├── contract/
├── integration/
└── unit/

# [REMOVE IF UNUSED] Option 2: Web application (when "frontend" + "backend" detected)
backend/
├── src/
│   ├── models/
│   ├── services/
│   └── api/
└── tests/

frontend/
├── src/
│   ├── components/
│   ├── pages/
│   └── services/
└── tests/

# [REMOVE IF UNUSED] Option 3: Mobile + API (when "iOS/Android" detected)
api/
└── [same as backend above]

ios/ or android/
└── [platform-specific structure: feature modules, UI flows, platform tests]
```

**Structure Decision**: [Document the selected structure and reference the real
directories captured above]

 **Structure Decision**: This repository is an API client SDK for Go. Add new SDK models to `models/vps/snapshots` and a new module to `modules/vps/snapshots`. Tests live in corresponding module; contract tests should be added to the repository's contract test suite.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |

### Phase 0: Research tasks (generated)

- R1: Confirm which shared models from `models/vps/common` to reuse for `Snapshot` relationships (Project, User). (owner: dev)
- R2: Confirm indexable list contract details with API team or `vps.yaml` (owner: dev)
- R3: Decide canonical time serialization (`createdAt`, `updatedAt`) mapping to Go `time.Time` and RFC3339 serialization (owner: dev)

### Phase 1: Design & Contracts tasks

- D1: Implement `models/vps/snapshots` with types and Validate() methods.
- D2: Implement `modules/vps/snapshots` client with methods: Create(ctx, req), List(ctx, opts), Get(ctx, id), Update(ctx, id, req), Delete(ctx, id).
 - D2: Implement `modules/vps/snapshots` client with methods: Create(ctx, req), List(ctx, opts), Get(ctx, id), Update(ctx, id, req), Delete(ctx, id). Follow `modules/vps/volumes` and `modules/vps/networks` client patterns: use `internal/http` request object, centralize authentication and retries, wrap errors with contextual messages, and ensure all public methods accept `context.Context`.
- D3: Add unit tests for client with request validation and response mapping; follow `modules/vps/volumes` patterns.
- D4: Add repository-level contract tests verifying create/list/get/update/delete including hot-snapshot and deletion behavior.

### Phase 2: Implementation tasks (brief)

- I1: Write `models/vps/snapshots` fetched from `swagger/vps.yaml` and add `SnapshotStatus` enumerations. The Go `Snapshot` struct must preserve Swagger field names including `project`, `user`, and `namespace` and serialize times per RFC3339.
- I2: Write `modules/vps/snapshots/client.go` similar to `modules/vps/volumes/client.go` and wire it to `internal/http`.
- I3: Add unit tests and contract tests, run `go test ./...` and `hack/test` for contract tests.
- I4: Update README and docs to include usage and code examples.
