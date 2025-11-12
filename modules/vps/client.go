// Package vps provides project-scoped VPS operations.
package vps

import (
	"net/http"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/internal/types"
	"github.com/Zillaforge/cloud-sdk/modules/vps/flavors"
	"github.com/Zillaforge/cloud-sdk/modules/vps/floatingips"
	"github.com/Zillaforge/cloud-sdk/modules/vps/keypairs"
	"github.com/Zillaforge/cloud-sdk/modules/vps/networks"
	"github.com/Zillaforge/cloud-sdk/modules/vps/quotas"
	"github.com/Zillaforge/cloud-sdk/modules/vps/securitygroups"
	"github.com/Zillaforge/cloud-sdk/modules/vps/servers"
)

// Client provides access to VPS operations for a specific project.
// All operations are scoped to the project ID provided at creation.
type Client struct {
	baseClient *internalhttp.Client
	projectID  string
	basePath   string
}

// NewClient creates a new project-scoped VPS client.
// This is typically called via cloudsdk.Client.Project(projectID).VPS().
func NewClient(baseURL, token, projectID string, httpClient *http.Client, logger types.Logger) *Client {
	basePath := "/api/v1/project/" + projectID

	return &Client{
		baseClient: internalhttp.NewClient(baseURL, token, httpClient, logger),
		projectID:  projectID,
		basePath:   basePath,
	}
}

// ProjectID returns the project ID this client is bound to.
func (c *Client) ProjectID() string {
	return c.projectID
}

// Networks returns the network operations client.
func (c *Client) Networks() *networks.Client {
	return networks.NewClient(c.baseClient, c.projectID)
}

// FloatingIPs returns the floating IP operations client.
func (c *Client) FloatingIPs() *floatingips.Client {
	return floatingips.NewClient(c.baseClient, c.projectID)
}

// Flavors returns the flavors operations client.
func (c *Client) Flavors() *flavors.Client {
	return flavors.NewClient(c.baseClient, c.projectID)
}

// Keypairs returns the keypairs operations client.
func (c *Client) Keypairs() *keypairs.Client {
	return keypairs.NewClient(c.baseClient, c.projectID)
}

// Quotas returns the quotas operations client.
func (c *Client) Quotas() *quotas.Client {
	return quotas.NewClient(c.baseClient, c.projectID)
}

// SecurityGroups returns the security groups operations client.
func (c *Client) SecurityGroups() *securitygroups.Client {
	return securitygroups.NewClient(c.baseClient, c.projectID)
}

// Servers returns the servers operations client.
func (c *Client) Servers() *servers.Client {
	return servers.NewClient(c.baseClient, c.projectID)
}
