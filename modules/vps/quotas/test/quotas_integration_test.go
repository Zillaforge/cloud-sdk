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

// TestQuotaRetrieval verifies the complete quota retrieval workflow
func TestQuotaRetrieval(t *testing.T) {
	// Mock quota response for a project
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
			Limit: 102400, // 100 GiB in MiB
			Usage: 24576,  // 24 GiB in MiB
		},
		GPU: quotas.QuotaDetail{
			Limit: 2,
			Usage: 0,
		},
		BlockSize: quotas.QuotaDetail{
			Limit: 1000, // 1000 GiB
			Usage: 250,  // 250 GiB used
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
			Limit: 500, // 500 GiB
			Usage: 100, // 100 GiB used
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		if r.Method != http.MethodGet {
			w.WriteHeader(http.StatusMethodNotAllowed)
			_ = json.NewEncoder(w).Encode(map[string]string{"error": "Method not allowed"})
			return
		}

		// Handle get quotas
		if r.URL.Path == "/vps/api/v1/project/proj-123/quotas" {
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(mockQuota)
			return
		}

		w.WriteHeader(http.StatusNotFound)
		_ = json.NewEncoder(w).Encode(map[string]string{"error": "Not found"})
	}))
	defer server.Close()

	// Initialize SDK client
	client := cloudsdk.NewClient(server.URL, "test-token")
	vpsClient := client.Project("proj-123").VPS()
	quotasClient := vpsClient.Quotas()

	ctx := context.Background()

	// Test Case 1: Retrieve project quotas
	t.Run("retrieve project quotas", func(t *testing.T) {
		quota, err := quotasClient.Get(ctx)
		if err != nil {
			t.Fatalf("failed to get quotas: %v", err)
		}

		// Verify VM quota
		if quota.VM.Limit != 10 {
			t.Errorf("expected VM limit 10, got %d", quota.VM.Limit)
		}
		if quota.VM.Usage != 3 {
			t.Errorf("expected VM usage 3, got %d", quota.VM.Usage)
		}

		// Verify VCPU quota
		if quota.VCPU.Limit != 40 {
			t.Errorf("expected VCPU limit 40, got %d", quota.VCPU.Limit)
		}
		if quota.VCPU.Usage != 12 {
			t.Errorf("expected VCPU usage 12, got %d", quota.VCPU.Usage)
		}

		// Verify RAM quota
		if quota.RAM.Limit != 102400 {
			t.Errorf("expected RAM limit 102400, got %d", quota.RAM.Limit)
		}
		if quota.RAM.Usage != 24576 {
			t.Errorf("expected RAM usage 24576, got %d", quota.RAM.Usage)
		}
	})

	// Test Case 2: Calculate remaining capacity
	t.Run("calculate remaining capacity", func(t *testing.T) {
		quota, err := quotasClient.Get(ctx)
		if err != nil {
			t.Fatalf("failed to get quotas: %v", err)
		}

		// Verify we can calculate remaining resources
		vmRemaining := quota.VM.Limit - quota.VM.Usage
		if vmRemaining != 7 {
			t.Errorf("expected 7 VMs remaining, got %d", vmRemaining)
		}

		vcpuRemaining := quota.VCPU.Limit - quota.VCPU.Usage
		if vcpuRemaining != 28 {
			t.Errorf("expected 28 VCPUs remaining, got %d", vcpuRemaining)
		}

		ramRemaining := quota.RAM.Limit - quota.RAM.Usage
		if ramRemaining != 77824 {
			t.Errorf("expected 77824 MiB RAM remaining, got %d", ramRemaining)
		}

		floatingIPRemaining := quota.FloatingIP.Limit - quota.FloatingIP.Usage
		if floatingIPRemaining != 6 {
			t.Errorf("expected 6 Floating IPs remaining, got %d", floatingIPRemaining)
		}
	})

	// Test Case 3: Handle unlimited quotas
	t.Run("handle unlimited quotas", func(t *testing.T) {
		unlimitedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.Header().Set("Content-Type", "application/json")
			unlimitedQuota := &quotas.Quota{
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
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(unlimitedQuota)
		}))
		defer unlimitedServer.Close()

		unlimitedClient := cloudsdk.NewClient(unlimitedServer.URL, "test-token")
		unlimitedVPS := unlimitedClient.Project("proj-unlimited").VPS()
		unlimitedQuotasClient := unlimitedVPS.Quotas()

		quota, err := unlimitedQuotasClient.Get(ctx)
		if err != nil {
			t.Fatalf("failed to get unlimited quotas: %v", err)
		}

		// Verify unlimited quotas (-1)
		if quota.VM.Limit != -1 {
			t.Errorf("expected unlimited VM quota (-1), got %d", quota.VM.Limit)
		}
		if quota.VCPU.Limit != -1 {
			t.Errorf("expected unlimited VCPU quota (-1), got %d", quota.VCPU.Limit)
		}

		// Verify usage is still tracked
		if quota.VM.Usage != 5 {
			t.Errorf("expected VM usage 5, got %d", quota.VM.Usage)
		}
		if quota.VCPU.Usage != 20 {
			t.Errorf("expected VCPU usage 20, got %d", quota.VCPU.Usage)
		}
	})
}

// TestQuotaRetrieval_EdgeCases verifies edge cases in quota retrieval
func TestQuotaRetrieval_EdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		setupServer  func() *httptest.Server
		expectError  bool
		validateFunc func(*testing.T, *quotas.Quota)
	}{
		{
			name: "zero usage across all resources",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&quotas.Quota{
						VM:         quotas.QuotaDetail{Limit: 10, Usage: 0},
						VCPU:       quotas.QuotaDetail{Limit: 40, Usage: 0},
						RAM:        quotas.QuotaDetail{Limit: 102400, Usage: 0},
						GPU:        quotas.QuotaDetail{Limit: 2, Usage: 0},
						BlockSize:  quotas.QuotaDetail{Limit: 1000, Usage: 0},
						Network:    quotas.QuotaDetail{Limit: 5, Usage: 0},
						Router:     quotas.QuotaDetail{Limit: 3, Usage: 0},
						FloatingIP: quotas.QuotaDetail{Limit: 10, Usage: 0},
					})
				}))
			},
			expectError: false,
			validateFunc: func(t *testing.T, quota *quotas.Quota) {
				if quota.VM.Usage != 0 || quota.VCPU.Usage != 0 {
					t.Errorf("expected zero usage, got VM=%d VCPU=%d", quota.VM.Usage, quota.VCPU.Usage)
				}
			},
		},
		{
			name: "quota at limit",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusOK)
					_ = json.NewEncoder(w).Encode(&quotas.Quota{
						VM:         quotas.QuotaDetail{Limit: 10, Usage: 10},
						VCPU:       quotas.QuotaDetail{Limit: 40, Usage: 40},
						RAM:        quotas.QuotaDetail{Limit: 102400, Usage: 102400},
						GPU:        quotas.QuotaDetail{Limit: 2, Usage: 2},
						BlockSize:  quotas.QuotaDetail{Limit: 1000, Usage: 1000},
						Network:    quotas.QuotaDetail{Limit: 5, Usage: 5},
						Router:     quotas.QuotaDetail{Limit: 3, Usage: 3},
						FloatingIP: quotas.QuotaDetail{Limit: 10, Usage: 10},
					})
				}))
			},
			expectError: false,
			validateFunc: func(t *testing.T, quota *quotas.Quota) {
				if quota.VM.Limit != quota.VM.Usage {
					t.Errorf("expected VM at limit, got limit=%d usage=%d", quota.VM.Limit, quota.VM.Usage)
				}
				if quota.VCPU.Limit != quota.VCPU.Usage {
					t.Errorf("expected VCPU at limit, got limit=%d usage=%d", quota.VCPU.Limit, quota.VCPU.Usage)
				}
			},
		},
		{
			name: "unauthorized request",
			setupServer: func() *httptest.Server {
				return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusUnauthorized)
					_ = json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
				}))
			},
			expectError: true,
			validateFunc: func(t *testing.T, quota *quotas.Quota) {
				if quota != nil {
					t.Errorf("expected nil quota on error, got %v", quota)
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := tt.setupServer()
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

			if tt.validateFunc != nil {
				tt.validateFunc(t, quota)
			}
		})
	}
}
