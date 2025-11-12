package tags

import (
	"encoding/json"
	"fmt"
	"time"
)

// Repository represents a virtual image repository (shallow copy from repositories package).
// This is a lightweight version used in tag responses to avoid circular dependencies.
// It contains basic repository metadata without the full tag list.
type Repository struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Namespace       string    `json:"namespace"`
	OperatingSystem string    `json:"operatingSystem"`
	Description     string    `json:"description,omitempty"`
	Count           int       `json:"count"`
	Creator         *IDName   `json:"creator,omitempty"`
	Project         *IDName   `json:"project,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// IDName represents a lightweight reference to an entity.
// It provides a minimal representation of entities like users or projects,
// containing ID and optional name/display information.
type IDName struct {
	ID          string `json:"id"`
	Name        string `json:"name,omitempty"`
	Account     string `json:"account,omitempty"`
	DisplayName string `json:"displayName,omitempty"`
}

// Tag represents a specific version/tag of an image within a repository.
// Each tag contains metadata about a specific image version, including
// its type, size, status, and any extra properties. Tags are associated
// with a repository and track creation/update timestamps.
type Tag struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	RepositoryID string                 `json:"repositoryID"`
	Type         string                 `json:"type"`
	Size         int64                  `json:"size"`
	Status       string                 `json:"status,omitempty"`
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
	if t.Type == "" {
		return fmt.Errorf("type cannot be empty")
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
	Tags  []*Tag `json:"tags"`
	Total int    `json:"total"`
}
