# Implementation Plan: Fix Floating IP Model

**Branch**: `007-fix-floating-ip-model` | **Date**: November 11, 2025 | **Spec**: `/specs/007-fix-floating-ip-model/spec.md`
**Input**: Feature specification from `/specs/007-fix-floating-ip-model/spec.md`

## Summary

This feature corrects the FloatingIP model and module implementation in the cloud-sdk to strictly conform to the VPS API specification (pb.FloatingIPInfo and pb.FIPListOutput from vps.yaml). The primary change is updating `models/vps/floatingips/floatingip.go` to include all fields from the API specification and updating `modules/vps/floatingips/client.go` to use proper Request/Response structs. Six operations are implemented: List, Create, Get, Update, Delete, and Disassociate. This is a breaking change with immediate removal of the "items" field in list responses.

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: Go standard library (`encoding/json`, `context`, `net/http`, `fmt`)  
**Storage**: N/A (API client SDK)  
**Testing**: Go `testing` package with unit and contract tests per Swagger/OpenAPI  
**Target Platform**: Linux server (deployed as SDK library)
**Project Type**: SDK client library (single package structure)  
**Performance Goals**: Standard SDK operations (no specific latency constraints documented)  
**Constraints**: Minimal dependencies, TDD mandatory, breaking change documented  
**Scale/Scope**: 6 API operations, 1 model with 19 fields, comprehensive test coverage

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The feature plan MUST satisfy all applicable Cloud SDK Constitution principles:

- ✅ **TDD mandatory**: Unit tests and contract tests will be written first for all 6 operations. Tests will verify model serialization/deserialization and API contracts against vps.yaml specification.
- ✅ **Public API shape**: FloatingIP model exposes idiomatic Go struct fields; Client methods follow pattern `Client.Operation(ctx, params)` returning `(*FloatingIP, error)` or `([]*FloatingIP, error)`; no raw HTTP exposure; all public methods accept `context.Context`.
- ✅ **Dependencies**: Only Go standard library used (`encoding/json`, `context`, `net/http`, `fmt`); internal HTTP client abstractions already in place.
- ✅ **Versioning**: Breaking change (removal of "items" field in ListResponse). Will be marked as MAJOR version bump with migration guide in release notes.
- ✅ **Observability**: Error handling follows established pattern with `fmt.Errorf("failed to {operation}: %w", err)` wrapping for context without forcing vendor.
- ✅ **Security**: No secrets in logs; HTTP client configuration inherited from SDK's centralized HTTP client setup.

## Project Structure

```
# Cloud SDK - Floating IP Fix (007-fix-floating-ip-model)

models/vps/floatingips/
├── floatingip.go                 # MODIFIED: Add all pb.FloatingIPInfo fields
├── floatingip_test.go            # NEW: Unit tests for model (TDD first)

modules/vps/floatingips/
├── client.go                     # MODIFIED: Update API operations with Request/Response structs
├── client_test.go                # MODIFIED: Update tests for new model
├── test/
│   ├── contract_test.go          # NEW: Contract tests against Swagger spec
│   └── integration_test.go        # NEW: Integration tests (if applicable)

specs/007-fix-floating-ip-model/
├── plan.md                       # This file
├── research.md                   # Phase 0: Research findings
├── data-model.md                 # Phase 1: Data model details
├── quickstart.md                 # Phase 1: Usage examples
├── contracts/
│   └── floatingip-openapi.yaml   # Phase 1: API contract from vps.yaml
└── checklists/
    └── requirements.md           # Spec quality checklist
```

**Structure Decision**: Single package structure for SDK client library. Reuses existing `models/vps/floatingips/` and `modules/vps/floatingips/` directories. No new external dependencies. Tests follow established patterns in sibling modules (flavors, keypairs).

## Complexity Tracking

| Item | Status | Notes |
|------|--------|-------|
| Constitution Check | ✅ PASS | All principles satisfied, TDD mandatory, no external dependencies needed |
| Unresolved Clarifications | ✅ NONE | All 5 clarifications from /speckit.clarify session resolved |
| Breaking Changes | ⚠️ DOCUMENTED | List response field renamed "items" → "floatingips"; MAJOR version bump required with migration guide |
| Dependencies | ✅ NONE | Standard library only; existing internal HTTP client abstractions reused |
| Test Coverage | 📋 PLANNED | Unit tests (model serialization), contract tests (Swagger compliance), integration tests (mock API responses) |

---

# Phase 0: Research & Clarifications

**Status**: ✅ COMPLETE (moved to research.md)

All clarifications resolved in `/speckit.clarify` session on 2025-11-11. Key decisions:

1. **List Response**: Breaking change approved. Remove "items", use "floatingips" per pb.FIPListOutput
2. **Timestamp Format**: ISO 8601 / RFC3339 strings (e.g., "2025-11-11T10:30:00Z")
3. **Null Handling**: Preserve API structure; leave project/user null when missing
4. **Status Field**: Read-only string, no helper methods
5. **Reserved Field**: Read-only, excluded from create/update requests

**Research Tasks**: None needed. Technical decisions clear from Swagger spec (vps.yaml) and clarifications.

---

# Phase 1: Design & Contracts

## Data Model

**Source**: `models/vps/floatingips/floatingip.go`

### FloatingIP Resource

From pb.FloatingIPInfo (vps.yaml):

```go
// FloatingIPStatus is the custom enum type for floating IP status
type FloatingIPStatus string

const (
    FloatingIPStatusActive   FloatingIPStatus = "ACTIVE"   // Ready for use
    FloatingIPStatusPending  FloatingIPStatus = "PENDING"  // Awaiting approval
    FloatingIPStatusDown     FloatingIPStatus = "DOWN"     // Service issue
    FloatingIPStatusRejected FloatingIPStatus = "REJECTED" // Rejected by admin
)

type FloatingIP struct {
    // Core identity
    ID              string    `json:"id"`
    UUID            string    `json:"uuid"`
    
    // Allocation details
    Address         string    `json:"address"`
    Name            string    `json:"name"`
    Description     string    `json:"description,omitempty"`
    
    // Association
    PortID          string    `json:"port_id,omitempty"`
    DeviceID        string    `json:"device_id,omitempty"`
    DeviceName      string    `json:"device_name,omitempty"`
    DeviceType      string    `json:"device_type,omitempty"`
    
    // Network info
    ProjectID       string    `json:"project_id"`
    Project         *IDName   `json:"project,omitempty"`
    Namespace       string    `json:"namespace,omitempty"`
    ExtNetID        string    `json:"extnet_id,omitempty"`
    
    // Ownership
    UserID          string    `json:"user_id"`
    User            *IDName   `json:"user,omitempty"`
    
    // State
    Status          FloatingIPStatus `json:"status"`       // Custom enum type
    StatusReason    string           `json:"status_reason,omitempty"`
    Reserved        bool             `json:"reserved"`
    
    // Lifecycle
    CreatedAt       string    `json:"createdAt"`
    UpdatedAt       string    `json:"updatedAt,omitempty"`
    ApprovedAt      string    `json:"approvedAt,omitempty"`
}

type IDName struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

### Request/Response Structs

```go
// List Response (pb.FIPListOutput)
type FloatingIPListResponse struct {
    FloatingIPs []*FloatingIP `json:"floatingips"`
}

// Create Request (FIPCreateInput)
type FloatingIPCreateRequest struct {
    Name        string `json:"name,omitempty"`
    Description string `json:"description,omitempty"`
}

// Create Response = FloatingIP

// Get Response = FloatingIP

// Update Request (FIPUpdateInput)
type FloatingIPUpdateRequest struct {
    Name        string `json:"name,omitempty"`
    Description string `json:"description,omitempty"`
    Reserved    bool   `json:"reserved,omitempty"`  // NOTE: reserved ignored during update per clarification
}

// Update Response = FloatingIP

// Delete Response = nil (no body)

// Disassociate Response = nil (no body)
```

**Validation Rules**:
- Name: required for create, optional for update
- Description: optional for create and update
- Reserved field in requests is ignored per clarifications (read-only)
- Status field read-only; controlled by API

---

## API Contract

**Source**: `/swagger/vps.yaml` definitions (pb.FloatingIPInfo, pb.FIPListOutput, FIPCreateInput, FIPUpdateInput)

### Operations

```
1. LIST   GET    /api/v1/project/{project-id}/floatingips
   Input:  Query params (name, user_id, status, reserved, device_type, device_id, extnet_id, address, detail)
   Output: FloatingIPListResponse []*FloatingIP
   
2. CREATE POST   /api/v1/project/{project-id}/floatingips
   Input:  FloatingIPCreateRequest
   Output: FloatingIP
   
3. GET    GET    /api/v1/project/{project-id}/floatingips/{fip-id}
   Input:  fip-id path param
   Output: FloatingIP
   
4. UPDATE PUT    /api/v1/project/{project-id}/floatingips/{fip-id}
   Input:  FloatingIPUpdateRequest (reserved field ignored)
   Output: FloatingIP
   
5. DELETE DELETE /api/v1/project/{project-id}/floatingips/{fip-id}
   Input:  fip-id path param
   Output: (empty)
   
6. DISASSOCIATE POST /api/v1/project/{project-id}/floatingips/{fip-id}/disassociate
   Input:  fip-id path param
   Output: (empty)
```

---

## Module Implementation

**Source**: `modules/vps/floatingips/client.go`

### Client Methods Signature

```go
type Client struct {
    baseClient *internalhttp.Client
    projectID  string
    basePath   string
}

// List(ctx, opts) -> ([]*FloatingIP, error)
// Create(ctx, req) -> (*FloatingIP, error)
// Get(ctx, fipID) -> (*FloatingIP, error)
// Update(ctx, fipID, req) -> (*FloatingIP, error)
// Delete(ctx, fipID) -> error
// Disassociate(ctx, fipID) -> error
```

**Error Handling Pattern**:
```go
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to {operation}: %w", err)
}
```

---

## Test Strategy

### Unit Tests (`floatingip_test.go`)

- JSON marshaling/unmarshaling with all fields
- Null/omitempty field handling
- Timestamp string validation
- Project/User IDName struct handling

### Contract Tests (`modules/vps/floatingips/test/contract_test.go`)

- Verify response structures match vps.yaml definitions
- Verify request encoding matches FIPCreateInput/FIPUpdateInput
- Verify list response uses "floatingips" array (not "items")
- Verify field names use camelCase per spec

### Integration Tests (`modules/vps/floatingips/test/integration_test.go`)

- Mock all 6 operations (List, Create, Get, Update, Delete, Disassociate)
- Verify proper error wrapping with context
- Verify context propagation in all operations

---

## Validation Gate

✅ Constitution Check PASSED:
- TDD approach: tests written first ✅
- Idiomatic Go: Client methods pattern ✅
- No external deps: stdlib only ✅
- Breaking change documented: MAJOR version bump + migration guide ✅
- Observability: error wrapping pattern ✅
- Security: no secret logging ✅

**Ready for Phase 2 (Tasks)**: YES
