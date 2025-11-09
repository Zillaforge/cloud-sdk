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

// TestServersCreate_Success verifies successful server creation
func TestServersCreate_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"id":          "svr-new",
		"name":        "new-server",
		"description": "Newly created server",
		"status":      "BUILD",
		"flavor_id":   "flv-1",
		"image_id":    "img-1",
		"project_id":  "proj-123",
		"created_at":  "2025-01-01T00:00:00Z",
		"updated_at":  "2025-01-01T00:00:00Z",
	}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/servers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer mockServer.Close()

	client, err := cloudsdk.New(mockServer.URL, "test-token")
	if err != nil {
		t.Fatalf("failed to create client: %v", err)
	}
	vpsClient := client.Project("proj-123").VPS()

	req := &servers.ServerCreateRequest{
		Name:        "new-server",
		Description: "Newly created server",
		FlavorID:    "flv-1",
		ImageID:     "img-1",
		NICs: []servers.ServerNICCreateRequest{
			{NetworkID: "net-1", SGIDs: []string{"sg-1"}},
		},
		SGIDs: []string{"sg-1"},
	}

	server, err := vpsClient.Servers().Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if server.ID != "svr-new" {
		t.Errorf("expected ID 'svr-new', got '%s'", server.ID)
	}
	if server.Name != "new-server" {
		t.Errorf("expected name 'new-server', got '%s'", server.Name)
	}
	if server.Status != "BUILD" {
		t.Errorf("expected status 'BUILD', got '%s'", server.Status)
	}
}

// TestServersCreate_Errors verifies error handling for server creation
func TestServersCreate_Errors(t *testing.T) {
	tests := []struct {
		name         string
		statusCode   int
		mockResponse interface{}
	}{
		{
			name:       "validation error - 400",
			statusCode: http.StatusBadRequest,
			mockResponse: map[string]interface{}{
				"message":    "Invalid request",
				"error_code": 400,
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
		{
			name:       "quota exceeded - 409",
			statusCode: http.StatusConflict,
			mockResponse: map[string]interface{}{
				"message":    "Quota exceeded",
				"error_code": 409,
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

			req := &servers.ServerCreateRequest{
				Name:     "test-server",
				FlavorID: "flv-1",
				ImageID:  "img-1",
			}

			_, err = vpsClient.Servers().Create(context.Background(), req)
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
