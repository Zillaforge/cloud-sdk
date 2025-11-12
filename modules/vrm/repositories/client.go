package repositories

import (
	"context"
	"fmt"
	"net/url"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	repmod "github.com/Zillaforge/cloud-sdk/models/vrm/repositories"
)

// Client provides access to repository operations for a specific project.
// It handles all CRUD operations for VRM repositories including namespace support
// for multi-tenant scenarios. All operations are scoped to the configured project.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new repository operations client.
// It is configured with the base HTTP client, project ID, and API base path
// for performing repository operations within the project scope.
func NewClient(baseClient *internalhttp.Client, projectID, basePath string) *Client {
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves all repositories in the project.
// GET /api/v1/project/{project-id}/repositories
// Supports pagination via limit/offset and filtering via where conditions.
// Uses default namespace handling (no X-Namespace header).
func (c *Client) List(ctx context.Context, opts *repmod.ListRepositoriesOptions) ([]*repmod.Repository, error) {
	if opts == nil {
		opts = &repmod.ListRepositoriesOptions{}
	}
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid list options: %w", err)
	}

	path := c.basePath + "/repositories"

	// Build query parameters
	query := url.Values{}
	if opts.Limit > 0 {
		query.Set("limit", fmt.Sprintf("%d", opts.Limit))
	}
	if opts.Offset > 0 {
		query.Set("offset", fmt.Sprintf("%d", opts.Offset))
	}
	for _, filter := range opts.Where {
		query.Add("where", filter)
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	headers := make(map[string]string)
	if opts.Namespace != "" {
		headers["X-Namespace"] = opts.Namespace
	}

	req := &internalhttp.Request{
		Method:  "GET",
		Path:    path,
		Headers: headers,
	}

	var repos []*repmod.Repository
	if err := c.baseClient.Do(ctx, req, &repos); err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	return repos, nil
}

// Create creates a new repository.
// POST /api/v1/project/{project-id}/repository
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Create(ctx context.Context, req *repmod.CreateRepositoryRequest) (*repmod.Repository, error) {
	return c.CreateWithNamespace(ctx, req, "")
}

// CreateWithNamespace creates a new repository with optional namespace header.
// POST /api/v1/project/{project-id}/repository
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) CreateWithNamespace(ctx context.Context, req *repmod.CreateRepositoryRequest, namespace string) (*repmod.Repository, error) {
	if req == nil {
		return nil, fmt.Errorf("create request cannot be nil")
	}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid create request: %w", err)
	}

	path := c.basePath + "/repository"

	headers := make(map[string]string)
	if namespace != "" {
		headers["X-Namespace"] = namespace
	}

	httpReq := &internalhttp.Request{
		Method:  "POST",
		Path:    path,
		Body:    req,
		Headers: headers,
	}

	var repo repmod.Repository
	if err := c.baseClient.Do(ctx, httpReq, &repo); err != nil {
		return nil, fmt.Errorf("failed to create repository: %w", err)
	}

	return &repo, nil
}

// Get retrieves a specific repository.
// GET /api/v1/project/{project-id}/repository/{repository-id}
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Get(ctx context.Context, repositoryID string) (*repmod.Repository, error) {
	return c.GetWithNamespace(ctx, repositoryID, "")
}

// GetWithNamespace retrieves a specific repository with optional namespace header.
// GET /api/v1/project/{project-id}/repository/{repository-id}
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) GetWithNamespace(ctx context.Context, repositoryID string, namespace string) (*repmod.Repository, error) {
	path := c.basePath + "/repository/" + url.PathEscape(repositoryID)

	headers := make(map[string]string)
	if namespace != "" {
		headers["X-Namespace"] = namespace
	}

	req := &internalhttp.Request{
		Method:  "GET",
		Path:    path,
		Headers: headers,
	}

	var repo repmod.Repository
	if err := c.baseClient.Do(ctx, req, &repo); err != nil {
		return nil, fmt.Errorf("failed to get repository %s: %w", repositoryID, err)
	}

	return &repo, nil
}

// Update updates an existing repository.
// PUT /api/v1/project/{project-id}/repository/{repository-id}
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Update(ctx context.Context, repositoryID string, req *repmod.UpdateRepositoryRequest) (*repmod.Repository, error) {
	return c.UpdateWithNamespace(ctx, repositoryID, req, "")
}

// UpdateWithNamespace updates an existing repository with optional namespace header.
// PUT /api/v1/project/{project-id}/repository/{repository-id}
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) UpdateWithNamespace(ctx context.Context, repositoryID string, req *repmod.UpdateRepositoryRequest, namespace string) (*repmod.Repository, error) {
	if req == nil {
		return nil, fmt.Errorf("update request cannot be nil")
	}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid update request: %w", err)
	}

	path := c.basePath + "/repository/" + url.PathEscape(repositoryID)

	headers := make(map[string]string)
	if namespace != "" {
		headers["X-Namespace"] = namespace
	}

	httpReq := &internalhttp.Request{
		Method:  "PUT",
		Path:    path,
		Body:    req,
		Headers: headers,
	}

	var repo repmod.Repository
	if err := c.baseClient.Do(ctx, httpReq, &repo); err != nil {
		return nil, fmt.Errorf("failed to update repository %s: %w", repositoryID, err)
	}

	return &repo, nil
}

// Delete deletes a repository.
// DELETE /api/v1/project/{project-id}/repository/{repository-id}
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Delete(ctx context.Context, repositoryID string) error {
	return c.DeleteWithNamespace(ctx, repositoryID, "")
}

// DeleteWithNamespace deletes a repository with optional namespace header.
// DELETE /api/v1/project/{project-id}/repository/{repository-id}
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) DeleteWithNamespace(ctx context.Context, repositoryID string, namespace string) error {
	path := c.basePath + "/repository/" + url.PathEscape(repositoryID)

	headers := make(map[string]string)
	if namespace != "" {
		headers["X-Namespace"] = namespace
	}

	req := &internalhttp.Request{
		Method:  "DELETE",
		Path:    path,
		Headers: headers,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete repository %s: %w", repositoryID, err)
	}

	return nil
}
