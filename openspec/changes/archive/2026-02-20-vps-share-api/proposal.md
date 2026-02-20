## Why

The VPS SDK is missing a Go client for the Share API, which provides managed NFS-like shared storage resources. All other major VPS resource types (volumes, snapshots, networks, servers, etc.) already have SDK coverage; shares are the final missing piece.

## What Changes

- New `models/vps/shares/` package with request/response types for shares and share rules
- New `modules/vps/shares/` package with a `Client` implementing CRUD + extend + rule management
- Share client registered on the VPS module client (following the pattern of volumes, networks, etc.)

## Capabilities

### New Capabilities
- `vps-shares`: Full CRUD (list, get, create, update, delete) + extend + access rule management (add/remove) for VPS shares

### Modified Capabilities
<!-- No existing specs are changing requirements -->

## Impact

- New files: `models/vps/shares/share.go`, `models/vps/shares/share_test.go`, `modules/vps/shares/client.go`, `modules/vps/shares/client_test.go`
- `modules/vps/core/client.go` — add `Shares` field and constructor wiring
- No breaking changes; purely additive
