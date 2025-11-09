package servers

// FloatingIPInfo contains floating IP information.
type FloatingIPInfo struct {
	Address      string  `json:"address"`
	ApprovedAt   string  `json:"approvedAt,omitempty"`
	CreatedAt    string  `json:"createdAt,omitempty"`
	Description  string  `json:"description,omitempty"`
	DeviceID     string  `json:"device_id,omitempty"`
	DeviceName   string  `json:"device_name,omitempty"`
	DeviceType   string  `json:"device_type,omitempty"`
	ExtnetID     string  `json:"extnet_id,omitempty"`
	ID           string  `json:"id"`
	Name         string  `json:"name,omitempty"`
	Namespace    string  `json:"namespace,omitempty"`
	PortID       string  `json:"port_id,omitempty"`
	Project      *IDName `json:"project,omitempty"`
	ProjectID    string  `json:"project_id,omitempty"`
	Reserved     bool    `json:"reserved,omitempty"`
	Status       string  `json:"status,omitempty"`
	StatusReason string  `json:"status_reason,omitempty"`
	UpdatedAt    string  `json:"updatedAt,omitempty"`
	User         *IDName `json:"user,omitempty"`
	UserID       string  `json:"user_id,omitempty"`
	UUID         string  `json:"uuid,omitempty"`
}

// ServerNIC represents a vNIC attached to a server (ServerNICInfo from pb).
type ServerNIC struct {
	ID             string          `json:"id"`
	MAC            string          `json:"mac"`
	NetworkID      string          `json:"network_id"`
	Network        *IDName         `json:"network,omitempty"`
	Addresses      []string        `json:"addresses,omitempty"`
	FloatingIP     *FloatingIPInfo `json:"floating_ip,omitempty"`
	SecurityGroups []*IDName       `json:"security_groups,omitempty"`
	SGIDs          []string        `json:"sg_ids,omitempty"`
	IsProviderNet  bool            `json:"is_provider_net,omitempty"`
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

// ServerNICsListResponse is the response from listing NICs.
type ServerNICsListResponse struct {
	NICs []*ServerNIC `json:"nics"`
}

// ServerNICAssociateFloatingIPRequest is the body for AssociateFloatingIPToNIC.
type ServerNICAssociateFloatingIPRequest struct {
	FIPID string `json:"fip_id"`
}
