## Purpose

Provides complete Swagger-aligned server models and typed operations for server lifecycle, actions, metrics, NICs, and volume attachments.

## Requirements

### Requirement: Complete server model
The Server model MUST include all documented fields, and status fields MUST use typed enum-like values with constants for valid documented states.

#### Scenario: Decode server response
- **WHEN** a documented server response is decoded
- **THEN** all fields and typed status values are available without data loss

### Requirement: Server lifecycle and actions
The server client MUST support documented CRUD and start, stop, reboot, resize, extend-root, password, approve, and reject actions with their required parameters.

#### Scenario: Perform server action
- **WHEN** a caller invokes a supported action with valid parameters
- **THEN** the client sends the documented request and returns the typed result or structured error

### Requirement: Server subresources
The client MUST support documented metrics, NIC management, floating-IP association, console URL, and volume attach/detach operations.

#### Scenario: Manage server attachments
- **WHEN** a caller attaches or detaches a documented NIC or volume
- **THEN** the typed request and response follow the API contract
