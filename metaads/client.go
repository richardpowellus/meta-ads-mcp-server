package metaads

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

const (
	defaultGraphURL = "https://graph.facebook.com"
	videoGraphURL   = "https://graph-video.facebook.com"
	defaultVersion  = "v25.0"
)

// APIError represents a Meta Graph API error.
type APIError struct {
	Message        string `json:"message"`
	Type           string `json:"type"`
	Code           int    `json:"code"`
	ErrorSubcode   int    `json:"error_subcode"`
	ErrorUserTitle string `json:"error_user_title"`
	ErrorUserMsg   string `json:"error_user_msg"`
	IsTransient    bool   `json:"is_transient"`
	FBTraceID      string `json:"fbtrace_id"`
}

func (e *APIError) Error() string {
	s := fmt.Sprintf("Meta API %d", e.Code)
	if e.ErrorSubcode > 0 {
		s += fmt.Sprintf("/%d", e.ErrorSubcode)
	}
	s += ": " + e.Message
	if e.FBTraceID != "" {
		s += " [fbtrace_id: " + e.FBTraceID + "]"
	}
	return s
}

func (e *APIError) IsRateLimit() bool {
	return e.Code == 4 || e.Code == 17 || e.Code == 613 || e.Code == 32
}

func (e *APIError) IsAuthError() bool {
	return e.Code == 190 || e.Code == 102
}

// Client wraps the Meta Graph API for advertising.
type Client struct {
	accessToken string
	appSecret   string
	adAccountID string
	businessID  string
	version     string
	http        *http.Client
	sem         chan struct{}
}

// NewClient creates a new Meta Graph API client.
func NewClient(accessToken, appSecret, adAccountID, businessID string) *Client {
	return &Client{
		accessToken: accessToken,
		appSecret:   appSecret,
		adAccountID: adAccountID,
		businessID:  businessID,
		version:     defaultVersion,
		http: &http.Client{
			Timeout: 120 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        20,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		sem: make(chan struct{}, 10),
	}
}

// AdAccountID returns the configured ad account ID (with act_ prefix).
func (c *Client) AdAccountID() string {
	id := c.adAccountID
	if !strings.HasPrefix(id, "act_") {
		id = "act_" + id
	}
	return id
}

// BusinessID returns the configured business ID.
func (c *Client) BusinessID() string { return c.businessID }

// appsecretProof computes HMAC-SHA256(app_secret, access_token) if app_secret is set.
func (c *Client) appsecretProof() string {
	if c.appSecret == "" {
		return ""
	}
	mac := hmac.New(sha256.New, []byte(c.appSecret))
	mac.Write([]byte(c.accessToken))
	return hex.EncodeToString(mac.Sum(nil))
}

// graphURL returns the base Graph API URL for the given path.
func (c *Client) graphURL(path string, isVideo bool) string {
	base := defaultGraphURL
	if isVideo {
		base = videoGraphURL
	}
	return fmt.Sprintf("%s/%s%s", base, c.version, path)
}

// doRequest executes an HTTP request with concurrency limiting, retry, and error parsing.
func (c *Client) doRequest(ctx context.Context, method, fullURL string, body io.Reader, contentType string) (json.RawMessage, error) {
	select {
	case c.sem <- struct{}{}:
		defer func() { <-c.sem }()
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	maxRetries := 3
	backoff := []time.Duration{2 * time.Second, 5 * time.Second, 10 * time.Second}
	var lastErr error

	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			wait := backoff[min(attempt-1, len(backoff)-1)]
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(wait):
			}
		}

		req, err := http.NewRequestWithContext(ctx, method, fullURL, body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+c.accessToken)
		if contentType != "" {
			req.Header.Set("Content-Type", contentType)
		}

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = err
			if method == "GET" {
				continue
			}
			return nil, err
		}

		respBody, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		if err != nil {
			return nil, fmt.Errorf("reading response: %w", err)
		}

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			if len(respBody) == 0 {
				return json.RawMessage(`{"success":true}`), nil
			}
			return json.RawMessage(respBody), nil
		}

		// Parse Graph API error
		var errResp struct {
			Error APIError `json:"error"`
		}
		_ = json.Unmarshal(respBody, &errResp)
		apiErr := &errResp.Error
		if apiErr.Message == "" {
			apiErr.Message = http.StatusText(resp.StatusCode)
			apiErr.Code = resp.StatusCode
		}

		// Retry on rate limits and transient errors
		if apiErr.IsRateLimit() || apiErr.IsTransient {
			lastErr = apiErr
			continue
		}
		if method == "GET" && resp.StatusCode >= 500 {
			lastErr = apiErr
			continue
		}

		return nil, apiErr
	}

	return nil, fmt.Errorf("request failed after %d retries: %w", maxRetries, lastErr)
}

// Get performs a GET request to the Graph API.
func (c *Client) Get(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	if proof := c.appsecretProof(); proof != "" {
		params.Set("appsecret_proof", proof)
	}
	fullURL := c.graphURL(path, false)
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}
	return c.doRequest(ctx, "GET", fullURL, nil, "")
}

// Post performs a POST request with form data.
func (c *Client) Post(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	params.Set("access_token", c.accessToken)
	if proof := c.appsecretProof(); proof != "" {
		params.Set("appsecret_proof", proof)
	}
	fullURL := c.graphURL(path, false)
	return c.doRequest(ctx, "POST", fullURL, strings.NewReader(params.Encode()), "application/x-www-form-urlencoded")
}

// PostJSON performs a POST request with a JSON body.
func (c *Client) PostJSON(ctx context.Context, path string, body json.RawMessage) (json.RawMessage, error) {
	fullURL := c.graphURL(path, false)
	params := url.Values{}
	if proof := c.appsecretProof(); proof != "" {
		params.Set("appsecret_proof", proof)
	}
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}
	return c.doRequest(ctx, "POST", fullURL, strings.NewReader(string(body)), "application/json")
}

// PostVideo performs a POST to the video upload endpoint.
func (c *Client) PostVideo(ctx context.Context, path string, body io.Reader, contentType string) (json.RawMessage, error) {
	fullURL := fmt.Sprintf("%s/%s%s", videoGraphURL, c.version, path)
	params := url.Values{}
	params.Set("access_token", c.accessToken)
	if proof := c.appsecretProof(); proof != "" {
		params.Set("appsecret_proof", proof)
	}
	fullURL += "?" + params.Encode()
	return c.doRequest(ctx, "POST", fullURL, body, contentType)
}

// Delete performs a DELETE request.
func (c *Client) Delete(ctx context.Context, path string) (json.RawMessage, error) {
	params := url.Values{}
	params.Set("access_token", c.accessToken)
	if proof := c.appsecretProof(); proof != "" {
		params.Set("appsecret_proof", proof)
	}
	fullURL := c.graphURL(path, false) + "?" + params.Encode()
	return c.doRequest(ctx, "DELETE", fullURL, nil, "")
}

// FetchAll auto-paginates a Graph API list endpoint using cursor-based pagination.
func (c *Client) FetchAll(ctx context.Context, path string, dataKey string, params url.Values) ([]json.RawMessage, error) {
	if params == nil {
		params = url.Values{}
	}
	if params.Get("limit") == "" {
		params.Set("limit", "100")
	}

	var allItems []json.RawMessage

	for {
		raw, err := c.Get(ctx, path, params)
		if err != nil {
			return nil, err
		}

		var wrapper map[string]json.RawMessage
		if err := json.Unmarshal(raw, &wrapper); err != nil {
			return nil, fmt.Errorf("parsing paginated response: %w", err)
		}

		dataRaw, ok := wrapper[dataKey]
		if !ok {
			break
		}

		var items []json.RawMessage
		if err := json.Unmarshal(dataRaw, &items); err != nil {
			return nil, fmt.Errorf("parsing data array: %w", err)
		}
		allItems = append(allItems, items...)

		// Check for next page cursor
		pagingRaw, ok := wrapper["paging"]
		if !ok {
			break
		}
		var paging struct {
			Cursors struct {
				After string `json:"after"`
			} `json:"cursors"`
			Next string `json:"next"`
		}
		if err := json.Unmarshal(pagingRaw, &paging); err != nil || paging.Next == "" {
			break
		}
		params.Set("after", paging.Cursors.After)
	}

	return allItems, nil
}

// GetWithPagination returns a single page of results with cursor info.
func (c *Client) GetWithPagination(ctx context.Context, path string, params url.Values) (json.RawMessage, error) {
	return c.Get(ctx, path, params)
}

// ParseIntParam safely parses an integer from a string parameter.
func ParseIntParam(s string) int {
	if s == "" {
		return 0
	}
	n, _ := strconv.Atoi(s)
	return n
}
