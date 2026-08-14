# Data Model

## Entities

### Server

**Fields** (from pb.ServerInfo):
- id: string (unique identifier)
- name: string
- description: string (optional)
- status: ServerStatus (custom type simulating enum: ACTIVE, BUILD, SHUTOFF, ERROR, etc.)
- status_reason: string (optional)
- flavor_id: string
- flavor: IDName (id, name)
- flavor_detail: FlavorInfo (detailed flavor info)
- image_id: string
- image: VRMImgInfo (repository info)
- project_id: string
- project: IDName
- user_id: string
- user: IDName
- keypair_id: string
- keypair: IDName
- metadatas: map[string]string
- private_ips: []string
- public_ips: []string
- az: string
- namespace: string
- root_disk_id: string
- root_disk_size: int
- boot_script: string
- approvedAt: string
- createdAt: string
- updatedAt: string
- uuid: string

**Relationships**:
- Belongs to Project
- Has many NICs (ServerNIC)
- Has Flavor, Image, Keypair

**Validation Rules**:
- id, name, status required
- flavor_id, image_id required for creation
- status in allowed values

**State Transitions**:
- BUILD → ACTIVE
- ACTIVE → SHUTOFF, REBOOT
- SHUTOFF → ACTIVE
- Any → ERROR

**Custom Types**:
```go
type ServerStatus string

const (
    ServerStatusActive   ServerStatus = "ACTIVE"
    ServerStatusBuild    ServerStatus = "BUILD"
    ServerStatusShutoff  ServerStatus = "SHUTOFF"
    ServerStatusError    ServerStatus = "ERROR"
    ServerStatusReboot   ServerStatus = "REBOOT"
    // Add other statuses as defined in API
)
```

### ServerCreateRequest

**Fields** (from SvrCreateInput in Swagger, struct name in Go: ServerCreateRequest):
- name: string
- description: string (optional)
- flavor_id: string
- image_id: string
- nics: []NicCreateInput
- sg_ids: []string
- keypair_id: string (optional)
- password: string (optional, base64)
- boot_script: string (optional, base64)
- volume_ids: []string (deprecated)
- volumes: []ServerDiskInput

**Validation**: name, flavor_id, image_id, nics required

### ServerUpdateRequest

**Fields** (from SvrUpdateInput in Swagger, struct name in Go: ServerUpdateRequest):
- name: string (optional)
- description: string (optional)

### ServerActionRequest

**Fields** (from SvrActionInput in Swagger, struct name in Go: ServerActionRequest):
- action: ServerAction (custom type simulating enum: stop, start, reboot, resize, approve, reject, extend_root, get_pwd)
- flavor_id: string (for resize)
- private_key: string (for get_pwd)
- reboot_type: RebootType (custom type: hard/soft for reboot)
- root_size: int (for extend_root)

**Validation**: action required

**Custom Types**:
```go
type ServerAction string

const (
    ServerActionStop       ServerAction = "stop"
    ServerActionStart      ServerAction = "start"
    ServerActionReboot     ServerAction = "reboot"
    ServerActionResize     ServerAction = "resize"
    ServerActionApprove    ServerAction = "approve"
    ServerActionReject     ServerAction = "reject"
    ServerActionExtendRoot ServerAction = "extend_root"
    ServerActionGetPwd     ServerAction = "get_pwd"
)

type RebootType string

const (
    RebootTypeHard RebootType = "hard"
    RebootTypeSoft RebootType = "soft"
)
```

### ServerActionResponse

**Fields** (from SvrActionOutput in Swagger, struct name in Go: ServerActionResponse):
- password: string (returned for get_pwd action)

### ServerMetricsRequest

**Fields** (query parameters, struct name in Go: ServerMetricsRequest):
- type: string (cpu, memory, disk, net, vgpu)
- granularity: int
- start: int (unix timestamp)
- direction: string (incoming/outgoing for net)
- rw: string (read/write for disk)

### ServerMetricsResponse

**Fields** (struct name in Go: ServerMetricsResponse, returns array of MetricInfo):
- measures: []MetricInfo (timestamp, value, granularity)

### ServerNIC

**Fields** (from pb.ServerNICInfo):
- id: string
- mac: string
- network_id: string
- network: IDName
- addresses: []string
- floating_ip: FloatingIPInfo (optional)
- security_groups: []IDName
- sg_ids: []string
- is_provider_net: bool

**Relationships**:
- Belongs to Server
- Belongs to Network
- Has Security Groups

**Operations** (Sub-resource pattern):
- List: GET /servers/{svr-id}/nics
- Create: POST /servers/{svr-id}/nics
- Update: PUT /servers/{svr-id}/nics/{nic-id}
- Delete: DELETE /servers/{svr-id}/nics/{nic-id}
- Associate FloatingIP: POST /servers/{svr-id}/nics/{nic-id}/floatingip

### ServerNICCreateRequest

**Fields** (from NicCreateInput in Swagger, struct name in Go: ServerNICCreateRequest):
- network_id: string
- sg_ids: []string
- fixed_ip: string (optional)

**Validation**: network_id, sg_ids required

### ServerNICUpdateRequest

**Fields** (from NicUpdateInput in Swagger, struct name in Go: ServerNICUpdateRequest):
- sg_ids: []string

**Validation**: sg_ids required

## List Options & Responses

### ServersListRequest

**Fields** (query parameters, struct name in Go: ServersListRequest):
- name: string
- user_id: string
- status: string
- flavor_id: string
- image_id: string
- detail: bool (query parameter to include detailed information)

### ServersListResponse

**Fields** (struct name in Go: ServersListResponse):
- servers: []Server

### ServerNICsListResponse

**Fields** (struct name in Go: ServerNICsListResponse):
- nics: []ServerNIC