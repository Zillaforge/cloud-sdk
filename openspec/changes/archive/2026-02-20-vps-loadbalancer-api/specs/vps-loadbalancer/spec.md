## ADDED Requirements

### Requirement: List load balancers
The SDK SHALL provide a `List` method on `loadbalancers.Client` that returns all load balancers for a project as `[]*LoadBalancerResource`, supporting optional filters.

#### Scenario: List all load balancers
- **WHEN** `List(ctx, nil)` is called
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/loadbalancers` and return a slice of `*LoadBalancerResource`, each with `Listeners()` and `Pools()` wired

#### Scenario: List load balancers with name filter
- **WHEN** `List(ctx, &ListLoadBalancersOptions{Name: "my-lb"})` is called
- **THEN** the client SHALL append `?name=my-lb` to the request URL

#### Scenario: List load balancers with detail
- **WHEN** `List(ctx, &ListLoadBalancersOptions{Detail: true})` is called
- **THEN** the client SHALL append `?detail=true` to the request URL

---

### Requirement: Get load balancer
The SDK SHALL provide a `Get` method that fetches a single load balancer by ID and returns it wrapped in a `LoadBalancerResource`.

#### Scenario: Get existing load balancer
- **WHEN** `Get(ctx, lbID)` is called with a valid ID
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/loadbalancers/{lb-id}` and return `*LoadBalancerResource` with `Listeners()` and `Pools()` wired

#### Scenario: Get non-existent load balancer
- **WHEN** `Get(ctx, lbID)` returns an HTTP error from the server
- **THEN** the client SHALL propagate the error wrapping the load balancer ID

---

### Requirement: Create load balancer
The SDK SHALL provide a `Create` method that creates a new load balancer and returns it as a `LoadBalancerResource`.

#### Scenario: Create load balancer with required fields
- **WHEN** `Create(ctx, req)` is called with `Name` and `NetworkID` set
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/loadbalancers` with the JSON body and return `*LoadBalancerResource`

#### Scenario: Create load balancer with listeners
- **WHEN** `Create(ctx, req)` is called with `Listeners` array populated
- **THEN** the client SHALL include the nested listener/pool/member structure in the POST body

#### Scenario: Create load balancer with missing name
- **WHEN** `Create(ctx, req)` is called without `Name`
- **THEN** the client SHALL return a validation error before making any HTTP call

#### Scenario: Create load balancer with missing network_id
- **WHEN** `Create(ctx, req)` is called without `NetworkID`
- **THEN** the client SHALL return a validation error before making any HTTP call

---

### Requirement: Update load balancer
The SDK SHALL provide an `Update` method that updates load balancer metadata (name, description).

#### Scenario: Update load balancer name
- **WHEN** `Update(ctx, lbID, req)` is called with `Name` set
- **THEN** the client SHALL issue `PUT /api/v1/project/{project-id}/loadbalancers/{lb-id}` and return updated `*LoadBalancerResource`

---

### Requirement: Delete load balancer
The SDK SHALL provide a `Delete` method that removes a load balancer.

#### Scenario: Delete existing load balancer
- **WHEN** `Delete(ctx, lbID)` is called
- **THEN** the client SHALL issue `DELETE /api/v1/project/{project-id}/loadbalancers/{lb-id}` and return nil on success

---

### Requirement: Associate floating IP to load balancer
The SDK SHALL provide an `AssociateFloatingIP` method on `loadbalancers.Client` that associates a floating IP to the LB's VIP.

#### Scenario: Associate with existing floating IP
- **WHEN** `AssociateFloatingIP(ctx, lbID, &AssociateFloatingIPRequest{FIPID: "fip-123"})` is called
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/loadbalancers/{lb-id}/floatingip` with `{"fip_id": "fip-123"}` and return `*floatingips.FloatingIP`

#### Scenario: Associate with new floating IP (empty FIPID)
- **WHEN** `AssociateFloatingIP(ctx, lbID, &AssociateFloatingIPRequest{})` is called without `FIPID`
- **THEN** the client SHALL issue the POST with empty body and the server creates a new floating IP

---

### Requirement: LoadBalancerResource wraps LoadBalancer with sub-clients
The SDK SHALL provide a `LoadBalancerResource` struct that embeds `*LoadBalancer` and exposes `Listeners()` and `Pools()` methods returning scoped sub-clients.

#### Scenario: Access listeners client from LoadBalancerResource
- **WHEN** a caller receives a `*LoadBalancerResource` from `Get`, `List`, or `Create`
- **THEN** calling `Listeners()` SHALL return a `ListenerOperations` bound to that LB's ID

#### Scenario: Access pools client from LoadBalancerResource
- **WHEN** a caller receives a `*LoadBalancerResource` from `Get`, `List`, or `Create`
- **THEN** calling `Pools()` SHALL return a `PoolOperations` bound to that LB's ID

---

### Requirement: LoadBalancer model and status
The SDK SHALL define a `LoadBalancer` struct matching `pb.LbInfo` with all fields and JSON tags, and `LoadBalancerStatus` type with known status constants.

#### Scenario: LoadBalancer has required fields
- **WHEN** a `LoadBalancer` struct is populated from API response
- **THEN** it SHALL contain fields: ID, Name, Description, Status, StatusReason, VIP, PortID, NetworkID, Network, ProjectID, Project, UserID, User, Namespace, FloatingIP, Listeners, CreatedAt, UpdatedAt

---

### Requirement: LoadBalancer request validation
The `LoadBalancerCreateRequest` model SHALL validate required fields.

#### Scenario: Valid create request
- **WHEN** `Validate()` is called with non-empty `Name` and `NetworkID`
- **THEN** it SHALL return nil

#### Scenario: Missing name
- **WHEN** `Validate()` is called with empty `Name`
- **THEN** it SHALL return an error containing "name is required"

#### Scenario: Missing network_id
- **WHEN** `Validate()` is called with empty `NetworkID`
- **THEN** it SHALL return an error containing "network_id is required"

---

### Requirement: VPS core client exposes load balancers
The VPS core `Client` SHALL expose a `LoadBalancers()` method returning `*loadbalancers.Client`.

#### Scenario: Access load balancers client
- **WHEN** `vpsClient.LoadBalancers()` is called
- **THEN** it SHALL return a non-nil `*loadbalancers.Client` bound to the same project ID

---

### Requirement: Wait for load balancer active
The VPS core SHALL provide a `WaitForLoadBalancerActive` convenience function that polls until the LB reaches ACTIVE status.

#### Scenario: Wait for LB to become active
- **WHEN** `WaitForLoadBalancerActive(ctx, client, lbID)` is called
- **THEN** it SHALL poll `loadbalancers.Client.Get()` until status is `ACTIVE`

#### Scenario: LB enters ERROR state while waiting
- **WHEN** the LB status becomes `ERROR` during polling
- **THEN** the waiter SHALL return an error immediately

#### Scenario: Timeout waiting for LB
- **WHEN** the max wait duration is exceeded
- **THEN** the waiter SHALL return a timeout error
