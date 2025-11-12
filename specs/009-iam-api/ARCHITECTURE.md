# IAM Client Architecture: Resources().Verb() Pattern

**Feature**: 009-iam-api  
**Date**: 2025-11-12  
**Pattern**: Resources().Verb()

## Overview

The IAM client implements a **Resources().Verb() pattern** where operations are organized by resource type (User, Projects), and each resource exposes specific verbs (Get, List).

## Architecture Pattern

### Two Resources

1. **User Resource** (singular) - Represents the current authenticated user
   - Single verb: `Get()`
   - Returns information about the authenticated user

2. **Projects Resource** (plural) - Represents the collection of projects
   - Two verbs: `List()` and `Get(projectID)`
   - Manages project listing and retrieval

### Usage Pattern

```go
// Initialize
client := cloudsdk.NewClient(baseURL, token)
iamClient := client.IAM()

// User resource - singular, one verb
user, err := iamClient.User().Get(ctx)

// Projects resource - plural, multiple verbs
projects, err := iamClient.Projects().List(ctx, iam.WithLimit(50))
project, err := iamClient.Projects().Get(ctx, projectID)
```

## Implementation Structure

### File Organization

```
modules/iam/
├── client.go       # Base IAM client
│   └── User()      # Returns *UserResource
│   └── Projects()  # Returns *ProjectsResource
├── user.go         # UserResource implementation
│   └── Get(ctx)    # GET /user
└── projects.go     # ProjectsResource implementation
    ├── List(ctx, opts...)  # GET /projects
    └── Get(ctx, id)        # GET /project/{id}
```

### Type Definitions

```go
// modules/iam/client.go
type Client struct {
    baseClient *http.Client
    baseURL    string
    token      string
}

func (c *Client) User() *UserResource {
    return &UserResource{client: c}
}

func (c *Client) Projects() *ProjectsResource {
    return &ProjectsResource{client: c}
}

// modules/iam/user.go
type UserResource struct {
    client *Client
}

func (r *UserResource) Get(ctx context.Context) (*User, error) {
    // GET /user
}

// modules/iam/projects.go
type ProjectsResource struct {
    client *Client
}

func (r *ProjectsResource) List(ctx context.Context, opts ...ListProjectsOption) (*ListProjectsResponse, error) {
    // GET /projects
}

func (r *ProjectsResource) Get(ctx context.Context, projectID string) (*GetProjectResponse, error) {
    // GET /project/{project-id}
}
```

## Rationale

### Why Resources().Verb()?

1. **RESTful Mapping**: Directly maps to REST resource model
   - `/user` → `User()` resource
   - `/projects` → `Projects()` resource

2. **Clear Separation**: Different resources are explicitly separated
   - User operations in `user.go`
   - Project operations in `projects.go`

3. **Self-Documenting**: Code reads like natural language
   - `client.IAM().User().Get()` - "get the IAM user"
   - `client.IAM().Projects().List()` - "list the IAM projects"

4. **Extensibility**: Easy to add verbs to resources without breaking changes
   - Add `User().Update()` in future
   - Add `Projects().Create()` for admin API

5. **Testability**: Resources can be mocked independently
   ```go
   type MockUserResource struct{}
   func (m *MockUserResource) Get(ctx) (*User, error) { ... }
   ```

### Comparison with Other Patterns

#### ❌ Flat Client Pattern (Rejected)
```go
// Mixes resource concerns
client.IAM().GetUser(ctx)
client.IAM().ListProjects(ctx)
client.IAM().GetProject(ctx, id)
```
**Problem**: No clear resource boundaries, harder to extend

#### ❌ Separate Clients Pattern (Rejected)
```go
// Requires multiple client initializations
userClient := client.User()
projectClient := client.Projects()
```
**Problem**: Inconsistent with existing SDK patterns

#### ✅ Resources().Verb() Pattern (Chosen)
```go
// Clear resource separation, intuitive API
client.IAM().User().Get(ctx)
client.IAM().Projects().List(ctx)
client.IAM().Projects().Get(ctx, id)
```
**Benefits**: Resource-oriented, extensible, self-documenting

## API Surface

### Complete API

```go
// Root client
client.IAM() *iam.Client

// User resource (1 verb)
client.IAM().User().Get(ctx context.Context) (*User, error)

// Projects resource (2 verbs)
client.IAM().Projects().List(ctx context.Context, opts ...ListProjectsOption) (*ListProjectsResponse, error)
client.IAM().Projects().Get(ctx context.Context, projectID string) (*GetProjectResponse, error)

// Pagination options for List
iam.WithOffset(offset int) ListProjectsOption
iam.WithLimit(limit int) ListProjectsOption
iam.WithOrder(order string) ListProjectsOption
```

### Examples

#### Get Current User
```go
user, err := client.IAM().User().Get(ctx)
if err != nil {
    return fmt.Errorf("failed to get user: %w", err)
}
fmt.Printf("Logged in as: %s\n", user.DisplayName)
```

#### List Projects with Pagination
```go
resp, err := client.IAM().Projects().List(ctx,
    iam.WithOffset(20),
    iam.WithLimit(10),
)
if err != nil {
    return err
}
for _, membership := range resp.Projects {
    fmt.Printf("- %s\n", membership.Project.DisplayName)
}
```

#### Get Specific Project
```go
project, err := client.IAM().Projects().Get(ctx, projectID)
if err != nil {
    return fmt.Errorf("project not found: %w", err)
}
fmt.Printf("Project: %s\n", project.DisplayName)
```

## Differences from VPS/VRM

### IAM (Non-Project-Scoped)
```go
client.IAM().User().Get(ctx)
client.IAM().Projects().List(ctx)
```
- Direct access via `client.IAM()`
- No project context required
- Cross-project operations

### VPS/VRM (Project-Scoped)
```go
client.Project(projectID).VPS().Flavors().List(ctx)
client.Project(projectID).VRM().Repositories().List(ctx)
```
- Requires project context
- Operations scoped to specific project
- Project-level resources

## Testing Strategy

### Unit Tests

```go
// Test User resource
func TestUserResource_Get(t *testing.T) {
    // Mock HTTP client
    // Call User().Get()
    // Verify request and response
}

// Test Projects resource
func TestProjectsResource_List(t *testing.T) {
    // Test pagination options
    // Test response parsing
}

func TestProjectsResource_Get(t *testing.T) {
    // Test project ID parameter
    // Test error cases
}
```

### Contract Tests

```go
func TestIAMContracts(t *testing.T) {
    // Validate GET /user response schema
    // Validate GET /projects response schema
    // Validate GET /project/{id} response schema
}
```

## Future Extensions

### Potential User Resource Verbs
- `User().Update(ctx, req)` - Update user profile
- `User().UpdatePassword(ctx, req)` - Change password
- `User().GetPublicKey(ctx)` - Get SSH public key

### Potential Projects Resource Verbs
- `Projects().Create(ctx, req)` - Create new project (admin)
- `Projects().Update(ctx, id, req)` - Update project (admin)
- `Projects().Delete(ctx, id)` - Delete project (admin)

All future additions maintain the Resources().Verb() pattern for consistency.

## Summary

The **Resources().Verb() pattern** provides:
- ✅ Clear resource boundaries (User, Projects)
- ✅ Self-documenting API surface
- ✅ Natural mapping to REST endpoints
- ✅ Easy extensibility per resource
- ✅ Independent resource testing
- ✅ Consistent with RESTful design principles

This pattern makes the IAM client intuitive to use and maintain while following Go best practices and SDK constitution requirements.
