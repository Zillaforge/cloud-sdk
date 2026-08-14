## Purpose

Provides typed VPS snapshot lifecycle operations for capturing and managing data-volume state according to the authoritative Swagger contract.

## Requirements

### Requirement: Create snapshots
The client MUST create snapshots for valid data volumes using the documented project-scoped request and return the typed snapshot result.

#### Scenario: Create snapshot
- **WHEN** an authenticated caller submits a valid snapshot request
- **THEN** the client sends the documented request and returns the created snapshot

### Requirement: Manage snapshots
The client MUST support documented snapshot list, get, update, delete, and lifecycle operations with Swagger-aligned request and response models.

#### Scenario: Retrieve snapshot
- **WHEN** a caller requests an existing snapshot by ID
- **THEN** the client returns its complete typed representation or the API's structured not-found error

### Requirement: Snapshot validation and authorization
Snapshot operations MUST require project-scoped bearer authentication, validate required input, and propagate context, authentication, and API errors.

#### Scenario: Invalid snapshot request
- **WHEN** a caller submits an invalid request or lacks authorization
- **THEN** the client returns a validation or structured API error before exposing an invalid result
