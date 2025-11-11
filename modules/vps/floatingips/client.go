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
// Returns a slice of FloatingIP directly (not wrapped in an object).
// This is a breaking change from the old "items" field wrapper.
func (c *Client) List(ctx context.Context, opts *floatingips.ListFloatingIPsOptions) ([]*floatingips.FloatingIP, error) {
	path := fmt.Sprintf("%s/floatingips", c.basePath)

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	// Add query parameters if provided
	if opts != nil {
		queryParams := ""
		if opts.Status != "" {
			queryParams += fmt.Sprintf("status=%s", opts.Status)
		}
		if opts.UserID != "" {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("user_id=%s", opts.UserID)
		}
		if opts.DeviceType != "" {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("device_type=%s", opts.DeviceType)
		}
		if opts.DeviceID != "" {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("device_id=%s", opts.DeviceID)
		}
		if opts.ExtNetID != "" {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("extnet_id=%s", opts.ExtNetID)
		}
		if opts.Address != "" {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("address=%s", opts.Address)
		}
		if opts.Name != "" {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += fmt.Sprintf("name=%s", opts.Name)
		}
		if opts.Detail {
			if queryParams != "" {
				queryParams += "&"
			}
			queryParams += "detail=true"
		}
		if queryParams != "" {
			req.Path = fmt.Sprintf("%s?%s", path, queryParams)
		}
	}

	// The API returns a wrapper object with "floatingips" field containing the array
	// We unmarshal into FloatingIPListResponse then return just the slice
	var response floatingips.FloatingIPListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list floatingips: %w", err)
	}

	return response.FloatingIPs, nil
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
		return nil, fmt.Errorf("failed to create floatingip: %w", err)
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
		return nil, fmt.Errorf("failed to get floatingip: %w", err)
	}

	return &floatingIP, nil
}

// Update updates floating IP fields (name, description).
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
		return nil, fmt.Errorf("failed to update floatingip: %w", err)
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

	if err := c.baseClient.Do(ctx, httpReq, nil); err != nil {
		return fmt.Errorf("failed to delete floatingip: %w", err)
	}
	return nil
}

// Approve approves a pending floating IP request (admin only).
// POST /api/v1/project/{project-id}/floatingips/{fip-id}/approve
func (c *Client) Approve(ctx context.Context, fipID string) error {
	path := fmt.Sprintf("%s/floatingips/%s/approve", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, httpReq, nil); err != nil {
		return fmt.Errorf("failed to approve floatingip: %w", err)
	}
	return nil
}

// Reject rejects a pending floating IP request (admin only).
// POST /api/v1/project/{project-id}/floatingips/{fip-id}/reject
func (c *Client) Reject(ctx context.Context, fipID string) error {
	path := fmt.Sprintf("%s/floatingips/%s/reject", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, httpReq, nil); err != nil {
		return fmt.Errorf("failed to reject floatingip: %w", err)
	}
	return nil
}

// Disassociate disassociates a floating IP from its port.
// POST /api/v1/project/{project-id}/floatingips/{fip-id}/disassociate
func (c *Client) Disassociate(ctx context.Context, fipID string) error {
	path := fmt.Sprintf("%s/floatingips/%s/disassociate", c.basePath, fipID)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
	}

	var floatingIP floatingips.FloatingIP
	if err := c.baseClient.Do(ctx, httpReq, &floatingIP); err != nil {
		return fmt.Errorf("failed to disassociate floatingip: %w", err)
	}
	return nil
}
