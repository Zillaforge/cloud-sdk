package servers

// ServerAction represents a server action type.
type ServerAction string

// Server action constants
const (
	ServerActionStop       ServerAction = "stop"
	ServerActionStart      ServerAction = "start"
	ServerActionReboot     ServerAction = "reboot"
	ServerActionResize     ServerAction = "resize"
	ServerActionApprove    ServerAction = "approve"
	ServerActionReject     ServerAction = "reject"
	ServerActionExtendRoot ServerAction = "extend_root"
	ServerActionGetPwd     ServerAction = "get_pwd"
)

// RebootType represents the type of reboot.
type RebootType string

// Reboot type constants
const (
	RebootTypeHard RebootType = "hard"
	RebootTypeSoft RebootType = "soft"
)

// ServerActionRequest is the body for Action.
type ServerActionRequest struct {
	Action     ServerAction `json:"action"`
	RebootType RebootType   `json:"reboot_type,omitempty"` // for reboot
	FlavorID   string       `json:"flavor_id,omitempty"`   // for resize
	RootSize   int          `json:"root_size,omitempty"`   // for extend_root
	PrivateKey string       `json:"private_key,omitempty"` // Base64 for get_pwd
}

// ServerActionResponse is the response from Action.
type ServerActionResponse struct {
	Password string `json:"password,omitempty"` // returned for get_pwd action
}
