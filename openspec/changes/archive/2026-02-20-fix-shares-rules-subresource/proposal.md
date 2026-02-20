## Why

Share rule operations (`AddRule`, `RemoveRule`) are implemented as flat methods on `shares.Client`, inconsistent with the SubResource pattern used by NIC (on `ServerResource`) and Port (on `NetworkResource`). Aligning rules to a dedicated `RulesClient` exposed via a `ShareResource` wrapper improves consistency and makes the SDK surface predictable.

## What Changes

- **BREAKING** Remove `AddRule(ctx, shareID, req)` and `RemoveRule(ctx, shareID, ruleID)` from `shares.Client`
- Introduce `ShareResource` struct (wrapping `*shares.Share`) with a `Rules()` accessor returning a `RuleOperations` interface
- Introduce `RulesClient` in `modules/vps/shares/rules.go` implementing `RuleOperations` with:
  - `Create(ctx, req) (*ShareRule, error)` — POST .../rules
  - `Delete(ctx, ruleID) error` — DELETE .../rules/{rule-id}
- Fix return types: `Create` returns `*ShareRule` (not `*Share`); `Delete` returns `error` (not `*Share`)
- Update `Get` on `shares.Client` to return `*ShareResource` (wiring the `RulesClient`)
- Add a `ShareRule` model `Validate()` method and a dedicated `AddShareRuleRequest` model
- Update existing `vps-shares` spec to reflect the removed flat methods and changed return types

## Capabilities

### New Capabilities
- `vps-shares-rules`: SubResource client for share access rule operations (`Create`, `Delete`) exposed via `ShareResource.Rules()`

### Modified Capabilities
- `vps-shares`: Requirements for `AddRule`/`RemoveRule` are removed; `Get` return type changes from `*Share` to `*ShareResource`

## Impact

- `modules/vps/shares/client.go` — remove `AddRule`, `RemoveRule`; add `ShareResource`, `RuleOperations` interface; change `Get` return type
- `modules/vps/shares/rules.go` — new file implementing `RulesClient`
- `models/vps/shares/share.go` — fix `AddRule` return type annotation; add `Validate()` to `AddShareRuleRequest`
- `modules/vps/shares/client_test.go` — update tests for removed methods and `Get` return type
- `openspec/specs/vps-shares/spec.md` — delta to remove flat-rule requirements
- No API contract changes — URL paths remain identical
