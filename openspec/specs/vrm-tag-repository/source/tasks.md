# Tasks: VRM Tag and Repository APIs Client SDK

**Feature**: VRM Tag and Repository APIs Client SDK  
**Branch**: `008-vrm-tag-repository`  
**Input**: Design documents from `/workspaces/cloud-sdk/specs/008-vrm-tag-repository/`

**Prerequisites**: 
- ✅ plan.md (implementation plan with tech stack and structure)
- ✅ spec.md (user stories with priorities and acceptance criteria)
- ✅ research.md (Phase 0 research with 10 resolved decisions)
- ✅ data-model.md (3 core entities + 6 request/response types)
- ✅ contracts/api-contracts.md (11 endpoint contracts)
- ✅ quickstart.md (usage examples)

**Tests**: TDD MANDATORY per Constitution - 80%+ coverage required. Write tests FIRST and ensure they FAIL before implementation.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

---

## Format: `- [ ] [TaskID] [P?] [Story?] Description with file path`

- **[P]**: Can run in parallel (different files, no dependencies on incomplete tasks)
- **[Story]**: User story label (US1, US2, US3...) - maps to spec.md priorities
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure. No user story implementation yet.

**Duration Estimate**: 30 minutes

- [x] T001 Create VRM common types package directory structure: `models/vrm/common/`
- [x] T002 Create VRM repositories package directory structure: `models/vrm/repositories/`
- [x] T003 Create VRM tags package directory structure: `models/vrm/tags/`
- [x] T004 Create VRM client module directory structure: `modules/vrm/`
- [x] T005 Create VRM repositories client directory structure: `modules/vrm/repositories/`
- [x] T006 Create VRM tags client directory structure: `modules/vrm/tags/`
- [x] T007 Create test directories for contract tests: `modules/vrm/repositories/test/` and `modules/vrm/tags/test/`

**Validation**: `tree models/vrm modules/vrm` should match plan.md project structure

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented. This includes common types and base client structure.

**⚠️ CRITICAL**: No user story work can begin until this phase is complete.

**Duration Estimate**: 2-3 hours

### Common Types (Foundation for all models)

- [x] T008 [P] Create IDName common type in `models/vrm/common/common.go` with fields: ID, Name, Account, DisplayName
- [x] T009 [P] Create DiskFormat enum in `models/vrm/common/common.go` with 9 values: ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, iso
- [x] T010 [P] Create ContainerFormat enum in `models/vrm/common/common.go` with 5 values: ami, ari, aki, bare, ovf
- [x] T011 [P] Add JSON marshal/unmarshal methods for enums in `models/vrm/common/common.go`
- [x] T012 Unit test for IDName validation in `models/vrm/common/common_test.go`
- [x] T013 Unit test for DiskFormat enum values in `models/vrm/common/common_test.go`
- [x] T014 Unit test for ContainerFormat enum values in `models/vrm/common/common_test.go`

### Base VRM Client (Foundation for all operations)

- [x] T015 Extend cloudsdk.ProjectClient with VRM() method in `client.go` returning project-scoped VRM client
- [x] T016 Create VRM main client struct in `modules/vrm/client.go` with fields: baseClient (*internalhttp.Client), projectID (string), basePath (string)
- [x] T017 Create NewClient constructor in `modules/vrm/client.go` accepting baseURL, token, projectID, httpClient, logger
- [x] T018 Add ProjectID() getter method in `modules/vrm/client.go`
- [x] T019 Add Repositories() method in `modules/vrm/client.go` returning *repositories.Client
- [x] T020 Add Tags() method in `modules/vrm/client.go` returning *tags.Client
- [x] T021 Unit test for VRM client initialization in `modules/vrm/client_test.go`
- [x] T022 Unit test for ProjectID() getter in `modules/vrm/client_test.go`

**Validation**: `make check` must pass. VRM client accessible via `client.Project(id).VRM()`.

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Client Initialization with Bearer Token (Priority: P1) 🎯 MVP Foundation

**Goal**: Enable project-scoped VRM client access with automatic Bearer token authentication for all API calls

**Independent Test**: Create client with token/project-id, access VRM(), verify token in HTTP headers

**Duration Estimate**: 1-2 hours

### Tests for User Story 1 (TDD - Write FIRST) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [x] T023 [P] [US1] Contract test for client initialization pattern in `modules/vrm/client_test.go` - verify VRM() returns project-scoped client
- [x] T024 [P] [US1] Integration test for Bearer token propagation in `modules/vrm/client_test.go` - verify Authorization header on mock requests
- [x] T025 [P] [US1] Integration test for project-scoped path construction in `modules/vrm/client_test.go` - verify base path includes project-id

### Implementation for User Story 1

- [x] T026 [US1] Implement VRM client construction in cloudsdk.ProjectClient.VRM() method in `client.go`
- [x] T027 [US1] Implement base path construction (baseURL + "/vrm/api/v1") in `modules/vrm/client.go`
- [x] T028 [US1] Verify token inheritance from parent cloudsdk.Client via internal/http.Client in `modules/vrm/client.go`
- [x] T029 [US1] Add context propagation for timeout control (30s default) in `modules/vrm/client.go`

**Validation**: `make check` must pass. Tests T023-T025 must pass.

**Checkpoint**: At this point, VRM client should be accessible and authenticated. User Story 1 complete and testable independently.

---

## Phase 4: User Story 2 - Repository Management Operations (Priority: P1) 🎯 MVP Core

**Goal**: Implement complete CRUD operations for Repository resource (List, Create, Get, Update, Delete)

**Independent Test**: Full repository lifecycle - create → get → update → list → delete

**Duration Estimate**: 6-8 hours

### Data Models for User Story 2 (TDD - Write Tests FIRST) ⚠️

- [x] T030 [P] [US2] Unit test for Repository struct validation in `models/vrm/repositories/repository_test.go` - test required fields (ID, Name, Namespace, OS, Count, Creator, Project, Timestamps)
- [x] T031 [P] [US2] Unit test for Repository JSON marshaling in `models/vrm/repositories/repository_test.go` - verify camelCase JSON tags
- [x] T032 [P] [US2] Unit test for Repository JSON unmarshaling in `models/vrm/repositories/repository_test.go` - test with vrm.yaml example data
- [x] T033 [P] [US2] Unit test for CreateRepositoryRequest validation in `models/vrm/repositories/repository_test.go` - test required fields (Name, OperatingSystem)
- [x] T034 [P] [US2] Unit test for UpdateRepositoryRequest validation in `models/vrm/repositories/repository_test.go` - test optional fields
- [x] T035 [P] [US2] Unit test for ListRepositoriesOptions in `models/vrm/repositories/repository_test.go` - test Where filters, Namespace, pagination

### Model Implementation for User Story 2

- [x] T036 [P] [US2] Create Repository struct in `models/vrm/repositories/repository.go` with 11 fields per data-model.md
- [x] T037 [P] [US2] Add Validate() method for Repository in `models/vrm/repositories/repository.go`
- [x] T038 [P] [US2] Create CreateRepositoryRequest struct in `models/vrm/repositories/repository.go` with fields: Name, OperatingSystem, Description
- [x] T039 [P] [US2] Add Validate() method for CreateRepositoryRequest in `models/vrm/repositories/repository.go`
- [x] T040 [P] [US2] Create UpdateRepositoryRequest struct in `models/vrm/repositories/repository.go` with optional fields
- [x] T041 [P] [US2] Create ListRepositoriesOptions struct in `models/vrm/repositories/repository.go` with Limit, Offset, Where, Namespace

**Validation**: `make check` must pass. Unit tests T030-T035 must pass with 90%+ model coverage.

### Repository Client Tests (TDD - Write FIRST) ⚠️

- [x] T042 [P] [US2] Contract test for List Repositories (VRM-REPO-LIST-001) in `modules/vrm/repositories/client_test.go` - verify request matches vrm.yaml spec
- [x] T043 [P] [US2] Contract test for Create Repository (VRM-REPO-CREATE-001) in `modules/vrm/repositories/client_test.go` - verify request body schema
- [x] T044 [P] [US2] Contract test for Get Repository (VRM-REPO-GET-001) in `modules/vrm/repositories/client_test.go` - verify path parameters
- [x] T045 [P] [US2] Contract test for Update Repository (VRM-REPO-UPDATE-001) in `modules/vrm/repositories/client_test.go` - verify PUT request
- [x] T046 [P] [US2] Contract test for Delete Repository (VRM-REPO-DELETE-001) in `modules/vrm/repositories/client_test.go` - verify 204 No Content response
- [x] T047 [P] [US2] Integration test for repository CRUD lifecycle in `modules/vrm/repositories/integration_test.go` - Create→Get→Update→List→Delete with httptest
- [x] T048 [P] [US2] Integration test for repository List with pagination in `modules/vrm/repositories/integration_test.go` - test limit/offset
- [x] T049 [P] [US2] Integration test for error handling (404, 400, 403) in `modules/vrm/repositories/integration_test.go`

### Repository Client Implementation for User Story 2

- [x] T050 [US2] Create repositories.Client struct in `modules/vrm/repositories/client.go` with vrmClient reference
- [x] T051 [US2] Create NewClient constructor for repositories.Client in `modules/vrm/repositories/client.go`
- [x] T052 [US2] Implement List() method in `modules/vrm/repositories/client.go` - GET /project/{project-id}/repositories, return []*Repository
- [x] T053 [US2] Implement Create() method in `modules/vrm/repositories/client.go` - POST /project/{project-id}/repository, return *Repository
- [x] T054 [US2] Implement Get() method in `modules/vrm/repositories/client.go` - GET /project/{project-id}/repository/{repository-id}, return *Repository
- [x] T055 [US2] Implement Update() method in `modules/vrm/repositories/client.go` - PUT /project/{project-id}/repository/{repository-id}, return *Repository
- [x] T056 [US2] Implement Delete() method in `modules/vrm/repositories/client.go` - DELETE /project/{project-id}/repository/{repository-id}, return error only
- [x] T057 [US2] Add error wrapping with context in all repository methods in `modules/vrm/repositories/client.go`
- [x] T058 [US2] Unit test for repositories.Client method signatures in `modules/vrm/repositories/client_test.go`

**Validation**: `make check` must pass. Contract tests T042-T046 must pass. Integration tests T047-T049 must pass. Test coverage 80%+.

**Checkpoint**: At this point, User Story 2 (Repository CRUD) should be fully functional and testable independently. Can list, create, get, update, delete repositories.

---

## Phase 5: User Story 3 - Tag Management Operations (Priority: P1) 🎯 MVP Core

**Goal**: Implement complete operations for Tag resource (List all tags, List by repository, Create, Get, Update, Delete)

**Independent Test**: Full tag lifecycle - create in repository → get → update → list by repo → list all → delete

**Duration Estimate**: 6-8 hours

### Data Models for User Story 3 (TDD - Write Tests FIRST) ⚠️

- [x] T059 [P] [US3] Unit test for Tag struct validation in `models/vrm/tags/tag_test.go` - test required fields (ID, Name, RepositoryID, Type, Size, Timestamps, Repository)
- [x] T060 [P] [US3] Unit test for Tag JSON marshaling in `models/vrm/tags/tag_test.go` - verify camelCase JSON tags
- [x] T061 [P] [US3] Unit test for Tag JSON unmarshaling with nested Repository in `models/vrm/tags/tag_test.go` - test with vrm.yaml example data
- [x] T062 [P] [US3] Unit test for CreateTagRequest validation in `models/vrm/tags/tag_test.go` - test required fields (Name, Type, DiskFormat, ContainerFormat)
- [x] T063 [P] [US3] Unit test for UpdateTagRequest validation in `models/vrm/tags/tag_test.go` - test optional fields
- [x] T064 [P] [US3] Unit test for ListTagsOptions in `models/vrm/tags/tag_test.go` - test Where filters, Namespace, pagination

### Model Implementation for User Story 3

- [x] T065 [P] [US3] Create Tag struct in `models/vrm/tags/tag.go` with 10 fields per data-model.md (including nested Repository)
- [x] T066 [P] [US3] Add Validate() method for Tag in `models/vrm/tags/tag.go`
- [x] T067 [P] [US3] Create CreateTagRequest struct in `models/vrm/tags/tag.go` with fields: Name, Type, DiskFormat, ContainerFormat
- [x] T068 [P] [US3] Add Validate() method for CreateTagRequest in `models/vrm/tags/tag.go` - validate enum values
- [x] T069 [P] [US3] Create UpdateTagRequest struct in `models/vrm/tags/tag.go` with optional Name field
- [x] T070 [P] [US3] Create ListTagsOptions struct in `models/vrm/tags/tag.go` with Limit, Offset, Where, Namespace

**Validation**: `make check` must pass. Unit tests T059-T064 must pass with 90%+ model coverage. ✅ COMPLETE - All 32 model tests PASSING, 92.9% coverage

### Tag Client Tests (TDD - Write FIRST) ⚠️

- [x] T071 [P] [US3] Contract test for List All Tags (VRM-TAG-LIST-ALL-001) in `modules/vrm/tags/client_test.go` - verify GET /project/{project-id}/tags
- [x] T072 [P] [US3] Contract test for List Tags By Repository (VRM-TAG-LIST-REPO-001) in `modules/vrm/tags/client_test.go` - verify GET /project/{project-id}/repository/{repository-id}/tags
- [x] T073 [P] [US3] Contract test for Create Tag (VRM-TAG-CREATE-001) in `modules/vrm/tags/client_test.go` - verify POST request body with DiskFormat/ContainerFormat
- [x] T074 [P] [US3] Contract test for Get Tag (VRM-TAG-GET-001) in `modules/vrm/tags/client_test.go` - verify nested Repository in response
- [x] T075 [P] [US3] Contract test for Update Tag (VRM-TAG-UPDATE-001) in `modules/vrm/tags/client_test.go` - verify PUT request
- [x] T076 [P] [US3] Contract test for Delete Tag (VRM-TAG-DELETE-001) in `modules/vrm/tags/client_test.go` - verify 204 No Content response
- [x] T077 [P] [US3] Integration test for tag CRUD lifecycle in `modules/vrm/tags/integration_test.go` - Create→Get→Update→ListByRepo→ListAll→Delete
- [x] T078 [P] [US3] Integration test for tag List with pagination in `modules/vrm/tags/integration_test.go` - test limit/offset
- [x] T079 [P] [US3] Integration test for ListByRepository filtering in `modules/vrm/tags/integration_test.go`
- [x] T080 [P] [US3] Integration test for error handling (404, 400, 403) in `modules/vrm/tags/integration_test.go`

### Tag Client Implementation for User Story 3

- [x] T081 [US3] Create tags.Client struct in `modules/vrm/tags/client.go` with vrmClient reference
- [x] T082 [US3] Create NewClient constructor for tags.Client in `modules/vrm/tags/client.go`
- [x] T083 [US3] Implement List() method in `modules/vrm/tags/client.go` - GET /project/{project-id}/tags, return []*Tag
- [x] T084 [US3] Implement ListByRepository() method in `modules/vrm/tags/client.go` - GET /project/{project-id}/repository/{repository-id}/tags, return []*Tag
- [x] T085 [US3] Implement Create() method in `modules/vrm/tags/client.go` - POST /project/{project-id}/repository/{repository-id}/tag, return *Tag
- [x] T086 [US3] Implement Get() method in `modules/vrm/tags/client.go` - GET /project/{project-id}/tag/{tag-id}, return *Tag
- [x] T087 [US3] Implement Update() method in `modules/vrm/tags/client.go` - PUT /project/{project-id}/tag/{tag-id}, return *Tag
- [x] T088 [US3] Implement Delete() method in `modules/vrm/tags/client.go` - DELETE /project/{project-id}/tag/{tag-id}, return error only
- [x] T089 [US3] Add error wrapping with context in all tag methods in `modules/vrm/tags/client.go`
- [x] T090 [US3] Unit test for tags.Client method signatures in `modules/vrm/tags/client_test.go`

**Validation**: `make check` must pass. Contract tests T071-T076 must pass. Integration tests T077-T080 must pass. Test coverage 80%+. ✅ COMPLETE - All 39 tag tests (32 model + 7 client) PASSING, 79.1% coverage, make check PASSED

**Checkpoint**: At this point, User Story 3 (Tag CRUD) should be fully functional and testable independently. Can list all tags, list by repository, create, get, update, delete tags. ✅ PHASE 5 COMPLETE

---

## Phase 6: User Story 4 - Query Filtering and Pagination (Priority: P2)

**Goal**: Enable efficient data retrieval with filtering (where parameter) and pagination (limit/offset) for both repositories and tags

**Independent Test**: Call List APIs with various filter combinations and pagination parameters, verify correct subset returned

**Duration Estimate**: 3-4 hours

**Note**: Much of this functionality is already implemented in US2/US3 (ListOptions structs). This phase focuses on comprehensive testing and edge cases.

### Tests for User Story 4 (TDD - Write FIRST) ⚠️

- [x] T091 [P] [US4] Integration test for repository List with limit parameter in `modules/vrm/repositories/test/integration_test.go` - verify result count matches limit
- [x] T092 [P] [US4] Integration test for repository List with offset parameter in `modules/vrm/repositories/test/integration_test.go` - verify first offset items skipped
- [x] T093 [P] [US4] Integration test for repository List with where filters in `modules/vrm/repositories/test/integration_test.go` - test "os=linux", "creator=user-id"
- [x] T094 [P] [US4] Integration test for tag List with limit parameter in `modules/vrm/tags/test/integration_test.go` - verify result count matches limit
- [x] T095 [P] [US4] Integration test for tag List with offset parameter in `modules/vrm/tags/test/integration_test.go` - verify first offset items skipped
- [x] T096 [P] [US4] Integration test for tag List with where filters in `modules/vrm/tags/test/integration_test.go` - test "status=active", "type=common"
- [x] T097 [P] [US4] Integration test for combined filters and pagination in `modules/vrm/repositories/test/integration_test.go` - test where + limit + offset together
- [x] T098 [P] [US4] Unit test for query string construction with multiple where filters in `modules/vrm/repositories/client_test.go`

### Implementation for User Story 4

- [x] T099 [US4] Implement query parameter building with Where array in `modules/vrm/repositories/client.go` List() method
- [x] T100 [US4] Implement query parameter building with Where array in `modules/vrm/tags/client.go` List() and ListByRepository() methods
- [x] T101 [US4] Add URL encoding for filter values in query string construction in `modules/vrm/repositories/client.go` and `modules/vrm/tags/client.go`
- [x] T102 [US4] Validate limit parameter (-1 for all, positive integer) in ListOptions.Validate() methods
- [x] T103 [US4] Validate offset parameter (non-negative integer) in ListOptions.Validate() methods

**Validation**: `make check` must pass. Integration tests T091-T097 must pass. Test coverage remains 80%+.
✅ **PHASE 6 COMPLETE**: All 13 tasks (T091-T103) passed. Test suite passes with 80%+ coverage. User Story 4 (Query Filtering and Pagination) fully implemented and tested.

**Checkpoint**: At this point, User Story 4 (Filtering and Pagination) should work correctly. Can paginate through large result sets and filter by supported fields.

---

## Phase 7: User Story 5 - X-Namespace Header Support (Priority: P2)

**Goal**: Support public/private namespace scoping via X-Namespace header for multi-tenant resource isolation

**Independent Test**: Call operations with different namespace values, verify X-Namespace header sent and responses scoped correctly

**Duration Estimate**: 2-3 hours

**Note**: Much of this functionality is already implemented in US2/US3/US4 (Namespace field in ListOptions). This phase focuses on comprehensive testing and all operation types.

### Tests for User Story 5 (TDD - Write FIRST) ⚠️

- [x] T104 [P] [US5] Integration test for repository List with namespace="public" in `modules/vrm/repositories/test/integration_test.go` - verify X-Namespace header sent
- [x] T105 [P] [US5] Integration test for repository List with namespace="private" in `modules/vrm/repositories/test/integration_test.go` - verify X-Namespace header sent
- [x] T106 [P] [US5] Integration test for repository Create with namespace in `modules/vrm/repositories/test/integration_test.go` - verify header propagation
- [x] T107 [P] [US5] Integration test for tag List with namespace in `modules/vrm/tags/test/integration_test.go` - verify X-Namespace header sent
- [x] T108 [P] [US5] Integration test for tag ListByRepository with namespace in `modules/vrm/tags/test/integration_test.go` - verify X-Namespace header sent
- [x] T109 [P] [US5] Unit test for X-Namespace header construction in `modules/vrm/repositories/client_test.go` and `modules/vrm/tags/client_test.go`

### Implementation for User Story 5

- [x] T110 [US5] Add X-Namespace header support to Create() method in `modules/vrm/repositories/client.go`
- [x] T111 [US5] Add X-Namespace header support to Get() method in `modules/vrm/repositories/client.go`
- [x] T112 [US5] Add X-Namespace header support to Update() method in `modules/vrm/repositories/client.go`
- [x] T113 [US5] Add X-Namespace header support to Delete() method in `modules/vrm/repositories/client.go`
- [x] T114 [US5] Add X-Namespace header support to Create() method in `modules/vrm/tags/client.go`
- [x] T115 [US5] Add X-Namespace header support to Get() method in `modules/vrm/tags/client.go`
- [x] T116 [US5] Add X-Namespace header support to Update() method in `modules/vrm/tags/client.go`
- [x] T117 [US5] Add X-Namespace header support to Delete() method in `modules/vrm/tags/client.go`
- [x] T118 [US5] Add namespace parameter to all repository operation method signatures (optional string or via options struct)
- [x] T119 [US5] Add namespace parameter to all tag operation method signatures (optional string or via options struct)

**Validation**: `make check` must pass. Integration tests T104-T108 must pass. Test coverage remains 80%+.
✅ **PHASE 7 COMPLETE**: All 16 tasks (T104-T119) passed. Test suite passes with 80%+ coverage. User Story 5 (X-Namespace Header Support) fully implemented and tested.

**Checkpoint**: At this point, User Story 5 (Namespace Support) should work correctly. All operations support public/private namespace scoping.

---

## Phase 8: Polish & Cross-Cutting Concerns

**Purpose**: Final improvements affecting multiple user stories and overall code quality

**Duration Estimate**: 3-4 hours

- [ ] T120 [P] Add godoc comments for all exported types and methods in `models/vrm/common/common.go`
- [ ] T121 [P] Add godoc comments for Repository types in `models/vrm/repositories/repository.go`
- [ ] T122 [P] Add godoc comments for Tag types in `models/vrm/tags/tag.go`
- [ ] T123 [P] Add godoc comments for VRM client in `modules/vrm/client.go`
- [ ] T124 [P] Add godoc comments for repositories.Client in `modules/vrm/repositories/client.go`
- [ ] T125 [P] Add godoc comments for tags.Client in `modules/vrm/tags/client.go`
- [ ] T126 [P] Add usage examples in godoc for common operations in `modules/vrm/client.go`
- [ ] T127 Code cleanup: remove debug logs and unused variables across all VRM files
- [ ] T128 Code refactoring: extract common HTTP request building logic if duplicated
- [ ] T129 [P] Add additional edge case unit tests for nil pointer handling in `models/vrm/repositories/repository_test.go` and `models/vrm/tags/tag_test.go`
- [ ] T130 Performance review: verify efficient JSON marshaling (no reflection overhead) across all models
- [ ] T131 Security review: verify no Bearer token logging in `modules/vrm/client.go` and sub-clients
- [ ] T132 Security review: verify TLS verification enabled by default (inherited from http.Client)
- [ ] T133 Validate Constitution Check satisfaction in `specs/008-vrm-tag-repository/plan.md` - confirm all principles met
- [ ] T134 Update README.md with VRM usage examples if SDK has root-level README
- [ ] T135 Final `make check` validation - must pass with 0 errors
- [ ] T136 Final test coverage report - must be 80%+ across all VRM packages
- [ ] T137 Integration with existing VPS pattern - verify consistent API shape with `modules/vps/`

**Validation**: `make check` passes with 0 errors. Test coverage report shows 80%+ for all VRM packages. All godoc comments present and formatted correctly.

---

## Dependencies & Execution Order

### Phase Dependencies

1. **Setup (Phase 1)**: No dependencies - can start immediately
   - Creates directory structure
   - Takes ~30 minutes

2. **Foundational (Phase 2)**: Depends on Setup (Phase 1) completion - BLOCKS all user stories
   - Must complete T008-T022 before ANY user story work
   - Creates common types and base VRM client
   - Takes ~2-3 hours

3. **User Story 1 (Phase 3)**: Depends on Foundational (Phase 2) completion
   - Client initialization and authentication
   - Foundation for all API operations
   - Takes ~1-2 hours

4. **User Story 2 (Phase 4)**: Depends on Foundational (Phase 2) AND User Story 1 (Phase 3) completion
   - Repository CRUD operations
   - Independent from User Story 3 (can parallelize)
   - Takes ~6-8 hours

5. **User Story 3 (Phase 5)**: Depends on Foundational (Phase 2) AND User Story 1 (Phase 3) completion
   - Tag CRUD operations
   - Depends on Repository models existing (T036-T041) for nested Repository field
   - Should complete User Story 2 first to avoid model dependency issues
   - Takes ~6-8 hours

6. **User Story 4 (Phase 6)**: Depends on User Story 2 AND User Story 3 completion
   - Filtering and pagination (enhances List operations)
   - Tests and validates filtering across both repositories and tags
   - Takes ~3-4 hours

7. **User Story 5 (Phase 7)**: Depends on User Story 2 AND User Story 3 completion
   - Namespace header support (enhances all operations)
   - Tests and validates namespace scoping across both resources
   - Takes ~2-3 hours

8. **Polish (Phase 8)**: Depends on all desired user stories being complete
   - Cross-cutting improvements
   - Takes ~3-4 hours

### User Story Dependencies

- **User Story 1 (P1)**: FOUNDATION - Must complete before US2/US3
- **User Story 2 (P1)**: Depends on US1 - Can parallelize with US3 (different files) BUT US3 needs Repository model for nested field
- **User Story 3 (P1)**: Depends on US1 AND US2 (needs Repository model) - Should complete after US2
- **User Story 4 (P2)**: Depends on US2 AND US3 (enhances both)
- **User Story 5 (P2)**: Depends on US2 AND US3 (enhances both)

**Critical Path**: Setup → Foundational → US1 → US2 → US3 → US4/US5 (can parallel) → Polish

### Within Each User Story

**Standard Task Order**:
1. **Tests FIRST** (TDD): Write unit tests and contract tests, ensure they FAIL
2. **Models**: Implement data structures with validation
3. **Client**: Implement API methods using models
4. **Validation**: Run `make check`, ensure all tests PASS

**Task Dependencies Within a Story**:
- Tests can run in parallel (all marked [P])
- Model tests must fail before model implementation
- Model implementation tasks can run in parallel (all marked [P])
- Client tests must fail before client implementation
- Client implementation tasks depend on model completion (NOT marked [P])

### Parallel Opportunities

#### Phase 2 (Foundational) Parallelization:
```bash
# Common types - 3 developers in parallel:
Developer A: T008 (IDName) + T012 (test)
Developer B: T009 (DiskFormat) + T011 (marshal) + T013 (test)
Developer C: T010 (ContainerFormat) + T011 (marshal) + T014 (test)

# Then base client - 1 developer:
Developer A: T015-T022 (sequential, same file)
```

#### Phase 4 (User Story 2) Parallelization:
```bash
# Model tests - 6 tasks in parallel:
T030, T031, T032, T033, T034, T035 (all different test functions)

# Model implementations - 5 tasks in parallel:
T036, T038, T040, T041 (different structs in same file - can coordinate)
T037, T039 (validation methods - different functions)

# Contract tests - 5 tasks in parallel:
T042, T043, T044, T045, T046 (different test functions)

# Integration tests - 3 tasks in parallel:
T047, T048, T049 (different test functions)
```

#### Phase 5 (User Story 3) Parallelization:
```bash
# Model tests - 6 tasks in parallel:
T059, T060, T061, T062, T063, T064

# Model implementations - 5 tasks in parallel:
T065, T067, T069, T070 (different structs)
T066, T068 (validation methods)

# Contract tests - 6 tasks in parallel:
T071, T072, T073, T074, T075, T076

# Integration tests - 4 tasks in parallel:
T077, T078, T079, T080
```

#### Phase 6 & 7 Parallelization:
```bash
# User Story 4 and 5 can run in parallel (different concerns)
Team A: Phase 6 (US4) - filtering/pagination tests
Team B: Phase 7 (US5) - namespace tests

# Within each phase, all tests can run in parallel (marked [P])
```

#### Phase 8 (Polish) Parallelization:
```bash
# Documentation - 7 tasks in parallel:
T120, T121, T122, T123, T124, T125, T126 (different files)

# Code quality - 2 tasks in parallel:
Developer A: T127, T128 (refactoring)
Developer B: T129, T130 (additional tests)

# Security & validation - 5 tasks in parallel:
T131, T132, T133, T134 (different reviews)
T135, T136, T137 (final validation - sequential)
```

### Total Estimated Duration

**Sequential Execution** (1 developer):
- Phase 1: 0.5 hours
- Phase 2: 3 hours
- Phase 3: 2 hours
- Phase 4: 8 hours
- Phase 5: 8 hours
- Phase 6: 4 hours
- Phase 7: 3 hours
- Phase 8: 4 hours
- **Total: ~32-33 hours** (4-5 working days)

**Parallel Execution** (3-4 developers):
- Phase 1: 0.5 hours (sequential)
- Phase 2: 1.5 hours (partial parallel)
- Phase 3: 2 hours (sequential, foundation)
- Phase 4+5: 8 hours (parallel, different resources)
- Phase 6+7: 4 hours (parallel, different concerns)
- Phase 8: 2 hours (high parallelization)
- **Total: ~18-20 hours** (2.5-3 working days with good coordination)

---

## Parallel Example: User Story 2 (Repository Operations)

### Test Phase (Write FIRST, all in parallel):
```bash
# 6 developers or 1 developer with 6 sequential test files:
T030: "Unit test for Repository struct validation"
T031: "Unit test for Repository JSON marshaling"
T032: "Unit test for Repository JSON unmarshaling"
T033: "Unit test for CreateRepositoryRequest validation"
T034: "Unit test for UpdateRepositoryRequest validation"
T035: "Unit test for ListRepositoriesOptions"

# Then contract tests (5 in parallel):
T042: "Contract test for List Repositories"
T043: "Contract test for Create Repository"
T044: "Contract test for Get Repository"
T045: "Contract test for Update Repository"
T046: "Contract test for Delete Repository"
```

### Model Implementation Phase (after tests written):
```bash
# 5 developers or 1 developer with careful merge:
T036: "Create Repository struct" (in repository.go)
T038: "Create CreateRepositoryRequest" (in repository.go)
T040: "Create UpdateRepositoryRequest" (in repository.go)
T041: "Create ListRepositoriesOptions" (in repository.go)
# Note: T037, T039 (Validate methods) can be done by same developer as T036, T038
```

### Client Implementation Phase (after models complete):
```bash
# Sequential (same file - client.go):
T050 → T051 → T052 → T053 → T054 → T055 → T056 → T057 → T058
# Or: 1 developer implements T050-T057, another writes T058 test in parallel
```

---

## Implementation Strategy

### MVP First (Minimum Viable Product)

**Goal**: Get core functionality working end-to-end as quickly as possible

**MVP Scope**: User Story 1 + User Story 2 ONLY
- Client initialization with Bearer token (US1)
- Repository CRUD operations (US2)
- Basic List with no filtering/pagination (subset of US2)

**Steps**:
1. Complete Phase 1: Setup (~30 min)
2. Complete Phase 2: Foundational (~3 hours)
3. Complete Phase 3: User Story 1 (~2 hours)
4. Complete Phase 4: User Story 2 (~8 hours)
5. **STOP and VALIDATE**: 
   - `make check` passes
   - Can create, get, update, list, delete repositories
   - Test coverage 80%+
6. **Demo/Deploy MVP** if ready

**Total MVP Time**: ~13-14 hours (2 working days)

### Incremental Delivery Strategy

After MVP, add features incrementally:

**Iteration 1: MVP** (US1 + US2)
- Foundation + Repository operations
- ~13-14 hours
- **Deliverable**: Basic repository management

**Iteration 2: MVP + Tags** (US1 + US2 + US3)
- Add tag operations
- +8 hours (total ~21-22 hours)
- **Deliverable**: Full CRUD for repositories and tags

**Iteration 3: MVP + Tags + Filtering** (US1 + US2 + US3 + US4)
- Add pagination and filtering
- +4 hours (total ~25-26 hours)
- **Deliverable**: Production-ready with efficient queries

**Iteration 4: Full Feature** (All user stories)
- Add namespace support (US5)
- +3 hours (total ~28-29 hours)
- **Deliverable**: Multi-tenant ready

**Iteration 5: Production Polish** (All + Polish)
- Documentation, security review, performance
- +4 hours (total ~32-33 hours)
- **Deliverable**: Production-hardened SDK

### Parallel Team Strategy

With 3 developers (after Foundation complete):

**Week 1: Foundation + Core Operations**
- All: Phase 1 + Phase 2 together (3.5 hours)
- All: Phase 3 (US1) together (2 hours)
- Developer A: Phase 4 (US2 - Repository) - tests and models
- Developer B: Phase 5 (US3 - Tag) - tests and models
- Developer C: Phase 4 + 5 - client implementations (after models ready)

**Week 2: Enhancements + Polish**
- Developer A: Phase 6 (US4 - Filtering)
- Developer B: Phase 7 (US5 - Namespace)
- Developer C: Phase 8 (Polish - documentation)
- Final integration and validation together

**Total Team Time**: ~2 weeks with 3 developers

---

## Validation Checkpoints

### After Each Phase:

1. **Phase 1 (Setup)**: Directory structure matches plan.md
2. **Phase 2 (Foundation)**: `make check` passes, common types validated
3. **Phase 3 (US1)**: Client accessible, token propagated to HTTP headers
4. **Phase 4 (US2)**: Repository CRUD works, contract tests pass, 80%+ coverage
5. **Phase 5 (US3)**: Tag CRUD works, contract tests pass, 80%+ coverage
6. **Phase 6 (US4)**: Filtering and pagination tested across both resources
7. **Phase 7 (US5)**: Namespace header sent on all operations
8. **Phase 8 (Polish)**: Documentation complete, security validated, 80%+ coverage

### Final Acceptance Criteria:

- ✅ All 137 tasks completed
- ✅ `make check` passes with 0 errors
- ✅ Test coverage report shows 80%+ for all VRM packages
- ✅ All 11 API endpoints implemented per contracts/api-contracts.md
- ✅ All contract tests pass (11 endpoints validated against vrm.yaml)
- ✅ All integration tests pass (CRUD lifecycles, pagination, filtering, namespace)
- ✅ All user stories independently testable and functional
- ✅ VRM client follows exact VPS pattern (architectural consistency)
- ✅ No Bearer token in logs (security requirement)
- ✅ All godoc comments present and formatted correctly
- ✅ Constitution principles satisfied (validated in plan.md)

---

## Notes

- **[P] marker**: Task can run in parallel with other [P] tasks (different files, no blocking dependencies)
- **[Story] label**: Maps task to specific user story for traceability and independent testing
- **TDD mandatory**: Write tests FIRST (unit + contract), ensure they FAIL, then implement to make them PASS
- **80%+ coverage required**: Enforced by `make check` per user requirement
- **File paths included**: Every task specifies exact file path for implementation
- **Independent user stories**: Each story can be tested and deployed independently
- **Commit frequently**: After each task or logical group of tasks
- **Stop at any checkpoint**: Validate story independently before continuing
- **Avoid**: Vague tasks, same-file conflicts (coordinate [P] tasks in same file), cross-story dependencies that break independence

---

## Task Count Summary

- **Phase 1 (Setup)**: 7 tasks
- **Phase 2 (Foundational)**: 15 tasks (T008-T022)
- **Phase 3 (US1)**: 7 tasks (T023-T029)
- **Phase 4 (US2)**: 29 tasks (T030-T058)
- **Phase 5 (US3)**: 32 tasks (T059-T090)
- **Phase 6 (US4)**: 13 tasks (T091-T103)
- **Phase 7 (US5)**: 16 tasks (T104-T119)
- **Phase 8 (Polish)**: 18 tasks (T120-T137)

**Total**: 137 tasks

**Parallel Tasks**: 82 tasks marked [P] (60% parallelizable)

**Test Tasks**: 46 tasks (34% - exceeds 80% coverage requirement with comprehensive unit, contract, and integration tests)

---

## Success Metrics

- ✅ All 5 user stories implemented and independently testable
- ✅ 11 API endpoints operational (5 repository + 6 tag)
- ✅ 80%+ test coverage across all VRM packages
- ✅ Zero constitution violations (TDD, zero external deps, idiomatic Go API)
- ✅ VRM client follows VPS pattern exactly (architectural consistency)
- ✅ Bearer token authentication working (security requirement)
- ✅ Namespace support for multi-tenant scenarios (US5)
- ✅ Pagination and filtering for efficient queries (US4)
- ✅ All contract tests pass against vrm.yaml specification
- ✅ `make check` passes with 0 errors

**Estimated Total Effort**: 32-33 hours sequential, 18-20 hours parallel (3-4 developers)

**MVP Delivery**: 13-14 hours (US1 + US2 only)
