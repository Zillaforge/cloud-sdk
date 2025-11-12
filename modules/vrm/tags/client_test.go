package tags

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vrm/tags"
)

// TestNewClient tests the NewClient constructor
// T082: Verify client initialization
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

// T071: Contract test for List All Tags
// Verify GET /project/{project-id}/tags
func TestClient_List(t *testing.T) {
	tests := []struct {
		name     string
		opts     *tags.ListTagsOptions
		mockTags []*tags.Tag
		expected int
	}{
		{
			name: "list all tags",
			opts: nil,
			mockTags: []*tags.Tag{
				{
					ID:           "tag-1",
					Name:         "v1",
					RepositoryID: "repo-123",
					Type:         "common",
					Size:         1024,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
			},
			expected: 1,
		},
		{
			name: "list with pagination",
			opts: &tags.ListTagsOptions{
				Limit:  5,
				Offset: 0,
			},
			mockTags: []*tags.Tag{},
			expected: 0,
		},
		{
			name:     "empty list",
			opts:     nil,
			mockTags: []*tags.Tag{},
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				if r.URL.Path != "/api/v1/project/proj-123/tags" {
					t.Errorf("expected path /api/v1/project/proj-123/tags, got %s", r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(marshalTags(tt.mockTags))
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			result, err := client.List(ctx, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != tt.expected {
				t.Errorf("expected %d tags, got %d", tt.expected, len(result))
			}
		})
	}
}

// T072: Contract test for List Tags By Repository
// Verify GET /project/{project-id}/repository/{repository-id}/tags
func TestClient_ListByRepository(t *testing.T) {
	tests := []struct {
		name         string
		repositoryID string
		opts         *tags.ListTagsOptions
		mockTags     []*tags.Tag
		expected     int
	}{
		{
			name:         "list tags in repository",
			repositoryID: "repo-456",
			opts:         nil,
			mockTags: []*tags.Tag{
				{
					ID:           "tag-1",
					Name:         "v1",
					RepositoryID: "repo-456",
					Type:         "common",
					Size:         2048,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
				{
					ID:           "tag-2",
					Name:         "v2",
					RepositoryID: "repo-456",
					Type:         "common",
					Size:         4096,
					CreatedAt:    time.Now(),
					UpdatedAt:    time.Now(),
				},
			},
			expected: 2,
		},
		{
			name:         "empty repository",
			repositoryID: "repo-456",
			opts:         nil,
			mockTags:     []*tags.Tag{},
			expected:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodGet {
					t.Errorf("expected GET, got %s", r.Method)
				}
				expectedPath := "/api/v1/project/proj-123/repository/" + tt.repositoryID + "/tags"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write(marshalTags(tt.mockTags))
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

			ctx := context.Background()
			result, err := client.ListByRepository(ctx, tt.repositoryID, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != tt.expected {
				t.Errorf("expected %d tags, got %d", tt.expected, len(result))
			}
		})
	}
}

// T073: Contract test for Create Tag
// Verify POST request body with DiskFormat/ContainerFormat
func TestClient_Create(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		expectedPath := "/api/v1/project/proj-123/repository/repo-456/tag"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify request body
		var req tags.CreateTagRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("failed to decode request: %v", err)
		}
		if req.Name != "v1" || req.DiskFormat != "qcow2" {
			t.Errorf("unexpected request body")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_, _ = w.Write([]byte(`{
			"id": "tag-123",
			"name": "v1",
			"repositoryID": "repo-456",
			"type": "common",
			"size": 0,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z"
		}`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	ctx := context.Background()
	req := &tags.CreateTagRequest{
		Name:            "v1",
		Type:            "common",
		DiskFormat:      "qcow2",
		ContainerFormat: "bare",
	}

	result, err := client.Create(ctx, "repo-456", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "tag-123" {
		t.Errorf("expected ID tag-123, got %s", result.ID)
	}
}

// T074: Contract test for Get Tag
// Verify nested Repository in response
func TestClient_Get(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		expectedPath := "/api/v1/project/proj-123/tag/tag-123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "tag-123",
			"name": "v1",
			"repositoryID": "repo-456",
			"type": "common",
			"size": 1024,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-01T00:00:00Z",
			"repository": {
				"id": "repo-456",
				"name": "ubuntu",
				"namespace": "public",
				"operatingSystem": "linux",
				"count": 1,
				"createdAt": "2024-01-01T00:00:00Z",
				"updatedAt": "2024-01-01T00:00:00Z"
			}
		}`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	ctx := context.Background()
	result, err := client.Get(ctx, "tag-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.ID != "tag-123" {
		t.Errorf("expected ID tag-123, got %s", result.ID)
	}
	if result.Repository == nil {
		t.Fatal("expected nested repository, got nil")
	}
	if result.Repository.ID != "repo-456" {
		t.Errorf("expected repository ID repo-456, got %s", result.Repository.ID)
	}
}

// T075: Contract test for Update Tag
// Verify PUT request
func TestClient_Update(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		expectedPath := "/api/v1/project/proj-123/tag/tag-123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"id": "tag-123",
			"name": "v2",
			"repositoryID": "repo-456",
			"type": "common",
			"size": 1024,
			"createdAt": "2024-01-01T00:00:00Z",
			"updatedAt": "2024-01-02T00:00:00Z"
		}`))
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	ctx := context.Background()
	req := &tags.UpdateTagRequest{
		Name: "v2",
	}

	result, err := client.Update(ctx, "tag-123", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result.Name != "v2" {
		t.Errorf("expected name v2, got %s", result.Name)
	}
}

// T076: Contract test for Delete Tag
// Verify 204 No Content response
func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}
		expectedPath := "/api/v1/project/proj-123/tag/tag-123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	ctx := context.Background()
	err := client.Delete(ctx, "tag-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Helper function to marshal tags
func marshalTags(tagList []*tags.Tag) []byte {
	if len(tagList) == 0 {
		return []byte(`{"tags":[],"total":0}`)
	}
	result := `{"tags":[`
	for i := range tagList {
		if i > 0 {
			result += ","
		}
		result += `{"id":"tag-` + string(rune(48+i)) + `","name":"v` + string(rune(48+i)) + `","repositoryID":"repo-123","type":"common","size":0}`
	}
	result += `],"total":` + string(rune(48+len(tagList))) + `}`
	return []byte(result)
}

// T109: Unit test for X-Namespace header construction in tag client
func TestTagNamespaceHeaderConstruction(t *testing.T) {
	tests := []struct {
		name              string
		opts              *tags.ListTagsOptions
		expectNamespace   string
		expectNoNamespace bool
	}{
		{
			name: "with public namespace",
			opts: &tags.ListTagsOptions{
				Namespace: "public",
			},
			expectNamespace: "public",
		},
		{
			name: "with private namespace",
			opts: &tags.ListTagsOptions{
				Namespace: "private",
			},
			expectNamespace: "private",
		},
		{
			name: "with empty namespace (should not set header)",
			opts: &tags.ListTagsOptions{
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
