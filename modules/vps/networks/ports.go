package networks

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// PortsClient handles port operations for a specific network.
type PortsClient struct {
	baseClient *internalhttp.Client
	projectID  string
	networkID  string
}

// List lists all ports on the network.
// GET /api/v1/project/{project-id}/networks/{net-id}/ports
func (c *PortsClient) List(ctx context.Context) ([]*networks.NetworkPort, error) {
	path := fmt.Sprintf("/api/v1/project/%s/networks/%s/ports", c.projectID, c.networkID)

	// Make request
	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var ports []*networks.NetworkPort
	if err := c.baseClient.Do(ctx, req, &ports); err != nil {
		return nil, err
	}

	return ports, nil
}
