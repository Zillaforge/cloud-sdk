## Why

Listener `Protocol`, Pool `Protocol`, and Pool `Method` fields currently use raw `string` types in request/response structs, making it possible to pass arbitrary invalid values and providing no discoverability. A typed `PoolMethod` const already exists but is not used consistently — the same pattern should be applied to `ListenerProtocol` and `PoolProtocol`, and all three types should be used everywhere in models and requests.

## What Changes

- Introduce `ListenerProtocol` typed const enum (`HTTP`, `HTTPS`, `TERMINATED_HTTPS`, `TCP`, `UDP`).
- Introduce `PoolProtocol` typed const enum (`HTTP`, `HTTPS`, `TCP`, `UDP`).
- Replace `string` with `ListenerProtocol` in `Listener.Protocol` and `ListenerCreateRequest.Protocol`.
- Replace `string` with `PoolProtocol` in `Pool.Protocol` and `PoolCreateRequest.Protocol`.
- Replace `string` with `PoolMethod` in `Pool.Method`, `PoolCreateRequest.Method`, and `PoolUpdateRequest.Method` (PoolMethod const already defined but unused in structs).
- Add `Valid()` methods for `ListenerProtocol` and `PoolProtocol`.
- Update `Validate()` on `ListenerCreateRequest` and `PoolCreateRequest` to use typed validation.
- Update tests to use typed constants instead of string literals.

## Capabilities

### New Capabilities

_(none — this change only modifies existing capabilities)_

### Modified Capabilities

- `vps-loadbalancer-listeners`: Listener model and `ListenerCreateRequest.Protocol` field changes from `string` to `ListenerProtocol` enum type; validation updated accordingly.
- `vps-loadbalancer-pools`: Pool model fields `Protocol` and `Method` change from `string` to typed enums; `PoolCreateRequest` and `PoolUpdateRequest` updated; validation updated.

## Impact

- **Models**: `models/vps/loadbalancers/listener.go`, `models/vps/loadbalancers/pool.go`
- **Tests**: `models/vps/loadbalancers/listener_test.go`, `models/vps/loadbalancers/pool_test.go`, `modules/vps/loadbalancers/listeners_test.go`, `modules/vps/loadbalancers/pools_test.go`
- **No API contract changes** — JSON serialization values remain identical (e.g. `"HTTP"`, `"TCP"`)
- **No breaking changes for callers** that already use string literals matching const values; callers assigning arbitrary strings will get a compile error (intentional enforcement)
