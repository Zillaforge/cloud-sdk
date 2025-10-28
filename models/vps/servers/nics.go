package servers

// ServerNIC represents a vNIC attached to a server.
type ServerNIC struct {
	ID        string   `json:"id"`
	NetworkID string   `json:"network_id"`
	FixedIPs  []string `json:"fixed_ips"`
	MACAddr   string   `json:"mac_address"`
	SGIDs     []string `json:"sg_ids"`
}

// ServerNICCreateRequest is the body for AddServerNIC.
type ServerNICCreateRequest struct {
	NetworkID string   `json:"network_id"`
	SGIDs     []string `json:"sg_ids"`
	FixedIP   string   `json:"fixed_ip,omitempty"`
}

// ServerNICUpdateRequest is the body for UpdateServerNIC.
type ServerNICUpdateRequest struct {
	SGIDs []string `json:"sg_ids"`
}

// FloatingIPAssociateRequest is the body for AssociateFloatingIPToNIC.
type FloatingIPAssociateRequest struct {
	FloatingIPID string `json:"fip_id,omitempty"` // existing FIP; omit to create new
}
