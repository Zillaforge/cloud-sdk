# Data Model: vps snapshot resource

## Entities

- Snapshot
  - id: string  (pb SnapshotInfo.id)
  - name: string  (pb SnapshotInfo.name)
  - description: string (optional) (pb SnapshotInfo.description)
  - volume_id: string (pb SnapshotInfo.volume_id) — references a `volumes.Volume`
  - size: int (pb SnapshotInfo.size) — size in GB
  - status: SnapshotStatus (enum) — use custom type, do not use raw strings
  - status_reason: string (pb SnapshotInfo.status_reason)
  - project_id: string (pb SnapshotInfo.project_id)
  - user_id: string (pb SnapshotInfo.user_id)
  - namespace: string (pb SnapshotInfo.namespace) — the namespace/tenant scope for the snapshot
  - project: pb.IDName (pb SnapshotInfo.project) — project reference (id + name)
  - user: pb.IDName (pb SnapshotInfo.user) — user reference (id + name)
  - createdAt: string/time (pb SnapshotInfo.createdAt) — serialized as ISO 8601
  - updatedAt: string/time (pb SnapshotInfo.updatedAt)

## SnapshotStatus enumerations

Define a custom Go type:

type SnapshotStatus string

const (
    SnapshotStatusCreating SnapshotStatus = "creating"
    SnapshotStatusAvailable SnapshotStatus = "available"
    SnapshotStatusDeleting SnapshotStatus = "deleting"
    SnapshotStatusError SnapshotStatus = "error"
)

## Relationships

- Snapshot is created from Volume (volume_id). Volume is defined in `models/vps/volumes`.
- Snapshot is scoped to a Project (project_id) and has a user (user_id) for auditing.

## Validation rules

- name: required, non-empty
- volume_id: required, must exist in the same project
- size: positive integer

## Notes

- Use types and shared entities from `models/vps/common` when possible. For example, `common.IDName` for Project/User references.
- Naming: use Request/Response suffix naming for models (e.g., `CreateSnapshotRequest`, `SnapshotResponse`) instead of Input/Output.
- `ListSnapshots` returns `[]*Snapshot` (indexable) and should map JSON arrays from `pb.SnapshotListOutput` in swagger.

## Go struct example

Map snapshot fields directly to the Swagger field names using Go exported fields and json tags. Example:

```go
type Snapshot struct {
  ID           string          `json:"id"`
  Name         string          `json:"name"`
  Description  string          `json:"description,omitempty"`
  VolumeID     string          `json:"volume_id"`
  Size         int             `json:"size"`
  Status       SnapshotStatus  `json:"status"`
  StatusReason string          `json:"status_reason,omitempty"`
  Project      common.IDName   `json:"project"`
  ProjectID    string          `json:"project_id"`
  User         common.IDName   `json:"user"`
  UserID       string          `json:"user_id"`
  Namespace    string          `json:"namespace,omitempty"`
  CreatedAt    time.Time       `json:"createdAt,omitempty"`
  UpdatedAt    time.Time       `json:"updatedAt,omitempty"`
}
```
