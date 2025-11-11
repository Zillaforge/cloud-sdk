package floatingips

import (
	"encoding/json"
	"testing"

	"github.com/Zillaforge/cloud-sdk/models/vps/common"
)

// TestFloatingIPStatusConstants verifies all status constants are defined correctly
func TestFloatingIPStatusConstants(t *testing.T) {
	tests := []struct {
		status   FloatingIPStatus
		expected string
	}{
		{FloatingIPStatusActive, "ACTIVE"},
		{FloatingIPStatusPending, "PENDING"},
		{FloatingIPStatusDown, "DOWN"},
		{FloatingIPStatusRejected, "REJECTED"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			if tt.status.String() != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, tt.status.String())
			}
		})
	}
}

// TestFloatingIPStatusValid verifies the Valid() method correctly identifies valid statuses
func TestFloatingIPStatusValid(t *testing.T) {
	tests := []struct {
		status   FloatingIPStatus
		isValid  bool
		testName string
	}{
		{FloatingIPStatusActive, true, "ACTIVE is valid"},
		{FloatingIPStatusPending, true, "PENDING is valid"},
		{FloatingIPStatusDown, true, "DOWN is valid"},
		{FloatingIPStatusRejected, true, "REJECTED is valid"},
		{FloatingIPStatus("INVALID"), false, "INVALID is not valid"},
		{FloatingIPStatus(""), false, "empty string is not valid"},
		{FloatingIPStatus("active"), false, "lowercase is not valid"},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			if tt.status.Valid() != tt.isValid {
				t.Errorf("expected Valid()=%v for %q", tt.isValid, tt.status)
			}
		})
	}
}

// TestFloatingIPUnmarshalingAllFields verifies FloatingIP unmarshals with all 19 fields
func TestFloatingIPUnmarshalingAllFields(t *testing.T) {
	jsonData := `{
		"id": "fip-123",
		"uuid": "550e8400-e29b-41d4-a716-446655440000",
		"name": "MyFloatingIP",
		"address": "203.0.113.42",
		"extnet_id": "net-ext-1",
		"port_id": "port-456",
		"project_id": "proj-789",
		"namespace": "prod",
		"user_id": "user-999",
		"user": {"id": "user-999", "name": "Alice"},
		"project": {"id": "proj-789", "name": "Project A"},
		"device_id": "server-111",
		"device_name": "web-server-01",
		"device_type": "server",
		"description": "Production Web IP",
		"status": "ACTIVE",
		"status_reason": "",
		"reserved": false,
		"createdAt": "2025-11-11T10:30:00Z",
		"updatedAt": "2025-11-11T11:00:00Z",
		"approvedAt": "2025-11-11T10:35:00Z"
	}`

	var fip FloatingIP
	err := json.Unmarshal([]byte(jsonData), &fip)
	if err != nil {
		t.Fatalf("failed to unmarshal FloatingIP: %v", err)
	}

	// Verify all fields are present and correctly typed
	if fip.ID != "fip-123" {
		t.Errorf("ID: expected %q, got %q", "fip-123", fip.ID)
	}
	if fip.UUID != "550e8400-e29b-41d4-a716-446655440000" {
		t.Errorf("UUID: expected %q, got %q", "550e8400-e29b-41d4-a716-446655440000", fip.UUID)
	}
	if fip.Name != "MyFloatingIP" {
		t.Errorf("Name: expected %q, got %q", "MyFloatingIP", fip.Name)
	}
	if fip.Address != "203.0.113.42" {
		t.Errorf("Address: expected %q, got %q", "203.0.113.42", fip.Address)
	}
	if fip.ExtNetID != "net-ext-1" {
		t.Errorf("ExtNetID: expected %q, got %q", "net-ext-1", fip.ExtNetID)
	}
	if fip.PortID != "port-456" {
		t.Errorf("PortID: expected %q, got %q", "port-456", fip.PortID)
	}
	if fip.ProjectID != "proj-789" {
		t.Errorf("ProjectID: expected %q, got %q", "proj-789", fip.ProjectID)
	}
	if fip.Namespace != "prod" {
		t.Errorf("Namespace: expected %q, got %q", "prod", fip.Namespace)
	}
	if fip.UserID != "user-999" {
		t.Errorf("UserID: expected %q, got %q", "user-999", fip.UserID)
	}
	if fip.User == nil || fip.User.ID != "user-999" || fip.User.Name != "Alice" {
		t.Errorf("User: expected IDName(ID=user-999, Name=Alice), got %+v", fip.User)
	}
	if fip.Project == nil || fip.Project.ID != "proj-789" || fip.Project.Name != "Project A" {
		t.Errorf("Project: expected IDName(ID=proj-789, Name=Project A), got %+v", fip.Project)
	}
	if fip.DeviceID != "server-111" {
		t.Errorf("DeviceID: expected %q, got %q", "server-111", fip.DeviceID)
	}
	if fip.DeviceName != "web-server-01" {
		t.Errorf("DeviceName: expected %q, got %q", "web-server-01", fip.DeviceName)
	}
	if fip.DeviceType != "server" {
		t.Errorf("DeviceType: expected %q, got %q", "server", fip.DeviceType)
	}
	if fip.Description != "Production Web IP" {
		t.Errorf("Description: expected %q, got %q", "Production Web IP", fip.Description)
	}
	if fip.Status != FloatingIPStatusActive {
		t.Errorf("Status: expected %v, got %v", FloatingIPStatusActive, fip.Status)
	}
	if fip.StatusReason != "" {
		t.Errorf("StatusReason: expected empty, got %q", fip.StatusReason)
	}
	if fip.Reserved != false {
		t.Errorf("Reserved: expected false, got %v", fip.Reserved)
	}
	if fip.CreatedAt != "2025-11-11T10:30:00Z" {
		t.Errorf("CreatedAt: expected %q, got %q", "2025-11-11T10:30:00Z", fip.CreatedAt)
	}
	if fip.UpdatedAt != "2025-11-11T11:00:00Z" {
		t.Errorf("UpdatedAt: expected %q, got %q", "2025-11-11T11:00:00Z", fip.UpdatedAt)
	}
	if fip.ApprovedAt != "2025-11-11T10:35:00Z" {
		t.Errorf("ApprovedAt: expected %q, got %q", "2025-11-11T10:35:00Z", fip.ApprovedAt)
	}
}

// TestFloatingIPMarshalingRoundtrip verifies FloatingIP marshaling preserves all field values
func TestFloatingIPMarshalingRoundtrip(t *testing.T) {
	original := &FloatingIP{
		ID:           "fip-123",
		UUID:         "550e8400-e29b-41d4-a716-446655440000",
		Name:         "MyFloatingIP",
		Address:      "203.0.113.42",
		ExtNetID:     "net-ext-1",
		PortID:       "port-456",
		ProjectID:    "proj-789",
		Namespace:    "prod",
		UserID:       "user-999",
		User:         &common.IDName{ID: "user-999", Name: "Alice"},
		Project:      &common.IDName{ID: "proj-789", Name: "Project A"},
		DeviceID:     "server-111",
		DeviceName:   "web-server-01",
		DeviceType:   "server",
		Description:  "Production Web IP",
		Status:       FloatingIPStatusActive,
		StatusReason: "",
		Reserved:     false,
		CreatedAt:    "2025-11-11T10:30:00Z",
		UpdatedAt:    "2025-11-11T11:00:00Z",
		ApprovedAt:   "2025-11-11T10:35:00Z",
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal FloatingIP: %v", err)
	}

	// Unmarshal back
	var result FloatingIP
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal FloatingIP: %v", err)
	}

	// Compare all fields
	if result.ID != original.ID {
		t.Errorf("ID roundtrip failed: expected %q, got %q", original.ID, result.ID)
	}
	if result.UUID != original.UUID {
		t.Errorf("UUID roundtrip failed: expected %q, got %q", original.UUID, result.UUID)
	}
	if result.Name != original.Name {
		t.Errorf("Name roundtrip failed: expected %q, got %q", original.Name, result.Name)
	}
	if result.Address != original.Address {
		t.Errorf("Address roundtrip failed: expected %q, got %q", original.Address, result.Address)
	}
	if result.Status != original.Status {
		t.Errorf("Status roundtrip failed: expected %v, got %v", original.Status, result.Status)
	}
	if result.CreatedAt != original.CreatedAt {
		t.Errorf("CreatedAt roundtrip failed: expected %q, got %q", original.CreatedAt, result.CreatedAt)
	}
}

// TestFloatingIPOptionalFieldsOmitted verifies optional fields are omitted when empty
func TestFloatingIPOptionalFieldsOmitted(t *testing.T) {
	jsonData := `{
		"id": "fip-minimal",
		"uuid": "550e8400-e29b-41d4-a716-446655440000",
		"name": "MinimalIP",
		"address": "203.0.113.99",
		"project_id": "proj-minimal",
		"user_id": "user-minimal",
		"status": "PENDING",
		"reserved": true,
		"createdAt": "2025-11-11T10:00:00Z"
	}`

	var fip FloatingIP
	err := json.Unmarshal([]byte(jsonData), &fip)
	if err != nil {
		t.Fatalf("failed to unmarshal FloatingIP: %v", err)
	}

	// Verify optional fields are empty/nil
	if fip.ExtNetID != "" {
		t.Errorf("ExtNetID should be empty, got %q", fip.ExtNetID)
	}
	if fip.PortID != "" {
		t.Errorf("PortID should be empty, got %q", fip.PortID)
	}
	if fip.Namespace != "" {
		t.Errorf("Namespace should be empty, got %q", fip.Namespace)
	}
	if fip.User != nil {
		t.Errorf("User should be nil, got %+v", fip.User)
	}
	if fip.Project != nil {
		t.Errorf("Project should be nil, got %+v", fip.Project)
	}
	if fip.DeviceID != "" {
		t.Errorf("DeviceID should be empty, got %q", fip.DeviceID)
	}
	if fip.DeviceName != "" {
		t.Errorf("DeviceName should be empty, got %q", fip.DeviceName)
	}
	if fip.DeviceType != "" {
		t.Errorf("DeviceType should be empty, got %q", fip.DeviceType)
	}
	if fip.Description != "" {
		t.Errorf("Description should be empty, got %q", fip.Description)
	}
	if fip.StatusReason != "" {
		t.Errorf("StatusReason should be empty, got %q", fip.StatusReason)
	}
	if fip.UpdatedAt != "" {
		t.Errorf("UpdatedAt should be empty, got %q", fip.UpdatedAt)
	}
	if fip.ApprovedAt != "" {
		t.Errorf("ApprovedAt should be empty, got %q", fip.ApprovedAt)
	}
}

// TestFloatingIPCamelCaseJSONTags verifies camelCase JSON field tags match API spec
func TestFloatingIPCamelCaseJSONTags(t *testing.T) {
	// Test that API response with camelCase timestamps unmarshals correctly
	jsonData := `{
		"id": "fip-123",
		"uuid": "550e8400-e29b-41d4-a716-446655440000",
		"name": "Test",
		"address": "203.0.113.42",
		"project_id": "proj-123",
		"user_id": "user-123",
		"status": "ACTIVE",
		"reserved": false,
		"createdAt": "2025-11-11T10:30:00Z",
		"updatedAt": "2025-11-11T11:00:00Z",
		"approvedAt": "2025-11-11T10:35:00Z"
	}`

	var fip FloatingIP
	err := json.Unmarshal([]byte(jsonData), &fip)
	if err != nil {
		t.Fatalf("failed to unmarshal with camelCase tags: %v", err)
	}

	if fip.CreatedAt != "2025-11-11T10:30:00Z" {
		t.Errorf("createdAt camelCase tag failed: expected %q, got %q", "2025-11-11T10:30:00Z", fip.CreatedAt)
	}
	if fip.UpdatedAt != "2025-11-11T11:00:00Z" {
		t.Errorf("updatedAt camelCase tag failed: expected %q, got %q", "2025-11-11T11:00:00Z", fip.UpdatedAt)
	}
	if fip.ApprovedAt != "2025-11-11T10:35:00Z" {
		t.Errorf("approvedAt camelCase tag failed: expected %q, got %q", "2025-11-11T10:35:00Z", fip.ApprovedAt)
	}
}

// TestFloatingIPIDNameNestedObjects verifies nested IDName objects marshal/unmarshal correctly
func TestFloatingIPIDNameNestedObjects(t *testing.T) {
	jsonData := `{
		"id": "fip-123",
		"uuid": "550e8400-e29b-41d4-a716-446655440000",
		"name": "Test",
		"address": "203.0.113.42",
		"project_id": "proj-789",
		"user_id": "user-999",
		"user": {"id": "user-999", "name": "Alice Smith"},
		"project": {"id": "proj-789", "name": "Main Project"},
		"status": "ACTIVE",
		"reserved": false,
		"createdAt": "2025-11-11T10:30:00Z"
	}`

	var fip FloatingIP
	err := json.Unmarshal([]byte(jsonData), &fip)
	if err != nil {
		t.Fatalf("failed to unmarshal FloatingIP with nested objects: %v", err)
	}

	// Verify nested User object
	if fip.User == nil {
		t.Fatal("User should not be nil")
	}
	if fip.User.ID != "user-999" {
		t.Errorf("User.ID: expected %q, got %q", "user-999", fip.User.ID)
	}
	if fip.User.Name != "Alice Smith" {
		t.Errorf("User.Name: expected %q, got %q", "Alice Smith", fip.User.Name)
	}

	// Verify nested Project object
	if fip.Project == nil {
		t.Fatal("Project should not be nil")
	}
	if fip.Project.ID != "proj-789" {
		t.Errorf("Project.ID: expected %q, got %q", "proj-789", fip.Project.ID)
	}
	if fip.Project.Name != "Main Project" {
		t.Errorf("Project.Name: expected %q, got %q", "Main Project", fip.Project.Name)
	}

	// Verify roundtrip
	data, err := json.Marshal(fip)
	if err != nil {
		t.Fatalf("failed to marshal FloatingIP: %v", err)
	}

	var result FloatingIP
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal after marshaling: %v", err)
	}

	if result.User.ID != "user-999" || result.User.Name != "Alice Smith" {
		t.Errorf("User roundtrip failed: expected (user-999, Alice Smith), got (%q, %q)", result.User.ID, result.User.Name)
	}
	if result.Project.ID != "proj-789" || result.Project.Name != "Main Project" {
		t.Errorf("Project roundtrip failed: expected (proj-789, Main Project), got (%q, %q)", result.Project.ID, result.Project.Name)
	}
}

// TestFloatingIPUpdateRequestMarshaling verifies JSON marshaling/unmarshaling of FloatingIPUpdateRequest
func TestFloatingIPUpdateRequestMarshaling(t *testing.T) {
	tests := []struct {
		name     string
		request  FloatingIPUpdateRequest
		expected string
	}{
		{
			name: "all fields set",
			request: FloatingIPUpdateRequest{
				Name:        "updated-fip",
				Description: "Updated floating IP",
				Reserved:    &[]bool{true}[0],
			},
			expected: `{"name":"updated-fip","description":"Updated floating IP","reserved":true}`,
		},
		{
			name: "reserved false",
			request: FloatingIPUpdateRequest{
				Name:        "fip-name",
				Description: "Description",
				Reserved:    &[]bool{false}[0],
			},
			expected: `{"name":"fip-name","description":"Description","reserved":false}`,
		},
		{
			name: "reserved nil (omitted)",
			request: FloatingIPUpdateRequest{
				Name:        "fip-name",
				Description: "Description",
				Reserved:    nil,
			},
			expected: `{"name":"fip-name","description":"Description"}`,
		},
		{
			name: "only reserved set",
			request: FloatingIPUpdateRequest{
				Reserved: &[]bool{true}[0],
			},
			expected: `{"reserved":true}`,
		},
		{
			name:     "empty request",
			request:  FloatingIPUpdateRequest{},
			expected: `{}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test marshaling
			data, err := json.Marshal(tt.request)
			if err != nil {
				t.Fatalf("failed to marshal FloatingIPUpdateRequest: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("marshaling failed: expected %q, got %q", tt.expected, string(data))
			}

			// Test unmarshaling
			var result FloatingIPUpdateRequest
			err = json.Unmarshal(data, &result)
			if err != nil {
				t.Fatalf("failed to unmarshal FloatingIPUpdateRequest: %v", err)
			}

			// Verify fields
			if result.Name != tt.request.Name {
				t.Errorf("Name: expected %q, got %q", tt.request.Name, result.Name)
			}
			if result.Description != tt.request.Description {
				t.Errorf("Description: expected %q, got %q", tt.request.Description, result.Description)
			}
			if tt.request.Reserved == nil && result.Reserved != nil {
				t.Errorf("Reserved: expected nil, got %v", result.Reserved)
			}
			if tt.request.Reserved != nil && result.Reserved == nil {
				t.Errorf("Reserved: expected %v, got nil", *tt.request.Reserved)
			}
			if tt.request.Reserved != nil && result.Reserved != nil && *result.Reserved != *tt.request.Reserved {
				t.Errorf("Reserved: expected %v, got %v", *tt.request.Reserved, *result.Reserved)
			}
		})
	}
}

// TestFloatingIPUpdateRequestRoundtrip verifies that marshaling and unmarshaling preserves data
func TestFloatingIPUpdateRequestRoundtrip(t *testing.T) {
	original := FloatingIPUpdateRequest{
		Name:        "test-fip",
		Description: "Test description",
		Reserved:    &[]bool{true}[0],
	}

	// Marshal to JSON
	data, err := json.Marshal(original)
	if err != nil {
		t.Fatalf("failed to marshal: %v", err)
	}

	// Unmarshal back
	var result FloatingIPUpdateRequest
	err = json.Unmarshal(data, &result)
	if err != nil {
		t.Fatalf("failed to unmarshal: %v", err)
	}

	// Verify all fields match
	if result.Name != original.Name {
		t.Errorf("Name mismatch: expected %q, got %q", original.Name, result.Name)
	}
	if result.Description != original.Description {
		t.Errorf("Description mismatch: expected %q, got %q", original.Description, result.Description)
	}
	if result.Reserved == nil || *result.Reserved != *original.Reserved {
		t.Errorf("Reserved mismatch: expected %v, got %v", *original.Reserved, result.Reserved)
	}
}
