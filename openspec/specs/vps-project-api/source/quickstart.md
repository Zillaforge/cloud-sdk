# VPS SDK Quickstart Guide

**Get started with the Cloud SDK for VPS in 5 minutes**

## Installation

```bash
go get github.com/Zillaforge/cloud-sdk@v0.1.0
```

## Basic Usage

### 1. Initialize the SDK

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    cloudsdk "github.com/Zillaforge/cloud-sdk"
)

func main() {
    // Initialize the root client with base URL and bearer token
    client := cloudsdk.NewClient(
        os.Getenv("CLOUD_SDK_BASE_URL"), // e.g., https://api.example.com
        os.Getenv("CLOUD_SDK_TOKEN"),    // Bearer token
    )

    // Get a project-scoped VPS client
    projectID := os.Getenv("PROJECT_ID")
    vpsClient := client.Project(projectID).VPS()

    ctx := context.Background()

    // List all servers in the project
    servers, err := vpsClient.Servers().List(ctx, nil)
    if err != nil {
        // Error handling: SDKError includes StatusCode, ErrorCode, Message, Meta
        if sdkErr, ok := err.(*cloudsdk.SDKError); ok {
            log.Fatalf("Failed to list servers: [%d] %s - %s",
                sdkErr.StatusCode, sdkErr.ErrorCode, sdkErr.Message)
        }
        log.Fatalf("Unexpected error: %v", err)
    }

    fmt.Printf("Found %d servers\n", len(servers.Items))
    for _, server := range servers.Items {
        fmt.Printf("- %s (ID: %s, Status: %s)\n", server.Name, server.ID, server.Status)
    }
}
```

### 2. Create a Server

```go
import (
    "time"
    cloudsdk "github.com/Zillaforge/cloud-sdk"
    "github.com/Zillaforge/cloud-sdk/modules/vps"
)

func createServer(vpsClient *vps.Client) {
    ctx := context.Background()

    // Create a server
    req := &vps.ServerCreateRequest{
        Name:     "my-web-server",
        FlavorID: "flavor-001",
        ImageID:  "ubuntu-22.04",
        Networks: []string{"network-001"},
    }

    server, err := vpsClient.Servers().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }

    fmt.Printf("Server created: %s (ID: %s, Status: %s)\n",
        server.Name, server.ID, server.Status)

    // Wait for server to become active (optional waiter helper)
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
    defer cancel()

    activeServer, err := vpsClient.Servers().WaitForServerStatus(ctx, server.ID, "ACTIVE", 10*time.Second)
    if err != nil {
        log.Fatalf("Timeout waiting for server: %v", err)
    }

    fmt.Printf("Server is now active: %s\n", activeServer.ID)
}
```

### 3. Manage Networks

```go
func listNetworks(vpsClient *vps.Client) {
    ctx := context.Background()

    networks, err := vpsClient.Networks().List(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list networks: %v", err)
    }

    for _, net := range networks.Items {
        fmt.Printf("Network: %s (CIDR: %s)\n", net.Name, net.CIDR)
    }
}
```

### 4. Associate Floating IP

```go
func associateFloatingIP(vpsClient *vps.Client, floatingIPID, serverID, nicID string) {
    ctx := context.Background()

    // Get the server first
    server, err := vpsClient.Servers().Get(ctx, serverID)
    if err != nil {
        log.Fatalf("Failed to get server: %v", err)
    }

    // Associate floating IP to a NIC
    req := &vps.FloatingIPAssociateRequest{
        FloatingIPID: floatingIPID,
    }
    err = server.NICs().AssociateFloatingIP(ctx, nicID, req)
    if err != nil {
        log.Fatalf("Failed to associate floating IP: %v", err)
    }

    fmt.Println("Floating IP associated successfully")
}
```

### 5. Manage Server NICs and Volumes

```go
func manageServerResources(vpsClient *vps.Client, serverID string) {
    ctx := context.Background()

    // Get the server
    server, err := vpsClient.Servers().Get(ctx, serverID)
    if err != nil {
        log.Fatalf("Failed to get server: %v", err)
    }

    // List NICs
    nics, err := server.NICs().List(ctx)
    if err != nil {
        log.Fatalf("Failed to list NICs: %v", err)
    }
    fmt.Printf("Found %d NICs\n", len(nics))

    // List volumes
    volumes, err := server.Volumes().List(ctx)
    if err != nil {
        log.Fatalf("Failed to list volumes: %v", err)
    }
    fmt.Printf("Found %d volumes\n", len(volumes))

    // Attach a volume
    err = server.Volumes().Attach(ctx, "volume-001")
    if err != nil {
        log.Fatalf("Failed to attach volume: %v", err)
    }
    fmt.Println("Volume attached successfully")
}
```

### 6. Manage Router Networks

```go
func manageRouterNetworks(vpsClient *vps.Client, routerID, networkID string) {
    ctx := context.Background()

    // Get the router
    router, err := vpsClient.Routers().Get(ctx, routerID)
    if err != nil {
        log.Fatalf("Failed to get router: %v", err)
    }

    // List associated networks
    networks, err := router.Networks().List(ctx)
    if err != nil {
        log.Fatalf("Failed to list router networks: %v", err)
    }
    fmt.Printf("Router has %d networks\n", len(networks))

    // Associate a network
    err = router.Networks().Associate(ctx, networkID)
    if err != nil {
        log.Fatalf("Failed to associate network: %v", err)
    }
    fmt.Println("Network associated successfully")
}
```

### 7. Manage Security Group Rules

```go
func manageSecurityGroupRules(vpsClient *vps.Client, sgID string) {
    ctx := context.Background()

    // Get the security group
    sg, err := vpsClient.SecurityGroups().Get(ctx, sgID)
    if err != nil {
        log.Fatalf("Failed to get security group: %v", err)
    }

    // Add a rule
    ruleReq := &vps.SecurityGroupRuleCreateRequest{
        Direction: "ingress",
        Protocol:  "tcp",
        PortMin:   80,
        PortMax:   80,
        RemoteCIDR: "0.0.0.0/0",
    }
    rule, err := sg.Rules().Add(ctx, ruleReq)
    if err != nil {
        log.Fatalf("Failed to add rule: %v", err)
    }
    fmt.Printf("Rule added: %s\n", rule.ID)

    // Delete a rule
    err = sg.Rules().Delete(ctx, rule.ID)
    if err != nil {
        log.Fatalf("Failed to delete rule: %v", err)
    }
    fmt.Println("Rule deleted successfully")
}
```

### 8. Error Handling

The SDK provides structured error information through `SDKError`:

```go
import cloudsdk "github.com/Zillaforge/cloud-sdk"

_, err := vpsClient.Servers().Get(ctx, "invalid-id")
if err != nil {
    if sdkErr, ok := err.(*cloudsdk.SDKError); ok {
        switch {
        case sdkErr.StatusCode == 404:
            fmt.Println("Server not found")
        case sdkErr.StatusCode == 0:
            // Client-side error (network issue, timeout, etc.)
            fmt.Printf("Client error: %s\n", sdkErr.Message)
        case sdkErr.StatusCode >= 500:
            // Server error - automatic retry may have been attempted
            fmt.Printf("Server error: %s (retries exhausted)\n", sdkErr.Message)
        default:
            fmt.Printf("API error [%d]: %s\n", sdkErr.StatusCode, sdkErr.Message)
        }
        // Access additional metadata if available
        if sdkErr.Meta != nil {
            fmt.Printf("Details: %+v\n", sdkErr.Meta)
        }
    }
}
```

## Configuration Options

### Custom Timeout

```go
// Per-request timeout via context
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

servers, err := vpsClient.Servers().List(ctx, nil)
```

### Custom HTTP Client

```go
import "net/http"

customClient := &http.Client{
    Timeout: 60 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns: 100,
    },
}

client := cloudsdk.NewClient(
    baseURL,
    token,
    cloudsdk.WithHTTPClient(customClient),
)
```

### Enable Debug Logging

```go
import "log"

client := cloudsdk.NewClient(
    baseURL,
    token,
    cloudsdk.WithLogger(log.Default()),
)
```

## Retry Behavior

The SDK automatically retries safe read operations (GET, HEAD) on transient errors:

- **Retry Codes**: 429 (Too Many Requests), 502 (Bad Gateway), 503 (Service Unavailable), 504 (Gateway Timeout)
- **Strategy**: Exponential backoff with jitter
- **Max Attempts**: 3 (initial + 2 retries)
- **Backoff**: 1s → 2s → 4s (with +/-20% jitter)

Retry attempts are logged at debug level when a logger is configured.

## Async Operations & Waiters

Long-running operations (server creation, deletion, etc.) return immediately. Use optional waiter helpers for convenience:

```go
// Non-blocking: returns immediately with server in BUILD state
server, _ := vpsClient.Servers().Create(ctx, req)

// Optional: wait for server to reach ACTIVE state
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

activeServer, err := vpsClient.Servers().WaitForServerStatus(ctx, server.ID, "ACTIVE", 10*time.Second)
// Polls every 10s until status matches or context times out
```

Available waiters:
- `WaitForServerStatus(ctx, serverID, targetStatus, pollInterval)`
- `WaitForFloatingIPActive(ctx, floatingIPID, pollInterval)`

## Next Steps

- **Full API Reference**: See `contracts/*.go` for all available operations
- **Data Models**: See `data-model.md` for detailed entity schemas
- **Error Handling**: See `spec.md` Error Handling Contract for complete error taxonomy
- **Testing**: Run `go test ./vps/...` for examples of SDK usage patterns

## Support & Resources

- **Swagger API Docs**: `swagger/vps.json` (source of truth)
- **Feature Spec**: `specs/1-vps-project-api/spec.md`
- **Implementation Plan**: `specs/1-vps-project-api/plan.md`
- **Constitution**: `.specify/memory/constitution.md` (TDD guidelines, versioning policy)

---

**Ready to build?** Dive into the VPS SDK and manage your infrastructure with confidence! 🚀
