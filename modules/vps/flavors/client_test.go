package flavors

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/flavors"
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

// TestClient_List tests flavor listing with various options
func TestClient_List(t *testing.T) {
	tests := []struct {
		name         string
		opts         *flavors.ListFlavorsOptions
		mockFlavors  []*flavors.Flavor
		expectedPath string
	}{
		{
			name: "list all flavors",
			opts: nil,
			mockFlavors: []*flavors.Flavor{
				{ID: "flav-1", Name: "small", VCPUs: 2, RAM: 4096, Disk: 20, Public: true},
				{ID: "flav-2", Name: "medium", VCPUs: 4, RAM: 8192, Disk: 40, Public: true},
			},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
		{
			name: "filter by name",
			opts: &flavors.ListFlavorsOptions{Name: "small"},
			mockFlavors: []*flavors.Flavor{
				{ID: "flav-1", Name: "small", VCPUs: 2, RAM: 4096, Disk: 20, Public: true},
			},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
		{
			name: "filter by public=true",
			opts: func() *flavors.ListFlavorsOptions {
				publicFlag := true
				return &flavors.ListFlavorsOptions{Public: &publicFlag}
			}(),
			mockFlavors: []*flavors.Flavor{
				{ID: "flav-1", Name: "public-flavor", VCPUs: 2, RAM: 4096, Disk: 20, Public: true},
			},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
		{
			name: "filter by tag",
			opts: &flavors.ListFlavorsOptions{Tag: "gpu"},
			mockFlavors: []*flavors.Flavor{
				{ID: "flav-gpu", Name: "gpu-large", VCPUs: 8, RAM: 16384, Disk: 100, Public: true, Tags: []string{"gpu"}},
			},
			expectedPath: "/api/v1/project/proj-123/flavors",
		},
		{
			name:         "empty list",
			opts:         nil,
			mockFlavors:  []*flavors.Flavor{},
			expectedPath: "/api/v1/project/proj-123/flavors",
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
					if tt.opts.Name != "" && query.Get("name") != tt.opts.Name {
						t.Errorf("expected name=%s, got %s", tt.opts.Name, query.Get("name"))
					}
					if tt.opts.Public != nil && query.Get("public") != "" {
						expected := "true"
						if !*tt.opts.Public {
							expected = "false"
						}
						if query.Get("public") != expected {
							t.Errorf("expected public=%s, got %s", expected, query.Get("public"))
						}
					}
					if tt.opts.Tag != "" && query.Get("tag") != tt.opts.Tag {
						t.Errorf("expected tag=%s, got %s", tt.opts.Tag, query.Get("tag"))
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := &flavors.FlavorListResponse{Items: tt.mockFlavors}
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			resp, err := client.List(ctx, tt.opts)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if resp == nil {
				t.Fatal("expected response, got nil")
			}
			if len(resp.Items) != len(tt.mockFlavors) {
				t.Errorf("expected %d flavors, got %d", len(tt.mockFlavors), len(resp.Items))
			}
		})
	}
}

// TestClient_List_Errors tests error handling for flavor listing
func TestClient_List_Errors(t *testing.T) {
	tests := []struct {
		name         string
		mockStatus   int
		mockResponse map[string]interface{}
	}{
		{
			name:       "unauthorized",
			mockStatus: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Invalid token",
			},
		},
		{
			name:       "internal server error",
			mockStatus: http.StatusInternalServerError,
			mockResponse: map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{}, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			_, err := client.List(ctx, nil)

			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

// TestClient_Get tests retrieving a specific flavor
func TestClient_Get(t *testing.T) {
	mockFlavor := &flavors.Flavor{
		ID:          "flav-123",
		Name:        "large",
		Description: "Large compute instance",
		VCPUs:       8,
		RAM:         16384,
		Disk:        80,
		Public:      true,
		Tags:        []string{"compute", "balanced"},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		expectedPath := "/api/v1/project/proj-123/flavors/flav-123"
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockFlavor)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	flavor, err := client.Get(ctx, "flav-123")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if flavor == nil {
		t.Fatal("expected flavor, got nil")
	}
	if flavor.ID != mockFlavor.ID {
		t.Errorf("expected flavor ID %s, got %s", mockFlavor.ID, flavor.ID)
	}
	if flavor.Name != mockFlavor.Name {
		t.Errorf("expected flavor name %s, got %s", mockFlavor.Name, flavor.Name)
	}
	if flavor.VCPUs != mockFlavor.VCPUs {
		t.Errorf("expected %d VCPUs, got %d", mockFlavor.VCPUs, flavor.VCPUs)
	}
}

// TestClient_Get_Errors tests error handling for flavor retrieval
func TestClient_Get_Errors(t *testing.T) {
	tests := []struct {
		name         string
		flavorID     string
		mockStatus   int
		mockResponse map[string]interface{}
	}{
		{
			name:       "not found",
			flavorID:   "nonexistent",
			mockStatus: http.StatusNotFound,
			mockResponse: map[string]interface{}{
				"error":   "Not Found",
				"message": "Flavor not found",
			},
		},
		{
			name:       "unauthorized",
			flavorID:   "flav-123",
			mockStatus: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Invalid token",
			},
		},
		{
			name:       "internal server error",
			flavorID:   "flav-123",
			mockStatus: http.StatusInternalServerError,
			mockResponse: map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "Something went wrong",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.mockStatus)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{}, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			_, err := client.Get(ctx, tt.flavorID)

			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}
