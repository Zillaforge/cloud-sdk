# Research Findings: Fix Floating IP Model

**Date**: November 11, 2025  
**Feature**: 007-fix-floating-ip-model  
**Status**: Complete

## Overview

All research items addressed through specification clarifications and Swagger analysis. No technical blockers identified.

---

## Decision 1: List Response Structure

**Decision**: Use `pb.FIPListOutput` structure with `floatingips` array (breaking change, immediate removal of "items" field)

**Rationale**:
- Aligns SDK with actual API contract in vps.yaml
- Confirmed by clarification Q1: Option B (immediate breaking change)
- Enables direct indexing: `response.FloatingIPs[i]` instead of wrapper
- Matches pattern in sibling modules (flavors, keypairs use pb.*ListOutput)

**Alternatives Considered**:
- Option A (dual-field transition period): Higher complexity, longer maintenance burden
- Option C (permanent alias): Perpetuates API divergence, confuses users

**Source**: vps.yaml `pb.FIPListOutput` definition

---

## Decision 2: Timestamp Format

**Decision**: ISO 8601 / RFC3339 string format (e.g., "2025-11-11T10:30:00Z")

**Rationale**:
- Standard Go practice for timestamp marshaling/unmarshaling
- Confirmed by clarification Q2
- Matches keypair model pattern (spec 006)
- Human-readable for debugging
- Supports accurate lifecycle tracking (createdAt, updatedAt, approvedAt)

**Alternatives Considered**:
- Unix timestamps: Less human-readable, harder to debug
- Custom format: Adds unnecessary parsing logic
- Time struct: Not part of API contract

**Source**: vps.yaml `pb.FloatingIPInfo` field types (all strings), clarification session

---

## Decision 3: Null Object Handling

**Decision**: Preserve API structure. Leave project/user IDName objects null/empty when not provided by API

**Rationale**:
- Confirmed by clarification Q3: Option A
- Prevents SDK-level assumptions about data presence
- Maintains contract fidelity with actual API responses
- Allows application code to handle optional references appropriately
- Consistent with how Go handles omitted struct fields (nil pointers)

**Alternatives Considered**:
- Option B (return error): Breaks valid use cases where user/project may not be populated
- Option C (populate minimal objects): Invents data not in API response, violating contract

**Implementation**: Use pointer types for optional nested objects (*IDName)

**Source**: Clarification session Q3

---

## Decision 4: Status Field Handling

**Decision**: No helper methods for status checks. Expose status as custom `FloatingIPStatus` enum type with constants instead of raw string.

**Rationale**:
- Confirmed by clarification Q4: Option B (no helper methods)
- Custom enum provides type safety and IDE autocompletion
- Constants prevent typos and invalid status values
- Easier to add status-related utilities in the future without breaking API
- Matches Go best practices for restricted value sets

**Implementation**: 
```go
type FloatingIPStatus string

const (
    FloatingIPStatusActive   FloatingIPStatus = "ACTIVE"
    FloatingIPStatusPending  FloatingIPStatus = "PENDING"
    FloatingIPStatusDown     FloatingIPStatus = "DOWN"
    FloatingIPStatusRejected FloatingIPStatus = "REJECTED"
)
```

**Alternatives Considered**:
- Option: Raw string field (initial design) - Less type-safe, easier to misuse
- Option: Helper methods (IsActive(), IsPending(), etc.) - Violates Q4 clarification, adds surface area

**Source**: Clarification session Q4, Go custom type best practices

---

## Decision 5: Reserved Field Behavior

**Decision**: Reserved field is read-only. Ignore if provided in create request, only returned by API in responses

**Rationale**:
- Confirmed by clarification Q5: Option A
- Reserved status is API-controlled (admin allocation indicator)
- Users cannot set reservation during creation
- Prevents expectation mismatch
- Matches API contract behavior

**Alternatives Considered**:
- Option B (user-settable during creation): Violates API semantics
- Option C (validation error if missing): Unnecessary strictness

**Implementation**: Exclude Reserved from FloatingIPCreateRequest; include in FloatingIPUpdateRequest but document as ignored

**Source**: Clarification session Q5, vps.yaml FIPUpdateInput

---

## Field Mapping: pb.FloatingIPInfo → FloatingIP Struct

| API Field | Go Field | JSON Tag | Type | Notes |
|-----------|----------|----------|------|-------|
| address | Address | address | string | Floating IP address |
| approvedAt | ApprovedAt | approvedAt | string | RFC3339 timestamp |
| createdAt | CreatedAt | createdAt | string | RFC3339 timestamp |
| description | Description | description | string | omitempty if not set |
| device_id | DeviceID | device_id | string | omitempty for unassociated |
| device_name | DeviceName | device_name | string | omitempty |
| device_type | DeviceType | device_type | string | omitempty |
| extnet_id | ExtNetID | extnet_id | string | omitempty |
| id | ID | id | string | Primary key |
| name | Name | name | string | User-given name |
| namespace | Namespace | namespace | string | omitempty |
| port_id | PortID | port_id | string | omitempty, network port association |
| project | Project | project | *IDName | omitempty if not provided |
| project_id | ProjectID | project_id | string | Project owner |
| reserved | Reserved | reserved | bool | Read-only, not in create requests |
| status | Status | status | string | ACTIVE, PENDING, DOWN, REJECTED |
| status_reason | StatusReason | status_reason | string | omitempty, error details |
| updatedAt | UpdatedAt | updatedAt | string | RFC3339 timestamp, omitempty |
| user | User | user | *IDName | omitempty if not provided |
| user_id | UserID | user_id | string | Resource owner |
| uuid | UUID | uuid | string | Unique identifier |

---

## Request Input Mapping

| Struct | API Reference | Fields | Notes |
|--------|---------------|--------|-------|
| FloatingIPCreateRequest | FIPCreateInput | name, description | Both optional per Swagger |
| FloatingIPUpdateRequest | FIPUpdateInput | name, description, reserved | Reserved ignored during update |
| ListFloatingIPsOptions | Query params | status, user_id, device_type, device_id, extnet_id, address, name, detail | Optional filters |

---

## Testing Approach

### Unit Tests

- **Location**: `models/vps/floatingips/floatingip_test.go`
- **Focus**: JSON marshaling/unmarshaling, struct validation
- **Test Cases**:
  - All fields populate correctly from JSON
  - Null/empty fields handled properly
  - Timestamp strings preserved as-is
  - camelCase JSON tags match API specification

### Contract Tests

- **Location**: `modules/vps/floatingips/test/contract_test.go`
- **Focus**: API contract compliance against vps.yaml
- **Test Cases**:
  - List response uses "floatingips" array (not "items")
  - Create request body matches FIPCreateInput structure
  - All response fields present per pb.FloatingIPInfo
  - Error responses follow SDK error pattern

### Integration Tests

- **Location**: `modules/vps/floatingips/test/integration_test.go`
- **Focus**: End-to-end operation flows with mock server
- **Test Cases**:
  - List with various filters
  - Create with minimal and full request
  - Get non-existent resource (error handling)
  - Update existing resource
  - Delete operation
  - Disassociate operation
  - Error propagation and wrapping

---

## Dependencies Analysis

### Required

- **Go 1.21+**: Already in project (from copilot-instructions.md)
- **encoding/json**: Standard library (no external dep)
- **context**: Standard library (no external dep)
- **net/http**: Standard library (via internal HTTP client)
- **fmt**: Standard library (error wrapping)

### Not Required

- Third-party JSON library (Go standard library sufficient)
- Custom timestamp parser (string handling sufficient per clarification)
- Validation library (simple field presence checks only)

### Existing Abstractions Reused

- `internalhttp.Client`: Existing HTTP abstraction in `internal/http/client.go`
- Error wrapping pattern: Established in sibling modules
- Test infrastructure: Existing test utilities in `modules/vps/*/test/`

---

## Versioning & Migration

### Breaking Change Classification

**Type**: MAJOR version bump (per semantic versioning)

**Change**: List response field name change
- **Old**: `FloatingIPListResponse.Items []*FloatingIP`
- **New**: `FloatingIPListResponse.FloatingIPs []*FloatingIP`

### Migration Guide

Users upgrading must:

```go
// BEFORE (v0.x)
resp, _ := client.List(ctx, opts)
for _, fip := range resp.Items {  // ← "Items" field
    // ...
}

// AFTER (v1.0+)
resp, _ := client.List(ctx, opts)
for _, fip := range resp.FloatingIPs {  // ← "FloatingIPs" field
    // ...
}

// Equivalent calculation
count := len(resp.FloatingIPs)  // Instead of resp.Total
```

### Release Notes Entry

```
## v1.0.0 - Breaking Changes

### Floating IP Model Alignment (007-fix-floating-ip-model)

The Floating IP model has been updated to exactly match the VPS API specification.

**Breaking Changes:**
- `FloatingIPListResponse.Items` renamed to `FloatingIPListResponse.FloatingIPs`
- New fields added to `FloatingIP`: UUID, Namespace, ExtNetID, DeviceName, DeviceType, StatusReason, and others
- Reserved field is now read-only (remove from create/update request code if present)

**Migration:**
- Replace `response.Items` with `response.FloatingIPs`
- Remove `response.Total` usage (use `len(response.FloatingIPs)` instead)
- Remove `Reserved` field from any create/update requests

**Why:** Ensures SDK contract exactly matches the actual API, reducing confusion and bugs.
```

---

## Deliverables Summary

| Item | Status | Location |
|------|--------|----------|
| Clarifications | ✅ Complete | spec.md (5 Q&A pairs) |
| Data Model Design | ✅ Complete | This file + plan.md |
| Field Mapping | ✅ Complete | This file (table) |
| Test Strategy | ✅ Complete | This file + plan.md |
| Dependencies | ✅ Analyzed | This file (stdlib only) |
| Versioning | ✅ Planned | This file + plan.md |
| API Contracts | 📋 To Generate | Phase 1 (contracts/) |

---

## Next Phase

Ready for **Phase 1: Design & Contracts** execution:
1. Generate OpenAPI contract in `contracts/floatingip-openapi.yaml`
2. Create `data-model.md` with detailed field descriptions
3. Create `quickstart.md` with usage examples
4. Update agent context via `update-agent-context.sh`
