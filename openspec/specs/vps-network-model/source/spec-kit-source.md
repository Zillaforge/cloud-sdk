# Feature Specification: Fix Network Model Definition

**Feature Branch**: `002-fix-network-model`  
**Created**: October 31, 2025  
**Status**: Draft  
**Input**: User description: "fix network model definition not match swagger file"

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods;
  callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - SDK Users Access Complete Network Information (Priority: P1)

Developers using the cloud-sdk need to access complete network information returned by the VPS API, including all fields defined in the Swagger specification such as gateway, nameservers, bonding status, router details, project/user references, and network state information.

**Why this priority**: This is the most critical fix because missing fields in the Network model prevent SDK users from accessing essential network configuration data returned by the API. Without these fields, users cannot properly manage networks, understand their configuration, or make informed decisions about network operations.

**Independent Test**: Can be fully tested by retrieving a network via the SDK and verifying that all fields from the Swagger specification are accessible and correctly populated in the returned Network struct.

**Acceptance Scenarios**:

1. **Given** a network exists in the VPS system with gateway configuration, **When** an SDK user retrieves the network details, **Then** the gateway address is accessible in the Network struct
2. **Given** a network has nameservers configured, **When** an SDK user retrieves the network, **Then** the nameservers array is populated with all configured DNS servers
3. **Given** a network is associated with a router, **When** an SDK user retrieves the network, **Then** the router information (ID and details) is accessible
4. **Given** a network has project and user ownership, **When** an SDK user retrieves the network, **Then** both project and user details (ID and name) are accessible
5. **Given** a network has status and bonding configuration, **When** an SDK user retrieves the network, **Then** the status, status_reason, bonding, and shared flags are accessible

---

### User Story 2 - SDK Users Create Networks with Full Configuration Options (Priority: P2)

Developers need to create networks with all supported configuration options defined in the Swagger specification, including gateway address and router association, not just the basic name, description, and CIDR.

**Why this priority**: This is important for advanced network configuration scenarios where users need to specify custom gateway addresses or associate networks with specific routers during creation. While basic network creation works without these fields, enterprise users require this functionality for complex network topologies.

**Independent Test**: Can be tested independently by creating a network with gateway and router_id specified, then verifying the network is created with the correct configuration.

**Acceptance Scenarios**:

1. **Given** valid network creation parameters including a custom gateway, **When** an SDK user creates a network, **Then** the network is created with the specified gateway address
2. **Given** valid network creation parameters including a router_id, **When** an SDK user creates a network, **Then** the network is created and associated with the specified router
3. **Given** network creation parameters with only required fields (name and CIDR), **When** an SDK user creates a network, **Then** the network is created successfully with default gateway

---

### User Story 3 - SDK Maintains Type Safety and Consistency (Priority: P3)

SDK users rely on strongly-typed models that accurately reflect the API contract, ensuring compile-time type safety and preventing runtime errors from accessing non-existent fields.

**Why this priority**: While this is a development quality concern rather than a feature delivery blocker, maintaining accurate type definitions is essential for a production-grade SDK. It prevents runtime errors and improves developer experience.

**Independent Test**: Can be tested by running all existing SDK tests to ensure no breaking changes were introduced, and that the expanded model still works with all existing code paths.

**Acceptance Scenarios**:

1. **Given** the updated Network model, **When** existing SDK code retrieves networks, **Then** all existing field access patterns continue to work without modification
2. **Given** the updated NetworkCreateRequest model, **When** existing SDK code creates networks, **Then** network creation succeeds with backward compatibility
3. **Given** JSON responses from the API, **When** the SDK unmarshals them into Network structs, **Then** all fields are correctly populated including newly added ones

---

### Edge Cases

- What happens when optional fields (like gateway, router, nameservers) are not present in the API response? The struct should handle omitempty tags appropriately and allow nil/zero values.
- How does the SDK handle deprecated fields like `gw_state` that are marked as deprecated in the Swagger spec? They should be included but documented as deprecated.
- What happens when nested objects (project, user, router) are returned with only partial data? The nested IDName and RouterInfo structs should handle optional fields gracefully.
- How does the system handle arrays that might be empty (nameservers)? The Go model should support empty arrays without errors.

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: Network model MUST include all fields defined in the pb.NetworkInfo Swagger definition including: bonding, gateway, gw_state, is_default, nameservers, namespace, project, router, router_id, shared, status, status_reason, subnet_id, user, and user_id
- **FR-002**: NetworkCreateRequest MUST include the gateway and router_id optional fields as defined in the NetCreateInput Swagger definition
- **FR-003**: Network model MUST properly handle nested objects (project, user, router) using appropriate struct types with ID and name fields
- **FR-004**: All timestamp fields (createdAt, updatedAt) MUST maintain consistent naming with the Swagger camelCase convention
- **FR-005**: Network model MUST use appropriate Go types for each field (string for IDs/text, bool for flags, []string for arrays, nested structs for objects)
- **FR-006**: JSON struct tags MUST exactly match the field names in the Swagger specification for proper serialization/deserialization
- **FR-007**: Optional fields MUST be marked with `omitempty` JSON tags to allow proper handling of absent data
- **FR-008**: NetworkListResponse and ListNetworksOptions MUST remain compatible with existing code while supporting the expanded Network model
- **FR-009**: NetworkUpdateRequest MUST remain unchanged as it already matches the NetUpdateInput Swagger definition (name and description only)

### SDK Contract Requirements (Go)

- Public methods follow the pattern: `Client.Resource.Operation(ctx, params)`.
- Responses are strongly-typed structs matching Swagger models.
- Errors wrap context (status code, service error fields) without exposing raw HTTP.
- Authentication, retries, and timeouts are centralized in the Client.
- Breaking changes MUST be documented with migration guidance (adding fields is generally non-breaking, but changing existing field types would be breaking).

### Key Entities

- **Network**: Represents a virtual network in the VPS infrastructure with complete configuration including CIDR, gateway, router association, ownership information, and operational state
- **IDName**: Simple reference object containing ID and name fields for related entities (projects, users)
- **RouterInfo**: Nested object containing router details associated with a network (structure defined in router models)
- **NetworkCreateRequest**: Input model for creating new networks with name, description, CIDR, optional gateway, and optional router association
- **NetworkUpdateRequest**: Input model for updating existing networks (name and description only)

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: All 19 fields from the pb.NetworkInfo Swagger definition are accessible in the SDK's Network struct
- **SC-002**: Network creation supports all 5 fields defined in the NetCreateInput Swagger specification (name, description, cidr, gateway, router_id)
- **SC-003**: All existing SDK tests pass without modification, demonstrating backward compatibility
- **SC-004**: SDK users can successfully access previously unavailable fields (gateway, nameservers, router details, project/user info) when retrieving network information
- **SC-005**: JSON serialization and deserialization of Network objects works correctly for all fields including nested objects and arrays

## Assumptions

- The Swagger specification at `/workspaces/cloud-sdk/swagger/vps.json` is the authoritative source of truth for the API contract
- The `pb.NetworkInfo` definition represents the complete response structure for network operations
- Nested types like `pb.IDName` and `pb.RouterInfo` are already defined or will be defined in their respective model files
- Adding new fields to existing structs is non-breaking for SDK consumers (they can ignore fields they don't use)
- The API consistently returns all fields as defined in the Swagger spec, or omits optional fields when not applicable
- Timestamp fields use string type as defined in Swagger (not time.Time) to match API serialization format
- The deprecated `gw_state` field should be included for completeness but marked as deprecated in comments
