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

func TestMembersClient_Create(t *testing.T) {
	mockPool := &lbmodels.Pool{
		ID:   "pool-1",
		Name: "backend",
		Members: []lbmodels.Member{
			{ID: "mem-1", Address: "10.0.0.1"},
			{ID: "mem-2", Address: "10.0.0.2"},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/pools/pool-1/members" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockPool)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &MembersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1", poolID: "pool-1"}

	req := &lbmodels.MemberCreateRequest{Address: "10.0.0.2"}
	pool, err := client.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if pool.ID != "pool-1" {
		t.Errorf("expected pool ID pool-1, got %s", pool.ID)
	}
	if len(pool.Members) != 2 {
		t.Errorf("expected 2 members, got %d", len(pool.Members))
	}
}

func TestMembersClient_Create_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := &MembersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1", poolID: "pool-1"}

	_, err := client.Create(context.Background(), &lbmodels.MemberCreateRequest{})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestMembersClient_BatchUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/pools/pool-1/members" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &MembersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1", poolID: "pool-1"}

	req := &lbmodels.MembersBatchUpdateRequest{
		Members: []lbmodels.MemberCreateRequest{
			{Address: "10.0.0.1", Name: "member-1"},
			{Address: "10.0.0.2", Name: "member-2"},
		},
	}
	if err := client.BatchUpdate(context.Background(), req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestMembersClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/loadbalancers/lb-1/pools/pool-1/members/mem-1" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := &MembersClient{baseClient: baseClient, projectID: "proj-123", lbID: "lb-1", poolID: "pool-1"}

	if err := client.Delete(context.Background(), "mem-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
