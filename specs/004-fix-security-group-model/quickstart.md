# Quick Start: Security Groups API

**Feature**: Fix Security Group Model (001-fix-security-group-model)  
**Date**: November 9, 2025  
**Related**: [spec.md](./spec.md) | [plan.md](./plan.md) | [data-model.md](./data-model.md) | [contracts/](./contracts/)

## Overview

This guide demonstrates complete security group operations using the Cloud SDK, including CRUD operations, rule management via sub-resource pattern, and filtering. All examples follow the corrected model definitions matching the VPS API specification.

---

## Prerequisites

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    cloudsdk "github.com/Zillaforge/cloud-sdk"
    "github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)
```

**Environment Setup**:
```bash
export CLOUD_API_TOKEN="your-api-token"
export CLOUD_API_URL="https://api.example.com"
export CLOUD_PROJECT_ID="proj-123"
```

---

## Step 1: Initialize Client

```go
func main() {
    ctx := context.Background()

    // Create SDK client with base URL and API token
    baseURL := "https://api.example.com"  // Or use os.Getenv("CLOUD_API_URL")
    token := "your-api-token"             // Or use os.Getenv("CLOUD_API_TOKEN")
    
    client, err := cloudsdk.New(baseURL, token)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    // Select project and VPS service
    projectID := "proj-123"  // Or use os.Getenv("CLOUD_PROJECT_ID")
    vps := client.Project(projectID).VPS()

    // Access security groups operations
    sgClient := vps.SecurityGroups()

    // ... operations follow
}
```

---

## Step 2: Create Security Group (Without Rules)

```go
// Create a basic security group without initial rules
func createBasicSecurityGroup(ctx context.Context, sgClient securitygroups.Client) (*securitygroups.SecurityGroup, error) {
    req := &securitygroups.SecurityGroupCreateRequest{
        Name:        "web-servers",
        Description: "Security group for web servers",
        // Rules: nil (add rules later via sub-resource)
    }

    sg, err := sgClient.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create security group: %w", err)
    }

    fmt.Printf("Created security group: %s (ID: %s)\n", sg.Name, sg.ID)
    return sg, nil
}
```

**Output**:
```
Created security group: web-servers (ID: sg-f4e3d2c1)
```

---

## Step 3: Create Security Group with Initial Rules

```go
// Create a security group with pre-defined rules
func createSecurityGroupWithRules(ctx context.Context, sgClient securitygroups.Client) (*securitygroups.SecurityGroup, error) {
    // Helper function to create int pointers
    intPtr := func(v int) *int { return &v }

    req := &securitygroups.SecurityGroupCreateRequest{
        Name:        "database-servers",
        Description: "Security group for PostgreSQL databases",
        Rules: []securitygroups.SecurityGroupRuleCreateRequest{
            {
                Direction:  "ingress",
                Protocol:   "tcp",
                PortMin:    intPtr(5432),  // PostgreSQL port
                PortMax:    intPtr(5432),
                RemoteCIDR: "10.0.0.0/16", // Internal network only
            },
            {
                Direction:  "egress",
                Protocol:   "tcp",
                PortMin:    intPtr(443),   // HTTPS for updates
                PortMax:    intPtr(443),
                RemoteCIDR: "0.0.0.0/0",
            },
        },
    }

    sg, err := sgClient.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("failed to create security group: %w", err)
    }

    fmt.Printf("Created security group: %s with %d initial rules\n", sg.Name, len(sg.Rules))
    return sg, nil
}
```

**Output**:
```
Created security group: database-servers with 2 initial rules
```

---

## Step 4: List Security Groups (Without Rules)

```go
// List all security groups without rules (efficient for overview)
func listSecurityGroups(ctx context.Context, sgClient securitygroups.Client) error {
    // No options = default behavior (no rules included)
    resp, err := sgClient.List(ctx, nil)
    if err != nil {
        return fmt.Errorf("failed to list security groups: %w", err)
    }

    fmt.Printf("Found %d security groups:\n", len(resp.SecurityGroups))
    for _, sg := range resp.SecurityGroups {
        fmt.Printf("  - %s (ID: %s): %s\n", sg.Name, sg.ID, sg.Description)
        fmt.Printf("    Created: %s | Rules: %d (use detail=true to see rules)\n",
            sg.CreatedAt.Format(time.RFC3339), len(sg.Rules))
    }

    return nil
}
```

**Output**:
```
Found 3 security groups:
  - web-servers (ID: sg-001): Security group for web servers
    Created: 2025-11-09T10:00:00Z | Rules: 0 (use detail=true to see rules)
  - database-servers (ID: sg-002): Security group for PostgreSQL databases
    Created: 2025-11-09T09:00:00Z | Rules: 0 (use detail=true to see rules)
  - default (ID: sg-003): Default security group
    Created: 2025-11-01T08:00:00Z | Rules: 0 (use detail=true to see rules)
```

---

## Step 5: List Security Groups with Rules (Detailed)

```go
// List security groups with full rule details
func listSecurityGroupsDetailed(ctx context.Context, sgClient securitygroups.Client) error {
    // Set detail=true to include rules array
    detail := true
    opts := &securitygroups.ListSecurityGroupsOptions{
        Detail: &detail,
    }

    resp, err := sgClient.List(ctx, opts)
    if err != nil {
        return fmt.Errorf("failed to list security groups: %w", err)
    }

    fmt.Printf("Found %d security groups (with rules):\n", len(resp.SecurityGroups))
    for _, sg := range resp.SecurityGroups {
        fmt.Printf("\n%s (ID: %s):\n", sg.Name, sg.ID)
        fmt.Printf("  Description: %s\n", sg.Description)
        fmt.Printf("  Rules (%d):\n", len(sg.Rules))
        for _, rule := range sg.Rules {
            fmt.Printf("    - [%s] %s %s ports %d-%d from %s (Rule ID: %s)\n",
                rule.Direction, rule.Protocol,
                rule.Protocol, rule.PortMin, rule.PortMax,
                rule.RemoteCIDR, rule.ID)
        }
    }

    return nil
}
```

**Output**:
```
Found 2 security groups (with rules):

web-servers (ID: sg-001):
  Description: Security group for web servers
  Rules (2):
    - [ingress] tcp tcp ports 80-80 from 0.0.0.0/0 (Rule ID: rule-001)
    - [ingress] tcp tcp ports 443-443 from 0.0.0.0/0 (Rule ID: rule-002)

database-servers (ID: sg-002):
  Description: Security group for PostgreSQL databases
  Rules (2):
    - [ingress] tcp tcp ports 5432-5432 from 10.0.0.0/16 (Rule ID: rule-003)
    - [egress] tcp tcp ports 443-443 from 0.0.0.0/0 (Rule ID: rule-004)
```

---

## Step 6: Filter Security Groups by Name

```go
// Filter security groups by exact name match
func filterSecurityGroupsByName(ctx context.Context, sgClient securitygroups.Client, name string) error {
    opts := &securitygroups.ListSecurityGroupsOptions{
        Name: &name,
    }

    resp, err := sgClient.List(ctx, opts)
    if err != nil {
        return fmt.Errorf("failed to filter security groups: %w", err)
    }

    if len(resp.SecurityGroups) == 0 {
        fmt.Printf("No security groups found with name '%s'\n", name)
        return nil
    }

    for _, sg := range resp.SecurityGroups {
        fmt.Printf("Found: %s (ID: %s)\n", sg.Name, sg.ID)
    }

    return nil
}

// Usage
func main() {
    // ...
    err := filterSecurityGroupsByName(ctx, sgClient, "web-servers")
    // ...
}
```

**Output**:
```
Found: web-servers (ID: sg-001)
```

---

## Step 7: Get Single Security Group and Add Rules (Sub-Resource Pattern)

```go
// Get security group and add rules using sub-resource pattern
func addRulesToSecurityGroup(ctx context.Context, sgClient securitygroups.Client, sgID string) error {
    // Step 7a: Get security group (returns resource wrapper)
    sg, err := sgClient.Get(ctx, sgID)
    if err != nil {
        return fmt.Errorf("failed to get security group: %w", err)
    }

    fmt.Printf("Retrieved security group: %s\n", sg.Name)
    fmt.Printf("Current rules: %d\n", len(sg.Rules))

    // Step 7b: Access rule sub-resource operations
    rules := sg.Rules()

    // Step 7c: Add HTTP rule
    httpRule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{
        Direction:  "ingress",
        Protocol:   "tcp",
        PortMin:    intPtr(80),
        PortMax:    intPtr(80),
        RemoteCIDR: "0.0.0.0/0",
    })
    if err != nil {
        return fmt.Errorf("failed to create HTTP rule: %w", err)
    }
    fmt.Printf("Added HTTP rule (ID: %s)\n", httpRule.ID)

    // Step 7d: Add HTTPS rule
    httpsRule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{
        Direction:  "ingress",
        Protocol:   "tcp",
        PortMin:    intPtr(443),
        PortMax:    intPtr(443),
        RemoteCIDR: "0.0.0.0/0",
    })
    if err != nil {
        return fmt.Errorf("failed to create HTTPS rule: %w", err)
    }
    fmt.Printf("Added HTTPS rule (ID: %s)\n", httpsRule.ID)

    // Step 7e: Add ICMP rule (ping)
    icmpRule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{
        Direction:  "ingress",
        Protocol:   "icmp",
        RemoteCIDR: "0.0.0.0/0",
        // PortMin/PortMax omitted for ICMP
    })
    if err != nil {
        return fmt.Errorf("failed to create ICMP rule: %w", err)
    }
    fmt.Printf("Added ICMP rule (ID: %s)\n", icmpRule.ID)

    return nil
}

// Helper function
func intPtr(v int) *int {
    return &v
}
```

**Output**:
```
Retrieved security group: web-servers
Current rules: 0
Added HTTP rule (ID: rule-abc001)
Added HTTPS rule (ID: rule-abc002)
Added ICMP rule (ID: rule-abc003)
```

---

## Step 8: Delete Rule from Security Group

```go
// Delete a specific rule from a security group
func deleteRule(ctx context.Context, sgClient securitygroups.Client, sgID, ruleID string) error {
    // Get security group to access rule operations
    sg, err := sgClient.Get(ctx, sgID)
    if err != nil {
        return fmt.Errorf("failed to get security group: %w", err)
    }

    // Access rule sub-resource
    rules := sg.Rules()

    // Delete rule
    if err := rules.Delete(ctx, ruleID); err != nil {
        return fmt.Errorf("failed to delete rule: %w", err)
    }

    fmt.Printf("Deleted rule %s from security group %s\n", ruleID, sg.Name)
    return nil
}
```

**Output**:
```
Deleted rule rule-abc001 from security group web-servers
```

---

## Step 9: Update Security Group Metadata

```go
// Update security group name and description
func updateSecurityGroup(ctx context.Context, sgClient securitygroups.Client, sgID string) error {
    newName := "web-servers-v2"
    newDescription := "Updated security group for web servers (production)"

    req := &securitygroups.SecurityGroupUpdateRequest{
        Name:        &newName,
        Description: &newDescription,
    }

    sg, err := sgClient.Update(ctx, sgID, req)
    if err != nil {
        return fmt.Errorf("failed to update security group: %w", err)
    }

    fmt.Printf("Updated security group: %s\n", sg.Name)
    fmt.Printf("  New description: %s\n", sg.Description)
    fmt.Printf("  Updated at: %s\n", sg.UpdatedAt.Format(time.RFC3339))

    return nil
}
```

**Output**:
```
Updated security group: web-servers-v2
  New description: Updated security group for web servers (production)
  Updated at: 2025-11-09T11:30:00Z
```

---

## Step 10: Delete Security Group

```go
// Delete a security group and all its rules
func deleteSecurityGroup(ctx context.Context, sgClient securitygroups.Client, sgID string) error {
    if err := sgClient.Delete(ctx, sgID); err != nil {
        return fmt.Errorf("failed to delete security group: %w", err)
    }

    fmt.Printf("Deleted security group: %s\n", sgID)
    return nil
}
```

**Output**:
```
Deleted security group: sg-001
```

---

## Complete Example: Full Lifecycle

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"

    cloudsdk "github.com/Zillaforge/cloud-sdk"
    "github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

func main() {
    ctx := context.Background()

    // Initialize client
    baseURL := os.Getenv("CLOUD_API_URL")
    token := os.Getenv("CLOUD_API_TOKEN")
    
    client, err := cloudsdk.New(baseURL, token)
    if err != nil {
        log.Fatalf("Failed to create client: %v", err)
    }

    projectID := os.Getenv("CLOUD_PROJECT_ID")
    sgClient := client.Project(projectID).VPS().SecurityGroups()

    // Step 1: Create security group
    sg, err := sgClient.Create(ctx, &securitygroups.SecurityGroupCreateRequest{
        Name:        "demo-web-sg",
        Description: "Demo security group for web servers",
    })
    if err != nil {
        log.Fatalf("Create failed: %v", err)
    }
    fmt.Printf("✓ Created: %s (ID: %s)\n", sg.Name, sg.ID)

    // Step 2: Add rules using sub-resource pattern
    sgResource, err := sgClient.Get(ctx, sg.ID)
    if err != nil {
        log.Fatalf("Get failed: %v", err)
    }

    rules := sgResource.Rules()

    // Add HTTP rule
    httpRule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{
        Direction:  "ingress",
        Protocol:   "tcp",
        PortMin:    intPtr(80),
        PortMax:    intPtr(80),
        RemoteCIDR: "0.0.0.0/0",
    })
    if err != nil {
        log.Fatalf("Create HTTP rule failed: %v", err)
    }
    fmt.Printf("✓ Added HTTP rule (ID: %s)\n", httpRule.ID)

    // Add HTTPS rule
    httpsRule, err := rules.Create(ctx, &securitygroups.SecurityGroupRuleCreateRequest{
        Direction:  "ingress",
        Protocol:   "tcp",
        PortMin:    intPtr(443),
        PortMax:    intPtr(443),
        RemoteCIDR: "0.0.0.0/0",
    })
    if err != nil {
        log.Fatalf("Create HTTPS rule failed: %v", err)
    }
    fmt.Printf("✓ Added HTTPS rule (ID: %s)\n", httpsRule.ID)

    // Step 3: List security groups with rules
    detail := true
    resp, err := sgClient.List(ctx, &securitygroups.ListSecurityGroupsOptions{
        Detail: &detail,
    })
    if err != nil {
        log.Fatalf("List failed: %v", err)
    }
    fmt.Printf("✓ Listed %d security groups\n", len(resp.SecurityGroups))

    for _, sg := range resp.SecurityGroups {
        if sg.Name == "demo-web-sg" {
            fmt.Printf("  - %s: %d rules\n", sg.Name, len(sg.Rules))
        }
    }

    // Step 4: Update security group
    newDesc := "Updated demo security group"
    updatedSG, err := sgClient.Update(ctx, sg.ID, &securitygroups.SecurityGroupUpdateRequest{
        Description: &newDesc,
    })
    if err != nil {
        log.Fatalf("Update failed: %v", err)
    }
    fmt.Printf("✓ Updated description: %s\n", updatedSG.Description)

    // Step 5: Delete one rule
    if err := rules.Delete(ctx, httpRule.ID); err != nil {
        log.Fatalf("Delete rule failed: %v", err)
    }
    fmt.Printf("✓ Deleted HTTP rule\n")

    // Step 6: Delete security group
    if err := sgClient.Delete(ctx, sg.ID); err != nil {
        log.Fatalf("Delete failed: %v", err)
    }
    fmt.Printf("✓ Deleted security group: %s\n", sg.ID)

    fmt.Println("\n✅ Full lifecycle complete!")
}

func intPtr(v int) *int {
    return &v
}
```

**Output**:
```
✓ Created: demo-web-sg (ID: sg-demo123)
✓ Added HTTP rule (ID: rule-http001)
✓ Added HTTPS rule (ID: rule-https001)
✓ Listed 4 security groups
  - demo-web-sg: 2 rules
✓ Updated description: Updated demo security group
✓ Deleted HTTP rule
✓ Deleted security group: sg-demo123

✅ Full lifecycle complete!
```

---

## Error Handling Examples

### Handle 404 Not Found

```go
sg, err := sgClient.Get(ctx, "sg-nonexistent")
if err != nil {
    // Check error message or type for 404 handling
    // Error wrapping preserves context from internal/http client
    fmt.Printf("Error: %v\n", err)
    return nil  // Handle gracefully
}
```

### Handle 409 Conflict (Duplicate Name)

```go
sg, err := sgClient.Create(ctx, &securitygroups.SecurityGroupCreateRequest{
    Name: "web-servers",  // Already exists
})
if err != nil {
    // Check error message for conflict indication
    fmt.Printf("Create failed: %v\n", err)
    // Consider using Update instead or choosing a different name
    return fmt.Errorf("create failed: %w", err)
}
```

### Handle Context Timeout

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

resp, err := sgClient.List(ctx, nil)
if err != nil {
    if ctx.Err() == context.DeadlineExceeded {
        fmt.Println("Request timed out after 5 seconds")
        return nil
    }
    return fmt.Errorf("list failed: %w", err)
}
```

---

## Key Takeaways

1. **Client Initialization**: Use `cloudsdk.New(baseURL, token)` to create the SDK client
2. **Sub-Resource Pattern**: Use `sg.Rules().Create()` and `sg.Rules().Delete()` for rule management
3. **Detail Parameter**: Set `Detail: &true` to include rules in List response (default: false for performance)
4. **Pointer Types**: Request models use pointers for optional fields (`*string`, `*int`)
5. **Breaking Changes**: `Total` field removed from list response; use `len(resp.SecurityGroups)`
6. **Error Handling**: Errors are wrapped with context from internal/http client
7. **Context Support**: All methods accept `context.Context` for cancellation/timeout

---

## Next Steps

- Review [data-model.md](./data-model.md) for complete model definitions
- Check [contracts/](./contracts/) for interface documentation
- See [spec.md](./spec.md) for acceptance criteria and requirements
- Run tests: `go test ./modules/vps/securitygroups/...`
- Check coverage: `make coverage`
