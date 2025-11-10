package networks_test

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
	networksmod "github.com/Zillaforge/cloud-sdk/modules/vps/networks"
)

// TestNetworksUpdate_Success verifies successful network updates
func TestNetworksUpdate_Success(t *testing.T) {
	tests := []struct {
		name            string
		networkID       string
		request         *networks.NetworkUpdateRequest
		mockResponse    interface{}
		validateRequest func(*testing.T, map[string]interface{})
		validateResult  func(*testing.T, *networksmod.NetworkResource)
	}{
		{
			name:      "update name and description",
			networkID: "net-123",
			request: &networks.NetworkUpdateRequest{
				Name:        "updated-network",
				Description: "Updated description",
			},
			mockResponse: map[string]interface{}{
				"id":          "net-123",
				"name":        "updated-network",
				"description": "Updated description",
				"cidr":        "10.0.0.0/24",
				"project_id":  "proj-123",
				"gateway":     "10.0.0.1",
				"router_id":   "router-1",
				"shared":      true,
				"is_default":  false,
				"created_at":  "2025-01-01T00:00:00Z",
				"updated_at":  "2025-01-02T00:00:00Z",
			},
			validateRequest: func(t *testing.T, reqBody map[string]interface{}) {
				if reqBody["name"] != "updated-network" {
					t.Errorf("expected name 'updated-network', got '%v'", reqBody["name"])
				}
				if reqBody["description"] != "Updated description" {
					t.Errorf("expected description 'Updated description', got '%v'", reqBody["description"])
				}
				for _, key := range []string{"gateway", "router_id", "shared", "bonding"} {
					if _, exists := reqBody[key]; exists {
						t.Errorf("did not expect key %q in update payload", key)
					}
				}
			},
			validateResult: func(t *testing.T, network *networksmod.NetworkResource) {
				if network.Network.Name != "updated-network" {
					t.Errorf("expected name 'updated-network', got '%s'", network.Network.Name)
				}
				if network.Network.Description != "Updated description" {
					t.Errorf("expected description 'Updated description', got '%s'", network.Network.Description)
				}
				if network.Network.Gateway != "10.0.0.1" {
					t.Errorf("expected gateway '10.0.0.1', got '%s'", network.Network.Gateway)
				}
				if network.Network.RouterID != "router-1" {
					t.Errorf("expected router_id 'router-1', got '%s'", network.Network.RouterID)
				}
				if !network.Network.Shared {
					t.Errorf("expected shared to be true")
				}
				if network.Network.IsDefault {
					t.Errorf("expected is_default to remain false")
				}
			},
		},
		{
			name:      "update name only",
			networkID: "net-456",
			request: &networks.NetworkUpdateRequest{
				Name: "new-name",
			},
			mockResponse: map[string]interface{}{
				"id":          "net-456",
				"name":        "new-name",
				"description": "Original description",
				"cidr":        "192.168.0.0/24",
				"project_id":  "proj-123",
				"shared":      false,
				"is_default":  true,
				"created_at":  "2025-01-01T00:00:00Z",
				"updated_at":  "2025-01-02T00:00:00Z",
			},
			validateRequest: func(t *testing.T, reqBody map[string]interface{}) {
				if reqBody["name"] != "new-name" {
					t.Errorf("expected name 'new-name', got '%v'", reqBody["name"])
				}
				if _, exists := reqBody["router_id"]; exists {
					t.Errorf("did not expect router_id in update payload")
				}
			},
			validateResult: func(t *testing.T, network *networksmod.NetworkResource) {
				if network.Network.Name != "new-name" {
					t.Errorf("expected name 'new-name', got '%s'", network.Network.Name)
				}
				if network.Network.Shared {
					t.Errorf("expected shared to be false")
				}
				if !network.Network.IsDefault {
					t.Errorf("expected is_default to be true")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request method and path
				if r.Method != http.MethodPut {
					t.Errorf("expected PUT request, got %s", r.Method)
				}
				expectedPath := "/vps/api/v1/project/proj-123/networks/" + tt.networkID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path '%s', got '%s'", expectedPath, r.URL.Path)
				}

				// Verify Content-Type
				if r.Header.Get("Content-Type") != "application/json" {
					t.Errorf("expected Content-Type 'application/json', got '%s'", r.Header.Get("Content-Type"))
				}

				// Read and validate request body
				body, err := io.ReadAll(r.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}
				var reqBody map[string]interface{}
				if err := json.Unmarshal(body, &reqBody); err != nil {
					t.Fatalf("failed to parse request body: %v", err)
				}

				if tt.validateRequest != nil {
					tt.validateRequest(t, reqBody)
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
			network, err := vpsClient.Networks().Update(ctx, tt.networkID, tt.request)

			// Verify results
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if network == nil {
				t.Fatal("expected non-nil network")
			}

			// Run custom validations
			if tt.validateResult != nil {
				tt.validateResult(t, network)
			}
		})
	}
}

// TestNetworksUpdate_Errors verifies error handling for network updates
func TestNetworksUpdate_Errors(t *testing.T) {
	tests := []struct {
		name             string
		networkID        string
		request          *networks.NetworkUpdateRequest
		statusCode       int
		responseBody     interface{}
		expectedErrorMsg string
	}{
		{
			name:      "validation error - 400",
			networkID: "net-123",
			request: &networks.NetworkUpdateRequest{
				Name: "", // invalid empty name
			},
			statusCode: http.StatusBadRequest,
			responseBody: map[string]interface{}{
				"error_code": "VALIDATION_ERROR",
				"message":    "Name cannot be empty",
			},
			expectedErrorMsg: "Name cannot be empty",
		},
		{
			name:      "unauthorized - 401",
			networkID: "net-123",
			request: &networks.NetworkUpdateRequest{
				Name: "updated-name",
			},
			statusCode: http.StatusUnauthorized,
			responseBody: map[string]interface{}{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid or expired token",
			},
			expectedErrorMsg: "Invalid or expired token",
		},
		{
			name:      "forbidden - 403",
			networkID: "net-123",
			request: &networks.NetworkUpdateRequest{
				Name: "updated-name",
			},
			statusCode: http.StatusForbidden,
			responseBody: map[string]interface{}{
				"error_code": "FORBIDDEN",
				"message":    "Insufficient permissions",
			},
			expectedErrorMsg: "Insufficient permissions",
		},
		{
			name:      "not found - 404",
			networkID: "nonexistent-network",
			request: &networks.NetworkUpdateRequest{
				Name: "updated-name",
			},
			statusCode: http.StatusNotFound,
			responseBody: map[string]interface{}{
				"error_code": "NOT_FOUND",
				"message":    "Network not found",
			},
			expectedErrorMsg: "Network not found",
		},
		{
			name:      "internal server error - 500",
			networkID: "net-123",
			request: &networks.NetworkUpdateRequest{
				Name: "updated-name",
			},
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
			network, err := vpsClient.Networks().Update(ctx, tt.networkID, tt.request)

			// Verify error
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if network != nil {
				t.Errorf("expected nil network on error, got %+v", network)
			}

			// Verify error is SDKError
			var sdkErr *cloudsdk.SDKError
			if !errors.As(err, &sdkErr) {
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
