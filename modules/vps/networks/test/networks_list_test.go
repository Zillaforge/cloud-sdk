package networks_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// TestNetworksList_Success verifies successful network listing
func TestNetworksList_Success(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   interface{}
		expectedCount  int
		opts           *networks.ListNetworksOptions
		validateResult func(*testing.T, *networks.NetworkListResponse)
	}{
		{
			name: "list all networks",
			mockResponse: map[string]interface{}{
				"networks": []map[string]interface{}{
					{
						"id":          "net-1",
						"name":        "default-network",
						"description": "Default network",
						"cidr":        "10.0.0.0/24",
						"project_id":  "proj-123",
						"created_at":  "2025-01-01T00:00:00Z",
						"updated_at":  "2025-01-01T00:00:00Z",
					},
					{
						"id":          "net-2",
						"name":        "private-network",
						"description": "Private network",
						"cidr":        "10.0.1.0/24",
						"project_id":  "proj-123",
						"created_at":  "2025-01-02T00:00:00Z",
						"updated_at":  "2025-01-02T00:00:00Z",
					},
				},
			},
			expectedCount: 2,
			opts:          nil,
			validateResult: func(t *testing.T, resp *networks.NetworkListResponse) {
				if len(resp.Networks) != 2 {
					t.Errorf("expected 2 networks, got %d", len(resp.Networks))
				}
				if resp.Networks[0].Name != "default-network" {
					t.Errorf("expected first network name 'default-network', got '%s'", resp.Networks[0].Name)
				}
				if resp.Networks[0].CIDR != "10.0.0.0/24" {
					t.Errorf("expected first network CIDR '10.0.0.0/24', got '%s'", resp.Networks[0].CIDR)
				}
			},
		},
		{
			name: "empty list",
			mockResponse: map[string]interface{}{
				"networks": []map[string]interface{}{},
			},
			expectedCount: 0,
			opts:          nil,
			validateResult: func(t *testing.T, resp *networks.NetworkListResponse) {
				if len(resp.Networks) != 0 {
					t.Errorf("expected 0 networks, got %d", len(resp.Networks))
				}
			},
		},
		{
			name: "filter by name",
			mockResponse: map[string]interface{}{
				"networks": []map[string]interface{}{
					{
						"id":          "net-1",
						"name":        "filtered-network",
						"description": "Filtered result",
						"cidr":        "10.0.0.0/24",
						"project_id":  "proj-123",
						"created_at":  "2025-01-01T00:00:00Z",
						"updated_at":  "2025-01-01T00:00:00Z",
					},
				},
			},
			expectedCount: 1,
			opts: &networks.ListNetworksOptions{
				Name: "filtered-network",
			},
			validateResult: func(t *testing.T, resp *networks.NetworkListResponse) {
				if len(resp.Networks) != 1 {
					t.Errorf("expected 1 network, got %d", len(resp.Networks))
				}
				if resp.Networks[0].Name != "filtered-network" {
					t.Errorf("expected network name 'filtered-network', got '%s'", resp.Networks[0].Name)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method and path
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if r.URL.Path != "/vps/api/v1/project/proj-123/networks" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				// Verify query parameters if opts provided
				if tt.opts != nil && tt.opts.Name != "" {
					if r.URL.Query().Get("name") != tt.opts.Name {
						t.Errorf("expected name query param '%s', got '%s'", tt.opts.Name, r.URL.Query().Get("name"))
					}
				}

				// Send mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Create client
			client := cloudsdk.NewClient(server.URL, "test-token")
			projectClient := client.Project("proj-123")
			vpsClient := projectClient.VPS()

			// Execute test
			ctx := context.Background()
			resp, err := vpsClient.Networks().List(ctx, tt.opts)

			// Verify results
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected non-nil response")
			}
			if len(resp.Networks) != tt.expectedCount {
				t.Errorf("expected %d networks, got %d", tt.expectedCount, len(resp.Networks))
			}

			// Run custom validations
			if tt.validateResult != nil {
				tt.validateResult(t, resp)
			}
		})
	}
}

// TestNetworksList_Errors verifies error handling for network listing
func TestNetworksList_Errors(t *testing.T) {
	tests := []struct {
		name             string
		statusCode       int
		responseBody     interface{}
		expectedErrorMsg string
	}{
		{
			name:       "unauthorized - 401",
			statusCode: http.StatusUnauthorized,
			responseBody: map[string]interface{}{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid or expired token",
			},
			expectedErrorMsg: "Invalid or expired token",
		},
		{
			name:       "forbidden - 403",
			statusCode: http.StatusForbidden,
			responseBody: map[string]interface{}{
				"error_code": "FORBIDDEN",
				"message":    "Insufficient permissions",
			},
			expectedErrorMsg: "Insufficient permissions",
		},
		{
			name:       "internal server error - 500",
			statusCode: http.StatusInternalServerError,
			responseBody: map[string]interface{}{
				"error_code": "INTERNAL_ERROR",
				"message":    "Internal server error",
			},
			expectedErrorMsg: "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Create client
			client := cloudsdk.NewClient(server.URL, "test-token")
			projectClient := client.Project("proj-123")
			vpsClient := projectClient.VPS()

			// Execute test
			ctx := context.Background()
			resp, err := vpsClient.Networks().List(ctx, nil)

			// Verify error
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if resp != nil {
				t.Errorf("expected nil response on error, got %+v", resp)
			}

			// Verify error is SDKError
			sdkErr, ok := err.(*cloudsdk.SDKError)
			if !ok {
				t.Fatalf("expected *cloudsdk.SDKError, got %T", err)
			}
			if sdkErr.StatusCode != tt.statusCode {
				t.Errorf("expected status code %d, got %d", tt.statusCode, sdkErr.StatusCode)
			}
			if sdkErr.Message != tt.expectedErrorMsg {
				t.Errorf("expected error message '%s', got '%s'", tt.expectedErrorMsg, sdkErr.Message)
			}
		})
	}
}
