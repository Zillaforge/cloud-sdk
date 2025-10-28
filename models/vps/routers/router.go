package routers

import (
	"context"
	"time"
)

// Router represents a VPS router resource.
// Swagger reference: pb.RouterInfo
type Router struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Description  string    `json:"description,omitempty"`
	State        bool      `json:"state"`  // true = enabled, false = disabled
	Status       string    `json:"status"` // operational status
	ProjectID    string    `json:"project_id"`
	UserID       string    `json:"user_id,omitempty"`
	Namespace    string    `json:"namespace,omitempty"`
	ExtNetworkID string    `json:"extnetwork_id,omitempty"`
	IsDefault    bool      `json:"is_default"`
	Shared       bool      `json:"shared"`
	Bonding      bool      `json:"bonding"`
	GatewayAddrs []string  `json:"gw_addrs,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at,omitempty"`
}

// RouterCreateRequest represents the request body for creating a router.
// Swagger reference: RouterCreateInput
type RouterCreateRequest struct {
	Name         string `json:"name"`
	Description  string `json:"description,omitempty"`
	ExtNetworkID string `json:"extnetwork_id,omitempty"`
}

// RouterUpdateRequest represents the request body for updating a router.
// Swagger reference: RouterUpdateInput
type RouterUpdateRequest struct {
	Name        *string `json:"name,omitempty"`
	Description *string `json:"description,omitempty"`
}

// RouterSetStateRequest represents the request body for setting router state.
// Swagger reference: RouterSetStateInput
type RouterSetStateRequest struct {
	State bool `json:"state"` // true = enabled, false = disabled
}

// ListRoutersOptions contains query parameters for listing routers.
type ListRoutersOptions struct {
	Name   string // Filter by router name
	UserID string // Filter by user ID
	Detail bool   // Include detailed information
}

// RouterListResponse represents the response from listing routers.
// Swagger reference: pb.RouterListOutput
type RouterListResponse struct {
	Routers []Router `json:"routers"`
	Total   int      `json:"total"`
}

// IDName is a helper struct for nested ID/Name pairs in responses.
type IDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// RouterResource wraps a Router with a method to access network operations.
// The Networks() accessor must be injected by the client that creates this resource.
type RouterResource struct {
	*Router
	NetworksOps RouterNetworkOperations
}

// RouterNetworkOperations is an interface for router network sub-resource operations.
// This interface is implemented in the modules/vps/routers package.
type RouterNetworkOperations interface {
	List(ctx context.Context) ([]*RouterNetwork, error)
	Associate(ctx context.Context, networkID string) error
	Disassociate(ctx context.Context, networkID string) error
}

// Networks returns the network operations for this router.
func (r *RouterResource) Networks() RouterNetworkOperations {
	return r.NetworksOps
}
