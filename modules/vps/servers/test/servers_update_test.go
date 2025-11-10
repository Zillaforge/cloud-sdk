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

// TestServersUpdate_Success verifies successful server update
func TestServersUpdate_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"id":          "svr-1",
		"name":        "updated-server",
		"description": "Updated description",
		"status":      "ACTIVE",
		"flavor_id":   "flv-1",
		"image_id":    "img-1",
		"project_id":  "proj-123",
		"created_at":  "2025-01-01T00:00:00Z",
		"updated_at":  "2025-01-02T00:00:00Z",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/servers/svr-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	req := &servers.ServerUpdateRequest{
		Name:        "updated-server",
		Description: "Updated description",
	}

	server, err := vpsClient.Servers().Update(context.Background(), "svr-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if server.Server.Name != "updated-server" {
		t.Errorf("expected name 'updated-server', got '%s'", server.Server.Name)
	}
	if server.Server.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got '%s'", server.Server.Description)
	}
}

// TestServersUpdate_Errors verifies error handling for server update
func TestServersUpdate_Errors(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"not found - 404", http.StatusNotFound},
		{"bad request - 400", http.StatusBadRequest},
		{"unauthorized - 401", http.StatusUnauthorized},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(map[string]interface{}{
					"message":    "Error",
					"error_code": tt.statusCode,
				})
			}))
			defer mockServer.Close()

			client, err := cloudsdk.New(mockServer.URL, "test-token")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			vpsClient := client.Project("proj-123").VPS()

			req := &servers.ServerUpdateRequest{Name: "updated"}
			_, err = vpsClient.Servers().Update(context.Background(), "svr-1", req)
			if err == nil {
				t.Fatal("expected error, got nil")
			}
		})
	}
}
