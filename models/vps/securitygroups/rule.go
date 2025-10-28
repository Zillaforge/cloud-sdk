// Package securitygroups provides data structures for VPS security group resources.
package securitygroups

// SecurityGroupRule represents a security group rule.
type SecurityGroupRule struct {
	ID         string `json:"id"`
	Direction  string `json:"direction"`
	Protocol   string `json:"protocol"`
	PortMin    *int   `json:"port_min,omitempty"`
	PortMax    *int   `json:"port_max,omitempty"`
	RemoteCIDR string `json:"remote_cidr"`
}

// SecurityGroupRuleCreateRequest represents the request body for creating a security group rule.
type SecurityGroupRuleCreateRequest struct {
	Direction  string `json:"direction"`
	Protocol   string `json:"protocol"`
	PortMin    *int   `json:"port_min,omitempty"`
	PortMax    *int   `json:"port_max,omitempty"`
	RemoteCIDR string `json:"remote_cidr"`
}
