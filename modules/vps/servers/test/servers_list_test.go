package servers_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
	serversmod "github.com/Zillaforge/cloud-sdk/modules/vps/servers"
)

// TestServersList_Success verifies successful server listing
func TestServersList_Success(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   interface{}
		expectedCount  int
		opts           *servers.ServersListRequest
		validateResult func(*testing.T, []*serversmod.ServerResource)
	}{
		{
			name: "list all servers",
			mockResponse: map[string]interface{}{
				"servers": []map[string]interface{}{
					{
						"id":          "svr-1",
						"name":        "web-server",
						"description": "Web application server",
						"status":      "ACTIVE",
						"flavor_id":   "flv-1",
						"image_id":    "img-1",
						"project_id":  "proj-123",
						"user_id":     "user-1",
						"created_at":  "2025-01-01T00:00:00Z",
						"updated_at":  "2025-01-01T00:00:00Z",
					},
					{
						"id":          "svr-2",
						"name":        "db-server",
						"description": "Database server",
						"status":      "ACTIVE",
						"flavor_id":   "flv-2",
						"image_id":    "img-2",
						"project_id":  "proj-123",
						"user_id":     "user-1",
						"created_at":  "2025-01-02T00:00:00Z",
						"updated_at":  "2025-01-02T00:00:00Z",
					},
				},
				"total": 2,
			},
			expectedCount: 2,
			opts:          nil,
			validateResult: func(t *testing.T, resp []*serversmod.ServerResource) {
				if len(resp) != 2 {
					t.Errorf("expected 2 servers, got %d", len(resp))
				}
				if resp[0].Name != "web-server" {
					t.Errorf("expected first server name 'web-server', got '%s'", resp[0].Name)
				}
			},
		},
		{
			name: "filter by status",
			mockResponse: map[string]interface{}{
				"servers": []map[string]interface{}{
					{
						"id":         "svr-3",
						"name":       "build-server",
						"status":     "BUILD",
						"flavor_id":  "flv-1",
						"image_id":   "img-1",
						"project_id": "proj-123",
						"created_at": "2025-01-03T00:00:00Z",
						"updated_at": "2025-01-03T00:00:00Z",
					},
				},
				"total": 1,
			},
			expectedCount: 1,
			opts: &servers.ServersListRequest{
				Status: "BUILD",
			},
			validateResult: func(t *testing.T, resp []*serversmod.ServerResource) {
				if len(resp) != 1 {
					t.Errorf("expected 1 server, got %d", len(resp))
				}
				if resp[0].Status != servers.ServerStatusBuild {
					t.Errorf("expected status 'BUILD', got '%s'", resp[0].Status)
				}
			},
		},
		{
			name: "empty list",
			mockResponse: map[string]interface{}{
				"servers": []map[string]interface{}{},
			},
			expectedCount: 0,
			opts:          nil,
			validateResult: func(t *testing.T, resp []*serversmod.ServerResource) {
				if len(resp) != 0 {
					t.Errorf("expected 0 servers, got %d", len(resp))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock server
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/vps/api/v1/project/proj-123/servers" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				// Verify query parameters if opts provided
				if tt.opts != nil && tt.opts.Status != "" {
					if r.URL.Query().Get("status") != tt.opts.Status {
						t.Errorf("expected status query param '%s', got '%s'", tt.opts.Status, r.URL.Query().Get("status"))
					}
				}

				w.Header().Set("Content-Type", "application/json")
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer mockServer.Close()

			// Create client
			client, err := cloudsdk.New(mockServer.URL, "test-token")
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}
			vpsClient := client.Project("proj-123").VPS()

			// Call List
			resp, err := vpsClient.Servers().List(context.Background(), tt.opts)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Validate
			tt.validateResult(t, resp)
		})
	}
}

// TestServersList_Errors verifies error handling for server listing
func TestServersList_Errors(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		mockResponse   interface{}
		expectedErrMsg string
	}{
		{
			name:       "unauthorized - 401",
			statusCode: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"message":    "Unauthorized",
				"error_code": 401,
			},
			expectedErrMsg: "Unauthorized",
		},
		{
			name:       "forbidden - 403",
			statusCode: http.StatusForbidden,
			mockResponse: map[string]interface{}{
				"message":    "Forbidden",
				"error_code": 403,
			},
			expectedErrMsg: "Forbidden",
		},
		{
			name:       "internal server error - 500",
			statusCode: http.StatusInternalServerError,
			mockResponse: map[string]interface{}{
				"message":    "Internal server error",
				"error_code": 500,
			},
			expectedErrMsg: "Internal server error",
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

			_, err = vpsClient.Servers().List(context.Background(), nil)
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
