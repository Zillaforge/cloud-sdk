package routers

// RouterNetwork represents a network associated with a router.
// Swagger reference: RouterNetworkInfo
type RouterNetwork struct {
	NetworkID   string `json:"network_id"`
	NetworkName string `json:"network_name,omitempty"`
	SubnetID    string `json:"subnet_id,omitempty"`
	PortID      string `json:"port_id,omitempty"`
}

// RouterNetworkAssociateRequest represents the request body for associating a network with a router.
type RouterNetworkAssociateRequest struct {
	NetworkID string `json:"network_id"`
}
