package routers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/routers"
)

// TestNewClient tests the NewClient constructor
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
	if client.basePath != "/api/v1/project/proj-123" {
		t.Errorf("expected basePath /api/v1/project/proj-123, got %s", client.basePath)
	}
}

// TestClient_List tests successful router listing
func TestClient_List(t *testing.T) {
	mockResponse := &routers.RouterListResponse{
		Routers: []routers.Router{
			{
				ID:        "router-1",
				Name:      "main-router",
				State:     true,
				Status:    "ACTIVE",
				ProjectID: "proj-123",
			},
			{
				ID:        "router-2",
				Name:      "backup-router",
				State:     false,
				Status:    "DOWN",
				ProjectID: "proj-123",
			},
		},
		Total: 2,
	}

	expectedPath := "/api/v1/project/proj-123/routers"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	response, err := client.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response == nil {
		t.Fatal("expected response, got nil")
	}
	if response.Total != 2 {
		t.Errorf("expected total 2, got %d", response.Total)
	}
	if len(response.Routers) != 2 {
		t.Errorf("expected 2 routers, got %d", len(response.Routers))
	}
}

// TestClient_List_WithFilters tests listing with filters
func TestClient_List_WithFilters(t *testing.T) {
	mockResponse := &routers.RouterListResponse{
		Routers: []routers.Router{
			{
				ID:        "router-main",
				Name:      "main-router",
				State:     true,
				ProjectID: "proj-123",
			},
		},
		Total: 1,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify query parameters
		query := r.URL.Query()
		if query.Get("name") != "main" {
			t.Errorf("expected name query param 'main', got '%s'", query.Get("name"))
		}
		if query.Get("user_id") != "user-1" {
			t.Errorf("expected user_id query param 'user-1', got '%s'", query.Get("user_id"))
		}
		if query.Get("detail") != "true" {
			t.Errorf("expected detail query param 'true', got '%s'", query.Get("detail"))
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	opts := &routers.ListRoutersOptions{
		Name:   "main",
		UserID: "user-1",
		Detail: true,
	}

	ctx := context.Background()
	response, err := client.List(ctx, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(response.Routers) != 1 {
		t.Errorf("expected 1 router, got %d", len(response.Routers))
	}
}

// TestClient_List_Error tests error handling for List
func TestClient_List_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	_, err := client.List(ctx, nil)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_Create tests successful router creation
func TestClient_Create(t *testing.T) {
	mockRouter := &routers.Router{
		ID:          "router-new",
		Name:        "test-router",
		Description: "Test router",
		State:       true,
		Status:      "ACTIVE",
		ProjectID:   "proj-123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		var req routers.RouterCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Name != "test-router" {
			t.Errorf("expected name 'test-router', got '%s'", req.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockRouter)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	req := &routers.RouterCreateRequest{
		Name:        "test-router",
		Description: "Test router",
	}

	ctx := context.Background()
	result, err := client.Create(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "router-new" {
		t.Errorf("expected ID 'router-new', got '%s'", result.ID)
	}
	if result.Name != "test-router" {
		t.Errorf("expected name 'test-router', got '%s'", result.Name)
	}
}

// TestClient_Create_Error tests error handling for Create
func TestClient_Create_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	req := &routers.RouterCreateRequest{
		Name: "",
	}

	ctx := context.Background()
	_, err := client.Create(ctx, req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_Get tests successful router retrieval
func TestClient_Get(t *testing.T) {
	mockRouter := &routers.Router{
		ID:          "router-123",
		Name:        "test-router",
		Description: "Test router",
		State:       true,
		Status:      "ACTIVE",
		ProjectID:   "proj-123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/routers/router-123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockRouter)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	result, err := client.Get(ctx, "router-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected RouterResource, got nil")
	}
	if result.ID != "router-123" {
		t.Errorf("expected ID 'router-123', got '%s'", result.ID)
	}
	if result.Networks() == nil {
		t.Error("expected Networks() to return non-nil operations")
	}
}

// TestClient_Get_Error tests error handling for Get
func TestClient_Get_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	_, err := client.Get(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_Update tests successful router update
func TestClient_Update(t *testing.T) {
	mockRouter := &routers.Router{
		ID:          "router-123",
		Name:        "updated-router",
		Description: "Updated description",
		State:       true,
		ProjectID:   "proj-123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}

		var req routers.RouterUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if req.Name == nil || *req.Name != "updated-router" {
			t.Error("expected name to be 'updated-router'")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockRouter)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	name := "updated-router"
	desc := "Updated description"
	req := &routers.RouterUpdateRequest{
		Name:        &name,
		Description: &desc,
	}

	ctx := context.Background()
	result, err := client.Update(ctx, "router-123", req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "updated-router" {
		t.Errorf("expected name 'updated-router', got '%s'", result.Name)
	}
}

// TestClient_Update_Error tests error handling for Update
func TestClient_Update_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	name := "updated"
	req := &routers.RouterUpdateRequest{
		Name: &name,
	}

	ctx := context.Background()
	_, err := client.Update(ctx, "nonexistent", req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_Delete tests successful router deletion
func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/routers/router-123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	err := client.Delete(ctx, "router-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestClient_Delete_Error tests error handling for Delete
func TestClient_Delete_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	err := client.Delete(ctx, "nonexistent")

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

// TestClient_SetState tests successful router state change
func TestClient_SetState(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		var req routers.RouterSetStateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}

		if !req.State {
			t.Error("expected state to be true")
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	req := &routers.RouterSetStateRequest{
		State: true,
	}

	ctx := context.Background()
	err := client.SetState(ctx, "router-123", req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestClient_SetState_Error tests error handling for SetState
func TestClient_SetState_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	req := &routers.RouterSetStateRequest{
		State: true,
	}

	ctx := context.Background()
	err := client.SetState(ctx, "nonexistent", req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
