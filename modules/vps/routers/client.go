package routers

import (
	"context"
	"fmt"
	"net/url"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/routers"
)

// Client provides operations for managing routers.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new routers client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves a list of routers with optional filtering.
// GET /api/v1/project/{project-id}/routers
func (c *Client) List(ctx context.Context, opts *routers.ListRoutersOptions) (*routers.RouterListResponse, error) {
	path := c.basePath + "/routers"

	// Build query parameters
	query := url.Values{}
	if opts != nil {
		if opts.Name != "" {
			query.Set("name", opts.Name)
		}
		if opts.UserID != "" {
			query.Set("user_id", opts.UserID)
		}
		if opts.Detail {
			query.Set("detail", "true")
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response routers.RouterListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list routers: %w", err)
	}

	return &response, nil
}

// Create creates a new router.
// POST /api/v1/project/{project-id}/routers
func (c *Client) Create(ctx context.Context, req *routers.RouterCreateRequest) (*routers.Router, error) {
	path := fmt.Sprintf("%s/routers", c.basePath)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var router routers.Router
	if err := c.baseClient.Do(ctx, httpReq, &router); err != nil {
		return nil, fmt.Errorf("failed to create router: %w", err)
	}

	return &router, nil
}

// Get retrieves details of a specific router by ID.
// GET /api/v1/project/{project-id}/routers/{router-id}
// Returns a RouterResource that provides access to sub-resource operations.
func (c *Client) Get(ctx context.Context, routerID string) (*routers.RouterResource, error) {
	path := fmt.Sprintf("%s/routers/%s", c.basePath, routerID)

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var router routers.Router
	if err := c.baseClient.Do(ctx, req, &router); err != nil {
		return nil, fmt.Errorf("failed to get router %s: %w", routerID, err)
	}

	// Wrap with network operations
	networksOps := newRouterNetworksClient(c.baseClient, c.projectID, routerID)
	return &routers.RouterResource{
		Router:      &router,
		NetworksOps: networksOps,
	}, nil
}

// Update updates router name and/or description.
// PUT /api/v1/project/{project-id}/routers/{router-id}
func (c *Client) Update(ctx context.Context, routerID string, req *routers.RouterUpdateRequest) (*routers.Router, error) {
	path := fmt.Sprintf("%s/routers/%s", c.basePath, routerID)

	httpReq := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   req,
	}

	var router routers.Router
	if err := c.baseClient.Do(ctx, httpReq, &router); err != nil {
		return nil, fmt.Errorf("failed to update router %s: %w", routerID, err)
	}

	return &router, nil
}

// Delete removes a router.
// DELETE /api/v1/project/{project-id}/routers/{router-id}
func (c *Client) Delete(ctx context.Context, routerID string) error {
	path := fmt.Sprintf("%s/routers/%s", c.basePath, routerID)

	req := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete router %s: %w", routerID, err)
	}

	return nil
}

// SetState enables or disables a router.
// POST /api/v1/project/{project-id}/routers/{router-id}/action
func (c *Client) SetState(ctx context.Context, routerID string, req *routers.RouterSetStateRequest) error {
	path := fmt.Sprintf("%s/routers/%s/action", c.basePath, routerID)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	if err := c.baseClient.Do(ctx, httpReq, nil); err != nil {
		return fmt.Errorf("failed to set router %s state: %w", routerID, err)
	}

	return nil
}
