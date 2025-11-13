# Volume Types API Contract

Extracted from: swagger/vps.yaml (volume_types tag)
Date: 2025-11-13

## Endpoint: List Volume Types

Method: GET
Path: /api/v1/project/{project-id}/volume_types
Tag: volume_types

### Description

List all available volume types in the system.

### Parameters

- Name: project-id
  Type: string
  Location: path
  Required: Yes
  Description: Project ID

### Authentication

Bearer Token via Authorization header (ApiKeyAuth)

### Request Headers

Authorization: Bearer {token}
Accept: application/json

### Request Example

curl -X 'GET' \
  'https://localhost:9999/api/v1/project/871a0333-ae33-43f2-81f9-2432fa15cde7/volume_types' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer {token}'

### Response: 200 OK

Description: Successfully retrieved list of volume types
Content-Type: application/json
Schema: pb.VolumeTypeListOutput

Example JSON:
{
  "volume_types": ["SSD", "HDD", "NVMe"]
}

Response Structure (Go):
type VolumeTypeListResponse struct {
    VolumeTypes []string `json:"volume_types"`
}

### Response: 400 Bad Request

Description: Invalid request (e.g., invalid project ID format)
Content-Type: application/json
Schema: vpserr.ErrResponse

Example JSON:
{
  "errorCode": "INVALID_PROJECT_ID",
  "message": "Invalid project ID format"
}

### Response: 500 Internal Server Error

Description: Server error
Content-Type: application/json
Schema: vpserr.ErrResponse

Example JSON:
{
  "errorCode": "INTERNAL_ERROR",
  "message": "Internal server error"
}

---

## Contract Test Requirements

1. Status Code Validation:
   - Verify 200 response with valid project ID
   - Verify 400 response with invalid project ID
   - Verify 500 response can be handled

2. Response Structure:
   - Response must have volume_types field
   - volume_types must be string array
   - Array may be empty (valid scenario)

3. Authentication:
   - Request without Bearer token returns 401
   - Request with invalid token returns 401
   - Request with expired token returns 401

4. Data Validation:
   - Each volume type is non-empty string
   - Common types: SSD, HDD, NVMe, etc.
   - Types are consistent across calls

---

## Implementation Notes

- Response Unwrapping: Client must unmarshal to VolumeTypeListResponse then extract VolumeTypes array
- No Pagination: Returns complete list in single response
- Caching: Volume types are relatively static, consider caching
- Error Handling: Wrap errors with context using fmt.Errorf
