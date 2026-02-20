## 1. Share Models

- [x] 1.1 Create `models/vps/shares/share.go` with `Share`, `ShareRule`, status constants, and list response types
- [x] 1.2 Add `CreateShareRequest` with validation (name, network_id, size required)
- [x] 1.3 Add `UpdateShareRequest` with validation (no required fields)
- [x] 1.4 Add `ExtendShareRequest` with validation (new_size > 0)
- [x] 1.5 Add `AddShareRuleRequest` with validation (cidr and level required)
- [x] 1.6 Add `ListSharesOptions` struct for query filters (name, user_id, status)
- [x] 1.7 Create `models/vps/shares/share_test.go` with unit tests for all Validate() methods

## 2. Share Client

- [x] 2.1 Create `modules/vps/shares/client.go` with `Client` struct and `NewClient` constructor
- [x] 2.2 Implement `List(ctx, opts)` — GET /api/v1/project/{project-id}/shares
- [x] 2.3 Implement `Get(ctx, shareID)` — GET /api/v1/project/{project-id}/shares/{share-id}
- [x] 2.4 Implement `Create(ctx, req)` — POST /api/v1/project/{project-id}/shares
- [x] 2.5 Implement `Update(ctx, shareID, req)` — PUT /api/v1/project/{project-id}/shares/{share-id}
- [x] 2.6 Implement `Delete(ctx, shareID)` — DELETE /api/v1/project/{project-id}/shares/{share-id}
- [x] 2.7 Implement `Extend(ctx, shareID, req)` — POST /api/v1/project/{project-id}/shares/{share-id}/extend
- [x] 2.8 Implement `AddRule(ctx, shareID, req)` — POST /api/v1/project/{project-id}/shares/{share-id}/rules
- [x] 2.9 Implement `RemoveRule(ctx, shareID, ruleID)` — DELETE /api/v1/project/{project-id}/shares/{share-id}/rules/{rule-id}
- [x] 2.10 Create `modules/vps/shares/client_test.go` with HTTP mock tests for all client methods

## 3. VPS Core Client Wiring

- [x] 3.1 Add `shares` import to `modules/vps/core/client.go`
- [x] 3.2 Add `Shares() *shares.Client` accessor method to the VPS core `Client`

## 4. Verification

- [x] 4.1 Run `make build` — ensure no compile errors
- [x] 4.2 Run `make test` — all tests pass including new share tests
- [x] 4.3 Run `make lint` — no lint errors
