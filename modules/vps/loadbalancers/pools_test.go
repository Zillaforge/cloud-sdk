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

func TestPoolsClient_List(t *testing.T) {
	mockResp := &lbmodels.PoolListResponse{
		Pools: []*lbmodels.Pool{
			{ID: "pool-1", Name: "backend-pool", Protocol: lbmodels.PoolProtocolHTTP, MemberPort: 8080, Method: lbmodels.PoolMethodRoundRobin},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/pools" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResp)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &PoolsClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	resources, err := client.List(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 1 {
		t.Fatalf("expected 1 pool, got %d", len(resources))
	}
	if resources[0].ID != "pool-1" {
		t.Errorf("expected ID pool-1, got %s", resources[0].ID)
	}
	// Verify Members accessor is wired
	if resources[0].Members() == nil {
		t.Error("expected Members() to be non-nil")
	}
}

func TestPoolsClient_Get(t *testing.T) {
	mockPool := &lbmodels.Pool{
		ID:       "pool-1",
		Name:     "backend-pool",
		Protocol: lbmodels.PoolProtocolHTTP,
		Method:   lbmodels.PoolMethodRoundRobin,
		Members: []lbmodels.Member{
			{ID: "mem-1", Address: "10.0.0.1"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/pools/pool-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockPool)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &PoolsClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	resource, err := client.Get(context.Background(), "pool-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource.ID != "pool-1" {
		t.Errorf("expected ID pool-1, got %s", resource.ID)
	}
	if len(resource.Pool.Members) != 1 {
		t.Errorf("expected 1 member, got %d", len(resource.Pool.Members))
	}
}

func TestPoolsClient_Create(t *testing.T) {
	mockPool := &lbmodels.Pool{ID: "pool-new", Name: "new-pool", Protocol: lbmodels.PoolProtocolHTTP, MemberPort: 8080, Method: lbmodels.PoolMethodRoundRobin}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/pools" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockPool)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &PoolsClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	req := &lbmodels.PoolCreateRequest{
		Name:       "new-pool",
		Protocol:   lbmodels.PoolProtocolHTTP,
		MemberPort: 8080,
		Method:     lbmodels.PoolMethodRoundRobin,
	}
	resource, err := client.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource.ID != "pool-new" {
		t.Errorf("expected ID pool-new, got %s", resource.ID)
	}
}

func TestPoolsClient_Create_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := &PoolsClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	_, err := client.Create(context.Background(), &lbmodels.PoolCreateRequest{})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestPoolsClient_Update(t *testing.T) {
	mockPool := &lbmodels.Pool{ID: "pool-1", Name: "updated-pool"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/pools/pool-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockPool)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &PoolsClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	req := &lbmodels.PoolUpdateRequest{Name: "updated-pool"}
	resource, err := client.Update(context.Background(), "pool-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource.Name != "updated-pool" {
		t.Errorf("expected name updated-pool, got %s", resource.Name)
	}
}

func TestPoolsClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/pools/pool-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &PoolsClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1"}

	if err := client.Delete(context.Background(), "pool-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestPoolResource_MembersAccessor(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	pool := &lbmodels.Pool{ID: "pool-1", Name: "test"}
	resource := newPoolResource(baseClient, "proj-123", "lb-1", pool)

	if resource.Pool != pool {
		t.Error("expected embedded Pool to match")
	}
	if resource.Members() == nil {
		t.Error("expected Members() to be non-nil")
	}
}
