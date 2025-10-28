package networks_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// TestNetworksCreate_Success verifies successful network creation
func TestNetworksCreate_Success(t *testing.T) {
	tests := []struct {
		name            string
		request         *networks.NetworkCreateRequest
		mockResponse    interface{}
		validateRequest func(*testing.T, map[string]interface{})
		validateResult  func(*testing.T, *networks.Network)
	}{
		{
			name: "create network with all fields",
			request: &networks.NetworkCreateRequest{
				Name:        "test-network",
				Description: "Test network description",
				CIDR:        "10.0.0.0/24",
			},
			mockResponse: map[string]interface{}{
				"id":          "net-123",
				"name":        "test-network",
				"description": "Test network description",
				"cidr":        "10.0.0.0/24",
				"project_id":  "proj-123",
				"created_at":  "2025-01-01T00:00:00Z",
				"updated_at":  "2025-01-01T00:00:00Z",
			},
			validateRequest: func(t *testing.T, reqBody map[string]interface{}) {
				if reqBody["name"] != "test-network" {
					t.Errorf("expected name 'test-network', got '%v'", reqBody["name"])
				}
				if reqBody["description"] != "Test network description" {
					t.Errorf("expected description 'Test network description', got '%v'", reqBody["description"])
				}
				if reqBody["cidr"] != "10.0.0.0/24" {
					t.Errorf("expected CIDR '10.0.0.0/24', got '%v'", reqBody["cidr"])
				}
			},
			validateResult: func(t *testing.T, network *networks.Network) {
				if network.ID != "net-123" {
					t.Errorf("expected ID 'net-123', got '%s'", network.ID)
				}
				if network.Name != "test-network" {
					t.Errorf("expected name 'test-network', got '%s'", network.Name)
				}
				if network.CIDR != "10.0.0.0/24" {
					t.Errorf("expected CIDR '10.0.0.0/24', got '%s'", network.CIDR)
				}
			},
		},
		{
			name: "create network without description",
			request: &networks.NetworkCreateRequest{
				Name: "minimal-network",
				CIDR: "192.168.0.0/24",
			},
			mockResponse: map[string]interface{}{
				"id":         "net-456",
				"name":       "minimal-network",
				"cidr":       "192.168.0.0/24",
				"project_id": "proj-123",
				"created_at": "2025-01-01T00:00:00Z",
				"updated_at": "2025-01-01T00:00:00Z",
			},
			validateRequest: func(t *testing.T, reqBody map[string]interface{}) {
				if reqBody["name"] != "minimal-network" {
					t.Errorf("expected name 'minimal-network', got '%v'", reqBody["name"])
				}
				if reqBody["description"] != nil && reqBody["description"] != "" {
					t.Errorf("expected empty description, got '%v'", reqBody["description"])
				}
			},
			validateResult: func(t *testing.T, network *networks.Network) {
				if network.Name != "minimal-network" {
					t.Errorf("expected name 'minimal-network', got '%s'", network.Name)
				}
				if network.Description != "" {
					t.Errorf("expected empty description, got '%s'", network.Description)
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
				if r.URL.Path != "/vps/api/v1/project/proj-123/networks" {
					t.Errorf("unexpected path: %s", r.URL.Path)
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
			network, err := vpsClient.Networks().Create(ctx, tt.request)

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

// TestNetworksCreate_Errors verifies error handling for network creation
func TestNetworksCreate_Errors(t *testing.T) {
	tests := []struct {
		name             string
		request          *networks.NetworkCreateRequest
		statusCode       int
		responseBody     interface{}
		expectedErrorMsg string
	}{
		{
			name: "validation error - 400",
			request: &networks.NetworkCreateRequest{
				Name: "", // invalid empty name
				CIDR: "10.0.0.0/24",
			},
			statusCode: http.StatusBadRequest,
			responseBody: map[string]interface{}{
				"error_code": "VALIDATION_ERROR",
				"message":    "Name is required",
			},
			expectedErrorMsg: "Name is required",
		},
		{
			name: "invalid CIDR format - 400",
			request: &networks.NetworkCreateRequest{
				Name: "test-network",
				CIDR: "invalid-cidr",
			},
			statusCode: http.StatusBadRequest,
			responseBody: map[string]interface{}{
				"error_code": "VALIDATION_ERROR",
				"message":    "Invalid CIDR format",
			},
			expectedErrorMsg: "Invalid CIDR format",
		},
		{
			name: "unauthorized - 401",
			request: &networks.NetworkCreateRequest{
				Name: "test-network",
				CIDR: "10.0.0.0/24",
			},
			statusCode: http.StatusUnauthorized,
			responseBody: map[string]interface{}{
				"error_code": "UNAUTHORIZED",
				"message":    "Invalid or expired token",
			},
			expectedErrorMsg: "Invalid or expired token",
		},
		{
			name: "forbidden - 403",
			request: &networks.NetworkCreateRequest{
				Name: "test-network",
				CIDR: "10.0.0.0/24",
			},
			statusCode: http.StatusForbidden,
			responseBody: map[string]interface{}{
				"error_code": "FORBIDDEN",
				"message":    "Insufficient permissions",
			},
			expectedErrorMsg: "Insufficient permissions",
		},
		{
			name: "quota exceeded - 409",
			request: &networks.NetworkCreateRequest{
				Name: "test-network",
				CIDR: "10.0.0.0/24",
			},
			statusCode: http.StatusConflict,
			responseBody: map[string]interface{}{
				"error_code": "QUOTA_EXCEEDED",
				"message":    "Network quota exceeded",
			},
			expectedErrorMsg: "Network quota exceeded",
		},
		{
			name: "CIDR conflict - 409",
			request: &networks.NetworkCreateRequest{
				Name: "test-network",
				CIDR: "10.0.0.0/24",
			},
			statusCode: http.StatusConflict,
			responseBody: map[string]interface{}{
				"error_code": "CIDR_CONFLICT",
				"message":    "CIDR overlaps with existing network",
			},
			expectedErrorMsg: "CIDR overlaps with existing network",
		},
		{
			name: "internal server error - 500",
			request: &networks.NetworkCreateRequest{
				Name: "test-network",
				CIDR: "10.0.0.0/24",
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
			network, err := vpsClient.Networks().Create(ctx, tt.request)

			// Verify error
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if network != nil {
				t.Errorf("expected nil network on error, got %+v", network)
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
