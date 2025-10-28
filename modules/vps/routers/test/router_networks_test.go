package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/routers"
)

func TestContract_ListRouterNetworks_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/vps/api/v1/project/proj-1/routers/router-123":
			// Handle Get router request
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(routers.Router{
				ID:   "router-123",
				Name: "test-router",
			})
		case "/vps/api/v1/project/proj-1/routers/router-123/networks":
			// Handle List networks request
			response := []routers.RouterNetwork{
				{
					NetworkID:   "net-123",
					NetworkName: "internal-net",
					SubnetID:    "subnet-123",
				},
				{
					NetworkID:   "net-456",
					NetworkName: "external-net",
					SubnetID:    "subnet-456",
				},
			}
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(response)
		default:
			t.Errorf("unexpected path: %s", r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerResource, err := vpsClient.Routers().Get(context.Background(), "router-123")
	if err != nil {
		t.Fatalf("failed to get router: %v", err)
	}

	networks, err := routerResource.Networks().List(context.Background())
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(networks) != 2 {
		t.Errorf("expected 2 networks, got %d", len(networks))
	}

	if networks[0].NetworkID != "net-123" {
		t.Errorf("expected network ID 'net-123', got '%s'", networks[0].NetworkID)
	}
}

func TestContract_AssociateRouterNetwork_Success(t *testing.T) {
	routerID := "router-123"
	var getCallCount int

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/routers/"+routerID:
			// Handle Get router request
			getCallCount++
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(routers.Router{
				ID:   routerID,
				Name: "test-router",
			})
		case r.Method == http.MethodPost:
			// Verify path for associate
			expectedPath := "/vps/api/v1/project/proj-1/routers/router-123/networks/net-123"
			if r.URL.Path != expectedPath {
				t.Errorf("expected path '%s', got '%s'", expectedPath, r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Logf("Unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerResource, err := vpsClient.Routers().Get(context.Background(), routerID)
	if err != nil {
		t.Fatalf("failed to get router: %v", err)
	}

	err = routerResource.Networks().Associate(context.Background(), "net-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestContract_DisassociateRouterNetwork_Success(t *testing.T) {
	routerID := "router-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/routers/"+routerID:
			// Handle Get router request
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(routers.Router{
				ID:   routerID,
				Name: "test-router",
			})
		case r.Method == http.MethodDelete:
			// Verify path for disassociate
			expectedPath := "/vps/api/v1/project/proj-1/routers/router-123/networks/net-123"
			if r.URL.Path != expectedPath {
				t.Errorf("expected path '%s', got '%s'", expectedPath, r.URL.Path)
			}
			w.WriteHeader(http.StatusNoContent)
		default:
			t.Logf("Unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerResource, err := vpsClient.Routers().Get(context.Background(), routerID)
	if err != nil {
		t.Fatalf("failed to get router: %v", err)
	}

	err = routerResource.Networks().Disassociate(context.Background(), "net-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestContract_RouterNetworks_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/vps/api/v1/project/proj-1/routers/nonexistent":
			// Router not found
			w.WriteHeader(http.StatusNotFound)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "router not found",
			})
		default:
			t.Logf("Unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()

	_, err := vpsClient.Routers().Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for not found router, got nil")
	}
}
