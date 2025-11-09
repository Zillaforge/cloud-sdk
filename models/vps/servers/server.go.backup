package servers

// Server represents a compute instance.
type Server struct {
	ID          string            `json:"id"`
	Name        string            `json:"name"`
	Description string            `json:"description,omitempty"`
	Status      string            `json:"status"` // ACTIVE, BUILD, SHUTOFF, ERROR, etc.
	FlavorID    string            `json:"flavor_id"`
	ImageID     string            `json:"image_id"`
	ProjectID   string            `json:"project_id"`
	UserID      string            `json:"user_id,omitempty"`
	CreatedAt   string            `json:"created_at"`
	UpdatedAt   string            `json:"updated_at"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// ListServersOptions contains filter/pagination options for List.
type ListServersOptions struct {
	Name     string
	UserID   string
	Status   string
	FlavorID string
	ImageID  string
	Limit    int
	Offset   int
}

// ServerListResponse is the response from List.
type ServerListResponse struct {
	Items []*Server `json:"items"`
	Total int       `json:"total,omitempty"`
}

// ServerCreateRequest is the body for Create.
type ServerCreateRequest struct {
	Name        string             `json:"name"`
	Description string             `json:"description,omitempty"`
	FlavorID    string             `json:"flavor_id"`
	ImageID     string             `json:"image_id"`
	NICs        []ServerNICRequest `json:"nics"`
	SGIDs       []string           `json:"sg_ids"`
	KeypairID   string             `json:"keypair_id,omitempty"`
	Password    string             `json:"password,omitempty"`    // Base64
	BootScript  string             `json:"boot_script,omitempty"` // Base64
}

// ServerNICRequest specifies a NIC for server creation.
type ServerNICRequest struct {
	NetworkID string `json:"network_id"`
	FixedIP   string `json:"fixed_ip,omitempty"`
}

// ServerUpdateRequest is the body for Update.
type ServerUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}
