# Tasks: IAM API Client

**Input**: Design documents from `/specs/009-iam-api/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are MANDATORY for public APIs per the Constitution. Write tests FIRST
and ensure they FAIL before implementation. All tests are placed in resource directories
alongside implementation code (consistent with VPS/VRM pattern).

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and directory structure

- [X] T001 Create directory structure: `models/iam/{common,users,projects}` and `modules/iam/{users,projects}`
- [X] T002 [P] Create package declarations for all IAM modules
- [X] T003 [P] Create `modules/iam/README.md` with API overview and usage
- [X] T004 [P] Create `modules/iam/EXAMPLES.md` placeholder

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared models and base client that ALL user stories depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

### Shared Models (Required by all stories)

- [X] T005 [P] Create `models/iam/common/common.go` with Permission and TenantRole types
- [X] T006 [P] Write tests for Permission in `models/iam/common/common_test.go` (JSON parsing, validation)
- [X] T007 [P] Write tests for TenantRole in `models/iam/common/common_test.go` (enum values, validation)

### Base IAM Client Infrastructure

- [X] T008 Create `modules/iam/client.go` with base Client struct (receives *internalhttp.Client)
- [X] T009 Implement `NewClient(baseClient *internalhttp.Client)` constructor in `modules/iam/client.go`
- [X] T010 Implement `Users()` factory method in `modules/iam/client.go` (returns *users.Client)
- [X] T011 Implement `Projects()` factory method in `modules/iam/client.go` (returns *projects.Client)
- [X] T012 Write tests for Client initialization in `modules/iam/client_test.go`
- [X] T013 Write tests for factory methods in `modules/iam/client_test.go`

### Root SDK Integration

- [X] T014 Add `IAM()` method to root Client in `client.go` (constructs IAM client with internalhttp.Client)
- [X] T015 Write tests for `IAM()` method in `client_test.go`

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Retrieve Current User Information (Priority: P1) 🎯 MVP

**Goal**: Enable developers to retrieve authenticated user information via `client.IAM().User().Get(ctx)`

**Independent Test**: Initialize IAM client with valid Bearer token, call `User().Get(ctx)`, verify returned User struct contains account, userId, displayName, namespace, timestamps

**API Contract**: GET /user → Returns GetUserResponse with all user fields

### Tests for User Story 1 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [X] T016 [P] [US1] Write test for User model JSON parsing in `models/iam/users/models_test.go` (valid response)
- [X] T017 [P] [US1] Write test for User model with unknown fields in `models/iam/users/models_test.go` (forward compatibility)
- [X] T018 [P] [US1] Write test for GetUserResponse parsing in `models/iam/users/models_test.go`
- [X] T019 [P] [US1] Write test for User client Get() success case in `modules/iam/users/client_test.go` (mock HTTP 200)
- [X] T020 [P] [US1] Write test for User client Get() with 403 error in `modules/iam/users/client_test.go` (invalid token)
- [X] T021 [P] [US1] Write test for User client Get() with context timeout in `modules/iam/users/client_test.go`

### Implementation for User Story 1

- [X] T022 [P] [US1] Create User struct in `models/iam/users/models.go` (all fields from Swagger: account, userId, displayName, description, extra, namespace, email, frozen, mfa, timestamps)
- [X] T023 [P] [US1] Create GetUserResponse struct in `models/iam/users/models.go`
- [X] T024 [US1] Create `modules/iam/users/client.go` with Client struct (has baseClient *internalhttp.Client)
- [X] T025 [US1] Implement `NewClient(baseClient *internalhttp.Client)` in `modules/iam/users/client.go`
- [X] T026 [US1] Implement `Get(ctx context.Context)` method in `modules/iam/users/client.go` (GET /user, returns *User, error)
- [X] T027 [US1] Add error wrapping with context in Get() method (fmt.Errorf with %w)
- [X] T028 [US1] Run tests in `models/iam/users/models_test.go` - verify all pass
- [X] T029 [US1] Run tests in `modules/iam/users/client_test.go` - verify all pass

**Checkpoint**: User Story 1 complete - `client.IAM().User().Get(ctx)` fully functional and tested

---

## Phase 4: User Story 2 - List User's Projects (Priority: P2)

**Goal**: Enable developers to retrieve all projects the user belongs to via `client.IAM().Projects().List(ctx, opts)`

**Independent Test**: Call `Projects().List(ctx, nil)` for default pagination, verify returned list contains ProjectMembership items with project details, permissions, tenantRole. Test with custom ListProjectsOptions for pagination.

**API Contract**: GET /projects?offset=X&limit=Y&order=Z → Returns ListProjectsResponse with projects array and total count

### Tests for User Story 2 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [X] T030 [P] [US2] Write test for Project model JSON parsing in `models/iam/projects/models_test.go`
- [X] T031 [P] [US2] Write test for ProjectMembership model JSON parsing in `models/iam/projects/models_test.go`
- [X] T032 [P] [US2] Write test for ListProjectsResponse parsing in `models/iam/projects/models_test.go` (array extraction, total count)
- [X] T033 [P] [US2] Write test for ListProjectsOptions struct in `models/iam/projects/models_test.go` (nil handling, field pointers)
- [X] T034 [P] [US2] Write test for Projects client List() success with nil options in `modules/iam/projects/client_test.go`
- [X] T035 [P] [US2] Write test for Projects client List() with pagination options in `modules/iam/projects/client_test.go`
- [X] T036 [P] [US2] Write test for Projects client List() with empty results in `modules/iam/projects/client_test.go`
- [X] T037 [P] [US2] Write test for Projects client List() query parameter encoding in `modules/iam/projects/client_test.go`

### Implementation for User Story 2

- [X] T038 [P] [US2] Create Project struct in `models/iam/projects/models.go` (projectId, displayName, description, extra, namespace, frozen, timestamps)
- [X] T039 [P] [US2] Create ProjectMembership struct in `models/iam/projects/models.go` (wraps Project, adds globalPermissionId, userPermissionId, permissions, tenantRole, frozen)
- [X] T040 [P] [US2] Create ListProjectsResponse struct in `models/iam/projects/models.go` (projects array, total count)
- [X] T041 [P] [US2] Create ListProjectsOptions struct in `models/iam/projects/options.go` (Offset, Limit, Order as pointers)
- [X] T042 [US2] Create `modules/iam/projects/client.go` with Client struct (has baseClient *internalhttp.Client)
- [X] T043 [US2] Implement `NewClient(baseClient *internalhttp.Client)` in `modules/iam/projects/client.go`
- [X] T044 [US2] Implement `List(ctx context.Context, opts *ListProjectsOptions)` in `modules/iam/projects/client.go`
- [X] T045 [US2] Add query parameter building logic in List() (handle nil opts, encode offset/limit/order)
- [X] T046 [US2] Add error wrapping with context in List() method
- [X] T047 [US2] Run tests in `models/iam/projects/models_test.go` - verify all pass
- [X] T048 [US2] Run tests in `modules/iam/projects/client_test.go` for List() - verify all pass

**Checkpoint**: User Story 2 complete - `client.IAM().Projects().List(ctx, opts)` fully functional and tested

---

## Phase 5: User Story 3 - Retrieve Specific Project Details (Priority: P2)

**Goal**: Enable developers to fetch specific project details via `client.IAM().Projects().Get(ctx, projectID)`

**Independent Test**: Call `Projects().Get(ctx, validProjectID)`, verify returned project details. Test with invalid projectID for 400 error, unauthorized projectID for 403/404 error.

**API Contract**: GET /project/{project-id} → Returns GetProjectResponse with project details and permissions

### Tests for User Story 3 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [X] T049 [P] [US3] Write test for GetProjectResponse parsing in `models/iam/projects/models_test.go`
- [X] T050 [P] [US3] Write test for Projects client Get() success case in `modules/iam/projects/client_test.go` (mock HTTP 200)
- [X] T051 [P] [US3] Write test for Projects client Get() with invalid projectID in `modules/iam/projects/client_test.go` (400 error)
- [X] T052 [P] [US3] Write test for Projects client Get() with unauthorized projectID in `modules/iam/projects/client_test.go` (403/404 error)
- [X] T053 [P] [US3] Write test for Projects client Get() with context timeout in `modules/iam/projects/client_test.go`

### Implementation for User Story 3

- [X] T054 [P] [US3] Create GetProjectResponse struct in `models/iam/projects/models.go` (ProjectID, DisplayName, Description, Extra, GlobalPermission, UserPermission, Namespace, timestamps)
- [X] T055 [US3] Implement `Get(ctx context.Context, projectID string)` in `modules/iam/projects/client.go` (GET /project/{id})
- [X] T056 [US3] Add URL path construction with projectID in Get() method
- [X] T057 [US3] Add error wrapping with context in Get() method
- [X] T058 [US3] Run tests in `models/iam/projects/models_test.go` for GetProjectResponse - verify all pass
- [X] T059 [US3] Run tests in `modules/iam/projects/client_test.go` for Get() - verify all pass

**Checkpoint**: User Story 3 complete - `client.IAM().Projects().Get(ctx, projectID)` fully functional and tested

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, examples, validation, and final quality checks

- [ ] T060 [P] Update `modules/iam/README.md` with complete API documentation (all methods, parameters, return types)
- [ ] T061 [P] Complete `modules/iam/EXAMPLES.md` with usage examples for all three operations (GetUser, ListProjects, GetProject)
- [ ] T062 [P] Add quickstart example in `modules/iam/EXAMPLES.md` (5 lines of code)
- [ ] T063 [P] Add pagination example in `modules/iam/EXAMPLES.md`
- [ ] T064 [P] Add error handling example in `modules/iam/EXAMPLES.md`
- [ ] T065 Run `make check` - verify all tests pass and coverage >= 80%
- [ ] T066 Manual testing against real IAM service (all three endpoints)
- [ ] T067 Verify unknown fields are ignored (forward compatibility test)
- [ ] T068 Verify error messages are clear and include context
- [ ] T069 Validate Constitution Check compliance in plan.md
- [ ] T070 [P] Code review and refactoring for consistency with VPS/VRM patterns

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup (Phase 1) - BLOCKS all user stories
- **User Stories (Phases 3-5)**: All depend on Foundational (Phase 2) completion
  - US1, US2, US3 can then proceed in parallel (if team capacity allows)
  - Or sequentially in priority order: US1 (P1) → US2 (P2) → US3 (P2)
- **Polish (Phase 6)**: Depends on all user stories (Phases 3-5) being complete

### User Story Dependencies

- **User Story 1 (P1)**: Can start after Foundational (Phase 2) - No dependencies on other stories
- **User Story 2 (P2)**: Can start after Foundational (Phase 2) - Independent of US1 (different resource)
- **User Story 3 (P2)**: Can start after Foundational (Phase 2) - Builds on US2 (same Projects client) but should be independently testable

### Within Each User Story (TDD Flow)

1. **Tests FIRST**: All test tasks (T016-T021 for US1, T030-T037 for US2, T049-T053 for US3) - MUST be written and FAIL before implementation
2. **Models**: Create data structures (T022-T023 for US1, T038-T041 for US2, T054 for US3)
3. **Client**: Create resource client and methods (T024-T027 for US1, T042-T046 for US2, T055-T057 for US3)
4. **Validation**: Run tests to verify implementation (T028-T029 for US1, T047-T048 for US2, T058-T059 for US3)

### Parallel Opportunities

#### Setup Phase (Phase 1)
- T002, T003, T004 can all run in parallel (different files)

#### Foundational Phase (Phase 2)
- T005, T006, T007 (common models) can run in parallel
- T008-T013 (base client) are sequential (same file)
- T014-T015 (root integration) can run in parallel with T008-T013 or after

#### User Story 1 (Phase 3)
- All tests (T016-T021) can be written in parallel
- Models (T022-T023) can be written in parallel
- Client implementation (T024-T027) is mostly sequential (same file)

#### User Story 2 (Phase 4)
- All tests (T030-T037) can be written in parallel
- Models (T038-T041) can be written in parallel
- Client implementation (T042-T046) is mostly sequential

#### User Story 3 (Phase 5)
- All tests (T049-T053) can be written in parallel
- Model (T054) single task
- Client implementation (T055-T057) is sequential

#### Across User Stories
- **US1, US2, US3 can all be implemented in parallel** after Foundational phase completes (different resource directories)
- US2 and US3 share the same Projects client, so T042-T048 must complete before T055-T059 starts

#### Polish Phase (Phase 6)
- T060-T064 (documentation tasks) can all run in parallel
- T065-T070 (validation) should run after documentation

---

## Parallel Execution Examples

### Example 1: Single Developer (Sequential by Priority)

Week 1:
- Complete Phase 1 (Setup) - 1 day
- Complete Phase 2 (Foundational) - 2 days
- Complete Phase 3 (US1 - P1) - 2 days

Week 2:
- Complete Phase 4 (US2 - P2) - 2 days
- Complete Phase 5 (US3 - P2) - 2 days
- Complete Phase 6 (Polish) - 1 day

**Total**: ~10 days for MVP + all features

### Example 2: Two Developers (Parallel User Stories)

Week 1:
- Both: Complete Phase 1 + Phase 2 (Setup + Foundational) - 3 days
- Dev 1: Start Phase 3 (US1 - P1) - 2 days
- Dev 2: Start Phase 4 (US2 - P2) in parallel - 2 days

Week 2:
- Dev 1: Complete Phase 5 (US3 - P2) - 2 days
- Dev 2: Start Phase 6 (Polish) - 1 day
- Both: Final validation and documentation - 2 days

**Total**: ~8 days with parallel execution

### Example 3: MVP-First Approach (US1 Only)

Week 1:
- Complete Phase 1 (Setup) - 1 day
- Complete Phase 2 (Foundational) - 2 days
- Complete Phase 3 (US1 - P1 MVP) - 2 days
- Minimal Phase 6 (Polish for US1 only) - 0.5 days

**MVP Delivery**: ~5-6 days for User().Get() functionality
**Then**: Iterate to add US2 and US3 in subsequent releases

---

## Implementation Strategy

### Recommended Approach: MVP-First (US1), Then Incremental

1. **Phase 1-2** (3 days): Setup + Foundational infrastructure
2. **Phase 3** (2 days): US1 (P1) - Get User - **MVP RELEASE**
3. **Phase 4** (2 days): US2 (P2) - List Projects - **Release v0.2**
4. **Phase 5** (2 days): US3 (P2) - Get Project - **Release v0.3**
5. **Phase 6** (1 day): Polish - **Final v1.0**

### Why This Strategy?

- **Fast MVP**: User().Get() is the most fundamental operation (P1 priority)
- **Independent Stories**: Each phase delivers a complete, testable feature
- **Risk Mitigation**: Early feedback on API design from US1 before building US2/US3
- **Flexible Scope**: Can stop after MVP if business priorities change
- **Clear Milestones**: Each user story completion is a release candidate

---

## Task Summary

**Total Tasks**: 70
- **Phase 1 (Setup)**: 4 tasks
- **Phase 2 (Foundational)**: 11 tasks (BLOCKING)
- **Phase 3 (US1 - P1 MVP)**: 14 tasks (6 tests + 8 implementation)
- **Phase 4 (US2 - P2)**: 19 tasks (9 tests + 10 implementation)
- **Phase 5 (US3 - P2)**: 11 tasks (5 tests + 6 implementation)
- **Phase 6 (Polish)**: 11 tasks

**Test Tasks**: 20 (29% of total - comprehensive TDD coverage)
**Parallel Tasks**: 28 marked with [P] (40% can run in parallel)

**MVP Scope** (Suggested): Phase 1 + Phase 2 + Phase 3 = **29 tasks** for User().Get() functionality

**Estimated Effort**:
- MVP (US1 only): 5-6 days (single developer)
- All Features (US1+US2+US3): 10 days (single developer) or 8 days (two developers in parallel)

---

## Validation Checklist

Before considering implementation complete:

- [ ] All 70 tasks completed and checked off
- [ ] All tests pass (`make check`)
- [ ] Code coverage >= 80%
- [ ] All three endpoints work with real IAM service
- [ ] Error handling tested (invalid token, invalid projectID, timeouts)
- [ ] Pagination tested (nil opts, custom opts, empty results)
- [ ] Unknown fields ignored (forward compatibility)
- [ ] Documentation complete (README, EXAMPLES)
- [ ] Constitution Check validated
- [ ] Code consistent with VPS/VRM patterns (internalhttp.Client, single options pointer, tests in resource dirs)
