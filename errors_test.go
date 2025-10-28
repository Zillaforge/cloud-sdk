package cloudsdk

import (
	"errors"
	"testing"
)

func TestSDKError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *SDKError
		expected string
	}{
		{
			name: "client-side error",
			err: &SDKError{
				StatusCode: 0,
				ErrorCode:  0,
				Message:    "network error",
			},
			expected: "SDK error: network error",
		},
		{
			name: "HTTP error with error code",
			err: &SDKError{
				StatusCode: 400,
				ErrorCode:  1001,
				Message:    "invalid request",
			},
			expected: "HTTP 400 (code 1001): invalid request",
		},
		{
			name: "HTTP error without error code",
			err: &SDKError{
				StatusCode: 500,
				ErrorCode:  0,
				Message:    "internal server error",
			},
			expected: "HTTP 500: internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.err.Error()
			if result != tt.expected {
				t.Errorf("expected '%s', got '%s'", tt.expected, result)
			}
		})
	}
}

func TestSDKError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	err := &SDKError{
		StatusCode: 500,
		Message:    "wrapper error",
		Cause:      cause,
	}

	unwrapped := err.Unwrap()
	if unwrapped != cause {
		t.Errorf("expected unwrapped error to be %v, got %v", cause, unwrapped)
	}
}

func TestSDKError_Is(t *testing.T) {
	err1 := &SDKError{
		StatusCode: 400,
		ErrorCode:  1001,
		Message:    "error 1",
	}

	err2 := &SDKError{
		StatusCode: 400,
		ErrorCode:  1001,
		Message:    "error 2", // Different message, same codes
	}

	err3 := &SDKError{
		StatusCode: 500,
		ErrorCode:  1001,
		Message:    "error 3",
	}

	if !err1.Is(err2) {
		t.Error("expected err1.Is(err2) to be true (same status and error codes)")
	}

	if err1.Is(err3) {
		t.Error("expected err1.Is(err3) to be false (different status codes)")
	}

	// Test with non-SDKError
	otherErr := errors.New("not an sdk error")
	if err1.Is(otherErr) {
		t.Error("expected err1.Is(otherErr) to be false")
	}
}

func TestNewSDKError(t *testing.T) {
	cause := errors.New("cause")
	meta := map[string]interface{}{"key": "value"}

	err := NewSDKError(400, 1001, "test message", meta, cause)

	if err.StatusCode != 400 {
		t.Errorf("expected StatusCode 400, got %d", err.StatusCode)
	}
	if err.ErrorCode != 1001 {
		t.Errorf("expected ErrorCode 1001, got %d", err.ErrorCode)
	}
	if err.Message != "test message" {
		t.Errorf("expected Message 'test message', got '%s'", err.Message)
	}
	if err.Meta["key"] != "value" {
		t.Errorf("expected Meta['key'] = 'value', got %v", err.Meta["key"])
	}
	if err.Cause != cause {
		t.Errorf("expected Cause to be %v, got %v", cause, err.Cause)
	}
}

func TestNewNetworkError(t *testing.T) {
	cause := errors.New("connection refused")
	err := NewNetworkError("connection failed", cause)

	if err.StatusCode != 0 {
		t.Errorf("expected StatusCode 0, got %d", err.StatusCode)
	}
	if err.ErrorCode != 0 {
		t.Errorf("expected ErrorCode 0, got %d", err.ErrorCode)
	}
	if err.Message != "network error: connection failed" {
		t.Errorf("expected message with network prefix, got '%s'", err.Message)
	}
	if err.Meta["category"] != "network" {
		t.Errorf("expected Meta category 'network', got %v", err.Meta["category"])
	}
	if err.Cause != cause {
		t.Errorf("expected Cause to be %v, got %v", cause, err.Cause)
	}
}

func TestNewTimeoutError(t *testing.T) {
	cause := errors.New("deadline exceeded")
	err := NewTimeoutError(cause)

	if err.StatusCode != 0 {
		t.Errorf("expected StatusCode 0, got %d", err.StatusCode)
	}
	if err.Message != "request timeout" {
		t.Errorf("expected 'request timeout', got '%s'", err.Message)
	}
	if err.Meta["category"] != "timeout" {
		t.Errorf("expected Meta category 'timeout', got %v", err.Meta["category"])
	}
	if err.Cause != cause {
		t.Errorf("expected Cause to be %v, got %v", cause, err.Cause)
	}
}

func TestNewCanceledError(t *testing.T) {
	cause := errors.New("context canceled")
	err := NewCanceledError(cause)

	if err.StatusCode != 0 {
		t.Errorf("expected StatusCode 0, got %d", err.StatusCode)
	}
	if err.Message != "request canceled" {
		t.Errorf("expected 'request canceled', got '%s'", err.Message)
	}
	if err.Meta["category"] != "canceled" {
		t.Errorf("expected Meta category 'canceled', got %v", err.Meta["category"])
	}
	if err.Cause != cause {
		t.Errorf("expected Cause to be %v, got %v", cause, err.Cause)
	}
}

func TestNewHTTPError(t *testing.T) {
	err := NewHTTPError(500, "Internal Server Error")

	if err.StatusCode != 500 {
		t.Errorf("expected StatusCode 500, got %d", err.StatusCode)
	}
	if err.ErrorCode != 0 {
		t.Errorf("expected ErrorCode 0, got %d", err.ErrorCode)
	}
	if err.Message != "HTTP 500" {
		t.Errorf("expected 'HTTP 500', got '%s'", err.Message)
	}
	if err.Meta["raw"] != "Internal Server Error" {
		t.Errorf("expected Meta raw body, got %v", err.Meta["raw"])
	}
	if err.Cause != nil {
		t.Errorf("expected nil Cause, got %v", err.Cause)
	}
}

func TestSDKError_ErrorsAs(t *testing.T) {
	err := NewSDKError(404, 1001, "not found", nil, nil)

	var sdkErr *SDKError
	if !errors.As(err, &sdkErr) {
		t.Error("expected errors.As to work with SDKError")
	}

	if sdkErr.StatusCode != 404 {
		t.Errorf("expected StatusCode 404, got %d", sdkErr.StatusCode)
	}
}
