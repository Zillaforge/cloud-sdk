# Implementation Tasks: vps snapshot resource

**Feature**: vps snapshot resource (spec: `specs/011-vps-snapshot/spec.md`)  
**Plan**: `specs/011-vps-snapshot/plan.md`  

## Summary

Tasks below are organized by user story from highest priority to lowest. Tests are written first (TDD). Each user story is independently testable, and core dependencies are in the Foundational phase.

## Phase 1 — Setup

- [x] T001 Create `models/vps/snapshots` directory and add `snapshot.go` file with package header `package snapshots` — `models/vps/snapshots/snapshot.go`
 - [ ] T001 Create `models/vps/snapshots` directory and add `snapshot.go` file with package header `package snapshots` — `models/vps/snapshots/snapshot.go`
 - [x] T001a [TDD] Add failing unit tests for CreateSnapshot (happy-path + invalid `volume_id`) before implementing Create() — `modules/vps/snapshots/client_test.go` and repository-level contract tests
- [x] T002 Add `modules/vps/snapshots` directory and client skeleton file `client.go` with package header `package snapshots` — `modules/vps/snapshots/client.go`
- [x] T003 [P] Add unit test package folders and placeholders: `modules/vps/snapshots/client_test.go` and `models/vps/snapshots/snapshot_test.go` — `modules/vps/snapshots/client_test.go`, `models/vps/snapshots/snapshot_test.go`

## Phase 2 — Foundational (blocking prerequisites)

- [x] T004 Create `Snapshot` Go struct in `models/vps/snapshots/snapshot.go` using Swagger field names, and import `models/vps/common` for `Project/User` types. Add `json` tags matching `swagger/vps.yaml`. — `models/vps/snapshots/snapshot.go`
- [x] T005 Implement `SnapshotStatus` custom type and constants; add `Validate()` and `String()` helpers for enum conversion. — `models/vps/snapshots/snapshot.go`
- [x] T006 Implement `CreateSnapshotRequest` and `SnapshotResponse` types (Request/Response naming) and add `Validate()` for create. — `models/vps/snapshots/snapshot.go`
 - [x] T007 [P] Add time parsing helpers and mapping functions for `createdAt`/`updatedAt` (RFC3339) and JSON marshalling tests. — `models/vps/snapshots/time_helpers_test.go`
- [x] T008 Add model conversion functions to map `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` to `models/vps/snapshots.Snapshot` and vice versa. — `models/vps/snapshots/convert.go` (TODO: implement mappings)

## Phase 3 — User Story 1: Create Snapshot (P1)

 - [ ] T009 [US1] Add contract tests to assert `POST /snapshots` returns `status=creating`, contains `project`, `user`, `namespace`, and createdAt. — repository-level contract tests
- [x] T010 [US1] Implement `Create(ctx, req *CreateSnapshotRequest) (*Snapshot, error)` in `modules/vps/snapshots/client.go` following `modules/vps/volumes` pattern — `modules/vps/snapshots/client.go`
- [x] T011 [US1] Unit test for Create() : Validate request body JSON and response mapping to `Snapshot` (mock server). — `modules/vps/snapshots/client_test.go`
 - [ ] T012 [US1] Integration/contract test: add and run repository-level contract tests for create success and error cases (invalid `volume_id` and cross-project)

## Phase 4 — User Story 2: List & Inspect Snapshots (P1)

 - [ ] T013 [US2] Add contract tests for `GET /snapshots` and `GET /snapshots/{id}` expecting `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` fields including `namespace`, `project`, `user`, etc. — repository-level contract tests
 - [ ] T013 [US2] Add contract tests for `GET /snapshots` and `GET /snapshots/{id}` expecting `pegasus-cloud_com_aes_virtualplatformserviceclient_pb.SnapshotInfo` fields including `namespace`, `project`, `user`, etc. — repository-level contract tests
 - [x] T013a [US2] Add List filter tests (name, volume_id, user_id, status) and pagination tests to repository-level contract tests and client tests. — repository-level contract tests, `modules/vps/snapshots/client_test.go`
 - [x] T014 [US2] Implement `List(ctx, opts *ListSnapshotsOptions) ([]*Snapshot, error)` in `modules/vps/snapshots/client.go`. Ensure mapping from `pb.SnapshotListOutput` to `[]*Snapshot`. — `modules/vps/snapshots/client.go`
 - [x] T015 [US2] Implement `Get(ctx, id string) (*Snapshot, error)` and add unit tests for List/Get mapping. — `modules/vps/snapshots/client.go`, `modules/vps/snapshots/client_test.go`

## Phase 5 — User Story 3: Update & Delete Snapshot (P2)

 - [ ] T016 [US3] Add contract tests for `PUT /snapshots/{id}` and `DELETE /snapshots/{id}` — repository-level contract tests
 - [x] T017 [US3] Implement `Update(ctx, id string, req *UpdateSnapshotRequest) (*Snapshot, error)` with validation and unit tests. — `modules/vps/snapshots/client.go`, `modules/vps/snapshots/client_test.go`
 - [ ] T018 [US3] Implement `Delete(ctx, id string) error` with integration/contract tests to assert deletion does not affect volumes previously created. — `modules/vps/snapshots/client.go` and repository-level contract tests

## Phase 6 — Cross‑cutting & Polish

- [ ] T019 Update `docs/vps/` or `docs/` README to include `CreateSnapshot`/`List`/`Get` examples referencing `quickstart.md` — `docs/vps/README.md` or `specs/011-vps-snapshot/quickstart.md`
- [ ] T020 [P] Add exhaustive unit tests to ensure at least 85% coverage for `models/vps/snapshots` and `modules/vps/snapshots` — coverage target measured by `go test -cover` — `models/vps/snapshots/*`, `modules/vps/snapshots/*`
- [ ] T021 Add migration notes and changelog entry for new Snapshot APIs if needed — `docs/CHANGES.md`
 - [ ] T022 Add SDK helper `WaitForSnapshot(ctx, id, timeout)` used by tests and clients, and add unit + contract tests for status transition from `creating` to `available` — `modules/vps/snapshots/client.go`, `modules/vps/snapshots/client_test.go`
 - [ ] T023 Add performance/load tests to validate List() meets SC-003 (95% within 1s at 1k snapshots) — `tests/perf/snapshots_test.go`
 - [ ] T024 Update repository-level contract tests to include contract cases for cross-project rejection, invalid `volume_id`, filter param behavior, and ensure `namespace` is present in responses. — repository-level contract tests

## Dependencies & Order

1. Setup tasks (T001–T003) must be completed before Foundational tasks (T004–T008).
2. Foundational tasks must be completed before User Story tasks for Create/List/Get (T009–T015).
3. Update/Delete story tasks may be implemented after Create/List/Get are available (T016–T018).

## Parallel execution opportunities

- [P] T003 and T007 can run in parallel (test scaffolding and time helpers).  
- [P] T010 and T014 (client Create and List methods) can be developed in parallel after conversion helpers exist (T008).  
- [P] T019 and T021 are documentation tasks and can run concurrently with tests finishing.

## Independent Tests / Acceptance criteria (per user story)

- US1: Contract test shows POST returns 201 with `status=creating`; subsequent GET returns 200 and shows `status=available` when ready. Unit tests validate request JSON (vol ID, name) and response mapping.  
- US2: Contract test shows GET /snapshots returns an indexable list (`snapshots` array); tests check `namespace`, `project`, `user`, `createdAt` fields; unit tests validate mapping to `[]*Snapshot`. 
- US3: Contract tests for PUT and DELETE succeed; deletion does not impact previously created volumes.

## Notes

- All types should follow naming: `CreateSnapshotRequest`, `UpdateSnapshotRequest`, `SnapshotResponse`.  
- `List()` returns `[]*Snapshot` and must be indexable.  
- Time formats must use RFC3339 and Go `time.Time` semantics in JSON mapping.  
- Use `common.IDName` for `project` and `user` fields where available.
