package servers

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

// VolumesClient handles volume-related operations for a server.
type VolumesClient struct {
	baseClient *internalhttp.Client
	projectID  string
	serverID   string
}

// List lists all volume attachments for the server.
// GET /api/v1/project/{project-id}/servers/{svr-id}/volumes
func (c *VolumesClient) List(ctx context.Context) ([]*servers.ServerVolume, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/volumes", c.projectID, c.serverID)

	// Make request
	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var response servers.ServerVolumesResponse
	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to list volumes for server %s: %w", c.serverID, err)
	}

	return response.Disks, nil
}

// Attach attaches a volume to the server.
// POST /api/v1/project/{project-id}/servers/{svr-id}/volumes/{vol-id}
func (c *VolumesClient) Attach(ctx context.Context, volumeID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/volumes/%s", c.projectID, c.serverID, volumeID)

	// Make request
	req := &internalhttp.Request{
		Method: "POST",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to attach volume %s to server %s: %w", volumeID, c.serverID, err)
	}

	return nil
}

// Detach detaches a volume from the server.
// DELETE /api/v1/project/{project-id}/servers/{svr-id}/volumes/{vol-id}
func (c *VolumesClient) Detach(ctx context.Context, volumeID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/volumes/%s", c.projectID, c.serverID, volumeID)

	// Make request
	req := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to detach volume %s from server %s: %w", volumeID, c.serverID, err)
	}

	return nil
}
