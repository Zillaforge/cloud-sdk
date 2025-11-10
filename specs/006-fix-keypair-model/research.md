# Research: Fix Keypair Model

**Feature**: 006-fix-keypair-model  
**Date**: 2025-11-10  
**Phase**: 0 - Outline & Research

## Research Questions

### 1. JSON Tag Patterns for API Compatibility

**Task**: Research correct JSON tag patterns for pb.KeypairInfo fields

**Decision**: Use snake_case for JSON tags matching Swagger spec exactly

**Rationale**: The Swagger spec pb.KeypairInfo uses camelCase for some fields (createdAt, updatedAt) but snake_case for others (public_key, private_key, user_id). Go SDK must match these exactly for proper serialization.

**Pattern Identified**:
```go
CreatedAt   string  `json:"createdAt"`    // camelCase
UpdatedAt   string  `json:"updatedAt"`    // camelCase  
PublicKey   string  `json:"public_key"`   // snake_case
PrivateKey  string  `json:"private_key,omitempty"` // snake_case + optional
UserID      string  `json:"user_id"`      // snake_case
```

**Source**: `/workspaces/cloud-sdk/swagger/vps.yaml` lines 915-937

### 2. User Reference Type Pattern

**Task**: Research how to represent pb.IDName reference in Go models

**Decision**: Create reusable `IDName` type in `internal/types/types.go`

**Rationale**: The pb.IDName pattern is used across multiple resources (servers, images, etc.). Creating a shared type promotes consistency and reduces duplication.

**Implementation**:
```go
// IDName represents a reference to another resource with ID and name
type IDName struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

**Alternatives Considered**:
- Inline struct: Rejected - causes duplication across models
- Pointer to external package: Rejected - adds unnecessary dependency

### 3. List Response Pattern

**Task**: Research SDK pattern for list operations across the codebase

**Decision**: Return `[]*Resource` slice directly, not wrapped in response struct

**Rationale**: Examining `modules/vps/flavors/client.go` shows the pattern:
```go
func (c *Client) List(ctx context.Context, opts *flavors.ListFlavorsOptions) ([]*flavors.Flavor, error) {
    var response flavors.FlavorListResponse
    if err := c.baseClient.Do(ctx, req, &response); err != nil {
        return nil, fmt.Errorf("failed to list flavors: %w", err)
    }
    return response.Flavors, nil  // Return slice directly
}
```

This allows callers to:
- Access items by index: `keypairs[0]`
- Get count: `len(keypairs)`
- Iterate: `for _, kp := range keypairs`

**Source**: `/workspaces/cloud-sdk/modules/vps/flavors/client.go` lines 32-67

### 4. Error Handling Pattern

**Task**: Research consistent error wrapping pattern for API calls

**Decision**: Use `fmt.Errorf` with `%w` verb for error wrapping

**Pattern**:
```go
if err := c.baseClient.Do(ctx, req, &response); err != nil {
    return nil, fmt.Errorf("failed to <operation> <resource>: %w", err)
}
```

**Rationale**: This pattern:
- Preserves error chain for `errors.Is` / `errors.As`
- Provides context about which operation failed
- Consistent across all modules in the SDK

**Source**: Multiple files in `/workspaces/cloud-sdk/modules/vps/*/client.go`

### 5. Timestamp Handling

**Task**: Research timestamp representation in Go models

**Decision**: Use `string` type with ISO 8601 / RFC3339 format

**Rationale**: 
- Swagger spec defines timestamps as `type: string`
- Go's `encoding/json` handles RFC3339 strings natively
- Callers can parse with `time.Parse(time.RFC3339, timestamp)` if needed
- Avoids custom UnmarshalJSON complexity

**Clarification Applied**: From spec session 2025-11-10 - ISO 8601 / RFC3339 format

**Example**:
```go
CreatedAt string `json:"createdAt"` // "2025-11-10T15:30:00Z"
UpdatedAt string `json:"updatedAt"` // "2025-11-10T15:30:00Z"
```

### 6. Naming Convention: Request/Response vs Input/Output

**Task**: Research SDK naming conventions for API parameters

**Decision**: Use Request/Response suffixes instead of Input/Output

**Rationale**:
- **Request** - Clear indication of data being sent to API
- **Response** - Clear indication of data received from API
- More intuitive for SDK users
- Consistent with HTTP terminology
- Avoids confusion with functional programming terminology

**Pattern Applied**:
```go
// Input parameters use Request suffix
type KeypairCreateRequest struct { ... }
type KeypairUpdateRequest struct { ... }

// Output structures use Response suffix
type KeypairListResponse struct { ... }

// Main entities have no suffix
type Keypair struct { ... }
```

**Swagger Mapping**:
- `KeypairCreateInput` (Swagger) → `KeypairCreateRequest` (SDK)
- `KeypairUpdateInput` (Swagger) → `KeypairUpdateRequest` (SDK)
- `pb.KeypairListOutput` (Swagger) → `KeypairListResponse` (SDK)

### 7. Optional Field Handling

**Task**: Research omitempty tag usage for optional fields

**Decision**: Apply `omitempty` to fields that may be absent in API responses

**Fields Requiring omitempty**:
- `description` - optional in all operations
- `private_key` - only present in Create response
- `user` - may be null when not populated by API

**Pattern**:
```go
Description string   `json:"description,omitempty"`
PrivateKey  string   `json:"private_key,omitempty"`
User        *IDName  `json:"user,omitempty"`  // Pointer for nullable
```

**Rationale**: Prevents empty strings/null values in marshaled JSON, matches API behavior

## SDK Naming Conventions

### Consistent Naming Across SDK

All API-related structs follow these conventions:

1. **Request Types** (Input to API):
   - Suffix: `Request`
   - Examples: `KeypairCreateRequest`, `KeypairUpdateRequest`
   - Purpose: Parameters sent to API operations

2. **Response Types** (Output from API):
   - Suffix: `Response`
   - Examples: `KeypairListResponse`
   - Purpose: Data structures returned by API
   - Note: Often used internally when List() returns slice directly

3. **Entity Types** (Core Resources):
   - No suffix
   - Examples: `Keypair`, `Server`, `Network`
   - Purpose: Primary resource representations

4. **Options Types** (Query Parameters):
   - Suffix: `Options`
   - Examples: `ListKeypairsOptions`
   - Purpose: Optional filters for list operations

## Best Practices Applied

### Go Model Design
1. **Field Naming**: PascalCase for exported fields
2. **JSON Tags**: Match Swagger spec exactly (mix of camelCase and snake_case)
3. **Optional Fields**: Use `omitempty` tag
4. **Nullable Objects**: Use pointers (`*IDName`)
5. **Documentation**: Add field-level comments for sensitive data

### Testing Strategy
1. **Unit Tests**: JSON marshal/unmarshal with sample data
2. **Contract Tests**: Validate against actual Swagger examples
3. **Edge Cases**: Test with missing optional fields, null values
4. **Breaking Changes**: Test migration path from old to new model

### Migration Strategy
1. **Version Bump**: MAJOR (breaking change)
2. **Migration Guide**: Document in CHANGES.md
3. **Example**: Show before/after code for Total field removal
```go
// Before
total := response.Total

// After
total := len(response)  // Direct slice access
```

## References

- Swagger Spec: `/workspaces/cloud-sdk/swagger/vps.yaml`
- Similar Implementation: `/workspaces/cloud-sdk/models/vps/flavors/flavor.go`
- Client Pattern: `/workspaces/cloud-sdk/modules/vps/flavors/client.go`
- Constitution: `/workspaces/cloud-sdk/.specify/memory/constitution.md`
