## Purpose

Represents floating IP data and operations using the documented API structure, metadata, timestamps, and lifecycle semantics.

## ADDED Requirements

### Requirement: Floating IP model contract
The FloatingIP model MUST match the documented fields and use the `floatingips` list property rather than the removed legacy `items` property.

#### Scenario: Decode floating IP list
- **WHEN** a documented list response is decoded
- **THEN** floating IPs are available through the `floatingips` collection and no legacy collection is required

### Requirement: Floating IP lifecycle
The client MUST support documented list, create, get, update, delete, approve, reject, and disassociate operations with typed results and structured errors.

#### Scenario: Disassociate floating IP
- **WHEN** a caller disassociates a valid floating IP
- **THEN** the client sends the documented request and returns the API-defined asynchronous result

### Requirement: Optional and read-only fields
RFC3339 timestamps and optional project/user objects MUST preserve the API response shape; `reserved` MUST be treated as response-only during creation, and no status helper semantics are imposed.

#### Scenario: Create with reserved value
- **WHEN** a create request includes a reserved field
- **THEN** the client does not treat it as writable and the API response remains authoritative
