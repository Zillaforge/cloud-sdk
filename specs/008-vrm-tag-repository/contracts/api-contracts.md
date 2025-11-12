# VRM API Contracts

**Phase**: 1 - Design & Contracts  
**Date**: 2025-11-12  
**Source**: vrm.yaml OpenAPI 3.0 specification

This document defines the API contracts for VRM Repository and Tag operations, derived from the vrm.yaml swagger specification. These contracts serve as the specification for contract tests.

---

## Repository API Contracts

### 1. List Repositories

**Contract ID**: `VRM-REPO-LIST-001`  
**Endpoint**: `GET /vrm/api/v1/project/{project-id}/repositories`  
**Authentication**: Bearer token required  
**Description**: Retrieve a paginated list of repositories within a project

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |

**Query Parameters**:
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| limit | integer | ❌ | 100 | Max items to return (-1 for all) |
| offset | integer | ❌ | 0 | Number of items to skip |
| where | array[string] | ❌ | [] | Field filters (e.g., "os=linux") |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| X-Namespace | string | ❌ | Namespace scope ("public" or "private") |

**Success Response** (200 OK):
```json
{
  "repositories": [
    {
      "id": "59e0e12a-c857-44a4-88b2-2aa8baec4e00",
      "name": "ubuntu",
      "namespace": "public",
      "operatingSystem": "linux",
      "description": "",
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
  ],
  "total": 1
}
```

**Error Responses**:
- `403 Forbidden`: Invalid or missing authentication
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) List(ctx context.Context, opts *ListRepositoriesOptions) ([]*Repository, error)
```

---

### 2. Create Repository

**Contract ID**: `VRM-REPO-CREATE-001`  
**Endpoint**: `POST /vrm/api/v1/project/{project-id}/repository`  
**Authentication**: Bearer token required  
**Description**: Create a new repository in the project

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| Content-Type | string | ✅ | application/json |
| X-Namespace | string | ❌ | Namespace scope ("public" or "private") |

**Request Body**:
```json
{
  "name": "ubuntu",
  "operatingSystem": "linux",
  "description": "Ubuntu base images"
}
```

**Request Schema**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | ✅ | Repository name |
| operatingSystem | string | ✅ | "linux" or "windows" |
| description | string | ❌ | Optional description |

**Success Response** (201 Created):
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

**Error Responses**:
- `400 Bad Request`: Invalid request body or missing required fields
- `403 Forbidden`: Invalid authentication or permissions
- `404 Not Found`: Project not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) Create(ctx context.Context, req *CreateRepositoryRequest) (*Repository, error)
```

---

### 3. Get Repository

**Contract ID**: `VRM-REPO-GET-001`  
**Endpoint**: `GET /vrm/api/v1/project/{project-id}/repository/{repository-id}`  
**Authentication**: Bearer token required  
**Description**: Retrieve details of a specific repository

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |
| repository-id | string | ✅ | Repository UUID |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| X-Namespace | string | ❌ | Namespace scope |

**Success Response** (200 OK):
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

**Error Responses**:
- `403 Forbidden`: Invalid authentication or permissions
- `404 Not Found`: Repository not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) Get(ctx context.Context, repositoryID string) (*Repository, error)
```

---

### 4. Update Repository

**Contract ID**: `VRM-REPO-UPDATE-001`  
**Endpoint**: `PUT /vrm/api/v1/project/{project-id}/repository/{repository-id}`  
**Authentication**: Bearer token required  
**Description**: Update repository metadata

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |
| repository-id | string | ✅ | Repository UUID |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| Content-Type | string | ✅ | application/json |
| X-Namespace | string | ❌ | Namespace scope |

**Request Body**:
```json
{
  "description": "Updated description"
}
```

**Request Schema**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | ❌ | New repository name |
| description | string | ❌ | New description |

**Success Response** (200 OK):
```json
{
  "id": "59e0e12a-c857-44a4-88b2-2aa8baec4e00",
  "name": "ubuntu",
  "namespace": "public",
  "operatingSystem": "linux",
  "description": "Updated description",
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
  "updatedAt": "2024-08-19T08:32:25Z"
}
```

**Error Responses**:
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Authentication failure
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Repository not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) Update(ctx context.Context, repositoryID string, req *UpdateRepositoryRequest) (*Repository, error)
```

---

### 5. Delete Repository

**Contract ID**: `VRM-REPO-DELETE-001`  
**Endpoint**: `DELETE /vrm/api/v1/project/{project-id}/repository/{repository-id}`  
**Authentication**: Bearer token required  
**Description**: Delete a repository (soft delete)

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |
| repository-id | string | ✅ | Repository UUID |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| X-Namespace | string | ❌ | Namespace scope |

**Success Response** (204 No Content):
- Empty response body

**Error Responses**:
- `403 Forbidden`: Invalid authentication or permissions
- `404 Not Found`: Repository not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) Delete(ctx context.Context, repositoryID string) error
```

---

## Tag API Contracts

### 6. List All Tags

**Contract ID**: `VRM-TAG-LIST-ALL-001`  
**Endpoint**: `GET /vrm/api/v1/project/{project-id}/tags`  
**Authentication**: Bearer token required  
**Description**: Retrieve all accessible tags in the project

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |

**Query Parameters**:
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| limit | integer | ❌ | 100 | Max items to return (-1 for all) |
| offset | integer | ❌ | 0 | Number of items to skip |
| where | array[string] | ❌ | [] | Field filters (e.g., "type=common") |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| X-Namespace | string | ❌ | Namespace scope |

**Success Response** (200 OK):
```json
{
  "tags": [
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
  ],
  "total": 1
}
```

**Error Responses**:
- `403 Forbidden`: Invalid authentication
- `404 Not Found`: Project not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) List(ctx context.Context, opts *ListTagsOptions) ([]*Tag, error)
```

---

### 7. List Repository Tags

**Contract ID**: `VRM-TAG-LIST-REPO-001`  
**Endpoint**: `GET /vrm/api/v1/project/{project-id}/repository/{repository-id}/tags`  
**Authentication**: Bearer token required  
**Description**: Retrieve tags for a specific repository

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |
| repository-id | string | ✅ | Repository UUID |

**Query Parameters**:
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| limit | integer | ❌ | 100 | Max items to return |
| offset | integer | ❌ | 0 | Number of items to skip |
| where | array[string] | ❌ | [] | Field filters (e.g., "status=active", "type=common") |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| X-Namespace | string | ❌ | Namespace scope |

**Success Response** (200 OK): Same structure as List All Tags

**Error Responses**:
- `403 Forbidden`: Invalid authentication
- `404 Not Found`: Repository not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) ListByRepository(ctx context.Context, repositoryID string, opts *ListTagsOptions) ([]*Tag, error)
```

---

### 8. Create Tag

**Contract ID**: `VRM-TAG-CREATE-001`  
**Endpoint**: `POST /vrm/api/v1/project/{project-id}/repository/{repository-id}/tag`  
**Authentication**: Bearer token required  
**Description**: Create a new tag in a repository

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |
| repository-id | string | ✅ | Repository UUID |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| Content-Type | string | ✅ | application/json |
| X-Namespace | string | ❌ | Namespace scope |

**Request Body**:
```json
{
  "name": "v1",
  "type": "common",
  "diskFormat": "raw",
  "containerFormat": "bare"
}
```

**Request Schema**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | ✅ | Tag name (version label) |
| type | string | ✅ | Tag type (e.g., "common", "increase") |
| diskFormat | string | ✅ | ami\|ari\|aki\|vhd\|vmdk\|raw\|qcow2\|vdi\|iso |
| containerFormat | string | ✅ | ami\|ari\|aki\|bare\|ovf |

**Success Response** (201 Created):
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

**Error Responses**:
- `400 Bad Request`: Invalid request body or missing required fields
- `403 Forbidden`: Invalid authentication or permissions
- `404 Not Found`: Repository not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) Create(ctx context.Context, repositoryID string, req *CreateTagRequest) (*Tag, error)
```

---

### 9. Get Tag

**Contract ID**: `VRM-TAG-GET-001`  
**Endpoint**: `GET /vrm/api/v1/project/{project-id}/tag/{tag-id}`  
**Authentication**: Bearer token required  
**Description**: Retrieve details of a specific tag

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |
| tag-id | string | ✅ | Tag UUID |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| X-Namespace | string | ❌ | Namespace scope |

**Success Response** (200 OK): Same structure as Create Tag response

**Error Responses**:
- `403 Forbidden`: Invalid authentication or permissions
- `404 Not Found`: Tag not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) Get(ctx context.Context, tagID string) (*Tag, error)
```

---

### 10. Update Tag

**Contract ID**: `VRM-TAG-UPDATE-001`  
**Endpoint**: `PUT /vrm/api/v1/project/{project-id}/tag/{tag-id}`  
**Authentication**: Bearer token required  
**Description**: Update tag metadata

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |
| tag-id | string | ✅ | Tag UUID |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| Content-Type | string | ✅ | application/json |
| X-Namespace | string | ❌ | Namespace scope |

**Request Body**:
```json
{
  "name": "v1-updated"
}
```

**Request Schema**:
| Field | Type | Required | Description |
|-------|------|----------|-------------|
| name | string | ❌ | New tag name |

**Success Response** (200 OK): Same structure as Get Tag response

**Error Responses**:
- `400 Bad Request`: Invalid request body
- `401 Unauthorized`: Authentication failure
- `403 Forbidden`: Insufficient permissions
- `404 Not Found`: Tag not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) Update(ctx context.Context, tagID string, req *UpdateTagRequest) (*Tag, error)
```

---

### 11. Delete Tag

**Contract ID**: `VRM-TAG-DELETE-001`  
**Endpoint**: `DELETE /vrm/api/v1/project/{project-id}/tag/{tag-id}`  
**Authentication**: Bearer token required  
**Description**: Delete a tag (soft delete)

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string | ✅ | Project UUID |
| tag-id | string | ✅ | Tag UUID |

**Headers**:
| Header | Type | Required | Description |
|--------|------|----------|-------------|
| Authorization | string | ✅ | Bearer {token} |
| X-Namespace | string | ❌ | Namespace scope |

**Success Response** (204 No Content):
- Empty response body

**Error Responses**:
- `403 Forbidden`: Invalid authentication or permissions
- `404 Not Found`: Tag not found
- `500 Internal Server Error`: Server error

**Go Client Method**:
```go
func (c *Client) Delete(ctx context.Context, tagID string) error
```

---

## Contract Testing Requirements

### Test Coverage

Each contract MUST be validated by automated tests:

1. **Request Validation**:
   - HTTP method matches specification
   - Path construction with parameters
   - Query parameters formatted correctly
   - Headers include Authorization and optional X-Namespace
   - Request body matches JSON schema

2. **Response Validation**:
   - Status code matches expected value
   - Response body deserializes to Go struct
   - All required fields present
   - Field types match specification
   - Nested objects properly populated

3. **Error Cases**:
   - 400 Bad Request for invalid inputs
   - 401/403 for authentication/authorization failures
   - 404 for non-existent resources
   - 500 for server errors

### Test Implementation

Contract tests located in:
- `modules/vrm/repositories/test/contract_test.go`
- `modules/vrm/tags/test/contract_test.go`

Each test MUST:
- Reference the contract ID
- Use `httptest` to mock API responses
- Validate request construction
- Validate response deserialization
- Cover success and error paths

### Success Criteria

- All 11 contracts have automated tests
- Tests pass with 100% success rate
- Request/response structures match vrm.yaml exactly
- Error handling covers all documented error codes
