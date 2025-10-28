package networks

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		mockResponse   []*networks.NetworkPort
		mockStatusCode int
		wantErr        bool
	}{
		{
			name:      "list ports successfully",
			networkID: "net-123",
			mockResponse: []*networks.NetworkPort{
				{
					ID:        "port-1",
					NetworkID: "net-123",
					FixedIPs:  []string{"10.0.1.10"},
					MACAddr:   "fa:16:3e:11:22:33",
				},
				{
					ID:        "port-2",
					NetworkID: "net-123",
					FixedIPs:  []string{"10.0.1.20"},
					MACAddr:   "fa:16:3e:44:55:66",
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "empty port list",
			networkID:      "net-456",
			mockResponse:   []*networks.NetworkPort{},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
		},
		{
			name:           "network not found",
			networkID:      "net-999",
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

			if len(result) != len(tt.mockResponse) {
				t.Errorf("expected %d ports, got %d", len(tt.mockResponse), len(result))
			}

			// Verify first port if available
			if len(result) > 0 && len(tt.mockResponse) > 0 {
				if result[0].ID != tt.mockResponse[0].ID {
					t.Errorf("expected port ID %s, got %s", tt.mockResponse[0].ID, result[0].ID)
				}
				if result[0].NetworkID != tt.mockResponse[0].NetworkID {
					t.Errorf("expected network ID %s, got %s", tt.mockResponse[0].NetworkID, result[0].NetworkID)
				}
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
