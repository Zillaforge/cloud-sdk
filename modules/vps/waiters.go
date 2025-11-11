package vps

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/types"
	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	floatingipsmodels "github.com/Zillaforge/cloud-sdk/models/vps/floatingips"
	serversmodels "github.com/Zillaforge/cloud-sdk/models/vps/servers"
	"github.com/Zillaforge/cloud-sdk/modules/vps/floatingips"
	"github.com/Zillaforge/cloud-sdk/modules/vps/servers"
)

// ServerWaiterConfig holds configuration for server state waiting.
type ServerWaiterConfig struct {
	// Client is the servers client used to poll server state
	Client *servers.Client

	// ServerID is the ID of the server to monitor
	ServerID string

	// TargetStatus is the desired server state
	TargetStatus serversmodels.ServerStatus

	// WaiterOptions are passed to the underlying waiter framework
	WaiterOptions []waiter.Option
}

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

// WaitForServerStatus polls a server until it reaches the target status.
// It returns an error if:
// - The server reaches ERROR status (unless that's the target)
// - The context is canceled
// - The maximum wait duration is exceeded
// - An error occurs during polling
//
// Example usage:
//
//	err := vps.WaitForServerStatus(ctx, vps.ServerWaiterConfig{
//	    Client:       serverClient,
//	    ServerID:     "svr-123",
//	    TargetStatus: vps.ServerStatusActive,
//	    WaiterOptions: []waiter.Option{
//	        waiter.WithInterval(5 * time.Second),
//	        waiter.WithMaxWait(10 * time.Minute),
//	        waiter.WithBackoff(1.5, 30 * time.Second),
//	    },
//	})
func WaitForServerStatus(ctx context.Context, cfg ServerWaiterConfig) error {
	if cfg.Client == nil {
		return fmt.Errorf("server client is required")
	}
	if cfg.ServerID == "" {
		return fmt.Errorf("server ID is required")
	}
	if cfg.TargetStatus == "" {
		return fmt.Errorf("target status is required")
	}

	// Default waiter options for servers (can be overridden)
	defaultOpts := []waiter.Option{
		waiter.WithInterval(5 * time.Second),
		waiter.WithMaxWait(10 * time.Minute),
		waiter.WithBackoff(1.2, 30*time.Second),
	}

	// Merge user options (user options take precedence)
	opts := append(defaultOpts, cfg.WaiterOptions...)

	// Create the state check function
	checkState := func(ctx context.Context) (bool, error) {
		// Get current server state
		serverResource, err := cfg.Client.Get(ctx, cfg.ServerID)
		if err != nil {
			return false, fmt.Errorf("failed to get server status: %w", err)
		}

		currentStatus := serversmodels.ServerStatus(serverResource.Server.Status)

		// Check if we've reached the target status
		if currentStatus == cfg.TargetStatus {
			return true, nil
		}

		// If server is in ERROR state and that's not our target, fail immediately
		if currentStatus == serversmodels.ServerStatusError && cfg.TargetStatus != serversmodels.ServerStatusError {
			return false, fmt.Errorf("server entered ERROR state while waiting for %s", cfg.TargetStatus)
		}

		// Continue polling
		return false, nil
	}

	// Use the generic waiter framework
	return waiter.Wait(ctx, checkState, opts...)
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

// WaitForServerActive is a convenience function that waits for a server to become ACTIVE.
func WaitForServerActive(ctx context.Context, client *servers.Client, serverID string, opts ...waiter.Option) error {
	return WaitForServerStatus(ctx, ServerWaiterConfig{
		Client:        client,
		ServerID:      serverID,
		TargetStatus:  serversmodels.ServerStatusActive,
		WaiterOptions: opts,
	})
}

// WaitForServerShutoff is a convenience function that waits for a server to become SHUTOFF.
func WaitForServerShutoff(ctx context.Context, client *servers.Client, serverID string, opts ...waiter.Option) error {
	return WaitForServerStatus(ctx, ServerWaiterConfig{
		Client:        client,
		ServerID:      serverID,
		TargetStatus:  serversmodels.ServerStatusShutoff,
		WaiterOptions: opts,
	})
}

// WaitForServerDeleted is a convenience function that waits for a server to be deleted.
// This function handles the case where Get returns a 404 (not found) error,
// which indicates the server has been successfully deleted.
func WaitForServerDeleted(ctx context.Context, client *servers.Client, serverID string, opts ...waiter.Option) error {
	// Default waiter options for deletion (can be overridden)
	defaultOpts := []waiter.Option{
		waiter.WithInterval(3 * time.Second),
		waiter.WithMaxWait(5 * time.Minute),
	}

	// Merge user options
	allOpts := append(defaultOpts, opts...)

	checkState := func(ctx context.Context) (bool, error) {
		_, err := client.Get(ctx, serverID)
		if err != nil {
			// Check if this is an SDKError with 404 status
			var sdkErr *types.SDKError
			if errors.As(err, &sdkErr) && sdkErr.StatusCode == 404 {
				return true, nil // Server is deleted
			}
			// Other errors should be propagated
			return false, fmt.Errorf("error checking server deletion status: %w", err)
		}
		// Server still exists, keep waiting
		return false, nil
	}

	return waiter.Wait(ctx, checkState, allOpts...)
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
