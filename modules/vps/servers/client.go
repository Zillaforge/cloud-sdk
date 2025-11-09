package servers

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

// Client provides operations for managing servers.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
}

// NewClient creates a new servers client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
	}
}

// List retrieves all servers for the project with optional filters.
// GET /api/v1/project/{project-id}/servers
func (c *Client) List(ctx context.Context, opts *servers.ServersListRequest) (*servers.ServersListResponse, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers", c.projectID)

	// Build query parameters
	if opts != nil {
		query := ""
		if opts.Name != "" {
			query += fmt.Sprintf("name=%s&", opts.Name)
		}
		if opts.UserID != "" {
			query += fmt.Sprintf("user_id=%s&", opts.UserID)
		}
		if opts.Status != "" {
			query += fmt.Sprintf("status=%s&", opts.Status)
		}
		if opts.FlavorID != "" {
			query += fmt.Sprintf("flavor_id=%s&", opts.FlavorID)
		}
		if opts.ImageID != "" {
			query += fmt.Sprintf("image_id=%s&", opts.ImageID)
		}
		if opts.Detail {
			query += "detail=true&"
		}
		if query != "" {
			path = fmt.Sprintf("%s?%s", path, query[:len(query)-1]) // Remove trailing &
		}
	}

	// Make request
	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response servers.ServersListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Create provisions a new server instance.
// POST /api/v1/project/{project-id}/servers
func (c *Client) Create(ctx context.Context, req *servers.ServerCreateRequest) (*servers.Server, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers", c.projectID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var server servers.Server
	if err := c.baseClient.Do(ctx, httpReq, &server); err != nil {
		return nil, err
	}

	return &server, nil
}

// Get retrieves a specific server with sub-resource operations.
// GET /api/v1/project/{project-id}/servers/{svr-id}
func (c *Client) Get(ctx context.Context, serverID string) (*ServerResource, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s", c.projectID, serverID)

	// Make request
	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var server servers.Server
	if err := c.baseClient.Do(ctx, req, &server); err != nil {
		return nil, err
	}

	// Wrap in ServerResource with sub-resource operations
	return &ServerResource{
		Server: &server,
		nicOps: &NICsClient{
			baseClient: c.baseClient,
			projectID:  c.projectID,
			serverID:   serverID,
		},
		volumeOps: &VolumesClient{
			baseClient: c.baseClient,
			projectID:  c.projectID,
			serverID:   serverID,
		},
	}, nil
}

// Update modifies server name/description.
// PUT /api/v1/project/{project-id}/servers/{svr-id}
func (c *Client) Update(ctx context.Context, serverID string, req *servers.ServerUpdateRequest) (*servers.Server, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s", c.projectID, serverID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   req,
	}

	var server servers.Server
	if err := c.baseClient.Do(ctx, httpReq, &server); err != nil {
		return nil, err
	}

	return &server, nil
}

// Delete removes a server instance.
// DELETE /api/v1/project/{project-id}/servers/{svr-id}
func (c *Client) Delete(ctx context.Context, serverID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s", c.projectID, serverID)

	// Make request
	req := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// Action performs a control action on a server.
// POST /api/v1/project/{project-id}/servers/{svr-id}/action
func (c *Client) Action(ctx context.Context, serverID string, req *servers.ServerActionRequest) error {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/action", c.projectID, serverID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	if err := c.baseClient.Do(ctx, httpReq, nil); err != nil {
		return err
	}

	return nil
}

// Metrics retrieves time-series metrics for a server.
// GET /api/v1/project/{project-id}/servers/{svr-id}/metric
func (c *Client) Metrics(ctx context.Context, serverID string, req *servers.ServerMetricsRequest) (*servers.ServerMetricsResponse, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/metric", c.projectID, serverID)

	// Build query parameters
	if req != nil {
		query := ""
		if req.Type != "" {
			query += fmt.Sprintf("type=%s&", req.Type)
		}
		if req.Start > 0 {
			query += fmt.Sprintf("start=%d&", req.Start)
		}
		if req.Direction != "" {
			query += fmt.Sprintf("direction=%s&", req.Direction)
		}
		if req.RW != "" {
			query += fmt.Sprintf("rw=%s&", req.RW)
		}
		if req.Granularity > 0 {
			query += fmt.Sprintf("granularity=%d&", req.Granularity)
		}
		if query != "" {
			path = fmt.Sprintf("%s?%s", path, query[:len(query)-1]) // Remove trailing &
		}
	}

	// Make request
	httpReq := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response servers.ServerMetricsResponse
	if err := c.baseClient.Do(ctx, httpReq, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// GetVNCConsoleURL retrieves the VNC console URL for a server.
// GET /api/v1/project/{project-id}/servers/{svr-id}/vnc_url
func (c *Client) GetVNCConsoleURL(ctx context.Context, serverID string) (*servers.ServerConsoleURLResponse, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/vnc_url", c.projectID, serverID)

	// Make request
	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response servers.ServerConsoleURLResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// ServerResource wraps a Server with sub-resource operations.
type ServerResource struct {
	*servers.Server
	nicOps    NICOperations
	volumeOps VolumeOperations
}

// NICs returns the NIC operations for this server.
func (sr *ServerResource) NICs() NICOperations {
	return sr.nicOps
}

// Volumes returns the volume operations for this server.
func (sr *ServerResource) Volumes() VolumeOperations {
	return sr.volumeOps
}

// NICOperations defines operations on server NICs (sub-resource).
type NICOperations interface {
	List(ctx context.Context) (*servers.ServerNICsListResponse, error)
	Add(ctx context.Context, req *servers.ServerNICCreateRequest) (*servers.ServerNIC, error)
	Update(ctx context.Context, nicID string, req *servers.ServerNICUpdateRequest) (*servers.ServerNIC, error)
	Delete(ctx context.Context, nicID string) error
	AssociateFloatingIP(ctx context.Context, nicID string, req *servers.ServerNICAssociateFloatingIPRequest) (*servers.FloatingIPInfo, error)
}

// VolumeOperations defines operations on server volumes (sub-resource).
type VolumeOperations interface {
	List(ctx context.Context) ([]*servers.VolumeAttachment, error)
	Attach(ctx context.Context, volumeID string) error
	Detach(ctx context.Context, volumeID string) error
}
