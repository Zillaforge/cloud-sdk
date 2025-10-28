package waiter

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestWait_ImmediateSuccess(t *testing.T) {
	ctx := context.Background()
	checkCount := 0

	checkState := func(_ context.Context) (bool, error) {
		checkCount++
		return true, nil // Target state reached immediately
	}

	err := Wait(ctx, checkState, WithInterval(100*time.Millisecond))
	if err != nil {
		t.Errorf("Wait() error = %v, want nil", err)
	}
	if checkCount != 1 {
		t.Errorf("checkState called %d times, want 1", checkCount)
	}
}

func TestWait_SuccessAfterPolling(t *testing.T) {
	ctx := context.Background()
	checkCount := 0
	targetChecks := 3

	checkState := func(_ context.Context) (bool, error) {
		checkCount++
		return checkCount >= targetChecks, nil
	}

	start := time.Now()
	err := Wait(ctx, checkState, WithInterval(50*time.Millisecond), WithMaxWait(5*time.Second))
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Wait() error = %v, want nil", err)
	}
	if checkCount != targetChecks {
		t.Errorf("checkState called %d times, want %d", checkCount, targetChecks)
	}

	// Should have waited at least 2 intervals (after initial check, wait for 2 more)
	minDuration := 2 * 50 * time.Millisecond
	if duration < minDuration {
		t.Errorf("Wait() duration %v, want at least %v", duration, minDuration)
	}
}

func TestWait_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	callCount := 0
	checkState := func(_ context.Context) (bool, error) {
		callCount++
		return false, nil
	}

	// Cancel after a short delay
	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel()
	}()

	err := Wait(ctx, checkState, WithInterval(50*time.Millisecond))
	if !errors.Is(err, context.Canceled) {
		t.Errorf("expected context.Canceled, got %v", err)
	}
	if callCount < 2 {
		t.Errorf("expected at least 2 calls before cancellation, got %d", callCount)
	}
}

func TestWait_MaxWaitTimeout(t *testing.T) {
	ctx := context.Background()
	checkCount := 0

	checkState := func(_ context.Context) (bool, error) {
		checkCount++
		return false, nil // Never reaches target state
	}

	start := time.Now()
	err := Wait(ctx, checkState,
		WithInterval(50*time.Millisecond),
		WithMaxWait(200*time.Millisecond))
	duration := time.Since(start)

	if !errors.Is(err, ErrWaitTimeout) {
		t.Errorf("Wait() error = %v, want ErrWaitTimeout", err)
	}

	// Should have timed out around 200ms
	if duration < 200*time.Millisecond {
		t.Errorf("Wait() duration %v, want at least 200ms", duration)
	}
	if duration > 400*time.Millisecond {
		t.Errorf("Wait() duration %v, want less than 400ms", duration)
	}
}

func TestWait_ContextDeadlineEarlier(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	checkState := func(_ context.Context) (bool, error) {
		return false, nil
	}

	start := time.Now()
	err := Wait(ctx, checkState,
		WithInterval(50*time.Millisecond),
		WithMaxWait(5*time.Second)) // MaxWait longer, but context deadline wins
	duration := time.Since(start)

	if !errors.Is(err, context.DeadlineExceeded) {
		t.Errorf("Wait() error = %v, want context.DeadlineExceeded", err)
	}

	// Should have timed out around 100ms (context deadline)
	if duration < 100*time.Millisecond {
		t.Errorf("Wait() duration %v, want at least 100ms", duration)
	}
	if duration > 300*time.Millisecond {
		t.Errorf("Wait() duration %v, want less than 300ms", duration)
	}
}

func TestWait_StateCheckError(t *testing.T) {
	ctx := context.Background()
	checkCount := 0
	testErr := errors.New("state check failed")

	checkState := func(_ context.Context) (bool, error) {
		checkCount++
		if checkCount == 2 {
			return false, testErr // Return error on second check
		}
		return false, nil
	}

	err := Wait(ctx, checkState, WithInterval(50*time.Millisecond))
	if !errors.Is(err, testErr) {
		t.Errorf("Wait() error = %v, want %v", err, testErr)
	}
	if checkCount != 2 {
		t.Errorf("checkState called %d times, want 2", checkCount)
	}
}

func TestWait_WithBackoff(t *testing.T) {
	ctx := context.Background()
	checkCount := 0
	intervals := []time.Duration{}
	lastCheck := time.Now()

	checkState := func(_ context.Context) (bool, error) {
		checkCount++
		now := time.Now()
		if checkCount > 1 {
			intervals = append(intervals, now.Sub(lastCheck))
		}
		lastCheck = now

		// Succeed after 4 checks to observe backoff
		return checkCount >= 4, nil
	}

	err := Wait(ctx, checkState,
		WithInterval(50*time.Millisecond),
		WithBackoff(1.5, 200*time.Millisecond), // 50, 75, 112.5, capped at 200
		WithMaxWait(5*time.Second))

	if err != nil {
		t.Errorf("Wait() error = %v, want nil", err)
	}
	if checkCount != 4 {
		t.Errorf("checkState called %d times, want 4", checkCount)
	}

	// Verify backoff is applied (each interval should be longer than the previous)
	if len(intervals) >= 2 {
		if intervals[1] <= intervals[0] {
			t.Errorf("Expected backoff: interval[1]=%v should be > interval[0]=%v",
				intervals[1], intervals[0])
		}
	}
}

func TestWait_BackoffCappedAtMaxInterval(t *testing.T) {
	ctx := context.Background()
	checkCount := 0
	maxInterval := 100 * time.Millisecond

	checkState := func(_ context.Context) (bool, error) {
		checkCount++
		return checkCount >= 5, nil // Let it backoff several times
	}

	start := time.Now()
	err := Wait(ctx, checkState,
		WithInterval(20*time.Millisecond),
		WithBackoff(2.0, maxInterval), // 20, 40, 80, 100 (capped), 100 (capped)
		WithMaxWait(5*time.Second))
	duration := time.Since(start)

	if err != nil {
		t.Errorf("Wait() error = %v, want nil", err)
	}

	// With capping, the total duration should be reasonable
	// 20 + 40 + 80 + 100 = 240ms minimum
	if duration < 240*time.Millisecond {
		t.Errorf("Wait() duration %v, want at least 240ms", duration)
	}
	// But not exponentially growing
	if duration > 1*time.Second {
		t.Errorf("Wait() duration %v, want less than 1s (backoff should be capped)", duration)
	}
}

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.Interval != 2*time.Second {
		t.Errorf("DefaultConfig().Interval = %v, want 2s", cfg.Interval)
	}
	if cfg.MaxWait != 5*time.Minute {
		t.Errorf("DefaultConfig().MaxWait = %v, want 5m", cfg.MaxWait)
	}
	if cfg.BackoffMultiplier != 1.0 {
		t.Errorf("DefaultConfig().BackoffMultiplier = %v, want 1.0", cfg.BackoffMultiplier)
	}
	if cfg.MaxInterval != 30*time.Second {
		t.Errorf("DefaultConfig().MaxInterval = %v, want 30s", cfg.MaxInterval)
	}
}

func TestWait_NoBackoffByDefault(t *testing.T) {
	var intervals []time.Duration
	var lastTime time.Time
	callCount := 0

	checkState := func(_ context.Context) (bool, error) {
		now := time.Now()
		if callCount > 0 {
			intervals = append(intervals, now.Sub(lastTime))
		}
		lastTime = now
		callCount++

		if callCount >= 5 {
			return true, nil
		}
		return false, nil
	}

	err := Wait(
		context.Background(),
		checkState,
		WithInterval(50*time.Millisecond),
		// No WithBackoff - should use constant interval
	)

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}

	// All intervals should be roughly the same (~50ms)
	for i, interval := range intervals {
		if interval < 40*time.Millisecond || interval > 80*time.Millisecond {
			t.Errorf("interval %d (%v) should be ~50ms", i, interval)
		}
	}
}
