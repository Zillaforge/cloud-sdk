// Package users provides the client for IAM user operations.
package users

import (
	"context"
	"fmt"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/models/iam/users"
)

// Client handles user operations for the IAM API.
type Client struct {
	baseClient *internalhttp.Client
	basePath   string
}

// NewClient creates a new users client with the provided HTTP client.
func NewClient(baseClient *internalhttp.Client, basePath string) *Client {
	return &Client{
		baseClient: baseClient,
		basePath:   basePath,
	}
}

// Get retrieves the current authenticated user's information.
func (c *Client) Get(ctx context.Context) (*users.User, error) {
	var response users.GetUserResponse

	req := &internalhttp.Request{
		Method: "GET",
		Path:   c.basePath + "user",
	}

	if err := c.baseClient.Do(ctx, req, &response); err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return response.User, nil
}
