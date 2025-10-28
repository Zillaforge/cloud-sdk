// Package types provides common types used across the SDK.
package types

import (
	"errors"
	"fmt"
)

// Logger defines the interface for logging SDK operations.
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
}

// SDKError represents a structured error returned by the SDK.
// It includes HTTP status codes, error codes from the API, and additional metadata.
type SDKError struct {
	// StatusCode is the HTTP status code (0 for client-side errors)
	StatusCode int

	// ErrorCode is the error code from the API response (0 if not available)
	ErrorCode int

	// Message is a human-readable error message
	Message string

	// Meta contains additional context about the error
	Meta map[string]interface{}

	// Cause is the underlying error that caused this error
	Cause error
}

// Error implements the error interface.
func (e *SDKError) Error() string {
	if e.StatusCode == 0 {
		return fmt.Sprintf("SDK error: %s", e.Message)
	}
	if e.ErrorCode != 0 {
		return fmt.Sprintf("HTTP %d (code %d): %s", e.StatusCode, e.ErrorCode, e.Message)
	}
	return fmt.Sprintf("HTTP %d: %s", e.StatusCode, e.Message)
}

// Unwrap returns the underlying cause error.
func (e *SDKError) Unwrap() error {
	return e.Cause
}

// Is allows error comparison using errors.Is.
func (e *SDKError) Is(target error) bool {
	var sdkErr *SDKError
	if errors.As(target, &sdkErr) {
		return e.StatusCode == sdkErr.StatusCode && e.ErrorCode == sdkErr.ErrorCode
	}
	return false
}

// NewSDKError creates a new SDKError.
func NewSDKError(statusCode, errorCode int, message string, meta map[string]interface{}, cause error) *SDKError {
	return &SDKError{
		StatusCode: statusCode,
		ErrorCode:  errorCode,
		Message:    message,
		Meta:       meta,
		Cause:      cause,
	}
}

// NewNetworkError creates an SDKError for network-related failures.
func NewNetworkError(message string, cause error) *SDKError {
	return &SDKError{
		StatusCode: 0,
		ErrorCode:  0,
		Message:    fmt.Sprintf("network error: %s", message),
		Meta:       map[string]interface{}{"category": "network"},
		Cause:      cause,
	}
}

// NewTimeoutError creates an SDKError for timeout failures.
func NewTimeoutError(cause error) *SDKError {
	return &SDKError{
		StatusCode: 0,
		ErrorCode:  0,
		Message:    "request timeout",
		Meta:       map[string]interface{}{"category": "timeout"},
		Cause:      cause,
	}
}

// NewCanceledError creates an SDKError for canceled requests.
func NewCanceledError(cause error) *SDKError {
	return &SDKError{
		StatusCode: 0,
		ErrorCode:  0,
		Message:    "request canceled",
		Meta:       map[string]interface{}{"category": "canceled"},
		Cause:      cause,
	}
}

// NewHTTPError creates an SDKError from an HTTP response without a parseable body.
func NewHTTPError(statusCode int, rawBody string) *SDKError {
	return &SDKError{
		StatusCode: statusCode,
		ErrorCode:  0,
		Message:    fmt.Sprintf("HTTP %d", statusCode),
		Meta:       map[string]interface{}{"raw": rawBody},
		Cause:      nil,
	}
}
