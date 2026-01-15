package clients

import (
	"actions-service/internal/observability"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type HttpBackendClient interface {
	DoGetRequest(ctx context.Context,  path string)(*http.Response, error)
	DoPostRequest(ctx context.Context, path string, body interface{})(*http.Response, error)
	DoPutRequest(ctx context.Context, path string, body interface{})(*http.Response, error)
}

type httpBackendClient struct {
	baseUrl string
	client *http.Client
	logger *slog.Logger
}

func NewHttpBackendClient(baseUrl string) HttpBackendClient {
	return &httpBackendClient{
		baseUrl: baseUrl,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
		logger: observability.NewLogger("info"),
	}
}

func (c *httpBackendClient) DoGetRequest(ctx context.Context,  path string)(*http.Response, error) {
	start := time.Now()
	url := fmt.Sprintf("%s%s", c.baseUrl, path)
	
	c.logger.InfoContext(ctx, "Backend HTTP Request",
		slog.String("method", "GET"),
		slog.String("url", url),
		slog.String("path", path),
	)
	
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		c.logger.ErrorContext(ctx, "Failed to create backend request",
			slog.String("method", "GET"),
			slog.String("url", url),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	
	resp, err := c.client.Do(req)
	duration := time.Since(start)
	
	if err != nil {
		c.logger.ErrorContext(ctx, "Backend HTTP Request Failed",
			slog.String("method", "GET"),
			slog.String("url", url),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	
	c.logger.InfoContext(ctx, "Backend HTTP Response",
		slog.String("method", "GET"),
		slog.String("url", url),
		slog.Int("status_code", resp.StatusCode),
		slog.Int64("duration_ms", duration.Milliseconds()),
	)
	
	return resp, nil
}

func (c *httpBackendClient) DoPostRequest(ctx context.Context,  path string, body interface{}) (*http.Response, error) {
	start := time.Now()
	url := fmt.Sprintf("%s%s", c.baseUrl, path)
	var req *http.Request
	var err error
	var jsonBody []byte

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			c.logger.ErrorContext(ctx, "Failed to marshal request body",
				slog.String("method", "POST"),
				slog.String("url", url),
				slog.String("error", err.Error()),
			)
			return nil, err
		}
		
		c.logger.InfoContext(ctx, "Backend HTTP Request",
			slog.String("method", "POST"),
			slog.String("url", url),
			slog.String("path", path),
			slog.String("payload", string(jsonBody)),
		)
		
		req, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			c.logger.ErrorContext(ctx, "Failed to create backend request",
				slog.String("method", "POST"),
				slog.String("url", url),
				slog.String("error", err.Error()),
			)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		c.logger.InfoContext(ctx, "Backend HTTP Request",
			slog.String("method", "POST"),
			slog.String("url", url),
			slog.String("path", path),
		)
		
		req, err = http.NewRequestWithContext(ctx, "POST", url, nil)
		if err != nil {
			c.logger.ErrorContext(ctx, "Failed to create backend request",
				slog.String("method", "POST"),
				slog.String("url", url),
				slog.String("error", err.Error()),
			)
			return nil, err
		}
	}

	resp, err := c.client.Do(req)
	duration := time.Since(start)
	
	if err != nil {
		c.logger.ErrorContext(ctx, "Backend HTTP Request Failed",
			slog.String("method", "POST"),
			slog.String("url", url),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	
	// Read response body for logging
	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		// Restore body for caller
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}
	
	c.logger.InfoContext(ctx, "Backend HTTP Response",
		slog.String("method", "POST"),
		slog.String("url", url),
		slog.Int("status_code", resp.StatusCode),
		slog.Int64("duration_ms", duration.Milliseconds()),
		slog.String("response_body", string(respBody)),
	)
	
	return resp, nil
}

func(c *httpBackendClient) DoPutRequest(ctx context.Context,  path string, body interface{}) (*http.Response, error) {
	start := time.Now()
	url := fmt.Sprintf("%s%s", c.baseUrl, path)
	var req *http.Request
	var err error
	var jsonBody []byte

	if body != nil {
		jsonBody, err = json.Marshal(body)
		if err != nil {
			c.logger.ErrorContext(ctx, "Failed to marshal request body",
				slog.String("method", "PUT"),
				slog.String("url", url),
				slog.String("error", err.Error()),
			)
			return nil, err
		}
		
		c.logger.InfoContext(ctx, "Backend HTTP Request",
			slog.String("method", "PUT"),
			slog.String("url", url),
			slog.String("path", path),
			slog.String("payload", string(jsonBody)),
		)
		
		req, err = http.NewRequestWithContext(ctx, "PUT", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			c.logger.ErrorContext(ctx, "Failed to create backend request",
				slog.String("method", "PUT"),
				slog.String("url", url),
				slog.String("error", err.Error()),
			)
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		c.logger.InfoContext(ctx, "Backend HTTP Request",
			slog.String("method", "PUT"),
			slog.String("url", url),
			slog.String("path", path),
		)
		
		req, err = http.NewRequestWithContext(ctx, "PUT", url, nil)
		if err != nil {
			c.logger.ErrorContext(ctx, "Failed to create backend request",
				slog.String("method", "PUT"),
				slog.String("url", url),
				slog.String("error", err.Error()),
			)
			return nil, err
		}
	}

	resp, err := c.client.Do(req)
	duration := time.Since(start)
	
	if err != nil {
		c.logger.ErrorContext(ctx, "Backend HTTP Request Failed",
			slog.String("method", "PUT"),
			slog.String("url", url),
			slog.Int64("duration_ms", duration.Milliseconds()),
			slog.String("error", err.Error()),
		)
		return nil, err
	}
	
	// Read response body for logging
	var respBody []byte
	if resp.Body != nil {
		respBody, _ = io.ReadAll(resp.Body)
		resp.Body.Close()
		// Restore body for caller
		resp.Body = io.NopCloser(bytes.NewBuffer(respBody))
	}
	
	c.logger.InfoContext(ctx, "Backend HTTP Response",
		slog.String("method", "PUT"),
		slog.String("url", url),
		slog.Int("status_code", resp.StatusCode),
		slog.Int64("duration_ms", duration.Milliseconds()),
		slog.String("response_body", string(respBody)),
	)
	
	return resp, nil
}
