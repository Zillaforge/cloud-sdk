package vps

import (
	"context"
	"fmt"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	sharesmodels "github.com/Zillaforge/cloud-sdk/models/vps/shares"
	"github.com/Zillaforge/cloud-sdk/modules/vps/shares"
)

// ShareWaiterConfig holds configuration for share state waiting.
type ShareWaiterConfig struct {
	// Client is the shares client used to poll share state
	Client *shares.Client

	// ShareID is the ID of the share to monitor
	ShareID string

	// TargetStatus is the desired share state
	TargetStatus sharesmodels.ShareStatus

	// WaiterOptions are passed to the underlying waiter framework
	WaiterOptions []waiter.Option
}

// WaitForShareStatus polls a share until it reaches the target status.
// It returns an error if:
// - The share reaches ERROR status (unless that's the target)
// - The context is canceled
// - The maximum wait duration is exceeded
// - An error occurs during polling
//
// Example usage:
//
//	err := vps.WaitForShareStatus(ctx, vps.ShareWaiterConfig{
//	    Client:       shareClient,
//	    ShareID:      "share-123",
//	    TargetStatus: sharesmodels.ShareStatusAvailable,
//	    WaiterOptions: []waiter.Option{
//	        waiter.WithInterval(5 * time.Second),
//	        waiter.WithMaxWait(10 * time.Minute),
//	        waiter.WithBackoff(1.5, 30 * time.Second),
//	    },
//	})
func WaitForShareStatus(ctx context.Context, cfg ShareWaiterConfig) error {
	if cfg.Client == nil {
		return fmt.Errorf("share client is required")
	}
	if cfg.ShareID == "" {
		return fmt.Errorf("share ID is required")
	}
	if cfg.TargetStatus == "" {
		return fmt.Errorf("target status is required")
	}

	// Default waiter options for shares (can be overridden)
	defaultOpts := []waiter.Option{
		waiter.WithInterval(5 * time.Second),
		waiter.WithMaxWait(10 * time.Minute),
		waiter.WithBackoff(1.2, 30*time.Second),
	}

	// Merge user options (user options take precedence)
	opts := append(defaultOpts, cfg.WaiterOptions...)

	// Create the state check function
	checkState := func(ctx context.Context) (bool, error) {
		// Get current share state
		share, err := cfg.Client.Get(ctx, cfg.ShareID)
		if err != nil {
			return false, fmt.Errorf("failed to get share status: %w", err)
		}

		currentStatus := sharesmodels.ShareStatus(share.Status)

		// Check if we've reached the target status
		if currentStatus == cfg.TargetStatus {
			return true, nil
		}

		// If share is in ERROR state and that's not our target, fail immediately
		if currentStatus == sharesmodels.ShareStatusError && cfg.TargetStatus != sharesmodels.ShareStatusError {
			return false, fmt.Errorf("share entered ERROR state while waiting for %s", cfg.TargetStatus)
		}

		// Continue polling
		return false, nil
	}

	// Use the generic waiter framework
	return waiter.Wait(ctx, checkState, opts...)
}

// WaitForShareAvailable is a convenience function that waits for a share to become AVAILABLE.
func WaitForShareAvailable(ctx context.Context, client *shares.Client, shareID string, opts ...waiter.Option) error {
	return WaitForShareStatus(ctx, ShareWaiterConfig{
		Client:        client,
		ShareID:       shareID,
		TargetStatus:  sharesmodels.ShareStatusAvailable,
		WaiterOptions: opts,
	})
}
