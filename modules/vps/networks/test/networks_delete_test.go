package networks_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

// TestNetworksDelete_Success verifies successful network deletion
func TestNetworksDelete_Success(t *testing.T) {
	networkID := "net-123"

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		expectedPath := "/vps/api/v1/project/proj-123/networks/" + networkID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path '%s', got '%s'", expectedPath, r.URL.Path)
		}
		if r.ContentLength > 0 {
			t.Errorf("expected delete request to have empty body, got content-length %d", r.ContentLength)
		}

		// Send 204 No Content response
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	// Create client
	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()

	// Execute test
	ctx := context.Background()
	err := vpsClient.Networks().Delete(ctx, networkID)

	// Verify no error
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestNetworksDelete_Errors verifies error handling for network deletion
func TestNetworksDelete_Errors(t *testing.T) {
	tests := []struct {
		name             string
		networkID        string
		statusCode       int
		responseBody     interface{}
		expectedErrorMsg string
	}{
		{
			name:       "unauthorized - 401",
			networkID:  "net-123",
			statusCode: http.StatusUnauthorized,
			responseBody: map[string]interface{}{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid or expired token",
			},
			expectedErrorMsg: "Invalid or expired token",
		},
		{
			name:       "forbidden - 403",
			networkID:  "net-123",
			statusCode: http.StatusForbidden,
			responseBody: map[string]interface{}{
				"error_code": "FORBIDDEN",
				"message":    "Insufficient permissions",
			},
			expectedErrorMsg: "Insufficient permissions",
		},
		{
			name:       "not found - 404",
			networkID:  "nonexistent-network",
			statusCode: http.StatusNotFound,
			responseBody: map[string]interface{}{
				"error_code": "NOT_FOUND",
				"message":    "Network not found",
			},
			expectedErrorMsg: "Network not found",
		},
		{
			name:       "ports attached conflict - 409",
			networkID:  "net-123",
			statusCode: http.StatusConflict,
			responseBody: map[string]interface{}{
				"error_code": "CONFLICT",
				"message":    "Cannot delete network with attached ports",
			},
			expectedErrorMsg: "Cannot delete network with attached ports",
		},
		{
			name:       "internal server error - 500",
			networkID:  "net-123",
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
			err := vpsClient.Networks().Delete(ctx, tt.networkID)

			// Verify error
			if err == nil {
				t.Fatal("expected error, got nil")
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
