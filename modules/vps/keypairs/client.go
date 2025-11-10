package keypairs

import (
	"context"
	"fmt"
	"net/url"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/keypairs"
)

// Client provides operations for managing keypairs.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new keypairs client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves a list of keypairs with optional filtering.
// GET /api/v1/project/{project-id}/keypairs
func (c *Client) List(ctx context.Context, opts *keypairs.ListKeypairsOptions) ([]*keypairs.Keypair, error) {
	path := c.basePath + "/keypairs"

	// Build query parameters
	query := url.Values{}
	if opts != nil {
		if opts.Name != "" {
			query.Set("name", opts.Name)
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response keypairs.KeypairListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list keypairs: %w", err)
	}

	// Convert slice of values to slice of pointers for direct access
	keypairs := make([]*keypairs.Keypair, len(response.Keypairs))
	for i := range response.Keypairs {
		keypairs[i] = &response.Keypairs[i]
	}

	return keypairs, nil
}

// Create creates a new keypair.
// POST /api/v1/project/{project-id}/keypairs
func (c *Client) Create(ctx context.Context, req *keypairs.KeypairCreateRequest) (*keypairs.Keypair, error) {
	path := fmt.Sprintf("%s/keypairs", c.basePath)

	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var keypair keypairs.Keypair
	if err := c.baseClient.Do(ctx, httpReq, &keypair); err != nil {
		return nil, fmt.Errorf("failed to create keypair: %w", err)
	}

	return &keypair, nil
}

// Get retrieves details of a specific keypair by ID.
// GET /api/v1/project/{project-id}/keypairs/{keypair-id}
func (c *Client) Get(ctx context.Context, keypairID string) (*keypairs.Keypair, error) {
	path := fmt.Sprintf("%s/keypairs/%s", c.basePath, keypairID)

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var keypair keypairs.Keypair
	if err := c.baseClient.Do(ctx, req, &keypair); err != nil {
		return nil, fmt.Errorf("failed to get keypair %s: %w", keypairID, err)
	}

	return &keypair, nil
}

// Update updates keypair description.
// PUT /api/v1/project/{project-id}/keypairs/{keypair-id}
func (c *Client) Update(ctx context.Context, keypairID string, req *keypairs.KeypairUpdateRequest) (*keypairs.Keypair, error) {
	path := fmt.Sprintf("%s/keypairs/%s", c.basePath, keypairID)

	httpReq := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   req,
	}

	var keypair keypairs.Keypair
	if err := c.baseClient.Do(ctx, httpReq, &keypair); err != nil {
		return nil, fmt.Errorf("failed to update keypair %s: %w", keypairID, err)
	}

	return &keypair, nil
}

// Delete deletes a keypair.
// DELETE /api/v1/project/{project-id}/keypairs/{keypair-id}
func (c *Client) Delete(ctx context.Context, keypairID string) error {
	path := fmt.Sprintf("%s/keypairs/%s", c.basePath, keypairID)

	httpReq := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, httpReq, nil); err != nil {
		return fmt.Errorf("failed to delete keypair %s: %w", keypairID, err)
	}

	return nil
}
