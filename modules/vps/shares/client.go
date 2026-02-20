package shares

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	sharesmodel "github.com/Zillaforge/cloud-sdk/models/vps/shares"
)

// Client provides operations for managing shares.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new shares client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	return &Client{baseClient: baseClient, projectID: projectID, basePath: "/api/v1/project/" + projectID}
}

// List retrieves a list of shares with optional filtering.
// GET /api/v1/project/{project-id}/shares
func (c *Client) List(ctx context.Context, opts *sharesmodel.ListSharesOptions) ([]*sharesmodel.Share, error) {
	path := c.basePath + "/shares"

	query := url.Values{}
	if opts != nil {
		if opts.Name != "" {
			query.Set("name", opts.Name)
		}
		if opts.UserID != "" {
			query.Set("user_id", opts.UserID)
		}
		if opts.Detail {
			query.Set("detail", strconv.FormatBool(opts.Detail))
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req := &internalhttp.Request{Method: "GET", Path: path}

	var response sharesmodel.ShareListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list shares: %w", err)
	}

	return response.Shares, nil
}

// Get retrieves details of a specific share by ID.
// GET /api/v1/project/{project-id}/shares/{share-id}
func (c *Client) Get(ctx context.Context, shareID string) (*sharesmodel.Share, error) {
	path := fmt.Sprintf("%s/shares/%s", c.basePath, shareID)

	req := &internalhttp.Request{Method: "GET", Path: path}

	var response sharesmodel.Share
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to get share %s: %w", shareID, err)
	}

	return &response, nil
}

// Create creates a new share.
// POST /api/v1/project/{project-id}/shares
func (c *Client) Create(ctx context.Context, request *sharesmodel.CreateShareRequest) (*sharesmodel.Share, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := c.basePath + "/shares"

	req := &internalhttp.Request{Method: "POST", Path: path, Body: request}

	var response sharesmodel.Share
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to create share: %w", err)
	}

	return &response, nil
}

// Update updates share metadata (name, description).
// PUT /api/v1/project/{project-id}/shares/{share-id}
func (c *Client) Update(ctx context.Context, shareID string, request *sharesmodel.UpdateShareRequest) (*sharesmodel.Share, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("%s/shares/%s", c.basePath, shareID)

	req := &internalhttp.Request{Method: "PUT", Path: path, Body: request}

	var response sharesmodel.Share
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to update share %s: %w", shareID, err)
	}

	return &response, nil
}

// Delete deletes a share.
// DELETE /api/v1/project/{project-id}/shares/{share-id}
func (c *Client) Delete(ctx context.Context, shareID string) error {
	path := fmt.Sprintf("%s/shares/%s", c.basePath, shareID)

	req := &internalhttp.Request{Method: "DELETE", Path: path}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete share %s: %w", shareID, err)
	}

	return nil
}

// Extend increases the size of a share.
// POST /api/v1/project/{project-id}/shares/{share-id}/extend
func (c *Client) Extend(ctx context.Context, shareID string, request *sharesmodel.ExtendShareRequest) error {
	if err := request.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("%s/shares/%s/extend", c.basePath, shareID)

	req := &internalhttp.Request{Method: "POST", Path: path, Body: request}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to extend share %s: %w", shareID, err)
	}

	return nil
}

// AddRule creates an access rule on a share.
// POST /api/v1/project/{project-id}/shares/{share-id}/rules
func (c *Client) AddRule(ctx context.Context, shareID string, request *sharesmodel.AddShareRuleRequest) (*sharesmodel.Share, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("%s/shares/%s/rules", c.basePath, shareID)

	req := &internalhttp.Request{Method: "POST", Path: path, Body: request}

	var response sharesmodel.Share
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to add rule to share %s: %w", shareID, err)
	}

	return &response, nil
}

// RemoveRule deletes an access rule from a share.
// DELETE /api/v1/project/{project-id}/shares/{share-id}/rules/{rule-id}
func (c *Client) RemoveRule(ctx context.Context, shareID, ruleID string) (*sharesmodel.Share, error) {
	path := fmt.Sprintf("%s/shares/%s/rules/%s", c.basePath, shareID, ruleID)

	req := &internalhttp.Request{Method: "DELETE", Path: path}

	var response sharesmodel.Share
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to remove rule %s from share %s: %w", ruleID, shareID, err)
	}

	return &response, nil
}
