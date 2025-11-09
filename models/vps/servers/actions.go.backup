package servers

// ServerActionRequest is the body for Action.
type ServerActionRequest struct {
	Action     string `json:"action"`                // start, stop, reboot, resize, extend_root, get_pwd, approve, reject
	RebootType string `json:"reboot_type,omitempty"` // hard or soft (for reboot)
	FlavorID   string `json:"flavor_id,omitempty"`   // for resize
	RootSize   int    `json:"root_size,omitempty"`   // for extend_root
	PrivateKey string `json:"private_key,omitempty"` // Base64 for get_pwd
}
