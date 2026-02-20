package vps

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/types"
	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	lbmodels "github.com/Zillaforge/cloud-sdk/models/vps/loadbalancers"
	"github.com/Zillaforge/cloud-sdk/modules/vps/loadbalancers"
)

// LoadBalancerWaiterConfig holds configuration for load balancer state waiting.
type LoadBalancerWaiterConfig struct {
	// Client is the load balancers client used to poll load balancer state
	Client *loadbalancers.Client

	// LoadBalancerID is the ID of the load balancer to monitor
	LoadBalancerID string

	// TargetStatus is the desired load balancer state
	TargetStatus lbmodels.LoadBalancerStatus

	// WaiterOptions are passed to the underlying waiter framework
	WaiterOptions []waiter.Option
}

// WaitForLoadBalancerStatus polls a load balancer until it reaches the target status.
// It returns an error if:
// - The load balancer reaches ERROR status (unless that's the target)
// - The context is canceled
// - The maximum wait duration is exceeded
// - An error occurs during polling
//
// Example usage:
//
//	err := vps.WaitForLoadBalancerStatus(ctx, vps.LoadBalancerWaiterConfig{
//	    Client:         lbClient,
//	    LoadBalancerID: "lb-123",
//	    TargetStatus:   lbmodels.LoadBalancerStatusActive,
//	    WaiterOptions: []waiter.Option{
//	        waiter.WithInterval(5 * time.Second),
//	        waiter.WithMaxWait(10 * time.Minute),
//	    },
//	})
func WaitForLoadBalancerStatus(ctx context.Context, cfg LoadBalancerWaiterConfig) error {
	if cfg.Client == nil {
		return fmt.Errorf("load balancer client is required")
	}
	if cfg.LoadBalancerID == "" {
		return fmt.Errorf("load balancer ID is required")
	}
	if cfg.TargetStatus == "" {
		return fmt.Errorf("target status is required")
	}

	// Default waiter options for load balancers (can be overridden)
	defaultOpts := []waiter.Option{
		waiter.WithInterval(5 * time.Second),
		waiter.WithMaxWait(10 * time.Minute),
		waiter.WithBackoff(1.2, 30*time.Second),
	}

	// Merge user options (user options take precedence)
	opts := append(defaultOpts, cfg.WaiterOptions...)

	// Create the state check function
	checkState := func(ctx context.Context) (bool, error) {
		lbResource, err := cfg.Client.Get(ctx, cfg.LoadBalancerID)
		if err != nil {
			return false, fmt.Errorf("failed to get load balancer status: %w", err)
		}

		currentStatus := lbResource.Status

		// Check if we've reached the target status
		if currentStatus == cfg.TargetStatus {
			return true, nil
		}

		// If LB is in ERROR state and that's not our target, fail immediately
		if currentStatus == lbmodels.LoadBalancerStatusError && cfg.TargetStatus != lbmodels.LoadBalancerStatusError {
			return false, fmt.Errorf("load balancer entered ERROR state while waiting for %s", cfg.TargetStatus)
		}

		// Continue polling
		return false, nil
	}

	// Use the generic waiter framework
	return waiter.Wait(ctx, checkState, opts...)
}

// WaitForLoadBalancerActive is a convenience function that waits for a load balancer to become ACTIVE.
func WaitForLoadBalancerActive(ctx context.Context, client *loadbalancers.Client, lbID string, opts ...waiter.Option) error {
	return WaitForLoadBalancerStatus(ctx, LoadBalancerWaiterConfig{
		Client:         client,
		LoadBalancerID: lbID,
		TargetStatus:   lbmodels.LoadBalancerStatusActive,
		WaiterOptions:  opts,
	})
}

// WaitForLoadBalancerDeleted is a convenience function that waits for a load balancer to be deleted.
// This function handles the case where Get returns a 404 (not found) error,
// which indicates the load balancer has been successfully deleted.
func WaitForLoadBalancerDeleted(ctx context.Context, client *loadbalancers.Client, lbID string, opts ...waiter.Option) error {
	defaultOpts := []waiter.Option{
		waiter.WithInterval(3 * time.Second),
		waiter.WithMaxWait(5 * time.Minute),
	}

	opts = append(defaultOpts, opts...)

	checkState := func(ctx context.Context) (bool, error) {
		lbResource, err := client.Get(ctx, lbID)
		if err != nil {
			// Check if this is an SDKError with 404 status
			var sdkErr *types.SDKError
			if errors.As(err, &sdkErr) && sdkErr.StatusCode == 404 {
				return true, nil // Load balancer is deleted
			}
			return false, fmt.Errorf("failed to get load balancer status: %w", err)
		}

		currentStatus := lbResource.Status

		if currentStatus == lbmodels.LoadBalancerStatusDeleted {
			return true, nil
		}

		if currentStatus == lbmodels.LoadBalancerStatusError {
			return false, fmt.Errorf("load balancer entered ERROR state while waiting for deletion")
		}

		return false, nil
	}

	return waiter.Wait(ctx, checkState, opts...)
}
