# IAM API Contracts

**Feature**: 009-iam-api  
**Date**: 2025-11-12  
**Source**: `/workspaces/cloud-sdk/swagger/iam.yaml`

## Overview

This document defines the API contracts for the three IAM endpoints being implemented. All contracts are derived from the IAM Swagger specification and validated through contract tests.

## Authentication

All endpoints require Bearer Token authentication:

```
Authorization: Bearer {token}
```

## Base URL

```
http://127.0.0.1:8084/iam/api/v1
```

## Endpoints

### 1. GET /user

Retrieve information about the currently authenticated user.

#### Request

**Method**: `GET`  
**Path**: `/user`  
**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
Accept: application/json
```
**Query Parameters**: None  
**Request Body**: None

#### Response

**Status Code**: `200 OK`

**Response Body**:
```json
{
  "userId": "string (UUID)",
  "account": "string",
  "displayName": "string",
  "description": "string",
  "extra": {},
  "namespace": "string",
  "email": "string",
  "frozen": false,
  "mfa": false,
  "createdAt": "string (ISO8601)",
  "updatedAt": "string (ISO8601)",
  "lastLoginAt": "string (ISO8601)"
}
```

**Field Descriptions**:
- `userId`: Unique user ID (UUID format)
- `account`: User account identifier (typically email)
- `displayName`: Display name for UI
- `description`: User description
- `extra`: Arbitrary metadata object
- `namespace`: Tenant namespace (e.g., "ci.asus.com")
- `email`: User email address
- `frozen`: Whether account is frozen/disabled
- `mfa`: Whether MFA is enabled
- `createdAt`: User creation timestamp
- `updatedAt`: Last update timestamp
- `lastLoginAt`: Most recent login timestamp

#### Error Responses

**400 Bad Request**:
```json
{
  "errorCode": 1001,
  "message": "Bad Request"
}
```

**403 Forbidden** (invalid/expired token):
```json
{
  "errorCode": 1003,
  "message": "Forbidden"
}
```

**500 Internal Server Error**:
```json
{
  "errorCode": 2001,
  "message": "Internal Server Error"
}
```

#### Contract Test Cases

1. **Valid Token**: Returns 200 with complete user object
2. **Expired Token**: Returns 403 Forbidden
3. **Invalid Token**: Returns 403 Forbidden
4. **Missing Authorization Header**: Returns 401/403
5. **Malformed Token**: Returns 403
6. **Frozen Account**: Response includes `frozen: true` flag
7. **MFA Enabled**: Response includes `mfa: true` flag
8. **Empty Extra**: Handles empty object `extra: {}`
9. **Nested Extra Metadata**: Handles nested objects in `extra` field
10. **Unknown Fields**: Client ignores unexpected fields in response

---

### 2. GET /projects

List all projects accessible by the authenticated user.

#### Request

**Method**: `GET`  
**Path**: `/projects`  
**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
Accept: application/json
```

**Query Parameters**:
| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| offset | integer | No | 0 | Skip N items |
| limit | integer | No | 20 | Return max N items (1-100) |
| order | string | No | - | Sort order |

**Request Body**: None

#### Response

**Status Code**: `200 OK`

**Response Body**:
```json
{
  "projects": [
    {
      "project": {
        "projectId": "string (UUID)",
        "displayName": "string",
        "description": "string",
        "extra": {},
        "namespace": "string",
        "frozen": false,
        "createdAt": "string (ISO8601)",
        "updatedAt": "string (ISO8601)"
      },
      "globalPermissionId": "string (UUID)",
      "globalPermission": {
        "id": "string (UUID)",
        "label": "string"
      },
      "userPermissionId": "string (UUID)",
      "userPermission": {
        "id": "string (UUID)",
        "label": "string"
      },
      "frozen": false,
      "tenantRole": "string",
      "extra": {}
    }
  ],
  "total": 0
}
```

**Field Descriptions**:
- `projects`: Array of project memberships
  - `project`: Project details object
    - `projectId`: Unique project ID (UUID)
    - `displayName`: Project display name
    - `description`: Project description
    - `extra`: Arbitrary metadata (e.g., iService projectSysCode)
    - `namespace`: Tenant namespace (e.g., "ci.asus.com")
    - `frozen`: Project frozen/disabled status
    - `createdAt`: Project creation timestamp
    - `updatedAt`: Last update timestamp
  - `globalPermissionId`: Project-wide permission template ID
  - `globalPermission`: Default permission details
  - `userPermissionId`: User-specific permission ID
  - `userPermission`: User's effective permission
  - `frozen`: User's access frozen status (membership-level)
  - `tenantRole`: User's role enum value (TENANT_MEMBER, TENANT_ADMIN, TENANT_OWNER)
  - `extra`: Additional membership metadata
- `total`: Total count of accessible projects

**Note**: `tenantRole` is parsed as a TenantRole enum type with constants:
- `TenantRoleMember` - Standard project member
- `TenantRoleAdmin` - Project administrator  
- `TenantRoleOwner` - Project owner

#### Error Responses

**400 Bad Request** (invalid pagination):
```json
{
  "errorCode": 1001,
  "message": "Invalid pagination parameters"
}
```

**403 Forbidden**:
```json
{
  "errorCode": 1003,
  "message": "Forbidden"
}
```

**500 Internal Server Error**:
```json
{
  "errorCode": 2001,
  "message": "Internal Server Error"
}
```

#### Contract Test Cases

1. **No Pagination**: Returns default page (offset=0, limit=20)
2. **Custom Pagination**: offset=10, limit=5 returns correct subset
3. **Empty Result**: User with no projects returns empty array, total=0
4. **Max Limit**: limit=100 accepted
5. **Invalid Offset**: offset=-1 returns 400
6. **Invalid Limit**: limit=0 or limit=101 returns 400
7. **Frozen Project**: project.frozen=true included in results
8. **Frozen Membership**: membership frozen=true (user access frozen)
9. **Tenant Roles**: Different tenantRole values (TENANT_MEMBER, TENANT_ADMIN, TENANT_OWNER)
10. **Nested Extra**: Handles nested metadata in project.extra (e.g., iService data)
11. **Unknown Fields**: Client ignores unexpected fields

---

### 3. GET /project/{project-id}

Retrieve detailed information about a specific project.

#### Request

**Method**: `GET`  
**Path**: `/project/{project-id}`  
**Headers**:
```
Authorization: Bearer {token}
Content-Type: application/json
Accept: application/json
```

**Path Parameters**:
| Parameter | Type | Required | Description |
|-----------|------|----------|-------------|
| project-id | string (UUID) | Yes | Project ID to retrieve |

**Query Parameters**: None  
**Request Body**: None

#### Response

**Status Code**: `200 OK`

**Response Body**:
```json
**Response Body**:
```json
{
  "projectId": "string (UUID)",
  "displayName": "string",
  "description": "string",
  "extra": {},
  "namespace": "string",
  "frozen": false,
  "globalPermission": {
    "id": "string (UUID)",
    "label": "string"
  },
  "userPermission": {
    "id": "string (UUID)",
    "label": "string"
  },
  "createdAt": "string (ISO8601)",
  "updatedAt": "string (ISO8601)"
}
```

**Field Descriptions**:
- `projectId`: Unique project identifier (UUID) - **direct field, not nested**
- `displayName`: Project display name
- `description`: Project description
- `extra`: Arbitrary metadata (e.g., iService data)
- `namespace`: Tenant namespace
- `frozen`: Project frozen status
- `globalPermission`: Project's default permission
- `userPermission`: User's effective permission in project
- `createdAt`: Project creation timestamp
- `updatedAt`: Last update timestamp

**Important Differences from ListProjects**:
- **Flat structure**: Project fields are at the top level (not nested in `project` object)
- **No membership fields**: Does NOT include `tenantRole`, membership `frozen`, membership `extra`
- **No permission IDs**: Only includes permission objects, not `globalPermissionId`/`userPermissionId`
- GetProject returns **project-level data only**, while ListProjects returns **project + membership data**
```

**Field Descriptions**:
- `project`: Project details object (same as in ListProjects)
  - `projectId`: Unique project identifier (UUID)
  - `displayName`: Project display name
  - `description`: Project description
  - `extra`: Arbitrary metadata (e.g., iService data)
  - `namespace`: Tenant namespace
  - `frozen`: Project frozen status
  - `createdAt`: Project creation timestamp
  - `updatedAt`: Last update timestamp
- `globalPermissionId`: Project-wide permission template ID
- `globalPermission`: Project's default permission
  - `userPermissionId`: User-specific permission ID
  - `userPermission`: User's effective permission in project
  - `frozen`: User's membership frozen status
  - `tenantRole`: User's role enum value (TENANT_MEMBER, TENANT_ADMIN, TENANT_OWNER)
  - `extra`: Additional membership metadata

**Note**: GetProject response structure is **identical** to a single ProjectMembership item from ListProjects for consistency

**TenantRole Enum**: Parsed as custom TenantRole type with validation:
- `TenantRoleMember` - Standard member
- `TenantRoleAdmin` - Administrator
- `TenantRoleOwner` - Owner (highest authority)#### Error Responses

**400 Bad Request** (invalid project ID format):
```json
{
  "errorCode": 1001,
  "message": "user project does not exist"
}
```

**403 Forbidden** (user not member of project):
```json
{
  "errorCode": 1003,
  "message": "Forbidden"
}
```

**404 Not Found** (project doesn't exist):
```json
{
  "errorCode": 1004,
  "message": "Not Found"
}
```

**500 Internal Server Error**:
```json
{
  "errorCode": 2001,
  "message": "IAM Server internal error"
}
```

#### Contract Test Cases

1. **Valid Project ID**: User is member, returns 200 with project details
2. **Valid ID, Not Member**: Returns 403 or 400
3. **Invalid UUID Format**: Returns 400
4. **Non-existent Project**: Returns 400/404
5. **Empty Project ID**: Returns 400
6. **Frozen Project**: frozen=true returned (project level only)
7. **Flat Structure**: projectId, displayName etc. are direct fields (not nested)
8. **No Membership Data**: tenantRole, membership extra NOT included
9. **Nested Extra**: Handles nested metadata in extra field
10. **Unknown Fields**: Client ignores unexpected fields

---

## HTTP Status Code Summary

| Code | Meaning | When Used |
|------|---------|-----------|
| 200 | OK | Successful request |
| 400 | Bad Request | Invalid project ID, bad pagination params |
| 401 | Unauthorized | Missing auth header (implementation dependent) |
| 403 | Forbidden | Invalid/expired token, no access to resource |
| 404 | Not Found | Resource doesn't exist (implementation dependent) |
| 429 | Too Many Requests | Rate limit exceeded (retryable) |
| 500 | Internal Server Error | Service failure (retryable) |
| 502 | Bad Gateway | Upstream service failure (retryable) |
| 503 | Service Unavailable | Service temporarily down (retryable) |
| 504 | Gateway Timeout | Request timeout (retryable) |

## Retry Policy

Automatic retry with exponential backoff for:
- Status codes: 429, 502, 503, 504
- Network timeouts
- Connection errors

Retry configuration:
- Initial interval: 100ms
- Max interval: 5 seconds
- Multiplier: 2.0 (exponential)
- Max attempts: 3
- Jitter: enabled (±25% randomization)

Only GET requests are retried (safe/idempotent operations).

## Content Type

- Request: `application/json`
- Response: `application/json`
- Character encoding: UTF-8

## Timestamp Format

All timestamps use ISO8601/RFC3339 format:
```
2021-04-12T07:54:47Z
```

## Forward Compatibility

Clients MUST ignore unknown fields in responses to support API evolution:
- Server may add new fields without breaking clients
- Clients parse only documented fields
- Unknown fields are silently discarded during JSON unmarshaling

## Testing Strategy

### Contract Tests

Validate each endpoint against Swagger spec:

```go
// tests/contract/iam_contract_test.go
func TestGetUserContract(t *testing.T) {
    // 1. Call API
    // 2. Verify status code
    // 3. Validate response schema matches Swagger
    // 4. Check all required fields present
    // 5. Verify field types correct
}
```

### Integration Tests

Test real API interactions:

```go
// tests/integration/iam_integration_test.go
func TestIAMWorkflow(t *testing.T) {
    // 1. Get current user
    // 2. List their projects
    // 3. Get details of first project
    // 4. Verify data consistency
}
```

## References

- Swagger Spec: `/workspaces/cloud-sdk/swagger/iam.yaml`
- GetUser API: Lines 591-621 in iam.yaml
- ListProjects API: Lines 300-338 in iam.yaml  
- GetProject API: Lines 340-375 in iam.yaml
