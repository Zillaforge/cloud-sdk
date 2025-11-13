package cloudsdk

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"time"

	internalhttp "github.com/Zillaforge/cloud-sdk/internal/http"
	"github.com/Zillaforge/cloud-sdk/internal/types"
	iamprojects "github.com/Zillaforge/cloud-sdk/models/iam/projects"
	iam "github.com/Zillaforge/cloud-sdk/modules/iam/core"
	vps "github.com/Zillaforge/cloud-sdk/modules/vps/core"
	vrm "github.com/Zillaforge/cloud-sdk/modules/vrm/core"
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

// Project creates a project-scoped client for the given project ID or project code.
// If projectIDOrCode is a project ID, it creates the client directly.
// If it's a project code, it resolves the project ID by listing projects and finding the one with matching extra.iservice.projectSysCode.
// Returns an error if no matching project is found or if multiple projects match the code.
func (c *Client) Project(ctx context.Context, projectIDOrCode string) (*ProjectClient, error) {
	iamClient := c.IAM()

	// Try to get the project directly as projectID
	_, err := iamClient.Projects().Get(ctx, projectIDOrCode)
	if err == nil {
		// It's a valid projectID
		return &ProjectClient{
			client:    c,
			projectID: projectIDOrCode,
		}, nil
	}

	// Assume it's a projectCode, list all projects to find matching projectSysCode
	projects, listErr := iamClient.Projects().List(ctx, nil)
	if listErr != nil {
		return nil, fmt.Errorf("failed to list projects: %w", listErr)
	}

	var matchingProjects []*iamprojects.ProjectMembership
	for _, pm := range projects {
		if pm.Project != nil && pm.Project.Extra != nil {
			if iservice, ok := pm.Project.Extra["iservice"].(map[string]interface{}); ok {
				if sysCode, ok := iservice["projectSysCode"].(string); ok && sysCode == projectIDOrCode {
					matchingProjects = append(matchingProjects, pm)
				}
			}
		}
	}

	if len(matchingProjects) == 0 {
		return nil, fmt.Errorf("no project found with projectSysCode %s, please use projectID instead", projectIDOrCode)
	}

	if len(matchingProjects) > 1 {
		return nil, fmt.Errorf("multiple projects found with projectSysCode %s, please use projectID instead", projectIDOrCode)
	}

	projectID := matchingProjects[0].Project.ProjectID
	return &ProjectClient{
		client:    c,
		projectID: projectID,
	}, nil
}

// VPS returns a project-scoped VPS service client.
// All VPS operations will be performed within the context of the bound project.
func (pc *ProjectClient) VPS() *vps.Client {
	// Append /vps to baseURL for VPS service endpoints
	vpsBaseURL := pc.client.baseURL + "/vps"
	return vps.NewClient(vpsBaseURL, pc.client.token, pc.projectID, pc.client.httpClient, pc.client.logger)
}

// VRM returns a project-scoped VRM service client.
// All VRM operations will be performed within the context of the bound project.
func (pc *ProjectClient) VRM() *vrm.Client {
	// Append /vrm to baseURL for VRM service endpoints
	vrmBaseURL := pc.client.baseURL + "/vrm"
	return vrm.NewClient(vrmBaseURL, pc.client.token, pc.projectID, pc.client.httpClient, pc.client.logger)
}

// IAM returns a non-project-scoped IAM service client.
// IAM operations are global to the authenticated user and don't require a project context.
func (c *Client) IAM() *iam.Client {
	// Append /iam to baseURL for IAM service endpoints
	iamBaseURL := c.baseURL + "/iam"

	// Create internal HTTP client with retry and error handling
	baseClient := internalhttp.NewClient(iamBaseURL, c.token, c.httpClient, c.logger)

	return iam.NewClient(baseClient)
}

// BaseURL returns the configured base URL.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// HTTPClient returns the underlying HTTP client.
func (c *Client) HTTPClient() *http.Client {
	return c.httpClient
}
