package tags

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vrm/tags"
)

// T077: TestTagCRUDLifecycle verifies create -> get -> update -> list -> delete workflow
func TestTagCRUDLifecycle(t *testing.T) {
	const (
		projectID       = "proj-123"
		repositoryID    = "repo-456"
		tagID           = "tag-789"
		diskFormat      = "qcow2"
		containerFormat = "bare"
	)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			// Create tag endpoint
			if strings.Contains(r.URL.Path, "/repository/"+repositoryID+"/tag") {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				fmt.Fprintf(w, `{
					"id": "%s",
					"name": "test-tag",
					"repositoryID": "%s",
					"type": "image",
					"size": 1024,
					"status": "active",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-01T00:00:00Z",
					"repository": {
						"id": "%s",
						"name": "test-repo",
						"namespace": "test-ns",
						"operatingSystem": "linux",
						"count": 5,
						"createdAt": "2025-01-01T00:00:00Z",
						"updatedAt": "2025-01-01T00:00:00Z"
					}
				}`, tagID, repositoryID, repositoryID)
			}

		case "GET":
			// Get tag endpoint
			if strings.HasSuffix(r.URL.Path, "/tag/"+tagID) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{
					"id": "%s",
					"name": "test-tag-updated",
					"repositoryID": "%s",
					"type": "image",
					"size": 2048,
					"status": "active",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-02T00:00:00Z",
					"repository": {
						"id": "%s",
						"name": "test-repo",
						"namespace": "test-ns",
						"operatingSystem": "linux",
						"count": 5,
						"createdAt": "2025-01-01T00:00:00Z",
						"updatedAt": "2025-01-01T00:00:00Z"
					}
				}`, tagID, repositoryID, repositoryID)
			} else if strings.Contains(r.URL.Path, "/tags") || strings.Contains(r.URL.Path, "/repository/"+repositoryID+"/tags") {
				// List tags endpoint
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{
					"tags": [{
						"id": "%s",
						"name": "test-tag-updated",
						"repositoryID": "%s",
						"type": "image",
						"size": 2048,
						"status": "active",
						"createdAt": "2025-01-01T00:00:00Z",
						"updatedAt": "2025-01-02T00:00:00Z",
						"repository": {
							"id": "%s",
							"name": "test-repo",
							"namespace": "test-ns",
							"operatingSystem": "linux",
							"count": 5,
							"createdAt": "2025-01-01T00:00:00Z",
							"updatedAt": "2025-01-01T00:00:00Z"
						}
					}],
					"total": 1
				}`, tagID, repositoryID, repositoryID)
			}

		case "PUT":
			// Update tag endpoint
			if strings.Contains(r.URL.Path, "/tag/"+tagID) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				fmt.Fprintf(w, `{
					"id": "%s",
					"name": "test-tag-updated",
					"repositoryID": "%s",
					"type": "image",
					"size": 2048,
					"status": "active",
					"createdAt": "2025-01-01T00:00:00Z",
					"updatedAt": "2025-01-02T00:00:00Z",
					"repository": {
						"id": "%s",
						"name": "test-repo",
						"namespace": "test-ns",
						"operatingSystem": "linux",
						"count": 5,
						"createdAt": "2025-01-01T00:00:00Z",
						"updatedAt": "2025-01-01T00:00:00Z"
					}
				}`, tagID, repositoryID, repositoryID)
			}

		case "DELETE":
			// Delete tag endpoint
			if strings.Contains(r.URL.Path, "/tag/"+tagID) {
				w.WriteHeader(http.StatusNoContent)
			}
		}
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")
	ctx := context.Background()

	// Create tag
	createReq := &tags.CreateTagRequest{
		Name:            "test-tag",
		Type:            "image",
		DiskFormat:      diskFormat,
		ContainerFormat: containerFormat,
	}
	created, err := client.Create(ctx, repositoryID, createReq)
	if err != nil {
		t.Fatalf("Create: %v", err)
	}
	if created.ID != tagID {
		t.Errorf("Create: got ID %q, want %q", created.ID, tagID)
	}

	// Get tag
	got, err := client.Get(ctx, tagID)
	if err != nil {
		t.Fatalf("Get: %v", err)
	}
	if got.ID != tagID || got.Name != "test-tag-updated" {
		t.Errorf("Get: got tag ID %q name %q, want ID %q name test-tag-updated", got.ID, got.Name, tagID)
	}

	// Update tag
	updateReq := &tags.UpdateTagRequest{
		Name: "test-tag-updated",
	}
	updated, err := client.Update(ctx, tagID, updateReq)
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if updated.Name != "test-tag-updated" {
		t.Errorf("Update: got name %q, want test-tag-updated", updated.Name)
	}

	// List by repository
	listOpts := &tags.ListTagsOptions{Limit: 10}
	listed, err := client.ListByRepository(ctx, repositoryID, listOpts)
	if err != nil {
		t.Fatalf("ListByRepository: %v", err)
	}
	if len(listed) != 1 || listed[0].ID != tagID {
		t.Errorf("ListByRepository: got %d tags, want 1 with ID %q", len(listed), tagID)
	}

	// List all tags
	allTags, err := client.List(ctx, listOpts)
	if err != nil {
		t.Fatalf("List: %v", err)
	}
	if len(allTags) != 1 || allTags[0].ID != tagID {
		t.Errorf("List: got %d tags, want 1 with ID %q", len(allTags), tagID)
	}

	// Delete tag
	if err := client.Delete(ctx, tagID); err != nil {
		t.Fatalf("Delete: %v", err)
	}
}

// T078: TestTagListPagination verifies limit and offset parameters work correctly
func TestTagListPagination(t *testing.T) {
	const projectID = "proj-123"

	tests := []struct {
		name      string
		limit     int
		offset    int
		wantCount int
	}{
		{
			name:      "list_with_limit_10",
			limit:     10,
			offset:    0,
			wantCount: 3,
		},
		{
			name:      "list_with_offset_5",
			limit:     10,
			offset:    5,
			wantCount: 2,
		},
		{
			name:      "list_all_no_pagination",
			limit:     0,
			offset:    0,
			wantCount: 5,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" && strings.Contains(r.URL.Path, "/tags") {
					limit := r.URL.Query().Get("limit")
					offset := r.URL.Query().Get("offset")

					var tagsJSON []string
					// Simulate pagination logic - limit="10" and offset="" (empty/zero) means first page
					if limit == "10" && offset == "" {
						// Return first 3 tags when limit=10 and offset not set
						for i := 1; i <= 3; i++ {
							tagsJSON = append(tagsJSON, fmt.Sprintf(`{
								"id": "tag-%d",
								"name": "tag-%d",
								"repositoryID": "repo-123",
								"type": "image",
								"size": 1024,
								"status": "active",
								"createdAt": "2025-01-01T00:00:00Z",
								"updatedAt": "2025-01-01T00:00:00Z",
								"repository": {
									"id": "repo-123",
									"name": "test-repo",
									"namespace": "test-ns",
									"operatingSystem": "linux",
									"count": 5,
									"createdAt": "2025-01-01T00:00:00Z",
									"updatedAt": "2025-01-01T00:00:00Z"
								}
							}`, i, i))
						}
					} else if limit == "10" && offset == "5" {
						// Return next 2 tags
						for i := 6; i <= 7; i++ {
							tagsJSON = append(tagsJSON, fmt.Sprintf(`{
								"id": "tag-%d",
								"name": "tag-%d",
								"repositoryID": "repo-123",
								"type": "image",
								"size": 1024,
								"status": "active",
								"createdAt": "2025-01-01T00:00:00Z",
								"updatedAt": "2025-01-01T00:00:00Z",
								"repository": {
									"id": "repo-123",
									"name": "test-repo",
									"namespace": "test-ns",
									"operatingSystem": "linux",
									"count": 5,
									"createdAt": "2025-01-01T00:00:00Z",
									"updatedAt": "2025-01-01T00:00:00Z"
								}
							}`, i, i))
						}
					} else if limit == "" && offset == "" {
						// Return all 5 tags
						for i := 1; i <= 5; i++ {
							tagsJSON = append(tagsJSON, fmt.Sprintf(`{
								"id": "tag-%d",
								"name": "tag-%d",
								"repositoryID": "repo-123",
								"type": "image",
								"size": 1024,
								"status": "active",
								"createdAt": "2025-01-01T00:00:00Z",
								"updatedAt": "2025-01-01T00:00:00Z",
								"repository": {
									"id": "repo-123",
									"name": "test-repo",
									"namespace": "test-ns",
									"operatingSystem": "linux",
									"count": 5,
									"createdAt": "2025-01-01T00:00:00Z",
									"updatedAt": "2025-01-01T00:00:00Z"
								}
							}`, i, i))
						}
					}

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					fmt.Fprintf(w, `{"tags":[%s],"total":%d}`, strings.Join(tagsJSON, ","), len(tagsJSON))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, projectID, "/api/v1/project/"+projectID)
			ctx := context.Background()

			opts := &tags.ListTagsOptions{
				Limit:  tt.limit,
				Offset: tt.offset,
			}
			result, err := client.List(ctx, opts)
			if err != nil {
				t.Fatalf("List: %v", err)
			}
			if len(result) != tt.wantCount {
				t.Errorf("List: got %d tags, want %d", len(result), tt.wantCount)
			}
		})
	}
}

// T079: TestListByRepositoryFiltering verifies repository-scoped list operations
func TestListByRepositoryFiltering(t *testing.T) {
	const projectID = "proj-123"

	tests := []struct {
		name           string
		repositoryID   string
		expectedTagIDs []string
	}{
		{
			name:           "list_repo1_tags",
			repositoryID:   "repo-1",
			expectedTagIDs: []string{"tag-repo1-1", "tag-repo1-2"},
		},
		{
			name:           "list_repo2_tags",
			repositoryID:   "repo-2",
			expectedTagIDs: []string{"tag-repo2-1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method == "GET" && strings.Contains(r.URL.Path, "/repository/"+tt.repositoryID+"/tags") {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)

					var tagsJSON []string
					for _, tagID := range tt.expectedTagIDs {
						tagsJSON = append(tagsJSON, fmt.Sprintf(`{
							"id": "%s",
							"name": "%s",
							"repositoryID": "%s",
							"type": "image",
							"size": 1024,
							"status": "active",
							"createdAt": "2025-01-01T00:00:00Z",
							"updatedAt": "2025-01-01T00:00:00Z",
							"repository": {
								"id": "%s",
								"name": "test-repo",
								"namespace": "test-ns",
								"operatingSystem": "linux",
								"count": 2,
								"createdAt": "2025-01-01T00:00:00Z",
								"updatedAt": "2025-01-01T00:00:00Z"
							}
						}`, tagID, tagID, tt.repositoryID, tt.repositoryID))
					}

					fmt.Fprintf(w, `{"tags":[%s],"total":%d}`, strings.Join(tagsJSON, ","), len(tagsJSON))
					return
				}
				w.WriteHeader(http.StatusNotFound)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, projectID, "/api/v1/project/"+projectID)
			ctx := context.Background()

			result, err := client.ListByRepository(ctx, tt.repositoryID, nil)
			if err != nil {
				t.Fatalf("ListByRepository: %v", err)
			}
			if len(result) != len(tt.expectedTagIDs) {
				t.Errorf("ListByRepository: got %d tags, want %d", len(result), len(tt.expectedTagIDs))
			}
			for i, tag := range result {
				if tag.ID != tt.expectedTagIDs[i] {
					t.Errorf("ListByRepository: tag %d got ID %q, want %q", i, tag.ID, tt.expectedTagIDs[i])
				}
			}
		})
	}
}

// T080: TestTagErrorHandling verifies error scenarios (404, 401, 403, 500)
func TestTagErrorHandling(t *testing.T) {
	const (
		projectID    = "proj-123"
		repositoryID = "repo-456"
		tagID        = "tag-789"
	)

	tests := []struct {
		name        string
		method      string
		path        string
		statusCode  int
		wantErrType string
	}{
		{
			name:        "get_tag_not_found",
			method:      "GET",
			path:        "/api/v1/project/proj-123/tag/tag-789",
			statusCode:  http.StatusNotFound,
			wantErrType: "404",
		},
		{
			name:        "create_tag_unauthorized",
			method:      "POST",
			path:        "/api/v1/project/proj-123/repository/repo-456/tag",
			statusCode:  http.StatusUnauthorized,
			wantErrType: "401",
		},
		{
			name:        "update_tag_forbidden",
			method:      "PUT",
			path:        "/api/v1/project/proj-123/tag/tag-789",
			statusCode:  http.StatusForbidden,
			wantErrType: "403",
		},
		{
			name:        "delete_tag_server_error",
			method:      "DELETE",
			path:        "/api/v1/project/proj-123/tag/tag-789",
			statusCode:  http.StatusInternalServerError,
			wantErrType: "500",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.statusCode)
				if tt.statusCode >= 400 {
					w.Header().Set("Content-Type", "application/json")
					fmt.Fprintf(w, `{"error": "error message"}`)
				}
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, projectID, "/api/v1/project/"+projectID)
			ctx := context.Background()

			var err error
			switch tt.method {
			case "GET":
				_, err = client.Get(ctx, tagID)
			case "POST":
				req := &tags.CreateTagRequest{
					Name:            "test",
					Type:            "image",
					DiskFormat:      "qcow2",
					ContainerFormat: "bare",
				}
				_, err = client.Create(ctx, repositoryID, req)
			case "PUT":
				req := &tags.UpdateTagRequest{Name: "test"}
				_, err = client.Update(ctx, tagID, req)
			case "DELETE":
				err = client.Delete(ctx, tagID)
			}

			if err == nil {
				t.Errorf("%s: got no error, want error with status %s", tt.name, tt.wantErrType)
			} else if !strings.Contains(err.Error(), tt.wantErrType) {
				t.Errorf("%s: got error %v, want error containing %s", tt.name, err, tt.wantErrType)
			}
		})
	}
}

// T094: Integration test - Tag List with limit parameter
func TestTagListWithLimit(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || !strings.Contains(r.URL.Path, "/tags") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		limit := r.URL.Query().Get("limit")
		if limit != "5" {
			t.Errorf("expected limit=5, got %s", limit)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"tags": [
				{"id":"tag-1","name":"v1","repositoryID":"repo-123","type":"image","size":1024,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","repository":{"id":"repo-123","name":"repo","namespace":"ns","operatingSystem":"linux","count":1,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}},
				{"id":"tag-2","name":"v2","repositoryID":"repo-123","type":"image","size":2048,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","repository":{"id":"repo-123","name":"repo","namespace":"ns","operatingSystem":"linux","count":1,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}}
			],
			"total": 2
		}`)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &tags.ListTagsOptions{Limit: 5}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 2 {
		t.Errorf("expected 2 tags, got %d", len(result))
	}
}

// T095: Integration test - Tag List with offset parameter
func TestTagListWithOffset(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || !strings.Contains(r.URL.Path, "/tags") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		offset := r.URL.Query().Get("offset")
		if offset != "10" {
			t.Errorf("expected offset=10, got %s", offset)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"tags": [
				{"id":"tag-11","name":"v11","repositoryID":"repo-123","type":"image","size":1024,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","repository":{"id":"repo-123","name":"repo","namespace":"ns","operatingSystem":"linux","count":1,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}}
			],
			"total": 1
		}`)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &tags.ListTagsOptions{Offset: 10}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 tag, got %d", len(result))
	}
}

// T096: Integration test - Tag List with where filters
func TestTagListWithWhereFilters(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || !strings.Contains(r.URL.Path, "/tags") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		where := r.URL.Query()["where"]
		if len(where) != 2 {
			t.Errorf("expected 2 where filters, got %d", len(where))
		}
		// Check filters are present
		hasStatus := false
		hasType := false
		for _, w := range where {
			if w == "status=active" {
				hasStatus = true
			}
			if w == "type=common" {
				hasType = true
			}
		}
		if !hasStatus || !hasType {
			t.Errorf("expected filters 'status=active' and 'type=common'")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{
			"tags": [
				{"id":"tag-active","name":"active-tag","repositoryID":"repo-123","type":"common","size":1024,"status":"active","createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","repository":{"id":"repo-123","name":"repo","namespace":"ns","operatingSystem":"linux","count":1,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}}
			],
			"total": 1
		}`)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &tags.ListTagsOptions{
		Where: []string{"status=active", "type=common"},
	}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 tag, got %d", len(result))
	}
}

// T107: Integration test - Tag List with namespace header
func TestTagListWithNamespace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || !strings.Contains(r.URL.Path, "/tags") {
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
		fmt.Fprintf(w, `{
			"tags": [
				{"id":"tag-ns","name":"ns-tag","repositoryID":"repo-123","type":"image","size":1024,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","repository":{"id":"repo-123","name":"repo","namespace":"private","operatingSystem":"linux","count":1,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}}
			],
			"total": 1
		}`)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &tags.ListTagsOptions{Namespace: "private"}
	result, err := client.List(context.Background(), opts)
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 tag, got %d", len(result))
	}
}

// T108: Integration test - Tag ListByRepository with namespace header
func TestTagListByRepositoryWithNamespace(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" || !strings.Contains(r.URL.Path, "/repository/") || !strings.Contains(r.URL.Path, "/tags") {
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
		fmt.Fprintf(w, `{
			"tags": [
				{"id":"tag-repo-ns","name":"repo-ns-tag","repositoryID":"repo-456","type":"image","size":2048,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z","repository":{"id":"repo-456","name":"repo2","namespace":"public","operatingSystem":"windows","count":2,"createdAt":"2025-01-01T00:00:00Z","updatedAt":"2025-01-01T00:00:00Z"}}
			],
			"total": 1
		}`)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123", "/api/v1/project/proj-123")

	opts := &tags.ListTagsOptions{Namespace: "public"}
	result, err := client.ListByRepository(context.Background(), "repo-456", opts)
	if err != nil {
		t.Fatalf("ListByRepository failed: %v", err)
	}
	if len(result) != 1 {
		t.Errorf("expected 1 tag, got %d", len(result))
	}
}
