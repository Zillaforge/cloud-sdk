package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/routers"
)

// TestRouterNetworksClient_List tests successful network listing
func TestRouterNetworksClient_List(t *testing.T) {
	mockNetworks := []routers.RouterNetwork{
		{
			NetworkID:   "net-1",
			NetworkName: "network-1",
			SubnetID:    "subnet-1",
			PortID:      "port-1",
		},
		{
			NetworkID:   "net-2",
			NetworkName: "network-2",
			SubnetID:    "subnet-2",
			PortID:      "port-2",
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/routers/router-123/networks"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockNetworks)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := newRouterNetworksClient(baseClient, "proj-123", "router-123")

	ctx := context.Background()
	networks, err := client.List(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(networks) != 2 {
		t.Errorf("expected 2 networks, got %d", len(networks))
	}
	if networks[0].NetworkID != "net-1" {
		t.Errorf("expected NetworkID 'net-1', got '%s'", networks[0].NetworkID)
	}
}

// TestRouterNetworksClient_List_Error tests error handling for List
func TestRouterNetworksClient_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := newRouterNetworksClient(baseClient, "proj-123", "router-123")

	ctx := context.Background()
	_, err := client.List(ctx)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestRouterNetworksClient_Associate tests successful network association
func TestRouterNetworksClient_Associate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/routers/router-123/networks/net-456"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := newRouterNetworksClient(baseClient, "proj-123", "router-123")

	ctx := context.Background()
	err := client.Associate(ctx, "net-456")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
} // TestRouterNetworksClient_Associate_Error tests error handling for Associate
func TestRouterNetworksClient_Associate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := newRouterNetworksClient(baseClient, "proj-123", "router-123")

	ctx := context.Background()
	err := client.Associate(ctx, "invalid-net")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestRouterNetworksClient_Disassociate tests successful network disassociation
func TestRouterNetworksClient_Disassociate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/routers/router-123/networks/net-789"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := newRouterNetworksClient(baseClient, "proj-123", "router-123")

	ctx := context.Background()
	err := client.Disassociate(ctx, "net-789")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestRouterNetworksClient_Disassociate_Error tests error handling for Disassociate
func TestRouterNetworksClient_Disassociate_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := newRouterNetworksClient(baseClient, "proj-123", "router-123")

	ctx := context.Background()
	err := client.Disassociate(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestNewRouterNetworksClient tests the constructor
func TestNewRouterNetworksClient(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	projectID := "proj-123"
	routerID := "router-456"

	client := newRouterNetworksClient(baseClient, projectID, routerID)

	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.projectID != projectID {
		t.Errorf("expected projectID %s, got %s", projectID, client.projectID)
	}
	if client.routerID != routerID {
		t.Errorf("expected routerID %s, got %s", routerID, client.routerID)
	}
}
