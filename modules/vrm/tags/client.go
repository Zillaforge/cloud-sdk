package tags

import (
	"context"
	"fmt"
	"net/url"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	tagmod "github.com/Zillaforge/cloud-sdk/models/vrm/tags"
)

// Client provides access to tag operations for a specific project.
// It handles all CRUD operations for VRM tags including namespace support
// for multi-tenant scenarios. All operations are scoped to the configured project.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new tag operations client.
// It is configured with the base HTTP client, project ID, and API base path
// for performing tag operations within the project scope.
func NewClient(baseClient *internalhttp.Client, projectID, basePath string) *Client {
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves all tags in the project.
// GET /api/v1/project/{project-id}/tags
// Supports pagination via limit/offset and filtering via where conditions.
// Uses default namespace handling (no X-Namespace header).
func (c *Client) List(ctx context.Context, opts *tagmod.ListTagsOptions) ([]*tagmod.Tag, error) {
	if opts == nil {
		opts = &tagmod.ListTagsOptions{}
	}
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid list options: %w", err)
	}

	path := c.basePath + "/tags"

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

	var resp tagmod.ListTagsResponse
	if err := c.baseClient.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to list tags: %w", err)
	}

	return resp.Tags, nil
}

// ListByRepository retrieves all tags in a specific repository.
// GET /api/v1/project/{project-id}/repository/{repository-id}/tags
// Supports pagination via limit/offset and filtering via where conditions.
// Uses default namespace handling (no X-Namespace header).
func (c *Client) ListByRepository(ctx context.Context, repositoryID string, opts *tagmod.ListTagsOptions) ([]*tagmod.Tag, error) {
	if repositoryID == "" {
		return nil, fmt.Errorf("repository ID cannot be empty")
	}
	if opts == nil {
		opts = &tagmod.ListTagsOptions{}
	}
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid list options: %w", err)
	}

	path := c.basePath + "/repository/" + url.PathEscape(repositoryID) + "/tags"

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

	var resp tagmod.ListTagsResponse
	if err := c.baseClient.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to list tags by repository %s: %w", repositoryID, err)
	}

	return resp.Tags, nil
}

// Create creates a new tag in a repository.
// POST /api/v1/project/{project-id}/repository/{repository-id}/tag
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Create(ctx context.Context, repositoryID string, req *tagmod.CreateTagRequest) (*tagmod.Tag, error) {
	return c.CreateWithNamespace(ctx, repositoryID, req, "")
}

// CreateWithNamespace creates a new tag in a repository with optional namespace header.
// POST /api/v1/project/{project-id}/repository/{repository-id}/tag
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) CreateWithNamespace(ctx context.Context, repositoryID string, req *tagmod.CreateTagRequest, namespace string) (*tagmod.Tag, error) {
	if repositoryID == "" {
		return nil, fmt.Errorf("repository ID cannot be empty")
	}
	if req == nil {
		return nil, fmt.Errorf("create request cannot be nil")
	}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid create request: %w", err)
	}

	path := c.basePath + "/repository/" + url.PathEscape(repositoryID) + "/tag"

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

	var tag tagmod.Tag
	if err := c.baseClient.Do(ctx, httpReq, &tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return &tag, nil
}

// Get retrieves a specific tag.
// GET /api/v1/project/{project-id}/tag/{tag-id}
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Get(ctx context.Context, tagID string) (*tagmod.Tag, error) {
	return c.GetWithNamespace(ctx, tagID, "")
}

// GetWithNamespace retrieves a specific tag with optional namespace header.
// GET /api/v1/project/{project-id}/tag/{tag-id}
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) GetWithNamespace(ctx context.Context, tagID string, namespace string) (*tagmod.Tag, error) {
	if tagID == "" {
		return nil, fmt.Errorf("tag ID cannot be empty")
	}

	path := c.basePath + "/tag/" + url.PathEscape(tagID)

	headers := make(map[string]string)
	if namespace != "" {
		headers["X-Namespace"] = namespace
	}

	req := &internalhttp.Request{
		Method:  "GET",
		Path:    path,
		Headers: headers,
	}

	var tag tagmod.Tag
	if err := c.baseClient.Do(ctx, req, &tag); err != nil {
		return nil, fmt.Errorf("failed to get tag %s: %w", tagID, err)
	}

	return &tag, nil
}

// Update updates an existing tag.
// PUT /api/v1/project/{project-id}/tag/{tag-id}
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Update(ctx context.Context, tagID string, req *tagmod.UpdateTagRequest) (*tagmod.Tag, error) {
	return c.UpdateWithNamespace(ctx, tagID, req, "")
}

// UpdateWithNamespace updates an existing tag with optional namespace header.
// PUT /api/v1/project/{project-id}/tag/{tag-id}
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) UpdateWithNamespace(ctx context.Context, tagID string, req *tagmod.UpdateTagRequest, namespace string) (*tagmod.Tag, error) {
	if tagID == "" {
		return nil, fmt.Errorf("tag ID cannot be empty")
	}
	if req == nil {
		return nil, fmt.Errorf("update request cannot be nil")
	}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid update request: %w", err)
	}

	path := c.basePath + "/tag/" + url.PathEscape(tagID)

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

	var tag tagmod.Tag
	if err := c.baseClient.Do(ctx, httpReq, &tag); err != nil {
		return nil, fmt.Errorf("failed to update tag %s: %w", tagID, err)
	}

	return &tag, nil
}

// Delete deletes a tag.
// DELETE /api/v1/project/{project-id}/tag/{tag-id}
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Delete(ctx context.Context, tagID string) error {
	return c.DeleteWithNamespace(ctx, tagID, "")
}

// DeleteWithNamespace deletes a tag with optional namespace header.
// DELETE /api/v1/project/{project-id}/tag/{tag-id}
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) DeleteWithNamespace(ctx context.Context, tagID string, namespace string) error {
	if tagID == "" {
		return fmt.Errorf("tag ID cannot be empty")
	}

	path := c.basePath + "/tag/" + url.PathEscape(tagID)

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
		return fmt.Errorf("failed to delete tag %s: %w", tagID, err)
	}

	return nil
}
