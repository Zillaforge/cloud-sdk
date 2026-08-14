# Data Model: VRM Tag and Repository APIs

**Phase**: 1 - Design & Contracts  
**Date**: 2025-11-12  
**Source**: vrm.yaml OpenAPI specification

## Entity Overview

This data model defines the core entities for Virtual Registry Management (VRM) API, focusing on Repository and Tag resources. All models match the vrm.yaml API specification exactly.

### Entity Relationship

```
Project (external)
    ├── Repository (1:N)
    │   ├── Tags (1:N)
    │   ├── Creator (User reference)
    │   └── Project (reference back)
    └── Tag (direct access, 1:1 with Repository)
        ├── Repository (reference)
        └── Extra metadata (map)

IDName (common reference type)
    ├── Used for: Creator, Project, User references
    └── Fields: id, name, account, displayName
```

---

## Core Entities

### 1. Repository

**Description**: A virtual image repository within a project, containing multiple tagged versions.

**Location**: `models/vrm/repositories/repository.go`

**Fields**:

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | string | `id` | ✅ | Unique repository identifier (UUID) |
| Name | string | `name` | ✅ | Repository name (alphanumeric, hyphens, underscores) |
| Namespace | string | `namespace` | ✅ | Visibility scope: "public" or "private" |
| OperatingSystem | string | `operatingSystem` | ✅ | OS type: "linux" or "windows" |
| Description | string | `description,omitempty` | ❌ | Optional description |
| Tags | []*Tag | `tags,omitempty` | ❌ | Associated tags (populated in Get/List responses) |
| Count | int | `count` | ✅ | Number of tags in repository |
| Creator | IDName | `creator` | ✅ | User who created the repository |
| Project | IDName | `project` | ✅ | Project this repository belongs to |
| CreatedAt | time.Time | `createdAt` | ✅ | Creation timestamp (ISO 8601) |
| UpdatedAt | time.Time | `updatedAt` | ✅ | Last update timestamp (ISO 8601) |

**Validation Rules**:
- `ID`: Non-empty, valid UUID format
- `Name`: Non-empty, max 255 characters, alphanumeric + hyphen + underscore
- `Namespace`: Must be "public" or "private"
- `OperatingSystem`: Must be "linux" or "windows"
- `Count`: Must be >= 0

**State Lifecycle**:
1. **Created**: Repository created with name and OS
2. **Active**: Can create/update/delete tags
3. **Deleted**: Soft delete via DELETE endpoint (removed from listings)

**Example**:
```json
{
  "id": "59e0e12a-c857-44a4-88b2-2aa8baec4e00",
  "name": "ubuntu",
  "namespace": "public",
  "operatingSystem": "linux",
  "description": "Ubuntu base images",
  "tags": [],
  "count": 0,
  "creator": {
    "id": "4990ccdb-a9b1-49e5-91df-67c921601d81",
    "name": "system",
    "account": "system"
  },
  "project": {
    "id": "14735dfa-5553-46cc-b4bd-405e711b223f",
    "displayName": "admin"
  },
  "createdAt": "2024-08-19T08:32:15Z",
  "updatedAt": "2024-08-19T08:32:15Z"
}
```

---

### 2. Tag

**Description**: A specific version/tag of an image within a repository, with type and format specifications.

**Location**: `models/vrm/tags/tag.go`

**Fields**:

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | string | `id` | ✅ | Unique tag identifier (UUID) |
| Name | string | `name` | ✅ | Tag name (version label, e.g., "v1", "24.04") |
| RepositoryID | string | `repositoryID` | ✅ | Parent repository UUID |
| Type | string | `type` | ✅ | Tag type: "common", "increase", etc. |
| Size | int64 | `size` | ✅ | Image size in bytes |
| Status | string | `status,omitempty` | ❌ | Tag status (optional, e.g., "active", "building") |
| Extra | map[string]interface{} | `extra,omitempty` | ❌ | Additional metadata (flexible key-value pairs) |
| CreatedAt | time.Time | `createdAt` | ✅ | Creation timestamp (ISO 8601) |
| UpdatedAt | time.Time | `updatedAt` | ✅ | Last update timestamp (ISO 8601) |
| Repository | Repository | `repository` | ✅ | Parent repository details (populated in responses) |

**Validation Rules**:
- `ID`: Non-empty, valid UUID format
- `Name`: Non-empty, max 255 characters
- `RepositoryID`: Non-empty, valid UUID format
- `Type`: Non-empty (specific enum values TBD by API behavior)
- `Size`: Must be >= 0

**Image Format Fields** (used in Create requests, per vrm.yaml):
- `DiskFormat`: Enum - ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, iso
- `ContainerFormat`: Enum - ami, ari, aki, bare, ovf

**State Lifecycle**:
1. **Created**: Tag created with name, type, formats
2. **Building**: Image being uploaded/processed (optional status)
3. **Active**: Ready for use
4. **Deleted**: Soft delete via DELETE endpoint

**Example**:
```json
{
  "id": "800c66b5-d03d-407f-9a1a-e27ecc5b03a0",
  "name": "v1",
  "repositoryID": "59e0e12a-c857-44a4-88b2-2aa8baec4e00",
  "type": "common",
  "size": 0,
  "status": "",
  "extra": {},
  "createdAt": "2024-08-19T08:32:25Z",
  "updatedAt": "2024-08-19T08:32:25Z",
  "repository": {
    "id": "59e0e12a-c857-44a4-88b2-2aa8baec4e00",
    "name": "ubuntu",
    "namespace": "public",
    "operatingSystem": "linux",
    "creator": {
      "id": "4990ccdb-a9b1-49e5-91df-67c921601d81",
      "name": "system",
      "account": "system"
    },
    "project": {
      "id": "14735dfa-5553-46cc-b4bd-405e711b223f",
      "displayName": "admin"
    }
  }
}
```

---

### 3. IDName (Common Reference Type)

**Description**: Lightweight reference type for nested entities (Creator, Project, User).

**Location**: `models/vrm/common/common.go`

**Fields**:

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | string | `id` | ✅ | Entity identifier (UUID) |
| Name | string | `name,omitempty` | ❌ | Entity name |
| Account | string | `account,omitempty` | ❌ | Account name (for users) |
| DisplayName | string | `displayName,omitempty` | ❌ | Display name (for projects) |

**Usage**:
- Repository.Creator (user reference)
- Repository.Project (project reference)
- Tag.Repository.Creator (nested user reference)
- Tag.Repository.Project (nested project reference)

**Example**:
```json
{
  "id": "4990ccdb-a9b1-49e5-91df-67c921601d81",
  "name": "system",
  "account": "system"
}
```

---

## Request/Response Types

### Repository Operations

#### CreateRepositoryRequest
**Location**: `models/vrm/repositories/repository.go`
**Used By**: POST `/project/{project-id}/repository`

```go
type CreateRepositoryRequest struct {
    Name            string `json:"name"`            // Required
    OperatingSystem string `json:"operatingSystem"` // Required: "linux" | "windows"
    Description     string `json:"description,omitempty"` // Optional
}
```

#### UpdateRepositoryRequest
**Location**: `models/vrm/repositories/repository.go`
**Used By**: PUT `/project/{project-id}/repository/{repository-id}`

```go
type UpdateRepositoryRequest struct {
    Name        string `json:"name,omitempty"`        // Optional
    Description string `json:"description,omitempty"` // Optional
}
```

#### ListRepositoriesOptions
**Location**: `models/vrm/repositories/repository.go`
**Used By**: GET `/project/{project-id}/repositories`

```go
type ListRepositoriesOptions struct {
    Limit     int      // Pagination limit (default 100, -1 for all)
    Offset    int      // Pagination offset (default 0)
    Where     []string // Field filters: ["os=linux", "creator=user-id", "project-id=uuid"]
    Namespace string   // "public" | "private" (sent as X-Namespace header)
}
```

---

### Tag Operations

#### CreateTagRequest
**Location**: `models/vrm/tags/tag.go`
**Used By**: POST `/project/{project-id}/repository/{repository-id}/tag`

```go
type CreateTagRequest struct {
    Name            string `json:"name"`            // Required
    Type            string `json:"type"`            // Required: "common" | "increase" | etc.
    DiskFormat      string `json:"diskFormat"`      // Required: ami|ari|aki|vhd|vmdk|raw|qcow2|vdi|iso
    ContainerFormat string `json:"containerFormat"` // Required: ami|ari|aki|bare|ovf
}
```

#### UpdateTagRequest
**Location**: `models/vrm/tags/tag.go`
**Used By**: PUT `/project/{project-id}/tag/{tag-id}`

```go
type UpdateTagRequest struct {
    Name string `json:"name,omitempty"` // Optional
}
```

#### ListTagsOptions
**Location**: `models/vrm/tags/tag.go`
**Used By**: 
- GET `/project/{project-id}/tags`
- GET `/project/{project-id}/repository/{repository-id}/tags`

```go
type ListTagsOptions struct {
    Limit     int      // Pagination limit (default 100, -1 for all)
    Offset    int      // Pagination offset (default 0)
    Where     []string // Field filters: ["status=active", "type=common", "project-id=uuid"]
    Namespace string   // "public" | "private" (sent as X-Namespace header)
}
```

---

## Enumerations

### DiskFormat
**Location**: `models/vrm/common/common.go`

```go
type DiskFormat string

const (
    DiskFormatAMI   DiskFormat = "ami"
    DiskFormatARI   DiskFormat = "ari"
    DiskFormatAKI   DiskFormat = "aki"
    DiskFormatVHD   DiskFormat = "vhd"
    DiskFormatVMDK  DiskFormat = "vmdk"
    DiskFormatRaw   DiskFormat = "raw"
    DiskFormatQcow2 DiskFormat = "qcow2"
    DiskFormatVDI   DiskFormat = "vdi"
    DiskFormatISO   DiskFormat = "iso"
)
```

### ContainerFormat
**Location**: `models/vrm/common/common.go`

```go
type ContainerFormat string

const (
    ContainerFormatAMI  ContainerFormat = "ami"
    ContainerFormatARI  ContainerFormat = "ari"
    ContainerFormatAKI  ContainerFormat = "aki"
    ContainerFormatBare ContainerFormat = "bare"
    ContainerFormatOVF  ContainerFormat = "ovf"
)
```

---

## Data Model Validation Matrix

| Entity | Required Fields | Optional Fields | Validation Method |
|--------|----------------|-----------------|-------------------|
| Repository | ID, Name, Namespace, OS, Count, Creator, Project, Timestamps | Description, Tags | `Validate()` method |
| Tag | ID, Name, RepositoryID, Type, Size, Timestamps, Repository | Status, Extra | `Validate()` method |
| IDName | ID | Name, Account, DisplayName | N/A (all optional except ID) |
| CreateRepositoryRequest | Name, OperatingSystem | Description | `Validate()` method |
| CreateTagRequest | Name, Type, DiskFormat, ContainerFormat | None | `Validate()` method |
| ListOptions | None (all optional) | Limit, Offset, Where, Namespace | `Validate()` method |

---

## Timestamp Handling

All timestamps follow ISO 8601 / RFC3339 format as specified in the API:
- Format: `2006-01-02T15:04:05Z`
- Timezone: UTC (Z suffix)
- Type: `time.Time` in Go
- JSON: String representation via standard encoding

**Example**:
```go
createdAt, _ := time.Parse(time.RFC3339, "2024-08-19T08:32:15Z")
```

---

## Nested Object Handling

Tags include full Repository details in responses:
- **Shallow Reference**: Tag.RepositoryID (UUID string)
- **Deep Reference**: Tag.Repository (full Repository object)

This allows clients to access repository metadata without additional API calls.

**Serialization Note**: When creating/updating tags, only RepositoryID is sent. The full Repository object is only populated in API responses.

---

## Summary

**Total Entities**: 3 core (Repository, Tag, IDName) + 6 request/response types  
**Total API Endpoints**: 11 (5 repository + 6 tag operations)  
**Validation Methods**: 5 (Repository, Tag, CreateRepositoryRequest, CreateTagRequest, ListOptions)  
**Enumerations**: 2 (DiskFormat with 9 values, ContainerFormat with 5 values)

All models are designed to:
1. Match vrm.yaml API specification exactly
2. Support JSON marshaling/unmarshaling
3. Provide validation for required fields and data types
4. Enable direct index access for List APIs (return `[]*Resource`)
5. Use Request/Response naming convention (not Input/Output)
