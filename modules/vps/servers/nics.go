package servers

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

// NICsClient handles NIC-related operations for a server.
type NICsClient struct {
	baseClient *internalhttp.Client
	projectID  string
	serverID   string
}

// List lists all network interfaces on the server.
// GET /api/v1/project/{project-id}/servers/{svr-id}/nics
func (c *NICsClient) List(ctx context.Context) ([]*servers.ServerNIC, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/nics", c.projectID, c.serverID)

	// Make request
	req := &internalhttp.Request{
		Method: "GET",
		Path:   path,
	}

	var nics []*servers.ServerNIC
	if err := c.baseClient.Do(ctx, req, &nics); err != nil {
		return nil, err
	}

	return nics, nil
}

// Add attaches a new vNIC to the server.
// POST /api/v1/project/{project-id}/servers/{svr-id}/nics
func (c *NICsClient) Add(ctx context.Context, req *servers.ServerNICCreateRequest) (*servers.ServerNIC, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/nics", c.projectID, c.serverID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	var nic servers.ServerNIC
	if err := c.baseClient.Do(ctx, httpReq, &nic); err != nil {
		return nil, err
	}

	return &nic, nil
}

// Update updates security groups on an existing vNIC.
// PUT /api/v1/project/{project-id}/servers/{svr-id}/nics/{nic-id}
func (c *NICsClient) Update(ctx context.Context, nicID string, req *servers.ServerNICUpdateRequest) (*servers.ServerNIC, error) {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/nics/%s", c.projectID, c.serverID, nicID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "PUT",
		Path:   path,
		Body:   req,
	}

	var nic servers.ServerNIC
	if err := c.baseClient.Do(ctx, httpReq, &nic); err != nil {
		return nil, err
	}

	return &nic, nil
}

// Delete detaches and removes a vNIC from the server.
// DELETE /api/v1/project/{project-id}/servers/{svr-id}/nics/{nic-id}
func (c *NICsClient) Delete(ctx context.Context, nicID string) error {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/nics/%s", c.projectID, c.serverID, nicID)

	// Make request
	req := &internalhttp.Request{
		Method: "DELETE",
		Path:   path,
	}

	if err := c.baseClient.Do(ctx, req, nil); err != nil {
		return err
	}

	return nil
}

// AssociateFloatingIP associates a floating IP to a specific vNIC.
// POST /api/v1/project/{project-id}/servers/{svr-id}/nics/{nic-id}/floatingip
func (c *NICsClient) AssociateFloatingIP(ctx context.Context, nicID string, req *servers.FloatingIPAssociateRequest) error {
	path := fmt.Sprintf("/api/v1/project/%s/servers/%s/nics/%s/floatingip", c.projectID, c.serverID, nicID)

	// Make request
	httpReq := &internalhttp.Request{
		Method: "POST",
		Path:   path,
		Body:   req,
	}

	if err := c.baseClient.Do(ctx, httpReq, nil); err != nil {
		return err
	}

	return nil
}
