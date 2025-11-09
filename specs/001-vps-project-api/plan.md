# Implementation Plan: VPS Project APIs SDK

**Branch**: `1-vps-project-api` | **Date**: 2025-10-26 | **Spec**: [spec.md](./spec.md)  
**Input**: Feature specification from `/specs/1-vps-project-api/spec.md`

## Summary

Implement a native Go SDK for VPS project-scoped APIs that:
- Exposes a typed, idiomatic Go client with per-service modules (VPS in dedicated `vps/` module)
- Wraps all `/api/v1/project/{project-id}/*` endpoints from `swagger/vps.json`
- Provides project-scoped VPS client (bind once, call methods without repeating project-id)
- Implements automatic retry for safe reads (429/502/503/504), exponential backoff + jitter
- Offers optional waiter helpers for async operations (202 Accepted)
- Returns structured errors with HTTP status code, detailed message, and meta (statusCode=0 for client-side failures)
- Enforces 30s default per-request timeout (overridable via context)
- Initializes with base URL + Bearer token
- Follows TDD: tests written first (unit + contract tests derived from Swagger)

## Technical Context

**Language/Version**: Go 1.21+  
**Primary Dependencies**: Standard library (`net/http`, `encoding/json`, `context`, `time`); consider code-gen from Swagger (e.g., `oapi-codegen` for models only, justified for type safety and maintainability)  
**Storage**: N/A (client SDK, no persistence)  
**Testing**: Go standard `testing` package; table-driven tests; httptest for contract tests  
**Target Platform**: Cross-platform (Linux, macOS, Windows); library consumed by Go applications  
**Project Type**: Single-module monorepo (github.com/Zillaforge/cloud-sdk) with service-specific packages under modules/ directory (each service gets a dedicated package folder)  
**Performance Goals**: <100ms SDK overhead per call; retry backoff tunable; non-blocking by default  
**Constraints**: <30s default timeout per request; max 3 retry attempts for safe reads; statusCode=0 for client-side errors; no secrets in logs  
**Scale/Scope**: 8 resource groups (Networks, Floating IPs, Servers, Keypairs, Routers, Security Groups, Flavors, Quotas); ~50 API operations total (16 for servers alone)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

The feature plan MUST satisfy all applicable Cloud SDK Constitution principles:

- ✅ **TDD mandatory**: All public methods will have tests written first (unit + contract from Swagger). Tests will initially fail, then pass after implementation.
- ✅ **Test Coverage**: Minimum 75% code coverage required. Unit tests MUST be in the same package as implementation (`package vps`) to ensure coverage is properly measured. Contract tests in separate packages (`package contract`) validate API contracts but don't contribute to package coverage metrics.
- ✅ **Public API shape**: Expose `cloudsdk.Client` (top-level) → `Project(projectID)` → `VPS()` returning a project-scoped `vps.Client` with typed methods (e.g., `Servers().List(ctx, opts)`, `Networks().Create(ctx, req)`). Resource operations follow `Resource.Verb()` pattern. Sub-resources follow `Resource.SubResource().Verb()` pattern (e.g., `server.NICs().List()`, `router.Networks().Associate()`, `securityGroup.Rules().Add()`). All methods accept `context.Context` and return typed structs + wrapped errors. No raw HTTP exposed.
- ✅ **Dependencies**: Prefer stdlib. Justify any external deps:
  - **oapi-codegen** (optional, models only): generates Go structs from Swagger, ensuring type safety and reducing manual boilerplate. Alternative: hand-write ~50 request/response structs (high maintenance cost, drift risk). **Decision**: Use for model generation only if justified by complexity; otherwise hand-write with Swagger as reference.
  - No other runtime dependencies planned.
- ✅ **Versioning**: Initial release v0.1.0 (MINOR). Future breaking changes (e.g., method signature changes) = MAJOR bump with migration notes in CHANGELOG.md.
- ✅ **Observability**: Provide `ClientOption` to inject a logger interface (e.g., `WithLogger(Logger)`). Retry attempts and waiter polls log at debug level without secrets.
- ✅ **Security**: Bearer token passed in `Authorization` header; never logged. TLS verification enabled by default (stdlib `http.Client`). Base URL + token configured explicitly at init.
- ✅ **Quality Gate**: **MANDATORY** - Run `make check` before claiming any implementation complete. This command executes:
  - Code formatting (`gofmt`, `goimports`)
  - Linting (`golangci-lint`) with strict error checking
  - All tests (unit, contract, integration)
  - Coverage validation
  - **Failure = Implementation NOT complete** - All errors must be fixed before proceeding to next task

**Gate Status**: ✅ PASS (no violations)

## Project Structure

### Documentation (this feature)

```text
specs/1-vps-project-api/
├── plan.md              # This file
├── research.md          # Phase 0: dependency choices, Go module layout, error wrapping patterns
├── data-model.md        # Phase 1: entities (Server, Network, etc.), error types, waiter state machines
├── quickstart.md        # Phase 1: "Get started in 5 minutes" with init + list servers example
├── contracts/           # Phase 1: extracted Swagger endpoints as Go interface signatures
│   ├── servers.go       # VPS server operations contract (test-first)
│   ├── networks.go
│   ├── floatingips.go
│   ├── keypairs.go
│   ├── routers.go
│   ├── securitygroups.go
│   ├── flavors.go
│   └── quotas.go
└── tasks.md             # Phase 2 (NOT created by /speckit.plan; created by /speckit.tasks)
```

### Source Code (repository root)

```text
/workspaces/cloud-sdk/
├── go.mod                          # Single module: github.com/Zillaforge/cloud-sdk
├── go.sum
├── README.md                       # SDK overview, installation, quick start link
├── CHANGELOG.md                    # Semantic versioning, migration notes
├── client.go                       # Top-level Client, Project() selector, ClientOption functional options (WithLogger, WithTimeout, WithHTTPClient)
├── client_test.go
├── errors.go                       # Structured error re-exports from internal/types
├── errors_test.go
├── models/                         # Model definitions organized by service
│   └── vps/                        # VPS models (import: github.com/Zillaforge/cloud-sdk/models/vps/...)
│       ├── servers/                # Server models and sub-resources
│       │   ├── server.go           # Server, ServerCreateRequest, ServerUpdateRequest, ServerListResponse, ListServersOptions
│       │   ├── server_test.go
│       │   ├── nics.go             # NIC (Network Interface Card) sub-resource models
│       │   ├── nics_test.go
│       │   ├── volumes.go          # Volume sub-resource models
│       │   └── volumes_test.go
│       ├── networks/               # Network models and sub-resources
│       │   ├── network.go          # Network, NetworkCreateRequest, NetworkUpdateRequest, NetworkListResponse, ListNetworksOptions
│       │   ├── network_test.go
│       │   ├── ports.go            # NetworkPort sub-resource models
│       │   └── ports_test.go
│       ├── floatingips/            # Floating IP models
│       │   ├── floatingip.go       # FloatingIP, FloatingIPCreateRequest, FloatingIPUpdateRequest, FloatingIPListResponse, ListFloatingIPsOptions
│       │   └── floatingip_test.go
│       ├── keypairs/               # Keypair models
│       │   ├── keypair.go          # Keypair, KeypairCreateRequest, KeypairListResponse, ListKeypairsOptions
│       │   └── keypair_test.go
│       ├── routers/                # Router models and sub-resources
│       │   ├── router.go           # Router, RouterCreateRequest, RouterUpdateRequest, RouterListResponse, ListRoutersOptions
│       │   ├── router_test.go
│       │   ├── networks.go         # Router-Network association sub-resource models
│       │   └── networks_test.go
│       ├── securitygroups/         # Security group models and sub-resources
│       │   ├── securitygroup.go    # SecurityGroup, SecurityGroupCreateRequest, SecurityGroupUpdateRequest, SecurityGroupListResponse, ListSecurityGroupsOptions
│       │   ├── securitygroup_test.go
│       │   ├── rules.go            # SecurityGroupRule sub-resource models
│       │   └── rules_test.go
│       ├── flavors/                # Flavor models
│       │   ├── flavor.go           # Flavor, FlavorListResponse, ListFlavorsOptions
│       │   └── flavor_test.go
│       └── quotas/                 # Quota models
│           ├── quota.go            # Quota, QuotaResponse
│           └── quota_test.go
├── modules/                        # Service-specific packages (not separate modules)
│   └── vps/                        # VPS service package (import: github.com/Zillaforge/cloud-sdk/modules/vps)
│       ├── client.go               # Project-scoped VPS client coordinator (bound to project-id); imports all resource submodules
│       ├── client_test.go
│       ├── flavors/                # Flavor resource submodule
│       │   ├── client.go           # Flavor list/get operations
│       │   ├── client_test.go      # Unit tests (same package for coverage)
│       │   └── test/               # Contract and integration tests (separate package)
│       │       └── [test files]
│       ├── floatingips/            # Floating IP resource submodule
│       │   ├── client.go           # Floating IP CRUD, approve/reject, disassociate operations
│       │   ├── client_test.go      # ✅ Unit tests (same package, 97.5% coverage)
│       │   └── test/               # ✅ Contract and integration tests (separate package)
│       │       ├── floatingips_approve_test.go
│       │       ├── floatingips_create_test.go
│       │       ├── floatingips_delete_test.go
│       │       ├── floatingips_disassociate_test.go
│       │       ├── floatingips_get_test.go
│       │       ├── floatingips_integration_test.go
│       │       ├── floatingips_list_test.go
│       │       ├── floatingips_reject_test.go
│       │       └── floatingips_update_test.go
│       ├── keypairs/               # Keypair resource submodule
│       │   ├── client.go           # Keypair CRUD operations
│       │   ├── client_test.go      # Unit tests (same package for coverage)
│       │   └── test/               # Contract and integration tests (separate package)
│       │       └── [test files]
│       ├── networks/               # Network resource submodule
│       │   ├── client.go           # Network CRUD operations; NetworkResource with Ports sub-resource
│       │   ├── client_test.go      # ✅ Unit tests (same package, 95.0% coverage)
│       │   ├── ports.go            # PortsClient for network port operations
│       │   ├── ports_test.go       # ✅ Unit tests for ports operations
│       │   └── test/               # ✅ Contract and integration tests (separate package)
│       │       ├── network_ports_test.go
│       │       ├── networks_create_test.go
│       │       ├── networks_delete_test.go
│       │       ├── networks_get_test.go
│       │       ├── networks_integration_test.go
│       │       ├── networks_list_test.go
│       │       └── networks_update_test.go
│       ├── quotas/                 # Quota resource submodule
│       │   ├── client.go           # Quota get operations
│       │   ├── client_test.go      # Unit tests (same package for coverage)
│       │   └── test/               # Contract and integration tests (separate package)
│       │       └── [test files]
│       ├── routers/                # Router resource submodule
│       │   ├── client.go           # Router CRUD, set_state; RouterResource with Networks sub-resource
│       │   ├── client_test.go      # Unit tests (same package for coverage)
│       │   └── test/               # Contract and integration tests (separate package)
│       │       └── [test files]
│       ├── securitygroups/         # Security group resource submodule
│       │   ├── client.go           # Security group CRUD; SecurityGroupResource with Rules sub-resource
│       │   ├── client_test.go      # Unit tests (same package for coverage)
│       │   └── test/               # Contract and integration tests (separate package)
│       │       └── [test files]
│       └── servers/                # Server resource submodule
│           ├── client.go           # Server CRUD, actions, metrics; ServerResource with NICs/Volumes sub-resources
│           ├── client_test.go      # Unit tests (same package for coverage)
│           └── test/               # Contract and integration tests (separate package)
│               └── [test files]
├── internal/                       # Internal helpers (not exported)
│   ├── http/
│   │   ├── client.go               # Shared HTTP executor (retry, timeout, error wrap)
│   │   └── client_test.go
│   ├── backoff/
│   │   ├── backoff.go              # Exponential backoff + jitter
│   │   └── backoff_test.go
│   └── waiter/
│       ├── waiter.go               # Generic waiter framework (poll state with context, backoff, max wait - reusable across all services)
│       └── waiter_test.go
└── testdata/                       # Swagger fixture, mock responses
    └── vps.json -> ../../swagger/vps.json (symlink or copy)
```

**Structure Decision**: 
- **Single-module monorepo**: One Go module (`github.com/Zillaforge/cloud-sdk`) with service-specific packages under `modules/` directory (`modules/vps/`, `modules/compute/`, `modules/storage/`, etc.). Each service is a Go package, not a separate module. This simplifies dependency management, versioning, and cross-service shared code while maintaining clear service boundaries.
- **Model organization**: Models are organized in `models/vps/` with subdirectories for each of the 8 resource types (servers, networks, floatingips, keypairs, routers, securitygroups, flavors, quotas). Sub-resource models (e.g., NICs, Volumes, Ports, Rules) are located within their parent resource directory. This structure provides clear organization while maintaining logical grouping.
- **VPS module organization** (✅ **IMPLEMENTED**): Each of the 8 VPS resource types has its own subpackage under `modules/vps/` (flavors, floatingips, keypairs, networks, quotas, routers, securitygroups, servers). The main `modules/vps/client.go` acts as a coordinator that imports all submodules and provides typed accessor methods (e.g., `Networks()`, `FloatingIPs()`, `Servers()`). This structure:
  - Provides clear separation of concerns with dedicated files per resource type
  - Makes it easier to maintain and test each resource independently
  - Follows Go best practices for package organization
  - Enables parallel development on different resource types
  - Keeps the codebase modular and scalable as more resources are added
- **Project-scoped pattern**: Top-level `Client` provides `Project(projectID)` → service client factory (e.g., `.VPS()`) to eliminate project-id repetition in method calls.
- **Service baseURL construction** (✅ **IMPLEMENTED**): Each service client factory method appends the service name to the baseURL before creating the service client. This pattern:
  - **URL Format**: `<HOST>/<SERVICE>/api/v1/project/<PROJECT_ID>` (e.g., `https://api.example.com/vps/api/v1/project/proj-123`)
  - **Implementation**: Service factory method (e.g., `VPS()`) appends `/vps` to baseURL: `vpsBaseURL := pc.client.baseURL + "/vps"`
  - **Benefits**: Clean separation of concerns; service name defined at client creation point; individual endpoint methods remain service-agnostic
  - **Reusability**: Easy pattern for future services (e.g., `storageBaseURL := pc.client.baseURL + "/storage"`)
  - **Example**: 
    ```go
    func (pc *ProjectClient) VPS() *vps.Client {
        vpsBaseURL := pc.client.baseURL + "/vps"
        return vps.NewClient(vpsBaseURL, pc.client.token, pc.projectID, pc.client.httpClient, pc.client.logger)
    }
    ```
- **Resource.Verb pattern**: Main resource operations use short method names (e.g., `List`, `Create`, `Get`, `Update`, `Delete`) where the resource type provides context (e.g., `Servers().List()`).
- **Sub-resource pattern**: Nested resources accessed via `Resource.SubResource().Verb()` (e.g., `server.NICs().List()`, `router.Networks().Associate()`, `securityGroup.Rules().Add()`). Get() operations return resource wrappers (e.g., `*ServerResource`, `*NetworkResource`) that provide sub-resource access methods.
- **Internal packages**: Shared HTTP/retry/backoff logic in `internal/` to avoid API surface bloat and enable reuse across services.
- **TDD structure**: Every `.go` file has a corresponding `_test.go` with table-driven tests; contract tests use `httptest` and Swagger fixtures.
- **Test organization** (✅ **IMPLEMENTED**): Each resource submodule follows a two-tier test structure:
  - **Unit tests** (same package, e.g., `package floatingips`): Located in the resource root directory (e.g., `modules/vps/floatingips/client_test.go`). These tests verify implementation logic and contribute to package coverage metrics. Required for meeting the 75% coverage threshold.
  - **Contract & Integration tests** (separate package, e.g., `package floatingips_test`): Located in a `test/` subdirectory within each resource (e.g., `modules/vps/floatingips/test/`). These tests validate API behavior against Swagger specs and full lifecycle scenarios. They use `httptest` for mocking and don't contribute to the parent package's coverage metrics.
  - **Rationale**: Separating test types by location and package enables:
    - Clear distinction between unit tests (white-box, implementation-focused) and contract tests (black-box, API-focused)
    - Accurate coverage measurement (unit tests in same package count toward coverage)
    - Organized test files per resource without cluttering the main package directory
    - Easy parallel test execution (`go test ./modules/vps/floatingips/...` runs all tests for that resource)
  - **Example structure**:
    ```
    modules/vps/floatingips/
    ├── client.go              # Implementation
    ├── client_test.go         # Unit tests (package floatingips)
    └── test/                  # Contract & integration tests
        ├── floatingips_create_test.go      (package floatingips_test)
        ├── floatingips_list_test.go        (package floatingips_test)
        └── floatingips_integration_test.go (package floatingips_test)
    ```

## Build & Quality Tooling

### Makefile Configuration

**Purpose**: Provide consistent commands for development workflows (build, test, lint, format, coverage)

**Critical requirement**: Exclude `specs/` directory from all Go tooling to prevent errors when scanning markdown files.

**Implementation patterns**:

```makefile
# Package discovery with specs exclusion
PACKAGES := $(shell go list ./... | grep -v '/specs')
GO_FILES := $(shell find . -name '*.go' -not -path './specs/*' -not -path './vendor/*')

# Quality targets
.PHONY: test
test:
	go test $(PACKAGES) -v -cover

.PHONY: lint
lint:
	@command -v golangci-lint >/dev/null 2>&1 || \
		(echo "golangci-lint not found. Run: make install-tools" && exit 1)
	golangci-lint run ./internal/... ./modules/... .

.PHONY: fmt
fmt:
	gofmt -s -w $(GO_FILES)
	@command -v goimports >/dev/null 2>&1 && \
		goimports -w $(GO_FILES) || true

.PHONY: build
build:
	go build $(PACKAGES)

.PHONY: coverage
coverage:
	go test $(PACKAGES) -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html

.PHONY: check
check: fmt lint test
	@echo "All checks passed!"

.PHONY: install-tools
install-tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install golang.org/x/tools/cmd/goimports@latest
```

**Key patterns**:
- **Package filtering**: Use `grep -v '/specs'` to exclude specs from `go list ./...`
- **File filtering**: Use `find` with `-not -path './specs/*'` for file-based operations
- **Explicit paths for linters**: golangci-lint gets explicit paths (`./internal/... ./modules/... .`) instead of `./...` to reliably exclude specs
- **Graceful degradation**: Missing tools show helpful error messages with install instructions
- **Tool installation**: Separate `install-tools` target for optional tooling setup

**Rationale**:
- `specs/` contains Markdown documentation that breaks Go parsing (`.md` files are not Go code)
- golangci-lint's `exclude-dirs` config doesn't work reliably; explicit path specification is more robust
- Optional tools (goimports) degrade gracefully without blocking the build
- Consistent target names (`test`, `lint`, `fmt`, `build`, `check`) align with Go community conventions

**Testing targets**:
- `make test`: Run all tests with coverage summary
- `make test-race`: Run tests with race detector (requires CGO, optional in Alpine)
- `make coverage`: Generate HTML coverage report
- `make check`: **MANDATORY quality gate** - Run before claiming implementation complete

**See also**: research.md section 7 for testing strategy details

## Quality Gate: `make check` (MANDATORY)

**Requirement**: Every implementation MUST pass `make check` before being marked as complete.

**What it does**:
1. **Format** (`gofmt`, `goimports`): Auto-formats all Go code to standard style
2. **Lint** (`golangci-lint`): Runs strict linters including:
   - `unused-parameter`: Detects unused function parameters (must use `_` for intentionally unused params)
   - `errcheck`: Ensures all errors are handled
   - `govet`: Catches common mistakes
   - And 20+ other linters
3. **Test**: Runs all unit, contract, and integration tests
4. **Coverage**: Validates minimum coverage thresholds

**Exit codes**:
- `0`: All checks passed ✅
- `1`: Formatting issues (auto-fixed by `make fmt`)
- `1`: Lint errors (must fix manually)
- `1`: Test failures (must fix implementation or tests)
- `2`: Make error (tooling issue)

**Workflow**:
```bash
# After implementing a feature
make check

# If it fails:
# 1. Read the error output carefully
# 2. Fix the issues (format, lint, test failures)
# 3. Re-run make check
# 4. Repeat until exit code 0

# Only after make check passes:
# Mark task as complete in tasks.md
```

**Common fixes**:
- **Unused parameter**: Rename `r *http.Request` to `_ *http.Request` in test handlers
- **Unhandled error**: Use `_ = json.NewEncoder(w).Encode(...)` pattern for test mocks
- **Test timeout**: Add `&http.Client{Timeout: 5 * time.Second}` to internal HTTP client in tests
- **Import ordering**: Run `make fmt` to auto-fix

**Why mandatory**:
- Ensures code quality and consistency
- Catches bugs early (especially in error handling)
- Prevents accumulation of technical debt
- Makes code review faster (no style/lint discussions)
- Builds confidence in implementation correctness

## Test Coverage Strategy

**Requirement**: Minimum 75% overall code coverage

**Implementation Success**:
- **Problem Solved**: Contract tests in separate package didn't contribute to coverage
- **Solution Implemented**: Created comprehensive unit tests in same package (`package vps`)

**Coverage Strategy** (Three-tier Testing):
1. **Contract tests** (separate package `contract`): Validate API contracts, HTTP behavior, error handling against Swagger spec
2. **Unit tests** (same package): Test implementation logic, ensure coverage metrics
   - ✅ **floatingips**: Comprehensive unit tests with 97.5% coverage (440 lines, 9 test functions)
   - ✅ **networks**: Comprehensive unit tests with 95.0% coverage (400 lines, 7 test functions)
   - ✅ **ports**: Unit tests for port operations (131 lines, 2 test functions)
   - ⏳ **servers, flavors, keypairs, quotas, routers, securitygroups**: Pending implementation
3. **Integration tests** (package `integration`): Validate full lifecycles end-to-end
   - ✅ **floatingips_integration_test.go**: Full lifecycle tests for floating IPs
   - ✅ **networks_integration_test.go**: Full lifecycle tests for networks

**Testing Pattern** (Table-driven tests):
```go
tests := []struct {
    name         string
    mockResponse interface{}
    mockStatus   int
    expectError  bool
}{
    {
        name:         "success case",
        mockResponse: /* valid response */,
        mockStatus:   http.StatusOK,
        expectError:  false,
    },
    {
        name:         "error case",
        mockResponse: /* error response */,
        mockStatus:   http.StatusBadRequest,
        expectError:  true,
    },
}
```

**Quality Assurance**:
- All test files use `httptest.NewServer` with proper timeout configuration (`&http.Client{Timeout: 5 * time.Second}`)
- Linter compliance enforced via `make check` (goimports + errcheck)
- Error handling in test mocks uses explicit ignore pattern (`_ = json.NewEncoder(w).Encode(...)`)

## Complexity Tracking

> No Constitution violations; this section is empty.
