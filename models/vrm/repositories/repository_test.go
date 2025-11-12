package repositories

import (
	"testing"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vrm/common"
)

// T030: Unit test for Repository struct validation
func TestRepositoryValidate(t *testing.T) {
	tests := []struct {
		name    string
		repo    *Repository
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid repository",
			repo: &Repository{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Name:            "ubuntu",
				Namespace:       "public",
				OperatingSystem: "linux",
				Description:     "Ubuntu base images",
				Count:           5,
				Creator:         &common.IDName{ID: "user-1", Name: "admin"},
				Project:         &common.IDName{ID: "proj-1", Name: "test"},
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: false,
		},
		{
			name: "valid repository without description",
			repo: &Repository{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Name:            "centos",
				Namespace:       "private",
				OperatingSystem: "linux",
				Count:           0,
				Creator:         &common.IDName{ID: "user-1", Name: "admin"},
				Project:         &common.IDName{ID: "proj-1", Name: "test"},
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: false,
		},
		{
			name:    "nil repository",
			repo:    nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
		{
			name: "missing ID",
			repo: &Repository{
				Name:            "ubuntu",
				Namespace:       "public",
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "ID is required",
		},
		{
			name: "missing Name",
			repo: &Repository{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Namespace:       "public",
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "invalid Namespace",
			repo: &Repository{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Name:            "ubuntu",
				Namespace:       "internal",
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "namespace must be",
		},
		{
			name: "invalid OperatingSystem",
			repo: &Repository{
				ID:              "550e8400-e29b-41d4-a716-446655440000",
				Name:            "ubuntu",
				Namespace:       "public",
				OperatingSystem: "macos",
			},
			wantErr: true,
			errMsg:  "operating system must be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// T031: Unit test for Repository JSON marshaling
func TestRepositoryJSONMarshaling(t *testing.T) {
	repo := &Repository{
		ID:              "550e8400-e29b-41d4-a716-446655440000",
		Name:            "ubuntu",
		Namespace:       "public",
		OperatingSystem: "linux",
		Description:     "Ubuntu base images",
		Count:           3,
		Creator: &common.IDName{
			ID:   "user-123",
			Name: "admin",
		},
		Project: &common.IDName{
			ID:   "proj-456",
			Name: "default",
		},
		CreatedAt: time.Date(2024, 8, 19, 8, 32, 15, 0, time.UTC),
		UpdatedAt: time.Date(2024, 8, 19, 8, 32, 15, 0, time.UTC),
	}

	// Test that JSON field names use camelCase
	data, err := repo.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Verify camelCase JSON tags
	jsonStr := string(data)
	expectedFields := []string{
		`"id"`,
		`"name"`,
		`"namespace"`,
		`"operatingSystem"`,
		`"description"`,
		`"count"`,
		`"creator"`,
		`"project"`,
		`"createdAt"`,
		`"updatedAt"`,
	}

	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON marshaling missing field %v", field)
		}
	}
}

// T032: Unit test for Repository JSON unmarshaling
func TestRepositoryJSONUnmarshaling(t *testing.T) {
	jsonData := []byte(`{
		"id": "550e8400-e29b-41d4-a716-446655440000",
		"name": "ubuntu",
		"namespace": "public",
		"operatingSystem": "linux",
		"description": "Ubuntu base images",
		"count": 3,
		"creator": {
			"id": "user-123",
			"name": "admin"
		},
		"project": {
			"id": "proj-456",
			"name": "default"
		},
		"createdAt": "2024-08-19T08:32:15Z",
		"updatedAt": "2024-08-19T08:32:15Z"
	}`)

	var repo Repository
	err := repo.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	if repo.ID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("ID = %v, want 550e8400-e29b-41d4-a716-446655440000", repo.ID)
	}
	if repo.Name != "ubuntu" {
		t.Errorf("Name = %v, want ubuntu", repo.Name)
	}
	if repo.OperatingSystem != "linux" {
		t.Errorf("OperatingSystem = %v, want linux", repo.OperatingSystem)
	}
	if repo.Count != 3 {
		t.Errorf("Count = %v, want 3", repo.Count)
	}
}

// T033: Unit test for CreateRepositoryRequest validation
func TestCreateRepositoryRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateRepositoryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid create request",
			req: &CreateRepositoryRequest{
				Name:            "ubuntu",
				OperatingSystem: "linux",
				Description:     "Ubuntu base images",
			},
			wantErr: false,
		},
		{
			name: "valid without description",
			req: &CreateRepositoryRequest{
				Name:            "centos",
				OperatingSystem: "windows",
			},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
		{
			name: "missing Name",
			req: &CreateRepositoryRequest{
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing OperatingSystem",
			req: &CreateRepositoryRequest{
				Name: "ubuntu",
			},
			wantErr: true,
			errMsg:  "operatingSystem is required",
		},
		{
			name: "invalid OperatingSystem",
			req: &CreateRepositoryRequest{
				Name:            "ubuntu",
				OperatingSystem: "macos",
			},
			wantErr: true,
			errMsg:  "operatingSystem must be",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// T034: Unit test for UpdateRepositoryRequest validation
func TestUpdateRepositoryRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *UpdateRepositoryRequest
		wantErr bool
	}{
		{
			name: "valid update with description",
			req: &UpdateRepositoryRequest{
				Description: "Updated description",
			},
			wantErr: false,
		},
		{
			name:    "valid update with empty fields",
			req:     &UpdateRepositoryRequest{},
			wantErr: false,
		},
		{
			name:    "nil request",
			req:     nil,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// T035: Unit test for ListRepositoriesOptions
func TestListRepositoriesOptionsValidate(t *testing.T) {
	tests := []struct {
		name    string
		opts    *ListRepositoriesOptions
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid options with limit",
			opts: &ListRepositoriesOptions{
				Limit: 50,
			},
			wantErr: false,
		},
		{
			name: "valid options with offset",
			opts: &ListRepositoriesOptions{
				Limit:  50,
				Offset: 100,
			},
			wantErr: false,
		},
		{
			name: "valid options with where",
			opts: &ListRepositoriesOptions{
				Where: []string{"namespace=public", "operatingSystem=linux"},
			},
			wantErr: false,
		},
		{
			name: "valid options with namespace",
			opts: &ListRepositoriesOptions{
				Namespace: "public",
			},
			wantErr: false,
		},
		{
			name:    "valid options empty",
			opts:    &ListRepositoriesOptions{},
			wantErr: false,
		},
		{
			name: "invalid negative limit",
			opts: &ListRepositoriesOptions{
				Limit: -10,
			},
			wantErr: true,
			errMsg:  "limit must be",
		},
		{
			name: "invalid negative offset",
			opts: &ListRepositoriesOptions{
				Offset: -1,
			},
			wantErr: true,
			errMsg:  "offset must be",
		},
		{
			name:    "nil options",
			opts:    nil,
			wantErr: true,
			errMsg:  "cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.opts.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" {
				if !contains(err.Error(), tt.errMsg) {
					t.Errorf("Validate() error message = %v, want to contain %v", err.Error(), tt.errMsg)
				}
			}
		})
	}
}

// Helper function
func contains(s, substr string) bool {
	for i := 0; i < len(s)-len(substr)+1; i++ {
		match := true
		for j := 0; j < len(substr); j++ {
			if s[i+j] != substr[j] {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
