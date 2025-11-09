// Package securitygroups provides data structures for VPS security group resources.
package securitygroups

// Protocol represents the network protocol for a security group rule.
type Protocol string

const (
	// ProtocolTCP represents the TCP protocol.
	ProtocolTCP Protocol = "tcp"
	// ProtocolUDP represents the UDP protocol.
	ProtocolUDP Protocol = "udp"
	// ProtocolICMP represents the ICMP protocol.
	ProtocolICMP Protocol = "icmp"
	// ProtocolAny represents any protocol (wildcard).
	ProtocolAny Protocol = "any"
)

// Direction represents the traffic direction for a security group rule.
type Direction string

const (
	// DirectionIngress represents inbound traffic.
	DirectionIngress Direction = "ingress"
	// DirectionEgress represents outbound traffic.
	DirectionEgress Direction = "egress"
)

// SecurityGroupRule represents a security group rule.
type SecurityGroupRule struct {
	ID         string    `json:"id"`
	Direction  Direction `json:"direction"`
	Protocol   Protocol  `json:"protocol"`
	PortMin    int       `json:"port_min"`
	PortMax    int       `json:"port_max"`
	RemoteCIDR string    `json:"remote_cidr"`
}

// SecurityGroupRuleCreateRequest represents the request body for creating a security group rule.
type SecurityGroupRuleCreateRequest struct {
	Direction  Direction `json:"direction"`
	Protocol   Protocol  `json:"protocol"`
	PortMin    *int      `json:"port_min,omitempty"`
	PortMax    *int      `json:"port_max,omitempty"`
	RemoteCIDR string    `json:"remote_cidr"`
}
