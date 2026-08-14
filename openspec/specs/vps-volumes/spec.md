## Purpose

Provides typed project-scoped access to VPS volume types and volume lifecycle operations defined by the VPS API.

## Requirements

### Requirement: Discover volume types
The client MUST list documented volume types for an authenticated project and decode them into typed results.

#### Scenario: List volume types
- **WHEN** a caller requests volume types with valid project context and authentication
- **THEN** the client returns the available types or a structured project/authentication error

### Requirement: Manage volumes
The client MUST support the documented volume list, create, get, update, delete, and action operations using context-aware typed models.

#### Scenario: Create and retrieve volume
- **WHEN** a caller creates a valid volume and then retrieves it
- **THEN** requests and responses match the VPS contract and expose the resulting volume data

### Requirement: Volume validation and errors
The client MUST validate documented required inputs before requests and preserve API errors for invalid project, authentication, conflict, or lifecycle states.

#### Scenario: Invalid volume request
- **WHEN** a caller omits a required field
- **THEN** the client returns a validation error without making an HTTP request
