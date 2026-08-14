# Feature Specification: Volumes API Client

**Feature Branch**: `008-volumes-api`  
**Created**: 2025-11-13  
**Status**: Draft  
**Input**: User description: "implement Volumes related restful api based on vps.yaml. Only implement api which has volumes and volume_types tag, and ignore others. All APIs are authorized through Bearer Token, and client has to provide the token when initialization."

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods;
  callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - List Available Volume Types (Priority: P1)

A developer wants to discover what types of storage volumes are available in the system before creating a volume, so they can select the appropriate type (e.g., SSD, HDD) for their workload requirements.

**Why this priority**: This is foundational knowledge required before any volume creation. Users must know available volume types to make informed decisions about storage provisioning.

**Independent Test**: Can be fully tested by calling the list volume types method and verifying the response contains a list of available storage types. Delivers immediate value by enabling users to discover available storage options.

**Acceptance Scenarios**:

1. **Given** a valid project ID and authentication token, **When** listing volume types, **Then** the system returns a list of available volume types (e.g., ["SSD", "HDD", "NVMe"])
2. **Given** an invalid project ID, **When** listing volume types, **Then** the system returns an appropriate error message
3. **Given** an expired authentication token, **When** listing volume types, **Then** the system returns an authentication error

---

### User Story 2 - Create and Manage Volumes (Priority: P1)

A developer needs to create storage volumes for their virtual machines, specifying name, size, and type. After creation, they should be able to update volume metadata (name, description) and delete volumes when no longer needed.

**Why this priority**: Core CRUD operations for volume lifecycle management. This is the primary use case for the volumes API and must work independently for the SDK to be useful.

**Independent Test**: Can be fully tested by creating a volume, verifying its creation, updating its metadata, retrieving the updated volume, and then deleting it. Delivers complete volume lifecycle management value.

**Acceptance Scenarios**:

1. **Given** valid volume parameters (name, size, type), **When** creating a volume, **Then** the system returns the created volume with assigned ID and status
2. **Given** a volume ID and new metadata, **When** updating the volume, **Then** the system updates the volume's name and/or description
3. **Given** a volume ID for an unattached volume, **When** deleting the volume, **Then** the system successfully removes the volume
4. **Given** a volume ID for an attached volume, **When** deleting the volume, **Then** the system returns an error indicating the volume is in use
5. **Given** a volume creation request exceeding quota, **When** creating the volume, **Then** the system returns a quota exceeded error

---

### User Story 3 - List and Retrieve Volume Information (Priority: P2)

A developer wants to list all volumes in their project with optional filtering by name, status, or type, and retrieve detailed information about specific volumes to monitor storage usage and status.

**Why this priority**: Essential for volume discovery and monitoring, but can be tested after basic CRUD operations are functional. Users need this for operational visibility.

**Independent Test**: Can be fully tested by creating several volumes with different attributes, then listing with various filters and retrieving individual volume details. Delivers value for volume inventory management.

**Acceptance Scenarios**:

1. **Given** a project with multiple volumes, **When** listing all volumes, **Then** the system returns all volumes in the project
2. **Given** a project with volumes, **When** listing with name filter, **Then** the system returns only volumes matching the name pattern
3. **Given** a project with volumes, **When** listing with status filter (e.g., "available", "in-use"), **Then** the system returns only volumes with matching status
4. **Given** a project with volumes, **When** listing with type filter, **Then** the system returns only volumes of the specified type
5. **Given** a project with volumes, **When** listing with detail=true, **Then** the system includes attachment information in the response
6. **Given** a valid volume ID, **When** retrieving volume details, **Then** the system returns complete volume information including status, size, attachments, and metadata

---

### User Story 4 - Perform Volume Actions (Priority: P2)

A developer needs to perform operations on volumes such as attaching to servers, detaching from servers, extending volume size, and reverting to snapshots.

**Why this priority**: Advanced operations that depend on volumes existing. Critical for production use but can be implemented and tested after basic CRUD operations.

**Independent Test**: Can be fully tested by creating a volume, performing actions (attach, extend, detach), and verifying state changes. Delivers value for advanced volume lifecycle operations.

**Acceptance Scenarios**:

1. **Given** a volume ID and server ID, **When** attaching the volume to the server, **Then** the system successfully attaches the volume and updates its status
2. **Given** an attached volume ID, **When** detaching the volume from the server, **Then** the system successfully detaches the volume and updates its status
3. **Given** a volume ID and new size greater than current size, **When** extending the volume, **Then** the system increases the volume size successfully
4. **Given** a volume ID and new size smaller than current size, **When** extending the volume, **Then** the system returns an error
5. **Given** a volume created from snapshot, **When** reverting the volume, **Then** the system restores the volume to the snapshot state

---

### User Story 5 - Create Volumes from Snapshots (Priority: P3)

A developer wants to create new volumes from existing snapshots to restore data or clone volumes for testing/development purposes.

**Why this priority**: Useful feature but depends on both volumes and snapshots existing. Lower priority as it's more advanced use case.

**Independent Test**: Can be fully tested by creating a volume from a snapshot ID and verifying the volume is created with snapshot data. Delivers value for backup/restore and cloning workflows.

**Acceptance Scenarios**:

1. **Given** a valid snapshot ID, **When** creating a volume from the snapshot, **Then** the system creates a volume with data from the snapshot
2. **Given** an invalid snapshot ID, **When** creating a volume from the snapshot, **Then** the system returns an appropriate error
3. **Given** a snapshot in a different project, **When** creating a volume from the snapshot, **Then** the system returns an authorization error

---

### Edge Cases

- What happens when attempting to delete a volume that is currently attached to a running server?
- How does the system handle volume creation when quota limits are exceeded?
- What happens when attempting to extend a volume to a size that exceeds project quotas?
- How does the system handle concurrent modification of the same volume (race conditions)?
- What happens when attempting to attach a volume that is already attached to another server?
- How are volume operations handled during server maintenance or outages?
- What happens when authentication token expires mid-operation?
- How does the list operation perform with very large numbers of volumes (no pagination available)?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST provide a method to list all available volume types in a project
- **FR-002**: System MUST provide a method to create volumes with required parameters (name, type, size) and optional parameters (description, snapshot_id)
- **FR-003**: System MUST provide a method to list volumes with optional filters (name, user_id, status, type, detail)
- **FR-004**: System MUST provide a method to retrieve detailed information about a specific volume by ID
- **FR-005**: System MUST provide a method to update volume metadata (name, description)
- **FR-006**: System MUST provide a method to delete volumes
- **FR-007**: System MUST provide a method to perform volume actions (attach, detach, extend, revert)
- **FR-008**: All API methods MUST accept a context.Context as the first parameter for cancellation and timeout control
- **FR-009**: All API methods MUST require project ID as a parameter to scope operations to the correct project
- **FR-010**: Client MUST be initialized with a Bearer Token for authentication
- **FR-011**: Client MUST automatically include Bearer Token in Authorization header for all API requests
- **FR-012**: System MUST return strongly-typed Go structs matching the Swagger/OpenAPI schema definitions
- **FR-013**: System MUST return descriptive errors that include HTTP status codes and service error details without exposing raw HTTP responses
- **FR-013a**: When authentication token expires, system MUST return an authentication error immediately without retry, allowing caller to handle token refresh
- **FR-014**: Volume action operations MUST validate action type (attach, detach, extend, revert) and require appropriate parameters for each action type
- **FR-015**: Volume creation MUST accept size parameter in GB; validation of size limits is delegated to server-side enforcement
- **FR-016**: Volume attachment action MUST require server_id parameter
- **FR-017**: Volume extend action MUST require new_size parameter
- **FR-018**: List volumes operation MUST support returning attachment information when detail=true parameter is provided
- **FR-019**: SDK MUST log HTTP request/response details (method, URL, status code) and retry attempts at DEBUG level using the logger provided during client initialization
- **FR-020**: List volumes operation returns complete result set in single response without pagination support

### SDK Contract Requirements (Go)

- Public methods follow the pattern: `Client.Volumes.Operation(ctx, projectID, params)` and `Client.VolumeTypes.Operation(ctx, projectID, params)`
- Responses are strongly-typed structs matching Swagger models (`VolumeInfo`, `VolumeListOutput`, `VolumeTypeListOutput`)
- Errors wrap context (status code, service error fields) without exposing raw HTTP
- Authentication via Bearer Token is centralized in the Client initialization
- Methods return `(*VolumeInfo, error)`, `(*VolumeListOutput, error)`, etc.
- All operations accept `context.Context` for timeout and cancellation
- Client initialization: `NewClient(baseURL, bearerToken string)` returns configured client
- Volume operations namespace: `client.Volumes.Create()`, `client.Volumes.List()`, `client.Volumes.Get()`, `client.Volumes.Update()`, `client.Volumes.Delete()`, `client.Volumes.Action()`
- Volume type operations namespace: `client.VolumeTypes.List()`

### Key Entities

- **Volume**: Represents a block storage volume with attributes including ID, name, description, size (in GB), type (storage class), status, attachments (list of servers), project information, user information, creation/update timestamps, and optional snapshot reference
- **VolumeType**: Represents an available storage type in the system (e.g., "SSD", "HDD", "NVMe") that can be used when creating volumes
- **VolumeAction**: Represents an operation to perform on a volume (attach, detach, extend, revert) with action-specific parameters
- **Project**: The organizational scope for volume operations, identified by project ID
- **Server**: Virtual machine that volumes can be attached to, referenced by server ID in attachment operations

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Developers can list available volume types in under 2 seconds
- **SC-002**: Developers can create a volume and receive confirmation in under 5 seconds
- **SC-003**: Developers can list volumes with filters and receive results in under 3 seconds
- **SC-004**: Developers can retrieve volume details in under 2 seconds
- **SC-005**: Developers can update volume metadata in under 3 seconds
- **SC-006**: Developers can delete an unattached volume in under 5 seconds
- **SC-007**: Developers can perform volume actions (attach/detach/extend) with confirmation in under 10 seconds
- **SC-008**: 100% of API endpoints tagged with "volumes" or "volume_types" in Swagger specification are implemented in the SDK
- **SC-009**: All SDK methods properly propagate context cancellation within 1 second
- **SC-010**: SDK provides clear, actionable error messages for all API error responses (4xx, 5xx)
- **SC-011**: Authentication errors are clearly distinguishable from other error types
- **SC-012**: List operations return complete volume list without pagination issues for projects with typical volume counts

### Quality Gates

- **QG-001**: All unit tests pass with 80%+ code coverage
- **QG-002**: All contract tests against Swagger specification pass
- **QG-003**: Integration tests validate end-to-end volume lifecycle (create, list, get, update, action, delete)
- **QG-004**: Error handling tests cover all documented error responses (400, 500, authentication failures)
- **QG-005**: Context cancellation tests validate proper cleanup and timeout behavior
- **QG-006**: Documentation includes code examples for each public method

## Assumptions

- The VPS API server is accessible and operational at the configured base URL
- Bearer tokens are obtained through a separate authentication mechanism (out of scope for this feature)
- Bearer tokens have sufficient permissions for all volume operations in the specified project
- Project IDs are valid and accessible to the authenticated user
- Volume types returned by the API are consistent and don't change during runtime
- The Swagger schema in vps.yaml accurately reflects the actual API behavior
- HTTP status codes follow RESTful conventions (200 OK, 201 Created, 204 No Content, 400 Bad Request, 500 Internal Server Error)
- Volume operations are asynchronous where appropriate (the API may return accepted status before completion)
- The SDK will use the existing internal HTTP client (if present) or create a minimal HTTP client wrapper
- The SDK will follow Go naming conventions and idiomatic patterns used in the existing codebase

## Dependencies

- Existing internal HTTP client implementation (`internal/http/client.go`) with automatic retry using exponential backoff
- Existing types and utilities (`internal/types`)
- Go standard library: `context`, `net/http`, `encoding/json`, `fmt`, `time`
- VPS API server endpoints as defined in `swagger/vps.yaml`
- Bearer token authentication infrastructure (assumed to be provided by caller)

## Out of Scope

- Authentication token generation or refresh mechanisms
- Snapshot management APIs (separate feature)
- Server management APIs (separate feature)
- Quota management APIs
- Billing or cost tracking for volumes
- Volume encryption configuration
- Volume backup scheduling
- Multi-region volume replication
- Volume performance monitoring and metrics
- Admin-level volume operations (creating volume types, managing backend storage)

## Clarifications

### Session 2025-11-13

- Q: HTTP Request Retry Strategy - How should the SDK handle transient failures (network errors, 5xx responses)? → A: Automatic retry with exponential backoff for 5xx and network errors (max 3 attempts) - consistent with existing internal/http/client.go implementation
- Q: Bearer Token Expiry Handling - How should the SDK handle authentication token expiry during operations? → A: Return authentication error immediately - caller responsible for token refresh
- Q: Volume Size Constraints - Should the SDK validate minimum/maximum volume size limits? → A: No client-side validation - rely entirely on server-side validation
- Q: Logging and Observability - What SDK operations should be logged for debugging? → A: Log HTTP requests/responses (method, URL, status) and retry attempts at DEBUG level - consistent with existing logger pattern passed to NewClient()
- Q: Pagination Strategy for Large Volume Lists - How should the SDK handle pagination for list operations? → A: No pagination support - API returns full list in single response
