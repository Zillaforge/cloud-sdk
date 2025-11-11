package floatingips_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
)

// TestFloatingIPLifecycle tests the complete floating IP lifecycle
func TestFloatingIPLifecycle(t *testing.T) {
	var createdFIPID string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Create floating IP
		if r.Method == http.MethodPost && r.URL.Path == "/vps/api/v1/project/proj-123/floatingips" {
			createdFIPID = "fip-test-1"
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdFIPID,
				"address":     "203.0.113.10",
				"status":      "PENDING",
				"project_id":  "proj-123",
				"description": "Test floating IP",
				"created_at":  "2025-01-01T00:00:00Z",
			})
			return
		}

		// Get floating IP
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/vps/api/v1/project/proj-123/floatingips/") {
			w.WriteHeader(http.StatusOK)
			status := "PENDING"
			portID := ""
			if strings.Contains(r.URL.Path, "/vps/api/v1/project/proj-123/floatingips/fip-test-1") {
				// After approve, status changes to ACTIVE
				if r.Header.Get("X-After-Approve") == "true" {
					status = "ACTIVE"
					portID = "port-1"
				}
			}
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdFIPID,
				"address":     "203.0.113.10",
				"status":      status,
				"project_id":  "proj-123",
				"port_id":     portID,
				"description": "Updated description",
				"created_at":  "2025-01-01T00:00:00Z",
			})
			return
		}

		// Update floating IP
		if r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/vps/api/v1/project/proj-123/floatingips/") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdFIPID,
				"address":     "203.0.113.10",
				"status":      "PENDING",
				"project_id":  "proj-123",
				"description": "Updated description",
				"created_at":  "2025-01-01T00:00:00Z",
			})
			return
		}

		// Approve floating IP
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/approve") {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Disassociate floating IP
		if r.Method == http.MethodPost && strings.HasSuffix(r.URL.Path, "/disassociate") {
			w.WriteHeader(http.StatusAccepted)
			return
		}

		// Delete floating IP
		if r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/vps/api/v1/project/proj-123/floatingips/") &&
			!strings.HasSuffix(r.URL.Path, "/disassociate") {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		// Default: not found
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error_code": "NOT_FOUND",
			"message":    "Resource not found",
		})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()
	ctx := context.Background()

	// Step 1: Create floating IP
	t.Log("Step 1: Creating floating IP...")
	fip, err := vpsClient.FloatingIPs().Create(ctx, &floatingips.FloatingIPCreateRequest{
		Description: "Test floating IP",
	})
	if err != nil {
		t.Fatalf("Failed to create floating IP: %v", err)
	}
	if fip == nil {
		t.Fatal("Expected non-nil floating IP")
	}
	if fip.ID == "" {
		t.Fatal("Expected floating IP to have an ID")
	}
	if fip.Status != "PENDING" {
		t.Errorf("Expected status PENDING, got %s", fip.Status)
	}
	t.Logf("Created floating IP: %s (status: %s)", fip.ID, fip.Status)

	// Step 2: Get floating IP
	t.Log("Step 2: Getting floating IP...")
	retrievedFIP, err := vpsClient.FloatingIPs().Get(ctx, fip.ID)
	if err != nil {
		t.Fatalf("Failed to get floating IP: %v", err)
	}
	if retrievedFIP.ID != fip.ID {
		t.Errorf("Expected ID %s, got %s", fip.ID, retrievedFIP.ID)
	}

	// Step 3: Update floating IP description
	t.Log("Step 3: Updating floating IP description...")
	updatedFIP, err := vpsClient.FloatingIPs().Update(ctx, fip.ID, &floatingips.FloatingIPUpdateRequest{
		Description: "Updated description",
	})
	if err != nil {
		t.Fatalf("Failed to update floating IP: %v", err)
	}
	if updatedFIP.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", updatedFIP.Description)
	}

	// Step 4: Approve floating IP
	t.Log("Step 4: Approving floating IP...")
	if err := vpsClient.FloatingIPs().Approve(ctx, fip.ID); err != nil {
		t.Fatalf("Failed to approve floating IP: %v", err)
	}

	// Step 5: Verify status changed (simulated - in real scenario would wait/poll)
	t.Log("Step 5: Verifying floating IP status...")
	// In a real scenario, we'd poll until status changes to ACTIVE

	// Step 6: Disassociate floating IP
	t.Log("Step 6: Disassociating floating IP...")
	err = vpsClient.FloatingIPs().Disassociate(ctx, fip.ID)
	if err != nil {
		t.Fatalf("Failed to disassociate floating IP: %v", err)
	}

	// Step 7: Delete floating IP
	t.Log("Step 7: Deleting floating IP...")
	if err := vpsClient.FloatingIPs().Delete(ctx, fip.ID); err != nil {
		t.Fatalf("Failed to delete floating IP: %v", err)
	}

	t.Log("Floating IP lifecycle test completed successfully!")
}

// TestFloatingIPLifecycle_ErrorHandling tests error scenarios in lifecycle
func TestFloatingIPLifecycle_ErrorHandling(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error_code": "NOT_FOUND",
			"message":    "Floating IP not found",
		})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()
	ctx := context.Background()

	// Test operations on non-existent floating IP
	t.Run("get non-existent floating IP", func(t *testing.T) {
		_, err := vpsClient.FloatingIPs().Get(ctx, "fip-nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent floating IP")
		}
	})

	t.Run("update non-existent floating IP", func(t *testing.T) {
		_, err := vpsClient.FloatingIPs().Update(ctx, "fip-nonexistent", &floatingips.FloatingIPUpdateRequest{
			Description: "test",
		})
		if err == nil {
			t.Error("Expected error for non-existent floating IP")
		}
	})

	t.Run("delete non-existent floating IP", func(t *testing.T) {
		err := vpsClient.FloatingIPs().Delete(ctx, "fip-nonexistent")
		if err == nil {
			t.Error("Expected error for non-existent floating IP")
		}
	})
}

// TestFloatingIPList_Integration tests List operation with mock API server
// returning pb.FIPListOutput structure (floatingips array)
func TestFloatingIPList_Integration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Handle List operation
		if r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-123/floatingips" {
			w.WriteHeader(http.StatusOK)

			// Build response based on query parameters
			var floatingIPs []map[string]interface{}

			if r.URL.Query().Get("status") == "ACTIVE" {
				floatingIPs = []map[string]interface{}{
					{
						"id":         "fip-1",
						"uuid":       "550e8400-e29b-41d4-a716-446655440001",
						"name":       "web-ip",
						"address":    "203.0.113.10",
						"status":     "ACTIVE",
						"project_id": "proj-123",
						"user_id":    "user-123",
						"reserved":   false,
						"createdAt":  "2025-01-01T00:00:00Z",
					},
				}
			} else {
				floatingIPs = []map[string]interface{}{
					{
						"id":         "fip-1",
						"uuid":       "550e8400-e29b-41d4-a716-446655440001",
						"name":       "web-ip",
						"address":    "203.0.113.10",
						"status":     "ACTIVE",
						"project_id": "proj-123",
						"user_id":    "user-123",
						"reserved":   false,
						"createdAt":  "2025-01-01T00:00:00Z",
					},
					{
						"id":         "fip-2",
						"uuid":       "550e8400-e29b-41d4-a716-446655440002",
						"name":       "db-ip",
						"address":    "203.0.113.11",
						"status":     "PENDING",
						"project_id": "proj-123",
						"user_id":    "user-123",
						"reserved":   false,
						"createdAt":  "2025-01-02T00:00:00Z",
					},
				}
			}

			response := map[string]interface{}{
				"floating_ips": floatingIPs,
			}
			_ = json.NewEncoder(w).Encode(response)
			return
		}
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()
	ctx := context.Background()

	// Test listing all floating IPs
	t.Run("list all floating IPs", func(t *testing.T) {
		result, err := vpsClient.FloatingIPs().List(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if result == nil {
			t.Fatal("expected non-nil result")
		}
		if len(result) != 2 {
			t.Errorf("expected 2 floating IPs, got %d", len(result))
		}
		if result[0].Status != floatingips.FloatingIPStatusActive {
			t.Errorf("expected first IP status ACTIVE, got %s", result[0].Status)
		}
	})

	// Test listing with status filter
	t.Run("list with status filter", func(t *testing.T) {
		opts := &floatingips.ListFloatingIPsOptions{Status: "ACTIVE"}
		result, err := vpsClient.FloatingIPs().List(ctx, opts)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if len(result) != 1 {
			t.Errorf("expected 1 floating IP with filter, got %d", len(result))
		}
		if result[0].ID != "fip-1" {
			t.Errorf("expected ID 'fip-1', got '%s'", result[0].ID)
		}
	})
}

// TestFloatingIPCRUDWorkflow_Integration tests the complete CRUD workflow:
// create → get → update → delete (T036)
func TestFloatingIPCRUDWorkflow_Integration(t *testing.T) {
	var createdID string

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// CREATE operation
		if r.Method == http.MethodPost && r.URL.Path == "/vps/api/v1/project/proj-123/floatingips" {
			createdID = "fip-workflow-1"
			w.WriteHeader(http.StatusCreated)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdID,
				"uuid":        "550e8400-e29b-41d4-a716-446655440000",
				"address":     "203.0.113.20",
				"status":      "PENDING",
				"project_id":  "proj-123",
				"user_id":     "user-123",
				"name":        "workflow-test-ip",
				"description": "Test IP for workflow",
				"reserved":    false,
				"createdAt":   "2025-01-01T00:00:00Z",
			})
			return
		}

		// GET operation
		if r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/vps/api/v1/project/proj-123/floatingips/fip-workflow-1") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdID,
				"uuid":        "550e8400-e29b-41d4-a716-446655440000",
				"address":     "203.0.113.20",
				"status":      "PENDING",
				"project_id":  "proj-123",
				"user_id":     "user-123",
				"name":        "workflow-test-ip",
				"description": "Test IP for workflow",
				"reserved":    false,
				"createdAt":   "2025-01-01T00:00:00Z",
			})
			return
		}

		// UPDATE operation
		if r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/vps/api/v1/project/proj-123/floatingips/fip-workflow-1") {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          createdID,
				"uuid":        "550e8400-e29b-41d4-a716-446655440000",
				"address":     "203.0.113.20",
				"status":      "PENDING",
				"project_id":  "proj-123",
				"user_id":     "user-123",
				"name":        "updated-workflow-ip",
				"description": "Updated description",
				"reserved":    false,
				"createdAt":   "2025-01-01T00:00:00Z",
				"updatedAt":   "2025-01-02T00:00:00Z",
			})
			return
		}

		// DELETE operation
		if r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/vps/api/v1/project/proj-123/floatingips/fip-workflow-1") {
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()
	ctx := context.Background()

	// CREATE
	t.Run("create floating IP", func(t *testing.T) {
		req := &floatingips.FloatingIPCreateRequest{
			Name:        "workflow-test-ip",
			Description: "Test IP for workflow",
		}
		fip, err := vpsClient.FloatingIPs().Create(ctx, req)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}
		if fip.ID != "fip-workflow-1" {
			t.Errorf("expected ID 'fip-workflow-1', got '%s'", fip.ID)
		}
		createdID = fip.ID
	})

	// GET
	t.Run("get floating IP", func(t *testing.T) {
		fip, err := vpsClient.FloatingIPs().Get(ctx, createdID)
		if err != nil {
			t.Fatalf("get failed: %v", err)
		}
		if fip.ID != createdID {
			t.Errorf("expected ID '%s', got '%s'", createdID, fip.ID)
		}
		if fip.Status != floatingips.FloatingIPStatusPending {
			t.Errorf("expected status PENDING, got %s", fip.Status)
		}
	})

	// UPDATE
	t.Run("update floating IP", func(t *testing.T) {
		updateReq := &floatingips.FloatingIPUpdateRequest{
			Name:        "updated-workflow-ip",
			Description: "Updated description",
		}
		fip, err := vpsClient.FloatingIPs().Update(ctx, createdID, updateReq)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}
		if fip.Name != "updated-workflow-ip" {
			t.Errorf("expected name 'updated-workflow-ip', got '%s'", fip.Name)
		}
		if fip.Description != "Updated description" {
			t.Errorf("expected description 'Updated description', got '%s'", fip.Description)
		}
	})

	// DELETE
	t.Run("delete floating IP", func(t *testing.T) {
		err := vpsClient.FloatingIPs().Delete(ctx, createdID)
		if err != nil {
			t.Fatalf("delete failed: %v", err)
		}
	})
}

// TestFloatingIPDisassociate_Integration tests the Disassociate operation (T037)
func TestFloatingIPDisassociate_Integration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// DISASSOCIATE operation
		if r.Method == http.MethodPost && strings.Contains(r.URL.Path, "/floatingips/fip-disassoc-1/disassociate") {
			w.WriteHeader(http.StatusOK)
			// Return floating IP with cleared device associations
			_ = json.NewEncoder(w).Encode(map[string]interface{}{
				"id":          "fip-disassoc-1",
				"uuid":        "550e8400-e29b-41d4-a716-446655440001",
				"address":     "203.0.113.30",
				"status":      "ACTIVE",
				"project_id":  "proj-123",
				"user_id":     "user-123",
				"name":        "test-disassoc-ip",
				"reserved":    false,
				"device_id":   nil,
				"device_name": nil,
				"device_type": nil,
				"port_id":     nil,
				"createdAt":   "2025-01-01T00:00:00Z",
			})
			return
		}
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()
	ctx := context.Background()

	t.Run("disassociate floating IP", func(t *testing.T) {
		err := vpsClient.FloatingIPs().Disassociate(ctx, "fip-disassoc-1")
		if err != nil {
			t.Fatalf("disassociate failed: %v", err)
		}
	})
}
