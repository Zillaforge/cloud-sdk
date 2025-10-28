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

func TestContract_ListRouters_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := routers.RouterListResponse{
			Routers: []routers.Router{
				{
					ID:          "router-123",
					Name:        "main-router",
					Description: "Main project router",
					State:       true,
					Status:      "ACTIVE",
					ProjectID:   "proj-1",
					IsDefault:   true,
				},
				{
					ID:          "router-456",
					Name:        "backup-router",
					Description: "Backup router",
					State:       false,
					Status:      "DOWN",
					ProjectID:   "proj-1",
					IsDefault:   false,
				},
			},
			Total: 2,
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	result, err := routerClient.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Routers) != 2 {
		t.Errorf("expected 2 routers, got %d", len(result.Routers))
	}

	if result.Total != 2 {
		t.Errorf("expected total 2, got %d", result.Total)
	}

	if result.Routers[0].ID != "router-123" {
		t.Errorf("expected ID router-123, got %s", result.Routers[0].ID)
	}

	if result.Routers[0].Name != "main-router" {
		t.Errorf("expected name 'main-router', got %s", result.Routers[0].Name)
	}

	if !result.Routers[0].State {
		t.Error("expected router 0 state to be true (enabled)")
	}

	if result.Routers[1].State {
		t.Error("expected router 1 state to be false (disabled)")
	}
}

func TestContract_ListRouters_WithNameFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify name query parameter
		name := r.URL.Query().Get("name")
		if name != "main" {
			t.Errorf("expected name query param 'main', got '%s'", name)
		}

		response := routers.RouterListResponse{
			Routers: []routers.Router{
				{
					ID:          "router-123",
					Name:        "main-router",
					Description: "Main project router",
					State:       true,
					ProjectID:   "proj-1",
				},
			},
			Total: 1,
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	opts := &routers.ListRoutersOptions{
		Name: "main",
	}

	result, err := routerClient.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Routers) != 1 {
		t.Errorf("expected 1 router, got %d", len(result.Routers))
	}

	if result.Routers[0].Name != "main-router" {
		t.Errorf("expected filtered router name 'main-router', got %s", result.Routers[0].Name)
	}
}

func TestContract_ListRouters_EmptyResult(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := routers.RouterListResponse{
			Routers: []routers.Router{},
			Total:   0,
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	result, err := routerClient.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(result.Routers) != 0 {
		t.Errorf("expected 0 routers, got %d", len(result.Routers))
	}

	if result.Total != 0 {
		t.Errorf("expected total 0, got %d", result.Total)
	}
}

func TestContract_ListRouters_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "internal server error",
		})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	_, err := routerClient.List(context.Background(), nil)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
