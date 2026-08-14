# Tasks: Fix Floating IP Model (007-fix-floating-ip-model)

**Input**: Design documents from `/specs/007-fix-floating-ip-model/`  
**Status**: Generated November 11, 2025  
**Branch**: `007-fix-floating-ip-model`

**Prerequisites**: 
- ✅ plan.md (defines tech stack, structure, Constitution check)
- ✅ spec.md (3 user stories with priorities P1, P2, P3)
- ✅ data-model.md (19 FloatingIP fields, custom FloatingIPStatus enum)
- ✅ research.md (5 design decisions documented)
- ✅ quickstart.md (usage examples for all 6 operations)
- ✅ contracts/floatingip-openapi.yaml (6 endpoints from vps.yaml)

**Tests**: Mandatory for all public APIs per SDK Constitution. Tests MUST be written FIRST (TDD) using Go `testing` package.

**Organization**: Tasks grouped by user story to enable independent implementation and testing of each story. Each story is independently testable and can be delivered as an MVP increment.

---

## Format: `- [ ] [TaskID] [P?] [Story?] Description`

- **Checkbox**: Always `- [ ]` at start
- **Task ID**: Sequential (T001, T002, T003...) in execution order
- **[P]**: Parallelizable if present (different files, no dependencies on incomplete tasks)
- **[Story]**: User story label (US1, US2, US3) for story-specific tasks only
- **Description**: Clear action with exact file path

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure for Floating IP model and module

**Expected Duration**: < 1 hour

- [X] T001 Review existing floatingips model and module structure in `models/vps/floatingips/` and `modules/vps/floatingips/`
- [X] T002 [P] Verify Go 1.21+ is installed and `make check` infrastructure is available
- [X] T003 Ensure branch `007-fix-floating-ip-model` is active and tracking origin

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core data types and enums that MUST be defined before any API operations can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase completes. This phase defines the custom `FloatingIPStatus` enum and supporting `IDName` type used by all stories.

**Estimated Duration**: 1-2 hours

### Type Definitions (Models)

- [X] T004 [P] Add `FloatingIPStatus` custom enum type to `models/vps/floatingips/floatingip.go` with 4 constants (Active, Pending, Down, Rejected) and `String()` method
- [X] T005 [P] Add `Valid()` method to `FloatingIPStatus` to validate against the 4 allowed values
- [X] T006 [P] Add `IDName` struct type to `models/vps/floatingips/floatingip.go` with `ID` and `Name` fields for project/user references
- [X] T007 Add unit tests for `FloatingIPStatus` enum constants and `Valid()` method in `models/vps/floatingips/floatingip_test.go`
- [X] T008 [P] Add unit tests for `IDName` struct marshaling/unmarshaling in `models/vps/floatingips/floatingip_test.go`

**Checkpoint**: Custom types defined and tested. Ready for user story implementation.

---

## Phase 3: User Story 1 - Correct FloatingIP Model Structure (Priority: P1) 🎯 MVP

**Goal**: Define the complete FloatingIP model with all 19 fields from pb.FloatingIPInfo matching API specification, enabling developers to access all floating IP attributes without data loss.

**Independent Test Criteria**:
- FloatingIP struct successfully unmarshals from API JSON response with all fields
- All 19 fields from pb.FloatingIPInfo are present and correctly typed
- Marshaling FloatingIP back to JSON preserves all field values
- Optional fields (omitempty) are handled correctly when absent from API response
- Timestamp fields maintain RFC3339 format
- Status field uses custom FloatingIPStatus enum type

### Tests for User Story 1 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementing the model

- [X] T009 [P] [US1] Create unit test for FloatingIP struct JSON unmarshaling with all 19 fields in `models/vps/floatingips/floatingip_test.go`
- [X] T010 [P] [US1] Create unit test for FloatingIP struct JSON marshaling roundtrip (unmarshal → marshal → compare) in `models/vps/floatingips/floatingip_test.go`
- [X] T011 [P] [US1] Create unit test for optional fields (description, device_id, device_name, device_type, extnet_id, namespace, updatedAt, approvedAt, statusReason, project, user) being omitted from JSON when empty in `models/vps/floatingips/floatingip_test.go`
- [X] T012 [P] [US1] Create unit test for camelCase JSON field tags matching API specification (createdAt not created_at, etc.) in `models/vps/floatingips/floatingip_test.go`
- [X] T013 [P] [US1] Create unit test for IDName nested objects in project and user fields in `models/vps/floatingips/floatingip_test.go`

### Implementation for User Story 1

- [X] T014 [US1] Add complete FloatingIP struct definition with all 19 fields to `models/vps/floatingips/floatingip.go`: ID, UUID, Name, Address, Description, PortID, DeviceID, DeviceName, DeviceType, ProjectID, Project, Namespace, ExtNetID, UserID, User, Status (FloatingIPStatus), StatusReason, Reserved, CreatedAt, UpdatedAt, ApprovedAt with proper JSON tags matching API spec
- [X] T015 [US1] Verify all JSON field tags use camelCase for API field names (e.g., `json:"createdAt"` not `json:"created_at"`)
- [X] T016 [US1] Verify optional/nullable fields use `omitempty` tag where appropriate and pointer types (*IDName) for optional references
- [X] T017 [US1] Run unit tests (T009-T013) and verify all tests PASS
- [X] T018 [US1] Validate FloatingIP model with `make check` to ensure formatting and linting pass

**Checkpoint**: User Story 1 complete. FloatingIP model fully defined with all 19 fields, custom enum type for status, and comprehensive unit tests passing. Model ready for API operations.

---

## Phase 4: User Story 2 - Floating IP List Response Structure (Priority: P2)

**Goal**: Implement List operation returning `[]*FloatingIP` slice matching pb.FIPListOutput structure (breaking change from old "items" wrapper), enabling developers to directly iterate through floating IPs without custom unwrapping.

**Independent Test Criteria**:
- List() operation returns `[]*FloatingIP` slice (not a wrapper object)
- Response structure matches pb.FIPListOutput from API specification
- Empty list returns empty slice without errors
- List operation can be tested independently from other stories
- Response correctly marshals/unmarshals with API contract

### Tests for User Story 2 (MANDATORY - TDD) ⚠️

- [X] T019 [P] [US2] Create contract test for List operation against vps.yaml pb.FIPListOutput spec in `modules/vps/floatingips/test/contract_test.go` (test that response contains "floatingips" array, not "items")
- [X] T020 [P] [US2] Create unit test for List() method returning []*FloatingIP in `modules/vps/floatingips/client_test.go`
- [X] T021 [P] [US2] Create unit test for empty list response in `modules/vps/floatingips/test/contract_test.go`
- [X] T022 [P] [US2] Create integration test for List() with mock API server returning pb.FIPListOutput structure in `modules/vps/floatingips/test/integration_test.go`

### Implementation for User Story 2

- [X] T023 [US2] Create `FloatingIPListOptions` struct in `modules/vps/floatingips/client.go` with optional filter fields (status, user_id, device_type, device_id, extnet_id, address, name, detail) matching API query parameters
- [X] T024 [US2] Implement `List(ctx context.Context, opts *FloatingIPListOptions) ([]*FloatingIP, error)` method in `modules/vps/floatingips/client.go` returning slice directly (not wrapper object with "items" field)
- [X] T025 [US2] Implement proper error wrapping in List operation following SDK pattern: `fmt.Errorf("failed to list floatingips: %w", err)`
- [X] T026 [US2] Add query parameter construction for optional filters in List operation
- [X] T027 [US2] Run contract tests (T019, T021) and verify API response structure compliance
- [X] T028 [US2] Run integration tests (T022) against mock API server
- [X] T029 [US2] Update existing tests in `modules/vps/floatingips/client_test.go` if needed to reflect new list response structure
- [X] T030 [US2] Validate with `make check` to ensure all tests pass

**Checkpoint**: User Story 2 complete. List operation implemented with correct pb.FIPListOutput structure (floatingips array). Breaking change from "items" wrapper now live. Stories 1 and 2 both working independently.

---

## Phase 5: User Story 3 - Floating IP Create/Update Requests (Priority: P3)

**Goal**: Implement Create, Get, Update, Delete, and Disassociate operations with proper Request/Response structs matching FIPCreateInput/FIPUpdateInput specifications, enabling developers to create and manage floating IPs with name and description fields.

**Independent Test Criteria**:
- Create operation accepts name and description fields per FIPCreateInput
- Update operation accepts name and description fields per FIPUpdateInput
- Get operation returns complete FloatingIP with all fields
- Delete operation succeeds without errors on valid ID
- Disassociate operation removes device association
- All operations follow SDK error pattern with proper context wrapping
- All 6 operations can be tested independently from Stories 1 and 2
- Request/Response structures match Swagger specification exactly

### Tests for User Story 3 (MANDATORY - TDD) ⚠️

- [X] T031 [P] [US3] Create contract test for Create operation against vps.yaml FIPCreateInput spec in `modules/vps/floatingips/test/contract_test.go` (verify request body has name, description fields)
- [X] T032 [P] [US3] Create contract test for Get operation in `modules/vps/floatingips/test/contract_test.go`
- [X] T033 [P] [US3] Create contract test for Update operation against FIPUpdateInput spec in `modules/vps/floatingips/test/contract_test.go`
- [X] T034 [P] [US3] Create contract test for Delete operation in `modules/vps/floatingips/test/contract_test.go`
- [X] T035 [P] [US3] Create contract test for Disassociate operation in `modules/vps/floatingips/test/contract_test.go`
- [X] T036 [P] [US3] Create integration tests for complete CRUD workflow (create → get → update → delete) in `modules/vps/floatingips/test/integration_test.go`
- [X] T037 [P] [US3] Create integration test for Disassociate operation in `modules/vps/floatingips/test/integration_test.go`
- [ ] T038 [P] [US3] Create unit tests for request body construction in `modules/vps/floatingips/client_test.go`

### Implementation for User Story 3

#### Request/Response Types

- [X] T039 [US3] Create `FloatingIPCreateRequest` struct in `modules/vps/floatingips/client.go` with fields: Name (required), Description (optional) matching FIPCreateInput (exclude Reserved field)
- [X] T040 [US3] Create `FloatingIPUpdateRequest` struct in `modules/vps/floatingips/client.go` with fields: Name, Description, Reserved matching FIPUpdateInput (note: Reserved is included but documented as ignored by API)
- [X] T041 [US3] Create `FloatingIPResponse` or use existing response type for single resource responses matching pb.FloatingIPInfo

#### Create Operation

- [X] T042 [US3] Implement `Create(ctx context.Context, req *FloatingIPCreateRequest) (*FloatingIP, error)` method in `modules/vps/floatingips/client.go` with:
  - POST request to `/floatingips` endpoint
  - Request body marshaling
  - Response unmarshaling into FloatingIP struct
  - Error wrapping: `fmt.Errorf("failed to create floatingip: %w", err)`
  - Input validation (Name field required)

#### Get Operation

- [X] T043 [US3] Implement `Get(ctx context.Context, id string) (*FloatingIP, error)` method in `modules/vps/floatingips/client.go` with:
  - GET request to `/floatingips/{id}` endpoint
  - Response unmarshaling into FloatingIP struct
  - Error wrapping with operation context
  - Proper HTTP 404 handling for non-existent resources

#### Update Operation

- [X] T044 [US3] Implement `Update(ctx context.Context, id string, req *FloatingIPUpdateRequest) (*FloatingIP, error)` method in `modules/vps/floatingips/client.go` with:
  - PUT request to `/floatingips/{id}` endpoint
  - Request body marshaling
  - Response unmarshaling into FloatingIP struct
  - Error wrapping with operation context
  - Documentation that Reserved field is ignored by API

#### Delete Operation

- [X] T045 [US3] Implement `Delete(ctx context.Context, id string) error` method in `modules/vps/floatingips/client.go` with:
  - DELETE request to `/floatingips/{id}` endpoint
  - Error wrapping with operation context
  - Success verification (204 No Content or 200 OK handling)

#### Disassociate Operation

- [X] T046 [US3] Implement `Disassociate(ctx context.Context, id string) (*FloatingIP, error)` method in `modules/vps/floatingips/client.go` with:
  - POST request to `/floatingips/{id}/disassociate` endpoint
  - Response unmarshaling into FloatingIP struct showing updated device association fields
  - Error wrapping with operation context

#### Error Handling & Testing

- [X] T047 [US3] Implement error response parsing following SDK pattern for all 6 operations
- [X] T048 [US3] Run contract tests (T031-T035) and verify request/response structure compliance
- [X] T049 [US3] Run integration tests (T036-T037) against mock API server for complete workflows
- [X] T050 [US3] Run unit tests for request body construction (T038)
- [X] T051 [US3] Update existing `modules/vps/floatingips/client_test.go` tests to reflect new Request/Response structs and all 6 operations
- [X] T052 [US3] Validate all operations with `make check` to ensure tests pass and code quality standards met

**Checkpoint**: User Story 3 complete. All 6 operations (List, Create, Get, Update, Delete, Disassociate) implemented with proper Request/Response structs. All three user stories (1, 2, 3) working independently and together.

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Documentation, migration guide, final validation, and breaking change communication

**Expected Duration**: 1-2 hours

### Documentation & Migration

- [ ] T053 Create migration guide in `specs/007-fix-floating-ip-model/MIGRATION.md` documenting breaking change from "items" field to direct slice response, with before/after code examples
- [ ] T054 Update `modules/vps/floatingips/README.md` (or create if missing) with complete API documentation for all 6 operations
- [ ] T055 Update `modules/vps/floatingips/EXAMPLES.md` (or create if missing) with 6 operation examples matching `quickstart.md`
- [ ] T056 [P] Add inline code documentation (godoc comments) for all public types and methods in `models/vps/floatingips/floatingip.go`
- [ ] T057 [P] Add inline code documentation (godoc comments) for Client and all methods in `modules/vps/floatingips/client.go`

### Final Validation

- [ ] T058 Run full test suite: `make check` to verify all tests pass (unit, contract, integration, linting, formatting)
- [ ] T059 Verify Constitution compliance checklist from `plan.md` is satisfied (all 7 principles)
- [ ] T060 [P] Validate test coverage is adequate for all 6 operations and 19 model fields
- [ ] T061 [P] Review error messages for clarity and debugging value following SDK logging pattern
- [ ] T062 Verify no breaking changes to sibling modules (flavors, keypairs, etc.)

### Code Quality & Cleanup

- [ ] T063 [P] Run `go fmt` on all modified files to ensure formatting compliance
- [ ] T064 [P] Run `go vet` on all modified files to check for common Go errors
- [ ] T065 Remove any temporary debug code or commented-out code
- [ ] T066 Ensure all type definitions and enums follow Go naming conventions (PascalCase for public, camelCase for private)

### Release Preparation

- [ ] T067 Create release notes entry in `RELEASE_NOTES.md` or `CHANGELOG.md` documenting:
  - Feature: Fixed FloatingIP model to match API specification (007-fix-floating-ip-model)
  - Breaking Change: List response now returns direct slice instead of wrapper object with "items" field
  - Migration: Link to MIGRATION.md with before/after examples
  - Version: Bump to next MAJOR version per semantic versioning
  - All 19 fields now available including timestamps, project/user references, device info
  - Status field now uses custom FloatingIPStatus enum type (type-safe)
  
- [ ] T068 Update agent context in `.github/copilot-instructions.md` to reflect 007-fix-floating-ip-model completion

### Final Checkpoint

- [ ] T069 All 69 tasks complete and validated
- [ ] T070 `make check` passes with no errors or warnings
- [ ] T071 All test suites pass: unit, contract, integration, cross-cutting
- [ ] T072 Breaking change migration guide published and accessible
- [ ] T073 Branch `007-fix-floating-ip-model` ready for code review and merge

**Checkpoint**: Polish complete. All documentation updated, breaking change communicated, full test suite passing. Ready for merge and release.

---

## Dependencies & Execution Order

### Phase Dependencies

| Phase | Depends On | Status |
|-------|-----------|--------|
| Phase 1: Setup | None | Can start immediately |
| Phase 2: Foundational | Setup ✅ | BLOCKS all user stories - must complete first |
| Phase 3: User Story 1 (P1) | Setup ✅, Foundational ✅ | MVP focus - can start after Foundational |
| Phase 4: User Story 2 (P2) | Setup ✅, Foundational ✅ | Can start after Foundational (parallel with US1 if team capacity) |
| Phase 5: User Story 3 (P3) | Setup ✅, Foundational ✅ | Can start after Foundational (parallel with US1/US2 if team capacity) |
| Phase 6: Polish | All stories ✅ | Final phase - only after all features complete |

### Within Phase 2 (Foundational)

```
T004, T005, T006 (Define custom types) [P] --parallel-->
                                         |
                                         v
                                    T007, T008 (Unit tests for types) [P]
```

**Gate**: All foundational types tested before proceeding to user stories.

### Within Phase 3 (User Story 1)

```
T009-T013 (Write failing unit tests) [P] --parallel-->
                                         |
                                         v
                                    T014-T018 (Implement model)
                                         |
                                         v
                                    make check (Validation)
```

**Gate**: Tests must FAIL before implementation starts (TDD requirement).

### Within Phase 4 (User Story 2)

```
T019-T022 (Write failing contract/integration tests) [P] --parallel-->
                                         |
                                         v
                                    T023-T029 (Implement List operation)
                                         |
                                         v
                                    T030: make check (Validation)
```

### Within Phase 5 (User Story 3)

```
T031-T038 (Write failing contract tests) [P] --parallel-->
                                         |
                                         v
                      T039-T046 (Define requests, implement 6 operations) [P] --parallel-->
                                         |
                                         v
                      T047-T052 (Error handling, validation, make check)
```

### Parallel Opportunities

**Phase 2 (Foundational)**: Tasks T004, T005, T006 can run in parallel (different types, no dependencies)

**Phase 3 (User Story 1)**: 
- Tests T009-T013 can run in parallel (all unit tests for same struct)
- Implementation T014-T016 sequential (must complete model definition first)

**Phase 4 (User Story 2)**:
- Tests T019-T022 can run in parallel (different test files)
- Implementation T023-T029 sequential (building List operation step-by-step)

**Phase 5 (User Story 3)**:
- All contract tests T031-T035 can run in parallel (different endpoints)
- Request struct definitions T039-T041 can run in parallel
- Operation implementations T042-T046 can run in parallel (different endpoints)
- Final validation T047-T052 sequential (make check must be last)

**Cross-Story Parallelization**:
Once Phase 2 (Foundational) completes:
- Developer A can work on User Story 1 (Phase 3)
- Developer B can work on User Story 2 (Phase 4)  
- Developer C can work on User Story 3 (Phase 5)
- All stories proceed independently until Phase 6

### Explicit Sequential Dependencies

1. T001-T003 (Setup) → T004-T008 (Foundational) → [US1, US2, US3 start]
2. T009-T013 (US1 tests) → MUST FAIL before T014 (US1 implementation)
3. T019-T022 (US2 tests) → MUST FAIL before T023 (US2 implementation)
4. T031-T038 (US3 tests) → MUST FAIL before T039 (US3 implementation)
5. All stories → T058 (Full make check) → T067-T073 (Polish & Release)

---

## Implementation Strategy

### MVP First (Recommended Approach)

For single developer or tight timeline:

1. **Phase 1**: Complete Setup (30 min)
2. **Phase 2**: Complete Foundational (60-90 min) - CRITICAL PATH
3. **Phase 3**: Complete User Story 1 (2-3 hours) - MVP deliverable
4. **STOP & VALIDATE**: Test US1 independently (30 min)
5. **Deploy/Demo**: Show working FloatingIP model
6. **Phase 4**: User Story 2 (1-2 hours)
7. **Phase 5**: User Story 3 (2-3 hours)
8. **Phase 6**: Polish & Release (1-2 hours)

**Total MVP Timeline**: ~4-5 hours to Phase 3 completion + validation

### Incremental Delivery

Recommended for continuous integration/deployment:

```
Sprint 1:
├── Phase 1: Setup ✅
├── Phase 2: Foundational ✅
└── Phase 3: User Story 1 ✅ → DEMO/DEPLOY

Sprint 2:
├── Phase 4: User Story 2 ✅ → DEMO/DEPLOY

Sprint 3:
├── Phase 5: User Story 3 ✅ → DEMO/DEPLOY

Sprint 4:
├── Phase 6: Polish & Release ✅ → PRODUCTION RELEASE
```

Each sprint delivers working increment without breaking previous functionality.

### Parallel Team Strategy (3+ Developers)

**Week 1:**
- All developers: Phase 1 Setup (30 min)
- All developers: Phase 2 Foundational (90 min)

**Week 2:** (Parallel development)
- Developer A: Phase 3 User Story 1
- Developer B: Phase 4 User Story 2
- Developer C: Phase 5 User Story 3

**Week 3:**
- All developers: Phase 6 Polish & Release

Stories proceed independently with no merge conflicts (different files):
- `models/vps/floatingips/floatingip.go` → Main source of contention (use branch strategy)
- `modules/vps/floatingips/client.go` → Single file, but different methods (unlikely conflict)
- Test files separate by story: contract_test.go, integration_test.go, client_test.go

---

## Validation Checkpoints

### After Each Phase

| Phase | Checkpoint Command | Expected Result |
|-------|-------------------|-----------------|
| Phase 2 | `make check` | All foundational type tests PASS |
| Phase 3 | `make check` | US1 model tests PASS, structure correct |
| Phase 4 | `make check` | US2 list operation tests PASS, API contract verified |
| Phase 5 | `make check` | US3 all 6 operations tests PASS, contract compliance verified |
| Phase 6 | `make check` | All tests PASS, no lint/format errors, ready to merge |

### End-to-End Validation

```bash
# Run before Phase 6 Polish
make check

# Expected output:
# ✓ go fmt
# ✓ go vet
# ✓ golint
# ✓ unit tests: 45 tests PASS (example count)
# ✓ contract tests: 12 tests PASS
# ✓ integration tests: 8 tests PASS
# All checks passed!
```

---

## Notes & Guidance

### For Implementation

- **TDD Mandatory**: Write tests FIRST, ensure they FAIL before implementation. This validates test quality.
- **One Task at a Time**: Complete a task fully before moving to next. Commit after each logical group.
- **Error Wrapping**: Always wrap errors with context: `fmt.Errorf("failed to {operation}: %w", err)`
- **Breaking Change**: Phase 4 (US2) implements the breaking change. Document clearly for users.
- **Enum Usage**: Custom `FloatingIPStatus` enum in all code examples, not raw strings.
- **JSON Tags**: Match API specification exactly (createdAt, not created_at). Validate with contracts.

### For Testing

- **Contract Tests**: Verify request/response structure matches vps.yaml exactly
- **Integration Tests**: Use mock API server returning actual API response structures
- **Unit Tests**: Verify individual methods in isolation
- **Test Data**: Use realistic examples from vps.yaml pb.FloatingIPInfo

### For Code Quality

- **Godoc Comments**: All public types and methods must have documentation comments
- **Error Messages**: Include operation name and context (e.g., "failed to list floatingips")
- **Naming**: Follow Go conventions (PascalCase public, camelCase private)
- **No External Dependencies**: Ensure only Go stdlib used (encoding/json, context, net/http, fmt)

### Parallel Execution Example

**Scenario**: Team of 3 developers, working on Phases 3-5 in parallel

```bash
# Terminal 1 - Developer A: User Story 1
# After Foundational complete:
make test -run TestFloatingIP    # Run US1 tests
# Implement T014-T018
make check

# Terminal 2 - Developer B: User Story 2  
# Parallel with Developer A:
make test -run TestList          # Run US2 tests
# Implement T023-T029
make check

# Terminal 3 - Developer C: User Story 3
# Parallel with A and B:
make test -run TestCreate        # Run US3 tests
make test -run TestGet           # etc.
# Implement T039-T052
make check

# All developers merge when ready (different files)
# Then Phase 6: Polish together
```

---

## Task Summary

**Total Tasks**: 73 tasks organized in 6 phases

| Phase | Task Count | Duration | Priority | Status |
|-------|-----------|----------|----------|--------|
| Phase 1: Setup | 3 | <1 hour | P0 | Not Started |
| Phase 2: Foundational | 5 | 1-2 hours | P0 (BLOCKER) | Not Started |
| Phase 3: User Story 1 (P1) | 10 | 2-3 hours | P1 (MVP) | Not Started |
| Phase 4: User Story 2 (P2) | 12 | 1-2 hours | P2 | Not Started |
| Phase 5: User Story 3 (P3) | 30 | 3-4 hours | P3 | Not Started |
| Phase 6: Polish | 13 | 1-2 hours | P0 | Not Started |
| **TOTAL** | **73** | **8-14 hours** | Mixed | Not Started |

### Parallelizable Tasks

- **Foundational (Phase 2)**: 3 type definition tasks [P]
- **User Story 1 (Phase 3)**: 5 test tasks [P]
- **User Story 2 (Phase 4)**: 4 test tasks [P]
- **User Story 3 (Phase 5)**: 8 test tasks [P], 5 request struct definitions [P]
- **Polish (Phase 6)**: 4 code quality tasks [P]

**Parallel Reduction**: With proper parallelization, timeline can be reduced from 8-14 hours to 4-6 hours with 3 developers.

---

## References

- **Specification**: `/specs/007-fix-floating-ip-model/spec.md` (user stories, requirements)
- **Design**: `/specs/007-fix-floating-ip-model/plan.md` (tech stack, structure)
- **Data Model**: `/specs/007-fix-floating-ip-model/data-model.md` (19 fields, custom enum)
- **Research**: `/specs/007-fix-floating-ip-model/research.md` (5 design decisions)
- **Usage**: `/specs/007-fix-floating-ip-model/quickstart.md` (operation examples)
- **Contracts**: `/specs/007-fix-floating-ip-model/contracts/floatingip-openapi.yaml` (API spec)
- **Branch**: `007-fix-floating-ip-model` on GitHub
- **SDK Constitution**: `.github/copilot-instructions.md` (7 principles)

---

## Generated By

- **Generator**: Copilot /speckit.tasks
- **Date**: November 11, 2025
- **Input**: Phase 1 planning documents (plan.md, spec.md, data-model.md, research.md, quickstart.md, contracts/)
- **Output**: This tasks.md file with 73 tasks organized in 6 phases, ready for implementation
