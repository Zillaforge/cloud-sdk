# Volumes API Contracts

Extracted from: swagger/vps.yaml (volumes tag)
Date: 2025-11-13

This document specifies the API contracts for all 6 volume-related endpoints.

---

## 1. List Volumes

**Method**: GET
**Path**: `/api/v1/project/{project-id}/volumes`

### Parameters
- project-id (path, required): Project ID
- name (query, optional): Filter by name
- user_id (query, optional): Filter by user ID  
- status (query, optional): Filter by status
- type (query, optional): Filter by volume type
- detail (query, optional): Include attachment details

### Responses
- 200: VolumeListResponse `{"volumes": [...]}`
- 400: Invalid request
- 500: Server error

### Contract Tests
- Verify response has volumes array
- Test filtering parameters
- Verify detail flag includes attachments
- Test empty results

---

## 2. Create Volume

**Method**: POST
**Path**: `/api/v1/project/{project-id}/volumes`

### Request Body (VolumeCreateInput)
```json
{
  "name": "my-volume",
  "type": "SSD",
  "size": 100,
  "description": "Optional description",
  "snapshot_id": "Optional snapshot ID"
}
```

### Responses
- 201: Volume created `{...Volume...}`
- 400: Invalid input
- 500: Server error

### Contract Tests
- Verify required fields (name, type)
- Test optional fields (size, description, snapshot_id)
- Verify returned Volume has ID and timestamps
- Test quota exceeded error

---

## 3. Get Volume

**Method**: GET
**Path**: `/api/v1/project/{project-id}/volumes/{vol-id}`

### Parameters
- project-id (path, required): Project ID
- vol-id (path, required): Volume ID

### Responses
- 200: Volume details `{...Volume...}`
- 400: Invalid ID
- 404: Volume not found
- 500: Server error

### Contract Tests
- Verify all Volume fields present
- Test with invalid volume ID
- Verify timestamps are valid ISO 8601
- Check IDName structures (Project, User, Attachments)

---

## 4. Update Volume

**Method**: PUT
**Path**: `/api/v1/project/{project-id}/volumes/{vol-id}`

### Request Body (VolumeUpdateInput)
```json
{
  "name": "new-name",
  "description": "new-description"
}
```

### Responses
- 200: Updated Volume `{...Volume...}`
- 400: Invalid input
- 404: Volume not found
- 500: Server error

### Contract Tests
- Verify name can be updated
- Verify description can be updated
- Verify size/type cannot be updated
- Test partial updates (only name or only description)

---

## 5. Delete Volume

**Method**: DELETE
**Path**: `/api/v1/project/{project-id}/volumes/{vol-id}`

### Parameters
- project-id (path, required): Project ID
- vol-id (path, required): Volume ID

### Responses
- 204: Volume deleted (no content)
- 400: Invalid request (e.g., volume in use)
- 404: Volume not found
- 500: Server error

### Contract Tests
- Verify 204 on successful delete
- Test delete of attached volume (should fail)
- Verify volume no longer accessible after delete
- Test double delete (404)

---

## 6. Volume Action

**Method**: POST
**Path**: `/api/v1/project/{project-id}/volumes/{vol-id}/action`

### Request Body (VolActionInput)
```json
{
  "action": "attach|detach|extend|revert",
  "server_id": "required for attach/detach",
  "new_size": 200
}
```

### Action Types
- **attach**: Attach volume to server (requires server_id)
- **detach**: Detach volume from server (requires server_id)
- **extend**: Extend volume size (requires new_size > current size)
- **revert**: Revert volume to snapshot

### Responses
- 202: Action accepted (async operation)
- 400: Invalid action or parameters
- 404: Volume or server not found
- 500: Server error

### Contract Tests
- Test all 4 action types
- Verify action-specific parameter requirements
- Test invalid action type
- Test extend with smaller size (should fail)
- Test attach already attached volume (should fail)
- Test detach unattached volume (should fail)

---

## Shared Response Structures

### Volume
```json
{
  "id": "vol-123",
  "name": "my-volume",
  "description": "Optional",
  "size": 100,
  "type": "SSD",
  "status": "available",
  "status_reason": "",
  "attachments": [
    {"id": "server-1", "name": "server-name"}
  ],
  "project": {"id": "proj-1", "name": "Project"},
  "project_id": "proj-1",
  "user": {"id": "user-1", "name": "User"},
  "user_id": "user-1",
  "namespace": "default",
  "createdAt": "2025-11-13T10:00:00Z",
  "updatedAt": "2025-11-13T10:30:00Z"
}
```

### Error Response
```json
{
  "errorCode": "ERROR_CODE",
  "message": "Human readable message"
}
```

---

## Authentication

All endpoints require Bearer Token:
```
Authorization: Bearer {token}
```

Unauthorized responses:
- 401: Missing or invalid token
- 403: Insufficient permissions

---

## Implementation Checklist

For each endpoint:
- [ ] Request struct defined
- [ ] Response struct defined
- [ ] Client method implemented
- [ ] Error wrapping applied
- [ ] Unit tests written (85%+ coverage)
- [ ] Contract test written
- [ ] Integration test written
- [ ] Documentation/examples added

---

## Validation Rules

### Client-Side (before API call)
- Required fields present
- Action-specific parameters validated
- Basic type checking

### Server-Side (API enforced)
- Size limits and quotas
- Volume status transitions
- Resource existence
- Permission checks

---

## Error Handling Patterns

```go
// List volumes
return nil, fmt.Errorf("failed to list volumes: %w", err)

// Create volume
return nil, fmt.Errorf("failed to create volume: %w", err)

// Get volume
return nil, fmt.Errorf("failed to get volume %s: %w", volumeID, err)

// Update volume
return nil, fmt.Errorf("failed to update volume %s: %w", volumeID, err)

// Delete volume
return fmt.Errorf("failed to delete volume %s: %w", volumeID, err)

// Volume action
return fmt.Errorf("failed to perform action on volume %s: %w", volumeID, err)
```

---

## Performance Expectations

Per spec.md Success Criteria:
- List volumes: < 3 seconds
- Create volume: < 5 seconds
- Get volume: < 2 seconds
- Update volume: < 3 seconds
- Delete volume: < 5 seconds
- Volume actions: < 10 seconds

---

## Reference

Full API documentation: swagger/vps.yaml (volumes tag, lines 7664-7960)
