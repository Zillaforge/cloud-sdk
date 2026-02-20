package loadbalancers

import (
	"errors"
)

// ListenerProtocol represents the protocol type for a listener.
type ListenerProtocol string

// ListenerProtocol constants for supported listener protocols.
const (
	ListenerProtocolHTTP            ListenerProtocol = "HTTP"
	ListenerProtocolHTTPS           ListenerProtocol = "HTTPS"
	ListenerProtocolTerminatedHTTPS ListenerProtocol = "TERMINATED_HTTPS"
	ListenerProtocolTCP             ListenerProtocol = "TCP"
	ListenerProtocolUDP             ListenerProtocol = "UDP"
)

// Valid returns true if the protocol is a valid ListenerProtocol value.
func (p ListenerProtocol) Valid() bool {
	switch p {
	case ListenerProtocolHTTP, ListenerProtocolHTTPS, ListenerProtocolTerminatedHTTPS, ListenerProtocolTCP, ListenerProtocolUDP:
		return true
	}
	return false
}

// Listener represents a load balancer listener resource.
// Matches pb.ListenerInfo from vps.yaml.
type Listener struct {
	ID                     string            `json:"id"`
	Name                   string            `json:"name,omitempty"`
	LbID                   string            `json:"lb_id,omitempty"`
	Protocol               ListenerProtocol  `json:"protocol"`
	ProtocolPort           int               `json:"protocol_port"`
	PoolID                 string            `json:"pool_id,omitempty"`
	Pool                   *Pool             `json:"pool,omitempty"`
	Status                 string            `json:"status,omitempty"`
	StatusReason           string            `json:"status_reason,omitempty"`
	AllowedCIDRs           []string          `json:"allowed_cidrs,omitempty"`
	InsertHeaders          map[string]string `json:"insert_headers,omitempty"`
	DefaultTLSContainerRef string            `json:"default_tls_container_ref,omitempty"`
	TimeoutClientData      int               `json:"timeout_client_data,omitempty"`
	TimeoutMemberConnect   int               `json:"timeout_member_connect,omitempty"`
	TimeoutMemberData      int               `json:"timeout_member_data,omitempty"`
	TimeoutTCPInspect      int               `json:"timeout_tcp_inspect,omitempty"`
}

// ListenerCreateRequest represents parameters for creating a listener.
// Matches ListenerCreateInput from vps.yaml.
type ListenerCreateRequest struct {
	Protocol             ListenerProtocol   `json:"protocol"`                         // Required: HTTP, HTTPS, TERMINATED_HTTPS, TCP, UDP
	ProtocolPort         int                `json:"protocol_port"`                    // Required: Port number
	Name                 string             `json:"name,omitempty"`                   // Optional: Listener name
	PoolID               string             `json:"pool_id,omitempty"`                // Optional: Existing pool ID
	Pool                 *PoolCreateRequest `json:"pool,omitempty"`                   // Optional: Inline pool creation
	Secret               *Secret            `json:"secret,omitempty"`                 // Optional: TLS cert for TERMINATED_HTTPS
	AllowedCIDRs         []string           `json:"allowed_cidrs,omitempty"`          // Optional: Allowed CIDRs
	InsertHeaders        map[string]string  `json:"insert_headers,omitempty"`         // Optional: Headers to insert
	TimeoutClientData    int                `json:"timeout_client_data,omitempty"`    // Optional: Client timeout (ms)
	TimeoutMemberConnect int                `json:"timeout_member_connect,omitempty"` // Optional: Member connect timeout (ms)
	TimeoutMemberData    int                `json:"timeout_member_data,omitempty"`    // Optional: Member data timeout (ms)
	TimeoutTCPInspect    int                `json:"timeout_tcp_inspect,omitempty"`    // Optional: TCP inspect timeout (ms)
}

// Validate checks if the ListenerCreateRequest has valid values.
func (r *ListenerCreateRequest) Validate() error {
	if r.Protocol == "" {
		return errors.New("protocol is required")
	}
	if !r.Protocol.Valid() {
		return errors.New("protocol must be one of HTTP, HTTPS, TERMINATED_HTTPS, TCP, UDP")
	}
	if r.ProtocolPort == 0 {
		return errors.New("protocol_port is required")
	}
	return nil
}

// ListenerUpdateRequest represents parameters for updating a listener.
// Matches user.ListenerUpdateInput from vps.yaml.
type ListenerUpdateRequest struct {
	Name                 string            `json:"name,omitempty"`                   // Optional: New name
	PoolID               string            `json:"pool_id,omitempty"`                // Optional: Pool ID to associate
	Secret               *Secret           `json:"secret,omitempty"`                 // Optional: TLS cert
	AllowedCIDRs         []string          `json:"allowed_cidrs,omitempty"`          // Optional: Allowed CIDRs
	InsertHeaders        map[string]string `json:"insert_headers,omitempty"`         // Optional: Headers
	TimeoutClientData    *int              `json:"timeout_client_data,omitempty"`    // Optional: Client timeout
	TimeoutMemberConnect *int              `json:"timeout_member_connect,omitempty"` // Optional: Member connect timeout
	TimeoutMemberData    *int              `json:"timeout_member_data,omitempty"`    // Optional: Member data timeout
	TimeoutTCPInspect    *int              `json:"timeout_tcp_inspect,omitempty"`    // Optional: TCP inspect timeout
}

// ListenerListResponse represents the response from listing listeners.
// Matches pb.ListenerListOutput from vps.yaml.
type ListenerListResponse struct {
	Listeners []*Listener `json:"listeners"`
}
