package networks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// TestPortsClient_List tests the List method for ports
func TestPortsClient_List(t *testing.T) {
	tests := []struct {
		name           string
		networkID      string
		mockResponse   interface{}
		mockStatusCode int
		expectedCount  int
		wantErr        bool
		validate       func(*testing.T, []*networks.NetworkPort)
	}{
		{
			name:      "list ports successfully",
			networkID: "net-full-001",
			mockResponse: []map[string]interface{}{
				{
					"id":        "port-1",
					"addresses": []string{"10.0.1.10"},
					"server": map[string]interface{}{
						"id":         "srv-1",
						"name":       "app-1",
						"status":     "ACTIVE",
						"project_id": "proj-123",
						"user_id":    "user-1",
					},
				},
				{
					"id":        "port-2",
					"addresses": []string{"10.0.1.20", "10.0.1.21"},
					"server": map[string]interface{}{
						"id":     "srv-2",
						"name":   "app-2",
						"status": "BUILD",
					},
				},
			},
			mockStatusCode: http.StatusOK,
			expectedCount:  2,
			validate: func(t *testing.T, ports []*networks.NetworkPort) {
				if len(ports) != 2 {
					t.Fatalf("expected 2 ports, got %d", len(ports))
				}

				assertStringField(t, ports[0], "ID", "port-1")
				assertStringSliceField(t, ports[0], "Addresses", []string{"10.0.1.10"})
				server := requirePointerStructField(t, ports[0], "Server")
				assertStringField(t, server.Interface(), "ID", "srv-1")
				assertStringField(t, server.Interface(), "Status", "ACTIVE")

				assertStringField(t, ports[1], "ID", "port-2")
				assertStringSliceField(t, ports[1], "Addresses", []string{"10.0.1.20", "10.0.1.21"})
				serverField := requireStructField(t, ports[1], "Server")
				if serverField.Kind() != reflect.Ptr {
					t.Fatalf("expected server field to be pointer, got %s", serverField.Kind())
				}
				if serverField.IsNil() {
					t.Fatal("expected second port server to be populated")
				}
			},
		},
		{
			name:           "empty port list",
			networkID:      "net-empty",
			mockResponse:   []map[string]interface{}{},
			mockStatusCode: http.StatusOK,
			expectedCount:  0,
			validate: func(t *testing.T, ports []*networks.NetworkPort) {
				if len(ports) != 0 {
					t.Fatalf("expected 0 ports, got %d", len(ports))
				}
			},
		},
		{
			name:      "port without server",
			networkID: "net-detached",
			mockResponse: []map[string]interface{}{
				{
					"id":        "port-3",
					"addresses": []string{"192.168.1.100"},
				},
			},
			mockStatusCode: http.StatusOK,
			expectedCount:  1,
			validate: func(t *testing.T, ports []*networks.NetworkPort) {
				if len(ports) != 1 {
					t.Fatalf("expected 1 port, got %d", len(ports))
				}

				assertStringField(t, ports[0], "ID", "port-3")
				assertStringSliceField(t, ports[0], "Addresses", []string{"192.168.1.100"})

				serverField := requireStructField(t, ports[0], "Server")
				if serverField.Kind() != reflect.Ptr {
					t.Fatalf("expected server field to be pointer, got %s", serverField.Kind())
				}
				if !serverField.IsNil() {
					t.Fatalf("expected nil server pointer, got %#v", serverField.Interface())
				}
			},
		},
		{
			name:           "network not found",
			networkID:      "net-missing",
			mockStatusCode: http.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}

				expectedPath := "/api/v1/project/proj-123/networks/" + tt.networkID + "/ports"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
			client := &PortsClient{
				baseClient: baseClient,
				projectID:  "proj-123",
				networkID:  tt.networkID,
			}

			result, err := client.List(context.Background())

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Fatal("expected result, got nil")
			}

			if len(result) != tt.expectedCount {
				t.Errorf("expected %d ports, got %d", tt.expectedCount, len(result))
			}

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestPortsClient_PathConstruction tests that the correct path is constructed
func TestPortsClient_PathConstruction(t *testing.T) {
	projectID := "my-project"
	networkID := "my-network"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/api/v1/project/" + projectID + "/networks/" + networkID + "/ports"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode([]*networks.NetworkPort{})
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := &PortsClient{
		baseClient: baseClient,
		projectID:  projectID,
		networkID:  networkID,
	}

	_, err := client.List(context.Background())

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
