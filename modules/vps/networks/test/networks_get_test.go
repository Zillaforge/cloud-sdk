package networks_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
)

// TestNetworksGet_Success verifies successful network retrieval
func TestNetworksGet_Success(t *testing.T) {
	mockResponse := map[string]interface{}{
		"id":          "net-123",
		"name":        "test-network",
		"description": "Test network",
		"cidr":        "10.0.0.0/24",
		"project_id":  "proj-123",
		"created_at":  "2025-01-01T00:00:00Z",
		"updated_at":  "2025-01-01T00:00:00Z",
	}

	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method and path
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/networks/net-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		// Send mock response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client
	client := cloudsdk.NewClient(server.URL, "test-token")
	projectClient := client.Project("proj-123")
	vpsClient := projectClient.VPS()

	// Execute test
	ctx := context.Background()
	resource, err := vpsClient.Networks().Get(ctx, "net-123")

	// Verify results
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource == nil {
		t.Fatal("expected non-nil network resource")
	}
	if resource.Network == nil {
		t.Fatal("expected non-nil network")
	}
	if resource.Network.ID != "net-123" {
		t.Errorf("expected ID 'net-123', got '%s'", resource.Network.ID)
	}
	if resource.Network.Name != "test-network" {
		t.Errorf("expected name 'test-network', got '%s'", resource.Network.Name)
	}
	if resource.Network.CIDR != "10.0.0.0/24" {
		t.Errorf("expected CIDR '10.0.0.0/24', got '%s'", resource.Network.CIDR)
	}

	// Verify resource has Ports() accessor
	if resource.Ports() == nil {
		t.Error("expected non-nil Ports() accessor")
	}
}

// TestNetworksGet_Errors verifies error handling for network retrieval
func TestNetworksGet_Errors(t *testing.T) {
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
			resource, err := vpsClient.Networks().Get(ctx, tt.networkID)

			// Verify error
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if resource != nil {
				t.Errorf("expected nil resource on error, got %+v", resource)
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
