## 1. Client Configuration

- [x] 1.1 Add identity URL state to the top-level client and implement `NewWithIdentityURL` with equivalent URL, token, default HTTP-client, and option handling to `New`.
- [x] 1.2 Update IAM client construction to select the configured identity URL and append the existing IAM service path, while preserving the base-URL fallback for `New` clients.

## 2. Regression Coverage

- [x] 2.1 Add constructor tests for valid identity URL configuration, option preservation, and invalid or scheme-less identity URLs.
- [x] 2.2 Add request-routing tests proving direct IAM calls and `Project` IAM lookups use the dedicated identity endpoint.
- [x] 2.3 Add compatibility and service-isolation tests proving `New` still routes IAM through `baseURL + "/iam"` and VPS/VRM remain based on `baseURL` when an identity URL is configured.

## 3. Verification

- [x] 3.1 Run the affected package tests and `go test ./...` to verify the new routing behavior and regressions.
