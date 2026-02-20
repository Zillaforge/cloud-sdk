package loadbalancers

import (
	"encoding/json"
	"testing"
)

func TestLoadBalancerJSONRoundtrip(t *testing.T) {
	original := `{
		"id": "lb-123",
		"name": "my-lb",
		"description": "test lb",
		"status": "ACTIVE",
		"vip": "192.168.1.100",
		"network_id": "net-456",
		"project_id": "proj-789"
	}`

	var lb LoadBalancer
	if err := json.Unmarshal([]byte(original), &lb); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if lb.ID != "lb-123" {
		t.Errorf("expected ID lb-123, got %s", lb.ID)
	}
	if lb.Name != "my-lb" {
		t.Errorf("expected Name my-lb, got %s", lb.Name)
	}
	if lb.Status != LoadBalancerStatusActive {
		t.Errorf("expected Status ACTIVE, got %s", lb.Status)
	}
	if lb.VIP != "192.168.1.100" {
		t.Errorf("expected VIP 192.168.1.100, got %s", lb.VIP)
	}

	// Re-marshal
	data, err := json.Marshal(&lb)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	var roundtrip LoadBalancer
	if err := json.Unmarshal(data, &roundtrip); err != nil {
		t.Fatalf("failed to unmarshal roundtrip: %v", err)
	}

	if roundtrip.ID != lb.ID || roundtrip.Name != lb.Name || roundtrip.Status != lb.Status {
		t.Error("roundtrip mismatch")
	}
}

func TestLoadBalancerCreateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     LoadBalancerCreateRequest
		wantErr string
	}{
		{
			name:    "valid request",
			req:     LoadBalancerCreateRequest{Name: "my-lb", NetworkID: "net-123"},
			wantErr: "",
		},
		{
			name:    "missing name",
			req:     LoadBalancerCreateRequest{NetworkID: "net-123"},
			wantErr: "name is required",
		},
		{
			name:    "missing network_id",
			req:     LoadBalancerCreateRequest{Name: "my-lb"},
			wantErr: "network_id is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}

func TestLoadBalancerUpdateRequestValidate(t *testing.T) {
	tests := []struct {
		name    string
		req     LoadBalancerUpdateRequest
		wantErr string
	}{
		{
			name:    "valid with name",
			req:     LoadBalancerUpdateRequest{Name: "new-name"},
			wantErr: "",
		},
		{
			name:    "valid with description",
			req:     LoadBalancerUpdateRequest{Description: "new desc"},
			wantErr: "",
		},
		{
			name:    "empty request",
			req:     LoadBalancerUpdateRequest{},
			wantErr: "at least one field (name or description) must be provided",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.req.Validate()
			if tt.wantErr == "" {
				if err != nil {
					t.Errorf("expected no error, got: %v", err)
				}
			} else {
				if err == nil {
					t.Errorf("expected error containing %q, got nil", tt.wantErr)
				} else if err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			}
		})
	}
}

func TestLoadBalancerStatusConstants(t *testing.T) {
	if LoadBalancerStatusActive.String() != "ACTIVE" {
		t.Errorf("expected ACTIVE, got %s", LoadBalancerStatusActive.String())
	}
	if LoadBalancerStatusError.String() != "ERROR" {
		t.Errorf("expected ERROR, got %s", LoadBalancerStatusError.String())
	}
}

func TestLoadBalancerListResponseJSON(t *testing.T) {
	jsonData := `{"load_balancers": [{"id": "lb-1", "name": "lb1", "status": "ACTIVE", "network_id": "net-1", "project_id": "proj-1"}]}`

	var response LoadBalancerListResponse
	if err := json.Unmarshal([]byte(jsonData), &response); err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	if len(response.LoadBalancers) != 1 {
		t.Fatalf("expected 1 LB, got %d", len(response.LoadBalancers))
	}
	if response.LoadBalancers[0].ID != "lb-1" {
		t.Errorf("expected ID lb-1, got %s", response.LoadBalancers[0].ID)
	}
}
