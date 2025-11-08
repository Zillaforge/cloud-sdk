# Tasks: Fix Flavors Model

**Input**: Design documents from `/specs/001-fix-flavors-model/`
**Prerequisites**: plan.md (required), spec.md (required for user stories), research.md, data-model.md, contracts/

**Tests**: Tests are MANDATORY for public APIs per the Constitution. Write tests FIRST
and ensure they FAIL before implementation. Include unit tests and contract tests
derived from Swagger/OpenAPI for each endpoint wrapper.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

- **Go SDK**: `models/` (data structures), `modules/` (client implementations), `internal/` (shared utilities)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [x] T001 Create project structure per implementation plan
- [x] T002 Initialize Go module dependencies (standard library only)
- [x] T003 [P] Configure Go linting and formatting tools

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [x] T004 Setup internal utilities for JSON handling and time parsing
- [x] T005 Create base error types for API responses
- [x] T006 Setup HTTP client infrastructure with context support

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Correct Flavor Model Structure (Priority: P1) 🎯 MVP

**Goal**: Update the Flavor struct to accurately reflect the API specification with all required fields, correct JSON tags, and proper Go types.

**Independent Test**: Can be tested by verifying that Flavor structs can be unmarshaled from API responses and marshaled back without data loss.

### Tests for User Story 1 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [x] T007 [P] [US1] Unit tests for Flavor struct JSON marshaling/unmarshaling in models/vps/flavors/flavor_test.go
- [x] T008 [P] [US1] Contract tests for flavor API responses in modules/vps/flavors/client_test.go

### Implementation for User Story 1

- [x] T009 [US1] Update Flavor struct with all pb.FlavorInfo fields in models/vps/flavors/flavor.go
- [x] T010 [US1] Add proper JSON tags matching API specification in models/vps/flavors/flavor.go
- [x] T011 [US1] Implement field name changes (VCPUs→VCPU, RAM→Memory) with migration notes
- [x] T012 [US1] Add validation for required fields and data types

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - GPU Support in Flavors (Priority: P2)

**Goal**: Add GPU information support to the Flavor model for GPU-enabled flavors.

**Independent Test**: Can be tested by creating Flavor instances with GPU info and verifying serialization.

### Tests for User Story 2 (MANDATORY - TDD) ⚠️

- [x] T013 [P] [US2] Unit tests for GPUInfo struct in models/vps/flavors/flavor_test.go
- [x] T014 [P] [US2] Contract tests for GPU flavor responses in modules/vps/flavors/client_test.go

### Implementation for User Story 2

- [x] T015 [US2] Create GPUInfo struct with count, model, and VGPU fields in models/vps/flavors/flavor.go
- [x] T016 [US2] Add GPU field to Flavor struct as optional pointer type
- [x] T017 [US2] Implement JSON marshaling for GPU configuration
- [x] T018 [US2] Add GPU field validation rules

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Flavor Timestamps (Priority: P3)

**Goal**: Add creation, update, and deletion timestamps to the Flavor model for lifecycle tracking.

**Independent Test**: Can be tested by verifying timestamp fields are properly handled in JSON operations.

### Tests for User Story 3 (MANDATORY - TDD) ⚠️

- [x] T019 [P] [US3] Unit tests for timestamp field handling in models/vps/flavors/flavor_test.go
- [x] T020 [P] [US3] Contract tests for timestamp serialization in modules/vps/flavors/client_test.go

### Implementation for User Story 3

- [x] T021 [US3] Add timestamp fields (createdAt, updatedAt, deletedAt) as *time.Time in models/vps/flavors/flavor.go
- [x] T022 [US3] Implement custom JSON marshaling for ISO 8601 timestamps
- [x] T023 [US3] Add timestamp validation and parsing logic
- [x] T024 [US3] Update Flavor struct documentation with timestamp usage

**Checkpoint**: All user stories should now be independently functional

---

## Phase 6: Client Implementation & Filtering

**Purpose**: Implement the client methods and filtering capabilities for the flavors API

- [x] T025 [P] Update ListFlavorsOptions struct with API query parameters in models/vps/flavors/flavor.go
- [x] T026 [P] Implement List method with query parameter building in modules/vps/flavors/client.go
- [x] T027 [P] Implement Get method for individual flavor retrieval in modules/vps/flavors/client.go
- [x] T028 [P] Add URL encoding support for multiple tag filtering
- [x] T029 [P] Add context.Context support to all public methods
- [x] T030 [P] Implement proper error wrapping and typed responses

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T031 [P] Update documentation and examples in README.md
- [ ] T032 Code cleanup and refactoring across all files
- [ ] T033 [P] Additional unit tests for edge cases in models/vps/flavors/flavor_test.go
- [ ] T034 [P] Integration tests for client methods in modules/vps/flavors/client_test.go
- [ ] T035 Validate Constitution Check requirements (TDD, error handling, context usage)
- [ ] T036 [P] Create migration guide for breaking changes in docs/migration.md
- [ ] T037 Run make check to validate all changes pass tests and linting

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order (P1 → P2 → P3)
- **Client Implementation (Phase 6)**: Depends on all model stories (3-5) being complete
- **Polish (Phase 7)**: Depends on all implementation phases being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of other stories
- **User Story 3 (P3)**: Can start after Foundational (Phase 2) - Independent of other stories

### Within Each User Story

- Tests (MANDATORY) MUST be written and FAIL before implementation
- Model updates before client implementation
- Core implementation before integration
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- Once Foundational phase completes, all user stories can start in parallel (if team capacity allows)
- All tests for a user story marked [P] can run in parallel
- Client implementation tasks marked [P] can run in parallel
- Polish tasks marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Unit tests for Flavor struct JSON marshaling/unmarshaling in models/vps/flavors/flavor_test.go"
Task: "Contract tests for flavor API responses in modules/vps/flavors/client_test.go"

# Launch implementation tasks sequentially (dependencies):
Task: "Update Flavor struct with all pb.FlavorInfo fields in models/vps/flavors/flavor.go"
Task: "Add proper JSON tags matching API specification in models/vps/flavors/flavor.go"
Task: "Implement field name changes (VCPUs→VCPU, RAM→Memory) with migration notes"
Task: "Add validation for required fields and data types"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP!)
3. Add User Story 2 → Test independently → Deploy/Demo
4. Add User Story 3 → Test independently → Deploy/Demo
5. Add Client Implementation → Test filtering and API calls
6. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers:

1. Team completes Setup + Foundational together
2. Once Foundational is done:
   - Developer A: User Story 1 (core model)
   - Developer B: User Story 2 (GPU support)
   - Developer C: User Story 3 (timestamps)
3. Stories complete and integrate independently
4. One developer: Client implementation and filtering
5. Team: Polish and cross-cutting concerns

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence</content>
<parameter name="filePath">/workspaces/cloud-sdk/specs/001-fix-flavors-model/tasks.md