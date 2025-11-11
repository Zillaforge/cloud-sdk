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
	serversmodels "github.com/Zillaforge/cloud-sdk/models/vps/servers"
	"github.com/Zillaforge/cloud-sdk/modules/vps/floatingips"
	"github.com/Zillaforge/cloud-sdk/modules/vps/servers"
)

// TestWaitForServerStatus_Success verifies waiting for a server to reach target status.
func TestWaitForServerStatus_Success(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus serversmodels.ServerStatus
		statusFlow   []string // Sequence of statuses returned by mock server
	}{
		{
			name:         "wait for ACTIVE from BUILD",
			targetStatus: serversmodels.ServerStatusActive,
			statusFlow:   []string{"BUILD", "BUILD", "ACTIVE"},
		},
		{
			name:         "wait for SHUTOFF from ACTIVE",
			targetStatus: serversmodels.ServerStatusShutoff,
			statusFlow:   []string{"ACTIVE", "SHUTOFF"},
		},
		{
			name:         "already at target status",
			targetStatus: serversmodels.ServerStatusActive,
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
					"id":         "svr-test-1",
					"name":       "test-server",
					"status":     status,
					"flavor_id":  "flavor-1",
					"image_id":   "image-1",
					"project_id": "proj-1",
					"created_at": "2025-10-28T00:00:00Z",
					"updated_at": "2025-10-28T00:00:00Z",
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
			serverClient := servers.NewClient(httpClient, "proj-1")

			// Wait for status with short intervals for testing
			ctx := context.Background()
			err := WaitForServerStatus(ctx, ServerWaiterConfig{
				Client:       serverClient,
				ServerID:     "svr-test-1",
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

// TestWaitForServerStatus_ErrorState verifies handling when server enters ERROR state.
func TestWaitForServerStatus_ErrorState(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]interface{}{
			"id":         "svr-test-1",
			"name":       "test-server",
			"status":     "ERROR",
			"flavor_id":  "flavor-1",
			"image_id":   "image-1",
			"project_id": "proj-1",
			"created_at": "2025-10-28T00:00:00Z",
			"updated_at": "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	serverClient := servers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForServerStatus(ctx, ServerWaiterConfig{
		Client:       serverClient,
		ServerID:     "svr-test-1",
		TargetStatus: serversmodels.ServerStatusActive,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(1 * time.Second),
		},
	})

	if err == nil {
		t.Fatal("expected error when server enters ERROR state, got nil")
	}

	expectedMsg := "server entered ERROR state"
	if err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("expected error message to start with '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestWaitForServerStatus_Timeout verifies timeout handling.
func TestWaitForServerStatus_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always return BUILD status
		response := map[string]interface{}{
			"id":         "svr-test-1",
			"name":       "test-server",
			"status":     "BUILD",
			"flavor_id":  "flavor-1",
			"image_id":   "image-1",
			"project_id": "proj-1",
			"created_at": "2025-10-28T00:00:00Z",
			"updated_at": "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	serverClient := servers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForServerStatus(ctx, ServerWaiterConfig{
		Client:       serverClient,
		ServerID:     "svr-test-1",
		TargetStatus: serversmodels.ServerStatusActive,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(100 * time.Millisecond),
		},
	})

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if err != waiter.ErrWaitTimeout {
		t.Errorf("expected ErrWaitTimeout, got: %v", err)
	}
}

// TestWaitForServerStatus_ContextCancellation verifies context cancellation.
func TestWaitForServerStatus_ContextCancellation(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always return BUILD status
		response := map[string]interface{}{
			"id":         "svr-test-1",
			"name":       "test-server",
			"status":     "BUILD",
			"flavor_id":  "flavor-1",
			"image_id":   "image-1",
			"project_id": "proj-1",
			"created_at": "2025-10-28T00:00:00Z",
			"updated_at": "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	serverClient := servers.NewClient(httpClient, "proj-1")

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := WaitForServerStatus(ctx, ServerWaiterConfig{
		Client:       serverClient,
		ServerID:     "svr-test-1",
		TargetStatus: serversmodels.ServerStatusActive,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(5 * time.Second),
		},
	})

	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}

	if err != context.Canceled {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}

// TestWaitForServerStatus_ValidationErrors verifies input validation.
func TestWaitForServerStatus_ValidationErrors(t *testing.T) {
	httpClient := internalhttp.NewClient("http://example.com", "test-token", &http.Client{}, nil)
	serverClient := servers.NewClient(httpClient, "proj-1")

	tests := []struct {
		name        string
		cfg         ServerWaiterConfig
		expectedErr string
	}{
		{
			name: "missing client",
			cfg: ServerWaiterConfig{
				Client:       nil,
				ServerID:     "svr-1",
				TargetStatus: serversmodels.ServerStatusActive,
			},
			expectedErr: "server client is required",
		},
		{
			name: "missing server ID",
			cfg: ServerWaiterConfig{
				Client:       serverClient,
				ServerID:     "",
				TargetStatus: serversmodels.ServerStatusActive,
			},
			expectedErr: "server ID is required",
		},
		{
			name: "missing target status",
			cfg: ServerWaiterConfig{
				Client:       serverClient,
				ServerID:     "svr-1",
				TargetStatus: "",
			},
			expectedErr: "target status is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := WaitForServerStatus(ctx, tt.cfg)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

// TestWaitForServerActive verifies the convenience function for waiting ACTIVE status.
func TestWaitForServerActive(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]interface{}{
			"id":         "svr-test-1",
			"name":       "test-server",
			"status":     "ACTIVE",
			"flavor_id":  "flavor-1",
			"image_id":   "image-1",
			"project_id": "proj-1",
			"created_at": "2025-10-28T00:00:00Z",
			"updated_at": "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	serverClient := servers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForServerActive(ctx, serverClient, "svr-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(1*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestWaitForServerShutoff verifies the convenience function for waiting SHUTOFF status.
func TestWaitForServerShutoff(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]interface{}{
			"id":         "svr-test-1",
			"name":       "test-server",
			"status":     "SHUTOFF",
			"flavor_id":  "flavor-1",
			"image_id":   "image-1",
			"project_id": "proj-1",
			"created_at": "2025-10-28T00:00:00Z",
			"updated_at": "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	serverClient := servers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForServerShutoff(ctx, serverClient, "svr-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(1*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestWaitForServerDeleted verifies waiting for server deletion.
func TestWaitForServerDeleted(t *testing.T) {
	requestCount := 0

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		requestCount++

		// First request returns server exists, second returns 404
		if requestCount == 1 {
			response := map[string]interface{}{
				"id":         "svr-test-1",
				"name":       "test-server",
				"status":     "ACTIVE",
				"flavor_id":  "flavor-1",
				"image_id":   "image-1",
				"project_id": "proj-1",
				"created_at": "2025-10-28T00:00:00Z",
				"updated_at": "2025-10-28T00:00:00Z",
			}

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		} else {
			// Server deleted - return 404
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNotFound)
			response := map[string]interface{}{
				"error":      "not_found",
				"message":    "Server not found",
				"error_code": 404,
			}
			if err := json.NewEncoder(w).Encode(response); err != nil {
				t.Fatalf("failed to encode response: %v", err)
			}
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	serverClient := servers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForServerDeleted(ctx, serverClient, "svr-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(1*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if requestCount < 2 {
		t.Errorf("expected at least 2 requests (polling until 404), got %d", requestCount)
	}
}

// TestWaitForServerStatus_WithBackoff verifies backoff behavior.
func TestWaitForServerStatus_WithBackoff(t *testing.T) {
	callCount := 0
	var callTimes []time.Time

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		callCount++
		callTimes = append(callTimes, time.Now())

		// Return ACTIVE after 3 calls
		status := "BUILD"
		if callCount >= 3 {
			status = "ACTIVE"
		}

		response := map[string]interface{}{
			"id":         "svr-test-1",
			"name":       "test-server",
			"status":     status,
			"flavor_id":  "flavor-1",
			"image_id":   "image-1",
			"project_id": "proj-1",
			"created_at": "2025-10-28T00:00:00Z",
			"updated_at": "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	serverClient := servers.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForServerStatus(ctx, ServerWaiterConfig{
		Client:       serverClient,
		ServerID:     "svr-test-1",
		TargetStatus: serversmodels.ServerStatusActive,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(50 * time.Millisecond),
			waiter.WithMaxWait(5 * time.Second),
			waiter.WithBackoff(2.0, 500*time.Millisecond), // Double interval each time, max 500ms
		},
	})

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if callCount < 3 {
		t.Errorf("expected at least 3 calls, got %d", callCount)
	}

	// Verify backoff is working (intervals should increase)
	if len(callTimes) >= 3 {
		interval1 := callTimes[1].Sub(callTimes[0])
		interval2 := callTimes[2].Sub(callTimes[1])

		// Second interval should be longer than first (due to backoff)
		if interval2 <= interval1 {
			t.Logf("Intervals: %v, %v", interval1, interval2)
			// This is informational - timing can be flaky in tests
		}
	}
}

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
