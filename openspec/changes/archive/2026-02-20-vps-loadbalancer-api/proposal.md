## Why

The VPS SDK currently lacks Load Balancer support. This is the last major resource listed in the TODO and is required to provide complete VPS API coverage. Without it, users must make raw HTTP calls to manage load balancers, listeners, pools, and members — breaking the ergonomic SDK experience offered for all other VPS resources.

## What Changes

- Add `loadbalancers` resource with full CRUD (List, Get, Create, Update, Delete)
- Add `AssociateFloatingIP` operation on LoadBalancer (matches existing server NIC pattern)
- Add `listeners` as a sub-resource of LoadBalancer (List, Get, Create, Update, Delete)
- Add `pools` as a sub-resource of LoadBalancer (List, Get, Create, Update, Delete)
- Add `members` as a sub-sub-resource of Pool (Create, BatchUpdate, Delete)
- Add `LoadBalancerResource` wrapper with `Listeners()` and `Pools()` sub-client accessors
- Add `PoolResource` wrapper with `Members()` sub-client accessor
- Add `WaitForLoadBalancerActive` waiter in VPS core
- Add `LoadBalancers()` accessor to VPS core `Client`
- Add models for all request/response types with validation
- Add unit tests for all modules and models
- Add end-to-end example in `cmd/main.go`

## Capabilities

### New Capabilities

- `vps-loadbalancer`: Load Balancer resource CRUD, AssociateFIP operation, and waiter
- `vps-loadbalancer-listeners`: Listener sub-resource CRUD under LoadBalancer
- `vps-loadbalancer-pools`: Pool sub-resource CRUD under LoadBalancer, with Member sub-sub-resource (Create, BatchUpdate, Delete)

### Modified Capabilities

_None — no existing specs are affected._

## Impact

- **New packages**: `models/vps/loadbalancers/`, `modules/vps/loadbalancers/`
- **Modified files**: `modules/vps/core/client.go` (add `LoadBalancers()` accessor), `modules/vps/core/waiters.go` (add LB waiter), `cmd/main.go` (add e2e example)
- **API paths covered**: All 19 endpoints under `/api/v1/project/{project-id}/loadbalancers/...` from vps.yaml
- **Dependencies**: Reuses existing `common.IDName`, `floatingips.FloatingIP`, `internalhttp.Client`, `waiter` framework
- **No breaking changes** to existing code
