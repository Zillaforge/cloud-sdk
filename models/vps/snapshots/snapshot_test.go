package snapshots

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

func TestSnapshotValidate(t *testing.T) {
	tests := []struct {
		name     string
		snapshot *Snapshot
		wantErr  bool
		errMsg   string
	}{
		{
			name: "valid snapshot",
			snapshot: &Snapshot{
				ID:       "snap-123",
				VolumeID: "vol-456",
				Name:     "test-snapshot",
				Size:     10,
				Status:   SnapshotStatusAvailable,
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			snapshot: &Snapshot{
				VolumeID: "vol-456",
				Name:     "test-snapshot",
			},
			wantErr: true,
			errMsg:  "id is required",
		},
		{
			name: "missing VolumeID",
			snapshot: &Snapshot{
				ID:   "snap-123",
				Name: "test-snapshot",
			},
			wantErr: true,
			errMsg:  "volume_id is required",
		},
		{
			name: "empty ID",
			snapshot: &Snapshot{
				ID:       "",
				VolumeID: "vol-456",
				Name:     "test-snapshot",
			},
			wantErr: true,
			errMsg:  "id is required",
		},
		{
			name: "empty VolumeID",
			snapshot: &Snapshot{
				ID:       "snap-123",
				VolumeID: "",
				Name:     "test-snapshot",
			},
			wantErr: true,
			errMsg:  "volume_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.snapshot.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestListSnapshotsOptionsValidate(t *testing.T) {
	opts := &ListSnapshotsOptions{}
	if err := opts.Validate(); err != nil {
		t.Errorf("Validate() should always return nil, got %v", err)
	}

	opts = &ListSnapshotsOptions{
		Name:     "test-snap",
		VolumeID: "vol-123",
		UserID:   "user-456",
		Status:   "available",
	}
	if err := opts.Validate(); err != nil {
		t.Errorf("Validate() should always return nil, got %v", err)
	}
}

func TestSnapshotStatusConstants(t *testing.T) {
	expectedStatuses := []SnapshotStatus{
		SnapshotStatusCreating,
		SnapshotStatusAvailable,
		SnapshotStatusDeleting,
		SnapshotStatusDeleted,
		SnapshotStatusError,
	}

	expectedValues := []string{"creating", "available", "deleting", "deleted", "error"}

	for i, status := range expectedStatuses {
		if string(status) != expectedValues[i] {
			t.Errorf("SnapshotStatus constant %v = %v, want %v", status, string(status), expectedValues[i])
		}
	}
}

func TestSnapshotJSONMarshaling(t *testing.T) {
	now := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	snapshot := &Snapshot{
		ID:           "snap-123",
		Name:         "test-snapshot",
		VolumeID:     "vol-456",
		Size:         10,
		Status:       SnapshotStatusAvailable,
		StatusReason: "success",
		Description:  "A test snapshot",
		Project:      common.IDName{ID: "proj-1", Name: "project1"},
		ProjectID:    "proj-1",
		User:         common.IDName{ID: "user-1", Name: "user1"},
		UserID:       "user-1",
		Namespace:    "public",
		CreatedAt:    &now,
		UpdatedAt:    &now,
	}

	data, err := json.Marshal(snapshot)
	if err != nil {
		t.Fatalf("MarshalJSON() error = %v", err)
	}

	// Check if JSON contains expected fields
	expectedFields := []string{
		`"id":"snap-123"`,
		`"name":"test-snapshot"`,
		`"volume_id":"vol-456"`,
		`"size":10`,
		`"status":"available"`,
		`"status_reason":"success"`,
		`"description":"A test snapshot"`,
		`"project":{"id":"proj-1","name":"project1"}`,
		`"project_id":"proj-1"`,
		`"user":{"id":"user-1","name":"user1"}`,
		`"user_id":"user-1"`,
		`"namespace":"public"`,
		`"createdAt":"2023-10-01T12:00:00Z"`,
		`"updatedAt":"2023-10-01T12:00:00Z"`,
	}

	jsonStr := string(data)
	for _, field := range expectedFields {
		if !contains(jsonStr, field) {
			t.Errorf("JSON does not contain expected field: %s", field)
		}
	}
}

func TestSnapshotJSONUnmarshaling(t *testing.T) {
	jsonData := `{
		"id": "snap-123",
		"name": "test-snapshot",
		"volume_id": "vol-456",
		"size": 10,
		"status": "available",
		"status_reason": "success",
		"description": "A test snapshot",
		"project": {"id": "proj-1", "name": "project1"},
		"project_id": "proj-1",
		"user": {"id": "user-1", "name": "user1"},
		"user_id": "user-1",
		"namespace": "public",
		"createdAt": "2023-10-01T12:00:00Z",
		"updatedAt": "2023-10-01T12:00:00Z"
	}`

	var snapshot Snapshot
	err := json.Unmarshal([]byte(jsonData), &snapshot)
	if err != nil {
		t.Fatalf("UnmarshalJSON() error = %v", err)
	}

	expectedTime := time.Date(2023, 10, 1, 12, 0, 0, 0, time.UTC)
	if snapshot.ID != "snap-123" || snapshot.Name != "test-snapshot" || snapshot.VolumeID != "vol-456" ||
		snapshot.Size != 10 || snapshot.Status != SnapshotStatusAvailable || snapshot.StatusReason != "success" ||
		snapshot.Description != "A test snapshot" || snapshot.ProjectID != "proj-1" || snapshot.UserID != "user-1" ||
		snapshot.Namespace != "public" || snapshot.CreatedAt == nil || *snapshot.CreatedAt != expectedTime ||
		snapshot.UpdatedAt == nil || *snapshot.UpdatedAt != expectedTime {
		t.Errorf("UnmarshalJSON() did not parse correctly")
	}

	if snapshot.Project.ID != "proj-1" || snapshot.Project.Name != "project1" ||
		snapshot.User.ID != "user-1" || snapshot.User.Name != "user1" {
		t.Errorf("UnmarshalJSON() did not parse nested objects correctly")
	}
}

func TestCreateSnapshotRequest_Validate(t *testing.T) {
	req := &CreateSnapshotRequest{}
	if err := req.Validate(); err == nil {
		t.Fatal("expected validation error for empty request, got nil")
	}

	req = &CreateSnapshotRequest{Name: "snap", VolumeID: "vol-1"}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestUpdateSnapshotRequest_Validate(t *testing.T) {
	req := &UpdateSnapshotRequest{}
	if err := req.Validate(); err == nil {
		t.Fatal("expected validation error for empty update request, got nil")
	}

	req = &UpdateSnapshotRequest{Name: "new-name"}
	if err := req.Validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Helper function to check if a string contains a substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
