package networks_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// TestNetworkPorts_List_Success verifies successful port listing
func TestNetworkPorts_List_Success(t *testing.T) {
	tests := []struct {
		name           string
		networkID      string
		mockResponse   interface{}
		expectedCount  int
		validateResult func(*testing.T, []*networks.NetworkPort)
	}{
		{
			name:      "list ports with results",
			networkID: "net-full-001",
			mockResponse: []map[string]interface{}{
				{
					"id":        "port-1",
					"addresses": []string{"10.0.0.10"},
					"server": map[string]interface{}{
						"id":         "srv-1",
						"name":       "frontend",
						"status":     "ACTIVE",
						"project_id": "proj-123",
						"user_id":    "user-1",
					},
				},
				{
					"id":        "port-2",
					"addresses": []string{"10.0.0.20", "10.0.0.21"},
				},
			},
			expectedCount: 2,
			validateResult: func(t *testing.T, ports []*networks.NetworkPort) {
				if len(ports) != 2 {
					t.Fatalf("expected 2 ports, got %d", len(ports))
				}
				assertStringField(t, ports[0], "ID", "port-1")
				assertStringSliceField(t, ports[0], "Addresses", []string{"10.0.0.10"})
				server := requirePointerStructField(t, ports[0], "Server")
				assertStringField(t, server.Interface(), "ID", "srv-1")
				assertStringField(t, server.Interface(), "Status", "ACTIVE")

				assertStringField(t, ports[1], "ID", "port-2")
				assertStringSliceField(t, ports[1], "Addresses", []string{"10.0.0.20", "10.0.0.21"})
				serverField := requireStructField(t, ports[1], "Server")
				if serverField.Kind() != reflect.Ptr {
					t.Fatalf("expected server field to be pointer, got %s", serverField.Kind())
				}
				if !serverField.IsNil() {
					t.Fatalf("expected nil server pointer for second port, got %#v", serverField.Interface())
				}
			},
		},
		{
			name:          "empty port list",
			networkID:     "net-456",
			mockResponse:  []map[string]interface{}{},
			expectedCount: 0,
			validateResult: func(t *testing.T, ports []*networks.NetworkPort) {
				if len(ports) != 0 {
					t.Errorf("expected 0 ports, got %d", len(ports))
				}
			},
		},
		{
			name:      "port without server",
			networkID: "net-789",
			mockResponse: []map[string]interface{}{
				{
					"id":        "port-3",
					"addresses": []string{"192.168.1.100"},
				},
			},
			expectedCount: 1,
			validateResult: func(t *testing.T, ports []*networks.NetworkPort) {
				if len(ports) != 1 {
					t.Fatalf("expected 1 port, got %d", len(ports))
				}
				assertStringField(t, ports[0], "ID", "port-3")
				assertStringSliceField(t, ports[0], "Addresses", []string{"192.168.1.100"})

				serverField := requireStructField(t, ports[0], "Server")
				if serverField.Kind() != reflect.Ptr {
					t.Fatalf("expected server pointer, got %s", serverField.Kind())
				}
				if !serverField.IsNil() {
					t.Fatalf("expected nil server pointer, got %#v", serverField.Interface())
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock HTTP server that handles both network GET and ports LIST
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				// Handle network GET request
				networkPath := "/vps/api/v1/project/proj-123/networks/" + tt.networkID
				if r.URL.Path == networkPath && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusOK)
					fixture := loadFixtureBytes(t, "network_full.json")
					_, _ = w.Write(fixture)
					return
				}

				// Handle ports LIST request
				portsPath := "/vps/api/v1/project/proj-123/networks/" + tt.networkID + "/ports"
				if r.URL.Path == portsPath && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
					return
				}

				// Unexpected request
				t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			// Create client and get network resource
			ctx := context.Background()
			client := cloudsdk.NewClient(server.URL, "test-token")

			// Get the network resource
			resource, err := client.Project("proj-123").VPS().Networks().Get(ctx, tt.networkID)
			if err != nil {
				t.Fatalf("failed to get network resource: %v", err)
			}

			// Execute test - list ports on the network
			ports, err := resource.Ports().List(ctx)

			// Verify results
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if ports == nil {
				t.Fatal("expected non-nil ports slice")
			}
			if len(ports) != tt.expectedCount {
				t.Errorf("expected %d ports, got %d", tt.expectedCount, len(ports))
			}

			// Run custom validations
			if tt.validateResult != nil {
				tt.validateResult(t, ports)
			}
		})
	}
}

// TestNetworkPorts_List_Errors verifies error handling for port listing
func TestNetworkPorts_List_Errors(t *testing.T) {
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
			name:       "network not found - 404",
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
			// Create mock HTTP server for ports endpoint
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.statusCode)
				_ = json.NewEncoder(w).Encode(tt.responseBody)
			}))
			defer server.Close()

			// Get network resource (mock the Get and error List operations)
			ctx := context.Background()

			// Server already handles the error response for ports LIST
			// We just need it to also handle network GET
			multiServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")

				// Handle network GET
				networkPath := "/vps/api/v1/project/proj-123/networks/" + tt.networkID
				if r.URL.Path == networkPath && r.Method == http.MethodGet {
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"id":         tt.networkID,
						"name":       "test-network",
						"cidr":       "10.0.0.0/24",
						"project_id": "proj-123",
					})
					return
				}

				// Handle ports LIST with error
				portsPath := "/vps/api/v1/project/proj-123/networks/" + tt.networkID + "/ports"
				if r.URL.Path == portsPath {
					w.WriteHeader(tt.statusCode)
					_ = json.NewEncoder(w).Encode(tt.responseBody)
					return
				}

				w.WriteHeader(http.StatusNotFound)
			}))
			defer multiServer.Close()

			client := cloudsdk.NewClient(multiServer.URL, "test-token")
			resource, err := client.Project("proj-123").VPS().Networks().Get(ctx, tt.networkID)
			if err != nil {
				t.Fatalf("failed to get network resource: %v", err)
			}

			// Execute test - list ports (should error)
			ports, err := resource.Ports().List(ctx)

			// Verify error
			if err == nil {
				t.Fatal("expected error, got nil")
			}
			if ports != nil {
				t.Errorf("expected nil ports on error, got %+v", ports)
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
