# Implementation Plan: Fix Flavors Model

**Branch**: `001-fix-flavors-model` | **Date**: November 8, 2025 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-fix-flavors-model/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Update the Flavor model in the Go SDK to match the pb.FlavorInfo definition from vps.yaml, including all required fields, correct JSON tags, and support for server-side filtering in List API. Implement breaking changes with migration notes for field name updates (VCPUs→VCPU, RAM→Memory, Items→Flavors in FlavorListResponse).

**Implementation Status**: 
- ✅ User Story 1 (T007-T012) - Correct Flavor Model Structure - Complete
- ✅ User Story 2 (T013-T018) - GPU Support in Flavors - Complete
- ✅ User Story 3 (T019-T024) - Flavor Timestamps - Complete
- ✅ Phase 6 (T025-T030) - Client Implementation & Filtering - Complete

See [CHANGES.md](CHANGES.md) for detailed change summary.

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: Standard library (net/http, encoding/json, context, time)  
**Storage**: N/A (API client SDK)  
**Testing**: Go testing framework with httptest for contract tests  
**Target Platform**: Cross-platform (Linux dev, but SDK works on any Go-supported platform)  
**Project Type**: SDK library  
**Performance Goals**: Standard API client performance (<200ms typical requests)  
**Constraints**: Idiomatic Go API, no raw HTTP exposure, context.Context usage  
**Scale/Scope**: Single model update with filtering enhancements

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The feature plan MUST satisfy all applicable Cloud SDK Constitution principles:

- TDD mandatory: tests are written first and initially fail; include unit tests and contract tests derived from the corresponding Swagger/OpenAPI. ✓ Satisfied - Unit tests for model and client, contract tests via httptest.
- Public API shape: expose idiomatic Go packages and a Client with methods per API; do not expose raw HTTP to callers; all public methods accept `context.Context` and return typed responses and wrapped errors. ✓ Satisfied - Client.List() and Client.Get() methods.
- Dependencies: prefer standard library; any external dependency MUST be justified here with clear value (security, maintainability, or functionality). ✓ Satisfied - Only standard library used.
- Versioning: document whether the change is MAJOR/MINOR/PATCH and include a migration note for any breaking change. ✓ Satisfied - MAJOR version bump required due to field name changes (VCPUs→VCPU, RAM→Memory, FlavorListResponse.Items→Flavors).
- Observability: provide hooks for logging/metrics without forcing a vendor. ✓ N/A - No new observability features.
- Security: do not log secrets; TLS verification enabled by default; configuration through env vars or explicit constructors. ✓ N/A - No security changes.

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

## Project Structure

### Documentation (this feature)

```text
specs/001-fix-flavors-model/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
models/vps/flavors/
├── flavor.go           # Flavor struct and related types
└── ...

modules/vps/flavors/
├── client.go           # Client implementation with List/Get methods
├── client_test.go      # Contract tests
└── test/
    └── flavors_list_test.go  # Unit tests for List functionality
```

**Structure Decision**: SDK library structure with separate models and modules directories. Models contain data structures, modules contain client implementations and tests.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

| Violation | Why Needed | Simpler Alternative Rejected Because |
|-----------|------------|-------------------------------------|
| [e.g., 4th project] | [current need] | [why 3 projects insufficient] |
| [e.g., Repository pattern] | [specific problem] | [why direct DB access insufficient] |
