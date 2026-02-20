# Spec: vps-shares

## Purpose

Defines the SDK requirements for the VPS Share API — NFS/CIFS managed file shares provided per-project. Covers CRUD operations, size extension, access rule management, request validation, and core client wiring.

---

## Requirements

### Requirement: List shares
The SDK SHALL provide a `List` method on `shares.Client` that returns all shares for a project, supporting optional filters.

#### Scenario: List all shares
- **WHEN** `List(ctx, nil)` is called
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/shares` and return a slice of `*Share`

#### Scenario: List shares with name filter
- **WHEN** `List(ctx, &ListSharesOptions{Name: "my-share"})` is called
- **THEN** the client SHALL append `?name=my-share` to the request URL

---

### Requirement: Get share
The SDK SHALL provide a `Get` method that fetches a single share by ID.

#### Scenario: Get existing share
- **WHEN** `Get(ctx, shareID)` is called with a valid ID
- **THEN** the client SHALL issue `GET /api/v1/project/{project-id}/shares/{share-id}` and return `*Share`

#### Scenario: Get non-existent share
- **WHEN** `Get(ctx, shareID)` returns an HTTP error from the server
- **THEN** the client SHALL propagate the error wrapping the share ID

---

### Requirement: Create share
The SDK SHALL provide a `Create` method that creates a new share.

#### Scenario: Create share with required fields
- **WHEN** `Create(ctx, req)` is called with `Name`, `NetworkID`, and `Size` set
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/shares` with the JSON body and return `*Share`

#### Scenario: Create share with missing required field
- **WHEN** `Create(ctx, req)` is called without `Name`
- **THEN** the client SHALL return a validation error before making any HTTP call

---

### Requirement: Update share
The SDK SHALL provide an `Update` method that updates share metadata (name, description).

#### Scenario: Update share name
- **WHEN** `Update(ctx, shareID, req)` is called with `Name` set
- **THEN** the client SHALL issue `PUT /api/v1/project/{project-id}/shares/{share-id}` and return updated `*Share`

---

### Requirement: Delete share
The SDK SHALL provide a `Delete` method that removes a share.

#### Scenario: Delete existing share
- **WHEN** `Delete(ctx, shareID)` is called
- **THEN** the client SHALL issue `DELETE /api/v1/project/{project-id}/shares/{share-id}` and return nil on success

---

### Requirement: Extend share size
The SDK SHALL provide an `Extend` method that increases a share's size.

#### Scenario: Extend share
- **WHEN** `Extend(ctx, shareID, req)` is called with `NewSize` greater than current size
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/shares/{share-id}/extend` with the JSON body

#### Scenario: Extend with invalid size
- **WHEN** `Extend(ctx, shareID, req)` is called with `NewSize` <= 0
- **THEN** the client SHALL return a validation error before making any HTTP call

---

### Requirement: Add share access rule
The SDK SHALL provide an `AddRule` method that creates an access rule on a share.

#### Scenario: Add rule with required fields
- **WHEN** `AddRule(ctx, shareID, req)` is called with `CIDR` and `Level` set
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/shares/{share-id}/rules` and return `*ShareRule`

#### Scenario: Add rule with missing CIDR
- **WHEN** `AddRule(ctx, shareID, req)` is called without `CIDR`
- **THEN** the client SHALL return a validation error before making any HTTP call

---

### Requirement: Remove share access rule
The SDK SHALL provide a `RemoveRule` method that deletes an access rule from a share.

#### Scenario: Remove rule
- **WHEN** `RemoveRule(ctx, shareID, ruleID)` is called
- **THEN** the client SHALL issue `DELETE /api/v1/project/{project-id}/shares/{share-id}/rules/{rule-id}` and return nil on success

---

### Requirement: Share model validation
The `Share` model SHALL validate required fields.

#### Scenario: Valid share
- **WHEN** a `Share` struct has non-empty `ID` and `Name`
- **THEN** `Validate()` SHALL return nil

#### Scenario: Invalid share missing ID
- **WHEN** a `Share` struct has empty `ID`
- **THEN** `Validate()` SHALL return an error

---

### Requirement: VPS core client exposes shares
The VPS core `Client` SHALL expose a `Shares()` method returning `*shares.Client`.

#### Scenario: Access shares client
- **WHEN** `vpsClient.Shares()` is called
- **THEN** it SHALL return a non-nil `*shares.Client` bound to the same project ID
