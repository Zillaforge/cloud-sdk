package keypairs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/keypairs"
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

// TestClient_List tests successful keypair listing
func TestClient_List(t *testing.T) {
	mockResponse := &keypairs.KeypairListResponse{
		Keypairs: []keypairs.Keypair{
			{
				ID:          "key-1",
				Name:        "keypair-one",
				Description: "First keypair",
				PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
				Fingerprint: "SHA256:abc123...",
				UserID:      "user-1",
			},
			{
				ID:          "key-2",
				Name:        "keypair-two",
				Description: "Second keypair",
				PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
				Fingerprint: "SHA256:def456...",
				UserID:      "user-1",
			},
		},
		Total: 2,
	}

	expectedPath := "/api/v1/project/proj-123/keypairs"

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
	if len(response.Keypairs) != 2 {
		t.Errorf("expected 2 keypairs, got %d", len(response.Keypairs))
	}
}

// TestClient_List_WithFilter tests listing with name filter
func TestClient_List_WithFilter(t *testing.T) {
	mockResponse := &keypairs.KeypairListResponse{
		Keypairs: []keypairs.Keypair{
			{
				ID:          "key-special",
				Name:        "special-keypair",
				Description: "Special keypair",
				PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
				Fingerprint: "SHA256:special...",
				UserID:      "user-1",
			},
		},
		Total: 1,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}

		// Verify query parameter
		nameParam := r.URL.Query().Get("name")
		if nameParam != "special" {
			t.Errorf("expected name query param 'special', got %s", nameParam)
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
	opts := &keypairs.ListKeypairsOptions{
		Name: "special",
	}
	response, err := client.List(ctx, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if response.Total != 1 {
		t.Errorf("expected total 1, got %d", response.Total)
	}
}

// TestClient_Create tests successful keypair creation
func TestClient_Create(t *testing.T) {
	mockKeypair := &keypairs.Keypair{
		ID:          "key-new",
		Name:        "new-keypair",
		Description: "New SSH key",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
		Fingerprint: "SHA256:new123...",
		UserID:      "user-1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST request, got %s", r.Method)
		}

		// Verify request body
		var req keypairs.KeypairCreateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Name != "new-keypair" {
			t.Errorf("expected name 'new-keypair', got %s", req.Name)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockKeypair)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	req := &keypairs.KeypairCreateRequest{
		Name:        "new-keypair",
		Description: "New SSH key",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
	}
	keypair, err := client.Create(ctx, req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keypair == nil {
		t.Fatal("expected keypair, got nil")
	}
	if keypair.ID != "key-new" {
		t.Errorf("expected ID key-new, got %s", keypair.ID)
	}
	if keypair.Name != "new-keypair" {
		t.Errorf("expected name new-keypair, got %s", keypair.Name)
	}
}

// TestClient_Get tests successful keypair retrieval
func TestClient_Get(t *testing.T) {
	mockKeypair := &keypairs.Keypair{
		ID:          "key-123",
		Name:        "my-keypair",
		Description: "My SSH key",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
		Fingerprint: "SHA256:abc123...",
		UserID:      "user-1",
	}

	expectedPath := "/api/v1/project/proj-123/keypairs/key-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockKeypair)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	keypair, err := client.Get(ctx, "key-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keypair == nil {
		t.Fatal("expected keypair, got nil")
	}
	if keypair.ID != "key-123" {
		t.Errorf("expected ID key-123, got %s", keypair.ID)
	}
	if keypair.Name != "my-keypair" {
		t.Errorf("expected name my-keypair, got %s", keypair.Name)
	}
}

// TestClient_Update tests successful keypair update
func TestClient_Update(t *testing.T) {
	mockKeypair := &keypairs.Keypair{
		ID:          "key-123",
		Name:        "my-keypair",
		Description: "Updated description",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2E...",
		Fingerprint: "SHA256:abc123...",
		UserID:      "user-1",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT request, got %s", r.Method)
		}

		// Verify request body
		var req keypairs.KeypairUpdateRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Description != "Updated description" {
			t.Errorf("expected description 'Updated description', got %s", req.Description)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockKeypair)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	req := &keypairs.KeypairUpdateRequest{
		Description: "Updated description",
	}
	keypair, err := client.Update(ctx, "key-123", req)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if keypair == nil {
		t.Fatal("expected keypair, got nil")
	}
	if keypair.Description != "Updated description" {
		t.Errorf("expected description 'Updated description', got %s", keypair.Description)
	}
}

// TestClient_Delete tests successful keypair deletion
func TestClient_Delete(t *testing.T) {
	expectedPath := "/api/v1/project/proj-123/keypairs/key-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE request, got %s", r.Method)
		}
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
	err := client.Delete(ctx, "key-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// TestClient_List_EmptyResult tests listing with no keypairs
func TestClient_List_EmptyResult(t *testing.T) {
	mockResponse := &keypairs.KeypairListResponse{
		Keypairs: []keypairs.Keypair{},
		Total:    0,
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
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
	if response.Total != 0 {
		t.Errorf("expected total 0, got %d", response.Total)
	}
	if len(response.Keypairs) != 0 {
		t.Errorf("expected empty list, got %d keypairs", len(response.Keypairs))
	}
}

// TestClient_Get_NotFound tests Get with non-existent keypair
func TestClient_Get_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Not Found",
			"message": "Keypair not found",
		})
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

// TestClient_Create_Conflict tests Create with duplicate name
func TestClient_Create_Conflict(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_ = json.NewEncoder(w).Encode(map[string]interface{}{
			"error":   "Conflict",
			"message": "Keypair with this name already exists",
		})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	req := &keypairs.KeypairCreateRequest{
		Name:      "existing-keypair",
		PublicKey: "ssh-rsa AAAAB3NzaC1yc2E...",
	}
	_, err := client.Create(ctx, req)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
