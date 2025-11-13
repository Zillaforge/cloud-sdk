# Tasks: Volumes API Client

**Input**: Design documents from `/workspaces/cloud-sdk/specs/008-volumes-api/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are MANDATORY for public APIs per the Constitution. Write tests FIRST
and ensure they FAIL before implementation. Include unit tests and contract tests
derived from Swagger/OpenAPI for each endpoint wrapper.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create directory structure: `models/vps/volumes/`, `models/vps/volumetypes/`, `modules/vps/volumes/`, `modules/vps/volumetypes/`
- [X] T002 Verify internal dependencies are accessible: `internal/http`, `internal/types`, `internal/backoff`, `models/vps/common`
- [X] T003 [P] Review existing VPS module patterns in `modules/vps/flavors/` and `modules/vps/keypairs/` for consistency

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core shared types that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 [P] Define VolumeStatus custom type and constants in `models/vps/volumes/volume.go`
- [X] T005 [P] Define VolumeAction custom type and constants in `models/vps/volumes/volume.go`
- [X] T006 [P] Define Volume struct (16 fields) in `models/vps/volumes/volume.go`
- [X] T007 [P] Define response wrapper structs (VolumeListResponse, VolumeResponse) in `models/vps/volumes/volume.go`
- [X] T008 [P] Define VolumeTypeListResponse struct in `models/vps/volumetypes/volumetype.go`
- [X] T009 Write unit tests for Volume struct JSON marshaling/unmarshaling in `models/vps/volumes/volume_test.go` (TDD - write first, ensure they fail)
- [X] T010 Write unit tests for VolumeStatus constants in `models/vps/volumes/volume_test.go` (TDD - write first, ensure they fail)
- [X] T011 Write unit tests for VolumeTypeListResponse JSON marshaling in `models/vps/volumetypes/volumetype_test.go` (TDD - write first, ensure they fail)
- [X] T012 Implement Volume struct to pass T009 tests in `models/vps/volumes/volume.go`
- [X] T013 Implement VolumeStatus constants to pass T010 tests in `models/vps/volumes/volume.go`
- [X] T014 Implement VolumeTypeListResponse to pass T011 tests in `models/vps/volumetypes/volumetype.go`
- [X] T015 Run `make check` to verify foundational models and tests pass

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - List Available Volume Types (Priority: P1) 🎯 MVP

**Goal**: Enable developers to discover available storage types before volume creation

**Independent Test**: Call VolumeTypes().List(ctx) and verify response contains []string of storage types

**Dependencies**: Requires Phase 2 complete (VolumeTypeListResponse model)

### Tests for User Story 1 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [X] T016 [P] [US1] Write contract test for GET /volume_types endpoint in `modules/vps/volumetypes/client_test.go` - verify response schema matches swagger, test 200/400/500 status codes
- [X] T017 [P] [US1] Write unit test for VolumeTypes().List() success case in `modules/vps/volumetypes/client_test.go` - mock HTTP client, verify returned []string
- [X] T018 [P] [US1] Write unit test for VolumeTypes().List() error cases in `modules/vps/volumetypes/client_test.go` - test network errors, 500 errors, auth errors
- [X] T019 [P] [US1] Write unit test for context cancellation in `modules/vps/volumetypes/client_test.go` - verify proper cleanup

### Implementation for User Story 1

- [X] T020 [US1] Create VolumeTypes client struct with NewClient constructor in `modules/vps/volumetypes/client.go` - Client struct should have fields: baseClient (*internalhttp.Client), projectID (string), basePath (string). NewClient(baseClient, projectID) stores projectID and builds basePath = "/api/v1/project/" + projectID. Follows pattern from modules/vps/flavors/client.go.
- [X] T021 [US1] Implement List(ctx) method in `modules/vps/volumetypes/client.go` - use internal HTTP client, handle VolumeTypeListResponse unwrapping, wrap errors with fmt.Errorf
- [X] T022 [US1] Add VolumeTypes() accessor method to VPS client in `modules/vps/client.go` - Return volumetypes.NewClient(c.baseClient, c.projectID). Follows pattern of existing Flavors(), Keypairs() accessor methods.
- [X] T023 [US1] Run tests T016-T019 and verify they pass
- [X] T024 [US1] Run `make check` to verify linting, formatting, and all tests pass

**Checkpoint**: At this point, User Story 1 should be fully functional - developers can list available volume types

---

## Phase 4: User Story 2 - Create and Manage Volumes (Priority: P1) 🎯 MVP

**Goal**: Enable complete volume lifecycle management (Create, Update, Delete)

**Independent Test**: Create volume, update metadata, delete volume - verify each operation works independently

**Dependencies**: Requires Phase 2 complete (Volume model, request/response types)

### Request/Response Models for User Story 2

- [X] T025 [P] [US2] Define CreateVolumeRequest struct in `models/vps/volumes/volume.go`
- [X] T026 [P] [US2] Define UpdateVolumeRequest struct in `models/vps/volumes/volume.go`
- [X] T027 [P] [US2] Write unit tests for CreateVolumeRequest validation in `models/vps/volumes/volume_test.go` (TDD - write first)
- [X] T028 [P] [US2] Write unit tests for UpdateVolumeRequest validation in `models/vps/volumes/volume_test.go` (TDD - write first)
- [X] T029 [US2] Implement request struct validation to pass T027-T028 tests

### Tests for User Story 2 (MANDATORY - TDD) ⚠️

- [X] T030 [P] [US2] Write contract test for POST /volumes (create) in `modules/vps/volumes/client_test.go` - verify 201 response, required fields, quota errors
- [X] T031 [P] [US2] Write contract test for PUT /volumes/{id} (update) in `modules/vps/volumes/client_test.go` - verify 200 response, partial updates work
- [X] T032 [P] [US2] Write contract test for DELETE /volumes/{id} (delete) in `modules/vps/volumes/client_test.go` - verify 204 response, attached volume error
- [X] T033 [P] [US2] Write unit test for Create() method success case in `modules/vps/volumes/client_test.go` - mock HTTP, verify *Volume returned
- [X] T034 [P] [US2] Write unit test for Create() method error cases in `modules/vps/volumes/client_test.go` - quota exceeded, invalid type, auth errors
- [X] T035 [P] [US2] Write unit test for Update() method in `modules/vps/volumes/client_test.go` - test name update, description update, partial updates
- [X] T036 [P] [US2] Write unit test for Delete() method in `modules/vps/volumes/client_test.go` - test success (204), volume in use error, 404 error
- [X] T037 [P] [US2] Write integration test for full lifecycle in `modules/vps/volumes/client_test.go` - create → update → delete

### Implementation for User Story 2

- [X] T038 [US2] Create Volumes client struct with NewClient constructor in `modules/vps/volumes/client.go` - Client struct should have fields: baseClient (*internalhttp.Client), projectID (string), basePath (string). NewClient(baseClient, projectID) stores projectID and builds basePath = "/api/v1/project/" + projectID. Follows pattern from modules/vps/flavors/client.go.
- [X] T039 [US2] Implement Create(ctx, *CreateVolumeRequest) method in `modules/vps/volumes/client.go` - POST endpoint, handle response, wrap errors
- [X] T040 [US2] Implement Update(ctx, volumeID, *UpdateVolumeRequest) method in `modules/vps/volumes/client.go` - PUT endpoint, handle partial updates
- [X] T041 [US2] Implement Delete(ctx, volumeID) method in `modules/vps/volumes/client.go` - DELETE endpoint, handle 204 response
- [X] T042 [US2] Add Volumes() accessor method to VPS client in `modules/vps/client.go` - Return volumes.NewClient(c.baseClient, c.projectID). Follows pattern of existing Flavors(), Keypairs() accessor methods.
- [X] T043 [US2] Run tests T030-T037 and verify they pass
- [X] T044 [US2] Run `make check` to verify linting, formatting, and all tests pass

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - volume CRUD operations functional

---

## Phase 5: User Story 3 - List and Retrieve Volume Information (Priority: P2)

**Goal**: Enable volume discovery with filtering and detailed information retrieval

**Independent Test**: Create volumes with different attributes, list with filters, retrieve individual details

**Dependencies**: Requires Phase 2 complete (Volume model), can work with or without User Story 2

### Query Options for User Story 3

- [X] T045 [P] [US3] Define ListVolumesOptions struct in `models/vps/volumes/volume.go`
- [X] T046 [P] [US3] Write unit tests for ListVolumesOptions in `models/vps/volumes/volume_test.go` - test filter combinations

### Tests for User Story 3 (MANDATORY - TDD) ⚠️

- [X] T047 [P] [US3] Write contract test for GET /volumes (list) in `modules/vps/volumes/client_test.go` - verify VolumeListResponse unwrapping, query parameters work
- [X] T048 [P] [US3] Write contract test for GET /volumes/{id} (get) in `modules/vps/volumes/client_test.go` - verify all Volume fields present, 404 handling
- [X] T049 [P] [US3] Write unit test for List() with no filters in `modules/vps/volumes/client_test.go` - mock HTTP, verify []*Volume returned
- [X] T050 [P] [US3] Write unit test for List() with name filter in `modules/vps/volumes/client_test.go` - verify query param added
- [X] T051 [P] [US3] Write unit test for List() with status filter in `modules/vps/volumes/client_test.go` - verify filtering works
- [X] T052 [P] [US3] Write unit test for List() with detail=true in `modules/vps/volumes/client_test.go` - verify attachments included
- [X] T053 [P] [US3] Write unit test for Get() method success case in `modules/vps/volumes/client_test.go` - verify *Volume returned with all fields
- [X] T054 [P] [US3] Write unit test for Get() method error cases in `modules/vps/volumes/client_test.go` - 404, invalid ID, auth errors

### Implementation for User Story 3

- [X] T055 [US3] Implement List(ctx, *ListVolumesOptions) method in `modules/vps/volumes/client.go` - GET endpoint with query params, unwrap VolumeListResponse
- [X] T056 [US3] Implement Get(ctx, volumeID) method in `modules/vps/volumes/client.go` - GET endpoint with volume ID, return *Volume
- [X] T057 [US3] Run tests T047-T054 and verify they pass
- [X] T058 [US3] Run `make check` to verify all tests pass with 85%+ coverage

**Checkpoint**: All core volume operations (US1, US2, US3) should now be independently functional

---

## Phase 6: User Story 4 - Perform Volume Actions (Priority: P2)

**Goal**: Enable attach, detach, extend, revert operations on volumes

**Independent Test**: Create volume, attach to server, extend size, detach, verify state changes

**Dependencies**: Requires Phase 2 complete (Volume model, VolumeAction type), ideally after User Story 2 (to create volumes for testing)

### Action Request Model for User Story 4

- [X] T059 [P] [US4] Define VolumeActionRequest struct in `models/vps/volumes/volume.go`
- [X] T060 [P] [US4] Write unit tests for VolumeActionRequest validation in `models/vps/volumes/volume_test.go` - test action-specific parameter requirements
- [X] T061 [US4] Implement VolumeActionRequest validation (attach requires ServerID, extend requires NewSize)

### Tests for User Story 4 (MANDATORY - TDD) ⚠️

- [X] T062 [P] [US4] Write contract test for POST /volumes/{id}/action in `modules/vps/volumes/client_test.go` - verify 202 response, all 4 action types
- [X] T063 [P] [US4] Write unit test for Action() with attach in `modules/vps/volumes/client_test.go` - verify ServerID required, request format correct
- [X] T064 [P] [US4] Write unit test for Action() with detach in `modules/vps/volumes/client_test.go` - verify detach logic
- [X] T065 [P] [US4] Write unit test for Action() with extend in `modules/vps/volumes/client_test.go` - verify NewSize required, size increase only
- [X] T066 [P] [US4] Write unit test for Action() with revert in `modules/vps/volumes/client_test.go` - verify snapshot requirement
- [X] T067 [P] [US4] Write unit test for Action() error cases in `modules/vps/volumes/client_test.go` - invalid action, missing params, volume not found
- [X] T068 [P] [US4] Write integration test for attach/extend/detach workflow in `modules/vps/volumes/client_test.go`

### Implementation for User Story 4

- [X] T069 [US4] Implement Action(ctx, volumeID, *VolumeActionRequest) method in `modules/vps/volumes/client.go` - POST endpoint, handle async 202 response
- [X] T070 [US4] Run tests T062-T068 and verify they pass
- [X] T071 [US4] Run `make check` to verify all tests pass

**Checkpoint**: All volume management features (US1-US4) should now be fully functional

---

## Phase 7: User Story 5 - Create Volumes from Snapshots (Priority: P3)

**Goal**: Enable volume creation from existing snapshots for restore/clone workflows

**Independent Test**: Create volume with snapshot_id parameter, verify volume created with snapshot data

**Dependencies**: Requires Phase 2 complete and User Story 2 (CreateVolumeRequest already supports SnapshotID field)

### Tests for User Story 5 (MANDATORY - TDD) ⚠️

- [X] T072 [P] [US5] Write contract test for POST /volumes with snapshot_id in `modules/vps/volumes/client_test.go` - verify snapshot parameter accepted
- [X] T073 [P] [US5] Write unit test for Create() with SnapshotID in `modules/vps/volumes/client_test.go` - verify snapshot_id passed correctly
- [X] T074 [P] [US5] Write unit test for invalid snapshot errors in `modules/vps/volumes/client_test.go` - test invalid snapshot ID, cross-project snapshot

### Implementation for User Story 5

- [X] T075 [US5] Verify Create() method handles SnapshotID field correctly in `modules/vps/volumes/client.go` (should already work from US2)
- [X] T076 [US5] Run tests T072-T074 and verify they pass
- [X] T077 [US5] Add snapshot creation example to documentation in `specs/008-volumes-api/quickstart.md`

**Checkpoint**: All user stories (US1-US5) should now be independently functional

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T078 [P] Add code examples to `modules/vps/volumes/README.md` - document all methods with usage examples
- [ ] T079 [P] Add code examples to `modules/vps/volumetypes/README.md` - document List method
- [ ] T080 [P] Update `modules/vps/EXAMPLES.md` with volume management workflows
- [ ] T081 [P] Add error handling examples in documentation - show how to detect quota errors, auth errors, volume-in-use errors
- [ ] T082 [P] Add logging examples in documentation - show DEBUG level HTTP logging
- [ ] T083 Code cleanup and refactoring - ensure consistent error messages, remove any dead code
- [ ] T084 [P] Add performance benchmarks in `modules/vps/volumes/client_test.go` - measure typical operation times
- [ ] T085 [P] Add edge case tests in `modules/vps/volumes/client_test.go` - empty names, extremely large sizes, concurrent modifications
- [ ] T086 Validate Constitution Check in plan.md - verify all principles satisfied (TDD, public API, dependencies, versioning, observability, security)
- [ ] T087 Run final `make check` to verify 85%+ test coverage achieved
- [ ] T088 Update main README.md with Volumes API client section

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (Phase 4)**: Depends on Foundational (Phase 2) - No dependencies on other stories (can run parallel with US1)
- **User Story 3 (Phase 5)**: Depends on Foundational (Phase 2) - No dependencies on other stories (can run parallel with US1, US2)
- **User Story 4 (Phase 6)**: Depends on Foundational (Phase 2) - Ideally after US2 (to have volumes to test actions on)
- **User Story 5 (Phase 7)**: Depends on Foundational (Phase 2) and US2 (uses Create method) - No other dependencies
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start immediately after Foundational - List volume types is completely independent
- **User Story 2 (P1)**: Can start immediately after Foundational - Create/Update/Delete volumes is independent
- **User Story 3 (P2)**: Can start immediately after Foundational - List/Get volumes is independent (though testing is easier if US2 exists)
- **User Story 4 (P2)**: Can start after Foundational - Volume actions are independent (though testing needs volumes from US2)
- **User Story 5 (P3)**: Uses Create method from US2 - light dependency but Create already supports SnapshotID

### Within Each User Story

- Request/Response models BEFORE tests (TDD requires models to exist to write tests against)
- Tests MUST be written and FAIL before implementation
- Implementation AFTER tests are written
- Verify tests pass after implementation
- Run `make check` after story complete

### Parallel Opportunities

**Phase 1 (Setup)**: All 3 tasks can run in parallel

**Phase 2 (Foundational)**: 
- T004-T008: All model definitions can run in parallel (different files)
- T009-T011: All unit test writing can run in parallel (TDD - write first)
- T012-T014: All implementations can run in parallel (different files)

**Phase 3 (User Story 1)**: 
- T016-T019: All test writing can run in parallel
- T020-T022: Implementation tasks are sequential (client struct → method → integration)

**Phase 4 (User Story 2)**:
- T025-T026: Model definitions in parallel
- T027-T028: Unit test writing in parallel
- T030-T037: All test writing can run in parallel (8 test tasks)
- T038-T041: Implementation tasks are sequential

**Phase 5 (User Story 3)**:
- T047-T054: All test writing can run in parallel (8 test tasks)
- T055-T056: Implementation tasks can run in parallel (different methods)

**Phase 6 (User Story 4)**:
- T062-T068: All test writing can run in parallel (7 test tasks)

**Phase 7 (User Story 5)**:
- T072-T074: All test writing can run in parallel (3 test tasks)

**Phase 8 (Polish)**:
- T078-T082, T084-T085: All documentation and additional testing can run in parallel

**MAXIMUM PARALLELISM**: After Phase 2 completes, User Stories 1, 2, and 3 can ALL be worked on simultaneously by different team members.

---

## Parallel Example: User Story 2

If working with multiple developers, User Story 2 tasks could be distributed:

**Developer A** (Test Writing):
- T027, T028 (unit tests for request models)
- T030, T031, T032 (contract tests)

**Developer B** (Test Writing):
- T033, T034, T035, T036 (unit tests for client methods)
- T037 (integration test)

**Developer C** (Implementation):
- Wait for tests to be written (T027-T037)
- Then implement T038-T041 sequentially
- Run tests to verify

This demonstrates **test-first TDD** with parallel test writing followed by implementation.

---

## Implementation Strategy

**MVP Scope**: User Stories 1 & 2 (P1 priority)
- US1: List volume types (foundation knowledge)
- US2: Create, Update, Delete volumes (core lifecycle)

**Incremental Delivery**:
1. **Phase 2** (Foundational): Block ~1-2 days, team works together on shared models
2. **Phase 3** (US1): ~1 day, can be done by one developer
3. **Phase 4** (US2): ~2-3 days, highest value CRUD operations
4. **Checkpoint**: MVP complete, SDK is usable for basic volume management
5. **Phase 5** (US3): ~1-2 days, adds discovery and monitoring capabilities
6. **Phase 6** (US4): ~2-3 days, adds advanced operations (attach/detach/extend)
7. **Phase 7** (US5): ~1 day, adds snapshot restore capability
8. **Phase 8** (Polish): ~1-2 days, documentation and final cleanup

**Estimated Total Time**: 8-12 development days with 1-2 developers

---

## Validation Checklist

Before marking this feature complete, verify:

- [ ] All 88 tasks completed
- [ ] All contract tests pass (7 endpoints tested against Swagger schema)
- [ ] All unit tests pass with 85%+ code coverage
- [ ] All integration tests pass (full volume lifecycle validated)
- [ ] `make check` passes (linting, formatting, tests)
- [ ] Constitution Check verified (TDD, public API, dependencies, versioning, observability, security)
- [ ] VPS client has Volumes() and VolumeTypes() accessor methods
- [ ] All 5 user stories independently testable and functional
- [ ] Error messages are descriptive and actionable
- [ ] Context cancellation works properly (< 1 second)
- [ ] Documentation includes code examples for all public methods
- [ ] README.md updated with Volumes API client information
- [ ] No raw HTTP exposed in public APIs
- [ ] Bearer Token never logged
- [ ] Logger interface used for DEBUG level request/response logging

---

**Total Tasks**: 88
- Setup: 3 tasks
- Foundational: 12 tasks (blocks all user stories)
- User Story 1: 9 tasks
- User Story 2: 20 tasks
- User Story 3: 14 tasks
- User Story 4: 13 tasks
- User Story 5: 6 tasks
- Polish: 11 tasks

**Critical Path**: Setup → Foundational → User Story 2 (longest implementation path)

**Parallel Opportunities**: 40+ tasks can run in parallel across different phases

**MVP Scope**: Setup + Foundational + US1 + US2 = 44 tasks (~50% of total)
