package main

import (
	"strings"
	"testing"
)

func validConfig() string {
	return `
servers:
  - frontend:
      port: 8080
    backend:
      hosts:
        - url: "http://localhost:8081"
`
}

func invalidConfigNoPort() string {
	return `
servers:
  - frontend:
      port: 0
    backend:
      hosts:
        - url: "http://localhost:8081"
`
}

func invalidConfigNoUrl() string {
	return `
servers:
  - frontend:
      port: 8080
    backend:
      hosts:
        - url: ""
`
}

func invalidConfigNoServers() string {
	return `
servers: []
`
}

func configWithMultipleServers() string {
	return `
servers:
  - frontend:
      port: 8080
    backend:
      hosts:
        - url: "http://localhost:8081"
  - frontend:
      port: 8081
    backend:
      hosts:
        - url: "http://localhost:8082"
`
}

func TestReadConfigValid(t *testing.T) {
	r := strings.NewReader(validConfig())
	config, err := ReadConfig(r)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(config.Servers) != 1 {
		t.Fatalf("expected 1 server, got %d", len(config.Servers))
	}
}

func TestReadConfigInvalidNoPort(t *testing.T) {
	r := strings.NewReader(invalidConfigNoPort())
	_, err := ReadConfig(r)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestReadConfigInvalidNoUrl(t *testing.T) {
	r := strings.NewReader(invalidConfigNoUrl())
	_, err := ReadConfig(r)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestReadConfigInvalidNoServers(t *testing.T) {
	r := strings.NewReader(invalidConfigNoServers())
	_, err := ReadConfig(r)
	if err == nil {
		t.Fatalf("expected error, got nil")
	}
}

func TestReadConfigMultipleServers(t *testing.T) {
	r := strings.NewReader(configWithMultipleServers())
	config, err := ReadConfig(r)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(config.Servers) != 2 {
		t.Fatalf("expected 2 servers, got %d", len(config.Servers))
	}
}
