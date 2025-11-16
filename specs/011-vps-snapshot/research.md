# Research: vps snapshot resource

## Decisions made

- Swagger (`swagger/vps.yaml`) will be the source of truth for models and endpoints; `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` is the canonical snapshot schema.
- Snapshot creation is asynchronous from a system perspective but returns 201 with a `status` field that clients can poll. This mirrors the repo's existing async patterns for volumes and actions.
- Snapshots may be taken while a volume is attached (hot-snapshot). An optional `quiesce` parameter is out-of-scope for v1 but noted for future work.
- Deleting a snapshot does not affect volumes created from it (volumes become independent); retention/lifecycle/TTL are out-of-scope.

## Rationale

 - Following the Swagger spec ensures backwards-compatible API contracts and enables repository-level contract tests to be precise and automated.
- Returning 201 with `status` reduces client blocking and is consistent with other resource create patterns in the SDK (ease of implementation and clarity).
- Hot-snapshot support is a low friction, higher-value user feature. `quiesce` will be deferred to avoid extra operational and API complexity in the initial delivery.

## Alternatives considered

- Synchronous create (block until snapshot available): rejected due to longer operation times and issues with timeouts.
- Force-protect lifecycle: rejected for v1 to keep scope small and avoid storage lifecycle changes.
