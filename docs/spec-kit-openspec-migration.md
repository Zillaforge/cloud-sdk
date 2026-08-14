# Spec-Kit to OpenSpec Migration

The repository's canonical capability specifications are maintained under `openspec/specs/`. Each capability contains an OpenSpec `spec.md` plus a `source/` snapshot of the complete spec-kit-generated directory for traceability. The legacy primary `spec.md` is named `source/spec-kit-source.md` so OpenSpec does not mistake it for another canonical capability. The former numbered `specs/` directory has been removed after content verification.

## Source Mapping

| Legacy source | Canonical OpenSpec capability |
| --- | --- |
| `specs/001-vps-project-api` | `openspec/specs/vps-project-api` |
| `specs/002-fix-network-model` | `openspec/specs/vps-network-model` |
| `specs/003-fix-flavors-model` | `openspec/specs/vps-flavor-model` |
| `specs/004-fix-security-group-model` | `openspec/specs/vps-security-group-model` |
| `specs/005-fix-server-model` | `openspec/specs/vps-server-model` |
| `specs/006-fix-keypair-model` | `openspec/specs/vps-keypair-model` |
| `specs/007-fix-floating-ip-model` | `openspec/specs/vps-floating-ip-model` |
| `specs/008-vrm-tag-repository` | `openspec/specs/vrm-tag-repository` |
| `specs/009-iam-api` | `openspec/specs/iam-api` |
| `specs/010-volumes-api` | `openspec/specs/vps-volumes` |
| `specs/011-vps-snapshot` | `openspec/specs/vps-snapshots` |

Supporting files such as plans, research, checklists, contracts, and source task lists are copied under each capability's `source/` directory. Their detailed API and test information must be consulted when implementing or changing the corresponding capability, with Swagger treated as authoritative when source documents conflict.

## Verification Coverage

The legacy test references remain the verification backlog for their mapped capabilities. They cover model JSON round trips, endpoint contracts, validation, authentication, context cancellation, retry/timeout behavior, asynchronous waiters, pagination, concurrency, and integration behavior. During implementation, tests should be added or updated alongside the affected Go package; the migration itself does not alter production behavior.

The source audit identified incomplete or conflicting legacy test plans for server VNC and metrics, keypair timestamp parsing, floating-IP disassociation, volume token-expiry/logging, and snapshot deletion/status semantics. These are intentionally retained as explicit follow-up decisions rather than silently normalized during format conversion.

## Compatibility Notes

- Flavor field and list names, security-group response fields, keypair list shape, and floating-IP list shape contain documented breaking changes from the legacy models.
- Project-scoped clients are authoritative for VPS and VRM APIs; individual methods do not require repeated project IDs.
- Legacy source documents contain unresolved inconsistencies around metric shapes, namespace inputs, timestamp parsing, floating-IP disassociation, and snapshot status codes. These require explicit API-contract decisions before implementation changes.
