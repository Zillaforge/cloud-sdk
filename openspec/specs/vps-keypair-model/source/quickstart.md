# Quickstart: Keypair Operations

**Feature**: 006-fix-keypair-model  
**Date**: 2025-11-10  
**Phase**: 1 - Design & Contracts

This guide demonstrates using the corrected keypair model with the VPS SDK.

## Setup

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/Zillaforge/cloud-sdk"
    "github.com/Zillaforge/cloud-sdk/modules/vps/keypairs"
)

func main() {
    // Initialize SDK client
    client, err := cloudsdk.NewClient(
        cloudsdk.WithToken("your-api-token"),
        cloudsdk.WithBaseURL("https://api.example.com"),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create keypair client for your project
    kpClient := client.VPS.Keypairs("project-id-123")
    
    ctx := context.Background()
    
    // Examples below...
}
```

## Example 1: Create New Keypair (Generate)

Generate a new SSH keypair on the server:

```go
// Create request
req := &keypairs.KeypairCreateRequest{
    Name:        "my-dev-keypair",
    Description: "Development environment key",
    // PublicKey omitted - server generates new pair
}

// Execute create
keypair, err := kpClient.Create(ctx, req)
if err != nil {
    log.Fatalf("Failed to create keypair: %v", err)
}

// ⚠️ IMPORTANT: Save private key immediately!
// Private key is ONLY returned during creation and cannot be retrieved later
if keypair.PrivateKey != "" {
    fmt.Println("🔑 Private Key (SAVE THIS NOW):")
    fmt.Println(keypair.PrivateKey)
    
    // Example: Save to file
    err = os.WriteFile("my-dev-keypair.pem", []byte(keypair.PrivateKey), 0600)
    if err != nil {
        log.Fatalf("Failed to save private key: %v", err)
    }
}

fmt.Printf("✅ Created keypair: %s (ID: %s)\n", keypair.Name, keypair.ID)
fmt.Printf("   Fingerprint: %s\n", keypair.Fingerprint)
fmt.Printf("   Created: %s\n", keypair.CreatedAt)
```

## Example 2: Import Existing Public Key

Import an existing SSH public key:

```go
// Read your existing public key
publicKey, err := os.ReadFile("~/.ssh/id_rsa.pub")
if err != nil {
    log.Fatalf("Failed to read public key: %v", err)
}

// Create request with public key
req := &keypairs.KeypairCreateRequest{
    Name:        "my-imported-key",
    Description: "Imported from local machine",
    PublicKey:   string(publicKey),
}

// Execute create
keypair, err := kpClient.Create(ctx, req)
if err != nil {
    log.Fatalf("Failed to import keypair: %v", err)
}

// Note: No PrivateKey in response when importing
fmt.Printf("✅ Imported keypair: %s\n", keypair.Name)
fmt.Printf("   Fingerprint: %s\n", keypair.Fingerprint)
```

## Example 3: List All Keypairs

List all keypairs in your project:

```go
// List all keypairs (no filters)
keypairs, err := kpClient.List(ctx, nil)
if err != nil {
    log.Fatalf("Failed to list keypairs: %v", err)
}

// ✨ NEW: Direct slice access (no .Total field)
fmt.Printf("Found %d keypairs:\n", len(keypairs))

for i, kp := range keypairs {
    fmt.Printf("\n%d. %s (ID: %s)\n", i+1, kp.Name, kp.ID)
    fmt.Printf("   Fingerprint: %s\n", kp.Fingerprint)
    fmt.Printf("   Created: %s\n", kp.CreatedAt)
    
    // User reference may be null
    if kp.User != nil {
        fmt.Printf("   Owner: %s (%s)\n", kp.User.Name, kp.User.ID)
    } else {
        fmt.Printf("   Owner ID: %s\n", kp.UserID)
    }
}
```

## Example 4: List with Filter

Filter keypairs by name:

```go
// List with name filter
opts := &keypairs.ListKeypairsOptions{
    Name: "my-dev-keypair",
}

keypairs, err := kpClient.List(ctx, opts)
if err != nil {
    log.Fatalf("Failed to list keypairs: %v", err)
}

if len(keypairs) == 0 {
    fmt.Println("No keypairs found matching filter")
    return
}

kp := keypairs[0]  // Direct index access
fmt.Printf("Found: %s\n", kp.Name)
```

## Example 5: Get Specific Keypair

Retrieve details of a specific keypair:

```go
keypairID := "kp-abc123"

keypair, err := kpClient.Get(ctx, keypairID)
if err != nil {
    log.Fatalf("Failed to get keypair: %v", err)
}

fmt.Printf("Keypair: %s\n", keypair.Name)
fmt.Printf("Description: %s\n", keypair.Description)
fmt.Printf("Fingerprint: %s\n", keypair.Fingerprint)
fmt.Printf("Public Key:\n%s\n", keypair.PublicKey)

// Note: PrivateKey NOT returned in Get operation
// It's only available during Create (when generated)
```

## Example 6: Update Keypair Description

Only the description field can be updated:

```go
keypairID := "kp-abc123"

req := &keypairs.KeypairUpdateRequest{
    Description: "Updated description for production use",
}

keypair, err := kpClient.Update(ctx, keypairID, req)
if err != nil {
    log.Fatalf("Failed to update keypair: %v", err)
}

fmt.Printf("✅ Updated keypair: %s\n", keypair.Name)
fmt.Printf("   New description: %s\n", keypair.Description)
fmt.Printf("   Last updated: %s\n", keypair.UpdatedAt)
```

## Example 7: Delete Keypair

Delete a keypair:

```go
keypairID := "kp-abc123"

err := kpClient.Delete(ctx, keypairID)
if err != nil {
    log.Fatalf("Failed to delete keypair: %v", err)
}

fmt.Printf("✅ Deleted keypair: %s\n", keypairID)
```

## Example 8: Parse Timestamps

Work with timestamp fields:

```go
import "time"

keypair, err := kpClient.Get(ctx, "kp-abc123")
if err != nil {
    log.Fatal(err)
}

// Parse ISO 8601 timestamp
createdTime, err := time.Parse(time.RFC3339, keypair.CreatedAt)
if err != nil {
    log.Fatalf("Failed to parse timestamp: %v", err)
}

// Calculate age
age := time.Since(createdTime)
fmt.Printf("Keypair created %v ago\n", age.Round(time.Second))

// Format for display
fmt.Printf("Created: %s\n", createdTime.Format("2006-01-02 15:04:05"))
```

## Migration from v0.x

### Before (v0.x)
```go
// Old API returned wrapped response with Total field
response, err := kpClient.List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

total := response.Total  // ❌ Field no longer exists
for _, kp := range response.Keypairs {
    fmt.Println(kp.Name)
}
```

### After (v1.0)
```go
// New API returns slice directly
keypairs, err := kpClient.List(ctx, nil)
if err != nil {
    log.Fatal(err)
}

total := len(keypairs)  // ✅ Calculate from slice length
for _, kp := range keypairs {
    fmt.Println(kp.Name)
}

// Bonus: Direct index access
first := keypairs[0]  // ✅ Access by index
```

## Error Handling

All operations return descriptive errors:

```go
keypair, err := kpClient.Get(ctx, "invalid-id")
if err != nil {
    // Error includes context about what failed
    // Example: "failed to get keypair invalid-id: HTTP 404: Not Found"
    fmt.Printf("Error: %v\n", err)
    
    // Check for specific error types
    var sdkErr *types.SDKError
    if errors.As(err, &sdkErr) {
        fmt.Printf("Status Code: %d\n", sdkErr.StatusCode)
        fmt.Printf("Error Code: %d\n", sdkErr.ErrorCode)
    }
    
    return
}
```

## Security Best Practices

### 1. Save Private Keys Immediately

```go
// ✅ GOOD: Save immediately after creation
keypair, err := kpClient.Create(ctx, req)
if err != nil {
    return err
}

if keypair.PrivateKey != "" {
    // Save with restricted permissions
    err = os.WriteFile("key.pem", []byte(keypair.PrivateKey), 0600)
    if err != nil {
        return fmt.Errorf("failed to save private key: %w", err)
    }
}

// ❌ BAD: Trying to retrieve later (not possible)
// Private key is only returned once during creation
```

### 2. Protect Private Key Files

```go
// Set restrictive file permissions (owner read/write only)
err = os.Chmod("key.pem", 0600)

// Verify permissions
info, _ := os.Stat("key.pem")
mode := info.Mode()
if mode != 0600 {
    log.Warning("Private key file has insecure permissions")
}
```

### 3. Don't Log Sensitive Data

```go
// ❌ BAD: Logging private key
log.Printf("Created keypair with private key: %s", keypair.PrivateKey)

// ✅ GOOD: Log without sensitive data
log.Printf("Created keypair %s (ID: %s)", keypair.Name, keypair.ID)
```

## Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/Zillaforge/cloud-sdk"
)

func main() {
    // Initialize
    client, err := cloudsdk.NewClient(
        cloudsdk.WithToken(os.Getenv("API_TOKEN")),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    kpClient := client.VPS.Keypairs(os.Getenv("PROJECT_ID"))
    ctx := context.Background()
    
    // Create new keypair
    keypair, err := kpClient.Create(ctx, &keypairs.KeypairCreateRequest{
        Name:        "quickstart-key",
        Description: "Created by quickstart guide",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    // Save private key
    if keypair.PrivateKey != "" {
        err = os.WriteFile("quickstart-key.pem", []byte(keypair.PrivateKey), 0600)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Println("✅ Private key saved to quickstart-key.pem")
    }
    
    fmt.Printf("✅ Created keypair: %s\n", keypair.Name)
    fmt.Printf("   ID: %s\n", keypair.ID)
    fmt.Printf("   Fingerprint: %s\n", keypair.Fingerprint)
    
    // List all keypairs
    keypairs, err := kpClient.List(ctx, nil)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("\n📋 Total keypairs: %d\n", len(keypairs))
}
```

## Next Steps

- Review [data-model.md](./data-model.md) for complete field reference
- Check [contracts/keypair-api.md](./contracts/keypair-api.md) for API specifications
- See [CHANGES.md] for migration guide from v0.x to v1.0
