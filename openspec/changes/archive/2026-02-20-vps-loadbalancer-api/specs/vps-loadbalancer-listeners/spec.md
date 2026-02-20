## ADDED Requirements

### Requirement: List listeners
The SDK SHALL provide a `List` method on `ListenerOperations` that returns all listeners for a load balancer.

#### Scenario: List all listeners
- **WHEN** `lbResource.Listeners().List(ctx)` is called
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners` and return `[]*Listener`

---

### Requirement: Get listener
The SDK SHALL provide a `Get` method on `ListenerOperations` that fetches a single listener by ID.

#### Scenario: Get existing listener
- **WHEN** `lbResource.Listeners().Get(ctx, listenerID)` is called
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners/{listener-id}` and return `*Listener`

#### Scenario: Get non-existent listener
- **WHEN** `lbResource.Listeners().Get(ctx, listenerID)` returns an HTTP error
- **THEN** the client SHALL propagate the error wrapping both the listener ID and LB ID

---

### Requirement: Create listener
The SDK SHALL provide a `Create` method on `ListenerOperations` that adds a listener to a load balancer.

#### Scenario: Create listener with required fields
- **WHEN** `lbResource.Listeners().Create(ctx, req)` is called with `Protocol` and `ProtocolPort` set
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners` and return `*Listener`

#### Scenario: Create listener with missing protocol
- **WHEN** `Create(ctx, req)` is called without `Protocol`
- **THEN** the client SHALL return a validation error before making any HTTP call

#### Scenario: Create listener with missing protocol_port
- **WHEN** `Create(ctx, req)` is called without `ProtocolPort`
- **THEN** the client SHALL return a validation error before making any HTTP call

---

### Requirement: Update listener
The SDK SHALL provide an `Update` method on `ListenerOperations` that updates listener attributes.

#### Scenario: Update listener name
- **WHEN** `lbResource.Listeners().Update(ctx, listenerID, req)` is called
- **THEN** the client SHALL issue `PUT /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners/{listener-id}` and return `*Listener`

---

### Requirement: Delete listener
The SDK SHALL provide a `Delete` method on `ListenerOperations` that removes a listener.

#### Scenario: Delete listener
- **WHEN** `lbResource.Listeners().Delete(ctx, listenerID)` is called
- **THEN** the client SHALL issue `DELETE /api/v1/project/{project-id}/loadbalancers/{lb-id}/listeners/{listener-id}` and return nil on success

---

### Requirement: ListenerOperations interface
The SDK SHALL define a `ListenerOperations` interface satisfied by `ListenersClient`.

#### Scenario: Interface compliance
- **WHEN** `ListenersClient` is compiled
- **THEN** it SHALL satisfy `ListenerOperations` at compile time (enforced via interface assertion)

---

### Requirement: Listener model
The SDK SHALL define a `Listener` struct matching `pb.ListenerInfo` with all fields and JSON tags.

#### Scenario: Listener has required fields
- **WHEN** a `Listener` struct is populated from API response
- **THEN** it SHALL contain fields: ID, Name, LbID, Protocol, ProtocolPort, PoolID, Pool, Status, StatusReason, AllowedCIDRs, InsertHeaders, DefaultTLSContainerRef, TimeoutClientData, TimeoutMemberConnect, TimeoutMemberData, TimeoutTCPInspect

---

### Requirement: ListenerCreateRequest validation
The `ListenerCreateRequest` model SHALL validate required fields.

#### Scenario: Valid create request
- **WHEN** `Validate()` is called with non-empty `Protocol` and non-zero `ProtocolPort`
- **THEN** it SHALL return nil

#### Scenario: Missing protocol
- **WHEN** `Validate()` is called with empty `Protocol`
- **THEN** it SHALL return an error containing "protocol is required"

#### Scenario: Missing protocol_port
- **WHEN** `Validate()` is called with zero `ProtocolPort`
- **THEN** it SHALL return an error containing "protocol_port is required"
