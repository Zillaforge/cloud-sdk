# Spec: vps-shares-rules

## Purpose

Defines the SDK requirements for share access rule operations exposed as a SubResource via `ShareResource.Rules()`. Covers creation and deletion of access rules, request validation, and sub-client wiring.

---

## ADDED Requirements

### Requirement: ShareResource wraps Share with Rules sub-client
The SDK SHALL provide a `ShareResource` struct that embeds `*shares.Share` and exposes a `Rules() RuleOperations` method returning a scoped `RulesClient` for that share.

#### Scenario: Access rules client from ShareResource
- **WHEN** a caller receives a `*ShareResource` from `Get`, `List`, or `Create`
- **THEN** calling `Rules()` SHALL return a `RuleOperations` bound to that share's ID

---

### Requirement: Create share access rule
The SDK SHALL provide a `Create` method on `RuleOperations` that adds an access rule to a share.

#### Scenario: Create rule with required fields
- **WHEN** `shareResource.Rules().Create(ctx, req)` is called with `CIDR` and `Level` set
- **THEN** the client SHALL issue `POST /api/v1/project/{project-id}/shares/{share-id}/rules` with the JSON body and return `*ShareRule`

#### Scenario: Create rule with missing CIDR
- **WHEN** `shareResource.Rules().Create(ctx, req)` is called without `CIDR`
- **THEN** the client SHALL return a validation error before making any HTTP call

#### Scenario: Create rule with invalid level
- **WHEN** `shareResource.Rules().Create(ctx, req)` is called with `Level` not in `{rw, ro}`
- **THEN** the client SHALL return a validation error before making any HTTP call

---

### Requirement: Delete share access rule
The SDK SHALL provide a `Delete` method on `RuleOperations` that removes an access rule from a share.

#### Scenario: Delete rule by ID
- **WHEN** `shareResource.Rules().Delete(ctx, ruleID)` is called with a valid rule ID
- **THEN** the client SHALL issue `DELETE /api/v1/project/{project-id}/shares/{share-id}/rules/{rule-id}` and return nil on success

#### Scenario: Delete non-existent rule
- **WHEN** `shareResource.Rules().Delete(ctx, ruleID)` returns an HTTP error from the server
- **THEN** the client SHALL propagate the error wrapping both the rule ID and share ID

---

### Requirement: RuleOperations interface
The SDK SHALL define a `RuleOperations` interface satisfied by `RulesClient`.

#### Scenario: Interface compliance
- **WHEN** `RulesClient` is compiled
- **THEN** it SHALL satisfy `RuleOperations` at compile time (enforced via interface assertion)

---

### Requirement: CreateShareRuleRequest validation
The SDK SHALL provide a `CreateShareRuleRequest` model with a `Validate()` method.

#### Scenario: Valid request passes validation
- **WHEN** `Validate()` is called on a request with non-empty `CIDR` and `Level` in `{rw, ro}`
- **THEN** it SHALL return nil

#### Scenario: Missing CIDR fails validation
- **WHEN** `Validate()` is called with empty `CIDR`
- **THEN** it SHALL return an error containing "cidr is required"

#### Scenario: Invalid level fails validation
- **WHEN** `Validate()` is called with a `Level` value not in `{rw, ro}`
- **THEN** it SHALL return an error containing "level must be"
