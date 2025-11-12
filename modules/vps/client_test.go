package vps

import (
	"net/http"
	"testing"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/types"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
)

// TestNewClient tests the NewClient constructor
func TestNewClient(t *testing.T) {
	baseURL := "https://api.example.com"
	token := "test-token"
	projectID := "proj-123"
	httpClient := &http.Client{Timeout: 10 * time.Second}

	client := NewClient(baseURL, token, projectID, httpClient, nil)

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

// TestClient_Networks tests that Networks() returns a NetworksClient
func TestClient_Networks(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, nil)

	networksClient := client.Networks()

	if networksClient == nil {
		t.Fatal("expected NetworksClient, got nil")
	}
}

// TestClient_FloatingIPs tests that FloatingIPs() returns a FloatingIPsClient
func TestClient_FloatingIPs(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, nil)

	floatingIPsClient := client.FloatingIPs()

	if floatingIPsClient == nil {
		t.Fatal("expected FloatingIPsClient, got nil")
	}
}

// TestClient_Flavors tests that Flavors() returns a FlavorsClient
func TestClient_Flavors(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, nil)

	flavorsClient := client.Flavors()

	if flavorsClient == nil {
		t.Fatal("expected FlavorsClient, got nil")
	}
}

// TestClient_Keypairs tests that Keypairs() returns a KeypairsClient
func TestClient_Keypairs(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, nil)

	keypairsClient := client.Keypairs()

	if keypairsClient == nil {
		t.Fatal("expected KeypairsClient, got nil")
	}
}

// TestClient_Quotas tests that Quotas() returns a QuotasClient
func TestClient_Quotas(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, nil)

	quotasClient := client.Quotas()

	if quotasClient == nil {
		t.Fatal("expected QuotasClient, got nil")
	}
}

// TestClient_SecurityGroups tests that SecurityGroups() returns a SecurityGroupsClient
func TestClient_SecurityGroups(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, nil)

	securityGroupsClient := client.SecurityGroups()

	if securityGroupsClient == nil {
		t.Fatal("expected SecurityGroupsClient, got nil")
	}
}

// TestClient_Servers tests that Servers() returns a ServersClient
func TestClient_Servers(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, nil)

	serversClient := client.Servers()

	if serversClient == nil {
		t.Fatal("expected ServersClient, got nil")
	}
}

// TestClient_AllAccessors tests all accessor methods in one test
func TestClient_AllAccessors(t *testing.T) {
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, nil)

	accessors := []struct {
		name   string
		result interface{}
	}{
		{"Networks", client.Networks()},
		{"FloatingIPs", client.FloatingIPs()},
		{"Flavors", client.Flavors()},
		{"Keypairs", client.Keypairs()},
		{"Quotas", client.Quotas()},
		{"SecurityGroups", client.SecurityGroups()},
		{"Servers", client.Servers()},
	}

	for _, accessor := range accessors {
		if accessor.result == nil {
			t.Errorf("%s() returned nil", accessor.name)
		}
	}
}

// TestClient_ProjectID tests that ProjectID() returns the correct project ID
func TestClient_ProjectID(t *testing.T) {
	projectID := "proj-123"
	client := NewClient("https://api.example.com", "test-token", projectID, &http.Client{}, nil)

	if client.ProjectID() != projectID {
		t.Errorf("expected project ID %s, got %s", projectID, client.ProjectID())
	}
}

// mockLogger implements types.Logger for testing
type mockLogger struct{}

func (m *mockLogger) Debug(_ string, _ ...interface{}) {}
func (m *mockLogger) Info(_ string, _ ...interface{})  {}
func (m *mockLogger) Warn(_ string, _ ...interface{})  {}
func (m *mockLogger) Error(_ string, _ ...interface{}) {}

// TestNewClient_WithLogger tests client creation with a logger
func TestNewClient_WithLogger(t *testing.T) {
	logger := &mockLogger{}
	client := NewClient("https://api.example.com", "test-token", "proj-123", &http.Client{}, logger)

	if client == nil {
		t.Fatal("expected client, got nil")
	}
	// Logger is passed to internal client, so we can't directly check it here
	// but we verify the client was created successfully
}

// TestNetworkListResponse tests the response structure
func TestNetworkListResponse(t *testing.T) {
	networkList := []*networks.Network{
		{
			ID:   "net-1",
			Name: "network-1",
			CIDR: "10.0.1.0/24",
		},
		{
			ID:   "net-2",
			Name: "network-2",
			CIDR: "10.0.2.0/24",
		},
	}

	response := &networks.NetworkListResponse{
		Networks: networkList,
	}

	if len(response.Networks) != 2 {
		t.Errorf("expected 2 networks, got %d", len(response.Networks))
	}
	if response.Networks[0].ID != "net-1" {
		t.Errorf("expected network ID net-1, got %s", response.Networks[0].ID)
	}
}

// TestNetworkCreateRequest tests the request structure
func TestNetworkCreateRequest(t *testing.T) {
	req := &networks.NetworkCreateRequest{
		Name:        "test-network",
		Description: "Test description",
		CIDR:        "10.0.1.0/24",
	}

	if req.Name != "test-network" {
		t.Errorf("expected name test-network, got %s", req.Name)
	}
	if req.CIDR != "10.0.1.0/24" {
		t.Errorf("expected CIDR 10.0.1.0/24, got %s", req.CIDR)
	}
}

// TestNetworkUpdateRequest tests the request structure
func TestNetworkUpdateRequest(t *testing.T) {
	name := "updated-network"
	desc := "Updated description"

	req := &networks.NetworkUpdateRequest{
		Name:        name,
		Description: desc,
	}

	if req.Name != name {
		t.Errorf("expected name %s, got %s", name, req.Name)
	}
	if req.Description != desc {
		t.Errorf("expected description %s, got %s", desc, req.Description)
	}
}

// TestListNetworksOptions tests the options structure
func TestListNetworksOptions(t *testing.T) {
	opts := &networks.ListNetworksOptions{
		Name: "test-filter",
	}

	if opts.Name != "test-filter" {
		t.Errorf("expected name filter test-filter, got %s", opts.Name)
	}
}

// TestNetworkPort tests the network port structure
func TestNetworkPort(t *testing.T) {
	port := &networks.NetworkPort{
		ID:        "port-123",
		Addresses: []string{"10.0.1.10"},
		Server: &networks.ServerSummary{
			ID:        "srv-123",
			Name:      "web",
			Status:    "ACTIVE",
			ProjectID: "proj-123",
			UserID:    "user-123",
		},
	}

	if port.ID != "port-123" {
		t.Errorf("expected port ID port-123, got %s", port.ID)
	}
	if len(port.Addresses) != 1 {
		t.Errorf("expected 1 address, got %d", len(port.Addresses))
	}
	if port.Server == nil {
		t.Fatal("expected non-nil server summary")
	}
}

// TestNetwork tests the network structure
func TestNetwork(t *testing.T) {
	network := &networks.Network{
		ID:          "net-123",
		Name:        "test-network",
		Description: "Test description",
		CIDR:        "10.0.1.0/24",
		ProjectID:   "proj-123",
		Gateway:     "10.0.1.1",
		Shared:      true,
		Bonding:     false,
		CreatedAt:   "2025-01-01T00:00:00Z",
		UpdatedAt:   "2025-01-02T00:00:00Z",
	}

	if network.ID != "net-123" {
		t.Errorf("expected ID net-123, got %s", network.ID)
	}
	if network.Name != "test-network" {
		t.Errorf("expected name test-network, got %s", network.Name)
	}
	if network.CIDR != "10.0.1.0/24" {
		t.Errorf("expected CIDR 10.0.1.0/24, got %s", network.CIDR)
	}
	if network.Gateway != "10.0.1.1" {
		t.Errorf("expected gateway 10.0.1.1, got %s", network.Gateway)
	}
	if !network.Shared {
		t.Error("expected shared flag true")
	}
}

var _ types.Logger = (*mockLogger)(nil)
