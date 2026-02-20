package loadbalancers

import (
	"errors"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

// PoolMethod represents the load balancing algorithm.
type PoolMethod string

// PoolMethod constants for load balancing algorithms.
const (
	PoolMethodRoundRobin       PoolMethod = "ROUND_ROBIN"
	PoolMethodLeastConnections PoolMethod = "LEAST_CONNECTIONS"
	PoolMethodSourceIP         PoolMethod = "SOURCE_IP"
)

// Valid returns true if the method is a valid PoolMethod value.
func (m PoolMethod) Valid() bool {
	switch m {
	case PoolMethodRoundRobin, PoolMethodLeastConnections, PoolMethodSourceIP:
		return true
	}
	return false
}

// PoolProtocol represents the protocol type for a pool.
type PoolProtocol string

// PoolProtocol constants for supported pool protocols.
const (
	PoolProtocolHTTP  PoolProtocol = "HTTP"
	PoolProtocolHTTPS PoolProtocol = "HTTPS"
	PoolProtocolTCP   PoolProtocol = "TCP"
	PoolProtocolUDP   PoolProtocol = "UDP"
)

// Valid returns true if the protocol is a valid PoolProtocol value.
func (p PoolProtocol) Valid() bool {
	switch p {
	case PoolProtocolHTTP, PoolProtocolHTTPS, PoolProtocolTCP, PoolProtocolUDP:
		return true
	}
	return false
}

// Pool represents a load balancer pool resource.
// Matches pb.PoolInfo from vps.yaml.
type Pool struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Protocol     PoolProtocol    `json:"protocol"`
	MemberPort   int             `json:"member_port"`
	Method       PoolMethod      `json:"method"`
	Status       string          `json:"status,omitempty"`
	StatusReason string          `json:"status_reason,omitempty"`
	ASG          *common.IDName  `json:"asg,omitempty"`
	Listeners    []common.IDName `json:"listeners,omitempty"`
	Members      []Member        `json:"members,omitempty"`
}

// Member represents a pool member.
// Matches pb.MemberInfo from vps.yaml.
type Member struct {
	ID           string `json:"id"`
	Name         string `json:"name,omitempty"`
	Address      string `json:"address"`
	ProtocolPort int    `json:"protocol_port,omitempty"`
	Status       string `json:"status,omitempty"`
}

// PoolCreateRequest represents parameters for creating a pool.
// Matches PoolCreateInput from vps.yaml.
type PoolCreateRequest struct {
	Name       string                `json:"name"`              // Required: Pool name
	Protocol   PoolProtocol          `json:"protocol"`          // Required: HTTP, HTTPS, TCP, UDP
	MemberPort int                   `json:"member_port"`       // Required: Backend member port
	Method     PoolMethod            `json:"method"`            // Required: ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP
	ASGID      string                `json:"asg_id,omitempty"`  // Optional: Auto-scale group ID
	Members    []MemberCreateRequest `json:"members,omitempty"` // Optional: Initial members
}

// Validate checks if the PoolCreateRequest has valid values.
func (r *PoolCreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Protocol == "" {
		return errors.New("protocol is required")
	}
	if !r.Protocol.Valid() {
		return errors.New("protocol must be one of HTTP, HTTPS, TCP, UDP")
	}
	if r.MemberPort == 0 {
		return errors.New("member_port is required")
	}
	if r.Method == "" {
		return errors.New("method is required")
	}
	if !r.Method.Valid() {
		return errors.New("method must be one of ROUND_ROBIN, LEAST_CONNECTIONS, SOURCE_IP")
	}
	return nil
}

// PoolUpdateRequest represents parameters for updating a pool.
// Matches user.PoolUpdateInput from vps.yaml.
type PoolUpdateRequest struct {
	Name       string     `json:"name,omitempty"`        // Optional: New pool name
	MemberPort *int       `json:"member_port,omitempty"` // Optional: New member port
	Method     PoolMethod `json:"method,omitempty"`      // Optional: New LB method
	ASGID      string     `json:"asg_id,omitempty"`      // Optional: Auto-scale group ID
}

// PoolListResponse represents the response from listing pools.
// Matches pb.PoolListOutput from vps.yaml.
type PoolListResponse struct {
	Pools []*Pool `json:"pools"`
}

// MemberCreateRequest represents parameters for creating a member.
// Matches MemberCreateInput from vps.yaml.
type MemberCreateRequest struct {
	Address string `json:"address"`        // Required: IPv4 address of backend member
	Name    string `json:"name,omitempty"` // Optional: Member name
}

// Validate checks if the MemberCreateRequest has valid values.
func (r *MemberCreateRequest) Validate() error {
	if r.Address == "" {
		return errors.New("address is required")
	}
	return nil
}

// MembersBatchUpdateRequest represents parameters for batch updating members.
// Matches MembersUpdateInput from vps.yaml.
// Note: This replaces the entire member list for the pool.
type MembersBatchUpdateRequest struct {
	Members []MemberCreateRequest `json:"members"` // Full list of desired members
}
