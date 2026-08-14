## Purpose

Provides Swagger-aligned security group and rule models and typed project-scoped operations for managing them.

## Requirements

### Requirement: Security group model contract
Security group list, detail, create, update, and delete payloads MUST match the documented API structures and preserve structured errors.

#### Scenario: Use security group CRUD
- **WHEN** an authorized caller lists or performs CRUD on a security group
- **THEN** requests and decoded responses conform to the API contract

### Requirement: Security group rule subresource
The client MUST support creating and deleting security group rules with documented protocol, port, CIDR, and related fields.

#### Scenario: Manage a rule
- **WHEN** an authorized caller adds or removes a valid rule
- **THEN** the rule subresource request succeeds and the resulting rule set reflects the operation

### Requirement: Authorization and conflicts
Security group operations MUST rely on project permissions and propagate API-level authorization and optimistic-lock conflict errors.

#### Scenario: Unauthorized or conflicting update
- **WHEN** a caller lacks permission or updates a stale representation
- **THEN** the client returns the corresponding structured API error without hiding the conflict
