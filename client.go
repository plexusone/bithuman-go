// Package bithuman provides a Go client for the bitHuman API.
//
// bitHuman is a real-time avatar animation platform that creates digital avatars
// with lip-sync to audio at 25 FPS. Audio in, animated video out.
//
// The client wraps the ogen-generated API client with a higher-level interface
// that handles authentication and provides convenient methods for common operations.
//
// # Quick Start
//
//	client, err := bithuman.NewClient(bithuman.WithAPIKey("your-api-key"))
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Create a session with an agent
//	session, err := client.Sessions().Create(ctx, &api.CreateSessionRequest{
//	    AgentId: "agent-id",
//	})
//
// # Environment Variables
//
// If no API key is provided, the client will look for BITHUMAN_API_KEY in the environment.
package bithuman

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/plexusone/bithuman-go/api"
)

// Version is the SDK version.
const Version = "0.1.0"

// DefaultBaseURL is the default bitHuman API base URL.
const DefaultBaseURL = "https://api.bithuman.ai"

// Client is the main bitHuman client for interacting with the API.
type Client struct {
	apiClient *api.Client
	apiKey    string
	baseURL   string

	// Domain-based service accessors
	agentsSvc   *AgentsService
	sessionsSvc *SessionsService
	ttsSvc      *TTSService
	videosSvc   *VideosService
	filesSvc    *FilesService
	billingSvc  *BillingService
	webhooksSvc *WebhooksService
}

// NewClient creates a new bitHuman client with the given options.
func NewClient(opts ...Option) (*Client, error) {
	options := defaultClientOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Try environment variable if API key not set
	if options.apiKey == "" {
		options.apiKey = os.Getenv("BITHUMAN_API_KEY")
	}

	// Create HTTP client
	httpClient := options.httpClient
	if httpClient == nil {
		httpClient = &http.Client{
			Timeout: options.timeout,
		}
	}

	// Create security source for authentication
	securitySource := &apiKeySecuritySource{apiKey: options.apiKey}

	// Create the ogen client
	apiClient, err := api.NewClient(
		options.baseURL,
		securitySource,
		api.WithClient(&sdkHTTPClient{client: httpClient}),
	)
	if err != nil {
		return nil, err
	}

	c := &Client{
		apiClient: apiClient,
		apiKey:    options.apiKey,
		baseURL:   options.baseURL,
	}

	// Initialize domain-based services
	c.agentsSvc = &AgentsService{client: apiClient}
	c.sessionsSvc = &SessionsService{client: apiClient}
	c.ttsSvc = &TTSService{client: apiClient}
	c.videosSvc = &VideosService{client: apiClient}
	c.filesSvc = &FilesService{client: apiClient}
	c.billingSvc = &BillingService{client: apiClient}
	c.webhooksSvc = &WebhooksService{client: apiClient}

	return c, nil
}

// apiKeySecuritySource implements api.SecuritySource for API key authentication.
type apiKeySecuritySource struct {
	apiKey string
}

// ApiSecret implements api.SecuritySource.
func (s *apiKeySecuritySource) ApiSecret(ctx context.Context, operationName api.OperationName) (api.ApiSecret, error) {
	return api.ApiSecret{APIKey: s.apiKey}, nil
}

// sdkHTTPClient wraps an http.Client to add SDK version headers.
type sdkHTTPClient struct {
	client *http.Client
}

// Do implements ht.Client interface.
func (c *sdkHTTPClient) Do(req *http.Request) (*http.Response, error) {
	// Add SDK version headers
	req.Header.Set("X-BitHuman-SDK-Version", Version)
	req.Header.Set("X-BitHuman-SDK-Lang", "go")

	// #nosec G704 -- URL is controlled by SDK configuration (BaseURL), not untrusted input
	return c.client.Do(req)
}

// API returns the underlying ogen-generated API client for advanced usage.
// Use this when you need access to API endpoints not covered by the
// high-level wrapper methods.
func (c *Client) API() *api.Client {
	return c.apiClient
}

// Agents returns the agents service for avatar agent management.
func (c *Client) Agents() *AgentsService {
	return c.agentsSvc
}

// Sessions returns the sessions service for real-time conversation sessions.
func (c *Client) Sessions() *SessionsService {
	return c.sessionsSvc
}

// TTS returns the TTS service for text-to-speech synthesis.
func (c *Client) TTS() *TTSService {
	return c.ttsSvc
}

// Videos returns the videos service for talking video generation.
func (c *Client) Videos() *VideosService {
	return c.videosSvc
}

// Files returns the files service for file uploads.
func (c *Client) Files() *FilesService {
	return c.filesSvc
}

// Billing returns the billing service for account balance.
func (c *Client) Billing() *BillingService {
	return c.billingSvc
}

// Webhooks returns the webhooks service for event notifications.
func (c *Client) Webhooks() *WebhooksService {
	return c.webhooksSvc
}

// APIKey returns the API key used by the client.
func (c *Client) APIKey() string {
	return c.apiKey
}

// BaseURL returns the base URL used by the client.
func (c *Client) BaseURL() string {
	return c.baseURL
}

// Validate checks if the API credentials are valid.
func (c *Client) Validate(ctx context.Context) (*api.ValidateResponse, error) {
	res, err := c.apiClient.Validate(ctx)
	if err != nil {
		return nil, err
	}
	if v, ok := res.(*api.ValidateResponse); ok {
		return v, nil
	}
	return nil, fmt.Errorf("unexpected response type: %T", res)
}

// clientOptions holds the options for creating a Client.
type clientOptions struct {
	apiKey     string
	baseURL    string
	httpClient *http.Client
	timeout    time.Duration
}

func defaultClientOptions() *clientOptions {
	return &clientOptions{
		baseURL: DefaultBaseURL,
		timeout: 120 * time.Second,
	}
}

// Option is a functional option for configuring the Client.
type Option func(*clientOptions)

// WithAPIKey sets the API key for authentication.
func WithAPIKey(apiKey string) Option {
	return func(o *clientOptions) {
		o.apiKey = apiKey
	}
}

// WithBaseURL sets the API base URL.
func WithBaseURL(baseURL string) Option {
	return func(o *clientOptions) {
		o.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) Option {
	return func(o *clientOptions) {
		o.httpClient = client
	}
}

// WithTimeout sets the request timeout.
func WithTimeout(timeout time.Duration) Option {
	return func(o *clientOptions) {
		o.timeout = timeout
	}
}
