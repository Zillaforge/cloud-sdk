package shares

import (
	"errors"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

// ShareStatus represents the current state of a share.
type ShareStatus string

// ShareStatus constants for share lifecycle states.
const (
	ShareStatusCreating  ShareStatus = "creating"  // Share is being created
	ShareStatusAvailable ShareStatus = "available" // Share created and ready for use
	ShareStatusDeleting  ShareStatus = "deleting"  // Share is being deleted
	ShareStatusDeleted   ShareStatus = "deleted"   // Share has been deleted
	ShareStatusError     ShareStatus = "error"     // Operation failed, check StatusReason
	ShareStatusInactive  ShareStatus = "inactive"  // Share is inactive
)

// ShareRuleLevel represents access level for a share rule.
type ShareRuleLevel string

// ShareRuleLevel constants for access levels.
const (
	ShareRuleLevelRW ShareRuleLevel = "rw" // Read-write access
	ShareRuleLevelRO ShareRuleLevel = "ro" // Read-only access
)

// Share represents a managed file share resource.
// Matches pb.ShareInfo from vps.yaml.
type Share struct {
	ID           string        `json:"id"`
	Name         string        `json:"name"`
	Description  string        `json:"description,omitempty"`
	Size         int           `json:"size"`                    // Size in GB
	Path         string        `json:"path,omitempty"`          // Share export path
	Status       ShareStatus   `json:"status"`                  // Share status
	StatusReason string        `json:"status_reason,omitempty"` // Status reason text
	Network      common.IDName `json:"network"`                 // Network reference
	NetworkID    string        `json:"network_id"`
	Project      common.IDName `json:"project"` // Project reference
	ProjectID    string        `json:"project_id"`
	User         common.IDName `json:"user"` // User reference
	UserID       string        `json:"user_id"`
	Namespace    string        `json:"namespace"`
	Rules        []ShareRule   `json:"rules,omitempty"`     // Access rules
	CreatedAt    *time.Time    `json:"createdAt,omitempty"` // Creation timestamp
	UpdatedAt    *time.Time    `json:"updatedAt,omitempty"` // Last update timestamp
}

// Validate checks if the Share has valid required fields.
func (s *Share) Validate() error {
	if s.ID == "" {
		return errors.New("id is required")
	}
	if s.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// ShareRule represents an access rule on a share.
// Matches pb.ShareRuleInfo from vps.yaml.
type ShareRule struct {
	ID     string `json:"id"`
	CIDR   string `json:"cidr"`
	Level  string `json:"level"`
	Status string `json:"status,omitempty"`
}

// CreateShareRequest represents parameters for creating a share.
// Matches user.ShareCreateInput from vps.yaml.
type CreateShareRequest struct {
	Name        string `json:"name"`                  // Required: Share name
	NetworkID   string `json:"network_id"`            // Required: Network UUID hosting the share
	Size        int    `json:"size"`                  // Required: Size in GB
	Description string `json:"description,omitempty"` // Optional: Share description
}

// Validate checks if the CreateShareRequest has valid values.
func (r *CreateShareRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.NetworkID == "" {
		return errors.New("network_id is required")
	}
	if r.Size <= 0 {
		return errors.New("size must be positive")
	}
	return nil
}

// UpdateShareRequest represents parameters for updating a share.
// Matches ShareUpdateInput from vps.yaml.
type UpdateShareRequest struct {
	Name        string `json:"name,omitempty"`        // Optional: New share name
	Description string `json:"description,omitempty"` // Optional: New share description
}

// Validate checks if the UpdateShareRequest has valid values.
func (r *UpdateShareRequest) Validate() error {
	if r.Name == "" && r.Description == "" {
		return errors.New("at least one field (name or description) must be provided")
	}
	return nil
}

// ExtendShareRequest represents parameters for extending a share's size.
// Matches ShareExtendInput from vps.yaml.
type ExtendShareRequest struct {
	NewSize int `json:"new_size"` // Required: New size in GB (must be larger than current)
}

// Validate checks if the ExtendShareRequest has valid values.
func (r *ExtendShareRequest) Validate() error {
	if r.NewSize <= 0 {
		return errors.New("new_size must be positive")
	}
	return nil
}

// AddShareRuleRequest represents parameters for adding an access rule to a share.
// Matches ShareRuleCreateInput from vps.yaml.
type AddShareRuleRequest struct {
	CIDR  string `json:"cidr"`  // Required: IP address or CIDR range (e.g., "10.50.0.0/24")
	Level string `json:"level"` // Required: Access level ("rw" or "ro")
}

// Validate checks if the AddShareRuleRequest has valid values.
func (r *AddShareRuleRequest) Validate() error {
	if r.CIDR == "" {
		return errors.New("cidr is required")
	}
	if r.Level == "" {
		return errors.New("level is required")
	}
	if r.Level != string(ShareRuleLevelRW) && r.Level != string(ShareRuleLevelRO) {
		return errors.New("level must be one of rw, ro")
	}
	return nil
}

// ListSharesOptions provides filtering options for listing shares.
type ListSharesOptions struct {
	Name   string // Filter by name
	UserID string // Filter by user ID (admin only)
	Detail bool   // If true, return with rules
}

// ShareListResponse represents the response from listing shares.
// Matches pb.ShareListOutput from vps.yaml.
type ShareListResponse struct {
	Shares []*Share `json:"shares"`
}
