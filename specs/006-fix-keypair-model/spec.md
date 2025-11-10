# Feature Specification: Fix Keypair Model

**Feature Branch**: `006-fix-keypair-model`  
**Created**: November 10, 2025  
**Status**: Draft  
**Input**: User description: "修正 keypair 的 model 與 module 符合文件定義"

## Clarifications

### Session 2025-11-10

- Q: How should existing SDK users migrate when removing the Total field from KeypairListResponse? → A: Immediate removal with migration guide in release notes
- Q: What format should the SDK expect for timestamp fields (createdAt, updatedAt)? → A: ISO 8601 / RFC3339 string format (e.g., "2025-11-10T15:30:00Z")
- Q: How should the SDK handle user object when it's null or missing from API response? → A: Leave user object as null/empty when not provided by API (keep result consistent with API response)
- Q: Should the SDK provide any guidance or warnings about handling the private_key field? → A: Document that private_key is sensitive and should be saved immediately by caller
- Q: How should the SDK handle timestamp strings that cannot be parsed as RFC3339 format? → A: Return deserialization error with details about which field failed

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods;
  callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Correct Keypair Model Structure (Priority: P1)

As a developer using the cloud SDK, I want the Keypair model to accurately reflect the API specification (pb.KeypairInfo) so that I can properly handle keypair data from VPS API responses including timestamps and user information.

**Why this priority**: This is the core functionality needed for the SDK to work correctly with keypair operations. Without accurate models, developers cannot reliably access all keypair information.

**Independent Test**: Can be tested by verifying that Keypair structs can be unmarshaled from API responses containing all fields (including createdAt, updatedAt, user, private_key) and marshaled back without data loss.

**Acceptance Scenarios**:

1. **Given** a VPS API response containing keypair data with all fields including timestamps, **When** unmarshaling into Keypair struct, **Then** all fields (id, name, description, fingerprint, public_key, private_key, user_id, user, createdAt, updatedAt) are correctly populated
2. **Given** a Keypair struct with all fields set, **When** marshaling to JSON, **Then** the output matches the API specification format with camelCase field names
3. **Given** a keypair response from the Create operation, **When** accessing the private_key field, **Then** the generated private key is available and SDK documentation warns it must be saved immediately

---

### User Story 2 - Keypair List Response Structure (Priority: P2)

As a developer listing keypairs, I want the list response to match the API specification (pb.KeypairListOutput) so that I can correctly iterate through all keypairs without relying on a custom "total" field.

**Why this priority**: The current ListResponse includes a "total" field not present in the API specification, causing potential confusion and incorrect expectations.

**Independent Test**: Can be tested by calling the List operation and verifying the response structure matches pb.KeypairListOutput with only the keypairs array.

**Acceptance Scenarios**:

1. **Given** a list of keypairs in the project, **When** calling List(), **Then** the response contains a keypairs array matching pb.KeypairListOutput
2. **Given** an empty keypairs list, **When** calling List(), **Then** the response contains an empty array without errors

---

### User Story 3 - Timestamp Handling (Priority: P3)

As a developer needing to track keypair lifecycle, I want the Keypair model to include creation and update timestamps so that I can monitor when keypairs were created or modified for security auditing.

**Why this priority**: Timestamps provide metadata for auditing and debugging but are not required for basic keypair operations.

**Independent Test**: Can be tested by verifying timestamp fields are properly handled in JSON operations and can be parsed into Go time.Time types.

**Acceptance Scenarios**:

1. **Given** a keypair with timestamps from the API, **When** unmarshaling, **Then** createdAt and updatedAt are parsed as string timestamps
2. **Given** a newly created keypair, **When** retrieving it, **Then** createdAt reflects the creation time

---

### Edge Cases

- What happens when optional fields (description, private_key) are missing in API response? (Model preserves API response exactly)
- How does system handle keypairs when user object is null/missing? (Preserved as null, user_id still available)
- What happens when timestamp strings cannot be parsed as RFC3339 format? (Return deserialization error with field details)
- What happens when public_key is provided during creation vs. generated? (Private key only returned when generated)
- What happens when very old or future timestamps exceed valid date ranges?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Keypair model MUST include all fields defined in pb.KeypairInfo from vps.yaml: id, name, description, fingerprint, public_key, private_key, user_id, user, createdAt, updatedAt
- **FR-002**: Field names MUST match the API specification naming conventions
- **FR-003**: The Keypair model MUST properly serialize to and deserialize from the API JSON format
- **FR-004**: User information MUST be represented as a nested object with id and name attributes, remaining null/empty when not provided by the API
- **FR-005**: Timestamp fields MUST accurately represent creation and update times from the API in ISO 8601 / RFC3339 format (e.g., "2025-11-10T15:30:00Z")
- **FR-006**: The PrivateKey field MUST be accessible when a new keypair is generated by the API (only returned during Create operation)
- **FR-007**: KeypairListResponse MUST match pb.KeypairListOutput structure with only the keypairs array (Total field removed immediately)
- **FR-011**: SDK documentation MUST warn that private_key is sensitive and returned only once during creation, requiring immediate secure storage by the caller
- **FR-008**: All optional fields (description, private_key) MUST be properly handled when absent from API responses
- **FR-009**: KeypairCreateRequest and KeypairUpdateRequest MUST match KeypairCreateInput and KeypairUpdateInput from the Swagger specification
- **FR-012**: When timestamp fields contain malformed or unparseable values, the SDK MUST return a deserialization error identifying the specific field that failed
- **FR-010**: Migration guide MUST be provided in release notes documenting the removal of Total field and how to calculate count from array length

### SDK Contract Requirements (Go)

- Public methods follow the pattern: `Client.Resource.Operation(ctx, params)`.
- Responses are strongly-typed structs matching Swagger models.
- Errors wrap context (status code, service error fields) without exposing raw HTTP.
- Authentication, retries, and timeouts are centralized in the Client.

### Key Entities *(include if feature involves data)*

- **Keypair**: Represents an SSH keypair for server access with attributes like id, name, description, fingerprint, public_key, private_key, user_id, user object, and timestamps (createdAt, updatedAt)
- **User Reference**: Contains user identification (id and name) referencing the keypair owner
- **KeypairCreateRequest**: Input for creating a new keypair with name, optional description, and optional public_key (if omitted, generates new key pair)
- **KeypairUpdateRequest**: Input for updating keypair description only
- **KeypairListResponse**: Output containing array of keypairs matching the API specification

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Keypair model successfully represents all required fields without data loss
- **SC-002**: All fields from pb.KeypairInfo are accessible in the Keypair model
- **SC-003**: Keypair data can be serialized to and deserialized from API JSON format without errors
- **SC-004**: KeypairListResponse matches pb.KeypairListOutput structure (Total field removed)
- **SC-005**: Private key from Create operation can be retrieved and saved by developers
- **SC-006**: Release notes include migration guide showing how to replace Total field usage with len(response.Keypairs)
- **SC-007**: All existing integration tests pass or are updated to reflect the corrected model structure
