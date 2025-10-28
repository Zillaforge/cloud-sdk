// Package backoff provides exponential backoff with jitter for retrying operations.
package backoff

import (
	"math"
	"math/rand"
	"time"
)

// Strategy defines the retry backoff configuration.
type Strategy struct {
	// InitialInterval is the initial wait time before the first retry
	InitialInterval time.Duration

	// MaxInterval is the maximum wait time between retries
	MaxInterval time.Duration

	// Multiplier is the factor by which the interval increases each retry
	Multiplier float64

	// MaxRetries is the maximum number of retry attempts
	MaxRetries int

	// Jitter adds randomization to prevent thundering herd
	Jitter bool
}

// DefaultStrategy returns a reasonable default retry strategy.
// - Initial interval: 100ms
// - Max interval: 5s
// - Multiplier: 2.0 (exponential)
// - Max retries: 3
// - Jitter: enabled
func DefaultStrategy() *Strategy {
	return &Strategy{
		InitialInterval: 100 * time.Millisecond,
		MaxInterval:     5 * time.Second,
		Multiplier:      2.0,
		MaxRetries:      3,
		Jitter:          true,
	}
}

// Duration calculates the backoff duration for the given attempt number.
// attempt is zero-indexed (0 = first retry).
func (s *Strategy) Duration(attempt int) time.Duration {
	if attempt < 0 {
		attempt = 0
	}

	// Calculate exponential backoff
	duration := float64(s.InitialInterval) * math.Pow(s.Multiplier, float64(attempt))

	// Cap at max interval
	if duration > float64(s.MaxInterval) {
		duration = float64(s.MaxInterval)
	}

	// Apply jitter if enabled (±25% randomization)
	if s.Jitter {
		//nolint:gosec // Non-cryptographic randomness is acceptable for backoff jitter
		jitterFactor := 0.75 + rand.Float64()*0.5 // Range: 0.75 to 1.25
		duration *= jitterFactor
	}

	return time.Duration(duration)
}

// ShouldRetry determines if another retry should be attempted.
func (s *Strategy) ShouldRetry(attempt int) bool {
	return attempt < s.MaxRetries
}

// IsRetryableStatusCode checks if an HTTP status code is retryable.
// Retryable codes: 429 (Too Many Requests), 502 (Bad Gateway),
// 503 (Service Unavailable), 504 (Gateway Timeout).
func IsRetryableStatusCode(statusCode int) bool {
	return statusCode == 429 || statusCode == 502 || statusCode == 503 || statusCode == 504
}

// IsRetryableMethod checks if an HTTP method is safe to retry.
// Only GET and HEAD are considered safe for automatic retry.
func IsRetryableMethod(method string) bool {
	return method == "GET" || method == "HEAD"
}
