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

func TestContract_UpdateRouter_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT method, got %s", r.Method)
		}

		// Parse request body
		var req routers.RouterUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Name == nil || *req.Name != "updated-router" {
			t.Error("expected name to be 'updated-router'")
		}

		response := routers.Router{
			ID:          "router-123",
			Name:        *req.Name,
			Description: *req.Description,
			State:       true,
			Status:      "ACTIVE",
			ProjectID:   "proj-1",
		}
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	name := "updated-router"
	desc := "Updated description"
	req := &routers.RouterUpdateRequest{
		Name:        &name,
		Description: &desc,
	}

	result, err := routerClient.Update(context.Background(), "router-123", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.Name != "updated-router" {
		t.Errorf("expected name 'updated-router', got '%s'", result.Name)
	}

	if result.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got '%s'", result.Description)
	}
}

func TestContract_UpdateRouter_NotFound(t *testing.T) {
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

	name := "updated"
	req := &routers.RouterUpdateRequest{
		Name: &name,
	}

	_, err := routerClient.Update(context.Background(), "nonexistent", req)
	if err == nil {
		t.Fatal("expected error for not found router, got nil")
	}
}
