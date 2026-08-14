## Context

The repository has eleven numbered spec-kit directories under `specs/`. Each contains a primary `spec.md` plus varying supporting artifacts such as plans, data models, contracts, research, checklists, and task lists. OpenSpec uses a change-local proposal, one delta spec per capability, design, and tasks, with capability paths independent of numeric ordering.

## Goals / Non-Goals

**Goals:**

- Preserve externally observable requirements and acceptance behavior from every source `spec.md`.
- Establish stable capability names and a one-to-one source mapping that OpenSpec can validate and archive.
- Keep implementation sequencing, source references, compatibility notes, and verification criteria explicit for the later apply workflow.

**Non-Goals:**

- Do not modify Go code, Swagger files, docs, or the existing spec-kit directories in this planning change.
- Do not infer new API behavior where a source document is ambiguous; retain the source contract and flag any unresolved behavior for implementation review.
- Do not convert every supporting spec-kit file into a separate OpenSpec capability when it is explanatory or implementation material.

## Decisions

### Use one migration change with eleven capability deltas

The primary `spec.md` in each numbered directory defines the externally observable contract, so each becomes one OpenSpec capability spec. Numeric prefixes are omitted from capability paths because OpenSpec paths identify durable domains rather than historical execution order. The mapping is recorded in proposal.md and tasks.md.

### Normalize structure, preserve semantics

Each delta uses `## Purpose`, `## ADDED Requirements`, normative SHALL/MUST language, and four-hash scenario headings. User stories and acceptance scenarios are condensed into testable behavior requirements; implementation-specific class names, code layout, and workflow prose stay out of specs. Existing compatibility decisions, such as removal of legacy fields, remain explicit in the relevant requirement.

### Treat supporting artifacts as source material

Plans, data models, contracts, research, checklists, and source task lists will be reviewed during apply. Their API details and test constraints will be reconciled against the new delta specs, while their implementation sequencing is represented by the migration tasks. This avoids creating duplicate or conflicting OpenSpec capabilities.

### Preserve source files until validation succeeds

The apply phase will create or update canonical `openspec/specs/<capability>/spec.md` files first, validate them, and only then decide whether old spec-kit files should be archived, linked, or removed. This permits rollback by leaving the source tree intact until parity is demonstrated.

### Verify semantic parity, not textual identity

Validation will compare source requirements, acceptance scenarios, API paths/models, error behavior, compatibility notes, and test expectations against the OpenSpec deltas. OpenSpec syntax validation is necessary but insufficient because format conversion intentionally changes headings and prose.

## Risks / Trade-offs

- [Requirement compression] → A concise delta may omit detail from a long source spec; apply tasks require per-source traceability review before source cleanup.
- [Capability naming drift] → Stable names may differ from numbered directories; the explicit mapping table and canonical-path checks prevent duplicate capabilities.
- [Conflicting source artifacts] → Supporting contracts or plans may disagree with `spec.md`; treat the authoritative Swagger and the primary requirement as inputs to a documented resolution rather than silently choosing one.
- [Breaking compatibility notes] → Fields intentionally removed by source specs could be reintroduced accidentally; migration review and contract tests must assert the new shape.
- [Partial migration] → A failed validation could leave canonical and source specs out of sync; retain originals and make cleanup the final gated task.

## Migration Plan

1. Re-read every source directory and map all supporting artifacts to its capability.
2. Create canonical OpenSpec specs from the delta files and run `openspec validate --strict`.
3. Reconcile Swagger/contracts and add traceability notes or corrections without changing intended behavior.
4. Run repository documentation/spec checks and review compatibility-impact notes.
5. Archive or remove spec-kit artifacts only after parity and validation pass; retain a migration reference if source history is needed.

Rollback consists of stopping before cleanup and removing only the newly generated canonical OpenSpec specs/change artifacts; the original numbered source directories remain untouched until step 5.
