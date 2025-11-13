package volumetypes

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	volumetypesmodel "github.com/Zillaforge/cloud-sdk/models/vps/volumetypes"
)

// Client provides operations for managing volume types.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new volume types client.
// Follows pattern from modules/vps/flavors/client.go.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// List retrieves a list of available volume types in the project.
// GET /api/v1/project/{project-id}/volume_types
func (c *Client) List(ctx context.Context) ([]string, error) {
	path := c.basePath + "/volume_types"

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response volumetypesmodel.VolumeTypeListResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list volume types: %w", err)
	}

	return response.VolumeTypes, nil
}
