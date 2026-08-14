# Research: IAM API Client

**Feature**: 009-iam-api  
**Date**: 2025-11-12  
**Phase**: 0 - Research & Technical Decisions

## Overview

This document consolidates research findings and technical decisions for implementing the IAM API client, resolving all technical unknowns identified during planning.

## Technical Decisions

### Decision 1: IAM Client Architecture Pattern

**Decision**: **Resources().Verb() pattern** - Two resources (User, Projects) each with specific verbs

**Rationale**:
- IAM has two distinct resource types:
  - **User** (singular): Represents current authenticated user, single verb `Get()`
  - **Projects** (plural): Represents project collection, verbs `List()` and `Get(projectID)`
- Follows RESTful resource-oriented design principles
- Clear separation: `client.IAM().User().Get()` vs `client.IAM().Projects().List()`
- Self-documenting API: resource type is explicit in call chain
- Unlike VPS/VRM which are project-scoped, IAM is user-scoped and accessed directly

**Alternatives Considered**:
- **Flat client methods** (`client.IAM().GetUser()`, `client.IAM().ListProjects()`): Rejected because it mixes different resource concerns
- **Project-scoped pattern**: Rejected because IAM APIs don't require project context
- **Separate standalone client**: Rejected to maintain consistency with SDK's unified client entry point

**Implementation Impact**:
- Add `IAM()` method to root `Client` struct in `client.go`
- Create `modules/iam/client.go` with IAM base client
- **Create `modules/iam/users/client.go`** with User resource client exposing `Get()` method (following VPS/VRM pattern with subdirectories)
- **Create `modules/iam/projects/client.go`** with Projects resource client exposing `List()` and `Get(projectID)` methods
- **Directory Structure**: Mirror VPS/VRM with separate subdirectories for each resource type:
  - `modules/iam/users/` (like `modules/vps/flavors/`)
  - `modules/iam/projects/` (like `modules/vps/servers/`)
  - `models/iam/users/` (like `models/vps/flavors/`)
  - `models/iam/projects/` (like `models/vps/servers/`)

```go
// Usage pattern
client := cloudsdk.NewClient(baseURL, token)

// User resource - only Get()
user, err := client.IAM().User().Get(ctx)

// Projects resource - List() and Get()
projects, err := client.IAM().Projects().List(ctx, iam.WithLimit(50))
project, err := client.IAM().Projects().Get(ctx, projectID)
```

### Decision 2: Model Structure and Naming

**Decision**: Use Swagger-aligned struct fields with Request/Response suffix for operation DTOs

**Rationale**:
- Field names must match Swagger spec exactly for automatic JSON marshaling
- Request/Response naming convention (not Input/Output) aligns with Go HTTP best practices
- Example: `GetUserResponse`, `ListProjectsRequest`, `ListProjectsResponse`
- Nested structs like `Permission`, `UsageTime` match Swagger schema names exactly

**Alternatives Considered**:
- **Custom field names with json tags**: Rejected for simplicity and maintainability
- **Input/Output suffix**: Rejected per project requirements

**Implementation Impact**:
- `models/iam/common/common.go` - shared types (Permission, UsageTime, timestamps)
- **`models/iam/users/models.go`** - User, GetUserResponse (following VPS/VRM naming: models.go)
- **`models/iam/projects/models.go`** - Project, ProjectMembership, ListProjectsResponse, GetProjectResponse
- **`models/iam/projects/options.go`** - ListProjectsOptions, functional option helpers
- All struct fields use exact Swagger field names (projectId, displayName, etc.)
- **Pattern**: Each resource has its own directory under models/iam/ (like vps/flavors/, vps/servers/)

### Decision 3: List API Response Parsing Strategy

**Decision**: Custom UnmarshalJSON for list responses to extract arrays from wrapper objects

**Rationale**:
- Swagger defines ListProjects response as `{"projects": [...], "total": N}`
- Direct unmarshal to `[]*Project` fails because JSON is wrapped in object
- Custom unmarshaler extracts array from "projects" field, provides both list and total count
- Pattern already established in VPS module (flavors, networks)

**Alternatives Considered**:
- **Return wrapper struct directly**: Rejected because it forces users to access `.Projects` field with range loops
- **Two separate methods**: Rejected as unnecessarily complex

**Implementation Impact**:
```go
type ListProjectsResponse struct {
    Projects []*ProjectMembership `json:"projects"`
    Total    int                  `json:"total"`
}

// Allow direct indexing: projects[0] while preserving total count
func (r *ListProjectsResponse) UnmarshalJSON(data []byte) error {
    // Parse wrapper object, extract projects array
}
```

### Decision 4: Error Handling Strategy

**Decision**: Use existing internal/types.SDKError with fmt.Errorf wrapping

**Rationale**:
- Consistent with VPS/VRM modules error handling
- Leverages existing retry logic in internal/http/client.go
- Provides status codes, error codes, and descriptive messages
- Pattern: `if err := c.baseClient.Do(ctx, req, &response); err != nil { return nil, fmt.Errorf("failed to get user: %w", err) }`

**Alternatives Considered**:
- **Custom IAM error types**: Rejected for unnecessary complexity and inconsistency
- **Direct error return**: Rejected as it loses context about operation that failed

**Implementation Impact**:
- No new error types needed
- All methods wrap errors with operation context
- Automatic retry for 429, 502, 503, 504 status codes (already implemented)

### Decision 5: Test Strategy and Coverage

**Decision**: TDD with unit, contract, and integration tests targeting 80%+ coverage

**Rationale**:
- Constitution requires test-first development
- Contract tests validate against Swagger spec (iam.yaml)
- Unit tests cover model parsing, options building, error cases
- Integration tests verify real API interactions
- 80% coverage threshold ensures comprehensive test suite

**Test Structure** (following VPS/VRM pattern):
```
models/iam/
├── common/
│   └── common_test.go      # Tests for shared types
├── users/
│   └── models_test.go      # Tests for user models (JSON parsing)
└── projects/
    └── models_test.go      # Tests for project models (JSON parsing, unmarshaler)

modules/iam/
├── client_test.go          # Tests for base IAM client
├── users/
│   └── client_test.go      # Tests for user resource operations
└── projects/
    └── client_test.go      # Tests for projects resource operations

tests/
├── contract/
│   └── iam_contract_test.go  # Swagger contract validation
└── integration/
    └── iam_integration_test.go  # Live API tests
```

**Implementation Approach**:
1. Write failing tests for each API method
2. Implement models to pass parsing tests
3. Implement client methods to pass API tests
4. Run `make check` to verify coverage

### Decision 6: Pagination Options Pattern

**Decision**: Single options struct pointer for ListProjects pagination (consistent with VPS/VRM)

**Rationale**:
- VPS/VRM use single struct pointer pattern: `List(ctx, *ListOptions)`
- Example: `flavors.List(ctx, *flavors.ListFlavorsOptions)`, `servers.List(ctx, *servers.ServersListRequest)`
- Simpler than functional options - just create struct with desired fields
- Consistent across all SDK modules
- nil options means use defaults

**Implementation**:
```go
type ListProjectsOptions struct {
    Offset *int
    Limit  *int
    Order  *string
}

func (c *Client) List(ctx context.Context, opts *ListProjectsOptions) (*ListProjectsResponse, error) {
    // Handle nil opts - use defaults
    // Build query parameters from opts fields
}

// Usage
opts := &ListProjectsOptions{Offset: intPtr(10), Limit: intPtr(5)}
projects, err := client.IAM().Projects().List(ctx, opts)

// Or pass nil for defaults
projects, err := client.IAM().Projects().List(ctx, nil)
```

**Alternatives Considered**:
- **Functional options**: Rejected for inconsistency with VPS/VRM
- **Separate methods**: Rejected as unnecessarily complex

### Decision 7: Base URL Configuration

**Decision**: Hard-code IAM base URL with override support via client options

**Rationale**:
- Spec defines default: http://127.0.0.1:8084/iam/api/v1
- IAM service has different base URL than VPS/VRM (/iam/api/v1 vs /vps/api/v1)
- Allow override for different environments (staging, production)
- Pattern: `NewClient(baseURL, token).IAM()` where IAM() constructs IAM-specific URL

**Implementation**:
```go
// In client.go
func (c *Client) IAM() *iam.Client {
    iamBaseURL := c.baseURL + "/iam/api/v1"
    return iam.NewClient(iamBaseURL, c.token, c.httpClient)
}
```

### Decision 8: HTTP Client and Logger Pattern

**Decision**: Use `internalhttp.Client` with embedded logger (consistent with VPS/VRM)

**Rationale**:
- VPS/VRM modules use `internalhttp.Client` for all HTTP operations
- `internalhttp.Client` encapsulates HTTP client, retry logic, and logger
- IAM client constructor receives `*internalhttp.Client` instead of separate parameters
- Logger is configured at root SDK client level, passed down through `internalhttp.Client`
- Maintains consistent behavior across all modules

**Implementation**:
```go
// modules/iam/client.go
type Client struct {
    baseClient *internalhttp.Client
}

func NewClient(baseClient *internalhttp.Client) *Client {
    return &Client{baseClient: baseClient}
}

// modules/iam/users/client.go
type Client struct {
    baseClient *internalhttp.Client
}

func NewClient(baseClient *internalhttp.Client) *Client {
    return &Client{baseClient: baseClient}
}

// modules/iam/projects/client.go
type Client struct {
    baseClient *internalhttp.Client
}

func NewClient(baseClient *internalhttp.Client) *Client {
    return &Client{baseClient: baseClient}
}
```

**Alternatives Considered**:
1. **Separate logger parameter** - Rejected for inconsistency with VPS/VRM pattern
2. **No logger** - Rejected as VPS/VRM use logger for internal operations
3. **Expose logger in IAM API** - Rejected to keep clean API surface

## Technology Stack Summary

| Component | Technology | Justification |
|-----------|-----------|---------------|
| Language | Go 1.21+ | Project standard |
| HTTP Client | internalhttp.Client | Consistent with VPS/VRM, encapsulates retry + logger |
| JSON Parsing | encoding/json (stdlib) | Sufficient for IAM API needs |
| Testing | Go testing (stdlib) in resource dirs | Constitution requirement: TDD with stdlib, tests alongside code |
| Contract Testing | Manual Swagger validation | Verify responses match iam.yaml schemas |
| Retry Logic | internal/backoff (via internalhttp.Client) | Existing exponential backoff with jitter |
| Error Types | internal/types.SDKError | Existing error wrapping, consistent with VPS/VRM |
| Logger | internal/types.Logger (via internalhttp.Client) | Consistent with VPS/VRM pattern |
| Pagination | Single options struct pointer | Consistent with VPS/VRM (flavors.ListFlavorsOptions, servers.ServersListRequest) |

## Best Practices Applied

### 1. Context Propagation
- All public methods accept `context.Context` as first parameter
- Respect context cancellation and timeouts
- Pass context through to internal/http client

### 2. Options Struct Pattern
- Single options struct pointer for List methods (consistent with VPS/VRM)
- nil options means use server defaults
- Options struct allows adding fields without breaking changes
- Example: `List(ctx, *ListProjectsOptions)` where opts can be nil

### 3. Pointer Semantics
- Return `*User`, `*Project` for nullable responses
- Use `*int`, `*string` for optional request fields
- List methods return `[]*Type` not `[]Type` for consistency

### 4. Error Wrapping
- Wrap all errors with operation context
- Preserve underlying SDKError for status code access
- Use `fmt.Errorf` with `%w` verb for error chains

### 5. No Global State
- All configuration passed via constructors
- Client is explicitly constructed and passed
- No package-level variables for configuration

## Open Questions

None - all technical decisions resolved.

## References

- IAM Swagger Spec: `/workspaces/cloud-sdk/swagger/iam.yaml`
- Existing VPS Client: `/workspaces/cloud-sdk/modules/vps/client.go`
- Existing VRM Client: `/workspaces/cloud-sdk/modules/vrm/client.go`
- Root Client: `/workspaces/cloud-sdk/client.go`
- Internal HTTP Client: `/workspaces/cloud-sdk/internal/http/client.go`
- Backoff Strategy: `/workspaces/cloud-sdk/internal/backoff/backoff.go`
- Constitution: `/workspaces/cloud-sdk/.specify/memory/constitution.md`
