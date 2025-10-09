package clients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type HttpBackendClient interface {
	DoGetRequest(ctx context.Context,  path string)(*http.Response, error)
	DoPostRequest(ctx context.Context, path string, body interface{})(*http.Response, error)
}

type httpBackendClient struct {
	baseUrl string
	client *http.Client
}

func NewHttpBackendClient(baseUrl string) HttpBackendClient {
	return &httpBackendClient{
		baseUrl: baseUrl,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *httpBackendClient) DoGetRequest(ctx context.Context,  path string)(*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseUrl, path)
	//log.Println("url: ", url)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c *httpBackendClient) DoPostRequest(ctx context.Context,  path string, body interface{}) (*http.Response, error) {
	url := fmt.Sprintf("%s%s", c.baseUrl, path)
	log.Println("url: ", url)
	var req *http.Request
	var err error

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		log.Printf("JSON sent: %s", string(jsonBody))
		req, err = http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonBody))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, "POST", url, nil)
	}

	if err != nil {
		return nil, err
	}
	
	return c.client.Do(req)
}
