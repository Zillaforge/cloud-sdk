# Feature Specification: IAM API Client

**Feature Branch**: `009-iam-api`  
**Created**: 2025-11-12  
**Status**: Draft  
**Input**: User description: "Implement IAM RESTful API with GET /user, GET /projects, and GET /project/<project_id> endpoints using Bearer Token authentication"

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods;
  callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Retrieve Current User Information (Priority: P1)

A developer needs to retrieve information about the currently authenticated user to display their profile or verify their identity in their application.

**Why this priority**: This is the most fundamental operation - knowing who the current user is forms the basis for all other IAM operations. Without this, no user-centric functionality can be built.

**Independent Test**: Can be fully tested by initializing an IAM client with a valid Bearer token and calling GetUser(). Should return user details including account, userId, displayName, namespace, and timestamps.

**Acceptance Scenarios**:

1. **Given** a valid Bearer token, **When** GetUser() is called, **Then** the system returns complete user information including account, userId, displayName, description, extra metadata, namespace, creation timestamp, update timestamp, last login timestamp, and permission details.
2. **Given** an expired Bearer token, **When** GetUser() is called, **Then** the system returns an authentication error with status 403.
3. **Given** an invalid Bearer token, **When** GetUser() is called, **Then** the system returns an authentication error.
4. **Given** a valid token with context timeout, **When** GetUser() is called, **Then** the system respects the context deadline and returns a timeout error.

---

### User Story 2 - List User's Projects (Priority: P2)

A developer needs to retrieve all projects that the authenticated user belongs to, enabling them to display project selection interfaces or validate project access.

**Why this priority**: After identifying the user, listing their projects is the next essential step for multi-tenant applications. This enables users to see what resources they can access.

**Independent Test**: Can be fully tested by calling ListProjects() with optional pagination parameters (offset, limit, order). Should return a list of projects with associated permissions and membership details.

**Acceptance Scenarios**:

1. **Given** a user with multiple project memberships, **When** ListProjects() is called without pagination, **Then** the system returns all projects with default pagination (offset=0, limit=20).
2. **Given** a user with multiple project memberships, **When** ListProjects() is called with custom offset and limit, **Then** the system returns the specified page of projects.
3. **Given** a user with no project memberships, **When** ListProjects() is called, **Then** the system returns an empty list with total count of 0.
4. **Given** a valid request, **When** ListProjects() returns results, **Then** each project includes projectId, displayName, description, extra metadata, namespace, timestamps, permissions (global and user), and frozen status.

---

### User Story 3 - Retrieve Specific Project Details (Priority: P2)

A developer needs to fetch detailed information about a specific project using its project ID, to display project configuration or validate access before performing project-scoped operations.

**Why this priority**: Once users select a project, applications need detailed information about that specific project. This is essential for project-scoped operations in multi-tenant systems.

**Independent Test**: Can be fully tested by calling GetProject(projectId) with a valid project ID that the user has access to. Should return complete project details including permissions.

**Acceptance Scenarios**:

1. **Given** a valid project ID that the user belongs to, **When** GetProject(projectId) is called, **Then** the system returns complete project details including projectID, displayName, description, extra metadata, namespace, timestamps, and permission information.
2. **Given** a valid project ID that the user does NOT belong to, **When** GetProject(projectId) is called, **Then** the system returns a not found or forbidden error.
3. **Given** an invalid or non-existent project ID, **When** GetProject(projectId) is called, **Then** the system returns a 400 Bad Request error.
4. **Given** a context with timeout, **When** GetProject(projectId) is called, **Then** the system respects the context deadline.

---

### Edge Cases

- What happens when the Bearer token expires during a long-running operation?
- How does the system handle rate limiting or throttling from the IAM service?
- What happens when the IAM service returns malformed JSON or unexpected fields? → Malformed JSON returns parse error; unexpected fields are ignored for forward compatibility
- How does the client handle network interruptions or connection timeouts?
- What happens when pagination parameters exceed allowed ranges (e.g., limit > 100)? → Server validates and returns 400 Bad Request with error details
- How does the system handle projects with null or missing fields in the response?
- What happens when extra metadata contains deeply nested or large JSON objects?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a Client struct that is initialized with a Bearer token for authentication
- **FR-002**: Client MUST accept standard http.Client configuration options for timeout, retry, and transport customization
- **FR-002a**: Client MUST implement automatic retry with exponential backoff for transient failures (429, 502, 503, 504 status codes) on GET requests, with default 3 retry attempts configurable during initialization
- **FR-003**: System MUST provide a GetUser method that retrieves current authenticated user information
- **FR-004**: GetUser MUST return user details including account, userId, displayName, description, extra metadata, namespace, createdAt, updatedAt, lastLoginAt, and permission information
- **FR-005**: System MUST provide a ListProjects method that retrieves all projects for the authenticated user
- **FR-006**: ListProjects MUST support optional pagination parameters (offset, limit, order)
- **FR-007**: ListProjects MUST return project list with total count, where each project includes detailed information and permission context
- **FR-008**: System MUST provide a GetProject method that retrieves a specific project by project ID
- **FR-009**: GetProject MUST return complete project details including metadata and permission information
- **FR-010**: All API methods MUST accept context.Context as the first parameter for cancellation and timeout control
- **FR-011**: All API methods MUST return strongly-typed Go structs matching the Swagger/OpenAPI specifications
- **FR-011a**: Client MUST ignore unknown/unexpected fields in API responses, parsing only known fields to support forward compatibility with API evolution
- **FR-012**: System MUST return descriptive errors that include HTTP status codes and service error messages
- **FR-012a**: Client MUST NOT include built-in logging to maintain consistency with other service clients; error information is conveyed through return values only
- **FR-013**: Client MUST send Bearer token in the Authorization header for all requests
- **FR-014**: Client MUST use base URL "http://127.0.0.1:8084/iam/api/v1" with support for configuration override
- **FR-014a**: Client MUST support both HTTP and HTTPS protocols without enforcement or warnings, allowing deployment flexibility
- **FR-015**: Client MUST properly encode query parameters for pagination (offset, limit, order)
- **FR-015a**: Client MUST pass pagination parameters to API without client-side validation; server validates and returns 400 error for invalid values
- **FR-016**: System MUST handle HTTP status codes: 200 (success), 400 (bad request), 403 (forbidden), 500 (internal server error)

### SDK Contract Requirements (Go)

- Public methods follow the Resources().Verb() pattern: `client.IAM().User().Get(ctx)`, `client.IAM().Projects().List(ctx, opts)`, `client.IAM().Projects().Get(ctx, projectID)`.
- Two resources: User (singular) with Get() verb, Projects (plural) with List() and Get() verbs.
- Responses are strongly-typed structs matching Swagger models (User, Project, ProjectMembership).
- Errors wrap context (status code, service error fields) without exposing raw HTTP.
- Authentication (Bearer token), retries, and timeouts are centralized in the Client.
- Pagination options are passed via functional options pattern (WithOffset, WithLimit, WithOrder).

### Key Entities *(include if feature involves data)*

- **User**: Represents an authenticated user with account, userId, displayName, description, extra metadata, namespace, timestamps (createdAt, updatedAt, lastLoginAt), email, frozen status, and MFA flag
- **Project**: Represents a project/tenant with projectId, displayName, description, extra metadata, namespace, frozen status, and timestamps
- **ProjectMembership**: Wraps a Project with additional context including globalPermissionId, userPermissionId, permission objects (global and user), tenantRole, and frozen status
- **Permission**: Represents access permissions with id and label fields
- **ListProjectsOptions**: Configuration for pagination including offset (default: 0), limit (default: 20, max: 100), and order

## Clarifications

### Session 2025-11-12

- Q: Should the IAM client implement automatic retry logic for transient failures? → A: Automatic retries with exponential backoff for 5xx errors and network timeouts (default: 3 attempts), configurable during client initialization
- Q: Should the IAM client enforce HTTPS for production environments? → A: No enforcement - allow both HTTP and HTTPS without warnings
- Q: What level of logging should the IAM client provide? → A: No built-in logging - maintain consistency with other service clients
- Q: How should the client handle invalid pagination parameters? → A: Pass through to API without validation - let server return 400 error
- Q: How should the client handle unexpected or additional fields in API responses? → A: Ignore unknown fields - parse known fields only, allowing API evolution

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can initialize an IAM client and retrieve current user information in under 5 lines of code
- **SC-002**: API calls complete within 2 seconds under normal network conditions (excluding network latency)
- **SC-003**: Client handles authentication errors gracefully with clear error messages indicating token expiry or invalidity
- **SC-004**: Pagination works correctly for project lists containing 1 to 100+ projects
- **SC-005**: All API methods respect context cancellation and timeout within 100ms
- **SC-006**: 100% of Swagger schema fields are represented in Go structs with appropriate types
- **SC-007**: Error responses include sufficient context for debugging without exposing sensitive internal details
- **SC-008**: Client can be tested in isolation using mock HTTP responses matching Swagger contracts

## Assumptions

- Base URL is consistent across environments or can be overridden during client initialization
- Both HTTP and HTTPS protocols are supported; transport security is the responsibility of the deployment environment
- No built-in logging mechanism for consistency with other service clients in the SDK; users implement their own observability as needed
- API responses may contain additional fields beyond those defined in current Swagger spec; client ignores unknown fields to maintain forward compatibility
- Bearer tokens are obtained through an external authentication flow (login is out of scope)
- Token refresh/renewal is handled externally or by a separate mechanism
- The IAM service follows the OpenAPI 3.0.0 specification provided in iam.yaml
- HTTP communication uses standard REST conventions (GET for retrieval, proper status codes)
- JSON is the only supported content type for request/response bodies
- Default pagination limits (20 items per page, max 100) are sufficient for most use cases
- Network retry logic uses exponential backoff strategy (100ms initial, 5s max interval, 2.0 multiplier) with jitter, retrying up to 3 times for 5xx errors and network timeouts on idempotent GET operations
- All timestamps are in RFC3339 format or similar standard time representation
- Extra metadata fields are unstructured JSON objects (map[string]interface{} in Go)

## Dependencies

- Go standard library: net/http, context, encoding/json, time
- IAM service API endpoint accessible at configured base URL
- Valid Bearer token for authentication

## Out of Scope

- User login/authentication flows (token generation)
- Token refresh or renewal mechanisms
- Admin-scoped APIs (only user-scoped endpoints)
- Project mutation operations (create, update, delete)
- User profile mutation operations
- Membership management
- Credential management (access/secret keys)
- MFA operations
- CLI secrets
- Personal access tokens
- SAML operations
- Password management
- Public key management
- Permission management beyond read-only access
- Frozen/usage time management
- Any endpoints not explicitly listed (GET /user, GET /projects, GET /project/{project-id})
