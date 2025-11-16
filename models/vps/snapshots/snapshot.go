package snapshots

import (
	"errors"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

// SnapshotStatus represents the current state of a snapshot.
type SnapshotStatus string

const (
	SnapshotStatusCreating  SnapshotStatus = "creating"
	SnapshotStatusAvailable SnapshotStatus = "available"
	SnapshotStatusDeleting  SnapshotStatus = "deleting"
	SnapshotStatusDeleted   SnapshotStatus = "deleted"
	SnapshotStatusError     SnapshotStatus = "error"
)

// Snapshot represents a point-in-time capture of a data volume.
type Snapshot struct {
	ID           string         `json:"id"`
	Name         string         `json:"name"`
	VolumeID     string         `json:"volume_id"`
	Size         int            `json:"size"`
	Status       SnapshotStatus `json:"status"`
	StatusReason string         `json:"status_reason,omitempty"`
	Description  string         `json:"description,omitempty"`
	Project      common.IDName  `json:"project"`
	ProjectID    string         `json:"project_id"`
	User         common.IDName  `json:"user"`
	UserID       string         `json:"user_id"`
	Namespace    string         `json:"namespace"`
	CreatedAt    *time.Time     `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time     `json:"updatedAt,omitempty"`
}

// Validate ensures required Snapshot fields are present.
func (s *Snapshot) Validate() error {
	if s.ID == "" {
		return errors.New("id is required")
	}
	if s.VolumeID == "" {
		return errors.New("volume_id is required")
	}
	return nil
}

// CreateSnapshotRequest represents the payload for creating a snapshot.
type CreateSnapshotRequest struct {
	Name     string `json:"name"`
	VolumeID string `json:"volume_id"`
}

// Validate checks required fields for snapshot creation.
func (r *CreateSnapshotRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	if r.VolumeID == "" {
		return errors.New("volume_id is required")
	}
	return nil
}

// UpdateSnapshotRequest represents rename/metadata updates for a snapshot.
type UpdateSnapshotRequest struct {
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
}

// Validate ensures at least one field is present for update.
func (r *UpdateSnapshotRequest) Validate() error {
	if r.Name == "" {
		return errors.New("name is required")
	}
	return nil
}

// ListSnapshotsOptions provides filter options for listing snapshots.
type ListSnapshotsOptions struct {
	Name     string
	VolumeID string
	UserID   string
	Status   string
}

// Validate ensures List options are valid (none are required).
func (o *ListSnapshotsOptions) Validate() error {
	return nil
}

// SnapshotListResponse models response list for snapshots.
type SnapshotListResponse struct {
	Snapshots []*Snapshot `json:"snapshots"`
}

// SnapshotResponse models response for a single snapshot.
type SnapshotResponse struct {
	Snapshot *Snapshot `json:"snapshot,omitempty"`
}
