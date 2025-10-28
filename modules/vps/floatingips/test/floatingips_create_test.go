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

// TestFloatingIPsCreate_Success verifies successful floating IP creation
func TestFloatingIPsCreate_Success(t *testing.T) {
	tests := []struct {
		name           string
		request        *floatingips.FloatingIPCreateRequest
		mockResponse   interface{}
		validateResult func(*testing.T, *floatingips.FloatingIP)
	}{
		{
			name: "create with description",
			request: &floatingips.FloatingIPCreateRequest{
				Description: "Web server IP",
			},
			mockResponse: map[string]interface{}{
				"id":          "fip-1",
				"address":     "203.0.113.10",
				"status":      "ACTIVE",
				"project_id":  "proj-123",
				"description": "Web server IP",
				"created_at":  "2025-01-01T00:00:00Z",
			},
			validateResult: func(t *testing.T, fip *floatingips.FloatingIP) {
				if fip.ID != "fip-1" {
					t.Errorf("expected ID 'fip-1', got '%s'", fip.ID)
				}
				if fip.Address != "203.0.113.10" {
					t.Errorf("expected address '203.0.113.10', got '%s'", fip.Address)
				}
				if fip.Status != "ACTIVE" {
					t.Errorf("expected status 'ACTIVE', got '%s'", fip.Status)
				}
				if fip.Description != "Web server IP" {
					t.Errorf("expected description 'Web server IP', got '%s'", fip.Description)
				}
			},
		},
		{
			name:    "create without description",
			request: &floatingips.FloatingIPCreateRequest{},
			mockResponse: map[string]interface{}{
				"id":         "fip-2",
				"address":    "203.0.113.11",
				"status":     "ACTIVE",
				"project_id": "proj-123",
				"created_at": "2025-01-01T00:00:00Z",
			},
			validateResult: func(t *testing.T, fip *floatingips.FloatingIP) {
				if fip.ID != "fip-2" {
					t.Errorf("expected ID 'fip-2', got '%s'", fip.ID)
				}
				if fip.Description != "" {
					t.Errorf("expected empty description, got '%s'", fip.Description)
				}
			},
		},
		{
			name: "create with pending status (requires approval)",
			request: &floatingips.FloatingIPCreateRequest{
				Description: "Pending IP",
			},
			mockResponse: map[string]interface{}{
				"id":          "fip-3",
				"address":     "203.0.113.12",
				"status":      "PENDING",
				"project_id":  "proj-123",
				"description": "Pending IP",
				"created_at":  "2025-01-01T00:00:00Z",
			},
			validateResult: func(t *testing.T, fip *floatingips.FloatingIP) {
				if fip.Status != "PENDING" {
					t.Errorf("expected status 'PENDING', got '%s'", fip.Status)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method and path
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != "/vps/api/v1/project/proj-123/floatingips" {
					t.Errorf("unexpected path: %s", r.URL.Path)
				}

				// Verify request body
				var body floatingips.FloatingIPCreateRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}
				if body.Description != tt.request.Description {
					t.Errorf("expected description '%s', got '%s'", tt.request.Description, body.Description)
				}

				// Send mock response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			// Create client
			client := cloudsdk.NewClient(server.URL, "test-token")
			projectClient := client.Project("proj-123")
			vpsClient := projectClient.VPS()

			// Execute test
			ctx := context.Background()
			fip, err := vpsClient.FloatingIPs().Create(ctx, tt.request)

			// Verify results
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fip == nil {
				t.Fatal("expected non-nil floating IP")
			}

			// Run custom validations
			if tt.validateResult != nil {
				tt.validateResult(t, fip)
			}
		})
	}
}

// TestFloatingIPsCreate_Errors verifies error handling for floating IP creation
func TestFloatingIPsCreate_Errors(t *testing.T) {
	tests := []struct {
		name           string
		request        *floatingips.FloatingIPCreateRequest
		mockResponse   interface{}
		mockStatusCode int
		expectError    bool
	}{
		{
			name:    "validation error - 400",
			request: &floatingips.FloatingIPCreateRequest{},
			mockResponse: map[string]interface{}{
				"error_code": "VALIDATION_ERROR",
				"message":    "Invalid request",
			},
			mockStatusCode: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:    "unauthorized - 401",
			request: &floatingips.FloatingIPCreateRequest{},
			mockResponse: map[string]interface{}{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid token",
			},
			mockStatusCode: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:    "forbidden - 403",
			request: &floatingips.FloatingIPCreateRequest{},
			mockResponse: map[string]interface{}{
				"error_code": "FORBIDDEN",
				"message":    "Access denied",
			},
			mockStatusCode: http.StatusForbidden,
			expectError:    true,
		},
		{
			name: "quota exceeded - 409",
			request: &floatingips.FloatingIPCreateRequest{
				Description: "Test",
			},
			mockResponse: map[string]interface{}{
				"error_code": "QUOTA_EXCEEDED",
				"message":    "Floating IP quota exceeded",
			},
			mockStatusCode: http.StatusConflict,
			expectError:    true,
		},
		{
			name:    "internal server error - 500",
			request: &floatingips.FloatingIPCreateRequest{},
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
			fip, err := vpsClient.FloatingIPs().Create(ctx, tt.request)

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

			if fip != nil {
				t.Error("expected nil floating IP on error")
			}
		})
	}
}
