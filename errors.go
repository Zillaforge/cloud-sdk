package cloudsdk

import (
	"github.com/Zillaforge/cloud-sdk/internal/types"
)

// SDKError is re-exported from internal/types for public API.
type SDKError = types.SDKError

// NewSDKError creates a new SDKError.
func NewSDKError(statusCode, errorCode int, message string, meta map[string]interface{}, cause error) *SDKError {
	return types.NewSDKError(statusCode, errorCode, message, meta, cause)
}

// NewNetworkError creates an SDKError for network-related failures.
func NewNetworkError(message string, cause error) *SDKError {
	return types.NewNetworkError(message, cause)
}

// NewTimeoutError creates an SDKError for timeout failures.
func NewTimeoutError(cause error) *SDKError {
	return types.NewTimeoutError(cause)
}

// NewCanceledError creates an SDKError for canceled requests.
func NewCanceledError(cause error) *SDKError {
	return types.NewCanceledError(cause)
}

// NewHTTPError creates an SDKError from an HTTP response without a parseable body.
func NewHTTPError(statusCode int, rawBody string) *SDKError {
	return types.NewHTTPError(statusCode, rawBody)
}
