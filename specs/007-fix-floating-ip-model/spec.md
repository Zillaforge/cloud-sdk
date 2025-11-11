# Feature Specification: Fix Floating IP Model

**Feature Branch**: `007-fix-floating-ip-model`  
**Created**: November 11, 2025  
**Status**: Draft  
**Input**: User description: "修正 floating ip 的 model 與 module 符合文件定義。spec 編號為 007"

## Clarifications

### Session 2025-11-11

- Q: How should existing SDK users handle the breaking change from "items" to "floatingips" array in the list response? → A: Immediate breaking change with migration guide only (remove "items" field entirely in next release)
- Q: What format should the SDK expect for timestamp fields (createdAt, updatedAt, approvedAt)? → A: ISO 8601 / RFC3339 string format (e.g., "2025-11-11T10:30:00Z")
- Q: How should the SDK handle floating IPs with null or missing project/user objects in API responses? → A: Leave project/user objects as null/empty when not provided by API (preserve API response structure)
- Q: Should the SDK provide helper methods to check floating IP status or readiness states? → A: No helper methods, expose raw status field only (developers check status string directly)
- Q: How should the SDK handle the "reserved" field behavior during create operations? → A: Read-only field; ignore if provided in create request, only returned by API in responses

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods;
  callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Correct FloatingIP Model Structure (Priority: P1)

As a developer using the cloud SDK, I want the FloatingIP model to accurately reflect the API specification (pb.FloatingIPInfo) so that I can properly handle floating IP data from VPS API responses including all metadata fields.

**Why this priority**: This is the core functionality needed for the SDK to work correctly with floating IP operations. Without accurate models, developers cannot reliably access all floating IP information.

**Independent Test**: Can be tested by verifying that FloatingIP structs can be unmarshaled from API responses containing all fields and marshaled back without data loss.

**Acceptance Scenarios**:

1. **Given** a VPS API response containing floating IP data with all fields including timestamps, project, user, device info, **When** unmarshaling into FloatingIP struct, **Then** all fields are correctly populated
2. **Given** a FloatingIP struct with all fields set, **When** marshaling to JSON, **Then** the output matches the API specification format with camelCase field names
3. **Given** a floating IP response from operations like Create or Get, **When** accessing fields like approvedAt, device_id, reserved, **Then** the values are available for application logic

---

### User Story 2 - Floating IP List Response Structure (Priority: P2)

As a developer listing floating IPs, I want the list response to match the API specification (pb.FIPListOutput) so that I can correctly iterate through all floating IPs without custom wrappers.

**Why this priority**: The current ListResponse uses a custom "items" wrapper not present in the API specification, causing potential confusion.

**Independent Test**: Can be tested by calling the List operation and verifying the response structure matches pb.FIPListOutput.

**Acceptance Scenarios**:

1. **Given** a list of floating IPs in the project, **When** calling List(), **Then** the response contains a floatingips array matching pb.FIPListOutput
2. **Given** an empty floating IPs list, **When** calling List(), **Then** the response contains an empty array without errors

---

### User Story 3 - Floating IP Create/Update Requests (Priority: P3)

As a developer creating or updating floating IPs, I want the request structures to match the API specification (FIPCreateInput/FIPUpdateInput) so that I can provide name and description fields.

**Why this priority**: The current request structures are missing the name field required by the API.

**Independent Test**: Can be tested by verifying create and update requests include name and description fields as per API spec.

**Acceptance Scenarios**:

1. **Given** create request with name and description, **When** calling Create(), **Then** the request matches FIPCreateInput structure
2. **Given** update request with name and description, **When** calling Update(), **Then** the request matches FIPUpdateInput structure

---

### Edge Cases

- What happens when optional fields (description, name, device info) are missing in API response? (Model preserves API response exactly)
- How does system handle floating IPs with null project or user objects? (Objects remain null/empty, preserving API response structure; project_id and user_id fields still available)
- What happens when timestamp strings cannot be parsed as RFC3339 format? (Return deserialization error with field details)
- How does system handle reserved vs non-reserved floating IPs? (Reserved field indicates special allocation status; read-only, not settable during creation)
- What happens when device_type or device_id are provided for associated floating IPs? (These fields indicate which resource the floating IP is attached to)
- What happens if a user attempts to set reserved field in create request? (SDK ignores the field; only API controls reservation status)

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The FloatingIP model MUST include all fields defined in pb.FloatingIPInfo from vps.yaml: address, approvedAt, createdAt, description, device_id, device_name, device_type, extnet_id, id, name, namespace, port_id, project, project_id, reserved, status, status_reason, updatedAt, user, user_id, uuid
- **FR-002**: Field names MUST match the API specification naming conventions (camelCase)
- **FR-003**: The FloatingIP model MUST properly serialize to and deserialize from the API JSON format
- **FR-004**: Project and user information MUST be represented as nested IDName objects when provided by the API; these objects remain null/empty when not provided, preserving the API response structure
- **FR-005**: Timestamp fields (createdAt, updatedAt, approvedAt) MUST accurately represent times from the API in ISO 8601 / RFC3339 string format (e.g., "2025-11-11T10:30:00Z")
- **FR-006**: Device-related fields (device_id, device_name, device_type) MUST be included for associated floating IPs
- **FR-007**: Reserved field MUST indicate whether the floating IP is reserved for special use; this is a read-only field returned by the API and not included in create requests
- **FR-007a**: The SDK MUST expose the raw status field as a string without helper methods; status interpretation remains the responsibility of application code
- **FR-008**: FloatingIPListResponse MUST match pb.FIPListOutput structure with floatingips array (breaking change: "items" field removed immediately)
- **FR-009**: FloatingIPCreateRequest and FloatingIPUpdateRequest MUST match FIPCreateInput and FIPUpdateInput from the Swagger specification including name and description fields (reserved field excluded from create/update requests)
- **FR-010**: All optional fields MUST be properly handled when absent from API responses
- **FR-011**: When timestamp fields contain malformed values, the SDK MUST return a deserialization error identifying the specific field
- **FR-012**: Migration guide MUST be provided in release notes documenting the breaking change from "items" to "floatingips" array field and showing how to update client code

### SDK Contract Requirements (Go)
```go
// Public methods follow the pattern: Client.Resource.Operation(ctx, params)
type Client struct {
    // ...
}

// Responses are strongly-typed structs matching Swagger models
type FloatingIP struct {
    // All fields from pb.FloatingIPInfo
}

// Errors wrap context (status code, service error fields) without exposing raw HTTP
```

## Key Entities *(include if feature involves data)*

- **FloatingIP**: Represents a floating IP address with attributes like id, address, status, project, user, device info, timestamps, and metadata
- **IDName**: Contains identification (id) and display name for project and user references
- **FloatingIPCreateRequest**: Input for creating a new floating IP with name and description
- **FloatingIPUpdateRequest**: Input for updating floating IP name and description
- **FloatingIPListResponse**: Output containing array of floating IPs matching the API specification

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: FloatingIP model successfully represents all required fields without data loss
- **SC-002**: All fields from pb.FloatingIPInfo are accessible in the FloatingIP model
- **SC-003**: FloatingIP data can be serialized to and deserialized from API JSON format without errors
- **SC-004**: FloatingIPListResponse matches pb.FIPListOutput structure
- **SC-005**: Create and update requests include name and description fields as per API specification
- **SC-006**: Release notes include migration guide showing how to adapt to corrected model structure
- **SC-007**: All existing integration tests pass or are updated to reflect the corrected model structure