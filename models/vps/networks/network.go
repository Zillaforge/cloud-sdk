package networks

// Network represents a virtual network in the VPS infrastructure.
// Swagger reference: pb.NetworkInfo
type Network struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CIDR        string `json:"cidr"`
	Bonding     bool   `json:"bonding,omitempty"`
	Gateway     string `json:"gateway,omitempty"`
	// Deprecated: Marked as deprecated in the upstream pb.NetworkInfo contract.
	GWState      bool        `json:"gw_state,omitempty"`
	IsDefault    bool        `json:"is_default,omitempty"`
	Nameservers  []string    `json:"nameservers,omitempty"`
	Namespace    string      `json:"namespace,omitempty"`
	Project      *IDName     `json:"project,omitempty"`
	ProjectID    string      `json:"project_id,omitempty"`
	Router       *RouterInfo `json:"router,omitempty"`
	RouterID     string      `json:"router_id,omitempty"`
	Shared       bool        `json:"shared,omitempty"`
	Status       string      `json:"status,omitempty"`
	StatusReason string      `json:"status_reason,omitempty"`
	SubnetID     string      `json:"subnet_id,omitempty"`
	User         *IDName     `json:"user,omitempty"`
	UserID       string      `json:"user_id,omitempty"`
	CreatedAt    string      `json:"createdAt"`
	UpdatedAt    string      `json:"updatedAt,omitempty"`
}

// IDName is a reusable nested identifier/name pair.
type IDName struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// ExtNetworkInfo represents the external network summary embedded in router info.
type ExtNetworkInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CIDR        string `json:"cidr,omitempty"`
	Namespace   string `json:"namespace,omitempty"`
	SegmentID   string `json:"segment_id,omitempty"`
	Type        string `json:"type,omitempty"`
	IsDefault   bool   `json:"is_default,omitempty"`
}

// RouterInfo is a lightweight router representation associated with a network.
// Swagger reference: pb.RouterInfo
type RouterInfo struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Bonding      bool            `json:"bonding,omitempty"`
	IsDefault    bool            `json:"is_default,omitempty"`
	Shared       bool            `json:"shared,omitempty"`
	State        bool            `json:"state,omitempty"`
	Status       string          `json:"status,omitempty"`
	StatusReason string          `json:"status_reason,omitempty"`
	Namespace    string          `json:"namespace,omitempty"`
	Project      *IDName         `json:"project,omitempty"`
	ProjectID    string          `json:"project_id,omitempty"`
	User         *IDName         `json:"user,omitempty"`
	UserID       string          `json:"user_id,omitempty"`
	ExtNetwork   *ExtNetworkInfo `json:"extnetwork,omitempty"`
	ExtNetworkID string          `json:"extnetwork_id,omitempty"`
	GWAddrs      []string        `json:"gw_addrs,omitempty"`
	CreatedAt    string          `json:"createdAt,omitempty"`
	UpdatedAt    string          `json:"updatedAt,omitempty"`
}

// NetworkCreateRequest represents the request to create a new network.
// Swagger reference: NetCreateInput
type NetworkCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CIDR        string `json:"cidr"`
	Gateway     string `json:"gateway,omitempty"`
	RouterID    string `json:"router_id,omitempty"`
}

// NetworkUpdateRequest represents the request to update an existing network.
type NetworkUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// NetworkListResponse represents the response from listing networks.
type NetworkListResponse struct {
	Networks []*Network `json:"networks"`
}

// ListNetworksOptions provides filtering options for network listing.
// Swagger reference: GET /api/v1/project/{project-id}/networks query parameters
type ListNetworksOptions struct {
	Name     string // Filter by network name
	UserID   string // Filter by user_id
	Status   string // Filter by network status
	RouterID string // Filter by router_id
	Detail   *bool  // Get detailed information (optional boolean)
}
