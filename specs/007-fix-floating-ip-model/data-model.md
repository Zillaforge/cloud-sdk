# Data Model: Floating IP Resource

**Date**: November 11, 2025  
**Feature**: 007-fix-floating-ip-model  
**Source**: pb.FloatingIPInfo from vps.yaml  

---

## Entity Overview

### FloatingIP

Represents a floating IP address resource allocated to a project and optionally associated with a network port or device.

**Responsibility**: Store and communicate floating IP allocation state including identity, association, and lifecycle metadata.

**Lifecycle**: Created → Active → Optionally Associated → Deleted

---

## Core Fields

### Identity

| Field | Type | JSON Tag | Required | Constraints |
|-------|------|----------|----------|-------------|
| `ID` | string | `id` | Yes | Unique within project; immutable |
| `UUID` | string | `uuid` | Yes | Global unique identifier; immutable |
| `Name` | string | `name` | Yes | User-given display name; 1-255 chars |

### Network & Allocation

| Field | Type | JSON Tag | Required | Constraints |
|-------|------|----------|----------|-------------|
| `Address` | string | `address` | Yes | IPv4 address; format: "a.b.c.d" |
| `ExtNetID` | string | `extnet_id` | No | External network where allocated |
| `PortID` | string | `port_id` | No | Associated network port (if associated) |
| `ProjectID` | string | `project_id` | Yes | Owning project; immutable |
| `Namespace` | string | `namespace` | No | Deployment namespace context |

### Ownership & References

| Field | Type | JSON Tag | Required | Constraints |
|-------|------|----------|----------|-------------|
| `UserID` | string | `user_id` | Yes | Owning user ID; immutable |
| `User` | *IDName | `user` | No | User object reference (optional per API) |
| `Project` | *IDName | `project` | No | Project object reference (optional per API) |

### Association Details

| Field | Type | JSON Tag | Required | Constraints |
|-------|------|----------|----------|-------------|
| `DeviceID` | string | `device_id` | No | Associated device/VM ID (if associated) |
| `DeviceName` | string | `device_name` | No | Associated device name |
| `DeviceType` | string | `device_type` | No | Type of associated device (e.g., "server", "lb") |

### Status & Metadata

| Field | Type | JSON Tag | Required | Constraints |
|-------|------|----------|----------|-------------|
| `Description` | string | `description` | No | User-provided description; 0-1000 chars |
| `Status` | FloatingIPStatus | `status` | Yes | ACTIVE \| PENDING \| DOWN \| REJECTED; custom enum; read-only |
| `StatusReason` | string | `status_reason` | No | Error details if status=REJECTED |
| `Reserved` | bool | `reserved` | Yes | True if reserved for special use; read-only |

### Lifecycle Timestamps

| Field | Type | JSON Tag | Format | Constraints |
|-------|------|----------|--------|-------------|
| `CreatedAt` | string | `createdAt` | RFC3339 | When allocated; immutable |
| `UpdatedAt` | string | `updatedAt` | RFC3339 | When last modified; omitted if unchanged |
| `ApprovedAt` | string | `approvedAt` | RFC3339 | When approved (if applicable); omitted if not approved |

**Timestamp Format Example**: `2025-11-11T10:30:00Z`

---

## Supporting Types

### FloatingIPStatus (Custom Enum)

Status of a floating IP resource with restricted set of valid values.

```go
type FloatingIPStatus string

const (
    FloatingIPStatusActive   FloatingIPStatus = "ACTIVE"   // Floating IP is allocated and ready for use
    FloatingIPStatusPending  FloatingIPStatus = "PENDING"  // Floating IP awaiting approval
    FloatingIPStatusDown     FloatingIPStatus = "DOWN"     // Service issue, not available
    FloatingIPStatusRejected FloatingIPStatus = "REJECTED" // Request rejected by admin
)

// String returns the string representation of the status
func (s FloatingIPStatus) String() string {
    return string(s)
}

// Valid returns true if the status is a valid FloatingIPStatus value
func (s FloatingIPStatus) Valid() bool {
    switch s {
    case FloatingIPStatusActive, FloatingIPStatusPending, FloatingIPStatusDown, FloatingIPStatusRejected:
        return true
    default:
        return false
    }
}
```

**JSON Handling**: 
- Marshaling: `FloatingIPStatus` marshals to JSON string (e.g., `"ACTIVE"`)
- Unmarshaling: JSON string values are validated and converted to `FloatingIPStatus` constant
- Invalid values: Deserialization error if API returns unknown status value

**Usage**: Use in FloatingIP model's Status field; IDE autocompletion for valid constants.

### IDName

Reference to a resource by ID and name.

```go
type IDName struct {
    ID   string `json:"id"`    // Unique identifier
    Name string `json:"name"`  // Display name
}
```

**Usage**: Optional project and user references in FloatingIP. May be null/omitted if not provided by API.

---

## Relationships

```
FloatingIP
├─ Project (1:1) [optional reference via project_id + Project object]
├─ User (1:1) [via user_id + optional User object]
├─ Device (0..1) [optional via device_id, device_name, device_type]
├─ Port (0..1) [optional via port_id]
└─ ExternalNetwork (1:1) [via extnet_id]
```

---

## State Transitions

```
Created
   ↓
┌──────────────┐
│   PENDING    │ ← Awaiting approval (if required)
└──────────────┘
   ↓
┌──────────────┐
│   ACTIVE     │ ← Ready for use
└──────────────┘
   │      │
   │      └─→ Associated with Port/Device (device_id, port_id populate)
   │
   └─→ DOWN ← Service issue
   │
   └─→ REJECTED ← Admin rejected (status_reason populated)
   
   ↓
Deleted (resource removed)
```

---

## Validation Rules

### Create Request (FloatingIPCreateRequest)

| Field | Validation |
|-------|-----------|
| `name` | Optional; if provided 1-255 chars, no null bytes |
| `description` | Optional; if provided 0-1000 chars |

### Update Request (FloatingIPUpdateRequest)

| Field | Validation |
|-------|-----------|
| `name` | Optional; if provided 1-255 chars, no null bytes |
| `description` | Optional; if provided 0-1000 chars |
| `reserved` | Present in struct but **ignored by SDK** (read-only per API) |

### Response Fields

| Field | Validation |
|-------|-----------|
| `status` | Must be one of: FloatingIPStatusActive, FloatingIPStatusPending, FloatingIPStatusDown, FloatingIPStatusRejected |
| All identity fields | Non-empty string |
| Timestamp fields | Valid RFC3339 format or omitted if not set |
| Optional associations | Null pointer if not set by API response |
| Reserved | Boolean (always present in response) |

---

## Field Behavior Across Operations

### LIST Operation

**Request**: Supports optional query filters (not modeled in go struct, handled by HTTP client)
**Response**: `FloatingIPListResponse` with array of FloatingIP (field: `floatingips`, not `items`)
**Fields Returned**: All fields populated by API (some may be null/empty)

### CREATE Operation

**Request**: FloatingIPCreateRequest (name, description only)
**Response**: FloatingIP (all fields including status=PENDING initially)
**Field Behavior**:
- Name from request, Description from request
- ID, UUID, ProjectID, UserID generated by API
- Status=PENDING, Reserved=false initially
- CreatedAt populated by API
- Device fields empty (not associated yet)

### GET Operation

**Request**: FloatingIP ID path parameter
**Response**: FloatingIP (full detail)
**Field Behavior**: Same as LIST, single resource

### UPDATE Operation

**Request**: FloatingIPUpdateRequest (name, description, reserved)
**Response**: FloatingIP (updated state)
**Field Behavior**:
- Name, Description updated if provided
- **Reserved field ignored** (read-only per clarification)
- UpdatedAt set by API to current timestamp
- Status may change due to API backend processing

### DELETE Operation

**Request**: FloatingIP ID path parameter
**Response**: (empty, HTTP 204)
**Effect**: Resource removed from project

### DISASSOCIATE Operation

**Request**: FloatingIP ID path parameter
**Response**: (empty, HTTP 204)
**Effect**: Detaches FloatingIP from associated device/port
- **Before**: device_id, device_name, device_type populated; port_id populated
- **After**: device_id, device_name, device_type emptied; port_id emptied; status may change

---

## Edge Cases & Null Handling

### Null/Empty Scenarios

| Scenario | Behavior | Model Representation |
|----------|----------|---------------------|
| User object not included in API response | SDK accepts it | `User` field is nil pointer |
| Project object not included in API response | SDK accepts it | `Project` field is nil pointer |
| FloatingIP not associated with device | Normal state | `DeviceID`, `DeviceName`, `DeviceType` are empty strings |
| UpdatedAt not set (no updates) | API omits field | JSON unmarshaling: empty string or omitted |
| ApprovedAt only if approved flow | Conditional | JSON unmarshaling: empty string or omitted if not applicable |
| StatusReason empty if status != REJECTED | Expected | JSON unmarshaling: empty string or omitted |

### Error Scenarios

| Error | HTTP Status | SDK Handling |
|-------|-------------|-----|
| Malformed ID | 400 | Propagated with context: "failed to get floatingips: {error}" |
| FloatingIP not found | 404 | Propagated with context |
| Permission denied | 403 | Propagated with context |
| Timestamp parse error | 500 (from API) or deserialization error | SDK error: malformed timestamp in field X |

---

## Example Instances

### Allocated but Unassociated

```json
{
  "id": "fip-abc123",
  "uuid": "12345678-1234-1234-1234-123456789012",
  "address": "203.0.113.45",
  "name": "web-server-fip",
  "description": "Floating IP for web server",
  "projectId": "proj-001",
  "project": {"id": "proj-001", "name": "Production"},
  "userId": "user-xyz",
  "user": {"id": "user-xyz", "name": "alice@example.com"},
  "status": "ACTIVE",
  "reserved": false,
  "createdAt": "2025-11-10T15:30:00Z",
  "updatedAt": "2025-11-10T15:30:00Z"
}
```

### Associated with Device

```json
{
  "id": "fip-def456",
  "uuid": "87654321-4321-4321-4321-210987654321",
  "address": "203.0.113.46",
  "name": "lb-fip",
  "description": "Load balancer floating IP",
  "projectId": "proj-001",
  "project": {"id": "proj-001", "name": "Production"},
  "portId": "port-123",
  "deviceId": "lb-456",
  "deviceName": "web-lb-prod",
  "deviceType": "lb",
  "userId": "user-xyz",
  "user": {"id": "user-xyz", "name": "alice@example.com"},
  "status": "ACTIVE",
  "reserved": false,
  "createdAt": "2025-11-09T10:00:00Z",
  "updatedAt": "2025-11-10T14:22:00Z"
}
```

### Pending Approval

```json
{
  "id": "fip-ghi789",
  "uuid": "11111111-2222-3333-4444-555555555555",
  "address": "203.0.113.47",
  "name": "staging-fip",
  "description": "Floating IP for staging environment",
  "projectId": "proj-002",
  "userId": "user-abc",
  "status": "PENDING",
  "statusReason": "Awaiting network admin approval",
  "reserved": false,
  "createdAt": "2025-11-11T09:00:00Z"
}
```

### Minimal Response (sparse API response)

```json
{
  "id": "fip-jkl012",
  "address": "203.0.113.48",
  "projectId": "proj-001",
  "userId": "user-xyz",
  "status": "ACTIVE",
  "reserved": false,
  "createdAt": "2025-11-11T10:00:00Z"
}
```

---

## Constraints & Invariants

1. **Identity**: ID and UUID are immutable and unique within the system
2. **Ownership**: ProjectID and UserID never change after creation
3. **Association**: A FloatingIP can be associated with at most one device/port at a time
4. **Status**: Only API can change status; SDK only reads it
5. **Reserved**: Only API admin can set/change; user requests ignore this field
6. **Timestamps**: CreatedAt never changes; UpdatedAt only moves forward; ApprovedAt only set once

---

## Design Decisions (Rationale)

| Decision | Rationale |
|----------|-----------|
| Use string for timestamps | Matches API contract; Go stdlib handles RFC3339 parsing in tests/app code |
| Pointer for optional *IDName | Go idiom; nil indicates "not provided by API" |
| Reserved as read-only in SDK | Per clarification Q5; prevents user confusion about what they can control |
| No status helper methods | Per clarification Q4; status interpretation is application concern |
| Direct field exposure (not getters) | Idiomatic Go; standard practice in SDK models |
| "FloatingIPs" array in list response | Per pb.FIPListOutput specification and clarification Q1 |

---

## Summary for Developers

The FloatingIP model represents a network resource with:
- **Identity**: ID, UUID, Name (user-assigned)
- **Network**: Address (IPv4), ExtNetID, associated Port
- **Ownership**: ProjectID, UserID (immutable)
- **Association**: Optional link to a device via DeviceID/DeviceName/DeviceType
- **Lifecycle**: CreatedAt → UpdatedAt → ApprovedAt (optional), with Status transitions
- **Control**: Description and Name can be updated; most other fields are read-only
- **Metadata**: StatusReason for error details, Reserved flag for special allocations

All fields are directly accessible on the struct with no helper methods; status interpretation and business logic belongs in application code.
