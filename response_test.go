package main

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

func newTestBackend() *backend {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return &backend{
		hosts:  make([]*hostContainer, 0),
		stop:   make(chan struct{}),
		recv:   make(chan *Response),
		send:   make(chan *http.Request),
		logger: logger,
	}
}

func newTestHostContainer(urlStr string) (*hostContainer, error) {
	host, err := newTestHost(urlStr)
	if err != nil {
		return nil, err
	}
	return NewHostContainer(host), nil
}

func TestNewBackendValidConfig(t *testing.T) {
	conf := BackendConfig{
		Hosts: []struct {
			Url string `yaml:"url"`
		}{
			{Url: "http://localhost:8080"},
		},
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	recv := make(chan *Response)
	send := make(chan *http.Request)

	backend, err := NewBackend(conf, logger, recv, send)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(backend.hosts) != 1 {
		t.Fatalf("expected 1 host, got %d", len(backend.hosts))
	}
}

func TestNewBackendInvalidConfig(t *testing.T) {
	conf := BackendConfig{
		Hosts: []struct {
			Url string `yaml:"url"`
		}{
			{Url: ""},
		},
	}
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	recv := make(chan *Response)
	send := make(chan *http.Request)

	_, err := NewBackend(conf, logger, recv, send)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestNextHostNoHostsAvailable(t *testing.T) {
	backend := newTestBackend()
	_, err := backend.nextHost()
	if !errors.Is(err, ErrNoBackendHostAvailable) {
		t.Fatalf("expected ErrNoBackendHostAvailable, got %v", err)
	}
}

func TestRequestNoHostsAvailable(t *testing.T) {
	backend := newTestBackend()
	req, _ := http.NewRequest("GET", "http://localhost:8080", nil)
	req.Header.Set("__requestID", "test-request-id")

	go backend.Request(req)
	resp := <-backend.recv
	if resp.Code != http.StatusBadGateway {
		t.Fatalf("expected status code 502, got %d", resp.Code)
	}
}
