package vps

import (
	"context"
	"fmt"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	snapshotsmodels "github.com/Zillaforge/cloud-sdk/models/vps/snapshots"
	"github.com/Zillaforge/cloud-sdk/modules/vps/snapshots"
)

// SnapshotWaiterConfig holds configuration for snapshot state waiting.
type SnapshotWaiterConfig struct {
	// Client is the snapshots client used to poll snapshot state
	Client *snapshots.Client

	// SnapshotID is the ID of the snapshot to monitor
	SnapshotID string

	// TargetStatus is the desired snapshot state
	TargetStatus snapshotsmodels.SnapshotStatus

	// WaiterOptions are passed to the underlying waiter framework
	WaiterOptions []waiter.Option
}

// WaitForSnapshotStatus polls a snapshot until it reaches the target status.
// It returns an error if:
// - The snapshot reaches ERROR status (unless that's the target)
// - The context is canceled
// - The maximum wait duration is exceeded
// - An error occurs during polling
//
// Example usage:
//
//	err := vps.WaitForSnapshotStatus(ctx, vps.SnapshotWaiterConfig{
//	    Client:       snapshotClient,
//	    SnapshotID:   "snap-123",
//	    TargetStatus: vps.SnapshotStatusAvailable,
//	    WaiterOptions: []waiter.Option{
//	        waiter.WithInterval(5 * time.Second),
//	        waiter.WithMaxWait(10 * time.Minute),
//	        waiter.WithBackoff(1.5, 30 * time.Second),
//	    },
//	})
func WaitForSnapshotStatus(ctx context.Context, cfg SnapshotWaiterConfig) error {
	if cfg.Client == nil {
		return fmt.Errorf("snapshot client is required")
	}
	if cfg.SnapshotID == "" {
		return fmt.Errorf("snapshot ID is required")
	}
	if cfg.TargetStatus == "" {
		return fmt.Errorf("target status is required")
	}

	// Default waiter options for snapshots (can be overridden)
	defaultOpts := []waiter.Option{
		waiter.WithInterval(5 * time.Second),
		waiter.WithMaxWait(10 * time.Minute),
		waiter.WithBackoff(1.2, 30*time.Second),
	}

	// Merge user options (user options take precedence)
	opts := append(defaultOpts, cfg.WaiterOptions...)

	// Create the state check function
	checkState := func(ctx context.Context) (bool, error) {
		// Get current snapshot state
		snapshot, err := cfg.Client.Get(ctx, cfg.SnapshotID)
		if err != nil {
			return false, fmt.Errorf("failed to get snapshot status: %w", err)
		}

		currentStatus := snapshotsmodels.SnapshotStatus(snapshot.Status)

		// Check if we've reached the target status
		if currentStatus == cfg.TargetStatus {
			return true, nil
		}

		// If snapshot is in ERROR state and that's not our target, fail immediately
		if currentStatus == snapshotsmodels.SnapshotStatusError && cfg.TargetStatus != snapshotsmodels.SnapshotStatusError {
			return false, fmt.Errorf("snapshot entered ERROR state while waiting for %s", cfg.TargetStatus)
		}

		// Continue polling
		return false, nil
	}

	// Use the generic waiter framework
	return waiter.Wait(ctx, checkState, opts...)
}

// WaitForSnapshotAvailable is a convenience function that waits for a snapshot to become AVAILABLE.
func WaitForSnapshotAvailable(ctx context.Context, client *snapshots.Client, snapshotID string, opts ...waiter.Option) error {
	return WaitForSnapshotStatus(ctx, SnapshotWaiterConfig{
		Client:        client,
		SnapshotID:    snapshotID,
		TargetStatus:  snapshotsmodels.SnapshotStatusAvailable,
		WaiterOptions: opts,
	})
}
