# Implementation Plan: VRM Tag and Repository APIs Client SDK

**Branch**: `008-vrm-tag-repository` | **Date**: 2025-11-12 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/008-vrm-tag-repository/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement a Go SDK client for Virtual Registry Management (VRM) API following the existing VPS service pattern. The implementation will provide project-scoped access to Repository and Tag operations (CRUD) with Bearer token authentication, pagination, filtering, and namespace support. Only User/Repository and User/Tag endpoints from vrm.yaml will be implemented, excluding all Admin/*, MemberAcl, ProjectAcl, Image, Export, and Snapshot endpoints.

**Technical Approach**:
- Follow identical architecture pattern as VPS service: `modules/vrm/client.go` with sub-packages `modules/vrm/repositories` and `modules/vrm/tags`
- Data models in `models/vrm/repositories` and `models/vrm/tags` matching API swagger spec
- Project-scoped access via `cloudsdk.Client.Project(projectID).VRM()`
- Leverage existing `internal/http.Client` for HTTP operations with built-in retry and 30s default timeout
- TDD approach with unit, contract, and integration tests targeting 80%+ coverage
- API Request/Response struct naming convention (not Input/Output)
- List APIs return `[]*Resource` for direct index access

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: 
- Go standard library (`net/http`, `encoding/json`, `context`, `time`)
- Internal packages: `internal/http` (HTTP client with retry), `internal/backoff`, `internal/types`
- Zero external dependencies (following constitution principle)

**Storage**: N/A (SDK client library, no persistence layer)  
**Testing**: Go testing package (`testing`, `net/http/httptest` for mocking)
- Unit tests: model validation, struct marshaling/unmarshaling
- Contract tests: validate against vrm.yaml OpenAPI specification
- Integration tests: end-to-end API call flows with test server

**Target Platform**: Cross-platform (Linux, macOS, Windows) - Go SDK library  
**Project Type**: Single project (Go SDK library following existing VPS pattern)  
**Performance Goals**: 
- HTTP client: 30-second default timeout (configurable)
- List operations: support pagination with limit/offset
- Concurrent requests: thread-safe by design

**Constraints**: 
- Must follow exact VPS service architectural pattern
- API paths match vrm.yaml specification exactly
- Field names in models match API swagger spec (camelCase in JSON tags)
- Request/Response naming (not Input/Output)
- List APIs return `[]*Resource` directly

**Scale/Scope**: 
- 2 resource types: Repository (5 operations), Tag (6 operations)
- 11 total API endpoints to implement
- Target: 80%+ test coverage
- Exclude: Admin/*, MemberAcl, ProjectAcl, Image, Export, Snapshot endpoints

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### Initial Check (Pre-Phase 0)

The feature plan MUST satisfy all applicable Cloud SDK Constitution principles:

✅ **TDD mandatory**: 
- Unit tests for all models (validation, marshaling/unmarshaling)
- Contract tests derived from vrm.yaml OpenAPI specification
- Integration tests for API operations
- Tests written first following TDD workflow
- Target: 80%+ coverage per user requirement

✅ **Public API shape**: 
- Idiomatic Go: `client.Project(projectID).VRM().Repositories().List(ctx, opts)`
- No raw HTTP exposure to callers
- All methods accept `context.Context`
- Return typed responses (`*Repository`, `[]*Tag`, etc.) and wrapped errors
- Follows existing VPS service pattern exactly

✅ **Dependencies**: 
- Zero external dependencies
- Uses Go standard library only (`net/http`, `encoding/json`, `context`, `time`)
- Leverages existing internal packages (`internal/http`, `internal/backoff`, `internal/types`)
- Justification: All functionality achievable with stdlib + internal infrastructure

✅ **Versioning**: 
- Change type: MINOR version bump (new VRM service addition, no breaking changes to existing APIs)
- Extends `cloudsdk.ProjectClient` with new `VRM()` method
- No breaking changes to existing VPS or core SDK APIs
- Migration: None required (additive feature)

✅ **Observability**: 
- Uses existing `types.Logger` interface (pluggable logging)
- HTTP client supports logger injection via `NewClient()`
- No vendor lock-in (follows existing VPS pattern)

✅ **Security**: 
- Bearer token passed via constructor, not logged
- Token included in HTTP Authorization header via `internal/http.Client`
- TLS verification enabled by default (inherited from `http.Client`)
- No secrets in logs (following internal/http implementation pattern)
- Configuration via explicit constructor parameters (baseURL, token)

### Post-Design Re-Evaluation (After Phase 1)

After completing research, data model, API contracts, and quickstart documentation:

✅ **TDD mandatory - VALIDATED**: 
- Data model includes validation methods (`Validate()` on Repository, Tag, Request types)
- Contract tests clearly specified in `contracts/api-contracts.md` with 11 endpoint contracts
- Unit test coverage plan: Model validation, JSON marshaling/unmarshaling, field constraints
- Contract test coverage plan: All 11 endpoints with request/response verification using httptest
- Integration test coverage plan: CRUD workflows, pagination, filtering, error handling
- Test-first workflow enforced: Write test → Implement → Validate → Refactor
- 80%+ coverage requirement tracked via `make check` per phase

✅ **Public API shape - VALIDATED**: 
- Architecture confirmed: `client.Project(projectID).VRM().Repositories()/Tags()`
- Method signatures verified in `contracts/api-contracts.md`:
  - All methods accept `context.Context` as first parameter
  - Repository methods: `List(ctx, opts) ([]*Repository, error)`, `Create(ctx, req) (*Repository, error)`, etc.
  - Tag methods: `List(ctx, opts) ([]*Tag, error)`, `ListByRepository(ctx, repoID, opts) ([]*Tag, error)`, etc.
- Naming convention validated: Request/Response (not Input/Output)
- Return types validated: `[]*Resource` for List (no wrapper struct), `*Resource` for single entities
- Idiomatic error handling: `(result, error)` pattern throughout

✅ **Dependencies - VALIDATED**: 
- Data model uses only Go standard library: `encoding/json`, `time`, `fmt`
- HTTP operations use only `internal/http.Client` (which wraps stdlib `net/http`)
- No external dependencies introduced in any Phase 1 design
- Confirmed in `data-model.md`: All types use stdlib + internal packages only

✅ **Versioning - VALIDATED**: 
- MINOR version bump confirmed (additive feature, no breaking changes)
- New methods added to ProjectClient: `VRM() *vrm.Client`
- Existing VPS methods unchanged
- API surface expansion documented in quickstart.md
- Backward compatibility: Existing SDK users unaffected

✅ **Observability - VALIDATED**: 
- Logger integration via existing `types.Logger` interface
- HTTP client inherits logging from `internal/http.Client`
- No new observability mechanisms required
- Follows VPS pattern exactly (verified in research.md)

✅ **Security - VALIDATED**: 
- Bearer token flow confirmed in `contracts/api-contracts.md`: "Authorization: Bearer {token}" header on all requests
- Token passed via SDK constructor, stored in ProjectClient, inherited by VRM client
- No token logging (inherits from internal/http.Client behavior)
- TLS verification default behavior inherited from stdlib `http.Client`
- Namespace security: Private/public namespaces supported via query parameters

**Constitution Compliance: PASS ✅**  
All principles satisfied in initial check AND post-design validation.  
No violations or justifications required.

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
# VRM Service Implementation (follows VPS pattern)

client.go                          # Extend ProjectClient with VRM() method

models/vrm/                        # Data models matching vrm.yaml spec
├── common/
│   └── common.go                  # Shared types (IDName, etc.)
├── repositories/
│   ├── repository.go              # Repository model, ListOptions, Request/Response structs
│   └── repository_test.go         # Model validation and marshaling tests
└── tags/
    ├── tag.go                     # Tag model, ListOptions, Request/Response structs
    └── tag_test.go                # Model validation and marshaling tests

modules/vrm/                       # API client implementation
├── client.go                      # Main VRM client (project-scoped)
├── client_test.go                 # VRM client tests
├── repositories/
│   ├── client.go                  # Repository operations (List, Create, Get, Update, Delete)
│   ├── client_test.go             # Unit tests for repository client
│   └── test/
│       ├── contract_test.go       # Contract tests vs vrm.yaml
│       └── integration_test.go    # Integration tests with test server
└── tags/
    ├── client.go                  # Tag operations (List, ListByRepository, Create, Get, Update, Delete)
    ├── client_test.go             # Unit tests for tag client
    └── test/
        ├── contract_test.go       # Contract tests vs vrm.yaml
        └── integration_test.go    # Integration tests with test server

swagger/
└── vrm.yaml                       # OpenAPI spec (already exists, used for contract tests)
```

**Structure Decision**: 

This implementation follows the **existing VPS service pattern exactly**:

1. **Models** (`models/vrm/`): Data structures matching API specification
   - `common/`: Shared types across VRM resources
   - `repositories/`: Repository models, list options, request/response types
   - `tags/`: Tag models, list options, request/response types

2. **Modules** (`modules/vrm/`): API client implementation
   - `client.go`: Main VRM client providing access to sub-clients
   - `repositories/client.go`: Repository CRUD operations
   - `tags/client.go`: Tag CRUD operations
   - Each sub-package includes unit tests and integration tests

3. **Root-level**: Extension of `client.go` to add `ProjectClient.VRM()` method

4. **Testing**:
   - Unit tests: `*_test.go` files alongside implementation
   - Contract tests: `test/contract_test.go` validating against vrm.yaml
   - Integration tests: `test/integration_test.go` with httptest server

This structure maintains consistency with VPS (`modules/vps/flavors`, `models/vps/flavors`) and allows users to discover the VRM API through the same pattern: `client.Project(id).VRM().Repositories()` / `.Tags()`

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

**No violations detected.** All constitution principles are satisfied:
- Zero external dependencies (stdlib only)
- TDD workflow with 80%+ coverage target
- Idiomatic Go API design following VPS pattern
- MINOR version bump (additive feature)
- Pluggable logging via existing types.Logger
- Security via Bearer token in HTTP Authorization header
