package quotas

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/vps/quotas"
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

// TestClient_Get tests successful quota retrieval
func TestClient_Get(t *testing.T) {
	mockQuota := &quotas.Quota{
		VM: quotas.QuotaDetail{
			Limit: 10,
			Usage: 3,
		},
		VCPU: quotas.QuotaDetail{
			Limit: 40,
			Usage: 12,
		},
		RAM: quotas.QuotaDetail{
			Limit: 102400,
			Usage: 24576,
		},
		GPU: quotas.QuotaDetail{
			Limit: 2,
			Usage: 0,
		},
		BlockSize: quotas.QuotaDetail{
			Limit: 1000,
			Usage: 250,
		},
		Network: quotas.QuotaDetail{
			Limit: 5,
			Usage: 2,
		},
		Router: quotas.QuotaDetail{
			Limit: 3,
			Usage: 1,
		},
		FloatingIP: quotas.QuotaDetail{
			Limit: 10,
			Usage: 4,
		},
		Share: quotas.QuotaDetail{
			Limit: 5,
			Usage: 1,
		},
		ShareSize: quotas.QuotaDetail{
			Limit: 500,
			Usage: 100,
		},
	}

	expectedPath := "/api/v1/project/proj-123/quotas"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != expectedPath {
			t.Errorf("expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			t.Error("expected Authorization header, got empty")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockQuota)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	quota, err := client.Get(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quota == nil {
		t.Fatal("expected quota, got nil")
	}

	// Verify all quota fields
	if quota.VM.Limit != mockQuota.VM.Limit {
		t.Errorf("expected VM limit %d, got %d", mockQuota.VM.Limit, quota.VM.Limit)
	}
	if quota.VM.Usage != mockQuota.VM.Usage {
		t.Errorf("expected VM usage %d, got %d", mockQuota.VM.Usage, quota.VM.Usage)
	}
	if quota.VCPU.Limit != mockQuota.VCPU.Limit {
		t.Errorf("expected VCPU limit %d, got %d", mockQuota.VCPU.Limit, quota.VCPU.Limit)
	}
	if quota.VCPU.Usage != mockQuota.VCPU.Usage {
		t.Errorf("expected VCPU usage %d, got %d", mockQuota.VCPU.Usage, quota.VCPU.Usage)
	}
	if quota.RAM.Limit != mockQuota.RAM.Limit {
		t.Errorf("expected RAM limit %d, got %d", mockQuota.RAM.Limit, quota.RAM.Limit)
	}
	if quota.RAM.Usage != mockQuota.RAM.Usage {
		t.Errorf("expected RAM usage %d, got %d", mockQuota.RAM.Usage, quota.RAM.Usage)
	}
	if quota.GPU.Limit != mockQuota.GPU.Limit {
		t.Errorf("expected GPU limit %d, got %d", mockQuota.GPU.Limit, quota.GPU.Limit)
	}
	if quota.GPU.Usage != mockQuota.GPU.Usage {
		t.Errorf("expected GPU usage %d, got %d", mockQuota.GPU.Usage, quota.GPU.Usage)
	}
	if quota.BlockSize.Limit != mockQuota.BlockSize.Limit {
		t.Errorf("expected BlockSize limit %d, got %d", mockQuota.BlockSize.Limit, quota.BlockSize.Limit)
	}
	if quota.BlockSize.Usage != mockQuota.BlockSize.Usage {
		t.Errorf("expected BlockSize usage %d, got %d", mockQuota.BlockSize.Usage, quota.BlockSize.Usage)
	}
	if quota.Network.Limit != mockQuota.Network.Limit {
		t.Errorf("expected Network limit %d, got %d", mockQuota.Network.Limit, quota.Network.Limit)
	}
	if quota.Network.Usage != mockQuota.Network.Usage {
		t.Errorf("expected Network usage %d, got %d", mockQuota.Network.Usage, quota.Network.Usage)
	}
	if quota.Router.Limit != mockQuota.Router.Limit {
		t.Errorf("expected Router limit %d, got %d", mockQuota.Router.Limit, quota.Router.Limit)
	}
	if quota.Router.Usage != mockQuota.Router.Usage {
		t.Errorf("expected Router usage %d, got %d", mockQuota.Router.Usage, quota.Router.Usage)
	}
	if quota.FloatingIP.Limit != mockQuota.FloatingIP.Limit {
		t.Errorf("expected FloatingIP limit %d, got %d", mockQuota.FloatingIP.Limit, quota.FloatingIP.Limit)
	}
	if quota.FloatingIP.Usage != mockQuota.FloatingIP.Usage {
		t.Errorf("expected FloatingIP usage %d, got %d", mockQuota.FloatingIP.Usage, quota.FloatingIP.Usage)
	}
}

// TestClient_Get_UnlimitedQuotas tests handling of unlimited quotas
func TestClient_Get_UnlimitedQuotas(t *testing.T) {
	mockQuota := &quotas.Quota{
		VM: quotas.QuotaDetail{
			Limit: -1, // unlimited
			Usage: 5,
		},
		VCPU: quotas.QuotaDetail{
			Limit: -1,
			Usage: 20,
		},
		RAM: quotas.QuotaDetail{
			Limit: -1,
			Usage: 40960,
		},
		GPU: quotas.QuotaDetail{
			Limit: 0,
			Usage: 0,
		},
		BlockSize: quotas.QuotaDetail{
			Limit: -1,
			Usage: 500,
		},
		Network: quotas.QuotaDetail{
			Limit: -1,
			Usage: 3,
		},
		Router: quotas.QuotaDetail{
			Limit: -1,
			Usage: 2,
		},
		FloatingIP: quotas.QuotaDetail{
			Limit: -1,
			Usage: 6,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockQuota)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	quota, err := client.Get(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quota == nil {
		t.Fatal("expected quota, got nil")
	}

	// Verify unlimited quotas (-1)
	if quota.VM.Limit != -1 {
		t.Errorf("expected VM limit -1 (unlimited), got %d", quota.VM.Limit)
	}
	if quota.VCPU.Limit != -1 {
		t.Errorf("expected VCPU limit -1 (unlimited), got %d", quota.VCPU.Limit)
	}
	if quota.RAM.Limit != -1 {
		t.Errorf("expected RAM limit -1 (unlimited), got %d", quota.RAM.Limit)
	}

	// Verify usage is still tracked
	if quota.VM.Usage != 5 {
		t.Errorf("expected VM usage 5, got %d", quota.VM.Usage)
	}
	if quota.VCPU.Usage != 20 {
		t.Errorf("expected VCPU usage 20, got %d", quota.VCPU.Usage)
	}
}

// TestClient_Get_ZeroUsage tests handling of quotas with zero usage
func TestClient_Get_ZeroUsage(t *testing.T) {
	mockQuota := &quotas.Quota{
		VM:         quotas.QuotaDetail{Limit: 10, Usage: 0},
		VCPU:       quotas.QuotaDetail{Limit: 40, Usage: 0},
		RAM:        quotas.QuotaDetail{Limit: 102400, Usage: 0},
		GPU:        quotas.QuotaDetail{Limit: 2, Usage: 0},
		BlockSize:  quotas.QuotaDetail{Limit: 1000, Usage: 0},
		Network:    quotas.QuotaDetail{Limit: 5, Usage: 0},
		Router:     quotas.QuotaDetail{Limit: 3, Usage: 0},
		FloatingIP: quotas.QuotaDetail{Limit: 10, Usage: 0},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockQuota)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	quota, err := client.Get(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quota == nil {
		t.Fatal("expected quota, got nil")
	}

	// Verify zero usage
	if quota.VM.Usage != 0 {
		t.Errorf("expected VM usage 0, got %d", quota.VM.Usage)
	}
	if quota.VCPU.Usage != 0 {
		t.Errorf("expected VCPU usage 0, got %d", quota.VCPU.Usage)
	}
	if quota.RAM.Usage != 0 {
		t.Errorf("expected RAM usage 0, got %d", quota.RAM.Usage)
	}
}

// TestClient_Get_AtLimit tests handling of quotas at their limit
func TestClient_Get_AtLimit(t *testing.T) {
	mockQuota := &quotas.Quota{
		VM:         quotas.QuotaDetail{Limit: 10, Usage: 10},
		VCPU:       quotas.QuotaDetail{Limit: 40, Usage: 40},
		RAM:        quotas.QuotaDetail{Limit: 102400, Usage: 102400},
		GPU:        quotas.QuotaDetail{Limit: 2, Usage: 2},
		BlockSize:  quotas.QuotaDetail{Limit: 1000, Usage: 1000},
		Network:    quotas.QuotaDetail{Limit: 5, Usage: 5},
		Router:     quotas.QuotaDetail{Limit: 3, Usage: 3},
		FloatingIP: quotas.QuotaDetail{Limit: 10, Usage: 10},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockQuota)
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	ctx := context.Background()
	quota, err := client.Get(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quota == nil {
		t.Fatal("expected quota, got nil")
	}

	// Verify quotas at limit
	if quota.VM.Limit != quota.VM.Usage {
		t.Errorf("expected VM at limit, got limit=%d usage=%d", quota.VM.Limit, quota.VM.Usage)
	}
	if quota.VCPU.Limit != quota.VCPU.Usage {
		t.Errorf("expected VCPU at limit, got limit=%d usage=%d", quota.VCPU.Limit, quota.VCPU.Usage)
	}
	if quota.RAM.Limit != quota.RAM.Usage {
		t.Errorf("expected RAM at limit, got limit=%d usage=%d", quota.RAM.Limit, quota.RAM.Usage)
	}
}

// TestClient_Get_Errors tests error handling for quota retrieval
func TestClient_Get_Errors(t *testing.T) {
	tests := []struct {
		name         string
		mockStatus   int
		mockResponse map[string]interface{}
		expectError  bool
	}{
		{
			name:       "unauthorized - 401",
			mockStatus: http.StatusUnauthorized,
			mockResponse: map[string]interface{}{
				"error":   "Unauthorized",
				"message": "Invalid or expired token",
			},
			expectError: true,
		},
		{
			name:       "forbidden - 403",
			mockStatus: http.StatusForbidden,
			mockResponse: map[string]interface{}{
				"error":   "Forbidden",
				"message": "Insufficient permissions",
			},
			expectError: true,
		},
		{
			name:       "not found - 404",
			mockStatus: http.StatusNotFound,
			mockResponse: map[string]interface{}{
				"error":   "Not Found",
				"message": "Project not found",
			},
			expectError: true,
		},
		{
			name:       "internal server error - 500",
			mockStatus: http.StatusInternalServerError,
			mockResponse: map[string]interface{}{
				"error":   "Internal Server Error",
				"message": "An unexpected error occurred",
			},
			expectError: true,
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

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			quota, err := client.Get(ctx)

			if tt.expectError && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if tt.expectError && quota != nil {
				t.Errorf("expected nil quota on error, got %v", quota)
			}
		})
	}
}

// TestClient_Get_ContextCancellation tests context cancellation handling
func TestClient_Get_ContextCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Simulate slow response
		time.Sleep(100 * time.Millisecond)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(&quotas.Quota{})
	}))
	defer server.Close()

	httpClient := &http.Client{Timeout: 5 * time.Second}
	baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
	client := NewClient(baseClient, "proj-123")

	// Create a context that will be cancelled immediately
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	quota, err := client.Get(ctx)

	if err == nil {
		t.Error("expected error due to context cancellation, got nil")
	}
	if quota != nil {
		t.Errorf("expected nil quota on error, got %v", quota)
	}
}

// TestClient_Get_OptionalFields tests handling of optional fields in quota response
func TestClient_Get_OptionalFields(t *testing.T) {
	tests := []struct {
		name      string
		mockQuota *quotas.Quota
	}{
		{
			name: "with share fields",
			mockQuota: &quotas.Quota{
				VM:         quotas.QuotaDetail{Limit: 10, Usage: 3},
				VCPU:       quotas.QuotaDetail{Limit: 40, Usage: 12},
				RAM:        quotas.QuotaDetail{Limit: 102400, Usage: 24576},
				GPU:        quotas.QuotaDetail{Limit: 2, Usage: 0},
				BlockSize:  quotas.QuotaDetail{Limit: 1000, Usage: 250},
				Network:    quotas.QuotaDetail{Limit: 5, Usage: 2},
				Router:     quotas.QuotaDetail{Limit: 3, Usage: 1},
				FloatingIP: quotas.QuotaDetail{Limit: 10, Usage: 4},
				Share:      quotas.QuotaDetail{Limit: 5, Usage: 1},
				ShareSize:  quotas.QuotaDetail{Limit: 500, Usage: 100},
			},
		},
		{
			name: "without share fields",
			mockQuota: &quotas.Quota{
				VM:         quotas.QuotaDetail{Limit: 10, Usage: 3},
				VCPU:       quotas.QuotaDetail{Limit: 40, Usage: 12},
				RAM:        quotas.QuotaDetail{Limit: 102400, Usage: 24576},
				GPU:        quotas.QuotaDetail{Limit: 2, Usage: 0},
				BlockSize:  quotas.QuotaDetail{Limit: 1000, Usage: 250},
				Network:    quotas.QuotaDetail{Limit: 5, Usage: 2},
				Router:     quotas.QuotaDetail{Limit: 3, Usage: 1},
				FloatingIP: quotas.QuotaDetail{Limit: 10, Usage: 4},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_ = json.NewEncoder(w).Encode(tt.mockQuota)
			}))
			defer server.Close()

			httpClient := &http.Client{Timeout: 5 * time.Second}
			baseClient := internalhttp.NewClient(server.URL, "test-token", httpClient, nil)
			client := NewClient(baseClient, "proj-123")

			ctx := context.Background()
			quota, err := client.Get(ctx)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if quota == nil {
				t.Fatal("expected quota, got nil")
			}

			// Verify required fields
			if quota.VM.Limit != tt.mockQuota.VM.Limit {
				t.Errorf("expected VM limit %d, got %d", tt.mockQuota.VM.Limit, quota.VM.Limit)
			}

			// Verify optional fields if present
			if tt.mockQuota.Share.Limit != 0 {
				if quota.Share.Limit != tt.mockQuota.Share.Limit {
					t.Errorf("expected Share limit %d, got %d", tt.mockQuota.Share.Limit, quota.Share.Limit)
				}
			}
			if tt.mockQuota.ShareSize.Limit != 0 {
				if quota.ShareSize.Limit != tt.mockQuota.ShareSize.Limit {
					t.Errorf("expected ShareSize limit %d, got %d", tt.mockQuota.ShareSize.Limit, quota.ShareSize.Limit)
				}
			}
		})
	}
}
