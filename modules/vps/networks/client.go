package networks

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// Client handles network-related operations.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new networks client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves all networks for the project with optional filters.
// GET /api/v1/project/{project-id}/networks
func (c *Client) List(ctx context.Context, opts *networks.ListNetworksOptions) ([]*NetworkResource, error) {
	path := c.basePath + "/networks"

	// Build query parameters
	query := url.Values{}
	var queryParts []string
	if opts != nil {
		if opts.Name != "" {
			query.Set("name", opts.Name)
			queryParts = append(queryParts, "name="+opts.Name)
		}
		if opts.UserID != "" {
			query.Set("user_id", opts.UserID)
			queryParts = append(queryParts, "user_id="+opts.UserID)
		}
		if opts.Status != "" {
			query.Set("status", opts.Status)
			queryParts = append(queryParts, "status="+opts.Status)
		}
		if opts.RouterID != "" {
			query.Set("router_id", opts.RouterID)
			queryParts = append(queryParts, "router_id="+opts.RouterID)
		}
		if opts.Detail != nil {
			if *opts.Detail {
				queryParts = append(queryParts, "detail=true")
			} else {
				queryParts = append(queryParts, "detail=false")
			}
		}
	}

	if len(queryParts) > 0 {
		path += "?" + strings.Join(queryParts, "&")
	}

	req := &internalhttp.Request{
		Method:  "GET",
		Path:    path,
		Headers: make(map[string]string),
	}

	var response networks.NetworkListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list networks: %w", err)
	}

	// Wrap networks in NetworkResource
	networkResources := make([]*NetworkResource, len(response.Networks))
	for i, network := range response.Networks {
		networkResources[i] = &NetworkResource{
			Network: network,
			portOps: &PortsClient{
				baseClient: c.baseClient,
				projectID:  c.projectID,
				networkID:  network.ID,
			},
		}
	}

	return networkResources, nil
}

// Create creates a new network.
// POST /api/v1/project/{project-id}/networks
func (c *Client) Create(ctx context.Context, req *networks.NetworkCreateRequest) (*NetworkResource, error) {
	path := fmt.Sprintf("/api/v1/project/%s/networks", c.projectID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var network networks.Network
	if err := c.baseClient.Do(ctx, httpReq, &network); err != nil {
		return nil, fmt.Errorf("failed to create network: %w", err)
	}

	// Wrap in NetworkResource with sub-resource operations
	return &NetworkResource{
		Network: &network,
		portOps: &PortsClient{
			baseClient: c.baseClient,
			projectID:  c.projectID,
			networkID:  network.ID,
		},
	}, nil
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
		return nil, fmt.Errorf("failed to get network %s: %w", networkID, err)
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
func (c *Client) Update(ctx context.Context, networkID string, req *networks.NetworkUpdateRequest) (*NetworkResource, error) {
	path := fmt.Sprintf("/api/v1/project/%s/networks/%s", c.projectID, networkID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   req,
	}

	var network networks.Network
	if err := c.baseClient.Do(ctx, httpReq, &network); err != nil {
		return nil, fmt.Errorf("failed to update network %s: %w", networkID, err)
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
		return fmt.Errorf("failed to delete network %s: %w", networkID, err)
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
