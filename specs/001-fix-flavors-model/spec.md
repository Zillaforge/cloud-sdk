# Feature Specification: Fix Flavors Model

**Feature Branch**: `001-fix-flavors-model`  
**Created**: November 8, 2025  
**Status**: Draft  
**Input**: User description: "修正 `models/flavors` 裡的 model"

## Clarifications

### Session 2025-11-08

- Q: How to handle backward compatibility when updating field names and JSON tags to match API specification? → A: Break backward compatibility - Use new field names (VCPU, Memory) matching API, document migration in release notes

- SDK surface MUST be an idiomatic Go package: construct a Client and call methods;
  callers MUST NOT manage raw HTTP details.
- All public APIs MUST accept `context.Context` and return typed results and errors.
- Tests are written first (unit + contract tests per Swagger/OpenAPI) and MUST pass.
- External dependencies MUST be minimized and justified.
- Breaking changes MUST be called out with migration notes; use semantic versioning.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Correct Flavor Model Structure (Priority: P1)

As a developer using the cloud SDK, I want the Flavor model to accurately reflect the API specification so that I can properly handle flavor data from VPS API responses.

**Why this priority**: This is the core functionality needed for the SDK to work correctly with flavor operations.

**Independent Test**: Can be tested by verifying that Flavor structs can be unmarshaled from API responses and marshaled back without data loss.

**Acceptance Scenarios**:

1. **Given** a VPS API response containing flavor data, **When** unmarshaling into Flavor struct, **Then** all fields are correctly populated
2. **Given** a Flavor struct with all fields set, **When** marshaling to JSON, **Then** the output matches the API specification format

---

### User Story 2 - GPU Support in Flavors (Priority: P2)

As a developer working with GPU-enabled flavors, I want the Flavor model to include GPU information so that I can access GPU details for flavor selection.

**Why this priority**: GPU flavors are a specific use case that requires additional fields.

**Independent Test**: Can be tested by creating Flavor instances with GPU info and verifying serialization.

**Acceptance Scenarios**:

1. **Given** a flavor with GPU configuration, **When** accessing GPU field, **Then** count, model, and VGPU status are available

---

### User Story 3 - Flavor Timestamps (Priority: P3)

As a developer needing to track flavor lifecycle, I want the Flavor model to include creation and update timestamps so that I can monitor flavor changes.

**Why this priority**: Timestamps provide metadata for auditing and debugging.

**Independent Test**: Can be tested by verifying timestamp fields are properly handled in JSON operations.

**Acceptance Scenarios**:

1. **Given** a flavor with timestamps, **When** unmarshaling, **Then** timestamps are parsed as time.Time

---

### Edge Cases

- What happens when optional fields are missing in API response?
- How does system handle flavors with no GPU configuration?
- What if project_ids array is empty?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: The Flavor struct MUST include all fields defined in pb.FlavorInfo from vps.yaml
- **FR-002**: Field names MUST use PascalCase for Go struct fields
- **FR-003**: JSON tags MUST match the camelCase field names from the API specification
- **FR-004**: GPU information MUST be represented as a separate GPUInfo struct
- **FR-005**: Timestamp fields MUST be of type *time.Time for proper JSON handling
- **FR-006**: Array fields like tags and project_ids MUST be []string
- **FR-007**: Boolean fields like public MUST be bool type

### SDK Contract Requirements (Go)

- Public methods follow the pattern: `Client.Resource.Operation(ctx, params)`.
- Responses are strongly-typed structs matching Swagger models.
- Errors wrap context (status code, service error fields) without exposing raw HTTP.
- Authentication, retries, and timeouts are centralized in the Client.

### Key Entities *(include if feature involves data)*

- **Flavor**: Represents a compute instance flavor/size with attributes like VCPU, memory, disk, GPU, tags, project restrictions, and timestamps
- **GPUInfo**: Contains GPU configuration details including count, model, and VGPU status

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Flavor struct compiles without errors in Go
- **SC-002**: All fields from pb.FlavorInfo are present and correctly typed
- **SC-003**: JSON marshaling/unmarshaling works correctly for all field types
- **SC-004**: GPU flavors can be properly represented and accessed
- **SC-005**: Migration notes MUST be provided for breaking changes to field names (VCPUs → VCPU, RAM → Memory)
