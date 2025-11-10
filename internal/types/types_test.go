package types

import (
	"errors"
	"testing"
)

func TestSDKError_Error(t *testing.T) {
	tests := []struct {
		name     string
		sdkError *SDKError
		expected string
	}{
		{
			name: "client-side error (StatusCode 0)",
			sdkError: &SDKError{
				StatusCode: 0,
				Message:    "client error",
			},
			expected: "SDK error: client error",
		},
		{
			name: "HTTP error with ErrorCode",
			sdkError: &SDKError{
				StatusCode: 404,
				ErrorCode:  1001,
				Message:    "not found",
			},
			expected: "HTTP 404 (code 1001): not found",
		},
		{
			name: "HTTP error without ErrorCode",
			sdkError: &SDKError{
				StatusCode: 500,
				ErrorCode:  0,
				Message:    "internal server error",
			},
			expected: "HTTP 500: internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.sdkError.Error(); got != tt.expected {
				t.Errorf("SDKError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestSDKError_Unwrap(t *testing.T) {
	cause := errors.New("underlying error")
	sdkError := &SDKError{
		Cause: cause,
	}

	if got := sdkError.Unwrap(); got != cause {
		t.Errorf("SDKError.Unwrap() = %v, want %v", got, cause)
	}
}

func TestSDKError_Is(t *testing.T) {
	sdkError1 := &SDKError{
		StatusCode: 404,
		ErrorCode:  1001,
	}
	sdkError2 := &SDKError{
		StatusCode: 404,
		ErrorCode:  1001,
	}
	sdkError3 := &SDKError{
		StatusCode: 500,
		ErrorCode:  1002,
	}
	regularError := errors.New("regular error")

	if !errors.Is(sdkError1, sdkError2) {
		t.Errorf("expected sdkError1 to be equal to sdkError2")
	}

	if errors.Is(sdkError1, sdkError3) {
		t.Errorf("expected sdkError1 not to be equal to sdkError3")
	}

	if errors.Is(sdkError1, regularError) {
		t.Errorf("expected sdkError1 not to be equal to regularError")
	}
}

func TestNewSDKError(t *testing.T) {
	meta := map[string]interface{}{"key": "value"}
	cause := errors.New("cause")
	sdkError := NewSDKError(400, 2001, "bad request", meta, cause)

	if sdkError.StatusCode != 400 {
		t.Errorf("StatusCode = %v, want 400", sdkError.StatusCode)
	}
	if sdkError.ErrorCode != 2001 {
		t.Errorf("ErrorCode = %v, want 2001", sdkError.ErrorCode)
	}
	if sdkError.Message != "bad request" {
		t.Errorf("Message = %v, want 'bad request'", sdkError.Message)
	}
	if sdkError.Meta["key"] != "value" {
		t.Errorf("Meta = %v, want key 'value'", sdkError.Meta)
	}
	if sdkError.Cause != cause {
		t.Errorf("Cause = %v, want %v", sdkError.Cause, cause)
	}
}

func TestNewNetworkError(t *testing.T) {
	cause := errors.New("network cause")
	sdkError := NewNetworkError("connection failed", cause)

	if sdkError.StatusCode != 0 {
		t.Errorf("StatusCode = %v, want 0", sdkError.StatusCode)
	}
	if sdkError.ErrorCode != 0 {
		t.Errorf("ErrorCode = %v, want 0", sdkError.ErrorCode)
	}
	if sdkError.Message != "network error: connection failed" {
		t.Errorf("Message = %v, want 'network error: connection failed'", sdkError.Message)
	}
	if sdkError.Meta["category"] != "network" {
		t.Errorf("Meta category = %v, want 'network'", sdkError.Meta["category"])
	}
	if sdkError.Cause != cause {
		t.Errorf("Cause = %v, want %v", sdkError.Cause, cause)
	}
}

func TestNewTimeoutError(t *testing.T) {
	cause := errors.New("timeout cause")
	sdkError := NewTimeoutError(cause)

	if sdkError.StatusCode != 0 {
		t.Errorf("StatusCode = %v, want 0", sdkError.StatusCode)
	}
	if sdkError.ErrorCode != 0 {
		t.Errorf("ErrorCode = %v, want 0", sdkError.ErrorCode)
	}
	if sdkError.Message != "request timeout" {
		t.Errorf("Message = %v, want 'request timeout'", sdkError.Message)
	}
	if sdkError.Meta["category"] != "timeout" {
		t.Errorf("Meta category = %v, want 'timeout'", sdkError.Meta["category"])
	}
	if sdkError.Cause != cause {
		t.Errorf("Cause = %v, want %v", sdkError.Cause, cause)
	}
}

func TestNewCanceledError(t *testing.T) {
	cause := errors.New("canceled cause")
	sdkError := NewCanceledError(cause)

	if sdkError.StatusCode != 0 {
		t.Errorf("StatusCode = %v, want 0", sdkError.StatusCode)
	}
	if sdkError.ErrorCode != 0 {
		t.Errorf("ErrorCode = %v, want 0", sdkError.ErrorCode)
	}
	if sdkError.Message != "request canceled" {
		t.Errorf("Message = %v, want 'request canceled'", sdkError.Message)
	}
	if sdkError.Meta["category"] != "canceled" {
		t.Errorf("Meta category = %v, want 'canceled'", sdkError.Meta["category"])
	}
	if sdkError.Cause != cause {
		t.Errorf("Cause = %v, want %v", sdkError.Cause, cause)
	}
}

func TestNewHTTPError(t *testing.T) {
	sdkError := NewHTTPError(502, "bad gateway")

	if sdkError.StatusCode != 502 {
		t.Errorf("StatusCode = %v, want 502", sdkError.StatusCode)
	}
	if sdkError.ErrorCode != 0 {
		t.Errorf("ErrorCode = %v, want 0", sdkError.ErrorCode)
	}
	if sdkError.Message != "HTTP 502" {
		t.Errorf("Message = %v, want 'HTTP 502'", sdkError.Message)
	}
	if sdkError.Meta["raw"] != "bad gateway" {
		t.Errorf("Meta raw = %v, want 'bad gateway'", sdkError.Meta["raw"])
	}
	if sdkError.Cause != nil {
		t.Errorf("Cause = %v, want nil", sdkError.Cause)
	}
}
