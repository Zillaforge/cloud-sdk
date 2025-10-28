package servers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/servers"
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
	if client.baseClient == nil {
		t.Error("expected baseClient to be initialized")
	}
}

func TestClient_List(t *testing.T) {
	tests := []struct {
		name           string
		opts           *servers.ListServersOptions
		mockResponse   *servers.ServerListResponse
		mockStatusCode int
		wantErr        bool
		checkPath      string
	}{
		{
			name: "list all servers",
			opts: nil,
			mockResponse: &servers.ServerListResponse{
				Items: []*servers.Server{
					{ID: "svr-1", Name: "server-1", Status: "ACTIVE"},
					{ID: "svr-2", Name: "server-2", Status: "STOPPED"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/servers",
		},
		{
			name: "filter by status",
			opts: &servers.ListServersOptions{Status: "ACTIVE"},
			mockResponse: &servers.ServerListResponse{
				Items: []*servers.Server{
					{ID: "svr-1", Name: "server-1", Status: "ACTIVE"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/servers?status=ACTIVE",
		},
		{
			name:           "server error",
			opts:           nil,
			mockResponse:   nil,
			mockStatusCode: http.StatusInternalServerError,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.checkPath != "" {
					actualPath := r.URL.Path
					if r.URL.RawQuery != "" {
						actualPath += "?" + r.URL.RawQuery
					}
					if actualPath != tt.checkPath {
						t.Errorf("expected path %s, got %s", tt.checkPath, actualPath)
					}
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
			client := NewClient(baseClient, "proj-123")

			result, err := client.List(context.Background(), tt.opts)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result == nil {
				t.Fatal("expected result, got nil")
			}

			if len(result.Items) != len(tt.mockResponse.Items) {
				t.Errorf("expected %d servers, got %d", len(tt.mockResponse.Items), len(result.Items))
			}
		})
	}
}

func TestClient_Create(t *testing.T) {
	mockResponse := &servers.Server{
		ID:        "svr-new",
		Name:      "new-server",
		Status:    "BUILDING",
		FlavorID:  "flv-1",
		ImageID:   "img-1",
		ProjectID: "proj-123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	req := &servers.ServerCreateRequest{
		Name:     "new-server",
		FlavorID: "flv-1",
		ImageID:  "img-1",
	}

	result, err := client.Create(context.Background(), req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.ID != mockResponse.ID {
		t.Errorf("expected ID %s, got %s", mockResponse.ID, result.ID)
	}
}

func TestClient_Get(t *testing.T) {
	mockResponse := &servers.Server{
		ID:        "svr-1",
		Name:      "test-server",
		Status:    "ACTIVE",
		FlavorID:  "flv-1",
		ImageID:   "img-1",
		ProjectID: "proj-123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	result, err := client.Get(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("expected ServerResource, got nil")
	}

	if result.Server.ID != mockResponse.ID {
		t.Errorf("expected ID %s, got %s", mockResponse.ID, result.Server.ID)
	}

	// Verify sub-resource accessors
	if result.NICs() == nil {
		t.Error("expected NICs() to return non-nil client")
	}
	if result.Volumes() == nil {
		t.Error("expected Volumes() to return non-nil client")
	}
}

func TestClient_Update(t *testing.T) {
	mockResponse := &servers.Server{
		ID:          "svr-1",
		Name:        "updated-server",
		Description: "updated description",
		Status:      "ACTIVE",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PUT" {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	req := &servers.ServerUpdateRequest{
		Name:        "updated-server",
		Description: "updated description",
	}

	result, err := client.Update(context.Background(), "svr-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.Name != mockResponse.Name {
		t.Errorf("expected name %s, got %s", mockResponse.Name, result.Name)
	}
}

func TestClient_Delete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	err := client.Delete(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Action(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("expected POST, got %s", r.Method)
		}

		w.WriteHeader(http.StatusAccepted)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	req := &servers.ServerActionRequest{
		Action: "start",
	}

	err := client.Action(context.Background(), "svr-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestClient_Metrics(t *testing.T) {
	mockResponse := &servers.ServerMetricsResponse{
		Type: "cpu",
		Series: []servers.MetricPoint{
			{Timestamp: time.Now().Unix(), Value: 25.5},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	req := &servers.ServerMetricsRequest{
		Type:  "cpu",
		Start: time.Now().Add(-1 * time.Hour).Unix(),
		End:   time.Now().Unix(),
	}

	result, err := client.Metrics(context.Background(), "svr-1", req)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(result.Series) != len(mockResponse.Series) {
		t.Errorf("expected %d metrics, got %d", len(mockResponse.Series), len(result.Series))
	}
}

func TestClient_VNCURL(t *testing.T) {
	mockResponse := &servers.VNCURLResponse{
		URL: "https://console.example.com/vnc?token=abc123",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	result, err := client.VNCURL(context.Background(), "svr-1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if result.URL != mockResponse.URL {
		t.Errorf("expected URL %s, got %s", mockResponse.URL, result.URL)
	}
}
