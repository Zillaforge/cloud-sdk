<!--
Sync Impact Report

Version change: N/A -> 1.0.0

Modified/Added Principles:
- Added: Library-First / Native Go SDK
- Added: Test-First (TDD) - NON-NEGOTIABLE
- Added: Direct Call Interfaces (no raw HTTP exposure)
- Added: Minimal Dependencies & Native Go
- Added: Versioning, Observability & Backwards Compatibility

Added sections:
- Security & Compliance
- Development Workflow

Removed sections:
- None

Templates checked:
- .specify/templates/plan-template.md: ✅ updated
- .specify/templates/spec-template.md: ✅ updated
- .specify/templates/tasks-template.md: ✅ updated
- .specify/templates/commands/*.md: ⚠ directory missing - manual review required

Follow-up TODOs:
- Review and adapt `.specify/templates/*` templates to reflect constitution gates (Constitution Check) where required.
-->

# Cloud SDK Constitution

## Core Principles

### Library-First / Native Go SDK
All public features MUST be provided as idiomatic, importable Go packages and types.
Packages MUST expose clear, unit-testable APIs that can be used directly from Go code
without requiring callers to construct or manipulate raw HTTP requests. Global state is
forbidden; clients and configuration MUST be explicit and passed via constructors.

Rationale: The project delivers a Go SDK. Prioritizing native packages makes the
library easy to adopt, test, and integrate in existing Go codebases.

### Test-First (TDD) (NON-NEGOTIABLE)
Every public API surface (each endpoint wrapper, model, and public helper) MUST have
tests written before implementation. Tests MUST include unit tests and a contract
test derived from the corresponding Swagger/OpenAPI definition. Tests are the
explicit specification of behavior; no PR may merge if its tests were not written
first and failing prior to implementation.

Rationale: The SDK must be reliable and maintainable across many cloud services.
TDD enforces clear contracts, prevents regressions, and drives design for testability.

### Direct Call Interfaces (no raw HTTP exposure)
The SDK MUST present direct function/method call interfaces: a constructed Client
type with methods for each API (e.g., client.Service.DoThing(ctx, params) (resp, err)).
Internals may use net/http, but callers MUST not manage HTTP details. All public
methods MUST accept context.Context and return typed responses and wrapped errors.

Rationale: Consumers asked for a Go-first SDK (callable directly). This reduces
boilerplate, improves ergonomics, and centralizes retry/auth handling.

### Minimal Dependencies & Native Go
Prefer the Go standard library and language features. External dependencies are
allowed only when they provide strong, well-justified value (e.g., code-gen or
proven test utilities). Dependency choices MUST be justified in the feature plan
and minimized to avoid security and maintenance burden.

Rationale: Smaller dependency surface reduces security risk, simplifies builds,
and aligns with the requirement to avoid unnecessary third-party packages.

### Versioning, Observability & Backwards Compatibility
SDKs MUST follow semantic versioning. Breaking changes to public APIs require a
MAJOR version bump and a documented migration path. Libraries MUST expose hooks
for logging / observability (pluggable logger interfaces or context-based hooks)
without forcing a particular vendor. Backwards-compatible enhancements SHOULD use
MINOR version bumps.

Rationale: Multi-service SDKs and long-lived clients demand clear versioning and
observability to allow safe consumption and upgrades.

## Security & Compliance
Credentials and secrets MUST NOT be logged. Configuration MUST favor environment
variables or explicit constructor parameters instead of implicit global files. TLS
certificate verification MUST be enabled by default. Any service-specific
compliance requirement (e.g., data residency) MUST be documented in the service's
spec and enforced by the service integration tests.

## Development Workflow
- Branch naming: `feature/<short>-<description>` or `fix/<issue-number>-<desc>`.
- All work follows TDD: tests written first, then implementation, then refactor.
- Every PR MUST include unit tests, and integration/contract tests when changing
	or adding an API surface. CI MUST run tests and linters; no PR merges on failing
	checks.
- Releases: Create annotated Git tag matching semantic version. Publish release
	notes that list breaking changes and migration steps.

## Governance
Amendments to this constitution require a documented proposal in a PR, approval by
the repository owners (or a two-owner majority when 3+ owners exist), and a
migration plan for any required changes to templates or automation. Minor edits
that clarify wording (non-semantic) are a PATCH bump; adding a principle or any
material expansion is a MINOR bump; removing or redefining an existing principle
in a backward-incompatible way is a MAJOR bump.

Conformance: All project plans MUST include a Constitution Check section that
validates the plan against the principles above before Phase 0 research completes.
The CI and PR review checklist MUST verify presence of tests and a passing
Constitution Check for feature work touching public APIs.

**Version**: 1.0.0 | **Ratified**: 2025-10-26 | **Last Amended**: 2025-10-26
