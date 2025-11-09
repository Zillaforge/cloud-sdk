# Feature Specification: Fix Security Group Model

**Feature Branch**: `001-fix-security-group-model`  
**Created**: November 9, 2025  
**Status**: Draft  
**Input**: User description: "修正 security group 的model，符合 API文件定義"

**Additional Requirements**:
- Create/Delete/Update/Get operations must conform to API file definition
- Design security group sub-resource `Rule` operations

## Clarifications

### Session 2025-11-09

- Q: What authentication and authorization requirements should be enforced for security group CRUD operations? → A: Operations require valid API token with project-level permissions; tenant admins can manage all groups, regular users only their own
- Q: What performance targets should be set for security group operations? → A: No specific performance requirements - rely on underlying API performance
- Q: What error handling expectations should be set for security group operations? → A: Use existing custom sdkError type consistent with network and flavor APIs
- Q: How should concurrent operations on security groups be handled? → A: Last-write-wins with optimistic locking (API-level conflict detection)
- Q: What scale assumptions should be made for security groups and rules? → A: Standard cloud service assumptions (hundreds of security groups per project, thousands of rules total)

## Constitution Constraints (must reflect in requirements & tests)

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods;
  callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Correct Security Group List Response Structure (Priority: P1)

As a developer using the cloud SDK, I want the SecurityGroupListResponse to accurately reflect the API specification so that I can properly handle list responses from VPS API.

**Why this priority**: This is the core functionality needed for the SDK to work correctly with security group list operations.

**Independent Test**: Can be tested by verifying that SecurityGroupListResponse structs can be unmarshaled from API responses without extra fields.

**Acceptance Scenarios**:

1. **Given** a VPS API list response containing security groups, **When** unmarshaling into SecurityGroupListResponse struct, **Then** all fields are correctly populated
2. **Given** a SecurityGroupListResponse struct with all fields set, **When** marshaling to JSON, **Then** the output matches the API specification format without extra fields

---

### User Story 2 - Accurate Security Group Model Fields (Priority: P2)

As a developer working with security group data, I want the SecurityGroup and SecurityGroupRule structs to match the API specification exactly so that I can access all available fields correctly.

**Why this priority**: Ensures complete compatibility with the VPS API data structures.

**Independent Test**: Can be tested by creating SecurityGroup instances and verifying serialization matches API format.

**Acceptance Scenarios**:

1. **Given** a security group with all fields, **When** accessing fields, **Then** all API-defined fields are available
2. **Given** security group rule data, **When** unmarshaling, **Then** port_min and port_max are handled correctly

---

### User Story 3 - Consistent JSON Serialization (Priority: P3)

As a developer integrating with the VPS API, I want consistent JSON field naming between request/response models and the API specification so that requests and responses are properly formatted.

**Why this priority**: Prevents serialization issues that could cause API failures.

**Independent Test**: Can be tested by marshaling/unmarshaling models and comparing with API examples.

**Acceptance Scenarios**:

1. **Given** create/update request structs, **When** marshaling, **Then** JSON field names match API specification
2. **Given** API response JSON, **When** unmarshaling into response structs, **Then** all data is correctly populated

---

### User Story 4 - Complete Security Group CRUD Operations (Priority: P1)

As a developer managing security groups, I want full Create, Read, Update, Delete (CRUD) operations that match the API specification so that I can perform all necessary security group management tasks.

**Why this priority**: Core functionality for security group lifecycle management.

**Independent Test**: Can be tested by performing each CRUD operation and verifying API compliance.

**Acceptance Scenarios**:

1. **Given** security group creation parameters, **When** calling create, **Then** the request matches SgCreateInput format
2. **Given** a security group ID, **When** calling get, **Then** the response matches pb.SgInfo format
3. **Given** security group update parameters, **When** calling update, **Then** the request matches SgUpdateInput format
4. **Given** a security group ID, **When** calling delete, **Then** the operation succeeds without errors

---

### User Story 5 - Security Group Listing with Filters (Priority: P2)

As a developer listing security groups, I want to use query parameters that match the API specification so that I can filter results appropriately.

**Why this priority**: Essential for managing large numbers of security groups.

**Independent Test**: Can be tested by listing with different filter combinations.

**Acceptance Scenarios**:

1. **Given** filter parameters (name, user_id, detail), **When** calling list, **Then** the query parameters match API specification
2. **Given** detail=true, **When** calling list, **Then** response includes rules in each security group

---

### User Story 6 - Security Group Rule Management (Priority: P1)

As a developer managing security group rules, I want to create and delete rules as a sub-resource so that I can control network access policies.

**Why this priority**: Rules are the core functionality of security groups for network security.

**Independent Test**: Can be tested by creating and deleting rules independently.

**Acceptance Scenarios**:

1. **Given** rule creation parameters, **When** calling create rule, **Then** the request matches SgRuleCreateInput format
2. **Given** a rule ID, **When** calling delete rule, **Then** the rule is removed from the security group
3. **Given** rule parameters with optional ports, **When** creating rule, **Then** port_min and port_max are handled correctly

---

### Edge Cases

- What happens when optional fields like port_min/port_max are missing in API response?
- How does system handle security groups with no rules?
- What if the list response contains empty security_groups array?
- What happens when creating a security group with invalid rule parameters?
- How does system handle concurrent rule operations on the same security group?
- What if deleting a security group that still has rules?
- What happens when listing with invalid filter parameters?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The SecurityGroupListResponse struct MUST include only the fields defined in pb.SgListOutput from vps.yaml
- **FR-002**: Remove the Total field from SecurityGroupListResponse as it is not present in the API specification
- **FR-003**: Field names MUST use PascalCase for Go struct fields
- **FR-004**: JSON tags MUST match the camelCase field names from the API specification
- **FR-005**: SecurityGroup and SecurityGroupRule structs MUST match pb.SgInfo and pb.SgRuleInfo exactly
- **FR-006**: Request structs (create/update) MUST match the corresponding Input definitions in the API
- **FR-007**: ListSecurityGroupsOptions MUST include all query parameters defined in the API (name, user_id, detail)
- **FR-008**: Security group CRUD operations MUST match the API endpoints and request/response formats
- **FR-009**: Security group rule sub-resource operations (create/delete) MUST match the API specification
- **FR-010**: Rule operations MUST be scoped under security group resources following REST conventions
- **FR-011**: Security group operations MUST enforce proper authentication and authorization: require valid API token in request headers, return 401 Unauthorized for missing/invalid tokens, return 403 Forbidden for insufficient project-level permissions (tenant admins can manage all groups, regular users only their own)
- **FR-012**: Error handling MUST use existing custom sdkError type from errors.go (fields: StatusCode int, Message string, Details string) consistent with network and flavor APIs
- **FR-013**: Concurrent operations MUST use last-write-wins with optimistic locking (API-level conflict detection)
- **FR-014**: Implementation MUST support standard cloud service scale assumptions (hundreds of security groups per project, thousands of rules total)
- **FR-015**: Protocol and Direction fields MUST use custom Go types (not plain strings) with predefined constants for type safety (ProtocolTCP/UDP/ICMP/Any, DirectionIngress/Egress)

### SDK Contract Requirements (Go)

- Public methods follow the pattern: `Client.Resource.Operation(ctx, params)`.
- Responses are strongly-typed structs matching Swagger models.
- Errors wrap context (status code, service error fields) without exposing raw HTTP.
- Authentication, retries, and timeouts are centralized in the Client.

### Key Entities *(include if feature involves data)*

- **SecurityGroup**: Represents a security group in the VPS service with ID, name, description, project/user info, rules, and timestamps
- **SecurityGroupRule**: Represents individual rules within a security group with direction, protocol, ports, and remote CIDR
- **SecurityGroupListResponse**: Contains the list of security groups returned by list operations
- **ListSecurityGroupsOptions**: Query parameters for filtering security group lists (name, user_id, detail)
- **SecurityGroupCreateRequest**: Request payload for creating new security groups
- **SecurityGroupUpdateRequest**: Request payload for updating existing security groups
- **SecurityGroupRuleCreateRequest**: Request payload for creating new security group rules

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: SecurityGroupListResponse struct compiles without errors in Go and matches pb.SgListOutput
- **SC-002**: All fields from pb.SgInfo and pb.SgRuleInfo are present and correctly typed in model structs
- **SC-003**: JSON marshaling/unmarshaling works correctly for all field types without extra fields
- **SC-004**: Request structs match API input specifications for create and update operations
- **SC-005**: Migration notes MUST be provided for the removal of the Total field from SecurityGroupListResponse
- **SC-006**: CRUD operations for security groups match API endpoint specifications
- **SC-007**: Security group rule create/delete operations match API specifications
- **SC-008**: ListSecurityGroupsOptions includes all API-defined query parameters
- **SC-009**: All operations handle optional parameters correctly according to API specification
- **SC-010**: Protocol and Direction custom types with constants are implemented and used consistently across all rule models
