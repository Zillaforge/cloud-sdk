# Research Findings: Fix Network Model Definition

## Decision 1: Align Network JSON tags with Swagger camelCase fields
- **Rationale**: `pb.NetworkInfo` defines fields such as `createdAt`, `updatedAt`, and `status_reason`. Using identical JSON tags ensures seamless (un)marshalling without custom mapping.
- **Alternatives considered**: Retain existing snake_case tags (e.g., `created_at`). Rejected because it would require manual translation and diverge from the contract.

## Decision 2: Represent nested references with shared IDName struct
- **Rationale**: Both projects and users appear across multiple models. Reusing the existing `models/vps/securitygroups.IDName` pattern (ID + Name) keeps consistency and reduces duplication.
- **Alternatives considered**: Introduce bespoke anonymous structs within the Network model. Rejected due to duplication and harder maintenance.

## Decision 3: Keep Router reference lightweight in Network model
- **Rationale**: Network responses embed router metadata via `pb.RouterInfo`. We will declare a compact `RouterInfo` struct within the networks model that mirrors the contract, avoiding cross-package circular dependencies.
- **Alternatives considered**: Reuse the broader `routers.Router` type. Rejected because it uses `time.Time` fields and additional properties absent in `pb.RouterInfo`, leading to incorrect JSON bindings.

## Decision 4: Expand port model to match `NetPort`
- **Rationale**: The SDK exposes router network port helpers. Aligning the `NetworkPort` struct with `NetPort` ensures port IDs, addresses, and embedded server references deserialize correctly for downstream operations.
- **Alternatives considered**: Leave existing minimalist port struct unchanged. Rejected because it would continue omitting `addresses` and `server` relationships required by the Swagger contract.
