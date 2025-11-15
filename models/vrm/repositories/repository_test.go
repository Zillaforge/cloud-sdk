package repositories

import (
	"encoding/json"
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

// T036: Unit test for CreateSnapshotRequest validation
func TestCreateSnapshotRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateSnapshotRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with new repository",
			req: &CreateSnapshotRequest{
				Version:         "v1",
				Name:            "snapshot-repo",
				OperatingSystem: "linux",
				Description:     "test",
			},
			wantErr: false,
		},
		{
			name: "valid request targeting existing repository",
			req: &CreateSnapshotRequest{
				Version:      "v2",
				RepositoryID: "repo-123",
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
			name: "missing version",
			req: &CreateSnapshotRequest{
				Name:            "snapshot-repo",
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "missing name when repositoryId absent",
			req: &CreateSnapshotRequest{
				Version:         "v1",
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing operating system when repositoryId absent",
			req: &CreateSnapshotRequest{
				Version: "v1",
				Name:    "snapshot-repo",
			},
			wantErr: true,
			errMsg:  "operatingSystem is required",
		},
		{
			name: "invalid operating system",
			req: &CreateSnapshotRequest{
				Version:         "v1",
				Name:            "snapshot-repo",
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

// T037: Unit test for UploadImageRequest validation
func TestUploadImageRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     *UploadImageRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid new repository upload",
			req: &UploadImageRequest{
				Name:            "ubuntu",
				OperatingSystem: "linux",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image",
			},
			wantErr: false,
		},
		{
			name: "valid repositoryId upload",
			req: &UploadImageRequest{
				RepositoryID:    "repo-1",
				Version:         "v2",
				Type:            "common",
				DiskFormat:      "raw",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image",
			},
			wantErr: false,
		},
		{
			name: "valid tagId upload",
			req: &UploadImageRequest{
				TagID:    "tag-1",
				Filepath: "s3://bucket/image",
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
			name: "missing filepath",
			req: &UploadImageRequest{
				Name:            "ubuntu",
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "filepath is required",
		},
		{
			name: "both repositoryId and tagId provided",
			req: &UploadImageRequest{
				RepositoryID: "repo-1",
				TagID:        "tag-1",
				Filepath:     "s3://bucket/image",
			},
			wantErr: true,
			errMsg:  "only one",
		},
		{
			name: "missing name for new repository mode",
			req: &UploadImageRequest{
				OperatingSystem: "linux",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "invalid operating system",
			req: &UploadImageRequest{
				Name:            "ubuntu",
				OperatingSystem: "macos",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image",
			},
			wantErr: true,
			errMsg:  "operatingSystem must be",
		},
		{
			name: "repository mode missing version",
			req: &UploadImageRequest{
				RepositoryID:    "repo-1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image",
			},
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "repository mode invalid disk format",
			req: &UploadImageRequest{
				RepositoryID:    "repo-1",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "invalid",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image",
			},
			wantErr: true,
			errMsg:  "invalid diskFormat",
		},
		{
			name: "base mode invalid container format",
			req: &UploadImageRequest{
				Name:            "ubuntu",
				OperatingSystem: "linux",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "invalid",
				Filepath:        "s3://bucket/image",
			},
			wantErr: true,
			errMsg:  "invalid containerFormat",
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

// T038: Unit test for Tag validation - Tag doesn't have Validate method, so this test is not applicable
// Skipping Tag validation tests as Tag struct doesn't implement Validate method

// T039: Unit test for Tag JSON marshaling - Tag doesn't have custom JSON methods, so this test is not applicable
// Skipping Tag JSON marshaling tests as Tag struct doesn't implement custom JSON methods

// T040: Unit test for Tag JSON unmarshaling - Tag doesn't have custom JSON methods, so this test is not applicable
// Skipping Tag JSON unmarshaling tests as Tag struct doesn't implement custom JSON methods

// T041: Unit test for CreateSnapshotFromNewRepositoryRequest
func TestCreateSnapshotFromNewRepositoryRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateSnapshotFromNewRepositoryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &CreateSnapshotFromNewRepositoryRequest{
				Name:            "snapshot-repo",
				OperatingSystem: "linux",
				Version:         "v1",
				Description:     "test snapshot",
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
			name: "missing name",
			req: &CreateSnapshotFromNewRepositoryRequest{
				OperatingSystem: "linux",
				Version:         "v1",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "invalid operating system",
			req: &CreateSnapshotFromNewRepositoryRequest{
				Name:            "snapshot-repo",
				OperatingSystem: "macos",
				Version:         "v1",
			},
			wantErr: true,
			errMsg:  "operatingSystem must be",
		},
		{
			name: "missing version",
			req: &CreateSnapshotFromNewRepositoryRequest{
				Name:            "snapshot-repo",
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "version is required",
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

	// Test ToCreateSnapshotRequest conversion
	t.Run("ToCreateSnapshotRequest conversion", func(t *testing.T) {
		req := &CreateSnapshotFromNewRepositoryRequest{
			Name:            "test-repo",
			OperatingSystem: "linux",
			Version:         "v1.0",
			Description:     "test description",
		}

		converted := req.ToCreateSnapshotRequest()
		if converted.Name != req.Name {
			t.Errorf("Name = %v, want %v", converted.Name, req.Name)
		}
		if converted.OperatingSystem != req.OperatingSystem {
			t.Errorf("OperatingSystem = %v, want %v", converted.OperatingSystem, req.OperatingSystem)
		}
		if converted.Version != req.Version {
			t.Errorf("Version = %v, want %v", converted.Version, req.Version)
		}
		if converted.Description != req.Description {
			t.Errorf("Description = %v, want %v", converted.Description, req.Description)
		}
	})
}

// T042: Unit test for CreateSnapshotFromExistingRepositoryRequest
func TestCreateSnapshotFromExistingRepositoryRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateSnapshotFromExistingRepositoryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &CreateSnapshotFromExistingRepositoryRequest{
				RepositoryID: "repo-123",
				Version:      "v2",
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
			name: "missing repository ID",
			req: &CreateSnapshotFromExistingRepositoryRequest{
				Version: "v2",
			},
			wantErr: true,
			errMsg:  "repositoryId is required",
		},
		{
			name: "missing version",
			req: &CreateSnapshotFromExistingRepositoryRequest{
				RepositoryID: "repo-123",
			},
			wantErr: true,
			errMsg:  "version is required",
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

	// Test ToCreateSnapshotRequest conversion
	t.Run("ToCreateSnapshotRequest conversion", func(t *testing.T) {
		req := &CreateSnapshotFromExistingRepositoryRequest{
			RepositoryID: "repo-456",
			Version:      "v3.0",
		}

		converted := req.ToCreateSnapshotRequest()
		if converted.RepositoryID != req.RepositoryID {
			t.Errorf("RepositoryID = %v, want %v", converted.RepositoryID, req.RepositoryID)
		}
		if converted.Version != req.Version {
			t.Errorf("Version = %v, want %v", converted.Version, req.Version)
		}
	})
}

// T043: Unit test for UploadToNewRepositoryRequest
func TestUploadToNewRepositoryRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *UploadToNewRepositoryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				OperatingSystem: "linux",
				Description:     "Ubuntu image",
				Filepath:        "s3://bucket/image.qcow2",
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
			name: "missing name",
			req: &UploadToNewRepositoryRequest{
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				OperatingSystem: "linux",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing operating system",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "operatingSystem is required",
		},
		{
			name: "invalid operating system",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				OperatingSystem: "macos",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "operatingSystem must be",
		},
		{
			name: "missing version",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				OperatingSystem: "linux",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "missing type",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				OperatingSystem: "linux",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "type is required",
		},
		{
			name: "missing disk format",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				Type:            "common",
				ContainerFormat: "bare",
				OperatingSystem: "linux",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "diskFormat is required",
		},
		{
			name: "invalid disk format",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "invalid",
				ContainerFormat: "bare",
				OperatingSystem: "linux",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "invalid diskFormat",
		},
		{
			name: "missing container format",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				OperatingSystem: "linux",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "containerFormat is required",
		},
		{
			name: "invalid container format",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "invalid",
				OperatingSystem: "linux",
				Filepath:        "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "invalid containerFormat",
		},
		{
			name: "missing filepath",
			req: &UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				OperatingSystem: "linux",
			},
			wantErr: true,
			errMsg:  "filepath is required",
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

	// Test ToUploadImageRequest conversion
	t.Run("ToUploadImageRequest conversion", func(t *testing.T) {
		req := &UploadToNewRepositoryRequest{
			Name:            "test-repo",
			Version:         "v1.0",
			Type:            "common",
			DiskFormat:      "qcow2",
			ContainerFormat: "bare",
			OperatingSystem: "linux",
			Description:     "test description",
			Filepath:        "s3://bucket/test.qcow2",
		}

		converted := req.ToUploadImageRequest()
		if converted.Name != req.Name {
			t.Errorf("Name = %v, want %v", converted.Name, req.Name)
		}
		if converted.Version != req.Version {
			t.Errorf("Version = %v, want %v", converted.Version, req.Version)
		}
		if converted.Filepath != req.Filepath {
			t.Errorf("Filepath = %v, want %v", converted.Filepath, req.Filepath)
		}
	})
}

// T044: Unit test for UploadToExistingRepositoryRequest
func TestUploadToExistingRepositoryRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *UploadToExistingRepositoryRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &UploadToExistingRepositoryRequest{
				RepositoryID:    "repo-123",
				Version:         "v2",
				Type:            "common",
				DiskFormat:      "raw",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image.raw",
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
			name: "missing repository ID",
			req: &UploadToExistingRepositoryRequest{
				Version:         "v2",
				Type:            "common",
				DiskFormat:      "raw",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image.raw",
			},
			wantErr: true,
			errMsg:  "repositoryId is required",
		},
		{
			name: "missing version",
			req: &UploadToExistingRepositoryRequest{
				RepositoryID:    "repo-123",
				Type:            "common",
				DiskFormat:      "raw",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image.raw",
			},
			wantErr: true,
			errMsg:  "version is required",
		},
		{
			name: "missing type",
			req: &UploadToExistingRepositoryRequest{
				RepositoryID:    "repo-123",
				Version:         "v2",
				DiskFormat:      "raw",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image.raw",
			},
			wantErr: true,
			errMsg:  "type is required",
		},
		{
			name: "missing disk format",
			req: &UploadToExistingRepositoryRequest{
				RepositoryID:    "repo-123",
				Version:         "v2",
				Type:            "common",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image.raw",
			},
			wantErr: true,
			errMsg:  "diskFormat is required",
		},
		{
			name: "invalid disk format",
			req: &UploadToExistingRepositoryRequest{
				RepositoryID:    "repo-123",
				Version:         "v2",
				Type:            "common",
				DiskFormat:      "invalid",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image.raw",
			},
			wantErr: true,
			errMsg:  "invalid diskFormat",
		},
		{
			name: "missing container format",
			req: &UploadToExistingRepositoryRequest{
				RepositoryID: "repo-123",
				Version:      "v2",
				Type:         "common",
				DiskFormat:   "raw",
				Filepath:     "s3://bucket/image.raw",
			},
			wantErr: true,
			errMsg:  "containerFormat is required",
		},
		{
			name: "invalid container format",
			req: &UploadToExistingRepositoryRequest{
				RepositoryID:    "repo-123",
				Version:         "v2",
				Type:            "common",
				DiskFormat:      "raw",
				ContainerFormat: "invalid",
				Filepath:        "s3://bucket/image.raw",
			},
			wantErr: true,
			errMsg:  "invalid containerFormat",
		},
		{
			name: "missing filepath",
			req: &UploadToExistingRepositoryRequest{
				RepositoryID:    "repo-123",
				Version:         "v2",
				Type:            "common",
				DiskFormat:      "raw",
				ContainerFormat: "bare",
			},
			wantErr: true,
			errMsg:  "filepath is required",
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

	// Test ToUploadImageRequest conversion
	t.Run("ToUploadImageRequest conversion", func(t *testing.T) {
		req := &UploadToExistingRepositoryRequest{
			RepositoryID:    "repo-456",
			Version:         "v3.0",
			Type:            "image",
			DiskFormat:      "qcow2",
			ContainerFormat: "bare",
			Filepath:        "s3://bucket/test.qcow2",
		}

		converted := req.ToUploadImageRequest()
		if converted.RepositoryID != req.RepositoryID {
			t.Errorf("RepositoryID = %v, want %v", converted.RepositoryID, req.RepositoryID)
		}
		if converted.Version != req.Version {
			t.Errorf("Version = %v, want %v", converted.Version, req.Version)
		}
		if converted.Filepath != req.Filepath {
			t.Errorf("Filepath = %v, want %v", converted.Filepath, req.Filepath)
		}
	})
}

// T045: Unit test for UploadToExistingTagRequest
func TestUploadToExistingTagRequest(t *testing.T) {
	tests := []struct {
		name    string
		req     *UploadToExistingTagRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			req: &UploadToExistingTagRequest{
				TagID:    "tag-123",
				Filepath: "s3://bucket/image.qcow2",
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
			name: "missing tag ID",
			req: &UploadToExistingTagRequest{
				Filepath: "s3://bucket/image.qcow2",
			},
			wantErr: true,
			errMsg:  "tagId is required",
		},
		{
			name: "missing filepath",
			req: &UploadToExistingTagRequest{
				TagID: "tag-123",
			},
			wantErr: true,
			errMsg:  "filepath is required",
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

	// Test ToUploadImageRequest conversion
	t.Run("ToUploadImageRequest conversion", func(t *testing.T) {
		req := &UploadToExistingTagRequest{
			TagID:    "tag-789",
			Filepath: "s3://bucket/test.qcow2",
		}

		converted := req.ToUploadImageRequest()
		if converted.TagID != req.TagID {
			t.Errorf("TagID = %v, want %v", converted.TagID, req.TagID)
		}
		if converted.Filepath != req.Filepath {
			t.Errorf("Filepath = %v, want %v", converted.Filepath, req.Filepath)
		}
	})
}

// T046: Unit test for response structures
func TestResponseStructures(t *testing.T) {
	// Test CreateSnapshotResponse
	snapshotResp := &CreateSnapshotResponse{
		Repository: &Repository{
			ID:              "repo-123",
			Name:            "snapshot-repo",
			Namespace:       "public",
			OperatingSystem: "linux",
			Count:           1,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		Tag: &Tag{
			ID:           "tag-456",
			Name:         "v1",
			RepositoryID: "repo-123",
			Type:         "common",
			Size:         1024,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	if snapshotResp.Repository == nil {
		t.Error("CreateSnapshotResponse.Repository should not be nil")
	}
	if snapshotResp.Tag == nil {
		t.Error("CreateSnapshotResponse.Tag should not be nil")
	}

	// Test UploadImageResponse
	uploadResp := &UploadImageResponse{
		Repository: &Repository{
			ID:              "repo-789",
			Name:            "upload-repo",
			Namespace:       "private",
			OperatingSystem: "windows",
			Count:           2,
			CreatedAt:       time.Now(),
			UpdatedAt:       time.Now(),
		},
		Tag: &Tag{
			ID:           "tag-101",
			Name:         "v2",
			RepositoryID: "repo-789",
			Type:         "image",
			Size:         2048,
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		},
	}

	if uploadResp.Repository == nil {
		t.Error("UploadImageResponse.Repository should not be nil")
	}
	if uploadResp.Tag == nil {
		t.Error("UploadImageResponse.Tag should not be nil")
	}

	// Test ListRepositoriesResponse
	listResp := &ListRepositoriesResponse{
		Repositories: []*Repository{
			{
				ID:              "repo-1",
				Name:            "repo1",
				Namespace:       "public",
				OperatingSystem: "linux",
				Count:           1,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              "repo-2",
				Name:            "repo2",
				Namespace:       "private",
				OperatingSystem: "windows",
				Count:           2,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
		Total: 2,
	}

	if len(listResp.Repositories) != 2 {
		t.Errorf("ListRepositoriesResponse should have 2 repositories, got %d", len(listResp.Repositories))
	}
	if listResp.Total != 2 {
		t.Errorf("ListRepositoriesResponse.Total should be 2, got %d", listResp.Total)
	}
}

// T047: Unit test for edge cases and boundary conditions
func TestEdgeCasesAndBoundaries(t *testing.T) {
	// Test Repository with negative count
	repo := &Repository{
		ID:              "test-repo",
		Name:            "test",
		Namespace:       "public",
		OperatingSystem: "linux",
		Count:           -1, // Negative count should be invalid
	}

	err := repo.Validate()
	if err == nil {
		t.Error("Repository with negative count should be invalid")
	}

	// Test Repository with zero count (should be valid)
	repo.Count = 0
	err = repo.Validate()
	if err != nil {
		t.Errorf("Repository with zero count should be valid, got error: %v", err)
	}

	// Test Repository with very large count
	repo.Count = 999999
	err = repo.Validate()
	if err != nil {
		t.Errorf("Repository with large count should be valid, got error: %v", err)
	}

	// Test Tag with zero size (Tag doesn't have validation, just basic struct)
	tag := &Tag{
		ID:           "test-tag",
		Name:         "test",
		RepositoryID: "test-repo",
		Type:         "common",
		Size:         0, // Zero size should be valid for Tag
	}

	// Tag doesn't have Validate method, so we just check the struct can be created
	if tag.ID != "test-tag" {
		t.Errorf("Tag ID should be test-tag, got %s", tag.ID)
	}

	// Test Tag with negative size (should be allowed as it's just a field)
	tag.Size = -1
	if tag.Size != -1 {
		t.Errorf("Tag size should be -1, got %d", tag.Size)
	}

	// Test empty extra map
	tag.Size = 1024
	tag.Extra = map[string]interface{}{}
	// No validation for Tag, just check the map is empty
	if len(tag.Extra) != 0 {
		t.Errorf("Tag extra should be empty, got %d items", len(tag.Extra))
	}

	// Test nil extra map
	tag.Extra = nil
	if tag.Extra != nil {
		t.Error("Tag extra should be nil")
	}
}

// T048: Unit test for JSON marshaling/unmarshaling edge cases
func TestJSONEdgeCases(t *testing.T) {
	// Test Repository with empty timestamps
	repo := &Repository{
		ID:              "test-repo",
		Name:            "test",
		Namespace:       "public",
		OperatingSystem: "linux",
		Count:           1,
		CreatedAt:       time.Time{}, // Zero time
		UpdatedAt:       time.Time{}, // Zero time
	}

	data, err := repo.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON with zero times failed: %v", err)
	}

	var unmarshaled Repository
	err = unmarshaled.UnmarshalJSON(data)
	if err != nil {
		t.Fatalf("UnmarshalJSON with zero times failed: %v", err)
	}

	// Test Tag with empty timestamps (using standard JSON)
	tag := &Tag{
		ID:           "test-tag",
		Name:         "test",
		RepositoryID: "test-repo",
		Type:         "common",
		Size:         1024,
		CreatedAt:    time.Time{}, // Zero time
		UpdatedAt:    time.Time{}, // Zero time
	}

	data, err = json.Marshal(tag)
	if err != nil {
		t.Fatalf("Tag JSON marshal with zero times failed: %v", err)
	}

	var unmarshaledTag Tag
	err = json.Unmarshal(data, &unmarshaledTag)
	if err != nil {
		t.Fatalf("Tag JSON unmarshal with zero times failed: %v", err)
	}

	// Test invalid JSON for Repository
	invalidJSON := []byte(`{"id": "test", "invalid": }`)
	var repo2 Repository
	err = repo2.UnmarshalJSON(invalidJSON)
	if err == nil {
		t.Error("UnmarshalJSON with invalid JSON should fail")
	}

	// Test invalid timestamp format
	invalidTimeJSON := []byte(`{
		"id": "test-repo",
		"name": "test",
		"namespace": "public",
		"operatingSystem": "linux",
		"count": 1,
		"createdAt": "invalid-timestamp",
		"updatedAt": "2024-01-01T00:00:00Z"
	}`)
	err = repo2.UnmarshalJSON(invalidTimeJSON)
	if err == nil {
		t.Error("UnmarshalJSON with invalid timestamp should fail")
	}
}

// T049: Unit test for interface compliance
func TestInterfaceCompliance(t *testing.T) {
	// Test that concrete types implement interfaces
	var _ SnapshotRequester = (*CreateSnapshotRequest)(nil)
	var _ SnapshotRequester = (*CreateSnapshotFromNewRepositoryRequest)(nil)
	var _ SnapshotRequester = (*CreateSnapshotFromExistingRepositoryRequest)(nil)

	var _ UploadRequester = (*UploadToNewRepositoryRequest)(nil)
	var _ UploadRequester = (*UploadToExistingRepositoryRequest)(nil)
	var _ UploadRequester = (*UploadToExistingTagRequest)(nil)

	// Test interface method calls
	snapshotReq := &CreateSnapshotFromNewRepositoryRequest{
		Name:            "test",
		OperatingSystem: "linux",
		Version:         "v1",
	}

	var snapshotRequester SnapshotRequester = snapshotReq
	converted := snapshotRequester.ToCreateSnapshotRequest()
	if converted.Name != "test" {
		t.Errorf("Interface method call failed: expected Name='test', got %s", converted.Name)
	}

	uploadReq := &UploadToNewRepositoryRequest{
		Name:            "test-repo",
		OperatingSystem: "linux",
		Version:         "v1",
		Type:            "common",
		DiskFormat:      "qcow2",
		ContainerFormat: "bare",
		Filepath:        "s3://bucket/test",
	}

	var uploadRequester UploadRequester = uploadReq
	convertedUpload := uploadRequester.ToUploadImageRequest()
	if convertedUpload.Name != "test-repo" {
		t.Errorf("Interface method call failed: expected Name='test-repo', got %s", convertedUpload.Name)
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
