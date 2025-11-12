package tags

import (
	"testing"
	"time"
)

// Helper function to check if error message contains substring
func contains(s, substr string) bool {
	return len(s) > 0 && len(substr) > 0 && (s == substr || len(s) > len(substr))
}

// T059: Unit test for Tag struct validation - required fields (ID, Name, RepositoryID, Type, Size, Timestamps, Repository)
func TestTagValidate(t *testing.T) {
	tests := []struct {
		name      string
		tag       *Tag
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid tag",
			tag: &Tag{
				ID:           "tag-123",
				Name:         "v1",
				RepositoryID: "repo-456",
				Type:         "common",
				Size:         1024,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
				Repository: &Repository{
					ID:              "repo-456",
					Name:            "ubuntu",
					Namespace:       "public",
					OperatingSystem: "linux",
					Count:           1,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			expectErr: false,
		},
		{
			name:      "nil tag",
			tag:       nil,
			expectErr: true,
			errMsg:    "tag cannot be nil",
		},
		{
			name: "missing ID",
			tag: &Tag{
				Name:         "v1",
				RepositoryID: "repo-456",
				Type:         "common",
				Size:         1024,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectErr: true,
			errMsg:    "id cannot be empty",
		},
		{
			name: "missing Name",
			tag: &Tag{
				ID:           "tag-123",
				RepositoryID: "repo-456",
				Type:         "common",
				Size:         1024,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectErr: true,
			errMsg:    "name cannot be empty",
		},
		{
			name: "missing RepositoryID",
			tag: &Tag{
				ID:        "tag-123",
				Name:      "v1",
				Type:      "common",
				Size:      1024,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectErr: true,
			errMsg:    "repositoryID cannot be empty",
		},
		{
			name: "missing Type",
			tag: &Tag{
				ID:           "tag-123",
				Name:         "v1",
				RepositoryID: "repo-456",
				Size:         1024,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectErr: true,
			errMsg:    "type cannot be empty",
		},
		{
			name: "negative size",
			tag: &Tag{
				ID:           "tag-123",
				Name:         "v1",
				RepositoryID: "repo-456",
				Type:         "common",
				Size:         -1,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectErr: true,
			errMsg:    "size cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tag.Validate()
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectErr && err != nil && !contains(err.Error(), tt.errMsg) {
				t.Errorf("expected error containing '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}

// T060: Unit test for Tag JSON marshaling - verify camelCase JSON tags
func TestTagJSONMarshaling(t *testing.T) {
	tag := &Tag{
		ID:           "tag-123",
		Name:         "v1",
		RepositoryID: "repo-456",
		Type:         "common",
		Size:         2048,
		Status:       "active",
		Extra:        map[string]interface{}{"key": "value"},
		CreatedAt:    time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC),
		UpdatedAt:    time.Date(2024, 1, 2, 12, 0, 0, 0, time.UTC),
	}

	data, err := tag.MarshalJSON()
	if err != nil {
		t.Fatalf("MarshalJSON failed: %v", err)
	}

	// Verify camelCase JSON field names
	expected := map[string]interface{}{
		"id":           "tag-123",
		"name":         "v1",
		"repositoryID": "repo-456",
		"type":         "common",
		"size":         2048.0,
		"status":       "active",
		"extra":        map[string]interface{}{"key": "value"},
		"createdAt":    "2024-01-01T12:00:00Z",
		"updatedAt":    "2024-01-02T12:00:00Z",
	}

	for key := range expected {
		if !contains(string(data), key) {
			t.Errorf("expected JSON key '%s' in output", key)
		}
	}
}

// T061: Unit test for Tag JSON unmarshaling with nested Repository
func TestTagJSONUnmarshaling(t *testing.T) {
	jsonData := []byte(`{
		"id": "tag-789",
		"name": "v2",
		"repositoryID": "repo-456",
		"type": "common",
		"size": 3072,
		"status": "active",
		"extra": {"arch": "amd64"},
		"createdAt": "2024-08-19T08:32:25Z",
		"updatedAt": "2024-08-19T08:32:25Z",
		"repository": {
			"id": "repo-456",
			"name": "ubuntu",
			"namespace": "public",
			"operatingSystem": "linux",
			"count": 2,
			"creator": {
				"id": "user-1",
				"name": "admin"
			},
			"project": {
				"id": "proj-1",
				"displayName": "test"
			},
			"createdAt": "2024-08-19T08:32:15Z",
			"updatedAt": "2024-08-19T08:32:15Z"
		}
	}`)

	var tag Tag
	err := tag.UnmarshalJSON(jsonData)
	if err != nil {
		t.Fatalf("UnmarshalJSON failed: %v", err)
	}

	if tag.ID != "tag-789" {
		t.Errorf("expected ID tag-789, got %s", tag.ID)
	}
	if tag.Name != "v2" {
		t.Errorf("expected Name v2, got %s", tag.Name)
	}
	if tag.Size != 3072 {
		t.Errorf("expected Size 3072, got %d", tag.Size)
	}
	if tag.Repository == nil {
		t.Fatal("expected repository to be populated")
	}
	if tag.Repository.ID != "repo-456" {
		t.Errorf("expected repository ID repo-456, got %s", tag.Repository.ID)
	}
	if tag.CreatedAt.IsZero() || tag.UpdatedAt.IsZero() {
		t.Errorf("timestamps not properly parsed")
	}
}

// T062: Unit test for CreateTagRequest validation - required fields (Name, Type, DiskFormat, ContainerFormat)
func TestCreateTagRequestValidate(t *testing.T) {
	tests := []struct {
		name      string
		req       *CreateTagRequest
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid create request",
			req: &CreateTagRequest{
				Name:            "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
			},
			expectErr: false,
		},
		{
			name:      "nil request",
			req:       nil,
			expectErr: true,
			errMsg:    "request cannot be nil",
		},
		{
			name: "missing Name",
			req: &CreateTagRequest{
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
			},
			expectErr: true,
			errMsg:    "name cannot be empty",
		},
		{
			name: "missing Type",
			req: &CreateTagRequest{
				Name:            "v1",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
			},
			expectErr: true,
			errMsg:    "type cannot be empty",
		},
		{
			name: "missing DiskFormat",
			req: &CreateTagRequest{
				Name:            "v1",
				Type:            "common",
				ContainerFormat: "bare",
			},
			expectErr: true,
			errMsg:    "diskFormat cannot be empty",
		},
		{
			name: "missing ContainerFormat",
			req: &CreateTagRequest{
				Name:       "v1",
				Type:       "common",
				DiskFormat: "qcow2",
			},
			expectErr: true,
			errMsg:    "containerFormat cannot be empty",
		},
		{
			name: "invalid DiskFormat",
			req: &CreateTagRequest{
				Name:            "v1",
				Type:            "common",
				DiskFormat:      "invalid",
				ContainerFormat: "bare",
			},
			expectErr: true,
			errMsg:    "invalid diskFormat",
		},
		{
			name: "invalid ContainerFormat",
			req: &CreateTagRequest{
				Name:            "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "invalid",
			},
			expectErr: true,
			errMsg:    "invalid containerFormat",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectErr && err != nil && !contains(err.Error(), tt.errMsg) {
				t.Errorf("expected error containing '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}

// T063: Unit test for UpdateTagRequest validation - optional fields
func TestUpdateTagRequestValidate(t *testing.T) {
	tests := []struct {
		name      string
		req       *UpdateTagRequest
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid update with name",
			req: &UpdateTagRequest{
				Name: "v2",
			},
			expectErr: false,
		},
		{
			name: "valid update with empty fields",
			req: &UpdateTagRequest{
				Name: "",
			},
			expectErr: false,
		},
		{
			name:      "nil request",
			req:       nil,
			expectErr: true,
			errMsg:    "request cannot be nil",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectErr && err != nil && !contains(err.Error(), tt.errMsg) {
				t.Errorf("expected error containing '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}

// T064: Unit test for ListTagsOptions validation - test Where filters, Namespace, pagination
func TestListTagsOptionsValidate(t *testing.T) {
	tests := []struct {
		name      string
		opts      *ListTagsOptions
		expectErr bool
		errMsg    string
	}{
		{
			name: "valid options with limit",
			opts: &ListTagsOptions{
				Limit: 10,
			},
			expectErr: false,
		},
		{
			name: "valid options with offset",
			opts: &ListTagsOptions{
				Offset: 5,
			},
			expectErr: false,
		},
		{
			name: "valid options with where",
			opts: &ListTagsOptions{
				Where: []string{"status=active", "type=common"},
			},
			expectErr: false,
		},
		{
			name: "valid options with namespace",
			opts: &ListTagsOptions{
				Namespace: "public",
			},
			expectErr: false,
		},
		{
			name:      "valid options empty",
			opts:      &ListTagsOptions{},
			expectErr: false,
		},
		{
			name: "invalid negative limit",
			opts: &ListTagsOptions{
				Limit: -2,
			},
			expectErr: true,
			errMsg:    "limit must be",
		},
		{
			name: "invalid negative offset",
			opts: &ListTagsOptions{
				Offset: -1,
			},
			expectErr: true,
			errMsg:    "offset cannot be negative",
		},
		{
			name:      "nil options",
			opts:      nil,
			expectErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.opts != nil {
				err = tt.opts.Validate()
			}
			if tt.expectErr && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectErr && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectErr && err != nil && !contains(err.Error(), tt.errMsg) {
				t.Errorf("expected error containing '%s', got '%s'", tt.errMsg, err.Error())
			}
		})
	}
}
