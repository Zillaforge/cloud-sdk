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
	volumesmodels "github.com/Zillaforge/cloud-sdk/models/vps/volumes"
	"github.com/Zillaforge/cloud-sdk/modules/vps/volumes"
)

// TestWaitForVolumeStatus_Success verifies waiting for a volume to reach target status.
func TestWaitForVolumeStatus_Success(t *testing.T) {
	tests := []struct {
		name         string
		targetStatus volumesmodels.VolumeStatus
		statusFlow   []string // Sequence of statuses returned by mock volume
	}{
		{
			name:         "wait for AVAILABLE from CREATING",
			targetStatus: volumesmodels.VolumeStatusAvailable,
			statusFlow:   []string{"creating", "creating", "available"},
		},
		{
			name:         "wait for IN-USE from AVAILABLE",
			targetStatus: volumesmodels.VolumeStatusInUse,
			statusFlow:   []string{"available", "in-use"},
		},
		{
			name:         "already at target status",
			targetStatus: volumesmodels.VolumeStatusAvailable,
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
					"id":          "vol-test-1",
					"name":        "test-volume",
					"status":      status,
					"size":        10,
					"type":        "ssd",
					"project_id":  "proj-1",
					"user_id":     "user-1",
					"attachments": []interface{}{},
					"created_at":  "2025-10-28T00:00:00Z",
					"updated_at":  "2025-10-28T00:00:00Z",
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
			volumeClient := volumes.NewClient(httpClient, "proj-1")

			// Wait for status with short intervals for testing
			ctx := context.Background()
			err := WaitForVolumeStatus(ctx, VolumeWaiterConfig{
				Client:       volumeClient,
				VolumeID:     "vol-test-1",
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

// TestWaitForVolumeStatus_ErrorState verifies handling when volume enters ERROR state.
func TestWaitForVolumeStatus_ErrorState(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		response := map[string]interface{}{
			"id":          "vol-test-1",
			"name":        "test-volume",
			"status":      "error",
			"size":        10,
			"type":        "ssd",
			"project_id":  "proj-1",
			"user_id":     "user-1",
			"attachments": []interface{}{},
			"created_at":  "2025-10-28T00:00:00Z",
			"updated_at":  "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	volumeClient := volumes.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForVolumeStatus(ctx, VolumeWaiterConfig{
		Client:       volumeClient,
		VolumeID:     "vol-test-1",
		TargetStatus: volumesmodels.VolumeStatusAvailable,
		WaiterOptions: []waiter.Option{
			waiter.WithInterval(10 * time.Millisecond),
			waiter.WithMaxWait(1 * time.Second),
		},
	})

	if err == nil {
		t.Fatal("expected error when volume enters ERROR state, got nil")
	}

	expectedMsg := "volume entered ERROR state"
	if err.Error()[:len(expectedMsg)] != expectedMsg {
		t.Errorf("expected error message to start with '%s', got '%s'", expectedMsg, err.Error())
	}
}

// TestWaitForVolumeStatus_Timeout verifies timeout handling.
func TestWaitForVolumeStatus_Timeout(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always return CREATING status
		response := map[string]interface{}{
			"id":          "vol-test-1",
			"name":        "test-volume",
			"status":      "creating",
			"size":        10,
			"type":        "ssd",
			"project_id":  "proj-1",
			"user_id":     "user-1",
			"attachments": []interface{}{},
			"created_at":  "2025-10-28T00:00:00Z",
			"updated_at":  "2025-10-28T00:00:00Z",
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
	volumeClient := volumes.NewClient(baseClient, "proj-1")

	ctx := context.Background()
	err := WaitForVolumeStatus(ctx, VolumeWaiterConfig{
		Client:       volumeClient,
		VolumeID:     "vol-test-1",
		TargetStatus: volumesmodels.VolumeStatusAvailable,
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

// TestWaitForVolumeStatus_ContextCancellation verifies context cancellation.
func TestWaitForVolumeStatus_ContextCancellation(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		// Always return CREATING status
		response := map[string]interface{}{
			"id":          "vol-test-1",
			"name":        "test-volume",
			"status":      "creating",
			"size":        10,
			"type":        "ssd",
			"project_id":  "proj-1",
			"user_id":     "user-1",
			"attachments": []interface{}{},
			"created_at":  "2025-10-28T00:00:00Z",
			"updated_at":  "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	volumeClient := volumes.NewClient(httpClient, "proj-1")

	ctx, cancel := context.WithCancel(context.Background())

	// Cancel after a short delay
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()

	err := WaitForVolumeStatus(ctx, VolumeWaiterConfig{
		Client:       volumeClient,
		VolumeID:     "vol-test-1",
		TargetStatus: volumesmodels.VolumeStatusAvailable,
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

// TestWaitForVolumeStatus_ValidationErrors verifies input validation.
func TestWaitForVolumeStatus_ValidationErrors(t *testing.T) {
	httpClient := internalhttp.NewClient("http://example.com", "test-token", &http.Client{}, nil)
	volumeClient := volumes.NewClient(httpClient, "proj-1")

	tests := []struct {
		name        string
		cfg         VolumeWaiterConfig
		expectedErr string
	}{
		{
			name: "missing client",
			cfg: VolumeWaiterConfig{
				Client:       nil,
				VolumeID:     "vol-1",
				TargetStatus: volumesmodels.VolumeStatusAvailable,
			},
			expectedErr: "volume client is required",
		},
		{
			name: "missing volume ID",
			cfg: VolumeWaiterConfig{
				Client:       volumeClient,
				VolumeID:     "",
				TargetStatus: volumesmodels.VolumeStatusAvailable,
			},
			expectedErr: "volume ID is required",
		},
		{
			name: "missing target status",
			cfg: VolumeWaiterConfig{
				Client:       volumeClient,
				VolumeID:     "vol-1",
				TargetStatus: "",
			},
			expectedErr: "target status is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			err := WaitForVolumeStatus(ctx, tt.cfg)

			if err == nil {
				t.Fatal("expected error, got nil")
			}

			if err.Error() != tt.expectedErr {
				t.Errorf("expected error '%s', got '%s'", tt.expectedErr, err.Error())
			}
		})
	}
}

// TestWaitForVolumeAvailable verifies the convenience function for waiting AVAILABLE status.
func TestWaitForVolumeAvailable(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		response := map[string]interface{}{
			"id":          "vol-test-1",
			"name":        "test-volume",
			"status":      "available",
			"size":        10,
			"type":        "ssd",
			"project_id":  "proj-1",
			"user_id":     "user-1",
			"attachments": []interface{}{},
			"created_at":  "2025-10-28T00:00:00Z",
			"updated_at":  "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	volumeClient := volumes.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForVolumeAvailable(ctx, volumeClient, "vol-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(2*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestWaitForVolumeInUse verifies the convenience function for waiting IN-USE status.
func TestWaitForVolumeInUse(t *testing.T) {
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("expected GET, got %s", r.Method)
		}

		response := map[string]interface{}{
			"id":          "vol-test-1",
			"name":        "test-volume",
			"status":      "in-use",
			"size":        10,
			"type":        "ssd",
			"project_id":  "proj-1",
			"user_id":     "user-1",
			"attachments": []interface{}{},
			"created_at":  "2025-10-28T00:00:00Z",
			"updated_at":  "2025-10-28T00:00:00Z",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Fatalf("failed to encode response: %v", err)
		}
	}))
	defer mockServer.Close()

	httpClient := internalhttp.NewClient(mockServer.URL, "test-token", &http.Client{}, nil)
	volumeClient := volumes.NewClient(httpClient, "proj-1")

	ctx := context.Background()
	err := WaitForVolumeInUse(ctx, volumeClient, "vol-test-1",
		waiter.WithInterval(10*time.Millisecond),
		waiter.WithMaxWait(2*time.Second),
	)

	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}
