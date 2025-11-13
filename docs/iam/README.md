# IAM API Client

Go client for IAM (Identity and Access Management) API operations.

## Overview

The IAM client provides non-project-scoped access to identity and access management operations:

- **User Operations**: Retrieve current authenticated user information
- **Project Operations**: List and retrieve project details and memberships

## Installation

```go
import (
    "github.com/Zillaforge/cloud-sdk"
    "github.com/Zillaforge/cloud-sdk/modules/iam"
)
```

## Quick Start

```go
// Initialize the SDK client
client, err := cloudsdk.NewClient(
    "https://api.example.com",
    "your-bearer-token",
)
if err != nil {
    log.Fatal(err)
}

// Get current user
ctx := context.Background()
user, err := client.IAM().Users().Get(ctx)
if err != nil {
    log.Fatal(err)
}
fmt.Printf("User: %s (%s)\n", user.DisplayName, user.Email)
```

## API Reference

### User Operations

#### Get Current User

```go
user, err := client.IAM().Users().Get(ctx)
```

Returns the authenticated user's information including account details, permissions, and timestamps.

### Project Operations

#### List Projects

```go
// List all projects (default pagination)
projects, err := client.IAM().Projects().List(ctx, nil)

// List with custom pagination
opts := &projects.ListProjectsOptions{
    Offset: intPtr(10),
    Limit:  intPtr(20),
    Order:  strPtr("displayName"),
}
projects, err := client.IAM().Projects().List(ctx, opts)
```

Returns `[]*ProjectMembership` containing all projects the user belongs to with membership context (permissions, role, frozen status).

#### Get Project Details

```go
project, err := client.IAM().Projects().Get(ctx, "project-uuid")
```

Returns specific project details including permissions for the current user.

## Error Handling

All operations return typed errors that can be inspected:

```go
user, err := client.IAM().Users().Get(ctx)
if err != nil {
    // Handle error
    log.Printf("Failed to get user: %v", err)
    return
}
```

See [EXAMPLES.md](EXAMPLES.md) for more detailed usage examples.
