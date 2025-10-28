package floatingips_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

// TestFloatingIPsGet_Success verifies successful floating IP retrieval
func TestFloatingIPsGet_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"id":          "fip-1",
		"address":     "203.0.113.10",
		"status":      "ACTIVE",
		"project_id":  "proj-123",
		"port_id":     "port-1",
		"description": "Web server IP",
		"created_at":  "2025-01-01T00:00:00Z",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/floatingips/fip-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()

	ctx := context.Background()
	fip, err := vpsClient.FloatingIPs().Get(ctx, "fip-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fip == nil {
		t.Fatal("expected non-nil floating IP")
	}
	if fip.ID != "fip-1" {
		t.Errorf("expected ID 'fip-1', got '%s'", fip.ID)
	}
	if fip.Address != "203.0.113.10" {
		t.Errorf("expected address '203.0.113.10', got '%s'", fip.Address)
	}
	if fip.Status != "ACTIVE" {
		t.Errorf("expected status 'ACTIVE', got '%s'", fip.Status)
	}
}

// TestFloatingIPsGet_Errors verifies error handling for floating IP retrieval
func TestFloatingIPsGet_Errors(t *testing.T) {
	tests := []struct {
		name           string
		fipID          string
		mockResponse   interface{}
		mockStatusCode int
		expectError    bool
	}{
		{
			name:  "not found - 404",
			fipID: "fip-999",
			mockResponse: map[string]interface{}{
				"error_code": "NOT_FOUND",
				"message":    "Floating IP not found",
			},
			mockStatusCode: http.StatusNotFound,
			expectError:    true,
		},
		{
			name:  "unauthorized - 401",
			fipID: "fip-1",
			mockResponse: map[string]interface{}{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid token",
			},
			mockStatusCode: http.StatusUnauthorized,
			expectError:    true,
		},
		{
			name:  "forbidden - 403",
			fipID: "fip-1",
			mockResponse: map[string]interface{}{
				"error_code": "FORBIDDEN",
				"message":    "Access denied",
			},
			mockStatusCode: http.StatusForbidden,
			expectError:    true,
		},
		{
			name:  "internal server error - 500",
			fipID: "fip-1",
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
			fip, err := vpsClient.FloatingIPs().Get(ctx, tt.fipID)

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
