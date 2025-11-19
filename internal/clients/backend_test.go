package clients

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHttpBackendClient_GetAndPost(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/get":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"message":"ok"}`))
		case "/post":
			if r.Method != http.MethodPost {
				w.WriteHeader(http.StatusMethodNotAllowed)
				return
			}
			var payload map[string]string
			_ = json.NewDecoder(r.Body).Decode(&payload)
			if payload["foo"] != "bar" {
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			w.WriteHeader(http.StatusCreated)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	client := NewHttpBackendClient(server.URL)
	ctx := context.Background()

	resp, err := client.DoGetRequest(ctx, "/get")
	if err != nil {
		t.Fatalf("DoGetRequest failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()

	resp, err = client.DoPostRequest(ctx, "/post", map[string]string{"foo": "bar"})
	if err != nil {
		t.Fatalf("DoPostRequest failed: %v", err)
	}
	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("expected 201, got %d", resp.StatusCode)
	}
	_ = resp.Body.Close()
}


