package floatingips

// FloatingIP represents a floating IP address.
type FloatingIP struct {
	ID          string `json:"id"`
	Address     string `json:"address"`
	Status      string `json:"status"` // ACTIVE, PENDING, DOWN, REJECTED
	ProjectID   string `json:"project_id"`
	PortID      string `json:"port_id,omitempty"`
	Description string `json:"description,omitempty"`
	CreatedAt   string `json:"created_at"`
}

// FloatingIPCreateRequest represents the request to create a new floating IP.
type FloatingIPCreateRequest struct {
	Description string `json:"description,omitempty"`
	// ExtNetworkID omitted if default
}

// FloatingIPUpdateRequest represents the request to update an existing floating IP.
type FloatingIPUpdateRequest struct {
	Description string `json:"description,omitempty"`
}

// FloatingIPListResponse represents the response from listing floating IPs.
type FloatingIPListResponse struct {
	Items []*FloatingIP `json:"items"`
}

// ListFloatingIPsOptions provides filtering options for floating IP listing.
type ListFloatingIPsOptions struct {
	Status string // filter by status
}
