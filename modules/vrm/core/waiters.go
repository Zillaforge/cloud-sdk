package vrm

import (
	"context"
	"fmt"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/waiter"
	commonmodels "github.com/Zillaforge/cloud-sdk/models/vrm/common"
	"github.com/Zillaforge/cloud-sdk/modules/vrm/tags"
)

// TagWaiterConfig holds configuration for tag state waiting.
type TagWaiterConfig struct {
	// Client is the tags client used to poll tag state
	Client *tags.Client

	// TagID is the ID of the tag to monitor
	TagID string

	// TargetStatus is the desired tag state
	TargetStatus commonmodels.TagStatus

	// WaiterOptions are passed to the underlying waiter framework
	WaiterOptions []waiter.Option
}

// WaitForTagStatus polls a tag until it reaches the target status.
// It returns an error if:
// - The tag reaches ERROR status (unless that's the target)
// - The context is canceled
// - The maximum wait duration is exceeded
// - An error occurs during polling
//
// Example usage:
//
//	err := vrm.WaitForTagStatus(ctx, vrm.TagWaiterConfig{
//	    Client:       tagClient,
//	    TagID:        "tag-123",
//	    TargetStatus: commonmodels.TagStatusActive,
//	    WaiterOptions: []waiter.Option{
//	        waiter.WithInterval(5 * time.Second),
//	        waiter.WithMaxWait(10 * time.Minute),
//	        waiter.WithBackoff(1.5, 30 * time.Second),
//	    },
//	})
func WaitForTagStatus(ctx context.Context, cfg TagWaiterConfig) error {
	if cfg.Client == nil {
		return fmt.Errorf("tag client is required")
	}
	if cfg.TagID == "" {
		return fmt.Errorf("tag ID is required")
	}
	if cfg.TargetStatus == "" {
		return fmt.Errorf("target status is required")
	}

	// Default waiter options for tags (can be overridden)
	defaultOpts := []waiter.Option{
		waiter.WithInterval(5 * time.Second),
		waiter.WithMaxWait(10 * time.Minute),
		waiter.WithBackoff(1.2, 30*time.Second),
	}

	// Merge user options (user options take precedence)
	opts := append(defaultOpts, cfg.WaiterOptions...)

	// Create the state check function
	checkState := func(ctx context.Context) (bool, error) {
		// Get current tag state
		tagResource, err := cfg.Client.Get(ctx, cfg.TagID)
		if err != nil {
			return false, fmt.Errorf("failed to get tag status: %w", err)
		}

		currentStatus := tagResource.Status

		// Check if we've reached the target status
		if currentStatus == cfg.TargetStatus {
			return true, nil
		}

		// If tag is in ERROR state and that's not our target, fail immediately
		if currentStatus == commonmodels.TagStatusError && cfg.TargetStatus != commonmodels.TagStatusError {
			return false, fmt.Errorf("tag entered ERROR state while waiting for %s", cfg.TargetStatus)
		}

		// Continue polling
		return false, nil
	}

	// Use the generic waiter framework
	return waiter.Wait(ctx, checkState, opts...)
}

// WaitForTagActive is a convenience function that waits for a tag to become ACTIVE.
func WaitForTagActive(ctx context.Context, client *tags.Client, tagID string, opts ...waiter.Option) error {
	return WaitForTagStatus(ctx, TagWaiterConfig{
		Client:        client,
		TagID:         tagID,
		TargetStatus:  commonmodels.TagStatusActive,
		WaiterOptions: opts,
	})
}

// WaitForTagAvailable is a convenience function that waits for a tag to become AVAILABLE.
func WaitForTagAvailable(ctx context.Context, client *tags.Client, tagID string, opts ...waiter.Option) error {
	return WaitForTagStatus(ctx, TagWaiterConfig{
		Client:        client,
		TagID:         tagID,
		TargetStatus:  commonmodels.TagStatusAvailable,
		WaiterOptions: opts,
	})
}
