# VPS Module Examples

This document provides comprehensive code examples for common VPS operations using the Cloud SDK.

## Table of Contents

1. [Setup and Initialization](#setup-and-initialization)
2. [Server Management](#server-management)
3. [Network Management](#network-management)
4. [Floating IP Management](#floating-ip-management)
5. [SSH Keypair Management](#ssh-keypair-management)
6. [Router Management](#router-management)
7. [Security Group Management](#security-group-management)
8. [Flavor Discovery](#flavor-discovery)
9. [Quota Management](#quota-management)
10. [Error Handling](#error-handling)
11. [Advanced Patterns](#advanced-patterns)

---

## Setup and Initialization

### Basic Client Setup

```go
package main

import (
    "context"
    "log"
    "time"
    
    cloudsdk "github.com/Zillaforge/cloud-sdk"
)

func main() {
    // Create SDK client
    client, err := cloudsdk.New(
        "https://api.cloud.example.com",
        "your-bearer-token-here",
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    // Get project-scoped VPS client
    vps := client.Project("your-project-id").VPS()
    
    // Create context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    // Now use vps for all operations...
}
```

### Client with Custom Logger

```go
import (
    "log"
    "os"
    
    cloudsdk "github.com/Zillaforge/cloud-sdk"
)

func main() {
    logger := log.New(os.Stdout, "[SDK] ", log.LstdFlags)
    
    client, err := cloudsdk.New(
        "https://api.cloud.example.com",
        "your-bearer-token-here",
        cloudsdk.WithLogger(logger),
    )
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }
    
    vps := client.Project("project-123").VPS()
    // Operations will now be logged
}
```

---

## Server Management

### Create a Server

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

func createServer(ctx context.Context, vps *vps.Client) {
    req := &servers.CreateRequest{
        Name:        "web-server-01",
        FlavorID:    "flavor-standard-4cpu-8gb",
        ImageID:     "ubuntu-22.04-lts",
        KeypairName: "my-ssh-key",
        Networks: []servers.NetworkAttachment{
            {NetworkID: "net-private-123"},
        },
        SecurityGroups: []string{"sg-web-tier"},
        UserData:       "#!/bin/bash\napt-get update\napt-get install -y nginx",
    }
    
    server, err := vps.Servers().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    
    fmt.Printf("Created server: %s (ID: %s, Status: %s)\n",
        server.Name, server.ID, server.Status)
}
```

### Wait for Server to Become Active

```go
import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/Zillaforge/cloud-sdk/modules/vps"
)

func waitForServerActive(ctx context.Context, vpsClient *vps.Client, serverID string) {
    // Wait up to 5 minutes for server to become active
    waitCtx, cancel := context.WithTimeout(ctx, 5*time.Minute)
    defer cancel()
    
    server, err := vps.WaitForServerActive(waitCtx, vpsClient, serverID)
    if err != nil {
        log.Fatalf("Server failed to become active: %v", err)
    }
    
    fmt.Printf("Server %s is now ACTIVE\n", server.ID)
}
```

### Perform Server Actions

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

func serverActions(ctx context.Context, vps *vps.Client, serverID string) {
    // Stop server
    err := vps.Servers().Action(ctx, serverID, &servers.ActionRequest{
        Action: "stop",
    })
    if err != nil {
        log.Fatalf("Failed to stop server: %v", err)
    }
    fmt.Println("Server stopped")
    
    // Wait for SHUTOFF status
    _, err = vps.WaitForServerShutoff(ctx, vps, serverID)
    if err != nil {
        log.Fatalf("Failed waiting for SHUTOFF: %v", err)
    }
    
    // Start server
    err = vps.Servers().Action(ctx, serverID, &servers.ActionRequest{
        Action: "start",
    })
    if err != nil {
        log.Fatalf("Failed to start server: %v", err)
    }
    fmt.Println("Server started")
    
    // Reboot server
    err = vps.Servers().Action(ctx, serverID, &servers.ActionRequest{
        Action: "reboot",
        Force:  false, // Graceful reboot
    })
    if err != nil {
        log.Fatalf("Failed to reboot server: %v", err)
    }
    fmt.Println("Server rebooting")
    
    // Resize server
    err = vps.Servers().Action(ctx, serverID, &servers.ActionRequest{
        Action:   "resize",
        FlavorID: "flavor-standard-8cpu-16gb",
    })
    if err != nil {
        log.Fatalf("Failed to resize server: %v", err)
    }
    fmt.Println("Server resized")
}
```

### Manage Server NICs

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

func manageServerNICs(ctx context.Context, vps *vps.Client, serverID string) {
    serverRes := vps.Servers().Resource(serverID)
    
    // List existing NICs
    nics, err := serverRes.NICs().List(ctx)
    if err != nil {
        log.Fatalf("Failed to list NICs: %v", err)
    }
    fmt.Printf("Server has %d NICs\n", len(nics.NICs))
    
    // Add a new NIC
    nic, err := serverRes.NICs().Add(ctx, &servers.AddNICRequest{
        NetworkID: "net-private-456",
    })
    if err != nil {
        log.Fatalf("Failed to add NIC: %v", err)
    }
    fmt.Printf("Added NIC: %s (IP: %v)\n", nic.ID, nic.IPAddresses)
    
    // Associate floating IP with NIC
    err = serverRes.NICs().AssociateFloatingIP(ctx, nic.ID, &servers.AssociateFloatingIPRequest{
        FloatingIPID: "fip-123",
    })
    if err != nil {
        log.Fatalf("Failed to associate floating IP: %v", err)
    }
    fmt.Println("Floating IP associated")
    
    // Update NIC
    updated, err := serverRes.NICs().Update(ctx, nic.ID, &servers.UpdateNICRequest{
        SecurityGroups: []string{"sg-web-tier", "sg-ssh-access"},
    })
    if err != nil {
        log.Fatalf("Failed to update NIC: %v", err)
    }
    fmt.Printf("Updated NIC security groups: %v\n", updated.SecurityGroups)
    
    // Delete NIC
    err = serverRes.NICs().Delete(ctx, nic.ID)
    if err != nil {
        log.Fatalf("Failed to delete NIC: %v", err)
    }
    fmt.Println("NIC deleted")
}
```

### Manage Server Volumes

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

func manageServerVolumes(ctx context.Context, vps *vps.Client, serverID string) {
    serverRes := vps.Servers().Resource(serverID)
    
    // List attached volumes
    volumes, err := serverRes.Volumes().List(ctx)
    if err != nil {
        log.Fatalf("Failed to list volumes: %v", err)
    }
    fmt.Printf("Server has %d volumes\n", len(volumes.Volumes))
    
    // Attach a volume
    volume, err := serverRes.Volumes().Attach(ctx, &servers.AttachVolumeRequest{
        VolumeID: "vol-data-001",
        Device:   "/dev/vdb",
    })
    if err != nil {
        log.Fatalf("Failed to attach volume: %v", err)
    }
    fmt.Printf("Attached volume: %s (%dGB) at %s\n",
        volume.ID, volume.Size, volume.Device)
    
    // Detach volume
    err = serverRes.Volumes().Detach(ctx, volume.ID)
    if err != nil {
        log.Fatalf("Failed to detach volume: %v", err)
    }
    fmt.Println("Volume detached")
}
```

### Get Server Metrics

```go
import (
    "context"
    "fmt"
    "log"
    "time"
)

func getServerMetrics(ctx context.Context, vps *vps.Client, serverID string) {
    // Get metrics for the last hour
    endTime := time.Now()
    startTime := endTime.Add(-1 * time.Hour)
    
    metrics, err := vps.Servers().Metrics(ctx, serverID, startTime, endTime)
    if err != nil {
        log.Fatalf("Failed to get metrics: %v", err)
    }
    
    fmt.Printf("Server Metrics (%s to %s):\n", startTime.Format(time.RFC3339), endTime.Format(time.RFC3339))
    fmt.Printf("  CPU Usage: %.2f%%\n", metrics.CPU.Average)
    fmt.Printf("  Memory Usage: %.2fMB (%.2f%%)\n", metrics.Memory.UsedMB, metrics.Memory.UsagePercent)
    fmt.Printf("  Disk Read: %.2fMB\n", metrics.Disk.ReadMB)
    fmt.Printf("  Disk Write: %.2fMB\n", metrics.Disk.WriteMB)
    fmt.Printf("  Network In: %.2fMB\n", metrics.Network.InMB)
    fmt.Printf("  Network Out: %.2fMB\n", metrics.Network.OutMB)
}
```

### Get VNC Console URL

```go
import (
    "context"
    "fmt"
    "log"
)

func getVNCConsole(ctx context.Context, vps *vps.Client, serverID string) {
    vncURL, err := vps.Servers().VNCURL(ctx, serverID)
    if err != nil {
        log.Fatalf("Failed to get VNC URL: %v", err)
    }
    
    fmt.Printf("VNC Console URL: %s\n", vncURL.URL)
    fmt.Printf("Access Type: %s\n", vncURL.Type)
    fmt.Println("Open this URL in a browser to access the console")
}
```

### Complete Server Lifecycle

```go
func completeServerLifecycle(ctx context.Context, vps *vps.Client) {
    // 1. Create server
    server, err := vps.Servers().Create(ctx, &servers.CreateRequest{
        Name:     "lifecycle-demo",
        FlavorID: "flavor-small",
        ImageID:  "ubuntu-22.04",
    })
    if err != nil {
        log.Fatalf("Create failed: %v", err)
    }
    serverID := server.ID
    fmt.Printf("Created server: %s\n", serverID)
    
    // 2. Wait for ACTIVE
    _, err = vps.WaitForServerActive(ctx, vps, serverID)
    if err != nil {
        log.Fatalf("Wait failed: %v", err)
    }
    fmt.Println("Server is ACTIVE")
    
    // 3. Update server
    updated, err := vps.Servers().Update(ctx, serverID, &servers.UpdateRequest{
        Name: "lifecycle-demo-renamed",
    })
    if err != nil {
        log.Fatalf("Update failed: %v", err)
    }
    fmt.Printf("Updated server name: %s\n", updated.Name)
    
    // 4. Stop server
    err = vps.Servers().Action(ctx, serverID, &servers.ActionRequest{Action: "stop"})
    if err != nil {
        log.Fatalf("Stop failed: %v", err)
    }
    _, err = vps.WaitForServerShutoff(ctx, vps, serverID)
    if err != nil {
        log.Fatalf("Wait for SHUTOFF failed: %v", err)
    }
    fmt.Println("Server is SHUTOFF")
    
    // 5. Delete server
    err = vps.Servers().Delete(ctx, serverID)
    if err != nil {
        log.Fatalf("Delete failed: %v", err)
    }
    fmt.Println("Server deleted")
}
```

---

## Network Management

### Create a Network

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

func createNetwork(ctx context.Context, vps *vps.Client) {
    req := &networks.CreateRequest{
        Name:        "private-network-web",
        CIDR:        "10.0.1.0/24",
        Description: "Private network for web tier",
        DHCPEnabled: true,
    }
    
    network, err := vps.Networks().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to create network: %v", err)
    }
    
    fmt.Printf("Created network: %s (ID: %s, CIDR: %s)\n",
        network.Name, network.ID, network.CIDR)
}
```

### List Networks with Filters

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

func listNetworks(ctx context.Context, vps *vps.Client) {
    // List all networks
    allNetworks, err := vps.Networks().List(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list networks: %v", err)
    }
    fmt.Printf("Total networks: %d\n", len(allNetworks.Networks))
    
    // Filter by name
    filtered, err := vps.Networks().List(ctx, &networks.ListOptions{
        Name: "private-network",
    })
    if err != nil {
        log.Fatalf("Failed to list filtered networks: %v", err)
    }
    
    for _, net := range filtered.Networks {
        fmt.Printf("  - %s (CIDR: %s, Status: %s)\n",
            net.Name, net.CIDR, net.Status)
    }
}
```

### List Network Ports

```go
import (
    "context"
    "fmt"
    "log"
)

func listNetworkPorts(ctx context.Context, vps *vps.Client, networkID string) {
    ports, err := vps.Networks().Resource(networkID).Ports().List(ctx)
    if err != nil {
        log.Fatalf("Failed to list ports: %v", err)
    }
    
    fmt.Printf("Network has %d ports:\n", len(ports.Ports))
    for _, port := range ports.Ports {
        fmt.Printf("  - Port %s: %v (Server: %s)\n",
            port.ID, port.IPAddresses, port.ServerID)
    }
}
```

### Complete Network Lifecycle

```go
func completeNetworkLifecycle(ctx context.Context, vps *vps.Client) {
    // Create
    network, err := vps.Networks().Create(ctx, &networks.CreateRequest{
        Name: "test-network",
        CIDR: "192.168.100.0/24",
    })
    if err != nil {
        log.Fatalf("Create failed: %v", err)
    }
    networkID := network.ID
    fmt.Printf("Created network: %s\n", networkID)
    
    // Get
    retrieved, err := vps.Networks().Get(ctx, networkID)
    if err != nil {
        log.Fatalf("Get failed: %v", err)
    }
    fmt.Printf("Retrieved network: %s (CIDR: %s)\n", retrieved.Name, retrieved.CIDR)
    
    // Update
    updated, err := vps.Networks().Update(ctx, networkID, &networks.UpdateRequest{
        Name:        "test-network-updated",
        Description: "Updated description",
    })
    if err != nil {
        log.Fatalf("Update failed: %v", err)
    }
    fmt.Printf("Updated network: %s\n", updated.Name)
    
    // List ports
    ports, err := vps.Networks().Resource(networkID).Ports().List(ctx)
    if err != nil {
        log.Fatalf("List ports failed: %v", err)
    }
    fmt.Printf("Network has %d ports\n", len(ports.Ports))
    
    // Delete
    err = vps.Networks().Delete(ctx, networkID)
    if err != nil {
        log.Fatalf("Delete failed: %v", err)
    }
    fmt.Println("Network deleted")
}
```

---

## Floating IP Management

### Create and Approve Floating IP

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
)

func createFloatingIP(ctx context.Context, vps *vps.Client) {
    // Create floating IP
    req := &floatingips.CreateRequest{
        Description: "Production web server public IP",
    }
    
    fip, err := vps.FloatingIPs().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to create floating IP: %v", err)
    }
    
    fmt.Printf("Created floating IP: %s (Status: %s)\n", fip.ID, fip.Status)
    
    // If status is PENDING, approve it
    if fip.Status == "PENDING" {
        err = vps.FloatingIPs().Approve(ctx, fip.ID)
        if err != nil {
            log.Fatalf("Failed to approve floating IP: %v", err)
        }
        fmt.Println("Floating IP approved")
        
        // Get updated status
        updated, err := vps.FloatingIPs().Get(ctx, fip.ID)
        if err != nil {
            log.Fatalf("Failed to get floating IP: %v", err)
        }
        fmt.Printf("Floating IP now ACTIVE: %s\n", updated.FloatingIP)
    }
}
```

### Associate Floating IP with Server

```go
func associateFloatingIP(ctx context.Context, vps *vps.Client, serverID, floatingIPID string) {
    // Get server's first NIC
    nics, err := vps.Servers().Resource(serverID).NICs().List(ctx)
    if err != nil {
        log.Fatalf("Failed to list NICs: %v", err)
    }
    if len(nics.NICs) == 0 {
        log.Fatal("Server has no NICs")
    }
    
    nicID := nics.NICs[0].ID
    
    // Associate floating IP
    err = vps.Servers().Resource(serverID).NICs().AssociateFloatingIP(
        ctx, nicID, &servers.AssociateFloatingIPRequest{
            FloatingIPID: floatingIPID,
        },
    )
    if err != nil {
        log.Fatalf("Failed to associate floating IP: %v", err)
    }
    
    fmt.Printf("Associated floating IP %s with server NIC %s\n", floatingIPID, nicID)
}
```

### Disassociate Floating IP

```go
func disassociateFloatingIP(ctx context.Context, vps *vps.Client, floatingIPID string) {
    err := vps.FloatingIPs().Disassociate(ctx, floatingIPID)
    if err != nil {
        log.Fatalf("Failed to disassociate floating IP: %v", err)
    }
    fmt.Printf("Floating IP %s disassociated\n", floatingIPID)
}
```

### Complete Floating IP Lifecycle

```go
func completeFloatingIPLifecycle(ctx context.Context, vps *vps.Client) {
    // Create
    fip, err := vps.FloatingIPs().Create(ctx, &floatingips.CreateRequest{
        Description: "Test floating IP",
    })
    if err != nil {
        log.Fatalf("Create failed: %v", err)
    }
    fipID := fip.ID
    fmt.Printf("Created floating IP: %s\n", fipID)
    
    // Approve if pending
    if fip.Status == "PENDING" {
        err = vps.FloatingIPs().Approve(ctx, fipID)
        if err != nil {
            log.Fatalf("Approve failed: %v", err)
        }
        fmt.Println("Approved")
    }
    
    // Update
    updated, err := vps.FloatingIPs().Update(ctx, fipID, &floatingips.UpdateRequest{
        Description: "Updated description",
    })
    if err != nil {
        log.Fatalf("Update failed: %v", err)
    }
    fmt.Printf("Updated: %s\n", updated.Description)
    
    // Disassociate (if associated)
    err = vps.FloatingIPs().Disassociate(ctx, fipID)
    if err != nil {
        // Ignore error if not associated
        fmt.Printf("Disassociate warning: %v\n", err)
    }
    
    // Delete
    err = vps.FloatingIPs().Delete(ctx, fipID)
    if err != nil {
        log.Fatalf("Delete failed: %v", err)
    }
    fmt.Println("Deleted")
}
```

---

## SSH Keypair Management

### Create Keypair with Public Key

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/keypairs"
)

func createKeypair(ctx context.Context, vps *vps.Client) {
    publicKey := `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABgQC... user@host`
    
    req := &keypairs.CreateRequest{
        Name:      "deploy-key-prod",
        PublicKey: publicKey,
    }
    
    keypair, err := vps.Keypairs().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to create keypair: %v", err)
    }
    
    fmt.Printf("Created keypair: %s (Fingerprint: %s)\n",
        keypair.Name, keypair.Fingerprint)
}
```

### Generate Keypair

```go
func generateKeypair(ctx context.Context, vps *vps.Client) {
    // Leave PublicKey empty to auto-generate
    req := &keypairs.CreateRequest{
        Name: "auto-generated-key",
    }
    
    keypair, err := vps.Keypairs().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to generate keypair: %v", err)
    }
    
    fmt.Printf("Generated keypair: %s\n", keypair.Name)
    fmt.Printf("Private Key:\n%s\n", keypair.PrivateKey)
    fmt.Println("SAVE THIS PRIVATE KEY - it won't be shown again!")
}
```

### List and Filter Keypairs

```go
func listKeypairs(ctx context.Context, vps *vps.Client) {
    // List all
    all, err := vps.Keypairs().List(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list keypairs: %v", err)
    }
    fmt.Printf("Total keypairs: %d\n", len(all.Keypairs))
    
    // Filter by name
    filtered, err := vps.Keypairs().List(ctx, &keypairs.ListOptions{
        Name: "deploy",
    })
    if err != nil {
        log.Fatalf("Failed to list filtered keypairs: %v", err)
    }
    
    for _, kp := range filtered.Keypairs {
        fmt.Printf("  - %s (Fingerprint: %s)\n", kp.Name, kp.Fingerprint)
    }
}
```

---

## Router Management

### Create Router with External Network

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/routers"
)

func createRouter(ctx context.Context, vps *vps.Client) {
    req := &routers.CreateRequest{
        Name:              "main-router",
        Description:       "Primary network router",
        ExternalNetworkID: "ext-net-public",
    }
    
    router, err := vps.Routers().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to create router: %v", err)
    }
    
    fmt.Printf("Created router: %s (ID: %s, Status: %s)\n",
        router.Name, router.ID, router.Status)
}
```

### Enable/Disable Router

```go
func setRouterState(ctx context.Context, vps *vps.Client, routerID string, enabled bool) {
    req := &routers.SetStateRequest{
        Enabled: enabled,
    }
    
    err := vps.Routers().SetState(ctx, routerID, req)
    if err != nil {
        log.Fatalf("Failed to set router state: %v", err)
    }
    
    state := "disabled"
    if enabled {
        state = "enabled"
    }
    fmt.Printf("Router %s\n", state)
}
```

### Associate Networks with Router

```go
func manageRouterNetworks(ctx context.Context, vps *vps.Client, routerID string) {
    routerRes := vps.Routers().Resource(routerID)
    
    // Associate a network
    err := routerRes.Networks().Associate(ctx, &routers.AssociateNetworkRequest{
        NetworkID: "net-private-web",
    })
    if err != nil {
        log.Fatalf("Failed to associate network: %v", err)
    }
    fmt.Println("Network associated")
    
    // List associated networks
    networks, err := routerRes.Networks().List(ctx)
    if err != nil {
        log.Fatalf("Failed to list networks: %v", err)
    }
    fmt.Printf("Router has %d networks:\n", len(networks.Networks))
    for _, net := range networks.Networks {
        fmt.Printf("  - %s (CIDR: %s)\n", net.Name, net.CIDR)
    }
    
    // Disassociate network
    err = routerRes.Networks().Disassociate(ctx, "net-private-web")
    if err != nil {
        log.Fatalf("Failed to disassociate network: %v", err)
    }
    fmt.Println("Network disassociated")
}
```

---

## Security Group Management

### Create Security Group with Rules

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

func createSecurityGroup(ctx context.Context, vps *vps.Client) {
    req := &securitygroups.CreateRequest{
        Name:        "web-tier-sg",
        Description: "Security group for web servers",
        Rules: []securitygroups.RuleRequest{
            // Allow HTTP from anywhere
            {
                Direction:  "ingress",
                Protocol:   "tcp",
                PortMin:    80,
                PortMax:    80,
                RemoteCIDR: "0.0.0.0/0",
            },
            // Allow HTTPS from anywhere
            {
                Direction:  "ingress",
                Protocol:   "tcp",
                PortMin:    443,
                PortMax:    443,
                RemoteCIDR: "0.0.0.0/0",
            },
            // Allow SSH from office network only
            {
                Direction:  "ingress",
                Protocol:   "tcp",
                PortMin:    22,
                PortMax:    22,
                RemoteCIDR: "203.0.113.0/24",
            },
            // Allow all outbound traffic
            {
                Direction:  "egress",
                Protocol:   "tcp",
                PortMin:    1,
                PortMax:    65535,
                RemoteCIDR: "0.0.0.0/0",
            },
        },
    }
    
    sg, err := vps.SecurityGroups().Create(ctx, req)
    if err != nil {
        log.Fatalf("Failed to create security group: %v", err)
    }
    
    fmt.Printf("Created security group: %s (ID: %s)\n", sg.Name, sg.ID)
    fmt.Printf("Rules: %d\n", len(sg.Rules))
}
```

### List Security Groups

```go
func listSecurityGroups(ctx context.Context, vps *vps.Client) {
    sgs, err := vps.SecurityGroups().List(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list security groups: %v", err)
    }
    
    fmt.Printf("Security Groups: %d\n", len(sgs.SecurityGroups))
    for _, sg := range sgs.SecurityGroups {
        fmt.Printf("  - %s: %d rules\n", sg.Name, len(sg.Rules))
    }
}
```

---

## Flavor Discovery

### List All Flavors

```go
import (
    "context"
    "fmt"
    "log"
)

func listFlavors(ctx context.Context, vps *vps.Client) {
    flavors, err := vps.Flavors().List(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list flavors: %v", err)
    }
    
    fmt.Printf("Available Flavors: %d\n", len(flavors.Flavors))
    for _, flavor := range flavors.Flavors {
        fmt.Printf("  - %s: %d vCPUs, %dMB RAM, %dGB disk\n",
            flavor.Name, flavor.VCPUs, flavor.RAM, flavor.Disk)
    }
}
```

### Find Flavor by Characteristics

```go
import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk/models/vps/flavors"
)

func findFlavor(ctx context.Context, vps *vps.Client, minVCPUs, minRAM int) *flavors.Flavor {
    allFlavors, err := vps.Flavors().List(ctx, nil)
    if err != nil {
        log.Fatalf("Failed to list flavors: %v", err)
    }
    
    for _, flavor := range allFlavors.Flavors {
        if flavor.VCPUs >= minVCPUs && flavor.RAM >= minRAM {
            fmt.Printf("Found suitable flavor: %s (%d vCPUs, %dMB RAM)\n",
                flavor.Name, flavor.VCPUs, flavor.RAM)
            return &flavor
        }
    }
    
    return nil
}
```

### Filter Public Flavors

```go
func listPublicFlavors(ctx context.Context, vps *vps.Client) {
    publicFlavors, err := vps.Flavors().List(ctx, &flavors.ListOptions{
        Public: true,
    })
    if err != nil {
        log.Fatalf("Failed to list public flavors: %v", err)
    }
    
    fmt.Printf("Public Flavors: %d\n", len(publicFlavors.Flavors))
    for _, flavor := range publicFlavors.Flavors {
        fmt.Printf("  - %s (public)\n", flavor.Name)
    }
}
```

---

## Quota Management

### View Current Quotas

```go
import (
    "context"
    "fmt"
    "log"
)

func viewQuotas(ctx context.Context, vps *vps.Client) {
    quotas, err := vps.Quotas().Get(ctx)
    if err != nil {
        log.Fatalf("Failed to get quotas: %v", err)
    }
    
    fmt.Println("Project Quotas:")
    fmt.Printf("  VMs:         %d / %d\n", quotas.VMs.Used, quotas.VMs.Limit)
    fmt.Printf("  vCPUs:       %d / %d\n", quotas.VCPUs.Used, quotas.VCPUs.Limit)
    fmt.Printf("  RAM:         %d / %d MB\n", quotas.RAM.Used, quotas.RAM.Limit)
    fmt.Printf("  Storage:     %d / %d GB\n", quotas.Storage.Used, quotas.Storage.Limit)
    fmt.Printf("  Networks:    %d / %d\n", quotas.Networks.Used, quotas.Networks.Limit)
    fmt.Printf("  Floating IPs: %d / %d\n", quotas.FloatingIPs.Used, quotas.FloatingIPs.Limit)
    fmt.Printf("  Routers:     %d / %d\n", quotas.Routers.Used, quotas.Routers.Limit)
    
    if quotas.GPUs != nil {
        fmt.Printf("  GPUs:        %d / %d\n", quotas.GPUs.Used, quotas.GPUs.Limit)
    }
}
```

### Check Capacity Before Creation

```go
func checkCapacity(ctx context.Context, vps *vps.Client, flavorID string) bool {
    // Get quotas
    quotas, err := vps.Quotas().Get(ctx)
    if err != nil {
        log.Fatalf("Failed to get quotas: %v", err)
    }
    
    // Get flavor details
    flavor, err := vps.Flavors().Get(ctx, flavorID)
    if err != nil {
        log.Fatalf("Failed to get flavor: %v", err)
    }
    
    // Check if we have capacity
    hasVMs := quotas.VMs.Used < quotas.VMs.Limit
    hasVCPUs := (quotas.VCPUs.Used + flavor.VCPUs) <= quotas.VCPUs.Limit
    hasRAM := (quotas.RAM.Used + flavor.RAM) <= quotas.RAM.Limit
    
    canCreate := hasVMs && hasVCPUs && hasRAM
    
    if !canCreate {
        fmt.Println("Insufficient capacity:")
        if !hasVMs {
            fmt.Printf("  VM limit reached: %d/%d\n", quotas.VMs.Used, quotas.VMs.Limit)
        }
        if !hasVCPUs {
            fmt.Printf("  vCPU limit would be exceeded: %d+%d > %d\n",
                quotas.VCPUs.Used, flavor.VCPUs, quotas.VCPUs.Limit)
        }
        if !hasRAM {
            fmt.Printf("  RAM limit would be exceeded: %d+%d > %d\n",
                quotas.RAM.Used, flavor.RAM, quotas.RAM.Limit)
        }
    }
    
    return canCreate
}
```

---

## Error Handling

### Structured Error Handling

```go
import (
    "context"
    "fmt"
    "log"
    
    cloudsdk "github.com/Zillaforge/cloud-sdk"
)

func handleErrors(ctx context.Context, vps *vps.Client) {
    server, err := vps.Servers().Get(ctx, "non-existent-id")
    if err != nil {
        // Type assert to SDK error
        if sdkErr, ok := err.(*cloudsdk.SDKError); ok {
            fmt.Printf("Error: %s\n", sdkErr.Message)
            fmt.Printf("Status Code: %d\n", sdkErr.StatusCode)
            
            // Handle specific error codes
            switch sdkErr.StatusCode {
            case 401:
                log.Fatal("Authentication failed - check your token")
            case 403:
                log.Fatal("Permission denied - insufficient privileges")
            case 404:
                fmt.Println("Server not found")
                return
            case 409:
                fmt.Println("Conflict - resource already exists or quota exceeded")
                return
            case 429:
                fmt.Println("Rate limited - automatic retry will occur")
            case 500, 502, 503, 504:
                fmt.Println("Server error - automatic retry will occur")
            default:
                log.Fatalf("Unexpected error: %v", err)
            }
        } else {
            // Generic error
            log.Fatalf("Request failed: %v", err)
        }
    }
    
    if server != nil {
        fmt.Printf("Server: %s\n", server.Name)
    }
}
```

### Context Timeout Handling

```go
import (
    "context"
    "fmt"
    "log"
    "time"
)

func handleTimeout(vps *vps.Client) {
    // Set aggressive timeout
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel()
    
    _, err := vps.Servers().List(ctx, nil)
    if err != nil {
        if ctx.Err() == context.DeadlineExceeded {
            fmt.Println("Request timed out - consider increasing timeout")
        } else {
            log.Fatalf("Request failed: %v", err)
        }
    }
}
```

---

## Advanced Patterns

### Parallel Resource Creation

```go
import (
    "context"
    "fmt"
    "sync"
)

func createMultipleServers(ctx context.Context, vps *vps.Client, count int) {
    var wg sync.WaitGroup
    results := make(chan string, count)
    errors := make(chan error, count)
    
    for i := 0; i < count; i++ {
        wg.Add(1)
        go func(index int) {
            defer wg.Done()
            
            server, err := vps.Servers().Create(ctx, &servers.CreateRequest{
                Name:     fmt.Sprintf("server-%d", index),
                FlavorID: "flavor-small",
                ImageID:  "ubuntu-22.04",
            })
            
            if err != nil {
                errors <- err
                return
            }
            
            results <- server.ID
        }(i)
    }
    
    wg.Wait()
    close(results)
    close(errors)
    
    // Process results
    fmt.Println("Created servers:")
    for id := range results {
        fmt.Printf("  - %s\n", id)
    }
    
    // Check errors
    for err := range errors {
        fmt.Printf("Error: %v\n", err)
    }
}
```

### Resource Cleanup with Defer

```go
func createTemporaryResources(ctx context.Context, vps *vps.Client) {
    // Create network
    network, err := vps.Networks().Create(ctx, &networks.CreateRequest{
        Name: "temp-network",
        CIDR: "192.168.1.0/24",
    })
    if err != nil {
        log.Fatalf("Failed to create network: %v", err)
    }
    defer func() {
        err := vps.Networks().Delete(ctx, network.ID)
        if err != nil {
            fmt.Printf("Failed to cleanup network: %v\n", err)
        }
    }()
    
    // Create server
    server, err := vps.Servers().Create(ctx, &servers.CreateRequest{
        Name:     "temp-server",
        FlavorID: "flavor-small",
        ImageID:  "ubuntu-22.04",
        Networks: []servers.NetworkAttachment{{NetworkID: network.ID}},
    })
    if err != nil {
        log.Fatalf("Failed to create server: %v", err)
    }
    defer func() {
        err := vps.Servers().Delete(ctx, server.ID)
        if err != nil {
            fmt.Printf("Failed to cleanup server: %v\n", err)
        }
    }()
    
    // Do work with temporary resources...
    fmt.Println("Using temporary resources...")
    
    // Cleanup happens automatically via defer
}
```

### Polling with Custom Waiter

```go
import (
    "context"
    "fmt"
    "time"
    
    "github.com/Zillaforge/cloud-sdk/internal/waiter"
)

func customWaiter(ctx context.Context, vps *vps.Client, serverID string) error {
    condition := func(ctx context.Context) (bool, error) {
        server, err := vps.Servers().Get(ctx, serverID)
        if err != nil {
            return false, err
        }
        
        fmt.Printf("Server status: %s\n", server.Status)
        return server.Status == "ACTIVE", nil
    }
    
    return waiter.Wait(ctx, condition, &waiter.Options{
        Interval:    5 * time.Second,
        MaxDuration: 10 * time.Minute,
    })
}
```

---

For more information, see the [README.md](./README.md) and [API documentation](https://godoc.org/github.com/Zillaforge/cloud-sdk/modules/vps).
