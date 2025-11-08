package networks

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// Client handles network-related operations.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
}

// NewClient creates a new networks client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
	}
}

// List retrieves all networks for the project with optional filters.
// GET /api/v1/project/{project-id}/networks
func (c *Client) List(ctx context.Context, opts *networks.ListNetworksOptions) (*networks.NetworkListResponse, error) {
	path := fmt.Sprintf("/api/v1/project/%s/networks", c.projectID)

	// Build query parameters
	if opts != nil {
		queryParams := []string{}
		if opts.Name != "" {
			queryParams = append(queryParams, fmt.Sprintf("name=%s", opts.Name))
		}
		if opts.UserID != "" {
			queryParams = append(queryParams, fmt.Sprintf("user_id=%s", opts.UserID))
		}
		if opts.Status != "" {
			queryParams = append(queryParams, fmt.Sprintf("status=%s", opts.Status))
		}
		if opts.RouterID != "" {
			queryParams = append(queryParams, fmt.Sprintf("router_id=%s", opts.RouterID))
		}
		if opts.Detail != nil {
			queryParams = append(queryParams, fmt.Sprintf("detail=%t", *opts.Detail))
		}
		if len(queryParams) > 0 {
			path = fmt.Sprintf("%s?%s", path, joinQueryParams(queryParams))
		}
	}

	// Make request
	req := &internalhttp.Request{
		Method:  "GET",
		Path:    path,
		Headers: make(map[string]string),
	}

	var response networks.NetworkListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// joinQueryParams joins query parameters with "&" separator.
func joinQueryParams(params []string) string {
	result := ""
	for i, param := range params {
		if i > 0 {
			result += "&"
		}
		result += param
	}
	return result
}

// Create creates a new network.
// POST /api/v1/project/{project-id}/networks
func (c *Client) Create(ctx context.Context, req *networks.NetworkCreateRequest) (*networks.Network, error) {
	path := fmt.Sprintf("/api/v1/project/%s/networks", c.projectID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var network networks.Network
	if err := c.baseClient.Do(ctx, httpReq, &network); err != nil {
		return nil, err
	}

	return &network, nil
}

// Get retrieves a specific network with sub-resource operations.
// GET /api/v1/project/{project-id}/networks/{net-id}
func (c *Client) Get(ctx context.Context, networkID string) (*NetworkResource, error) {
	path := fmt.Sprintf("/api/v1/project/%s/networks/%s", c.projectID, networkID)

	// Make request
	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var network networks.Network
	if err := c.baseClient.Do(ctx, req, &network); err != nil {
		return nil, err
	}

	// Wrap in NetworkResource with sub-resource operations
	return &NetworkResource{
		Network: &network,
		portOps: &PortsClient{
			baseClient: c.baseClient,
			projectID:  c.projectID,
			networkID:  networkID,
		},
	}, nil
}

// Update updates network name/description.
// PUT /api/v1/project/{project-id}/networks/{net-id}
func (c *Client) Update(ctx context.Context, networkID string, req *networks.NetworkUpdateRequest) (*networks.Network, error) {
	path := fmt.Sprintf("/api/v1/project/%s/networks/%s", c.projectID, networkID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   req,
	}

	var network networks.Network
	if err := c.baseClient.Do(ctx, httpReq, &network); err != nil {
		return nil, err
	}

	return &network, nil
}

// Delete deletes a network.
// DELETE /api/v1/project/{project-id}/networks/{net-id}
func (c *Client) Delete(ctx context.Context, networkID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/networks/%s", c.projectID, networkID)

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

// NetworkResource wraps a Network with sub-resource operations.
type NetworkResource struct {
	*networks.Network
	portOps PortOperations
}

// Ports returns the port operations for this network.
func (nr *NetworkResource) Ports() PortOperations {
	return nr.portOps
}

// PortOperations defines operations on network ports (sub-resource).
type PortOperations interface {
	List(ctx context.Context) ([]*networks.NetworkPort, error)
}
