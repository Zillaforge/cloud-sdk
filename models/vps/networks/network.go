package networks

// Network represents a virtual network in the VPS infrastructure.
type Network struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CIDR        string `json:"cidr"`
	ProjectID   string `json:"project_id"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// NetworkCreateRequest represents the request to create a new network.
type NetworkCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	CIDR        string `json:"cidr"`
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
type ListNetworksOptions struct {
	Name string
	// Additional pagination fields can be added here as needed
}
