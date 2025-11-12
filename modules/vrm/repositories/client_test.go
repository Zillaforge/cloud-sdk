package repositories

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vrm/common"
	"github.com/Zillaforge/cloud-sdk/models/vrm/repositories"
)

// TestNewClient tests the NewClient constructor
// T050: Verify client initialization
func TestNewClient(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	projectID := "proj-123"
	basePath := "/api/v1/project/proj-123"

	client := NewClient(baseClient, projectID, basePath)

	if client == nil {
		t.Fatal("expected client, got nil")
	}
	if client.projectID != projectID {
		t.Errorf("expected projectID %s, got %s", projectID, client.projectID)
	}
	if client.basePath != basePath {
		t.Errorf("expected basePath %s, got %s", basePath, client.basePath)
	}
}

// T042: Contract test for List Repositories
// Verify GET /api/v1/project/{project-id}/repositories
func TestClient_List(t *testing.T) {
	tests := []struct {
		name             string
		opts             *repositories.ListRepositoriesOptions
		mockRepositories []*repositories.Repository
		expectedPath     string
		expectedQuery    string
	}{
		{
			name: "list all repositories",
			opts: nil,
			mockRepositories: []*repositories.Repository{
				{
					ID:              "repo-1",
					Name:            "ubuntu",
					Namespace:       "public",
					OperatingSystem: "linux",
					Count:           5,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			expectedPath: "/api/v1/project/proj-123/repositories",
		},
		{
			name: "list with limit and offset",
			opts: &repositories.ListRepositoriesOptions{
				Limit:  10,
				Offset: 5,
			},
			mockRepositories: []*repositories.Repository{},
			expectedPath:     "/api/v1/project/proj-123/repositories",
		},
		{
			name: "list with namespace filter",
			opts: &repositories.ListRepositoriesOptions{
				Namespace: "public",
			},
			mockRepositories: []*repositories.Repository{
				{
					ID:              "repo-public",
					Name:            "public-repo",
					Namespace:       "public",
					OperatingSystem: "linux",
					Count:           0,
					CreatedAt:       time.Now(),
					UpdatedAt:       time.Now(),
				},
			},
			expectedPath: "/api/v1/project/proj-123/repositories",
		},
		{
			name:             "empty list",
			opts:             nil,
			mockRepositories: []*repositories.Repository{},
			expectedPath:     "/api/v1/project/proj-123/repositories",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				if r.URL.Path != tt.expectedPath {
					t.Errorf("expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				// Verify query parameters if opts provided
				if tt.opts != nil {
					query := r.URL.Query()
					if tt.opts.Limit > 0 {
						limitStr := query.Get("limit")
						if limitStr == "" {
							t.Errorf("expected limit in query parameters")
						}
					}
					if tt.opts.Offset > 0 {
						offsetStr := query.Get("offset")
						if offsetStr == "" {
							t.Errorf("expected offset in query parameters")
						}
					}
				}

				// Verify namespace header if provided
				if tt.opts != nil && tt.opts.Namespace != "" {
					nsHeader := r.Header.Get("X-Namespace")
					if nsHeader != tt.opts.Namespace {
						t.Errorf("expected namespace header %s, got %s", tt.opts.Namespace, nsHeader)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(tt.mockRepositories)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			repos, err := client.List(ctx, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repos == nil && tt.mockRepositories != nil {
				t.Fatal("expected repositories, got nil")
			}
			if len(repos) != len(tt.mockRepositories) {
				t.Errorf("expected %d repositories, got %d", len(tt.mockRepositories), len(repos))
			}
		})
	}
}

// T043: Contract test for Create Repository
// Verify POST /api/v1/project/{project-id}/repository with request body
func TestClient_Create(t *testing.T) {
	tests := []struct {
		name        string
		req         *repositories.CreateRepositoryRequest
		mockRepo    *repositories.Repository
		expectError bool
	}{
		{
			name: "create repository with minimal fields",
			req: &repositories.CreateRepositoryRequest{
				Name:            "test-repo",
				OperatingSystem: "linux",
			},
			mockRepo: &repositories.Repository{
				ID:              "repo-123",
				Name:            "test-repo",
				Namespace:       "public",
				OperatingSystem: "linux",
				Count:           0,
				Creator: &common.IDName{
					ID:   "user-1",
					Name: "admin",
				},
				Project: &common.IDName{
					ID:   "proj-123",
					Name: "test-project",
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
		},
		{
			name: "create repository with description",
			req: &repositories.CreateRepositoryRequest{
				Name:            "ubuntu-repo",
				OperatingSystem: "linux",
				Description:     "Ubuntu repository",
			},
			mockRepo: &repositories.Repository{
				ID:              "repo-456",
				Name:            "ubuntu-repo",
				Namespace:       "public",
				OperatingSystem: "linux",
				Description:     "Ubuntu repository",
				Count:           0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			expectError: false,
		},
		{
			name:        "create with nil request",
			req:         nil,
			mockRepo:    nil,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/repository"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify request body if not nil request
				if tt.req != nil {
					var reqBody repositories.CreateRepositoryRequest
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Fatalf("failed to decode request body: %v", err)
					}
					if reqBody.Name != tt.req.Name {
						t.Errorf("expected name %s, got %s", tt.req.Name, reqBody.Name)
					}
					if reqBody.OperatingSystem != tt.req.OperatingSystem {
						t.Errorf("expected OS %s, got %s", tt.req.OperatingSystem, reqBody.OperatingSystem)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(tt.mockRepo)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			repo, err := client.Create(ctx, tt.req)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if repo == nil {
					t.Fatal("expected repository, got nil")
				}
				if repo.Name != tt.mockRepo.Name {
					t.Errorf("expected name %s, got %s", tt.mockRepo.Name, repo.Name)
				}
			}
		})
	}
}

// T044: Contract test for Get Repository
// Verify GET /api/v1/project/{project-id}/repository/{repository-id}
func TestClient_Get(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		mockRepo     *repositories.Repository
		expectError  bool
	}{
		{
			name:         "get existing repository",
			repositoryID: "repo-123",
			mockRepo: &repositories.Repository{
				ID:              "repo-123",
				Name:            "ubuntu",
				Namespace:       "public",
				OperatingSystem: "linux",
				Count:           5,
				Creator: &common.IDName{
					ID:   "user-1",
					Name: "admin",
				},
				Project: &common.IDName{
					ID:   "proj-123",
					Name: "test-project",
				},
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			},
			expectError: false,
		},
		{
			name:         "get repository with special chars in ID",
			repositoryID: "repo-with-special-123",
			mockRepo: &repositories.Repository{
				ID:              "repo-with-special-123",
				Name:            "special-repo",
				Namespace:       "private",
				OperatingSystem: "windows",
				Count:           0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET request, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/repository/" + tt.repositoryID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(tt.mockRepo)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			repo, err := client.Get(ctx, tt.repositoryID)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if repo == nil {
					t.Fatal("expected repository, got nil")
				}
				if repo.ID != tt.mockRepo.ID {
					t.Errorf("expected ID %s, got %s", tt.mockRepo.ID, repo.ID)
				}
			}
		})
	}
}

// T045: Contract test for Update Repository
// Verify PUT /api/v1/project/{project-id}/repository/{repository-id}
func TestClient_Update(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		req          *repositories.UpdateRepositoryRequest
		mockRepo     *repositories.Repository
		expectError  bool
	}{
		{
			name:         "update repository description",
			repositoryID: "repo-123",
			req: &repositories.UpdateRepositoryRequest{
				Description: "Updated description",
			},
			mockRepo: &repositories.Repository{
				ID:              "repo-123",
				Name:            "ubuntu",
				Namespace:       "public",
				OperatingSystem: "linux",
				Description:     "Updated description",
				Count:           5,
				CreatedAt:       time.Now().Add(-1 * time.Hour),
				UpdatedAt:       time.Now(),
			},
			expectError: false,
		},
		{
			name:         "update with nil request",
			repositoryID: "repo-123",
			req:          nil,
			mockRepo:     nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("expected PUT request, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/repository/" + tt.repositoryID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify request body if not nil request
				if tt.req != nil {
					var reqBody repositories.UpdateRepositoryRequest
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Fatalf("failed to decode request body: %v", err)
					}
					if tt.req.Description != "" && reqBody.Description != tt.req.Description {
						t.Errorf("expected description %s, got %s", tt.req.Description, reqBody.Description)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(tt.mockRepo)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			repo, err := client.Update(ctx, tt.repositoryID, tt.req)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if repo == nil {
					t.Fatal("expected repository, got nil")
				}
				if repo.ID != tt.mockRepo.ID {
					t.Errorf("expected ID %s, got %s", tt.mockRepo.ID, repo.ID)
				}
			}
		})
	}
}

// T046: Contract test for Delete Repository
// Verify DELETE /api/v1/project/{project-id}/repository/{repository-id} returns 204
func TestClient_Delete(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		mockStatus   int
		expectError  bool
	}{
		{
			name:         "delete repository success",
			repositoryID: "repo-123",
			mockStatus:   http.StatusNoContent,
			expectError:  false,
		},
		{
			name:         "delete non-existent repository",
			repositoryID: "repo-nonexistent",
			mockStatus:   http.StatusNotFound,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/repository/" + tt.repositoryID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatus)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			err := client.Delete(ctx, tt.repositoryID)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			}
		})
	}
}

// T098: Unit test for query string construction with multiple where filters
func TestQueryStringConstructionWithMultipleWhereFilters(t *testing.T) {
	tests := []struct {
		name          string
		opts          *repositories.ListRepositoriesOptions
		shouldHave    []string // substrings that should be in query
		shouldNotHave []string
	}{
		{
			name: "single where filter",
			opts: &repositories.ListRepositoriesOptions{
				Where: []string{"os=linux"},
			},
			shouldHave:    []string{"where=os%3Dlinux"},
			shouldNotHave: nil,
		},
		{
			name: "multiple where filters",
			opts: &repositories.ListRepositoriesOptions{
				Where: []string{"os=linux", "creator=user-1"},
			},
			shouldHave:    []string{"where=os%3Dlinux", "where=creator%3Duser-1"},
			shouldNotHave: nil,
		},
		{
			name: "limit, offset, and multiple where filters",
			opts: &repositories.ListRepositoriesOptions{
				Limit:  10,
				Offset: 5,
				Where:  []string{"os=windows", "namespace=public"},
			},
			shouldHave:    []string{"limit=10", "offset=5", "where=os%3Dwindows", "where=namespace%3Dpublic"},
			shouldNotHave: nil,
		},
		{
			name: "namespace header only",
			opts: &repositories.ListRepositoriesOptions{
				Namespace: "public",
			},
			shouldHave:    nil,
			shouldNotHave: []string{"where=", "limit=", "offset="},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				query := r.URL.RawQuery

				// Check URL-encoded query parameters
				for _, expected := range tt.shouldHave {
					if !findSubstring(query, expected) {
						t.Errorf("expected query to contain %q, got %q", expected, query)
					}
				}

				for _, unexpected := range tt.shouldNotHave {
					if findSubstring(query, unexpected) {
						t.Errorf("expected query NOT to contain %q, but it did: %q", unexpected, query)
					}
				}

				// Check X-Namespace header
				if tt.opts != nil && tt.opts.Namespace != "" {
					if r.Header.Get("X-Namespace") != tt.opts.Namespace {
						t.Errorf("expected X-Namespace header %q, got %q", tt.opts.Namespace, r.Header.Get("X-Namespace"))
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("[]"))
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			_, _ = client.List(context.Background(), tt.opts)
		})
	}
}

// T109: Unit test for X-Namespace header construction in repository client
func TestNamespaceHeaderConstruction(t *testing.T) {
	tests := []struct {
		name              string
		opts              *repositories.ListRepositoriesOptions
		expectNamespace   string
		expectNoNamespace bool
	}{
		{
			name: "with public namespace",
			opts: &repositories.ListRepositoriesOptions{
				Namespace: "public",
			},
			expectNamespace: "public",
		},
		{
			name: "with private namespace",
			opts: &repositories.ListRepositoriesOptions{
				Namespace: "private",
			},
			expectNamespace: "private",
		},
		{
			name: "with empty namespace (should not set header)",
			opts: &repositories.ListRepositoriesOptions{
				Namespace: "",
			},
			expectNoNamespace: true,
		},
		{
			name:              "with nil options (should not set header)",
			opts:              nil,
			expectNoNamespace: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.expectNoNamespace {
					if r.Header.Get("X-Namespace") != "" {
						t.Errorf("expected no X-Namespace header, got %q", r.Header.Get("X-Namespace"))
					}
				} else {
					if r.Header.Get("X-Namespace") != tt.expectNamespace {
						t.Errorf("expected X-Namespace %q, got %q", tt.expectNamespace, r.Header.Get("X-Namespace"))
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte("[]"))
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			_, _ = client.List(context.Background(), tt.opts)
		})
	}
}

// Helper function
func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
