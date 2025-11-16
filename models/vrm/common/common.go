package common

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// IDName represents a lightweight reference to an entity with id and name.
// Used for Creator, Project, User references throughout the VRM API.
type IDName struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Account     string `json:"account,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// Validate validates the IDName structure.
// At minimum, ID must be present.
func (n *IDName) Validate() error {
	if n == nil {
		return fmt.Errorf("IDName cannot be nil")
	}
	if strings.TrimSpace(n.ID) == "" {
		return fmt.Errorf("IDName.ID is required and must not be empty")
	}
	return nil
}

// Repository represents a virtual image repository within a project.
// It contains metadata about the repository including its name, namespace,
// operating system, description, and associated tags. The repository
// tracks creation and update timestamps, as well as creator and project information.
type Repository struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	OperatingSystem string    `json:"operatingSystem"`
	Description     string    `json:"description,omitempty"`
	Tags            []*Tag    `json:"tags,omitempty"`
	Count           int       `json:"count"`
	Creator         *IDName   `json:"creator"`
	Project         *IDName   `json:"project"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// Validate validates the Repository structure.
// It checks that required fields are present and have valid values,
// including ID, name, namespace, operating system, and count constraints.
func (r *Repository) Validate() error {
	if r == nil {
		return fmt.Errorf("repository cannot be nil")
	}
	if strings.TrimSpace(r.ID) == "" {
		return fmt.Errorf("ID is required and must not be empty")
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required and must not be empty")
	}
	if r.Namespace != "public" && r.Namespace != "private" {
		return fmt.Errorf("namespace must be 'public' or 'private'")
	}
	if r.OperatingSystem != "linux" && r.OperatingSystem != "windows" {
		return fmt.Errorf("operating system must be 'linux' or 'windows'")
	}
	if r.Count < 0 {
		return fmt.Errorf("count must be >= 0")
	}
	return nil
}

// MarshalJSON implements json.Marshaler for Repository.
// It customizes JSON marshaling to format timestamps as RFC3339 strings.
func (r *Repository) MarshalJSON() ([]byte, error) {
	type Alias Repository
	return json.Marshal(&struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias:     (*Alias)(r),
		CreatedAt: r.CreatedAt.Format(time.RFC3339),
		UpdatedAt: r.UpdatedAt.Format(time.RFC3339),
	})
}

// UnmarshalJSON implements json.Unmarshaler for Repository.
// It customizes JSON unmarshaling to parse RFC3339 timestamp strings into time.Time values.
func (r *Repository) UnmarshalJSON(data []byte) error {
	type Alias Repository
	aux := &struct {
		*Alias
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
	}{
		Alias: (*Alias)(r),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	if aux.CreatedAt != "" {
		t, err := time.Parse(time.RFC3339, aux.CreatedAt)
		if err != nil {
			return err
		}
		r.CreatedAt = t
	}
	if aux.UpdatedAt != "" {
		t, err := time.Parse(time.RFC3339, aux.UpdatedAt)
		if err != nil {
			return err
		}
		r.UpdatedAt = t
	}
	return nil
}

// TagType represents the type of a tag in the repository.
// Valid values are "common" and "increase".
type TagType string

const (
	TagTypeCommon   TagType = "common"
	TagTypeIncrease TagType = "increase"
)

// IsValid checks if the TagType is a valid value.
func (t TagType) IsValid() bool {
	switch t {
	case TagTypeCommon, TagTypeIncrease:
		return true
	}
	return false
}

// String returns the string representation of TagType.
func (t TagType) String() string {
	return string(t)
}

// TagStatus represents the status of a tag in the repository.
// Valid values include various states like queued, active, error, etc.
type TagStatus string

const (
	TagStatusQueued        TagStatus = "queued"
	TagStatusSaving        TagStatus = "saving"
	TagStatusImporting     TagStatus = "importing"
	TagStatusCreating      TagStatus = "creating"
	TagStatusRestoring     TagStatus = "restoring"
	TagStatusActive        TagStatus = "active"
	TagStatusKilled        TagStatus = "killed"
	TagStatusPendingDelete TagStatus = "pending_delete"
	TagStatusDeactivated   TagStatus = "deactivated"
	TagStatusAvailable     TagStatus = "available"
	TagStatusBackingUp     TagStatus = "backing-up"
	TagStatusDeleting      TagStatus = "deleting"
	TagStatusError         TagStatus = "error"
	TagStatusUnmanaging    TagStatus = "unmanaging"
	TagStatusErrorDeleting TagStatus = "error_deleting"
	TagStatusDeleted       TagStatus = "deleted"
)

// IsValid checks if the TagStatus is a valid value.
func (s TagStatus) IsValid() bool {
	switch s {
	case TagStatusQueued, TagStatusSaving, TagStatusImporting, TagStatusCreating, TagStatusRestoring,
		TagStatusActive, TagStatusKilled, TagStatusPendingDelete, TagStatusDeactivated, TagStatusAvailable,
		TagStatusBackingUp, TagStatusDeleting, TagStatusError, TagStatusUnmanaging, TagStatusErrorDeleting, TagStatusDeleted:
		return true
	}
	return false
}

// String returns the string representation of TagStatus.
func (s TagStatus) String() string {
	return string(s)
}

// Tag represents a tagged version of an image within a repository.
// Each tag contains metadata about a specific version of an image,
// including its size, status, and any extra properties. Tags are
// associated with a specific repository and track creation/update timestamps.
type Tag struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	RepositoryID string                 `json:"repositoryID"`
	Type         TagType                `json:"type"`
	Size         int64                  `json:"size"`
	Status       TagStatus              `json:"status,omitempty"`
	Extra        map[string]interface{} `json:"extra,omitempty"`
	CreatedAt    time.Time              `json:"createdAt"`
	UpdatedAt    time.Time              `json:"updatedAt"`
	Repository   *Repository            `json:"repository"`
}

// Validate validates the Tag struct.
// It ensures all required fields are present and have valid values,
// including non-empty ID, name, repositoryID, type, and non-negative size.
func (t *Tag) Validate() error {
	if t == nil {
		return fmt.Errorf("tag cannot be nil")
	}
	if t.ID == "" {
		return fmt.Errorf("id cannot be empty")
	}
	if t.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if t.RepositoryID == "" {
		return fmt.Errorf("repositoryID cannot be empty")
	}
	if !t.Type.IsValid() {
		return fmt.Errorf("type must be 'common' or 'increase'")
	}
	if t.Size < 0 {
		return fmt.Errorf("size cannot be negative")
	}
	return nil
}

// MarshalJSON marshals Tag to JSON with RFC3339 timestamps.
// It customizes JSON marshaling to format timestamp fields as RFC3339 strings.
func (t *Tag) MarshalJSON() ([]byte, error) {
	type Alias Tag
	return json.Marshal(&struct {
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		*Alias
	}{
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
		UpdatedAt: t.UpdatedAt.Format(time.RFC3339),
		Alias:     (*Alias)(t),
	})
}

// UnmarshalJSON unmarshals JSON to Tag with RFC3339 timestamp parsing.
// It customizes JSON unmarshaling to parse RFC3339 timestamp strings into time.Time values.
func (t *Tag) UnmarshalJSON(data []byte) error {
	type Alias Tag
	aux := &struct {
		CreatedAt string `json:"createdAt"`
		UpdatedAt string `json:"updatedAt"`
		*Alias
	}{
		Alias: (*Alias)(t),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	var err error
	if aux.CreatedAt != "" {
		t.CreatedAt, err = time.Parse(time.RFC3339, aux.CreatedAt)
		if err != nil {
			return fmt.Errorf("failed to parse createdAt: %w", err)
		}
	}
	if aux.UpdatedAt != "" {
		t.UpdatedAt, err = time.Parse(time.RFC3339, aux.UpdatedAt)
		if err != nil {
			return fmt.Errorf("failed to parse updatedAt: %w", err)
		}
	}
	return nil
}

// DiskFormat represents the disk image format for virtual machine images.
// Valid values include ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, and iso.
type DiskFormat string

const (
	DiskFormatAMI   DiskFormat = "ami"
	DiskFormatARI   DiskFormat = "ari"
	DiskFormatAKI   DiskFormat = "aki"
	DiskFormatVHD   DiskFormat = "vhd"
	DiskFormatVMDK  DiskFormat = "vmdk"
	DiskFormatRaw   DiskFormat = "raw"
	DiskFormatQcow2 DiskFormat = "qcow2"
	DiskFormatVDI   DiskFormat = "vdi"
	DiskFormatISO   DiskFormat = "iso"
)

// IsValid checks if the DiskFormat is a valid value.
func (f DiskFormat) IsValid() bool {
	switch f {
	case DiskFormatAMI, DiskFormatARI, DiskFormatAKI, DiskFormatVHD, DiskFormatVMDK,
		DiskFormatRaw, DiskFormatQcow2, DiskFormatVDI, DiskFormatISO:
		return true
	}
	return false
}

// String returns the string representation of DiskFormat.
func (f DiskFormat) String() string {
	return string(f)
}

// ContainerFormat represents the container image format for virtual machine images.
// Valid values include ami, ari, aki, bare, and ovf.
type ContainerFormat string

const (
	ContainerFormatAMI  ContainerFormat = "ami"
	ContainerFormatARI  ContainerFormat = "ari"
	ContainerFormatAKI  ContainerFormat = "aki"
	ContainerFormatBare ContainerFormat = "bare"
	ContainerFormatOVF  ContainerFormat = "ovf"
)

// IsValid checks if the ContainerFormat is a valid value.
func (f ContainerFormat) IsValid() bool {
	switch f {
	case ContainerFormatAMI, ContainerFormatARI, ContainerFormatAKI, ContainerFormatBare, ContainerFormatOVF:
		return true
	}
	return false
}

// String returns the string representation of ContainerFormat.
func (f ContainerFormat) String() string {
	return string(f)
}
