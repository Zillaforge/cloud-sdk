# IAM API Client Examples

Comprehensive usage examples for the IAM client.

## Table of Contents

- [Basic Setup](#basic-setup)
- [User Operations](#user-operations)
- [Project Operations](#project-operations)
- [Error Handling](#error-handling)
- [Pagination](#pagination)

## Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"

    cloudsdk "github.com/Zillaforge/cloud-sdk"
)

func main() {
    // Create SDK client with base URL and bearer token
    client, err := cloudsdk.NewClient(
        "https://api.example.com",
        "your-bearer-token",
    )
    if err != nil {
        log.Fatal(err)
    }

    ctx := context.Background()

    // Use IAM client
    user, err := client.IAM().Users().Get(ctx)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Authenticated as: %s\n", user.DisplayName)
}
```

## User Operations

### Get Current User

```go
// Retrieve authenticated user information
user, err := client.IAM().Users().Get(ctx)
if err != nil {
    log.Printf("Failed to get user: %v", err)
    return
}

fmt.Printf("User ID: %s\n", user.UserID)
fmt.Printf("Account: %s\n", user.Account)
fmt.Printf("Display Name: %s\n", user.DisplayName)
fmt.Printf("Email: %s\n", user.Email)
fmt.Printf("Namespace: %s\n", user.Namespace)
fmt.Printf("MFA Enabled: %v\n", user.MFA)
fmt.Printf("Account Frozen: %v\n", user.Frozen)
fmt.Printf("Last Login: %s\n", user.LastLoginAt)
```

## Project Operations

### List All Projects

```go
// List all projects with default pagination
projects, err := client.IAM().Projects().List(ctx, nil)
if err != nil {
    log.Printf("Failed to list projects: %v", err)
    return
}

fmt.Printf("Found %d projects\n", len(projects))
for _, membership := range projects {
    fmt.Printf("Project: %s (Role: %s)\n",
        membership.Project.DisplayName,
        membership.TenantRole,
    )
}
```

### Get Specific Project

```go
// Retrieve specific project details
projectID := "91457b61-0b92-4aa8-b136-b03d88f04946"
project, err := client.IAM().Projects().Get(ctx, projectID)
if err != nil {
    log.Printf("Failed to get project: %v", err)
    return
}

fmt.Printf("Project ID: %s\n", project.ProjectID)
fmt.Printf("Display Name: %s\n", project.DisplayName)
fmt.Printf("Namespace: %s\n", project.Namespace)
fmt.Printf("User Permission: %s\n", project.UserPermission.Label)
```

## Error Handling

```go
user, err := client.IAM().Users().Get(ctx)
if err != nil {
    // Handle specific error types
    log.Printf("Error: %v", err)
    
    // Check for context cancellation
    if ctx.Err() == context.Canceled {
        log.Println("Request was canceled")
        return
    }
    
    // Check for timeout
    if ctx.Err() == context.DeadlineExceeded {
        log.Println("Request timed out")
        return
    }
    
    return
}

// Success
fmt.Printf("User: %s\n", user.DisplayName)
```

## Pagination

### Custom Pagination Options

```go
// Helper functions for pointer types
func intPtr(i int) *int {
    return &i
}

func strPtr(s string) *string {
    return &s
}

// List projects with custom pagination
opts := &projects.ListProjectsOptions{
    Offset: intPtr(10),
    Limit:  intPtr(20),
    Order:  strPtr("displayName"),
}

result, err := client.IAM().Projects().List(ctx, opts)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Found %d projects\n", len(result))
```

### Iterate Through All Projects

```go
// Fetch all projects with pagination
const pageSize = 50
offset := 0
allProjects := []*projects.ProjectMembership{}

for {
    opts := &projects.ListProjectsOptions{
        Offset: intPtr(offset),
        Limit:  intPtr(pageSize),
    }
    
    projects, err := client.IAM().Projects().List(ctx, opts)
    if err != nil {
        log.Fatal(err)
    }
    
    allProjects = append(allProjects, projects...)
    
    // Check if we've fetched all projects (empty result means no more)
    if len(projects) == 0 {
        break
    }
    
    offset += pageSize
}

fmt.Printf("Fetched %d total projects\n", len(allProjects))
```

---

For more information, see the [README.md](README.md).
