# Tasks: Fix Network Model Definition

**Input**: Design documents from `/specs/002-fix-network-model/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

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

**Purpose**: Confirm source of truth before making changes

- [X] T001 Review pb.NetworkInfo, NetCreateInput, and NetPort definitions in `swagger/vps.json`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared fixtures and helpers required by all user stories

- [X] T002 Create reusable full-network fixture JSON at `modules/vps/networks/test/testdata/network_full.json`
- [X] T003 Add fixture loader helper in `modules/vps/networks/test/helpers.go` for expanded network payloads

**Checkpoint**: Fixture assets in place for story-specific tests

---

## Phase 3: User Story 1 - SDK Users Access Complete Network Information (Priority: P1) 🎯 MVP

**Goal**: Expose every pb.NetworkInfo field (including nested router/project/user data) via SDK retrieval

**Independent Test**: Retrieve a network through `vpsClient.Networks().Get` and assert all pb.NetworkInfo fields and port details are populated

### Tests for User Story 1 (MANDATORY - TDD) ⚠️

- [X] T004 [P] [US1] Add JSON round-trip tests covering all pb.NetworkInfo fields in `models/vps/networks/network_test.go`
- [X] T005 [P] [US1] Expand retrieval scenarios in `modules/vps/networks/client_test.go` to assert gateway, nameservers, router, project, and user fields
- [X] T006 [P] [US1] Update port client expectations in `modules/vps/networks/ports_test.go` for NetPort addresses and server summary
- [X] T007 [P] [US1] Update HTTP fixture tests in `modules/vps/networks/test/networks_get_test.go` and `modules/vps/networks/test/networks_list_test.go` to validate new response fields
- [X] T008 [P] [US1] Extend port sub-resource scenarios in `modules/vps/networks/test/network_ports_test.go` with address list and server assertions
- [X] T009 [P] [US1] Enrich lifecycle coverage in `modules/vps/networks/test/networks_integration_test.go` to exercise full network payloads

### Implementation for User Story 1

- [X] T010 [US1] Implement complete pb.NetworkInfo field set and nested types in `models/vps/networks/network.go`
- [X] T011 [US1] Implement NetPort struct with server summary in `models/vps/networks/ports.go`

**Checkpoint**: Network retrieval exposes all Swagger-defined fields with passing tests

---

## Phase 4: User Story 2 - SDK Users Create Networks with Full Configuration Options (Priority: P2)

**Goal**: Allow optional gateway and router association when creating networks

**Independent Test**: Create a network through SDK with custom gateway and router ID; verify response echoes configuration and defaults work when omitted

### Tests for User Story 2 (MANDATORY - TDD) ⚠️

- [X] T012 [P] [US2] Add NetworkCreateRequest serialization tests for gateway and router_id in `models/vps/networks/network_test.go`
- [X] T013 [P] [US2] Update create flow expectations in `modules/vps/networks/client_test.go` to assert request payload includes optional fields
- [X] T014 [P] [US2] Update create contract tests in `modules/vps/networks/test/networks_create_test.go` for full configuration and defaults

### Implementation for User Story 2

- [X] T015 [US2] Add optional gateway and router_id fields with correct tags to `models/vps/networks/network.go`

**Checkpoint**: Network creation supports all Swagger-defined options with passing tests

---

## Phase 5: User Story 3 - SDK Maintains Type Safety and Consistency (Priority: P3)

**Goal**: Preserve backwards compatibility and ensure optional fields remain safe to omit

**Independent Test**: Run existing update/delete flows and ensure omitting new optional fields produces identical behavior as before

### Tests for User Story 3 (MANDATORY - TDD) ⚠️

- [X] T016 [P] [US3] Add regression tests ensuring Network unmarshals correctly when optional fields are absent in `models/vps/networks/network_test.go`
- [X] T017 [P] [US3] Update stability checks in `modules/vps/networks/test/networks_update_test.go` and `modules/vps/networks/test/networks_delete_test.go` to confirm unchanged behavior

### Implementation for User Story 3

- [X] T018 [US3] Audit JSON tags and add deprecation notes/omitempty handling in `models/vps/networks/network.go` to maintain type safety

**Checkpoint**: Existing SDK consumers remain unaffected; optional fields handled safely

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final validation, formatting, and manual verification

- [ ] T019 [P] Run `gofmt` on updated Go files in `models/vps/networks/` and `modules/vps/networks/`
- [ ] T020 Execute `make check` from repository root and ensure no errors
- [ ] T021 Perform six network module operations via `curl` against live base URL (provided at test time) and compare responses with SDK outputs

---

## Dependencies & Execution Order

1. Phase 1 → Phase 2 → User Story phases → Phase 6
2. User stories proceed in priority order: **US1 (P1)** → **US2 (P2)** → **US3 (P3)**
3. Within each user story: Tests (TDD) must complete before implementation tasks
4. Phase 6 tasks run after all desired user stories finish

## Parallel Opportunities

- Tasks marked [P] operate on distinct files and can run concurrently once their phase prerequisites are met
- After Phase 2, US1/US2/US3 testing tasks can be split among teammates while respecting story priority for delivery
- Polish task T019 can run in parallel with T021 after all code updates but before `make check`

## Implementation Strategy

### MVP (User Story 1 Only)
1. Complete Phases 1–2
2. Execute US1 test tasks (T004–T009) then implementation (T010–T011)
3. Run T019–T020 and validate retrieval fields before proceeding

### Incremental Delivery
1. Deliver US1 (network retrieval)
2. Deliver US2 (creation enhancements)
3. Deliver US3 (type-safety regressions)
4. Finish with Phase 6 polish and manual `curl` comparison (T021)

### Notes
- Maintain TDD discipline: ensure each new test fails before implementation
- Keep changes scoped to `models/vps/networks` and related tests per plan constraints
- Coordinate manual `curl` verification once base URL and token are provided
