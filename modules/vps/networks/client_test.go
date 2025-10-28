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

// TestNewClient tests the NewClient constructor
func TestNewClient(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	projectID := "proj-123"

	client := NewClient(baseClient, projectID)

	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.projectID != projectID {
		t.Errorf("expected projectID %s, got %s", projectID, client.projectID)
	}
	if client.baseClient == nil {
		t.Error("expected baseClient to be initialized")
	}
}

// TestClient_List tests the List method
func TestClient_List(t *testing.T) {
	tests := []struct {
		name           string
		opts           *networks.ListNetworksOptions
		mockResponse   *networks.NetworkListResponse
		mockStatusCode int
		wantErr        bool
		checkPath      string
	}{
		{
			name: "list all networks",
			opts: nil,
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "network-1", CIDR: "10.0.1.0/24"},
					{ID: "net-2", Name: "network-2", CIDR: "10.0.2.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks",
		},
		{
			name: "list with name filter",
			opts: &networks.ListNetworksOptions{Name: "test-network"},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "test-network", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?name=test-network",
		},
		{
			name:           "server error",
			opts:           nil,
			mockResponse:   nil,
			mockStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkPath != "" {
					actualPath := r.URL.Path
					if r.URL.RawQuery != "" {
						actualPath += "?" + r.URL.RawQuery
					}
					if actualPath != tt.checkPath {
						t.Errorf("expected path %s, got %s", tt.checkPath, actualPath)
					}
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
			client := NewClient(baseClient, "proj-123")

			result, err := client.List(context.Background(), tt.opts)

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

			if len(result.Networks) != len(tt.mockResponse.Networks) {
				t.Errorf("expected %d networks, got %d", len(tt.mockResponse.Networks), len(result.Networks))
			}
		})
	}
}

// TestClient_Create tests the Create method
func TestClient_Create(t *testing.T) {
	tests := []struct {
		name           string
		request        *networks.NetworkCreateRequest
		mockResponse   *networks.Network
		mockStatusCode int
		wantErr        bool
	}{
		{
			name: "create network successfully",
			request: &networks.NetworkCreateRequest{
				Name:        "test-network",
				Description: "Test network",
				CIDR:        "10.0.1.0/24",
			},
			mockResponse: &networks.Network{
				ID:          "net-123",
				Name:        "test-network",
				Description: "Test network",
				CIDR:        "10.0.1.0/24",
				ProjectID:   "proj-123",
			},
			mockStatusCode: http.StatusCreated,
			wantErr:        false,
		},
		{
			name: "create without description",
			request: &networks.NetworkCreateRequest{
				Name: "simple-network",
				CIDR: "10.0.2.0/24",
			},
			mockResponse: &networks.Network{
				ID:   "net-456",
				Name: "simple-network",
				CIDR: "10.0.2.0/24",
			},
			mockStatusCode: http.StatusCreated,
			wantErr:        false,
		},
		{
			name:           "server error",
			request:        &networks.NetworkCreateRequest{},
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}

				expectedPath := "/api/v1/project/proj-123/networks"
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
			client := NewClient(baseClient, "proj-123")

			result, err := client.Create(context.Background(), tt.request)

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

			if result.ID != tt.mockResponse.ID {
				t.Errorf("expected ID %s, got %s", tt.mockResponse.ID, result.ID)
			}
		})
	}
}

// TestClient_Get tests the Get method
func TestClient_Get(t *testing.T) {
	tests := []struct {
		name           string
		networkID      string
		mockResponse   *networks.Network
		mockStatusCode int
		wantErr        bool
	}{
		{
			name:      "get network successfully",
			networkID: "net-123",
			mockResponse: &networks.Network{
				ID:   "net-123",
				Name: "test-network",
				CIDR: "10.0.1.0/24",
			},
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
				expectedPath := "/api/v1/project/proj-123/networks/" + tt.networkID
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
			client := NewClient(baseClient, "proj-123")

			result, err := client.Get(context.Background(), tt.networkID)

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

			if result.Network.ID != tt.mockResponse.ID {
				t.Errorf("expected ID %s, got %s", tt.mockResponse.ID, result.Network.ID)
			}

			// Verify ports operations are available
			if result.Ports() == nil {
				t.Error("expected ports operations to be available")
			}
		})
	}
}

// TestClient_Update tests the Update method
func TestClient_Update(t *testing.T) {
	networkID := "net-123"
	request := &networks.NetworkUpdateRequest{
		Name:        "updated-network",
		Description: "Updated description",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/networks/" + networkID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&networks.Network{
			ID:          networkID,
			Name:        "updated-network",
			Description: "Updated description",
		})
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	result, err := client.Update(context.Background(), networkID, request)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if result.Name != "updated-network" {
		t.Errorf("expected name 'updated-network', got %s", result.Name)
	}
}

// TestClient_Delete tests the Delete method
func TestClient_Delete(t *testing.T) {
	networkID := "net-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/networks/" + networkID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	err := client.Delete(context.Background(), networkID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestNetworkResource_Ports tests that NetworkResource provides port operations
func TestNetworkResource_Ports(t *testing.T) {
	network := &networks.Network{
		ID:   "net-123",
		Name: "test-network",
		CIDR: "10.0.1.0/24",
	}

	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)

	resource := &NetworkResource{
		Network: network,
		portOps: &PortsClient{
			baseClient: baseClient,
			projectID:  "proj-123",
			networkID:  "net-123",
		},
	}

	if resource.Ports() == nil {
		t.Error("expected ports operations, got nil")
	}

	if resource.Network.ID != "net-123" {
		t.Errorf("expected network ID net-123, got %s", resource.Network.ID)
	}
}
