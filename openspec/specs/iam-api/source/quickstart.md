# Quickstart: IAM API Client

**Feature**: 009-iam-api  
**Date**: 2025-11-12

## Overview

This quickstart guide demonstrates how to use the IAM API client to retrieve user information and list/access projects. The IAM client provides a simple, idiomatic Go interface for interacting with the IAM service.

## Prerequisites

- Go 1.21 or later
- Valid IAM Bearer token
- IAM service accessible at configured base URL

## Installation

```bash
go get github.com/Zillaforge/cloud-sdk
```

## Basic Usage

### 1. Initialize the Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/Zillaforge/cloud-sdk"
)

func main() {
    // Create root client with bearer token
    client, err := cloudsdk.New(
        "http://127.0.0.1:8084",  // Base URL
        "your-bearer-token-here",  // IAM token
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Get IAM client (non-project-scoped)
    iamClient := client.IAM()
    
    // Now use iamClient for IAM operations
}
```

### 2. Get Current User Information

```go
// Retrieve authenticated user details
ctx := context.Background()
user, err := iamClient.User().Get(ctx)
if err != nil {
    log.Fatalf("Failed to get user: %v", err)
}

fmt.Printf("User: %s (%s)\n", user.DisplayName, user.Account)
fmt.Printf("User ID: %s\n", user.UserID)
fmt.Printf("Namespace: %s\n", user.Namespace)
fmt.Printf("Last Login: %s\n", user.LastLoginAt)

// Access optional fields safely
if user.Permission != nil {
    fmt.Printf("Permission: %s\n", user.Permission.Label)
}

// Access extra metadata
if dept, ok := user.Extra["department"].(string); ok {
    fmt.Printf("Department: %s\n", dept)
}
```

### 3. List User's Projects

```go
// List all projects with default pagination (offset=0, limit=20)
projectsResp, err := iamClient.Projects().List(ctx)
if err != nil {
    log.Fatalf("Failed to list projects: %v", err)
}

fmt.Printf("Total projects: %d\n", projectsResp.Total)

// Direct array access
for i, membership := range projectsResp.Projects {
    project := membership.Project
    fmt.Printf("%d. %s (%s)\n", i+1, project.DisplayName, project.ProjectID)
    fmt.Printf("   User Permission: %s\n", membership.UserPermission.Label)
    fmt.Printf("   Tenant Role: %s\n", membership.TenantRole)
    fmt.Printf("   Frozen: %v\n", membership.Frozen)
}
```

### 4. List Projects with Pagination

```go
import "github.com/Zillaforge/cloud-sdk/models/iam/projects"

// Helper function for creating int pointers
func intPtr(i int) *int { return &i }
func strPtr(s string) *string { return &s }

// Custom pagination using options struct (consistent with VPS/VRM)
opts := &projects.ListProjectsOptions{
    Offset: intPtr(10),        // Skip first 10 projects
    Limit:  intPtr(5),         // Return max 5 projects
    Order:  strPtr("name"),    // Sort by name
}
projectsResp, err := iamClient.Projects().List(ctx, opts)
if err != nil {
    log.Fatalf("Failed to list projects: %v", err)
}

// Process paginated results
fmt.Printf("Showing %d of %d projects\n", 
    len(projectsResp.Projects), projectsResp.Total)
```

### 5. Get Specific Project Details

```go
// Get project by ID
projectID := "14735dfa-5553-46cc-b4bd-405e711b223f"
project, err := iamClient.Projects().Get(ctx, projectID)
if err != nil {
    log.Fatalf("Failed to get project: %v", err)
}

fmt.Printf("Project ID: %s\n", project.ProjectID)
fmt.Printf("Name: %s\n", project.DisplayName)
fmt.Printf("Description: %s\n", project.Description)
fmt.Printf("Namespace: %s\n", project.Namespace)
fmt.Printf("Created: %s\n", project.CreatedAt)
fmt.Printf("Frozen: %v\n", project.Frozen)

// Check permissions
if project.UserPermission != nil {
    fmt.Printf("Your permission: %s (%s)\n", 
        project.UserPermission.Label,
        project.UserPermission.ID)
}

// Access custom metadata
if extra, ok := project.Extra["iservice"].(map[string]interface{}); ok {
    if sysCode, ok := extra["projectSysCode"].(string); ok {
        fmt.Printf("System Code: %s\n", sysCode)
    }
}

// Note: GetProject does NOT include tenantRole or membership data
// Use ListProjects to get tenantRole information
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/Zillaforge/cloud-sdk"
    "github.com/Zillaforge/cloud-sdk/models/iam/common"
    "github.com/Zillaforge/cloud-sdk/modules/iam"
)

func main() {
    // Get token from environment
    token := os.Getenv("IAM_TOKEN")
    if token == "" {
        log.Fatal("IAM_TOKEN environment variable required")
    }

    // Initialize client
    client, err := cloudsdk.New("http://127.0.0.1:8084", token)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    ctx := context.Background()
    iamClient := client.IAM()

    // 1. Get current user
    fmt.Println("=== Current User ===")
    user, err := iamClient.User().Get(ctx)
    if err != nil {
        log.Fatalf("GetUser failed: %v", err)
    }
    fmt.Printf("%s (%s)\n\n", user.DisplayName, user.Account)

    // 2. List projects
    fmt.Println("=== Your Projects ===")
    projectsResp, err := iamClient.Projects().List(ctx, iam.WithLimit(10))
    if err != nil {
        log.Fatalf("ListProjects failed: %v", err)
    }
    
    fmt.Printf("Found %d projects:\n", projectsResp.Total)
    for i, membership := range projectsResp.Projects {
        project := membership.Project
        roleIcon := "👤" // member
        if membership.TenantRole == common.TenantRoleAdmin {
            roleIcon = "🔧" // admin
        } else if membership.TenantRole == common.TenantRoleOwner {
            roleIcon = "👑" // owner
        }
        fmt.Printf("%d. %s %s (%s) - %s\n", 
            i+1,
            roleIcon,
            project.DisplayName,
            membership.TenantRole,
            project.ProjectID)
    }

    // 3. Get details of first project
    if len(projectsResp.Projects) > 0 {
        firstMembership := projectsResp.Projects[0]
        firstProjectID := firstMembership.Project.ProjectID
        
        fmt.Printf("\n=== Project Details: %s ===\n", firstProjectID)
        
        // Option 1: Use membership data from List (has tenantRole)
        fmt.Printf("Name: %s\n", firstMembership.Project.DisplayName)
        fmt.Printf("Description: %s\n", firstMembership.Project.Description)
        fmt.Printf("Your Role: %s\n", firstMembership.TenantRole)
        fmt.Printf("Created: %s\n", firstMembership.Project.CreatedAt)
        fmt.Printf("Updated: %s\n", firstMembership.Project.UpdatedAt)
        
        // Option 2: Get fresh project data (no tenantRole, but latest data)
        resp, err := iamClient.Projects().Get(ctx, firstProjectID)
        if err != nil {
            log.Fatalf("GetProject failed: %v", err)
        }
        fmt.Printf("Latest permission: %s\n", resp.UserPermission.Label)
    }
}
```

## Error Handling

### Handle Authentication Errors

```go
user, err := iamClient.User().Get(ctx)
if err != nil {
    // Check for specific error types
    var sdkErr *types.SDKError
    if errors.As(err, &sdkErr) {
        switch sdkErr.StatusCode {
        case 403:
            log.Fatal("Authentication failed: token expired or invalid")
        case 500:
            log.Fatal("Service error: try again later")
        default:
            log.Fatalf("API error: %d - %s", sdkErr.StatusCode, sdkErr.Message)
        }
    }
    log.Fatalf("Unexpected error: %v", err)
}
```

### Handle Project Not Found

```go
project, err := iamClient.Projects().Get(ctx, projectID)
if err != nil {
    var sdkErr *types.SDKError
    if errors.As(err, &sdkErr) {
        if sdkErr.StatusCode == 400 {
            fmt.Printf("Project %s not found or no access\n", projectID)
            return
        }
    }
    log.Fatalf("Failed to get project: %v", err)
}
```

### Context Timeout

```go
// Set timeout for operation
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

user, err := iamClient.GetUser(ctx)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        log.Fatal("Request timed out after 5 seconds")
    }
    log.Fatalf("Request failed: %v", err)
}
```

## Advanced Configuration

### Custom HTTP Client

```go
import (
    "net/http"
    "time"
)

// Custom HTTP client with longer timeout
httpClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        MaxIdleConnsPerHost: 10,
        IdleConnTimeout:     90 * time.Second,
    },
}

client, err := cloudsdk.New(
    "http://127.0.0.1:8084",
    token,
    cloudsdk.WithHTTPClient(httpClient),
)
```

### Override Base URL

```go
// For different environments
prodClient, err := cloudsdk.New(
    "https://api.production.example.com",
    token,
)

stagingClient, err := cloudsdk.New(
    "https://api.staging.example.com",
    token,
)
```

## Pagination Best Practices

### Iterate Through All Projects

```go
const pageSize = 50
offset := 0

for {
    resp, err := iamClient.Projects().List(ctx,
        iam.WithOffset(offset),
        iam.WithLimit(pageSize),
    )
    if err != nil {
        return err
    }

    // Process current page
    for _, membership := range resp.Projects {
        processProject(membership.Project)
    }

    // Check if more pages exist
    offset += len(resp.Projects)
    if offset >= resp.Total {
        break
    }
}
```

### Load-on-Demand

```go
// Show initial page
resp, err := iamClient.Projects().List(ctx, iam.WithLimit(20))
if err != nil {
    return err
}

// User requests more...
if userWantsMore && resp.Total > len(resp.Projects) {
    nextPage, err := iamClient.Projects().List(ctx,
        iam.WithOffset(20),
        iam.WithLimit(20),
    )
    // ...
}
```

## Testing Your Code

### Unit Test with Mock

```go
func TestMyFunction(t *testing.T) {
    // Use httptest to mock IAM API
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.URL.Path == "/iam/api/v1/user" {
            json.NewEncoder(w).Encode(map[string]interface{}{
                "account": "test@example.com",
                "userId": "test-uuid",
                "displayName": "Test User",
                // ... other fields
            })
        }
    }))
    defer server.Close()

    client, _ := cloudsdk.New(server.URL, "test-token")
    iamClient := client.IAM()

    user, err := iamClient.User().Get(context.Background())
    assert.NoError(t, err)
    assert.Equal(t, "test@example.com", user.Account)
}
```

## Common Patterns

### Check Project Access

```go
func hasProjectAccess(ctx context.Context, iamClient *iam.Client, projectID string) bool {
    _, err := iamClient.Projects().Get(ctx, projectID)
    return err == nil
}

// Note: GetProject only checks if you can access the project,
// but does NOT return your tenantRole. Use ListProjects for role info.
```

### Find Project by Name

```go
func findProjectByName(ctx context.Context, iamClient *iam.Client, name string) (*common.Project, error) {
    resp, err := iamClient.Projects().List(ctx)
    if err != nil {
        return nil, err
    }

    for _, membership := range resp.Projects {
        if membership.Project.DisplayName == name {
            return membership.Project, nil
        }
    }
    return nil, fmt.Errorf("project %s not found", name)
}
```

### Display User Context

```go
func printUserContext(ctx context.Context, iamClient *iam.Client) error {
    user, err := iamClient.User().Get(ctx)
    if err != nil {
        return err
    }

    projectsResp, err := iamClient.Projects().List(ctx)
    if err != nil {
        return err
    }

    fmt.Printf("Logged in as: %s (%s)\n", user.DisplayName, user.Account)
    fmt.Printf("Access to %d projects\n", projectsResp.Total)
    
    return nil
}
```

### Check User Role in Project

```go
import "github.com/Zillaforge/cloud-sdk/models/iam/common"

func canAdministerProject(ctx context.Context, iamClient *iam.Client, projectID string) (bool, error) {
    // Must use ListProjects to get tenantRole - GetProject doesn't include it
    resp, err := iamClient.Projects().List(ctx)
    if err != nil {
        return false, err
    }

    // Find the project in the list
    for _, membership := range resp.Projects {
        if membership.Project.ProjectID == projectID {
            // Check if user is admin or owner
            switch membership.TenantRole {
            case common.TenantRoleAdmin, common.TenantRoleOwner:
                return true, nil
            default:
                return false, nil
            }
        }
    }
    
    return false, fmt.Errorf("project not found in user's project list")
}
```

### Filter Projects by Role

```go
func getOwnedProjects(ctx context.Context, iamClient *iam.Client) ([]*common.Project, error) {
    resp, err := iamClient.Projects().List(ctx)
    if err != nil {
        return nil, err
    }

    var owned []*common.Project
    for _, membership := range resp.Projects {
        if membership.TenantRole == common.TenantRoleOwner {
            owned = append(owned, membership.Project)
        }
    }
    return owned, nil
}
```

## Troubleshooting

### Token Issues
- Ensure token is valid and not expired
- Check token has `Bearer ` prefix removed (client adds it)
- Verify token has appropriate permissions

### Connection Issues
- Verify base URL is correct and accessible
- Check network connectivity to IAM service
- Confirm firewall/proxy settings

### Pagination Issues
- offset must be >= 0
- limit must be 1-100
- Invalid values return 400 error from server

## Next Steps

- Review [API Contracts](contracts/api-contracts.md) for detailed API documentation
- See [Data Model](data-model.md) for complete type definitions
- Check [Implementation Plan](plan.md) for architecture details
- Read SDK documentation for other modules (VPS, VRM)

## Support

For issues or questions:
- Review error messages carefully (include status codes)
- Check IAM service logs if available
- Verify Swagger spec compatibility
