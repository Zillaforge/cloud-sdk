package volumes

import (
	"context"
	"fmt"
	"net/url"
	"strconv"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	volumesmodel "github.com/Zillaforge/cloud-sdk/models/vps/volumes"
)

// Client provides operations for managing volumes.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new volumes client.
// Follows pattern from modules/vps/flavors/client.go.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// Create creates a new volume.
// POST /api/v1/project/{project-id}/volumes
func (c *Client) Create(ctx context.Context, request *volumesmodel.CreateVolumeRequest) (*volumesmodel.Volume, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := c.basePath + "/volumes"

	req := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   request,
	}

	var response volumesmodel.Volume
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to create volume: %w", err)
	}

	return &response, nil
}

// Update updates volume metadata (name, description).
// PUT /api/v1/project/{project-id}/volumes/{volume-id}
func (c *Client) Update(ctx context.Context, volumeID string, request *volumesmodel.UpdateVolumeRequest) (*volumesmodel.Volume, error) {
	if err := request.Validate(); err != nil {
		return nil, fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("%s/volumes/%s", c.basePath, volumeID)

	req := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   request,
	}

	var response volumesmodel.Volume
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to update volume %s: %w", volumeID, err)
	}

	return &response, nil
}

// Delete deletes a volume.
// DELETE /api/v1/project/{project-id}/volumes/{volume-id}
func (c *Client) Delete(ctx context.Context, volumeID string) error {
	path := fmt.Sprintf("%s/volumes/%s", c.basePath, volumeID)

	req := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to delete volume %s: %w", volumeID, err)
	}

	return nil
}

// List retrieves a list of volumes with optional filtering.
// GET /api/v1/project/{project-id}/volumes
func (c *Client) List(ctx context.Context, opts *volumesmodel.ListVolumesOptions) ([]*volumesmodel.Volume, error) {
	path := c.basePath + "/volumes"

	// Build query parameters
	query := url.Values{}
	if opts != nil {
		if opts.Name != "" {
			query.Set("name", opts.Name)
		}
		if opts.UserID != "" {
			query.Set("user_id", opts.UserID)
		}
		if opts.Status != "" {
			query.Set("status", opts.Status)
		}
		if opts.Type != "" {
			query.Set("type", opts.Type)
		}
		if opts.Detail {
			query.Set("detail", strconv.FormatBool(opts.Detail))
		}
	}

	if len(query) > 0 {
		path += "?" + query.Encode()
	}

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response volumesmodel.VolumeListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	return response.Volumes, nil
}

// Get retrieves details of a specific volume by ID.
// GET /api/v1/project/{project-id}/volumes/{volume-id}
func (c *Client) Get(ctx context.Context, volumeID string) (*volumesmodel.Volume, error) {
	path := fmt.Sprintf("%s/volumes/%s", c.basePath, volumeID)

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response volumesmodel.Volume
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to get volume %s: %w", volumeID, err)
	}

	return &response, nil
}

// Action performs an action on a volume (attach, detach, extend, revert).
// POST /api/v1/project/{project-id}/volumes/{volume-id}/action
// Returns 202 Accepted for async operations.
func (c *Client) Action(ctx context.Context, volumeID string, request *volumesmodel.VolumeActionRequest) error {
	if err := request.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	path := fmt.Sprintf("%s/volumes/%s/action", c.basePath, volumeID)

	req := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   request,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to perform action %s on volume %s: %w", request.Action, volumeID, err)
	}

	return nil
}
