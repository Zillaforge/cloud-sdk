// Package http provides internal HTTP client utilities for the SDK.
package http

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/backoff"
	"github.com/Zillaforge/cloud-sdk/internal/types"
)

// Client wraps the standard HTTP client with retry, timeout, and error handling.
type Client struct {
	baseURL       string
	token         string
	httpClient    *http.Client
	logger        types.Logger
	retryStrategy *backoff.Strategy
}

// NewClient creates a new internal HTTP client.
func NewClient(baseURL, token string, httpClient *http.Client, logger types.Logger) *Client {
	return &Client{
		baseURL:       baseURL,
		token:         token,
		httpClient:    httpClient,
		logger:        logger,
		retryStrategy: backoff.DefaultStrategy(),
	}
}

// Request represents an HTTP request to be executed.
type Request struct {
	Method  string
	Path    string
	Body    interface{}
	Headers map[string]string
}

// Do executes an HTTP request with retry logic and error handling.
// The context can override the default timeout.
func (c *Client) Do(ctx context.Context, req *Request, result interface{}) error {
	// Apply default timeout if context doesn't have a deadline
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.httpClient.Timeout)
		defer cancel()
	}

	attempt := 0

	for {
		err := c.doOnce(ctx, req, result)
		if err == nil {
			return nil
		}

		// Check if we should retry
		var sdkErr *types.SDKError
		if !errors.As(err, &sdkErr) {
			return err // Not an SDKError, don't retry
		}

		// Only retry for specific status codes and safe methods
		if !backoff.IsRetryableStatusCode(sdkErr.StatusCode) || !backoff.IsRetryableMethod(req.Method) {
			return err
		}

		// Check if we can retry
		if !c.retryStrategy.ShouldRetry(attempt) {
			if c.logger != nil {
				c.logger.Debug("max retries reached", "method", req.Method, "path", req.Path, "attempts", attempt+1)
			}
			return err
		}

		// Calculate backoff duration
		duration := c.retryStrategy.Duration(attempt)
		if c.logger != nil {
			c.logger.Debug("retrying request", "method", req.Method, "path", req.Path, "attempt", attempt+1, "backoff", duration)
		}

		// Wait with context cancellation support
		select {
		case <-time.After(duration):
			attempt++
		case <-ctx.Done():
			return types.NewCanceledError(ctx.Err())
		}
	}
}

// doOnce executes a single HTTP request without retry.
func (c *Client) doOnce(ctx context.Context, req *Request, result interface{}) error {
	// Build full URL
	url := c.baseURL + req.Path

	// Marshal request body if present
	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		if err != nil {
			return types.NewSDKError(0, 0, "failed to marshal request body", nil, err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, req.Method, url, bodyReader)
	if err != nil {
		return types.NewSDKError(0, 0, "failed to create request", nil, err)
	}

	// Set headers
	httpReq.Header.Set("Authorization", "Bearer "+c.token)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "application/json")
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// Execute request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		// Check for context errors
		if ctx.Err() == context.DeadlineExceeded {
			return types.NewTimeoutError(err)
		}
		if ctx.Err() == context.Canceled {
			return types.NewCanceledError(err)
		}
		return types.NewNetworkError(err.Error(), err)
	}
	defer resp.Body.Close()

	// Read response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return types.NewSDKError(resp.StatusCode, 0, "failed to read response body", nil, err)
	}

	// Handle error responses
	if resp.StatusCode >= 400 {
		return c.parseErrorResponse(resp.StatusCode, bodyBytes)
	}

	// Parse success response
	if result != nil && len(bodyBytes) > 0 {
		if err := json.Unmarshal(bodyBytes, result); err != nil {
			return types.NewSDKError(resp.StatusCode, 0, "failed to parse response", nil, err)
		}
	}

	return nil
}

// ErrorResponse represents the standard error response format.
type ErrorResponse struct {
	ErrorCode int                    `json:"errorCode"`
	Message   string                 `json:"message"`
	Meta      map[string]interface{} `json:"meta,omitempty"`
}

// parseErrorResponse attempts to parse an error response from the API.
func (c *Client) parseErrorResponse(statusCode int, body []byte) error {
	var errResp ErrorResponse
	if err := json.Unmarshal(body, &errResp); err != nil {
		// Failed to parse as structured error, return raw body
		return types.NewHTTPError(statusCode, string(body))
	}

	return types.NewSDKError(statusCode, errResp.ErrorCode, errResp.Message, errResp.Meta, nil)
}
