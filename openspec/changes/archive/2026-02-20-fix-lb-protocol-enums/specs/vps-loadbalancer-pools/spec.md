## ADDED Requirements

### Requirement: PoolProtocol enum type
The SDK SHALL define a `PoolProtocol` typed string constant with a `Valid()` method.

#### Scenario: Valid pool protocols
- **WHEN** `PoolProtocol.Valid()` is called on `PoolProtocolHTTP`, `PoolProtocolHTTPS`, `PoolProtocolTCP`, or `PoolProtocolUDP`
- **THEN** it SHALL return `true`

#### Scenario: Invalid pool protocol
- **WHEN** `PoolProtocol.Valid()` is called on any other value
- **THEN** it SHALL return `false`

## MODIFIED Requirements

### Requirement: Pool model
The SDK SHALL define a `Pool` struct matching `pb.PoolInfo` with all fields and JSON tags.

#### Scenario: Pool has required fields
- **WHEN** a `Pool` struct is populated from API response
- **THEN** it SHALL contain fields: ID, Name, Protocol, MemberPort, Method, Status, StatusReason, ASG, Listeners, Members

#### Scenario: Pool Protocol field uses PoolProtocol type
- **WHEN** `Pool.Protocol` is declared
- **THEN** its Go type SHALL be `PoolProtocol` (not `string`)

#### Scenario: Pool Method field uses PoolMethod type
- **WHEN** `Pool.Method` is declared
- **THEN** its Go type SHALL be `PoolMethod` (not `string`)

---

### Requirement: PoolCreateRequest validation
The `PoolCreateRequest` model SHALL validate required fields.

#### Scenario: Valid create request
- **WHEN** `Validate()` is called with `Name`, `Protocol`, `MemberPort`, and `Method` set to valid typed values
- **THEN** it SHALL return nil

#### Scenario: Missing name
- **WHEN** `Validate()` is called with empty `Name`
- **THEN** it SHALL return an error containing "name is required"

#### Scenario: Missing protocol
- **WHEN** `Validate()` is called with empty `Protocol`
- **THEN** it SHALL return an error containing "protocol is required"

#### Scenario: Invalid protocol
- **WHEN** `Validate()` is called with a `Protocol` value not in `{HTTP, HTTPS, TCP, UDP}`
- **THEN** it SHALL return an error containing "protocol must be one of"

#### Scenario: Missing member_port
- **WHEN** `Validate()` is called with zero `MemberPort`
- **THEN** it SHALL return an error containing "member_port is required"

#### Scenario: Missing method
- **WHEN** `Validate()` is called with empty `Method`
- **THEN** it SHALL return an error containing "method is required"

#### Scenario: Invalid method
- **WHEN** `Validate()` is called with `Method` not in `{ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP}`
- **THEN** it SHALL return an error containing "method must be one of"

#### Scenario: PoolCreateRequest Protocol field uses PoolProtocol type
- **WHEN** `PoolCreateRequest.Protocol` is declared
- **THEN** its Go type SHALL be `PoolProtocol` (not `string`)

#### Scenario: PoolCreateRequest Method field uses PoolMethod type
- **WHEN** `PoolCreateRequest.Method` is declared
- **THEN** its Go type SHALL be `PoolMethod` (not `string`)

#### Scenario: PoolUpdateRequest Method field uses PoolMethod type
- **WHEN** `PoolUpdateRequest.Method` is declared
- **THEN** its Go type SHALL be `PoolMethod` (not `string`)
