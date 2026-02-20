## Context

The VPS SDK follows a consistent two-layer pattern for each resource type:
- `models/vps/<resource>/` — plain Go structs matching the Swagger schema (request inputs, response outputs, constants)
- `modules/vps/<resource>/` — a `Client` struct with methods for each API endpoint, backed by `internal/http.Client`

The VPS core client (`modules/vps/core/client.go`) wires sub-clients together and exposes them via accessor methods.

Shares are a first-class VPS resource with additional sub-resources (access rules). No share SDK code exists yet.

## Goals / Non-Goals

**Goals:**
- Add `models/vps/shares/` package with all request/response types matching `vps.yaml`
- Add `modules/vps/shares/` package with a `Client` covering all 8 share endpoints
- Register `Shares()` accessor on the VPS core client
- Unit tests for models (validation) and client (HTTP round-trip with mock)

**Non-Goals:**
- Admin-only share endpoints (none exist in Swagger)
- Pagination / streaming — API returns full list in one call
- Waiter / polling helpers for async share operations

## Decisions

### Follow existing module conventions exactly
Every existing VPS module (`volumes`, `snapshots`, `networks`, etc.) uses the same layout. Deviating would harm discoverability and maintainability.
**Decision**: `models/vps/shares/share.go` + `modules/vps/shares/client.go`, mirroring the snapshots package.

### Share rules as methods on the share Client (not a separate sub-client)
The security groups module exposes rules as methods directly on `securitygroups.Client`. Shares have a similar pattern (rules are only meaningful in the context of a share).
**Decision**: `AddRule()` and `RemoveRule()` as methods on `shares.Client`, not a separate `rules.Client`.

### Extend as a dedicated method
The API exposes `POST /shares/{id}/extend` separately from the generic update path; it changes size not metadata.
**Decision**: `Extend(ctx, shareID, req)` method — keeps intent clear, avoids overloading `Update`.

### Status constants
Follow the `VolumeStatus` / `SnapshotStatus` pattern — typed string constants for all known share states.

## Risks / Trade-offs

- [Share rule add returns 204 No Content] → Return only `error`; no body to parse (same as Delete).
- [API may evolve] → All fields in models are kept optional except those marked `required` in Swagger, so new optional fields can be added without breaking callers.

## Migration Plan

Purely additive — no existing code changes except `core/client.go`. No migration needed.
