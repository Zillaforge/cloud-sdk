## Context

The LB model currently uses raw `string` for `Listener.Protocol`, `ListenerCreateRequest.Protocol`, `Pool.Protocol`, `Pool.Method`, `PoolCreateRequest.Protocol`, `PoolCreateRequest.Method`, and `PoolUpdateRequest.Method`. A `PoolMethod` typed const was introduced previously but the struct fields were never updated to use it. `ListenerProtocol` and `PoolProtocol` types do not yet exist. This creates an inconsistent and unsafe API surface.

## Goals / Non-Goals

**Goals:**
- Introduce `ListenerProtocol` and `PoolProtocol` typed string enums with `Valid()` methods.
- Update all model structs (`Listener`, `ListenerCreateRequest`, `Pool`, `PoolCreateRequest`, `PoolUpdateRequest`) to use the typed enums.
- Update `Validate()` to leverage typed validation (`Valid()`) instead of empty-string checks alone.
- Update all affected tests to use typed constant values instead of raw string literals.

**Non-Goals:**
- Changing the JSON wire format — serialization values remain the same (e.g. `"HTTP"`, `"TCP"`).
- Modifying the HTTP client or module-level code — only model types change.
- Adding new protocol values not already supported by the API.

## Decisions

### Decision: Define types in `pool.go` / `listener.go` alongside related structs
**Rationale**: Keeps enum type, constants, and structs in one file per domain. Mirrors the existing pattern where `PoolMethod` is defined in `pool.go`. Alternative of a shared `enums.go` file was considered but rejected — it offers no advantage when each type belongs clearly to one file.

### Decision: Use typed `string` as the underlying kind (not `iota int`)
**Rationale**: JSON marshaling must match string API values (`"HTTP"`, `"ROUND_ROBIN"`, etc.) transparently. String-based enums serialize/deserialize without custom marshalers. This is the same approach used by `PoolMethod`.

### Decision: Replace `string` fields directly (not a wrapper/pointer)
**Rationale**: No zero-value ambiguity issues since validation already rejects empty protocols. Direct replacement is simplest and keeps backward-compat for callers who already pass the correct string literals as typed values (Go allows untyped string constants to satisfy typed string aliases without explicit conversion when used as literals).

### Decision: Add `Valid()` method to both new types
**Rationale**: Consistent with `PoolMethod.Valid()`. Enables `Validate()` to check not just presence but correctness of value.

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Callers assigning arbitrary `string` variables to `Protocol` fields will get compile errors | This is intentional enforcement; all known internal usages are in tests and are updated in this change |
| Forgetting to update a test file causes build failure | `make build` + `make test` will catch it immediately before merge |

## Migration Plan

1. Add `ListenerProtocol` type + consts to `listener.go`.
2. Add `PoolProtocol` type + consts to `pool.go`; update `Pool.Protocol` and all request fields.
3. Update `Validate()` methods to use `.Valid()` checks.
4. Update `*_test.go` files to use typed constants.
5. Run `make build && make test` to confirm no regressions.
