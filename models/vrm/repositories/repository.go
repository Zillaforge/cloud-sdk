// Package repositories provides data models for VRM repository operations.
package repositories

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vrm/common"
)

// Repository represents a virtual image repository within a project.
// It contains metadata about the repository including its name, namespace,
// operating system, description, and associated tags. The repository
// tracks creation and update timestamps, as well as creator and project information.
type Repository struct {
	ID              string         `json:"id"`
	Name            string         `json:"name"`
	Namespace       string         `json:"namespace"`
	OperatingSystem string         `json:"operatingSystem"`
	Description     string         `json:"description,omitempty"`
	Tags            []*Tag         `json:"tags,omitempty"`
	Count           int            `json:"count"`
	Creator         *common.IDName `json:"creator"`
	Project         *common.IDName `json:"project"`
	CreatedAt       time.Time      `json:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt"`
}

// Tag represents a tagged version of an image within a repository.
// Each tag contains metadata about a specific version of an image,
// including its size, status, and any extra properties. Tags are
// associated with a specific repository and track creation/update timestamps.
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

// CreateRepositoryRequest represents a request to create a new repository.
// It includes the required name and operating system, with an optional description.
// The operating system must be either "linux" or "windows".
type CreateRepositoryRequest struct {
	Name            string `json:"name"`
	OperatingSystem string `json:"operatingSystem"`
	Description     string `json:"description,omitempty"`
}

// Validate validates the CreateRepositoryRequest.
// It ensures that name and operating system are provided and valid,
// with operating system restricted to "linux" or "windows".
func (r *CreateRepositoryRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("createRepositoryRequest cannot be nil")
	}
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("name is required and must not be empty")
	}
	if strings.TrimSpace(r.OperatingSystem) == "" {
		return fmt.Errorf("operatingSystem is required and must not be empty")
	}
	if r.OperatingSystem != "linux" && r.OperatingSystem != "windows" {
		return fmt.Errorf("operatingSystem must be 'linux' or 'windows'")
	}
	return nil
}

// UpdateRepositoryRequest represents a request to update an existing repository.
// Currently only supports updating the repository description.
type UpdateRepositoryRequest struct {
	Description string `json:"description,omitempty"`
}

// Validate validates the UpdateRepositoryRequest.
// Currently performs basic nil checks but may be extended for future validation rules.
func (r *UpdateRepositoryRequest) Validate() error {
	if r == nil {
		return fmt.Errorf("updateRepositoryRequest cannot be nil")
	}
	return nil
}

// ListRepositoriesOptions represents options for listing repositories.
// It supports pagination with limit and offset, filtering with where conditions,
// and namespace specification for multi-tenant operations.
type ListRepositoriesOptions struct {
	Limit     int      // -1 for all, positive integer for limit
	Offset    int      // non-negative integer for offset
	Where     []string // Filter conditions (e.g., "namespace=public")
	Namespace string   // Namespace for X-Namespace header
}

// Validate validates the ListRepositoriesOptions.
// It ensures that limit and offset values are within acceptable ranges.
func (o *ListRepositoriesOptions) Validate() error {
	if o == nil {
		return fmt.Errorf("listRepositoriesOptions cannot be nil")
	}
	if o.Limit < -1 {
		return fmt.Errorf("limit must be >= -1")
	}
	if o.Offset < 0 {
		return fmt.Errorf("offset must be >= 0")
	}
	return nil
}
