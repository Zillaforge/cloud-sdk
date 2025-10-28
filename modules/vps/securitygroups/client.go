// Package securitygroups provides operations for managing VPS security groups.
package securitygroups

import (
	"context"
	"fmt"
	"net/url"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/securitygroups"
)

// Client provides operations for managing security groups.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new security groups client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves a list of security groups with optional filtering.
// GET /api/v1/project/{project-id}/security_groups
func (c *Client) List(ctx context.Context, opts *securitygroups.ListSecurityGroupsOptions) (*securitygroups.SecurityGroupListResponse, error) {
	path := c.basePath + "/security_groups"

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

	var response securitygroups.SecurityGroupListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list security groups: %w", err)
	}

	return &response, nil
}

// Create creates a new security group.
// POST /api/v1/project/{project-id}/security_groups
func (c *Client) Create(ctx context.Context, req securitygroups.SecurityGroupCreateRequest) (*securitygroups.SecurityGroup, error) {
	path := fmt.Sprintf("%s/security_groups", c.basePath)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var sg securitygroups.SecurityGroup
	if err := c.baseClient.Do(ctx, httpReq, &sg); err != nil {
		return nil, fmt.Errorf("failed to create security group: %w", err)
	}

	return &sg, nil
}

// Get retrieves details of a specific security group by ID.
// GET /api/v1/project/{project-id}/security_groups/{sg-id}
func (c *Client) Get(ctx context.Context, sgID string) (*securitygroups.SecurityGroup, error) {
	path := fmt.Sprintf("%s/security_groups/%s", c.basePath, sgID)

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var sg securitygroups.SecurityGroup
	if err := c.baseClient.Do(ctx, req, &sg); err != nil {
		return nil, fmt.Errorf("failed to get security group %s: %w", sgID, err)
	}

	return &sg, nil
}

// Update updates security group name and/or description.
// PUT /api/v1/project/{project-id}/security_groups/{sg-id}
func (c *Client) Update(ctx context.Context, sgID string, req securitygroups.SecurityGroupUpdateRequest) (*securitygroups.SecurityGroup, error) {
	path := fmt.Sprintf("%s/security_groups/%s", c.basePath, sgID)

	httpReq := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   req,
	}

	var sg securitygroups.SecurityGroup
	if err := c.baseClient.Do(ctx, httpReq, &sg); err != nil {
		return nil, fmt.Errorf("failed to update security group %s: %w", sgID, err)
	}

	return &sg, nil
}

// Delete removes a security group.
// DELETE /api/v1/project/{project-id}/security_groups/{sg-id}
func (c *Client) Delete(ctx context.Context, sgID string) error {
	path := fmt.Sprintf("%s/security_groups/%s", c.basePath, sgID)

	req := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete security group %s: %w", sgID, err)
	}

	return nil
}
