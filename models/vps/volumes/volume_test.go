package volumes

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

func TestVolume_JSONMarshaling(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name    string
		volume  *Volume
		wantErr bool
	}{
		{
			name: "complete volume",
			volume: &Volume{
				ID:           "vol-123",
				Name:         "test-volume",
				Description:  "Test volume description",
				Size:         100,
				Type:         "SSD",
				Status:       VolumeStatusAvailable,
				StatusReason: "",
				Attachments:  []common.IDName{{ID: "srv-456", Name: "server-1"}},
				Project:      common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID:    "proj-789",
				User:         common.IDName{ID: "user-111", Name: "testuser"},
				UserID:       "user-111",
				Namespace:    "default",
				CreatedAt:    &now,
				UpdatedAt:    &now,
			},
			wantErr: false,
		},
		{
			name: "minimal volume",
			volume: &Volume{
				ID:        "vol-456",
				Name:      "minimal",
				Size:      50,
				Type:      "HDD",
				Status:    VolumeStatusCreating,
				Project:   common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID: "proj-789",
				User:      common.IDName{ID: "user-111", Name: "testuser"},
				UserID:    "user-111",
				Namespace: "default",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Marshal to JSON
			data, err := json.Marshal(tt.volume)
			if err != nil {
				t.Fatalf("failed to marshal volume: %v", err)
			}

			// Unmarshal back
			var unmarshaled Volume
			err = json.Unmarshal(data, &unmarshaled)
			if (err != nil) != tt.wantErr {
				t.Errorf("Unmarshal() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				// Verify key fields
				if unmarshaled.ID != tt.volume.ID {
					t.Errorf("ID mismatch: got %v, want %v", unmarshaled.ID, tt.volume.ID)
				}
				if unmarshaled.Name != tt.volume.Name {
					t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, tt.volume.Name)
				}
				if unmarshaled.Size != tt.volume.Size {
					t.Errorf("Size mismatch: got %v, want %v", unmarshaled.Size, tt.volume.Size)
				}
				if unmarshaled.Status != tt.volume.Status {
					t.Errorf("Status mismatch: got %v, want %v", unmarshaled.Status, tt.volume.Status)
				}
			}
		})
	}
}

func TestVolume_Validate(t *testing.T) {
	tests := []struct {
		name    string
		volume  *Volume
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid volume",
			volume: &Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Size:      100,
				Type:      "SSD",
				Status:    VolumeStatusAvailable,
				Project:   common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID: "proj-789",
				User:      common.IDName{ID: "user-111", Name: "testuser"},
				UserID:    "user-111",
				Namespace: "default",
			},
			wantErr: false,
		},
		{
			name: "missing ID",
			volume: &Volume{
				Name:      "test-volume",
				Size:      100,
				Type:      "SSD",
				Status:    VolumeStatusAvailable,
				Project:   common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID: "proj-789",
				User:      common.IDName{ID: "user-111", Name: "testuser"},
				UserID:    "user-111",
				Namespace: "default",
			},
			wantErr: true,
			errMsg:  "id is required",
		},
		{
			name: "missing name",
			volume: &Volume{
				ID:        "vol-123",
				Size:      100,
				Type:      "SSD",
				Status:    VolumeStatusAvailable,
				Project:   common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID: "proj-789",
				User:      common.IDName{ID: "user-111", Name: "testuser"},
				UserID:    "user-111",
				Namespace: "default",
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "negative size",
			volume: &Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Size:      -10,
				Type:      "SSD",
				Status:    VolumeStatusAvailable,
				Project:   common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID: "proj-789",
				User:      common.IDName{ID: "user-111", Name: "testuser"},
				UserID:    "user-111",
				Namespace: "default",
			},
			wantErr: true,
			errMsg:  "size must be >= 0",
		},
		{
			name: "missing type",
			volume: &Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Size:      100,
				Status:    VolumeStatusAvailable,
				Project:   common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID: "proj-789",
				User:      common.IDName{ID: "user-111", Name: "testuser"},
				UserID:    "user-111",
				Namespace: "default",
			},
			wantErr: true,
			errMsg:  "type is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.volume.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestVolumeStatus_Constants(t *testing.T) {
	tests := []struct {
		name   string
		status VolumeStatus
		want   string
	}{
		{"creating", VolumeStatusCreating, "creating"},
		{"available", VolumeStatusAvailable, "available"},
		{"in-use", VolumeStatusInUse, "in-use"},
		{"detaching", VolumeStatusDetaching, "detaching"},
		{"extending", VolumeStatusExtending, "extending"},
		{"deleting", VolumeStatusDeleting, "deleting"},
		{"deleted", VolumeStatusDeleted, "deleted"},
		{"error", VolumeStatusError, "error"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.status) != tt.want {
				t.Errorf("VolumeStatus constant = %v, want %v", tt.status, tt.want)
			}
		})
	}
}

func TestVolumeAction_Constants(t *testing.T) {
	tests := []struct {
		name   string
		action VolumeAction
		want   string
	}{
		{"attach", VolumeActionAttach, "attach"},
		{"detach", VolumeActionDetach, "detach"},
		{"extend", VolumeActionExtend, "extend"},
		{"revert", VolumeActionRevert, "revert"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if string(tt.action) != tt.want {
				t.Errorf("VolumeAction constant = %v, want %v", tt.action, tt.want)
			}
		})
	}
}

func TestCreateVolumeRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *CreateVolumeRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request",
			request: &CreateVolumeRequest{
				Name: "test-volume",
				Type: "SSD",
				Size: 100,
			},
			wantErr: false,
		},
		{
			name: "valid request without size",
			request: &CreateVolumeRequest{
				Name: "test-volume",
				Type: "SSD",
			},
			wantErr: false,
		},
		{
			name: "missing name",
			request: &CreateVolumeRequest{
				Type: "SSD",
				Size: 100,
			},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name: "missing type",
			request: &CreateVolumeRequest{
				Name: "test-volume",
				Size: 100,
			},
			wantErr: true,
			errMsg:  "type is required",
		},
		{
			name: "negative size",
			request: &CreateVolumeRequest{
				Name: "test-volume",
				Type: "SSD",
				Size: -10,
			},
			wantErr: true,
			errMsg:  "size must be >= 0 if provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUpdateVolumeRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *UpdateVolumeRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid request with name",
			request: &UpdateVolumeRequest{
				Name: "new-name",
			},
			wantErr: false,
		},
		{
			name: "valid request with description",
			request: &UpdateVolumeRequest{
				Description: "new description",
			},
			wantErr: false,
		},
		{
			name: "valid request with both",
			request: &UpdateVolumeRequest{
				Name:        "new-name",
				Description: "new description",
			},
			wantErr: false,
		},
		{
			name:    "empty request",
			request: &UpdateVolumeRequest{},
			wantErr: true,
			errMsg:  "at least one field (name or description) must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestVolumeActionRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		request *VolumeActionRequest
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid attach",
			request: &VolumeActionRequest{
				Action:   VolumeActionAttach,
				ServerID: "srv-123",
			},
			wantErr: false,
		},
		{
			name: "valid detach",
			request: &VolumeActionRequest{
				Action:   VolumeActionDetach,
				ServerID: "srv-123",
			},
			wantErr: false,
		},
		{
			name: "valid extend",
			request: &VolumeActionRequest{
				Action:  VolumeActionExtend,
				NewSize: 200,
			},
			wantErr: false,
		},
		{
			name: "valid revert",
			request: &VolumeActionRequest{
				Action: VolumeActionRevert,
			},
			wantErr: false,
		},
		{
			name: "attach missing server_id",
			request: &VolumeActionRequest{
				Action: VolumeActionAttach,
			},
			wantErr: true,
			errMsg:  "server_id is required for attach/detach actions",
		},
		{
			name: "detach missing server_id",
			request: &VolumeActionRequest{
				Action: VolumeActionDetach,
			},
			wantErr: true,
			errMsg:  "server_id is required for attach/detach actions",
		},
		{
			name: "extend missing new_size",
			request: &VolumeActionRequest{
				Action: VolumeActionExtend,
			},
			wantErr: true,
			errMsg:  "new_size must be positive for extend action",
		},
		{
			name: "extend negative new_size",
			request: &VolumeActionRequest{
				Action:  VolumeActionExtend,
				NewSize: -100,
			},
			wantErr: true,
			errMsg:  "new_size must be positive for extend action",
		},
		{
			name: "invalid action",
			request: &VolumeActionRequest{
				Action: "invalid",
			},
			wantErr: true,
			errMsg:  "invalid action: must be one of attach, detach, extend, revert",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.errMsg {
				t.Errorf("Validate() error message = %v, want %v", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestListVolumesOptions_Validate(t *testing.T) {
	tests := []struct {
		name    string
		options *ListVolumesOptions
		wantErr bool
	}{
		{
			name: "valid options with all fields",
			options: &ListVolumesOptions{
				Name:   "test",
				UserID: "user-123",
				Status: "available",
				Type:   "SSD",
				Detail: true,
			},
			wantErr: false,
		},
		{
			name:    "empty options",
			options: &ListVolumesOptions{},
			wantErr: false,
		},
		{
			name:    "nil options",
			options: nil,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			if tt.options != nil {
				err = tt.options.Validate()
			}
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestVolumeListResponse_JSONUnmarshaling(t *testing.T) {
	jsonData := `{
		"volumes": [
			{
				"id": "vol-123",
				"name": "volume-1",
				"size": 100,
				"type": "SSD",
				"status": "available",
				"project": {"id": "proj-789", "name": "test-project"},
				"project_id": "proj-789",
				"user": {"id": "user-111", "name": "testuser"},
				"user_id": "user-111",
				"namespace": "default"
			},
			{
				"id": "vol-456",
				"name": "volume-2",
				"size": 50,
				"type": "HDD",
				"status": "in-use",
				"project": {"id": "proj-789", "name": "test-project"},
				"project_id": "proj-789",
				"user": {"id": "user-111", "name": "testuser"},
				"user_id": "user-111",
				"namespace": "default"
			}
		]
	}`

	var response VolumeListResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(response.Volumes) != 2 {
		t.Errorf("expected 2 volumes, got %d", len(response.Volumes))
	}

	if response.Volumes[0].ID != "vol-123" {
		t.Errorf("expected vol-123, got %s", response.Volumes[0].ID)
	}

	if response.Volumes[1].Status != VolumeStatusInUse {
		t.Errorf("expected in-use status, got %s", response.Volumes[1].Status)
	}
}
