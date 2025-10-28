package networks_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
			networkID: "net-123",
			mockResponse: []map[string]interface{}{
				{
					"id":          "port-1",
					"network_id":  "net-123",
					"fixed_ips":   []string{"10.0.0.10"},
					"mac_address": "fa:16:3e:11:22:33",
					"server_id":   "srv-1",
				},
				{
					"id":          "port-2",
					"network_id":  "net-123",
					"fixed_ips":   []string{"10.0.0.20", "10.0.0.21"},
					"mac_address": "fa:16:3e:44:55:66",
					"server_id":   "srv-2",
				},
			},
			expectedCount: 2,
			validateResult: func(t *testing.T, ports []*networks.NetworkPort) {
				if len(ports) != 2 {
					t.Fatalf("expected 2 ports, got %d", len(ports))
				}
				if ports[0].ID != "port-1" {
					t.Errorf("expected first port ID 'port-1', got '%s'", ports[0].ID)
				}
				if ports[0].NetworkID != "net-123" {
					t.Errorf("expected network ID 'net-123', got '%s'", ports[0].NetworkID)
				}
				if len(ports[0].FixedIPs) != 1 || ports[0].FixedIPs[0] != "10.0.0.10" {
					t.Errorf("expected fixed IP '10.0.0.10', got %v", ports[0].FixedIPs)
				}
				if ports[1].ServerID != "srv-2" {
					t.Errorf("expected server ID 'srv-2', got '%s'", ports[1].ServerID)
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
					"id":          "port-3",
					"network_id":  "net-789",
					"fixed_ips":   []string{"192.168.1.100"},
					"mac_address": "fa:16:3e:77:88:99",
				},
			},
			expectedCount: 1,
			validateResult: func(t *testing.T, ports []*networks.NetworkPort) {
				if len(ports) != 1 {
					t.Fatalf("expected 1 port, got %d", len(ports))
				}
				if ports[0].ServerID != "" {
					t.Errorf("expected empty server ID, got '%s'", ports[0].ServerID)
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
					_ = json.NewEncoder(w).Encode(map[string]interface{}{
						"id":         tt.networkID,
						"name":       "test-network",
						"cidr":       "10.0.0.0/24",
						"project_id": "proj-123",
					})
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
