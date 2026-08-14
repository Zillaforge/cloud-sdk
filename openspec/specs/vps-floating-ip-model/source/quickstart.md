# Quickstart: Floating IP Operations

**Feature**: 007-fix-floating-ip-model  
**Updated**: November 11, 2025

---

## Overview

The Floating IP module provides operations to manage floating IP addresses in your project. This guide demonstrates the six core operations: **List**, **Create**, **Get**, **Update**, **Delete**, and **Disassociate**.

---

## Setup

```go
import (
    "context"
    "log"
    "github.com/Zillaforge/cloud-sdk/modules/vps"
)

// Initialize VPS client with project
vpsClient := vps.NewClient(httpClient, "project-123")
fipClient := vpsClient.FloatingIPs // Access floating IPs submodule
```

---

## Operation 1: List Floating IPs

Retrieve all floating IPs in your project with optional filtering.

### Basic List

```go
ctx := context.Background()
response, err := fipClient.List(ctx, nil)
if err != nil {
    log.Fatalf("failed to list floatingips: %v", err)
}

// Response is []*FloatingIP with Status as FloatingIPStatus enum
for _, fip := range response {
    log.Printf("Floating IP: %s (%s)", fip.Name, fip.Address)
    
    // Status is type-safe enum
    switch fip.Status {
    case floatingips.FloatingIPStatusActive:
        log.Printf("  Status: Active")
    case floatingips.FloatingIPStatusPending:
        log.Printf("  Status: Pending approval")
    case floatingips.FloatingIPStatusDown:
        log.Printf("  Status: Service unavailable")
    case floatingips.FloatingIPStatusRejected:
        log.Printf("  Status: Rejected - %s", fip.StatusReason)
    }
    
    log.Printf("  Created: %s", fip.CreatedAt)
}
```

### List with Filters

```go
opts := &floatingips.ListFloatingIPsOptions{
    Status: "ACTIVE",
}
response, err := fipClient.List(ctx, opts)
if err != nil {
    log.Fatalf("failed to list floatingips: %v", err)
}

// Iterate through results
count := len(response)
log.Printf("Found %d active floating IPs", count)
```

### Key Points

- Response is a slice `[]*FloatingIP` (direct indexing possible)
- No wrapper object; removed "items" field from old API
- Filters are optional; omit `opts` for full list
- Status is a `FloatingIPStatus` enum type (type-safe constants)

### Status Values

FloatingIP status is represented as a custom enum with type-safe constants:

```go
// Available status constants
floatingips.FloatingIPStatusActive     // ACTIVE - Ready for use
floatingips.FloatingIPStatusPending    // PENDING - Awaiting approval
floatingips.FloatingIPStatusDown       // DOWN - Service issue
floatingips.FloatingIPStatusRejected   // REJECTED - Admin rejected

// Usage example
for _, fip := range response {
    if fip.Status == floatingips.FloatingIPStatusActive {
        log.Printf("IP %s is active", fip.Address)
    }
}

---

## Operation 2: Create a Floating IP

Allocate a new floating IP for your project.

### Simple Create

```go
req := &floatingips.FloatingIPCreateRequest{
    Name: "web-server-fip",
}

fip, err := fipClient.Create(ctx, req)
if err != nil {
    log.Fatalf("failed to create floatingip: %v", err)
}

log.Printf("Created Floating IP: %s", fip.Address)
log.Printf("  ID: %s", fip.ID)
log.Printf("  Status: %s", fip.Status)
log.Printf("  Created at: %s", fip.CreatedAt)
```

### Create with Description

```go
req := &floatingips.FloatingIPCreateRequest{
    Name:        "staging-fip",
    Description: "Floating IP for staging environment",
}

fip, err := fipClient.Create(ctx, req)
if err != nil {
    log.Fatalf("failed to create floatingip: %v", err)
}

log.Printf("Floating IP %s allocated successfully", fip.Name)
log.Printf("  Address: %s", fip.Address)
log.Printf("  Description: %s", fip.Description)
```

### Key Points

- Name and description are optional
- Floating IP starts in PENDING or ACTIVE status depending on approval settings
- Address is assigned by the system
- **Important**: If allocated with approval required, status will be PENDING until approved

---

## Operation 3: Get a Specific Floating IP

Retrieve details of a specific floating IP.

```go
fipID := "fip-abc123"

fip, err := fipClient.Get(ctx, fipID)
if err != nil {
    log.Fatalf("failed to get floatingip: %v", err)
}

log.Printf("Floating IP Details:")
log.Printf("  Name: %s", fip.Name)
log.Printf("  Address: %s", fip.Address)

// Status is now type-safe enum constant
switch fip.Status {
case floatingips.FloatingIPStatusActive:
    log.Printf("  Status: ✓ Active")
case floatingips.FloatingIPStatusPending:
    log.Printf("  Status: ⧗ Pending approval")
case floatingips.FloatingIPStatusDown:
    log.Printf("  Status: ✗ Down - %s", fip.StatusReason)
case floatingips.FloatingIPStatusRejected:
    log.Printf("  Status: ✗ Rejected - %s", fip.StatusReason)
}

log.Printf("  Created: %s", fip.CreatedAt)

// Check if associated with a device
if fip.DeviceID != "" {
    log.Printf("  Associated with: %s (%s)", fip.DeviceName, fip.DeviceType)
} else {
    log.Printf("  Not currently associated with any device")
}

// Check ownership references (may be nil if not populated by API)
if fip.User != nil {
    log.Printf("  Owner: %s (%s)", fip.User.Name, fip.User.ID)
}
```

### Key Points

- Returns full FloatingIP struct with all fields
- Optional fields (User, Project, Device details) may be null/empty
- Timestamps are RFC3339 strings

---

## Operation 4: Update a Floating IP

Modify name and description of an existing floating IP.

### Update Description

```go
fipID := "fip-abc123"

req := &floatingips.FloatingIPUpdateRequest{
    Description: "Updated description for web server",
}

fip, err := fipClient.Update(ctx, fipID, req)
if err != nil {
    log.Fatalf("failed to update floatingip: %v", err)
}

log.Printf("Updated Floating IP %s", fip.Name)
log.Printf("  New description: %s", fip.Description)
log.Printf("  Last updated: %s", fip.UpdatedAt)
```

### Update Name

```go
req := &floatingips.FloatingIPUpdateRequest{
    Name: "prod-api-fip",
}

fip, err := fipClient.Update(ctx, fipID, req)
if err != nil {
    log.Fatalf("failed to update floatingip: %v", err)
}

log.Printf("Renamed to: %s", fip.Name)
```

### Update Multiple Fields

```go
req := &floatingips.FloatingIPUpdateRequest{
    Name:        "primary-api",
    Description: "Primary API gateway floating IP",
}

fip, err := fipClient.Update(ctx, fipID, req)
if err != nil {
    log.Fatalf("failed to update floatingip: %v", err)
}

log.Printf("Updated: %s - %s", fip.Name, fip.Description)
```

### Important Notes

- Only Name and Description can be updated
- Reserved field is read-only (ignored if present in request)
- Status is controlled by the API, not by update requests
- UpdatedAt timestamp is set by the API

---

## Operation 5: Delete a Floating IP

Release and remove a floating IP from your project.

```go
fipID := "fip-abc123"

err := fipClient.Delete(ctx, fipID)
if err != nil {
    log.Fatalf("failed to delete floatingip: %v", err)
}

log.Printf("Floating IP %s deleted successfully", fipID)
```

### Important Notes

- Disassociate before delete if associated with a device
- Deletion is permanent
- Address is released back to the pool
- Error if floating IP is still in use

---

## Operation 6: Disassociate Floating IP

Detach a floating IP from its associated device or port.

```go
fipID := "fip-abc123"

err := fipClient.Disassociate(ctx, fipID)
if err != nil {
    log.Fatalf("failed to disassociate floatingip: %v", err)
}

log.Printf("Floating IP %s disassociated from device", fipID)

// Verify disassociation
fip, err := fipClient.Get(ctx, fipID)
if err == nil {
    log.Printf("  DeviceID: %s (should be empty)", fip.DeviceID)
    log.Printf("  PortID: %s (should be empty)", fip.PortID)
}
```

### Use Cases

- Move floating IP to a different device
- Clean up before deleting the floating IP
- Recover from accidental association

---

## Error Handling

All operations follow consistent error handling patterns:

```go
result, err := fipClient.List(ctx, nil)
if err != nil {
    // Error is wrapped with operation context
    // Example: "failed to list floatingips: {details}"
    log.Fatalf("Error: %v", err)
}

result, err := fipClient.Create(ctx, req)
if err != nil {
    // Same pattern for create
    log.Fatalf("Error: %v", err)
}

err := fipClient.Delete(ctx, fipID)
if err != nil {
    // Even single operations wrapped with context
    log.Fatalf("Error: %v", err)
}
```

### Common Errors

| Error | Cause | Resolution |
|-------|-------|-----------|
| "failed to list floatingips: 404" | Project ID invalid | Verify project ID |
| "failed to create floatingip: 400" | Invalid request (missing name, etc.) | Check request fields |
| "failed to get floatingips: 404" | Floating IP not found | Verify floating IP ID |
| "failed to delete floatingip: 409" | Cannot delete while associated | Call Disassociate first |
| "failed to update floatingip: 403" | Permission denied | Check user permissions |

---

## Complete Example Workflow

```go
package main

import (
    "context"
    "log"
    "github.com/Zillaforge/cloud-sdk/modules/vps"
    "github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
)

func main() {
    ctx := context.Background()
    
    // Initialize client
    vpsClient := vps.NewClient(httpClient, "project-123")
    fipClient := vpsClient.FloatingIPs
    
    // 1. List existing floating IPs
    log.Println("=== Listing Floating IPs ===")
    existing, err := fipClient.List(ctx, nil)
    if err != nil {
        log.Fatalf("List failed: %v", err)
    }
    log.Printf("Found %d floating IPs", len(existing))
    
    // 2. Create a new floating IP
    log.Println("\n=== Creating Floating IP ===")
    createReq := &floatingips.FloatingIPCreateRequest{
        Name:        "example-fip",
        Description: "Example floating IP for demonstration",
    }
    newFIP, err := fipClient.Create(ctx, createReq)
    if err != nil {
        log.Fatalf("Create failed: %v", err)
    }
    log.Printf("Created: %s (%s)", newFIP.Name, newFIP.Address)
    log.Printf("Status: %s", newFIP.Status)
    
    // 3. Get the floating IP
    log.Println("\n=== Getting Floating IP Details ===")
    fip, err := fipClient.Get(ctx, newFIP.ID)
    if err != nil {
        log.Fatalf("Get failed: %v", err)
    }
    log.Printf("Retrieved: %s", fip.Name)
    
    // 4. Update the floating IP
    log.Println("\n=== Updating Floating IP ===")
    updateReq := &floatingips.FloatingIPUpdateRequest{
        Description: "Updated description",
    }
    updated, err := fipClient.Update(ctx, fip.ID, updateReq)
    if err != nil {
        log.Fatalf("Update failed: %v", err)
    }
    log.Printf("Updated: %s", updated.Description)
    
    // 5. Disassociate if associated (for demonstration)
    if fip.DeviceID != "" {
        log.Println("\n=== Disassociating Floating IP ===")
        err := fipClient.Disassociate(ctx, fip.ID)
        if err != nil {
            log.Fatalf("Disassociate failed: %v", err)
        }
        log.Println("Disassociated successfully")
    }
    
    // 6. Delete the floating IP
    log.Println("\n=== Deleting Floating IP ===")
    err = fipClient.Delete(ctx, fip.ID)
    if err != nil {
        log.Fatalf("Delete failed: %v", err)
    }
    log.Println("Deleted successfully")
}
```

---

## Breaking Change Notice

### Migrating from Previous Version

The Floating IP model has been updated to match the VPS API specification exactly.

**Before (old version)**:
```go
response, _ := fipClient.List(ctx, nil)
count := response.Total                        // ← REMOVED
for _, fip := range response.Items {           // ← RENAMED to FloatingIPs
    // ...
}
```

**After (new version)**:
```go
response, _ := fipClient.List(ctx, nil)
count := len(response)                         // ← Use slice length
for _, fip := range response {                 // ← Direct slice iteration
    // ...
}
```

**Model Changes**:
- Response type is now `[]*FloatingIP` (slice) instead of wrapper struct
- Use `len(response)` to count items instead of `response.Total`
- Access fields directly: `response[i].Name`, `response[i].Address`, etc.

---

## Additional Resources

- **Data Model**: See `data-model.md` for field definitions and relationships
- **API Contract**: See `contracts/floatingip-openapi.yaml` for OpenAPI specification
- **Research**: See `research.md` for design decisions and alternatives considered
