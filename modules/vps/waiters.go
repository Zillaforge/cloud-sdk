package vps

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/types"
	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	"github.com/Zillaforge/cloud-sdk/modules/vps/servers"
)

// ServerStatus represents the possible states of a server.
type ServerStatus string

const (
	// ServerStatusBuild indicates the server is being created.
	ServerStatusBuild ServerStatus = "BUILD"

	// ServerStatusActive indicates the server is running.
	ServerStatusActive ServerStatus = "ACTIVE"

	// ServerStatusShutoff indicates the server is stopped.
	ServerStatusShutoff ServerStatus = "SHUTOFF"

	// ServerStatusError indicates the server encountered an error.
	ServerStatusError ServerStatus = "ERROR"

	// ServerStatusDeleted indicates the server has been deleted.
	ServerStatusDeleted ServerStatus = "DELETED"

	// ServerStatusRebooting indicates the server is rebooting.
	ServerStatusRebooting ServerStatus = "REBOOT"

	// ServerStatusResizing indicates the server is being resized.
	ServerStatusResizing ServerStatus = "RESIZE"
)

// ServerWaiterConfig holds configuration for server state waiting.
type ServerWaiterConfig struct {
	// Client is the servers client used to poll server state
	Client *servers.Client

	// ServerID is the ID of the server to monitor
	ServerID string

	// TargetStatus is the desired server state
	TargetStatus ServerStatus

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

		currentStatus := ServerStatus(serverResource.Server.Status)

		// Check if we've reached the target status
		if currentStatus == cfg.TargetStatus {
			return true, nil
		}

		// If server is in ERROR state and that's not our target, fail immediately
		if currentStatus == ServerStatusError && cfg.TargetStatus != ServerStatusError {
			return false, fmt.Errorf("server entered ERROR state while waiting for %s", cfg.TargetStatus)
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
		TargetStatus:  ServerStatusActive,
		WaiterOptions: opts,
	})
}

// WaitForServerShutoff is a convenience function that waits for a server to become SHUTOFF.
func WaitForServerShutoff(ctx context.Context, client *servers.Client, serverID string, opts ...waiter.Option) error {
	return WaitForServerStatus(ctx, ServerWaiterConfig{
		Client:        client,
		ServerID:      serverID,
		TargetStatus:  ServerStatusShutoff,
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
