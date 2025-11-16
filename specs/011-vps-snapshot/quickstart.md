# Quickstart: Using the vps snapshot client

Example: Create a snapshot from a volume, then poll until available.

1. Create the client

 - Initialize `internalhttp.Client` with baseURL and auth (see `modules/vps/core/client.go`).
 - Create a snapshot client: `snapClient := snapshots.NewClient(baseClient, projectID)`

2. Create a snapshot (Go example)

```go
req := &snapshots.CreateSnapshotRequest{
    Name: "backup-2025-11-16",
    VolumeID: "vol-123",
}

snap, err := snapClient.Create(ctx, req)
// snap.Status == snapshots.SnapshotStatusCreating

// Poll until available
for { 
    s, _ := snapClient.Get(ctx, snap.ID)
    if s.Status == snapshots.SnapshotStatusAvailable { break }
    time.Sleep(1 * time.Second)
}
```

3. List snapshots and access by index

```go
list, _ := snapClient.List(ctx, nil)
first := list[0]
```
