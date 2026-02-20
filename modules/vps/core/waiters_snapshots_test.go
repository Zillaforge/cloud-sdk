package vps

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	snapshotsmodels "github.com/Zillaforge/cloud-sdk/models/vps/snapshots"
	"github.com/Zillaforge/cloud-sdk/modules/vps/snapshots"
)

// TestWaitForSnapshotStatus_Success verifies waiting for a snapshot to reach target status.
func TestWaitForSnapshotStatus_Success(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus snapshotsmodels.SnapshotStatus
		statusFlow   []string // Sequence of statuses returned by mock snapshot
	}{
		{
			name:         "wait for AVAILABLE from CREATING",
			targetStatus: snapshotsmodels.SnapshotStatusAvailable,
			statusFlow:   []string{"creating", "creating", "available"},
		},
		{
			name:         "wait for DELETED from DELETING",
			targetStatus: snapshotsmodels.SnapshotStatusDeleted,
			statusFlow:   []string{"deleting", "deleted"},
		},
		{
			name:         "already at target status",
			targetStatus: snapshotsmodels.SnapshotStatusAvailable,
			statusFlow:   []string{"available"},
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
					"id":            "snap-test-1",
					"name":          "test-snapshot",
					"volume_id":     "vol-1",
					"size":          10,
					"status":        status,
					"status_reason": "",
					"description":   "test snapshot",
					"project":       map[string]interface{}{"id": "proj-1", "name": "test-project"},
					"project_id":    "proj-1",
					"user":          map[string]interface{}{"id": "user-1", "name": "test-user"},
					"user_id":       "user-1",
					"namespace":     "test-namespace",
					"createdAt":     "2025-10-28T00:00:00Z",
					"updatedAt":     "2025-10-28T00:00:00Z",
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
			snapshotClient := snapshots.NewClient(httpClient, "proj-1")

			// Wait for status with short intervals for testing
			ctx := context.Background()
			err := WaitForSnapshotStatus(ctx, SnapshotWaiterConfig{
				Client:       snapshotClient,
				SnapshotID:   "snap-test-1",
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

// TestWaitForSnapshotStatus_ErrorState verifies handling when snapshot enters ERROR state.
func TestWaitForSnapshotStatus_ErrorState(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]interface{}{
			"id":            "snap-test-1",
			"name":          "test-snapshot",
			"volume_id":     "vol-1",
			"size":          10,
			"status":        "error",
			"status_reason": "creation failed",
			"description":   "test snapshot",
			"project":       map[string]interface{}{"id": "proj-1", "name": "test-project"},
			"project_id":    "proj-1",
			"user":          map[string]interface{}{"id": "user-1", "name": "test-user"},
			"user_id":       "user-1",
			"namespace":     "test-namespace",
			"createdAt":     "2025-10-28T00:00:00Z",
			"updatedAt":     "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	snapshotClient := snapshots.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForSnapshotStatus(ctx, SnapshotWaiterConfig{
		Client:       snapshotClient,
		SnapshotID:   "snap-test-1",
		TargetStatus: snapshotsmodels.SnapshotStatusAvailable,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(1 * time.Second),
		},
	})

	if err == nil {
		t.Fatal("expected error when snapshot enters ERROR state, got nil")
	}

	expectedMsg := "snapshot entered ERROR state"
	if err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("expected error message to start with '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestWaitForSnapshotStatus_Timeout verifies timeout handling.
func TestWaitForSnapshotStatus_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always return CREATING status
		response := map[string]interface{}{
			"id":            "snap-test-1",
			"name":          "test-snapshot",
			"volume_id":     "vol-1",
			"size":          10,
			"status":        "creating",
			"status_reason": "",
			"description":   "test snapshot",
			"project":       map[string]interface{}{"id": "proj-1", "name": "test-project"},
			"project_id":    "proj-1",
			"user":          map[string]interface{}{"id": "user-1", "name": "test-user"},
			"user_id":       "user-1",
			"namespace":     "test-namespace",
			"createdAt":     "2025-10-28T00:00:00Z",
			"updatedAt":     "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	// Create HTTP client with longer timeout to avoid request timeout
	httpClient := &http.Client{Timeout: 30 * time.Second}
	baseClient := internalhttp.NewClient(mockServer.URL, "test-token", httpClient, nil)
	snapshotClient := snapshots.NewClient(baseClient, "proj-1")

	ctx := context.Background()
	err := WaitForSnapshotStatus(ctx, SnapshotWaiterConfig{
		Client:       snapshotClient,
		SnapshotID:   "snap-test-1",
		TargetStatus: snapshotsmodels.SnapshotStatusAvailable,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(50 * time.Millisecond),
			waiter.WithMaxWait(200 * time.Millisecond),
		},
	})

	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}

	if err != waiter.ErrWaitTimeout {
		t.Errorf("expected ErrWaitTimeout, got: %v", err)
	}
}

// TestWaitForSnapshotStatus_ContextCancellation verifies context cancellation.
func TestWaitForSnapshotStatus_ContextCancellation(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always return CREATING status
		response := map[string]interface{}{
			"id":            "snap-test-1",
			"name":          "test-snapshot",
			"volume_id":     "vol-1",
			"size":          10,
			"status":        "creating",
			"status_reason": "",
			"description":   "test snapshot",
			"project":       map[string]interface{}{"id": "proj-1", "name": "test-project"},
			"project_id":    "proj-1",
			"user":          map[string]interface{}{"id": "user-1", "name": "test-user"},
			"user_id":       "user-1",
			"namespace":     "test-namespace",
			"createdAt":     "2025-10-28T00:00:00Z",
			"updatedAt":     "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	snapshotClient := snapshots.NewClient(httpClient, "proj-1")

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := WaitForSnapshotStatus(ctx, SnapshotWaiterConfig{
		Client:       snapshotClient,
		SnapshotID:   "snap-test-1",
		TargetStatus: snapshotsmodels.SnapshotStatusAvailable,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(5 * time.Second),
		},
	})

	if err == nil {
		t.Fatal("expected context cancellation error, got nil")
	}

	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got: %v", err)
	}
}

// TestWaitForSnapshotStatus_ValidationErrors verifies input validation.
func TestWaitForSnapshotStatus_ValidationErrors(t *testing.T) {
	httpClient := internalhttp.NewClient("http://example.com", "test-token", &http.Client{}, nil)
	snapshotClient := snapshots.NewClient(httpClient, "proj-1")

	tests := []struct {
		name        string
		cfg         SnapshotWaiterConfig
		expectedErr string
	}{
		{
			name: "missing client",
			cfg: SnapshotWaiterConfig{
				Client:       nil,
				SnapshotID:   "snap-1",
				TargetStatus: snapshotsmodels.SnapshotStatusAvailable,
			},
			expectedErr: "snapshot client is required",
		},
		{
			name: "missing snapshot ID",
			cfg: SnapshotWaiterConfig{
				Client:       snapshotClient,
				SnapshotID:   "",
				TargetStatus: snapshotsmodels.SnapshotStatusAvailable,
			},
			expectedErr: "snapshot ID is required",
		},
		{
			name: "missing target status",
			cfg: SnapshotWaiterConfig{
				Client:       snapshotClient,
				SnapshotID:   "snap-1",
				TargetStatus: "",
			},
			expectedErr: "target status is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := WaitForSnapshotStatus(ctx, tt.cfg)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

// TestWaitForSnapshotAvailable verifies the convenience function for waiting AVAILABLE status.
func TestWaitForSnapshotAvailable(t *testing.T) {
	currentStatusIndex := 0
	statusFlow := []string{"creating", "available"}

	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		status := statusFlow[currentStatusIndex]
		if currentStatusIndex < len(statusFlow)-1 {
			currentStatusIndex++
		}

		response := map[string]interface{}{
			"id":            "snap-test-1",
			"name":          "test-snapshot",
			"volume_id":     "vol-1",
			"size":          10,
			"status":        status,
			"status_reason": "",
			"description":   "test snapshot",
			"project":       map[string]interface{}{"id": "proj-1", "name": "test-project"},
			"project_id":    "proj-1",
			"user":          map[string]interface{}{"id": "user-1", "name": "test-user"},
			"user_id":       "user-1",
			"namespace":     "test-namespace",
			"createdAt":     "2025-10-28T00:00:00Z",
			"updatedAt":     "2025-10-28T00:00:00Z",
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
	snapshotClient := snapshots.NewClient(httpClient, "proj-1")

	// Test convenience function
	ctx := context.Background()
	err := WaitForSnapshotAvailable(ctx, snapshotClient, "snap-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(2*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
