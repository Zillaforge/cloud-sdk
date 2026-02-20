## 1. Models

- [x] 1.1 Create `models/vps/loadbalancers/loadbalancer.go` — LoadBalancer struct, LoadBalancerStatus type/constants, LoadBalancerCreateRequest (with Validate), LoadBalancerUpdateRequest, ListLoadBalancersOptions, LoadBalancerListResponse, AssociateFloatingIPRequest, Secret struct
- [x] 1.2 Create `models/vps/loadbalancers/listener.go` — Listener struct, ListenerCreateRequest (with Validate), ListenerUpdateRequest, ListenerListResponse
- [x] 1.3 Create `models/vps/loadbalancers/pool.go` — Pool struct, PoolCreateRequest (with Validate), PoolUpdateRequest, PoolListResponse, Member struct, MemberCreateRequest (with Validate), MembersBatchUpdateRequest, PoolMethod type/constants
- [x] 1.4 Create `models/vps/loadbalancers/loadbalancer_test.go` — JSON roundtrip tests and Validate() tests for LoadBalancerCreateRequest
- [x] 1.5 Create `models/vps/loadbalancers/listener_test.go` — JSON roundtrip tests and Validate() tests for ListenerCreateRequest
- [x] 1.6 Create `models/vps/loadbalancers/pool_test.go` — JSON roundtrip tests and Validate() tests for PoolCreateRequest, MemberCreateRequest

## 2. Module — LoadBalancer Client

- [x] 2.1 Create `modules/vps/loadbalancers/client.go` — Client struct, NewClient, LoadBalancerResource (with Listeners/Pools sub-client accessors), ListenerOperations/PoolOperations interfaces, List/Get/Create/Update/Delete methods, AssociateFloatingIP method
- [x] 2.2 Create `modules/vps/loadbalancers/client_test.go` — Unit tests for all Client methods with httptest mock servers

## 3. Module — Listeners Sub-Resource

- [x] 3.1 Create `modules/vps/loadbalancers/listeners.go` — ListenersClient struct, List/Get/Create/Update/Delete methods
- [x] 3.2 Create `modules/vps/loadbalancers/listeners_test.go` — Unit tests for all ListenersClient methods

## 4. Module — Pools Sub-Resource

- [x] 4.1 Create `modules/vps/loadbalancers/pools.go` — PoolsClient struct, PoolResource (with Members sub-client accessor), MemberOperations interface, List/Get/Create/Update/Delete methods
- [x] 4.2 Create `modules/vps/loadbalancers/pools_test.go` — Unit tests for all PoolsClient methods

## 5. Module — Members Sub-Sub-Resource

- [x] 5.1 Create `modules/vps/loadbalancers/members.go` — MembersClient struct, Create/BatchUpdate/Delete methods
- [x] 5.2 Create `modules/vps/loadbalancers/members_test.go` — Unit tests for all MembersClient methods

## 6. VPS Core Integration

- [x] 6.1 Add `LoadBalancers()` accessor to `modules/vps/core/client.go`
- [x] 6.2 Add `LoadBalancerWaiterConfig`, `WaitForLoadBalancerStatus`, and `WaitForLoadBalancerActive` to `modules/vps/core/waiters.go`
- [x] 6.3 Add waiter tests in `modules/vps/core/waiters_test.go`
- [x] 6.4 Update `modules/vps/core/client_test.go` to cover `LoadBalancers()` accessor

## 7. End-to-End Example

- [x] 7.1 Add loadbalancer lifecycle methods to `cmd/main.go` — create LB with listeners+pools, wait for active, associate FIP, add members, then tear down

## 8. Verification

- [x] 8.1 Run `make check` (fmt + lint + test) to verify all code compiles and tests pass
