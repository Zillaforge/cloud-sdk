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

func TestContract_SetStateRouter_Enable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify method
		if r.Method != http.MethodPost {
			t.Errorf("expected POST method, got %s", r.Method)
		}

		// Parse request body
		var req routers.RouterSetStateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if !req.State {
			t.Error("expected state to be true (enabled)")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	req := &routers.RouterSetStateRequest{
		State: true,
	}

	err := routerClient.SetState(context.Background(), "router-123", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestContract_SetStateRouter_Disable(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Parse request body
		var req routers.RouterSetStateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.State {
			t.Error("expected state to be false (disabled)")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-1").VPS()
	routerClient := vpsClient.Routers()

	req := &routers.RouterSetStateRequest{
		State: false,
	}

	err := routerClient.SetState(context.Background(), "router-123", req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestContract_SetStateRouter_NotFound(t *testing.T) {
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

	req := &routers.RouterSetStateRequest{
		State: true,
	}

	err := routerClient.SetState(context.Background(), "nonexistent", req)
	if err == nil {
		t.Fatal("expected error for not found router, got nil")
	}
}
