# Feature Specification: VPS Project APIs SDK

**Feature Branch**: `[1-vps-project-api]`  
**Created**: 2025-10-26  
**Status**: Draft  
**Input**: User description: "implement virtual platform service (vps) restful api based on #file:vps.json. Only implement api which has endpoint start with `/api/v1/project/`, and ignore others. All APIs are authorized through Bearer Token, and client has to provide the token when initilization."

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST provide a client-based, typed interface; callers MUST NOT manage raw HTTP details.
- Public APIs MUST be consistently typed and return structured results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## Clarifications

### Session 2025-10-26

- Q: What SDK retry policy should apply for transient errors? → A: Auto-retry only safe reads (GET/HEAD) on 429/502/503/504 with exponential backoff + jitter, max 3 attempts.
 - Q: How should the SDK handle asynchronous operations (202 Accepted)? → A: Offer optional Waiter helpers to poll resource state with context/backoff; default calls return immediately.
 - Q: How should client-side failures without an HTTP response be represented? → A: Use statusCode = 0 for client-side failures (timeout/canceled/network), with detailed message and meta.
 - Q: What project scoping model should the VPS client use? → A: Bind VPS client to a specific Project ID (e.g., sdk.Project(projectID).VPS()); methods do not take project-id.
 - Q: What default per-request timeout should the SDK use? → A: Default 30s timeout; override via context and client options.
 - Q: How should client-side failures without an HTTP response be represented? → A: Use statusCode = 0 for client-side failures (timeout/canceled/network), with detailed message and meta.

## User Scenarios & Testing (mandatory)

### User Story 1 - Manage Networks (Priority: P1)

As a project user or admin, I need to list/create/view/update/delete networks and inspect ports on a network.

**Independent Test**: The story is complete if a user can create a network, list/get/update/delete it, and list ports for a given network.

**Acceptance Scenarios**:
1. Given a valid Project ID, When I list networks with filters, Then I receive 200 with matching items.
2. Given valid input, When I create a network, Then I receive 201 with network info.
3. Given a network ID, When I get/update/delete it, Then I receive 200/204 accordingly.
4. Given a network ID, When I list ports, Then I receive 200 with port items.

---

### User Story 2 - Manage Floating IPs (Priority: P1)

As a project user or admin, I need to list, create, view, update, delete, approve/reject (admin), and disassociate floating IPs.

**Independent Test**: The story is complete if a user can allocate a FIP, retrieve/update it, associate/disassociate with devices, and an admin can approve/reject pending creation requests.

**Acceptance Scenarios**:
1. Given a valid Project ID, When I list FIPs with filters, Then I receive 200 with matching items.
2. Given valid input, When I create a FIP, Then I receive 201 with FIP details.
3. Given a FIP ID, When I get/update/delete it, Then I receive 200/204 and state reflects the change.
4. Given a FIP ID, When I disassociate it, Then I receive 202 and the association is removed.
5. Given a pending FIP request, When a tenant admin approves/rejects, Then I receive 202 and the request is resolved.

---

### User Story 3 - Manage Servers (Priority: P1)

As a project user or admin, I need to list/create/view/update/delete servers, perform server actions (stop/start/reboot/resize/extend_root/get_pwd; approve/reject by admin), and retrieve metrics.

**Independent Test**: The story is complete if a user can create a server, list/get/update/delete it, perform the supported actions with required parameters, and fetch metrics successfully.

**Acceptance Scenarios**:
1. Given a valid Project ID, When I list servers with filters, Then I receive 200 with items.
2. Given valid server input, When I create, Then I receive 201 with server info.
3. Given a server ID, When I get/update/delete it, Then I receive 200/202/204 accordingly.
4. Given a server ID, When I perform an action, Then I receive 202 with action result. Supported actions: start, stop, reboot (with reboot_type), resize (with flavor_id), extend_root (with root_size), get_pwd (with private_key), approve, reject.
5. Given a server ID, When I fetch metrics with type and time range, Then I receive 200 with metric series.
6. Given a server ID, When I list NICs, Then I receive 200 with NIC items.
7. Given a server ID and network/sg parameters, When I add a NIC, Then I receive 201/200 as defined and the NIC appears attached.
8. Given a server ID and NIC ID, When I update the NIC security groups, Then I receive 200 and the NIC reflects the new sg_ids.
9. Given a server ID and NIC ID, When I delete the NIC, Then I receive 204 and the NIC is removed.
10. Given a server ID and NIC ID, When I associate a floating IP to that NIC, Then I receive 200/202 as defined and the association exists.
11. Given a server ID, When I request a VNC URL, Then I receive 200 with a console URL.
12. Given a server ID, When I list volume attachments, Then I receive 200 with attached volumes.
13. Given a server ID and Volume ID, When I attach the volume, Then I receive 200/202 as defined and the volume is attached.
14. Given a server ID and Volume ID, When I detach the volume, Then I receive 204/202 as defined and the volume is detached.

---

### User Story 4 - Manage Keypairs (Priority: P2)

As a project user, I need to list, create, view, update description, and delete keypairs.

**Independent Test**: The story is complete if a user can import or generate a keypair, list/get it, update description, and delete it.

**Acceptance Scenarios**:
1. Given a valid Project ID, When I list keypairs (optionally by name), Then I receive 200 with items.
2. Given valid input, When I create a keypair, Then I receive 201 with keypair info.
3. Given a keypair ID, When I get/update/delete it, Then I receive 200/204 accordingly.

---

### User Story 5 - Manage Routers (Priority: P2)

As a project user or admin, I need to list/create/view/update/delete routers, enable/disable them, and associate/disassociate networks.

**Independent Test**: The story is complete if a user can create a router, toggle state via action, list associated networks, and attach/detach a network.

**Acceptance Scenarios**:
1. Given a valid Project ID, When I list routers with filters, Then I receive 200 with items.
2. Given valid input, When I create a router, Then I receive 201 with router info.
3. Given a router ID, When I get/update/delete it, Then I receive 200/204 accordingly.
4. Given a router ID, When I call action set_state, Then I receive 202 and state reflects the change.
5. Given a router and network IDs, When I associate or disassociate, Then I receive 204 and relationships update.

---

### User Story 6 - Manage Security Groups (Priority: P2)

As a project user or admin, I need to list/create/view/update/delete security groups and add/remove rules.

**Independent Test**: The story is complete if a user can create a security group, view/update/delete it, and add/remove rules with valid protocol/port/cidr specs.

**Acceptance Scenarios**:
1. Given a valid Project ID, When I list security groups with filters, Then I receive 200 with items.
2. Given valid input, When I create a security group (with or without rules), Then I receive 201 with group info.
3. Given a security group ID, When I get/update/delete it, Then I receive 200/204 accordingly.
4. Given a security group ID, When I add or delete a rule, Then I receive 200 and rule set reflects the change.

---

### User Story 7 - Discover Flavors (Priority: P3)

As a project user, I need to list and view flavors with filters (name, public, tag) to select instance sizes.

**Independent Test**: The story is complete if a user can list and get flavors and apply filters as documented.

**Acceptance Scenarios**:
1. Given a valid Project ID, When I list flavors with filters, Then I receive 200 with filtered items.
2. Given a flavor ID, When I get it, Then I receive 200 with details.

### User Story 8 - View Project Quotas (Priority: P2)

As a project user or admin, I need to view my project's quota limits and current usage so I can understand remaining capacity for resources (VMs, vCPU, RAM, GPU, block storage, networks, routers, floating IPs, shares).

**Independent Test**: The story is complete if a user can retrieve quotas for a project and see both limits and usage fields populated per Swagger model.

**Acceptance Scenarios**:
1. Given a valid Project ID, When I get project quotas, Then I receive 200 with `QuotaInfo` containing limits and usage for each resource type.
2. Given an unauthorized or invalid token, When I get project quotas, Then I receive 401 with a structured error.
3. Given a Project ID I cannot access, When I get project quotas, Then I receive 403/404 with a structured error.

## Edge Cases

- Invalid or expired Bearer token → 401 Unauthorized with structured error.
- Project ID not found or access denied → 403/404 with structured error.
- Quota exceeded (Servers, Floating IPs, Networks) → 400 with explicit error payload.
- Asynchronous operations (202 Accepted) require polling or eventual consistency in tests (e.g., server/routers actions).
- Parameter validation errors (e.g., thresholds out of bounds) → 400 with details.
- Conflict states (e.g., delete resource in use) → 409 if defined, otherwise 400 with context.
- Router state toggles or network associations invalid for current state → 400 with details.
- Security group rule protocol/port/cidr invalid combinations → 400 with details.
- Transient errors on safe read operations (429/502/503/504) trigger automatic retries with exponential backoff + jitter, up to 2 retries (max 3 attempts total).
 - Network/timeout/client-canceled failures without HTTP response map to structured errors with statusCode=0, a non-empty message, and meta describing the cause.

## Requirements (mandatory)

### Functional Requirements

- FR-000: SDK MUST be initialized with a Base URL (target API server) and a Bearer token; all requests derive from these settings.
- FR-001: System MUST support Bearer token authorization for all project-scoped endpoints.
- FR-002: System MUST provide Network operations under `/api/v1/project/{project-id}/networks*` (list, create, get, update, delete; list ports).
- FR-003: System MUST provide Floating IP operations under `/api/v1/project/{project-id}/floatingips*` (list, create, get, update, delete; approve; reject; disassociate).
 - FR-004: System MUST provide Server operations under `/api/v1/project/{project-id}/servers*`, including:
	 - CRUD: list, create, get, update, delete
	 - Actions: start, stop, reboot (with reboot_type), resize (with flavor_id), extend_root (with root_size), get_pwd (with private_key), approve, reject
	 - Metrics: get server metrics
	 - NICs: list NICs; add NIC; update NIC security groups; delete NIC
	 - Floating IP: associate floating IP to a specific NIC
	 - Console: get VNC URL
	 - Volumes: list attachments; attach volume; detach volume
- FR-005: System MUST provide Keypair operations under `/api/v1/project/{project-id}/keypairs*` (list, create, get, update, delete).
- FR-006: System MUST provide Router operations under `/api/v1/project/{project-id}/routers*` (list, create, get, update, delete; action set_state; list associated networks; associate/disassociate network).
- FR-007: System MUST provide Security Group operations under `/api/v1/project/{project-id}/security_groups*` (list, create, get, update, delete; rule create/delete).
- FR-008: System MUST provide Flavor discovery under `/api/v1/project/{project-id}/flavors*` (list with filters; get).
- FR-009: System MUST provide Quota retrieval under `/api/v1/project/{project-id}/quotas` (get).
- FR-010: System MUST ignore non-project-scoped endpoints (e.g., `/api/v1/admin/*`, root features) in this feature.
- FR-011: Error responses MUST be mapped to a structured error that includes the HTTP status code and a detailed error message, preserving server fields from Swagger definitions (e.g., `vpserr.ErrResponse` with `errorCode`, `message`, `meta`).
- FR-012: Request/response payloads MUST align with Swagger models for each endpoint.
- FR-013: Pagination/filters (where defined) MUST be supported and tested.
- FR-014: SDK MUST expose per-service clients. For this feature, a VPS client MUST encapsulate all VPS operations; other services will follow the same pattern in future features.
- FR-015: Read-only operations (GET/HEAD) MUST auto-retry on 429/502/503/504 with exponential backoff and jitter, up to a maximum of 3 attempts (initial + up to 2 retries). Non-idempotent operations MUST NOT be retried automatically.
- FR-016: SDK MUST provide optional waiter helpers for asynchronous operations that return 202, enabling callers/tests to poll for desired states using context, with configurable polling/backoff, without blocking by default.
 - FR-017: VPS client MUST be project-scoped: callers obtain a project-scoped client (e.g., via selecting a project) and then invoke VPS methods without passing project-id parameters. Multiple project-scoped clients MAY be created from the same SDK instance.
 - FR-018: The SDK MUST enforce a default per-request timeout of 30 seconds. Callers MAY override timeouts via context deadlines or client-level options; the earliest applicable deadline wins. On timeout, return a structured error per contract.

 

### SDK Initialization & Client Structure

- The SDK provides a top-level initializer that accepts:
	- Base URL of the target API server (including scheme, e.g., https://api.example.com)
	- Bearer token used for authorization on every request
- The SDK exposes a dedicated VPS client to perform all operations covered by this specification; callers do not construct URLs or set headers manually.
- Service clients inherit the base URL and authorization from the top-level SDK initialization.
- Public methods are consistent and typed, and accept a request context for cancellation and deadlines.
 - Project scoping: The SDK provides a way to select a project (e.g., `Project(<project-id>)`) that returns a project-scoped handle from which service clients (e.g., VPS) are obtained. All VPS methods operate within that project and do not accept a project-id argument. Creating multiple project-scoped handles from the same SDK instance is supported.

### Error Handling Contract

- All non-2xx responses return a structured error object to callers that contains at minimum:
	- statusCode: integer HTTP status code from the response
	- errorCode: integer service error code (from `vpserr.ErrResponse.errorCode`, if present)
	- message: detailed error message for human consumption (from `vpserr.ErrResponse.message`, if present)
	- meta: optional key-value object with additional fields from the server (from `vpserr.ErrResponse.meta`, if present)
- If the server body is missing or malformed, the SDK MUST still populate `statusCode` and a non-empty `message` explaining the failure (e.g., parse error, empty body).
 - Client-side failures with no HTTP response (e.g., network error, DNS failure, connection reset, TLS handshake failure, timeout, or context canceled) must be represented as:
	 - statusCode = 0
	 - message is non-empty and human-readable
	 - meta includes the underlying error string and a coarse category (e.g., "timeout", "canceled", "network") when available

### Resiliency & Retry Policy

- Scope: Automatic retries apply only to safe, read-only requests (GET/HEAD).
- Triggers: HTTP 429, 502, 503, 504.
- Limits: Max 3 attempts total (initial try + up to 2 retries).
- Backoff: Exponential backoff with jitter between attempts.
- Exclusions: Non-idempotent operations (e.g., resource creation, server action POSTs) are never auto-retried by default.
- Observability: Each retry attempt should be visible via debug logs or hooks without leaking secrets.

### Timeouts

- Default: 30 seconds per request if no explicit context deadline is provided.
- Overrides: Per-call via context.WithTimeout/WithDeadline; per-client via configurable option. The effective timeout is the earliest bound.
- Behavior: On timeout, the request is canceled, and a structured error is returned with statusCode=0, message indicating timeout, and meta containing cause.
- Waiters: Respect context deadlines; do not exceed caller-specified timeouts.

### Asynchronous Operations & Waiters

- Default behavior: Methods return immediately; the SDK does not auto-wait for completion.
- Waiters: Provide resource-specific helpers (e.g., WaitForServerStatus) that poll for a terminal or target state.
- Control: Accept context for cancellation/deadlines; support configurable initial delay, interval/backoff, and maximum wait duration.
- Errors: If the target state is not reached within the allowed time or the operation fails, return a structured error with a non-empty message.
- Tests: Contract/integration tests may use waiters to validate eventual state transitions deterministically.

### Key Entities (include if feature involves data)

- Project: identifies tenant scope for operations.
- Server: compute instance with lifecycle operations, actions, and metrics.
- Network: L2 network with CIDR, ports, and metadata.
- Router: L3 routing resource with state and network associations.
- SecurityGroup: security policy container with ingress/egress rules.
- FloatingIP: public address resource attachable to project resources.
- Keypair: SSH key resource for instances.
- Flavor: compute shape metadata used to size servers.
- Quota: resource limits/usage for the project.
- ErrorResponse: structured error containing `statusCode`, `errorCode`, `message`, and optional `meta`, derived from HTTP status and `vpserr.ErrResponse`.

## Success Criteria (mandatory)

### Measurable Outcomes

- SC-001: All project-scoped endpoints covered by this feature have passing unit and contract tests derived from Swagger.
- SC-002: Error handling demonstrates correct mapping for at least 10 representative negative cases across resources, including correct HTTP status code and a non-empty detailed message.
- SC-003: Security: 0 instances of secrets logged during automated test runs.
- SC-004: Developers can initialize the SDK with a Base URL and Bearer token, obtain a VPS client instance, and call at least one operation without writing any manual HTTP code.
- SC-005: Under injected 429/502/503/504 on read operations, the SDK performs up to 2 retries (max 3 attempts) with exponential backoff + jitter and either succeeds or returns a structured error with correct status and message.
- SC-006: For actions returning 202, waiter helpers can poll to the expected target state within a bounded time using context and backoff, or return a structured timeout/error without hanging.
 - SC-007: Under simulated client-side failures (timeout, canceled, network unreachable), the SDK returns structured errors with statusCode=0, a non-empty message, and preserves cause in meta.
 - SC-008: Developers can obtain a project-scoped VPS client and perform operations without supplying project-id per call; switching projects is achieved by creating another project-scoped client.
 - SC-009: Without an explicit context deadline, requests time out at ~30s and return structured timeout errors; with a shorter caller-provided deadline, calls honor the shorter deadline.
 - SC-008: Developers can obtain a project-scoped VPS client and perform operations without supplying project-id per call; switching projects is achieved by creating another project-scoped client.
 - SC-007: Under simulated client-side failures (timeout, canceled, network unreachable), the SDK returns structured errors with statusCode=0, a non-empty message, and preserves cause in meta.

## Assumptions

- The service base URL and basePath (`/vps`) are configurable; initialization values are provided by the caller. Feature scope remains project-scoped endpoints only.
- Authorization is via HTTP Bearer token set at client initialization and sent on every request.
- Some operations are asynchronous (202 Accepted) and may require polling in tests (contract/integration) rather than immediate state changes.
