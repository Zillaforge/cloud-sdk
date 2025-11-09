# Research: Fix Flavors Model

**Date**: November 8, 2025
**Feature**: 001-fix-flavors-model

## Research Tasks

### Task 1: pb.FlavorInfo Field Definitions
**Query**: What are the exact field names, types, and requirements in pb.FlavorInfo from vps.yaml?

**Findings**:
- id: string (required)
- name: string
- description: string
- vcpu: integer
- memory: integer (MiB)
- disk: integer (GiB)
- gpu: pb.GPUInfo (optional)
- public: boolean
- tags: []string
- project_ids: []string
- az: string
- createdAt: string (ISO 8601)
- updatedAt: string (ISO 8601)
- deletedAt: string (ISO 8601)

**Decision**: Map directly to Go struct with appropriate types.
**Rationale**: Ensures API compatibility.
**Alternatives**: Custom mapping - rejected for complexity.

### Task 2: Go JSON Handling for Optional Fields
**Query**: Best practices for optional fields in Go structs with JSON.

**Findings**:
- Use pointer types (*Type) for optional fields
- Use omitempty tag for fields that may be absent
- time.Time fields need custom marshaling for ISO 8601

**Decision**: Use *time.Time for timestamp fields, *GPUInfo for GPU.
**Rationale**: Proper JSON handling without zero values.
**Alternatives**: Custom types - rejected for simplicity.

### Task 3: URL Query Encoding for Multiple Tags
**Query**: How to encode multiple tag values in URL query parameters.

**Findings**:
- Use url.Values.Add() for multiple values with same key
- Results in tag=val1&tag=val2
- Go's url.Values handles this correctly

**Decision**: Use query.Add("tag", tag) in loop.
**Rationale**: Standard HTTP practice.
**Alternatives**: Single tag string - rejected for API compliance.

### Task 4: Breaking Changes Migration
**Query**: How to document migration for field name changes.

**Findings**:
- Document old vs new field names
- Provide code examples
- Note compilation errors

**Decision**: Include migration section in release notes.
**Rationale**: Helps users update code.
**Alternatives**: Deprecation warnings - rejected for Go limitations.</content>
<parameter name="filePath">/workspaces/cloud-sdk/specs/001-fix-flavors-model/research.md