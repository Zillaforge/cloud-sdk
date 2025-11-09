package securitygroups

import (
	"encoding/json"
	"testing"
)

// TestSecurityGroupListResponseUnmarshaling tests that SecurityGroupListResponse unmarshals correctly from API JSON.
func TestSecurityGroupListResponseUnmarshaling(t *testing.T) {
	// This test should initially fail if Total field exists, then pass after removal
	jsonData := `{
		"security_groups": [
			{
				"id": "sg-123",
				"name": "test-sg",
				"description": "Test security group",
				"project_id": "proj-456",
				"user_id": "user-789",
				"namespace": "default",
				"rules": [],
				"createdAt": "2025-11-09T10:00:00Z",
				"updatedAt": "2025-11-09T10:00:00Z"
			}
		]
	}`

	var resp SecurityGroupListResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if len(resp.SecurityGroups) != 1 {
		t.Errorf("len(SecurityGroups) = %d, want 1", len(resp.SecurityGroups))
	}

	if resp.SecurityGroups[0].ID != "sg-123" {
		t.Errorf("SecurityGroups[0].ID = %s, want sg-123", resp.SecurityGroups[0].ID)
	}

	// Verify we can get count from len() instead of Total field
	count := len(resp.SecurityGroups)
	if count != 1 {
		t.Errorf("count = %d, want 1", count)
	}
}

// TestSecurityGroupListResponseEmptyArray tests unmarshaling empty security_groups array.
func TestSecurityGroupListResponseEmptyArray(t *testing.T) {
	jsonData := `{"security_groups": []}`

	var resp SecurityGroupListResponse
	err := json.Unmarshal([]byte(jsonData), &resp)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	if len(resp.SecurityGroups) != 0 {
		t.Errorf("len(SecurityGroups) = %d, want 0", len(resp.SecurityGroups))
	}

	// Verify len() works for empty array
	count := len(resp.SecurityGroups)
	if count != 0 {
		t.Errorf("count = %d, want 0", count)
	}
}

// TestSecurityGroupListResponseMarshaling tests that SecurityGroupListResponse marshals correctly to JSON.
func TestSecurityGroupListResponseMarshaling(t *testing.T) {
	resp := SecurityGroupListResponse{
		SecurityGroups: []SecurityGroup{
			{
				ID:          "sg-123",
				Name:        "test-sg",
				Description: "Test",
				ProjectID:   "proj-456",
				UserID:      "user-789",
				Namespace:   "default",
				Rules:       []SecurityGroupRule{},
				CreatedAt:   "2025-11-09T10:00:00Z",
				UpdatedAt:   "2025-11-09T10:00:00Z",
			},
		},
	}

	data, err := json.Marshal(resp)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Verify the JSON contains security_groups but NOT total field
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map error = %v", err)
	}

	if _, exists := raw["total"]; exists {
		t.Error("JSON output should not contain 'total' field")
	}

	if _, exists := raw["security_groups"]; !exists {
		t.Error("JSON output should contain 'security_groups' field")
	}
}

// TestSecurityGroupFieldsMatchAPISpec verifies all SecurityGroup fields match pb.SgInfo from swagger/vps.yaml.
func TestSecurityGroupFieldsMatchAPISpec(t *testing.T) {
	// This JSON represents a complete pb.SgInfo response
	jsonData := `{
		"id": "sg-abc123",
		"name": "web-servers",
		"description": "Security group for web servers",
		"project_id": "proj-123",
		"user_id": "user-456",
		"namespace": "default",
		"rules": [
			{
				"id": "rule-001",
				"direction": "ingress",
				"protocol": "tcp",
				"port_min": 80,
				"port_max": 80,
				"remote_cidr": "0.0.0.0/0"
			}
		],
		"createdAt": "2025-11-09T10:00:00Z",
		"updatedAt": "2025-11-09T11:00:00Z",
		"project": {
			"id": "proj-123",
			"name": "production"
		},
		"user": {
			"id": "user-456",
			"name": "admin@example.com"
		}
	}`

	var sg SecurityGroup
	err := json.Unmarshal([]byte(jsonData), &sg)
	if err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}

	// Verify all fields are accessible and correct
	if sg.ID != "sg-abc123" {
		t.Errorf("ID = %s, want sg-abc123", sg.ID)
	}
	if sg.Name != "web-servers" {
		t.Errorf("Name = %s, want web-servers", sg.Name)
	}
	if sg.Description != "Security group for web servers" {
		t.Errorf("Description = %s, want 'Security group for web servers'", sg.Description)
	}
	if sg.ProjectID != "proj-123" {
		t.Errorf("ProjectID = %s, want proj-123", sg.ProjectID)
	}
	if sg.UserID != "user-456" {
		t.Errorf("UserID = %s, want user-456", sg.UserID)
	}
	if sg.Namespace != "default" {
		t.Errorf("Namespace = %s, want default", sg.Namespace)
	}
	if len(sg.Rules) != 1 {
		t.Errorf("len(Rules) = %d, want 1", len(sg.Rules))
	}
	if sg.CreatedAt != "2025-11-09T10:00:00Z" {
		t.Errorf("CreatedAt = %s, want 2025-11-09T10:00:00Z", sg.CreatedAt)
	}
	if sg.UpdatedAt != "2025-11-09T11:00:00Z" {
		t.Errorf("UpdatedAt = %s, want 2025-11-09T11:00:00Z", sg.UpdatedAt)
	}
	if sg.Project == nil {
		t.Error("Project should not be nil")
	} else {
		if sg.Project.ID != "proj-123" {
			t.Errorf("Project.ID = %s, want proj-123", sg.Project.ID)
		}
		if sg.Project.Name != "production" {
			t.Errorf("Project.Name = %s, want production", sg.Project.Name)
		}
	}
	if sg.User == nil {
		t.Error("User should not be nil")
	} else {
		if sg.User.ID != "user-456" {
			t.Errorf("User.ID = %s, want user-456", sg.User.ID)
		}
		if sg.User.Name != "admin@example.com" {
			t.Errorf("User.Name = %s, want admin@example.com", sg.User.Name)
		}
	}
}

// TestSecurityGroupCreateRequestJSONMarshaling tests SecurityGroupCreateRequest marshals correctly.
func TestSecurityGroupCreateRequestJSONMarshaling(t *testing.T) {
	req := SecurityGroupCreateRequest{
		Name:        "database-servers",
		Description: "Security group for PostgreSQL databases",
		Rules: []SecurityGroupRuleCreateRequest{
			{
				Direction:  DirectionIngress,
				Protocol:   ProtocolTCP,
				PortMin:    intPtr(5432),
				PortMax:    intPtr(5432),
				RemoteCIDR: "10.0.0.0/16",
			},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("Marshal() error = %v", err)
	}

	// Verify JSON structure matches API spec
	var raw map[string]interface{}
	if err := json.Unmarshal(data, &raw); err != nil {
		t.Fatalf("Unmarshal to map error = %v", err)
	}

	// Verify required fields present
	if _, exists := raw["name"]; !exists {
		t.Error("JSON should contain 'name' field")
	}
	if _, exists := raw["description"]; !exists {
		t.Error("JSON should contain 'description' field")
	}
	if _, exists := raw["rules"]; !exists {
		t.Error("JSON should contain 'rules' field")
	}

	// Unmarshal back to verify round-trip
	var got SecurityGroupCreateRequest
	if err := json.Unmarshal(data, &got); err != nil {
		t.Fatalf("Round-trip unmarshal error = %v", err)
	}

	if got.Name != req.Name {
		t.Errorf("Name = %s, want %s", got.Name, req.Name)
	}
	if got.Description != req.Description {
		t.Errorf("Description = %s, want %s", got.Description, req.Description)
	}
	if len(got.Rules) != 1 {
		t.Errorf("len(Rules) = %d, want 1", len(got.Rules))
	}
}

// TestSecurityGroupUpdateRequestJSONMarshaling tests SecurityGroupUpdateRequest with optional fields.
func TestSecurityGroupUpdateRequestJSONMarshaling(t *testing.T) {
	tests := []struct {
		name         string
		req          SecurityGroupUpdateRequest
		wantFields   []string
		unwantFields []string
	}{
		{
			name: "Both fields set",
			req: SecurityGroupUpdateRequest{
				Name:        strPtr("web-servers-v2"),
				Description: strPtr("Updated description"),
			},
			wantFields:   []string{"name", "description"},
			unwantFields: []string{},
		},
		{
			name: "Only name set",
			req: SecurityGroupUpdateRequest{
				Name: strPtr("web-servers-v3"),
			},
			wantFields:   []string{"name"},
			unwantFields: []string{"description"},
		},
		{
			name: "Only description set",
			req: SecurityGroupUpdateRequest{
				Description: strPtr("New description"),
			},
			wantFields:   []string{"description"},
			unwantFields: []string{"name"},
		},
		{
			name:         "No fields set (empty update)",
			req:          SecurityGroupUpdateRequest{},
			wantFields:   []string{},
			unwantFields: []string{"name", "description"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("Marshal() error = %v", err)
			}

			var raw map[string]interface{}
			if err := json.Unmarshal(data, &raw); err != nil {
				t.Fatalf("Unmarshal to map error = %v", err)
			}

			// Verify expected fields present
			for _, field := range tt.wantFields {
				if _, exists := raw[field]; !exists {
					t.Errorf("JSON should contain '%s' field", field)
				}
			}

			// Verify unwanted fields absent (due to omitempty)
			for _, field := range tt.unwantFields {
				if _, exists := raw[field]; exists {
					t.Errorf("JSON should NOT contain '%s' field (omitempty)", field)
				}
			}
		})
	}
}

// Helper functions for pointer types
func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
