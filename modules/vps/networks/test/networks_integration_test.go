package networks_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sync/atomic"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// TestNetworkLifecycle verifies the complete network lifecycle
// Create → Get → Update → List Ports → Delete
func TestNetworkLifecycle(t *testing.T) {
	var requestCount int32
	createdNetworkID := "net-full-001"
	fixture := loadFixtureBytes(t, "network_full.json")

	// Create mock HTTP server that handles all network operations
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&requestCount, 1)
		w.Header().Set("Content-Type", "application/json")

		// Route based on method and path
		switch {
		// CREATE: POST /api/v1/project/{project-id}/networks
		case r.Method == http.MethodPost && r.URL.Path == "/vps/api/v1/project/proj-123/networks":
			// Parse create request
			var createReq map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&createReq)

			// Validate request
			if createReq["name"] == "" {
				w.WriteHeader(http.StatusBadRequest)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"error_code": "VALIDATION_ERROR",
					"message":    "Name is required",
				})
				return
			}

			var payload map[string]interface{}
			_ = json.Unmarshal(fixture, &payload)
			payload["id"] = createdNetworkID
			if name, ok := createReq["name"]; ok {
				payload["name"] = name
			}
			if desc, ok := createReq["description"]; ok {
				payload["description"] = desc
			}
			if cidr, ok := createReq["cidr"]; ok {
				payload["cidr"] = cidr
			}

			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(payload)

		// GET: GET /api/v1/project/{project-id}/networks/{net-id}
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-123/networks/"+createdNetworkID:
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write(fixture)

		// UPDATE: PUT /api/v1/project/{project-id}/networks/{net-id}
		case r.Method == http.MethodPut && r.URL.Path == "/vps/api/v1/project/proj-123/networks/"+createdNetworkID:
			// Parse update request
			var updateReq map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&updateReq)

			var payload map[string]interface{}
			_ = json.Unmarshal(fixture, &payload)
			payload["id"] = createdNetworkID
			if name, ok := updateReq["name"]; ok {
				payload["name"] = name
			}
			if desc, ok := updateReq["description"]; ok {
				payload["description"] = desc
			}
			payload["updatedAt"] = "2025-01-07T00:00:00Z"

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(payload)

		// LIST PORTS: GET /api/v1/project/{project-id}/networks/{net-id}/ports
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-123/networks/"+createdNetworkID+"/ports":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{
					"id":        "port-1",
					"addresses": []string{"10.10.0.10"},
					"server": map[string]interface{}{
						"id":         "srv-1",
						"name":       "integration-node-1",
						"status":     "ACTIVE",
						"project_id": "proj-123",
						"user_id":    "user-567",
					},
				},
				{
					"id":        "port-2",
					"addresses": []string{"10.10.0.20"},
				},
			})

		// DELETE: DELETE /api/v1/project/{project-id}/networks/{net-id}
		case r.Method == http.MethodDelete && r.URL.Path == "/vps/api/v1/project/proj-123/networks/"+createdNetworkID:
			w.WriteHeader(http.StatusNoContent)

		default:
			t.Logf("Unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"error_code": "NOT_FOUND",
				"message":    "Endpoint not found",
			})
		}
	}))
	defer server.Close()

	// Create client
	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()
	ctx := context.Background()

	// Step 1: Create network
	t.Log("Step 1: Creating network...")
	createReq := &networks.NetworkCreateRequest{
		Name:        "integration-test-network",
		Description: "Network for integration testing",
		CIDR:        "10.10.0.0/24",
	}
	network, err := vpsClient.Networks().Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Step 1 failed - create network: %v", err)
	}
	assertStringField(t, network.Network, "ID", createdNetworkID)
	assertStringField(t, network.Network, "Name", "integration-test-network")
	assertStringField(t, network.Network, "Gateway", "10.42.0.1")
	assertStringSliceField(t, network.Network, "Nameservers", []string{"1.1.1.1", "8.8.8.8"})
	project := requirePointerStructField(t, network.Network, "Project")
	assertStringField(t, project.Interface(), "ID", "proj-001")
	router := requirePointerStructField(t, network.Network, "Router")
	assertStringField(t, router.Interface(), "ID", "router-123")
	t.Logf("✓ Network created with ID: %s", network.Network.ID)

	// Step 2: Get network
	t.Log("Step 2: Retrieving network...")
	resource, err := vpsClient.Networks().Get(ctx, createdNetworkID)
	if err != nil {
		t.Fatalf("Step 2 failed - get network: %v", err)
	}
	assertStringField(t, resource.Network, "ID", createdNetworkID)
	assertStringField(t, resource.Network, "CIDR", "10.42.0.0/24")
	assertBoolField(t, resource.Network, "Shared", true)
	assertStringField(t, resource.Network, "Status", "ACTIVE")
	t.Logf("✓ Network retrieved: %s (%s)", resource.Network.Name, resource.Network.CIDR)

	// Step 3: Update network
	t.Log("Step 3: Updating network...")
	updateReq := &networks.NetworkUpdateRequest{
		Name:        "updated-integration-network",
		Description: "Updated description",
	}
	updatedNetwork, err := vpsClient.Networks().Update(ctx, createdNetworkID, updateReq)
	if err != nil {
		t.Fatalf("Step 3 failed - update network: %v", err)
	}
	assertStringField(t, updatedNetwork.Network, "Name", "updated-integration-network")
	assertStringField(t, updatedNetwork.Network, "Description", "Updated description")
	assertStringField(t, updatedNetwork.Network, "RouterID", "router-123")
	t.Logf("✓ Network updated: %s", updatedNetwork.Network.Name)

	// Step 4: List ports on network
	t.Log("Step 4: Listing network ports...")
	ports, err := resource.Ports().List(ctx)
	if err != nil {
		t.Fatalf("Step 4 failed - list ports: %v", err)
	}
	if len(ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(ports))
	}
	assertStringField(t, ports[0], "ID", "port-1")
	assertStringSliceField(t, ports[0], "Addresses", []string{"10.10.0.10"})
	serverSummary := requirePointerStructField(t, ports[0], "Server")
	assertStringField(t, serverSummary.Interface(), "Status", "ACTIVE")

	assertStringField(t, ports[1], "ID", "port-2")
	assertStringSliceField(t, ports[1], "Addresses", []string{"10.10.0.20"})
	serverField := requireStructField(t, ports[1], "Server")
	if serverField.Kind() != reflect.Ptr {
		t.Fatalf("expected pointer server field, got %s", serverField.Kind())
	}
	if !serverField.IsNil() {
		t.Fatalf("expected nil server pointer for second port, got %#v", serverField.Interface())
	}
	t.Logf("✓ Found %d ports on network", len(ports))

	// Step 5: Delete network
	t.Log("Step 5: Deleting network...")
	err = vpsClient.Networks().Delete(ctx, createdNetworkID)
	if err != nil {
		t.Fatalf("Step 5 failed - delete network: %v", err)
	}
	t.Logf("✓ Network deleted successfully")

	// Verify all expected requests were made
	finalCount := atomic.LoadInt32(&requestCount)
	expectedRequests := int32(5) // Create, Get, Update, ListPorts, Delete
	if finalCount != expectedRequests {
		t.Errorf("expected %d requests, got %d", expectedRequests, finalCount)
	}

	t.Log("✅ Network lifecycle test completed successfully")
}

// TestNetworkLifecycle_ErrorHandling verifies error handling in lifecycle
func TestNetworkLifecycle_ErrorHandling(t *testing.T) {
	// Test creating network with invalid data
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error_code": "VALIDATION_ERROR",
			"message":    "Name is required",
		})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()
	ctx := context.Background()

	// Attempt to create with invalid request
	createReq := &networks.NetworkCreateRequest{
		Name: "", // Invalid empty name
		CIDR: "10.0.0.0/24",
	}
	network, err := vpsClient.Networks().Create(ctx, createReq)

	// Verify error
	if err == nil {
		t.Fatal("expected error for invalid create request, got nil")
	}
	if network != nil {
		t.Errorf("expected nil network on error, got %+v", network)
	}

	sdkErr, ok := err.(*cloudsdk.SDKError)
	if !ok {
		t.Fatalf("expected *cloudsdk.SDKError, got %T", err)
	}
	if sdkErr.StatusCode != http.StatusBadRequest {
		t.Errorf("expected status 400, got %d", sdkErr.StatusCode)
	}
	if sdkErr.Message != "Name is required" {
		t.Errorf("expected message 'Name is required', got '%s'", sdkErr.Message)
	}

	t.Log("✅ Error handling test completed successfully")
}
