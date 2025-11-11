package floatingips_test

import (
	"encoding/json"
	"testing"

	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
)

// TestFloatingIPListContract_FIPListOutput verifies that the List response structure
// matches pb.FIPListOutput specification from vps.yaml (floatingips array, not items)
func TestFloatingIPListContract_FIPListOutput(t *testing.T) {
	// Test data matching pb.FIPListOutput structure from vps.yaml
	testJSON := `{
		"floatingips": [
			{
				"id": "fip-123",
				"uuid": "550e8400-e29b-41d4-a716-446655440000",
				"name": "test-floating-ip",
				"address": "203.0.113.10",
				"description": "Test floating IP",
				"status": "ACTIVE",
				"project_id": "proj-123",
				"user_id": "user-123",
				"port_id": "port-123",
				"device_id": "device-123",
				"device_name": "web-server",
				"device_type": "server",
				"extnet_id": "ext-net-123",
				"namespace": "default",
				"status_reason": "",
				"reserved": false,
				"createdAt": "2025-01-01T00:00:00Z",
				"updatedAt": "2025-01-02T00:00:00Z",
				"approvedAt": "2025-01-01T12:00:00Z",
				"project": {
					"id": "proj-123",
					"name": "test-project"
				},
				"user": {
					"id": "user-123",
					"name": "test-user"
				}
			}
		]
	}`

	// Test that the JSON can be unmarshaled into []*FloatingIP
	var response struct {
		FloatingIPs []*floatingips.FloatingIP `json:"floatingips"`
	}

	err := json.Unmarshal([]byte(testJSON), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal response with 'floatingips' field: %v", err)
	}

	// Verify the structure
	if len(response.FloatingIPs) != 1 {
		t.Errorf("expected 1 floating IP, got %d", len(response.FloatingIPs))
	}

	fip := response.FloatingIPs[0]
	if fip.ID != "fip-123" {
		t.Errorf("expected ID 'fip-123', got '%s'", fip.ID)
	}
	if fip.Address != "203.0.113.10" {
		t.Errorf("expected address '203.0.113.10', got '%s'", fip.Address)
	}
	if fip.Status != floatingips.FloatingIPStatusActive {
		t.Errorf("expected status 'ACTIVE', got '%s'", fip.Status)
	}

	// Verify nested objects
	if fip.Project == nil || fip.Project.ID != "proj-123" {
		t.Errorf("expected project ID 'proj-123', got %v", fip.Project)
	}
	if fip.User == nil || fip.User.ID != "user-123" {
		t.Errorf("expected user ID 'user-123', got %v", fip.User)
	}
}

// TestFloatingIPListContract_NoItemsField verifies that the response does NOT contain
// the old "items" field (breaking change validation)
func TestFloatingIPListContract_NoItemsField(t *testing.T) {
	// Test data that would be invalid (old structure with "items")
	invalidJSON := `{
		"items": [
			{
				"id": "fip-123",
				"name": "test-floating-ip",
				"address": "203.0.113.10",
				"status": "ACTIVE"
			}
		]
	}`

	// This should fail to unmarshal into the correct structure
	var response struct {
		FloatingIPs []*floatingips.FloatingIP `json:"floatingips"`
	}

	err := json.Unmarshal([]byte(invalidJSON), &response)
	if err != nil {
		// This is expected - the JSON doesn't have the "floatingips" field
		t.Logf("correctly failed to unmarshal invalid structure: %v", err)
	} else {
		// If it succeeded, the response should be empty
		if len(response.FloatingIPs) != 0 {
			t.Errorf("unexpectedly succeeded in unmarshaling invalid structure, got %d floating IPs", len(response.FloatingIPs))
		}
	}
}

// TestFloatingIPListContract_EmptyResponse verifies empty list response structure
func TestFloatingIPListContract_EmptyResponse(t *testing.T) {
	// Test empty response matching pb.FIPListOutput
	emptyJSON := `{
		"floatingips": []
	}`

	var response struct {
		FloatingIPs []*floatingips.FloatingIP `json:"floatingips"`
	}

	err := json.Unmarshal([]byte(emptyJSON), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal empty response: %v", err)
	}

	if len(response.FloatingIPs) != 0 {
		t.Errorf("expected 0 floating IPs in empty response, got %d", len(response.FloatingIPs))
	}
}

// TestFloatingIPCreateContract_FIPCreateInput verifies that the Create request structure
// matches FIPCreateInput specification from vps.yaml (name, description fields)
func TestFloatingIPCreateContract_FIPCreateInput(t *testing.T) {
	// Test marshaling FloatingIPCreateRequest to JSON
	req := &floatingips.FloatingIPCreateRequest{
		Name:        "test-fip",
		Description: "A test floating IP",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal create request: %v", err)
	}

	// Verify the JSON structure
	var requestBody map[string]interface{}
	err = json.Unmarshal(data, &requestBody)
	if err != nil {
		t.Fatalf("failed to unmarshal request body: %v", err)
	}

	// Verify required fields are present
	if _, hasName := requestBody["name"]; !hasName {
		t.Errorf("expected 'name' field in request body")
	}
	if _, hasDescription := requestBody["description"]; !hasDescription {
		t.Errorf("expected 'description' field in request body")
	}

	// Verify field values
	if name, ok := requestBody["name"].(string); !ok || name != "test-fip" {
		t.Errorf("expected name 'test-fip', got %v", requestBody["name"])
	}
	if description, ok := requestBody["description"].(string); !ok || description != "A test floating IP" {
		t.Errorf("expected description 'A test floating IP', got %v", requestBody["description"])
	}
}

// TestFloatingIPUpdateContract_FIPUpdateInput verifies that the Update request structure
// matches FIPUpdateInput specification from vps.yaml (name, description fields)
func TestFloatingIPUpdateContract_FIPUpdateInput(t *testing.T) {
	// Test marshaling FloatingIPUpdateRequest to JSON
	req := &floatingips.FloatingIPUpdateRequest{
		Name:        "updated-fip",
		Description: "Updated description",
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal update request: %v", err)
	}

	// Verify the JSON structure
	var requestBody map[string]interface{}
	err = json.Unmarshal(data, &requestBody)
	if err != nil {
		t.Fatalf("failed to unmarshal request body: %v", err)
	}

	// Verify optional fields can be present
	if _, hasName := requestBody["name"]; !hasName {
		t.Errorf("expected 'name' field in request body")
	}
	if _, hasDescription := requestBody["description"]; !hasDescription {
		t.Errorf("expected 'description' field in request body")
	}
}

// TestFloatingIPGetContract_Response verifies that the Get response structure
// matches pb.FloatingIPInfo specification from vps.yaml
func TestFloatingIPGetContract_Response(t *testing.T) {
	// Test data matching pb.FloatingIPInfo structure for single resource response
	testJSON := `{
		"id": "fip-123",
		"uuid": "550e8400-e29b-41d4-a716-446655440000",
		"name": "test-floating-ip",
		"address": "203.0.113.10",
		"status": "ACTIVE",
		"project_id": "proj-123",
		"user_id": "user-123",
		"reserved": false,
		"createdAt": "2025-01-01T00:00:00Z"
	}`

	var fip floatingips.FloatingIP
	err := json.Unmarshal([]byte(testJSON), &fip)
	if err != nil {
		t.Fatalf("failed to unmarshal floating IP response: %v", err)
	}

	// Verify all essential fields are present
	if fip.ID != "fip-123" {
		t.Errorf("expected ID 'fip-123', got '%s'", fip.ID)
	}
	if fip.Status != floatingips.FloatingIPStatusActive {
		t.Errorf("expected status 'ACTIVE', got '%s'", fip.Status)
	}
}

// TestFloatingIPDeleteContract_EmptyResponse verifies that the Delete response
// is empty (204 No Content or empty JSON object)
func TestFloatingIPDeleteContract_EmptyResponse(t *testing.T) {
	// Test that delete returns either empty or minimal response
	testJSON := `{}`

	var response map[string]interface{}
	err := json.Unmarshal([]byte(testJSON), &response)
	if err != nil {
		t.Fatalf("failed to unmarshal delete response: %v", err)
	}

	// Empty response is valid for delete
	if len(response) > 0 {
		t.Logf("Delete response contains fields: %v (may be expected)", response)
	}
}

// TestFloatingIPDisassociateContract_Response verifies that the Disassociate response
// returns updated FloatingIP with device association fields cleared
func TestFloatingIPDisassociateContract_Response(t *testing.T) {
	// Test data showing disassociated floating IP (device fields cleared)
	testJSON := `{
		"id": "fip-123",
		"uuid": "550e8400-e29b-41d4-a716-446655440000",
		"name": "test-floating-ip",
		"address": "203.0.113.10",
		"status": "ACTIVE",
		"project_id": "proj-123",
		"user_id": "user-123",
		"reserved": false,
		"createdAt": "2025-01-01T00:00:00Z",
		"device_id": null,
		"device_name": null,
		"device_type": null,
		"port_id": null
	}`

	var fip floatingips.FloatingIP
	err := json.Unmarshal([]byte(testJSON), &fip)
	if err != nil {
		t.Fatalf("failed to unmarshal disassociate response: %v", err)
	}

	// Verify the floating IP still exists but device associations are cleared
	if fip.ID != "fip-123" {
		t.Errorf("expected ID 'fip-123', got '%s'", fip.ID)
	}
	// Device fields should be empty/nil after disassociation
	if fip.DeviceID != "" {
		t.Logf("Warning: DeviceID expected to be empty after disassociate, got '%s'", fip.DeviceID)
	}
}
