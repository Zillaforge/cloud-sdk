# Keypair API Contract

**Extracted from**: `/workspaces/cloud-sdk/swagger/vps.yaml`  
**Feature**: 006-fix-keypair-model  
**Date**: 2025-11-10

## Data Models

### pb.IDName
```yaml
pb.IDName:
  properties:
    id:
      type: string
    name:
      type: string
  type: object
```

### pb.KeypairInfo
```yaml
pb.KeypairInfo:
  properties:
    createdAt:
      type: string
    description:
      type: string
    fingerprint:
      type: string
    id:
      type: string
    name:
      type: string
    private_key:
      type: string
    public_key:
      type: string
    updatedAt:
      type: string
    user:
      $ref: '#/definitions/pb.IDName'
    user_id:
      type: string
  type: object
```

### pb.KeypairListOutput
```yaml
pb.KeypairListOutput:
  properties:
    keypairs:
      items:
        $ref: '#/definitions/pb.KeypairInfo'
      type: array
  type: object
```

### KeypairCreateInput
```yaml
KeypairCreateInput:
  properties:
    description:
      type: string
    name:
      type: string
    public_key:
      type: string
  required:
  - name
  type: object
```

### KeypairUpdateInput
```yaml
KeypairUpdateInput:
  properties:
    description:
      type: string
  type: object
```

## API Operations

### List Keypairs
```yaml
/api/v1/project/{project-id}/keypairs:
  get:
    summary: List keypairs
    parameters:
    - in: path
      name: project-id
      required: true
      type: string
    - in: query
      name: name
      description: Filter by keypair name
      required: false
      type: string
    responses:
      200:
        description: Success
        schema:
          $ref: '#/definitions/pb.KeypairListOutput'
```

### Create Keypair
```yaml
/api/v1/project/{project-id}/keypairs:
  post:
    summary: Create keypair
    parameters:
    - in: path
      name: project-id
      required: true
      type: string
    - in: body
      name: body
      required: true
      schema:
        $ref: '#/definitions/KeypairCreateInput'
    responses:
      200:
        description: Success
        schema:
          $ref: '#/definitions/pb.KeypairInfo'
```

### Get Keypair
```yaml
/api/v1/project/{project-id}/keypairs/{keypair-id}:
  get:
    summary: Get keypair details
    parameters:
    - in: path
      name: project-id
      required: true
      type: string
    - in: path
      name: keypair-id
      required: true
      type: string
    responses:
      200:
        description: Success
        schema:
          $ref: '#/definitions/pb.KeypairInfo'
```

### Update Keypair
```yaml
/api/v1/project/{project-id}/keypairs/{keypair-id}:
  put:
    summary: Update keypair
    parameters:
    - in: path
      name: project-id
      required: true
      type: string
    - in: path
      name: keypair-id
      required: true
      type: string
    - in: body
      name: body
      required: true
      schema:
        $ref: '#/definitions/KeypairUpdateInput'
    responses:
      200:
        description: Success
        schema:
          $ref: '#/definitions/pb.KeypairInfo'
```

### Delete Keypair
```yaml
/api/v1/project/{project-id}/keypairs/{keypair-id}:
  delete:
    summary: Delete keypair
    parameters:
    - in: path
      name: project-id
      required: true
      type: string
    - in: path
      name: keypair-id
      required: true
      type: string
    responses:
      204:
        description: No content
```

## Contract Test Cases

### Test 1: List Response Structure
**Given**: API returns list of keypairs  
**Expect**: Response matches `pb.KeypairListOutput` schema  
**Verify**: 
- `keypairs` array present
- NO `total` field present
- Each item matches `pb.KeypairInfo` schema

### Test 2: Create Response with Generated Key
**Given**: Create request without `public_key`  
**Expect**: Response includes `private_key` field  
**Verify**:
- All `pb.KeypairInfo` fields present
- `private_key` field non-empty
- Timestamps in ISO 8601 format

### Test 3: Get Response Structure
**Given**: Get existing keypair  
**Expect**: Response matches `pb.KeypairInfo` schema  
**Verify**:
- All required fields present
- NO `private_key` field (or empty)
- `user` may be null
- Timestamps parseable as RFC3339

### Test 4: Optional Fields Handling
**Given**: API response with null/missing optional fields  
**Expect**: Model deserializes successfully  
**Verify**:
- `description` empty when not provided
- `user` null when not provided
- `private_key` empty in Get/List operations

### Test 5: User Reference
**Given**: Response with populated `user` object  
**Expect**: User reference matches `pb.IDName` schema  
**Verify**:
- `user.id` matches `user_id`
- `user.name` present when user object populated
