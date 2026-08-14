## Why

Deployments can expose IAM on an identity-specific host rather than beneath the general API host. The SDK currently derives every IAM request from `baseURL`, preventing callers from targeting that separate endpoint while keeping VPS and VRM requests on the existing API host.

## What Changes

- Add an `identityURL` configuration value to the top-level SDK client.
- Add `NewWithIdentityURL` to construct a client with separate base and identity URLs while retaining existing client options.
- Route all IAM operations, including IAM calls made by `Project`, through the configured identity URL when the new constructor is used.
- Preserve `New` behavior: IAM continues to use `baseURL + "/iam"` when no identity URL is configured.
- Validate the supplied identity URL using the same URL requirements as the base URL.

## Capabilities

### New Capabilities
- `identity-endpoint-routing`: Configures a dedicated IAM endpoint for a Cloud SDK client without changing non-IAM service routing.

### Modified Capabilities
- `iam-api`: IAM operations can use a configured identity endpoint while retaining the existing derived IAM endpoint as the compatibility default.

## Impact

- Affects the public constructors and internal state in `client.go`.
- Affects IAM client construction used directly by `Client.IAM()` and indirectly by `Client.Project()`.
- Requires client-level tests for constructor validation, IAM URL selection, and legacy routing compatibility.
