package tags

import (
	"fmt"
	"strings"

	"github.com/Zillaforge/cloud-sdk/models/vrm/common"
)

// Tag is an alias for common.Tag
type Tag = common.Tag

// Repository is an alias for common.Repository
type Repository = common.Repository

// CreateTagRequest represents a request to create a new tag.
// It includes the tag name, type, and format specifications (disk and container).
// Disk format must be one of: ami, ari, aki, vhd, vmdk, raw, qcow2, vdi, iso.
// Container format must be one of: ami, ari, aki, bare, ovf.
type CreateTagRequest struct {
	Name            string `json:"name"`
	Type            string `json:"type"`
	DiskFormat      string `json:"diskFormat"`
	ContainerFormat string `json:"containerFormat"`
}

// Validate validates the CreateTagRequest.
// It ensures required fields are present and validates that disk and container
// formats are within the allowed enumerated values.
func (r *CreateTagRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if r.Name == "" {
		return fmt.Errorf("name cannot be empty")
	}
	if r.Type == "" {
		return fmt.Errorf("type cannot be empty")
	}
	if r.DiskFormat == "" {
		return fmt.Errorf("diskFormat cannot be empty")
	}
	// Validate disk format enum
	validDiskFormats := map[string]bool{
		"ami":   true,
		"ari":   true,
		"aki":   true,
		"vhd":   true,
		"vmdk":  true,
		"raw":   true,
		"qcow2": true,
		"vdi":   true,
		"iso":   true,
	}
	if !validDiskFormats[r.DiskFormat] {
		return fmt.Errorf("invalid diskFormat: %s", r.DiskFormat)
	}
	if r.ContainerFormat == "" {
		return fmt.Errorf("containerFormat cannot be empty")
	}
	// Validate container format enum
	validContainerFormats := map[string]bool{
		"ami":  true,
		"ari":  true,
		"aki":  true,
		"bare": true,
		"ovf":  true,
	}
	if !validContainerFormats[r.ContainerFormat] {
		return fmt.Errorf("invalid containerFormat: %s", r.ContainerFormat)
	}
	return nil
}

// UpdateTagRequest represents a request to update an existing tag.
// Currently supports updating the tag name.
type UpdateTagRequest struct {
	Name string `json:"name,omitempty"`
}

// Validate validates the UpdateTagRequest.
// Currently performs basic nil checks but may be extended for future validation rules.
func (r *UpdateTagRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("request cannot be nil")
	}
	return nil
}

// DownloadTagRequest represents a request to download an image for a tag into cloud storage.
// The filepath must point to a valid destination such as dss-public://bucket/path.
type DownloadTagRequest struct {
	Filepath string `json:"filepath"`
}

// Validate ensures the download request contains the required filepath destination.
func (r *DownloadTagRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("request cannot be nil")
	}
	if strings.TrimSpace(r.Filepath) == "" {
		return fmt.Errorf("filepath cannot be empty")
	}
	return nil
}

// ListTagsOptions represents options for listing tags.
// It supports pagination with limit and offset, filtering with where conditions,
// and namespace specification for multi-tenant operations.
type ListTagsOptions struct {
	Limit     int
	Offset    int
	Where     []string
	Namespace string
}

// Validate validates the ListTagsOptions.
// It ensures that limit and offset values are within acceptable ranges,
// where limit can be -1 (all) or a positive integer, and offset is non-negative.
func (o *ListTagsOptions) Validate() error {
	if o == nil {
		return nil
	}
	if o.Limit < -1 {
		return fmt.Errorf("limit must be >= -1 (where -1 means all)")
	}
	if o.Offset < 0 {
		return fmt.Errorf("offset cannot be negative")
	}
	return nil
}

// ListTagsResponse represents the response from listing tags.
// It contains the array of tags and the total count for pagination.
type ListTagsResponse struct {
	Tags  []*common.Tag `json:"tags"`
	Total int           `json:"total"`
}
