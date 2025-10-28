package routers

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/routers"
)

// RouterNetworkOperations defines operations on router networks (sub-resource).
type RouterNetworkOperations interface {
	List(ctx context.Context) ([]*routers.RouterNetwork, error)
	Associate(ctx context.Context, networkID string) error
	Disassociate(ctx context.Context, networkID string) error
}

// routerNetworksClient implements RouterNetworkOperations.
type routerNetworksClient struct {
	baseClient *internalhttp.Client
	projectID  string
	routerID   string
	basePath   string
}

// newRouterNetworksClient creates a new router networks client.
func newRouterNetworksClient(baseClient *internalhttp.Client, projectID, routerID string) *routerNetworksClient {
	basePath := fmt.Sprintf("/api/v1/project/%s/routers/%s", projectID, routerID)
	return &routerNetworksClient{
		baseClient: baseClient,
		projectID:  projectID,
		routerID:   routerID,
		basePath:   basePath,
	}
}

// List retrieves all networks associated with the router.
// GET /api/v1/project/{project-id}/routers/{router-id}/networks
func (c *routerNetworksClient) List(ctx context.Context) ([]*routers.RouterNetwork, error) {
	path := c.basePath + "/networks"

	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var networks []*routers.RouterNetwork
	if err := c.baseClient.Do(ctx, req, &networks); err != nil {
		return nil, fmt.Errorf("failed to list router networks: %w", err)
	}

	return networks, nil
}

// Associate associates a network with the router.
// POST /api/v1/project/{project-id}/routers/{router-id}/networks/{network-id}
func (c *routerNetworksClient) Associate(ctx context.Context, networkID string) error {
	path := fmt.Sprintf("%s/networks/%s", c.basePath, networkID)

	req := &internalhttp.Request{
		Method: "POST",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to associate network %s with router: %w", networkID, err)
	}

	return nil
}

// Disassociate disassociates a network from the router.
// DELETE /api/v1/project/{project-id}/routers/{router-id}/networks/{network-id}
func (c *routerNetworksClient) Disassociate(ctx context.Context, networkID string) error {
	path := fmt.Sprintf("%s/networks/%s", c.basePath, networkID)

	req := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return fmt.Errorf("failed to disassociate network %s from router: %w", networkID, err)
	}

	return nil
}
