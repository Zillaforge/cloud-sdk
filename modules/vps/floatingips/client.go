package floatingips

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
)

// Client provides operations for managing floating IPs.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new floating IPs client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves all floating IPs for the project.
// GET /api/v1/project/{project-id}/floatingips
func (c *Client) List(ctx context.Context, opts *floatingips.ListFloatingIPsOptions) (*floatingips.FloatingIPListResponse, error) {
	path := fmt.Sprintf("%s/floatingips", c.basePath)

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	// Add query parameters if provided
	if opts != nil && opts.Status != "" {
		req.Path = fmt.Sprintf("%s?status=%s", path, opts.Status)
	}

	var response floatingips.FloatingIPListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, err
	}

	return &response, nil
}

// Create allocates a new floating IP.
// POST /api/v1/project/{project-id}/floatingips
func (c *Client) Create(ctx context.Context, req *floatingips.FloatingIPCreateRequest) (*floatingips.FloatingIP, error) {
	path := fmt.Sprintf("%s/floatingips", c.basePath)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var floatingIP floatingips.FloatingIP
	if err := c.baseClient.Do(ctx, httpReq, &floatingIP); err != nil {
		return nil, err
	}

	return &floatingIP, nil
}

// Get retrieves a specific floating IP.
// GET /api/v1/project/{project-id}/floatingips/{fip-id}
func (c *Client) Get(ctx context.Context, fipID string) (*floatingips.FloatingIP, error) {
	path := fmt.Sprintf("%s/floatingips/%s", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var floatingIP floatingips.FloatingIP
	if err := c.baseClient.Do(ctx, httpReq, &floatingIP); err != nil {
		return nil, err
	}

	return &floatingIP, nil
}

// Update updates floating IP description.
// PUT /api/v1/project/{project-id}/floatingips/{fip-id}
func (c *Client) Update(ctx context.Context, fipID string, req *floatingips.FloatingIPUpdateRequest) (*floatingips.FloatingIP, error) {
	path := fmt.Sprintf("%s/floatingips/%s", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   req,
	}

	var floatingIP floatingips.FloatingIP
	if err := c.baseClient.Do(ctx, httpReq, &floatingIP); err != nil {
		return nil, err
	}

	return &floatingIP, nil
}

// Delete releases a floating IP.
// DELETE /api/v1/project/{project-id}/floatingips/{fip-id}
func (c *Client) Delete(ctx context.Context, fipID string) error {
	path := fmt.Sprintf("%s/floatingips/%s", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	return c.baseClient.Do(ctx, httpReq, nil)
}

// Approve approves a pending floating IP request (admin only).
// POST /api/v1/project/{project-id}/floatingips/{fip-id}/approve
func (c *Client) Approve(ctx context.Context, fipID string) error {
	path := fmt.Sprintf("%s/floatingips/%s/approve", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
	}

	return c.baseClient.Do(ctx, httpReq, nil)
}

// Reject rejects a pending floating IP request (admin only).
// POST /api/v1/project/{project-id}/floatingips/{fip-id}/reject
func (c *Client) Reject(ctx context.Context, fipID string) error {
	path := fmt.Sprintf("%s/floatingips/%s/reject", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
	}

	return c.baseClient.Do(ctx, httpReq, nil)
}

// Disassociate disassociates a floating IP from its port.
// POST /api/v1/project/{project-id}/floatingips/{fip-id}/disassociate
func (c *Client) Disassociate(ctx context.Context, fipID string) error {
	path := fmt.Sprintf("%s/floatingips/%s/disassociate", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
	}

	return c.baseClient.Do(ctx, httpReq, nil)
}
