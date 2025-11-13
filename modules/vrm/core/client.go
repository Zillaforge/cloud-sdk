// Package vrm provides project-scoped VRM operations.
package vrm

import (
	"net/http"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/internal/types"
	"github.com/Zillaforge/cloud-sdk/modules/vrm/repositories"
	"github.com/Zillaforge/cloud-sdk/modules/vrm/tags"
)

// Client provides access to VRM operations for a specific project.
// All operations are scoped to the project ID provided at creation.
// It provides sub-clients for repositories and tags management.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new project-scoped VRM client.
// This is typically called via cloudsdk.Client.Project(projectID).VRM().
// The client is configured with the provided base URL, authentication token,
// and project scope. It uses the provided HTTP client and logger for operations.
func NewClient(baseURL, token, projectID string, httpClient *http.Client, logger types.Logger) *Client {
	basePath := "/api/v1/project/" + projectID

	return &Client{
		baseClient: internalhttp.NewClient(baseURL, token, httpClient, logger),
		projectID:  projectID,
		basePath:   basePath,
	}
}

// ProjectID returns the project ID this client is bound to.
// This is the project scope for all VRM operations performed by this client.
func (c *Client) ProjectID() string {
	return c.projectID
}

// Repositories returns the repository operations client.
// Use this client to perform CRUD operations on VRM repositories within the project scope.
func (c *Client) Repositories() *repositories.Client {
	return repositories.NewClient(c.baseClient, c.projectID, c.basePath)
}

// Tags returns the tag operations client.
// Use this client to perform CRUD operations on VRM tags within the project scope.
func (c *Client) Tags() *tags.Client {
	return tags.NewClient(c.baseClient, c.projectID, c.basePath)
}
