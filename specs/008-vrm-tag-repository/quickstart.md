# Quickstart: VRM Tag and Repository APIs

**Feature**: VRM Tag and Repository APIs Client SDK  
**Date**: 2025-11-12  
**Audience**: Go developers integrating VRM functionality

## Overview

The VRM (Virtual Registry Management) SDK provides Go developers with idiomatic access to repository and tag management APIs. This quickstart demonstrates common workflows for managing virtual image registries.

---

## Prerequisites

- Go 1.21 or later
- Cloud SDK installed: `go get github.com/Zillaforge/cloud-sdk`
- Valid bearer token for authentication
- Project ID for scoped operations

---

## Installation

```bash
go get github.com/Zillaforge/cloud-sdk@latest
```

---

## Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk"
)

func main() {
    // Create SDK client with base URL and bearer token
    client, err := cloudsdk.New(
        "https://api.example.com",  // Base API URL
        "your-bearer-token-here",   // Authentication token
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Get project-scoped VRM client
    projectID := "14735dfa-5553-46cc-b4bd-405e711b223f"
    vrm := client.Project(projectID).VRM()
    
    // Now you can access Repositories() and Tags()
    fmt.Println("VRM client ready!")
}
```

---

## Repository Management

### Create a Repository

```go
import (
    "github.com/Zillaforge/cloud-sdk/models/vrm/repositories"
)

func createRepository(ctx context.Context, vrm *vrm.Client) {
    req := &repositories.CreateRepositoryRequest{
        Name:            "ubuntu",
        OperatingSystem: "linux",
        Description:     "Ubuntu base images for containers",
    }
    
    repo, err := vrm.Repositories().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to create repository: %v", err)
    }
    
    fmt.Printf("Created repository: %s (ID: %s)\n", repo.Name, repo.ID)
}
```

### List Repositories

```go
func listRepositories(ctx context.Context, vrm *vrm.Client) {
    // List with filters
    opts := &repositories.ListRepositoriesOptions{
        Limit:     10,
        Offset:    0,
        Where:     []string{"os=linux"},
        Namespace: "public",
    }
    
    repos, err := vrm.Repositories().List(ctx, opts)
    if err != nil {
        log.Fatalf("Failed to list repositories: %v", err)
    }
    
    fmt.Printf("Found %d repositories:\n", len(repos))
    for i, repo := range repos {
        fmt.Printf("%d. %s (%s, %d tags)\n", 
            i+1, repo.Name, repo.OperatingSystem, repo.Count)
    }
}
```

### Get Repository Details

```go
func getRepository(ctx context.Context, vrm *vrm.Client, repositoryID string) {
    repo, err := vrm.Repositories().Get(ctx, repositoryID)
    if err != nil {
        log.Fatalf("Failed to get repository: %v", err)
    }
    
    fmt.Printf("Repository: %s\n", repo.Name)
    fmt.Printf("  OS: %s\n", repo.OperatingSystem)
    fmt.Printf("  Namespace: %s\n", repo.Namespace)
    fmt.Printf("  Tags: %d\n", repo.Count)
    fmt.Printf("  Creator: %s\n", repo.Creator.Name)
    fmt.Printf("  Created: %s\n", repo.CreatedAt)
}
```

### Update Repository

```go
func updateRepository(ctx context.Context, vrm *vrm.Client, repositoryID string) {
    req := &repositories.UpdateRepositoryRequest{
        Description: "Updated: Ubuntu LTS images",
    }
    
    repo, err := vrm.Repositories().Update(ctx, repositoryID, req)
    if err != nil {
        log.Fatalf("Failed to update repository: %v", err)
    }
    
    fmt.Printf("Updated repository: %s\n", repo.Name)
}
```

### Delete Repository

```go
func deleteRepository(ctx context.Context, vrm *vrm.Client, repositoryID string) {
    err := vrm.Repositories().Delete(ctx, repositoryID)
    if err != nil {
        log.Fatalf("Failed to delete repository: %v", err)
    }
    
    fmt.Println("Repository deleted successfully")
}
```

---

## Tag Management

### Create a Tag

```go
import (
    "github.com/Zillaforge/cloud-sdk/models/vrm/tags"
)

func createTag(ctx context.Context, vrm *vrm.Client, repositoryID string) {
    req := &tags.CreateTagRequest{
        Name:            "24.04",
        Type:            "common",
        DiskFormat:      "qcow2",
        ContainerFormat: "bare",
    }
    
    tag, err := vrm.Tags().Create(ctx, repositoryID, req)
    if err != nil {
        log.Fatalf("Failed to create tag: %v", err)
    }
    
    fmt.Printf("Created tag: %s (ID: %s)\n", tag.Name, tag.ID)
    fmt.Printf("  Repository: %s\n", tag.Repository.Name)
}
```

### List All Tags

```go
func listAllTags(ctx context.Context, vrm *vrm.Client) {
    opts := &tags.ListTagsOptions{
        Limit:  20,
        Offset: 0,
        Where:  []string{"type=common"},
    }
    
    tagList, err := vrm.Tags().List(ctx, opts)
    if err != nil {
        log.Fatalf("Failed to list tags: %v", err)
    }
    
    fmt.Printf("Found %d tags:\n", len(tagList))
    for _, tag := range tagList {
        fmt.Printf("- %s:%s (Repository: %s)\n", 
            tag.Repository.Name, tag.Name, tag.RepositoryID)
    }
}
```

### List Tags by Repository

```go
func listRepositoryTags(ctx context.Context, vrm *vrm.Client, repositoryID string) {
    opts := &tags.ListTagsOptions{
        Limit:  10,
        Offset: 0,
    }
    
    tagList, err := vrm.Tags().ListByRepository(ctx, repositoryID, opts)
    if err != nil {
        log.Fatalf("Failed to list repository tags: %v", err)
    }
    
    fmt.Printf("Repository has %d tags:\n", len(tagList))
    for _, tag := range tagList {
        fmt.Printf("- %s (Type: %s, Size: %d bytes)\n", 
            tag.Name, tag.Type, tag.Size)
    }
}
```

### Get Tag Details

```go
func getTag(ctx context.Context, vrm *vrm.Client, tagID string) {
    tag, err := vrm.Tags().Get(ctx, tagID)
    if err != nil {
        log.Fatalf("Failed to get tag: %v", err)
    }
    
    fmt.Printf("Tag: %s\n", tag.Name)
    fmt.Printf("  Repository: %s\n", tag.Repository.Name)
    fmt.Printf("  Type: %s\n", tag.Type)
    fmt.Printf("  Size: %d bytes\n", tag.Size)
    fmt.Printf("  Status: %s\n", tag.Status)
    fmt.Printf("  Created: %s\n", tag.CreatedAt)
}
```

### Update Tag

```go
func updateTag(ctx context.Context, vrm *vrm.Client, tagID string) {
    req := &tags.UpdateTagRequest{
        Name: "24.04-lts",
    }
    
    tag, err := vrm.Tags().Update(ctx, tagID, req)
    if err != nil {
        log.Fatalf("Failed to update tag: %v", err)
    }
    
    fmt.Printf("Updated tag: %s\n", tag.Name)
}
```

### Delete Tag

```go
func deleteTag(ctx context.Context, vrm *vrm.Client, tagID string) {
    err := vrm.Tags().Delete(ctx, tagID)
    if err != nil {
        log.Fatalf("Failed to delete tag: %v", err)
    }
    
    fmt.Println("Tag deleted successfully")
}
```

---

## Complete CRUD Workflow

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/Zillaforge/cloud-sdk"
    "github.com/Zillaforge/cloud-sdk/models/vrm/repositories"
    "github.com/Zillaforge/cloud-sdk/models/vrm/tags"
)

func main() {
    // Setup
    ctx := context.Background()
    client, _ := cloudsdk.New("https://api.example.com", "bearer-token")
    vrm := client.Project("project-id").VRM()
    
    // 1. Create repository
    fmt.Println("=== Creating Repository ===")
    repoReq := &repositories.CreateRepositoryRequest{
        Name:            "demo-app",
        OperatingSystem: "linux",
        Description:     "Demo application images",
    }
    repo, err := vrm.Repositories().Create(ctx, repoReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("✓ Created: %s (ID: %s)\n\n", repo.Name, repo.ID)
    
    // 2. Create tag in repository
    fmt.Println("=== Creating Tag ===")
    tagReq := &tags.CreateTagRequest{
        Name:            "v1.0.0",
        Type:            "common",
        DiskFormat:      "qcow2",
        ContainerFormat: "bare",
    }
    tag, err := vrm.Tags().Create(ctx, repo.ID, tagReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("✓ Created: %s (ID: %s)\n\n", tag.Name, tag.ID)
    
    // 3. List repository tags
    fmt.Println("=== Listing Repository Tags ===")
    tagList, err := vrm.Tags().ListByRepository(ctx, repo.ID, nil)
    if err != nil {
        log.Fatal(err)
    }
    for i, t := range tagList {
        fmt.Printf("%d. %s (Type: %s)\n", i+1, t.Name, t.Type)
    }
    fmt.Println()
    
    // 4. Update tag
    fmt.Println("=== Updating Tag ===")
    updateReq := &tags.UpdateTagRequest{
        Name: "v1.0.1",
    }
    updatedTag, err := vrm.Tags().Update(ctx, tag.ID, updateReq)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("✓ Updated: %s -> %s\n\n", tag.Name, updatedTag.Name)
    
    // 5. Get tag details
    fmt.Println("=== Getting Tag Details ===")
    fetchedTag, err := vrm.Tags().Get(ctx, tag.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Name: %s\n", fetchedTag.Name)
    fmt.Printf("Repository: %s\n", fetchedTag.Repository.Name)
    fmt.Printf("Created: %s\n\n", fetchedTag.CreatedAt.Format(time.RFC3339))
    
    // 6. Delete tag
    fmt.Println("=== Deleting Tag ===")
    err = vrm.Tags().Delete(ctx, tag.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("✓ Tag deleted\n")
    
    // 7. Delete repository
    fmt.Println("=== Deleting Repository ===")
    err = vrm.Repositories().Delete(ctx, repo.ID)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("✓ Repository deleted\n")
    
    fmt.Println("=== Workflow Complete ===")
}
```

---

## Advanced Usage

### Pagination

```go
func paginateRepositories(ctx context.Context, vrm *vrm.Client) {
    pageSize := 10
    offset := 0
    
    for {
        opts := &repositories.ListRepositoriesOptions{
            Limit:  pageSize,
            Offset: offset,
        }
        
        repos, err := vrm.Repositories().List(ctx, opts)
        if err != nil {
            log.Fatal(err)
        }
        
        if len(repos) == 0 {
            break // No more results
        }
        
        for _, repo := range repos {
            fmt.Printf("- %s\n", repo.Name)
        }
        
        offset += pageSize
    }
}
```

### Context with Timeout

```go
func createWithTimeout(parentCtx context.Context, vrm *vrm.Client) {
    // Set 10-second timeout for this operation
    ctx, cancel := context.WithTimeout(parentCtx, 10*time.Second)
    defer cancel()
    
    req := &repositories.CreateRepositoryRequest{
        Name:            "timeout-demo",
        OperatingSystem: "linux",
    }
    
    repo, err := vrm.Repositories().Create(ctx, req)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("Operation timed out")
        } else {
            fmt.Printf("Error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Created: %s\n", repo.Name)
}
```

### Namespace Filtering

```go
func listPublicRepositories(ctx context.Context, vrm *vrm.Client) {
    opts := &repositories.ListRepositoriesOptions{
        Namespace: "public",
    }
    
    repos, err := vrm.Repositories().List(ctx, opts)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Public repositories: %d\n", len(repos))
}

func listPrivateRepositories(ctx context.Context, vrm *vrm.Client) {
    opts := &repositories.ListRepositoriesOptions{
        Namespace: "private",
    }
    
    repos, err := vrm.Repositories().List(ctx, opts)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Private repositories: %d\n", len(repos))
}
```

### Error Handling

```go
import (
    "errors"
    "github.com/Zillaforge/cloud-sdk/internal/types"
)

func handleErrors(ctx context.Context, vrm *vrm.Client) {
    repo, err := vrm.Repositories().Get(ctx, "non-existent-id")
    if err != nil {
        var sdkErr *types.SDKError
        if errors.As(err, &sdkErr) {
            switch sdkErr.StatusCode {
            case 404:
                fmt.Println("Repository not found")
            case 403:
                fmt.Println("Permission denied")
            case 401:
                fmt.Println("Authentication failed")
            default:
                fmt.Printf("API error: %s (code: %d)\n", 
                    sdkErr.Message, sdkErr.ErrorCode)
            }
        } else {
            fmt.Printf("Client error: %v\n", err)
        }
        return
    }
    
    fmt.Printf("Found repository: %s\n", repo.Name)
}
```

---

## Testing

### Unit Test Example

```go
package main

import (
    "context"
    "testing"
    
    "github.com/Zillaforge/cloud-sdk/models/vrm/repositories"
)

func TestRepositoryValidation(t *testing.T) {
    repo := &repositories.Repository{
        ID:              "test-id",
        Name:            "test-repo",
        Namespace:       "public",
        OperatingSystem: "linux",
        Count:           0,
    }
    
    if err := repo.Validate(); err != nil {
        t.Errorf("Valid repository failed validation: %v", err)
    }
}
```

---

## Next Steps

- **API Reference**: See `contracts/api-contracts.md` for complete API specification
- **Data Models**: See `data-model.md` for detailed field descriptions
- **Implementation Plan**: See `plan.md` for architecture and design decisions
- **Testing**: Run `go test ./...` to execute the test suite

---

## Common Pitfalls

1. **Missing Context**: Always pass a valid `context.Context` to API methods
2. **Token Management**: Bearer token is set at client creation and used for all requests
3. **Project Scoping**: VRM operations are project-scoped via `.Project(projectID).VRM()`
4. **List Return Type**: List methods return `[]*Resource` directly, not wrapped in a response struct
5. **Error Wrapping**: Check error types using `errors.As()` for detailed error information

---

## Support

For issues and questions:
- GitHub Issues: https://github.com/Zillaforge/cloud-sdk/issues
- Documentation: https://github.com/Zillaforge/cloud-sdk/tree/main/docs
