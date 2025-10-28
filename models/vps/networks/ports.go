package networks

// NetworkPort represents a port attached to a network.
type NetworkPort struct {
	ID        string   `json:"id"`
	NetworkID string   `json:"network_id"`
	FixedIPs  []string `json:"fixed_ips"`
	MACAddr   string   `json:"mac_address"`
	ServerID  string   `json:"server_id,omitempty"`
}
