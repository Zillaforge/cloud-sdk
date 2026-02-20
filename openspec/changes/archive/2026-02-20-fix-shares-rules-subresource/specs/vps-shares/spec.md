# Spec Delta: vps-shares

## MODIFIED Requirements

### Requirement: Get share
The SDK SHALL provide a `Get` method that fetches a single share by ID and returns it wrapped in a `ShareResource`.

#### Scenario: Get existing share
- **WHEN** `Get(ctx, shareID)` is called with a valid ID
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/shares/{share-id}` and return `*ShareResource` with `Rules()` wired

#### Scenario: Get non-existent share
- **WHEN** `Get(ctx, shareID)` returns an HTTP error from the server
- **THEN** the client SHALL propagate the error wrapping the share ID

---

### Requirement: List shares
The SDK SHALL provide a `List` method on `shares.Client` that returns all shares for a project as `[]*ShareResource`, supporting optional filters.

#### Scenario: List all shares
- **WHEN** `List(ctx, nil)` is called
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/shares` and return a slice of `*ShareResource`, each with `Rules()` wired

#### Scenario: List shares with name filter
- **WHEN** `List(ctx, &ListSharesOptions{Name: "my-share"})` is called
- **THEN** the client SHALL append `?name=my-share` to the request URL

---

### Requirement: Create share
The SDK SHALL provide a `Create` method that creates a new share and returns it as a `ShareResource`.

#### Scenario: Create share with required fields
- **WHEN** `Create(ctx, req)` is called with `Name`, `NetworkID`, and `Size` set
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/shares` with the JSON body and return `*ShareResource` with `Rules()` wired

#### Scenario: Create share with missing required field
- **WHEN** `Create(ctx, req)` is called without `Name`
- **THEN** the client SHALL return a validation error before making any HTTP call

---

## REMOVED Requirements

### Requirement: Add share access rule
**Reason**: Rule operations are now a SubResource accessed via `ShareResource.Rules().Create()`. Flat `AddRule` on `shares.Client` is removed to align with the SubResource pattern used by NIC and Port.
**Migration**: Replace `client.AddRule(ctx, shareID, req)` with `shareResource.Rules().Create(ctx, req)` where `shareResource` is obtained from `Get`, `List`, or `Create`.

---

### Requirement: Remove share access rule
**Reason**: Rule operations are now a SubResource accessed via `ShareResource.Rules().Delete()`. Flat `RemoveRule` on `shares.Client` is removed to align with the SubResource pattern used by NIC and Port.
**Migration**: Replace `client.RemoveRule(ctx, shareID, ruleID)` with `shareResource.Rules().Delete(ctx, ruleID)` where `shareResource` is obtained from `Get`, `List`, or `Create`.
