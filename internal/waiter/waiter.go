// Package waiter provides a generic framework for polling resource state changes.
// It supports context-based cancellation, configurable intervals, backoff strategies,
// and maximum wait durations. This package is reusable across all cloud services.
package waiter

import (
	"context"
	"fmt"
	"time"
)

// StateCheckFunc is a function that checks if a resource has reached the desired state.
// It should return:
// - done=true if the target state is reached
// - err if an error occurred during the check
// - done=false, err=nil to continue polling
type StateCheckFunc func(ctx context.Context) (done bool, err error)

// Config holds the configuration for a waiter.
type Config struct {
	// Interval is the duration between state checks
	Interval time.Duration

	// MaxWait is the maximum duration to wait for the target state
	MaxWait time.Duration

	// BackoffMultiplier is applied to Interval after each failed check (1.0 = no backoff)
	BackoffMultiplier float64

	// MaxInterval caps the interval when using backoff
	MaxInterval time.Duration
}

// Option is a functional option for configuring a waiter.
type Option func(*Config)

// WithInterval sets the polling interval between state checks.
func WithInterval(interval time.Duration) Option {
	return func(c *Config) {
		c.Interval = interval
	}
}

// WithMaxWait sets the maximum time to wait for the target state.
func WithMaxWait(maxWait time.Duration) Option {
	return func(c *Config) {
		c.MaxWait = maxWait
	}
}

// WithBackoff enables exponential backoff between polling attempts.
// multiplier is applied to the interval after each attempt (e.g., 1.5 for 50% increase).
// maxInterval caps the backoff interval.
func WithBackoff(multiplier float64, maxInterval time.Duration) Option {
	return func(c *Config) {
		c.BackoffMultiplier = multiplier
		c.MaxInterval = maxInterval
	}
}

// DefaultConfig returns the default waiter configuration.
func DefaultConfig() *Config {
	return &Config{
		Interval:          2 * time.Second,
		MaxWait:           5 * time.Minute,
		BackoffMultiplier: 1.0, // No backoff by default
		MaxInterval:       30 * time.Second,
	}
}

// Wait polls a resource state using the provided StateCheckFunc until:
// - The state check returns done=true (success)
// - The context is canceled (returns ctx.Err())
// - The maximum wait duration is exceeded (returns ErrWaitTimeout)
// - An error occurs during state check (returns the error)
//
// The polling interval can optionally use exponential backoff if configured.
func Wait(ctx context.Context, checkState StateCheckFunc, opts ...Option) error {
	cfg := DefaultConfig()
	for _, opt := range opts {
		opt(cfg)
	}

	// Create a derived context with MaxWait timeout if needed
	// This ensures we don't exceed MaxWait even if parent context has no deadline
	waitCtx := ctx
	var cancel context.CancelFunc
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		waitCtx, cancel = context.WithTimeout(ctx, cfg.MaxWait)
		defer cancel()
	}

	currentInterval := cfg.Interval
	ticker := time.NewTicker(currentInterval)
	defer ticker.Stop()

	// Check state immediately before first tick
	done, err := checkState(waitCtx)
	if err != nil {
		return err
	}
	if done {
		return nil
	}

	for {
		select {
		case <-waitCtx.Done():
			// Context deadline or cancellation
			err := waitCtx.Err()
			// If this is our internal timeout context and the parent isn't done, return timeout error
			if err == context.DeadlineExceeded && ctx.Err() == nil {
				return ErrWaitTimeout
			}
			return err

		case <-ticker.C:
			// Check the state
			done, err := checkState(waitCtx)
			if err != nil {
				return err
			}
			if done {
				return nil
			}

			// Apply backoff if configured
			if cfg.BackoffMultiplier > 1.0 {
				currentInterval = time.Duration(float64(currentInterval) * cfg.BackoffMultiplier)
				if currentInterval > cfg.MaxInterval {
					currentInterval = cfg.MaxInterval
				}
				ticker.Reset(currentInterval)
			}
		}
	}
}

// ErrWaitTimeout is returned when the maximum wait duration is exceeded.
var ErrWaitTimeout = fmt.Errorf("wait timeout: maximum wait duration exceeded")
