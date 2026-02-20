package loadbalancers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/common"
	"github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
	lbmodels "github.com/Zillaforge/cloud-sdk/models/vps/loadbalancers"
)

func TestNewClient(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	client := NewClient(baseClient, "proj-123")
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.projectID != "proj-123" {
		t.Errorf("expected projectID proj-123, got %s", client.projectID)
	}
}

func TestClient_List(t *testing.T) {
	mockList := &lbmodels.LoadBalancerListResponse{
		LoadBalancers: []*lbmodels.LoadBalancer{
			{ID: "lb-1", Name: "test-lb", Status: lbmodels.LoadBalancerStatusActive, NetworkID: "net-1"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers" {
			t.Errorf("expected path /api/v1/project/proj-123/loadbalancers, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockList)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	resources, err := client.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resources) != 1 {
		t.Fatalf("expected 1 resource, got %d", len(resources))
	}
	if resources[0].ID != "lb-1" {
		t.Errorf("expected ID lb-1, got %s", resources[0].ID)
	}
	// Verify sub-resource accessors are wired
	if resources[0].Listeners() == nil {
		t.Error("expected Listeners() to be non-nil")
	}
	if resources[0].Pools() == nil {
		t.Error("expected Pools() to be non-nil")
	}
}

func TestClient_List_WithFilter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("name") != "my-lb" {
			t.Errorf("expected name=my-lb, got %s", q.Get("name"))
		}
		if q.Get("detail") != "true" {
			t.Errorf("expected detail=true, got %s", q.Get("detail"))
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(&lbmodels.LoadBalancerListResponse{})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	opts := &lbmodels.ListLoadBalancersOptions{Name: "my-lb", Detail: true}
	_, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Get(t *testing.T) {
	mockLB := &lbmodels.LoadBalancer{
		ID:        "lb-1",
		Name:      "test-lb",
		Status:    lbmodels.LoadBalancerStatusActive,
		NetworkID: "net-1",
		Network:   common.IDName{ID: "net-1", Name: "test-net"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockLB)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	resource, err := client.Get(context.Background(), "lb-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource.ID != "lb-1" {
		t.Errorf("expected ID lb-1, got %s", resource.ID)
	}
}

func TestClient_Create(t *testing.T) {
	mockLB := &lbmodels.LoadBalancer{
		ID:        "lb-new",
		Name:      "my-lb",
		Status:    lbmodels.LoadBalancerStatusPendingCreate,
		NetworkID: "net-1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockLB)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	req := &lbmodels.LoadBalancerCreateRequest{Name: "my-lb", NetworkID: "net-1"}
	resource, err := client.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource.ID != "lb-new" {
		t.Errorf("expected ID lb-new, got %s", resource.ID)
	}
}

func TestClient_Create_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	_, err := client.Create(context.Background(), &lbmodels.LoadBalancerCreateRequest{})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestClient_Update(t *testing.T) {
	mockLB := &lbmodels.LoadBalancer{ID: "lb-1", Name: "updated-name", Status: lbmodels.LoadBalancerStatusActive}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockLB)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	req := &lbmodels.LoadBalancerUpdateRequest{Name: "updated-name"}
	resource, err := client.Update(context.Background(), "lb-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resource.Name != "updated-name" {
		t.Errorf("expected name updated-name, got %s", resource.Name)
	}
}

func TestClient_Update_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	_, err := client.Update(context.Background(), "lb-1", &lbmodels.LoadBalancerUpdateRequest{})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	if err := client.Delete(context.Background(), "lb-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_AssociateFloatingIP(t *testing.T) {
	mockFIP := &floatingips.FloatingIP{
		ID:     "fip-1",
		Status: floatingips.FloatingIPStatusActive,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/floatingip" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockFIP)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	req := &lbmodels.AssociateFloatingIPRequest{FIPID: "fip-1"}
	fip, err := client.AssociateFloatingIP(context.Background(), "lb-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fip.ID != "fip-1" {
		t.Errorf("expected FIP ID fip-1, got %s", fip.ID)
	}
}

func TestClient_AssociateFloatingIP_NewFIP(t *testing.T) {
	mockFIP := &floatingips.FloatingIP{ID: "fip-new", Status: floatingips.FloatingIPStatusPending}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockFIP)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	req := &lbmodels.AssociateFloatingIPRequest{} // empty FIPID → creates new
	fip, err := client.AssociateFloatingIP(context.Background(), "lb-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if fip.ID != "fip-new" {
		t.Errorf("expected FIP ID fip-new, got %s", fip.ID)
	}
}

func TestLoadBalancerResource_SubResourceAccessors(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	lb := &lbmodels.LoadBalancer{ID: "lb-1", Name: "test"}
	resource := newLoadBalancerResource(baseClient, "proj-123", lb)

	if resource.LoadBalancer != lb {
		t.Error("expected embedded LoadBalancer to match")
	}
	if resource.Listeners() == nil {
		t.Error("expected Listeners() to be non-nil")
	}
	if resource.Pools() == nil {
		t.Error("expected Pools() to be non-nil")
	}
}
