# Implementation Plan: Volumes API Client

**Branch**: `008-volumes-api` | **Date**: 2025-11-13 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/008-volumes-api/spec.md`

**Note**: This template is filled in by the `/speckit.plan` command. See `.specify/templates/commands/plan.md` for the execution workflow.

## Summary

Implement Go SDK client for VPS Volumes and VolumeTypes APIs based on swagger/vps.yaml specification. The implementation provides two resource clients: `Volumes` for managing block storage volumes (CRUD operations and actions like attach/detach/extend) and `VolumeTypes` for listing available storage types. Both clients follow the established VPS module patterns with models in `models/vps/volumes` and `models/vps/volumetypes`, client implementations in `modules/vps/volumes` and `modules/vps/volumetypes`, using the existing internal HTTP client for consistent retry/error handling and Bearer Token authentication.

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: 
- Go standard library (`context`, `net/http`, `encoding/json`, `fmt`, `time`)
- Internal packages: `internal/http` (HTTP client with retry), `internal/types` (common types), `internal/backoff` (retry strategy)
- Swagger specification: `swagger/vps.yaml` (volume and volume_types tags)

**Storage**: N/A (API client library, no persistence layer)  
**Testing**: Go testing (`testing` package), table-driven tests, contract tests against Swagger spec  
**Target Platform**: Cross-platform Go library (Linux, macOS, Windows)  
**Project Type**: Single project (SDK library)  
**Performance Goals**: 
- API response handling within specified timeouts (2-10 seconds per spec)
- Support projects with large volume counts (full list retrieval without pagination)

**Constraints**: 
- No client-side size validation (server-side enforcement)
- No pagination support (API returns complete list)
- Bearer Token authentication (caller-managed lifecycle)
- Retry with exponential backoff (consistent with existing `internal/http/client.go`)

**Scale/Scope**: 
- 2 resources: Volumes (6 operations), VolumeTypes (1 operation)
- 7 API endpoints total
- Unit test coverage target: 85%+
- Contract tests for all endpoints

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The feature plan MUST satisfy all applicable Cloud SDK Constitution principles:

✅ **TDD mandatory**: Tests written first per TDD workflow:
- Unit tests for all models (validation, serialization) - target 85%+ coverage
- Unit tests for all client methods (6 volumes operations + 1 volume types operation) - target 85%+ coverage
- Contract tests derived from swagger/vps.yaml (volumes and volume_types tags)
- Tests written before implementation, initially failing

✅ **Public API shape**: Idiomatic Go packages and Client pattern:
- Exposes `modules/vps/volumes.Client` and `modules/vps/volumetypes.Client`
- Methods: `volumes.List(ctx, opts)`, `volumes.Create(ctx, input)`, `volumes.Get(ctx, id)`, `volumes.Update(ctx, id, input)`, `volumes.Delete(ctx, id)`, `volumes.Action(ctx, id, input)`
- Methods: `volumetypes.List(ctx)` returning `[]string`
- All methods accept `context.Context` as first parameter
- Returns typed responses (`*Volume`, `[]*Volume`, `[]string`) and wrapped errors
- No raw HTTP exposed to callers
- Follows existing VPS module patterns (flavors, keypairs, networks, etc.)

✅ **Dependencies**: Minimal, justified dependencies:
- Go standard library only (no external dependencies)
- Reuses existing internal packages (`internal/http`, `internal/types`, `internal/backoff`)
- Justification: Consistency with existing VPS modules, proven retry/error handling infrastructure

✅ **Versioning**: MINOR version bump (new functionality, no breaking changes):
- Adds new `Volumes` and `VolumeTypes` modules to existing VPS client
- No changes to existing APIs
- No migration required for existing users
- Semantic version: v1.X.0 → v1.(X+1).0

✅ **Observability**: Logger interface support:
- Uses existing `types.Logger` interface passed through client initialization
- Logs HTTP requests/responses at DEBUG level (method, URL, status code)
- Logs retry attempts with backoff duration
- No vendor lock-in (pluggable logger interface)
- Consistent with existing VPS modules (IAM, VRM)

✅ **Security**: Compliant with security requirements:
- Bearer Token never logged
- Token passed via constructor, stored in client struct
- TLS verification enabled by default (handled by net/http)
- No secrets in error messages
- Configuration through explicit constructor parameters (`NewClient(baseURL, token, projectID, ...)`)

**Constitution Compliance**: ✅ PASS - All principles satisfied

## Project Structure

### Documentation (this feature)

```text
specs/008-volumes-api/
├── plan.md              # This file (/speckit.plan command output)
├── spec.md              # Feature specification (from /speckit.specify)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   ├── volumes.yaml     # Extracted from swagger/vps.yaml (volumes tag)
│   └── volume_types.yaml # Extracted from swagger/vps.yaml (volume_types tag)
├── checklists/          # Quality checklists
│   └── requirements.md  # Specification quality checklist (completed)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# Models (data structures matching Swagger definitions)
models/vps/
├── common/
│   └── common.go              # Shared types (IDName) - already exists
├── volumes/
│   ├── volume.go              # Volume model, validation, list options, action types
│   └── volume_test.go         # Unit tests for volume models (85%+ coverage)
└── volumetypes/
    ├── volumetype.go          # VolumeType list response model
    └── volumetype_test.go     # Unit tests for volume type models (85%+ coverage)

# Client implementations (API operations)
modules/vps/
├── client.go                  # VPS client (add Volumes() and VolumeTypes() accessors)
├── client_test.go             # Update with new module tests
├── volumes/
│   ├── client.go              # Volumes client with List/Create/Get/Update/Delete/Action
│   └── client_test.go         # Unit tests for volumes client (85%+ coverage)
└── volumetypes/
    ├── client.go              # VolumeTypes client with List
    └── client_test.go         # Unit tests for volume types client (85%+ coverage)

# Integration/Contract tests
tests/
├── contract/
│   ├── volumes_test.go        # Contract tests against swagger/vps.yaml
│   └── volumetypes_test.go    # Contract tests for volume types
└── integration/
    ├── volumes_test.go        # End-to-end volume lifecycle tests
    └── volumetypes_test.go    # End-to-end volume types tests

# Swagger specification (reference only, not modified)
swagger/
└── vps.yaml                   # VPS API specification (volumes, volume_types tags)
```

**Structure Decision**: Following existing single-project structure with VPS module pattern. Models in `models/vps/{resource}/`, clients in `modules/vps/{resource}/`, maintaining consistency with flavors, keypairs, networks, etc. Uses shared `common.IDName` type for resource references in Volume model.

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

No constitution violations. This feature follows all established patterns and principles.

---

## Phase 0: Research & Technical Decisions

**Objective**: Resolve technical unknowns and establish implementation patterns by researching existing codebase patterns and Swagger specification details.

### Research Tasks

1. **Volume API Response Format Analysis**
   - **Question**: How does the Volumes List API JSON response map to `[]*Volume` array?
   - **Investigation**: Examine `swagger/vps.yaml` → `pb.VolumeListOutput` structure
   - **Expected Outcome**: Response wrapper structure with `volumes` array field

2. **VolumeTypes API Response Format Analysis**
   - **Question**: How does the VolumeTypes List API JSON response map to `[]string` array?
   - **Investigation**: Examine `swagger/vps.yaml` → `pb.VolumeTypeListOutput` structure
   - **Expected Outcome**: Response wrapper structure with `volume_types` array field

3. **Existing VPS Client Integration Pattern**
   - **Question**: How to integrate new Volumes/VolumeTypes clients into existing VPS client?
   - **Investigation**: Review `modules/vps/client.go` accessor pattern (Flavors(), Keypairs(), etc.)
   - **Expected Outcome**: Add `Volumes()` and `VolumeTypes()` accessor methods

4. **Error Handling Consistency**
   - **Question**: What error wrapping patterns are used in existing VPS modules?
   - **Investigation**: Review error handling in `modules/vps/flavors/client.go`, `modules/vps/keypairs/client.go`
   - **Expected Outcome**: Standard `fmt.Errorf("failed to {action} {resource}: %w", err)` pattern

5. **Contract Test Framework**
   - **Question**: How are Swagger contract tests structured in the project?
   - **Investigation**: Search for existing contract test examples in `tests/contract/`
   - **Expected Outcome**: Test framework pattern for validating responses against Swagger schema

### Research Deliverable: `research.md`

Document findings for each research task with:
- **Decision**: Chosen approach
- **Rationale**: Why this approach fits the project
- **Alternatives Considered**: Other options evaluated
- **Code References**: Point to existing examples in codebase

---

## Phase 1: Design & Implementation Structure ✅

**Prerequisites**: Phase 0 research complete, all NEEDS CLARIFICATION resolved

### 1.1 Data Model Design (`data-model.md`) ✅

Extract and document entity structures from `swagger/vps.yaml`:

#### Volume Entity
- **Source**: `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.VolumeInfo`
- **Fields**:
  - `ID` (string): Unique volume identifier
  - `Name` (string): Volume name
  - `Description` (string, optional): Volume description
  - `Size` (int): Volume size in GB
  - `Type` (string): Volume type (SSD, HDD, etc.)
  - `Status` (VolumeStatus): Volume status enum (VolumeStatusAvailable, VolumeStatusInUse, etc.)
  - `StatusReason` (string, optional): Status reason text
  - `Attachments` ([]common.IDName): List of servers this volume is attached to
  - `Project` (common.IDName): Project reference
  - `ProjectID` (string): Project ID
  - `User` (common.IDName): User reference
  - `UserID` (string): User ID
  - `Namespace` (string): Namespace
  - `CreatedAt` (*time.Time): Creation timestamp
  - `UpdatedAt` (*time.Time): Last update timestamp

#### CreateVolumeRequest
- **Source**: `VolumeCreateInput` from swagger
- **Fields**:
  - `Name` (string, required): Volume name
  - `Type` (string, required): Volume type
  - `Size` (int, optional): Size in GB (validation server-side)
  - `Description` (string, optional): Description
  - `SnapshotID` (string, optional): Create from snapshot

#### UpdateVolumeRequest
- **Source**: `VolumeUpdateInput` from swagger
- **Fields**:
  - `Name` (string, optional): New name
  - `Description` (string, optional): New description

#### VolumeActionRequest
- **Source**: `VolActionInput` from swagger
- **Fields**:
  - `Action` (VolumeAction, required): Action type enum (VolumeActionAttach, VolumeActionDetach, VolumeActionExtend, VolumeActionRevert)
  - `ServerID` (string, optional): Required for attach/detach
  - `NewSize` (int, optional): Required for extend

#### ListVolumesOptions
- **Query Parameters**:
  - `Name` (string, optional): Filter by name
  - `UserID` (string, optional): Filter by user ID
  - `Status` (string, optional): Filter by status
  - `Type` (string, optional): Filter by type
  - `Detail` (bool, optional): Include attachment details

#### VolumeListResponse
- **Source**: `pb.VolumeListOutput`
- **Structure**: Wrapper with `Volumes` field containing `[]*Volume`

#### VolumeTypeListResponse
- **Source**: `pb.VolumeTypeListOutput`
- **Structure**: Wrapper with `VolumeTypes` field containing `[]string`

### 1.2 API Contracts (`contracts/`) ✅

Extract endpoint specifications from `swagger/vps.yaml` (volumes and volume_types tags):

**File**: `contracts/volumes.md` ✅
- 6 endpoints with volumes tag:
- GET    /api/v1/project/{project-id}/volumes
- POST   /api/v1/project/{project-id}/volumes
- GET    /api/v1/project/{project-id}/volumes/{vol-id}
- PUT    /api/v1/project/{project-id}/volumes/{vol-id}
- DELETE /api/v1/project/{project-id}/volumes/{vol-id}
- POST   /api/v1/project/{project-id}/volumes/{vol-id}/action

**File**: `contracts/volume_types.md` ✅
- 1 endpoint with volume_types tag:
- GET /api/v1/project/{project-id}/volume_types

### 1.3 Quickstart Guide (`quickstart.md`) ✅

Create developer-focused examples:

```go
// Initialize VPS client
client, _ := cloudsdk.New("https://api.example.com", "bearer-token")
projectClient, _ := client.Project(ctx, "project-id")
vpsClient := projectClient.VPS()

// List volume types
volumeTypesClient := vpsClient.VolumeTypes()
types, _ := volumeTypesClient.List(ctx)
fmt.Println("Available types:", types) // ["SSD", "HDD", "NVMe"]

// Create volume
volumesClient := vpsClient.Volumes()
volume, _ := volumesClient.Create(ctx, &volumes.CreateInput{
    Name: "my-volume",
    Type: "SSD",
    Size: 100, // GB
})

// List volumes with filters
vols, _ := volumesClient.List(ctx, &volumes.ListOptions{
    Status: "available",
    Detail: true,
})

// Attach volume to server
_ = volumesClient.Action(ctx, volume.ID, &volumes.ActionInput{
    Action:   "attach",
    ServerID: "server-123",
})

// Extend volume
_ = volumesClient.Action(ctx, volume.ID, &volumes.ActionInput{
    Action:  "extend",
    NewSize: 200,
})

// Delete volume
_ = volumesClient.Delete(ctx, volume.ID)
```

### 1.4 Agent Context Update ✅

Run `.specify/scripts/bash/update-agent-context.sh copilot` to update `.github/copilot-instructions.md`:
- Add "Go 1.21+ + Go standard library (net/http, encoding/json, context)" to Active Technologies
- Add "N/A (API client SDK)" to Active Technologies (no persistence)
- Preserve existing manual additions between markers

---

## Phase 1 Validation Gate ✅

After completing Phase 1 design artifacts:

1. **Re-run Constitution Check** ✅:
   - Verified all public APIs will accept `context.Context`
   - Verified typed responses and wrapped errors in design
   - Verified no raw HTTP exposure in public interfaces

2. **Run `make check`** ✅:
   - Go formatting validated
   - Linters passed
   - All existing tests passing

3. **Review Against Spec** ✅:
   - All 7 API endpoints covered in contracts
   - All functional requirements (FR-001 through FR-020) addressed in data model
   - Success criteria measurable via contract tests

**Gate Pass Criteria**: ✅ All checks passed, no constitution violations, design complete

**Phase 1 Completed**: 2025-11-13

---

## Phase 2: Task Breakdown (executed by `/speckit.tasks`)

*Note: Phase 2 is NOT performed by `/speckit.plan`. Use `/speckit.tasks` command to generate detailed task list.*

Expected task categories:
1. Models implementation (`models/vps/volumes/`, `models/vps/volumetypes/`)
2. Client implementation (`modules/vps/volumes/`, `modules/vps/volumetypes/`)
3. VPS client integration (`modules/vps/client.go`)
4. Unit tests (85%+ coverage target)
5. Contract tests (validate against Swagger)
6. Integration tests (end-to-end workflows)
7. Documentation updates

---

## Verification Steps

Each phase must pass `make check` before proceeding:

```bash
# After Phase 0 (research)
make check  # Should pass (no code changes yet)

# After Phase 1 (design artifacts)
make check  # Should pass (only documentation added)

# After Phase 2 (implementation - via /speckit.tasks)
make check          # Lint, format, vet
make test           # All unit tests
make test-coverage  # Verify 85%+ coverage
```

---

## Implementation Notes

### Response Unmarshaling

Volumes List API returns:
```json
{
  "volumes": [
    { "id": "vol-1", "name": "volume1", ... },
    { "id": "vol-2", "name": "volume2", ... }
  ]
}
```

Client must unmarshal to intermediate `VolumeListResponse` struct, then return `[]*Volume`:

```go
var response volumes.VolumeListResponse
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to list volumes: %w", err)
}
return response.Volumes, nil
```

VolumeTypes List API returns:
```json
{
  "volume_types": ["SSD", "HDD", "NVMe"]
}
```

Client must unmarshal to intermediate `VolumeTypeListResponse` struct, then return `[]string`:

```go
var response volumetypes.VolumeTypeListResponse
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to list volume types: %w", err)
}
return response.VolumeTypes, nil
```

### IDName Usage

Volume model uses `common.IDName` for references:
```go
import "github.com/Zillaforge/cloud-sdk/models/vps/common"

type Volume struct {
    // ...
    Attachments []common.IDName `json:"attachments,omitempty"`
    Project     common.IDName   `json:"project"`
    User        common.IDName   `json:"user"`
}
```

### Error Handling Consistency

Follow existing VPS module patterns:
```go
// Wrap errors with context
return nil, fmt.Errorf("failed to create volume: %w", err)
return nil, fmt.Errorf("failed to get volume %s: %w", volumeID, err)
return fmt.Errorf("failed to delete volume %s: %w", volumeID, err)
```

### Test Coverage Target

- Unit tests: 85%+ coverage
- All public methods tested
- Error paths tested
- Validation logic tested
- Edge cases covered (see spec.md Edge Cases section)

---

## Dependencies

- `internal/http` - HTTP client with retry and error handling
- `internal/types` - Common SDK types (Logger, SDKError)
- `internal/backoff` - Retry strategy (exponential backoff)
- `models/vps/common` - Shared types (IDName)
- `swagger/vps.yaml` - API contract specification (reference only)

---

## Success Criteria Mapping

| Success Criterion | Verification Method |
|-------------------|---------------------|
| SC-001: List volume types <2s | Performance test |
| SC-002: Create volume <5s | Performance test |
| SC-003: List volumes <3s | Performance test |
| SC-004: Get volume <2s | Performance test |
| SC-005: Update volume <3s | Performance test |
| SC-006: Delete volume <5s | Performance test |
| SC-007: Volume actions <10s | Performance test |
| SC-008: 100% API coverage | Contract tests |
| SC-009: Context propagation <1s | Unit tests |
| SC-010: Clear error messages | Unit tests |
| SC-011: Auth error distinction | Unit tests |
| SC-012: Handle large lists | Integration test |
| QG-001: 85%+ coverage | `make test-coverage` |
| QG-002: Contract tests pass | `make test-contract` |
| QG-003: Integration tests pass | `make test-integration` |
| QG-004: Error handling tests | Unit tests |
| QG-005: Context cancellation | Unit tests |
| QG-006: Documentation | Code examples in docs |

---

## Timeline Estimate

- Phase 0 (Research): 2-3 hours
- Phase 1 (Design): 3-4 hours
- Phase 2 (Implementation - via `/speckit.tasks`): 12-16 hours
- Total: ~18-23 hours

---

## Next Steps

1. Complete Phase 0 research tasks → generate `research.md`
2. Complete Phase 1 design artifacts → generate `data-model.md`, `contracts/`, `quickstart.md`
3. Update agent context → run update script
4. Re-validate Constitution Check
5. Run `make check` to ensure design phase complete
6. Execute `/speckit.tasks` to generate detailed task breakdown for Phase 2 implementation
