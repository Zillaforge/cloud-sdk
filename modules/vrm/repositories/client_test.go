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
	tagmod "github.com/Zillaforge/cloud-sdk/models/vrm/tags"
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
				response := map[string]interface{}{
					"repositories": tt.mockRepositories,
					"total":        len(tt.mockRepositories),
				}
				_ = json.NewEncoder(w).Encode(response)
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

// T111: Contract test for Upload Image
// Verify POST /api/v1/project/{project-id}/upload with multiple request variants
func TestClient_Upload(t *testing.T) {
	tests := []struct {
		name      string
		req       repositories.UploadRequester
		namespace string
	}{
		{
			name: "upload to new repository",
			req: &repositories.UploadToNewRepositoryRequest{
				Name:            "ubuntu",
				OperatingSystem: "linux",
				Version:         "v1",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image",
			},
		},
		{
			name:      "upload to existing repository with namespace",
			namespace: "public",
			req: &repositories.UploadToExistingRepositoryRequest{
				RepositoryID:    "repo-123",
				Version:         "v2",
				Type:            "common",
				DiskFormat:      "raw",
				ContainerFormat: "bare",
				Filepath:        "s3://bucket/image",
			},
		},
		{
			name: "upload to existing tag",
			req: &repositories.UploadToExistingTagRequest{
				TagID:    "tag-123",
				Filepath: "s3://bucket/image",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/upload"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}
				if tt.namespace != "" {
					if r.Header.Get("X-Namespace") != tt.namespace {
						t.Errorf("expected namespace header %s, got %s", tt.namespace, r.Header.Get("X-Namespace"))
					}
				} else if r.Header.Get("X-Namespace") != "" {
					t.Errorf("did not expect namespace header, got %s", r.Header.Get("X-Namespace"))
				}

				var body repositories.UploadImageRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(&repositories.UploadImageResponse{
					Repository: &repositories.Repository{
						ID:              "repo-999",
						Name:            "ubuntu",
						Namespace:       "public",
						OperatingSystem: "linux",
						Count:           1,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
					Tag: &repositories.Tag{
						ID:           "tag-1",
						Name:         "v1",
						RepositoryID: "repo-999",
						Type:         "common",
						Size:         10,
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
				})
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			var (
				resp *repositories.UploadImageResponse
				err  error
			)
			if tt.namespace != "" {
				resp, err = client.UploadWithNamespace(ctx, tt.req, tt.namespace)
			} else {
				resp, err = client.Upload(ctx, tt.req)
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected upload response, got nil")
			}
			if resp.Repository == nil || len(resp.Repository.Tags) == 0 {
				t.Fatalf("expected tag appended to repository, got none")
			}
			if resp.Repository.Tags[0].ID != "tag-1" {
				t.Errorf("expected tag ID tag-1, got %s", resp.Repository.Tags[0].ID)
			}
		})
	}
}

// T112: Validation test for UploadImageRequest nil handling
func TestClient_Upload_InvalidRequest(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	_, err := client.Upload(context.Background(), nil)
	if err == nil {
		t.Fatalf("expected error for nil request, got nil")
	}
}

// T047: Contract test for Snapshot creation
// Verify POST /api/v1/project/{project-id}/server/{server-id}/snapshot
func TestClient_Snapshot(t *testing.T) {
	tests := []struct {
		name            string
		namespace       string
		req             *repositories.CreateSnapshotRequest
		expectTagInRepo bool
	}{
		{
			name: "snapshot creates new repository",
			req: &repositories.CreateSnapshotRequest{
				Version:         "v1",
				Name:            "repo-from-server",
				OperatingSystem: "linux",
			},
			expectTagInRepo: true,
		},
		{
			name:      "snapshot into existing repository with namespace",
			namespace: "public",
			req: &repositories.CreateSnapshotRequest{
				Version:      "v2",
				RepositoryID: "repo-123",
			},
			expectTagInRepo: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/server/server-abc/snapshot"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.namespace != "" {
					if r.Header.Get("X-Namespace") != tt.namespace {
						t.Errorf("expected namespace %s, got %s", tt.namespace, r.Header.Get("X-Namespace"))
					}
				} else if r.Header.Get("X-Namespace") != "" {
					t.Errorf("did not expect namespace header, got %s", r.Header.Get("X-Namespace"))
				}

				var body repositories.CreateSnapshotRequest
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Fatalf("failed to decode request body: %v", err)
				}

				if tt.req.RepositoryID != "" {
					if body.RepositoryID != tt.req.RepositoryID {
						t.Errorf("expected repositoryId %s, got %s", tt.req.RepositoryID, body.RepositoryID)
					}
				} else {
					if body.Name != tt.req.Name {
						t.Errorf("expected name %s, got %s", tt.req.Name, body.Name)
					}
					if body.OperatingSystem != tt.req.OperatingSystem {
						t.Errorf("expected OS %s, got %s", tt.req.OperatingSystem, body.OperatingSystem)
					}
				}

				resp := &repositories.CreateSnapshotResponse{
					Repository: &repositories.Repository{
						ID:              "repo-from-snapshot",
						Name:            "repo-from-server",
						Namespace:       "public",
						OperatingSystem: "linux",
						Count:           1,
						CreatedAt:       time.Now(),
						UpdatedAt:       time.Now(),
					},
					Tag: &repositories.Tag{
						ID:           "tag-1",
						Name:         "v1",
						RepositoryID: "repo-from-snapshot",
						Type:         "common",
						CreatedAt:    time.Now(),
						UpdatedAt:    time.Now(),
					},
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(resp)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			var (
				resp *repositories.CreateSnapshotResponse
				err  error
			)
			if tt.namespace != "" {
				resp, err = client.SnapshotWithNamespace(ctx, "server-abc", tt.req, tt.namespace)
			} else {
				resp, err = client.Snapshot(ctx, "server-abc", tt.req)
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected snapshot response, got nil")
			}
			if resp.Repository.ID != "repo-from-snapshot" {
				t.Errorf("expected repository ID repo-from-snapshot, got %s", resp.Repository.ID)
			}
			if tt.expectTagInRepo {
				if resp.Repository == nil || len(resp.Repository.Tags) == 0 {
					t.Fatalf("expected tags in repository, got none")
				}
				if resp.Repository.Tags[0].ID != "tag-1" {
					t.Errorf("expected tag ID tag-1, got %s", resp.Repository.Tags[0].ID)
				}
			}
		})
	}
}

// T048: Validation tests for snapshot verb
func TestClient_SnapshotValidation(t *testing.T) {
	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient("https://api.example.com", "token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	_, err := client.Snapshot(context.Background(), "", &repositories.CreateSnapshotRequest{
		Version:         "v1",
		Name:            "repo",
		OperatingSystem: "linux",
	})
	if err == nil {
		t.Fatal("expected error for empty server ID, got nil")
	}

	_, err = client.Snapshot(context.Background(), "server-abc", nil)
	if err == nil {
		t.Fatal("expected error for nil request, got nil")
	}

	_, err = client.Snapshot(context.Background(), "server-abc", &repositories.CreateSnapshotRequest{
		Version:         "",
		Name:            "repo",
		OperatingSystem: "linux",
	})
	if err == nil {
		t.Fatal("expected error for invalid request, got nil")
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
				response := map[string]interface{}{
					"repositories": []interface{}{},
					"total":        0,
				}
				_ = json.NewEncoder(w).Encode(response)
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
				response := map[string]interface{}{
					"repositories": []interface{}{},
					"total":        0,
				}
				_ = json.NewEncoder(w).Encode(response)
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

// T110: Test RepositoryResource Tags method
func TestRepositoryResource_Tags(t *testing.T) {
	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)
	repo := &repositories.Repository{
		ID:   "repo-123",
		Name: "test-repo",
	}

	repoResource := &RepositoryResource{
		Repository: repo,
		tagOps: &TagsClient{
			baseClient:   baseClient,
			repositoryID: repo.ID,
			basePath:     "/api/v1/project/proj-123",
		},
	}

	tagsOps := repoResource.Tags()
	if tagsOps == nil {
		t.Fatal("expected TagOperations, got nil")
	}

	// Verify it's the correct type
	tagsClient, ok := tagsOps.(*TagsClient)
	if !ok {
		t.Fatal("expected TagsClient, got different type")
	}

	if tagsClient.repositoryID != repo.ID {
		t.Errorf("expected repositoryID %s, got %s", repo.ID, tagsClient.repositoryID)
	}
}

// T111: Test TagsClient List method
func TestTagsClient_List(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		opts         *tagmod.ListTagsOptions
		mockTags     []*tagmod.Tag
		expectedPath string
	}{
		{
			name:         "list tags for repository",
			repositoryID: "repo-123",
			opts:         nil,
			mockTags: []*tagmod.Tag{
				{
					ID:           "tag-1",
					Name:         "v1.0",
					RepositoryID: "repo-123",
					Type:         "image",
					Size:         1024,
					Status:       "active",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
					Repository: &tagmod.Repository{
						ID:   "repo-123",
						Name: "test-repo",
					},
				},
			},
			expectedPath: "/api/v1/project/proj-123/repository/repo-123/tags",
		},
		{
			name:         "list tags with limit and offset",
			repositoryID: "repo-456",
			opts: &tagmod.ListTagsOptions{
				Limit:  5,
				Offset: 10,
			},
			mockTags:     []*tagmod.Tag{},
			expectedPath: "/api/v1/project/proj-123/repository/repo-456/tags",
		},
		{
			name:         "list tags with namespace",
			repositoryID: "repo-789",
			opts: &tagmod.ListTagsOptions{
				Namespace: "private",
			},
			mockTags:     []*tagmod.Tag{},
			expectedPath: "/api/v1/project/proj-123/repository/repo-789/tags",
		},
		{
			name:         "list tags with where filters",
			repositoryID: "repo-999",
			opts: &tagmod.ListTagsOptions{
				Where: []string{"status=active", "type=image"},
			},
			mockTags: []*tagmod.Tag{
				{
					ID:           "tag-active",
					Name:         "active-tag",
					RepositoryID: "repo-999",
					Type:         "image",
					Status:       "active",
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
			},
			expectedPath: "/api/v1/project/proj-123/repository/repo-999/tags",
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

				// Verify query parameters
				if tt.opts != nil {
					query := r.URL.Query()
					if tt.opts.Limit > 0 {
						if query.Get("limit") != "5" {
							t.Errorf("expected limit=5, got %s", query.Get("limit"))
						}
					}
					if tt.opts.Offset > 0 {
						if query.Get("offset") != "10" {
							t.Errorf("expected offset=10, got %s", query.Get("offset"))
						}
					}
					if len(tt.opts.Where) > 0 {
						whereParams := query["where"]
						if len(whereParams) != len(tt.opts.Where) {
							t.Errorf("expected %d where params, got %d", len(tt.opts.Where), len(whereParams))
						}
					}
				}

				// Verify namespace header
				if tt.opts != nil && tt.opts.Namespace != "" {
					if r.Header.Get("X-Namespace") != tt.opts.Namespace {
						t.Errorf("expected X-Namespace %s, got %s", tt.opts.Namespace, r.Header.Get("X-Namespace"))
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := tagmod.ListTagsResponse{Tags: tt.mockTags}
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			tagsClient := &TagsClient{
				baseClient:   baseClient,
				repositoryID: tt.repositoryID,
				basePath:     "/api/v1/project/proj-123",
			}

			ctx := context.Background()
			tags, err := tagsClient.List(ctx, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(tags) != len(tt.mockTags) {
				t.Errorf("expected %d tags, got %d", len(tt.mockTags), len(tags))
			}
		})
	}
}

// T112: Test TagsClient Create method
func TestTagsClient_Create(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		req          *tagmod.CreateTagRequest
		mockTag      *tagmod.Tag
		expectError  bool
	}{
		{
			name:         "create tag successfully",
			repositoryID: "repo-123",
			req: &tagmod.CreateTagRequest{
				Name:            "v1.0",
				Type:            "image",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
			},
			mockTag: &tagmod.Tag{
				ID:           "tag-123",
				Name:         "v1.0",
				RepositoryID: "repo-123",
				Type:         "image",
				Size:         2048,
				Status:       "active",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectError: false,
		},
		{
			name:         "create tag with nil request",
			repositoryID: "repo-123",
			req:          nil,
			mockTag:      nil,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/repository/" + tt.repositoryID + "/tag"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify request body if not nil
				if tt.req != nil {
					var reqBody tagmod.CreateTagRequest
					if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
						t.Fatalf("failed to decode request body: %v", err)
					}
					if reqBody.Name != tt.req.Name {
						t.Errorf("expected name %s, got %s", tt.req.Name, reqBody.Name)
					}
				}

				w.Header().Set("Content-Type", "application/json")
				if tt.expectError {
					w.WriteHeader(http.StatusBadRequest)
				} else {
					w.WriteHeader(http.StatusCreated)
					_ = json.NewEncoder(w).Encode(tt.mockTag)
				}
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			tagsClient := &TagsClient{
				baseClient:   baseClient,
				repositoryID: tt.repositoryID,
				basePath:     "/api/v1/project/proj-123",
			}

			ctx := context.Background()
			tag, err := tagsClient.Create(ctx, tt.req)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tag == nil {
					t.Fatal("expected tag, got nil")
				}
				if tag.Name != tt.mockTag.Name {
					t.Errorf("expected name %s, got %s", tt.mockTag.Name, tag.Name)
				}
			}
		})
	}
}

// T113: Test TagsClient CreateWithNamespace method
func TestTagsClient_CreateWithNamespace(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		req          *tagmod.CreateTagRequest
		namespace    string
		mockTag      *tagmod.Tag
		expectError  bool
	}{
		{
			name:         "create tag with namespace",
			repositoryID: "repo-123",
			req: &tagmod.CreateTagRequest{
				Name:            "v2.0",
				Type:            "image",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
			},
			namespace: "private",
			mockTag: &tagmod.Tag{
				ID:           "tag-456",
				Name:         "v2.0",
				RepositoryID: "repo-123",
				Type:         "image",
				Size:         4096,
				Status:       "active",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectError: false,
		},
		{
			name:         "create tag with empty namespace",
			repositoryID: "repo-456",
			req: &tagmod.CreateTagRequest{
				Name:            "v3.0",
				Type:            "common",
				DiskFormat:      "qcow2",
				ContainerFormat: "bare",
			},
			namespace: "",
			mockTag: &tagmod.Tag{
				ID:           "tag-789",
				Name:         "v3.0",
				RepositoryID: "repo-456",
				Type:         "common",
				Size:         1024,
				Status:       "active",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/repository/" + tt.repositoryID + "/tag"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				// Verify namespace header
				if tt.namespace != "" {
					if r.Header.Get("X-Namespace") != tt.namespace {
						t.Errorf("expected X-Namespace %s, got %s", tt.namespace, r.Header.Get("X-Namespace"))
					}
				} else {
					if r.Header.Get("X-Namespace") != "" {
						t.Errorf("expected no X-Namespace header, got %s", r.Header.Get("X-Namespace"))
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(tt.mockTag)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			tagsClient := &TagsClient{
				baseClient:   baseClient,
				repositoryID: tt.repositoryID,
				basePath:     "/api/v1/project/proj-123",
			}

			ctx := context.Background()
			tag, err := tagsClient.CreateWithNamespace(ctx, tt.req, tt.namespace)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if tag == nil {
					t.Fatal("expected tag, got nil")
				}
				if tag.Name != tt.mockTag.Name {
					t.Errorf("expected name %s, got %s", tt.mockTag.Name, tag.Name)
				}
			}
		})
	}
}

// T114: Test CreateWithNamespace method for repositories
func TestClient_CreateWithNamespace(t *testing.T) {
	tests := []struct {
		name      string
		req       *repositories.CreateRepositoryRequest
		namespace string
		mockRepo  *repositories.Repository
	}{
		{
			name: "create repository with namespace",
			req: &repositories.CreateRepositoryRequest{
				Name:            "test-repo-ns",
				OperatingSystem: "linux",
			},
			namespace: "private",
			mockRepo: &repositories.Repository{
				ID:              "repo-ns-123",
				Name:            "test-repo-ns",
				Namespace:       "private",
				OperatingSystem: "linux",
				Count:           0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
		{
			name: "create repository with empty namespace",
			req: &repositories.CreateRepositoryRequest{
				Name:            "test-repo-no-ns",
				OperatingSystem: "windows",
			},
			namespace: "",
			mockRepo: &repositories.Repository{
				ID:              "repo-no-ns-456",
				Name:            "test-repo-no-ns",
				Namespace:       "public",
				OperatingSystem: "windows",
				Count:           0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
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

				// Verify namespace header
				if tt.namespace != "" {
					if r.Header.Get("X-Namespace") != tt.namespace {
						t.Errorf("expected X-Namespace %s, got %s", tt.namespace, r.Header.Get("X-Namespace"))
					}
				} else {
					if r.Header.Get("X-Namespace") != "" {
						t.Errorf("expected no X-Namespace header, got %s", r.Header.Get("X-Namespace"))
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
			repo, err := client.CreateWithNamespace(ctx, tt.req, tt.namespace)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo == nil {
				t.Fatal("expected repository, got nil")
			}
			if repo.Name != tt.mockRepo.Name {
				t.Errorf("expected name %s, got %s", tt.mockRepo.Name, repo.Name)
			}
		})
	}
}

// T115: Test GetWithNamespace method
func TestClient_GetWithNamespace(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		namespace    string
		mockRepo     *repositories.Repository
	}{
		{
			name:         "get repository with namespace",
			repositoryID: "repo-123",
			namespace:    "private",
			mockRepo: &repositories.Repository{
				ID:              "repo-123",
				Name:            "private-repo",
				Namespace:       "private",
				OperatingSystem: "linux",
				Count:           3,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		},
		{
			name:         "get repository with empty namespace",
			repositoryID: "repo-456",
			namespace:    "",
			mockRepo: &repositories.Repository{
				ID:              "repo-456",
				Name:            "public-repo",
				Namespace:       "public",
				OperatingSystem: "windows",
				Count:           0,
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
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

				// Verify namespace header
				if tt.namespace != "" {
					if r.Header.Get("X-Namespace") != tt.namespace {
						t.Errorf("expected X-Namespace %s, got %s", tt.namespace, r.Header.Get("X-Namespace"))
					}
				} else {
					if r.Header.Get("X-Namespace") != "" {
						t.Errorf("expected no X-Namespace header, got %s", r.Header.Get("X-Namespace"))
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
			repo, err := client.GetWithNamespace(ctx, tt.repositoryID, tt.namespace)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo == nil {
				t.Fatal("expected repository, got nil")
			}
			if repo.ID != tt.mockRepo.ID {
				t.Errorf("expected ID %s, got %s", tt.mockRepo.ID, repo.ID)
			}
		})
	}
}

// T116: Test UpdateWithNamespace method
func TestClient_UpdateWithNamespace(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		req          *repositories.UpdateRepositoryRequest
		namespace    string
		mockRepo     *repositories.Repository
	}{
		{
			name:         "update repository with namespace",
			repositoryID: "repo-123",
			req: &repositories.UpdateRepositoryRequest{
				Description: "Updated with namespace",
			},
			namespace: "private",
			mockRepo: &repositories.Repository{
				ID:              "repo-123",
				Name:            "test-repo",
				Namespace:       "private",
				OperatingSystem: "linux",
				Description:     "Updated with namespace",
				Count:           5,
				CreatedAt:       time.Now().Add(-1 * time.Hour),
				UpdatedAt:       time.Now(),
			},
		},
		{
			name:         "update repository with empty namespace",
			repositoryID: "repo-456",
			req: &repositories.UpdateRepositoryRequest{
				Description: "Updated without namespace",
			},
			namespace: "",
			mockRepo: &repositories.Repository{
				ID:              "repo-456",
				Name:            "test-repo-2",
				Namespace:       "public",
				OperatingSystem: "windows",
				Description:     "Updated without namespace",
				Count:           2,
				CreatedAt:       time.Now().Add(-2 * time.Hour),
				UpdatedAt:       time.Now(),
			},
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

				// Verify namespace header
				if tt.namespace != "" {
					if r.Header.Get("X-Namespace") != tt.namespace {
						t.Errorf("expected X-Namespace %s, got %s", tt.namespace, r.Header.Get("X-Namespace"))
					}
				} else {
					if r.Header.Get("X-Namespace") != "" {
						t.Errorf("expected no X-Namespace header, got %s", r.Header.Get("X-Namespace"))
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
			repo, err := client.UpdateWithNamespace(ctx, tt.repositoryID, tt.req, tt.namespace)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if repo == nil {
				t.Fatal("expected repository, got nil")
			}
			if repo.ID != tt.mockRepo.ID {
				t.Errorf("expected ID %s, got %s", tt.mockRepo.ID, repo.ID)
			}
		})
	}
}

// T117: Test DeleteWithNamespace method
func TestClient_DeleteWithNamespace(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		namespace    string
		mockStatus   int
		expectError  bool
	}{
		{
			name:         "delete repository with namespace success",
			repositoryID: "repo-123",
			namespace:    "private",
			mockStatus:   http.StatusNoContent,
			expectError:  false,
		},
		{
			name:         "delete repository with empty namespace success",
			repositoryID: "repo-456",
			namespace:    "",
			mockStatus:   http.StatusNoContent,
			expectError:  false,
		},
		{
			name:         "delete repository with namespace not found",
			repositoryID: "repo-nonexistent",
			namespace:    "private",
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

				// Verify namespace header
				if tt.namespace != "" {
					if r.Header.Get("X-Namespace") != tt.namespace {
						t.Errorf("expected X-Namespace %s, got %s", tt.namespace, r.Header.Get("X-Namespace"))
					}
				} else {
					if r.Header.Get("X-Namespace") != "" {
						t.Errorf("expected no X-Namespace header, got %s", r.Header.Get("X-Namespace"))
					}
				}

				w.WriteHeader(tt.mockStatus)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			err := client.DeleteWithNamespace(ctx, tt.repositoryID, tt.namespace)

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
