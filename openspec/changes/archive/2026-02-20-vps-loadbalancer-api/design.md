## Context

The SDK currently provides full coverage for VPS resources (servers, networks, floating IPs, keypairs, security groups, volumes, snapshots, shares) with a consistent pattern: models package for types, modules package for operations, resource wrappers for sub-resources, and waiters in VPS core.

Load Balancer is the last missing VPS resource. The API surface (19 endpoints from swagger) features a 3-level hierarchy:
- **LoadBalancer** → top-level CRUD + AssociateFIP
- **Listener** → sub-resource of LoadBalancer (CRUD)
- **Pool** → sub-resource of LoadBalancer (CRUD)
- **Member** → sub-resource of Pool (Create, BatchUpdate, Delete)

The existing `shares` package demonstrates the 2-level sub-resource pattern (Share → Rules). The `servers` package shows multi-sub-resource wrappers (Server → NICs + Volumes). This implementation will combine both approaches and extend to a 3-level hierarchy (LB → Pool → Member).

## Goals / Non-Goals

**Goals:**
- Full SDK coverage for all 19 loadbalancer endpoints
- Follow existing resource/sub-resource patterns (ShareResource, ServerResource)
- `LoadBalancerResource` with `Listeners()` and `Pools()` sub-client accessors
- `PoolResource` with `Members()` sub-client accessor (3-level nesting)
- `AssociateFloatingIP` operation on LoadBalancer (mirrors NIC pattern)
- `WaitForLoadBalancerActive` waiter in VPS core
- Request validation on all create/update types
- Unit tests for all modules and models
- End-to-end example in `cmd/main.go`

**Non-Goals:**
- Admin-level LB endpoints (if any exist)
- TLS certificate management (SecretCreateInput) — exposed as plain struct, no special handling
- Auto-scaling group integration — wire `asg_id` field but no ASG-specific logic
- Retry/backoff on individual LB operations (handled by HTTP client layer)

## Decisions

### 1. Package Structure

**Decision**: Single `loadbalancers` package under both `models/vps/loadbalancers/` and `modules/vps/loadbalancers/`.

**Rationale**: Matches the existing pattern where each resource has its own package (e.g., `shares`, `servers`). All sub-resource code (listeners, pools, members) lives in the same package since they share the LB context.

**Files**:
```
models/vps/loadbalancers/
  loadbalancer.go    # LB model, status, create/update/list types
  listener.go        # Listener model, create/update types
  pool.go            # Pool model, create/update types, member model + create/delete types

modules/vps/loadbalancers/
  client.go          # Client, LoadBalancerResource, AssociateFIP
  listeners.go       # ListenersClient (CRUD)
  pools.go           # PoolsClient, PoolResource (CRUD)
  members.go         # MembersClient (Create, BatchUpdate, Delete)
  client_test.go     # Tests for client
  listeners_test.go  # Tests for listeners
  pools_test.go      # Tests for pools
  members_test.go    # Tests for members
```

**Alternative considered**: Separate packages per sub-resource (`modules/vps/listeners/`). Rejected because the swagger groups them under loadbalancer, and the path hierarchy `/loadbalancers/{id}/listeners/` makes them inherently scoped.

### 2. Resource Wrapper Pattern (3-level)

**Decision**: Use `LoadBalancerResource` wrapping `*LoadBalancer` with `Listeners()` and `Pools()` sub-client accessors. `PoolsClient.Get/Create` returns `*PoolResource` wrapping `*Pool` with `Members()` accessor.

```go
// Level 1: LoadBalancerResource
type LoadBalancerResource struct {
    *loadbalancers.LoadBalancer
    listenerOps ListenerOperations
    poolOps     PoolOperations
}

// Level 2: PoolResource
type PoolResource struct {
    *loadbalancers.Pool
    memberOps MemberOperations
}
```

**Rationale**: Follows `ServerResource` pattern (NICs + Volumes). Extending to `PoolResource` keeps the member operations scoped to a specific pool, matching the API's `/pools/{pool-id}/members` hierarchy.

### 3. AssociateFloatingIP as Method on Client

**Decision**: `AssociateFloatingIP(ctx, lbID, req)` as a method on `loadbalancers.Client`, not on `LoadBalancerResource`.

**Rationale**: This is a direct operation on the LB (POST to `/loadbalancers/{lb-id}/floatingip`), not a sub-resource lifecycle. Making it a client method avoids needing a resource wrapper just for this operation, consistent with how snapshot/extend operations work on shares.

### 4. Waiter for LoadBalancer only

**Decision**: Only add `WaitForLoadBalancerActive` waiter. No separate listener/pool waiters.

**Rationale**: The swagger shows `status` fields on LB, listeners, pools, and members. However, the primary lifecycle waiter needed is for LB creation (which is an async operation). Listener/pool status is typically resolved by the time the LB reaches ACTIVE. If needed later, additional waiters can be added.

### 5. Member operations: Create, BatchUpdate, Delete (no List/Get)

**Decision**: Mirror the API exactly — members have `Create`, `BatchUpdate` (PUT), and `Delete`. No separate `List` or `Get` since the API doesn't expose them (members are embedded in Pool response).

**Rationale**: The swagger only defines POST (create), PUT (batch update), DELETE (delete) for members. Members appear within the `Pool.Members` slice from Get/List pool responses.

## Risks / Trade-offs

**[3-level nesting complexity]** → The Pool → Member pattern adds a third level of resource wrapping. Mitigated by keeping MembersClient simple (3 methods) and following the proven RulesClient pattern from shares.

**[LB status values unknown]** → The swagger doesn't enumerate LB status constants explicitly. Use observed statuses from similar resources: `PENDING_CREATE`, `ACTIVE`, `ERROR`, `PENDING_UPDATE`, `PENDING_DELETE`. Mitigated by using string type and documenting known values.

**[BatchUpdate members replaces entire list]** → The PUT endpoint for members is a full replace (not partial update). SDK users must provide all desired members in the batch. Documented in method godoc.
