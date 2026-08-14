## 1. Inventory and traceability

- [x] 1.1 Build a source-to-capability mapping for all eleven numbered directories.
- [x] 1.2 Read each source `spec.md`, `plan.md`, `data-model.md`, `research.md`, `tasks.md`, checklist, and contract file where present.
- [x] 1.3 Record source requirement, API path/model, compatibility, and test references for each migrated requirement.

## 2. Canonical OpenSpec conversion

- [x] 2.1 Create `openspec/specs/vps-project-api/spec.md` from `001-vps-project-api` and reconcile its full endpoint scope.
- [x] 2.2 Create `openspec/specs/vps-network-model/spec.md` from `002-fix-network-model`.
- [x] 2.3 Create `openspec/specs/vps-flavor-model/spec.md` from `003-fix-flavors-model`.
- [x] 2.4 Create `openspec/specs/vps-security-group-model/spec.md` from `004-fix-security-group-model`.
- [x] 2.5 Create `openspec/specs/vps-server-model/spec.md` from `005-fix-server-model`.
- [x] 2.6 Create `openspec/specs/vps-keypair-model/spec.md` from `006-fix-keypair-model`.
- [x] 2.7 Create `openspec/specs/vps-floating-ip-model/spec.md` from `007-fix-floating-ip-model`.
- [x] 2.8 Create `openspec/specs/vrm-tag-repository/spec.md` from `008-vrm-tag-repository`.
- [x] 2.9 Create `openspec/specs/iam-api/spec.md` from `009-iam-api`.
- [x] 2.10 Create `openspec/specs/vps-volumes/spec.md` from `010-volumes-api`.
- [x] 2.11 Create `openspec/specs/vps-snapshots/spec.md` from `011-vps-snapshot`.

## 3. Validation and parity

- [x] 3.1 Run `openspec validate "migrate-speckit-specs-to-openspec" --type change --strict` and fix syntax, naming, and scenario errors.
- [x] 3.2 Verify every requirement has at least one WHEN/THEN scenario and every source capability has a canonical spec.
- [x] 3.3 Compare API paths, model fields, enum/serialization rules, error behavior, authentication, timeout/retry, and async behavior against Swagger and source contracts.
- [x] 3.4 Add or update contract/unit test plans for any source behavior not represented by the canonical specs.
- [x] 3.5 Review all intentional breaking changes and prepare release-note migration guidance.

## 4. Source disposition

- [x] 4.1 Decide whether each spec-kit supporting artifact is migrated into docs, retained as historical reference, or removed.
- [x] 4.2 Update repository references and contributor guidance to point to OpenSpec canonical specs.
- [x] 4.3 Copy each complete numbered source directory under the mapped OpenSpec capability's `source/` directory, then remove the now-redundant numbered source directories after verification.
- [x] 4.4 Re-run OpenSpec validation and repository checks after source disposition; migrated specs pass strict validation and repository tests/lint pass, while unrelated pre-existing load-balancer spec failures remain outside this change.
