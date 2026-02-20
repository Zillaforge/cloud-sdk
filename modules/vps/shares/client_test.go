package shares

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/common"
	sharesmodel "github.com/Zillaforge/cloud-sdk/models/vps/shares"
)

func TestNewClient(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	projectID := "proj-123"
	client := NewClient(baseClient, projectID)
	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.projectID != projectID {
		t.Errorf("expected projectID %s, got %s", projectID, client.projectID)
	}
}

func TestClient_Create(t *testing.T) {
	mockShare := &sharesmodel.Share{
		ID:        "share-1",
		Name:      "test-share",
		Size:      10,
		Status:    sharesmodel.ShareStatusCreating,
		Network:   common.IDName{ID: "net-1", Name: "test-net"},
		NetworkID: "net-1",
		Project:   common.IDName{ID: "proj-1", Name: "test-project"},
		ProjectID: "proj-1",
		User:      common.IDName{ID: "user-1", Name: "test-user"},
		UserID:    "user-1",
		Namespace: "default",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares" {
			t.Errorf("expected path /api/v1/project/proj-123/shares, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockShare)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	req := &sharesmodel.CreateShareRequest{Name: "test-share", NetworkID: "net-1", Size: 10}
	resp, err := client.Create(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != mockShare.ID {
		t.Errorf("expected ID %s, got %s", mockShare.ID, resp.ID)
	}
}

func TestClient_Create_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	_, err := client.Create(ctx, &sharesmodel.CreateShareRequest{Name: "", NetworkID: "", Size: 0})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestClient_List(t *testing.T) {
	mockList := &sharesmodel.ShareListResponse{Shares: []*sharesmodel.Share{
		{ID: "share-1", Name: "s1", Size: 10, ProjectID: "proj-1"},
	}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares" {
			t.Errorf("expected path /api/v1/project/proj-123/shares, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockList)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	resp, err := client.List(ctx, nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(resp) != 1 || resp[0].ID != "share-1" {
		t.Fatalf("unexpected response: %v", resp)
	}
}

func TestClient_List_Filter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares" {
			t.Errorf("expected path /api/v1/project/proj-123/shares, got %s", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("name") != "my-share" {
			t.Errorf("expected name=my-share, got %s", query.Get("name"))
		}
		if query.Get("user_id") != "user-1" {
			t.Errorf("expected user_id=user-1, got %s", query.Get("user_id"))
		}
		if query.Get("detail") != "true" {
			t.Errorf("expected detail=true, got %s", query.Get("detail"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&sharesmodel.ShareListResponse{Shares: []*sharesmodel.Share{}})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	opts := &sharesmodel.ListSharesOptions{Name: "my-share", UserID: "user-1", Detail: true}
	_, err := client.List(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Get(t *testing.T) {
	mockShare := &sharesmodel.Share{ID: "share-1", Name: "test-share", Size: 10, ProjectID: "proj-1"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares/share-1" {
			t.Errorf("expected path /api/v1/project/proj-123/shares/share-1, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockShare)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	share, err := client.Get(ctx, "share-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if share.ID != "share-1" {
		t.Fatalf("unexpected id: %s", share.ID)
	}
}

func TestClient_Update(t *testing.T) {
	mockShare := &sharesmodel.Share{ID: "share-1", Name: "new-name", Size: 10}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares/share-1" {
			t.Errorf("expected path /api/v1/project/proj-123/shares/share-1, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockShare)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	updateReq := &sharesmodel.UpdateShareRequest{Name: "new-name"}
	share, err := client.Update(ctx, "share-1", updateReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if share.Name != "new-name" {
		t.Fatalf("unexpected name: %s", share.Name)
	}
}

func TestClient_Update_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	_, err := client.Update(ctx, "share-1", &sharesmodel.UpdateShareRequest{})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares/share-1" {
			t.Errorf("expected path /api/v1/project/proj-123/shares/share-1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	if err := client.Delete(ctx, "share-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Extend(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/shares/share-1/extend" {
			t.Errorf("expected path /api/v1/project/proj-123/shares/share-1/extend, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	err := client.Extend(ctx, "share-1", &sharesmodel.ExtendShareRequest{NewSize: 100})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Extend_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	err := client.Extend(ctx, "share-1", &sharesmodel.ExtendShareRequest{NewSize: 0})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}
