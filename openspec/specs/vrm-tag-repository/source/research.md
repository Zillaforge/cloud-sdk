# Research: VRM Tag and Repository APIs

**Phase**: 0 - Outline & Research  
**Date**: 2025-11-12  
**Status**: Complete

## Research Tasks

### 1. VPS Service Pattern Analysis

**Decision**: Follow VPS service architecture pattern exactly  
**Rationale**: 
- Maintains consistency across all service clients in the SDK
- Proven pattern already in production (VPS flavors, networks, floating IPs, etc.)
- Users expect same mental model across services
- Reduces learning curve and cognitive load

**Pattern Components**:
```go
// Architecture layers:
1. Root client extension: cloudsdk.ProjectClient.VRM()
2. Service client: modules/vrm/client.go (project-scoped)
3. Resource clients: modules/vrm/repositories/client.go, modules/vrm/tags/client.go
4. Data models: models/vrm/repositories/, models/vrm/tags/
5. Shared utilities: internal/http, internal/backoff, internal/types
```

**Key Decisions**:
- Use `internal/http.Client` for all HTTP operations (built-in retry, timeout, error wrapping)
- Project-scoped access pattern matches VPS: `client.Project(projectID).VRM()`
- Sub-clients for each resource type: `.Repositories()` and `.Tags()`
- Request/Response naming convention (not Input/Output per user requirement)

**Alternatives Considered**:
- Standalone VRM client: Rejected - breaks consistency with VPS pattern
- Global service client: Rejected - VRM is project-scoped per API spec
- Direct endpoint methods on VRM client: Rejected - violates separation of concerns

---

### 2. VRM API Specification Analysis (vrm.yaml)

**Decision**: Implement only User/Repository and User/Tag endpoints  
**Rationale**: 
- User requirement explicitly excludes Admin/*, MemberAcl, ProjectAcl, Image, Export, Snapshot
- Focus on core registry management operations (CRUD for repositories and tags)
- Reduces initial scope while providing complete functionality for target use cases

**Endpoints to Implement**:

**Repositories (5 operations)**:
- `GET /project/{project-id}/repositories` - List repositories with filters
- `POST /project/{project-id}/repository` - Create repository
- `GET /project/{project-id}/repository/{repository-id}` - Get repository
- `PUT /project/{project-id}/repository/{repository-id}` - Update repository
- `DELETE /project/{project-id}/repository/{repository-id}` - Delete repository

**Tags (6 operations)**:
- `GET /project/{project-id}/tags` - List all accessible tags
- `GET /project/{project-id}/repository/{repository-id}/tags` - List repository tags
- `POST /project/{project-id}/repository/{repository-id}/tag` - Create tag
- `GET /project/{project-id}/tag/{tag-id}` - Get tag
- `PUT /project/{project-id}/tag/{tag-id}` - Update tag
- `DELETE /project/{project-id}/tag/{tag-id}` - Delete tag

**Query Parameters**:
- Pagination: `limit` (default 100, -1 for all), `offset` (default 0)
- Filtering: `where` array for field-based filters
- Namespace: `X-Namespace` header for public/private scope

**Authentication**:
- Bearer token via `Authorization: Bearer {token}` header
- Inherited from parent cloudsdk.Client via internal/http.Client

---

### 3. Data Model Design from vrm.yaml

**Decision**: Match API specification field names exactly in JSON tags  
**Rationale**:
- Ensures correct serialization/deserialization with API
- Follows VPS model pattern (camelCase in JSON tags matching API)
- User requirement: struct fields match swagger spec

**Repository Model Fields** (from vrm.yaml examples):
```go
type Repository struct {
    ID              string    `json:"id"`
    Name            string    `json:"name"`
    Namespace       string    `json:"namespace"`        // "public" | "private"
    OperatingSystem string    `json:"operatingSystem"`  // "linux" | "windows"
    Description     string    `json:"description,omitempty"`
    Tags            []*Tag    `json:"tags,omitempty"`
    Count           int       `json:"count"`            // Number of tags
    Creator         IDName    `json:"creator"`
    Project         IDName    `json:"project"`
    CreatedAt       time.Time `json:"createdAt"`
    UpdatedAt       time.Time `json:"updatedAt"`
}
```

**Tag Model Fields** (from vrm.yaml examples):
```go
type Tag struct {
    ID           string                 `json:"id"`
    Name         string                 `json:"name"`
    RepositoryID string                 `json:"repositoryID"`
    Type         string                 `json:"type"`         // "common" | "increase" | etc.
    Size         int64                  `json:"size"`
    Status       string                 `json:"status,omitempty"`
    Extra        map[string]interface{} `json:"extra,omitempty"`
    CreatedAt    time.Time              `json:"createdAt"`
    UpdatedAt    time.Time              `json:"updatedAt"`
    Repository   Repository             `json:"repository"`
}
```

**Common Types**:
```go
type IDName struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Account     string `json:"account,omitempty"`
    DisplayName string `json:"displayName,omitempty"`
}
```

**Disk Format Enum** (from vrm.yaml):
- Values: ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, iso

**Container Format Enum** (from vrm.yaml):
- Values: ami, ari, aki, bare, ovf

---

### 4. Request/Response Struct Naming

**Decision**: Use Request/Response suffix (not Input/Output)  
**Rationale**: 
- User requirement explicitly states "API 輸入/輸出 的struct 名稱用 Request/Response 命名，不要用Input/Ouput"
- Consistent with HTTP terminology (HTTP request/response)
- Clear distinction from domain models

**Naming Pattern**:
```go
// Create operations
type CreateRepositoryRequest struct { ... }
type CreateRepositoryResponse struct { ... } // or just *Repository

// Update operations
type UpdateRepositoryRequest struct { ... }
type UpdateRepositoryResponse struct { ... } // or just *Repository

// List operations
type ListRepositoriesOptions struct { ... }   // Filter/pagination options
type ListRepositoriesResponse struct { ... }  // Contains []*Repository

// Note: List returns []*Resource directly per user requirement
// "List API 回傳 []*Resource，可以直接用 index 存取 元素"
func (c *Client) List(ctx context.Context, opts *ListRepositoriesOptions) ([]*Repository, error)
```

---

### 5. Testing Strategy (80%+ Coverage)

**Decision**: Three-tier testing approach matching VPS pattern  
**Rationale**:
- User requirement: "採TDD開發，要有unit test、contract test、integration test，test coverage要80%以上"
- Follows TDD constitution principle
- Comprehensive coverage across all test types

**Test Types**:

**Unit Tests** (`*_test.go` alongside implementation):
- Model validation (required fields, data types, constraints)
- JSON marshaling/unmarshaling
- Client method logic (parameter handling, path construction)
- Error handling paths
- Edge cases (nil pointers, empty values, etc.)

**Contract Tests** (`test/contract_test.go`):
- Validate requests match vrm.yaml OpenAPI spec
- Validate responses can deserialize from vrm.yaml examples
- Verify HTTP methods, paths, headers, query parameters
- Check required vs optional fields
- Enum value validation

**Integration Tests** (`test/integration_test.go`):
- End-to-end API flows with `httptest` mock server
- Full CRUD lifecycle (Create → Get → Update → List → Delete)
- Pagination and filtering
- Error response handling (404, 400, 401, etc.)
- Namespace header behavior
- Concurrent request safety

**Coverage Targets**:
- Models: 90%+ (simple validation logic, easy to test)
- Clients: 80%+ (HTTP logic covered by mocking)
- Overall: 80%+ minimum per requirement

---

### 6. Error Handling Pattern

**Decision**: Follow VPS error handling pattern with wrapped errors  
**Rationale**:
- User requirement: "API error 處理和 vps相同"
- Consistent error handling across SDK
- Uses `internal/http.Client` which already wraps errors in SDKError types

**Pattern**:
```go
func (c *Client) List(ctx context.Context, opts *ListOptions) ([]*Repository, error) {
    // ... build request ...
    
    var response RepositoryListResponse
    if err := c.baseClient.Do(ctx, req, &response); err != nil {
        return nil, fmt.Errorf("failed to list repositories: %w", err)
    }
    
    return response.Repositories, nil
}
```

**Error Types** (from `internal/types`):
- `SDKError`: General SDK error with status code, error code, message, metadata
- `HTTPError`: HTTP-specific error with status code
- `NetworkError`: Network connectivity issues
- `TimeoutError`: Request timeout
- `CanceledError`: Context cancellation

**Error Handling Features**:
- Automatic retry for 503, 429, network errors (via internal/http backoff)
- Context cancellation support
- Wrapped errors preserve stack trace
- Status codes propagated from API responses

---

### 7. List API Design

**Decision**: Return `[]*Resource` directly (not wrapped in response struct)  
**Rationale**:
- User requirement: "List API 回傳 []*Resource，可以直接用 index 存取 元素"
- Simpler API surface for consumers
- Direct array access without unwrapping

**Implementation**:
```go
// Internal response struct (not exported)
type repositoryListResponse struct {
    Repositories []*Repository `json:"repositories"`
    Total        int           `json:"total,omitempty"`
}

// Public method returns slice directly
func (c *Client) List(ctx context.Context, opts *ListRepositoriesOptions) ([]*Repository, error) {
    var response repositoryListResponse
    if err := c.baseClient.Do(ctx, req, &response); err != nil {
        return nil, fmt.Errorf("failed to list repositories: %w", err)
    }
    return response.Repositories, nil
}

// Usage: Direct index access
repos, err := client.Repositories().List(ctx, opts)
repo := repos[0]  // Direct access
```

**List Options Pattern**:
```go
type ListRepositoriesOptions struct {
    Limit     int      // Pagination limit (default 100, -1 for all)
    Offset    int      // Pagination offset
    Where     []string // Field filters (e.g., ["os=linux", "creator=user-id"])
    Namespace string   // "public" | "private"
}
```

---

### 8. Namespace Header Support

**Decision**: Pass namespace as optional parameter to each operation  
**Rationale**:
- vrm.yaml specifies `X-Namespace` header for public/private scope
- Optional parameter matches API optional header
- Allows per-request namespace control

**Implementation**:
```go
func (c *Client) List(ctx context.Context, opts *ListRepositoriesOptions) ([]*Repository, error) {
    req := &internalhttp.Request{
        Method: "GET",
        Path:   path,
        Headers: make(map[string]string),
    }
    
    if opts != nil && opts.Namespace != "" {
        req.Headers["X-Namespace"] = opts.Namespace
    }
    
    // ...
}
```

---

### 9. Filter Query Parameter Design

**Decision**: Support `where` array parameter in ListOptions  
**Rationale**:
- User requirement: "swagger 文件中，List API 有 filter的 query parameter 要提供在 ListOption struct中"
- vrm.yaml specifies `where` as array of strings for field-based filtering

**Supported Filters** (from vrm.yaml):
- Repositories: `os`, `creator`, `project-id`
- Tags: `status`, `type`, `project-id`

**Implementation**:
```go
type ListRepositoriesOptions struct {
    Limit     int
    Offset    int
    Where     []string  // e.g., ["os=linux", "creator=user-123"]
    Namespace string
}

// Build query string
query := url.Values{}
for _, filter := range opts.Where {
    query.Add("where", filter)
}
```

---

### 10. Phase Validation with `make check`

**Decision**: Validate each phase with `make check` before proceeding  
**Rationale**:
- User requirement: "每個phase 的設計要經過 `make check`的驗證"
- Ensures code quality and correctness at each step
- Catches issues early in development cycle

**Make Check Targets** (inferred from Go project):
```bash
make check  # Likely runs:
  - go fmt ./...
  - go vet ./...
  - go test -race ./...
  - golangci-lint run
  - go mod tidy && go mod verify
```

**Phase Validation Points**:
- Phase 0: N/A (research only, no code)
- Phase 1: Validate data models and contracts
- Phase 2: Validate implementation (post-/speckit.tasks)

---

## Research Summary

All technical unknowns have been resolved through analysis of:
1. Existing VPS service implementation pattern
2. vrm.yaml OpenAPI specification
3. User requirements and naming conventions
4. Constitution principles and testing standards

**Key Findings**:
- Zero external dependencies required (stdlib + internal packages sufficient)
- Exact pattern match with VPS service ensures consistency
- 11 API endpoints (5 repository + 6 tag operations)
- Request/Response naming convention
- List APIs return slices directly for index access
- 80%+ test coverage achievable with three-tier testing
- `make check` validation at each phase

**Ready for Phase 1**: Data model and contract design can proceed with confidence.
