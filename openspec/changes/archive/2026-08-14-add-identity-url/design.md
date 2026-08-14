## Context

`Client` currently retains one base URL, and `IAM()` always constructs its service client from `baseURL + "/iam"`. `Project()` depends on `IAM()` for direct project retrieval and project-code resolution, while `ProjectClient.VPS()` and `ProjectClient.VRM()` independently derive their URLs from `baseURL`. See `proposal.md` for motivation and the delta specifications for required behavior.

## Goals / Non-Goals

**Goals:**
- Add a public, option-compatible constructor for clients with separate general and identity endpoints.
- Keep existing `New` callers and their IAM routing behavior unchanged.
- Centralize IAM endpoint choice so direct IAM calls and project resolution cannot diverge.

**Non-Goals:**
- Add dedicated endpoint configuration for VPS, VRM, or other services.
- Change IAM operation paths, authentication, retries, HTTP-client configuration, or error semantics.
- Modify CLI environment-variable configuration in this change.

## Decisions

### Store the optional identity URL on the top-level client

Extend `Client` with an identity URL value populated only by `NewWithIdentityURL`. Keep `New` as the existing two-URL-parameter API so source and behavior compatibility are retained. Both constructors will share the existing client initialization and option application path after URL validation.

An option-based identity URL was considered, but a separate constructor matches the requested public API and makes the distinct endpoint configuration explicit at construction time.

### Derive the IAM base in one place

Make `Client.IAM()` select the configured identity URL when present; otherwise it will preserve its current base-URL-derived endpoint. Since `Project()` acquires its IAM client through this method, its IAM traffic follows the same selection automatically.

Passing endpoint selection into each IAM operation was considered, but would duplicate routing logic and risk inconsistent behavior between direct and project-resolution calls.

### Preserve service-path composition

The selected host continues to receive the IAM service suffix, yielding an IAM service base equivalent to `<selected-url>/iam`. VPS and VRM continue to append their existing suffixes only to `baseURL`.

Treating `identityURL` as a fully suffixed IAM URL was considered, but conflicts with the current constructor convention where service paths are appended by the SDK.

## Risks / Trade-offs

- [Trailing slashes can produce doubled path separators] -> Preserve current URL concatenation behavior and cover the expected no-trailing-slash inputs in routing tests; broader normalization is outside this change.
- [New constructor validation can differ from `New`] -> Reuse equivalent URL scheme validation and test invalid identity URLs.
- [Project resolution could accidentally retain the legacy endpoint] -> Use separate mock API and identity servers to verify all project-resolution IAM requests reach the identity server.

## Migration Plan

1. Release the additive constructor without changing `New` or existing call sites.
2. Consumers with separately deployed IAM endpoints opt in by constructing clients with `NewWithIdentityURL`.
3. Roll back by returning affected consumers to `New`; no persisted state or wire-contract migration is required.
