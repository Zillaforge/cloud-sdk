# Feature Specification: vps snapshot resource

**Feature Branch**: `011-vps-snapshot`  
**Created**: 2025-11-16  
**Status**: Draft  
**Input**: User description: "$ARGUMENTS"

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods;
  callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
 - Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass. The Swagger file `swagger/vps.yaml` is the authoritative source of truth for API endpoints and models used by contract tests.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## User Scenarios & Testing *(mandatory)*

<!--
  IMPORTANT: User stories should be PRIORITIZED as user journeys ordered by importance.
  Each user story/journey must be INDEPENDENTLY TESTABLE - meaning if you implement just ONE of them,
  you should still have a viable MVP (Minimum Viable Product) that delivers value.
  
  Assign priorities (P1, P2, P3, etc.) to each story, where P1 is the most critical.
  Think of each story as a standalone slice of functionality that can be:
  - Developed independently
  - Tested independently
  - Deployed independently
  - Demonstrated to users independently
-->

### User Story 1 - Create Snapshot (Priority: P1)

As a cloud user I need to create a snapshot of a data volume so that I can safely capture the current disk state for backups and restoring later.

**Why this priority**: Snapshots are a fundamental data-protection capability and are required for backup/restore workflows and for creating volumes from snapshots.

**Independent Test**: Call `POST /api/v1/project/{project-id}/snapshots` with valid `volume_id` and receive a 201 Created response with snapshot details; the created snapshot can then be retrieved via GET.

**Acceptance Scenarios**:

1. **Given** the user has an existing data volume in the same project, **When** they call POST /snapshots with a valid `volume_id` and a `name`, **Then** the API returns 201 Created and a snapshot resource with id, name, volume_id, status, createdAt.
2. **Given** the user attempts to create a snapshot for a non-existent volume or a root/unsupported volume, **When** they call POST /snapshots, **Then** the API returns a 400 Bad Request or 403 Forbidden depending on cross-project scope.

 

---

### User Story 2 - List & Inspect Snapshots (Priority: P1)

As a cloud user I need to list snapshots in my project and inspect details so I can manage backups and find restore candidates.

**Why this priority**: Listing and inspecting snapshots is core to snapshot lifecycle management and must be available for most user flows.

**Independent Test**: Call `GET /api/v1/project/{project-id}/snapshots` to get a list; call `GET /snapshots/{snapshot-id}` to get details.

**Acceptance Scenarios**:

1. **Given** snapshots exist in the project, **When** the user calls GET /snapshots, **Then** the API returns 200 OK and a list of snapshots including id, name, volume_id, status, project id and timestamps.
 1a. **Given** multiple snapshots exist in the project with varying names/volume_id/status, **When** the user calls GET /snapshots with query filters (name, volume_id, user_id, status), **Then** the API returns a filtered list according to the parameters (and supports pagination).
2. **Given** a snapshot id exists, **When** the user calls GET /snapshots/{id}, **Then** the API returns 200 OK with snapshot details.

 

---

### User Story 3 - Update and Delete Snapshot (Priority: P2)

As a cloud user I want to update snapshot metadata (such as rename) and delete snapshots when no longer needed.

**Why this priority**: Renaming is a convenience for administration; deletion is necessary to reclaim storage and enforce lifecycle policies.

**Independent Test**: Call `PUT /snapshots/{id}` to update the name; call `DELETE /snapshots/{id}` to delete and ensure list/update reflect deletion.

**Acceptance Scenarios**:

1. **Given** a snapshot id exists, **When** the user calls PUT /snapshots/{id} with a new name, **Then** the API returns 200 OK and snapshot shows the updated name.
2. **Given** a snapshot id exists and user calls DELETE /snapshots/{id}, **Then** the API returns 204 No Content (or 200/202 depending on async semantics) and subsequent GET returns 404.
 2. **Given** a snapshot id exists and user calls DELETE /snapshots/{id}, **Then** the API returns 204 No Content (synchronous delete for this release) and subsequent GET returns 404. Asynchronous delete (202 + Location) is out-of-scope for the initial delivery.

 

---

[Add more user stories as needed, each with an assigned priority]

### Edge Cases

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right edge cases.
-->

 - What happens when the user tries to snapshot a root disk or snapshot a volume from a different project? (Expect 403 or 400 depending on policy)
 - How does system handle concurrent snapshot creation for the same volume? (service should serialize or provide meaningful errors)
 - What happens when the snapshot service has storage quota exhaustion? (return 412 Precondition Failed or 429/507 with clear message)
 - What happens when the user tries to snapshot a root disk or snapshot a volume from a different project? (Expect 403 or 400 depending on policy)
 - How does system handle concurrent snapshot creation for the same volume? (service should serialize or provide meaningful errors)
 - What happens when the snapshot service has storage quota exhaustion? (return 412 Precondition Failed or 429/507 with clear message)
 - Does deleting a snapshot affect volumes created from it? (Answer: No — volumes remain independent)
 - Does the snapshot require the volume to be detached for consistent state? (Answer: No — hot-snapshot is allowed; tests verify volume is still attachable and operations continue.)

## Requirements *(mandatory)*

<!--
  ACTION REQUIRED: The content in this section represents placeholders.
  Fill them out with the right functional requirements.
-->

### Functional Requirements

 - **FR-001**: System MUST provide an API to create a snapshot for a data volume in the same project: POST /api/v1/project/{project-id}/snapshots. Request must accept `name` and `volume_id`. Response returns the new snapshot record.
 - **FR-001**: The API returns 201 Created immediately and the returned snapshot MUST include a `status` field (e.g., `creating`, `available`); clients may poll until `available`.
 - **FR-002**: System MUST provide an API to list snapshots with optional filters: GET /api/v1/project/{project-id}/snapshots (support filters: name, volume_id, user_id, status).
 - **FR-003**: System MUST provide an API to retrieve snapshot metadata: GET /api/v1/project/{project-id}/snapshots/{snapshot-id}.
 - **FR-004**: System MUST provide an API to update snapshot metadata: PUT /api/v1/project/{project-id}/snapshots/{snapshot-id} (e.g., rename).
 - **FR-005**: System MUST provide an API to delete a snapshot: DELETE /api/v1/project/{project-id}/snapshots/{snapshot-id}. Deleting a snapshot MUST NOT delete volumes previously created from that snapshot; those volumes remain independent. For this release, the API MUST return 204 No Content on successful synchronous delete. If async deletion is required later, the API should return 202 and include a `Location` header for polling (future enhancement).
 - **FR-006**: SDK: Provide a typed model `models/vps/snapshots` that represents snapshot fields (id, name, volume_id, status, project, user, createdAt, updatedAt, size, statusReason). Models MUST follow shapes defined in `swagger/vps.yaml` (e.g., `SnapshotCreateInput`, `SnapshotUpdateInput`, `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo`) and preserve JSON field names and types (including `project`, `user`, `namespace`). Use shared types from `models/vps/common` where applicable (e.g., `common.IDName`) to avoid duplicate definitions.
 - **FR-007**: SDK: Provide a module `modules/vps/snapshots` with a client exposing Create, List, Get, Update, Delete operations. The client must use endpoints and parameters as defined in `swagger/vps.yaml` and mirror the request/response formats declared by the API contract definitions. `List()` MUST return `[]*models/vps/snapshots.Snapshot` (indexable) and perform any necessary marshalling from `pb.SnapshotListOutput`.
 - **FR-008**: System MUST validate that `volume_id` exists and belongs to same project on creation, returning 400 or 403 on violation.
 - **FR-009**: System MUST expose appropriate contract tests validating existence of endpoints and common error cases. Tests MUST read the Swagger `vps.yaml` schemas and validate contract behavior accordingly (create/list/get/update/delete, invalid `volume_id`, cross-project access) and assert that deletion does not break volumes created from snapshots.
 - **FR-010**: System SHOULD allow snapshot creation when a volume is attached (hot-snapshot); the API MUST not require a detach by default. An optional `quiesce` parameter may be accepted in future to request application-consistent snapshots.
 - **FR-010**: System SHOULD allow snapshot creation when a volume is attached (hot-snapshot); the API MUST not require a detach by default. An optional `quiesce` parameter may be accepted in future to request application-consistent snapshots.
 - **FR-011**: No lifecycle/retention features (`protected` flag, TTL-based auto-delete) will be implemented in the initial delivery. These are explicitly out of scope and may be added in follow-up features if desired.
 - **FR-009**: (See above) Ensure repository-level contract tests contain checks for cross-project and invalid input cases, filter behavior, and deletion semantics.

### Testing & Contract Tests

- Add repository-level contract tests; these tests must include create/list/get/update/delete and error cases such as invalid `volume_id` and cross-project access.
- Unit tests for `modules/vps/snapshots` must include happy-path request serialization/deserialization and error-handling; follow patterns used in `modules/vps/volumes/client_test.go`.
- Integration/contract test expectations: service returns 201 Created on successful create, 400/403 on invalid input, 404 on missing snapshot, and 200 for get/list.
 - Integration/contract test expectations: service returns 201 Created on successful create, 400/403 on invalid input, 404 on missing snapshot, and 200 for get/list. Tests MUST compare JSON response fields against the types defined in `swagger/vps.yaml` (Schema validation). Use `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` as the canonical contract for snapshot responses.
 - Unit tests for `modules/vps/snapshots` must include happy-path request serialization/deserialization and error-handling; follow patterns used in `modules/vps/volumes/client_test.go`.
 - Integration/contract test expectations: service returns 201 Created on successful create, 400/403 on invalid input, 404 on missing snapshot, and 200 for get/list.
 - Contract tests MUST include creating a snapshot while a volume is attached and assert the volume remains available and the returned snapshot indicates `creating` state.
 - Contract tests MUST include creating a snapshot while a volume is attached and assert the volume remains available and the returned snapshot indicates `creating` state.
 - Contract tests MUST validate snapshot deletion is immediate and not subject to TTL or delayed expiry for this release (i.e., no automatic lifecycle is present).
- **FR-002**: System MUST [specific capability, e.g., "validate email addresses"]  
- **FR-003**: Users MUST be able to [key interaction, e.g., "reset their password"]
- **FR-004**: System MUST [data requirement, e.g., "persist user preferences"]
- **FR-005**: System MUST [behavior, e.g., "log all security events"]

### SDK Contract Requirements (Go)

- Public methods follow the pattern: `Client.Snapshots.Operation(ctx, params)` consistent with other modules (e.g., `modules/vps/volumes`).
- Response types MUST be strongly-typed structs: `Snapshot`, `CreateSnapshotRequest`, `UpdateSnapshotRequest`, `ListSnapshotsOptions`, `SnapshotListResponse`.
 - Response types MUST be strongly-typed structs and match the Swagger definitions: `Snapshot` should match `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` in `swagger/vps.yaml`. Request/response models MUST be compatible with the API in `swagger/vps.yaml`.
 - Response types MUST be strongly-typed structs and match the Swagger definitions: `Snapshot` should match `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` in `swagger/vps.yaml`. Request/response models MUST be compatible with the API in `swagger/vps.yaml`. Use `Request`/`Response` naming convention in the SDK models instead of `Input`/`Output` (e.g., `CreateSnapshotRequest`, `SnapshotResponse`).
- Errors returned by the module MUST wrap service error details (message, error code) but not leak internal HTTP details.
- Authentication, retries, and timeouts are centralized in the Client.
 - Authentication, retries, and timeouts are centralized in the Client.
 - **FR-012**: The SDK must include a small, documented mapping or conversion function when the Swagger `pb` representation differs from the Go struct naming (e.g., mapping `createdAt` <-> `CreatedAt`) to ensure contract tests pass and fields serialize/deserialize exactly as required.
 - **FR-013**: The SDK MUST implement `SnapshotStatus` as a custom Go type (enum-like constants) rather than returning raw strings to callers. Client serialization should still use the string values defined in `swagger/vps.yaml` when marshalling.

*Example of marking unclear requirements:* (none)

### Key Entities *(include if feature involves data)*

- **[Entity 1]**: [What it represents, key attributes without implementation]
- **[Entity 2]**: [What it represents, relationships to other entities]

### Key Entities

- **Snapshot**: A point-in-time capture of a data volume. Attributes include `id`, `name`, `volume_id`, `size`, `status`, `status_reason`, `project_id`, `user_id`, `namespace`, `createdAt`, `updatedAt`.
- **Snapshot**: A point-in-time capture of a data volume. Attributes include `id`, `name`, `volume_id`, `size`, `status`, `status_reason`, `project`, `project_id`, `user`, `user_id`, `namespace`, `createdAt`, `updatedAt`.
- **Volume**: Reference to `volumes.Volume` (existing model). In this context `volume_id` is required for snapshot creation and links to the volume to be captured.
- **Project**: Namespace grouping resources and permissions; snapshots are scoped to a project.
- **User**: Owner/creator of the snapshot; used for auditing and access control.

### Assumptions

- Snapshots are supported only for data volumes (not root OS volumes) by service policy; the service returns a 400/403 for unsupported snapshot requests.
- Snapshot creation is logically synchronous for API contract (returns the snapshot record) but the snapshot may be `creating` then `available` - test code expects eventual readiness.
 - Snapshot creation is logically synchronous for API contract (returns the snapshot record) but the snapshot may be `creating` then `available` - test code expects eventual readiness.
 - Snapshots may be taken while the volume is attached (hot-snapshot support). If application consistency is required, service can provide an optional `quiesce` parameter in future enhancements — tests should validate snapshot integrity and volume availability when this occurs.
- Deleting a snapshot removes its metadata and storage. If the snapshot is referenced by other resources (e.g., last created volume) behavior depends on the service; current assumption is the service will allow deletion and may fail if protected by policy.
- No retention or lifecycle automation is required as part of this spec; lifecycle policies are out of scope for initial delivery.
 - No retention or lifecycle automation is required as part of this spec; lifecycle policies are out of scope for initial delivery. For clarity: `protected` flags or TTL-based expiration MUST NOT be implemented in the initial delivery.

## Success Criteria *(mandatory)*

<!--
  ACTION REQUIRED: Define measurable success criteria.
  These must be technology-agnostic and measurable.
-->

### Measurable Outcomes

- **SC-001**: Users can create a snapshot and then retrieve it via GET within 3 seconds (end-to-end observable for basic operations in test environment).
- **SC-002**: System returns correct status codes for common errors (400 for invalid input, 403 for cross-project violations, 404 for missing snapshot).
 - **SC-003**: 95% of snapshot list calls return within 1 second in testing environment with 1k snapshot objects. Performance tests must be added to `tests/perf` to assert this behavior (see `T023`).
- **SC-004**: Users can create a volume from a snapshot (CreateVolume with snapshot_id) and validate that snapshot_id restrictions are enforced.

- **SC-006**: Deleting a snapshot does not make previously-created volumes inaccessible or lose their data; tests validate that volumes remain available and their data consistent after snapshot delete.

 - **SC-005**: Snapshot creation API returns 201 Created immediately, and the snapshot `status` transitions from `creating` to `available`; SDK must surface this `status` and provide polling helpers, and tests must validate state transition semantics (see `T022`).

 - **SC-007**: Creating a snapshot while the volume is attached must not interrupt or make the volume unavailable; tests will assert volume access before/during/after snapshot creation.

- **SC-008**: Snapshot deletion is immediate (no TTL/soft-delete) in this release; tests must show that no retention/expiry is applied after DELETE.

### Clarifications

#### Session 2025-11-16

 - Q: Should snapshots be allowed while a volume is attached? → A: Option A — Allow hot-snapshot creation; the service may optionally provide a `quiesce` parameter for application-consistent snapshots.


### Non-functional Criteria

- The spec must contain model and module names for the SDK; contract tests should be added to the repository's contract test suite.
- The module must be consistent with existing naming patterns and client model behaviors.
- **SC-002**: [Measurable metric, e.g., "System handles 1000 concurrent users without degradation"]
 - **SC-002**: [Measurable metric, e.g., "System handles 1000 concurrent users without degradation"]
 - **SC-003**: [Customer satisfaction metric placeholder — remove or realign to product KPIs]
 - **SC-004**: [Business metric placeholder — remove or realign to product KPIs]
