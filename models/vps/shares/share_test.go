package shares

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

func TestShare_JSONMarshaling(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name  string
		share *Share
	}{
		{
			name: "complete share",
			share: &Share{
				ID:           "share-123",
				Name:         "test-share",
				Description:  "Test share description",
				Size:         100,
				Path:         "10.0.0.1:/share-123",
				Status:       ShareStatusAvailable,
				StatusReason: "",
				Network:      common.IDName{ID: "net-456", Name: "test-network"},
				NetworkID:    "net-456",
				Project:      common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID:    "proj-789",
				User:         common.IDName{ID: "user-111", Name: "testuser"},
				UserID:       "user-111",
				Namespace:    "default",
				Rules: []ShareRule{
					{ID: "rule-1", CIDR: "10.0.0.0/24", Level: "rw", Status: "active"},
				},
				CreatedAt: &now,
				UpdatedAt: &now,
			},
		},
		{
			name: "minimal share",
			share: &Share{
				ID:        "share-456",
				Name:      "minimal",
				Size:      50,
				Status:    ShareStatusCreating,
				Network:   common.IDName{ID: "net-789", Name: "net"},
				NetworkID: "net-789",
				Project:   common.IDName{ID: "proj-789", Name: "test-project"},
				ProjectID: "proj-789",
				User:      common.IDName{ID: "user-111", Name: "testuser"},
				UserID:    "user-111",
				Namespace: "default",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.share)
			if err != nil {
				t.Fatalf("failed to marshal share: %v", err)
			}

			var unmarshaled Share
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("failed to unmarshal share: %v", err)
			}

			if unmarshaled.ID != tt.share.ID {
				t.Errorf("ID mismatch: got %v, want %v", unmarshaled.ID, tt.share.ID)
			}
			if unmarshaled.Name != tt.share.Name {
				t.Errorf("Name mismatch: got %v, want %v", unmarshaled.Name, tt.share.Name)
			}
			if unmarshaled.Size != tt.share.Size {
				t.Errorf("Size mismatch: got %v, want %v", unmarshaled.Size, tt.share.Size)
			}
			if unmarshaled.Status != tt.share.Status {
				t.Errorf("Status mismatch: got %v, want %v", unmarshaled.Status, tt.share.Status)
			}
		})
	}
}

func TestShare_Validate(t *testing.T) {
	tests := []struct {
		name    string
		share   *Share
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid share",
			share:   &Share{ID: "share-1", Name: "test"},
			wantErr: false,
		},
		{
			name:    "missing id",
			share:   &Share{ID: "", Name: "test"},
			wantErr: true,
			errMsg:  "id is required",
		},
		{
			name:    "missing name",
			share:   &Share{ID: "share-1", Name: ""},
			wantErr: true,
			errMsg:  "name is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.share.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestCreateShareRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateShareRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid request",
			req:     &CreateShareRequest{Name: "share-1", NetworkID: "net-1", Size: 10},
			wantErr: false,
		},
		{
			name:    "missing name",
			req:     &CreateShareRequest{Name: "", NetworkID: "net-1", Size: 10},
			wantErr: true,
			errMsg:  "name is required",
		},
		{
			name:    "missing network_id",
			req:     &CreateShareRequest{Name: "share-1", NetworkID: "", Size: 10},
			wantErr: true,
			errMsg:  "network_id is required",
		},
		{
			name:    "zero size",
			req:     &CreateShareRequest{Name: "share-1", NetworkID: "net-1", Size: 0},
			wantErr: true,
			errMsg:  "size must be positive",
		},
		{
			name:    "negative size",
			req:     &CreateShareRequest{Name: "share-1", NetworkID: "net-1", Size: -5},
			wantErr: true,
			errMsg:  "size must be positive",
		},
		{
			name:    "with description",
			req:     &CreateShareRequest{Name: "share-1", NetworkID: "net-1", Size: 10, Description: "desc"},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestUpdateShareRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *UpdateShareRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "name only",
			req:     &UpdateShareRequest{Name: "new-name"},
			wantErr: false,
		},
		{
			name:    "description only",
			req:     &UpdateShareRequest{Description: "new-desc"},
			wantErr: false,
		},
		{
			name:    "both fields",
			req:     &UpdateShareRequest{Name: "new-name", Description: "new-desc"},
			wantErr: false,
		},
		{
			name:    "empty request",
			req:     &UpdateShareRequest{},
			wantErr: true,
			errMsg:  "at least one field (name or description) must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestExtendShareRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *ExtendShareRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid",
			req:     &ExtendShareRequest{NewSize: 100},
			wantErr: false,
		},
		{
			name:    "zero size",
			req:     &ExtendShareRequest{NewSize: 0},
			wantErr: true,
			errMsg:  "new_size must be positive",
		},
		{
			name:    "negative size",
			req:     &ExtendShareRequest{NewSize: -1},
			wantErr: true,
			errMsg:  "new_size must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestCreateShareRuleRequest_Validate(t *testing.T) {
	tests := []struct {
		name    string
		req     *CreateShareRuleRequest
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid rw",
			req:     &CreateShareRuleRequest{CIDR: "10.0.0.0/24", Level: "rw"},
			wantErr: false,
		},
		{
			name:    "valid ro",
			req:     &CreateShareRuleRequest{CIDR: "0.0.0.0/0", Level: "ro"},
			wantErr: false,
		},
		{
			name:    "missing cidr",
			req:     &CreateShareRuleRequest{CIDR: "", Level: "rw"},
			wantErr: true,
			errMsg:  "cidr is required",
		},
		{
			name:    "missing level",
			req:     &CreateShareRuleRequest{CIDR: "10.0.0.0/24", Level: ""},
			wantErr: true,
			errMsg:  "level is required",
		},
		{
			name:    "invalid level",
			req:     &CreateShareRuleRequest{CIDR: "10.0.0.0/24", Level: "admin"},
			wantErr: true,
			errMsg:  "level must be one of rw, ro",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err != nil && tt.errMsg != "" && err.Error() != tt.errMsg {
				t.Errorf("Validate() error = %q, want %q", err.Error(), tt.errMsg)
			}
		})
	}
}

func TestShareListResponse_JSONMarshaling(t *testing.T) {
	resp := &ShareListResponse{
		Shares: []*Share{
			{ID: "share-1", Name: "s1", Size: 10, Status: ShareStatusAvailable},
			{ID: "share-2", Name: "s2", Size: 20, Status: ShareStatusCreating},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var unmarshaled ShareListResponse
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(unmarshaled.Shares) != 2 {
		t.Fatalf("expected 2 shares, got %d", len(unmarshaled.Shares))
	}
	if unmarshaled.Shares[0].ID != "share-1" {
		t.Errorf("expected share-1, got %s", unmarshaled.Shares[0].ID)
	}
}
