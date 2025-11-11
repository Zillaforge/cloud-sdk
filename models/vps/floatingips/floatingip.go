package floatingips

import (
	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

// FloatingIPStatus represents the status of a floating IP resource.
type FloatingIPStatus string

const (
	FloatingIPStatusActive   FloatingIPStatus = "ACTIVE"
	FloatingIPStatusPending  FloatingIPStatus = "PENDING"
	FloatingIPStatusDown     FloatingIPStatus = "DOWN"
	FloatingIPStatusRejected FloatingIPStatus = "REJECTED"
)

// String returns the string representation of the status.
func (s FloatingIPStatus) String() string {
	return string(s)
}

// Valid returns true if the status is a valid FloatingIPStatus value.
func (s FloatingIPStatus) Valid() bool {
	switch s {
	case FloatingIPStatusActive, FloatingIPStatusPending, FloatingIPStatusDown, FloatingIPStatusRejected:
		return true
	}
	return false
}

// FloatingIP represents a floating IP address resource.
type FloatingIP struct {
	// Identity
	ID   string `json:"id"`
	UUID string `json:"uuid"`
	Name string `json:"name"`

	// Network & Allocation
	Address   string `json:"address"`
	ExtNetID  string `json:"extnet_id,omitempty"`
	PortID    string `json:"port_id,omitempty"`
	ProjectID string `json:"project_id"`
	Namespace string `json:"namespace,omitempty"`

	// Ownership & References
	UserID  string         `json:"user_id"`
	User    *common.IDName `json:"user,omitempty"`
	Project *common.IDName `json:"project,omitempty"`

	// Association Details
	DeviceID   string `json:"device_id,omitempty"`
	DeviceName string `json:"device_name,omitempty"`
	DeviceType string `json:"device_type,omitempty"`

	// Status & Metadata
	Description  string           `json:"description,omitempty"`
	Status       FloatingIPStatus `json:"status"`
	StatusReason string           `json:"status_reason,omitempty"`
	Reserved     bool             `json:"reserved"`

	// Lifecycle Timestamps
	CreatedAt  string `json:"createdAt"`
	UpdatedAt  string `json:"updatedAt,omitempty"`
	ApprovedAt string `json:"approvedAt,omitempty"`
}

// FloatingIPCreateRequest represents the request to create a new floating IP.
type FloatingIPCreateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	ExtNetID    string `json:"extnet_id,omitempty"`
}

// FloatingIPUpdateRequest represents the request to update an existing floating IP.
type FloatingIPUpdateRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	Reserved    *bool  `json:"reserved,omitempty"`
}

// FloatingIPListResponse represents the response from listing floating IPs (deprecated - use direct slice).
type FloatingIPListResponse struct {
	FloatingIPs []*FloatingIP `json:"floating_ips"`
}

// ListFloatingIPsOptions provides filtering options for floating IP listing.
type ListFloatingIPsOptions struct {
	Status     string // filter by status
	UserID     string // filter by user
	DeviceType string // filter by device type
	DeviceID   string // filter by device id
	ExtNetID   string // filter by external network
	Address    string // filter by address
	Name       string // filter by name
	Detail     bool   // return detailed response
}
