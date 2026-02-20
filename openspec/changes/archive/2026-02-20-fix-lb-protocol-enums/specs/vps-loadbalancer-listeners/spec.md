## ADDED Requirements

### Requirement: ListenerProtocol enum type
The SDK SHALL define a `ListenerProtocol` typed string constant with a `Valid()` method.

#### Scenario: Valid listener protocols
- **WHEN** `ListenerProtocol.Valid()` is called on `ListenerProtocolHTTP`, `ListenerProtocolHTTPS`, `ListenerProtocolTerminatedHTTPS`, `ListenerProtocolTCP`, or `ListenerProtocolUDP`
- **THEN** it SHALL return `true`

#### Scenario: Invalid listener protocol
- **WHEN** `ListenerProtocol.Valid()` is called on any other value
- **THEN** it SHALL return `false`

## MODIFIED Requirements

### Requirement: Listener model
The SDK SHALL define a `Listener` struct matching `pb.ListenerInfo` with all fields and JSON tags.

#### Scenario: Listener has required fields
- **WHEN** a `Listener` struct is populated from API response
- **THEN** it SHALL contain fields: ID, Name, LbID, Protocol, ProtocolPort, PoolID, Pool, Status, StatusReason, AllowedCIDRs, InsertHeaders, DefaultTLSContainerRef, TimeoutClientData, TimeoutMemberConnect, TimeoutMemberData, TimeoutTCPInspect

#### Scenario: Listener Protocol field uses ListenerProtocol type
- **WHEN** `Listener.Protocol` is declared
- **THEN** its Go type SHALL be `ListenerProtocol` (not `string`)

---

### Requirement: ListenerCreateRequest validation
The `ListenerCreateRequest` model SHALL validate required fields.

#### Scenario: Valid create request
- **WHEN** `Validate()` is called with non-empty `Protocol` and non-zero `ProtocolPort`
- **THEN** it SHALL return nil

#### Scenario: Missing protocol
- **WHEN** `Validate()` is called with empty `Protocol`
- **THEN** it SHALL return an error containing "protocol is required"

#### Scenario: Invalid protocol
- **WHEN** `Validate()` is called with a `Protocol` value not in `{HTTP, HTTPS, TERMINATED_HTTPS, TCP, UDP}`
- **THEN** it SHALL return an error containing "protocol must be one of"

#### Scenario: Missing protocol_port
- **WHEN** `Validate()` is called with zero `ProtocolPort`
- **THEN** it SHALL return an error containing "protocol_port is required"

#### Scenario: ListenerCreateRequest Protocol field uses ListenerProtocol type
- **WHEN** `ListenerCreateRequest.Protocol` is declared
- **THEN** its Go type SHALL be `ListenerProtocol` (not `string`)
