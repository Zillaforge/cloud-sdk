package common

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

func TestIDNameValidate(t *testing.T) {
	tests := []struct {
		name    string
		idname  *IDName
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid IDName",
			idname: &IDName{
				ID:   "test-id-123",
				Name: "test",
			},
			wantErr: false,
		},
		{
			name: "valid IDName with all fields",
			idname: &IDName{
				ID:          "user-456",
				Name:        "john",
				Account:     "john@example.com",
				DisplayName: "John Doe",
			},
			wantErr: false,
		},
		{
			name: "valid IDName with only ID",
			idname: &IDName{
				ID: "repo-789",
			},
			wantErr: false,
		},
		{
			name:    "nil IDName",
			idname:  nil,
			wantErr: true,
			errMsg:  "IDName cannot be nil",
		},
		{
			name: "empty ID",
			idname: &IDName{
				ID:   "",
				Name: "test",
			},
			wantErr: true,
			errMsg:  "IDName.ID is required and must not be empty",
		},
		{
			name: "whitespace-only ID",
			idname: &IDName{
				ID:   "   ",
				Name: "test",
			},
			wantErr: true,
			errMsg:  "IDName.ID is required and must not be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.idname.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestRepositoryValidate(t *testing.T) {
	tests := []struct {
		name    string
		repo    *Repository
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid Repository",
			repo: &Repository{
				ID:              "repo-123",
				Name:            "test-repo",
				Namespace:       "public",
				OperatingSystem: "linux",
				Count:           5,
				Creator:         &IDName{ID: "user-1"},
				Project:         &IDName{ID: "proj-1"},
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			wantErr: false,
		},
		{
			name:    "nil Repository",
			repo:    nil,
			wantErr: true,
			errMsg:  "repository cannot be nil",
		},
		{
			name: "empty ID",
			repo: &Repository{
				ID:              "",
				Name:            "test-repo",
				Namespace:       "public",
				OperatingSystem: "linux",
				Count:           0,
			},
			wantErr: true,
			errMsg:  "ID is required and must not be empty",
		},
		{
			name: "empty Name",
			repo: &Repository{
				ID:              "repo-123",
				Name:            "",
				Namespace:       "public",
				OperatingSystem: "linux",
				Count:           0,
			},
			wantErr: true,
			errMsg:  "name is required and must not be empty",
		},
		{
			name: "invalid Namespace",
			repo: &Repository{
				ID:              "repo-123",
				Name:            "test-repo",
				Namespace:       "invalid",
				OperatingSystem: "linux",
				Count:           0,
			},
			wantErr: true,
			errMsg:  "namespace must be 'public' or 'private'",
		},
		{
			name: "invalid OperatingSystem",
			repo: &Repository{
				ID:              "repo-123",
				Name:            "test-repo",
				Namespace:       "public",
				OperatingSystem: "invalid",
				Count:           0,
			},
			wantErr: true,
			errMsg:  "operating system must be 'linux' or 'windows'",
		},
		{
			name: "negative Count",
			repo: &Repository{
				ID:              "repo-123",
				Name:            "test-repo",
				Namespace:       "public",
				OperatingSystem: "linux",
				Count:           -1,
			},
			wantErr: true,
			errMsg:  "count must be >= 0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestTagValidate(t *testing.T) {
	tests := []struct {
		name    string
		tag     *Tag
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid Tag",
			tag: &Tag{
				ID:           "tag-123",
				Name:         "v1.0",
				RepositoryID: "repo-123",
				Type:         TagTypeCommon,
				Size:         1024,
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			wantErr: false,
		},
		{
			name:    "nil Tag",
			tag:     nil,
			wantErr: true,
			errMsg:  "tag cannot be nil",
		},
		{
			name: "empty ID",
			tag: &Tag{
				ID:           "",
				Name:         "v1.0",
				RepositoryID: "repo-123",
				Type:         TagTypeCommon,
				Size:         1024,
			},
			wantErr: true,
			errMsg:  "id cannot be empty",
		},
		{
			name: "empty Name",
			tag: &Tag{
				ID:           "tag-123",
				Name:         "",
				RepositoryID: "repo-123",
				Type:         TagTypeCommon,
				Size:         1024,
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "empty RepositoryID",
			tag: &Tag{
				ID:           "tag-123",
				Name:         "v1.0",
				RepositoryID: "",
				Type:         TagTypeCommon,
				Size:         1024,
			},
			wantErr: true,
			errMsg:  "repositoryID cannot be empty",
		},
		{
			name: "invalid Type",
			tag: &Tag{
				ID:           "tag-123",
				Name:         "v1.0",
				RepositoryID: "repo-123",
				Type:         TagType("invalid"),
				Size:         1024,
			},
			wantErr: true,
			errMsg:  "type must be 'common' or 'increase'",
		},
		{
			name: "negative Size",
			tag: &Tag{
				ID:           "tag-123",
				Name:         "v1.0",
				RepositoryID: "repo-123",
				Type:         TagTypeCommon,
				Size:         -1,
			},
			wantErr: true,
			errMsg:  "size cannot be negative",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.tag.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestRepositoryMarshalJSON(t *testing.T) {
	now := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	repo := &Repository{
		ID:              "repo-123",
		Name:            "test-repo",
		Namespace:       "public",
		OperatingSystem: "linux",
		Description:     "A test repository",
		Count:           5,
		Creator:         &IDName{ID: "user-1", Name: "creator"},
		Project:         &IDName{ID: "proj-1", Name: "project"},
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	data, err := json.Marshal(repo)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	expected := `{"id":"repo-123","name":"test-repo","namespace":"public","operatingSystem":"linux","description":"A test repository","count":5,"creator":{"id":"user-1","name":"creator"},"project":{"id":"proj-1","name":"project"},"createdAt":"2023-10-01T12:00:00Z","updatedAt":"2023-10-01T12:00:00Z"}`
	if string(data) != expected {
		t.Errorf("MarshalJSON() = %v, want %v", string(data), expected)
	}
}

func TestRepositoryUnmarshalJSON(t *testing.T) {
	jsonData := `{"id":"repo-123","name":"test-repo","namespace":"public","operatingSystem":"linux","description":"A test repository","count":5,"creator":{"id":"user-1","name":"creator"},"project":{"id":"proj-1","name":"project"},"createdAt":"2023-10-01T12:00:00Z","updatedAt":"2023-10-01T12:00:00Z"}`

	var repo Repository
	err := json.Unmarshal([]byte(jsonData), &repo)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	expectedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	if repo.ID != "repo-123" || repo.Name != "test-repo" || repo.CreatedAt != expectedTime || repo.UpdatedAt != expectedTime {
		t.Errorf("UnmarshalJSON() did not parse correctly")
	}
}

func TestTagMarshalJSON(t *testing.T) {
	now := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	tag := &Tag{
		ID:           "tag-123",
		Name:         "v1.0",
		RepositoryID: "repo-123",
		Type:         TagTypeCommon,
		Size:         1024,
		Status:       TagStatusActive,
		Extra:        map[string]interface{}{"key": "value"},
		CreatedAt:    now,
		UpdatedAt:    now,
		Repository: &Repository{
			ID:   "repo-123",
			Name: "test-repo",
		},
	}

	data, err := json.Marshal(tag)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Check if timestamps are formatted correctly
	if !strings.Contains(string(data), `"createdAt":"2023-10-01T12:00:00Z"`) {
		t.Errorf("MarshalJSON() did not format createdAt correctly")
	}
	if !strings.Contains(string(data), `"updatedAt":"2023-10-01T12:00:00Z"`) {
		t.Errorf("MarshalJSON() did not format updatedAt correctly")
	}
}

func TestTagUnmarshalJSON(t *testing.T) {
	jsonData := `{"id":"tag-123","name":"v1.0","repositoryID":"repo-123","type":"common","size":1024,"status":"active","extra":{"key":"value"},"createdAt":"2023-10-01T12:00:00Z","updatedAt":"2023-10-01T12:00:00Z","repository":{"id":"repo-123","name":"test-repo"}}`

	var tag Tag
	err := json.Unmarshal([]byte(jsonData), &tag)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	expectedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	if tag.ID != "tag-123" || tag.Name != "v1.0" || tag.CreatedAt != expectedTime || tag.UpdatedAt != expectedTime {
		t.Errorf("UnmarshalJSON() did not parse correctly")
	}
	if tag.Repository == nil || tag.Repository.ID != "repo-123" {
		t.Errorf("UnmarshalJSON() did not parse repository correctly")
	}
}

func TestDiskFormatIsValid(t *testing.T) {
	tests := []struct {
		format DiskFormat
		valid  bool
	}{
		{DiskFormatAMI, true},
		{DiskFormatARI, true},
		{DiskFormatAKI, true},
		{DiskFormatVHD, true},
		{DiskFormatVMDK, true},
		{DiskFormatRaw, true},
		{DiskFormatQcow2, true},
		{DiskFormatVDI, true},
		{DiskFormatISO, true},
		{DiskFormat("invalid"), false},
		{DiskFormat(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestContainerFormatIsValid(t *testing.T) {
	tests := []struct {
		format ContainerFormat
		valid  bool
	}{
		{ContainerFormatAMI, true},
		{ContainerFormatARI, true},
		{ContainerFormatAKI, true},
		{ContainerFormatBare, true},
		{ContainerFormatOVF, true},
		{ContainerFormat("invalid"), false},
		{ContainerFormat(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestDiskFormatString(t *testing.T) {
	tests := []struct {
		format   DiskFormat
		expected string
	}{
		{DiskFormatAMI, "ami"},
		{DiskFormatQcow2, "qcow2"},
		{DiskFormatISO, "iso"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestContainerFormatString(t *testing.T) {
	tests := []struct {
		format   ContainerFormat
		expected string
	}{
		{ContainerFormatAMI, "ami"},
		{ContainerFormatBare, "bare"},
		{ContainerFormatOVF, "ovf"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			if got := tt.format.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTagTypeIsValid(t *testing.T) {
	tests := []struct {
		tagType TagType
		valid   bool
	}{
		{TagTypeCommon, true},
		{TagTypeIncrease, true},
		{TagType("invalid"), false},
		{TagType(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.tagType), func(t *testing.T) {
			if got := tt.tagType.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestTagTypeString(t *testing.T) {
	tests := []struct {
		tagType  TagType
		expected string
	}{
		{TagTypeCommon, "common"},
		{TagTypeIncrease, "increase"},
	}

	for _, tt := range tests {
		t.Run(string(tt.tagType), func(t *testing.T) {
			if got := tt.tagType.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestTagStatusIsValid(t *testing.T) {
	tests := []struct {
		status TagStatus
		valid  bool
	}{
		{TagStatusQueued, true},
		{TagStatusSaving, true},
		{TagStatusImporting, true},
		{TagStatusCreating, true},
		{TagStatusRestoring, true},
		{TagStatusActive, true},
		{TagStatusKilled, true},
		{TagStatusPendingDelete, true},
		{TagStatusDeactivated, true},
		{TagStatusAvailable, true},
		{TagStatusBackingUp, true},
		{TagStatusDeleting, true},
		{TagStatusError, true},
		{TagStatusUnmanaging, true},
		{TagStatusErrorDeleting, true},
		{TagStatusDeleted, true},
		{TagStatus("invalid"), false},
		{TagStatus(""), false},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.IsValid(); got != tt.valid {
				t.Errorf("IsValid() = %v, want %v", got, tt.valid)
			}
		})
	}
}

func TestTagStatusString(t *testing.T) {
	tests := []struct {
		status   TagStatus
		expected string
	}{
		{TagStatusQueued, "queued"},
		{TagStatusSaving, "saving"},
		{TagStatusImporting, "importing"},
		{TagStatusCreating, "creating"},
		{TagStatusRestoring, "restoring"},
		{TagStatusActive, "active"},
		{TagStatusKilled, "killed"},
		{TagStatusPendingDelete, "pending_delete"},
		{TagStatusDeactivated, "deactivated"},
		{TagStatusAvailable, "available"},
		{TagStatusBackingUp, "backing-up"},
		{TagStatusDeleting, "deleting"},
		{TagStatusError, "error"},
		{TagStatusUnmanaging, "unmanaging"},
		{TagStatusErrorDeleting, "error_deleting"},
		{TagStatusDeleted, "deleted"},
	}

	for _, tt := range tests {
		t.Run(string(tt.status), func(t *testing.T) {
			if got := tt.status.String(); got != tt.expected {
				t.Errorf("String() = %v, want %v", got, tt.expected)
			}
		})
	}
}
