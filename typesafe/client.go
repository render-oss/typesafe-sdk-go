package typesafe

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

const defaultBaseURL = "https://api.typesafe.ai/v1/systemone"

// Client is a TypeSafe Jev API client.
type Client struct {
	APIKey  string
	BaseURL string
	Model   string // Default model to use (e.g., "jev-latest")

	// MaxRetries is the number of retries after the first attempt. Default 2.
	MaxRetries int

	httpClient *http.Client
	observe    func(Attempt)
}

// Attempt describes one HTTP attempt, for logging and metrics.
type Attempt struct {
	N          int // 0 for the first attempt
	StatusCode int // 0 if the request never completed
	Duration   time.Duration
	Err        error
	WillRetry  bool
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

// WithMaxRetries sets how many times a failed request is retried. 0 disables retries.
func WithMaxRetries(n int) ClientOption {
	return func(c *Client) {
		c.MaxRetries = n
	}
}

// WithObserver registers a callback invoked once per HTTP attempt. Use it to log
// requests or record metrics. It runs on the calling goroutine, so keep it quick.
func WithObserver(fn func(Attempt)) ClientOption {
	return func(c *Client) {
		c.observe = fn
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
		MaxRetries: 2,
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
func (c *Client) Classify(ctx context.Context, state any, questions map[string]Question) (*Response, error) {
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
func (c *Client) ClassifyJSON(ctx context.Context, state any, questions map[string]Question) (*Response, error) {
	return c.Classify(ctx, state, questions)
}

// doRequest executes the HTTP request to TypeSafe.
func (c *Client) doRequest(ctx context.Context, req Request) (*Response, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	for n := 0; ; n++ {
		httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL, bytes.NewReader(body))
		if err != nil {
			return nil, fmt.Errorf("create request: %w", err)
		}
		httpReq.Header.Set("Content-Type", "application/json")
		if c.APIKey != "" {
			httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)
		}

		resp, retryAfter, err := c.attempt(httpReq, n)
		if err == nil {
			return resp, nil
		}
		if retryAfter < 0 || n >= c.MaxRetries {
			return nil, err
		}
		if retryAfter == 0 {
			retryAfter = backoff(n)
		}
		select {
		case <-ctx.Done():
			return nil, fmt.Errorf("retry aborted: %w", ctx.Err())
		case <-time.After(retryAfter):
		}
	}
}

// attempt performs one request. A negative retryAfter means the error is not
// retryable; zero means retry after the default backoff.
func (c *Client) attempt(httpReq *http.Request, n int) (resp *Response, retryAfter time.Duration, err error) {
	start := time.Now()
	status := 0
	defer func() {
		if c.observe != nil {
			c.observe(Attempt{
				N:          n,
				StatusCode: status,
				Duration:   time.Since(start),
				Err:        err,
				WillRetry:  err != nil && retryAfter >= 0 && n < c.MaxRetries,
			})
		}
	}()

	httpResp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, 0, fmt.Errorf("execute request: %w", err)
	}
	if httpResp == nil {
		// A custom http.Client can supply a RoundTripper that breaks this contract.
		return nil, -1, errors.New("http client returned no response and no error")
	}
	defer httpResp.Body.Close()
	status = httpResp.StatusCode

	respBody, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("read response: %w", err)
	}

	if status < 200 || status >= 300 {
		err := fmt.Errorf("typesafe returned %d: %s", status, string(respBody))
		if status != http.StatusTooManyRequests && status != http.StatusRequestTimeout && status < 500 {
			return nil, -1, err
		}
		after := time.Duration(0)
		// Retry-After is seconds or an HTTP date; anything else falls back to backoff.
		if v := httpResp.Header.Get("Retry-After"); v != "" {
			if secs, cerr := strconv.Atoi(v); cerr == nil {
				after = time.Duration(secs) * time.Second
			} else if t, cerr := http.ParseTime(v); cerr == nil {
				after = time.Until(t)
			}
			if after < 0 {
				after = 0
			}
		}
		return nil, after, err
	}

	var out Response
	if err := json.Unmarshal(respBody, &out); err != nil {
		return nil, -1, fmt.Errorf("decode response: %w", err)
	}
	return &out, 0, nil
}

// backoff returns the delay before retry n: 200ms, 400ms, 800ms... plus jitter.
func backoff(n int) time.Duration {
	d := min(200*time.Millisecond<<n, 10*time.Second)
	return d + time.Duration(rand.Int63n(int64(d/2)))
}
