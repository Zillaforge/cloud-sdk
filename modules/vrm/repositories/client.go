package repositories

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	repmod "github.com/Zillaforge/cloud-sdk/models/vrm/repositories"
	tagmod "github.com/Zillaforge/cloud-sdk/models/vrm/tags"
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
func (c *Client) List(ctx context.Context, opts *repmod.ListRepositoriesOptions) ([]*RepositoryResource, error) {
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

	var resp repmod.ListRepositoriesResponse
	if err := c.baseClient.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to list repositories: %w", err)
	}

	repos := resp.Repositories

	// Wrap repositories in RepositoryResource
	repoResources := make([]*RepositoryResource, len(repos))
	for i, repo := range repos {
		repoResources[i] = &RepositoryResource{
			Repository: repo,
			tagOps: &TagsClient{
				baseClient:   c.baseClient,
				repositoryID: repo.ID,
				basePath:     c.basePath,
			},
		}
	}

	return repoResources, nil
}

// Create creates a new repository.
// POST /api/v1/project/{project-id}/repository
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Create(ctx context.Context, req *repmod.CreateRepositoryRequest) (*RepositoryResource, error) {
	return c.CreateWithNamespace(ctx, req, "")
}

// CreateWithNamespace creates a new repository with optional namespace header.
// POST /api/v1/project/{project-id}/repository
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) CreateWithNamespace(ctx context.Context, req *repmod.CreateRepositoryRequest, namespace string) (*RepositoryResource, error) {
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

	// Wrap in RepositoryResource with sub-resource operations
	return &RepositoryResource{
		Repository: &repo,
		tagOps: &TagsClient{
			baseClient:   c.baseClient,
			repositoryID: repo.ID,
			basePath:     c.basePath,
		},
	}, nil
}

// Get retrieves a specific repository.
// GET /api/v1/project/{project-id}/repository/{repository-id}
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Get(ctx context.Context, repositoryID string) (*RepositoryResource, error) {
	return c.GetWithNamespace(ctx, repositoryID, "")
}

// GetWithNamespace retrieves a specific repository with optional namespace header.
// GET /api/v1/project/{project-id}/repository/{repository-id}
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) GetWithNamespace(ctx context.Context, repositoryID string, namespace string) (*RepositoryResource, error) {
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

	// Wrap in RepositoryResource with sub-resource operations
	return &RepositoryResource{
		Repository: &repo,
		tagOps: &TagsClient{
			baseClient:   c.baseClient,
			repositoryID: repositoryID,
			basePath:     c.basePath,
		},
	}, nil
}

// Update updates an existing repository.
// PUT /api/v1/project/{project-id}/repository/{repository-id}
// Uses default namespace handling (no X-Namespace header).
func (c *Client) Update(ctx context.Context, repositoryID string, req *repmod.UpdateRepositoryRequest) (*RepositoryResource, error) {
	return c.UpdateWithNamespace(ctx, repositoryID, req, "")
}

// UpdateWithNamespace updates an existing repository with optional namespace header.
// PUT /api/v1/project/{project-id}/repository/{repository-id}
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (c *Client) UpdateWithNamespace(ctx context.Context, repositoryID string, req *repmod.UpdateRepositoryRequest, namespace string) (*RepositoryResource, error) {
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

	// Wrap in RepositoryResource with sub-resource operations
	return &RepositoryResource{
		Repository: &repo,
		tagOps: &TagsClient{
			baseClient:   c.baseClient,
			repositoryID: repositoryID,
			basePath:     c.basePath,
		},
	}, nil
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

// Snapshot triggers a snapshot operation for a server and returns the associated repository resource.
// POST /api/v1/project/{project-id}/server/{server-id}/snapshot
func (c *Client) Snapshot(ctx context.Context, serverID string, req repmod.SnapshotRequester) (*RepositoryResource, error) {
	return c.SnapshotWithNamespace(ctx, serverID, req, "")
}

// SnapshotWithNamespace triggers a snapshot operation with optional namespace header for multi-tenant use cases.
// POST /api/v1/project/{project-id}/server/{server-id}/snapshot
func (c *Client) SnapshotWithNamespace(ctx context.Context, serverID string, req repmod.SnapshotRequester, namespace string) (*RepositoryResource, error) {
	if strings.TrimSpace(serverID) == "" {
		return nil, fmt.Errorf("server ID cannot be empty")
	}
	if req == nil {
		return nil, fmt.Errorf("snapshot request cannot be nil")
	}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid snapshot request: %w", err)
	}

	createReq := req.ToCreateSnapshotRequest()

	path := c.basePath + "/server/" + url.PathEscape(serverID) + "/snapshot"

	headers := make(map[string]string)
	if namespace != "" {
		headers["X-Namespace"] = namespace
	}

	httpReq := &internalhttp.Request{
		Method:  "POST",
		Path:    path,
		Body:    &createReq,
		Headers: headers,
	}

	var resp repmod.CreateSnapshotResponse
	if err := c.baseClient.Do(ctx, httpReq, &resp); err != nil {
		return nil, fmt.Errorf("failed to snapshot server %s: %w", serverID, err)
	}
	if resp.Repository == nil {
		return nil, fmt.Errorf("snapshot response missing repository data")
	}
	if resp.Tag != nil {
		resp.Repository.Tags = append(resp.Repository.Tags, resp.Tag)
	}

	return &RepositoryResource{
		Repository: resp.Repository,
		tagOps: &TagsClient{
			baseClient:   c.baseClient,
			repositoryID: resp.Repository.ID,
			basePath:     c.basePath,
		},
	}, nil
}

// Upload uploads an image into VRM and returns the affected repository resource.
// POST /api/v1/project/{project-id}/upload
// Supports three modes: creating a new repository, uploading to an existing repository, or uploading to an existing tag.
func (c *Client) Upload(ctx context.Context, req repmod.UploadRequester) (*RepositoryResource, error) {
	return c.UploadWithNamespace(ctx, req, "")
}

// UploadWithNamespace uploads an image with optional namespace scoping via X-Namespace header.
// POST /api/v1/project/{project-id}/upload
func (c *Client) UploadWithNamespace(ctx context.Context, req repmod.UploadRequester, namespace string) (*RepositoryResource, error) {
	if req == nil {
		return nil, fmt.Errorf("upload request cannot be nil")
	}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid upload request: %w", err)
	}

	uploadReq := req.ToUploadImageRequest()

	path := c.basePath + "/upload"

	headers := make(map[string]string)
	if namespace != "" {
		headers["X-Namespace"] = namespace
	}

	httpReq := &internalhttp.Request{
		Method:  "POST",
		Path:    path,
		Body:    &uploadReq,
		Headers: headers,
	}

	var resp repmod.UploadImageResponse
	if err := c.baseClient.Do(ctx, httpReq, &resp); err != nil {
		return nil, fmt.Errorf("failed to upload image: %w", err)
	}
	if resp.Repository == nil {
		return nil, fmt.Errorf("upload response missing repository data")
	}
	if resp.Tag != nil {
		resp.Repository.Tags = append(resp.Repository.Tags, resp.Tag)
	}

	return &RepositoryResource{
		Repository: resp.Repository,
		tagOps: &TagsClient{
			baseClient:   c.baseClient,
			repositoryID: resp.Repository.ID,
			basePath:     c.basePath,
		},
	}, nil
}

// RepositoryResource wraps a Repository with sub-resource operations.
type RepositoryResource struct {
	*repmod.Repository
	tagOps TagOperations
}

// Tags returns the tag operations for this repository.
func (rr *RepositoryResource) Tags() TagOperations {
	return rr.tagOps
}

// TagOperations defines operations on repository tags (sub-resource).
type TagOperations interface {
	List(ctx context.Context, opts *tagmod.ListTagsOptions) ([]*tagmod.Tag, error)
	Create(ctx context.Context, req *tagmod.CreateTagRequest) (*tagmod.Tag, error)
	CreateWithNamespace(ctx context.Context, req *tagmod.CreateTagRequest, namespace string) (*tagmod.Tag, error)
}

// TagsClient implements TagOperations for a specific repository.
type TagsClient struct {
	baseClient   *internalhttp.Client
	repositoryID string
	basePath     string
}

// List retrieves all tags in the repository.
// GET /api/v1/project/{project-id}/repository/{repository-id}/tags
func (tc *TagsClient) List(ctx context.Context, opts *tagmod.ListTagsOptions) ([]*tagmod.Tag, error) {
	if opts == nil {
		opts = &tagmod.ListTagsOptions{}
	}
	if err := opts.Validate(); err != nil {
		return nil, fmt.Errorf("invalid list options: %w", err)
	}

	path := tc.basePath + "/repository/" + url.PathEscape(tc.repositoryID) + "/tags"

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
	if err := tc.baseClient.Do(ctx, req, &resp); err != nil {
		return nil, fmt.Errorf("failed to list tags by repository %s: %w", tc.repositoryID, err)
	}

	return resp.Tags, nil
}

// Create creates a new tag in the repository.
// POST /api/v1/project/{project-id}/repository/{repository-id}/tag
// Uses default namespace handling (no X-Namespace header).
func (tc *TagsClient) Create(ctx context.Context, req *tagmod.CreateTagRequest) (*tagmod.Tag, error) {
	return tc.CreateWithNamespace(ctx, req, "")
}

// CreateWithNamespace creates a new tag in the repository with optional namespace header.
// POST /api/v1/project/{project-id}/repository/{repository-id}/tag
// If namespace is provided, sets the X-Namespace header for multi-tenant operations.
func (tc *TagsClient) CreateWithNamespace(ctx context.Context, req *tagmod.CreateTagRequest, namespace string) (*tagmod.Tag, error) {
	if req == nil {
		return nil, fmt.Errorf("create request cannot be nil")
	}
	if err := req.Validate(); err != nil {
		return nil, fmt.Errorf("invalid create request: %w", err)
	}

	path := tc.basePath + "/repository/" + url.PathEscape(tc.repositoryID) + "/tag"

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
	if err := tc.baseClient.Do(ctx, httpReq, &tag); err != nil {
		return nil, fmt.Errorf("failed to create tag: %w", err)
	}

	return &tag, nil
}
