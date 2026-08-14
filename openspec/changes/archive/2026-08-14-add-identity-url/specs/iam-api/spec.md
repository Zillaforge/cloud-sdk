## MODIFIED Requirements

### Requirement: IAM endpoint selection
The IAM client MUST send requests to the IAM endpoint derived from the Cloud SDK client's base URL by default. When the client is constructed with a dedicated identity URL, the IAM client MUST instead send requests to the IAM endpoint derived from that identity URL.

#### Scenario: Default IAM endpoint remains compatible
- **WHEN** a caller constructs a client with the existing base-URL constructor and invokes an IAM operation
- **THEN** the request is sent to the IAM endpoint derived from the base URL

#### Scenario: Configured identity endpoint overrides default
- **WHEN** a caller constructs a client with a dedicated identity URL and invokes an IAM operation
- **THEN** the request is sent to the IAM endpoint derived from the identity URL
