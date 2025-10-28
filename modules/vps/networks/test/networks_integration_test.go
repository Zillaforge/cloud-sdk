package networks_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// TestNetworkLifecycle verifies the complete network lifecycle
// Create → Get → Update → List Ports → Delete
func TestNetworkLifecycle(t *testing.T) {
	var requestCount int32
	createdNetworkID := "net-lifecycle-123"

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

			// Return created network
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdNetworkID,
				"name":        createReq["name"],
				"description": createReq["description"],
				"cidr":        createReq["cidr"],
				"project_id":  "proj-123",
				"created_at":  "2025-01-01T00:00:00Z",
				"updated_at":  "2025-01-01T00:00:00Z",
			})

		// GET: GET /api/v1/project/{project-id}/networks/{net-id}
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-123/networks/"+createdNetworkID:
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdNetworkID,
				"name":        "integration-test-network",
				"description": "Network for integration testing",
				"cidr":        "10.10.0.0/24",
				"project_id":  "proj-123",
				"created_at":  "2025-01-01T00:00:00Z",
				"updated_at":  "2025-01-01T00:00:00Z",
			})

		// UPDATE: PUT /api/v1/project/{project-id}/networks/{net-id}
		case r.Method == http.MethodPut && r.URL.Path == "/vps/api/v1/project/proj-123/networks/"+createdNetworkID:
			// Parse update request
			var updateReq map[string]interface{}
			_ = json.NewDecoder(r.Body).Decode(&updateReq)

			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdNetworkID,
				"name":        updateReq["name"],
				"description": updateReq["description"],
				"cidr":        "10.10.0.0/24",
				"project_id":  "proj-123",
				"created_at":  "2025-01-01T00:00:00Z",
				"updated_at":  "2025-01-02T00:00:00Z",
			})

		// LIST PORTS: GET /api/v1/project/{project-id}/networks/{net-id}/ports
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-123/networks/"+createdNetworkID+"/ports":
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{
					"id":          "port-1",
					"network_id":  createdNetworkID,
					"fixed_ips":   []string{"10.10.0.10"},
					"mac_address": "fa:16:3e:11:22:33",
					"server_id":   "srv-1",
				},
				{
					"id":          "port-2",
					"network_id":  createdNetworkID,
					"fixed_ips":   []string{"10.10.0.20"},
					"mac_address": "fa:16:3e:44:55:66",
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
	if network.ID != createdNetworkID {
		t.Errorf("expected network ID '%s', got '%s'", createdNetworkID, network.ID)
	}
	if network.Name != "integration-test-network" {
		t.Errorf("expected network name 'integration-test-network', got '%s'", network.Name)
	}
	t.Logf("✓ Network created with ID: %s", network.ID)

	// Step 2: Get network
	t.Log("Step 2: Retrieving network...")
	resource, err := vpsClient.Networks().Get(ctx, createdNetworkID)
	if err != nil {
		t.Fatalf("Step 2 failed - get network: %v", err)
	}
	if resource.Network.ID != createdNetworkID {
		t.Errorf("expected network ID '%s', got '%s'", createdNetworkID, resource.Network.ID)
	}
	if resource.Network.CIDR != "10.10.0.0/24" {
		t.Errorf("expected CIDR '10.10.0.0/24', got '%s'", resource.Network.CIDR)
	}
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
	if updatedNetwork.Name != "updated-integration-network" {
		t.Errorf("expected updated name 'updated-integration-network', got '%s'", updatedNetwork.Name)
	}
	if updatedNetwork.Description != "Updated description" {
		t.Errorf("expected updated description 'Updated description', got '%s'", updatedNetwork.Description)
	}
	t.Logf("✓ Network updated: %s", updatedNetwork.Name)

	// Step 4: List ports on network
	t.Log("Step 4: Listing network ports...")
	ports, err := resource.Ports().List(ctx)
	if err != nil {
		t.Fatalf("Step 4 failed - list ports: %v", err)
	}
	if len(ports) != 2 {
		t.Errorf("expected 2 ports, got %d", len(ports))
	}
	for i, port := range ports {
		if port.NetworkID != createdNetworkID {
			t.Errorf("port %d has wrong network ID: expected '%s', got '%s'", i, createdNetworkID, port.NetworkID)
		}
		t.Logf("  Port %d: %s (IP: %v)", i+1, port.ID, port.FixedIPs)
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
