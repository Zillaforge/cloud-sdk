# Tasks: Fix Keypair Model

**Input**: Design documents from `/specs/006-fix-keypair-model/`
**Prerequisites**: plan.md (✅), spec.md (✅), research.md (✅), data-model.md (✅), contracts/ (✅)

**Tests**: Tests are MANDATORY for public APIs per the Constitution. Write tests FIRST
and ensure they FAIL before implementation. Include unit tests and contract tests
derived from Swagger/OpenAPI for each endpoint wrapper.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

## Path Conventions

Go SDK (single project): `models/`, `modules/`, `internal/`, tests co-located with implementation

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

✅ **COMPLETE**: All Phase 1 work was completed during design phase:
- Created feature specification in `/workspaces/cloud-sdk/specs/006-fix-keypair-model/spec.md`
- Created implementation plan in `/workspaces/cloud-sdk/specs/006-fix-keypair-model/plan.md`
- Created research documentation in `/workspaces/cloud-sdk/specs/006-fix-keypair-model/research.md`
- Created data model documentation in `/workspaces/cloud-sdk/specs/006-fix-keypair-model/data-model.md`
- Created contracts in `/workspaces/cloud-sdk/specs/006-fix-keypair-model/contracts/keypair-api.md`
- Created quickstart guide in `/workspaces/cloud-sdk/specs/006-fix-keypair-model/quickstart.md`
- Created breaking changes guide in `/workspaces/cloud-sdk/specs/006-fix-keypair-model/CHANGES.md`

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

✅ **COMPLETE**: All Phase 2 implementation tasks completed:
- Added IDName type to `models/vps/keypairs` package
- Updated Keypair model with all missing fields (PrivateKey, User, CreatedAt, UpdatedAt)
- Removed Total field from KeypairListResponse (BREAKING CHANGE)
- Updated List() method to return []*Keypair directly
- Applied Request/Response naming convention
- All code compiles successfully

### Implementation Tasks for Phase 2

- [x] T001 [P] Add IDName type to models/vps/keypairs/keypair.go for user reference
- [x] T002 [P] Add PrivateKey field to Keypair struct with proper JSON tag
- [x] T003 [P] Add User field to Keypair struct as pointer to IDName
- [x] T004 [P] Add CreatedAt field to Keypair struct with camelCase JSON tag
- [x] T005 [P] Add UpdatedAt field to Keypair struct with camelCase JSON tag
- [x] T006 Remove Total field from KeypairListResponse struct
- [x] T007 Update List() method return type from *KeypairListResponse to []*Keypair
- [x] T008 Update List() method implementation to return slice directly
- [x] T009 Verify all code compiles successfully after changes
- [x] T010 Run go fmt on modified files

**Checkpoint**: Foundation ready - user story testing can now begin

---

## Phase 3: User Story 1 - Correct Keypair Model Structure (Priority: P1) 🎯 MVP

**Goal**: Ensure the Keypair model accurately reflects the API specification (pb.KeypairInfo) so developers can properly handle keypair data including timestamps and user information.

**Independent Test**: Verify that Keypair structs can be unmarshaled from API responses containing all fields (including createdAt, updatedAt, user, private_key) and marshaled back without data loss.

### Tests for User Story 1 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [x] T011 [P] [US1] Create unit test for Keypair JSON unmarshaling in models/vps/keypairs/keypair_test.go
- [x] T012 [P] [US1] Create unit test for Keypair JSON marshaling in models/vps/keypairs/keypair_test.go
- [x] T013 [P] [US1] Create test for optional field handling (nil User, empty PrivateKey) in models/vps/keypairs/keypair_test.go
- [x] T014 [P] [US1] Create test for private_key field presence in Create response in models/vps/keypairs/keypair_test.go
- [x] T015 [P] [US1] Create test for private_key field absence in Get/List responses in models/vps/keypairs/keypair_test.go
- [x] T016 [P] [US1] Create contract test validating against pb.KeypairInfo schema in models/vps/keypairs/keypair_test.go
- [x] T017 [P] [US1] Create test for IDName struct JSON serialization in models/vps/keypairs/keypair_test.go

### Validation for User Story 1

- [x] T018 [US1] Run all Keypair model tests and verify they pass
- [x] T019 [US1] Verify Keypair struct matches pb.KeypairInfo field-by-field from swagger/vps.yaml
- [x] T020 [US1] Validate JSON tags match Swagger spec exactly (camelCase for createdAt/updatedAt, snake_case for public_key/private_key/user_id)

**Checkpoint**: At this point, User Story 1 should be fully functional - Keypair model correctly serializes/deserializes all fields

---

## Phase 4: User Story 2 - Keypair List Response Structure (Priority: P2)

**Goal**: Ensure the list response matches the API specification (pb.KeypairListOutput) so developers can correctly iterate through keypairs without relying on a non-existent "total" field.

**Independent Test**: Call the List operation and verify the response structure matches pb.KeypairListOutput with direct slice access (no Total field).

### Tests for User Story 2 (MANDATORY - TDD) ⚠️

- [x] T021 [P] [US2] Create unit test for List() return type in modules/vps/keypairs/client_test.go
- [x] T022 [P] [US2] Create test for direct slice access (len, index) in modules/vps/keypairs/client_test.go
- [x] T023 [P] [US2] Create test for empty keypairs list handling in modules/vps/keypairs/client_test.go
- [x] T024 [P] [US2] Create test for List with name filter in modules/vps/keypairs/client_test.go
- [x] T025 [P] [US2] Create contract test validating List response against pb.KeypairListOutput schema in modules/vps/keypairs/client_test.go

### Validation for User Story 2

- [x] T026 [US2] Run all List operation tests and verify they pass
- [x] T027 [US2] Verify List() returns []*Keypair directly (not wrapped struct)
- [x] T028 [US2] Confirm KeypairListResponse has no Total field
- [x] T029 [US2] Test that len(result) works for counting keypairs

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently - List operations return correct structure

---

## Phase 5: User Story 3 - Timestamp Handling (Priority: P3)

**Goal**: Ensure the Keypair model includes creation and update timestamps so developers can monitor when keypairs were created or modified for security auditing.

**Independent Test**: Verify timestamp fields are properly handled in JSON operations and can be parsed into Go time.Time types.

### Tests for User Story 3 (MANDATORY - TDD) ⚠️

- [x] T030 [P] [US3] Create test for timestamp parsing from API response in models/vps/keypairs/keypair_test.go
- [x] T031 [P] [US3] Create test for RFC3339 format validation in models/vps/keypairs/keypair_test.go
- [x] T032 [P] [US3] Create test for timestamp conversion to time.Time in models/vps/keypairs/keypair_test.go
- [x] T033 [P] [US3] Create test for createdAt immutability (same after update) in models/vps/keypairs/keypair_test.go
- [x] T034 [P] [US3] Create test for updatedAt changes after update in models/vps/keypairs/keypair_test.go

### Validation for User Story 3

- [x] T035 [US3] Run all timestamp-related tests and verify they pass
- [x] T036 [US3] Verify timestamps are preserved as strings (not converted)
- [x] T037 [US3] Test parsing timestamps with time.Parse(time.RFC3339, timestamp)
- [x] T038 [US3] Validate error handling for malformed timestamp strings

**Checkpoint**: All user stories should now be independently functional - complete timestamp support

---

## Phase 6: Integration & Contract Validation

**Purpose**: End-to-end testing across all operations

### Integration Tests

- [x] T039 [P] Create integration test for Create→Get→Update→Delete flow in modules/vps/keypairs/client_test.go
- [x] T040 [P] Create integration test for keypair import (existing public key) in modules/vps/keypairs/client_test.go
- [x] T041 [P] Create integration test for keypair generation (no public key) in modules/vps/keypairs/client_test.go
- [x] T042 Create test for private_key availability only during Create in modules/vps/keypairs/client_test.go

### Contract Tests Against Swagger

- [x] T043 [P] Validate Create request body matches KeypairCreateInput schema in modules/vps/keypairs/client_test.go
- [x] T044 [P] Validate Update request body matches KeypairUpdateInput schema in modules/vps/keypairs/client_test.go
- [x] T045 [P] Validate Create response matches pb.KeypairInfo schema in modules/vps/keypairs/client_test.go
- [x] T046 [P] Validate Get response matches pb.KeypairInfo schema in modules/vps/keypairs/client_test.go
- [x] T047 [P] Validate List response matches pb.KeypairListOutput schema in modules/vps/keypairs/client_test.go

### Edge Cases & Error Handling

- [ ] T048 [P] Test error handling for invalid keypair ID in modules/vps/keypairs/client_test.go
- [ ] T049 [P] Test error handling for missing required fields in modules/vps/keypairs/client_test.go
- [ ] T050 [P] Test error handling for invalid public key format in modules/vps/keypairs/client_test.go
- [ ] T051 Test context cancellation during operations in modules/vps/keypairs/client_test.go

**Checkpoint**: Complete integration validation across all keypair operations

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

### Documentation

- [ ] T052 [P] Update root CHANGELOG.md with v1.0.0 breaking changes entry
- [ ] T053 [P] Add migration section to main README.md for Total field removal
- [ ] T054 [P] Generate godoc comments for all public types in models/vps/keypairs/keypair.go
- [ ] T055 [P] Generate godoc comments for all public methods in modules/vps/keypairs/client.go
- [ ] T056 Update examples in modules/vps/keypairs/EXAMPLES.md with new Request/Response naming

### Code Quality

- [ ] T057 [P] Run go fmt on all modified files
- [ ] T058 [P] Run go vet on models/vps/keypairs and modules/vps/keypairs packages
- [ ] T059 Run golangci-lint on all modified code
- [ ] T060 Verify no new external dependencies added

### Security Review

- [ ] T061 Review private_key handling in all code examples and tests
- [ ] T062 Verify sensitive field documentation warns about immediate storage
- [ ] T063 Confirm no private keys logged or exposed in error messages

### Constitution Compliance

- [ ] T064 Validate all public APIs accept context.Context
- [ ] T065 Verify error wrapping follows "failed to <action>: %w" pattern
- [ ] T066 Confirm breaking changes documented with migration guide
- [ ] T067 Verify semantic versioning reflects MAJOR bump
- [ ] T068 Final check: All constitution gates from plan.md pass

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: ✅ COMPLETE
- **Foundational (Phase 2)**: ✅ COMPLETE - BLOCKS all user stories
- **User Stories (Phase 3-5)**: All depend on Foundational phase completion
  - Can proceed in parallel if team capacity allows
  - Or sequentially in priority order (P1 → P2 → P3)
- **Integration (Phase 6)**: Depends on all user stories (Phase 3-5) being complete
- **Polish (Phase 7)**: Depends on Integration (Phase 6) completion

### User Story Dependencies

- **User Story 1 (P1)**: 🚧 Implementation PENDING (Phase 2 complete, ready for testing)
- **User Story 2 (P2)**: 🚧 Implementation PENDING (Phase 2 complete, ready for testing)
- **User Story 3 (P3)**: 🚧 Implementation PENDING (Phase 2 complete, ready for testing)

Phase 2 foundation is complete - ready to begin TDD testing for user stories

### Within Each User Story Phase

1. Write ALL tests for the story (T001-T007 for US1)
2. Run tests and verify they FAIL appropriately
3. Validate implementation (already complete)
4. Run tests and verify they PASS
5. Story checkpoint reached

### Parallel Opportunities

**Within User Story 1**:
- T001-T007: All test creation tasks can run in parallel (different test functions)

**Within User Story 2**:
- T011-T015: All test creation tasks can run in parallel (different test functions)

**Within User Story 3**:
- T020-T024: All test creation tasks can run in parallel (different test functions)

**Phase 6 Integration**:
- T029-T032: Integration test tasks can run in parallel
- T033-T037: Contract test tasks can run in parallel
- T038-T041: Error handling test tasks can run in parallel

**Phase 7 Polish**:
- T042-T046: Documentation tasks can run in parallel
- T047-T049: Code quality tasks can run in parallel
- T051-T053: Security review tasks can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task T011: "Create unit test for Keypair JSON unmarshaling in models/vps/keypairs/keypair_test.go"
Task T012: "Create unit test for Keypair JSON marshaling in models/vps/keypairs/keypair_test.go"
Task T013: "Create test for optional field handling in models/vps/keypairs/keypair_test.go"
Task T014: "Create test for private_key field presence in Create response in models/vps/keypairs/keypair_test.go"
Task T015: "Create test for private_key field absence in Get/List responses in models/vps/keypairs/keypair_test.go"
Task T016: "Create contract test validating against pb.KeypairInfo schema in models/vps/keypairs/keypair_test.go"
Task T017: "Create test for IDName struct JSON serialization in models/vps/keypairs/keypair_test.go"

# Then validation sequentially:
Task T018: "Run all Keypair model tests and verify they pass"
Task T019: "Verify Keypair struct matches pb.KeypairInfo field-by-field"
Task T020: "Validate JSON tags match Swagger spec exactly"
```

---

## Implementation Strategy

### Current Status: Foundation Complete, Ready for TDD Testing

✅ **Completed**:
1. Phase 1: Setup - All documentation created
2. Phase 2: Foundational - All implementation tasks complete
   - IDName type added to models/vps/keypairs
   - Keypair model updated with all pb.KeypairInfo fields
   - Total field removed (BREAKING CHANGE)
   - List() returns []*Keypair directly
   - All code compiles successfully

🚧 **Next Steps**:
1. **Phase 3**: Write and run tests for User Story 1 (T011-T020)
2. **Phase 4**: Write and run tests for User Story 2 (T021-T029)
3. **Phase 5**: Write and run tests for User Story 3 (T030-T038)
4. **Phase 6**: Integration and contract validation (T039-T051)
5. **Phase 7**: Polish, documentation, and final validation (T052-T068)

### Next Steps

1. **Phase 2**: Complete implementation tasks T001-T010
2. **Phase 3**: Write and run tests for User Story 1 (T011-T020)
3. **Phase 4**: Write and run tests for User Story 2 (T021-T029)
4. **Phase 5**: Write and run tests for User Story 3 (T030-T038)
5. **Phase 6**: Integration and contract validation (T039-T051)
6. **Phase 7**: Polish, documentation, and final validation (T052-T068)

### MVP First (User Story 1 Only)

1. ✅ Complete Phase 1: Setup
2. ✅ Complete Phase 2: Foundational
3. 🚧 Complete Phase 3: User Story 1 (ready for TDD testing)
4. **STOP and VALIDATE**: Test User Story 1 independently
5. Deploy/demo if ready

### Incremental Delivery

1. ✅ Complete Setup + Foundational → Foundation ready
2. 🚧 Add User Story 1 tests → Validate independently → MVP ready
3. Add User Story 2 tests → Validate independently → Enhanced release
4. Add User Story 3 tests → Validate independently → Complete feature
5. Integration tests → Final validation
6. Polish → Production ready

### Testing Focus

Since foundation is complete, focus on:
1. **TDD Validation**: Write tests that validate the implemented code
2. **Contract Compliance**: Ensure all responses match Swagger exactly
3. **Edge Cases**: Test optional fields, null values, error conditions
4. **Migration**: Verify breaking changes are properly handled

---

## Summary

**Total Tasks**: 68
- Phase 1 (Setup): ✅ Complete (done during design)
- Phase 2 (Foundational): ✅ Complete (implementation complete)
- Phase 3 (User Story 1): 10 tasks (8 tests + 2 validation)
- Phase 4 (User Story 2): 9 tasks (5 tests + 4 validation)
- Phase 5 (User Story 3): 9 tasks (5 tests + 4 validation)
- Phase 6 (Integration): 13 tasks (4 integration + 5 contract + 4 error)
- Phase 7 (Polish): 17 tasks (5 docs + 4 quality + 3 security + 5 compliance)

**Parallel Opportunities**: 55 tasks marked [P] can run in parallel within their phase

**Independent Test Criteria**:
- US1: Keypair structs serialize/deserialize all fields correctly
- US2: List() returns []*Keypair with direct slice access
- US3: Timestamps parse correctly as RFC3339 strings

**MVP Scope**: User Story 1 only (correct Keypair model structure)

**Next Action**: Begin Phase 3 by writing tests T011-T017 for User Story 1 (ensure they FAIL first, then validate implementation)
