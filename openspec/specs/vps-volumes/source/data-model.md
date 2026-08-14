# Data Model: Volumes API

**Date**: 2025-11-13  
**Feature**: Volumes API Client (`008-volumes-api`)  
**Source**: `swagger/vps.yaml` (volumes and volume_types tags)

## Overview

This document defines the data structures for the Volumes API client implementation. All models map directly to Swagger definitions in `swagger/vps.yaml`.

### Project Scoping

**Pattern**: The VPS client follows a project-scoped pattern where `projectID` is provided at client initialization and embedded in the URL path.

- **Client Initialization**: `vpsClient := cloudsdk.Client.Project(projectID).VPS()`
- **Sub-client Creation**: `volumesClient := vpsClient.Volumes()` (inherits projectID from parent)
- **URL Path Format**: `/api/v1/project/{project-id}/volumes`
- **Method Signatures**: Methods do NOT require explicit `projectID` parameters (e.g., `Create(ctx, *CreateVolumeRequest)`)

This satisfies **FR-009** ("All API methods MUST require project ID as a parameter to scope operations to the correct project") by scoping at the client level rather than per-method. This pattern is consistent with existing VPS modules (flavors, keypairs, networks, etc.).

---

## Core Entities

### Volume

Represents a block storage volume resource.

**Source**: `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.VolumeInfo`  
**Package**: `models/vps/volumes`

```go
package volumes

import (
    "time"
    "github.com/Zillaforge/cloud-sdk/models/vps/common"
)

// VolumeStatus represents the current state of a volume.
type VolumeStatus string

// VolumeStatus constants for volume lifecycle states.
const (
    VolumeStatusCreating  VolumeStatus = "creating"  // Volume is being created
    VolumeStatusAvailable VolumeStatus = "available" // Volume created and ready to attach
    VolumeStatusInUse     VolumeStatus = "in-use"    // Volume attached to server
    VolumeStatusDetaching VolumeStatus = "detaching" // Volume being detached from server
    VolumeStatusExtending VolumeStatus = "extending" // Volume size is being increased
    VolumeStatusDeleting  VolumeStatus = "deleting"  // Volume is being deleted
    VolumeStatusDeleted   VolumeStatus = "deleted"   // Volume deleted
    VolumeStatusError     VolumeStatus = "error"     // Operation failed, check StatusReason
)

// Volume represents a block storage volume.
// Matches pegasus-cloud_com_aes_virtualplatformserviceclient_pb.VolumeInfo from vps.yaml.
type Volume struct {
    ID           string           `json:"id"`
    Name         string           `json:"name"`
    Description  string           `json:"description,omitempty"`
    Size         int              `json:"size"`                   // Size in GB
    Type         string           `json:"type"`                   // Storage type (SSD, HDD, etc.)
    Status       VolumeStatus     `json:"status"`                 // Volume status
    StatusReason string           `json:"status_reason,omitempty"` // Status reason text
    Attachments  []common.IDName  `json:"attachments,omitempty"` // Servers this volume is attached to
    Project      common.IDName    `json:"project"`               // Project reference
    ProjectID    string           `json:"project_id"`
    User         common.IDName    `json:"user"`                  // User reference
    UserID       string           `json:"user_id"`
    Namespace    string           `json:"namespace"`
    CreatedAt    *time.Time       `json:"createdAt,omitempty"`   // Creation timestamp (ISO 8601)
    UpdatedAt    *time.Time       `json:"updatedAt,omitempty"`   // Last update timestamp (ISO 8601)
}
```

**Field Descriptions**:
- `ID`: Unique volume identifier (assigned by server)
- `Name`: User-defined volume name
- `Description`: Optional user-defined description
- `Size`: Volume size in gigabytes
- `Type`: Storage class (from VolumeTypes list: SSD, HDD, NVMe, etc.)
- `Status`: Current volume state (typed enum: VolumeStatusAvailable, VolumeStatusInUse, etc.)
- `StatusReason`: Additional context for status (especially for error states)
- `Attachments`: List of servers this volume is attached to (empty if unattached)
- `Project`: Project this volume belongs to (ID and name)
- `ProjectID`: Project ID (duplicate for convenience)
- `User`: User who created the volume (ID and name)
- `UserID`: User ID (duplicate for convenience)
- `Namespace`: Organizational namespace
- `CreatedAt`: Volume creation timestamp
- `UpdatedAt`: Last modification timestamp

**Validation Rules**:
- ID: Required, non-empty
- Name: Required, non-empty
- Size: Non-negative integer
- Type: Required, non-empty, must match available volume type
- Status: Must be one of defined VolumeStatus constants

---

## Request/Response Models

### CreateVolumeRequest

Request for creating a new volume.

**Source**: `VolumeCreateInput` from swagger  
**Package**: `models/vps/volumes`

```go
// CreateVolumeRequest represents parameters for creating a volume.
// Matches VolumeCreateInput from vps.yaml.
type CreateVolumeRequest struct {
    Name        string `json:"name"`                   // Required: Volume name
    Type        string `json:"type"`                   // Required: Volume type (from VolumeTypes)
    Size        int    `json:"size,omitempty"`        // Optional: Size in GB (server validates)
    Description string `json:"description,omitempty"` // Optional: Volume description
    SnapshotID  string `json:"snapshot_id,omitempty"` // Optional: Create from snapshot
}
```

**Field Descriptions**:
- `Name`: User-defined name for the volume (required)
- `Type`: Storage class identifier, must be from VolumeTypes list (required)
- `Size`: Volume size in GB, validated by server (optional - may have default)
- `Description`: User-defined description (optional)
- `SnapshotID`: If provided, creates volume from snapshot (optional)

**Validation Rules** (client-side):
- Name: Required, non-empty
- Type: Required, non-empty
- Size: If provided, must be positive integer
- Note: Server performs additional validation (quotas, limits, snapshot compatibility)

---

### UpdateVolumeRequest

Request for updating volume metadata.

**Source**: `VolumeUpdateInput` from swagger  
**Package**: `models/vps/volumes`

```go
// UpdateVolumeRequest represents parameters for updating a volume.
// Matches VolumeUpdateInput from vps.yaml.
type UpdateVolumeRequest struct {
    Name        string `json:"name,omitempty"`        // Optional: New volume name
    Description string `json:"description,omitempty"` // Optional: New volume description
}
```

**Field Descriptions**:
- `Name`: New name for the volume (optional)
- `Description`: New description for the volume (optional)

**Validation Rules**:
- At least one field should be provided (both optional but at least one recommended)
- Name: If provided, non-empty

**Note**: Size and Type cannot be changed via Update (use VolumeActionRequest with "extend" for size changes)

---

### VolumeActionRequest

Request for performing operations on a volume.

**Source**: `VolActionInput` from swagger  
**Package**: `models/vps/volumes`

```go
// VolumeActionRequest represents parameters for volume actions.
// Matches VolActionInput from vps.yaml.
type VolumeActionRequest struct {
    Action   VolumeAction `json:"action"`              // Required: Action type
    ServerID string       `json:"server_id,omitempty"` // Required for attach/detach
    NewSize  int          `json:"new_size,omitempty"`  // Required for extend
}

// VolumeAction represents supported volume action types.
type VolumeAction string

// VolumeAction constants for supported volume actions.
const (
    VolumeActionAttach VolumeAction = "attach" // Attach volume to server
    VolumeActionDetach VolumeAction = "detach" // Detach volume from server
    VolumeActionExtend VolumeAction = "extend" // Extend volume size
    VolumeActionRevert VolumeAction = "revert" // Revert volume to snapshot
)
```

**Field Descriptions**:
- `Action`: Type of operation to perform (required, must be one of: attach, detach, extend, revert)
- `ServerID`: Target server ID for attach/detach operations (required for attach/detach)
- `NewSize`: New size in GB for extend operation (required for extend, must be > current size)

**Validation Rules** (client-side):
- Action: Required, must be one of VolumeActionAttach, VolumeActionDetach, VolumeActionExtend, VolumeActionRevert
- ServerID: Required when Action is "attach" or "detach"
- NewSize: Required when Action is "extend", must be positive integer

**Action-Specific Requirements**:
- **attach**: Requires ServerID, volume must be "available" status
- **detach**: Requires ServerID, volume must be "in-use" status
- **extend**: Requires NewSize > current Size, volume can be attached or detached
- **revert**: Volume must have been created from a snapshot

---

## Query Options

### ListVolumesOptions

Options for filtering volume list queries.

**Source**: API query parameters from swagger  
**Package**: `models/vps/volumes`

```go
// ListVolumesOptions provides filtering options for listing volumes.
type ListVolumesOptions struct {
    Name   string // Filter by name (partial match)
    UserID string // Filter by user ID
    Status string // Filter by status (available, in-use, etc.)
    Type   string // Filter by volume type
    Detail bool   // Include attachment details
}
```

**Field Descriptions**:
- `Name`: Filter volumes by name (may support partial matching depending on API)
- `UserID`: Filter volumes by creator user ID
- `Status`: Filter volumes by current status
- `Type`: Filter volumes by storage type
- `Detail`: If true, API returns full attachment information (server details)

**Validation Rules**:
- All fields optional
- Multiple filters can be combined (AND logic)

---

## Response Models

### VolumeResponse

Response wrapper for single volume operations (Create, Get, Update).

**Package**: `models/vps/volumes`

```go
// VolumeResponse represents the response containing a single volume.
type VolumeResponse struct {
    Volume *Volume `json:"volume,omitempty"`
}
```

**Purpose**: Some API endpoints may wrap single volume in a response object. If not needed, can return Volume directly.

---

### VolumeListResponse

Response wrapper for list volumes API.

**Source**: `pb.VolumeListOutput` from swagger  
**Package**: `models/vps/volumes`

```go
// VolumeListResponse represents the response from listing volumes.
// Matches pb.VolumeListOutput from vps.yaml.
type VolumeListResponse struct {
    Volumes []*Volume `json:"volumes"`
}
```

**Purpose**: Intermediate struct for JSON unmarshaling. API returns `{"volumes": [...]}`, this struct unwraps to `[]*Volume`.

---

### VolumeTypeListResponse

Response wrapper for list volume types API.

**Source**: `pb.VolumeTypeListOutput` from swagger  
**Package**: `models/vps/volumetypes`

```go
package volumetypes

// VolumeTypeListResponse represents the response from listing volume types.
// Matches pb.VolumeTypeListOutput from vps.yaml.
type VolumeTypeListResponse struct {
    VolumeTypes []string `json:"volume_types"`
}
```

**Purpose**: Intermediate struct for JSON unmarshaling. API returns `{"volume_types": ["SSD", "HDD", ...]}`, this struct unwraps to `[]string`.

**Note**: VolumeTypes are simple strings representing storage classes. No separate VolumeType entity needed.

---

## Shared Types

### IDName (from common package)

Reference to another resource with ID and name.

**Source**: `pb.IDName` from swagger  
**Package**: `models/vps/common` (already exists)

```go
package common

// IDName represents a reference to another resource with ID and name.
type IDName struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

**Usage in Volume**:
- `Attachments`: List of servers (ID and name)
- `Project`: Project reference
- `User`: User reference

---

## State Transitions

### Volume Status Lifecycle

```
[creating] → [available] ⇄ [in-use] → [detaching] → [available]
                ↓                                        ↓
            [deleting] ← ← ← ← ← ← ← ← ← ← ← ← ← ← [deleting]
                ↓
            [deleted]
```

**Status Values** (VolumeStatus type):
- `VolumeStatusCreating`: Volume is being created
- `VolumeStatusAvailable`: Volume created and ready to attach
- `VolumeStatusInUse`: Volume attached to server
- `VolumeStatusDetaching`: Volume being detached from server
- `VolumeStatusExtending`: Volume size is being increased
- `VolumeStatusDeleting`: Volume is being deleted
- `VolumeStatusDeleted`: Volume deleted (may not be returned by API)
- `VolumeStatusError`: Operation failed, check StatusReason

**Valid Action Transitions**:
- `VolumeStatusAvailable` → `VolumeStatusInUse`: attach action
- `VolumeStatusInUse` → `VolumeStatusDetaching` → `VolumeStatusAvailable`: detach action
- `VolumeStatusAvailable` or `VolumeStatusInUse` → `VolumeStatusExtending` → same: extend action
- `VolumeStatusAvailable` → `VolumeStatusDeleting` → `VolumeStatusDeleted`: delete operation
- Cannot delete volume in `VolumeStatusInUse` status

---

## Entity Relationships

```
Volume
├── Project (common.IDName)    # Belongs to one project
├── User (common.IDName)       # Created by one user
└── Attachments ([]common.IDName) # Attached to zero or more servers

VolumeTypes ([]string)         # Available storage classes (independent)
```

**Constraints**:
- Volume belongs to exactly one project
- Volume created by exactly one user
- Volume can be attached to multiple servers (depending on volume type/configuration)
- Volume must be detached before deletion
- Volume type cannot be changed after creation

---

## Validation Summary

| Entity | Required Fields | Constraints |
|--------|----------------|-------------|
| Volume | ID, Name, Type, Size, Status | Size ≥ 0, Type from VolumeTypes, Status is VolumeStatus type |
| CreateVolumeRequest | Name, Type | Size > 0 if provided |
| UpdateVolumeRequest | - (both optional) | At least one field recommended |
| VolumeActionRequest | Action | ServerID for attach/detach, NewSize for extend |
| ListVolumesOptions | - (all optional) | Combine filters with AND logic |

---

## Swagger Reference Mapping

| Model | Swagger Definition | Location |
|-------|-------------------|----------|
| Volume | `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.VolumeInfo` | lines 1717-1752 |
| CreateVolumeRequest | `VolumeCreateInput` | lines 610-627 |
| UpdateVolumeRequest | `VolumeUpdateInput` | lines 628-634 |
| VolumeActionRequest | `VolActionInput` | lines 593-609 |
| VolumeListResponse | `pb.VolumeListOutput` | lines 1672-1678 |
| VolumeTypeListResponse | `pb.VolumeTypeListOutput` | lines 1679-1686 |
| IDName | `pb.IDName` | existing in common package |

---

## Implementation Files

### `models/vps/volumes/volume.go`
- Volume struct
- VolumeStatus type and constants
- CreateVolumeRequest struct
- UpdateVolumeRequest struct
- VolumeActionRequest struct
- VolumeAction type and constants
- ListVolumesOptions struct
- VolumeResponse struct (if needed)
- VolumeListResponse struct
- Validate() methods

### `models/vps/volumes/volume_test.go`
- Unit tests for all structs
- Validation tests
- JSON marshaling/unmarshaling tests
- Edge case tests
- VolumeStatus constant tests
- Target: 85%+ coverage

### `models/vps/volumetypes/volumetype.go`
- VolumeTypeListResponse struct

### `models/vps/volumetypes/volumetype_test.go`
- Unit tests for VolumeTypeListResponse
- JSON marshaling/unmarshaling tests
- Target: 85%+ coverage

### `modules/vps/volumes/client.go`
- Client struct with fields: baseClient, projectID, basePath
- NewClient(baseClient *internalhttp.Client, projectID string) *Client constructor
- Client implementation with methods:
  - List(ctx, *ListVolumesOptions) ([]*Volume, error)
  - Create(ctx, *CreateVolumeRequest) (*Volume, error)
  - Get(ctx, volumeID) (*Volume, error)
  - Update(ctx, volumeID, *UpdateVolumeRequest) (*Volume, error)
  - Delete(ctx, volumeID) error
  - Action(ctx, volumeID, *VolumeActionRequest) error

**Note**: projectID is stored in the Client struct and used to build URL paths (`/api/v1/project/{project-id}/volumes`), NOT passed as a method parameter. This follows the established VPS module pattern (see flavors, keypairs, networks).

### `modules/vps/volumes/client_test.go`
- Unit tests for all client methods
- Error handling tests
- Mock HTTP client tests
- Target: 85%+ coverage

### `modules/vps/volumetypes/client.go`
- Client struct with fields: baseClient, projectID, basePath
- NewClient(baseClient *internalhttp.Client, projectID string) *Client constructor
- Client implementation with methods:
  - List(ctx) ([]string, error)

**Note**: projectID is inherited from parent VPS client and embedded in URL paths.

### `modules/vps/volumetypes/client_test.go`
- Unit tests for List method
- Target: 85%+ coverage

---

## Notes

1. **Project Scoping Pattern**: projectID is provided at VPS client initialization, stored in sub-clients (volumes, volumetypes), and embedded in URL paths (`/api/v1/project/{project-id}/...`). Method signatures do NOT include explicit projectID parameters. This satisfies FR-009 via client-level scoping.
2. **Request/Response Naming**: Using Request/Response pattern for API inputs/outputs (not Input/Output)
3. **Test File Location**: Test files located in same directory as implementation (e.g., `models/vps/volumes/volume_test.go`)
4. **Volume Status Type**: Using custom VolumeStatus type for type safety (not plain string)
5. **No Client-Side Size Validation**: Size limits are validated server-side only
6. **No Pagination**: List operations return complete result set in single response
7. **IDName Reuse**: Leverages existing `common.IDName` type for consistency
8. **Action Constants**: Provides type-safe VolumeAction constants to prevent typos
9. **Optional Timestamps**: CreatedAt/UpdatedAt are pointers to handle missing values

---

## Future Considerations

- If pagination is added later, extend ListVolumesOptions with Limit/Offset/Marker fields
- If client-side validation is desired, add Validate() methods to request structs
- If additional volume types emerge, consider VolumeType struct (currently just strings)
- Volume encryption fields not present in current API, may be added in future
- Additional VolumeStatus values may be added as API evolves
