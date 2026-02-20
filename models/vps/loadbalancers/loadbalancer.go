package loadbalancers

import (
	"errors"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
)

// LoadBalancerStatus represents the current state of a load balancer.
type LoadBalancerStatus string

// LoadBalancerStatus constants for load balancer lifecycle states.
const (
	LoadBalancerStatusPendingCreate LoadBalancerStatus = "PENDING_CREATE" // LB is being created
	LoadBalancerStatusActive        LoadBalancerStatus = "ACTIVE"         // LB is active and ready
	LoadBalancerStatusPendingUpdate LoadBalancerStatus = "PENDING_UPDATE" // LB is being updated
	LoadBalancerStatusPendingDelete LoadBalancerStatus = "PENDING_DELETE" // LB is being deleted
	LoadBalancerStatusDeleted       LoadBalancerStatus = "DELETED"        // LB has been deleted
	LoadBalancerStatusError         LoadBalancerStatus = "ERROR"          // Operation failed
)

// String returns the string representation of the status.
func (s LoadBalancerStatus) String() string {
	return string(s)
}

// LoadBalancer represents a load balancer resource.
// Matches pb.LbInfo from vps.yaml.
type LoadBalancer struct {
	ID           string                  `json:"id"`
	Name         string                  `json:"name"`
	Description  string                  `json:"description,omitempty"`
	Status       LoadBalancerStatus      `json:"status"`
	StatusReason string                  `json:"status_reason,omitempty"`
	VIP          string                  `json:"vip,omitempty"`
	PortID       string                  `json:"port_id,omitempty"`
	NetworkID    string                  `json:"network_id"`
	Network      common.IDName           `json:"network"`
	ProjectID    string                  `json:"project_id"`
	Project      common.IDName           `json:"project"`
	UserID       string                  `json:"user_id"`
	User         common.IDName           `json:"user"`
	Namespace    string                  `json:"namespace,omitempty"`
	FloatingIP   *floatingips.FloatingIP `json:"floating_ip,omitempty"`
	Listeners    []Listener              `json:"listeners,omitempty"`
	CreatedAt    *time.Time              `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time              `json:"updatedAt,omitempty"`
}

// LoadBalancerCreateRequest represents parameters for creating a load balancer.
// Matches LbCreateInput from vps.yaml.
type LoadBalancerCreateRequest struct {
	Name        string                  `json:"name"`                  // Required: LB name
	NetworkID   string                  `json:"network_id"`            // Required: Network UUID
	Description string                  `json:"description,omitempty"` // Optional: LB description
	Listeners   []ListenerCreateRequest `json:"listeners,omitempty"`   // Optional: Inline listener creation
}

// Validate checks if the LoadBalancerCreateRequest has valid values.
func (r *LoadBalancerCreateRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.NetworkID == "" {
		return errors.New("network_id is required")
	}
	return nil
}

// LoadBalancerUpdateRequest represents parameters for updating a load balancer.
// Matches LbUpdateInput from vps.yaml.
type LoadBalancerUpdateRequest struct {
	Name        string `json:"name,omitempty"`        // Optional: New LB name
	Description string `json:"description,omitempty"` // Optional: New LB description
}

// Validate checks if the LoadBalancerUpdateRequest has valid values.
func (r *LoadBalancerUpdateRequest) Validate() error {
	if r.Name == "" && r.Description == "" {
		return errors.New("at least one field (name or description) must be provided")
	}
	return nil
}

// AssociateFloatingIPRequest represents parameters for associating a floating IP to a load balancer.
// Matches FIPAssoInput from vps.yaml.
type AssociateFloatingIPRequest struct {
	FIPID string `json:"fip_id,omitempty"` // Optional: Existing floating IP ID. If empty, a new one is created.
}

// Secret represents TLS certificate and key for TERMINATED_HTTPS listeners.
// Matches user.SecretCreateInput from vps.yaml.
type Secret struct {
	Certificate string `json:"certificate,omitempty"` // PKCS12 format certificate (Base64 encoded)
	Key         string `json:"key,omitempty"`         // Private key (Base64 encoded)
}

// ListLoadBalancersOptions provides filtering options for listing load balancers.
type ListLoadBalancersOptions struct {
	Name   string // Filter by name
	UserID string // Filter by user ID (admin only)
	Detail bool   // If true, return LB with listeners and pools
}

// LoadBalancerListResponse represents the response from listing load balancers.
// Matches pb.LbListOutput from vps.yaml.
type LoadBalancerListResponse struct {
	LoadBalancers []*LoadBalancer `json:"load_balancers"`
}
