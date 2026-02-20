package vps

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	floatingipsmodels "github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
	"github.com/Zillaforge/cloud-sdk/modules/vps/floatingips"
)

// TestWaitForFloatingIPStatus_Success verifies waiting for a floating IP to reach target status.
func TestWaitForFloatingIPStatus_Success(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus floatingipsmodels.FloatingIPStatus
		statusFlow   []string // Sequence of statuses returned by mock floating IP
	}{
		{
			name:         "wait for ACTIVE from PENDING",
			targetStatus: floatingipsmodels.FloatingIPStatusActive,
			statusFlow:   []string{"PENDING", "PENDING", "ACTIVE"},
		},
		{
			name:         "wait for DOWN from ACTIVE",
			targetStatus: floatingipsmodels.FloatingIPStatusDown,
			statusFlow:   []string{"ACTIVE", "DOWN"},
		},
		{
			name:         "already at target status",
			targetStatus: floatingipsmodels.FloatingIPStatusActive,
			statusFlow:   []string{"ACTIVE"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			currentStatusIndex := 0

			// Create mock server that cycles through status flow
			mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != "GET" {
					t.Errorf("expected GET, got %s", r.Method)
				}

				status := tt.statusFlow[currentStatusIndex]
				if currentStatusIndex < len(tt.statusFlow)-1 {
					currentStatusIndex++
				}

				response := map[string]interface{}{
					"id":         "fip-test-1",
					"uuid":       "fip-uuid-1",
					"name":       "test-floating-ip",
					"address":    "192.168.1.100",
					"extnet_id":  "ext-net-1",
					"project_id": "proj-1",
					"user_id":    "user-1",
					"status":     status,
					"reserved":   false,
					"createdAt":  "2025-10-28T00:00:00Z",
					"updatedAt":  "2025-10-28T00:00:00Z",
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				if err := json.NewEncoder(w).Encode(response); err != nil {
					t.Fatalf("failed to encode response: %v", err)
				}
			}))
			defer mockServer.Close()

			// Create client
			httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
			floatingIPClient := floatingips.NewClient(httpClient, "proj-1")

			// Wait for status with short intervals for testing
			ctx := context.Background()
			err := WaitForFloatingIPStatus(ctx, FloatingIPWaiterConfig{
				Client:       floatingIPClient,
				FloatingIPID: "fip-test-1",
				TargetStatus: tt.targetStatus,
				WaiterOptions: []waiter.Option{
					waiter.WithInterval(10 * time.Millisecond),
					waiter.WithMaxWait(2 * time.Second),
				},
			})

			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestWaitForFloatingIPStatus_RejectedState verifies handling when floating IP enters REJECTED state.
func TestWaitForFloatingIPStatus_RejectedState(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]interface{}{
			"id":         "fip-test-1",
			"uuid":       "fip-uuid-1",
			"name":       "test-floating-ip",
			"address":    "192.168.1.100",
			"extnet_id":  "ext-net-1",
			"project_id": "proj-1",
			"user_id":    "user-1",
			"status":     "REJECTED",
			"reserved":   false,
			"createdAt":  "2025-10-28T00:00:00Z",
			"updatedAt":  "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	// Create client
	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	floatingIPClient := floatingips.NewClient(httpClient, "proj-1")

	// Wait for ACTIVE status - should fail when REJECTED is encountered
	ctx := context.Background()
	err := WaitForFloatingIPStatus(ctx, FloatingIPWaiterConfig{
		Client:       floatingIPClient,
		FloatingIPID: "fip-test-1",
		TargetStatus: floatingipsmodels.FloatingIPStatusActive,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(1 * time.Second),
		},
	})

	if err == nil {
		t.Error("expected error when floating IP enters REJECTED state")
	}
	if err.Error() != "floating IP entered REJECTED state while waiting for ACTIVE" {
		t.Errorf("unexpected error message: %v", err)
	}
}

// TestWaitForFloatingIPActive convenience function.
func TestWaitForFloatingIPActive(t *testing.T) {
	currentStatusIndex := 0
	statusFlow := []string{"PENDING", "ACTIVE"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		status := statusFlow[currentStatusIndex]
		if currentStatusIndex < len(statusFlow)-1 {
			currentStatusIndex++
		}

		response := map[string]interface{}{
			"id":         "fip-test-1",
			"uuid":       "fip-uuid-1",
			"name":       "test-floating-ip",
			"address":    "192.168.1.100",
			"extnet_id":  "ext-net-1",
			"project_id": "proj-1",
			"user_id":    "user-1",
			"status":     status,
			"reserved":   false,
			"createdAt":  "2025-10-28T00:00:00Z",
			"updatedAt":  "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	// Create client
	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	floatingIPClient := floatingips.NewClient(httpClient, "proj-1")

	// Test convenience function
	ctx := context.Background()
	err := WaitForFloatingIPActive(ctx, floatingIPClient, "fip-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(2*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
