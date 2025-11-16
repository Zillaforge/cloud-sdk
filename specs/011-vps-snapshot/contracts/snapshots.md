# API Contracts: Snapshot endpoints (derived from swagger/vps.yaml)

This file documents the core contracts for snapshot resources. The authoritative source is `swagger/vps.yaml` and tests rely on `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` definitions.

Endpoints

1. POST /api/v1/project/{project-id}/snapshots
   - Request: `SnapshotCreateRequest` {
       name string (required)
       volume_id string (required)
       description string (optional)
     }
   - Response: 201 Created, body: `SnapshotResponse` (snapshot fields: id, name, volume_id, project, project_id, user, user_id, namespace, status, size, createdAt, updatedAt) with `status` set to `creating`.

2. GET /api/v1/project/{project-id}/snapshots
   - Request: query filters (name, volume_id, user_id, status)
   - Response: 200 OK, body: `SnapshotListResponse` { snapshots: [ Snapshot ] }
   - `List()` returns []*Snapshot in SDK; tests must verify indexable access.

3. GET /api/v1/project/{project-id}/snapshots/{snapshot-id}
   - Response: 200 OK, body: Snapshot

4. PUT /api/v1/project/{project-id}/snapshots/{snapshot-id}
   - Request: `SnapshotUpdateRequest` { name, description }
   - Response: 200 OK, body: Snapshot

5. DELETE /api/v1/project/{project-id}/snapshots/{snapshot-id}
   - Response: 204 No Content (or 200/202 depending on service) — contract tests will accept 204.
   - Deleting a snapshot does not affect volumes created from it.

Error cases
- 400 Bad Request: invalid input (missing name/volume_id)
- 403 Forbidden: cross-project or permission denied
- 404 Not Found: snapshot or volume not found
- 409 Conflict: deletion blocked (not used in v1 due to out-of-scope lifecycle)

Schema references
- Snapshot: `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` from `swagger/vps.yaml`
- SnapshotCreateRequest: `SnapshotCreateInput` in `swagger/vps.yaml` but in SDK naming use `CreateSnapshotRequest`
