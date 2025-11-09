# Implementation Plan: Fix Server Model & Module

**Branch**: `005-fix-server-model` | **Date**: 2025-11-09 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/workspaces/cloud-sdk/specs/005-fix-server-model/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Fix server model and module to align with documentation and Swagger spec. Use custom types for status/enum fields (e.g., `type ServerStatus string` with const definitions) instead of plain strings. Implement server CRUD, actions, metrics, NIC management with Resource.SubResource().Verb() pattern. Support query parameter filtering for List APIs. Achieve 80%+ test coverage passing make check.

## Technical Context

<!--
  ACTION REQUIRED: Replace the content in this section with the technical details
  for the project. The structure here is presented in advisory capacity to guide
  the iteration process.
-->

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: Go standard library (`net/http`, `encoding/json`, `context`, `time`)  
**Storage**: N/A (API client SDK)  
**Testing**: Go testing framework with contract tests  
**Target Platform**: Linux  
**Project Type**: SDK library  
**Performance Goals**: <500ms p95 response time, 100 req/s throughput  
**Constraints**: <100MB memory usage, 100 concurrent requests  
**Scale/Scope**: 10k servers, 10k NICs

## Constitution Check

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The feature plan satisfies all applicable Cloud SDK Constitution principles:

- TDD mandatory: Unit tests and contract tests will be written first for the updated server model and NIC sub-resource APIs.
- Public API shape: The SDK exposes a Client with methods like client.Servers().List(), client.Servers().NICs().List(); no raw HTTP exposure.
- Dependencies: Uses only Go standard library; no external dependencies added.
- Versioning: This is a PATCH version change (fixing model alignment without breaking public APIs).
- Observability: The SDK supports pluggable logging via context or interfaces.
- Security: Credentials not logged; TLS verification enabled by default; configuration via constructors or env vars.

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

```text
models/vps/servers/
├── server.go          # Server model with custom type for status
├── actions.go         # ServerAction custom type and constants
├── console.go         # Console types
├── metrics.go         # Metrics types
├── nics.go           # ServerNIC model
└── volumes.go        # Volume attachment types

modules/vps/servers/
├── client.go         # Server CRUD operations
├── client_test.go    # Unit tests
├── actions.go        # Server actions (start, stop, reboot, etc.)
├── metrics.go        # Metrics retrieval
├── nics.go          # NIC sub-resource operations
├── nics_test.go     # NIC tests
└── test/            # Contract tests
```

**Structure Decision**: Single SDK project with models and modules separation. Server and NIC models use custom types for status/enum fields (e.g., `type ServerStatus string` with const values). All structs use Request/Response suffix for consistency (e.g., ServerCreateRequest, ServerActionResponse). Sub-resource pattern: `client.Servers().NICs().List()`.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
