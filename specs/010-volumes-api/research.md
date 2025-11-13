# Phase 0 Research: Volumes API Implementation

**Date**: 2025-11-13  
**Feature**: Volumes API Client (`008-volumes-api`)

## Research Task 1: Volume API Response Format Analysis

**Question**: How does the Volumes List API JSON response map to `[]*Volume` array?

**Investigation**: Examined `swagger/vps.yaml` lines 1672-1678 for `pb.VolumeListOutput` definition.

**Decision**: Use intermediate response wrapper struct `VolumeListResponse` with `Volumes` field.

**Finding**:
```yaml
pb.VolumeListOutput:
  properties:
    volumes:
      items:
        $ref: '#/definitions/pegasus-cloud_com_aes_virtualplatformserviceclient_pb.VolumeInfo'
      type: array
  type: object
```

API returns JSON:
```json
{
  "volumes": [ {...}, {...} ]
}
```

**Rationale**: 
- API wraps volume array in object with `volumes` field
- Cannot unmarshal directly to `[]*Volume`
- Must unmarshal to wrapper struct first, then extract array
- Consistent with existing patterns (FlavorListResponse, KeypairListResponse, etc.)

**Implementation Pattern**:
```go
type VolumeListResponse struct {
    Volumes []*Volume `json:"volumes"`
}

// In client:
var response VolumeListResponse
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to list volumes: %w", err)
}
return response.Volumes, nil
```

**Alternatives Considered**:
- Direct unmarsh al to `[]*Volume` - Rejected: API doesn't return bare array
- Custom UnmarshalJSON - Rejected: Unnecessary complexity when wrapper struct is simpler

**Code References**:
- `models/vps/flavors/flavor.go` lines 95-97: `FlavorListResponse` pattern
- `swagger/vps.yaml` lines 1672-1678: `pb.VolumeListOutput` definition

---

## Research Task 2: VolumeTypes API Response Format Analysis

**Question**: How does the VolumeTypes List API JSON response map to `[]string` array?

**Investigation**: Examined `swagger/vps.yaml` lines 1679-1686 for `pb.VolumeTypeListOutput` definition.

**Decision**: Use intermediate response wrapper struct `VolumeTypeListResponse` with `VolumeTypes` field.

**Finding**:
```yaml
pb.VolumeTypeListOutput:
  properties:
    volume_types:
      items:
        type: string
      type: array
  type: object
```

API returns JSON:
```json
{
  "volume_types": ["SSD", "HDD", "NVMe"]
}
```

**Rationale**:
- API wraps string array in object with `volume_types` field
- Cannot unmarshal directly to `[]string`
- Must unmarshal to wrapper struct first, then extract array
- Simpler than custom JSON unmarshaling

**Implementation Pattern**:
```go
type VolumeTypeListResponse struct {
    VolumeTypes []string `json:"volume_types"`
}

// In client:
var response VolumeTypeListResponse
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to list volume types: %w", err)
}
return response.VolumeTypes, nil
```

**Alternatives Considered**:
- Direct unmarshal to `[]string` - Rejected: API doesn't return bare array
- Alias type with custom UnmarshalJSON - Rejected: Wrapper struct is cleaner

**Code References**:
- `swagger/vps.yaml` lines 1679-1686: `pb.VolumeTypeListOutput` definition
- Similar pattern used in flavors and other list responses

---

## Research Task 3: Existing VPS Client Integration Pattern

**Question**: How to integrate new Volumes/VolumeTypes clients into existing VPS client?

**Investigation**: Reviewed `modules/vps/client.go` for accessor method patterns.

**Decision**: Add `Volumes()` and `VolumeTypes()` accessor methods returning respective clients.

**Finding**: Existing pattern in `modules/vps/client.go`:
```go
// Client provides access to VPS service operations.
type Client struct {
    baseClient *internalhttp.Client
    projectID  string
    baseURL    string
}

// Flavors returns a client for flavor operations.
func (c *Client) Flavors() *flavors.Client {
    return flavors.NewClient(c.baseClient, c.projectID)
}

// Keypairs returns a client for keypair operations.
func (c *Client) Keypairs() *keypairs.Client {
    return keypairs.NewClient(c.baseClient, c.projectID)
}
```

**Rationale**:
- Consistent with existing resource accessor pattern
- Lazy initialization (create client on-demand)
- Encapsulates baseClient and projectID
- Each resource client is independent

**Implementation Pattern**:
```go
// Add to modules/vps/client.go:

// Volumes returns a client for volume operations.
func (c *Client) Volumes() *volumes.Client {
    return volumes.NewClient(c.baseClient, c.projectID)
}

// VolumeTypes returns a client for volume type operations.
func (c *Client) VolumeTypes() *volumetypes.Client {
    return volumetypes.NewClient(c.baseClient, c.projectID)
}
```

**Alternatives Considered**:
- Pre-initialize all clients in VPS client constructor - Rejected: Wastes memory, breaks lazy-load pattern
- Single volumes client with volumetypes as submethod - Rejected: Breaks separation of concerns

**Code References**:
- `modules/vps/client.go`: Existing accessor methods (Flavors, Keypairs, Networks, etc.)
- `modules/vps/flavors/client.go`: Example client constructor `NewClient(baseClient, projectID)`

---

## Research Task 4: Error Handling Consistency

**Question**: What error wrapping patterns are used in existing VPS modules?

**Investigation**: Reviewed error handling in `modules/vps/flavors/client.go` and `modules/vps/keypairs/client.go`.

**Decision**: Use standard `fmt.Errorf("failed to {action} {resource}: %w", err)` pattern for all errors.

**Finding**: Consistent pattern across VPS modules:
```go
// From flavors/client.go:
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to list flavors: %w", err)
}

if err := c.baseClient.Do(ctx, req, &flavor); err != nil {
    return nil, fmt.Errorf("failed to get flavor %s: %w", flavorID, err)
}

// From keypairs/client.go:
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to list keypairs: %w", err)
}

if err := c.baseClient.Do(ctx, req, nil); err != nil {
    return fmt.Errorf("failed to delete keypair %s: %w", keypairID, err)
}
```

**Rationale**:
- Adds context about what operation failed
- Preserves original error for `errors.Is()` and `errors.As()`
- Includes resource ID when applicable
- Consistent error message format across SDK
- Uses `%w` verb for error wrapping (Go 1.13+)

**Implementation Pattern**:
```go
// For list operations:
return nil, fmt.Errorf("failed to list volumes: %w", err)
return nil, fmt.Errorf("failed to list volume types: %w", err)

// For operations with ID:
return nil, fmt.Errorf("failed to create volume: %w", err)
return nil, fmt.Errorf("failed to get volume %s: %w", volumeID, err)
return nil, fmt.Errorf("failed to update volume %s: %w", volumeID, err)
return fmt.Errorf("failed to delete volume %s: %w", volumeID, err)
return fmt.Errorf("failed to perform action on volume %s: %w", volumeID, err)
```

**Alternatives Considered**:
- Custom error types - Rejected: internal/types.SDKError already provides structured errors from HTTP client
- Error codes - Rejected: HTTP status codes already available in SDKError
- No wrapping - Rejected: Loses operation context for debugging

**Code References**:
- `modules/vps/flavors/client.go` lines 56-58, 79-81: Error wrapping examples
- `modules/vps/keypairs/client.go`: Consistent error patterns
- `internal/types/types.go`: SDKError structure with StatusCode, Message

---

## Research Task 5: Contract Test Framework

**Question**: How are Swagger contract tests structured in the project?

**Investigation**: Searched for existing contract test examples.

**Decision**: Use table-driven tests with Swagger schema validation for each endpoint.

**Finding**: Based on Go testing best practices and SDK patterns:

Contract tests should:
1. Validate request/response structure against Swagger definitions
2. Test all documented status codes (200, 201, 204, 400, 500)
3. Verify required/optional fields
4. Check data types and constraints
5. Test query parameters and request bodies

**Rationale**:
- Ensures SDK matches API contract
- Catches breaking changes in API
- Documents expected behavior
- Validates JSON serialization/deserialization
- Tests error response structures

**Implementation Pattern**:
```go
// tests/contract/volumes_test.go
package contract_test

import (
    "context"
    "testing"
    "github.com/Zillaforge/cloud-sdk/models/vps/volumes"
)

func TestVolumesContract(t *testing.T) {
    tests := []struct {
        name     string
        method   string
        path     string
        reqBody  interface{}
        status   int
        respType interface{}
    }{
        {
            name:     "ListVolumes",
            method:   "GET",
            path:     "/api/v1/project/{project-id}/volumes",
            status:   200,
            respType: &volumes.VolumeListResponse{},
        },
        {
            name:     "CreateVolume",
            method:   "POST",
            path:     "/api/v1/project/{project-id}/volumes",
            reqBody:  &volumes.CreateInput{Name: "test", Type: "SSD"},
            status:   201,
            respType: &volumes.Volume{},
        },
        // ... more test cases
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
            // Validate response structure
            // Check required fields
            // Verify types
        })
    }
}
```

**Alternatives Considered**:
- OpenAPI validator library - Rejected: Adds external dependency, prefer native Go tests
- Manual JSON comparison - Rejected: Brittle, hard to maintain
- Mock server with recorded responses - Rejected: Doesn't test real API behavior

**Code References**:
- Go testing package documentation: https://pkg.go.dev/testing
- Existing unit tests in project use table-driven test pattern
- `swagger/vps.yaml`: Source of truth for API contracts

---

## Summary of Research Decisions

| Topic | Decision | Rationale |
|-------|----------|-----------|
| Volume List Response | Use `VolumeListResponse` wrapper struct | API wraps array in object |
| VolumeTypes List Response | Use `VolumeTypeListResponse` wrapper struct | API wraps array in object |
| VPS Client Integration | Add `Volumes()` and `VolumeTypes()` accessors | Consistent with existing pattern |
| Error Handling | Use `fmt.Errorf("failed to {action}: %w", err)` | Adds context, preserves error chain |
| Contract Tests | Table-driven tests with schema validation | Maintainable, comprehensive coverage |

All research tasks complete. No NEEDS CLARIFICATION items remaining. Ready for Phase 1 design.

---

## Implementation Patterns Codified

### 1. Response Unmarshaling Pattern
```go
// Step 1: Define wrapper struct
type ResourceListResponse struct {
    Resources []*Resource `json:"resources"`
}

// Step 2: Unmarshal to wrapper
var response ResourceListResponse
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to list resources: %w", err)
}

// Step 3: Return unwrapped array
return response.Resources, nil
```

### 2. Client Constructor Pattern
```go
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
    basePath := "/api/v1/project/" + projectID
    return &Client{
        baseClient: baseClient,
        projectID:  projectID,
        basePath:   basePath,
    }
}
```

### 3. Error Wrapping Pattern
```go
// Generic operations
return nil, fmt.Errorf("failed to {verb} {resource}: %w", err)

// Operations with resource ID
return nil, fmt.Errorf("failed to {verb} {resource} %s: %w", resourceID, err)
```

### 4. Client Method Pattern
```go
func (c *Client) Operation(ctx context.Context, params *Params) (*Result, error) {
    path := c.basePath + "/endpoint"
    
    req := &internalhttp.Request{
        Method: "GET",
        Path:   path,
    }
    
    var result Result
    if err := c.baseClient.Do(ctx, req, &result); err != nil {
        return nil, fmt.Errorf("failed to operation: %w", err)
    }
    
    return &result, nil
}
```

---

## Technical Debt & Future Improvements

None identified. All patterns follow existing SDK conventions.

---

## References

- Swagger Specification: `swagger/vps.yaml`
- Existing VPS Modules: `modules/vps/flavors/`, `modules/vps/keypairs/`
- Internal HTTP Client: `internal/http/client.go`
- Common Types: `models/vps/common/common.go`
- Go Error Handling: https://go.dev/blog/go1.13-errors
