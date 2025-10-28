package quotas_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	cloudsdk "github.com/Zillaforge/cloud-sdk"
	"github.com/Zillaforge/cloud-sdk/models/vps/quotas"
)

// TestQuotasGet_Success verifies successful quota retrieval
func TestQuotasGet_Success(t *testing.T) {
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

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET request, got %s", r.Method)
		}
		if r.URL.Path != "/vps/api/v1/project/proj-123/quotas" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockQuota)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	quotasClient := vpsClient.Quotas()

	ctx := context.Background()
	quota, err := quotasClient.Get(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quota == nil {
		t.Fatal("expected quota, got nil")
	}

	// Verify VM quota
	if quota.VM.Limit != mockQuota.VM.Limit {
		t.Errorf("expected VM limit %d, got %d", mockQuota.VM.Limit, quota.VM.Limit)
	}
	if quota.VM.Usage != mockQuota.VM.Usage {
		t.Errorf("expected VM usage %d, got %d", mockQuota.VM.Usage, quota.VM.Usage)
	}

	// Verify VCPU quota
	if quota.VCPU.Limit != mockQuota.VCPU.Limit {
		t.Errorf("expected VCPU limit %d, got %d", mockQuota.VCPU.Limit, quota.VCPU.Limit)
	}
	if quota.VCPU.Usage != mockQuota.VCPU.Usage {
		t.Errorf("expected VCPU usage %d, got %d", mockQuota.VCPU.Usage, quota.VCPU.Usage)
	}

	// Verify RAM quota
	if quota.RAM.Limit != mockQuota.RAM.Limit {
		t.Errorf("expected RAM limit %d, got %d", mockQuota.RAM.Limit, quota.RAM.Limit)
	}
	if quota.RAM.Usage != mockQuota.RAM.Usage {
		t.Errorf("expected RAM usage %d, got %d", mockQuota.RAM.Usage, quota.RAM.Usage)
	}

	// Verify GPU quota
	if quota.GPU.Limit != mockQuota.GPU.Limit {
		t.Errorf("expected GPU limit %d, got %d", mockQuota.GPU.Limit, quota.GPU.Limit)
	}
	if quota.GPU.Usage != mockQuota.GPU.Usage {
		t.Errorf("expected GPU usage %d, got %d", mockQuota.GPU.Usage, quota.GPU.Usage)
	}

	// Verify BlockSize quota
	if quota.BlockSize.Limit != mockQuota.BlockSize.Limit {
		t.Errorf("expected BlockSize limit %d, got %d", mockQuota.BlockSize.Limit, quota.BlockSize.Limit)
	}
	if quota.BlockSize.Usage != mockQuota.BlockSize.Usage {
		t.Errorf("expected BlockSize usage %d, got %d", mockQuota.BlockSize.Usage, quota.BlockSize.Usage)
	}

	// Verify Network quota
	if quota.Network.Limit != mockQuota.Network.Limit {
		t.Errorf("expected Network limit %d, got %d", mockQuota.Network.Limit, quota.Network.Limit)
	}
	if quota.Network.Usage != mockQuota.Network.Usage {
		t.Errorf("expected Network usage %d, got %d", mockQuota.Network.Usage, quota.Network.Usage)
	}

	// Verify Router quota
	if quota.Router.Limit != mockQuota.Router.Limit {
		t.Errorf("expected Router limit %d, got %d", mockQuota.Router.Limit, quota.Router.Limit)
	}
	if quota.Router.Usage != mockQuota.Router.Usage {
		t.Errorf("expected Router usage %d, got %d", mockQuota.Router.Usage, quota.Router.Usage)
	}

	// Verify FloatingIP quota
	if quota.FloatingIP.Limit != mockQuota.FloatingIP.Limit {
		t.Errorf("expected FloatingIP limit %d, got %d", mockQuota.FloatingIP.Limit, quota.FloatingIP.Limit)
	}
	if quota.FloatingIP.Usage != mockQuota.FloatingIP.Usage {
		t.Errorf("expected FloatingIP usage %d, got %d", mockQuota.FloatingIP.Usage, quota.FloatingIP.Usage)
	}

	// Verify Share quota (optional field)
	if quota.Share.Limit != mockQuota.Share.Limit {
		t.Errorf("expected Share limit %d, got %d", mockQuota.Share.Limit, quota.Share.Limit)
	}
	if quota.Share.Usage != mockQuota.Share.Usage {
		t.Errorf("expected Share usage %d, got %d", mockQuota.Share.Usage, quota.Share.Usage)
	}

	// Verify ShareSize quota (optional field)
	if quota.ShareSize.Limit != mockQuota.ShareSize.Limit {
		t.Errorf("expected ShareSize limit %d, got %d", mockQuota.ShareSize.Limit, quota.ShareSize.Limit)
	}
	if quota.ShareSize.Usage != mockQuota.ShareSize.Usage {
		t.Errorf("expected ShareSize usage %d, got %d", mockQuota.ShareSize.Usage, quota.ShareSize.Usage)
	}
}

// TestQuotasGet_UnlimitedQuota verifies handling of unlimited quotas (-1 limit)
func TestQuotasGet_UnlimitedQuota(t *testing.T) {
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

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	quotasClient := vpsClient.Quotas()

	ctx := context.Background()
	quota, err := quotasClient.Get(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quota == nil {
		t.Fatal("expected quota, got nil")
	}

	// Verify unlimited quotas are properly set to -1
	if quota.VM.Limit != -1 {
		t.Errorf("expected VM limit -1 (unlimited), got %d", quota.VM.Limit)
	}
	if quota.VCPU.Limit != -1 {
		t.Errorf("expected VCPU limit -1 (unlimited), got %d", quota.VCPU.Limit)
	}
	if quota.RAM.Limit != -1 {
		t.Errorf("expected RAM limit -1 (unlimited), got %d", quota.RAM.Limit)
	}

	// Verify usage is still tracked even with unlimited quota
	if quota.VM.Usage != 5 {
		t.Errorf("expected VM usage 5, got %d", quota.VM.Usage)
	}
}

// TestQuotasGet_Errors verifies error handling for quota retrieval
func TestQuotasGet_Errors(t *testing.T) {
	tests := []struct {
		name         string
		mockStatus   int
		mockResponse interface{}
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

			client := cloudsdk.NewClient(server.URL, "test-token")
			vpsClient := client.Project("proj-123").VPS()
			quotasClient := vpsClient.Quotas()

			ctx := context.Background()
			quota, err := quotasClient.Get(ctx)

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

// TestQuotasGet_ZeroUsage verifies handling of quotas with zero usage
func TestQuotasGet_ZeroUsage(t *testing.T) {
	mockQuota := &quotas.Quota{
		VM: quotas.QuotaDetail{
			Limit: 10,
			Usage: 0,
		},
		VCPU: quotas.QuotaDetail{
			Limit: 40,
			Usage: 0,
		},
		RAM: quotas.QuotaDetail{
			Limit: 102400,
			Usage: 0,
		},
		GPU: quotas.QuotaDetail{
			Limit: 2,
			Usage: 0,
		},
		BlockSize: quotas.QuotaDetail{
			Limit: 1000,
			Usage: 0,
		},
		Network: quotas.QuotaDetail{
			Limit: 5,
			Usage: 0,
		},
		Router: quotas.QuotaDetail{
			Limit: 3,
			Usage: 0,
		},
		FloatingIP: quotas.QuotaDetail{
			Limit: 10,
			Usage: 0,
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_ = json.NewEncoder(w).Encode(mockQuota)
	}))
	defer server.Close()

	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	quotasClient := vpsClient.Quotas()

	ctx := context.Background()
	quota, err := quotasClient.Get(ctx)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if quota == nil {
		t.Fatal("expected quota, got nil")
	}

	// Verify zero usage is properly handled
	if quota.VM.Usage != 0 {
		t.Errorf("expected VM usage 0, got %d", quota.VM.Usage)
	}
	if quota.VCPU.Usage != 0 {
		t.Errorf("expected VCPU usage 0, got %d", quota.VCPU.Usage)
	}
}
