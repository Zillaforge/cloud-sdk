package snapshots

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/common"
	snapshotsmodel "github.com/Zillaforge/cloud-sdk/models/vps/snapshots"
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
	mockSnapshot := &snapshotsmodel.Snapshot{
		ID:        "snap-1",
		Name:      "snap-one",
		VolumeID:  "vol-123",
		Status:    snapshotsmodel.SnapshotStatusCreating,
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
		if r.URL.Path != "/api/v1/project/proj-123/snapshots" {
			t.Errorf("expected path /api/v1/project/proj-123/snapshots, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockSnapshot)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	req := &snapshotsmodel.CreateSnapshotRequest{Name: "snap-one", VolumeID: "vol-123"}
	resp, err := client.Create(ctx, req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if resp.ID != mockSnapshot.ID {
		t.Errorf("expected ID %s, got %s", mockSnapshot.ID, resp.ID)
	}
}

func TestClient_Create_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	_, err := client.Create(ctx, &snapshotsmodel.CreateSnapshotRequest{Name: "", VolumeID: ""})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
	}
}

func TestClient_List(t *testing.T) {
	mockList := &snapshotsmodel.SnapshotListResponse{Snapshots: []*snapshotsmodel.Snapshot{
		{ID: "snap-1", Name: "snap-one", VolumeID: "vol-123", ProjectID: "proj-1"},
	}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/snapshots" {
			t.Errorf("expected path /api/v1/project/proj-123/snapshots, got %s", r.URL.Path)
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
	if len(resp) != 1 || resp[0].ID != "snap-1" {
		t.Fatalf("unexpected response: %v", resp)
	}
}

func TestClient_List_Filter(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/snapshots" {
			t.Errorf("expected path /api/v1/project/proj-123/snapshots, got %s", r.URL.Path)
		}
		query := r.URL.Query()
		if query.Get("name") != "snap" || query.Get("volume_id") != "vol-1" || query.Get("user_id") != "user-1" || query.Get("status") != "creating" {
			t.Errorf("unexpected query params: %v", query)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&snapshotsmodel.SnapshotListResponse{Snapshots: []*snapshotsmodel.Snapshot{}})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	opts := &snapshotsmodel.ListSnapshotsOptions{Name: "snap", VolumeID: "vol-1", UserID: "user-1", Status: "creating"}
	_, err := client.List(ctx, opts)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Get(t *testing.T) {
	mockSnapshot := &snapshotsmodel.Snapshot{ID: "snap-1", Name: "snap-one", VolumeID: "vol-123", ProjectID: "proj-1"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/snapshots/snap-1" {
			t.Errorf("expected path /api/v1/project/proj-123/snapshots/snap-1, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockSnapshot)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	snap, err := client.Get(ctx, "snap-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.ID != "snap-1" {
		t.Fatalf("unexpected id: %s", snap.ID)
	}
}

func TestClient_Update(t *testing.T) {
	mockSnapshot := &snapshotsmodel.Snapshot{ID: "snap-1", Name: "new-name", VolumeID: "vol-123"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/snapshots/snap-1" {
			t.Errorf("expected path /api/v1/project/proj-123/snapshots/snap-1, got %s", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockSnapshot)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	updateReq := &snapshotsmodel.UpdateSnapshotRequest{Name: "new-name"}
	snap, err := client.Update(ctx, "snap-1", updateReq)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if snap.Name != "new-name" {
		t.Fatalf("unexpected name: %s", snap.Name)
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
		if r.URL.Path != "/api/v1/project/proj-123/snapshots/snap-1" {
			t.Errorf("expected path /api/v1/project/proj-123/snapshots/snap-1, got %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	if err := client.Delete(ctx, "snap-1"); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}
