package floatingips_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

// TestFloatingIPsDisassociate_Success verifies successful floating IP disassociation
func TestFloatingIPsDisassociate_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/floatingips/fip-1/disassociate" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()

	ctx := context.Background()
	err := vpsClient.FloatingIPs().Disassociate(ctx, "fip-1")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestFloatingIPsDisassociate_Errors verifies error handling for floating IP disassociation
func TestFloatingIPsDisassociate_Errors(t *testing.T) {
	tests := []struct {
		name           string
		fipID          string
		mockResponse   interface{}
		mockStatusCode int
		expectError    bool
	}{
		{
			name:  "not associated - 400",
			fipID: "fip-1",
			mockResponse: map[string]interface{}{
				"error_code": "INVALID_STATE",
				"message":    "Floating IP is not associated with a port",
			},
			mockStatusCode: http.StatusBadRequest,
			expectError:    true,
		},
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
			err := vpsClient.FloatingIPs().Disassociate(ctx, tt.fipID)

			if !tt.expectError {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
				return
			}

			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
