# Data Model: Keypair

**Feature**: 006-fix-keypair-model  
**Date**: 2025-11-10  
**Phase**: 1 - Design & Contracts

## Entity: Keypair

**Purpose**: Represents an SSH keypair for server access in the VPS service

**Source**: Swagger `pb.KeypairInfo` (vps.yaml:915-937)

### Fields

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | string | `id` | Yes | Unique identifier for the keypair |
| Name | string | `name` | Yes | User-defined name for the keypair |
| Description | string | `description` | No | Optional description (omitempty) |
| Fingerprint | string | `fingerprint` | Yes | SSH public key fingerprint |
| PublicKey | string | `public_key` | Yes | SSH public key content |
| PrivateKey | string | `private_key` | No | SSH private key (only in Create response, omitempty) |
| UserID | string | `user_id` | Yes | ID of the user who owns this keypair |
| User | *IDName | `user` | No | User reference object (may be null, omitempty) |
| CreatedAt | string | `createdAt` | Yes | ISO 8601 timestamp of creation |
| UpdatedAt | string | `updatedAt` | Yes | ISO 8601 timestamp of last update |

### Field Notes

**PrivateKey**:
- Only returned during Create operation when keypair is generated
- Never returned in List or Get operations
- SENSITIVE DATA - must be saved immediately by caller
- omitempty tag prevents inclusion when empty

**User**:
- Pointer type `*IDName` to handle null values from API
- May be null even when UserID is present
- SDK preserves API response exactly (no population from UserID)

**Timestamps**:
- String type matching Swagger `type: string`
- Format: ISO 8601 / RFC3339 (e.g., "2025-11-10T15:30:00Z")
- Callers parse with `time.Parse(time.RFC3339, timestamp)` if needed

## Entity: IDName

**Purpose**: Reusable reference type for linked resources

**Source**: Swagger `pb.IDName` (vps.yaml:908-913)

### Fields

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| ID | string | `id` | Yes | Resource identifier |
| Name | string | `name` | Yes | Resource display name |

**Usage**: Referenced by Keypair.User, Server.Image, Server.Flavor, etc.

## Request Models

### KeypairCreateRequest

**Purpose**: Request input for creating a new keypair

**Source**: Swagger `KeypairCreateInput` (vps.yaml:218-228)  
**SDK Convention**: Uses Request suffix for input parameters

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Name | string | `name` | Yes | Name for the new keypair |
| Description | string | `description` | No | Optional description (omitempty) |
| PublicKey | string | `public_key` | No | Import existing public key; omit to generate new pair (omitempty) |

**Behavior**:
- If `PublicKey` provided: Import existing key, no PrivateKey in response
- If `PublicKey` omitted: Generate new pair, PrivateKey returned in response

### KeypairUpdateRequest

**Purpose**: Request input for updating keypair description

**Source**: Swagger `KeypairUpdateInput` (vps.yaml:229-232)  
**SDK Convention**: Uses Request suffix for input parameters

| Field | Type | JSON Tag | Required | Description |
|-------|------|----------|----------|-------------|
| Description | string | `description` | No | New description (omitempty) |

**Note**: Only description is mutable; other fields are immutable

## Response Models

### KeypairListResponse

**Purpose**: Response struct for unmarshaling list API response

**Source**: Swagger `pb.KeypairListOutput` (vps.yaml:938-943)  
**SDK Convention**: Uses Response suffix for output structures

| Field | Type | JSON Tag | Description |
|-------|------|----------|-------------|
| Keypairs | []*Keypair | `keypairs` | Array of keypair objects |

**Breaking Change**: Removed `Total` field (not in Swagger spec)

**Client Method Signature**:
```go
func (c *Client) List(ctx context.Context, opts *ListKeypairsOptions) ([]*Keypair, error)
```

Returns slice directly, not wrapped struct. Callers use `len(result)` for count.

### ListKeypairsOptions

**Purpose**: Optional filters for List operation

| Field | Type | Description |
|-------|------|-------------|
| Name | string | Filter by keypair name (optional) |

**Usage**: Query parameter `?name=value` when provided

## Relationships

```text
Keypair
├── User (IDName reference, nullable)
│   ├── ID (user_id) - always present
│   └── Name - may be absent if user object null
└── Timestamps
    ├── CreatedAt (ISO 8601 string)
    └── UpdatedAt (ISO 8601 string)
```

## State Transitions

Keypairs have no explicit state field. Lifecycle:

1. **Create**: POST /keypairs → Keypair (with PrivateKey if generated)
2. **Active**: GET /keypairs/{id} → Keypair (no PrivateKey)
3. **Update**: PUT /keypairs/{id} → Keypair (description only)
4. **Delete**: DELETE /keypairs/{id} → No content

## Validation Rules

From Swagger `required` fields:
- Create: `name` required
- Update: All fields optional (description-only update)
- Response: `id`, `name`, `fingerprint`, `public_key`, `user_id`, `createdAt`, `updatedAt` always present

## JSON Examples

### Create Request (Generate New)
```json
{
  "name": "my-keypair",
  "description": "Development keypair"
}
```

### Create Response (Generated)
```json
{
  "id": "kp-abc123",
  "name": "my-keypair",
  "description": "Development keypair",
  "fingerprint": "SHA256:abc...xyz",
  "public_key": "ssh-rsa AAAAB3...",
  "private_key": "-----BEGIN RSA PRIVATE KEY-----\n...",
  "user_id": "user-123",
  "user": {
    "id": "user-123",
    "name": "john@example.com"
  },
  "createdAt": "2025-11-10T15:30:00Z",
  "updatedAt": "2025-11-10T15:30:00Z"
}
```

### Get Response (No Private Key)
```json
{
  "id": "kp-abc123",
  "name": "my-keypair",
  "description": "Development keypair",
  "fingerprint": "SHA256:abc...xyz",
  "public_key": "ssh-rsa AAAAB3...",
  "user_id": "user-123",
  "user": {
    "id": "user-123",
    "name": "john@example.com"
  },
  "createdAt": "2025-11-10T15:30:00Z",
  "updatedAt": "2025-11-10T15:30:00Z"
}
```

### List Response
```json
{
  "keypairs": [
    {
      "id": "kp-abc123",
      "name": "my-keypair",
      "fingerprint": "SHA256:abc...xyz",
      "public_key": "ssh-rsa AAAAB3...",
      "user_id": "user-123",
      "createdAt": "2025-11-10T15:30:00Z",
      "updatedAt": "2025-11-10T15:30:00Z"
    }
  ]
}
```

## Migration Notes

### Breaking Changes

**Removed Field**: `KeypairListResponse.Total`

**Before (v0.x)**:
```go
response, err := client.Keypairs.List(ctx, nil)
if err != nil {
    // handle error
}
total := response.Total
for _, kp := range response.Keypairs {
    // process keypair
}
```

**After (v1.0)**:
```go
keypairs, err := client.Keypairs.List(ctx, nil)
if err != nil {
    // handle error
}
total := len(keypairs)  // Calculate from slice length
for _, kp := range keypairs {
    // process keypair
}
```

**Rationale**: Aligns with Swagger spec `pb.KeypairListOutput` which has no `total` field. Consistent with industry-standard REST patterns where clients calculate count from array length.
