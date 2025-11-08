package networks

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/networks"
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
	if client.baseClient == nil {
		t.Error("expected baseClient to be initialized")
	}
}

// TestClient_List tests the List method
func TestClient_List(t *testing.T) {
	tests := []struct {
		name           string
		opts           *networks.ListNetworksOptions
		mockResponse   *networks.NetworkListResponse
		mockStatusCode int
		wantErr        bool
		checkPath      string
	}{
		{
			name: "list all networks",
			opts: nil,
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "network-1", CIDR: "10.0.1.0/24"},
					{ID: "net-2", Name: "network-2", CIDR: "10.0.2.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks",
		},
		{
			name: "list with name filter",
			opts: &networks.ListNetworksOptions{Name: "test-network"},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "test-network", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?name=test-network",
		},
		{
			name: "list with user_id filter",
			opts: &networks.ListNetworksOptions{UserID: "user-123"},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "network-1", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?user_id=user-123",
		},
		{
			name: "list with status filter",
			opts: &networks.ListNetworksOptions{Status: "ACTIVE"},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "network-1", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?status=ACTIVE",
		},
		{
			name: "list with router_id filter",
			opts: &networks.ListNetworksOptions{RouterID: "router-456"},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "network-1", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?router_id=router-456",
		},
		{
			name: "list with detail true",
			opts: &networks.ListNetworksOptions{Detail: func() *bool { b := true; return &b }()},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "network-1", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?detail=true",
		},
		{
			name: "list with detail false",
			opts: &networks.ListNetworksOptions{Detail: func() *bool { b := false; return &b }()},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "network-1", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?detail=false",
		},
		{
			name: "list with multiple filters",
			opts: &networks.ListNetworksOptions{
				Name:     "test-network",
				Status:   "ACTIVE",
				RouterID: "router-123",
			},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "test-network", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?name=test-network&status=ACTIVE&router_id=router-123",
		},
		{
			name: "list with all filters including detail",
			opts: &networks.ListNetworksOptions{
				Name:     "full-test",
				UserID:   "user-999",
				Status:   "BUILD",
				RouterID: "router-999",
				Detail:   func() *bool { b := true; return &b }(),
			},
			mockResponse: &networks.NetworkListResponse{
				Networks: []*networks.Network{
					{ID: "net-1", Name: "full-test", CIDR: "10.0.1.0/24"},
				},
			},
			mockStatusCode: http.StatusOK,
			wantErr:        false,
			checkPath:      "/api/v1/project/proj-123/networks?name=full-test&user_id=user-999&status=BUILD&router_id=router-999&detail=true",
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

			if len(result.Networks) != len(tt.mockResponse.Networks) {
				t.Errorf("expected %d networks, got %d", len(tt.mockResponse.Networks), len(result.Networks))
			}
		})
	}
}

// TestClient_Create tests the Create method
func TestClient_Create(t *testing.T) {
	tests := []struct {
		name            string
		request         *networks.NetworkCreateRequest
		mockResponse    *networks.Network
		mockStatusCode  int
		wantErr         bool
		validateRequest func(*testing.T, map[string]interface{})
	}{
		{
			name: "create network successfully",
			request: &networks.NetworkCreateRequest{
				Name:        "test-network",
				Description: "Test network",
				CIDR:        "10.0.1.0/24",
				Gateway:     "10.0.1.1",
				RouterID:    "router-123",
			},
			mockResponse: &networks.Network{
				ID:          "net-123",
				Name:        "test-network",
				Description: "Test network",
				CIDR:        "10.0.1.0/24",
				ProjectID:   "proj-123",
				RouterID:    "router-123",
				Gateway:     "10.0.1.1",
			},
			mockStatusCode: http.StatusCreated,
			wantErr:        false,
			validateRequest: func(t *testing.T, payload map[string]interface{}) {
				if payload["gateway"] != "10.0.1.1" {
					t.Fatalf("expected gateway '10.0.1.1', got %v", payload["gateway"])
				}
				if payload["router_id"] != "router-123" {
					t.Fatalf("expected router_id 'router-123', got %v", payload["router_id"])
				}
			},
		},
		{
			name: "create without description",
			request: &networks.NetworkCreateRequest{
				Name: "simple-network",
				CIDR: "10.0.2.0/24",
			},
			mockResponse: &networks.Network{
				ID:   "net-456",
				Name: "simple-network",
				CIDR: "10.0.2.0/24",
			},
			mockStatusCode: http.StatusCreated,
			wantErr:        false,
			validateRequest: func(t *testing.T, payload map[string]interface{}) {
				if _, ok := payload["gateway"]; ok {
					t.Fatalf("unexpected gateway present in payload: %v", payload["gateway"])
				}
				if _, ok := payload["router_id"]; ok {
					t.Fatalf("unexpected router_id present in payload: %v", payload["router_id"])
				}
			},
		},
		{
			name:           "server error",
			request:        &networks.NetworkCreateRequest{},
			mockStatusCode: http.StatusBadRequest,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost {
					t.Errorf("expected POST, got %s", r.Method)
				}

				expectedPath := "/api/v1/project/proj-123/networks"
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				if tt.validateRequest != nil {
					body, err := io.ReadAll(r.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}
					_ = r.Body.Close()

					var payload map[string]interface{}
					if err := json.Unmarshal(body, &payload); err != nil {
						t.Fatalf("failed to decode request body: %v", err)
					}
					tt.validateRequest(t, payload)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockResponse != nil {
					_ = json.NewEncoder(w).Encode(tt.mockResponse)
				}
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
			client := NewClient(baseClient, "proj-123")

			result, err := client.Create(context.Background(), tt.request)

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

			if result.ID != tt.mockResponse.ID {
				t.Errorf("expected ID %s, got %s", tt.mockResponse.ID, result.ID)
			}
		})
	}
}

// TestClient_Get tests the Get method
func TestClient_Get(t *testing.T) {
	fixture := loadClientNetworkFixture(t)

	tests := []struct {
		name           string
		networkID      string
		mockBody       []byte
		mockStatusCode int
		wantErr        bool
		validate       func(*testing.T, *NetworkResource)
	}{
		{
			name:           "get network successfully",
			networkID:      "net-full-001",
			mockBody:       fixture,
			mockStatusCode: http.StatusOK,
			validate: func(t *testing.T, resource *NetworkResource) {
				if resource == nil {
					t.Fatal("expected non-nil resource")
				}

				assertStringField(t, resource.Network, "ID", "net-full-001")
				assertStringField(t, resource.Network, "Gateway", "10.42.0.1")
				assertStringSliceField(t, resource.Network, "Nameservers", []string{"1.1.1.1", "8.8.8.8"})
				assertStringField(t, resource.Network, "Status", "ACTIVE")
				assertStringField(t, resource.Network, "StatusReason", "OK")

				project := requirePointerStructField(t, resource.Network, "Project")
				assertStringField(t, project.Interface(), "ID", "proj-001")
				assertStringField(t, project.Interface(), "Name", "Tenant Alpha")

				user := requirePointerStructField(t, resource.Network, "User")
				assertStringField(t, user.Interface(), "ID", "user-123")
				assertStringField(t, user.Interface(), "Name", "Alice Ops")

				router := requirePointerStructField(t, resource.Network, "Router")
				assertStringField(t, router.Interface(), "ID", "router-123")
				assertStringField(t, router.Interface(), "Status", "ACTIVE")

				if resource.Ports() == nil {
					t.Error("expected ports operations to be available")
				}
			},
		},
		{
			name:           "network not found",
			networkID:      "net-missing",
			mockStatusCode: http.StatusNotFound,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				expectedPath := "/api/v1/project/proj-123/networks/" + tt.networkID
				if r.URL.Path != expectedPath {
					t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
				}

				w.WriteHeader(tt.mockStatusCode)
				if tt.mockBody != nil {
					w.Header().Set("Content-Type", "application/json")
					_, _ = w.Write(tt.mockBody)
				}
			}))
			defer server.Close()

			baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
			client := NewClient(baseClient, "proj-123")

			result, err := client.Get(context.Background(), tt.networkID)

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

			if tt.validate != nil {
				tt.validate(t, result)
			}
		})
	}
}

// TestClient_Update tests the Update method
func TestClient_Update(t *testing.T) {
	networkID := "net-123"
	request := &networks.NetworkUpdateRequest{
		Name:        "updated-network",
		Description: "Updated description",
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/networks/" + networkID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&networks.Network{
			ID:          networkID,
			Name:        "updated-network",
			Description: "Updated description",
		})
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	result, err := client.Update(context.Background(), networkID, request)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
		return
	}

	if result.Name != "updated-network" {
		t.Errorf("expected name 'updated-network', got %s", result.Name)
	}
}

// TestClient_Delete tests the Delete method
func TestClient_Delete(t *testing.T) {
	networkID := "net-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("expected DELETE, got %s", r.Method)
		}

		expectedPath := "/api/v1/project/proj-123/networks/" + networkID
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	baseClient := internalhttp.NewClient(server.URL, "test-token", &http.Client{Timeout: 5 * time.Second}, nil)
	client := NewClient(baseClient, "proj-123")

	err := client.Delete(context.Background(), networkID)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestNetworkResource_Ports tests that NetworkResource provides port operations
func TestNetworkResource_Ports(t *testing.T) {
	network := &networks.Network{
		ID:   "net-123",
		Name: "test-network",
		CIDR: "10.0.1.0/24",
	}

	baseClient := internalhttp.NewClient("https://api.example.com", "test-token", &http.Client{}, nil)

	resource := &NetworkResource{
		Network: network,
		portOps: &PortsClient{
			baseClient: baseClient,
			projectID:  "proj-123",
			networkID:  "net-123",
		},
	}

	if resource.Ports() == nil {
		t.Error("expected ports operations, got nil")
	}

	if resource.Network.ID != "net-123" {
		t.Errorf("expected network ID net-123, got %s", resource.Network.ID)
	}
}

func loadClientNetworkFixture(t *testing.T) []byte {
	t.Helper()

	path := filepath.Join("test", "testdata", "network_full.json")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read network fixture: %v", err)
	}
	return data
}

func requireStructValue(t *testing.T, obj interface{}) reflect.Value {
	t.Helper()

	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		if v.IsNil() {
			t.Fatal("nil pointer encountered when struct value expected")
		}
		v = v.Elem()
	}
	if v.Kind() != reflect.Struct {
		t.Fatalf("expected struct value, got %s", v.Kind())
	}
	return v
}

func requireStructField(t *testing.T, obj interface{}, fieldName string) reflect.Value {
	t.Helper()

	v := requireStructValue(t, obj)
	field := v.FieldByName(fieldName)
	if !field.IsValid() {
		t.Fatalf("missing field %s", fieldName)
	}
	return field
}

func requirePointerStructField(t *testing.T, obj interface{}, fieldName string) reflect.Value {
	t.Helper()

	field := requireStructField(t, obj, fieldName)
	if field.Kind() != reflect.Ptr {
		t.Fatalf("expected pointer field %s, got %s", fieldName, field.Kind())
	}
	if field.IsNil() {
		t.Fatalf("expected field %s to be non-nil", fieldName)
	}

	elem := field.Elem()
	if elem.Kind() != reflect.Struct {
		t.Fatalf("expected struct pointer for field %s, got %s", fieldName, elem.Kind())
	}
	return elem
}

func assertStringField(t *testing.T, obj interface{}, fieldName, expected string) {
	t.Helper()

	field := requireStructField(t, obj, fieldName)
	if field.Kind() != reflect.String {
		t.Fatalf("expected string field %s, got %s", fieldName, field.Kind())
	}
	if field.String() != expected {
		t.Fatalf("field %s mismatch: expected %s, got %s", fieldName, expected, field.String())
	}
}

func assertStringSliceField(t *testing.T, obj interface{}, fieldName string, expected []string) {
	t.Helper()

	field := requireStructField(t, obj, fieldName)
	if field.Kind() != reflect.Slice {
		t.Fatalf("expected slice field %s, got %s", fieldName, field.Kind())
	}

	actual := make([]string, field.Len())
	for i := 0; i < field.Len(); i++ {
		elem := field.Index(i)
		if elem.Kind() != reflect.String {
			t.Fatalf("expected string slice element for field %s, got %s", fieldName, elem.Kind())
		}
		actual[i] = elem.String()
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Fatalf("field %s mismatch: expected %v, got %v", fieldName, expected, actual)
	}
}
