---
description: "Implementation tasks for fixing server model and module"
---

# Tasks: Fix Server Model & Module

**Branch**: `005-fix-server-model`
**Input**: Design documents from `/workspaces/cloud-sdk/specs/005-fix-server-model/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are MANDATORY for public APIs per the Constitution. Write tests FIRST
and ensure they FAIL before implementation. Include unit tests and contract tests
derived from Swagger/OpenAPI for each endpoint wrapper.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3, US4)
- Include exact file paths in descriptions

## Path Conventions

- **Models**: `/workspaces/cloud-sdk/models/vps/servers/`
- **Modules**: `/workspaces/cloud-sdk/modules/vps/servers/`
- **Tests**: Test files alongside implementation (`*_test.go`)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Review existing server model structure in `/workspaces/cloud-sdk/models/vps/servers/`
- [X] T002 Review existing server module structure in `/workspaces/cloud-sdk/modules/vps/servers/`
- [X] T003 [P] Create backup of existing server.go and actions.go for reference

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core types and infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 [P] Define ServerStatus custom type with constants in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T005 [P] Define ServerAction custom type with constants in `/workspaces/cloud-sdk/models/vps/servers/actions.go`
- [X] T006 [P] Define RebootType custom type with constants in `/workspaces/cloud-sdk/models/vps/servers/actions.go`
- [X] T007 [P] Create IDName helper type in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T008 [P] Create FlavorInfo type in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T009 [P] Create VRMImgInfo type in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T010 Update Server struct with all pb.ServerInfo fields in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T011 Add JSON tags and validation rules to Server struct in `/workspaces/cloud-sdk/models/vps/servers/server.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - List Servers (Priority: P1) 🎯 MVP

**Goal**: Enable listing servers with all documented fields and query parameter filtering

**Independent Test**: Call `client.Servers().List(ctx, projectID, options)` with various filters and verify response contains all documented fields

### Tests for User Story 1 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [X] T012 [P] [US1] Write contract test for GET /servers in `/workspaces/cloud-sdk/modules/vps/servers/test/list_servers_test.go`
- [X] T013 [P] [US1] Write unit test for List method in `/workspaces/cloud-sdk/modules/vps/servers/client_test.go`
- [X] T014 [P] [US1] Write test for ServersListRequest query parameter handling in `/workspaces/cloud-sdk/modules/vps/servers/client_test.go`

### Implementation for User Story 1

- [X] T015 [P] [US1] Create ServersListRequest struct in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T016 [P] [US1] Create ServersListResponse struct in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T017 [US1] Implement List() method in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T018 [US1] Add query parameter encoding logic for filters in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T019 [US1] Verify all tests pass for User Story 1

**Checkpoint**: ✅ User Story 1 is fully functional and testable independently - All 41 tests passing

---

## Phase 4: User Story 2 - Server CRUD Operations (Priority: P2)

**Goal**: Enable create, get, update, and delete server operations with proper request/response types

**Independent Test**: Create server, get it, update its name/description, delete it - all operations should work with documented fields

### Tests for User Story 2 (MANDATORY - TDD) ⚠️

- [X] T020 [P] [US2] Write contract test for POST /servers in `/workspaces/cloud-sdk/modules/vps/servers/test/create_server_test.go`
- [X] T021 [P] [US2] Write contract test for GET /servers/{id} in `/workspaces/cloud-sdk/modules/vps/servers/test/get_server_test.go`
- [X] T022 [P] [US2] Write contract test for PUT /servers/{id} in `/workspaces/cloud-sdk/modules/vps/servers/test/update_server_test.go`
- [X] T023 [P] [US2] Write contract test for DELETE /servers/{id} in `/workspaces/cloud-sdk/modules/vps/servers/test/delete_server_test.go`
- [X] T024 [P] [US2] Write unit tests for Create/Get/Update/Delete methods in `/workspaces/cloud-sdk/modules/vps/servers/client_test.go`

### Implementation for User Story 2

- [X] T025 [P] [US2] Create ServerCreateRequest struct with validation in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T026 [P] [US2] Create ServerUpdateRequest struct in `/workspaces/cloud-sdk/models/vps/servers/server.go`
- [X] T027 [P] [US2] Create ServerDiskRequest struct in `/workspaces/cloud-sdk/models/vps/servers/volumes.go`
- [X] T028 [US2] Implement Create() method in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T029 [US2] Implement Get() method in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T030 [US2] Implement Update() method in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T031 [US2] Implement Delete() method in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T032 [US2] Add validation logic for ServerCreateRequest required fields in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T033 [US2] Verify all tests pass for User Story 2

**Checkpoint**: ✅ User Stories 1 AND 2 are both fully functional and independently testable

---

## Phase 5: User Story 3 - Server Actions (Priority: P3)

**Goal**: Enable server action operations (start, stop, reboot, resize, approve, reject, extend_root, get_pwd)

**Independent Test**: Perform each server action with appropriate parameters and verify action completes successfully

### Tests for User Story 3 (MANDATORY - TDD) ⚠️

- [X] T034 [P] [US3] Write contract test for POST /servers/{id}/action in `/workspaces/cloud-sdk/modules/vps/servers/test/server_actions_test.go`
- [X] T035 [P] [US3] Write unit tests for each action type in `/workspaces/cloud-sdk/modules/vps/servers/client_test.go`
- [X] T036 [P] [US3] Write test for action parameter validation in `/workspaces/cloud-sdk/modules/vps/servers/client_test.go`

### Implementation for User Story 3

- [X] T037 [P] [US3] Create ServerActionRequest struct in `/workspaces/cloud-sdk/models/vps/servers/actions.go`
- [X] T038 [P] [US3] Create ServerActionResponse struct in `/workspaces/cloud-sdk/models/vps/servers/actions.go`
- [X] T039 [US3] Implement Action() base method in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T040 [US3] Implement Stop() helper method via Action() in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T041 [US3] Implement Start() helper method via Action() in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T042 [US3] Implement Reboot() helper method with RebootType via Action() in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T043 [US3] Implement Resize() helper method via Action() in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T044 [US3] Implement Approve() helper method via Action() in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T045 [US3] Implement Reject() helper method via Action() in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T046 [US3] Implement ExtendRoot() helper method via Action() in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T047 [US3] Implement GetPassword() helper method via Action() in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T048 [US3] Add action parameter validation logic in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T049 [US3] Verify all tests pass for User Story 3

**Checkpoint**: ✅ User Stories 1, 2, AND 3 are all fully functional and independently testable

---

## Phase 6: User Story 4 - Metrics & NIC Management (Priority: P4)

**Goal**: Enable server metrics retrieval and NIC sub-resource management (List, Create, Update, Delete, Associate FloatingIP)

**Independent Test**: Get server metrics, list NICs, create NIC, update NIC security groups, associate floating IP, delete NIC - all should work independently

### Tests for User Story 4 (MANDATORY - TDD) ⚠️

- [X] T050 [P] [US4] Write contract test for GET /servers/{id}/metric in `/workspaces/cloud-sdk/modules/vps/servers/test/servers_metrics_test.go`
- [X] T051 [P] [US4] Write contract test for GET /servers/{id}/nics in `/workspaces/cloud-sdk/modules/vps/servers/test/server_nics_test.go`
- [X] T052 [P] [US4] Write contract test for POST /servers/{id}/nics in `/workspaces/cloud-sdk/modules/vps/servers/test/server_nics_test.go`
- [X] T053 [P] [US4] Write contract test for PUT /servers/{id}/nics/{nic-id} in `/workspaces/cloud-sdk/modules/vps/servers/test/server_nics_test.go`
- [X] T054 [P] [US4] Write contract test for DELETE /servers/{id}/nics/{nic-id} in `/workspaces/cloud-sdk/modules/vps/servers/test/server_nics_test.go`
- [X] T055 [P] [US4] Write contract test for POST /servers/{id}/nics/{nic-id}/floatingip in `/workspaces/cloud-sdk/modules/vps/servers/test/server_nics_test.go`
- [X] T056 [P] [US4] Write unit tests for metrics and NIC operations in `/workspaces/cloud-sdk/modules/vps/servers/client_test.go`

### Implementation for User Story 4

#### Metrics Sub-Tasks

- [X] T057 [P] [US4] Create MetricInfo struct in `/workspaces/cloud-sdk/models/vps/servers/metrics.go`
- [X] T058 [P] [US4] Create Measure struct in `/workspaces/cloud-sdk/models/vps/servers/metrics.go`
- [X] T059 [P] [US4] Create ServerMetricsRequest options in `/workspaces/cloud-sdk/models/vps/servers/metrics.go`
- [X] T060 [US4] Implement Metrics() method in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T061 [US4] Add query parameter encoding for metrics filters in `/workspaces/cloud-sdk/modules/vps/servers/client.go`

#### NIC Sub-Resource Sub-Tasks

- [X] T062 [P] [US4] Create ServerNIC struct in `/workspaces/cloud-sdk/models/vps/servers/nics.go`
- [X] T063 [P] [US4] Create ServerNICsListResponse struct in `/workspaces/cloud-sdk/models/vps/servers/nics.go`
- [X] T064 [P] [US4] Create ServerNICCreateRequest struct in `/workspaces/cloud-sdk/models/vps/servers/nics.go`
- [X] T065 [P] [US4] Create ServerNICUpdateRequest struct in `/workspaces/cloud-sdk/models/vps/servers/nics.go`
- [X] T066 [P] [US4] Create ServerNICAssociateFloatingIPRequest struct in `/workspaces/cloud-sdk/models/vps/servers/nics.go`
- [X] T067 [P] [US4] Create FloatingIPInfo struct in `/workspaces/cloud-sdk/models/vps/servers/nics.go`
- [X] T068 [US4] Create NICsClient sub-resource struct in `/workspaces/cloud-sdk/modules/vps/servers/nics.go`
- [X] T069 [US4] Implement NICs() sub-resource accessor in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [X] T070 [US4] Implement List() method for NICs in `/workspaces/cloud-sdk/modules/vps/servers/nics.go`
- [X] T071 [US4] Implement Add() method for NICs in `/workspaces/cloud-sdk/modules/vps/servers/nics.go`
- [X] T072 [US4] Implement Update() method for NICs in `/workspaces/cloud-sdk/modules/vps/servers/nics.go`
- [X] T073 [US4] Implement Delete() method for NICs in `/workspaces/cloud-sdk/modules/vps/servers/nics.go`
- [X] T074 [US4] Implement AssociateFloatingIP() method for NICs in `/workspaces/cloud-sdk/modules/vps/servers/nics.go`
- [X] T075 [US4] Add validation for NIC create/update requests in `/workspaces/cloud-sdk/modules/vps/servers/nics.go`
- [X] T076 [US4] Verify all tests pass for User Story 4

**Checkpoint**: ✅ All user stories (1-4) are now fully functional and independently testable

---

## Phase 7: Console & Volume Operations (Priority: P5)

**Goal**: Enable VNC console URL retrieval (console operations are documented in contracts)

**Independent Test**: Get VNC console URL for a server and verify URL is returned

### Tests for User Story 5 (MANDATORY - TDD) ⚠️

- [ ] T077 [P] [US5] Write contract test for GET /servers/{id}/vnc_url in `/workspaces/cloud-sdk/modules/vps/servers/test/console_test.go`
- [ ] T078 [P] [US5] Write unit test for GetVNCConsoleURL in `/workspaces/cloud-sdk/modules/vps/servers/client_test.go`

### Implementation for User Story 5

- [ ] T079 [P] [US5] Create ServerConsoleURLResponse struct in `/workspaces/cloud-sdk/models/vps/servers/console.go`
- [ ] T080 [US5] Implement GetVNCConsoleURL() method in `/workspaces/cloud-sdk/modules/vps/servers/client.go`
- [ ] T081 [US5] Verify all tests pass for User Story 5

**Checkpoint**: All core functionality complete

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T082 [P] Update module README.md with new API examples in `/workspaces/cloud-sdk/modules/vps/README.md`
- [ ] T083 [P] Update EXAMPLES.md with server operations in `/workspaces/cloud-sdk/modules/vps/EXAMPLES.md`
- [ ] T084 Run `make check` to verify 80%+ test coverage
- [ ] T085 Fix any linting issues reported by make check
- [ ] T086 [P] Add godoc comments to all public types and methods in models/vps/servers/
- [ ] T087 [P] Add godoc comments to all public methods in modules/vps/servers/
- [ ] T088 Validate Constitution Check items in plan.md (TDD, API shape, no deps, versioning, observability, security)
- [ ] T089 Update go.mod if any new internal dependencies added
- [ ] T090 Run integration tests against real API (if available) or mock server

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phases 3-7)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (US1 → US2 → US3 → US4 → US5)
- **Polish (Phase 8)**: Depends on all desired user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of US1
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Independent of US1/US2
- **User Story 4 (P4)**: Can start after Foundational (Phase 2) - Independent of US1/US2/US3
- **User Story 5 (P5)**: Can start after Foundational (Phase 2) - Independent of all others

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- Models before services
- Services before endpoints
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Models within a story marked [P] can run in parallel
- Different user stories can be worked on in parallel by different team members

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
T012: "Write contract test for GET /servers"
T013: "Write unit test for List method"
T014: "Write test for ServersListOptions query parameter handling"

# Launch all models for User Story 1 together:
T015: "Create ServersListRequest struct"
T016: "Create ServersListResponse struct"

# Then sequential implementation:
T017: "Implement List() method"
T018: "Add query parameter encoding logic"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup (T001-T003)
2. Complete Phase 2: Foundational (T004-T011) - CRITICAL - blocks all stories
3. Complete Phase 3: User Story 1 (T012-T019)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Run `make check` - should pass with basic coverage

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (List Servers) → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 (CRUD) → Test independently → Deploy/Demo
4. Add User Story 3 (Actions) → Test independently → Deploy/Demo
5. Add User Story 4 (Metrics & NICs) → Test independently → Deploy/Demo
6. Add User Story 5 (Console) → Test independently → Deploy/Demo
7. Complete Polish phase → Final release

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (List)
   - Developer B: User Story 2 (CRUD)
   - Developer C: User Story 3 (Actions)
   - Developer D: User Story 4 (Metrics/NICs)
   - Developer E: User Story 5 (Console)
3. Stories complete and integrate independently

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Constitution requires: TDD, clean API (no raw HTTP), no extra deps, typed errors, pluggable logging
- Run `make check` after each user story to ensure coverage stays above 80%
