package volumes

import (
	"errors"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

// VolumeStatus represents the current state of a volume.
type VolumeStatus string

// VolumeStatus constants for volume lifecycle states.
const (
	VolumeStatusCreating  VolumeStatus = "creating"  // Volume is being created
	VolumeStatusAvailable VolumeStatus = "available" // Volume created and ready to attach
	VolumeStatusInUse     VolumeStatus = "in-use"    // Volume attached to server
	VolumeStatusDetaching VolumeStatus = "detaching" // Volume being detached from server
	VolumeStatusExtending VolumeStatus = "extending" // Volume size is being increased
	VolumeStatusDeleting  VolumeStatus = "deleting"  // Volume is being deleted
	VolumeStatusDeleted   VolumeStatus = "deleted"   // Volume deleted
	VolumeStatusError     VolumeStatus = "error"     // Operation failed, check StatusReason
)

// VolumeAction represents supported volume action types.
type VolumeAction string

// VolumeAction constants for supported volume actions.
const (
	VolumeActionAttach VolumeAction = "attach" // Attach volume to server
	VolumeActionDetach VolumeAction = "detach" // Detach volume from server
	VolumeActionExtend VolumeAction = "extend" // Extend volume size
	VolumeActionRevert VolumeAction = "revert" // Revert volume to snapshot
)

// Volume represents a block storage volume.
// Matches pegasus-cloud_com_aes_virtualplatformserviceclient_pb.VolumeInfo from vps.yaml.
type Volume struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description,omitempty"`
	Size         int             `json:"size"`                    // Size in GB
	Type         string          `json:"type"`                    // Storage type (SSD, HDD, etc.)
	Status       VolumeStatus    `json:"status"`                  // Volume status
	StatusReason string          `json:"status_reason,omitempty"` // Status reason text
	Attachments  []common.IDName `json:"attachments,omitempty"`   // Servers this volume is attached to
	Project      common.IDName   `json:"project"`                 // Project reference
	ProjectID    string          `json:"project_id"`
	User         common.IDName   `json:"user"` // User reference
	UserID       string          `json:"user_id"`
	Namespace    string          `json:"namespace"`
	CreatedAt    *time.Time      `json:"createdAt,omitempty"` // Creation timestamp (ISO 8601)
	UpdatedAt    *time.Time      `json:"updatedAt,omitempty"` // Last update timestamp (ISO 8601)
}

// Validate checks if the Volume has valid required fields and data types.
func (v *Volume) Validate() error {
	if v.ID == "" {
		return errors.New("id is required")
	}
	if v.Name == "" {
		return errors.New("name is required")
	}
	if v.Size < 0 {
		return errors.New("size must be >= 0")
	}
	if v.Type == "" {
		return errors.New("type is required")
	}
	return nil
}

// CreateVolumeRequest represents parameters for creating a volume.
// Matches VolumeCreateInput from vps.yaml.
type CreateVolumeRequest struct {
	Name        string `json:"name"`                  // Required: Volume name
	Type        string `json:"type"`                  // Required: Volume type (from VolumeTypes)
	Size        int    `json:"size,omitempty"`        // Optional: Size in GB (server validates)
	Description string `json:"description,omitempty"` // Optional: Volume description
	SnapshotID  string `json:"snapshot_id,omitempty"` // Optional: Create from snapshot
}

// Validate checks if the CreateVolumeRequest has valid values.
func (r *CreateVolumeRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.Type == "" {
		return errors.New("type is required")
	}
	if r.Size < 0 {
		return errors.New("size must be >= 0 if provided")
	}
	return nil
}

// UpdateVolumeRequest represents parameters for updating a volume.
// Matches VolumeUpdateInput from vps.yaml.
type UpdateVolumeRequest struct {
	Name        string `json:"name,omitempty"`        // Optional: New volume name
	Description string `json:"description,omitempty"` // Optional: New volume description
}

// Validate checks if the UpdateVolumeRequest has valid values.
func (r *UpdateVolumeRequest) Validate() error {
	// At least one field should be provided
	if r.Name == "" && r.Description == "" {
		return errors.New("at least one field (name or description) must be provided")
	}
	return nil
}

// VolumeActionRequest represents parameters for volume actions.
// Matches VolActionInput from vps.yaml.
type VolumeActionRequest struct {
	Action   VolumeAction `json:"action"`              // Required: Action type
	ServerID string       `json:"server_id,omitempty"` // Required for attach/detach
	NewSize  int          `json:"new_size,omitempty"`  // Required for extend
}

// Validate checks if the VolumeActionRequest has valid action-specific parameters.
func (r *VolumeActionRequest) Validate() error {
	switch r.Action {
	case VolumeActionAttach, VolumeActionDetach:
		if r.ServerID == "" {
			return errors.New("server_id is required for attach/detach actions")
		}
	case VolumeActionExtend:
		if r.NewSize <= 0 {
			return errors.New("new_size must be positive for extend action")
		}
	case VolumeActionRevert:
		// No additional parameters required
	default:
		return errors.New("invalid action: must be one of attach, detach, extend, revert")
	}
	return nil
}

// ListVolumesOptions provides filtering options for listing volumes.
type ListVolumesOptions struct {
	Name   string // Filter by name (partial match)
	UserID string // Filter by user ID
	Status string // Filter by status (available, in-use, etc.)
	Type   string // Filter by volume type
	Detail bool   // Include attachment details
}

// Validate checks if the ListVolumesOptions has valid values.
func (o *ListVolumesOptions) Validate() error {
	// All fields are optional, no validation needed
	return nil
}

// VolumeListResponse represents the response from listing volumes.
// Matches pb.VolumeListOutput from vps.yaml.
type VolumeListResponse struct {
	Volumes []*Volume `json:"volumes"`
}

// VolumeResponse represents the response containing a single volume.
type VolumeResponse struct {
	Volume *Volume `json:"volume,omitempty"`
}
