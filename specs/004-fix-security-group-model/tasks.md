# Implementation Tasks: Fix Security Group Model

**Feature**: Fix Security Group Model (001-fix-security-group-model)  
**Branch**: `001-fix-security-group-model`  
**Date**: November 9, 2025  
**Related**: [spec.md](./spec.md) | [plan.md](./plan.md) | [data-model.md](./data-model.md)

## Overview

This document breaks down the implementation of security group model corrections and CRUD operations into concrete, executable tasks. Tasks are organized by user story to enable independent implementation and testing per the TDD workflow.

**Total Estimated Tasks**: 81  
**Parallel Opportunities**: 15 tasks can be executed in parallel  
**MVP Scope**: User Story 1 + User Story 4 (Core model corrections + CRUD operations)

---

## Implementation Strategy

### Delivery Approach

1. **MVP First**: Complete User Story 1 (model corrections) + User Story 4 (CRUD operations) for immediate value
2. **Incremental Stories**: Each user story is independently testable and deliverable
3. **Parallel Execution**: Multiple tasks within each story can be done simultaneously (marked with [P])
4. **TDD Workflow**: Tests written before implementation for all stories (80%+ coverage target)

### Story Dependencies

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational)
    ↓
[US1] Model Corrections ────────┐
    ↓                           ↓
[US4] CRUD Operations ──────→ [US6] Rule Management
    ↓
[US5] List with Filters
    ↓
[US2] Model Field Accuracy ──→ [US3] JSON Serialization
    ↓
Polish & Documentation
```

**Key Insight**: US1 is foundational for all others. US4 and US6 can be developed in parallel after US1. US2 and US3 are validation/polish stories.

---

## Phase 1: Setup & Prerequisites

**Objective**: Prepare development environment and validate existing structure.

**Duration Estimate**: 30 minutes

**Independent Test**: Can verify by running `make check` and confirming test infrastructure works.

### Tasks

- [X] T001 Verify Go 1.21+ installation and module dependencies via `go version && go mod verify`
- [X] T002 Confirm existing test infrastructure works via `go test ./models/vps/securitygroups/... -v`
- [X] T003 Review swagger/vps.yaml security group definitions (pb.SgInfo, pb.SgListOutput, pb.SgRuleInfo)
- [X] T004 Create feature branch 001-fix-security-group-model if not exists via `git checkout -b 001-fix-security-group-model`
- [X] T005 Validate make check baseline via `make check` (expect existing issues, document them)

---

## Phase 2: Foundational Tasks

**Objective**: Set up shared infrastructure needed by all user stories.

**Duration Estimate**: 1 hour

**Independent Test**: Unit tests for model types compile and pass.

### Tasks

- [X] T006 [P] Create Protocol custom type with constants (ProtocolTCP="tcp", ProtocolUDP="udp", ProtocolICMP="icmp", ProtocolAny="any") as `type Protocol string` in models/vps/securitygroups/rule.go
- [X] T007 [P] Create Direction custom type with constants (DirectionIngress="ingress", DirectionEgress="egress") as `type Direction string` in models/vps/securitygroups/rule.go
- [X] T008 [P] Write unit tests for Protocol type marshaling/unmarshaling (verify JSON serialization to/from string values) in models/vps/securitygroups/rule_test.go
- [X] T009 [P] Write unit tests for Direction type marshaling/unmarshaling (verify JSON serialization to/from string values) in models/vps/securitygroups/rule_test.go
- [X] T010 Run foundational tests via `go test ./models/vps/securitygroups/rule_test.go -v`

---

## Phase 3: [US1] Correct Security Group List Response Structure (P1)

**User Story**: As a developer using the cloud SDK, I want the SecurityGroupListResponse to accurately reflect the API specification so that I can properly handle list responses from VPS API.

**Why Priority P1**: This is the core functionality needed for the SDK to work correctly with security group list operations.

**Independent Test**: Verify SecurityGroupListResponse unmarshals from API JSON without Total field and matches pb.SgListOutput.

**Acceptance Criteria**:
- SecurityGroupListResponse struct compiles without errors
- JSON unmarshaling from API responses works correctly
- No extra fields present (Total field removed)
- Migration notes documented

**Duration Estimate**: 2 hours

### Tasks

- [X] T011 [P] [US1] Write failing unit test for SecurityGroupListResponse unmarshaling in models/vps/securitygroups/securitygroup_test.go
- [X] T012 [US1] Remove Total field from SecurityGroupListResponse in models/vps/securitygroups/securitygroup.go
- [X] T013 [P] [US1] Add test case for empty security_groups array unmarshaling in models/vps/securitygroups/securitygroup_test.go
- [X] T014 [P] [US1] Add test case for marshaling SecurityGroupListResponse to JSON in models/vps/securitygroups/securitygroup_test.go
- [X] T015 [US1] Run unit tests via `go test ./models/vps/securitygroups/securitygroup_test.go -v -run TestSecurityGroupListResponse`
- [X] T016 [US1] Document breaking change in CHANGELOG.md with migration notes (use len(resp.SecurityGroups) instead of resp.Total)

**Verification**: `go test ./models/vps/securitygroups/ -v -cover` shows 100% coverage for SecurityGroupListResponse.

---

## Phase 4: [US2] Accurate Security Group Model Fields (P2)

**User Story**: As a developer working with security group data, I want the SecurityGroup and SecurityGroupRule structs to match the API specification exactly so that I can access all available fields correctly.

**Why Priority P2**: Ensures complete compatibility with the VPS API data structures.

**Independent Test**: Create SecurityGroup instances and verify all pb.SgInfo fields are accessible and serialize correctly.

**Acceptance Criteria**:
- All pb.SgInfo fields present in SecurityGroup struct
- SecurityGroupRule matches pb.SgRuleInfo exactly
- Port fields (PortMin/PortMax) have correct types
- Time fields use time.Time type

**Duration Estimate**: 2 hours

### Tasks

- [X] T017 [P] [US2] Write unit test verifying all SecurityGroup fields match pb.SgInfo in models/vps/securitygroups/securitygroup_test.go
- [X] T018 [US2] Update SecurityGroupRule to use Protocol and Direction custom types in models/vps/securitygroups/rule.go
- [X] T019 [US2] Verify PortMin/PortMax are *int in SecurityGroupRuleCreateRequest in models/vps/securitygroups/rule.go
- [X] T020 [US2] Change PortMin/PortMax to int (not *int) in SecurityGroupRule response model in models/vps/securitygroups/rule.go
- [X] T021 [P] [US2] Write unit test for SecurityGroupRule with all fields populated in models/vps/securitygroups/rule_test.go
- [X] T022 [P] [US2] Write unit test for SecurityGroupRule with ICMP protocol (ports should be 0) in models/vps/securitygroups/rule_test.go
- [X] T023 [US2] Run unit tests via `go test ./models/vps/securitygroups/ -v -cover`

**Verification**: All model fields accessible, marshaling/unmarshaling works, 100% field coverage in tests.

---

## Phase 5: [US3] Consistent JSON Serialization (P3)

**User Story**: As a developer integrating with the VPS API, I want consistent JSON field naming between request/response models and the API specification so that requests and responses are properly formatted.

**Why Priority P3**: Prevents serialization issues that could cause API failures.

**Independent Test**: Marshal request structs and compare JSON output with API examples from swagger/vps.yaml.

**Acceptance Criteria**:
- All JSON tags match API spec field names
- Create/Update requests marshal correctly
- Response models unmarshal correctly from API JSON
- Optional fields use omitempty correctly

**Duration Estimate**: 1.5 hours

### Tasks

- [X] T024 [P] [US3] Write unit test for SecurityGroupCreateRequest JSON marshaling in models/vps/securitygroups/securitygroup_test.go
- [X] T025 [P] [US3] Write unit test for SecurityGroupUpdateRequest JSON marshaling with optional fields in models/vps/securitygroups/securitygroup_test.go
- [X] T026 [P] [US3] Write unit test for SecurityGroupRuleCreateRequest JSON marshaling in models/vps/securitygroups/rule_test.go
- [X] T027 [US3] Verify all JSON tags in SecurityGroup struct match API spec in models/vps/securitygroups/securitygroup.go
- [X] T028 [US3] Verify all JSON tags in SecurityGroupRule struct match API spec in models/vps/securitygroups/rule.go
- [X] T029 [US3] Run serialization tests via `go test ./models/vps/securitygroups/ -v -run TestJSON`

**Verification**: All marshaling/unmarshaling tests pass, JSON output matches API examples.

---

## Phase 6: [US4] Complete Security Group CRUD Operations (P1)

**User Story**: As a developer managing security groups, I want full Create, Read, Update, Delete (CRUD) operations that match the API specification so that I can perform all necessary security group management tasks.

**Why Priority P1**: Core functionality for security group lifecycle management.

**Independent Test**: Perform each CRUD operation with mocked HTTP responses and verify request/response handling.

**Acceptance Criteria**:
- Create operation sends SgCreateInput format
- Get operation returns SecurityGroupResource with Rules() accessor
- Update operation sends SgUpdateInput format
- Delete operation completes without errors
- All operations accept context.Context
- Errors use sdkError type

**Duration Estimate**: 4 hours

### Tasks

- [x] T030 [P] [US4] Write contract test for Create operation in modules/vps/securitygroups/test/contract_test.go
- [x] T031 [P] [US4] Write contract test for Get operation in modules/vps/securitygroups/test/contract_test.go
- [x] T032 [P] [US4] Write contract test for Update operation in modules/vps/securitygroups/test/contract_test.go
- [x] T033 [P] [US4] Write contract test for Delete operation in modules/vps/securitygroups/test/contract_test.go
- [x] T034 [US4] Create modules/vps/securitygroups/client.go with Client struct (baseClient, projectID fields)
- [x] T035 [US4] Implement Create(ctx, req) method in modules/vps/securitygroups/client.go
- [x] T036 [US4] Implement Get(ctx, id) method returning *SecurityGroupResource in modules/vps/securitygroups/client.go
- [x] T037 [US4] Implement Update(ctx, id, req) method in modules/vps/securitygroups/client.go
- [x] T038 [US4] Implement Delete(ctx, id) method in modules/vps/securitygroups/client.go
- [x] T039 [P] [US4] Write unit tests for Create method in modules/vps/securitygroups/client_test.go
- [x] T040 [P] [US4] Write unit tests for Get method in modules/vps/securitygroups/client_test.go
- [x] T041 [P] [US4] Write unit tests for Update method in modules/vps/securitygroups/client_test.go
- [x] T042 [P] [US4] Write unit tests for Delete method in modules/vps/securitygroups/client_test.go
- [x] T043 [US4] Define SecurityGroupResource struct with embedded SecurityGroup and rulesOps field in modules/vps/securitygroups/client.go
- [x] T044 [US4] Implement Rules() method on SecurityGroupResource in modules/vps/securitygroups/client.go
- [x] T045 [US4] Create test fixtures for CRUD operations in modules/vps/securitygroups/test/fixtures.go
- [x] T046 [US4] Run CRUD tests via `go test ./modules/vps/securitygroups/... -v -run TestCRUD`

**Verification**: All CRUD operations work with mocked HTTP, contract tests pass, 80%+ coverage.

---

## Phase 7: [US5] Security Group Listing with Filters (P2)

**User Story**: As a developer listing security groups, I want to use query parameters that match the API specification so that I can filter results appropriately.

**Why Priority P2**: Essential for managing large numbers of security groups.

**Independent Test**: List with different filter combinations and verify query string generation.

**Acceptance Criteria**:
- List operation accepts ListSecurityGroupsOptions
- Query parameters (name, user_id, detail) encoded correctly
- detail=true includes rules in response
- Empty options work (lists all)
- url.Values used for query encoding

**Duration Estimate**: 2.5 hours

### Tasks

- [x] T047 [P] [US5] Write contract test for List operation without filters in modules/vps/securitygroups/test/contract_test.go
- [x] T048 [P] [US5] Write contract test for List with name filter in modules/vps/securitygroups/test/contract_test.go
- [x] T049 [P] [US5] Write contract test for List with detail=true in modules/vps/securitygroups/test/contract_test.go
- [x] T050 [US5] Create ListSecurityGroupsOptions struct in models/vps/securitygroups/securitygroup.go
- [x] T051 [US5] Implement List(ctx, opts) method in modules/vps/securitygroups/client.go
- [x] T052 [US5] Implement query parameter encoding using url.Values in modules/vps/securitygroups/client.go
- [x] T053 [P] [US5] Write unit test for List without options in modules/vps/securitygroups/client_test.go
- [x] T054 [P] [US5] Write unit test for List with all filter options in modules/vps/securitygroups/client_test.go
- [x] T055 [P] [US5] Write unit test verifying query string format in modules/vps/securitygroups/client_test.go
- [x] T056 [US5] Run list tests via `go test ./modules/vps/securitygroups/... -v -run TestList`

**Verification**: List operation works with all filter combinations, query encoding correct.

---

## Phase 8: [US6] Security Group Rule Management (P1)

**User Story**: As a developer managing security group rules, I want to create and delete rules as a sub-resource so that I can control network access policies.

**Why Priority P1**: Rules are the core functionality of security groups for network security.

**Independent Test**: Create and delete rules independently via sub-resource pattern.

**Acceptance Criteria**:
- Rules accessed via sg.Rules() method
- Create rule sends SgRuleCreateInput format
- Delete rule completes successfully
- Port parameters optional (for ICMP)
- Protocol and Direction use custom types
- Sub-resource pattern follows existing conventions

**Duration Estimate**: 3.5 hours

### Tasks

- [x] T057 [P] [US6] Write contract test for rule Create operation in modules/vps/securitygroups/test/contract_test.go
- [x] T058 [P] [US6] Write contract test for rule Delete operation in modules/vps/securitygroups/test/contract_test.go
- [x] T059 [US6] Create modules/vps/securitygroups/rules.go with RulesClient struct
- [x] T060 [US6] Implement Create(ctx, req) method in modules/vps/securitygroups/rules.go
- [x] T061 [US6] Implement Delete(ctx, ruleID) method in modules/vps/securitygroups/rules.go
- [x] T062 [P] [US6] Write unit test for rule Create with TCP protocol in modules/vps/securitygroups/rules_test.go
- [x] T063 [P] [US6] Write unit test for rule Create with ICMP protocol (no ports) in modules/vps/securitygroups/rules_test.go
- [x] T064 [P] [US6] Write unit test for rule Create with port range in modules/vps/securitygroups/rules_test.go
- [x] T065 [P] [US6] Write unit test for rule Delete in modules/vps/securitygroups/rules_test.go
- [x] T066 [US6] Define RuleOperations interface in modules/vps/securitygroups/rules.go
- [x] T067 [US6] Update SecurityGroupResource to initialize rulesOps in Get() method in modules/vps/securitygroups/client.go
- [x] T068 [US6] Run rule tests via `go test ./modules/vps/securitygroups/... -v -run TestRules`

**Verification**: Sub-resource pattern works, all rule operations functional, tests pass.

---

## Phase 9: Polish & Cross-Cutting Concerns

**Objective**: Final validation, documentation, and quality checks.

**Duration Estimate**: 2 hours

**Independent Test**: Run make check and verify all quality gates pass.

### Tasks

- [ ] T069 Run full test suite via `go test ./models/vps/securitygroups/... ./modules/vps/securitygroups/... -v -cover`
- [ ] T070 Verify 80%+ test coverage via `go test ./modules/vps/securitygroups/... -coverprofile=coverage.out && go tool cover -func=coverage.out`
- [ ] T071 Run gofmt on all files via `gofmt -s -w models/vps/securitygroups/ modules/vps/securitygroups/`
- [ ] T072 Run golangci-lint via `golangci-lint run ./models/vps/securitygroups/... ./modules/vps/securitygroups/...`
- [ ] T073 Run make check via `make check`
- [ ] T074 Update CHANGELOG.md with breaking changes, migration notes, and new features
- [ ] T075 Update README.md or module docs with security group examples
- [ ] T076 Review quickstart.md examples and verify all code compiles
- [ ] T077 Create PR with summary: breaking changes, test coverage, migration guide link
- [ ] T078 Write test for concurrent rule operations (last-write-wins behavior per FR-013)
- [ ] T079 Write test for deleting security group with existing rules (verify API behavior: success or 409 conflict)
- [ ] T080 Write test for list operation with invalid filter parameters (verify API behavior: 400 error or graceful ignore)
- [ ] T081 Write test for context cancellation propagation (verify timeout context cancels HTTP requests correctly)

---

## Dependencies Graph

### Story Completion Order

```
Setup (T001-T005)
    ↓
Foundational (T006-T010) ← Custom types
    ↓
    ├─→ [US1] Model Corrections (T011-T016) ← MVP Part 1
    │       ↓
    │   [US4] CRUD Operations (T030-T046) ← MVP Part 2
    │       ↓
    │   [US5] List with Filters (T047-T056)
    │       ↓
    ├─→ [US6] Rule Management (T057-T068) ← Can run parallel to US5
    │       ↓
    ├─→ [US2] Field Accuracy (T017-T023) ← Validation story
    │       ↓
    └─→ [US3] JSON Serialization (T024-T029) ← Validation story
            ↓
    Polish (T069-T077)
```

### Critical Path

**Longest dependency chain**: T001 → T006 → T011 → T030 → T047 → T069 → T077 (15 tasks)

**Parallel Opportunities**:
- After T010 (Foundational): US1, US2, US3 can start
- After T016 (US1): US4 and US6 can run in parallel
- After T046 (US4): US5 can start
- US2 and US3 can run anytime after US1
- Within each story: Contract tests ([P] marked) can be written in parallel

---

## Parallel Execution Examples

### Story-Level Parallelism

**After Foundational Phase Complete**:

| Developer A | Developer B | Developer C |
|-------------|-------------|-------------|
| US1: Model corrections (T011-T016) | US2: Field accuracy validation (T017-T023) | US3: JSON serialization tests (T024-T029) |

**After US1 Complete**:

| Developer A | Developer B |
|-------------|-------------|
| US4: CRUD operations (T030-T046) | US6: Rule management (T057-T068) |

### Task-Level Parallelism (Within Stories)

**US4 CRUD Operations**:
- Write all 4 contract tests in parallel (T030-T033)
- Write all 4 unit tests in parallel (T039-T042)
- Implement Create and Delete in parallel (T035, T038)
- Implement Get and Update in parallel (T036, T037)

**US5 List with Filters**:
- Write all 3 contract tests in parallel (T047-T049)
- Write all 3 unit tests in parallel (T053-T055)

**US6 Rule Management**:
- Write contract tests in parallel (T057-T058)
- Write all 4 unit tests in parallel (T062-T065)

---

## Test Coverage Strategy

### Target: 80%+ Coverage

**Coverage Breakdown by Component**:

| Component | Target | Strategy |
|-----------|--------|----------|
| models/vps/securitygroups/*.go | 100% | Unit tests for all struct methods, marshaling |
| modules/vps/securitygroups/client.go | 85%+ | Unit tests + contract tests for all CRUD operations |
| modules/vps/securitygroups/rules.go | 85%+ | Unit tests + contract tests for rule operations |

**Test Types**:

1. **Unit Tests** (same package): Implementation logic, error paths, edge cases
2. **Contract Tests** (test/ subdirectory): API compliance via httptest
3. **Integration Tests** (optional): Full lifecycle with real/staging API

**Quality Gate**: All tests must pass via `make check` before PR merge.

---

## MVP Definition

**Minimum Viable Product**: User Story 1 + User Story 4

**Rationale**: 
- US1 fixes critical model bug (Total field removal)
- US4 provides core CRUD operations
- Together they enable basic security group management
- Other stories add filters, validation, and rule management

**MVP Task Count**: 21 tasks (T001-T010, T011-T016, T030-T046)

**MVP Delivery Time**: ~8 hours (setup + foundational + US1 + US4)

**Post-MVP Priorities**:
1. US6 (Rule Management) - Completes security group functionality
2. US5 (List Filters) - Improves usability
3. US2/US3 (Validation stories) - Polish and edge cases

---

## Validation Checklist

Before marking this feature complete, verify:

- [ ] All 81 tasks completed
- [ ] 80%+ test coverage achieved
- [ ] make check passes (formatting, linting, tests)
- [ ] CHANGELOG.md updated with breaking changes
- [ ] Migration notes documented
- [ ] All contract tests pass (API compliance verified)
- [ ] Sub-resource pattern works (sg.Rules().Create())
- [ ] Custom types (Protocol, Direction) defined with constants and used consistently across all models
- [ ] Protocol constants: ProtocolTCP, ProtocolUDP, ProtocolICMP, ProtocolAny
- [ ] Direction constants: DirectionIngress, DirectionEgress
- [ ] No Total field in SecurityGroupListResponse
- [ ] All CRUD operations functional
- [ ] Query parameter encoding correct (url.Values)
- [ ] Error handling uses sdkError type
- [ ] Documentation updated (README, quickstart examples)

---

## Notes

- **TDD Mandatory**: All tests must be written before implementation per constitution
- **Breaking Change**: Total field removal requires MINOR/MAJOR version bump
- **Custom Types**: Protocol and Direction provide type safety for API enums
- **Sub-Resource Pattern**: Follows existing server.Volumes(), router.Networks() conventions
- **Query Encoding**: Use url.Values instead of manual string concatenation
- **Coverage Target**: 80%+ (exceeds constitution's 75% minimum)
- **Quality Gate**: make check must pass (gofmt, golangci-lint, all tests)

**Next Step**: Start with Phase 1 (Setup) tasks T001-T005 to validate environment.
