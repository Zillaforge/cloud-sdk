package backoff

import (
	"testing"
	"time"
)

func TestDefaultStrategy(t *testing.T) {
	strategy := DefaultStrategy()

	if strategy.InitialInterval != 100*time.Millisecond {
		t.Errorf("expected InitialInterval 100ms, got %v", strategy.InitialInterval)
	}
	if strategy.MaxInterval != 5*time.Second {
		t.Errorf("expected MaxInterval 5s, got %v", strategy.MaxInterval)
	}
	if strategy.Multiplier != 2.0 {
		t.Errorf("expected Multiplier 2.0, got %v", strategy.Multiplier)
	}
	if strategy.MaxRetries != 3 {
		t.Errorf("expected MaxRetries 3, got %v", strategy.MaxRetries)
	}
	if !strategy.Jitter {
		t.Error("expected Jitter to be enabled")
	}
}

func TestStrategy_Duration(t *testing.T) {
	tests := []struct {
		name              string
		strategy          *Strategy
		attempt           int
		minDuration       time.Duration
		maxDuration       time.Duration
		expectWithinRange bool
	}{
		{
			name: "first attempt without jitter",
			strategy: &Strategy{
				InitialInterval: 100 * time.Millisecond,
				MaxInterval:     5 * time.Second,
				Multiplier:      2.0,
				MaxRetries:      3,
				Jitter:          false,
			},
			attempt:           0,
			minDuration:       100 * time.Millisecond,
			maxDuration:       100 * time.Millisecond,
			expectWithinRange: true,
		},
		{
			name: "second attempt without jitter",
			strategy: &Strategy{
				InitialInterval: 100 * time.Millisecond,
				MaxInterval:     5 * time.Second,
				Multiplier:      2.0,
				MaxRetries:      3,
				Jitter:          false,
			},
			attempt:           1,
			minDuration:       200 * time.Millisecond,
			maxDuration:       200 * time.Millisecond,
			expectWithinRange: true,
		},
		{
			name: "third attempt without jitter",
			strategy: &Strategy{
				InitialInterval: 100 * time.Millisecond,
				MaxInterval:     5 * time.Second,
				Multiplier:      2.0,
				MaxRetries:      3,
				Jitter:          false,
			},
			attempt:           2,
			minDuration:       400 * time.Millisecond,
			maxDuration:       400 * time.Millisecond,
			expectWithinRange: true,
		},
		{
			name: "first attempt with jitter",
			strategy: &Strategy{
				InitialInterval: 100 * time.Millisecond,
				MaxInterval:     5 * time.Second,
				Multiplier:      2.0,
				MaxRetries:      3,
				Jitter:          true,
			},
			attempt:           0,
			minDuration:       75 * time.Millisecond,  // 100ms * 0.75
			maxDuration:       125 * time.Millisecond, // 100ms * 1.25
			expectWithinRange: true,
		},
		{
			name: "capped at max interval",
			strategy: &Strategy{
				InitialInterval: 1 * time.Second,
				MaxInterval:     2 * time.Second,
				Multiplier:      2.0,
				MaxRetries:      5,
				Jitter:          false,
			},
			attempt:           3, // Would be 8s without cap
			minDuration:       2 * time.Second,
			maxDuration:       2 * time.Second,
			expectWithinRange: true,
		},
		{
			name: "negative attempt handled",
			strategy: &Strategy{
				InitialInterval: 100 * time.Millisecond,
				MaxInterval:     5 * time.Second,
				Multiplier:      2.0,
				MaxRetries:      3,
				Jitter:          false,
			},
			attempt:           -1,
			minDuration:       100 * time.Millisecond,
			maxDuration:       100 * time.Millisecond,
			expectWithinRange: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			duration := tt.strategy.Duration(tt.attempt)

			if tt.expectWithinRange {
				if duration < tt.minDuration || duration > tt.maxDuration {
					t.Errorf("duration %v not in expected range [%v, %v]", duration, tt.minDuration, tt.maxDuration)
				}
			}
		})
	}
}

func TestStrategy_ShouldRetry(t *testing.T) {
	strategy := &Strategy{
		MaxRetries: 3,
	}

	tests := []struct {
		attempt     int
		shouldRetry bool
	}{
		{0, true},
		{1, true},
		{2, true},
		{3, false},
		{4, false},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := strategy.ShouldRetry(tt.attempt)
			if result != tt.shouldRetry {
				t.Errorf("attempt %d: expected ShouldRetry=%v, got %v", tt.attempt, tt.shouldRetry, result)
			}
		})
	}
}

func TestIsRetryableStatusCode(t *testing.T) {
	tests := []struct {
		statusCode int
		retryable  bool
	}{
		{200, false},
		{400, false},
		{404, false},
		{429, true}, // Too Many Requests
		{500, false},
		{502, true}, // Bad Gateway
		{503, true}, // Service Unavailable
		{504, true}, // Gateway Timeout
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			result := IsRetryableStatusCode(tt.statusCode)
			if result != tt.retryable {
				t.Errorf("status %d: expected retryable=%v, got %v", tt.statusCode, tt.retryable, result)
			}
		})
	}
}

func TestIsRetryableMethod(t *testing.T) {
	tests := []struct {
		method    string
		retryable bool
	}{
		{"GET", true},
		{"HEAD", true},
		{"POST", false},
		{"PUT", false},
		{"PATCH", false},
		{"DELETE", false},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			result := IsRetryableMethod(tt.method)
			if result != tt.retryable {
				t.Errorf("method %s: expected retryable=%v, got %v", tt.method, tt.retryable, result)
			}
		})
	}
}
