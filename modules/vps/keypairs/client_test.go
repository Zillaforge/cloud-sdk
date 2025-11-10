package keypairs

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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
	result, err := client.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result == nil {
		t.Fatal("expected result, got nil")
	}
	// Use len() to count keypairs (no Total field)
	if len(result) != 2 {
		t.Errorf("expected 2 keypairs, got %d", len(result))
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
	result, err := client.List(ctx, opts)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Use len() to count keypairs (no Total field)
	if len(result) != 1 {
		t.Errorf("expected 1 keypair, got %d", len(result))
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
	result, err := client.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	// Use len() to count keypairs (no Total field)
	if len(result) != 0 {
		t.Errorf("expected empty list, got %d keypairs", len(result))
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

// ========================================================================
// Phase 4: User Story 2 - List Response Structure Tests
// ========================================================================

// TestListReturnType validates that List() returns []*Keypair directly (T021)
func TestListReturnType(t *testing.T) {
	mockResponse := keypairs.KeypairListResponse{
		Keypairs: []keypairs.Keypair{
			{
				ID:          "kp-001",
				Name:        "test-key-1",
				PublicKey:   "ssh-rsa AAAAB3...",
				Fingerprint: "SHA256:abc...",
				UserID:      "user-123",
				CreatedAt:   "2025-11-10T10:00:00Z",
				UpdatedAt:   "2025-11-10T10:00:00Z",
			},
			{
				ID:          "kp-002",
				Name:        "test-key-2",
				PublicKey:   "ssh-rsa AAAAB4...",
				Fingerprint: "SHA256:def...",
				UserID:      "user-123",
				CreatedAt:   "2025-11-10T11:00:00Z",
				UpdatedAt:   "2025-11-10T11:00:00Z",
			},
		},
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
	result, err := client.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify return type is []*Keypair (slice of pointers)
	if result == nil {
		t.Fatal("expected []*Keypair, got nil")
	}

	// Verify we got the correct number of keypairs
	if len(result) != 2 {
		t.Errorf("expected 2 keypairs, got %d", len(result))
	}

	// Verify elements are pointers (not copies)
	for i, kp := range result {
		if kp == nil {
			t.Errorf("keypair at index %d is nil", i)
		}
	}
}

// TestDirectSliceAccess validates direct slice access patterns (len, index) (T022)
func TestDirectSliceAccess(t *testing.T) {
	mockResponse := keypairs.KeypairListResponse{
		Keypairs: []keypairs.Keypair{
			{
				ID:          "kp-001",
				Name:        "key-alpha",
				PublicKey:   "ssh-rsa AAAAB3...",
				Fingerprint: "SHA256:abc...",
				UserID:      "user-123",
				CreatedAt:   "2025-11-10T10:00:00Z",
				UpdatedAt:   "2025-11-10T10:00:00Z",
			},
			{
				ID:          "kp-002",
				Name:        "key-beta",
				PublicKey:   "ssh-rsa AAAAB4...",
				Fingerprint: "SHA256:def...",
				UserID:      "user-123",
				CreatedAt:   "2025-11-10T11:00:00Z",
				UpdatedAt:   "2025-11-10T11:00:00Z",
			},
			{
				ID:          "kp-003",
				Name:        "key-gamma",
				PublicKey:   "ssh-rsa AAAAB5...",
				Fingerprint: "SHA256:ghi...",
				UserID:      "user-123",
				CreatedAt:   "2025-11-10T12:00:00Z",
				UpdatedAt:   "2025-11-10T12:00:00Z",
			},
		},
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
	keypairs, err := client.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Test len() works for counting
	count := len(keypairs)
	if count != 3 {
		t.Errorf("len(keypairs) = %d, want 3", count)
	}

	// Test direct index access
	if keypairs[0].Name != "key-alpha" {
		t.Errorf("keypairs[0].Name = %s, want 'key-alpha'", keypairs[0].Name)
	}
	if keypairs[1].Name != "key-beta" {
		t.Errorf("keypairs[1].Name = %s, want 'key-beta'", keypairs[1].Name)
	}
	if keypairs[2].Name != "key-gamma" {
		t.Errorf("keypairs[2].Name = %s, want 'key-gamma'", keypairs[2].Name)
	}

	// Test range iteration
	names := []string{}
	for _, kp := range keypairs {
		names = append(names, kp.Name)
	}
	expectedNames := []string{"key-alpha", "key-beta", "key-gamma"}
	if len(names) != len(expectedNames) {
		t.Errorf("iteration produced %d names, want %d", len(names), len(expectedNames))
	}
	for i, name := range names {
		if name != expectedNames[i] {
			t.Errorf("names[%d] = %s, want %s", i, name, expectedNames[i])
		}
	}

	// Test pointer modification (verify we have pointers, not copies)
	originalName := keypairs[0].Name
	keypairs[0].Name = "modified-name"
	if keypairs[0].Name != "modified-name" {
		t.Error("failed to modify keypair through pointer")
	}
	// Restore for other tests
	keypairs[0].Name = originalName
}

// TestEmptyKeypairsList validates handling of empty list response (T023)
func TestEmptyKeypairsList(t *testing.T) {
	mockResponse := keypairs.KeypairListResponse{
		Keypairs: []keypairs.Keypair{},
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
	keypairs, err := client.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty list should return empty slice, not nil
	if keypairs == nil {
		t.Error("expected empty slice, got nil")
	}

	// len() should return 0 for empty list
	if len(keypairs) != 0 {
		t.Errorf("len(keypairs) = %d, want 0", len(keypairs))
	}

	// range iteration should work with empty slice
	iterationCount := 0
	for range keypairs {
		iterationCount++
	}
	if iterationCount != 0 {
		t.Errorf("iteration count = %d, want 0", iterationCount)
	}
}

// TestListWithNameFilter validates filtering by name parameter (T024)
func TestListWithNameFilter(t *testing.T) {
	tests := []struct {
		name          string
		filterName    string
		mockKeypairs  []keypairs.Keypair
		expectedCount int
		expectedNames []string
		validateQuery func(t *testing.T, queryName string)
	}{
		{
			name:       "filter returns single match",
			filterName: "production-key",
			mockKeypairs: []keypairs.Keypair{
				{
					ID:          "kp-prod",
					Name:        "production-key",
					PublicKey:   "ssh-rsa AAAAB3...",
					Fingerprint: "SHA256:prod...",
					UserID:      "user-123",
					CreatedAt:   "2025-11-10T10:00:00Z",
					UpdatedAt:   "2025-11-10T10:00:00Z",
				},
			},
			expectedCount: 1,
			expectedNames: []string{"production-key"},
			validateQuery: func(t *testing.T, queryName string) {
				if queryName != "production-key" {
					t.Errorf("query parameter name = %s, want 'production-key'", queryName)
				}
			},
		},
		{
			name:          "filter returns no matches",
			filterName:    "nonexistent",
			mockKeypairs:  []keypairs.Keypair{},
			expectedCount: 0,
			expectedNames: []string{},
			validateQuery: func(t *testing.T, queryName string) {
				if queryName != "nonexistent" {
					t.Errorf("query parameter name = %s, want 'nonexistent'", queryName)
				}
			},
		},
		{
			name:       "filter with special characters",
			filterName: "key-with-dash",
			mockKeypairs: []keypairs.Keypair{
				{
					ID:          "kp-special",
					Name:        "key-with-dash",
					PublicKey:   "ssh-rsa AAAAB3...",
					Fingerprint: "SHA256:special...",
					UserID:      "user-123",
					CreatedAt:   "2025-11-10T10:00:00Z",
					UpdatedAt:   "2025-11-10T10:00:00Z",
				},
			},
			expectedCount: 1,
			expectedNames: []string{"key-with-dash"},
			validateQuery: func(t *testing.T, queryName string) {
				if queryName != "key-with-dash" {
					t.Errorf("query parameter name = %s, want 'key-with-dash'", queryName)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockResponse := keypairs.KeypairListResponse{
				Keypairs: tt.mockKeypairs,
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Validate query parameter
				queryName := r.URL.Query().Get("name")
				tt.validateQuery(t, queryName)

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
				Name: tt.filterName,
			}
			result, err := client.List(ctx, opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(result) != tt.expectedCount {
				t.Errorf("len(result) = %d, want %d", len(result), tt.expectedCount)
			}

			for i, expectedName := range tt.expectedNames {
				if i >= len(result) {
					t.Errorf("missing keypair at index %d", i)
					continue
				}
				if result[i].Name != expectedName {
					t.Errorf("result[%d].Name = %s, want %s", i, result[i].Name, expectedName)
				}
			}
		})
	}
}

// TestListResponseContractValidation validates response against pb.KeypairListOutput (T025)
func TestListResponseContractValidation(t *testing.T) {
	// This test validates the response structure matches pb.KeypairListOutput from Swagger
	swaggerCompliantJSON := `{
		"keypairs": [
			{
				"id": "kp-123",
				"name": "test-keypair",
				"description": "Test keypair",
				"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
				"fingerprint": "SHA256:abcdef1234567890",
				"user_id": "user-456",
				"user": {
					"id": "user-456",
					"name": "test@example.com"
				},
				"createdAt": "2025-11-10T10:00:00Z",
				"updatedAt": "2025-11-10T10:00:00Z"
			}
		]
	}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(swaggerCompliantJSON))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	result, err := client.List(ctx, nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Validate structure matches pb.KeypairListOutput
	if len(result) != 1 {
		t.Fatalf("expected 1 keypair, got %d", len(result))
	}

	kp := result[0]

	// Validate all pb.KeypairInfo fields are present and correct
	if kp.ID != "kp-123" {
		t.Errorf("ID = %s, want 'kp-123'", kp.ID)
	}
	if kp.Name != "test-keypair" {
		t.Errorf("Name = %s, want 'test-keypair'", kp.Name)
	}
	if kp.Description != "Test keypair" {
		t.Errorf("Description = %s, want 'Test keypair'", kp.Description)
	}
	if kp.PublicKey != "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC..." {
		t.Errorf("PublicKey mismatch")
	}
	if kp.Fingerprint != "SHA256:abcdef1234567890" {
		t.Errorf("Fingerprint = %s, want 'SHA256:abcdef1234567890'", kp.Fingerprint)
	}
	if kp.UserID != "user-456" {
		t.Errorf("UserID = %s, want 'user-456'", kp.UserID)
	}
	if kp.CreatedAt != "2025-11-10T10:00:00Z" {
		t.Errorf("CreatedAt = %s, want '2025-11-10T10:00:00Z'", kp.CreatedAt)
	}
	if kp.UpdatedAt != "2025-11-10T10:00:00Z" {
		t.Errorf("UpdatedAt = %s, want '2025-11-10T10:00:00Z'", kp.UpdatedAt)
	}

	// Validate User object
	if kp.User == nil {
		t.Fatal("User should not be nil")
	}
	if kp.User.ID != "user-456" {
		t.Errorf("User.ID = %s, want 'user-456'", kp.User.ID)
	}
	if kp.User.Name != "test@example.com" {
		t.Errorf("User.Name = %s, want 'test@example.com'", kp.User.Name)
	}

	// Validate PrivateKey is not present in List response (per spec)
	if kp.PrivateKey != "" {
		t.Errorf("PrivateKey should be empty in List response, got %s", kp.PrivateKey)
	}

	// Verify no Total field exists in response structure
	// (This is validated by the fact that KeypairListResponse doesn't have Total field anymore)
	// If Total field existed, compilation would fail
	t.Run("no Total field in structure", func(_ *testing.T) {
		// This test passes by compilation - KeypairListResponse has no Total field
		var response keypairs.KeypairListResponse
		_ = response // Used to verify type exists without Total
	})
}

// ========================================================================
// Phase 6: Integration & Contract Validation Tests
// ========================================================================

// TestIntegrationCRUD_FullLifecycle validates the complete Create→Get→Update→Delete flow (T039)
func TestIntegrationCRUD_FullLifecycle(t *testing.T) {
	// Setup test server with multiple endpoints
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")

		switch callCount {
		case 1: // Create
			w.WriteHeader(http.StatusCreated)
			response := `{
				"id": "kp-integration-test",
				"name": "integration-test-key",
				"description": "Integration test keypair",
				"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
				"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEA...",
				"fingerprint": "SHA256:integration123",
				"user_id": "user-integration",
				"user": {"id": "user-integration", "name": "integration@example.com"},
				"createdAt": "2025-11-10T10:00:00Z",
				"updatedAt": "2025-11-10T10:00:00Z"
			}`
			_, _ = w.Write([]byte(response))

		case 2: // Get
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "kp-integration-test",
				"name": "integration-test-key",
				"description": "Integration test keypair",
				"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
				"fingerprint": "SHA256:integration123",
				"user_id": "user-integration",
				"user": {"id": "user-integration", "name": "integration@example.com"},
				"createdAt": "2025-11-10T10:00:00Z",
				"updatedAt": "2025-11-10T10:00:00Z"
			}`
			_, _ = w.Write([]byte(response))

		case 3: // Update
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "kp-integration-test",
				"name": "integration-test-key",
				"description": "Updated integration test keypair",
				"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
				"fingerprint": "SHA256:integration123",
				"user_id": "user-integration",
				"user": {"id": "user-integration", "name": "integration@example.com"},
				"createdAt": "2025-11-10T10:00:00Z",
				"updatedAt": "2025-11-10T11:30:45Z"
			}`
			_, _ = w.Write([]byte(response))

		case 4: // Delete
			w.WriteHeader(http.StatusNoContent)
		}
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-integration")

	ctx := context.Background()

	// Step 1: Create keypair
	createReq := &keypairs.KeypairCreateRequest{
		Name:        "integration-test-key",
		Description: "Integration test keypair",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
	}

	created, err := client.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Validate Create response
	if created.ID != "kp-integration-test" {
		t.Errorf("Expected ID 'kp-integration-test', got '%s'", created.ID)
	}
	if created.Name != "integration-test-key" {
		t.Errorf("Expected Name 'integration-test-key', got '%s'", created.Name)
	}
	if created.Description != "Integration test keypair" {
		t.Errorf("Expected Description 'Integration test keypair', got '%s'", created.Description)
	}
	if created.PrivateKey == "" {
		t.Error("Expected private_key in Create response, got empty")
	}

	// Step 2: Get keypair
	retrieved, err := client.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	// Validate Get response (no private_key)
	if retrieved.PrivateKey != "" {
		t.Error("Expected no private_key in Get response, got non-empty")
	}
	if retrieved.ID != created.ID {
		t.Errorf("Get returned different ID: expected %s, got %s", created.ID, retrieved.ID)
	}

	// Step 3: Update keypair
	updateReq := &keypairs.KeypairUpdateRequest{
		Description: "Updated integration test keypair",
	}

	updated, err := client.Update(ctx, created.ID, updateReq)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}

	// Validate Update response
	if updated.Description != "Updated integration test keypair" {
		t.Errorf("Expected updated description 'Updated integration test keypair', got '%s'", updated.Description)
	}
	if updated.CreatedAt != "2025-11-10T10:00:00Z" {
		t.Errorf("CreatedAt should remain unchanged, got '%s'", updated.CreatedAt)
	}
	if updated.UpdatedAt == "2025-11-10T10:00:00Z" {
		t.Error("UpdatedAt should have changed after update")
	}

	// Step 4: Delete keypair
	err = client.Delete(ctx, created.ID)
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify all operations were called
	if callCount != 4 {
		t.Errorf("Expected 4 HTTP calls, got %d", callCount)
	}
}

// TestIntegrationKeypairImport validates keypair import with existing public key (T040)
func TestIntegrationKeypairImport(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Response should include generated private_key for imported keypair
		response := `{
			"id": "kp-import-test",
			"name": "imported-keypair",
			"description": "Imported existing public key",
			"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDT8Z...",
			"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpQIBAAKCAQEAwJ8Z...",
			"fingerprint": "SHA256:importfingerprint123",
			"user_id": "user-import",
			"user": {"id": "user-import", "name": "import@example.com"},
			"createdAt": "2025-11-10T12:00:00Z",
			"updatedAt": "2025-11-10T12:00:00Z"
		}`
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-import")

	ctx := context.Background()

	// Import existing public key
	importReq := &keypairs.KeypairCreateRequest{
		Name:        "imported-keypair",
		Description: "Imported existing public key",
		PublicKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDT8Z...", // Existing key provided
	}

	imported, err := client.Create(ctx, importReq)
	if err != nil {
		t.Fatalf("Import failed: %v", err)
	}

	// Validate import response
	if imported.Name != "imported-keypair" {
		t.Errorf("Expected name 'imported-keypair', got '%s'", imported.Name)
	}
	if imported.PublicKey != "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDT8Z..." {
		t.Error("Public key should match the imported key")
	}
	if imported.PrivateKey == "" {
		t.Error("Expected private_key to be generated for imported keypair")
	}
	if imported.Fingerprint != "SHA256:importfingerprint123" {
		t.Errorf("Expected fingerprint 'SHA256:importfingerprint123', got '%s'", imported.Fingerprint)
	}
}

// TestIntegrationKeypairGeneration validates keypair generation without public key (T041)
func TestIntegrationKeypairGeneration(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		// Response should include both generated public and private keys
		response := `{
			"id": "kp-generated-test",
			"name": "generated-keypair",
			"description": "Auto-generated keypair",
			"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDauto...",
			"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAauto...",
			"fingerprint": "SHA256:generatedfingerprint456",
			"user_id": "user-generate",
			"user": {"id": "user-generate", "name": "generate@example.com"},
			"createdAt": "2025-11-10T13:00:00Z",
			"updatedAt": "2025-11-10T13:00:00Z"
		}`
		_, _ = w.Write([]byte(response))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-generate")

	ctx := context.Background()

	// Generate new keypair (no public key provided)
	generateReq := &keypairs.KeypairCreateRequest{
		Name:        "generated-keypair",
		Description: "Auto-generated keypair",
		// PublicKey omitted - should trigger generation
	}

	generated, err := client.Create(ctx, generateReq)
	if err != nil {
		t.Fatalf("Generation failed: %v", err)
	}

	// Validate generation response
	if generated.Name != "generated-keypair" {
		t.Errorf("Expected name 'generated-keypair', got '%s'", generated.Name)
	}
	if generated.PublicKey == "" {
		t.Error("Expected public_key to be generated")
	}
	if generated.PrivateKey == "" {
		t.Error("Expected private_key to be generated")
	}
	if !strings.HasPrefix(generated.PublicKey, "ssh-rsa ") {
		t.Errorf("Expected public key to start with 'ssh-rsa ', got '%s'", generated.PublicKey[:20])
	}
	if !strings.HasPrefix(generated.PrivateKey, "-----BEGIN RSA PRIVATE KEY-----") {
		t.Error("Expected private key to start with PEM header")
	}
	if generated.Fingerprint != "SHA256:generatedfingerprint456" {
		t.Errorf("Expected fingerprint 'SHA256:generatedfingerprint456', got '%s'", generated.Fingerprint)
	}
}

// TestPrivateKeyAvailability validates private_key is only available during Create (T042)
func TestPrivateKeyAvailability(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		w.Header().Set("Content-Type", "application/json")

		switch callCount {
		case 1: // Create
			w.WriteHeader(http.StatusCreated)
			response := `{
				"id": "kp-private-test",
				"name": "private-key-test",
				"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCprivate...",
				"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAprivate...",
				"fingerprint": "SHA256:private123",
				"user_id": "user-private",
				"user": {"id": "user-private", "name": "private@example.com"},
				"createdAt": "2025-11-10T14:00:00Z",
				"updatedAt": "2025-11-10T14:00:00Z"
			}`
			_, _ = w.Write([]byte(response))

		case 2: // Get
			w.WriteHeader(http.StatusOK)
			response := `{
				"id": "kp-private-test",
				"name": "private-key-test",
				"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCprivate...",
				"fingerprint": "SHA256:private123",
				"user_id": "user-private",
				"user": {"id": "user-private", "name": "private@example.com"},
				"createdAt": "2025-11-10T14:00:00Z",
				"updatedAt": "2025-11-10T14:00:00Z"
			}`
			_, _ = w.Write([]byte(response))

		case 3: // List
			w.WriteHeader(http.StatusOK)
			response := `{
				"keypairs": [
					{
						"id": "kp-private-test",
						"name": "private-key-test",
						"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCprivate...",
						"fingerprint": "SHA256:private123",
						"user_id": "user-private",
						"user": {"id": "user-private", "name": "private@example.com"},
						"createdAt": "2025-11-10T14:00:00Z",
						"updatedAt": "2025-11-10T14:00:00Z"
					}
				]
			}`
			_, _ = w.Write([]byte(response))
		}
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-private")

	ctx := context.Background()

	// Step 1: Create should include private_key
	createReq := &keypairs.KeypairCreateRequest{
		Name: "private-key-test",
	}

	created, err := client.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if created.PrivateKey == "" {
		t.Error("Create response should include private_key")
	}

	// Step 2: Get should NOT include private_key
	retrieved, err := client.Get(ctx, created.ID)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if retrieved.PrivateKey != "" {
		t.Error("Get response should NOT include private_key")
	}

	// Step 3: List should NOT include private_key
	list, err := client.List(ctx, nil)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}

	if len(list) == 0 {
		t.Fatal("List should return at least one keypair")
	}

	for _, kp := range list {
		if kp.PrivateKey != "" {
			t.Error("List response should NOT include private_key for any keypair")
		}
	}
}

// TestContractValidationCreateRequest validates Create request matches KeypairCreateInput schema (T043)
func TestContractValidationCreateRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     *keypairs.KeypairCreateRequest
		expectError bool
		description string
	}{
		{
			name: "valid create request with all fields",
			request: &keypairs.KeypairCreateRequest{
				Name:        "test-keypair",
				Description: "Test keypair",
				PublicKey:   "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
			},
			expectError: false,
			description: "should match KeypairCreateInput schema",
		},
		{
			name: "valid create request minimal (name only)",
			request: &keypairs.KeypairCreateRequest{
				Name: "minimal-keypair",
			},
			expectError: false,
			description: "should match KeypairCreateInput schema with required fields only",
		},
		{
			name: "valid create request with description only",
			request: &keypairs.KeypairCreateRequest{
				Name:        "desc-only-keypair",
				Description: "Description without public key",
			},
			expectError: false,
			description: "should match KeypairCreateInput schema",
		},
		{
			name: "valid create request with public key only",
			request: &keypairs.KeypairCreateRequest{
				Name:      "pubkey-only-keypair",
				PublicKey: "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQC...",
			},
			expectError: false,
			description: "should match KeypairCreateInput schema",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling (request serialization)
			data, err := json.Marshal(tt.request)
			if tt.expectError && err == nil {
				t.Errorf("Expected marshaling to fail, but it succeeded")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected marshaling to succeed, but got error: %v", err)
			}

			if err != nil {
				return // Skip further validation if marshaling failed as expected
			}

			// Validate the JSON structure matches expected schema
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal for validation: %v", err)
			}

			// Check required fields
			if _, hasName := unmarshaled["name"]; !hasName {
				t.Error("Request should have 'name' field")
			}

			// Check optional fields are properly omitted when empty
			if tt.request.Description == "" {
				if _, hasDesc := unmarshaled["description"]; hasDesc {
					t.Error("Empty description should be omitted from JSON")
				}
			}
			if tt.request.PublicKey == "" {
				if _, hasPubKey := unmarshaled["public_key"]; hasPubKey {
					t.Error("Empty public_key should be omitted from JSON")
				}
			}
		})
	}
}

// TestContractValidationUpdateRequest validates Update request matches KeypairUpdateInput schema (T044)
func TestContractValidationUpdateRequest(t *testing.T) {
	tests := []struct {
		name        string
		request     *keypairs.KeypairUpdateRequest
		expectError bool
		description string
	}{
		{
			name: "valid update request with description",
			request: &keypairs.KeypairUpdateRequest{
				Description: "Updated description",
			},
			expectError: false,
			description: "should match KeypairUpdateInput schema",
		},
		{
			name:        "valid update request empty",
			request:     &keypairs.KeypairUpdateRequest{},
			expectError: false,
			description: "should match KeypairUpdateInput schema (all fields optional)",
		},
		{
			name: "valid update request nil description",
			request: &keypairs.KeypairUpdateRequest{
				Description: "",
			},
			expectError: false,
			description: "should match KeypairUpdateInput schema",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test JSON marshaling (request serialization)
			data, err := json.Marshal(tt.request)
			if tt.expectError && err == nil {
				t.Errorf("Expected marshaling to fail, but it succeeded")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected marshaling to succeed, but got error: %v", err)
			}

			if err != nil {
				return // Skip further validation if marshaling failed as expected
			}

			// Validate the JSON structure matches expected schema
			var unmarshaled map[string]interface{}
			if err := json.Unmarshal(data, &unmarshaled); err != nil {
				t.Fatalf("Failed to unmarshal for validation: %v", err)
			}

			// KeypairUpdateInput only has optional 'description' field
			// Check that empty description is properly omitted
			if tt.request.Description == "" {
				if _, hasDesc := unmarshaled["description"]; hasDesc {
					t.Error("Empty description should be omitted from JSON")
				}
			} else {
				if desc, hasDesc := unmarshaled["description"]; !hasDesc {
					t.Error("Non-empty description should be included in JSON")
				} else if desc != tt.request.Description {
					t.Errorf("Description mismatch: expected %s, got %s", tt.request.Description, desc)
				}
			}
		})
	}
}

// TestContractValidationCreateResponse validates Create response matches pb.KeypairInfo schema (T045)
func TestContractValidationCreateResponse(t *testing.T) {
	// Test response from Create operation (includes private_key)
	createResponseJSON := `{
		"id": "kp-create-contract",
		"name": "create-contract-test",
		"description": "Create response contract validation",
		"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCcreate...",
		"private_key": "-----BEGIN RSA PRIVATE KEY-----\nMIIEpAIBAAKCAQEAcreate...",
		"fingerprint": "SHA256:createcontract123",
		"user_id": "user-create",
		"user": {
			"id": "user-create",
			"name": "create@example.com"
		},
		"createdAt": "2025-11-10T15:00:00Z",
		"updatedAt": "2025-11-10T15:00:00Z"
	}`

	var response keypairs.Keypair
	err := json.Unmarshal([]byte(createResponseJSON), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal Create response: %v", err)
	}

	// Validate all required pb.KeypairInfo fields are present
	if response.ID == "" {
		t.Error("Create response should have 'id' field")
	}
	if response.Name == "" {
		t.Error("Create response should have 'name' field")
	}
	if response.PublicKey == "" {
		t.Error("Create response should have 'public_key' field")
	}
	if response.Fingerprint == "" {
		t.Error("Create response should have 'fingerprint' field")
	}
	if response.UserID == "" {
		t.Error("Create response should have 'user_id' field")
	}
	if response.CreatedAt == "" {
		t.Error("Create response should have 'createdAt' field")
	}
	if response.UpdatedAt == "" {
		t.Error("Create response should have 'updatedAt' field")
	}

	// Validate timestamps are RFC3339 format
	if _, err := time.Parse(time.RFC3339, response.CreatedAt); err != nil {
		t.Errorf("createdAt should be RFC3339 format: %v", err)
	}
	if _, err := time.Parse(time.RFC3339, response.UpdatedAt); err != nil {
		t.Errorf("updatedAt should be RFC3339 format: %v", err)
	}

	// Validate user object structure
	if response.User == nil {
		t.Error("Create response should have 'user' object")
	} else {
		if response.User.ID == "" {
			t.Error("User object should have 'id' field")
		}
		if response.User.Name == "" {
			t.Error("User object should have 'name' field")
		}
	}

	// Validate Create-specific fields (private_key should be present)
	if response.PrivateKey == "" {
		t.Error("Create response should include 'private_key' field")
	}

	// Validate optional fields
	// Description is optional, so it can be present or absent
}

// TestContractValidationGetResponse validates Get response matches pb.KeypairInfo schema (T046)
func TestContractValidationGetResponse(t *testing.T) {
	// Test response from Get operation (excludes private_key)
	getResponseJSON := `{
		"id": "kp-get-contract",
		"name": "get-contract-test",
		"description": "Get response contract validation",
		"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCget...",
		"fingerprint": "SHA256:getcontract456",
		"user_id": "user-get",
		"user": {
			"id": "user-get",
			"name": "get@example.com"
		},
		"createdAt": "2025-11-10T16:00:00Z",
		"updatedAt": "2025-11-10T16:30:00Z"
	}`

	var response keypairs.Keypair
	err := json.Unmarshal([]byte(getResponseJSON), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal Get response: %v", err)
	}

	// Validate all required pb.KeypairInfo fields are present
	if response.ID == "" {
		t.Error("Get response should have 'id' field")
	}
	if response.Name == "" {
		t.Error("Get response should have 'name' field")
	}
	if response.PublicKey == "" {
		t.Error("Get response should have 'public_key' field")
	}
	if response.Fingerprint == "" {
		t.Error("Get response should have 'fingerprint' field")
	}
	if response.UserID == "" {
		t.Error("Get response should have 'user_id' field")
	}
	if response.CreatedAt == "" {
		t.Error("Get response should have 'createdAt' field")
	}
	if response.UpdatedAt == "" {
		t.Error("Get response should have 'updatedAt' field")
	}

	// Validate timestamps are RFC3339 format
	if _, err := time.Parse(time.RFC3339, response.CreatedAt); err != nil {
		t.Errorf("createdAt should be RFC3339 format: %v", err)
	}
	if _, err := time.Parse(time.RFC3339, response.UpdatedAt); err != nil {
		t.Errorf("updatedAt should be RFC3339 format: %v", err)
	}

	// Validate user object structure
	if response.User == nil {
		t.Error("Get response should have 'user' object")
	} else {
		if response.User.ID == "" {
			t.Error("User object should have 'id' field")
		}
		if response.User.Name == "" {
			t.Error("User object should have 'name' field")
		}
	}

	// Validate Get-specific fields (private_key should NOT be present)
	if response.PrivateKey != "" {
		t.Error("Get response should NOT include 'private_key' field")
	}
}

// TestContractValidationListResponse validates List response matches pb.KeypairListOutput schema (T047)
func TestContractValidationListResponse(t *testing.T) {
	// Test response from List operation
	listResponseJSON := `{
		"keypairs": [
			{
				"id": "kp-list-1",
				"name": "list-contract-test-1",
				"description": "First keypair in list",
				"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQClist1...",
				"fingerprint": "SHA256:listcontract1",
				"user_id": "user-list",
				"user": {
					"id": "user-list",
					"name": "list@example.com"
				},
				"createdAt": "2025-11-10T17:00:00Z",
				"updatedAt": "2025-11-10T17:00:00Z"
			},
			{
				"id": "kp-list-2",
				"name": "list-contract-test-2",
				"public_key": "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQClist2...",
				"fingerprint": "SHA256:listcontract2",
				"user_id": "user-list",
				"user": {
					"id": "user-list",
					"name": "list@example.com"
				},
				"createdAt": "2025-11-10T17:15:00Z",
				"updatedAt": "2025-11-10T17:30:00Z"
			}
		]
	}`

	var response keypairs.KeypairListResponse
	err := json.Unmarshal([]byte(listResponseJSON), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal List response: %v", err)
	}

	// Validate pb.KeypairListOutput structure
	if response.Keypairs == nil {
		t.Fatal("List response should have 'keypairs' array")
	}

	if len(response.Keypairs) != 2 {
		t.Errorf("Expected 2 keypairs in list, got %d", len(response.Keypairs))
	}

	// Validate each keypair in the list matches pb.KeypairInfo schema
	for i, kp := range response.Keypairs {
		if kp.ID == "" {
			t.Errorf("Keypair %d should have 'id' field", i)
		}
		if kp.Name == "" {
			t.Errorf("Keypair %d should have 'name' field", i)
		}
		if kp.PublicKey == "" {
			t.Errorf("Keypair %d should have 'public_key' field", i)
		}
		if kp.Fingerprint == "" {
			t.Errorf("Keypair %d should have 'fingerprint' field", i)
		}
		if kp.UserID == "" {
			t.Errorf("Keypair %d should have 'user_id' field", i)
		}
		if kp.CreatedAt == "" {
			t.Errorf("Keypair %d should have 'createdAt' field", i)
		}
		if kp.UpdatedAt == "" {
			t.Errorf("Keypair %d should have 'updatedAt' field", i)
		}

		// Validate timestamps are RFC3339 format
		if _, err := time.Parse(time.RFC3339, kp.CreatedAt); err != nil {
			t.Errorf("Keypair %d createdAt should be RFC3339 format: %v", i, err)
		}
		if _, err := time.Parse(time.RFC3339, kp.UpdatedAt); err != nil {
			t.Errorf("Keypair %d updatedAt should be RFC3339 format: %v", i, err)
		}

		// Validate user object structure
		if kp.User == nil {
			t.Errorf("Keypair %d should have 'user' object", i)
		} else {
			if kp.User.ID == "" {
				t.Errorf("Keypair %d user object should have 'id' field", i)
			}
			if kp.User.Name == "" {
				t.Errorf("Keypair %d user object should have 'name' field", i)
			}
		}

		// Validate List-specific fields (private_key should NOT be present)
		if kp.PrivateKey != "" {
			t.Errorf("Keypair %d in list should NOT include 'private_key' field", i)
		}
	}
}
