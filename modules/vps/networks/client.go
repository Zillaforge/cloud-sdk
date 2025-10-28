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
	headers := make(map[string]string)
	if opts != nil && opts.Name != "" {
		// Add query params via path - internal client will handle URL encoding
		path = fmt.Sprintf("%s?name=%s", path, opts.Name)
	}

	// Make request
	req := &internalhttp.Request{
		Method:  "GET",
		Path:    path,
		Headers: headers,
	}

	var response networks.NetworkListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
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
