## Purpose

Provides bearer-authenticated typed IAM access for current-user and project discovery operations.

## Requirements

### Requirement: IAM endpoint selection
The IAM client MUST send requests to the IAM endpoint derived from the Cloud SDK client's base URL by default. When the client is constructed with a dedicated identity URL, the IAM client MUST instead send requests to the IAM endpoint derived from that identity URL.

#### Scenario: Default IAM endpoint remains compatible
- **WHEN** a caller constructs a client with the existing base-URL constructor and invokes an IAM operation
- **THEN** the request is sent to the IAM endpoint derived from the base URL

#### Scenario: Configured identity endpoint overrides default
- **WHEN** a caller constructs a client with a dedicated identity URL and invokes an IAM operation
- **THEN** the request is sent to the IAM endpoint derived from the identity URL

### Requirement: Current user retrieval
The IAM client MUST provide a context-aware operation that retrieves the authenticated user's documented account, identity, metadata, namespace, timestamps, and permissions.

#### Scenario: Retrieve current user
- **WHEN** a caller invokes the user operation with a valid token
- **THEN** the client returns the complete typed user response

### Requirement: Project retrieval
The IAM client MUST provide typed operations for listing accessible projects and retrieving a project by ID.

#### Scenario: Retrieve accessible projects
- **WHEN** a caller requests project information with valid authentication
- **THEN** the client returns the documented project collection or a structured authorization error

### Requirement: IAM request errors
IAM requests MUST honor context deadlines and propagate invalid or expired bearer-token responses as structured errors.

#### Scenario: IAM request timeout
- **WHEN** the context deadline expires during an IAM request
- **THEN** the client returns a timeout or cancellation error without hanging
