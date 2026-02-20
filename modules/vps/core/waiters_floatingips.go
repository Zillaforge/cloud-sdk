package vps

import (
	"context"
	"fmt"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	floatingipsmodels "github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
	"github.com/Zillaforge/cloud-sdk/modules/vps/floatingips"
)

// FloatingIPWaiterConfig holds configuration for floating IP state waiting.
type FloatingIPWaiterConfig struct {
	// Client is the floating IPs client used to poll floating IP state
	Client *floatingips.Client

	// FloatingIPID is the ID of the floating IP to monitor
	FloatingIPID string

	// TargetStatus is the desired floating IP state
	TargetStatus floatingipsmodels.FloatingIPStatus

	// WaiterOptions are passed to the underlying waiter framework
	WaiterOptions []waiter.Option
}

// WaitForFloatingIPStatus polls a floating IP until it reaches the target status.
// It returns an error if:
// - The floating IP reaches REJECTED status (unless that's the target)
// - The context is canceled
// - The maximum wait duration is exceeded
// - An error occurs during polling
//
// Example usage:
//
//	err := vps.WaitForFloatingIPStatus(ctx, vps.FloatingIPWaiterConfig{
//	    Client:        floatingIPClient,
//	    FloatingIPID:  "fip-123",
//	    TargetStatus:  vps.FloatingIPStatusActive,
//	    WaiterOptions: []waiter.Option{
//	        waiter.WithInterval(5 * time.Second),
//	        waiter.WithMaxWait(10 * time.Minute),
//	        waiter.WithBackoff(1.5, 30 * time.Second),
//	    },
//	})
func WaitForFloatingIPStatus(ctx context.Context, cfg FloatingIPWaiterConfig) error {
	if cfg.Client == nil {
		return fmt.Errorf("floating IP client is required")
	}
	if cfg.FloatingIPID == "" {
		return fmt.Errorf("floating IP ID is required")
	}
	if cfg.TargetStatus == "" {
		return fmt.Errorf("target status is required")
	}

	// Default waiter options for floating IPs (can be overridden)
	defaultOpts := []waiter.Option{
		waiter.WithInterval(5 * time.Second),
		waiter.WithMaxWait(10 * time.Minute),
		waiter.WithBackoff(1.2, 30*time.Second),
	}

	// Merge user options (user options take precedence)
	opts := append(defaultOpts, cfg.WaiterOptions...)

	// Create the state check function
	checkState := func(ctx context.Context) (bool, error) {
		// Get current floating IP state
		floatingIP, err := cfg.Client.Get(ctx, cfg.FloatingIPID)
		if err != nil {
			return false, fmt.Errorf("failed to get floating IP status: %w", err)
		}

		currentStatus := floatingipsmodels.FloatingIPStatus(floatingIP.Status)

		// Check if we've reached the target status
		if currentStatus == cfg.TargetStatus {
			return true, nil
		}

		// If floating IP is in REJECTED state and that's not our target, fail immediately
		if currentStatus == floatingipsmodels.FloatingIPStatusRejected && cfg.TargetStatus != floatingipsmodels.FloatingIPStatusRejected {
			return false, fmt.Errorf("floating IP entered REJECTED state while waiting for %s", cfg.TargetStatus)
		}

		// Continue polling
		return false, nil
	}

	// Use the generic waiter framework
	return waiter.Wait(ctx, checkState, opts...)
}

// WaitForFloatingIPActive is a convenience function that waits for a floating IP to become ACTIVE.
func WaitForFloatingIPActive(ctx context.Context, client *floatingips.Client, floatingIPID string, opts ...waiter.Option) error {
	return WaitForFloatingIPStatus(ctx, FloatingIPWaiterConfig{
		Client:        client,
		FloatingIPID:  floatingIPID,
		TargetStatus:  floatingipsmodels.FloatingIPStatusActive,
		WaiterOptions: opts,
	})
}
