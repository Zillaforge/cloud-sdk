package servers_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

// TestServersGet_Success verifies successful server retrieval
func TestServersGet_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"id":          "svr-1",
		"name":        "web-server",
		"description": "Web application server",
		"status":      "ACTIVE",
		"flavor_id":   "flv-1",
		"image_id":    "img-1",
		"project_id":  "proj-123",
		"created_at":  "2025-01-01T00:00:00Z",
		"updated_at":  "2025-01-01T00:00:00Z",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
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

	serverResource, err := vpsClient.Servers().Get(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if serverResource.Server.ID != "svr-1" {
		t.Errorf("expected ID 'svr-1', got '%s'", serverResource.Server.ID)
	}
	if serverResource.Server.Name != "web-server" {
		t.Errorf("expected name 'web-server', got '%s'", serverResource.Server.Name)
	}

	// Verify sub-resource accessors exist
	if serverResource.NICs() == nil {
		t.Error("expected NICs() to return non-nil")
	}
	if serverResource.Volumes() == nil {
		t.Error("expected Volumes() to return non-nil")
	}
}

// TestServersGet_Errors verifies error handling for server retrieval
func TestServersGet_Errors(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		mockResponse interface{}
	}{
		{
			name:       "not found - 404",
			statusCode: http.StatusNotFound,
			mockResponse: map[string]interface{}{
				"message":    "Server not found",
				"error_code": 404,
			},
		},
		{
			name:       "unauthorized - 401",
			statusCode: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"message":    "Unauthorized",
				"error_code": 401,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer mockServer.Close()

			client, err := cloudsdk.New(mockServer.URL, "test-token")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			vpsClient := client.Project("proj-123").VPS()

			_, err = vpsClient.Servers().Get(context.Background(), "svr-1")
			if err == nil {
				t.Fatal("expected error, got nil")
			}

			var sdkErr *cloudsdk.SDKError
			if !errors.As(err, &sdkErr) {
				t.Fatalf("expected *cloudsdk.SDKError, got %T", err)
			}

			if sdkErr.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, sdkErr.StatusCode)
			}
		})
	}
}
