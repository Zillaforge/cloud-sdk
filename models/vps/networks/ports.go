package networks

// NetworkPort represents a port attached to a network.
// Swagger reference: NetPort
type NetworkPort struct {
	ID        string         `json:"id"`
	Addresses []string       `json:"addresses,omitempty"`
	Server    *ServerSummary `json:"server,omitempty"`
}

// ServerSummary captures the subset of server information returned with a network port.
type ServerSummary struct {
	ID        string `json:"id"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
	ProjectID string `json:"project_id,omitempty"`
	UserID    string `json:"user_id,omitempty"`
}
