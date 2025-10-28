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

// TestFloatingIPsUpdate_Success verifies successful floating IP update
func TestFloatingIPsUpdate_Success(t *testing.T) {
	tests := []struct {
		name           string
		fipID          string
		request        *floatingips.FloatingIPUpdateRequest
		mockResponse   interface{}
		validateResult func(*testing.T, *floatingips.FloatingIP)
	}{
		{
			name:  "update description",
			fipID: "fip-1",
			request: &floatingips.FloatingIPUpdateRequest{
				Description: "Updated description",
			},
			mockResponse: map[string]interface{}{
				"id":          "fip-1",
				"address":     "203.0.113.10",
				"status":      "ACTIVE",
				"project_id":  "proj-123",
				"description": "Updated description",
				"created_at":  "2025-01-01T00:00:00Z",
			},
			validateResult: func(t *testing.T, fip *floatingips.FloatingIP) {
				if fip.Description != "Updated description" {
					t.Errorf("expected description 'Updated description', got '%s'", fip.Description)
				}
			},
		},
		{
			name:  "clear description",
			fipID: "fip-1",
			request: &floatingips.FloatingIPUpdateRequest{
				Description: "",
			},
			mockResponse: map[string]interface{}{
				"id":         "fip-1",
				"address":    "203.0.113.10",
				"status":     "ACTIVE",
				"project_id": "proj-123",
				"created_at": "2025-01-01T00:00:00Z",
			},
			validateResult: func(t *testing.T, fip *floatingips.FloatingIP) {
				if fip.Description != "" {
					t.Errorf("expected empty description, got '%s'", fip.Description)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("expected PUT request, got %s", r.Method)
				}
				expectedPath := "/vps/api/v1/project/proj-123/floatingips/" + tt.fipID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				var body floatingips.FloatingIPUpdateRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			client := cloudsdk.NewClient(server.URL, "test-token")
			projectClient := client.Project("proj-123")
			vpsClient := projectClient.VPS()

			ctx := context.Background()
			fip, err := vpsClient.FloatingIPs().Update(ctx, tt.fipID, tt.request)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if fip == nil {
				t.Fatal("expected non-nil floating IP")
			}

			if tt.validateResult != nil {
				tt.validateResult(t, fip)
			}
		})
	}
}

// TestFloatingIPsUpdate_Errors verifies error handling for floating IP update
func TestFloatingIPsUpdate_Errors(t *testing.T) {
	tests := []struct {
		name           string
		fipID          string
		request        *floatingips.FloatingIPUpdateRequest
		mockResponse   interface{}
		mockStatusCode int
		expectError    bool
	}{
		{
			name:  "validation error - 400",
			fipID: "fip-1",
			request: &floatingips.FloatingIPUpdateRequest{
				Description: string(make([]byte, 300)), // Too long
			},
			mockResponse: map[string]interface{}{
				"error_code": "VALIDATION_ERROR",
				"message":    "Description too long",
			},
			mockStatusCode: http.StatusBadRequest,
			expectError:    true,
		},
		{
			name:    "not found - 404",
			fipID:   "fip-999",
			request: &floatingips.FloatingIPUpdateRequest{Description: "Test"},
			mockResponse: map[string]interface{}{
				"error_code": "NOT_FOUND",
				"message":    "Floating IP not found",
			},
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:    "unauthorized - 401",
			fipID:   "fip-1",
			request: &floatingips.FloatingIPUpdateRequest{Description: "Test"},
			mockResponse: map[string]interface{}{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid token",
			},
			mockStatusCode: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:    "forbidden - 403",
			fipID:   "fip-1",
			request: &floatingips.FloatingIPUpdateRequest{Description: "Test"},
			mockResponse: map[string]interface{}{
				"error_code": "FORBIDDEN",
				"message":    "Access denied",
			},
			mockStatusCode: http.StatusForbidden,
			expectError:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatusCode)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			client := cloudsdk.NewClient(server.URL, "test-token")
			projectClient := client.Project("proj-123")
			vpsClient := projectClient.VPS()

			ctx := context.Background()
			fip, err := vpsClient.FloatingIPs().Update(ctx, tt.fipID, tt.request)

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
