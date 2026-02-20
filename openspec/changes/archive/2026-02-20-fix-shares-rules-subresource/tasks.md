## 1. Models

- [x] 1.1 Rename `AddShareRuleRequest` to `CreateShareRuleRequest` in `models/vps/shares/share.go`
- [x] 1.2 Add `Validate()` to `CreateShareRuleRequest` — require non-empty `CIDR`; require `Level` in `{rw, ro}`
- [x] 1.3 Update any existing references to `AddShareRuleRequest` in test files

## 2. RulesClient

- [x] 2.1 Create `modules/vps/shares/rules.go` with `RulesClient` struct holding `baseClient`, `projectID`, `shareID`
- [x] 2.2 Implement `Create(ctx, req) (*sharesmodel.ShareRule, error)` — POST `.../shares/{share-id}/rules`, call `req.Validate()` before request
- [x] 2.3 Implement `Delete(ctx, ruleID) error` — DELETE `.../shares/{share-id}/rules/{rule-id}`

## 3. RuleOperations Interface & ShareResource

- [x] 3.1 Define `RuleOperations` interface in `modules/vps/shares/client.go` with `Create` and `Delete` signatures
- [x] 3.2 Add compile-time interface assertion: `var _ RuleOperations = (*RulesClient)(nil)`
- [x] 3.3 Define `ShareResource` struct embedding `*sharesmodel.Share` with `ruleOps RuleOperations` field
- [x] 3.4 Add `Rules() RuleOperations` method on `ShareResource`
- [x] 3.5 Add helper `newShareResource(baseClient, projectID, share)` to construct `ShareResource` with wired `RulesClient`

## 4. Update shares.Client Return Types

- [x] 4.1 Update `Get` to return `*ShareResource` — wrap response with `newShareResource`
- [x] 4.2 Update `List` to return `[]*ShareResource` — wrap each share in `newShareResource`
- [x] 4.3 Update `Create` to return `*ShareResource` — wrap response with `newShareResource`
- [x] 4.4 Remove `AddRule` method from `shares.Client`
- [x] 4.5 Remove `RemoveRule` method from `shares.Client`

## 5. Tests

- [x] 5.1 Create `modules/vps/shares/rules_test.go` — test `Create` happy path, `Create` validation errors (missing CIDR, invalid level), `Delete` happy path
- [x] 5.2 Update `modules/vps/shares/client_test.go` — update `Get`, `List`, `Create` assertions to use `*ShareResource`
- [x] 5.3 Update `modules/vps/shares/client_test.go` — remove test cases for `AddRule` and `RemoveRule`
- [x] 5.4 Update `models/vps/shares/share_test.go` — add `Validate()` tests for `CreateShareRuleRequest`

## 6. Verify

- [x] 6.1 Run `make build` — confirm no compile errors
- [x] 6.2 Run `make test` — confirm all tests pass
- [x] 6.3 Run `make lint` — confirm no lint warnings
