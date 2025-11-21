package llm

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIClient_Chat(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req chatRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			w.WriteHeader(400)
			return
		}
		_ = req
		_ = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"choices":[{"message":{"role":"assistant","content":"hi"}}]}`))
	}))
	defer ts.Close()

	c := &OpenAIClient{EndpointURL: ts.URL, Model: "m", APIKey: "k"}
	got, err := c.Chat(context.Background(), "hello")
	if err != nil || got != "hi" {
		t.Fatalf("got=%q err=%v", got, err)
	}
}

func TestOpenAIClient_Non200(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(401)
		_, _ = w.Write([]byte("nope"))
	}))
	defer ts.Close()

	c := &OpenAIClient{EndpointURL: ts.URL, Model: "m", APIKey: "k"}
	_, err := c.Chat(context.Background(), "hello")
	if err == nil {
		t.Fatalf("expected error for non-200")
	}
}
