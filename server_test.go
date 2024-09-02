package main

import (
	"errors"
	"io"
	"log/slog"
	"net/http"
	"testing"
)

func newTestServer() *Server {
	logger := slog.New(slog.NewJSONHandler(io.Discard, nil))
	return &Server{
		err:    make(chan error),
		stop:   make(chan struct{}),
		logger: logger,
	}
}

func newTestFrontend() Frontend {
	return &mockFrontend{}
}

func newTestBackendZeroed() Backend {
	return &mockBackend{}
}

type mockFrontend struct{}

func (m *mockFrontend) Serve() error {
	return nil
}

type mockBackend struct{}

func (m *mockBackend) Request(request *http.Request) {}

func TestNewServerValidConfig(t *testing.T) {
	conf := ServerConfig{
		Frontend: FrontendConfig{Port: 8080},
		Backend: BackendConfig{
			Hosts: []struct {
				Url string `yaml:"url"`
			}{
				{Url: "http://localhost:8080"},
			},
		},
	}
	server, err := NewServer(conf)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if server.conf.Frontend.Port != 8080 {
		t.Fatalf("expected frontend port 8080, got %d", server.conf.Frontend.Port)
	}
}

func TestRunServer(t *testing.T) {
	server := newTestServer()
	server.frontend = newTestFrontend()
	server.backend = newTestBackendZeroed()

	go func() {
		server.stop <- struct{}{}
	}()

	err := server.Run()
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestRunServerWithError(t *testing.T) {
	server := newTestServer()
	server.frontend = &mockFrontendWithError{}
	server.backend = newTestBackendZeroed()

	go func() {
		server.err <- errors.New("frontend error")
	}()

	err := server.Run()
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

type mockFrontendWithError struct{}

func (m *mockFrontendWithError) Serve() error {
	return errors.New("frontend error")
}
