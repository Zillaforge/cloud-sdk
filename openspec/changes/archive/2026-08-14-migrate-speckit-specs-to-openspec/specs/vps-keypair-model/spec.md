## Purpose

Represents VPS keypairs and their operations according to the API contract, including timestamps, optional users, and sensitive private-key output.

## ADDED Requirements

### Requirement: Keypair model and operations
The keypair model and client MUST match documented list, create/import, get, update-description, and delete payloads; list responses MUST not expose the removed legacy Total field.

#### Scenario: Manage keypair
- **WHEN** a caller performs a documented keypair operation
- **THEN** the request and typed response match the API structure

### Requirement: Keypair timestamps and optional data
Timestamp fields MUST use RFC3339-compatible values, deserialization failures MUST identify the invalid field, and absent user data MUST remain empty or null as returned.

#### Scenario: Decode keypair metadata
- **WHEN** a response contains valid timestamps or omits its user object
- **THEN** timestamps decode correctly and the optional user remains absent

### Requirement: Sensitive private key handling
The API model MUST expose private-key output where documented and identify it as sensitive data that callers must save immediately.

#### Scenario: Receive generated keypair
- **WHEN** keypair creation returns a private key
- **THEN** the typed response exposes it without silently discarding it
