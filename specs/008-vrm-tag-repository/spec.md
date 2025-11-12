# Feature Specification: VRM Tag and Repository APIs Client SDK

**Feature Branch**: `008-vrm-tag-repository`  
**Created**: November 12, 2025  
**Status**: Draft  
**Input**: User description: "implement virtual registry service (vps) restful api based on #file:vrm.yaml . Only implement api which has `User/Tag` and `User/Repository`, and ignore others. All APIs are authorized through Bearer Token, and client has to provide the token when initilization. 編號用008"

## Clarifications

### Session 2025-11-12

- Q: What HTTP request timeout should the SDK client use for API calls to prevent indefinite hangs? → A: 30-60 seconds with configurable override (following existing VPS client pattern with 30s default)
- Q: Should the VRM client support concurrent API calls from multiple goroutines safely? → A: Thread-safe by design using shared internal HTTP client (matches VPS pattern)
- Q: How should the VRM client be instantiated and accessed by SDK users? → A: Project-scoped via cloudsdk.Client.Project(projectID).VRM() (maintains architectural consistency with VPS service)

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods; callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Bearer Token authentication MUST be configured at client initialization and automatically included in all requests.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Implementation MUST follow the OpenAPI 3.0 specification defined in vrm.yaml for User/Tag and User/Repository APIs only.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Client Initialization with Bearer Token (Priority: P1)

As an application developer, I want to access the VRM client through the project-scoped pattern (cloudsdk.Client.Project(projectID).VRM()) so that all subsequent API calls automatically include authentication and project context without requiring me to manage tokens or project IDs in each request.

**Why this priority**: Without proper client initialization, developers cannot authenticate with the VRM API. This is a prerequisite for all other operations. Using the project-scoped pattern maintains architectural consistency with the existing VPS service.

**Independent Test**: Can be tested by creating a client with a token and project ID, accessing VRM via the project-scoped pattern, and verifying the token is included in HTTP headers for subsequent requests.

**Acceptance Scenarios**:

1. **Given** a valid Bearer token and base URL, **When** creating a cloudsdk.Client and accessing VRM via client.Project(projectID).VRM(), **Then** the VRM client is created successfully and ready for API calls
2. **Given** a project-scoped VRM client, **When** making an API request (e.g., List Repositories), **Then** the request includes `Authorization: Bearer {token}` header and uses the correct project-scoped path
3. **Given** a project-scoped VRM client, **When** calling multiple API operations sequentially, **Then** all requests include the same token and project context without requiring re-specification
4. **Given** multiple projects, **When** creating different project-scoped VRM clients via client.Project(project1).VRM() and client.Project(project2).VRM(), **Then** each client operates within its respective project scope

---

### User Story 2 - Repository Management Operations (Priority: P1)

As a developer managing virtual repositories, I want to list, create, retrieve, update, and delete repositories within a project so that I can organize and manage image registries for my application deployments.

**Why this priority**: Repository operations are the fundamental resource management operations required for any VRM integration.

**Independent Test**: Can be tested through CRUD operations with validation that responses match the API specification.

**Acceptance Scenarios**:

1. **Given** a valid project ID and authentication token, **When** calling ListRepositories(), **Then** a list of repositories with their metadata (id, name, namespace, operatingSystem, creator, etc.) is returned
2. **Given** a project ID and repository creation parameters (name, operatingSystem), **When** calling CreateRepository(), **Then** a new repository is created and returned with complete metadata including timestamps
3. **Given** a project ID and repository ID, **When** calling GetRepository(), **Then** the repository details including nested project and creator information are returned
4. **Given** a project ID, repository ID, and update parameters, **When** calling UpdateRepository(), **Then** the repository is updated and the updated instance is returned
5. **Given** a project ID and repository ID, **When** calling DeleteRepository(), **Then** the repository is removed without returning a response body

---

### User Story 3 - Tag Management Operations (Priority: P1)

As a developer managing image versions, I want to list, create, retrieve, update, and delete tags within a repository so that I can organize different versions of images.

**Why this priority**: Tag operations enable version management of images, which is essential for container lifecycle management.

**Independent Test**: Can be tested through CRUD operations with validation that responses match the API specification.

**Acceptance Scenarios**:

1. **Given** a project ID and authentication token, **When** calling ListTags(), **Then** a list of all accessible tags is returned with their metadata (id, name, type, size, status, etc.)
2. **Given** a project ID, repository ID, and filter criteria, **When** calling ListRepositoryTags(), **Then** only tags belonging to that repository are returned
3. **Given** a project ID, repository ID, and tag creation parameters (name, type, diskFormat, containerFormat), **When** calling CreateTag(), **Then** a new tag is created and returned with complete metadata
4. **Given** a project ID and tag ID, **When** calling GetTag(), **Then** the tag details including nested repository and creator information are returned
5. **Given** a project ID, tag ID, and update parameters, **When** calling UpdateTag(), **Then** the tag is updated and the updated instance is returned
6. **Given** a project ID and tag ID, **When** calling DeleteTag(), **Then** the tag is removed without returning a response body

---

### User Story 4 - Query Filtering and Pagination (Priority: P2)

As a developer retrieving large lists of repositories or tags, I want to filter results and paginate through them so that I can efficiently retrieve only the data I need and handle large datasets.

**Why this priority**: Pagination and filtering are important for performance and usability when dealing with potentially large numbers of repositories and tags.

**Independent Test**: Can be tested by calling list operations with various filter and pagination parameters and validating the response subset.

**Acceptance Scenarios**:

1. **Given** ListRepositories called with limit parameter, **When** the result exceeds the limit, **Then** only the specified number of items are returned
2. **Given** ListRepositories called with offset parameter, **When** offset > 0, **Then** the first offset items are skipped in the response
3. **Given** ListRepositories called with where filter (e.g., operatingSystem=linux), **When** the filter is applied, **Then** only matching repositories are returned
4. **Given** ListTags called with various filter parameters, **When** filters are provided, **Then** only matching tags are returned based on supported filter fields

---

### User Story 5 - X-Namespace Header Support (Priority: P2)

As a developer working with both public and private namespaces, I want the ability to specify which namespace (public/private) a request targets so that I can manage resources in different visibility scopes.

**Why this priority**: Namespace support is important for multi-tenant scenarios and resource isolation.

**Independent Test**: Can be tested by verifying that requests include the X-Namespace header when provided and that responses are scoped to the correct namespace.

**Acceptance Scenarios**:

1. **Given** a namespace value in the request context, **When** making an API call, **Then** the X-Namespace header is included in the HTTP request
2. **Given** operations performed in different namespaces, **When** retrieving resources, **Then** resources are correctly scoped to their respective namespaces

---

### Edge Cases

- What happens when an invalid project ID is provided to Repository operations? (API returns 404 error; SDK propagates as appropriate error)
- How does the system handle repositories or tags that no longer exist when accessing them? (API returns 404 error; SDK propagates as appropriate error)
- What happens when creating a repository with a name that already exists within the same project? (API returns 400 error; SDK propagates as appropriate error)
- What happens when special characters or Unicode characters are included in repository/tag names? (Names are properly URL-encoded in requests; API response parsing handles Unicode correctly)
- How should the SDK handle malformed JSON responses from the API? (SDK returns a deserialization error with context)
- What happens if the Bearer token expires during a session? (SDK does not automatically refresh; callers must handle 401 responses and re-authenticate)
- How does pagination work when total count exceeds the maximum limit? (Standard offset-based pagination; offset + limit determines the window)
- What happens if required fields are missing from API responses? (SDK returns a deserialization error if required fields are absent)

## Requirements *(mandatory)*

### Functional Requirements

#### Client Initialization & Authentication
- **FR-001**: VRM Client MUST be accessed via the project-scoped pattern: cloudsdk.Client.Project(projectID).VRM() following the existing VPS service architecture
- **FR-002**: All API requests MUST automatically include the `Authorization: Bearer {token}` header inherited from the parent cloudsdk.Client
- **FR-003**: The VRM client MUST use the base URL from the parent cloudsdk.Client with /vrm path appended (e.g., baseURL + "/vrm")
- **FR-004**: The VRM client MUST accept optional X-Namespace header for namespace-scoped operations

#### Repository Operations
- **FR-005**: ListRepositories() MUST accept project ID, optional limit, optional offset, optional filters (where parameter), and return a list matching ListRepositoriesOutput schema
- **FR-006**: CreateRepository() MUST accept project ID and required input parameters (name, operatingSystem) and return a repository matching CreateRepositoryOutput schema
- **FR-007**: GetRepository() MUST accept project ID and repository ID and return a single repository matching GetRepositoryOutput schema
- **FR-008**: UpdateRepository() MUST accept project ID, repository ID, and update parameters (e.g., description) and return the updated repository matching UpdateRepositoryOutput schema
- **FR-009**: DeleteRepository() MUST accept project ID and repository ID and remove the repository (returns 204 No Content response)
- **FR-010**: Repository operations MUST include optional X-Namespace header when specified by the caller

#### Tag Operations
- **FR-011**: ListTags() MUST accept project ID, optional limit, optional offset, optional filters, and return a list matching ListTagsOutput schema with all accessible tags
- **FR-012**: ListRepositoryTags() MUST accept project ID and repository ID, optional limit, optional offset, optional filters (status, type), and return a list matching ListTagsOutput schema for tags within that repository
- **FR-013**: CreateTag() MUST accept project ID, repository ID, and required parameters (name, type, diskFormat, containerFormat) and return a tag matching CreateTagOutput schema
- **FR-014**: GetTag() MUST accept project ID and tag ID and return a single tag matching GetTagOutput schema
- **FR-015**: UpdateTag() MUST accept project ID, tag ID, and update parameters (e.g., name) and return the updated tag matching UpdateTagOutput schema
- **FR-016**: DeleteTag() MUST accept project ID and tag ID and remove the tag (returns 204 No Content response)
- **FR-017**: Tag operations MUST include optional X-Namespace header when specified by the caller

#### Data Model Requirements
- **FR-018**: Repository model MUST include all fields from the API response: id, name, namespace, operatingSystem, description, tags array, count, creator (nested IDName), project (nested IDName), createdAt, updatedAt
- **FR-019**: Tag model MUST include all fields from the API response: id, name, repositoryID, type, size, status, extra (map), createdAt, updatedAt, repository (nested Repository)
- **FR-020**: All timestamp fields (createdAt, updatedAt) MUST be represented as time.Time in Go or ISO 8601 string format as received from API
- **FR-021**: Nested objects (creator, project, repository) MUST be properly structured with id and name fields (IDName pattern)

#### Error Handling
- **FR-022**: HTTP error responses (4xx, 5xx) MUST be wrapped in SDK-specific error types with appropriate status codes and error messages
- **FR-023**: Malformed API responses that cannot be deserialized MUST result in a clear deserialization error
- **FR-024**: Missing required context (token, project ID) MUST result in a validation error before making HTTP requests

#### Performance & Reliability
- **FR-029**: The SDK client MUST use a default HTTP timeout of 30 seconds for all API requests
- **FR-030**: The SDK MUST support configurable timeout override via WithTimeout() client option following the existing VPS client pattern
- **FR-031**: The SDK MUST use the shared internal/http.Client with built-in retry logic for transient failures (503, 429, network errors)
- **FR-032**: The VRM client MUST be safe for concurrent use by multiple goroutines without external synchronization (thread-safe by design using the underlying thread-safe http.Client)

#### API Compliance
- **FR-025**: All Repository API endpoints MUST match the User/Repository paths and methods from vrm.yaml specification
- **FR-026**: All Tag API endpoints MUST match the User/Tag paths and methods from vrm.yaml specification
- **FR-027**: All Admin/* endpoints from vrm.yaml MUST be explicitly excluded from this implementation
- **FR-028**: All non-Repository and non-Tag endpoints (MemberAcl, ProjectAcl, Image, Export, Snapshot, etc.) MUST be explicitly excluded from this implementation

#### Architecture Pattern Compliance
- **FR-033**: The VRM module MUST follow the exact same architectural pattern as the VPS module: modules/vrm/client.go as main client, with sub-packages modules/vrm/repositories and modules/vrm/tags
- **FR-034**: The cloudsdk.ProjectClient MUST be extended with a VRM() method that returns a project-scoped VRM client, analogous to the existing VPS() method
- **FR-035**: The VRM base URL MUST be constructed as parent baseURL + "/vrm" (e.g., "https://api.example.com/vrm")

### SDK Contract Requirements (Go)
```go
// Architecture follows existing VPS pattern:
// cloudsdk.Client.Project(projectID).VRM() returns project-scoped VRM client

// VRM Client provides access to VRM operations for a specific project
type Client struct {
    // Private fields: baseClient *internalhttp.Client, projectID string, basePath string
}

// NewClient creates a new project-scoped VRM client (internal use)
// Called via cloudsdk.Client.Project(projectID).VRM()
func NewClient(baseURL, token, projectID string, httpClient *http.Client, logger types.Logger) *Client

// ProjectID returns the project ID this client is bound to
func (c *Client) ProjectID() string

// Repositories returns the repository operations sub-client
func (c *Client) Repositories() *repositories.Client

// Tags returns the tag operations sub-client
func (c *Client) Tags() *tags.Client

// Data structures
type Repository struct {
    ID              string
    Name            string
    Namespace       string
    OperatingSystem string
    Description     string
    Tags            []Tag
    Count           int
    Creator         IDName
    Project         IDName
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

type Tag struct {
    ID           string
    Name         string
    RepositoryID string
    Type         string
    Size         int64
    Status       string
    Extra        map[string]interface{}
    CreatedAt    time.Time
    UpdatedAt    time.Time
    Repository   Repository
}

type IDName struct {
    ID   string
    Name string
}

// Repository operations (modules/vrm/repositories package)
type repositories.Client struct {
    // Methods:
    // List(ctx, limit, offset, filters, namespace) (*ListRepositoriesOutput, error)
    // Create(ctx, input, namespace) (*Repository, error)
    // Get(ctx, repositoryID, namespace) (*Repository, error)
    // Update(ctx, repositoryID, input, namespace) (*Repository, error)
    // Delete(ctx, repositoryID, namespace) error
}

// Tag operations (modules/vrm/tags package)
type tags.Client struct {
    // Methods:
    // List(ctx, limit, offset, filters, namespace) (*ListTagsOutput, error)
    // ListByRepository(ctx, repositoryID, limit, offset, filters, namespace) (*ListTagsOutput, error)
    // Create(ctx, repositoryID, input, namespace) (*Tag, error)
    // Get(ctx, tagID, namespace) (*Tag, error)
    // Update(ctx, tagID, input, namespace) (*Tag, error)
    // Delete(ctx, tagID, namespace) error
}

// Usage example:
// client, _ := cloudsdk.New("https://api.example.com", "token")
// vrmClient := client.Project("project-123").VRM()
// repos, _ := vrmClient.Repositories().List(ctx, 10, 0, nil, "public")
// tag, _ := vrmClient.Tags().Get(ctx, "tag-id", "public")
```

## Key Entities *(include if feature involves data)*

- **Client**: VRM API client managing HTTP communication, authentication, and request/response handling
- **Repository**: Virtual repository resource with metadata, creator, project affiliation, and tag collection
- **Tag**: Image tag/version resource with type, format specifications, and repository relationship
- **IDName**: Common structure for nested resource references containing id and name
- **ListRepositoriesOutput**: Paginated list response containing repositories array and total count
- **ListTagsOutput**: Paginated list response containing tags array and total count
- **CreateRepositoryInput**: Input parameters for creating a repository (name, operatingSystem, optional description)
- **CreateTagInput**: Input parameters for creating a tag (name, type, diskFormat, containerFormat)
- **UpdateRepositoryInput**: Input parameters for updating a repository (optional fields like description)
- **UpdateTagInput**: Input parameters for updating a tag (optional fields like name)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: VRM Client successfully initializes with Bearer token and includes authentication in all requests
- **SC-002**: All 5 Repository CRUD operations (List, Create, Get, Update, Delete) work correctly and return data matching the API specification
- **SC-003**: All 6 Tag operations (List project tags, List repository tags, Create, Get, Update, Delete) work correctly and return data matching the API specification
- **SC-004**: Client properly handles pagination with limit and offset parameters for List operations
- **SC-005**: Client properly handles optional where filters for repository and tag queries
- **SC-006**: Client includes optional X-Namespace header in requests when provided
- **SC-007**: HTTP error responses are wrapped in appropriate SDK error types with status codes and messages
- **SC-008**: All data model fields are correctly deserialized from API JSON responses
- **SC-009**: Nested objects (creator, project, repository) are properly resolved to IDName structures
- **SC-010**: All optional fields are correctly handled when absent from API responses
- **SC-011**: Repository paths exclude all non-User/Repository endpoints (Admin/*, MemberAcl, ProjectAcl, Image, Export, Snapshot, etc.)
- **SC-012**: Tag paths exclude all non-User/Tag endpoints (Admin/*, MemberAcl, ProjectAcl, Image, Export, Snapshot, etc.)
- **SC-013**: Unit tests covering all operations pass with 100% success rate
- **SC-014**: Contract tests validating API specification compliance pass
