package vps

import (
	"context"
	"fmt"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	volumesmodels "github.com/Zillaforge/cloud-sdk/models/vps/volumes"
	"github.com/Zillaforge/cloud-sdk/modules/vps/volumes"
)

// VolumeWaiterConfig holds configuration for volume state waiting.
type VolumeWaiterConfig struct {
	// Client is the volumes client used to poll volume state
	Client *volumes.Client

	// VolumeID is the ID of the volume to monitor
	VolumeID string

	// TargetStatus is the desired volume state
	TargetStatus volumesmodels.VolumeStatus

	// WaiterOptions are passed to the underlying waiter framework
	WaiterOptions []waiter.Option
}

// WaitForVolumeStatus polls a volume until it reaches the target status.
// It returns an error if:
// - The volume reaches ERROR status (unless that's the target)
// - The context is canceled
// - The maximum wait duration is exceeded
// - An error occurs during polling
//
// Example usage:
//
//	err := vps.WaitForVolumeStatus(ctx, vps.VolumeWaiterConfig{
//	    Client:       volumeClient,
//	    VolumeID:     "vol-123",
//	    TargetStatus: vps.VolumeStatusAvailable,
//	    WaiterOptions: []waiter.Option{
//	        waiter.WithInterval(5 * time.Second),
//	        waiter.WithMaxWait(10 * time.Minute),
//	        waiter.WithBackoff(1.5, 30 * time.Second),
//	    },
//	})
func WaitForVolumeStatus(ctx context.Context, cfg VolumeWaiterConfig) error {
	if cfg.Client == nil {
		return fmt.Errorf("volume client is required")
	}
	if cfg.VolumeID == "" {
		return fmt.Errorf("volume ID is required")
	}
	if cfg.TargetStatus == "" {
		return fmt.Errorf("target status is required")
	}

	// Default waiter options for volumes (can be overridden)
	defaultOpts := []waiter.Option{
		waiter.WithInterval(5 * time.Second),
		waiter.WithMaxWait(10 * time.Minute),
		waiter.WithBackoff(1.2, 30*time.Second),
	}

	// Merge user options (user options take precedence)
	opts := append(defaultOpts, cfg.WaiterOptions...)

	// Create the state check function
	checkState := func(ctx context.Context) (bool, error) {
		// Get current volume state
		volume, err := cfg.Client.Get(ctx, cfg.VolumeID)
		if err != nil {
			return false, fmt.Errorf("failed to get volume status: %w", err)
		}

		currentStatus := volumesmodels.VolumeStatus(volume.Status)

		// Check if we've reached the target status
		if currentStatus == cfg.TargetStatus {
			return true, nil
		}

		// If volume is in ERROR state and that's not our target, fail immediately
		if currentStatus == volumesmodels.VolumeStatusError && cfg.TargetStatus != volumesmodels.VolumeStatusError {
			return false, fmt.Errorf("volume entered ERROR state while waiting for %s", cfg.TargetStatus)
		}

		// Continue polling
		return false, nil
	}

	// Use the generic waiter framework
	return waiter.Wait(ctx, checkState, opts...)
}

// WaitForVolumeAvailable is a convenience function that waits for a volume to become AVAILABLE.
func WaitForVolumeAvailable(ctx context.Context, client *volumes.Client, volumeID string, opts ...waiter.Option) error {
	return WaitForVolumeStatus(ctx, VolumeWaiterConfig{
		Client:        client,
		VolumeID:      volumeID,
		TargetStatus:  volumesmodels.VolumeStatusAvailable,
		WaiterOptions: opts,
	})
}

// WaitForVolumeInUse is a convenience function that waits for a volume to become IN-USE.
func WaitForVolumeInUse(ctx context.Context, client *volumes.Client, volumeID string, opts ...waiter.Option) error {
	return WaitForVolumeStatus(ctx, VolumeWaiterConfig{
		Client:        client,
		VolumeID:      volumeID,
		TargetStatus:  volumesmodels.VolumeStatusInUse,
		WaiterOptions: opts,
	})
}
