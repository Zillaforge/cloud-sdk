# VPS Module

The VPS module provides a comprehensive Go SDK for managing Virtual Private Server (VPS) resources through a RESTful API. It offers type-safe, idiomatic Go interfaces for all VPS operations with automatic retry, structured errors, and full context support.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Quick Start](#quick-start)
- [Architecture](#architecture)
- [Resources](#resources)
  - [Servers](#servers)
  - [Networks](#networks)
  - [Floating IPs](#floating-ips)
  - [Keypairs](#keypairs)
  - [Routers](#routers)
  - [Security Groups](#security-groups)
  - [Flavors](#flavors)
  - [Quotas](#quotas)
- [Error Handling](#error-handling)
- [Retry Logic](#retry-logic)
- [Timeouts and Context](#timeouts-and-context)
- [Testing](#testing)
- [Examples](#examples)

## Features

- 🚀 **Full VPS Resource Coverage**: Manage servers, networks, floating IPs, keypairs, routers, security groups, flavors, and quotas
- 🔒 **Type-Safe**: Idiomatic Go interfaces with strong typing and compile-time safety
- 📦 **Project-Scoped**: Bind operations to a project once, no repetition
- 🔄 **Automatic Retry**: Exponential backoff with jitter for transient failures (429, 502, 503, 504)
- ❌ **Structured Errors**: Rich error information with HTTP status codes and metadata
- ⏱️ **Context Support**: Full timeout control and cancellation via context.Context
- 🧪 **Well-Tested**: 95%+ code coverage with unit, integration, and contract tests
- 🎯 **Resource-Oriented Design**: Clean, consistent API across all resources
- 🔌 **Sub-Resource Support**: Access related resources like NICs, volumes, ports, and networks

## Installation

```bash
go get github.com/Zillaforge/cloud-sdk
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    cloudsdk "github.com/Zillaforge/cloud-sdk"
)

func main() {
    // Create client with base URL and bearer token
    client, err := cloudsdk.New("https://api.example.com", "your-bearer-token")
    if err != nil {
        log.Fatal(err)
    }
    
    // Get project-scoped VPS client
    vps := client.Project("your-project-id").VPS()
    
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // List all servers
    servers, err := vps.Servers().List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Found %d servers\n", len(servers.Servers))
}
```

## Architecture

The VPS module follows a **coordinator pattern** with resource-specific sub-modules:

```
modules/vps/
├── client.go           # VPS coordinator (entry point)
├── waiters.go          # Server status waiters
├── flavors/            # Flavor discovery operations
├── floatingips/        # Floating IP management
├── keypairs/           # SSH keypair management
├── networks/           # Network management
│   └── ports.go        # Network ports sub-resource
├── quotas/             # Project quota viewing
├── routers/            # Router management
│   └── networks.go     # Router networks sub-resource
├── securitygroups/     # Security group management
└── servers/            # Server management
    ├── nics.go         # Server NICs sub-resource
    └── volumes.go      # Server volumes sub-resource
```

Each resource module provides:
- **Client**: Resource-specific operations (List, Create, Get, Update, Delete)
- **Models**: Strongly-typed request/response structures
- **Tests**: Unit tests with 90%+ coverage
- **Contract Tests**: API contract validation against Swagger specs

## Resources

### Servers

Manage virtual server instances with full lifecycle operations.

```go
vps := client.Project("project-id").VPS()

// Create a server
server, err := vps.Servers().Create(ctx, &servers.CreateRequest{
    Name:     "web-server-01",
    FlavorID: "flavor-123",
    ImageID:  "ubuntu-20.04",
})

// Perform actions
err = vps.Servers().Action(ctx, serverID, &servers.ActionRequest{
    Action: "start",
})

// Access sub-resources
nics, err := vps.Servers().Resource(serverID).NICs().List(ctx)
volumes, err := vps.Servers().Resource(serverID).Volumes().List(ctx)

// Get metrics
metrics, err := vps.Servers().Metrics(ctx, serverID, time.Now().Add(-1*time.Hour), time.Now())

// Get VNC console URL
vncURL, err := vps.Servers().VNCURL(ctx, serverID)
```

**Operations**: List, Create, Get, Update, Delete, Action, Metrics, VNCURL  
**Sub-resources**: NICs (List, Add, Update, Delete, AssociateFloatingIP), Volumes (List, Attach, Detach)

### Networks

Manage private networks with CIDR configuration.

```go
vps := client.Project("project-id").VPS()

// Create a network
network, err := vps.Networks().Create(ctx, &networks.CreateRequest{
    Name:        "private-net",
    CIDR:        "10.0.0.0/24",
    Description: "Private network for web tier",
})

// List network ports
ports, err := vps.Networks().Resource(networkID).Ports().List(ctx)
```

**Operations**: List, Create, Get, Update, Delete  
**Sub-resources**: Ports (List)

### Floating IPs

Manage floating IP addresses with approval workflows.

```go
vps := client.Project("project-id").VPS()

// Create a floating IP
fip, err := vps.FloatingIPs().Create(ctx, &floatingips.CreateRequest{
    Description: "Web server public IP",
})

// Approve pending floating IP
err = vps.FloatingIPs().Approve(ctx, fipID)

// Disassociate from resource
err = vps.FloatingIPs().Disassociate(ctx, fipID)
```

**Operations**: List, Create, Get, Update, Delete, Approve, Reject, Disassociate

### Keypairs

Manage SSH keypairs for server access.

```go
vps := client.Project("project-id").VPS()

// Create a keypair
keypair, err := vps.Keypairs().Create(ctx, &keypairs.CreateRequest{
    Name:      "deploy-key",
    PublicKey: "ssh-rsa AAAAB3NzaC1...",
})

// Or generate a new keypair
generated, err := vps.Keypairs().Create(ctx, &keypairs.CreateRequest{
    Name: "auto-generated-key",
    // Leave PublicKey empty to auto-generate
})
```

**Operations**: List, Create, Get, Update, Delete

### Routers

Manage routers with network associations.

```go
vps := client.Project("project-id").VPS()

// Create a router
router, err := vps.Routers().Create(ctx, &routers.CreateRequest{
    Name:              "main-router",
    ExternalNetworkID: "ext-net-123",
})

// Set router state
err = vps.Routers().SetState(ctx, routerID, &routers.SetStateRequest{
    Enabled: true,
})

// Associate networks
err = vps.Routers().Resource(routerID).Networks().Associate(ctx, &routers.AssociateNetworkRequest{
    NetworkID: "net-456",
})
```

**Operations**: List, Create, Get, Update, Delete, SetState  
**Sub-resources**: Networks (List, Associate, Disassociate)

### Security Groups

Manage security groups with firewall rules.

```go
vps := client.Project("project-id").VPS()

// Create a security group
sg, err := vps.SecurityGroups().Create(ctx, &securitygroups.CreateRequest{
    Name:        "web-tier",
    Description: "Allow HTTP/HTTPS traffic",
    Rules: []securitygroups.RuleRequest{
        {
            Direction:   "ingress",
            Protocol:    "tcp",
            PortMin:     80,
            PortMax:     80,
            RemoteCIDR:  "0.0.0.0/0",
        },
        {
            Direction:   "ingress",
            Protocol:    "tcp",
            PortMin:     443,
            PortMax:     443,
            RemoteCIDR:  "0.0.0.0/0",
        },
    },
})
```

**Operations**: List, Create, Get, Update, Delete

### Flavors

Discover available instance sizes and configurations.

```go
vps := client.Project("project-id").VPS()

// List all flavors
flavors, err := vps.Flavors().List(ctx, nil)

// Filter by characteristics
flavors, err = vps.Flavors().List(ctx, &flavors.ListOptions{
    Name:   "m1.large",
    Public: true,
    Tag:    "production",
})

// Get specific flavor
flavor, err := vps.Flavors().Get(ctx, "flavor-123")
```

**Operations**: List, Get

### Quotas

View project resource quotas and usage.

```go
vps := client.Project("project-id").VPS()

// Get current quotas
quotas, err := vps.Quotas().Get(ctx)

fmt.Printf("VMs: %d/%d\n", quotas.VMs.Used, quotas.VMs.Limit)
fmt.Printf("vCPUs: %d/%d\n", quotas.VCPUs.Used, quotas.VCPUs.Limit)
fmt.Printf("RAM: %d/%d MB\n", quotas.RAM.Used, quotas.RAM.Limit)
```

**Operations**: Get

## Error Handling

All operations return structured errors with HTTP status codes and detailed messages:

```go
servers, err := vps.Servers().List(ctx, nil)
if err != nil {
    if sdkErr, ok := err.(*cloudsdk.SDKError); ok {
        fmt.Printf("HTTP %d: %s\n", sdkErr.StatusCode, sdkErr.Message)
        if sdkErr.StatusCode == 401 {
            // Handle authentication error
        }
    }
    return err
}
```

Common error codes:
- **400**: Validation error (check request parameters)
- **401**: Unauthorized (invalid or expired token)
- **403**: Forbidden (insufficient permissions)
- **404**: Resource not found
- **409**: Conflict (resource already exists or quota exceeded)
- **429**: Rate limited (automatic retry with backoff)
- **500**: Internal server error
- **502/503/504**: Service unavailable (automatic retry with backoff)

## Retry Logic

The SDK automatically retries transient failures with exponential backoff and jitter:

- **Retry on**: HTTP 429, 502, 503, 504 for GET and HEAD requests
- **Base delay**: 100ms
- **Max delay**: 5 seconds
- **Max attempts**: 3
- **Jitter**: Random ±10% to prevent thundering herd

```go
// Retry is automatic - no configuration needed
servers, err := vps.Servers().List(ctx, nil)
// Will retry automatically on 429/502/503/504
```

## Timeouts and Context

All operations accept `context.Context` for timeout control and cancellation:

```go
// Set a 10-second timeout
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

servers, err := vps.Servers().List(ctx, nil)
if err != nil {
    // Check for timeout
    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("Request timed out")
    }
}
```

Default timeout: **30 seconds** (configurable via context)

## Testing

The VPS module has comprehensive test coverage:

- **Unit Tests**: 95%+ coverage for all resource clients
- **Integration Tests**: Full lifecycle tests for each resource
- **Contract Tests**: API contract validation against Swagger specs

Run tests:

```bash
# Run all tests
make test

# Run tests with coverage
make coverage

# Run tests for specific resource
go test ./modules/vps/servers/...
```

## Examples

See [EXAMPLES.md](./EXAMPLES.md) for detailed code samples covering:

- Complete server lifecycle (create, start, stop, resize, delete)
- Network management with port inspection
- Floating IP association and approval workflows
- SSH keypair management
- Router configuration with network associations
- Security group creation with firewall rules
- Flavor discovery and selection
- Quota monitoring and capacity planning

## API Documentation

Full API documentation is available via godoc:

```bash
godoc -http=:6060
# Visit http://localhost:6060/pkg/github.com/Zillaforge/cloud-sdk/modules/vps/
```

## Contributing

Contributions are welcome! Please ensure:

- All tests pass (`make test`)
- Code coverage remains above 75% (`make coverage`)
- Code is formatted (`make fmt`)
- Code passes linting (`make lint`)

## License

See [LICENSE](../../LICENSE) file for details.

## Support

For issues, questions, or feature requests, please open an issue on GitHub.
