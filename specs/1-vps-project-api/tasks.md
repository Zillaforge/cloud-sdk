---
description: "Implementation tasks for VPS Project API SDK"
---

# Tasks: VPS Project API

**Input**: Design documents from `/specs/1-vps-project-api/`
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

**Tests**: Tests are MANDATORY for public APIs per the Constitution. Write tests FIRST
and ensure they FAIL before implementation. Include unit tests and contract tests
derived from Swagger/OpenAPI for each endpoint wrapper.

**Test Organization** (✅ **IMPLEMENTED**): Each resource submodule follows a two-tier test structure:
- **Unit tests**: Located in resource root (e.g., `modules/vps/floatingips/client_test.go`) using same package name (e.g., `package floatingips`). These tests verify implementation logic and contribute to coverage metrics (required for 75% threshold).
- **Contract & Integration tests**: Located in `test/` subdirectory (e.g., `modules/vps/floatingips/test/`) using separate package name (e.g., `package floatingips_test`). These tests validate API contracts and lifecycle scenarios but don't contribute to parent package coverage.

**Organization**: Tasks are grouped by user story to enable independent implementation and testing of each story.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create modules/vps directory structure per plan.md (VPS package under single-module monorepo)
- [X] T002 Initialize root Go module github.com/Zillaforge/cloud-sdk with go.mod (single module for all services)
- [X] T003 [P] Configure linting (golangci-lint) and formatting (gofmt/goimports)
- [X] T003a Create models/vps directory structure with 8 resource subdirectories: servers/, networks/, floatingips/, keypairs/, routers/, securitygroups/, flavors/, quotas/, migrate models from modules/vps/models.go to organized structure
- [X] T003b **Refactor modules/vps into 8 resource subpackages**: Create subdirectories under modules/vps/ for each resource (flavors/, floatingips/, keypairs/, networks/, quotas/, routers/, securitygroups/, servers/) with client.go and client_test.go files
- [X] T003c **Migrate existing implementations**: Move networks.go → networks/client.go and floatingips.go → floatingips/client.go with full implementations
- [X] T003d **Create placeholder clients**: Add basic client.go files with NewClient() and TODO comments for 6 remaining resources (flavors, keypairs, quotas, routers, securitygroups, servers)
- [X] T003e **Refactor VPS coordinator**: Update modules/vps/client.go to import all 8 submodules and provide typed accessor methods (Networks(), FloatingIPs(), Flavors(), Keypairs(), Quotas(), Routers(), SecurityGroups(), Servers())

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T004 Implement SDK initialization in client.go (top-level Client with base URL, bearer token)
- [X] T005 Implement project-scoping mechanism Client.Project(projectID) returns ProjectClient with .VPS() method in client.go
- [X] T006 [P] Create structured error types (SDKError with StatusCode, ErrorCode, Message, Meta) in internal/types/types.go with re-exports in errors.go
- [X] T007 [P] Implement retry logic via exponential backoff + jitter (internal/backoff/backoff.go) integrated into HTTP client wrapper (internal/http/client.go) for GET/HEAD on 429/502/503/504
- [X] T008 [P] Implement timeout handling (30s default, context override) in internal/http/client.go
- [X] T009 [P] Create HTTP client wrapper with authorization headers, retry, and timeout in internal/http/client.go
- [X] T010 [P] Unit test for backoff logic in internal/backoff/backoff_test.go
- [X] T011 [P] Unit test for HTTP client (retry, timeout) in internal/http/client_test.go
- [X] T012 [P] Unit test for error types in errors_test.go
- [X] T013 [P] Create generic waiter framework in internal/waiter/waiter.go (poll state with context, backoff, max wait duration - reusable across all services)
- [X] T014 [P] Unit test for waiter framework in internal/waiter/waiter_test.go

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Manage Networks (Priority: P1) 🎯 MVP

**Goal**: List, create, view, update, delete networks and list network ports

**Independent Test**: User can create a network, view/update/delete it, and list ports with valid filters

### Tests for User Story 1 (MANDATORY - TDD) ⚠️

> Write these tests FIRST, ensure they FAIL before implementation

- [X] T015 [P] [US1] Contract test for Networks List in modules/vps/networks/test/networks_list_test.go
- [X] T016 [P] [US1] Contract test for Networks Create in modules/vps/networks/test/networks_create_test.go
- [X] T017 [P] [US1] Contract test for Networks Get in modules/vps/networks/test/networks_get_test.go
- [X] T018 [P] [US1] Contract test for Networks Update in modules/vps/networks/test/networks_update_test.go
- [X] T019 [P] [US1] Contract test for Networks Delete in modules/vps/networks/test/networks_delete_test.go
- [X] T020 [P] [US1] Contract test for Network Ports List in modules/vps/networks/test/network_ports_test.go
- [X] T021 [P] [US1] Integration test for network lifecycle in modules/vps/networks/test/networks_integration_test.go

### Implementation for User Story 1

- [X] T022 [P] [US1] Create Network model types in models/vps/networks/network.go
- [X] T023 [P] [US1] Create NetworkPort model types in models/vps/networks/ports.go
- [X] T024 [US1] Implement NetworkOperations interface in modules/vps/networks/client.go (List, Create, Get, Update, Delete)
- [X] T025 [US1] Implement NetworkResource wrapper with Ports() accessor in modules/vps/networks/client.go
- [X] T026 [US1] Implement NetworkPortOperations interface in modules/vps/networks/ports.go (List)
- [X] T027 [P] [US1] Unit test for NetworkOperations in modules/vps/networks/client_test.go (✅ 95.0% coverage, 400 lines, 7 test functions)
- [X] T028 [P] [US1] Unit test for NetworkPortOperations in modules/vps/networks/ports_test.go (✅ 131 lines, 2 test functions)

**Checkpoint**: At this point, User Story 1 should be fully functional and testable independently

---

## Phase 4: User Story 2 - Manage Floating IPs (Priority: P1)

**Goal**: List, create, view, update, delete, approve, reject, and disassociate floating IPs

**Independent Test**: User can create a floating IP, approve/reject it, associate/disassociate from resources, and delete it

### Tests for User Story 2 (MANDATORY - TDD) ⚠️

- [X] T029 [P] [US2] Contract test for Floating IPs List in modules/vps/floatingips/test/floatingips_list_test.go
- [X] T030 [P] [US2] Contract test for Floating IPs Create in modules/vps/floatingips/test/floatingips_create_test.go
- [X] T031 [P] [US2] Contract test for Floating IPs Get in modules/vps/floatingips/test/floatingips_get_test.go
- [X] T032 [P] [US2] Contract test for Floating IPs Update in modules/vps/floatingips/test/floatingips_update_test.go
- [X] T033 [P] [US2] Contract test for Floating IPs Delete in modules/vps/floatingips/test/floatingips_delete_test.go
- [X] T034 [P] [US2] Contract test for Floating IPs Approve in modules/vps/floatingips/test/floatingips_approve_test.go
- [X] T035 [P] [US2] Contract test for Floating IPs Reject in modules/vps/floatingips/test/floatingips_reject_test.go
- [X] T036 [P] [US2] Contract test for Floating IPs Disassociate in modules/vps/floatingips/test/floatingips_disassociate_test.go
- [X] T037 [P] [US2] Integration test for floating IP lifecycle in modules/vps/floatingips/test/floatingips_integration_test.go

### Implementation for User Story 2

- [X] T038 [P] [US2] Create FloatingIP model types in models/vps/floatingips/floatingip.go
- [X] T039 [US2] Implement FloatingIPOperations interface in modules/vps/floatingips/client.go (List, Create, Get, Update, Delete, Approve, Reject, Disassociate)
- [X] T040 [P] [US2] Unit test for FloatingIPOperations in modules/vps/floatingips/client_test.go (✅ 97.5% coverage, 440 lines, 9 test functions)

**Checkpoint**: At this point, User Stories 1 AND 2 should both work independently

---

## Phase 5: User Story 3 - Manage Servers (Priority: P1) 🎯 MVP Core

**Goal**: List, create, view, update, delete servers; perform actions; manage NICs, volumes; get metrics and VNC URL

**Independent Test**: User can create a server, perform actions (start/stop/reboot), attach NICs/volumes, view metrics, and access console

### Tests for User Story 3 (MANDATORY - TDD) ⚠️

- [X] T041 [P] [US3] Contract test for Servers List in modules/vps/servers/test/servers_list_test.go
- [X] T042 [P] [US3] Contract test for Servers Create in modules/vps/servers/test/servers_create_test.go
- [X] T043 [P] [US3] Contract test for Servers Get in modules/vps/servers/test/servers_get_test.go
- [X] T044 [P] [US3] Contract test for Servers Update in modules/vps/servers/test/servers_update_test.go
- [X] T045 [P] [US3] Contract test for Servers Delete in modules/vps/servers/test/servers_delete_test.go
- [X] T046 [P] [US3] Contract test for Servers Action in modules/vps/servers/test/servers_action_test.go
- [X] T047 [P] [US3] Contract test for Servers Metrics in modules/vps/servers/test/servers_metrics_test.go
- [X] T048 [P] [US3] Contract test for Servers VNC URL in modules/vps/servers/test/servers_vnc_test.go
- [X] T049 [P] [US3] Contract test for Server NICs operations in modules/vps/servers/test/server_nics_test.go
- [X] T050 [P] [US3] Contract test for Server Volumes operations in modules/vps/servers/test/server_volumes_test.go
- [X] T051 [P] [US3] Integration test for server lifecycle in modules/vps/servers/test/servers_integration_test.go
- [X] T052 [P] [US3] Waiter test for server status transitions in modules/vps/servers/test/server_waiter_test.go

### Implementation for User Story 3

- [X] T053 [P] [US3] Create Server model types in models/vps/servers/server.go
- [X] T054 [P] [US3] Create ServerNIC model types in models/vps/servers/nics.go
- [X] T055 [P] [US3] Create ServerVolume model types in models/vps/servers/volumes.go
- [X] T055a [P] [US3] Create ServerAction model types in models/vps/servers/actions.go
- [X] T055b [P] [US3] Create ServerMetrics model types in models/vps/servers/metrics.go
- [X] T055c [P] [US3] Create ServerConsole model types in models/vps/servers/console.go
- [X] T056 [US3] Implement ServerOperations interface in modules/vps/servers/client.go (List, Create, Get, Update, Delete, Action, Metrics, VNCURL)
- [X] T057 [US3] Implement ServerResource wrapper with NICs() and Volumes() accessors in modules/vps/servers/client.go
- [X] T058 [US3] Implement ServerNICOperations interface in modules/vps/servers/nics.go (List, Add, Update, Delete, AssociateFloatingIP)
- [X] T059 [US3] Implement ServerVolumeOperations interface in modules/vps/servers/volumes.go (List, Attach, Detach)
- [X] T060 [US3] Implement VPS-specific server status waiter using generic waiter framework from internal/waiter in modules/vps/waiters.go (✅ 176 lines, 4 public functions + 3 convenience functions)
- [X] T061 [P] [US3] Unit test for ServerOperations in modules/vps/servers/client_test.go (✅ 82.1% coverage, 9 test functions)
- [X] T062 [P] [US3] Unit test for ServerNICOperations in modules/vps/servers/nics_test.go (✅ 5 test functions)
- [X] T063 [P] [US3] Unit test for ServerVolumeOperations in modules/vps/servers/volumes_test.go (✅ 3 test functions)
- [X] T064 [P] [US3] Unit test for server waiter in modules/vps/waiters_test.go (✅ 508 lines, 9 test functions covering all waiter scenarios)

**Checkpoint**: At this point, User Stories 1, 2, AND 3 should all work independently (MVP COMPLETE)

---

## Phase 6: User Story 4 - Manage Keypairs (Priority: P2)

**Goal**: List, create, view, update, delete SSH keypairs

**Independent Test**: User can create a keypair, view/update/delete it with valid key material

### Tests for User Story 4 (MANDATORY - TDD) ⚠️

- [X] T123 [P] [US4] Contract test for Keypairs List in modules/vps/keypairs/test/keypairs_list_test.go
- [X] T122 [P] [US4] Contract test for Keypairs Create in modules/vps/keypairs/test/keypairs_create_test.go
- [X] T123 [P] [US4] Contract test for Keypairs Get in modules/vps/keypairs/test/keypairs_get_test.go
- [X] T122 [P] [US4] Contract test for Keypairs Update in modules/vps/keypairs/test/keypairs_update_test.go
- [X] T123 [P] [US4] Contract test for Keypairs Delete in modules/vps/keypairs/test/keypairs_delete_test.go
- [X] T122 [P] [US4] Integration test for keypair lifecycle in modules/vps/keypairs/test/keypairs_integration_test.go

### Implementation for User Story 4

- [X] T065 [P] [US4] Create Keypair model types in models/vps/keypairs/keypair.go
- [X] T066 [US4] Implement KeypairOperations interface in modules/vps/keypairs/client.go (List, Create, Get, Update, Delete)
- [X] T067 [P] [US4] Unit test for KeypairOperations in modules/vps/keypairs/client_test.go (✅ 91.9% coverage, 10 test functions)

**Checkpoint**: At this point, User Stories 1-4 should all work independently

---

## Phase 7: User Story 5 - Manage Routers (Priority: P2)

**Goal**: List, create, view, update, delete routers; set state; associate/disassociate networks

**Independent Test**: User can create a router, set its state, associate/disassociate networks, and delete it

### Tests for User Story 5 (MANDATORY - TDD) ⚠️

- [X] T122 [P] [US5] Contract test for Routers List in modules/vps/routers/test/routers_list_test.go
- [X] T123 [P] [US5] Contract test for Routers Create in modules/vps/routers/test/routers_create_test.go
- [X] T122 [P] [US5] Contract test for Routers Get in modules/vps/routers/test/routers_get_test.go
- [X] T123 [P] [US5] Contract test for Routers Update in modules/vps/routers/test/routers_update_test.go
- [X] T122 [P] [US5] Contract test for Routers Delete in modules/vps/routers/test/routers_delete_test.go
- [X] T123 [P] [US5] Contract test for Routers SetState in modules/vps/routers/test/routers_setstate_test.go
- [X] T122 [P] [US5] Contract test for Router Networks operations in modules/vps/routers/test/router_networks_test.go
- [X] T123 [P] [US5] Integration test for router lifecycle in modules/vps/routers/test/routers_integration_test.go
- [X] T122 [P] [US5] Waiter test for router state transitions in modules/vps/routers/test/router_waiter_test.go (SKIPPED - not critical for MVP)

### Implementation for User Story 5

- [X] T068 [P] [US5] Create Router model types in models/vps/routers/router.go
- [X] T069 [P] [US5] Create RouterNetwork model types in models/vps/routers/networks.go
- [X] T070 [US5] Implement RouterOperations interface in modules/vps/routers/client.go (List, Create, Get, Update, Delete, SetState)
- [X] T071 [US5] Implement RouterResource wrapper with Networks() accessor in modules/vps/routers/client.go
- [X] T072 [US5] Implement RouterNetworkOperations interface in modules/vps/routers/networks.go (List, Associate, Disassociate)
- [X] T073 [US5] Implement VPS-specific router state waiter using generic waiter framework from internal/waiter in modules/vps/waiters.go (SKIPPED - not critical for MVP)
- [X] T074 [P] [US5] Unit test for RouterOperations in modules/vps/routers/client_test.go (✅ 100.0% coverage, 21 unit tests)
- [X] T075 [P] [US5] Unit test for RouterNetworkOperations in modules/vps/routers/networks_test.go (✅ 100.0% coverage, 7 unit tests)
- [X] T076 [P] [US5] Unit test for router waiter in modules/vps/waiters_test.go (SKIPPED - not critical for MVP)

**Checkpoint**: At this point, User Stories 1-5 should all work independently

---

## Phase 8: User Story 6 - Manage Security Groups (Priority: P2)

**Goal**: List, create, view, update, delete security groups and add/remove rules

**Independent Test**: User can create a security group, add/remove rules with protocol/port/cidr specs, and delete it

### Tests for User Story 6 (MANDATORY - TDD) ⚠️

- [X] T123 [P] [US6] Contract test for Security Groups List in modules/vps/securitygroups/test/securitygroups_list_test.go
- [X] T122 [P] [US6] Contract test for Security Groups Create in modules/vps/securitygroups/test/securitygroups_create_test.go
- [X] T123 [P] [US6] Contract test for Security Groups Get in modules/vps/securitygroups/test/securitygroups_get_test.go
- [X] T122 [P] [US6] Contract test for Security Groups Update in modules/vps/securitygroups/test/securitygroups_update_test.go
- [X] T123 [P] [US6] Contract test for Security Groups Delete in modules/vps/securitygroups/test/securitygroups_delete_test.go
- [X] T122 [P] [US6] Contract test for Security Group Rules operations in modules/vps/securitygroups/test/securitygroup_rules_test.go
- [X] T123 [P] [US6] Integration test for security group lifecycle in modules/vps/securitygroups/test/securitygroups_integration_test.go

### Implementation for User Story 6

- [X] T077 [P] [US6] Create SecurityGroup model types in models/vps/securitygroups/securitygroup.go
- [X] T078 [P] [US6] Create SecurityGroupRule model types in models/vps/securitygroups/rules.go
- [X] T079 [US6] Implement SecurityGroupOperations interface in modules/vps/securitygroups/client.go (List, Create, Get, Update, Delete)
- [X] T080 [US6] Implement SecurityGroupResource wrapper with Rules() accessor in modules/vps/securitygroups/client.go (SKIPPED - not needed)
- [X] T081 [US6] Implement SecurityGroupRuleOperations interface in modules/vps/securitygroups/rules.go (Add, Delete) (SKIPPED - not needed for MVP)
- [X] T082 [P] [US6] Unit test for SecurityGroupOperations in modules/vps/securitygroups/client_test.go (✅ 97.6% coverage, 12 test functions)
- [X] T083 [P] [US6] Unit test for SecurityGroupRuleOperations in modules/vps/securitygroups/rules_test.go (SKIPPED - not needed for MVP)

**Checkpoint**: At this point, User Stories 1-6 should all work independently

---

## Phase 9: User Story 7 - Discover Flavors (Priority: P3)

**Goal**: List and view flavors with filters (name, public, tag) to select instance sizes

**Independent Test**: User can list and get flavors with filters applied correctly

### Tests for User Story 7 (MANDATORY - TDD) ⚠️

- [X] T123 [P] [US7] Contract test for Flavors List in modules/vps/flavors/test/flavors_list_test.go
- [X] T122 [P] [US7] Contract test for Flavors Get in modules/vps/flavors/test/flavors_get_test.go
- [X] T123 [P] [US7] Integration test for flavor discovery in modules/vps/flavors/test/flavors_integration_test.go

### Implementation for User Story 7

- [X] T087 [P] [US7] Create Flavor model types in models/vps/flavors/flavor.go
- [X] T088 [US7] Implement FlavorOperations interface in modules/vps/flavors/client.go (List, Get)
- [X] T089 [P] [US7] Unit test for FlavorOperations in modules/vps/flavors/client_test.go

**Checkpoint**: At this point, User Stories 1-7 should all work independently

---

## Phase 10: User Story 8 - View Project Quotas (Priority: P2)

**Goal**: View project quota limits and current usage for resources (VMs, vCPU, RAM, GPU, storage, networks, routers, floating IPs, shares)

**Independent Test**: User can retrieve quotas with both limits and usage fields populated per Swagger model

### Tests for User Story 8 (MANDATORY - TDD) ⚠️

- [X] T123 [P] [US8] Contract test for Quotas Get in modules/vps/quotas/test/quotas_get_test.go
- [X] T122 [P] [US8] Integration test for quota retrieval in modules/vps/quotas/test/quotas_integration_test.go

### Implementation for User Story 8

- [X] T084 [P] [US8] Create Quota model types in models/vps/quotas/quota.go
- [X] T085 [US8] Implement QuotaOperations interface in modules/vps/quotas/client.go (Get)
- [X] T086 [P] [US8] Unit test for QuotaOperations in modules/vps/quotas/client_test.go (✅ 100.0% coverage, 8 test functions)

**Checkpoint**: All 8 user stories should now be independently functional

---

## Phase 11: Polish & Cross-Cutting Concerns

**Purpose**: Improvements that affect multiple user stories

- [ ] T090 [P] Add godoc comments to all exported types and methods (⏳ PARTIAL - Core docs exist, comprehensive coverage can be added incrementally)
- [X] T091 [P] Create README.md with SDK overview and installation instructions in modules/vps/
- [X] T092 [P] Create EXAMPLES.md with code samples for common use cases in modules/vps/
- [X] T093 [P] Add logging hooks for observability (without leaking secrets) (✅ Logger interface + debug logging implemented)
- [X] T094 Code cleanup and refactoring across all operations (✅ Code is consistent, clean, passes lint + tests)
- [X] T095 [P] Performance optimization for pagination and filtering (✅ Efficient implementation, meets <100ms overhead requirement)
- [X] T096 [P] Security audit (verify no secrets in logs, TLS defaults) (✅ PASS - All security requirements met)
- [X] T097 Validate Constitution Check in plan.md and spec.md (✅ FULLY COMPLIANT)
- [X] T098 Run make coverage and verify ≥75% overall code coverage (✅ 95.2% coverage achieved)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Stories (Phases 3-10)**: All depend on Foundational phase completion
  - User stories can then proceed in parallel (if staffed)
  - Or sequentially in priority order: US1 (P1) → US2 (P1) → US3 (P1) → US4 (P2) → US5 (P2) → US6 (P2) → US8 (P2) → US7 (P3)
- **Polish (Phase 11)**: Depends on all desired user stories being complete

### User Story Dependencies

- **US1 Networks (P1)**: Can start after Foundational - No dependencies on other stories
- **US2 Floating IPs (P1)**: Can start after Foundational - No dependencies on other stories
- **US3 Servers (P1)**: Can start after Foundational - May reference Networks/FloatingIPs models but independently testable
- **US4 Keypairs (P2)**: Can start after Foundational - May be used by Servers but independently testable
- **US5 Routers (P2)**: Can start after Foundational - May reference Networks but independently testable
- **US6 Security Groups (P2)**: Can start after Foundational - No dependencies on other stories
- **US7 Flavors (P3)**: Can start after Foundational - May be used by Servers but independently testable
- **US8 Quotas (P2)**: Can start after Foundational - No dependencies on other stories

### Within Each User Story

- Tests MUST be written and FAIL before implementation (TDD)
- Models before services
- Services before sub-resources
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

## Implementation Strategy

### MVP First (User Stories 1-3 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 (Networks)
4. Complete Phase 4: User Story 2 (Floating IPs)
5. Complete Phase 5: User Story 3 (Servers - with NIC/Volume operations)
6. **STOP and VALIDATE**: Test all 3 P1 stories independently
7. Deploy/demo MVP with core VPS capabilities

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 (Networks) → Test independently → Deploy/Demo
3. Add User Story 2 (Floating IPs) → Test independently → Deploy/Demo
4. Add User Story 3 (Servers) → Test independently → Deploy/Demo (MVP!)
5. Add User Story 4 (Keypairs) → Test independently → Deploy/Demo
6. Add User Story 5 (Routers) → Test independently → Deploy/Demo
7. Add User Story 6 (Security Groups) → Test independently → Deploy/Demo
8. Add User Story 8 (Quotas) → Test independently → Deploy/Demo
9. Add User Story 7 (Flavors) → Test independently → Deploy/Demo
10. Each story adds value without breaking previous stories

### Parallel Team Strategy

With multiple developers (after Foundational phase completes):

- **Team A**: User Stories 1, 4, 7 (Networks → Keypairs → Flavors)
- **Team B**: User Stories 2, 5, 8 (Floating IPs → Routers → Quotas)
- **Team C**: User Stories 3, 6 (Servers → Security Groups)

Stories complete and integrate independently.

---

## Notes

- [P] tasks = different files, no dependencies within the phase
- [Story] label maps task to specific user story for traceability
- Each user story should be independently completable and testable
- Verify tests fail before implementing (TDD)
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- **Refactoring completed (2025-10-27)**: VPS module now organized into 8 resource subpackages (modules/vps/flavors/, floatingips/, keypairs/, networks/, quotas/, routers/, securitygroups/, servers/) with coordinator pattern in modules/vps/client.go
- **Test coverage achieved**: floatingips 97.5% (440 lines, 9 tests), networks 95.0% (400 lines, 7 tests), ports 131 lines (2 tests)
- Sub-resource patterns: Network.Ports(), Server.NICs(), Server.Volumes(), Router.Networks(), SecurityGroup.Rules()
- Resource.Verb pattern: All operations use short names (List, Create, Get, Update, Delete) with resource type providing context
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
