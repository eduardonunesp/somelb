package main

import (
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

func newTestHost(urlStr string) (*Host, error) {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return NewHostFromString(urlStr, logger)
}

func TestNewHostFromStringValidURL(t *testing.T) {
	host, err := newTestHost("http://localhost:8080")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if host.url.String() != "http://localhost:8080" {
		t.Fatalf("expected url to be http://localhost:8080, got %s", host.url.String())
	}
}

func TestIsAliveInitiallyFalse(t *testing.T) {
	host, _ := newTestHost("http://localhost:8080")
	if host.IsAlive() {
		t.Fatalf("expected host to be initially not alive")
	}
}

func TestSendRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	host, _ := newTestHost(server.URL)
	req, _ := http.NewRequest("GET", server.URL, nil)
	req.Header.Set("__requestID", "test-request-id")

	resp, err := host.SendRequest(req)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if resp.Code != http.StatusOK {
		t.Fatalf("expected status code 200, got %d", resp.Code)
	}
}
