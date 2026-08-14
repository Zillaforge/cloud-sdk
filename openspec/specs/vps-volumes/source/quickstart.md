# Volumes API Quickstart

Quick start guide for using the VPS Volumes and VolumeTypes API client.

---

## Installation

```bash
go get github.com/your-org/cloud-sdk
```

---

## Basic Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/your-org/cloud-sdk"
    "github.com/your-org/cloud-sdk/modules/vps"
    "github.com/your-org/cloud-sdk/modules/vps/volumes"
    "github.com/your-org/cloud-sdk/modules/vps/volumetypes"
)

func main() {
    // Initialize SDK client
    sdkClient := cloudsdk.NewClient(
        "https://api.example.com",
        "your-bearer-token",
        nil, // optional logger
    )
    
    // Create VPS client
    vpsClient := vps.NewClient(sdkClient, "project-123")
    
    // Access volumes client
    volumesClient := vpsClient.Volumes()
    
    // Access volume types client
    volumeTypesClient := vpsClient.VolumeTypes()
    
    ctx := context.Background()
    
    // Your code here...
}
```

---

## 1. List Available Volume Types

Before creating volumes, check what types are available:

```go
func listVolumeTypes(ctx context.Context, client *volumetypes.Client) error {
    types, err := client.List(ctx)
    if err != nil {
        return fmt.Errorf("failed to list volume types: %w", err)
    }
    
    fmt.Println("Available volume types:")
    for _, volumeType := range types {
        fmt.Printf("  - %s\n", volumeType)
    }
    
    return nil
}
```

**Expected Output:**
```
Available volume types:
  - SSD
  - HDD
  - NVMe
```

---

## 2. Create a Volume

Create a new volume with specified type and size:

```go
func createVolume(ctx context.Context, client *volumes.Client) (*volumes.Volume, error) {
    request := &volumes.CreateVolumeRequest{
        Name:        "my-data-volume",
        Type:        "SSD",
        Size:        100, // GB
        Description: "Application data storage",
    }
    
    volume, err := client.Create(ctx, request)
    if err != nil {
        return nil, fmt.Errorf("failed to create volume: %w", err)
    }
    
    fmt.Printf("Created volume: %s (ID: %s)\n", volume.Name, volume.ID)
    fmt.Printf("  Status: %s\n", volume.Status)
    fmt.Printf("  Size: %d GB\n", volume.Size)
    
    return volume, nil
}
```

---

## 3. Create Volume from Snapshot

Restore a volume from an existing snapshot:

```go
func createVolumeFromSnapshot(ctx context.Context, client *volumes.Client) (*volumes.Volume, error) {
    request := &volumes.CreateVolumeRequest{
        Name:       "restored-volume",
        Type:       "SSD",
        SnapshotID: "snap-abc123",
    }
    
    volume, err := client.Create(ctx, request)
    if err != nil {
        return nil, fmt.Errorf("failed to create volume from snapshot: %w", err)
    }
    
    fmt.Printf("Restored volume %s from snapshot\n", volume.ID)
    return volume, nil
}
```

---

## 4. List Volumes with Filters

Query volumes with various filters:

```go
func listVolumes(ctx context.Context, client *volumes.Client) error {
    // List all volumes
    allVolumes, err := client.List(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to list volumes: %w", err)
    }
    fmt.Printf("Total volumes: %d\n", len(allVolumes))
    
    // Filter by status
    opts := &volumes.ListVolumesOptions{
        Status: "available",
    }
    availableVolumes, err := client.List(ctx, opts)
    if err != nil {
        return fmt.Errorf("failed to list available volumes: %w", err)
    }
    fmt.Printf("Available volumes: %d\n", len(availableVolumes))
    
    // Filter by type with details
    optsDetailed := &volumes.ListVolumesOptions{
        Type:   "SSD",
        Detail: true,
    }
    ssdVolumes, err := client.List(ctx, optsDetailed)
    if err != nil {
        return fmt.Errorf("failed to list SSD volumes: %w", err)
    }
    
    for _, vol := range ssdVolumes {
        fmt.Printf("Volume: %s (%s)\n", vol.Name, vol.ID)
        fmt.Printf("  Type: %s, Size: %d GB, Status: %s\n", 
            vol.Type, vol.Size, vol.Status)
        if len(vol.Attachments) > 0 {
            fmt.Printf("  Attached to: %s\n", vol.Attachments[0].Name)
        }
    }
    
    return nil
}
```

---

## 5. Get Volume Details

Retrieve detailed information about a specific volume:

```go
func getVolumeDetails(ctx context.Context, client *volumes.Client, volumeID string) error {
    volume, err := client.Get(ctx, volumeID)
    if err != nil {
        return fmt.Errorf("failed to get volume: %w", err)
    }
    
    fmt.Printf("Volume Details:\n")
    fmt.Printf("  ID: %s\n", volume.ID)
    fmt.Printf("  Name: %s\n", volume.Name)
    fmt.Printf("  Type: %s\n", volume.Type)
    fmt.Printf("  Size: %d GB\n", volume.Size)
    fmt.Printf("  Status: %s\n", volume.Status)
    fmt.Printf("  Project: %s (%s)\n", volume.Project.Name, volume.Project.ID)
    fmt.Printf("  Owner: %s (%s)\n", volume.User.Name, volume.User.ID)
    fmt.Printf("  Created: %s\n", volume.CreatedAt)
    fmt.Printf("  Updated: %s\n", volume.UpdatedAt)
    
    if len(volume.Attachments) > 0 {
        fmt.Printf("  Attachments:\n")
        for _, att := range volume.Attachments {
            fmt.Printf("    - %s (%s)\n", att.Name, att.ID)
        }
    }
    
    return nil
}
```

---

## 6. Update Volume

Update volume name or description:

```go
func updateVolume(ctx context.Context, client *volumes.Client, volumeID string) error {
    request := &volumes.UpdateVolumeRequest{
        Name:        "updated-volume-name",
        Description: "Updated description",
    }
    
    volume, err := client.Update(ctx, volumeID, request)
    if err != nil {
        return fmt.Errorf("failed to update volume: %w", err)
    }
    
    fmt.Printf("Updated volume: %s\n", volume.Name)
    return nil
}
```

---

## 7. Attach Volume to Server

Attach a volume to a running server:

```go
func attachVolume(ctx context.Context, client *volumes.Client, 
    volumeID, serverID string) error {
    
    request := &volumes.VolumeActionRequest{
        Action:   volumes.VolumeActionAttach,
        ServerID: serverID,
    }
    
    err := client.Action(ctx, volumeID, request)
    if err != nil {
        return fmt.Errorf("failed to attach volume: %w", err)
    }
    
    fmt.Printf("Volume %s attached to server %s\n", volumeID, serverID)
    return nil
}
```

---

## 8. Detach Volume from Server

Detach a volume before deletion or reattachment:

```go
func detachVolume(ctx context.Context, client *volumes.Client, 
    volumeID, serverID string) error {
    
    request := &volumes.VolumeActionRequest{
        Action:   volumes.VolumeActionDetach,
        ServerID: serverID,
    }
    
    err := client.Action(ctx, volumeID, request)
    if err != nil {
        return fmt.Errorf("failed to detach volume: %w", err)
    }
    
    fmt.Printf("Volume %s detached from server %s\n", volumeID, serverID)
    return nil
}
```

---

## 9. Extend Volume Size

Increase volume capacity (cannot be decreased):

```go
func extendVolume(ctx context.Context, client *volumes.Client, 
    volumeID string, newSize int) error {
    
    request := &volumes.VolumeActionRequest{
        Action:  volumes.VolumeActionExtend,
        NewSize: newSize,
    }
    
    err := client.Action(ctx, volumeID, request)
    if err != nil {
        return fmt.Errorf("failed to extend volume: %w", err)
    }
    
    fmt.Printf("Volume %s extended to %d GB\n", volumeID, newSize)
    return nil
}
```

---

## 10. Delete Volume

Delete a volume (must be detached first):

```go
func deleteVolume(ctx context.Context, client *volumes.Client, volumeID string) error {
    err := client.Delete(ctx, volumeID)
    if err != nil {
        return fmt.Errorf("failed to delete volume: %w", err)
    }
    
    fmt.Printf("Volume %s deleted\n", volumeID)
    return nil
}
```

---

## Complete Workflow Example

A complete lifecycle example:

```go
func completeVolumeLifecycle(ctx context.Context, vpsClient *vps.Client) error {
    volumesClient := vpsClient.Volumes()
    
    // 1. List available types
    typesClient := vpsClient.VolumeTypes()
    types, err := typesClient.List(ctx)
    if err != nil {
        return err
    }
    fmt.Printf("Available types: %v\n", types)
    
    // 2. Create volume
    createRequest := &volumes.CreateVolumeRequest{
        Name: "temp-volume",
        Type: "SSD",
        Size: 50,
    }
    volume, err := volumesClient.Create(ctx, createRequest)
    if err != nil {
        return err
    }
    volumeID := volume.ID
    fmt.Printf("Created volume: %s\n", volumeID)
    
    // 3. Update volume metadata
    updateRequest := &volumes.UpdateVolumeRequest{
        Description: "Temporary storage",
    }
    _, err = volumesClient.Update(ctx, volumeID, updateRequest)
    if err != nil {
        return err
    }
    fmt.Println("Updated volume metadata")
    
    // 4. Attach to server
    attachRequest := &volumes.VolumeActionRequest{
        Action:   volumes.VolumeActionAttach,
        ServerID: "server-123",
    }
    err = volumesClient.Action(ctx, volumeID, attachRequest)
    if err != nil {
        return err
    }
    fmt.Println("Attached volume to server")
    
    // 5. Extend volume
    extendRequest := &volumes.VolumeActionRequest{
        Action:  volumes.VolumeActionExtend,
        NewSize: 100,
    }
    err = volumesClient.Action(ctx, volumeID, extendRequest)
    if err != nil {
        return err
    }
    fmt.Println("Extended volume to 100GB")
    
    // 6. Detach from server
    detachRequest := &volumes.VolumeActionRequest{
        Action:   volumes.VolumeActionDetach,
        ServerID: "server-123",
    }
    err = volumesClient.Action(ctx, volumeID, detachRequest)
    if err != nil {
        return err
    }
    fmt.Println("Detached volume from server")
    
    // 7. Delete volume
    err = volumesClient.Delete(ctx, volumeID)
    if err != nil {
        return err
    }
    fmt.Println("Deleted volume")
    
    return nil
}
```

---

## Error Handling Best Practices

```go
func robustVolumeOperation(ctx context.Context, client *volumes.Client) {
    volume, err := client.Create(ctx, &volumes.CreateVolumeRequest{
        Name: "my-volume",
        Type: "SSD",
        Size: 100,
    })
    
    if err != nil {
        // Check for specific error types
        var sdkErr *types.SDKError
        if errors.As(err, &sdkErr) {
            switch sdkErr.Code {
            case "QUOTA_EXCEEDED":
                log.Printf("Quota exceeded, reduce size or delete old volumes")
            case "INVALID_TYPE":
                log.Printf("Invalid volume type specified")
            default:
                log.Printf("API error: %s - %s", sdkErr.Code, sdkErr.Message)
            }
        } else {
            log.Printf("Unexpected error: %v", err)
        }
        return
    }
    
    log.Printf("Successfully created volume: %s", volume.ID)
}
```

---

## Testing Your Integration

```go
func TestVolumeIntegration(t *testing.T) {
    // Use test credentials
    client := setupTestClient(t)
    ctx := context.Background()
    
    // Create volume
    volume, err := client.Volumes().Create(ctx, &volumes.CreateVolumeRequest{
        Name: "test-volume",
        Type: "SSD",
        Size: 10,
    })
    require.NoError(t, err)
    defer client.Volumes().Delete(ctx, volume.ID)
    
    // Verify volume exists
    retrieved, err := client.Volumes().Get(ctx, volume.ID)
    require.NoError(t, err)
    assert.Equal(t, volume.ID, retrieved.ID)
    assert.Equal(t, volumes.VolumeStatusAvailable, retrieved.Status)
}
```

---

## Next Steps

- Review full API documentation in `EXAMPLES.md`
- Check contract tests in `contracts/volumes.md`
- See advanced patterns in `modules/vps/volumes/client_test.go`
- Explore waiters for async operations in `modules/vps/waiters.go`

---

## Support

For issues or questions:
- API Reference: swagger/vps.yaml
- SDK Documentation: README.md
- Examples: modules/vps/EXAMPLES.md
