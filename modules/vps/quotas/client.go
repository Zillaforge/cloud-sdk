package quotas

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/quotas"
)

// Client provides operations for managing quotas.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new quotas client.
func NewClient(baseClient *internalhttp.Client, projectID string) *Client {
	basePath := "/api/v1/project/" + projectID
	return &Client{
		baseClient: baseClient,
		projectID:  projectID,
		basePath:   basePath,
	}
}

// Get retrieves the project's quota information.
// GET /api/v1/project/{project-id}/quotas
func (c *Client) Get(ctx context.Context) (*quotas.Quota, error) {
	path := c.basePath + "/quotas"

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var quota quotas.Quota
	if err := c.baseClient.Do(ctx, req, &quota); err != nil {
		return nil, fmt.Errorf("failed to get quotas: %w", err)
	}

	return &quota, nil
}
