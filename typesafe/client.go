package typesafe

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const defaultBaseURL = "https://api.typesafe.ai/v1/systemone"

// Client is a TypeSafe Jev API client.
type Client struct {
	APIKey  string
	BaseURL string
	Model   string // Default model to use (e.g., "jev-latest")

	httpClient *http.Client
}

// ClientOption configures a Client.
type ClientOption func(*Client)

// WithAPIKey sets the TypeSafe API key.
func WithAPIKey(key string) ClientOption {
	return func(c *Client) {
		c.APIKey = key
	}
}

// WithBaseURL sets the TypeSafe API endpoint.
func WithBaseURL(url string) ClientOption {
	return func(c *Client) {
		c.BaseURL = url
	}
}

// WithModel sets the default model (e.g., "jev-latest").
func WithModel(model string) ClientOption {
	return func(c *Client) {
		c.Model = model
	}
}

// WithHTTPClient sets a custom HTTP client.
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// New creates a new TypeSafe Jev API client.
func New(opts ...ClientOption) *Client {
	c := &Client{
		BaseURL:    defaultBaseURL,
		Model:      "jev-latest",
		httpClient: http.DefaultClient,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

// Classify sends a classification request to TypeSafe.
// state can be a string, object, or array.
// questions is a map of question ID to Question.
func (c *Client) Classify(ctx context.Context, state interface{}, questions map[string]Question) (*Response, error) {
	req := Request{
		State:     state,
		Model:     c.Model,
		Questions: questions,
	}
	return c.doRequest(ctx, req)
}

// ClassifyString is a convenience method for string state.
func (c *Client) ClassifyString(ctx context.Context, state string, questions map[string]Question) (*Response, error) {
	return c.Classify(ctx, state, questions)
}

// ClassifyJSON sends a classification request with JSON-encoded state.
func (c *Client) ClassifyJSON(ctx context.Context, state interface{}, questions map[string]Question) (*Response, error) {
	return c.Classify(ctx, state, questions)
}

// doRequest executes the HTTP request to TypeSafe.
func (c *Client) doRequest(ctx context.Context, req Request) (*Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if c.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
	}

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("execute request: %w", err)
	}
	defer httpResp.Body.Close()

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("typesafe returned %d: %s", httpResp.StatusCode, string(respBody))
	}

	var resp Response
	if err := json.Unmarshal(respBody, &resp); err != nil {
		return nil, fmt.Errorf("decode response: %w", err)
	}

	return &resp, nil
}
