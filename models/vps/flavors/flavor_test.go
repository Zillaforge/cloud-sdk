package flavors

import (
	"encoding/json"
	"testing"
	"time"
)

// TestFlavorJSONMarshaling tests the basic JSON marshaling/unmarshaling for Flavor struct.
func TestFlavorJSONMarshaling(t *testing.T) {
	// Create a sample timestamp
	now := time.Now().UTC()

	tests := []struct {
		name    string
		flavor  *Flavor
		wantErr bool
	}{
		{
			name: "complete flavor with all fields",
			flavor: &Flavor{
				ID:          "flavor-123",
				Name:        "m1.large",
				Description: "Large flavor with GPU",
				VCPU:        8,
				Memory:      16384,
				Disk:        100,
				GPU: &GPUInfo{
					Count:  2,
					IsVGPU: false,
					Model:  "NVIDIA Tesla V100",
				},
				Public:     true,
				Tags:       []string{"gpu", "large"},
				ProjectIDs: []string{"proj-1", "proj-2"},
				AZ:         "az-1",
				CreatedAt:  &now,
				UpdatedAt:  &now,
			},
			wantErr: false,
		},
		{
			name: "minimal flavor with required fields only",
			flavor: &Flavor{
				ID:     "flavor-456",
				Name:   "m1.small",
				VCPU:   2,
				Memory: 2048,
				Disk:   20,
				Public: true,
			},
			wantErr: false,
		},
		{
			name: "flavor without GPU",
			flavor: &Flavor{
				ID:     "flavor-789",
				Name:   "m1.medium",
				VCPU:   4,
				Memory: 8192,
				Disk:   50,
				Public: false,
				Tags:   []string{"standard"},
			},
			wantErr: false,
		},
		{
			name: "flavor with only timestamps",
			flavor: &Flavor{
				ID:        "flavor-timestamp",
				Name:      "m1.timestamp",
				VCPU:      2,
				Memory:    4096,
				Disk:      20,
				Public:    true,
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: &now,
			},
			wantErr: false,
		},
		{
			name: "flavor with empty slices",
			flavor: &Flavor{
				ID:         "flavor-empty",
				Name:       "m1.empty",
				VCPU:       2,
				Memory:     4096,
				Disk:       20,
				Public:     true,
				Tags:       []string{},
				ProjectIDs: []string{},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.flavor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Unmarshal back to struct
			var unmarshaled Flavor
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}

			// Compare fields
			if unmarshaled.ID != tt.flavor.ID {
				t.Errorf("ID mismatch: got %v, want %v", unmarshaled.ID, tt.flavor.ID)
			}
			if unmarshaled.Name != tt.flavor.Name {
				t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, tt.flavor.Name)
			}
			if unmarshaled.VCPU != tt.flavor.VCPU {
				t.Errorf("VCPU mismatch: got %v, want %v", unmarshaled.VCPU, tt.flavor.VCPU)
			}
			if unmarshaled.Memory != tt.flavor.Memory {
				t.Errorf("Memory mismatch: got %v, want %v", unmarshaled.Memory, tt.flavor.Memory)
			}
			if unmarshaled.Disk != tt.flavor.Disk {
				t.Errorf("Disk mismatch: got %v, want %v", unmarshaled.Disk, tt.flavor.Disk)
			}
			if unmarshaled.Public != tt.flavor.Public {
				t.Errorf("Public mismatch: got %v, want %v", unmarshaled.Public, tt.flavor.Public)
			}
		})
	}
}

// TestFlavorJSONTags verifies JSON tags match API specification.
func TestFlavorJSONTags(t *testing.T) {
	flavor := &Flavor{
		ID:          "test-id",
		Name:        "test-name",
		Description: "test-desc",
		VCPU:        4,
		Memory:      8192,
		Disk:        50,
		GPU: &GPUInfo{
			Count:  1,
			IsVGPU: true,
			Model:  "Test GPU",
		},
		Public:     true,
		Tags:       []string{"tag1", "tag2"},
		ProjectIDs: []string{"proj1"},
		AZ:         "az-1",
	}

	data, err := json.Marshal(flavor)
	if err != nil {
		t.Fatalf("Failed to marshal flavor: %v", err)
	}

	// Verify expected JSON keys are present
	expectedKeys := []string{
		"\"id\":",
		"\"name\":",
		"\"description\":",
		"\"vcpu\":",
		"\"memory\":",
		"\"disk\":",
		"\"gpu\":",
		"\"public\":",
		"\"tags\":",
		"\"project_ids\":",
		"\"az\":",
	}

	jsonStr := string(data)
	for _, key := range expectedKeys {
		if !contains(jsonStr, key) {
			t.Errorf("Expected JSON to contain key %s, but it doesn't. JSON: %s", key, jsonStr)
		}
	}
}

// TestGPUInfoJSONMarshaling tests GPU info JSON marshaling.
func TestGPUInfoJSONMarshaling(t *testing.T) {
	tests := []struct {
		name    string
		gpu     *GPUInfo
		wantErr bool
	}{
		{
			name: "complete GPU info",
			gpu: &GPUInfo{
				Count:  4,
				IsVGPU: true,
				Model:  "NVIDIA A100",
			},
			wantErr: false,
		},
		{
			name: "GPU without vGPU",
			gpu: &GPUInfo{
				Count:  1,
				IsVGPU: false,
				Model:  "NVIDIA Tesla T4",
			},
			wantErr: false,
		},
		{
			name: "GPU with zero count",
			gpu: &GPUInfo{
				Count:  0,
				IsVGPU: false,
				Model:  "NVIDIA Tesla T4",
			},
			wantErr: false,
		},
		{
			name: "GPU with large count",
			gpu: &GPUInfo{
				Count:  16,
				IsVGPU: true,
				Model:  "NVIDIA A100",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.gpu)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			var unmarshaled GPUInfo
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}

			if unmarshaled.Count != tt.gpu.Count {
				t.Errorf("Count mismatch: got %v, want %v", unmarshaled.Count, tt.gpu.Count)
			}
			if unmarshaled.IsVGPU != tt.gpu.IsVGPU {
				t.Errorf("IsVGPU mismatch: got %v, want %v", unmarshaled.IsVGPU, tt.gpu.IsVGPU)
			}
			if unmarshaled.Model != tt.gpu.Model {
				t.Errorf("Model mismatch: got %v, want %v", unmarshaled.Model, tt.gpu.Model)
			}
		})
	}
}

// TestTimestampMarshaling tests timestamp field handling.
func TestTimestampMarshaling(t *testing.T) {
	now := time.Date(2025, 11, 8, 12, 0, 0, 0, time.UTC)
	past := time.Date(2025, 10, 1, 10, 0, 0, 0, time.UTC)
	future := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)

	tests := []struct {
		name    string
		flavor  *Flavor
		wantErr bool
	}{
		{
			name: "all timestamps present",
			flavor: &Flavor{
				ID:        "test-id",
				Name:      "test-name",
				VCPU:      2,
				Memory:    4096,
				Disk:      20,
				Public:    true,
				CreatedAt: &now,
				UpdatedAt: &now,
				DeletedAt: &now,
			},
			wantErr: false,
		},
		{
			name: "only created and updated timestamps",
			flavor: &Flavor{
				ID:        "test-id",
				Name:      "test-name",
				VCPU:      2,
				Memory:    4096,
				Disk:      20,
				Public:    true,
				CreatedAt: &past,
				UpdatedAt: &now,
			},
			wantErr: false,
		},
		{
			name: "only deleted timestamp",
			flavor: &Flavor{
				ID:        "test-id",
				Name:      "test-name",
				VCPU:      2,
				Memory:    4096,
				Disk:      20,
				Public:    true,
				DeletedAt: &future,
			},
			wantErr: false,
		},
		{
			name: "no timestamps",
			flavor: &Flavor{
				ID:     "test-id",
				Name:   "test-name",
				VCPU:   2,
				Memory: 4096,
				Disk:   20,
				Public: true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.flavor)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			var unmarshaled Flavor
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}

			// Check timestamps are preserved
			if tt.flavor.CreatedAt != nil {
				if unmarshaled.CreatedAt == nil {
					t.Error("CreatedAt should not be nil after unmarshal")
				} else if !unmarshaled.CreatedAt.Equal(*tt.flavor.CreatedAt) {
					t.Errorf("CreatedAt mismatch: got %v, want %v", unmarshaled.CreatedAt, tt.flavor.CreatedAt)
				}
			}

			if tt.flavor.UpdatedAt != nil {
				if unmarshaled.UpdatedAt == nil {
					t.Error("UpdatedAt should not be nil after unmarshal")
				} else if !unmarshaled.UpdatedAt.Equal(*tt.flavor.UpdatedAt) {
					t.Errorf("UpdatedAt mismatch: got %v, want %v", unmarshaled.UpdatedAt, tt.flavor.UpdatedAt)
				}
			}

			if tt.flavor.DeletedAt != nil {
				if unmarshaled.DeletedAt == nil {
					t.Error("DeletedAt should not be nil after unmarshal")
				} else if !unmarshaled.DeletedAt.Equal(*tt.flavor.DeletedAt) {
					t.Errorf("DeletedAt mismatch: got %v, want %v", unmarshaled.DeletedAt, tt.flavor.DeletedAt)
				}
			}
		})
	}
}

// TestFlavorValidation tests field validation rules.
func TestFlavorValidation(t *testing.T) {
	tests := []struct {
		name    string
		flavor  *Flavor
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid flavor",
			flavor: &Flavor{
				ID:     "test-id",
				Name:   "test-name",
				VCPU:   2,
				Memory: 4096,
				Disk:   20,
				Public: true,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			flavor: &Flavor{
				Name:   "test-name",
				VCPU:   2,
				Memory: 4096,
				Disk:   20,
				Public: true,
			},
			wantErr: true,
			errMsg:  "id is required",
		},
		{
			name: "missing name",
			flavor: &Flavor{
				ID:     "test-id",
				VCPU:   2,
				Memory: 4096,
				Disk:   20,
				Public: true,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "negative VCPU",
			flavor: &Flavor{
				ID:     "test-id",
				Name:   "test-name",
				VCPU:   -1,
				Memory: 4096,
				Disk:   20,
				Public: true,
			},
			wantErr: true,
			errMsg:  "vcpu must be >= 0",
		},
		{
			name: "negative memory",
			flavor: &Flavor{
				ID:     "test-id",
				Name:   "test-name",
				VCPU:   2,
				Memory: -100,
				Disk:   20,
				Public: true,
			},
			wantErr: true,
			errMsg:  "memory must be >= 0",
		},
		{
			name: "negative disk",
			flavor: &Flavor{
				ID:     "test-id",
				Name:   "test-name",
				VCPU:   2,
				Memory: 4096,
				Disk:   -10,
				Public: true,
			},
			wantErr: true,
			errMsg:  "disk must be >= 0",
		},
		{
			name: "GPU with negative count",
			flavor: &Flavor{
				ID:     "test-id",
				Name:   "test-name",
				VCPU:   2,
				Memory: 4096,
				Disk:   20,
				Public: true,
				GPU: &GPUInfo{
					Count:  -1,
					IsVGPU: false,
					Model:  "NVIDIA Tesla T4",
				},
			},
			wantErr: true,
			errMsg:  "gpu count must be >= 0",
		},
		{
			name: "GPU with empty model",
			flavor: &Flavor{
				ID:     "test-id",
				Name:   "test-name",
				VCPU:   2,
				Memory: 4096,
				Disk:   20,
				Public: true,
				GPU: &GPUInfo{
					Count:  1,
					IsVGPU: false,
					Model:  "",
				},
			},
			wantErr: true,
			errMsg:  "gpu model is required when GPU is present",
		},
		{
			name: "valid flavor with GPU",
			flavor: &Flavor{
				ID:     "test-id",
				Name:   "test-name",
				VCPU:   2,
				Memory: 4096,
				Disk:   20,
				Public: true,
				GPU: &GPUInfo{
					Count:  2,
					IsVGPU: true,
					Model:  "NVIDIA A100",
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.flavor.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// TestListFlavorsOptionsValidation tests the ListFlavorsOptions validation.
func TestListFlavorsOptionsValidation(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ListFlavorsOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid options with all fields",
			opts: &ListFlavorsOptions{
				Name:           "m1.large",
				Public:         boolPtr(true),
				Tags:           []string{"gpu", "large"},
				ResizeServerID: "server-123",
			},
			wantErr: false,
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: false,
		},
		{
			name: "valid options with tags only",
			opts: &ListFlavorsOptions{
				Tags: []string{"standard", "compute"},
			},
			wantErr: false,
		},
		{
			name: "empty tag in tags array",
			opts: &ListFlavorsOptions{
				Tags: []string{"valid", "", "another"},
			},
			wantErr: true,
			errMsg:  "tags must not contain empty strings",
		},
		{
			name: "multiple empty tags",
			opts: &ListFlavorsOptions{
				Tags: []string{"", "", ""},
			},
			wantErr: true,
			errMsg:  "tags must not contain empty strings",
		},
		{
			name: "single empty tag",
			opts: &ListFlavorsOptions{
				Tags: []string{""},
			},
			wantErr: true,
			errMsg:  "tags must not contain empty strings",
		},
		{
			name: "empty string in middle of tags",
			opts: &ListFlavorsOptions{
				Tags: []string{"first", "", "third"},
			},
			wantErr: true,
			errMsg:  "tags must not contain empty strings",
		},
		{
			name: "public false filter",
			opts: &ListFlavorsOptions{
				Public: boolPtr(false),
			},
			wantErr: false,
		},
		{
			name: "resize server ID only",
			opts: &ListFlavorsOptions{
				ResizeServerID: "server-456",
			},
			wantErr: false,
		},
		{
			name: "nil tags slice",
			opts: &ListFlavorsOptions{
				Tags: nil,
			},
			wantErr: false,
		},
		{
			name: "empty tags slice",
			opts: &ListFlavorsOptions{
				Tags: []string{},
			},
			wantErr: false,
		},
		{
			name: "multiple valid tags",
			opts: &ListFlavorsOptions{
				Tags: []string{"tag1", "tag2", "tag3", "tag4", "tag5"},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.opts != nil {
				err = tt.opts.Validate()
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr && err != nil {
				if err.Error() != tt.errMsg {
					t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// Helper functions

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsMiddle(s, substr)))
}

func containsMiddle(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func boolPtr(b bool) *bool {
	return &b
}

// TestFlavorListResponseJSONMarshaling tests FlavorListResponse JSON marshaling.
func TestFlavorListResponseJSONMarshaling(t *testing.T) {
	now := time.Now().UTC()

	tests := []struct {
		name     string
		response *FlavorListResponse
		wantErr  bool
	}{
		{
			name: "response with multiple flavors",
			response: &FlavorListResponse{
				Flavors: []*Flavor{
					{
						ID:     "flavor-1",
						Name:   "m1.small",
						VCPU:   1,
						Memory: 1024,
						Disk:   10,
						Public: true,
					},
					{
						ID:          "flavor-2",
						Name:        "m1.large",
						VCPU:        8,
						Memory:      16384,
						Disk:        100,
						Public:      true,
						Description: "Large flavor",
						GPU: &GPUInfo{
							Count:  2,
							IsVGPU: false,
							Model:  "NVIDIA Tesla V100",
						},
						Tags:       []string{"gpu", "large"},
						ProjectIDs: []string{"proj-1"},
						AZ:         "az-1",
						CreatedAt:  &now,
						UpdatedAt:  &now,
					},
				},
			},
			wantErr: false,
		},
		{
			name: "empty response",
			response: &FlavorListResponse{
				Flavors: []*Flavor{},
			},
			wantErr: false,
		},
		{
			name: "nil flavors slice",
			response: &FlavorListResponse{
				Flavors: nil,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("Marshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.wantErr {
				return
			}

			// Unmarshal back to struct
			var unmarshaled FlavorListResponse
			err = json.Unmarshal(data, &unmarshaled)
			if err != nil {
				t.Errorf("Unmarshal() error = %v", err)
				return
			}

			// Compare lengths
			if len(unmarshaled.Flavors) != len(tt.response.Flavors) {
				t.Errorf("Flavors length mismatch: got %d, want %d", len(unmarshaled.Flavors), len(tt.response.Flavors))
				return
			}

			// Compare each flavor
			for i, flavor := range tt.response.Flavors {
				if unmarshaled.Flavors[i].ID != flavor.ID {
					t.Errorf("Flavor[%d].ID mismatch: got %v, want %v", i, unmarshaled.Flavors[i].ID, flavor.ID)
				}
				if unmarshaled.Flavors[i].Name != flavor.Name {
					t.Errorf("Flavor[%d].Name mismatch: got %v, want %v", i, unmarshaled.Flavors[i].Name, flavor.Name)
				}
			}
		})
	}
}

// TestFlavorListResponseJSONTags verifies FlavorListResponse JSON tags.
func TestFlavorListResponseJSONTags(t *testing.T) {
	response := &FlavorListResponse{
		Flavors: []*Flavor{
			{
				ID:   "test-flavor",
				Name: "test-name",
				VCPU: 2,
			},
		},
	}

	data, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	jsonStr := string(data)
	if !contains(jsonStr, "\"flavors\":") {
		t.Errorf("Expected JSON to contain key \"flavors\":, but it doesn't. JSON: %s", jsonStr)
	}
}
