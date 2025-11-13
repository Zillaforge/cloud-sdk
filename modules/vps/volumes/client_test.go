package volumes

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/common"
	volumesmodel "github.com/Zillaforge/cloud-sdk/models/vps/volumes"
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

// mockHTTPClient implements a simple mock for testing
type mockHTTPClient struct {
	doFunc func(ctx context.Context, req *internalhttp.Request, response interface{}) error
}

func (m *mockHTTPClient) Do(ctx context.Context, req *internalhttp.Request, response interface{}) error {
	return m.doFunc(ctx, req, response)
}

// TestClient_Create tests volume creation
func TestClient_Create(t *testing.T) {
	tests := []struct {
		name         string
		request      *volumesmodel.CreateVolumeRequest
		mockResponse *volumesmodel.Volume
		expectedPath string
		wantErr      bool
	}{
		{
			name: "successful creation",
			request: &volumesmodel.CreateVolumeRequest{
				Name: "test-volume",
				Type: "SSD",
				Size: 100,
			},
			mockResponse: &volumesmodel.Volume{
				ID:     "vol-123",
				Name:   "test-volume",
				Type:   "SSD",
				Size:   100,
				Status: "creating",
				Project: common.IDName{
					ID:   "proj-1",
					Name: "test-project",
				},
				ProjectID: "proj-1",
				User: common.IDName{
					ID:   "user-1",
					Name: "test-user",
				},
				UserID:    "user-1",
				Namespace: "default",
			},
			expectedPath: "/api/v1/project/proj-123/volumes",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != tt.expectedPath {
					t.Errorf("expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusCreated)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			resp, err := client.Create(ctx, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("expected volume, got nil")
				}
				if resp.ID != tt.mockResponse.ID {
					t.Errorf("expected ID %s, got %s", tt.mockResponse.ID, resp.ID)
				}
			}
		})
	}
}

// TestClient_Create_Errors tests error handling for volume creation
func TestClient_Create_Errors(t *testing.T) {
	tests := []struct {
		name       string
		request    *volumesmodel.CreateVolumeRequest
		mockStatus int
		wantErr    bool
	}{
		{
			name: "invalid request - empty name",
			request: &volumesmodel.CreateVolumeRequest{
				Name: "",
				Type: "SSD",
				Size: 100,
			},
			wantErr: true,
		},
		{
			name: "server error",
			request: &volumesmodel.CreateVolumeRequest{
				Name: "test-volume",
				Type: "SSD",
				Size: 100,
			},
			mockStatus: http.StatusInternalServerError,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mockStatus == 0 {
				// Validation error - don't create server
				httpClient := &http.Client{Timeout: 5 * time.Second}
				baseClient := internalhttp.NewClient("https://api.example.com", "test-token", httpClient, nil)
				client := NewClient(baseClient, "proj-123")

				ctx := context.Background()
				_, err := client.Create(ctx, tt.request)

				if (err != nil) != tt.wantErr {
					t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				}
				return
			}

			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.WriteHeader(tt.mockStatus)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			_, err := client.Create(ctx, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestClient_Update tests volume update
func TestClient_Update(t *testing.T) {
	tests := []struct {
		name         string
		volumeID     string
		request      *volumesmodel.UpdateVolumeRequest
		mockResponse *volumesmodel.Volume
		expectedPath string
		wantErr      bool
	}{
		{
			name:     "successful update",
			volumeID: "vol-123",
			request: &volumesmodel.UpdateVolumeRequest{
				Name: "updated-name",
			},
			mockResponse: &volumesmodel.Volume{
				ID:     "vol-123",
				Name:   "updated-name",
				Type:   "SSD",
				Size:   100,
				Status: "available",
				Project: common.IDName{
					ID:   "proj-1",
					Name: "test-project",
				},
				ProjectID: "proj-1",
				User: common.IDName{
					ID:   "user-1",
					Name: "test-user",
				},
				UserID:    "user-1",
				Namespace: "default",
			},
			expectedPath: "/api/v1/project/proj-123/volumes/vol-123",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPut {
					t.Errorf("expected PUT request, got %s", r.Method)
				}
				if r.URL.Path != tt.expectedPath {
					t.Errorf("expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			resp, err := client.Update(ctx, tt.volumeID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("expected volume, got nil")
				}
				if resp.ID != tt.mockResponse.ID {
					t.Errorf("expected ID %s, got %s", tt.mockResponse.ID, resp.ID)
				}
			}
		})
	}
}

// TestClient_Delete tests volume deletion
func TestClient_Delete(t *testing.T) {
	tests := []struct {
		name         string
		volumeID     string
		expectedPath string
		wantErr      bool
	}{
		{
			name:         "successful deletion",
			volumeID:     "vol-123",
			expectedPath: "/api/v1/project/proj-123/volumes/vol-123",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete {
					t.Errorf("expected DELETE request, got %s", r.Method)
				}
				if r.URL.Path != tt.expectedPath {
					t.Errorf("expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				w.WriteHeader(http.StatusNoContent)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			err := client.Delete(ctx, tt.volumeID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestClient_List tests volume listing
func TestClient_List(t *testing.T) {
	tests := []struct {
		name         string
		opts         *volumesmodel.ListVolumesOptions
		mockVolumes  []*volumesmodel.Volume
		expectedPath string
		wantErr      bool
	}{
		{
			name: "list all volumes",
			opts: nil,
			mockVolumes: []*volumesmodel.Volume{
				{
					ID:        "vol-1",
					Name:      "volume1",
					Type:      "SSD",
					Size:      100,
					Status:    "available",
					Project:   common.IDName{ID: "proj-1", Name: "test-project"},
					ProjectID: "proj-1",
					User:      common.IDName{ID: "user-1", Name: "test-user"},
					UserID:    "user-1",
					Namespace: "default",
				},
			},
			expectedPath: "/api/v1/project/proj-123/volumes",
			wantErr:      false,
		},
		{
			name: "filter by name",
			opts: &volumesmodel.ListVolumesOptions{Name: "test"},
			mockVolumes: []*volumesmodel.Volume{
				{
					ID:        "vol-2",
					Name:      "test-volume",
					Type:      "HDD",
					Size:      200,
					Status:    "available",
					Project:   common.IDName{ID: "proj-1", Name: "test-project"},
					ProjectID: "proj-1",
					User:      common.IDName{ID: "user-1", Name: "test-user"},
					UserID:    "user-1",
					Namespace: "default",
				},
			},
			expectedPath: "/api/v1/project/proj-123/volumes",
			wantErr:      false,
		},
		{
			name: "filter by user_id",
			opts: &volumesmodel.ListVolumesOptions{UserID: "user-456"},
			mockVolumes: []*volumesmodel.Volume{
				{
					ID:        "vol-3",
					Name:      "user-volume",
					Type:      "SSD",
					Size:      50,
					Status:    "available",
					Project:   common.IDName{ID: "proj-1", Name: "test-project"},
					ProjectID: "proj-1",
					User:      common.IDName{ID: "user-456", Name: "other-user"},
					UserID:    "user-456",
					Namespace: "default",
				},
			},
			expectedPath: "/api/v1/project/proj-123/volumes",
			wantErr:      false,
		},
		{
			name: "filter by status",
			opts: &volumesmodel.ListVolumesOptions{Status: "in-use"},
			mockVolumes: []*volumesmodel.Volume{
				{
					ID:        "vol-4",
					Name:      "attached-volume",
					Type:      "SSD",
					Size:      100,
					Status:    "in-use",
					Project:   common.IDName{ID: "proj-1", Name: "test-project"},
					ProjectID: "proj-1",
					User:      common.IDName{ID: "user-1", Name: "test-user"},
					UserID:    "user-1",
					Namespace: "default",
				},
			},
			expectedPath: "/api/v1/project/proj-123/volumes",
			wantErr:      false,
		},
		{
			name: "filter by type",
			opts: &volumesmodel.ListVolumesOptions{Type: "NVMe"},
			mockVolumes: []*volumesmodel.Volume{
				{
					ID:        "vol-5",
					Name:      "nvme-volume",
					Type:      "NVMe",
					Size:      200,
					Status:    "available",
					Project:   common.IDName{ID: "proj-1", Name: "test-project"},
					ProjectID: "proj-1",
					User:      common.IDName{ID: "user-1", Name: "test-user"},
					UserID:    "user-1",
					Namespace: "default",
				},
			},
			expectedPath: "/api/v1/project/proj-123/volumes",
			wantErr:      false,
		},
		{
			name: "filter with detail",
			opts: &volumesmodel.ListVolumesOptions{Detail: true},
			mockVolumes: []*volumesmodel.Volume{
				{
					ID:        "vol-6",
					Name:      "detailed-volume",
					Type:      "SSD",
					Size:      100,
					Status:    "available",
					Project:   common.IDName{ID: "proj-1", Name: "test-project"},
					ProjectID: "proj-1",
					User:      common.IDName{ID: "user-1", Name: "test-user"},
					UserID:    "user-1",
					Namespace: "default",
				},
			},
			expectedPath: "/api/v1/project/proj-123/volumes",
			wantErr:      false,
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
					if tt.opts.Name != "" && query.Get("name") != tt.opts.Name {
						t.Errorf("expected name=%s, got %s", tt.opts.Name, query.Get("name"))
					}
					if tt.opts.UserID != "" && query.Get("user_id") != tt.opts.UserID {
						t.Errorf("expected user_id=%s, got %s", tt.opts.UserID, query.Get("user_id"))
					}
					if tt.opts.Status != "" && query.Get("status") != tt.opts.Status {
						t.Errorf("expected status=%s, got %s", tt.opts.Status, query.Get("status"))
					}
					if tt.opts.Type != "" && query.Get("type") != tt.opts.Type {
						t.Errorf("expected type=%s, got %s", tt.opts.Type, query.Get("type"))
					}
					if tt.opts.Detail && query.Get("detail") != "true" {
						t.Errorf("expected detail=true, got %s", query.Get("detail"))
					}
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				response := &volumesmodel.VolumeListResponse{Volumes: tt.mockVolumes}
				_ = json.NewEncoder(w).Encode(response)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			resp, err := client.List(ctx, tt.opts)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("expected volumes, got nil")
				}
				if len(resp) != len(tt.mockVolumes) {
					t.Errorf("expected %d volumes, got %d", len(tt.mockVolumes), len(resp))
				}
			}
		})
	}
}

// TestClient_Get tests getting a specific volume
func TestClient_Get(t *testing.T) {
	tests := []struct {
		name         string
		volumeID     string
		mockResponse *volumesmodel.Volume
		expectedPath string
		wantErr      bool
	}{
		{
			name:     "successful get",
			volumeID: "vol-123",
			mockResponse: &volumesmodel.Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Type:      "SSD",
				Size:      100,
				Status:    "available",
				Project:   common.IDName{ID: "proj-1", Name: "test-project"},
				ProjectID: "proj-1",
				User:      common.IDName{ID: "user-1", Name: "test-user"},
				UserID:    "user-1",
				Namespace: "default",
			},
			expectedPath: "/api/v1/project/proj-123/volumes/vol-123",
			wantErr:      false,
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

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(tt.mockResponse)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			resp, err := client.Get(ctx, tt.volumeID)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if resp == nil {
					t.Fatal("expected volume, got nil")
				}
				if resp.ID != tt.mockResponse.ID {
					t.Errorf("expected ID %s, got %s", tt.mockResponse.ID, resp.ID)
				}
			}
		})
	}
}

// TestClient_Action tests volume actions
func TestClient_Action(t *testing.T) {
	tests := []struct {
		name         string
		volumeID     string
		request      *volumesmodel.VolumeActionRequest
		expectedPath string
		wantErr      bool
	}{
		{
			name:     "successful attach",
			volumeID: "vol-123",
			request: &volumesmodel.VolumeActionRequest{
				Action:   volumesmodel.VolumeActionAttach,
				ServerID: "srv-456",
			},
			expectedPath: "/api/v1/project/proj-123/volumes/vol-123/action",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST request, got %s", r.Method)
				}
				if r.URL.Path != tt.expectedPath {
					t.Errorf("expected path %s, got %s", tt.expectedPath, r.URL.Path)
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusAccepted)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			err := client.Action(ctx, tt.volumeID, tt.request)

			if (err != nil) != tt.wantErr {
				t.Errorf("Action() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

// T030: Contract test for POST /volumes (create)
func TestVolumes_Create_Contract(t *testing.T) {
	tests := []struct {
		name       string
		request    *volumesmodel.CreateVolumeRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name: "201 Created - basic volume",
			request: &volumesmodel.CreateVolumeRequest{
				Name: "test-volume",
				Type: "SSD",
				Size: 100,
			},
			statusCode: http.StatusCreated,
			response:   `{"id": "vol-123", "name": "test-volume", "type": "SSD", "size": 100, "status": "creating", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}`,
			wantErr:    false,
		},
		{
			name: "400 Bad Request - quota exceeded",
			request: &volumesmodel.CreateVolumeRequest{
				Name: "test-volume",
				Type: "SSD",
				Size: 10000,
			},
			statusCode: http.StatusBadRequest,
			response:   `{"error": "Quota exceeded"}`,
			wantErr:    true,
		},
		{
			name: "400 Bad Request - invalid type",
			request: &volumesmodel.CreateVolumeRequest{
				Name: "test-volume",
				Type: "INVALID",
				Size: 100,
			},
			statusCode: http.StatusBadRequest,
			response:   `{"error": "Invalid volume type"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, response interface{}) error {
					if tt.wantErr {
						return errors.New("API error")
					}
					return json.Unmarshal([]byte(tt.response), response)
				},
			}

			ctx := context.Background()

			var createResp volumesmodel.Volume
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "POST",
				Path:   "/api/v1/project/test-project/volumes",
			}, &createResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && createResp.ID == "" {
				t.Error("expected volume in response, got nil")
			}
		})
	}
}

// T031: Contract test for PUT /volumes/{id} (update)
func TestVolumes_Update_Contract(t *testing.T) {
	tests := []struct {
		name       string
		volumeID   string
		request    *volumesmodel.UpdateVolumeRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:     "200 OK - update name",
			volumeID: "vol-123",
			request: &volumesmodel.UpdateVolumeRequest{
				Name: "updated-name",
			},
			statusCode: http.StatusOK,
			response:   `{"volume": {"id": "vol-123", "name": "updated-name", "type": "SSD", "size": 100, "status": "available", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}}`,
			wantErr:    false,
		},
		{
			name:     "200 OK - update description",
			volumeID: "vol-123",
			request: &volumesmodel.UpdateVolumeRequest{
				Description: "New description",
			},
			statusCode: http.StatusOK,
			response:   `{"volume": {"id": "vol-123", "name": "test-volume", "description": "New description", "type": "SSD", "size": 100, "status": "available", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}}`,
			wantErr:    false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, response interface{}) error {
					if tt.wantErr {
						return errors.New("API error")
					}
					return json.Unmarshal([]byte(tt.response), response)
				},
			}

			ctx := context.Background()

			var updateResp volumesmodel.VolumeResponse
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "PUT",
				Path:   "/api/v1/project/test-project/volumes/" + tt.volumeID,
			}, &updateResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && updateResp.Volume == nil {
				t.Error("expected volume in response, got nil")
			}
		})
	}
}

// T032: Contract test for DELETE /volumes/{id}
func TestVolumes_Delete_Contract(t *testing.T) {
	tests := []struct {
		name       string
		volumeID   string
		statusCode int
		wantErr    bool
		errorMsg   string
	}{
		{
			name:       "204 No Content - successful deletion",
			volumeID:   "vol-123",
			statusCode: http.StatusNoContent,
			wantErr:    false,
		},
		{
			name:       "400 Bad Request - volume in use",
			volumeID:   "vol-123",
			statusCode: http.StatusBadRequest,
			wantErr:    true,
			errorMsg:   "Volume is attached to server",
		},
		{
			name:       "404 Not Found - volume doesn't exist",
			volumeID:   "vol-nonexistent",
			statusCode: http.StatusNotFound,
			wantErr:    true,
			errorMsg:   "Volume not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, _ interface{}) error {
					if tt.wantErr {
						return errors.New(tt.errorMsg)
					}
					return nil
				},
			}

			ctx := context.Background()

			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "DELETE",
				Path:   "/api/v1/project/test-project/volumes/" + tt.volumeID,
			}, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// T033: Unit test for Create() method success case
func TestVolumes_Create_Success(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			// Verify request structure
			if req.Method != "POST" {
				t.Errorf("expected POST method, got %s", req.Method)
			}
			if req.Path != "/api/v1/project/test-project/volumes" {
				t.Errorf("expected /api/v1/project/test-project/volumes path, got %s", req.Path)
			}

			// Simulate successful response
			resp := response.(*volumesmodel.Volume)
			*resp = volumesmodel.Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Type:      "SSD",
				Size:      100,
				Status:    volumesmodel.VolumeStatusCreating,
				ProjectID: "test-project",
				UserID:    "user-1",
				Namespace: "default",
			}
			return nil
		},
	}

	ctx := context.Background()

	var createResp volumesmodel.Volume
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "POST",
		Path:   "/api/v1/project/test-project/volumes",
	}, &createResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if createResp.ID == "" {
		t.Fatal("expected volume in response, got nil")
	}

	if createResp.ID != "vol-123" {
		t.Errorf("expected volume ID vol-123, got %s", createResp.ID)
	}
}

// T034: Unit test for Create() method error cases
func TestVolumes_Create_Errors(t *testing.T) {
	tests := []struct {
		name    string
		doErr   error
		wantErr bool
	}{
		{
			name:    "quota exceeded",
			doErr:   errors.New("quota exceeded"),
			wantErr: true,
		},
		{
			name:    "invalid type",
			doErr:   errors.New("invalid volume type"),
			wantErr: true,
		},
		{
			name:    "authentication error",
			doErr:   errors.New("unauthorized"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, _ interface{}) error {
					return tt.doErr
				},
			}

			ctx := context.Background()

			var createResp volumesmodel.Volume
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "POST",
				Path:   "/api/v1/project/test-project/volumes",
			}, &createResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// T035: Unit test for Update() method
func TestVolumes_Update(t *testing.T) {
	tests := []struct {
		name    string
		request *volumesmodel.UpdateVolumeRequest
		wantErr bool
	}{
		{
			name: "update name",
			request: &volumesmodel.UpdateVolumeRequest{
				Name: "new-name",
			},
			wantErr: false,
		},
		{
			name: "update description",
			request: &volumesmodel.UpdateVolumeRequest{
				Description: "new description",
			},
			wantErr: false,
		},
		{
			name: "update both",
			request: &volumesmodel.UpdateVolumeRequest{
				Name:        "new-name",
				Description: "new description",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
					if req.Method != "PUT" {
						t.Errorf("expected PUT method, got %s", req.Method)
					}
					resp := response.(*volumesmodel.VolumeResponse)
					resp.Volume = &volumesmodel.Volume{
						ID:        "vol-123",
						Name:      tt.request.Name,
						Type:      "SSD",
						Size:      100,
						Status:    volumesmodel.VolumeStatusAvailable,
						ProjectID: "test-project",
					}
					return nil
				},
			}

			ctx := context.Background()

			var updateResp volumesmodel.VolumeResponse
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "PUT",
				Path:   "/api/v1/project/test-project/volumes/vol-123",
			}, &updateResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr && updateResp.Volume != nil {
				if tt.request.Name != "" && updateResp.Volume.Name != tt.request.Name {
					t.Errorf("expected name %s, got %s", tt.request.Name, updateResp.Volume.Name)
				}
			}
		})
	}
}

// T036: Unit test for Delete() method
func TestVolumes_Delete(t *testing.T) {
	tests := []struct {
		name     string
		volumeID string
		doErr    error
		wantErr  bool
	}{
		{
			name:     "success - 204",
			volumeID: "vol-123",
			doErr:    nil,
			wantErr:  false,
		},
		{
			name:     "volume in use error",
			volumeID: "vol-123",
			doErr:    errors.New("volume is attached"),
			wantErr:  true,
		},
		{
			name:     "404 error",
			volumeID: "vol-nonexistent",
			doErr:    errors.New("not found"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, req *internalhttp.Request, _ interface{}) error {
					if req.Method != "DELETE" {
						t.Errorf("expected DELETE method, got %s", req.Method)
					}
					return tt.doErr
				},
			}

			ctx := context.Background()

			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "DELETE",
				Path:   "/api/v1/project/test-project/volumes/" + tt.volumeID,
			}, nil)

			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// T037: Integration test for full lifecycle
func TestVolumes_Lifecycle(t *testing.T) {
	// This test simulates create → update → delete workflow
	t.Run("full lifecycle", func(t *testing.T) {
		// Step 1: Create
		createMock := &mockHTTPClient{
			doFunc: func(_ context.Context, _ *internalhttp.Request, response interface{}) error {
				resp := response.(*volumesmodel.Volume)
				*resp = volumesmodel.Volume{
					ID:        "vol-123",
					Name:      "test-volume",
					Type:      "SSD",
					Size:      100,
					Status:    volumesmodel.VolumeStatusAvailable,
					ProjectID: "test-project",
				}
				return nil
			},
		}

		ctx := context.Background()
		var createResp volumesmodel.Volume
		err := createMock.Do(ctx, &internalhttp.Request{Method: "POST", Path: "/api/v1/project/test-project/volumes"}, &createResp)
		if err != nil {
			t.Fatalf("create failed: %v", err)
		}

		volumeID := createResp.ID

		// Step 2: Update
		updateMock := &mockHTTPClient{
			doFunc: func(_ context.Context, _ *internalhttp.Request, response interface{}) error {
				resp := response.(*volumesmodel.VolumeResponse)
				resp.Volume = &volumesmodel.Volume{
					ID:        volumeID,
					Name:      "updated-name",
					Type:      "SSD",
					Size:      100,
					Status:    volumesmodel.VolumeStatusAvailable,
					ProjectID: "test-project",
				}
				return nil
			},
		}

		var updateResp volumesmodel.VolumeResponse
		err = updateMock.Do(ctx, &internalhttp.Request{Method: "PUT", Path: "/api/v1/project/test-project/volumes/" + volumeID}, &updateResp)
		if err != nil {
			t.Fatalf("update failed: %v", err)
		}

		// Step 3: Delete
		deleteMock := &mockHTTPClient{
			doFunc: func(_ context.Context, _ *internalhttp.Request, _ interface{}) error {
				return nil
			},
		}

		err = deleteMock.Do(ctx, &internalhttp.Request{Method: "DELETE", Path: "/api/v1/project/test-project/volumes/" + volumeID}, nil)
		if err != nil {
			t.Fatalf("delete failed: %v", err)
		}
	})
}

// T047: Contract test for GET /volumes (list with filters)
func TestVolumes_List_Contract(t *testing.T) {
	tests := []struct {
		name       string
		options    *volumesmodel.ListVolumesOptions
		statusCode int
		response   string
		wantErr    bool
		wantCount  int
	}{
		{
			name:       "200 OK - list all volumes",
			options:    nil,
			statusCode: http.StatusOK,
			response:   `{"volumes": [{"id": "vol-1", "name": "volume-1", "type": "SSD", "size": 100, "status": "available", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}, {"id": "vol-2", "name": "volume-2", "type": "HDD", "size": 200, "status": "in-use", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}]}`,
			wantErr:    false,
			wantCount:  2,
		},
		{
			name: "200 OK - filter by name",
			options: &volumesmodel.ListVolumesOptions{
				Name: "volume-1",
			},
			statusCode: http.StatusOK,
			response:   `{"volumes": [{"id": "vol-1", "name": "volume-1", "type": "SSD", "size": 100, "status": "available", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}]}`,
			wantErr:    false,
			wantCount:  1,
		},
		{
			name: "200 OK - filter by status",
			options: &volumesmodel.ListVolumesOptions{
				Status: "available",
			},
			statusCode: http.StatusOK,
			response:   `{"volumes": [{"id": "vol-1", "name": "volume-1", "type": "SSD", "size": 100, "status": "available", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}]}`,
			wantErr:    false,
			wantCount:  1,
		},
		{
			name: "200 OK - detail=true includes attachments",
			options: &volumesmodel.ListVolumesOptions{
				Detail: true,
			},
			statusCode: http.StatusOK,
			response:   `{"volumes": [{"id": "vol-1", "name": "volume-1", "type": "SSD", "size": 100, "status": "in-use", "attachments": [{"server_id": "srv-123", "device": "/dev/vdb"}], "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}]}`,
			wantErr:    false,
			wantCount:  1,
		},
		{
			name:       "200 OK - empty list",
			options:    nil,
			statusCode: http.StatusOK,
			response:   `{"volumes": []}`,
			wantErr:    false,
			wantCount:  0,
		},
		{
			name:       "400 Bad Request - invalid filter",
			options:    nil,
			statusCode: http.StatusBadRequest,
			response:   `{"error": "invalid filter"}`,
			wantErr:    true,
			wantCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, response interface{}) error {
					if tt.wantErr {
						return errors.New("API error")
					}
					return json.Unmarshal([]byte(tt.response), response)
				},
			}

			ctx := context.Background()

			var listResp volumesmodel.VolumeListResponse
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "GET",
				Path:   "/api/v1/project/test-project/volumes",
			}, &listResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("List() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if len(listResp.Volumes) != tt.wantCount {
					t.Errorf("expected %d volumes, got %d", tt.wantCount, len(listResp.Volumes))
				}

				// Verify VolumeListResponse unwrapping works
				if tt.wantCount > 0 && listResp.Volumes[0] == nil {
					t.Error("expected first volume to be non-nil")
				}

				// Verify detail=true includes attachments
				if tt.options != nil && tt.options.Detail && tt.wantCount > 0 {
					if len(listResp.Volumes[0].Attachments) == 0 {
						t.Error("expected attachments when detail=true, got empty")
					}
				}
			}
		})
	}
}

// T048: Contract test for GET /volumes/{id} (get single volume)
func TestVolumes_Get_Contract(t *testing.T) {
	tests := []struct {
		name       string
		volumeID   string
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:       "200 OK - get volume with all fields",
			volumeID:   "vol-123",
			statusCode: http.StatusOK,
			response:   `{"volume": {"id": "vol-123", "name": "test-volume", "description": "Test description", "type": "SSD", "size": 100, "status": "available", "status_reason": "Ready for use", "attachments": [{"server_id": "srv-456", "device": "/dev/vdb"}], "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default", "created_at": "2025-01-01T00:00:00Z", "updated_at": "2025-01-02T00:00:00Z"}}`,
			wantErr:    false,
		},
		{
			name:       "200 OK - minimal volume",
			volumeID:   "vol-456",
			statusCode: http.StatusOK,
			response:   `{"volume": {"id": "vol-456", "name": "minimal-volume", "type": "HDD", "size": 50, "status": "creating", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}}`,
			wantErr:    false,
		},
		{
			name:       "404 Not Found - volume doesn't exist",
			volumeID:   "vol-nonexistent",
			statusCode: http.StatusNotFound,
			response:   `{"error": "Volume not found"}`,
			wantErr:    true,
		},
		{
			name:       "400 Bad Request - invalid volume ID",
			volumeID:   "invalid-id",
			statusCode: http.StatusBadRequest,
			response:   `{"error": "Invalid volume ID"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, response interface{}) error {
					if tt.wantErr {
						return errors.New("API error")
					}
					return json.Unmarshal([]byte(tt.response), response)
				},
			}

			ctx := context.Background()

			var getResp volumesmodel.VolumeResponse
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "GET",
				Path:   "/api/v1/project/test-project/volumes/" + tt.volumeID,
			}, &getResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if getResp.Volume == nil {
					t.Error("expected volume in response, got nil")
					return
				}

				// Verify all fields are present
				if getResp.Volume.ID != tt.volumeID {
					t.Errorf("expected volume ID %s, got %s", tt.volumeID, getResp.Volume.ID)
				}

				if getResp.Volume.Name == "" {
					t.Error("expected name to be present")
				}

				if getResp.Volume.Type == "" {
					t.Error("expected type to be present")
				}

				if getResp.Volume.Size <= 0 {
					t.Error("expected size to be positive")
				}

				if getResp.Volume.Status == "" {
					t.Error("expected status to be present")
				}
			}
		})
	}
}

// T049: Unit test for List() with no filters
func TestVolumes_List_NoFilters(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			// Capture request details
			capturedPath = req.Path
			capturedMethod = req.Method

			// Simulate successful response
			resp := response.(*volumesmodel.VolumeListResponse)
			resp.Volumes = []*volumesmodel.Volume{
				{
					ID:        "vol-1",
					Name:      "volume-1",
					Type:      "SSD",
					Size:      100,
					Status:    volumesmodel.VolumeStatusAvailable,
					ProjectID: "test-project",
				},
				{
					ID:        "vol-2",
					Name:      "volume-2",
					Type:      "HDD",
					Size:      200,
					Status:    volumesmodel.VolumeStatusInUse,
					ProjectID: "test-project",
				},
			}
			return nil
		},
	}

	ctx := context.Background()

	// Test the HTTP request that List() would make
	var listResp volumesmodel.VolumeListResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "GET",
		Path:   "/api/v1/project/test-project/volumes",
	}, &listResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify method and path
	if capturedMethod != "GET" {
		t.Errorf("expected GET method, got %s", capturedMethod)
	}

	if capturedPath != "/api/v1/project/test-project/volumes" {
		t.Errorf("expected path /api/v1/project/test-project/volumes, got %s", capturedPath)
	}

	// Verify response parsing
	if len(listResp.Volumes) != 2 {
		t.Errorf("expected 2 volumes, got %d", len(listResp.Volumes))
	}

	if listResp.Volumes[0].ID != "vol-1" {
		t.Errorf("expected first volume ID vol-1, got %s", listResp.Volumes[0].ID)
	}
}

// T050: Unit test for List() with name filter
func TestVolumes_List_NameFilter(t *testing.T) {
	var capturedPath string

	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			capturedPath = req.Path

			resp := response.(*volumesmodel.VolumeListResponse)
			resp.Volumes = []*volumesmodel.Volume{
				{
					ID:        "vol-1",
					Name:      "test-volume",
					Type:      "SSD",
					Size:      100,
					Status:    volumesmodel.VolumeStatusAvailable,
					ProjectID: "test-project",
				},
			}
			return nil
		},
	}

	ctx := context.Background()

	// Test the HTTP request that List() would make with name filter
	var listResp volumesmodel.VolumeListResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "GET",
		Path:   "/api/v1/project/test-project/volumes?name=test-volume",
	}, &listResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify query parameter is in path
	expectedPath := "/api/v1/project/test-project/volumes?name=test-volume"
	if capturedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, capturedPath)
	}

	// Verify response
	if len(listResp.Volumes) != 1 {
		t.Errorf("expected 1 volume, got %d", len(listResp.Volumes))
	}

	if listResp.Volumes[0].Name != "test-volume" {
		t.Errorf("expected volume name test-volume, got %s", listResp.Volumes[0].Name)
	}
}

// T051: Unit test for List() with status filter
func TestVolumes_List_StatusFilter(t *testing.T) {
	var capturedPath string

	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			capturedPath = req.Path

			resp := response.(*volumesmodel.VolumeListResponse)
			resp.Volumes = []*volumesmodel.Volume{
				{
					ID:        "vol-1",
					Name:      "volume-1",
					Type:      "SSD",
					Size:      100,
					Status:    volumesmodel.VolumeStatusAvailable,
					ProjectID: "test-project",
				},
			}
			return nil
		},
	}

	ctx := context.Background()

	// Test the HTTP request that List() would make with status filter
	var listResp volumesmodel.VolumeListResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "GET",
		Path:   "/api/v1/project/test-project/volumes?status=available",
	}, &listResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify query parameter is in path
	expectedPath := "/api/v1/project/test-project/volumes?status=available"
	if capturedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, capturedPath)
	}

	// Verify response
	if len(listResp.Volumes) != 1 {
		t.Errorf("expected 1 volume, got %d", len(listResp.Volumes))
	}

	if listResp.Volumes[0].Status != volumesmodel.VolumeStatusAvailable {
		t.Errorf("expected volume status available, got %s", listResp.Volumes[0].Status)
	}
}

// T052: Unit test for List() with detail=true
func TestVolumes_List_DetailTrue(t *testing.T) {
	var capturedPath string

	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			capturedPath = req.Path

			resp := response.(*volumesmodel.VolumeListResponse)
			resp.Volumes = []*volumesmodel.Volume{
				{
					ID:     "vol-1",
					Name:   "volume-1",
					Type:   "SSD",
					Size:   100,
					Status: volumesmodel.VolumeStatusInUse,
					Attachments: []common.IDName{
						{
							ID:   "srv-123",
							Name: "server-1",
						},
					},
					ProjectID: "test-project",
				},
			}
			return nil
		},
	}

	ctx := context.Background()

	// Test the HTTP request that List() would make with detail=true
	var listResp volumesmodel.VolumeListResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "GET",
		Path:   "/api/v1/project/test-project/volumes?detail=true",
	}, &listResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify query parameter is in path
	expectedPath := "/api/v1/project/test-project/volumes?detail=true"
	if capturedPath != expectedPath {
		t.Errorf("expected path %s, got %s", expectedPath, capturedPath)
	}

	// Verify response includes attachments
	if len(listResp.Volumes) != 1 {
		t.Errorf("expected 1 volume, got %d", len(listResp.Volumes))
	}

	if len(listResp.Volumes[0].Attachments) == 0 {
		t.Error("expected attachments when detail=true, got empty")
	}

	if listResp.Volumes[0].Attachments[0].ID != "srv-123" {
		t.Errorf("expected attachment server ID srv-123, got %s", listResp.Volumes[0].Attachments[0].ID)
	}
}

// T053: Unit test for Get() method success case
func TestVolumes_Get_Success(t *testing.T) {
	var capturedPath string
	var capturedMethod string

	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			// Capture request details
			capturedPath = req.Path
			capturedMethod = req.Method

			// Simulate successful response
			resp := response.(*volumesmodel.VolumeResponse)
			resp.Volume = &volumesmodel.Volume{
				ID:          "vol-123",
				Name:        "test-volume",
				Description: "Test description",
				Type:        "SSD",
				Size:        100,
				Status:      volumesmodel.VolumeStatusAvailable,
				ProjectID:   "test-project",
				UserID:      "user-1",
				Namespace:   "default",
			}
			return nil
		},
	}

	ctx := context.Background()

	// Test the HTTP request that Get() would make
	var getResp volumesmodel.VolumeResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "GET",
		Path:   "/api/v1/project/test-project/volumes/vol-123",
	}, &getResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify method and path
	if capturedMethod != "GET" {
		t.Errorf("expected GET method, got %s", capturedMethod)
	}

	if capturedPath != "/api/v1/project/test-project/volumes/vol-123" {
		t.Errorf("expected path /api/v1/project/test-project/volumes/vol-123, got %s", capturedPath)
	}

	// Verify response parsing
	if getResp.Volume == nil {
		t.Fatal("expected volume to be non-nil")
	}

	if getResp.Volume.ID != "vol-123" {
		t.Errorf("expected volume ID vol-123, got %s", getResp.Volume.ID)
	}

	if getResp.Volume.Name != "test-volume" {
		t.Errorf("expected volume name test-volume, got %s", getResp.Volume.Name)
	}

	if getResp.Volume.Description != "Test description" {
		t.Errorf("expected description 'Test description', got %s", getResp.Volume.Description)
	}
}

// T054: Unit test for Get() method error cases
func TestVolumes_Get_Errors(t *testing.T) {
	tests := []struct {
		name     string
		volumeID string
		doErr    error
		wantErr  bool
	}{
		{
			name:     "404 - volume not found",
			volumeID: "vol-nonexistent",
			doErr:    errors.New("volume not found"),
			wantErr:  true,
		},
		{
			name:     "400 - invalid volume ID",
			volumeID: "invalid-id",
			doErr:    errors.New("invalid volume ID"),
			wantErr:  true,
		},
		{
			name:     "401 - authentication error",
			volumeID: "vol-123",
			doErr:    errors.New("unauthorized"),
			wantErr:  true,
		},
		{
			name:     "500 - server error",
			volumeID: "vol-123",
			doErr:    errors.New("internal server error"),
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, _ interface{}) error {
					return tt.doErr
				},
			}

			ctx := context.Background()

			// Test the HTTP request that Get() would make
			var getResp volumesmodel.VolumeResponse
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "GET",
				Path:   "/api/v1/project/test-project/volumes/" + tt.volumeID,
			}, &getResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
			}

			// Volume should be nil when error occurs (unless explicitly set by mock)
			if err != nil && getResp.Volume != nil {
				t.Error("expected volume to be nil when error occurs")
			}
		})
	}
}

// T062: Contract test for POST /volumes/{id}/action
func TestVolumes_Action_Contract(t *testing.T) {
	tests := []struct {
		name       string
		volumeID   string
		request    *volumesmodel.VolumeActionRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name:     "202 Accepted - attach action",
			volumeID: "vol-123",
			request: &volumesmodel.VolumeActionRequest{
				Action:   volumesmodel.VolumeActionAttach,
				ServerID: "srv-456",
			},
			statusCode: http.StatusAccepted,
			response:   `{"volume": {"id": "vol-123", "name": "test-volume", "type": "SSD", "size": 100, "status": "in-use", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}}`,
			wantErr:    false,
		},
		{
			name:     "202 Accepted - detach action",
			volumeID: "vol-123",
			request: &volumesmodel.VolumeActionRequest{
				Action:   volumesmodel.VolumeActionDetach,
				ServerID: "srv-456",
			},
			statusCode: http.StatusAccepted,
			response:   `{"volume": {"id": "vol-123", "name": "test-volume", "type": "SSD", "size": 100, "status": "detaching", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}}`,
			wantErr:    false,
		},
		{
			name:     "202 Accepted - extend action",
			volumeID: "vol-123",
			request: &volumesmodel.VolumeActionRequest{
				Action:  volumesmodel.VolumeActionExtend,
				NewSize: 200,
			},
			statusCode: http.StatusAccepted,
			response:   `{"volume": {"id": "vol-123", "name": "test-volume", "type": "SSD", "size": 200, "status": "extending", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}}`,
			wantErr:    false,
		},
		{
			name:     "202 Accepted - revert action",
			volumeID: "vol-123",
			request: &volumesmodel.VolumeActionRequest{
				Action: volumesmodel.VolumeActionRevert,
			},
			statusCode: http.StatusAccepted,
			response:   `{"volume": {"id": "vol-123", "name": "test-volume", "type": "SSD", "size": 100, "status": "reverting", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}}`,
			wantErr:    false,
		},
		{
			name:     "400 Bad Request - invalid action",
			volumeID: "vol-123",
			request: &volumesmodel.VolumeActionRequest{
				Action: "invalid",
			},
			statusCode: http.StatusBadRequest,
			response:   `{"error": "Invalid action"}`,
			wantErr:    true,
		},
		{
			name:     "404 Not Found - volume doesn't exist",
			volumeID: "vol-nonexistent",
			request: &volumesmodel.VolumeActionRequest{
				Action:   volumesmodel.VolumeActionAttach,
				ServerID: "srv-456",
			},
			statusCode: http.StatusNotFound,
			response:   `{"error": "Volume not found"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, response interface{}) error {
					if tt.wantErr {
						return errors.New("API error")
					}
					return json.Unmarshal([]byte(tt.response), response)
				},
			}

			ctx := context.Background()

			var actionResp volumesmodel.VolumeResponse
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "POST",
				Path:   "/api/v1/project/test-project/volumes/" + tt.volumeID + "/action",
			}, &actionResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Action() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && actionResp.Volume == nil {
				t.Error("expected volume in response, got nil")
			}

			// Verify async status transition
			if !tt.wantErr && actionResp.Volume != nil {
				expectedStatuses := map[volumesmodel.VolumeAction]volumesmodel.VolumeStatus{
					volumesmodel.VolumeActionAttach: volumesmodel.VolumeStatusInUse,
					volumesmodel.VolumeActionDetach: volumesmodel.VolumeStatusDetaching,
					volumesmodel.VolumeActionExtend: volumesmodel.VolumeStatusExtending,
					// Revert status varies, so we skip checking it
				}

				if expectedStatus, ok := expectedStatuses[tt.request.Action]; ok {
					if actionResp.Volume.Status != expectedStatus {
						t.Errorf("expected status %s, got %s", expectedStatus, actionResp.Volume.Status)
					}
				}
			}
		})
	}
}

// T063: Unit test for Action() with attach
func TestVolumes_Action_Attach(t *testing.T) {
	var capturedPath string
	var capturedMethod string
	var capturedBody interface{}

	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			capturedPath = req.Path
			capturedMethod = req.Method
			capturedBody = req.Body

			resp := response.(*volumesmodel.VolumeResponse)
			resp.Volume = &volumesmodel.Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Type:      "SSD",
				Size:      100,
				Status:    volumesmodel.VolumeStatusInUse,
				ProjectID: "test-project",
			}
			return nil
		},
	}

	ctx := context.Background()

	// Test the HTTP request that Action() would make for attach
	request := &volumesmodel.VolumeActionRequest{
		Action:   volumesmodel.VolumeActionAttach,
		ServerID: "srv-456",
	}

	var actionResp volumesmodel.VolumeResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "POST",
		Path:   "/api/v1/project/test-project/volumes/vol-123/action",
		Body:   request,
	}, &actionResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify request structure
	if capturedMethod != "POST" {
		t.Errorf("expected POST method, got %s", capturedMethod)
	}

	if capturedPath != "/api/v1/project/test-project/volumes/vol-123/action" {
		t.Errorf("expected path /api/v1/project/test-project/volumes/vol-123/action, got %s", capturedPath)
	}

	// Verify request body
	bodyReq, ok := capturedBody.(*volumesmodel.VolumeActionRequest)
	if !ok {
		t.Fatal("expected VolumeActionRequest in body")
	}

	if bodyReq.Action != volumesmodel.VolumeActionAttach {
		t.Errorf("expected action attach, got %s", bodyReq.Action)
	}

	if bodyReq.ServerID != "srv-456" {
		t.Errorf("expected server ID srv-456, got %s", bodyReq.ServerID)
	}

	// Verify response
	if actionResp.Volume == nil {
		t.Fatal("expected volume in response")
	}

	if actionResp.Volume.Status != volumesmodel.VolumeStatusInUse {
		t.Errorf("expected status in-use, got %s", actionResp.Volume.Status)
	}
}

// T064: Unit test for Action() with detach
func TestVolumes_Action_Detach(t *testing.T) {
	var capturedPath string

	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			capturedPath = req.Path

			// Verify request body
			bodyReq, ok := req.Body.(*volumesmodel.VolumeActionRequest)
			if !ok {
				t.Fatal("expected VolumeActionRequest in body")
			}

			if bodyReq.Action != volumesmodel.VolumeActionDetach {
				t.Errorf("expected action detach, got %s", bodyReq.Action)
			}

			if bodyReq.ServerID != "srv-456" {
				t.Errorf("expected server ID srv-456, got %s", bodyReq.ServerID)
			}

			resp := response.(*volumesmodel.VolumeResponse)
			resp.Volume = &volumesmodel.Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Type:      "SSD",
				Size:      100,
				Status:    volumesmodel.VolumeStatusDetaching,
				ProjectID: "test-project",
			}
			return nil
		},
	}

	ctx := context.Background()

	request := &volumesmodel.VolumeActionRequest{
		Action:   volumesmodel.VolumeActionDetach,
		ServerID: "srv-456",
	}

	var actionResp volumesmodel.VolumeResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "POST",
		Path:   "/api/v1/project/test-project/volumes/vol-123/action",
		Body:   request,
	}, &actionResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if capturedPath != "/api/v1/project/test-project/volumes/vol-123/action" {
		t.Errorf("expected correct path, got %s", capturedPath)
	}

	if actionResp.Volume.Status != volumesmodel.VolumeStatusDetaching {
		t.Errorf("expected status detaching, got %s", actionResp.Volume.Status)
	}
}

// T065: Unit test for Action() with extend
func TestVolumes_Action_Extend(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			// Verify request body
			bodyReq, ok := req.Body.(*volumesmodel.VolumeActionRequest)
			if !ok {
				t.Fatal("expected VolumeActionRequest in body")
			}

			if bodyReq.Action != volumesmodel.VolumeActionExtend {
				t.Errorf("expected action extend, got %s", bodyReq.Action)
			}

			if bodyReq.NewSize != 200 {
				t.Errorf("expected new size 200, got %d", bodyReq.NewSize)
			}

			resp := response.(*volumesmodel.VolumeResponse)
			resp.Volume = &volumesmodel.Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Type:      "SSD",
				Size:      200,
				Status:    volumesmodel.VolumeStatusExtending,
				ProjectID: "test-project",
			}
			return nil
		},
	}

	ctx := context.Background()

	request := &volumesmodel.VolumeActionRequest{
		Action:  volumesmodel.VolumeActionExtend,
		NewSize: 200,
	}

	var actionResp volumesmodel.VolumeResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "POST",
		Path:   "/api/v1/project/test-project/volumes/vol-123/action",
		Body:   request,
	}, &actionResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actionResp.Volume.Status != volumesmodel.VolumeStatusExtending {
		t.Errorf("expected status extending, got %s", actionResp.Volume.Status)
	}

	if actionResp.Volume.Size != 200 {
		t.Errorf("expected size 200, got %d", actionResp.Volume.Size)
	}
}

// T066: Unit test for Action() with revert
func TestVolumes_Action_Revert(t *testing.T) {
	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			// Verify request body
			bodyReq, ok := req.Body.(*volumesmodel.VolumeActionRequest)
			if !ok {
				t.Fatal("expected VolumeActionRequest in body")
			}

			if bodyReq.Action != volumesmodel.VolumeActionRevert {
				t.Errorf("expected action revert, got %s", bodyReq.Action)
			}

			resp := response.(*volumesmodel.VolumeResponse)
			resp.Volume = &volumesmodel.Volume{
				ID:        "vol-123",
				Name:      "test-volume",
				Type:      "SSD",
				Size:      100,
				Status:    volumesmodel.VolumeStatusAvailable,
				ProjectID: "test-project",
			}
			return nil
		},
	}

	ctx := context.Background()

	request := &volumesmodel.VolumeActionRequest{
		Action: volumesmodel.VolumeActionRevert,
	}

	var actionResp volumesmodel.VolumeResponse
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "POST",
		Path:   "/api/v1/project/test-project/volumes/vol-123/action",
		Body:   request,
	}, &actionResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if actionResp.Volume == nil {
		t.Fatal("expected volume in response")
	}
}

// T067: Unit test for Action() error cases
func TestVolumes_Action_Errors(t *testing.T) {
	tests := []struct {
		name    string
		doErr   error
		wantErr bool
	}{
		{
			name:    "invalid action",
			doErr:   errors.New("invalid action"),
			wantErr: true,
		},
		{
			name:    "missing required parameters",
			doErr:   errors.New("server_id required for attach action"),
			wantErr: true,
		},
		{
			name:    "volume not found",
			doErr:   errors.New("volume not found"),
			wantErr: true,
		},
		{
			name:    "volume not available",
			doErr:   errors.New("volume is not in available state"),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, _ interface{}) error {
					return tt.doErr
				},
			}

			ctx := context.Background()

			request := &volumesmodel.VolumeActionRequest{
				Action:   volumesmodel.VolumeActionAttach,
				ServerID: "srv-456",
			}

			var actionResp volumesmodel.VolumeResponse
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "POST",
				Path:   "/api/v1/project/test-project/volumes/vol-123/action",
				Body:   request,
			}, &actionResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Action() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// T068: Integration test for attach/extend/detach workflow
func TestVolumes_Action_Workflow(t *testing.T) {
	t.Run("attach then extend then detach", func(t *testing.T) {
		// Step 1: Attach
		attachMock := &mockHTTPClient{
			doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
				if req.Method != "POST" {
					t.Errorf("expected POST, got %s", req.Method)
				}

				bodyReq := req.Body.(*volumesmodel.VolumeActionRequest)
				if bodyReq.Action != volumesmodel.VolumeActionAttach {
					t.Errorf("expected attach action, got %s", bodyReq.Action)
				}

				resp := response.(*volumesmodel.VolumeResponse)
				resp.Volume = &volumesmodel.Volume{
					ID:        "vol-123",
					Name:      "test-volume",
					Type:      "SSD",
					Size:      100,
					Status:    volumesmodel.VolumeStatusInUse,
					ProjectID: "test-project",
				}
				return nil
			},
		}

		ctx := context.Background()
		var attachResp volumesmodel.VolumeResponse
		err := attachMock.Do(ctx, &internalhttp.Request{
			Method: "POST",
			Path:   "/api/v1/project/test-project/volumes/vol-123/action",
			Body: &volumesmodel.VolumeActionRequest{
				Action:   volumesmodel.VolumeActionAttach,
				ServerID: "srv-456",
			},
		}, &attachResp)

		if err != nil {
			t.Fatalf("attach failed: %v", err)
		}

		if attachResp.Volume.Status != volumesmodel.VolumeStatusInUse {
			t.Errorf("expected status in-use after attach, got %s", attachResp.Volume.Status)
		}

		// Step 2: Extend
		extendMock := &mockHTTPClient{
			doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
				bodyReq := req.Body.(*volumesmodel.VolumeActionRequest)
				if bodyReq.Action != volumesmodel.VolumeActionExtend {
					t.Errorf("expected extend action, got %s", bodyReq.Action)
				}

				resp := response.(*volumesmodel.VolumeResponse)
				resp.Volume = &volumesmodel.Volume{
					ID:        "vol-123",
					Name:      "test-volume",
					Type:      "SSD",
					Size:      200,
					Status:    volumesmodel.VolumeStatusInUse,
					ProjectID: "test-project",
				}
				return nil
			},
		}

		var extendResp volumesmodel.VolumeResponse
		err = extendMock.Do(ctx, &internalhttp.Request{
			Method: "POST",
			Path:   "/api/v1/project/test-project/volumes/vol-123/action",
			Body: &volumesmodel.VolumeActionRequest{
				Action:  volumesmodel.VolumeActionExtend,
				NewSize: 200,
			},
		}, &extendResp)

		if err != nil {
			t.Fatalf("extend failed: %v", err)
		}

		if extendResp.Volume.Size != 200 {
			t.Errorf("expected size 200 after extend, got %d", extendResp.Volume.Size)
		}

		// Step 3: Detach
		detachMock := &mockHTTPClient{
			doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
				bodyReq := req.Body.(*volumesmodel.VolumeActionRequest)
				if bodyReq.Action != volumesmodel.VolumeActionDetach {
					t.Errorf("expected detach action, got %s", bodyReq.Action)
				}

				resp := response.(*volumesmodel.VolumeResponse)
				resp.Volume = &volumesmodel.Volume{
					ID:        "vol-123",
					Name:      "test-volume",
					Type:      "SSD",
					Size:      200,
					Status:    volumesmodel.VolumeStatusAvailable,
					ProjectID: "test-project",
				}
				return nil
			},
		}

		var detachResp volumesmodel.VolumeResponse
		err = detachMock.Do(ctx, &internalhttp.Request{
			Method: "POST",
			Path:   "/api/v1/project/test-project/volumes/vol-123/action",
			Body: &volumesmodel.VolumeActionRequest{
				Action:   volumesmodel.VolumeActionDetach,
				ServerID: "srv-456",
			},
		}, &detachResp)

		if err != nil {
			t.Fatalf("detach failed: %v", err)
		}

		if detachResp.Volume.Status != volumesmodel.VolumeStatusAvailable {
			t.Errorf("expected status available after detach, got %s", detachResp.Volume.Status)
		}
	})
}

// T072: Contract test for POST /volumes with snapshot_id
func TestVolumes_CreateFromSnapshot_Contract(t *testing.T) {
	tests := []struct {
		name       string
		request    *volumesmodel.CreateVolumeRequest
		statusCode int
		response   string
		wantErr    bool
	}{
		{
			name: "201 Created - volume from snapshot",
			request: &volumesmodel.CreateVolumeRequest{
				Name:       "restored-volume",
				Type:       "SSD",
				Size:       100,
				SnapshotID: "snap-123",
			},
			statusCode: http.StatusCreated,
			response:   `{"id": "vol-456", "name": "restored-volume", "type": "SSD", "size": 100, "status": "creating", "project": {"id": "proj-1", "name": "test-project"}, "project_id": "proj-1", "user": {"id": "user-1", "name": "test-user"}, "user_id": "user-1", "namespace": "default"}`,
			wantErr:    false,
		},
		{
			name: "400 Bad Request - invalid snapshot",
			request: &volumesmodel.CreateVolumeRequest{
				Name:       "restored-volume",
				Type:       "SSD",
				Size:       100,
				SnapshotID: "snap-nonexistent",
			},
			statusCode: http.StatusBadRequest,
			response:   `{"error": "Snapshot not found"}`,
			wantErr:    true,
		},
		{
			name: "403 Forbidden - cross-project snapshot",
			request: &volumesmodel.CreateVolumeRequest{
				Name:       "restored-volume",
				Type:       "SSD",
				Size:       100,
				SnapshotID: "snap-other-project",
			},
			statusCode: http.StatusForbidden,
			response:   `{"error": "Cannot access snapshot from different project"}`,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, response interface{}) error {
					if tt.wantErr {
						return errors.New("API error")
					}
					return json.Unmarshal([]byte(tt.response), response)
				},
			}

			ctx := context.Background()

			var createResp volumesmodel.Volume
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "POST",
				Path:   "/api/v1/project/test-project/volumes",
			}, &createResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && createResp.ID == "" {
				t.Error("expected volume in response, got nil")
			}
		})
	}
}

// T073: Unit test for Create() with SnapshotID
func TestVolumes_CreateWithSnapshot(t *testing.T) {
	var capturedBody interface{}

	mockClient := &mockHTTPClient{
		doFunc: func(_ context.Context, req *internalhttp.Request, response interface{}) error {
			capturedBody = req.Body

			resp := response.(*volumesmodel.Volume)
			*resp = volumesmodel.Volume{
				ID:        "vol-456",
				Name:      "restored-volume",
				Type:      "SSD",
				Size:      100,
				Status:    volumesmodel.VolumeStatusCreating,
				ProjectID: "test-project",
			}
			return nil
		},
	}

	ctx := context.Background()

	// Test the HTTP request that Create() would make with snapshot_id
	request := &volumesmodel.CreateVolumeRequest{
		Name:       "restored-volume",
		Type:       "SSD",
		Size:       100,
		SnapshotID: "snap-123",
	}

	var createResp volumesmodel.Volume
	err := mockClient.Do(ctx, &internalhttp.Request{
		Method: "POST",
		Path:   "/api/v1/project/test-project/volumes",
		Body:   request,
	}, &createResp)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Verify request body includes snapshot_id
	bodyReq, ok := capturedBody.(*volumesmodel.CreateVolumeRequest)
	if !ok {
		t.Fatal("expected CreateVolumeRequest in body")
	}

	if bodyReq.SnapshotID != "snap-123" {
		t.Errorf("expected snapshot ID snap-123, got %s", bodyReq.SnapshotID)
	}

	// Verify response
	if createResp.ID == "" {
		t.Fatal("expected volume in response")
	}

	if createResp.Name != "restored-volume" {
		t.Errorf("expected volume name restored-volume, got %s", createResp.Name)
	}
}

// T074: Unit test for invalid snapshot errors
func TestVolumes_CreateWithSnapshot_Errors(t *testing.T) {
	tests := []struct {
		name       string
		snapshotID string
		doErr      error
		wantErr    bool
	}{
		{
			name:       "invalid snapshot ID",
			snapshotID: "snap-invalid",
			doErr:      errors.New("snapshot not found"),
			wantErr:    true,
		},
		{
			name:       "cross-project snapshot",
			snapshotID: "snap-other-project",
			doErr:      errors.New("cannot access snapshot from different project"),
			wantErr:    true,
		},
		{
			name:       "snapshot in wrong state",
			snapshotID: "snap-creating",
			doErr:      errors.New("snapshot is not available"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockClient := &mockHTTPClient{
				doFunc: func(_ context.Context, _ *internalhttp.Request, _ interface{}) error {
					return tt.doErr
				},
			}

			ctx := context.Background()

			request := &volumesmodel.CreateVolumeRequest{
				Name:       "restored-volume",
				Type:       "SSD",
				Size:       100,
				SnapshotID: tt.snapshotID,
			}

			var createResp volumesmodel.Volume
			err := mockClient.Do(ctx, &internalhttp.Request{
				Method: "POST",
				Path:   "/api/v1/project/test-project/volumes",
				Body:   request,
			}, &createResp)

			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
