// Package contracts defines the VPS network operations interface.
// Swagger reference: /api/v1/project/{project-id}/networks*
package contracts

import "context"

// NetworkOperations defines all network-related API operations.
type NetworkOperations interface {
	// List retrieves all networks for the project with optional filters.
	// GET /api/v1/project/{project-id}/networks
	// Returns: NetworkListResponse
	// Errors: 401, 403, 500
	List(ctx context.Context, opts *ListNetworksOptions) (*NetworkListResponse, error)

	// Create creates a new network.
	// POST /api/v1/project/{project-id}/networks
	// Body: NetworkCreateRequest
	// Returns: 201 + Network
	// Errors: 400, 401, 403, 409 (quota/CIDR conflict), 500
	Create(ctx context.Context, req *NetworkCreateRequest) (*Network, error)

	// Get retrieves a specific network with sub-resource operations.
	// GET /api/v1/project/{project-id}/networks/{net-id}
	// Returns: 200 + NetworkResource
	// Errors: 401, 403, 404, 500
	Get(ctx context.Context, networkID string) (*NetworkResource, error)

	// Update updates network name/description.
	// PUT /api/v1/project/{project-id}/networks/{net-id}
	// Body: NetworkUpdateRequest
	// Returns: 200 + Network
	// Errors: 400, 401, 403, 404, 500
	Update(ctx context.Context, networkID string, req *NetworkUpdateRequest) (*Network, error)

	// Delete deletes a network.
	// DELETE /api/v1/project/{project-id}/networks/{net-id}
	// Returns: 204
	// Errors: 401, 403, 404, 409 (ports attached), 500
	Delete(ctx context.Context, networkID string) error
}

type ListNetworksOptions struct {
	Name string
	// Pagination if defined
}

type NetworkListResponse struct {
	Items []*Network `json:"items"`
}

type Network struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CIDR        string `json:"cidr"`
	ProjectID   string `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// NetworkResource wraps a Network with sub-resource operations.
type NetworkResource struct {
	*Network
	portOps NetworkPortOperations
}

// Ports returns the port operations for this network.
func (n *NetworkResource) Ports() NetworkPortOperations {
	return n.portOps
}

// NetworkPortOperations defines operations on network ports (sub-resource).
type NetworkPortOperations interface {
	// List lists all ports on the network.
	// GET /api/v1/project/{project-id}/networks/{net-id}/ports
	// Returns: 200 + array of NetworkPort
	// Errors: 401, 403, 404, 500
	List(ctx context.Context) ([]*NetworkPort, error)
}

type NetworkCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CIDR        string `json:"cidr"`
}

type NetworkUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

type NetworkPort struct {
	ID        string   `json:"id"`
	NetworkID string   `json:"network_id"`
	FixedIPs  []string `json:"fixed_ips"`
	MACAddr   string   `json:"mac_address"`
	ServerID  string   `json:"server_id,omitempty"`
}
