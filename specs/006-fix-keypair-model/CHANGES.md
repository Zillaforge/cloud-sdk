# Changes: Fix Keypair Model

**Feature**: 006-fix-keypair-model  
**Version**: v1.0.0 (BREAKING CHANGE)  
**Date**: 2025-11-10

## Breaking Changes

### 1. KeypairListResponse.Total Field Removed

**Reason**: Field not present in Swagger specification `pb.KeypairListOutput`

**Migration**:

Before (v0.x):
```go
response, err := client.Keypairs.List(ctx, nil)
if err != nil {
    return err
}

total := response.Total
keypairs := response.Keypairs
```

After (v1.0):
```go
keypairs, err := client.Keypairs.List(ctx, nil)
if err != nil {
    return err
}

total := len(keypairs)  // Calculate from slice length
```

### 2. List() Return Type Changed

**Reason**: Direct slice access improves ergonomics

**Before**: `func List(...) (*KeypairListResponse, error)`  
**After**: `func List(...) ([]*Keypair, error)`

**Benefits**:
- Direct index access: `keypairs[0]`
- Built-in length: `len(keypairs)`
- Simpler iteration: `for _, kp := range keypairs`

## New Fields

### Keypair Model

Added fields to match `pb.KeypairInfo`:

- `PrivateKey string` - SSH private key (only in Create response)
- `User *IDName` - User reference object (may be null)
- `CreatedAt string` - ISO 8601 creation timestamp
- `UpdatedAt string` - ISO 8601 update timestamp

## New Types

### IDName

Reusable type for resource references:

```go
type IDName struct {
    ID   string `json:"id"`
    Name string `json:"name"`
}
```

Used by:
- `Keypair.User`
- Future: Server.Image, Server.Flavor, etc.

## Field Changes

### Optional Field Handling

Updated fields now use `omitempty` correctly:

- `Description` - Optional in all operations
- `PrivateKey` - Only present in Create response
- `User` - May be null even when UserID present

## Semantic Versioning

This release increments the **MAJOR** version due to breaking API changes:

- v0.x → v1.0.0

## Upgrade Path

### Step 1: Update Import

```bash
go get github.com/Zillaforge/cloud-sdk@v1.0.0
```

### Step 2: Update List Calls

Search your codebase for:
```go
response, err := *.Keypairs.List(
```

Replace with:
```go
keypairs, err := *.Keypairs.List(
```

And replace:
```go
response.Total
```

With:
```go
len(keypairs)
```

### Step 3: Handle New Fields

If you marshal/unmarshal keypairs to JSON/database, update your schemas to include:
- `private_key` (optional)
- `user` (optional nested object)
- `createdAt` (required string)
- `updatedAt` (required string)

### Step 4: Test

Run your tests to verify:
```bash
go test ./...
```

## Deprecation Notices

None. All changes are immediate removals due to spec alignment.

## Migration Script Example

```go
// migrate_keypair_code.go
package main

import (
    "bytes"
    "go/ast"
    "go/parser"
    "go/printer"
    "go/token"
    "log"
    "os"
)

func migrateFile(filename string) error {
    fset := token.NewFileSet()
    node, err := parser.ParseFile(fset, filename, nil, parser.ParseComments)
    if err != nil {
        return err
    }
    
    // AST transformation logic
    ast.Inspect(node, func(n ast.Node) bool {
        // Find assignments like: response, err := ...List(
        // Transform to: keypairs, err := ...List(
        
        // Find: response.Total
        // Transform to: len(response)
        
        return true
    })
    
    // Write back
    var buf bytes.Buffer
    printer.Fprint(&buf, fset, node)
    return os.WriteFile(filename, buf.Bytes(), 0644)
}
```

## Testing Recommendations

After upgrading:

1. **Unit Tests**: Verify JSON marshaling/unmarshaling
2. **Integration Tests**: Test against real API
3. **Contract Tests**: Validate Swagger compliance
4. **Migration Tests**: Run old test suite with new SDK

## Support

Questions or issues? 
- File an issue: https://github.com/Zillaforge/cloud-sdk/issues
- See: [quickstart.md](./quickstart.md) for examples
