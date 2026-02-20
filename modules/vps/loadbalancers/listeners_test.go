package loadbalancers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	lbmodels "github.com/Zillaforge/cloud-sdk/models/vps/loadbalancers"
)

func TestListenersClient_List(t *testing.T) {
	mockResp := &lbmodels.ListenerListResponse{
		Listeners: []*lbmodels.Listener{
			{ID: "lis-1", Name: "http-listener", Protocol: lbmodels.ListenerProtocolHTTP, ProtocolPort: 80},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/listeners" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &ListenersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	listeners, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(listeners) != 1 {
		t.Fatalf("expected 1 listener, got %d", len(listeners))
	}
	if listeners[0].ID != "lis-1" {
		t.Errorf("expected ID lis-1, got %s", listeners[0].ID)
	}
}

func TestListenersClient_Get(t *testing.T) {
	mockListener := &lbmodels.Listener{ID: "lis-1", Protocol: lbmodels.ListenerProtocolHTTP, ProtocolPort: 80}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/listeners/lis-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockListener)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &ListenersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	listener, err := client.Get(context.Background(), "lis-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if listener.ID != "lis-1" {
		t.Errorf("expected ID lis-1, got %s", listener.ID)
	}
}

func TestListenersClient_Create(t *testing.T) {
	mockListener := &lbmodels.Listener{ID: "lis-new", Protocol: lbmodels.ListenerProtocolHTTP, ProtocolPort: 80}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/listeners" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockListener)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &ListenersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	req := &lbmodels.ListenerCreateRequest{Protocol: lbmodels.ListenerProtocolHTTP, ProtocolPort: 80}
	listener, err := client.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if listener.ID != "lis-new" {
		t.Errorf("expected ID lis-new, got %s", listener.ID)
	}
}

func TestListenersClient_Create_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := &ListenersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	_, err := client.Create(context.Background(), &lbmodels.ListenerCreateRequest{})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestListenersClient_Update(t *testing.T) {
	mockListener := &lbmodels.Listener{ID: "lis-1", Name: "updated", Protocol: lbmodels.ListenerProtocolHTTP, ProtocolPort: 80}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/listeners/lis-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockListener)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &ListenersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	req := &lbmodels.ListenerUpdateRequest{Name: "updated"}
	listener, err := client.Update(context.Background(), "lis-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if listener.Name != "updated" {
		t.Errorf("expected name updated, got %s", listener.Name)
	}
}

func TestListenersClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/listeners/lis-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &ListenersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	if err := client.Delete(context.Background(), "lis-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
