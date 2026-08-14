# Implementation Plan: IAM API Client

**Branch**: `009-iam-api` | **Date**: 2025-11-12 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/009-iam-api/spec.md`

## Summary

Implement a Go SDK client for IAM (Identity and Access Management) API providing three non-project-scoped endpoints: GET /user (retrieve current user information), GET /projects (list user's accessible projects with pagination), and GET /project/{id} (get specific project details). The client uses Bearer token authentication, supports automatic retry with exponential backoff for transient failures, and follows TDD methodology with 80%+ test coverage. Implementation uses Go standard library only (net/http, encoding/json, context) with no external dependencies, following the existing SDK patterns established by VPS and VRM modules.

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: Go standard library only (net/http, encoding/json, context, time, fmt)  
**Storage**: N/A (API client, no persistence layer)  
**Testing**: Go testing package (stdlib), httptest for mocking  
**Target Platform**: Cross-platform (Linux, macOS, Windows) - Go compiled binaries  
**Project Type**: Library/SDK module (Go package)  
**Performance Goals**: 
- API calls complete within 2 seconds (excluding network latency)
- Context cancellation response within 100ms
- Support pagination up to 100 items per request

**Constraints**:
- Zero external dependencies (stdlib only per constitution)
- No built-in logging (consistency with other service clients)
- No client-side validation of pagination parameters (server validates)
- Ignore unknown JSON fields for forward compatibility
- HTTP and HTTPS support without enforcement

**Scale/Scope**:
- 3 API endpoints (GetUser, ListProjects, GetProject)
- 6 model structs (User, Project, ProjectMembership, Permission, UsageTime, responses)
- Target 80%+ test coverage
- Support pagination for projects (1-100+ items)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

✅ **Status: PASSED** (Initial and Post-Design)

The feature plan satisfies all applicable Cloud SDK Constitution principles:

- ✅ **TDD Mandatory**: 
  - Tests will be written first for all public methods (GetUser, ListProjects, GetProject)
  - Unit tests for model parsing and JSON unmarshaling
  - Contract tests validating responses against iam.yaml Swagger spec
  - Integration tests for real API interactions
  - Target 80%+ code coverage enforced via `make check`
  
- ✅ **Public API Shape**: 
  - Idiomatic Go package: `modules/iam`
  - Client accessed via `client.IAM()` (non-project-scoped)
  - **Resources().Verb() pattern**: Two resources (User, Projects)
  - User resource: `client.IAM().User().Get(ctx)` - retrieves current user
  - Projects resource: `client.IAM().Projects().List(ctx, opts)` and `.Get(ctx, projectID)` - manages project collection
  - **List signature**: `List(ctx, *ListProjectsOptions)` - single options pointer (consistent with VPS/VRM)
  - All methods accept `context.Context` as first parameter
  - Return typed structs: `*User`, `*ListProjectsResponse`, `*GetProjectResponse`
  - Errors wrapped with context: `fmt.Errorf("failed to get user: %w", err)`
  - No raw HTTP exposure to callers
  - Uses `internalhttp.Client` for HTTP operations (consistent with VPS/VRM)
  
- ✅ **Dependencies**: 
  - **ZERO external dependencies** - uses only Go standard library
  - `net/http` - HTTP client
  - `encoding/json` - JSON parsing
  - `context` - context propagation
  - `time`, `fmt` - utilities
  - Reuses existing internal packages: `internal/http`, `internal/types`, `internal/backoff`
  
- ✅ **Versioning**: 
  - Change type: **MINOR** (adds new module, no breaking changes)
  - No migration needed (new functionality)
  - Follows semantic versioning (0.x.y -> 0.(x+1).0)
  
- ✅ **Observability**: 
  - Logger interface passed via `internal/http.Client` (consistent with VPS/VRM)
  - IAM client receives `internalhttp.Client` which includes logger
  - Logger used for internal operations, not exposed to IAM client API
  - Users configure logger when creating root SDK client
  
- ✅ **Security**: 
  - Bearer token passed via constructor, never logged
  - Token sent in Authorization header only
  - Support both HTTP and HTTPS (no enforcement per clarification)
  - TLS verification enabled by default when HTTPS used (stdlib behavior)
  - No global state - configuration through explicit constructors
  - Extra metadata (map[string]interface{}) may contain sensitive data; not logged

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
# IAM Module Implementation (following VPS/VRM pattern)
models/iam/                           # IAM data models
├── common/
│   ├── common.go                     # Shared types (Permission, timestamps, metadata)
│   └── common_test.go                # Unit tests for common types
├── users/
│   ├── models.go                     # User, GetUserResponse
│   └── models_test.go                # User model tests (JSON parsing, validation)
└── projects/
    ├── models.go                     # Project, ProjectMembership, ListProjectsResponse, GetProjectResponse
    ├── models_test.go                # Project model tests (JSON parsing, custom unmarshaler)
    └── options.go                    # ListProjectsOptions, functional option helpers

modules/iam/                          # IAM client implementation (non-project-scoped)
├── client.go                         # Base IAM client (NewClient receives *internalhttp.Client, Users(), Projects() factories)
├── client_test.go                    # Client initialization tests
├── users/
│   ├── client.go                     # User operations client with Get(ctx) method
│   └── client_test.go                # User client tests
├── projects/
│   ├── client.go                     # Projects operations client with List(ctx, *opts) and Get(ctx, id) methods
│   └── client_test.go                # Projects client tests (including pagination)
├── EXAMPLES.md                       # Usage examples
└── README.md                         # Module documentation

# Root-level integration
client.go                             # Add IAM() method to root Client
```

**Structure Decision - Resources().Verb() Pattern**: 

This uses a **resource-oriented architecture** following the established VPS/VRM patterns:

**Consistency with VPS/VRM**:
- `models/iam/` mirrors `models/vps/` and `models/vrm/` structure
- `modules/iam/` mirrors `modules/vps/` and `modules/vrm/` structure
- Both have subdirectories for each resource type (users/, projects/ like flavors/, servers/, etc.)
- Each resource has its own client.go with operations
- **Client construction**: Uses `internalhttp.Client` (same as VPS/VRM)
- **Logger handling**: Logger passed via `internalhttp.Client`, not exposed in IAM API
- **List method signature**: `List(ctx, *Options)` - single options pointer (same as VPS/VRM)
- **Tests location**: Tests in each resource directory (users/, projects/), not in root tests/

**Key Difference - Non-Project-Scoped**:
- **VPS/VRM**: `client.Project(projectID).VPS().Flavors().List(ctx, opts)` - project-scoped
- **IAM**: `client.IAM().User().Get(ctx)` - NOT project-scoped (no projectID required)
- IAM provides identity and access management operations that are global to the authenticated user
- No project context needed for IAM operations

**Resource Structure**:
- **Two resources**: User (singular) and Projects (plural)
- **User resource**: `client.IAM().User().Get(ctx)` - single verb for current user
- **Projects resource**: `client.IAM().Projects().List(ctx, opts)` and `.Get(ctx, projectID)` - two verbs for project collection
- `models/iam/users/` - User models with tests (like `models/vps/flavors/`)
- `models/iam/projects/` - Project models with tests (like `models/vps/servers/`)
- `models/iam/common/` - Shared types (like `models/vps/common/`)
- `modules/iam/client.go` - Base IAM client receives `*internalhttp.Client` (like `modules/vps/client.go`)
- `modules/iam/users/client.go` - User operations with tests (like `modules/vps/flavors/client.go`)
- `modules/iam/projects/client.go` - Project operations with tests (like `modules/vps/servers/client.go`)
- Root `client.go` provides entry point: `client.IAM()` returns `*iam.Client`

This pattern provides clear separation between different resource types and makes the API self-documenting, while maintaining **full consistency** with the existing SDK architecture.

## Complexity Tracking

> **No violations** - Constitution Check fully passed.

## Phase 0: Research Summary

All technical decisions documented in [research.md](research.md):

1. **Client Architecture**: Non-project-scoped pattern via `client.IAM()`
2. **Model Structure**: Swagger-aligned fields with Request/Response suffix
3. **List API Parsing**: Custom UnmarshalJSON for array extraction
4. **Error Handling**: Reuse existing `internal/types.SDKError`
5. **Test Strategy**: TDD with 80%+ coverage (unit, contract, integration)
6. **Pagination**: Functional options pattern for clean API
7. **Base URL**: IAM-specific URL construction from root baseURL

## Phase 1: Design Artifacts

Generated documentation:

- ✅ [data-model.md](data-model.md) - Complete type definitions with validation rules
- ✅ [contracts/api-contracts.md](contracts/api-contracts.md) - API contracts from Swagger spec
- ✅ [quickstart.md](quickstart.md) - Usage guide with code examples

Key design decisions:

### Data Model
- 6 core types: User, Project, ProjectMembership, Permission, UsageTime, plus operation responses
- All fields match Swagger spec exactly (projectId, displayName, etc.)
- Custom JSON unmarshaler for ListProjectsResponse to enable direct array access
- Forward-compatible: unknown fields ignored during parsing

### API Contracts
- GET /user: No parameters, returns User object
- GET /projects: Optional pagination (offset, limit, order), returns wrapped array + total
- GET /project/{id}: Project ID parameter, returns project with permissions

### Client Interface
```go
type Client struct {
    baseClient *http.Client
    baseURL    string
    token      string
}

func (c *Client) GetUser(ctx context.Context) (*User, error)
func (c *Client) ListProjects(ctx context.Context, opts ...ListProjectsOption) (*ListProjectsResponse, error)
func (c *Client) GetProject(ctx context.Context, projectID string) (*GetProjectResponse, error)
```

## Implementation Checklist

### Models (TDD - Tests First)
- [ ] `models/iam/common/common.go` - Shared types (Permission, TenantRole, timestamps helpers)
- [ ] `models/iam/common/common_test.go` - JSON parsing, field validation (tests in same directory)
- [ ] `models/iam/users/models.go` - User, GetUserResponse
- [ ] `models/iam/users/models_test.go` - User model tests (JSON parsing, unknown fields) (tests in same directory)
- [ ] `models/iam/projects/models.go` - Project, ProjectMembership, ListProjectsResponse, GetProjectResponse
- [ ] `models/iam/projects/models_test.go` - Project model tests + custom unmarshaler (tests in same directory)
- [ ] `models/iam/projects/options.go` - ListProjectsOptions struct (single options pointer, like VPS/VRM)

### Client (TDD - Tests First)
- [ ] `modules/iam/client.go` - IAM base client (NewClient receives *internalhttp.Client), Users() and Projects() factory methods
- [ ] `modules/iam/client_test.go` - Client initialization tests (tests in same directory)
- [ ] `modules/iam/users/client.go` - User resource client with Get(ctx) method, uses baseClient *internalhttp.Client
- [ ] `modules/iam/users/client_test.go` - User client method tests (tests in same directory)
- [ ] `modules/iam/projects/client.go` - Projects resource client with List(ctx, *opts) and Get(ctx, id) methods, uses baseClient *internalhttp.Client
- [ ] `modules/iam/projects/client_test.go` - Projects client tests including pagination (tests in same directory)
- [ ] `modules/iam/README.md` - API documentation
- [ ] `modules/iam/EXAMPLES.md` - Usage examples

### Root Integration
- [ ] Update `client.go` to add IAM() method that creates IAM client with internalhttp.Client
- [ ] Test IAM() method returns properly configured client

### Validation
- [ ] Run `make check` - verify tests pass and coverage >= 80%
- [ ] Manual API testing with real IAM service
- [ ] Review error messages for clarity
- [ ] Verify unknown fields are ignored (forward compatibility)

## Next Steps

After plan completion:
1. Run `/speckit.tasks` to generate task breakdown
2. Begin TDD implementation starting with model tests
3. Implement each component with tests first, code second
4. Validate with `make check` after each phase
5. Update EXAMPLES.md with real usage patterns
````
