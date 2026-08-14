# Data Model: IAM API Client

**Feature**: 009-iam-api  
**Date**: 2025-11-12  
**Phase**: 1 - Data Structures & Types

## Overview

This document defines all data structures for the IAM API client, extracted from the Swagger specification and functional requirements. All types support JSON marshaling/unmarshaling matching the IAM API contract.

## Core Entities

### User

Represents an authenticated user in the IAM system.

**Location**: `models/iam/users/models.go` (following VPS/VRM pattern)

```go
type User struct {
    UserID      string                 `json:"userId"`
    Account     string                 `json:"account"`
    DisplayName string                 `json:"displayName"`
    Description string                 `json:"description"`
    Extra       map[string]interface{} `json:"extra"`
    Namespace   string                 `json:"namespace"`
    Email       string                 `json:"email"`
    Frozen      bool                   `json:"frozen"`
    MFA         bool                   `json:"mfa"`
    CreatedAt   string                 `json:"createdAt"`
    UpdatedAt   string                 `json:"updatedAt"`
    LastLoginAt string                 `json:"lastLoginAt"`
}
```

**Field Descriptions**:
- `UserID`: Unique user identifier (UUID format)
- `Account`: User account identifier (typically email)
- `DisplayName`: Human-readable name for display
- `Description`: User description
- `Extra`: Unstructured metadata (forward compatibility)
- `Namespace`: Tenant/namespace identifier (e.g., "ci.asus.com")
- `Email`: User email address
- `Frozen`: Whether user account is frozen/disabled
- `MFA`: Whether MFA (Multi-Factor Authentication) is enabled
- `CreatedAt`: ISO8601 timestamp of creation
- `UpdatedAt`: ISO8601 timestamp of last update
- `LastLoginAt`: ISO8601 timestamp of last login

**Example JSON**:
```json
{
  "userId": "4990ccdb-a9b1-49e5-91df-67c921601d81",
  "account": "system@ci.asus.com",
  "displayName": "system",
  "description": "This is system account",
  "extra": {},
  "namespace": "ci.asus.com",
  "email": "system@ci.asus.com",
  "frozen": false,
  "mfa": false,
  "createdAt": "2025-11-11T15:18:36Z",
  "updatedAt": "2025-11-11T15:32:32Z",
  "lastLoginAt": "2025-11-11T15:32:32Z"
}
```

**Validation Rules**:
- `UserID` must be valid UUID format (server-side validation)
- `Account` and `Email` typically match (server-side validation)
- Timestamps must be RFC3339 format (e.g., "2025-11-11T15:18:36Z")
- `Extra` may contain deeply nested objects; client ignores unknown fields
- `Frozen` determines account access eligibility
- `MFA` indicates multi-factor authentication status

**State Transitions**: N/A (read-only entity from client perspective)

---

### Project

Represents a project/tenant in the multi-tenant system.

**Location**: `models/iam/projects/models.go` (following VPS/VRM pattern)

```go
type Project struct {
    ProjectID   string                 `json:"projectId"`
    DisplayName string                 `json:"displayName"`
    Description string                 `json:"description"`
    Extra       map[string]interface{} `json:"extra"`
    Namespace   string                 `json:"namespace"`
    Frozen      bool                   `json:"frozen"`
    CreatedAt   string                 `json:"createdAt"`
    UpdatedAt   string                 `json:"updatedAt"`
}
```

**Field Descriptions**:
- `ProjectID`: Unique project identifier (UUID format)
- `DisplayName`: Human-readable project name
- `Description`: Project description
- `Extra`: Unstructured project metadata (e.g., iService projectSysCode)
- `Namespace`: Tenant namespace (e.g., "ci.asus.com")
- `Frozen`: Whether project is frozen/disabled
- `CreatedAt`: ISO8601 timestamp of creation
- `UpdatedAt`: ISO8601 timestamp of last update

**Example JSON**:
```json
{
  "projectId": "91457b61-0b92-4aa8-b136-b03d88f04946",
  "displayName": "prj1762875055667",
  "description": "",
  "extra": {
    "iservice": {
      "projectSysCode": "TCI111222"
    }
  },
  "namespace": "ci.asus.com",
  "frozen": false,
  "createdAt": "2025-11-11T15:31:02Z",
  "updatedAt": "2025-11-11T15:31:06Z"
}
```

**Validation Rules**:
- `ProjectID` must be valid UUID format
- `DisplayName` typically required (server-side validation)
- Timestamps must be RFC3339 format (e.g., "2025-11-11T15:31:02Z")
- `Extra` may contain nested objects like iService metadata
- `Frozen` determines project access eligibility

---

### ProjectMembership

Wraps a Project with user-specific membership context (permissions, access status, role).

**Location**: `models/iam/projects/models.go` (following VPS/VRM pattern)

```go
type ProjectMembership struct {
    Project            *Project               `json:"project"`
    GlobalPermissionID string                 `json:"globalPermissionId"`
    GlobalPermission   *Permission            `json:"globalPermission"`
    UserPermissionID   string                 `json:"userPermissionId"`
    UserPermission     *Permission            `json:"userPermission"`
    Frozen             bool                   `json:"frozen"`
    TenantRole         TenantRole             `json:"tenantRole"`
    Extra              map[string]interface{} `json:"extra"`
}
```

**Field Descriptions**:
- `Project`: Embedded project details
- `GlobalPermissionID`: ID of project-wide permission template
- `GlobalPermission`: Project-level permission details (default permission)
- `UserPermissionID`: ID of user-specific permission override
- `UserPermission`: User-specific permission details (effective permission)
- `Frozen`: Whether user's access to this project is frozen
- `TenantRole`: User's role enum (TenantRoleMember, TenantRoleAdmin, TenantRoleOwner)
- `Extra`: Additional membership metadata

**Example JSON**:
```json
{
  "project": {
    "projectId": "91457b61-0b92-4aa8-b136-b03d88f04946",
    "displayName": "prj1762875055667",
    "description": "",
    "extra": {
      "iservice": {
        "projectSysCode": "TCI111222"
      }
    },
    "namespace": "ci.asus.com",
    "frozen": false,
    "createdAt": "2025-11-11T15:31:02Z",
    "updatedAt": "2025-11-11T15:31:06Z"
  },
  "globalPermissionId": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
  "globalPermission": {
    "id": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
    "label": "DEFAULT"
  },
  "userPermissionId": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
  "userPermission": {
    "id": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
    "label": "DEFAULT"
  },
  "frozen": false,
  "tenantRole": "TENANT_MEMBER",
  "extra": {}
}
```

**Relationships**:
- `Project` (1:1): Each membership relates to exactly one project
- `GlobalPermission` (1:1): Project's default permission template
- `UserPermission` (1:1): User's effective permission in project (may differ from global)

**Validation Rules**:
- `Project` must not be nil
- Permission IDs typically non-empty (server-side)
- `Frozen` determines access eligibility
- `TenantRole` must be one of the valid enum values (validated via IsValid() method)
- Role hierarchy: TenantRoleMember < TenantRoleAdmin < TenantRoleOwner
- Both project-level and membership-level frozen flags must be checked for access

**Usage Examples**:
```go
// Check role
if membership.TenantRole == common.TenantRoleOwner {
    // User is project owner
}

// Check if role is valid
if membership.TenantRole.IsValid() {
    // Handle valid role
}

// Convert to string
roleStr := membership.TenantRole.String()
```

---

### Permission

Represents access control permissions.

**Location**: `models/iam/common/common.go`

```go
type Permission struct {
    ID    string `json:"id"`
    Label string `json:"label"`
}
```

**Field Descriptions**:
- `ID`: Unique permission identifier (UUID format)
- `Label`: Human-readable permission name (e.g., "DEFAULT", "ADMIN", "READONLY")

**Validation Rules**:
- Both fields typically non-empty
- Permission labels define access levels within projects

**Common Permission Labels**:
- `DEFAULT`: Standard project member access
- `ADMIN`: Administrative privileges
- `READONLY`: Read-only access

---

### TenantRole

Enum type representing user's role within a project/tenant.

**Location**: `models/iam/common/common.go`

```go
type TenantRole string

const (
    TenantRoleMember TenantRole = "TENANT_MEMBER"
    TenantRoleAdmin  TenantRole = "TENANT_ADMIN"
    TenantRoleOwner  TenantRole = "TENANT_OWNER"
)

// String returns the string representation of the TenantRole
func (r TenantRole) String() string {
    return string(r)
}

// IsValid checks if the TenantRole is a valid value
func (r TenantRole) IsValid() bool {
    switch r {
    case TenantRoleMember, TenantRoleAdmin, TenantRoleOwner:
        return true
    }
    return false
}
```

**Role Hierarchy** (ascending authority):
1. `TenantRoleMember` - Standard project member
2. `TenantRoleAdmin` - Project administrator
3. `TenantRoleOwner` - Project owner (highest authority)

**Validation Rules**:
- Only three valid values
- Case-sensitive string matching
- Unknown values should be handled gracefully (forward compatibility)

## Operation Models

### GetUser Operation

**Request**: None (no parameters)

**Response**: `models/iam/users/user.go`

```go
type GetUserResponse struct {
    *common.User
}
```

**Usage**:
```go
user, err := client.IAM().User().Get(ctx)
// user.Account, user.UserID, etc.
```

---

### ListProjects Operation

**Request Options**: `models/iam/projects/options.go` (consistent with VPS/VRM pattern)

```go
// Single options struct pointer pattern (like VPS/VRM)
type ListProjectsOptions struct {
    Offset *int
    Limit  *int
    Order  *string
}

// Usage:
// opts := &ListProjectsOptions{Offset: intPtr(10), Limit: intPtr(20)}
// projects, err := client.Projects().List(ctx, opts)
//
// Or pass nil for defaults:
// projects, err := client.Projects().List(ctx, nil)
```

**Response**: `models/iam/projects/project.go`

```go
type ListProjectsResponse struct {
    Projects []*common.ProjectMembership `json:"projects"`
    Total    int                         `json:"total"`
}
```

**Custom Unmarshaler**:
```go
func (r *ListProjectsResponse) UnmarshalJSON(data []byte) error {
    var wrapper struct {
        Projects []*common.ProjectMembership `json:"projects"`
        Total    int                         `json:"total"`
    }
    if err := json.Unmarshal(data, &wrapper); err != nil {
        return err
    }
    r.Projects = wrapper.Projects
    r.Total = wrapper.Total
    return nil
}
```

**Usage**:
```go
resp, err := client.IAM().Projects().List(ctx, 
    iam.WithOffset(20),
    iam.WithLimit(50),
)
// Direct array access: resp.Projects[0]
// Total count: resp.Total
```

**Validation Rules**:
- `offset` must be >= 0 (server-side validation)
- `limit` must be 1-100 (server-side validation)
- `order` format determined by server

---

### GetProject Operation

**Request**: Project ID as method parameter

**Response**: `models/iam/projects/project.go`

```go
type GetProjectResponse struct {
    ProjectID        string                 `json:"projectId"`
    DisplayName      string                 `json:"displayName"`
    Description      string                 `json:"description"`
    Extra            map[string]interface{} `json:"extra"`
    Namespace        string                 `json:"namespace"`
    Frozen           bool                   `json:"frozen"`
    GlobalPermission *common.Permission     `json:"globalPermission"`
    UserPermission   *common.Permission     `json:"userPermission"`
    CreatedAt        string                 `json:"createdAt"`
    UpdatedAt        string                 `json:"updatedAt"`
}
```

**Note**: GetProject response has a **flat structure** (project fields at top level), which differs from ListProjects where project details are nested inside a `project` object within ProjectMembership. This is the actual API design.

**Usage**:
```go
resp, err := client.IAM().Projects().Get(ctx, "project-uuid")
// resp.ProjectID, resp.DisplayName - direct access
// resp.UserPermission.Label - permission details
```

**Validation Rules**:
- `projectID` parameter must be non-empty valid UUID
- Response does NOT include membership-level fields (frozen, tenantRole, extra from membership)
- Only includes project details and permissions

## Error Models

Uses existing `internal/types.SDKError`:

```go
// No new error types - reuse existing SDK error infrastructure
type SDKError struct {
    StatusCode int
    ErrorCode  int
    Message    string
    Meta       map[string]interface{}
    Cause      error
}
```

**HTTP Status Mapping**:
- `200`: Success - parse response
- `400`: Bad Request - invalid project ID, pagination params
- `403`: Forbidden - invalid/expired token
- `500`: Internal Server Error - service failure

## Type Relationships

```
Client (root)
  └─> IAM Client
       ├─> User Resource
       │    └─> Get() -> User
       │         ├─> userId (UUID)
       │         ├─> account, email
       │         ├─> frozen, mfa (bools)
       │         └─> extra (metadata)
       │
       └─> Projects Resource
            ├─> List() -> ListProjectsResponse
            │    ├─> Projects: []*ProjectMembership (nested structure)
            │    │    ├─> Project (nested object)
            │    │    │    ├─> projectId (UUID)
            │    │    │    ├─> displayName, description
            │    │    │    ├─> frozen (bool)
            │    │    │    └─> extra (metadata)
            │    │    ├─> GlobalPermission
            │    │    ├─> UserPermission
            │    │    ├─> tenantRole (TenantRole enum)
            │    │    ├─> frozen (membership frozen)
            │    │    └─> extra (membership extra)
            │    └─> Total: int
            │
            └─> Get(id) -> GetProjectResponse (flat structure)
                 ├─> projectId (direct field, not nested)
                 ├─> displayName, description
                 ├─> frozen (project frozen)
                 ├─> GlobalPermission
                 ├─> UserPermission
                 └─> createdAt, updatedAt
                 
                 Note: Get() does NOT include:
                 - tenantRole (not in response)
                 - membership extra (not in response)
                 - membership frozen (only project frozen)

Shared Types:
  - Permission (referenced by ProjectMembership, GetProjectResponse)
  - Extra (map[string]interface{} in User, Project, ProjectMembership, GetProjectResponse)
```

## JSON Examples

### GetUser Response
```json
{
  "userId": "4990ccdb-a9b1-49e5-91df-67c921601d81",
  "account": "system@ci.asus.com",
  "displayName": "system",
  "description": "This is system account",
  "extra": {},
  "namespace": "ci.asus.com",
  "email": "system@ci.asus.com",
  "frozen": false,
  "mfa": false,
  "createdAt": "2025-11-11T15:18:36Z",
  "updatedAt": "2025-11-11T15:32:32Z",
  "lastLoginAt": "2025-11-11T15:32:32Z"
}
```

### ListProjects Response
```json
{
  "projects": [
    {
      "project": {
        "projectId": "91457b61-0b92-4aa8-b136-b03d88f04946",
        "displayName": "prj1762875055667",
        "description": "",
        "extra": {
          "iservice": {
            "projectSysCode": "TCI111222"
          }
        },
        "namespace": "ci.asus.com",
        "frozen": false,
        "createdAt": "2025-11-11T15:31:02Z",
        "updatedAt": "2025-11-11T15:31:06Z"
      },
      "globalPermissionId": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
      "globalPermission": {
        "id": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
        "label": "DEFAULT"
      },
      "userPermissionId": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
      "userPermission": {
        "id": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
        "label": "DEFAULT"
      },
      "frozen": false,
      "tenantRole": "TENANT_MEMBER",
      "extra": {}
    }
  ],
  "total": 1
}
```

### GetProject Response
```json
{
  "projectId": "91457b61-0b92-4aa8-b136-b03d88f04946",
  "displayName": "prj1762875055667",
  "description": "",
  "extra": {
    "iservice": {
      "projectSysCode": "TCI111222"
    }
  },
  "namespace": "ci.asus.com",
  "frozen": false,
  "globalPermission": {
    "id": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
    "label": "DEFAULT"
  },
  "userPermission": {
    "id": "e763477c-c6a2-4eff-a1ee-2d6a02b05a36",
    "label": "DEFAULT"
  },
  "createdAt": "2025-11-11T15:31:02Z",
  "updatedAt": "2025-11-11T15:31:06Z"
}
```

## File Organization

```
models/iam/
├── common/
│   ├── common.go          # Shared types: Permission, timestamps helpers
│   └── common_test.go     # Unit tests for common types
├── users/
│   ├── models.go          # User, GetUserResponse (following VPS/VRM pattern: models.go)
│   └── models_test.go     # Unit tests for user models (JSON parsing)
└── projects/
    ├── models.go          # Project, ProjectMembership, ListProjectsResponse, GetProjectResponse
    ├── models_test.go     # Unit tests for project models (JSON parsing, custom unmarshaler)
    └── options.go         # ListProjectsOptions struct (single options pointer pattern, like VPS/VRM)
```

**Core Types** (organized by directory following VPS/VRM pattern):

**In `models/iam/common/common.go`** (shared types):
- `Permission`: 2 fields (id, label)
- `TenantRole`: Custom string type enum with 3 constants (TenantRoleMember, TenantRoleAdmin, TenantRoleOwner)

**In `models/iam/users/models.go`** (user-specific types):
- `User`: 12 fields (userId, account, displayName, description, extra, namespace, email, frozen, mfa, createdAt, updatedAt, lastLoginAt)
- `GetUserResponse`: Contains User and additional operation-specific fields

**In `models/iam/projects/models.go`** (project-specific types):
- `Project`: 8 fields (projectId, displayName, description, extra, namespace, frozen, createdAt, updatedAt)
- `ProjectMembership`: 8 fields (project, globalPermissionId, globalPermission, userPermissionId, userPermission, frozen, tenantRole, extra)
- `ListProjectsResponse`: Contains []ProjectMembership and total count
- `GetProjectResponse`: Contains Project and permission details

**In `models/iam/projects/options.go`** (pagination options):
- `ListProjectsOptions`: Simple struct with Offset, Limit, Order fields (pointer types)
- Single options struct pattern (consistent with VPS/VRM: `flavors.ListFlavorsOptions`, `servers.ServersListRequest`)

## Testing Requirements

Each model MUST have unit tests verifying:
1. **JSON Unmarshaling**: Parse valid Swagger examples successfully
2. **Required Fields**: Detect missing required fields
3. **Optional Fields**: Handle missing optional fields gracefully
4. **Unknown Fields**: Ignore unexpected fields (forward compatibility)
5. **Type Conversions**: Proper handling of timestamps, booleans, nested objects
6. **Edge Cases**: Empty strings, null values, large Extra objects

Target: 80%+ test coverage for all model packages.
