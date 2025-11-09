# Research: Security Group Model and Sub-Resource Pattern

**Feature**: Fix Security Group Model (001-fix-security-group-model)  
**Date**: November 9, 2025  
**Related**: [spec.md](./spec.md) | [plan.md](./plan.md)

## Overview

This document records architectural decisions, pattern analysis, and dependency choices for implementing security group CRUD operations with rule sub-resource support. All decisions align with existing VPS module patterns established in `modules/vps/servers/`, `modules/vps/networks/`, and `modules/vps/routers/`.

---

## 1. Sub-Resource Pattern Implementation

### Pattern Analysis

The SDK uses a **Resource Wrapper** pattern to expose sub-resource operations. After examining existing implementations:

**Reference Implementation**: `modules/vps/servers/client.go` + `volumes.go` + `nics.go`

```go
// Pattern: Get() returns resource wrapper with sub-resource methods
type ServerResource struct {
    *servers.Server              // Embedded full server model
    nicOps    NICOperations      // Interface for sub-resource operations
    volumeOps VolumeOperations
}

func (sr *ServerResource) NICs() NICOperations {
    return sr.nicOps
}

func (sr *ServerResource) Volumes() VolumeOperations {
    return sr.volumeOps
}

// Sub-resource client holds parent IDs for path construction
type VolumesClient struct {
    baseClient *internalhttp.Client
    projectID  string
    serverID   string  // Parent resource ID
}

func (c *VolumesClient) Attach(ctx context.Context, volumeID string) error {
    path := fmt.Sprintf("/api/v1/project/%s/servers/%s/volumes/%s", 
                        c.projectID, c.serverID, volumeID)
    // ... HTTP request using path
}
```

**Usage Flow**:
```go
// 1. Get parent resource (returns wrapper)
server, err := client.Servers().Get(ctx, "server-123")

// 2. Access sub-resource operations
volumes := server.Volumes()

// 3. Perform sub-resource operation
err = volumes.Attach(ctx, "volume-456")
```

### Decision: Apply Pattern to Security Groups

**Implementation Plan**:

```go
// SecurityGroupResource wrapper (in modules/vps/securitygroups/client.go)
type SecurityGroupResource struct {
    *securitygroups.SecurityGroup  // Embedded model from models/vps/securitygroups/
    rulesOps RuleOperations         // Interface for rule operations
}

func (sgr *SecurityGroupResource) Rules() RuleOperations {
    return sgr.rulesOps
}

// RulesClient implementation (in modules/vps/securitygroups/rules.go)
type RulesClient struct {
    baseClient       *internalhttp.Client
    projectID        string
    securityGroupID  string  // Parent resource ID
}

func (c *RulesClient) Create(ctx context.Context, req *securitygroups.SecurityGroupRuleCreateRequest) (*securitygroups.SecurityGroupRule, error) {
    path := fmt.Sprintf("/api/v1/project/%s/security_groups/%s/rules", 
                        c.projectID, c.securityGroupID)
    // ... POST request
}

func (c *RulesClient) Delete(ctx context.Context, ruleID string) error {
    path := fmt.Sprintf("/api/v1/project/%s/security_groups/%s/rules/%s", 
                        c.projectID, c.securityGroupID, ruleID)
    // ... DELETE request
}
```

**Client.Get() Implementation**:
```go
func (c *Client) Get(ctx context.Context, sgID string) (*SecurityGroupResource, error) {
    path := fmt.Sprintf("/api/v1/project/%s/security_groups/%s", c.projectID, sgID)
    
    req := &internalhttp.Request{
        Method: "GET",
        Path:   path,
    }
    
    var sg securitygroups.SecurityGroup
    if err := c.baseClient.Do(ctx, req, &sg); err != nil {
        return nil, err
    }
    
    // Wrap with sub-resource operations
    return &SecurityGroupResource{
        SecurityGroup: &sg,
        rulesOps: &RulesClient{
            baseClient:      c.baseClient,
            projectID:       c.projectID,
            securityGroupID: sgID,
        },
    }, nil
}
```

### Alternative Considered: Direct Parent ID Parameter

**Rejected Approach**:
```go
// REJECTED: Pass parent ID to every sub-resource operation
client.SecurityGroupRules(sgID).Create(ctx, req)
client.SecurityGroupRules(sgID).Delete(ctx, ruleID)
```

**Why Rejected**:
- Inconsistent with existing `server.Volumes()`, `router.Networks()` pattern
- Requires repeating parent ID for every operation
- Less ergonomic API (more verbose)
- Breaks established SDK conventions

**Conclusion**: Use Resource Wrapper pattern for consistency.

---

## 2. Query Parameter Handling

### Pattern Analysis

**Reference Implementation**: `modules/vps/servers/client.go` (List method)

```go
func (c *Client) List(ctx context.Context, opts *servers.ListServersOptions) (*servers.ServerListResponse, error) {
    path := fmt.Sprintf("/api/v1/project/%s/servers", c.projectID)
    
    // Build query parameters
    if opts != nil {
        query := ""
        if opts.Name != "" {
            query += fmt.Sprintf("name=%s&", opts.Name)
        }
        if opts.UserID != "" {
            query += fmt.Sprintf("user_id=%s&", opts.UserID)
        }
        // ... more parameters
        if query != "" {
            path = fmt.Sprintf("%s?%s", path, query[:len(query)-1]) // Remove trailing &
        }
    }
    // ... make request
}
```

**Issues with Current Pattern**:
1. Manual string concatenation (error-prone, no URL encoding)
2. Trailing `&` removal hack
3. No handling of special characters in parameter values

### Decision: Use `url.Values` from Standard Library

**Improved Implementation**:

```go
import "net/url"

func (c *Client) List(ctx context.Context, opts *securitygroups.ListSecurityGroupsOptions) (*securitygroups.SecurityGroupListResponse, error) {
    path := fmt.Sprintf("/api/v1/project/%s/security_groups", c.projectID)
    
    // Build query parameters using url.Values
    if opts != nil {
        query := url.Values{}
        if opts.Name != nil {
            query.Set("name", *opts.Name)
        }
        if opts.UserID != nil {
            query.Set("user_id", *opts.UserID)
        }
        if opts.Detail != nil {
            query.Set("detail", fmt.Sprintf("%t", *opts.Detail))
        }
        
        if len(query) > 0 {
            path = fmt.Sprintf("%s?%s", path, query.Encode())  // Automatic URL encoding
        }
    }
    
    req := &internalhttp.Request{
        Method: "GET",
        Path:   path,
    }
    // ... make request
}
```

**Model Definition**:
```go
// ListSecurityGroupsOptions defines optional filters for listing security groups.
type ListSecurityGroupsOptions struct {
    Name   *string  // Filter by name (exact match)
    UserID *string  // Filter by user_id
    Detail *bool    // Include rules array in response (default: false)
}
```

**Benefits**:
- Automatic URL encoding (handles special characters correctly)
- Cleaner code (no manual string concatenation)
- Consistent with Go idiomatic practices
- Pointer types distinguish "not set" from "empty string"/"false"

---

## 3. Breaking Change Migration Strategy

### Issue: `Total` Field Not in API Specification

**Current Implementation** (`models/vps/securitygroups/securitygroup.go`):
```go
type SecurityGroupListResponse struct {
    SecurityGroups []*SecurityGroup `json:"security_groups"`
    Total          int              `json:"total"`  // ❌ NOT IN API SPEC
}
```

**API Specification** (`swagger/vps.yaml`, `pb.SgListOutput`):
```yaml
pb.SgListOutput:
  type: object
  properties:
    security_groups:
      type: array
      items:
        $ref: '#/definitions/pb.SgInfo'
  # No 'total' field
```

### Decision: Remove Field + Semantic Versioning

**Corrected Model**:
```go
type SecurityGroupListResponse struct {
    SecurityGroups []*SecurityGroup `json:"security_groups"`
    // Total field removed - not present in API specification
}
```

**Version Bump Strategy**:

**If SDK version < 1.0** (e.g., v0.5.2 → v0.6.0):
- **MINOR** version bump (breaking changes allowed in 0.x per semantic versioning)
- Migration note: "Breaking change acceptable in pre-1.0 versions"

**If SDK version ≥ 1.0** (e.g., v1.2.3 → v2.0.0):
- **MAJOR** version bump (breaking changes require major version increment)
- Follow full deprecation cycle if desired (mark deprecated in v1.x, remove in v2.0)

**Migration Notes Template** (for CHANGELOG.md):

```markdown
## [0.6.0] - 2025-11-09

### ⚠️ Breaking Changes

#### Security Group Model Corrections

**Removed Field**: `SecurityGroupListResponse.Total`

**Reason**: The `Total` field was never present in the VPS API specification (`pb.SgListOutput` in swagger/vps.yaml). The field was incorrectly included in the SDK model and was never populated by API responses.

**Migration Guide**:

Before (v0.5.x):
```go
resp, err := client.SecurityGroups().List(ctx, opts)
if err != nil {
    return err
}
fmt.Printf("Total: %d\n", resp.Total)  // Always 0, never populated
```

After (v0.6.0):
```go
resp, err := client.SecurityGroups().List(ctx, opts)
if err != nil {
    return err
}
fmt.Printf("Count: %d\n", len(resp.SecurityGroups))  // Actual count
```

**Action Required**: 
- Remove any references to `.Total` from your code
- Use `len(response.SecurityGroups)` to count results
- If pagination is needed, track offsets manually

**Impact**: Low - field was always zero/uninitialized, unlikely to be used in production
```

### Deprecation Alternative (Optional)

If gradual migration is preferred:

```go
// In v0.5.x
type SecurityGroupListResponse struct {
    SecurityGroups []*SecurityGroup `json:"security_groups"`
    Total          int              `json:"total" deprecated:"true"` // Deprecated: Not in API spec, will be removed in v0.6.0
}

// In v0.6.0: Remove field entirely
```

**Recommendation**: Direct removal is acceptable because:
1. Field was never populated by API (always zero value)
2. Removes confusion about unused fields
3. Aligns model with actual API specification
4. Pre-1.0 versions allow breaking changes

---

## 4. Port Range Type Decision

### API Specification Analysis

**From `swagger/vps.yaml`** (`pb.SgRuleInfo`):

```yaml
pb.SgRuleInfo:
  type: object
  properties:
    id:
      type: string
    direction:
      type: string
    protocol:
      type: string
    port_min:
      type: integer
      format: int32
    port_max:
      type: integer
      format: int32
    remote_cidr:
      type: string
```

**From `swagger/vps.yaml`** (`SgRuleCreateInput`):

```yaml
SgRuleCreateInput:
  type: object
  required:
    - direction
    - protocol
    - remote_cidr
  properties:
    direction:
      type: string
    protocol:
      type: string
    port_min:
      type: integer
      format: int32
    port_max:
      type: integer
      format: int32
    remote_cidr:
      type: string
```

**Observations**:
- `port_min` and `port_max` are **optional** in `SgRuleCreateInput` (not in `required` list)
- In `pb.SgRuleInfo` (response), ports are **always present** (integers, may be 0)
- Protocol can be `tcp`, `udp`, `icmp` (ICMP rules typically ignore port fields)

### Decision: Use Pointer Types for Optional Fields

**Model Definition** (`models/vps/securitygroups/rule.go`):

```go
// SecurityGroupRule represents a security group rule.
type SecurityGroupRule struct {
    ID         string `json:"id"`
    Direction  string `json:"direction"`   // "ingress" or "egress"
    Protocol   string `json:"protocol"`    // "tcp", "udp", "icmp"
    PortMin    int    `json:"port_min"`    // Response always includes (0 if not applicable)
    PortMax    int    `json:"port_max"`    // Response always includes (0 if not applicable)
    RemoteCIDR string `json:"remote_cidr"` // IP range in CIDR notation
}

// SecurityGroupRuleCreateRequest represents a rule creation request.
type SecurityGroupRuleCreateRequest struct {
    Direction  string `json:"direction"`          // Required: "ingress" or "egress"
    Protocol   string `json:"protocol"`           // Required: "tcp", "udp", "icmp"
    PortMin    *int   `json:"port_min,omitempty"` // Optional: starting port (TCP/UDP only)
    PortMax    *int   `json:"port_max,omitempty"` // Optional: ending port (TCP/UDP only)
    RemoteCIDR string `json:"remote_cidr"`        // Required: IP range in CIDR notation
}
```

**Rationale**:

1. **Request Model** (Create):
   - `*int` pointers allow distinguishing "not set" (nil) from "port 0" (valid but uncommon)
   - `omitempty` JSON tag excludes nil fields from request body
   - Supports ICMP rules where ports are irrelevant (send nil values)

2. **Response Model** (Rule):
   - Plain `int` since API always returns port values (0 if not applicable)
   - Simplifies consumption (no nil checks needed)
   - Matches API behavior (always populated in responses)

### Current Implementation Verification

**Existing Code** (`models/vps/securitygroups/rule.go`):

```go
type SecurityGroupRule struct {
    ID         string `json:"id"`
    Direction  string `json:"direction"`
    Protocol   string `json:"protocol"`
    PortMin    *int   `json:"port_min,omitempty"`  // ❓ Pointer in response model
    PortMax    *int   `json:"port_max,omitempty"`  // ❓ Pointer in response model
    RemoteCIDR string `json:"remote_cidr"`
}

type SecurityGroupRuleCreateRequest struct {
    Direction  string `json:"direction"`
    Protocol   string `json:"protocol"`
    PortMin    *int   `json:"port_min,omitempty"`  // ✅ Correct (optional)
    PortMax    *int   `json:"port_max,omitempty"`  // ✅ Correct (optional)
    RemoteCIDR string `json:"remote_cidr"`
}
```

**Issue Found**: Response model (`SecurityGroupRule`) uses `*int` but API always returns integers.

**Correction Required**: Change response model to plain `int`:

```go
type SecurityGroupRule struct {
    ID         string `json:"id"`
    Direction  string `json:"direction"`
    Protocol   string `json:"protocol"`
    PortMin    int    `json:"port_min"`    // Changed from *int
    PortMax    int    `json:"port_max"`    // Changed from *int
    RemoteCIDR string `json:"remote_cidr"`
}
```

**Migration Impact**: This is a **breaking change** (pointer to value), include in same version bump as `Total` field removal.

---

## 5. Error Handling Patterns

### Pattern Analysis

**Reference Implementation**: Existing VPS modules use `sdkError` from `errors.go`

```go
// errors.go (root package)
type sdkError struct {
    StatusCode int
    Message    string
    Details    string
}

func (e *sdkError) Error() string {
    return fmt.Sprintf("SDK error (HTTP %d): %s - %s", e.StatusCode, e.Message, e.Details)
}
```

**Usage in Clients**:
```go
// internal/http/client.go handles HTTP errors and wraps them
func (c *Client) Do(ctx context.Context, req *Request, result interface{}) error {
    // ... HTTP request logic
    if resp.StatusCode >= 400 {
        return &sdkError{
            StatusCode: resp.StatusCode,
            Message:    "Request failed",
            Details:    bodyString,
        }
    }
    // ... parse response
}
```

### Decision: Reuse Existing Pattern

**Implementation**:

```go
// modules/vps/securitygroups/client.go
func (c *Client) Create(ctx context.Context, req *securitygroups.SecurityGroupCreateRequest) (*securitygroups.SecurityGroup, error) {
    path := fmt.Sprintf("/api/v1/project/%s/security_groups", c.projectID)
    
    httpReq := &internalhttp.Request{
        Method: "POST",
        Path:   path,
        Body:   req,
    }
    
    var sg securitygroups.SecurityGroup
    if err := c.baseClient.Do(ctx, httpReq, &sg); err != nil {
        return nil, fmt.Errorf("failed to create security group: %w", err)  // Wrap error
    }
    
    return &sg, nil
}
```

**Error Wrapping Strategy**:
- Let `internal/http` client handle HTTP status codes and create `sdkError`
- Wrap errors with context using `fmt.Errorf` and `%w` verb
- Preserve error chain for debugging (supports `errors.Unwrap`)

**Context Cancellation**:
- `internal/http` client already respects `ctx.Done()` via `req.WithContext(ctx)`
- No additional handling needed in service-level clients

---

## 6. Test Strategy

### Test Organization

**Structure** (following established pattern):

```
modules/vps/securitygroups/
├── client.go                # CRUD operations
├── client_test.go           # Unit tests (package securitygroups)
├── rules.go                 # Rule sub-resource operations
├── rules_test.go            # Unit tests for rules (package securitygroups)
└── test/
    ├── fixtures.go          # Shared test fixtures
    ├── contract_test.go     # Contract tests vs API spec (package securitygroups_test)
    └── integration_test.go  # Full lifecycle tests (package securitygroups_test)
```

### Unit Tests (Same Package)

**Purpose**: Verify implementation logic, contribute to package coverage

**Approach**:
- Table-driven tests for each operation
- Mock `internal/http.Client` using custom test transport
- Test error paths (400, 404, 500 responses)
- Validate request body/query parameter construction

**Example** (`client_test.go`):
```go
func TestClient_List(t *testing.T) {
    tests := []struct {
        name     string
        opts     *securitygroups.ListSecurityGroupsOptions
        wantPath string
        wantErr  bool
    }{
        {
            name:     "list all",
            opts:     nil,
            wantPath: "/api/v1/project/proj-123/security_groups",
        },
        {
            name: "filter by name",
            opts: &securitygroups.ListSecurityGroupsOptions{
                Name: strPtr("web-sg"),
            },
            wantPath: "/api/v1/project/proj-123/security_groups?name=web-sg",
        },
        // ... more cases
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Mock HTTP transport, verify request, return fixture response
            // ...
        })
    }
}
```

### Contract Tests (Separate Package)

**Purpose**: Validate API behavior against Swagger specification

**Approach**:
- Use `httptest.Server` to mock API responses
- Load fixture JSON from Swagger examples or manually crafted samples
- Verify request/response marshaling matches API specification
- One test per endpoint

**Example** (`test/contract_test.go`):
```go
package securitygroups_test

func TestSecurityGroupsContract_Create(t *testing.T) {
    // Mock server with Swagger-compliant response
    server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if r.Method != "POST" {
            t.Errorf("expected POST, got %s", r.Method)
        }
        w.Header().Set("Content-Type", "application/json")
        w.WriteHeader(http.StatusCreated)
        fmt.Fprint(w, `{"id":"sg-123","name":"test-sg","description":"Test security group",` +
            `"project_id":"proj-123","user_id":"user-456","namespace":"default",` +
            `"rules":[],"createdAt":"2025-11-09T10:00:00Z","updatedAt":"2025-11-09T10:00:00Z"}`)
    }))
    defer server.Close()
    
    // Test against mock server
    // ...
}
```

### Integration Tests (Optional)

**Purpose**: Test full lifecycle against real or staging API

**Approach**:
- Use build tags (`//go:build integration`) to exclude from regular runs
- Require real API credentials via environment variables
- Create → Get → Update → Delete flow
- Clean up resources in defer/cleanup functions

**Target**: 80%+ coverage (unit tests) + complete contract test coverage (all Swagger endpoints)

---

## 7. Dependencies

### Summary

**Zero External Dependencies** - All requirements met by:

1. **Go Standard Library**:
   - `net/http`: HTTP client
   - `encoding/json`: JSON marshaling
   - `context`: Request cancellation and timeout
   - `time`: Timestamp handling
   - `net/url`: Query parameter encoding
   - `fmt`: String formatting and error wrapping
   - `testing`: Test framework

2. **Existing Internal Packages**:
   - `internal/http`: Shared HTTP client with retry/backoff
   - `internal/backoff`: Exponential backoff for retries
   - `internal/waiter`: Polling for async operations (not needed for security groups)
   - `internal/types`: Shared error types

3. **Existing Model Packages**:
   - `models/vps/securitygroups`: Shared data structures

**No new dependencies required** - consistent with constitution principle of minimal dependencies.

---

## Conclusion

All research questions resolved:

1. ✅ **Sub-resource pattern**: Use Resource Wrapper with embedded model + interface methods (consistent with `server.Volumes()`, `router.Networks()`)
2. ✅ **Query parameters**: Use `url.Values` for cleaner, safer query string construction
3. ✅ **Breaking changes**: Remove `Total` field + change `PortMin`/`PortMax` to non-pointers in response model; MINOR version bump (pre-1.0) or MAJOR (1.0+)
4. ✅ **Port types**: Pointers in request model (`*int`), plain integers in response model (`int`)
5. ✅ **Error handling**: Reuse existing `sdkError` type via `internal/http` client, wrap with context using `fmt.Errorf`
6. ✅ **Testing**: Unit tests (same package) + contract tests (test/ subdirectory) targeting 80%+ coverage
7. ✅ **Dependencies**: Zero external dependencies, stdlib + internal packages only

**Next Step**: Proceed to Phase 1 (Design) to create `data-model.md`, `contracts/`, and `quickstart.md`.
