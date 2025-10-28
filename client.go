package cloudsdk

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/Zillaforge/cloud-sdk/internal/types"
	"github.com/Zillaforge/cloud-sdk/modules/vps"
)

// Logger defines the interface for logging SDK operations.
type Logger = types.Logger

// Client is the main entry point for the Cloud SDK.
// It manages authentication, base URL, and HTTP client configuration.
type Client struct {
	baseURL    string
	token      string
	httpClient *http.Client
	logger     Logger
}

// ClientOption is a functional option for configuring the Client.
type ClientOption func(*Client)

// WithLogger sets a custom logger for the client.
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithTimeout sets a custom default timeout for all requests.
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// New creates a new Cloud SDK client.
// baseURL must be a valid URL with scheme (e.g., "https://api.example.com").
// token must be a non-empty bearer token for authentication.
func New(baseURL, token string, opts ...ClientOption) (*Client, error) {
	// Validate base URL
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}
	if parsedURL.Scheme == "" {
		return nil, fmt.Errorf("base URL must include scheme (e.g., https://)")
	}

	// Validate token
	if token == "" {
		return nil, fmt.Errorf("token cannot be empty")
	}

	// Create client with defaults
	client := &Client{
		baseURL: baseURL,
		token:   token,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // Default 30s timeout
		},
	}

	// Apply options
	for _, opt := range opts {
		opt(client)
	}

	return client, nil
}

// NewClient is a convenience function that creates a new client without error handling.
// This is primarily for testing purposes. Use New() for production code.
func NewClient(baseURL, token string) *Client {
	client, _ := New(baseURL, token)
	return client
}

// ProjectClient provides access to project-scoped service clients.
type ProjectClient struct {
	client    *Client
	projectID string
}

// Project creates a project-scoped client for the given project ID.
// All subsequent service client calls will be scoped to this project.
func (c *Client) Project(projectID string) *ProjectClient {
	return &ProjectClient{
		client:    c,
		projectID: projectID,
	}
}

// VPS returns a project-scoped VPS service client.
// All VPS operations will be performed within the context of the bound project.
func (pc *ProjectClient) VPS() *vps.Client {
	// Append /vps to baseURL for VPS service endpoints
	vpsBaseURL := pc.client.baseURL + "/vps"
	return vps.NewClient(vpsBaseURL, pc.client.token, pc.projectID, pc.client.httpClient, pc.client.logger)
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// HTTPClient returns the underlying HTTP client.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}
