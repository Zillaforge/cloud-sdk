package flavors

import (
	"errors"
	"time"
)

// Flavor represents a compute instance flavor/size.
// Matches pb.FlavorInfo from vps.yaml swagger specification.
//
// BREAKING CHANGES (v2.0.0):
// - VCPUs field renamed to VCPU (matches API specification)
// - RAM field renamed to Memory (matches API specification)
// - Added GPU field for GPU-enabled flavors
// - Added ProjectIDs, AZ, and timestamp fields
//
// Migration notes:
// - Replace .VCPUs with .VCPU
// - Replace .RAM with .Memory
type Flavor struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	VCPU        int        `json:"vcpu"`          // Number of virtual CPUs (renamed from VCPUs)
	Memory      int        `json:"memory"`        // RAM in MiB (renamed from RAM)
	Disk        int        `json:"disk"`          // Disk size in GiB
	GPU         *GPUInfo   `json:"gpu,omitempty"` // GPU configuration (optional)
	Public      bool       `json:"public"`
	Tags        []string   `json:"tags,omitempty"`
	ProjectIDs  []string   `json:"project_ids,omitempty"` // Restricted project IDs
	AZ          string     `json:"az,omitempty"`          // Availability zone
	CreatedAt   *time.Time `json:"createdAt,omitempty"`   // Creation timestamp (ISO 8601)
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`   // Last update timestamp (ISO 8601)
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`   // Deletion timestamp (ISO 8601)
}

// GPUInfo represents GPU configuration for flavors.
type GPUInfo struct {
	Count  int    `json:"count"`   // Number of GPUs
	IsVGPU bool   `json:"is_vgpu"` // Whether vGPU is used
	Model  string `json:"model"`   // GPU model name
}

// Validate checks if the Flavor has valid required fields and data types.
func (f *Flavor) Validate() error {
	if f.ID == "" {
		return errors.New("id is required")
	}
	if f.Name == "" {
		return errors.New("name is required")
	}
	if f.VCPU < 0 {
		return errors.New("vcpu must be >= 0")
	}
	if f.Memory < 0 {
		return errors.New("memory must be >= 0")
	}
	if f.Disk < 0 {
		return errors.New("disk must be >= 0")
	}
	if f.GPU != nil {
		if f.GPU.Count < 0 {
			return errors.New("gpu count must be >= 0")
		}
		if f.GPU.Model == "" {
			return errors.New("gpu model is required when GPU is present")
		}
	}
	return nil
}

// ListFlavorsOptions provides filtering options for listing flavors.
type ListFlavorsOptions struct {
	Name           string   // Filter by name
	Public         *bool    // nil = all, true = public only, false = private only
	Tags           []string // Filter by tags (multiple allowed, sent as separate query params)
	ResizeServerID string   // Filter flavors available for server resize
}

// Validate checks if the ListFlavorsOptions has valid values.
func (o *ListFlavorsOptions) Validate() error {
	if o == nil {
		return nil
	}
	for _, tag := range o.Tags {
		if tag == "" {
			return errors.New("tags must not contain empty strings")
		}
	}
	return nil
}

// FlavorListResponse represents the response from listing flavors.
type FlavorListResponse struct {
	Flavors []*Flavor `json:"flavors"` // Changed from "items" to "flavors" to match API contract
}
