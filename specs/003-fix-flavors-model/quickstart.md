# Quickstart: Fix Flavors Model

**Date**: November 8, 2025
**Feature**: 001-fix-flavors-model

## Overview
The flavors model has been updated to match the VPS API specification. This includes new fields, corrected field names, and enhanced filtering capabilities.

## Breaking Changes
- Field `VCPUs` renamed to `VCPU`
- Field `RAM` renamed to `Memory`
- JSON tags updated to match API ("vcpu", "memory")

## Usage Examples

### Basic Usage
```go
client := cloudsdk.NewClient("https://api.example.com", "token")
vps := client.Project("project-123").VPS()
flavors := vps.Flavors()

// List all flavors
allFlavors, err := flavors.List(ctx, nil)

// Get specific flavor
flavor, err := flavors.Get(ctx, "flavor-123")
```

### Filtering Flavors
```go
opts := &flavors.ListFlavorsOptions{
    Name: "gpu",
    Public: &[]bool{true}[0], // pointer to true
    Tags: []string{"compute", "gpu"},
    ProjectID: "project-456",
}

filteredFlavors, err := flavors.List(ctx, opts)
```

### Accessing New Fields
```go
flavor := filteredFlavors.Items[0]

// New fields
if flavor.GPU != nil {
    fmt.Printf("GPU: %s x%d\n", flavor.GPU.Model, flavor.GPU.Count)
}

if flavor.CreatedAt != nil {
    fmt.Printf("Created: %s\n", flavor.CreatedAt.Format(time.RFC3339))
}
```

## Migration Guide
Update your code to use new field names:

```go
// Old code
flavor.VCPUs    // ❌ No longer available
flavor.RAM      // ❌ No longer available

// New code
flavor.VCPU     // ✅
flavor.Memory   // ✅
```

## Testing
Run tests to verify compatibility:
```bash
go test ./modules/vps/flavors/...
make check
```</content>
<parameter name="filePath">/workspaces/cloud-sdk/specs/001-fix-flavors-model/quickstart.md