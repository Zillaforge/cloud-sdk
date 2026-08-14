# Data Model: Security Groups and Rules

**Feature**: Fix Security Group Model (001-fix-security-group-model)  
**Date**: November 9, 2025  
**Related**: [spec.md](./spec.md) | [plan.md](./plan.md) | [research.md](./research.md)

## Overview

This document defines the exact data structures for security group resources, matching the VPS API specification (`swagger/vps.yaml`). All models conform to `pb.SgInfo`, `pb.SgRuleInfo`, `pb.SgListOutput`, `SgCreateInput`, `SgUpdateInput`, and `SgRuleCreateInput` definitions.

---

## Model Hierarchy

```
SecurityGroup
├── ID                     (string)
├── Name                   (string)
├── Description            (string)
├── ProjectID              (string)
├── UserID                 (string)
├── Namespace              (string)
├── Rules                  ([]SecurityGroupRule)
├── CreatedAt              (time.Time)
├── UpdatedAt              (time.Time)
├── Project (IDName)       (nested object)
└── User (IDName)          (nested object)

SecurityGroupRule
├── ID                     (string)
├── Direction              (string: "ingress" | "egress")
├── Protocol               (string: "tcp" | "udp" | "icmp")
├── PortMin                (int)
├── PortMax                (int)
└── RemoteCIDR             (string: CIDR notation)
```

---

## Custom Types (Type-Safe Enums)

### Protocol Type

**Location**: `models/vps/securitygroups/rule.go`

```go
package securitygroups

// Protocol represents the network protocol for a security group rule.
type Protocol string

const (
    ProtocolTCP  Protocol = "tcp"
    ProtocolUDP  Protocol = "udp"
    ProtocolICMP Protocol = "icmp"
    ProtocolAny  Protocol = "any"
)
```

**Values**:
- `tcp`: Transmission Control Protocol (connection-oriented)
- `udp`: User Datagram Protocol (connectionless)
- `icmp`: Internet Control Message Protocol (ping, traceroute)
- `any`: All protocols (wildcard)

**Rationale**: Type-safe enum prevents typos and provides IDE autocomplete. Compile-time validation ensures only valid protocol values are used.

### Direction Type

**Location**: `models/vps/securitygroups/rule.go`

```go
package securitygroups

// Direction represents the traffic direction for a security group rule.
type Direction string

const (
    DirectionIngress Direction = "ingress"
    DirectionEgress  Direction = "egress"
)
```

**Values**:
- `ingress`: Inbound traffic (from remote CIDR to instances)
- `egress`: Outbound traffic (from instances to remote CIDR)

**Rationale**: Type-safe enum prevents incorrect direction values. Clear semantic meaning improves code readability.

---

## Core Models

### SecurityGroup (Response Model)

**Location**: `models/vps/securitygroups/securitygroup.go`

**API Mapping**: `pb.SgInfo` from `swagger/vps.yaml`

```go
package securitygroups

import "time"

// SecurityGroup represents a security group resource.
type SecurityGroup struct {
    ID          string               `json:"id"`
    Name        string               `json:"name"`
    Description string               `json:"description"`
    ProjectID   string               `json:"project_id"`
    UserID      string               `json:"user_id"`
    Namespace   string               `json:"namespace"`
    Rules       []SecurityGroupRule  `json:"rules"`
    CreatedAt   time.Time            `json:"createdAt"`
    UpdatedAt   time.Time            `json:"updatedAt"`
    Project     *IDName              `json:"project,omitempty"`
    User        *IDName              `json:"user,omitempty"`
}

// IDName represents a nested identifier-name pair.
type IDName struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

**Field Descriptions**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `ID` | string | Yes | Unique security group identifier (e.g., "sg-abc123") |
| `Name` | string | Yes | Security group name (user-defined) |
| `Description` | string | No | Optional description |
| `ProjectID` | string | Yes | Parent project identifier |
| `UserID` | string | Yes | Owner user identifier |
| `Namespace` | string | Yes | Kubernetes-style namespace (e.g., "default") |
| `Rules` | []SecurityGroupRule | Yes | Array of security rules (empty if detail=false) |
| `CreatedAt` | time.Time | Yes | Resource creation timestamp (ISO 8601) |
| `UpdatedAt` | time.Time | Yes | Last modification timestamp (ISO 8601) |
| `Project` | *IDName | No | Nested project details (id + name) |
| `User` | *IDName | No | Nested user details (id + name) |

**Example JSON Response**:
```json
{
  "id": "sg-f4e3d2c1",
  "name": "web-servers",
  "description": "Security group for web servers",
  "project_id": "proj-123",
  "user_id": "user-456",
  "namespace": "default",
  "rules": [
    {
      "id": "rule-001",
      "direction": "ingress",
      "protocol": "tcp",
      "port_min": 80,
      "port_max": 80,
      "remote_cidr": "0.0.0.0/0"
    },
    {
      "id": "rule-002",
      "direction": "ingress",
      "protocol": "tcp",
      "port_min": 443,
      "port_max": 443,
      "remote_cidr": "0.0.0.0/0"
    }
  ],
  "createdAt": "2025-11-09T10:00:00Z",
  "updatedAt": "2025-11-09T10:00:00Z",
  "project": {
    "id": "proj-123",
    "name": "production"
  },
  "user": {
    "id": "user-456",
    "name": "admin@example.com"
  }
}
```

---

### SecurityGroupRule (Response Model)

**Location**: `models/vps/securitygroups/rule.go`

**API Mapping**: `pb.SgRuleInfo` from `swagger/vps.yaml`

```go
package securitygroups

// SecurityGroupRule represents a security group rule.
type SecurityGroupRule struct {
    ID         string    `json:"id"`
    Direction  Direction `json:"direction"`   // DirectionIngress or DirectionEgress
    Protocol   Protocol  `json:"protocol"`    // ProtocolTCP, ProtocolUDP, ProtocolICMP, ProtocolAny
    PortMin    int       `json:"port_min"`    // Starting port (0 for ICMP)
    PortMax    int       `json:"port_max"`    // Ending port (0 for ICMP)
    RemoteCIDR string    `json:"remote_cidr"` // IP range in CIDR notation
}
```

**Field Descriptions**:

| Field | Type | Values | Description |
|-------|------|--------|-------------|
| `ID` | string | - | Unique rule identifier (e.g., "rule-xyz789") |
| `Direction` | Direction | `DirectionIngress`, `DirectionEgress` | Traffic direction (inbound/outbound) |
| `Protocol` | Protocol | `ProtocolTCP`, `ProtocolUDP`, `ProtocolICMP`, `ProtocolAny` | Network protocol |
| `PortMin` | int | 0-65535 | Starting port (0 for ICMP or all ports) |
| `PortMax` | int | 0-65535 | Ending port (0 for ICMP or all ports) |
| `RemoteCIDR` | string | - | Source/destination IP range (e.g., "10.0.0.0/8") |

**Example JSON Response**:
```json
{
  "id": "rule-abc123",
  "direction": "ingress",
  "protocol": "tcp",
  "port_min": 22,
  "port_max": 22,
  "remote_cidr": "192.168.1.0/24"
}
```

**Protocol-Specific Behavior**:

| Protocol | Port Fields | Example Use Case |
|----------|-------------|------------------|
| `ProtocolTCP` | Required | SSH (22), HTTP (80), HTTPS (443) |
| `ProtocolUDP` | Required | DNS (53), NTP (123) |
| `ProtocolICMP` | Ignored (0) | Ping, traceroute |
| `ProtocolAny` | Optional | Allow all protocols |

---

## Request Models

### SecurityGroupCreateRequest

**Location**: `models/vps/securitygroups/securitygroup.go`

**API Mapping**: `SgCreateInput` from `swagger/vps.yaml`

```go
package securitygroups

// SecurityGroupCreateRequest represents a security group creation request.
type SecurityGroupCreateRequest struct {
    Name        string                           `json:"name"`                  // Required
    Description string                           `json:"description,omitempty"` // Optional
    Rules       []SecurityGroupRuleCreateRequest `json:"rules,omitempty"`       // Optional
}
```

**Field Descriptions**:

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `Name` | string | **Yes** | Non-empty | Security group name |
| `Description` | string | No | - | Optional description |
| `Rules` | []SecurityGroupRuleCreateRequest | No | - | Initial rules (can add later via sub-resource) |

**Example JSON Request**:
```json
{
  "name": "database-servers",
  "description": "Security group for PostgreSQL databases",
  "rules": [
    {
      "direction": "ingress",
      "protocol": "tcp",
      "port_min": 5432,
      "port_max": 5432,
      "remote_cidr": "10.0.0.0/16"
    }
  ]
}
```

**Validation Rules**:
- `Name`: Must not be empty
- `Rules`: Each rule must satisfy `SecurityGroupRuleCreateRequest` validation

---

### SecurityGroupUpdateRequest

**Location**: `models/vps/securitygroups/securitygroup.go`

**API Mapping**: `SgUpdateInput` from `swagger/vps.yaml`

```go
package securitygroups

// SecurityGroupUpdateRequest represents a security group update request.
type SecurityGroupUpdateRequest struct {
    Name        *string `json:"name,omitempty"`        // Optional
    Description *string `json:"description,omitempty"` // Optional
}
```

**Field Descriptions**:

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `Name` | *string | No | New security group name (nil = no change) |
| `Description` | *string | No | New description (nil = no change) |

**Notes**:
- All fields are optional (pointer types distinguish "not set" from "empty string")
- At least one field should be non-nil (API may reject empty update)
- Rules cannot be updated via Update endpoint (use sub-resource operations)

**Example JSON Request**:
```json
{
  "name": "web-servers-v2",
  "description": "Updated description for web servers"
}
```

---

### SecurityGroupRuleCreateRequest

**Location**: `models/vps/securitygroups/rule.go`

**API Mapping**: `SgRuleCreateInput` from `swagger/vps.yaml`

```go
package securitygroups

// SecurityGroupRuleCreateRequest represents a rule creation request.
type SecurityGroupRuleCreateRequest struct {
    Direction  Direction `json:"direction"`          // Required: DirectionIngress or DirectionEgress
    Protocol   Protocol  `json:"protocol"`           // Required: ProtocolTCP, ProtocolUDP, ProtocolICMP, ProtocolAny
    PortMin    *int      `json:"port_min,omitempty"` // Optional: starting port (TCP/UDP only)
    PortMax    *int      `json:"port_max,omitempty"` // Optional: ending port (TCP/UDP only)
    RemoteCIDR string    `json:"remote_cidr"`        // Required: IP range in CIDR notation
}
```

**Field Descriptions**:

| Field | Type | Required | Constraints | Description |
|-------|------|----------|-------------|-------------|
| `Direction` | Direction | **Yes** | `DirectionIngress` or `DirectionEgress` | Traffic direction |
| `Protocol` | Protocol | **Yes** | `ProtocolTCP`, `ProtocolUDP`, `ProtocolICMP`, `ProtocolAny` | Network protocol |
| `PortMin` | *int | No | 0-65535 | Starting port (required for TCP/UDP) |
| `PortMax` | *int | No | 0-65535 | Ending port (required for TCP/UDP) |
| `RemoteCIDR` | string | **Yes** | Valid CIDR | Source/destination IP range |

**Validation Rules**:
- `Direction`: Must be DirectionIngress or DirectionEgress (enforced by type system)
- `Protocol`: Must be ProtocolTCP, ProtocolUDP, ProtocolICMP, or ProtocolAny (enforced by type system)
- `PortMin`/`PortMax`: Required for TCP/UDP, ignored for ICMP
- `PortMax` must be >= `PortMin` (if both set)
- `RemoteCIDR`: Must be valid CIDR notation (e.g., "0.0.0.0/0", "10.0.0.0/8")

**Example JSON Requests**:

**TCP Rule (Port Range)**:
```json
{
  "direction": "ingress",
  "protocol": "tcp",
  "port_min": 8000,
  "port_max": 9000,
  "remote_cidr": "192.168.0.0/16"
}
```
*Go code: Use `Direction: DirectionIngress, Protocol: ProtocolTCP`*

**UDP Rule (Single Port)**:
```json
{
  "direction": "egress",
  "protocol": "udp",
  "port_min": 53,
  "port_max": 53,
  "remote_cidr": "8.8.8.8/32"
}
```
*Go code: Use `Direction: DirectionEgress, Protocol: ProtocolUDP`*

**ICMP Rule (No Ports)**:
```json
{
  "direction": "ingress",
  "protocol": "icmp",
  "remote_cidr": "0.0.0.0/0"
}
```
*Go code: Use `Direction: DirectionIngress, Protocol: ProtocolICMP`*

---

## List and Filter Models

### ListSecurityGroupsOptions

**Location**: `models/vps/securitygroups/securitygroup.go`

**API Mapping**: Query parameters from `GET /api/v1/project/{project-id}/security_groups`

```go
package securitygroups

// ListSecurityGroupsOptions defines optional filters for listing security groups.
type ListSecurityGroupsOptions struct {
    Name   *string // Filter by name (exact match)
    UserID *string // Filter by user_id (admin only)
    Detail *bool   // Include rules array in response (default: false)
}
```

**Field Descriptions**:

| Field | Type | API Param | Description |
|-------|------|-----------|-------------|
| `Name` | *string | `name` | Filter by exact security group name |
| `UserID` | *string | `user_id` | Filter by owner user ID (admin users only) |
| `Detail` | *bool | `detail` | Include rules in response (default: false) |

**Query Parameter Encoding**:
```go
// If opts.Detail = true
// URL: /api/v1/project/proj-123/security_groups?detail=true

// If opts.Name = "web-sg" and opts.Detail = true
// URL: /api/v1/project/proj-123/security_groups?name=web-sg&detail=true
```

**Example Usage**:
```go
// List all security groups (no rules)
opts := nil
resp, err := client.SecurityGroups().List(ctx, opts)

// List security groups with rules
detail := true
opts := &ListSecurityGroupsOptions{Detail: &detail}
resp, err := client.SecurityGroups().List(ctx, opts)

// Filter by name and include rules
name := "web-servers"
detail := true
opts := &ListSecurityGroupsOptions{
    Name:   &name,
    Detail: &detail,
}
resp, err := client.SecurityGroups().List(ctx, opts)
```

---

### SecurityGroupListResponse

**Location**: `models/vps/securitygroups/securitygroup.go`

**API Mapping**: `pb.SgListOutput` from `swagger/vps.yaml`

```go
package securitygroups

// SecurityGroupListResponse represents a list of security groups.
type SecurityGroupListResponse struct {
    SecurityGroups []*SecurityGroup `json:"security_groups"`
    // Total field removed - not present in API specification
}
```

**Field Descriptions**:

| Field | Type | Description |
|-------|------|-------------|
| `SecurityGroups` | []*SecurityGroup | Array of security group objects |

**⚠️ Breaking Change**:
- **Removed**: `Total int` field (was never present in API spec, always 0)
- **Migration**: Use `len(response.SecurityGroups)` to count results

**Example JSON Response** (with `detail=false`):
```json
{
  "security_groups": [
    {
      "id": "sg-001",
      "name": "web-servers",
      "description": "Web server security group",
      "project_id": "proj-123",
      "user_id": "user-456",
      "namespace": "default",
      "rules": [],
      "createdAt": "2025-11-09T10:00:00Z",
      "updatedAt": "2025-11-09T10:00:00Z"
    },
    {
      "id": "sg-002",
      "name": "database-servers",
      "description": "Database security group",
      "project_id": "proj-123",
      "user_id": "user-456",
      "namespace": "default",
      "rules": [],
      "createdAt": "2025-11-09T09:00:00Z",
      "updatedAt": "2025-11-09T09:00:00Z"
    }
  ]
}
```

**Example JSON Response** (with `detail=true`):
```json
{
  "security_groups": [
    {
      "id": "sg-001",
      "name": "web-servers",
      "description": "Web server security group",
      "project_id": "proj-123",
      "user_id": "user-456",
      "namespace": "default",
      "rules": [
        {
          "id": "rule-001",
          "direction": "ingress",
          "protocol": "tcp",
          "port_min": 80,
          "port_max": 80,
          "remote_cidr": "0.0.0.0/0"
        },
        {
          "id": "rule-002",
          "direction": "ingress",
          "protocol": "tcp",
          "port_min": 443,
          "port_max": 443,
          "remote_cidr": "0.0.0.0/0"
        }
      ],
      "createdAt": "2025-11-09T10:00:00Z",
      "updatedAt": "2025-11-09T10:00:00Z"
    }
  ]
}
```

---

## Model Relationships

### Entity-Relationship Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│ SecurityGroup                                                   │
│ ─────────────────────────────────────────────────────────────── │
│ + ID: string                                                    │
│ + Name: string                                                  │
│ + Description: string                                           │
│ + ProjectID: string ─────────────────────┐                      │
│ + UserID: string ────────────────────┐   │                      │
│ + Namespace: string                  │   │                      │
│ + Rules: []SecurityGroupRule ────────┼───┼──────┐               │
│ + CreatedAt: time.Time               │   │      │               │
│ + UpdatedAt: time.Time               │   │      │               │
│ + Project: *IDName ───────────────────┼───┘      │               │
│ + User: *IDName ──────────────────────┘          │               │
└──────────────────────────────────────────────────┼───────────────┘
                                                   │
                                                   │ 1:N
                                                   │
┌──────────────────────────────────────────────────┼───────────────┐
│ SecurityGroupRule                                │               │
│ ────────────────────────────────────────────────────────────── │
│ + ID: string                                     │               │
│ + Direction: string ("ingress" | "egress")                      │
│ + Protocol: string ("tcp" | "udp" | "icmp")                     │
│ + PortMin: int (0-65535)                                        │
│ + PortMax: int (0-65535)                                        │
│ + RemoteCIDR: string (CIDR notation)                            │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│ IDName (Nested Reference)                                       │
│ ─────────────────────────────────────────────────────────────── │
│ + ID: string                                                    │
│ + Name: string                                                  │
└─────────────────────────────────────────────────────────────────┘
```

### Relationship Rules

1. **SecurityGroup → SecurityGroupRule**: One-to-Many
   - A security group contains 0 to N rules
   - Rules are embedded in the `Rules` array
   - Rules are returned only when `detail=true` in List operation

2. **SecurityGroup → Project (IDName)**: Many-to-One
   - Each security group belongs to one project
   - `ProjectID` stores the foreign key
   - `Project` (IDName) provides denormalized name for convenience

3. **SecurityGroup → User (IDName)**: Many-to-One
   - Each security group has one owner user
   - `UserID` stores the foreign key
   - `User` (IDName) provides denormalized name for convenience

---

## Breaking Changes from Previous Version

### 1. Removed `Total` Field from `SecurityGroupListResponse`

**Before (Incorrect)**:
```go
type SecurityGroupListResponse struct {
    SecurityGroups []*SecurityGroup `json:"security_groups"`
    Total          int              `json:"total"`  // ❌ Not in API spec
}
```

**After (Correct)**:
```go
type SecurityGroupListResponse struct {
    SecurityGroups []*SecurityGroup `json:"security_groups"`
    // Total field removed - not present in API specification
}
```

**Migration**:
```go
// Before
resp, err := client.SecurityGroups().List(ctx, opts)
count := resp.Total  // Always 0, never populated

// After
resp, err := client.SecurityGroups().List(ctx, opts)
count := len(resp.SecurityGroups)  // Actual count
```

### 2. Changed `PortMin`/`PortMax` Types in `SecurityGroupRule`

**Before (Incorrect)**:
```go
type SecurityGroupRule struct {
    // ...
    PortMin *int `json:"port_min,omitempty"`  // ❌ Pointer in response model
    PortMax *int `json:"port_max,omitempty"`  // ❌ Pointer in response model
}
```

**After (Correct)**:
```go
type SecurityGroupRule struct {
    // ...
    PortMin int `json:"port_min"`  // ✅ Plain int (API always returns value)
    PortMax int `json:"port_max"`  // ✅ Plain int (API always returns value)
}
```

**Migration**:
```go
// Before
if rule.PortMin != nil {
    port := *rule.PortMin
}

// After
port := rule.PortMin  // Direct access, no nil check needed
```

**Note**: `SecurityGroupRuleCreateRequest` still uses `*int` (correctly, since ports are optional in requests).

---

## Validation Summary

| Model | Required Fields | Optional Fields | Constraints |
|-------|-----------------|-----------------|-------------|
| SecurityGroupCreateRequest | `Name` | `Description`, `Rules` | Name non-empty |
| SecurityGroupUpdateRequest | - | `Name`, `Description` | At least one field recommended |
| SecurityGroupRuleCreateRequest | `Direction`, `Protocol`, `RemoteCIDR` | `PortMin`, `PortMax` | Ports required for TCP/UDP; PortMax >= PortMin; Valid CIDR |
| ListSecurityGroupsOptions | - | `Name`, `UserID`, `Detail` | All optional filters |

---

## Conclusion

All models are now aligned with the VPS API specification (`swagger/vps.yaml`):

- ✅ `SecurityGroup` matches `pb.SgInfo` exactly
- ✅ `SecurityGroupRule` matches `pb.SgRuleInfo` (corrected port types)
- ✅ `SecurityGroupListResponse` matches `pb.SgListOutput` (removed Total field)
- ✅ `SecurityGroupCreateRequest` matches `SgCreateInput`
- ✅ `SecurityGroupUpdateRequest` matches `SgUpdateInput`
- ✅ `SecurityGroupRuleCreateRequest` matches `SgRuleCreateInput`
- ✅ Breaking changes documented with migration guides

**Next Step**: Create contract interfaces in `contracts/` directory and quick-start example in `quickstart.md`.
