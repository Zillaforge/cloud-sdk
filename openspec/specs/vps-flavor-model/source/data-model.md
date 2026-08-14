# Data Model: Fix Flavors Model

**Date**: November 8, 2025
**Feature**: 001-fix-flavors-model

## Entities

### Flavor
Represents a compute instance flavor/size in the VPS API.

**Fields**:
- `id` (string): Unique identifier
- `name` (string): Human-readable name
- `description` (string, optional): Description text
- `vcpu` (int): Number of virtual CPUs
- `memory` (int): RAM in MiB
- `disk` (int): Disk size in GiB
- `gpu` (*GPUInfo, optional): GPU configuration
- `public` (bool): Whether flavor is public
- `tags` ([]string): Associated tags
- `project_ids` ([]string): Restricted project IDs
- `az` (string, optional): Availability zone
- `createdAt` (*time.Time, optional): Creation timestamp
- `updatedAt` (*time.Time, optional): Last update timestamp
- `deletedAt` (*time.Time, optional): Deletion timestamp

**Relationships**:
- None (standalone entity)

**Validation Rules**:
- `id`: Required, non-empty
- `name`: Required, non-empty
- `vcpu`: >= 0
- `memory`: >= 0
- `disk`: >= 0
- `tags`: Array of strings
- `project_ids`: Array of strings
- Timestamps: Valid ISO 8601 when present

### GPUInfo
GPU configuration for flavors.

**Fields**:
- `count` (int): Number of GPUs
- `is_vgpu` (bool): Whether VGPU is used
- `model` (string): GPU model name

**Validation Rules**:
- `count`: >= 0
- `model`: Required when GPU present

### ListFlavorsOptions
Filtering options for flavor listing.

**Fields**:
- `name` (string): Filter by name
- `public` (*bool): Filter by visibility
- `tags` ([]string): Filter by tags (multiple allowed)
- `resize_server_id` (string): Filter flavors available for server resize

**Validation Rules**:
- `tags`: Array of non-empty strings

## State Transitions
N/A - Flavors are static configuration entities.</content>
<parameter name="filePath">/workspaces/cloud-sdk/specs/001-fix-flavors-model/data-model.md