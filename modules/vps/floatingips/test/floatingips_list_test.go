package floatingips_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
)

// TestFloatingIPsList_Success verifies successful floating IP listing
func TestFloatingIPsList_Success(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   interface{}
		expectedCount  int
		opts           *floatingips.ListFloatingIPsOptions
		validateResult func(*testing.T, *floatingips.FloatingIPListResponse)
	}{
		{
			name: "list all floating IPs",
			mockResponse: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":          "fip-1",
						"address":     "203.0.113.10",
						"status":      "ACTIVE",
						"project_id":  "proj-123",
						"port_id":     "port-1",
						"description": "Web server IP",
						"created_at":  "2025-01-01T00:00:00Z",
					},
					{
						"id":          "fip-2",
						"address":     "203.0.113.11",
						"status":      "PENDING",
						"project_id":  "proj-123",
						"description": "Database IP",
						"created_at":  "2025-01-02T00:00:00Z",
					},
				},
			},
			expectedCount: 2,
			opts:          nil,
			validateResult: func(t *testing.T, resp *floatingips.FloatingIPListResponse) {
				if len(resp.Items) != 2 {
					t.Errorf("expected 2 floating IPs, got %d", len(resp.Items))
				}
				if resp.Items[0].Address != "203.0.113.10" {
					t.Errorf("expected first IP address '203.0.113.10', got '%s'", resp.Items[0].Address)
				}
				if resp.Items[0].Status != "ACTIVE" {
					t.Errorf("expected first IP status 'ACTIVE', got '%s'", resp.Items[0].Status)
				}
				if resp.Items[1].Status != "PENDING" {
					t.Errorf("expected second IP status 'PENDING', got '%s'", resp.Items[1].Status)
				}
			},
		},
		{
			name: "empty list",
			mockResponse: map[string]interface{}{
				"items": []map[string]interface{}{},
			},
			expectedCount: 0,
			opts:          nil,
			validateResult: func(t *testing.T, resp *floatingips.FloatingIPListResponse) {
				if len(resp.Items) != 0 {
					t.Errorf("expected 0 floating IPs, got %d", len(resp.Items))
				}
			},
		},
		{
			name: "filter by status",
			mockResponse: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":         "fip-1",
						"address":    "203.0.113.10",
						"status":     "ACTIVE",
						"project_id": "proj-123",
						"created_at": "2025-01-01T00:00:00Z",
					},
				},
			},
			expectedCount: 1,
			opts: &floatingips.ListFloatingIPsOptions{
				Status: "ACTIVE",
			},
			validateResult: func(t *testing.T, resp *floatingips.FloatingIPListResponse) {
				if len(resp.Items) != 1 {
					t.Errorf("expected 1 floating IP, got %d", len(resp.Items))
				}
				if resp.Items[0].Status != "ACTIVE" {
					t.Errorf("expected IP status 'ACTIVE', got '%s'", resp.Items[0].Status)
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
				if r.URL.Path != "/vps/api/v1/project/proj-123/floatingips" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				// Verify query parameters if opts provided
				if tt.opts != nil && tt.opts.Status != "" {
					if r.URL.Query().Get("status") != tt.opts.Status {
						t.Errorf("expected status query param '%s', got '%s'", tt.opts.Status, r.URL.Query().Get("status"))
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
			resp, err := vpsClient.FloatingIPs().List(ctx, tt.opts)

			// Verify results
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected non-nil response")
			}
			if len(resp.Items) != tt.expectedCount {
				t.Errorf("expected %d items, got %d", tt.expectedCount, len(resp.Items))
			}

			// Run custom validations
			if tt.validateResult != nil {
				tt.validateResult(t, resp)
			}
		})
	}
}

// TestFloatingIPsList_Errors verifies error handling for floating IP listing
func TestFloatingIPsList_Errors(t *testing.T) {
	tests := []struct {
		name           string
		mockResponse   interface{}
		mockStatusCode int
		expectError    bool
	}{
		{
			name: "unauthorized - 401",
			mockResponse: map[string]interface{}{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid token",
			},
			mockStatusCode: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name: "forbidden - 403",
			mockResponse: map[string]interface{}{
				"error_code": "FORBIDDEN",
				"message":    "Access denied",
			},
			mockStatusCode: http.StatusForbidden,
			expectError:    true,
		},
		{
			name: "internal server error - 500",
			mockResponse: map[string]interface{}{
				"error_code": "INTERNAL_ERROR",
				"message":    "Internal server error",
			},
			mockStatusCode: http.StatusInternalServerError,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Create client
			client := cloudsdk.NewClient(server.URL, "test-token")
			projectClient := client.Project("proj-123")
			vpsClient := projectClient.VPS()

			// Execute test
			ctx := context.Background()
			resp, err := vpsClient.FloatingIPs().List(ctx, nil)

			// Verify error
			if !tt.expectError {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Error("expected error, got nil")
				return
			}

			if resp != nil {
				t.Error("expected nil response on error")
			}
		})
	}
}
