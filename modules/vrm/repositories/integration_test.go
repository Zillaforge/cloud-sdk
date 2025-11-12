package repositories

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vrm/repositories"
)

// T047: Integration test - Repository CRUD lifecycle
// Verify Create -> Read -> Update -> Delete workflow
func TestRepositoryCRUDLifecycle(t *testing.T) {
	// Track request order
	requestLog := make([]string, 0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the request
		requestLog = append(requestLog, r.Method+" "+r.URL.Path)

		switch {
		// Create
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/project/proj-123/repository":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusCreated)
			_, _ = w.Write([]byte(`{
				"id": "repo-123",
				"name": "test-repo",
				"namespace": "public",
				"operatingSystem": "linux",
				"description": "Test repository",
				"count": 0,
				"creator": {"id": "user-1", "name": "admin"},
				"project": {"id": "proj-123", "name": "test-project"},
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z"
			}`))

		// Read (Get)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/project/proj-123/repository/repo-123":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": "repo-123",
				"name": "test-repo",
				"namespace": "public",
				"operatingSystem": "linux",
				"description": "Test repository",
				"count": 0,
				"creator": {"id": "user-1", "name": "admin"},
				"project": {"id": "proj-123", "name": "test-project"},
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z"
			}`))

		// Update
		case r.Method == http.MethodPut && r.URL.Path == "/api/v1/project/proj-123/repository/repo-123":
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"id": "repo-123",
				"name": "test-repo",
				"namespace": "public",
				"operatingSystem": "linux",
				"description": "Updated description",
				"count": 0,
				"creator": {"id": "user-1", "name": "admin"},
				"project": {"id": "proj-123", "name": "test-project"},
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-02T00:00:00Z"
			}`))

		// Delete
		case r.Method == http.MethodDelete && r.URL.Path == "/api/v1/project/proj-123/repository/repo-123":
			w.WriteHeader(http.StatusNoContent)

		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error": "not found"}`))
		}
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")
	ctx := context.Background()

	// Phase 1: Create
	createReq := &repositories.CreateRepositoryRequest{
		Name:            "test-repo",
		OperatingSystem: "linux",
		Description:     "Test repository",
	}
	repo, err := client.Create(ctx, createReq)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if repo.ID != "repo-123" {
		t.Errorf("Expected ID repo-123, got %s", repo.ID)
	}

	// Phase 2: Read
	repo, err = client.Get(ctx, "repo-123")
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}
	if repo.Description != "Test repository" {
		t.Errorf("Expected description 'Test repository', got '%s'", repo.Description)
	}

	// Phase 3: Update
	updateReq := &repositories.UpdateRepositoryRequest{
		Description: "Updated description",
	}
	repo, err = client.Update(ctx, "repo-123", updateReq)
	if err != nil {
		t.Fatalf("Update failed: %v", err)
	}
	if repo.Description != "Updated description" {
		t.Errorf("Expected description 'Updated description', got '%s'", repo.Description)
	}

	// Phase 4: Delete
	err = client.Delete(ctx, "repo-123")
	if err != nil {
		t.Fatalf("Delete failed: %v", err)
	}

	// Verify request sequence
	expectedSequence := []string{
		"POST /api/v1/project/proj-123/repository",
		"GET /api/v1/project/proj-123/repository/repo-123",
		"PUT /api/v1/project/proj-123/repository/repo-123",
		"DELETE /api/v1/project/proj-123/repository/repo-123",
	}
	if len(requestLog) != len(expectedSequence) {
		t.Errorf("Expected %d requests, got %d", len(expectedSequence), len(requestLog))
	}
	for i, expected := range expectedSequence {
		if i < len(requestLog) && requestLog[i] != expected {
			t.Errorf("Request %d: expected %s, got %s", i+1, expected, requestLog[i])
		}
	}
}

// T048: Integration test - Pagination workflow
// Verify List with pagination parameters
func TestRepositoryListPagination(t *testing.T) {
	tests := []struct {
		name   string
		limit  int
		offset int
		count  int
	}{
		{name: "first page", limit: 10, offset: 0, count: 10},
		{name: "second page", limit: 10, offset: 10, count: 10},
		{name: "partial page", limit: 10, offset: 20, count: 5},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Extract query parameters
				query := r.URL.Query()
				_ = query.Get("limit")
				_ = query.Get("offset")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)

				// Return appropriate number of items
				items := make([]repositories.Repository, tt.count)
				for i := 0; i < tt.count; i++ {
					items[i] = repositories.Repository{
						ID:              "repo-" + string(rune(i)),
						Name:            "repo",
						Namespace:       "public",
						OperatingSystem: "linux",
						Count:           0,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					}
				}

				// Marshal and write response
				_, _ = w.Write(marshalRepositories(items))
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")
			ctx := context.Background()

			opts := &repositories.ListRepositoriesOptions{
				Limit:  tt.limit,
				Offset: tt.offset,
			}
			repos, err := client.List(ctx, opts)
			if err != nil {
				t.Fatalf("List failed: %v", err)
			}
			if len(repos) != tt.count {
				t.Errorf("Expected %d repositories, got %d", tt.count, len(repos))
			}
		})
	}
}

// T049: Integration test - Error handling
// Verify proper error handling for various scenarios
func TestRepositoryErrorHandling(t *testing.T) {
	tests := []struct {
		name        string
		method      string
		path        string
		status      int
		expectError bool
		operation   func(client *Client, ctx context.Context) error
	}{
		{
			name:        "get non-existent repository",
			method:      http.MethodGet,
			path:        "/api/v1/project/proj-123/repository/non-existent",
			status:      http.StatusNotFound,
			expectError: true,
			operation: func(client *Client, ctx context.Context) error {
				_, err := client.Get(ctx, "non-existent")
				return err
			},
		},
		{
			name:        "unauthorized access",
			method:      http.MethodGet,
			path:        "/api/v1/project/proj-123/repositories",
			status:      http.StatusUnauthorized,
			expectError: true,
			operation: func(client *Client, ctx context.Context) error {
				_, err := client.List(ctx, nil)
				return err
			},
		},
		{
			name:        "forbidden access",
			method:      http.MethodDelete,
			path:        "/api/v1/project/proj-123/repository/repo-123",
			status:      http.StatusForbidden,
			expectError: true,
			operation: func(client *Client, ctx context.Context) error {
				return client.Delete(ctx, "repo-123")
			},
		},
		{
			name:        "server error",
			method:      http.MethodGet,
			path:        "/api/v1/project/proj-123/repositories",
			status:      http.StatusInternalServerError,
			expectError: true,
			operation: func(client *Client, ctx context.Context) error {
				_, err := client.List(ctx, nil)
				return err
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.status)
				if tt.status >= 400 {
					_, _ = w.Write([]byte(`{"error": "error occurred"}`))
				}
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")
			ctx := context.Background()

			err := tt.operation(client, ctx)

			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

// Helper function to marshal repositories (used in pagination test)
func marshalRepositories(repos []repositories.Repository) []byte {
	// Simple JSON marshaling for test
	if len(repos) == 0 {
		return []byte("[]")
	}
	result := "["
	for i := range repos {
		if i > 0 {
			result += ","
		}
		result += `{"id":"repo-` + string(rune(48+i)) + `","name":"repo","namespace":"public","operatingSystem":"linux","count":0}`
	}
	result += "]"
	return []byte(result)
}

// Verify that Create returns proper Repository struct with all fields
func TestCreateRepositoryResponseStructure(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id": "repo-789",
			"name": "complete-repo",
			"namespace": "public",
			"operatingSystem": "windows",
			"description": "A complete repository",
			"count": 42,
			"creator": {"id": "user-1", "name": "admin"},
			"project": {"id": "proj-123", "name": "test-project"},
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-02T12:30:45Z"
		}`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")
	ctx := context.Background()

	repo, err := client.Create(ctx, &repositories.CreateRepositoryRequest{
		Name:            "complete-repo",
		OperatingSystem: "windows",
		Description:     "A complete repository",
	})

	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	// Verify all fields are populated
	if repo.ID != "repo-789" {
		t.Errorf("ID mismatch: expected repo-789, got %s", repo.ID)
	}
	if repo.Name != "complete-repo" {
		t.Errorf("Name mismatch: expected complete-repo, got %s", repo.Name)
	}
	if repo.Namespace != "public" {
		t.Errorf("Namespace mismatch: expected public, got %s", repo.Namespace)
	}
	if repo.OperatingSystem != "windows" {
		t.Errorf("OS mismatch: expected windows, got %s", repo.OperatingSystem)
	}
	if repo.Description != "A complete repository" {
		t.Errorf("Description mismatch")
	}
	if repo.Count != 42 {
		t.Errorf("Count mismatch: expected 42, got %d", repo.Count)
	}
	if repo.Creator == nil || repo.Creator.ID != "user-1" {
		t.Errorf("Creator mismatch")
	}
	if repo.Project == nil || repo.Project.Name != "test-project" {
		t.Errorf("Project mismatch")
	}
	if repo.CreatedAt.IsZero() || repo.UpdatedAt.IsZero() {
		t.Errorf("Timestamps not properly parsed")
	}
}

// T091: Integration test - Repository List with limit parameter
func TestRepositoryListWithLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/project/proj-123/repositories" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		limit := r.URL.Query().Get("limit")
		if limit != "5" {
			t.Errorf("expected limit=5, got %s", limit)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"id":"repo-1","name":"repo1","namespace":"public","operatingSystem":"linux","count":0,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"},
			{"id":"repo-2","name":"repo2","namespace":"public","operatingSystem":"linux","count":0,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"},
			{"id":"repo-3","name":"repo3","namespace":"public","operatingSystem":"linux","count":0,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &repositories.ListRepositoriesOptions{Limit: 5}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 3 {
		t.Errorf("expected 3 repos, got %d", len(result))
	}
}

// T092: Integration test - Repository List with offset parameter
func TestRepositoryListWithOffset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/project/proj-123/repositories" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		offset := r.URL.Query().Get("offset")
		if offset != "10" {
			t.Errorf("expected offset=10, got %s", offset)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"id":"repo-11","name":"repo11","namespace":"public","operatingSystem":"linux","count":0,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"},
			{"id":"repo-12","name":"repo12","namespace":"public","operatingSystem":"linux","count":0,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &repositories.ListRepositoriesOptions{Offset: 10}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 repos, got %d", len(result))
	}
}

// T093: Integration test - Repository List with where filters
func TestRepositoryListWithWhereFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/project/proj-123/repositories" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		where := r.URL.Query()["where"]
		if len(where) != 2 {
			t.Errorf("expected 2 where filters, got %d", len(where))
		}
		// Check that filters are present (order may vary)
		hasOS := false
		hasCreator := false
		for _, w := range where {
			if w == "os=linux" {
				hasOS = true
			}
			if w == "creator=user-123" {
				hasCreator = true
			}
		}
		if !hasOS || !hasCreator {
			t.Errorf("expected filters 'os=linux' and 'creator=user-123'")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"id":"repo-1","name":"linux-repo","namespace":"public","operatingSystem":"linux","count":0,"creator":{"id":"user-123","name":"admin"},"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &repositories.ListRepositoriesOptions{
		Where: []string{"os=linux", "creator=user-123"},
	}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 repo, got %d", len(result))
	}
}

// T097: Integration test - Repository List with combined filters and pagination
func TestRepositoryListWithCombinedFiltersAndPagination(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/project/proj-123/repositories" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		limit := r.URL.Query().Get("limit")
		offset := r.URL.Query().Get("offset")
		where := r.URL.Query()["where"]

		if limit != "10" {
			t.Errorf("expected limit=10, got %s", limit)
		}
		if offset != "5" {
			t.Errorf("expected offset=5, got %s", offset)
		}
		if len(where) != 1 || where[0] != "os=windows" {
			t.Errorf("expected where=[os=windows], got %v", where)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"id":"repo-win-1","name":"windows-repo","namespace":"public","operatingSystem":"windows","count":5,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &repositories.ListRepositoriesOptions{
		Limit:  10,
		Offset: 5,
		Where:  []string{"os=windows"},
	}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 repo, got %d", len(result))
	}
	if result[0].OperatingSystem != "windows" {
		t.Errorf("expected OS=windows, got %s", result[0].OperatingSystem)
	}
}

// T104: Integration test - Repository List with namespace="public"
func TestRepositoryListWithNamespacePublic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/project/proj-123/repositories" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verify X-Namespace header is sent
		namespace := r.Header.Get("X-Namespace")
		if namespace != "public" {
			t.Errorf("expected X-Namespace: public, got %s", namespace)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"id":"pub-repo-1","name":"public-repo","namespace":"public","operatingSystem":"linux","count":3,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &repositories.ListRepositoriesOptions{Namespace: "public"}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 repo, got %d", len(result))
	}
}

// T105: Integration test - Repository List with namespace="private"
func TestRepositoryListWithNamespacePrivate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/project/proj-123/repositories" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Verify X-Namespace header is sent
		namespace := r.Header.Get("X-Namespace")
		if namespace != "private" {
			t.Errorf("expected X-Namespace: private, got %s", namespace)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`[
			{"id":"priv-repo-1","name":"private-repo","namespace":"private","operatingSystem":"linux","count":5,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"}
		]`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &repositories.ListRepositoriesOptions{Namespace: "private"}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 repo, got %d", len(result))
	}
}

// T106: Integration test - Repository Create with namespace header support
func TestRepositoryCreateWithNamespace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/api/v1/project/proj-123/repository" {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// For now, just verify Create is called (namespace header will be added in T110)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id":"repo-ns","name":"ns-repo","namespace":"public","operatingSystem":"linux","count":0,"createdAt":"2024-01-01T00:00:00Z","updatedAt":"2024-01-01T00:00:00Z"
		}`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	req := &repositories.CreateRepositoryRequest{
		Name:            "ns-repo",
		OperatingSystem: "linux",
	}
	result, err := client.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	if result.ID != "repo-ns" {
		t.Errorf("expected repo ID repo-ns, got %s", result.ID)
	}
}
