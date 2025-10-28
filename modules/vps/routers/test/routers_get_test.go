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

func TestContract_GetRouter_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify path contains router ID
		if r.URL.Path != "/vps/api/v1/project/proj-1/routers/router-123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		response := routers.Router{
			ID:          "router-123",
			Name:        "test-router",
			Description: "Test router",
			State:       true,
			Status:      "ACTIVE",
			ProjectID:   "proj-1",
			IsDefault:   false,
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	result, err := routerClient.Get(context.Background(), "router-123")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "router-123" {
		t.Errorf("expected ID 'router-123', got '%s'", result.ID)
	}

	if result.Name != "test-router" {
		t.Errorf("expected name 'test-router', got '%s'", result.Name)
	}

	if !result.State {
		t.Error("expected state to be true (enabled)")
	}
}

func TestContract_GetRouter_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "router not found",
		})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	_, err := routerClient.Get(context.Background(), "nonexistent")
	if err == nil {
		t.Fatal("expected error for not found router, got nil")
	}
}
