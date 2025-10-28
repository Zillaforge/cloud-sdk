// Package contracts defines the VPS server operations interface.
// This file is generated from swagger/vps.json and serves as a test-first contract.
//
// All methods accept context.Context for cancellation/timeout.
// All methods return (*Response, error) where error is *cloudsdk.SDKError on failure.
//
// Swagger reference: /api/v1/project/{project-id}/servers*
package contracts

import "context"

// ServerOperations defines all server-related API operations.
// Implementation: vps.Client
type ServerOperations interface {
	// List retrieves all servers for the project with optional filters.
	// GET /api/v1/project/{project-id}/servers
	// Query params: name, user_id, status, flavor_id, image_id
	// Returns: ServerListResponse (items + pagination)
	// Errors: 401 (unauthorized), 403 (forbidden), 500
	List(ctx context.Context, opts *ListServersOptions) (*ServerListResponse, error)

	// Create provisions a new server instance.
	// POST /api/v1/project/{project-id}/servers
	// Body: ServerCreateRequest
	// Returns: 201 + ServerInfo
	// Errors: 400 (validation), 401, 403, 409 (quota exceeded), 500
	Create(ctx context.Context, req *ServerCreateRequest) (*Server, error)

	// Get retrieves details for a specific server.
	// GET /api/v1/project/{project-id}/servers/{svr-id}
	// Returns: 200 + ServerInfo
	// Errors: 401, 403, 404 (not found), 500
	Get(ctx context.Context, serverID string) (*ServerResource, error)

	// Update modifies server name/description.
	// PUT /api/v1/project/{project-id}/servers/{svr-id}
	// Body: ServerUpdateRequest (name, description)
	// Returns: 200 + ServerInfo
	// Errors: 400, 401, 403, 404, 500
	Update(ctx context.Context, serverID string, req *ServerUpdateRequest) (*Server, error)

	// Delete removes a server instance.
	// DELETE /api/v1/project/{project-id}/servers/{svr-id}
	// Returns: 204 (no content)
	// Errors: 401, 403, 404, 409 (in use), 500
	Delete(ctx context.Context, serverID string) error

	// Action performs a control action on a server.
	// POST /api/v1/project/{project-id}/servers/{svr-id}/action
	// Body: ServerActionRequest (action, params)
	// Actions: start, stop, reboot, resize, extend_root, get_pwd, approve, reject
	// Returns: 202 (async action accepted)
	// Errors: 400 (invalid action/params), 401, 403, 404, 500
	Action(ctx context.Context, serverID string, req *ServerActionRequest) error

	// Metrics retrieves time-series metrics for a server.
	// GET /api/v1/project/{project-id}/servers/{svr-id}/metric
	// Query params: type (cpu/memory/disk/network), start, end, granularity
	// Returns: 200 + ServerMetricsResponse
	// Errors: 400 (invalid range), 401, 403, 404, 500
	Metrics(ctx context.Context, serverID string, req *ServerMetricsRequest) (*ServerMetricsResponse, error)

	// VNCURL retrieves the VNC console URL for a server.
	// GET /api/v1/project/{project-id}/servers/{svr-id}/vnc_url
	// Returns: 200 + VNCURLResponse (url string)
	// Errors: 401, 403, 404, 500
	VNCURL(ctx context.Context, serverID string) (*VNCURLResponse, error)
}

// ListServersOptions contains filter/pagination options for List.
type ListServersOptions struct {
	Name     string
	UserID   string
	Status   string
	FlavorID string
	ImageID  string
	// Pagination fields (if Swagger defines them):
	Limit  int
	Offset int
}

// ServerListResponse is the response from List.
type ServerListResponse struct {
	Items []*Server `json:"items"`
	Total int       `json:"total,omitempty"`
}

// Server represents a compute instance (from pb.ServerInfo).
type Server struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Status      string            `json:"status"` // ACTIVE, BUILD, SHUTOFF, ERROR, etc.
	FlavorID    string            `json:"flavor_id"`
	ImageID     string            `json:"image_id"`
	ProjectID   string            `json:"project_id"`
	UserID      string            `json:"user_id,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
	// Additional fields from Swagger as needed
}

// ServerResource wraps a Server with sub-resource operations.
type ServerResource struct {
	*Server
	nicOps    ServerNICOperations
	volumeOps ServerVolumeOperations
}

// NICs returns the NIC operations for this server.
func (s *ServerResource) NICs() ServerNICOperations {
	return s.nicOps
}

// Volumes returns the volume operations for this server.
func (s *ServerResource) Volumes() ServerVolumeOperations {
	return s.volumeOps
}

// ServerNICOperations defines operations on server NICs (sub-resource).
type ServerNICOperations interface {
	// List lists all network interfaces on the server.
	// GET /api/v1/project/{project-id}/servers/{svr-id}/nics
	// Returns: 200 + array of NICInfo
	// Errors: 401, 403, 404, 500
	List(ctx context.Context) ([]*ServerNIC, error)

	// Add attaches a new vNIC to the server.
	// POST /api/v1/project/{project-id}/servers/{svr-id}/nics
	// Body: ServerNICCreateRequest (network_id, sg_ids, fixed_ip)
	// Returns: 201 + NICInfo
	// Errors: 400, 401, 403, 404, 409 (quota/network conflict), 500
	Add(ctx context.Context, req *ServerNICCreateRequest) (*ServerNIC, error)

	// Update updates security groups on an existing vNIC.
	// PUT /api/v1/project/{project-id}/servers/{svr-id}/nics/{nic-id}
	// Body: ServerNICUpdateRequest (sg_ids)
	// Returns: 200 + NICInfo
	// Errors: 400, 401, 403, 404, 500
	Update(ctx context.Context, nicID string, req *ServerNICUpdateRequest) (*ServerNIC, error)

	// Delete detaches and removes a vNIC from the server.
	// DELETE /api/v1/project/{project-id}/servers/{svr-id}/nics/{nic-id}
	// Returns: 204
	// Errors: 401, 403, 404, 409 (last NIC), 500
	Delete(ctx context.Context, nicID string) error

	// AssociateFloatingIP associates a floating IP to a specific vNIC.
	// POST /api/v1/project/{project-id}/servers/{svr-id}/nics/{nic-id}/floatingip
	// Body: FloatingIPAssociateRequest (fip_id or create new)
	// Returns: 200 or 202 (if pending approval)
	// Errors: 400, 401, 403, 404, 409 (already associated), 500
	AssociateFloatingIP(ctx context.Context, nicID string, req *FloatingIPAssociateRequest) error
}

// ServerVolumeOperations defines operations on server volumes (sub-resource).
type ServerVolumeOperations interface {
	// List lists all volume attachments for the server.
	// GET /api/v1/project/{project-id}/servers/{svr-id}/volumes
	// Returns: 200 + array of VolumeAttachment
	// Errors: 401, 403, 404, 500
	List(ctx context.Context) ([]*VolumeAttachment, error)

	// Attach attaches a volume to the server.
	// POST /api/v1/project/{project-id}/servers/{svr-id}/volumes/{vol-id}
	// Returns: 200 or 202 (async)
	// Errors: 400, 401, 403, 404, 409 (volume in use), 500
	Attach(ctx context.Context, volumeID string) error

	// Detach detaches a volume from the server.
	// DELETE /api/v1/project/{project-id}/servers/{svr-id}/volumes/{vol-id}
	// Returns: 204 or 202 (async)
	// Errors: 401, 403, 404, 500
	Detach(ctx context.Context, volumeID string) error
}

// ServerCreateRequest is the body for Create.
type ServerCreateRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	FlavorID    string             `json:"flavor_id"`
	ImageID     string             `json:"image_id"`
	NICs        []ServerNICRequest `json:"nics"`
	SGIDs       []string           `json:"sg_ids"`
	KeypairID   string             `json:"keypair_id,omitempty"`
	Password    string             `json:"password,omitempty"`    // Base64
	BootScript  string             `json:"boot_script,omitempty"` // Base64
	// Volume fields if applicable
}

// ServerNICRequest specifies a NIC for server creation.
type ServerNICRequest struct {
	NetworkID string `json:"network_id"`
	FixedIP   string `json:"fixed_ip,omitempty"`
}

// ServerUpdateRequest is the body for Update.
type ServerUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// ServerActionRequest is the body for Action.
type ServerActionRequest struct {
	Action     string `json:"action"`                // start, stop, reboot, resize, extend_root, get_pwd, approve, reject
	RebootType string `json:"reboot_type,omitempty"` // hard or soft (for reboot)
	FlavorID   string `json:"flavor_id,omitempty"`   // for resize
	RootSize   int    `json:"root_size,omitempty"`   // for extend_root
	PrivateKey string `json:"private_key,omitempty"` // Base64 for get_pwd
}

// ServerMetricsRequest specifies query parameters for Metrics.
type ServerMetricsRequest struct {
	Type        string // cpu, memory, disk, network
	Start       int64  // Unix timestamp
	End         int64  // Unix timestamp
	Granularity int    // seconds
}

// ServerMetricsResponse is the response from Metrics.
type ServerMetricsResponse struct {
	Type   string        `json:"type"`
	Series []MetricPoint `json:"series"`
}

// MetricPoint is a single time-series data point.
type MetricPoint struct {
	Timestamp int64   `json:"timestamp"`
	Value     float64 `json:"value"`
}

// ServerNIC represents a vNIC attached to a server.
type ServerNIC struct {
	ID        string   `json:"id"`
	NetworkID string   `json:"network_id"`
	FixedIPs  []string `json:"fixed_ips"`
	MACAddr   string   `json:"mac_address"`
	SGIDs     []string `json:"sg_ids"`
}

// ServerNICCreateRequest is the body for AddServerNIC.
type ServerNICCreateRequest struct {
	NetworkID string   `json:"network_id"`
	SGIDs     []string `json:"sg_ids"`
	FixedIP   string   `json:"fixed_ip,omitempty"`
}

// ServerNICUpdateRequest is the body for UpdateServerNIC.
type ServerNICUpdateRequest struct {
	SGIDs []string `json:"sg_ids"`
}

// FloatingIPAssociateRequest is the body for AssociateFloatingIPToNIC.
type FloatingIPAssociateRequest struct {
	FloatingIPID string `json:"fip_id,omitempty"` // existing FIP; omit to create new
}

// VNCURLResponse is the response from VNCURL.
type VNCURLResponse struct {
	URL string `json:"url"`
}

// VolumeAttachment represents a volume attached to a server.
type VolumeAttachment struct {
	VolumeID string `json:"volume_id"`
	ServerID string `json:"server_id"`
	Device   string `json:"device,omitempty"`
}
