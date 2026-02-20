## Context

The VPS SDK has an established SubResource pattern for child resources that belong to a parent:

- **`modules/vps/servers`** — `ServerResource` wraps `*servers.Server` and exposes `NICs()` + `Volumes()` sub-clients
- **`modules/vps/networks`** — `NetworkResource` wraps `*networks.Network` and exposes `Ports()`

Share access rules (`ShareRule`) are a natural child resource of a `Share`, but are currently exposed as flat `AddRule`/`RemoveRule` methods on `shares.Client`. This diverges from the pattern and requires callers to keep track of `shareID` separately when operating on rules.

## Goals / Non-Goals

**Goals:**
- Introduce `ShareResource` + `RulesClient` + `RuleOperations` following the same shape as `ServerResource`/`NICsClient`
- Change `Client.Get`, `Client.List`, and `Client.Create` to return `*ShareResource` instead of `*Share`
- Remove `AddRule`/`RemoveRule` from `shares.Client`; expose `Create`/`Delete` via `ShareResource.Rules()`
- Fix the `Create` rule return type to `*ShareRule` (current `AddRule` incorrectly returns `*Share`)

**Non-Goals:**
- Adding a `List` rule operation (not in the current API surface; can be added later)
- Changing URL paths or HTTP semantics
- Modifying other share methods (`Update`, `Delete`, `Extend`)

## Decisions

### Decision 1: Mirror the `ServerResource`/`NICsClient` shape exactly

`ShareResource` embeds `*shares.Share` and holds a `ruleOps RuleOperations` field, exposed via `Rules() RuleOperations`. `RulesClient` is a separate struct in `modules/vps/shares/rules.go` with `shareID` baked in at construction time.

**Alternative considered**: Keep flat methods but rename them. Rejected — callers still need to thread `shareID` manually and the surface is inconsistent with the rest of the SDK.

### Decision 2: `Get`, `List`, and `Create` all return `*ShareResource`

Following the servers client, all methods that return a share model are updated to wrap it in `ShareResource`. This ensures that any retrieved share immediately has `Rules()` available.

**Alternative considered**: Only `Get` returns `*ShareResource`. Rejected — the servers precedent wraps at every return site; inconsistency within `shares.Client` would be confusing.

### Decision 3: `RuleOperations` interface with `Create` and `Delete` only

```go
type RuleOperations interface {
    Create(ctx context.Context, req *sharesmodel.CreateShareRuleRequest) (*sharesmodel.ShareRule, error)
    Delete(ctx context.Context, ruleID string) error
}
```

Verb naming (`Create`/`Delete`) aligns with REST semantics used elsewhere in the SDK (e.g., `servers.Client.Create`). The previous `Add`/`Remove` naming was non-standard.

### Decision 4: Rename request model to `CreateShareRuleRequest`

The existing `AddShareRuleRequest` is renamed to `CreateShareRuleRequest` to match the new verb and align with the `CreateShareRequest` naming convention in the same package.

## Risks / Trade-offs

- **BREAKING change in `shares.Client`** — callers using `AddRule`/`RemoveRule` on the flat client must migrate. Mitigation: document in CHANGELOG; the method removal is a compile-time error so migration is forced and clear.
- **`Get` return type change** — callers storing `*shares.Share` from `Get` must update their type. Mitigation: same — compile-time error, no silent breakage.
- **`List` return type change** — `[]*shares.Share` → `[]*ShareResource`. Callers that iterate and pass shares elsewhere will need updating. Mitigation: compile-time error.
