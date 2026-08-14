````markdown
# Implementation Plan: Fix Security Group Model

**Branch**: `001-fix-security-group-model` | **Date**: November 9, 2025 | **Spec**: [spec.md](./spec.md)
**Input**: Feature specification from `/specs/001-fix-security-group-model/spec.md`

## Summary

Fix the security group model to accurately match the VPS API specification and implement complete CRUD operations plus rule sub-resource management. The implementation will:
- Remove the non-existent `Total` field from `SecurityGroupListResponse` (breaking change)
- Ensure all model fields match `pb.SgInfo`, `pb.SgRuleInfo`, and `pb.SgListOutput` from swagger/vps.yaml
- Add custom types for `Protocol` (tcp/udp/icmp/any) and `Direction` (ingress/egress) to provide type-safe enums
- Implement security group CRUD operations (Create, Read, Update, Delete, List)
- Implement rule sub-resource operations using `SecurityGroup.Rules().Create()` and `SecurityGroup.Rules().Delete()` pattern
- Support server-side filtering via query parameters (name, user_id, detail)
- Use existing `sdkError` type consistent with network and flavor APIs
- Achieve 80%+ test coverage with unit, contract, and integration tests
- Pass `make check` quality gate

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: Standard library (`net/http`, `encoding/json`, `context`, `time`); existing internal packages (`internal/http`, `internal/backoff`, `internal/waiter`)  
**Storage**: N/A (client SDK, no persistence)  
**Testing**: Go standard `testing` package; table-driven tests; `httptest` for contract tests  
**Target Platform**: Cross-platform (Linux, macOS, Windows); library consumed by Go applications  
**Project Type**: Single-module SDK with service-specific packages under `modules/vps/securitygroups/`  
**Performance Goals**: No specific performance requirements - rely on underlying API performance (per clarification session)  
**Constraints**: 
- Use existing `sdkError` type consistent with network and flavor APIs
- Last-write-wins with optimistic locking for concurrent operations
- 30s default per-request timeout (inherited from project client)
- Max 3 retry attempts for safe reads (inherited from http client)
**Scale/Scope**: 
- Standard cloud service assumptions: hundreds of security groups per project, thousands of rules total (per clarification session)
- 5 CRUD endpoints for security groups + 2 rule sub-resource endpoints
- Query parameter filtering: name, user_id, detail

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The feature plan MUST satisfy all applicable Cloud SDK Constitution principles:

- TDD mandatory: tests are written first and initially fail; include unit tests and
  contract tests derived from the corresponding Swagger/OpenAPI.
- Public API shape: expose idiomatic Go packages and a Client with methods per API;
  do not expose raw HTTP to callers; all public methods accept `context.Context`
  and return typed responses and wrapped errors.
- Dependencies: prefer standard library; any external dependency MUST be justified
  here with clear value (security, maintainability, or functionality).
- Versioning: document whether the change is MAJOR/MINOR/PATCH and include a
  migration note for any breaking change.
- Observability: provide hooks for logging/metrics without forcing a vendor.
- Security: do not log secrets; TLS verification enabled by default; configuration
  through env vars or explicit constructors.

### Constitution Compliance Assessment

* **TDD Compliance**: ✅ **MANDATORY** - Tests must be written before implementation (unit tests alongside code, contract/integration tests in `test/` subdirectory). Target 80%+ test coverage as specified in clarification session (exceeds constitution's 75% minimum).

* **Minimal Dependencies**: ✅ **COMPLIANT** - Zero external dependencies. Uses only Go standard library (`net/http`, `encoding/json`, `context`, `time`) and existing internal packages (`internal/http`, `internal/backoff`, `internal/waiter`, `internal/types`). Pattern established by network and flavor implementations.

* **Direct Call Interfaces**: ✅ **COMPLIANT** - Exposes ergonomic client API following `Resource.SubResource().Verb()` pattern:
  - `client.SecurityGroups().List(ctx, opts)` 
  - `client.SecurityGroups().Create(ctx, req)`
  - `sg.Rules().Create(ctx, req)` (sub-resource pattern)
  - No raw HTTP exposed to consumers

* **Observability**: ✅ **COMPLIANT** - Consistent error handling using existing `sdkError` type with HTTP status codes. Error wrapping follows `fmt.Errorf` conventions for traceable error chains. No explicit logging (per project pattern - consumers handle logging).

* **Versioning**: ⚠️ **BREAKING CHANGE REQUIRED** - Removal of `Total` field from `SecurityGroupListResponse` is a breaking change. Requires:
  - **MINOR** version bump (per semantic versioning for backward-incompatible changes in pre-1.0 versions, or **MAJOR** if already 1.0+)
  - Migration notes documenting the removed field
  - CHANGELOG entry with migration guidance
  - Deprecated field marker with version removal notice (if phased deprecation preferred)

* **Code Quality**: ✅ **COMPLIANT** - Must pass `make check` validation:
  - `gofmt` formatting
  - `go vet` static analysis
  - `golangci-lint` linting
  - All unit + contract tests passing
  - Coverage thresholds met

## Project Structure

### Documentation (this feature)

```text
specs/[###-feature]/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

This feature follows the **single-module SDK structure** (Option 1) established by the existing VPS module pattern:

```text
models/vps/securitygroups/
├── securitygroup.go       # SecurityGroup model, CreateRequest, UpdateRequest, ListResponse, ListOptions
├── securitygroup_test.go  # Unit tests for model validation/marshaling
├── rule.go                # SecurityGroupRule model, RuleCreateRequest
└── rule_test.go           # Unit tests for rule model

modules/vps/securitygroups/
├── client.go              # SecurityGroupsClient with List/Create/Get/Update/Delete
├── client_test.go         # Unit tests for client methods
├── rules.go               # RulesClient with Create/Delete (sub-resource pattern)
├── rules_test.go          # Unit tests for rule methods
└── test/
    ├── contract_test.go   # Contract tests against API specification
    ├── fixtures.go        # Shared test fixtures and mock responses
    └── integration_test.go # Integration tests (optional, marked with build tags)
```

**Structure Decision**: Single-module monorepo structure is used. Models live in `models/vps/securitygroups/` for shared data types. Client implementation resides in `modules/vps/securitygroups/` with sub-resource pattern support. Tests are co-located with implementation (unit tests) or placed in `test/` subdirectory (contract/integration tests). This mirrors the established patterns in `modules/vps/flavors/`, `modules/vps/networks/`, etc.

**Key Files to Modify**:
- `models/vps/securitygroups/securitygroup.go` - Remove `Total` field from `SecurityGroupListResponse`
- `models/vps/securitygroups/rule.go` - Add custom types `Protocol` and `Direction`; verify `PortMin`/`PortMax` types match API spec (int vs *int)
- `modules/vps/securitygroups/client.go` - Implement CRUD methods
- `modules/vps/securitygroups/rules.go` - Implement sub-resource pattern (NEW FILE)

## Complexity Tracking

> **Fill ONLY if Constitution Check has violations that must be justified**

**No violations detected** - This feature fully complies with all constitution principles. The only complexity note is the sub-resource pattern implementation (`sg.Rules().Create(ctx, req)`), which is explicitly required by the user and already established in the project (e.g., `server.Volumes()`, `router.Networks()`).

---

## Iteration Plan

The implementation follows a three-phase process to ensure proper research, design, and task breakdown before coding begins.

### Phase 0: Research (`/speckit.plan` creates `research.md`)

**Objective**: Document architectural decisions, pattern choices, and dependencies before design begins.

**Deliverable**: `specs/001-fix-security-group-model/research.md`

**Research Topics**:

1. **Sub-Resource Pattern Implementation**
   - Examine existing implementations: `modules/vps/servers/volumes.go`, `modules/vps/routers/networks.go`, `modules/vps/networks/ports.go`
   - Document pattern: `Get()` returns resource wrapper (e.g., `*SecurityGroupResource`) with sub-resource methods
   - Example: `sg := client.SecurityGroups().Get(ctx, "sg-id")` then `sg.Rules().Create(ctx, req)`
   - Alternative considered: Pass parent ID to sub-resource operations (`Rules(sgID).Create(ctx, req)`) - rejected for consistency with existing pattern

2. **Query Parameter Handling**
   - Document `ListSecurityGroupsOptions` struct with optional fields: `Name *string`, `UserID *string`, `Detail *bool`
   - URL encoding: Use `url.Values` from standard library
   - Test strategy: Contract tests verify correct query parameter serialization

3. **Breaking Change Migration Strategy**
   - `Total` field removal from `SecurityGroupListResponse`
   - Semantic versioning: **MINOR** bump (breaking change in pre-1.0) or **MAJOR** (if >=1.0)
   - Migration notes template:
     ```
     ### Breaking Changes
     - `SecurityGroupListResponse.Total` field removed (not present in API spec)
     - Migration: Remove references to `.Total`; use `len(response.SecurityGroups)` for count
     - Rationale: Field never populated by API, caused confusion
     ```

4. **Port Range Type Decision**
   - API spec: `port_min` and `port_max` as integers (optional in request, always present in response)
   - Go model: `PortMin *int` and `PortMax *int` to distinguish "not set" from "0" in requests
   - Verify current implementation in `models/vps/securitygroups/rule.go`

5. **Custom Types for Protocol and Direction (Type-Safe Enums)**
   - Define custom types based on string to simulate enums for better type safety
   - `Protocol` type with constants: `ProtocolTCP`, `ProtocolUDP`, `ProtocolICMP`, `ProtocolAny`
   - `Direction` type with constants: `DirectionIngress`, `DirectionEgress`
   - Pattern:
     ```go
     type Protocol string
     const (
         ProtocolTCP  Protocol = "tcp"
         ProtocolUDP  Protocol = "udp"
         ProtocolICMP Protocol = "icmp"
         ProtocolAny  Protocol = "any"
     )
     
     type Direction string
     const (
         DirectionIngress Direction = "ingress"
         DirectionEgress  Direction = "egress"
     )
     ```
   - Benefits: Type safety, autocomplete in IDEs, compile-time validation, better documentation
   - Apply to both `SecurityGroupRule` and `SecurityGroupRuleCreateRequest` models

6. **Error Handling Patterns**
   - Reuse existing `sdkError` type from `errors.go`
   - HTTP status code mapping: 404 → NotFound, 400 → BadRequest, 409 → Conflict
   - Context cancellation: Respect `ctx.Done()` via `internal/http` client

**Dependencies**: None - all patterns already established in codebase.

### Phase 1: Design (`/speckit.plan` creates `data-model.md`, `contracts/`, `quickstart.md`)

**Objective**: Define exact models, interfaces, and quick-start examples before implementation.

**Deliverables**:

1. **`specs/001-fix-security-group-model/data-model.md`**
   - Complete Go struct definitions with JSON tags matching swagger/vps.yaml
   - Custom type definitions for `Protocol` (tcp/udp/icmp/any) and `Direction` (ingress/egress)
   - Model relationships diagram (SecurityGroup contains Rules array)
   - Breaking change documentation (removal of `Total` field)
   - Example JSON payloads for each operation (create, update, list with detail=true/false)

2. **`specs/001-fix-security-group-model/contracts/`**
   - `security-groups-client.go`: Interface definitions (package contracts)
     - Contains all interface definitions: `Client`, `SecurityGroupResource`, `RuleOperations`
     - Full GoDoc documentation with examples
     - HTTP status codes and error handling patterns
     ```go
     package contracts
     
     type Client interface {
         List(ctx context.Context, opts *securitygroups.ListSecurityGroupsOptions) (*securitygroups.SecurityGroupListResponse, error)
         Create(ctx context.Context, req *securitygroups.SecurityGroupCreateRequest) (*securitygroups.SecurityGroup, error)
         Get(ctx context.Context, id string) (*SecurityGroupResource, error)
         Update(ctx context.Context, id string, req *securitygroups.SecurityGroupUpdateRequest) (*securitygroups.SecurityGroup, error)
         Delete(ctx context.Context, id string) error
     }
     
     type SecurityGroupResource struct {
         *securitygroups.SecurityGroup
         rulesOps RuleOperations
     }
     
     type RuleOperations interface {
         Create(ctx context.Context, req *securitygroups.SecurityGroupRuleCreateRequest) (*securitygroups.SecurityGroupRule, error)
         Delete(ctx context.Context, ruleID string) error
     }
     ```

3. **`specs/001-fix-security-group-model/quickstart.md`**
   - Complete working example covering:
     - Client initialization
     - Creating a security group with rules
     - Listing security groups with filtering
     - Adding rules to existing group (sub-resource pattern)
     - Updating and deleting security groups
   - Error handling examples
   - Code comments explaining each step

**Success Criteria**:
- All models match swagger/vps.yaml exactly (verified by line-by-line comparison)
- Contract interfaces complete and idiomatic (peer review)
- Quick-start example compiles and demonstrates all operations

### Phase 2: Task Breakdown (Separate `/speckit.tasks` command creates `tasks.md`)

**Objective**: Generate implementation task list with acceptance criteria.

**Note**: Phase 2 is NOT part of the `/speckit.plan` command output. After completing Phase 0 and Phase 1, the user must run a separate `/speckit.tasks` command to generate `specs/001-fix-security-group-model/tasks.md`.

**Expected Task Structure** (for reference, not created here):
- Task 1: Remove `Total` field from SecurityGroupListResponse + unit tests
- Task 2: Implement SecurityGroupsClient CRUD methods + unit tests
- Task 3: Implement RulesClient sub-resource pattern + unit tests
- Task 4: Create contract tests with httptest and Swagger fixtures
- Task 5: Integration tests with full lifecycle scenarios
- Task 6: Documentation updates (CHANGELOG, migration notes, README)
- Task 7: Quality gate validation (make check, coverage report)

---

## Acceptance Criteria Summary

**This plan is complete when**:
1. ✅ Summary section documents feature scope and breaking change
2. ✅ Technical Context captures language, dependencies, constraints, and scale
3. ✅ Constitution Check validates compliance with all project principles
4. ✅ Project Structure defines file layout matching existing VPS module pattern
5. ✅ Complexity Tracking confirms no unjustified violations
6. ✅ Phase 0 (Research) identifies all architectural decisions and patterns
7. ✅ Phase 1 (Design) creates data-model.md, contracts/, and quickstart.md with complete interface definitions

**Next Steps**:
1. ✅ Execute research (Phase 0) - examine sub-resource pattern implementations → **COMPLETE**: `research.md` created
2. ✅ Execute design (Phase 1) - create data models, contracts, and quick-start example → **COMPLETE**: `data-model.md`, `contracts/security-groups-client.go`, `quickstart.md` created
3. ⏳ Run `/speckit.tasks` command to generate Phase 2 task breakdown → **PENDING**
4. ⏳ Begin TDD implementation following generated tasks → **PENDING**

---

## Phase Completion Status

### ✅ Phase 0: Research (COMPLETE)

**Deliverable**: [`research.md`](./research.md)

**Completed Research**:
- ✅ Sub-resource pattern implementation (examined `servers/volumes.go`, `routers/networks.go`)
- ✅ Query parameter handling (decided to use `url.Values` for clean encoding)
- ✅ Breaking change migration strategy (remove `Total` field, MINOR/MAJOR version bump)
- ✅ Port range type decision (pointers in requests, plain ints in responses)
- ✅ Error handling patterns (reuse existing `sdkError` type)
- ✅ Test strategy (unit + contract tests, 80%+ coverage target)
- ✅ Dependencies confirmed (zero external deps, stdlib only)

**Key Decisions**:
- Use Resource Wrapper pattern: `Get()` returns `*SecurityGroupResource` with `Rules()` accessor
- Query parameter encoding via `url.Values` (cleaner than manual string concatenation)
- Breaking changes: Remove `Total` field + change `PortMin`/`PortMax` to non-pointers in response model
- Migration notes required for CHANGELOG.md

### ✅ Phase 1: Design (COMPLETE)

**Deliverables**:
- ✅ [`data-model.md`](./data-model.md) - Complete Go struct definitions with JSON tags
- ✅ [`contracts/security-groups-client.go`](./contracts/security-groups-client.go) - Interface definitions with full documentation
- ✅ [`quickstart.md`](./quickstart.md) - Working examples demonstrating all operations

**Completed Design**:
- ✅ All models match `swagger/vps.yaml` exactly (verified field-by-field)
- ✅ Interface contracts defined for `Client` and `RuleOperations`
- ✅ Breaking changes documented with migration examples
- ✅ Quick-start covers full lifecycle: create → add rules → list → update → delete
- ✅ Error handling examples for 404, 409, timeout scenarios
- ✅ Sub-resource pattern demonstrated: `sg.Rules().Create(ctx, req)`

**Model Corrections Documented**:
- `SecurityGroupListResponse`: Removed `Total` field (not in API spec)
- `SecurityGroupRule`: Changed `PortMin`/`PortMax` from `*int` to `int` (API always returns values)
- `ListSecurityGroupsOptions`: Pointer types for optional filters (`*string`, `*bool`)

### ⏳ Phase 2: Task Breakdown (PENDING)

**Action Required**: Run `/speckit.tasks` command to generate `tasks.md`

This phase is intentionally separate and must be triggered after design approval.

---

## Implementation Readiness Checklist

- ✅ Constitution principles validated (no violations)
- ✅ Technical context documented (Go 1.21+, stdlib only, 80%+ coverage)
- ✅ Project structure defined (models/, modules/, test/ subdirectories)
- ✅ Research complete (patterns, decisions, dependencies)
- ✅ Data models designed (matching API spec exactly)
- ✅ Contracts defined (interfaces with full documentation)
- ✅ Quick-start example created (full lifecycle demonstration)
- ⏳ Tasks not yet generated (awaiting `/speckit.tasks` command)

**Status**: **DESIGN PHASE COMPLETE** - Ready for task generation and TDD implementation.
