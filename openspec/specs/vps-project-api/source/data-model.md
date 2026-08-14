# Data Model: VPS Project APIs SDK

**Phase 1 Output** | **Date**: 2025-10-26  
**Purpose**: Define entities, request/response types, error structures, and state machines extracted from feature spec and Swagger

## Core SDK Types

### 1. Client (Top-Level)

**Package**: `cloudsdk` (root)  
**Purpose**: Entry point for SDK; manages base URL, auth, HTTP client  
**Fields**:
```go
type Client struct {
    baseURL    string          // e.g., "https://api.example.com"
    token      string          // Bearer token
    httpClient *http.Client    // with timeout, TLS config
    logger     Logger          // optional, for retry/debug logs
}
```
**Methods**:
- `New(baseURL, token string, opts ...ClientOption) *Client`
- `Project(projectID string) *ProjectClient`

**Validation**:
- `baseURL` must be valid URL with scheme (`https://`)
- `token` must be non-empty

---

### 2. ProjectClient (Project Selector)

**Package**: `cloudsdk` (root)  
**Purpose**: Project-scoped handle; creates service clients bound to project ID  
**Fields**:
```go
type ProjectClient struct {
    client    *Client
    projectID string
}
```
**Methods**:
- `VPS() *vps.Client` (returns project-scoped VPS client)
- Future: `Compute() *compute.Client`, `Storage() *storage.Client`

---

### 3. SDKError (Structured Error)

**Package**: `cloudsdk` (root `errors.go`)  
**Purpose**: Unified error type per spec Error Handling Contract  
**Fields**:
```go
type SDKError struct {
    StatusCode int                    // HTTP status (0 for client-side)
    ErrorCode  int                    // from vpserr.ErrResponse.errorCode
    Message    string                 // human-readable
    Meta       map[string]interface{} // additional context
    Cause      error                  // underlying error (wrapped)
}
```
**Methods**:
- `Error() string` → returns `Message`
- `Unwrap() error` → returns `Cause`
- `Is(target error) bool` → allows `errors.Is` checks

**Mapping Rules**:
| Scenario | StatusCode | ErrorCode | Message | Meta |
|----------|-----------|-----------|---------|------|
| HTTP 400-5xx with JSON body | HTTP status | from body | from body | from body |
| HTTP error, no body | HTTP status | 0 | "HTTP {status}" | {"raw": rawBody} |
| Network/DNS error | 0 | 0 | "network error: ..." | {"category": "network"} |
| Context timeout | 0 | 0 | "request timeout" | {"category": "timeout"} |
| Context canceled | 0 | 0 | "request canceled" | {"category": "canceled"} |

---

## VPS Service Types

### 4. VPS Client (Project-Scoped)

**Package**: `vps`  
**Purpose**: Encapsulates all VPS operations for a project; no project-id params in methods  
**Fields**:
```go
type Client struct {
    baseClient *cloudsdk.Client
    projectID  string
    basePath   string // "/vps/api/v1/project/{project-id}"
}
```
**Methods** (grouped by resource):
- **Servers**: `List`, `Create`, `Get`, `Update`, `Delete`, `Action`, `Metrics`, `VNCURL`
  - Server sub-resource NICs: `server.NICs().List()`, `server.NICs().Add()`, `server.NICs().Update()`, `server.NICs().Delete()`, `server.NICs().AssociateFloatingIP()`
  - Server sub-resource Volumes: `server.Volumes().List()`, `server.Volumes().Attach()`, `server.Volumes().Detach()`
- **Networks**: `List`, `Create`, `Get`, `Update`, `Delete`
  - Network sub-resource: `network.Ports().List()`
- **Floating IPs**: `List`, `Create`, `Get`, `Update`, `Delete`, `Approve`, `Reject`, `Disassociate`
- **Keypairs**: `List`, `Create`, `Get`, `Update`, `Delete`
- **Routers**: `List`, `Create`, `Get`, `Update`, `Delete`, `SetState`
  - Router sub-resource: `router.Networks().List()`, `router.Networks().Associate()`, `router.Networks().Disassociate()`
- **Security Groups**: `List`, `Create`, `Get`, `Update`, `Delete`
  - SecurityGroup sub-resource: `securityGroup.Rules().Add()`, `securityGroup.Rules().Delete()`
- **Flavors**: `List`, `Get`
- **Quotas**: `Get`

---

### 5. Entity Models (Extracted from Swagger)

#### Server
```go
type Server struct {
    ID          string            `json:"id"`
    Name        string            `json:"name"`
    Description string            `json:"description,omitempty"`
    Status      string            `json:"status"` // e.g., "ACTIVE", "BUILD", "ERROR"
    FlavorID    string            `json:"flavor_id"`
    ImageID     string            `json:"image_id"`
    ProjectID   string            `json:"project_id"`
    UserID      string            `json:"user_id,omitempty"`
    CreatedAt   string            `json:"created_at"` // ISO 8601
    UpdatedAt   string            `json:"updated_at"`
    Metadata    map[string]string `json:"metadata,omitempty"`
    // ... additional fields from pb.ServerInfo
}

type ServerCreateRequest struct {
    Name        string              `json:"name"`
    Description string              `json:"description,omitempty"`
    FlavorID    string              `json:"flavor_id"`
    ImageID     string              `json:"image_id"`
    NICs        []ServerNICRequest  `json:"nics"`
    SGIDs       []string            `json:"sg_ids"`
    KeypairID   string              `json:"keypair_id,omitempty"`
    Password    string              `json:"password,omitempty"` // Base64
    BootScript  string              `json:"boot_script,omitempty"` // Base64
    // ... volume attachment fields
}

type ServerActionRequest struct {
    Action      string `json:"action"` // "start","stop","reboot","resize","extend_root","get_pwd","approve","reject"
    RebootType  string `json:"reboot_type,omitempty"` // "hard" or "soft"
    FlavorID    string `json:"flavor_id,omitempty"` // for resize
    RootSize    int    `json:"root_size,omitempty"` // for extend_root
    PrivateKey  string `json:"private_key,omitempty"` // Base64 for get_pwd
}

type ServerMetricsRequest struct {
    Type        string `json:"type"` // "cpu","memory","disk","network"
    Start       int64  `json:"start"` // Unix timestamp
    End         int64  `json:"end,omitempty"`
    Granularity int    `json:"granularity,omitempty"` // seconds
}

type ServerMetricsResponse struct {
    Type   string          `json:"type"`
    Series []MetricPoint   `json:"series"`
}

type MetricPoint struct {
    Timestamp int64   `json:"timestamp"`
    Value     float64 `json:"value"`
}
```

**State Transitions**:
- `BUILD` → `ACTIVE` (on successful creation)
- `ACTIVE` ↔ `SHUTOFF` (start/stop actions)
- `ACTIVE` → `REBOOT` → `ACTIVE` (reboot action)
- `ACTIVE` → `RESIZE` → `VERIFY_RESIZE` → `ACTIVE` (resize action)
- Any → `ERROR` (on failure)

---

#### Network
```go
type Network struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    CIDR        string   `json:"cidr"` // e.g., "192.168.1.0/24"
    ProjectID   string   `json:"project_id"`
    CreatedAt   string   `json:"created_at"`
    UpdatedAt   string   `json:"updated_at"`
}

type NetworkCreateRequest struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    CIDR        string `json:"cidr"`
}

type NetworkPort struct {
    ID        string   `json:"id"`
    NetworkID string   `json:"network_id"`
    FixedIPs  []string `json:"fixed_ips"`
    MACAddr   string   `json:"mac_address"`
    ServerID  string   `json:"server_id,omitempty"`
}
```

---

#### FloatingIP
```go
type FloatingIP struct {
    ID          string `json:"id"`
    Address     string `json:"address"` // public IP
    Status      string `json:"status"` // "ACTIVE", "PENDING", "DOWN"
    ProjectID   string `json:"project_id"`
    PortID      string `json:"port_id,omitempty"` // attached NIC
    CreatedAt   string `json:"created_at"`
}

type FloatingIPCreateRequest struct {
    Description string `json:"description,omitempty"`
    // ExtNetworkID omitted if using default external network
}
```

**State Transitions**:
- `PENDING` → `ACTIVE` (after admin approval)
- `PENDING` → `REJECTED` (admin rejects)
- `ACTIVE` → `DOWN` (disassociated)

---

#### Router
```go
type Router struct {
    ID          string   `json:"id"`
    Name        string   `json:"name"`
    Description string   `json:"description,omitempty"`
    State       string   `json:"state"` // "enabled", "disabled"
    ProjectID   string   `json:"project_id"`
    Networks    []string `json:"networks,omitempty"` // associated network IDs
}

type RouterCreateRequest struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
}

type RouterSetStateRequest struct {
    State string `json:"state"` // "enabled" or "disabled"
}
```

---

#### SecurityGroup
```go
type SecurityGroup struct {
    ID          string              `json:"id"`
    Name        string              `json:"name"`
    Description string              `json:"description,omitempty"`
    ProjectID   string              `json:"project_id"`
    Rules       []SecurityGroupRule `json:"rules,omitempty"`
}

type SecurityGroupRule struct {
    ID          string `json:"id"`
    Direction   string `json:"direction"` // "ingress" or "egress"
    Protocol    string `json:"protocol"` // "tcp", "udp", "icmp"
    PortMin     int    `json:"port_range_min,omitempty"`
    PortMax     int    `json:"port_range_max,omitempty"`
    RemoteCIDR  string `json:"remote_ip_prefix,omitempty"`
    RemoteGroup string `json:"remote_group_id,omitempty"`
}

type SecurityGroupCreateRequest struct {
    Name        string              `json:"name"`
    Description string              `json:"description,omitempty"`
    Rules       []SecurityGroupRule `json:"rules,omitempty"`
}
```

---

#### Keypair
```go
type Keypair struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    PublicKey   string `json:"public_key"`
    Fingerprint string `json:"fingerprint"`
    UserID      string `json:"user_id"`
}

type KeypairCreateRequest struct {
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    PublicKey   string `json:"public_key,omitempty"` // import existing; omit to generate
}
```

---

#### Flavor
```go
type Flavor struct {
    ID          string `json:"id"`
    Name        string `json:"name"`
    Description string `json:"description,omitempty"`
    VCPUs       int    `json:"vcpus"`
    RAM         int    `json:"ram"` // MiB
    Disk        int    `json:"disk"` // GiB
    Public      bool   `json:"public"`
    Tags        []string `json:"tags,omitempty"`
}
```

---

#### Quota
```go
type Quota struct {
    VM         QuotaDetail `json:"vm"`
    VCPU       QuotaDetail `json:"vcpu"`
    RAM        QuotaDetail `json:"ram"`
    GPU        QuotaDetail `json:"gpu"`
    BlockSize  QuotaDetail `json:"block_size"`
    Network    QuotaDetail `json:"network"`
    Router     QuotaDetail `json:"router"`
    FloatingIP QuotaDetail `json:"floating_ip"`
    Share      QuotaDetail `json:"share,omitempty"`
    ShareSize  QuotaDetail `json:"share_size,omitempty"`
}

type QuotaDetail struct {
    Limit int `json:"limit"` // -1 = unlimited
    Usage int `json:"usage"`
}
```

---

### 6. Waiter State Machines

#### WaitForServerStatus
**Purpose**: Poll server until it reaches target status or timeout  
**States**:
- `BUILD` → `ACTIVE` (typical creation flow)
- `SHUTOFF` → `ACTIVE` (start action)
- `ACTIVE` → `SHUTOFF` (stop action)
- `*` → `ERROR` (failure state, exit)

**Algorithm**:
1. Initial delay: 2s
2. Poll interval: 5s (configurable via `WaiterOption`)
3. Max wait: 5 minutes (configurable)
4. Backoff: none (constant interval) unless specified
5. Exit conditions:
   - Status == target → success
   - Status == "ERROR" → error
   - Context canceled/timeout → error
   - Max wait exceeded → `ErrWaiterTimeout`

**Go Signature**:
```go
func (c *Client) WaitForServerStatus(ctx context.Context, serverID, targetStatus string, opts ...WaiterOption) error
```

#### WaitForFloatingIPActive
**Purpose**: Poll FIP until `ACTIVE` after creation (admin approval flow)  
**States**: `PENDING` → `ACTIVE` (or `REJECTED`)

**Algorithm**: Similar to server waiter; max wait 10 minutes (approval latency)

---

## Validation Rules

| Entity | Field | Rule |
|--------|-------|------|
| Server | Name | Required, max 255 chars |
| Server | FlavorID | Required, must be valid UUID |
| Server | Password | If present, must be Base64 |
| Network | CIDR | Required, valid IPv4 CIDR (e.g., `192.168.0.0/24`) |
| SecurityGroupRule | Protocol | One of: tcp, udp, icmp |
| SecurityGroupRule | PortMin/Max | Valid range 1-65535 |
| Router | State | One of: enabled, disabled |
| FloatingIP | Address | Read-only (assigned by server) |
| Quota | Limit | -1 or >= 0 |

**Validation Strategy**: Perform client-side validation in request builder methods (e.g., `NewServerCreateRequest`) to fail fast; server will also validate (definitive).

---

## Error Taxonomy

| Error Type | StatusCode | Example Message |
|------------|-----------|-----------------|
| Validation | 400 | "Invalid CIDR format" |
| Unauthorized | 401 | "Bearer token expired" |
| Forbidden | 403 | "Project access denied" |
| Not Found | 404 | "Server {id} not found" |
| Conflict | 409 | "Network CIDR overlaps" |
| Rate Limit | 429 | "Too many requests" |
| Server Error | 500 | "Internal server error" |
| Bad Gateway | 502 | "Upstream unavailable" |
| Service Unavailable | 503 | "Service temporarily unavailable" |
| Timeout (client) | 0 | "Request timeout" |
| Network (client) | 0 | "Network error: connection refused" |

---

## Pagination Model (Future Enhancement)

**Note**: Current spec (FR-013) mentions pagination/filters MUST be supported where defined. Swagger endpoints likely support limit/offset or marker-based pagination. Deferred to tasks phase to inspect Swagger and implement if present.

Placeholder pattern:
```go
type ListOptions struct {
    Limit  int    `json:"limit,omitempty"`
    Offset int    `json:"offset,omitempty"`
    // or Marker string `json:"marker,omitempty"` for cursor-based
}

type ListResponse struct {
    Items []T        `json:"items"`
    Total int        `json:"total,omitempty"`
    Next  string     `json:"next,omitempty"` // URL or marker
}
```

---

## Summary

- ✅ Core SDK types: `Client`, `ProjectClient`, `SDKError`
- ✅ VPS entities: Server, Network, FloatingIP, Router, SecurityGroup, Keypair, Flavor, Quota
- ✅ Request/response models aligned with Swagger
- ✅ State machines for waiters (Server, FloatingIP)
- ✅ Validation rules per entity
- ✅ Error taxonomy mapped to HTTP status codes

**Next**: Generate contracts/ (Go interfaces for each resource group)
