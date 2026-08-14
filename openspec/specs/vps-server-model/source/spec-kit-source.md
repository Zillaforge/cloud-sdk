# Feature Specification: 修正 server 的model 與 module 符合文件的定義

**Feature Branch**: `005-fix-server-model`
**Created**: 2025-11-09
**Status**: Draft

## Clarifications

### Session 2025-11-09

- Q: Which of the following should be the approach for defining the Server model fields? → A: A

## User Scenarios & Testing

- As a developer, I need the server model and module to match the documented definitions so that all server-related operations are consistent, reliable, and testable.
- As a user of the SDK, I expect server CRUD, actions, metrics, NIC, and volume operations to behave as described in the documentation and Swagger spec.

**Acceptance Scenarios:**
1. Given a valid Project ID, when I list servers, I receive a response with all documented fields.
2. When I create, update, or delete a server, the request and response match the documented model.
3. When I perform server actions (start, stop, reboot, resize, extend_root, get_pwd, approve, reject), the module supports all required parameters and returns expected results.
4. When I fetch metrics, manage NICs, or attach/detach volumes, the model and module behave as defined in the documentation.

## Functional Requirements

- The server model MUST include all fields defined in the documentation and Swagger spec.
- Status fields MUST use custom types (e.g., `type ServerStatus string`) to simulate enums instead of plain strings, with const definitions for valid values.
- The server module MUST support all documented operations: CRUD, actions, metrics, NIC management, volume management.
- All request/response types MUST align with the documented data model.
- Error handling MUST follow the structured error contract.
- All changes MUST be covered by unit and contract tests.

## Success Criteria

- All server operations pass unit and contract tests.
- The server model and module match the documentation and Swagger spec.
- No missing or extra fields in server-related request/response types.
- All acceptance scenarios are covered by tests.
- No implementation details (languages, frameworks, APIs) leak into the specification.

## Key Entities

- Server
- ServerCreateRequest
- ServerUpdateRequest
- ServerActionRequest
- ServerActionResponse
- ServerMetricsRequest
- ServerMetricsResponse
- ServerNIC
- ServerNICCreateRequest
- ServerNICUpdateRequest
- VolumeAttachment

## Data Model

### Server

The Server model includes all fields from pb.ServerInfo in the Swagger spec to ensure complete alignment with the documented API.

**Naming Convention**: All Go structs use Request/Response suffix (e.g., ServerCreateRequest, ServerActionResponse) for consistency, even though Swagger uses Input/Output naming.

## Assumptions

- Documentation and Swagger spec are the source of truth for server model and module definitions.
- Reasonable defaults are used for unspecified details (e.g., error handling, timeouts).
- Only project-scoped endpoints are in scope.

## Edge Cases

- Invalid or missing fields in requests return structured errors.
- Unsupported server actions return errors.
- Asynchronous operations (202 Accepted) require polling for completion.
- Network or client-side errors are handled per the error contract.

## Dependencies

- Documentation and Swagger spec for server definitions.
- Existing SDK infrastructure for HTTP, error handling, and testing.
