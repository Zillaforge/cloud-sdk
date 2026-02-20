## ADDED Requirements

### Requirement: PoolProtocol enum type
The SDK SHALL define a `PoolProtocol` typed string constant with a `Valid()` method.

#### Scenario: Valid pool protocols
- **WHEN** `PoolProtocol.Valid()` is called on `PoolProtocolHTTP`, `PoolProtocolHTTPS`, `PoolProtocolTCP`, or `PoolProtocolUDP`
- **THEN** it SHALL return `true`

#### Scenario: Invalid pool protocol
- **WHEN** `PoolProtocol.Valid()` is called on any other value
- **THEN** it SHALL return `false`

---

### Requirement: List pools
The SDK SHALL provide a `List` method on `PoolOperations` that returns all pools for a load balancer.

#### Scenario: List all pools
- **WHEN** `lbResource.Pools().List(ctx)` is called
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools` and return `[]*PoolResource`, each with `Members()` wired

---

### Requirement: Get pool
The SDK SHALL provide a `Get` method on `PoolOperations` that fetches a single pool by ID.

#### Scenario: Get existing pool
- **WHEN** `lbResource.Pools().Get(ctx, poolID)` is called
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}` and return `*PoolResource` with `Members()` wired

#### Scenario: Get non-existent pool
- **WHEN** `lbResource.Pools().Get(ctx, poolID)` returns an HTTP error
- **THEN** the client SHALL propagate the error wrapping both the pool ID and LB ID

---

### Requirement: Create pool
The SDK SHALL provide a `Create` method on `PoolOperations` that adds a pool to a load balancer.

#### Scenario: Create pool with required fields
- **WHEN** `lbResource.Pools().Create(ctx, req)` is called with `Name`, `Protocol`, `MemberPort`, and `Method` set
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools` and return `*PoolResource`

#### Scenario: Create pool with missing name
- **WHEN** `Create(ctx, req)` is called without `Name`
- **THEN** the client SHALL return a validation error before making any HTTP call

#### Scenario: Create pool with missing method
- **WHEN** `Create(ctx, req)` is called without `Method`
- **THEN** the client SHALL return a validation error before making any HTTP call

---

### Requirement: Update pool
The SDK SHALL provide an `Update` method on `PoolOperations` that updates pool attributes.

#### Scenario: Update pool
- **WHEN** `lbResource.Pools().Update(ctx, poolID, req)` is called
- **THEN** the client SHALL issue `PUT /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}` and return `*PoolResource`

---

### Requirement: Delete pool
The SDK SHALL provide a `Delete` method on `PoolOperations` that removes a pool.

#### Scenario: Delete pool
- **WHEN** `lbResource.Pools().Delete(ctx, poolID)` is called
- **THEN** the client SHALL issue `DELETE /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}` and return nil on success

---

### Requirement: PoolOperations interface
The SDK SHALL define a `PoolOperations` interface satisfied by `PoolsClient`.

#### Scenario: Interface compliance
- **WHEN** `PoolsClient` is compiled
- **THEN** it SHALL satisfy `PoolOperations` at compile time (enforced via interface assertion)

---

### Requirement: PoolResource wraps Pool with Members sub-client
The SDK SHALL provide a `PoolResource` struct that embeds `*Pool` and exposes a `Members() MemberOperations` method.

#### Scenario: Access members client from PoolResource
- **WHEN** a caller receives a `*PoolResource` from pool `Get`, `List`, or `Create`
- **THEN** calling `Members()` SHALL return a `MemberOperations` bound to that pool's ID and LB's ID

---

### Requirement: Create member
The SDK SHALL provide a `Create` method on `MemberOperations` that adds a member to a pool.

#### Scenario: Create member with address
- **WHEN** `poolResource.Members().Create(ctx, req)` is called with `Address` set
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}/members` and return `*Pool`

#### Scenario: Create member with missing address
- **WHEN** `Create(ctx, req)` is called without `Address`
- **THEN** the client SHALL return a validation error before making any HTTP call

---

### Requirement: Batch update members
The SDK SHALL provide a `BatchUpdate` method on `MemberOperations` that replaces the full members list.

#### Scenario: Batch update members
- **WHEN** `poolResource.Members().BatchUpdate(ctx, req)` is called with `Members` array
- **THEN** the client SHALL issue `PUT /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}/members` with the JSON body

---

### Requirement: Delete member
The SDK SHALL provide a `Delete` method on `MemberOperations` that removes a member from a pool.

#### Scenario: Delete member
- **WHEN** `poolResource.Members().Delete(ctx, memberID)` is called
- **THEN** the client SHALL issue `DELETE /api/v1/project/{project-id}/loadbalancers/{lb-id}/pools/{pool-id}/members/{member-id}` and return nil on success

---

### Requirement: MemberOperations interface
The SDK SHALL define a `MemberOperations` interface satisfied by `MembersClient`.

#### Scenario: Interface compliance
- **WHEN** `MembersClient` is compiled
- **THEN** it SHALL satisfy `MemberOperations` at compile time (enforced via interface assertion)

---

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

### Requirement: Member model
The SDK SHALL define a `Member` struct matching `pb.MemberInfo` with all fields and JSON tags.

#### Scenario: Member has required fields
- **WHEN** a `Member` struct is populated from API response
- **THEN** it SHALL contain fields: ID, Name, Address, ProtocolPort, Status

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

---

### Requirement: MemberCreateRequest validation
The `MemberCreateRequest` model SHALL validate required fields.

#### Scenario: Valid create request
- **WHEN** `Validate()` is called with non-empty `Address`
- **THEN** it SHALL return nil

#### Scenario: Missing address
- **WHEN** `Validate()` is called with empty `Address`
- **THEN** it SHALL return an error containing "address is required"
