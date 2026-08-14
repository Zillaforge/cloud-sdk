## Purpose

Provides project-scoped, bearer-authenticated typed access to the VRM User/Tag and User/Repository APIs defined by the VRM OpenAPI contract.

## Requirements

### Requirement: Project-scoped VRM client
The SDK MUST expose VRM through a project-scoped client, automatically include the bearer token, and support concurrent calls safely.

#### Scenario: Initialize VRM client
- **WHEN** a caller selects a project and accesses VRM
- **THEN** subsequent supported requests include the token and project context without caller-managed headers

### Requirement: Tag and repository operations
The VRM client MUST implement only the documented User/Tag and User/Repository operations with typed context-aware requests and responses.

#### Scenario: Call supported VRM API
- **WHEN** a caller invokes a supported tag or repository operation
- **THEN** the request and response conform to the VRM OpenAPI contract

### Requirement: Request timeout
VRM requests MUST have a configurable timeout with a default of approximately 30 seconds and MUST respect earlier context deadlines.

#### Scenario: Context deadline expires
- **WHEN** a VRM request exceeds its context deadline
- **THEN** it terminates and returns a timeout or cancellation error
