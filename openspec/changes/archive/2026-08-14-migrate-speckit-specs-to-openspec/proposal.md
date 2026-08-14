## Why

The repository contains eleven numbered planning directories generated with spec-kit, while the repository now uses OpenSpec for capability specifications and change tracking. The existing documents are valuable source requirements, but their format, metadata, and organization are not directly consumable by OpenSpec tooling, making validation, dependency tracking, and future implementation workflows inconsistent.

## What Changes

- Convert each existing `specs/00x-*/spec.md` into an OpenSpec delta spec under this change, using requirement/scenario sections and preserving the documented behavior.
- Map the numbered feature names to stable kebab-case capability paths without changing API scope or requirements.
- Record the source-directory mapping and any compatibility-sensitive decisions, including the existing breaking-change notes.
- Define a repeatable migration approach for the remaining spec-kit supporting documents (`plan.md`, `tasks.md`, models, contracts, research, and checklists) so their useful material is represented in OpenSpec artifacts or retained as source references.
- Keep the original spec-kit directories unchanged during planning; implementation will perform the eventual repository reorganization and validation.

## Capabilities

### New Capabilities

- `vps-project-api`: Project-scoped VPS APIs from `001-vps-project-api`.
- `vps-network-model`: Swagger-aligned network model from `002-fix-network-model`.
- `vps-flavor-model`: Swagger-aligned flavor model from `003-fix-flavors-model`.
- `vps-security-group-model`: Security group and rule model/module behavior from `004-fix-security-group-model`.
- `vps-server-model`: Server model and module behavior from `005-fix-server-model`.
- `vps-keypair-model`: Keypair model/module behavior from `006-fix-keypair-model`.
- `vps-floating-ip-model`: Floating IP model/module behavior from `007-fix-floating-ip-model`.
- `vrm-tag-repository`: VRM tag and repository APIs from `008-vrm-tag-repository`.
- `iam-api`: IAM user and project APIs from `009-iam-api`.
- `vps-volumes`: Volume and volume-type APIs from `010-volumes-api`.
- `vps-snapshots`: VPS snapshot resource from `011-vps-snapshot`.

### Modified Capabilities

None. The OpenSpec capabilities are being introduced from repository-local spec-kit source documents; this change does not alter the intended API behavior.

## Impact

- Planning artifacts under `openspec/changes/migrate-speckit-specs-to-openspec/` and, during a later apply, the repository's canonical `openspec/specs/` tree.
- OpenSpec validation, status, apply, and archive workflows.
- No production Go code, generated Swagger contract, public API, or runtime dependency changes are authorized by this planning change.
- Existing `specs/` directories remain available as migration source until a separate implementation/apply request approves their disposition.
