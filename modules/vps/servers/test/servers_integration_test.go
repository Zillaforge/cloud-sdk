package servers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
)

// TestServerLifecycle tests the complete server lifecycle
func TestServerLifecycle(t *testing.T) {
	var serverID = "svr-lifecycle-123"

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Create server
		if r.Method == "POST" && r.URL.Path == "/vps/api/v1/project/proj-123/servers" {
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          serverID,
				"name":        "lifecycle-test-server",
				"description": "Server for lifecycle testing",
				"status":      "BUILD",
				"flavor_id":   "flv-1",
				"image_id":    "img-1",
				"project_id":  "proj-123",
			})
			return
		}

		// Get server
		if r.Method == "GET" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/"+serverID {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          serverID,
				"name":        "lifecycle-test-server",
				"description": "Server for lifecycle testing",
				"status":      "ACTIVE",
				"flavor_id":   "flv-1",
				"image_id":    "img-1",
				"project_id":  "proj-123",
			})
			return
		}

		// Update server
		if r.Method == "PUT" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/"+serverID {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          serverID,
				"name":        "updated-server",
				"description": "Updated description",
				"status":      "ACTIVE",
				"flavor_id":   "flv-1",
				"image_id":    "img-1",
				"project_id":  "proj-123",
			})
			return
		}

		// Server action
		if r.Method == "POST" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/"+serverID+"/action" {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Get metrics
		if r.Method == "GET" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/"+serverID+"/metric" {
			_ = json.NewEncoder(w).Encode([]map[string]interface{}{
				{
					"name": "cpu",
					"measures": []map[string]interface{}{
						{"granularity": 3600, "timestamp": 1704067200, "value": 45.5},
					},
				},
			})
			return
		}

		// Get VNC URL
		if r.Method == "GET" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/"+serverID+"/vnc_url" {
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"url": "https://console.example.com/vnc?token=test123",
			})
			return
		}

		// Delete server
		if r.Method == "DELETE" && r.URL.Path == "/vps/api/v1/project/proj-123/servers/"+serverID {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		t.Logf("Unhandled request: %s %s", r.Method, r.URL.Path)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	t.Log("Step 1: Creating server...")
	createReq := &servers.ServerCreateRequest{
		Name:        "lifecycle-test-server",
		Description: "Server for lifecycle testing",
		FlavorID:    "flv-1",
		ImageID:     "img-1",
		NICs: []servers.ServerNICCreateRequest{
			{NetworkID: "net-1", SGIDs: []string{"sg-1"}},
		},
		SGIDs: []string{"sg-1"},
	}
	server, err := vpsClient.Servers().Create(context.Background(), createReq)
	if err != nil {
		t.Fatalf("failed to create server: %v", err)
	}
	if server.ID != serverID {
		t.Errorf("expected server ID '%s', got '%s'", serverID, server.ID)
	}
	t.Logf("✓ Server created with ID: %s", server.ID)

	t.Log("Step 2: Getting server...")
	serverResource, err := vpsClient.Servers().Get(context.Background(), serverID)
	if err != nil {
		t.Fatalf("failed to get server: %v", err)
	}
	if serverResource.Server.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", serverResource.Server.Status)
	}
	t.Log("✓ Server retrieved successfully")

	t.Log("Step 3: Updating server...")
	updateReq := &servers.ServerUpdateRequest{
		Name:        "updated-server",
		Description: "Updated description",
	}
	updatedServer, err := vpsClient.Servers().Update(context.Background(), serverID, updateReq)
	if err != nil {
		t.Fatalf("failed to update server: %v", err)
	}
	if updatedServer.Name != "updated-server" {
		t.Errorf("expected name 'updated-server', got '%s'", updatedServer.Name)
	}
	t.Log("✓ Server updated successfully")

	t.Log("Step 4: Performing server action...")
	actionReq := &servers.ServerActionRequest{
		Action: "reboot",
	}
	err = vpsClient.Servers().Action(context.Background(), serverID, actionReq)
	if err != nil {
		t.Fatalf("failed to perform action: %v", err)
	}
	t.Log("✓ Server action completed")

	t.Log("Step 5: Getting server metrics...")
	metricsReq := &servers.ServerMetricsRequest{
		Type:        "cpu",
		Start:       1704067200,
		Granularity: 3600,
	}
	metrics, err := vpsClient.Servers().Metrics(context.Background(), serverID, metricsReq)
	if err != nil {
		t.Fatalf("failed to get metrics: %v", err)
	}
	if len(*metrics) == 0 {
		t.Error("expected metrics data, got empty")
	}
	t.Log("✓ Server metrics retrieved")

	t.Log("Step 6: Getting VNC URL...")
	vncResp, err := vpsClient.Servers().GetVNCConsoleURL(context.Background(), serverID)
	if err != nil {
		t.Fatalf("failed to get VNC URL: %v", err)
	}
	if vncResp.URL == "" {
		t.Error("expected VNC URL, got empty")
	}
	t.Log("✓ VNC URL retrieved")

	t.Log("Step 7: Deleting server...")
	err = vpsClient.Servers().Delete(context.Background(), serverID)
	if err != nil {
		t.Fatalf("failed to delete server: %v", err)
	}
	t.Log("✓ Server deleted successfully")

	t.Log("✅ Server lifecycle test completed successfully!")
}

// TestServerLifecycle_ErrorHandling tests error scenarios
func TestServerLifecycle_ErrorHandling(t *testing.T) {
	tests := []struct {
		name       string
		operation  string
		statusCode int
	}{
		{"get non-existent server", "GET", http.StatusNotFound},
		{"update non-existent server", "PUT", http.StatusNotFound},
		{"delete non-existent server", "DELETE", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"message":    "Resource not found",
					"error_code": tt.statusCode,
				})
			}))
			defer mockServer.Close()

			client, err := cloudsdk.New(mockServer.URL, "test-token")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			vpsClient := client.Project("proj-123").VPS()

			switch tt.operation {
			case "GET":
				_, err = vpsClient.Servers().Get(context.Background(), "non-existent")
			case "PUT":
				req := &servers.ServerUpdateRequest{Name: "test"}
				_, err = vpsClient.Servers().Update(context.Background(), "non-existent", req)
			case "DELETE":
				err = vpsClient.Servers().Delete(context.Background(), "non-existent")
			}

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			sdkErr, ok := err.(*cloudsdk.SDKError)
			if !ok {
				t.Fatalf("expected *cloudsdk.SDKError, got %T", err)
			}

			if sdkErr.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, sdkErr.StatusCode)
			}
		})
	}
}
