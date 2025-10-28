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

func TestContract_CreateRouter_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Parse request body
		var req routers.RouterCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Name != "test-router" {
			t.Errorf("expected name 'test-router', got '%s'", req.Name)
		}

		response := routers.Router{
			ID:          "router-new",
			Name:        req.Name,
			Description: req.Description,
			State:       true,
			Status:      "ACTIVE",
			ProjectID:   "proj-1",
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	req := &routers.RouterCreateRequest{
		Name:        "test-router",
		Description: "Test router description",
	}

	result, err := routerClient.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ID != "router-new" {
		t.Errorf("expected ID 'router-new', got '%s'", result.ID)
	}

	if result.Name != "test-router" {
		t.Errorf("expected name 'test-router', got '%s'", result.Name)
	}

	if result.Description != "Test router description" {
		t.Errorf("expected description 'Test router description', got '%s'", result.Description)
	}
}

func TestContract_CreateRouter_WithExtNetwork(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req routers.RouterCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.ExtNetworkID != "ext-net-123" {
			t.Errorf("expected ext_network_id 'ext-net-123', got '%s'", req.ExtNetworkID)
		}

		response := routers.Router{
			ID:           "router-new",
			Name:         req.Name,
			ExtNetworkID: req.ExtNetworkID,
			State:        true,
			Status:       "ACTIVE",
			ProjectID:    "proj-1",
		}
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	req := &routers.RouterCreateRequest{
		Name:         "gateway-router",
		ExtNetworkID: "ext-net-123",
	}

	result, err := routerClient.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if result.ExtNetworkID != "ext-net-123" {
		t.Errorf("expected ext_network_id 'ext-net-123', got '%s'", result.ExtNetworkID)
	}
}

func TestContract_CreateRouter_ValidationError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(map[string]string{
			"error": "name is required",
		})
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	req := &routers.RouterCreateRequest{
		Name: "", // Invalid: empty name
	}

	_, err := routerClient.Create(context.Background(), req)
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}
