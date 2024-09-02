package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/url"

	"gopkg.in/yaml.v2"
)

// FrontendConfig holds the configuration for the frontend server.
type FrontendConfig struct {
	Port uint16 `yaml:"port"`
}

// BackendConfig holds the configuration for the backend servers.
type BackendConfig struct {
	Hosts []struct {
		Url string `yaml:"url"`
	}
}

// ServerConfig holds the configuration for a server, including frontend and backend.
type ServerConfig struct {
	Backend  BackendConfig  `yaml:"backend"`
	Frontend FrontendConfig `yaml:"frontend"`
}

// Config holds the overall configuration for the load balancer.
type Config struct {
	Servers []ServerConfig `yaml:"servers"`
}

// Validate checks if the FrontendConfig is valid.
func (f FrontendConfig) Validate() error {
	if f.Port == 0 {
		return errors.New("port cannot be 0")
	}
	return nil
}

// Validate checks if the BackendConfig is valid.
func (b BackendConfig) Validate() error {
	for _, b := range b.Hosts {
		if b.Url == "" {
			return errors.New("url address cannot be empty")
		}
		if _, err := url.Parse(b.Url); err != nil {
			return errors.New("invalid URL")
		}
	}
	return nil
}

// Validate checks if the Config is valid.
func (c Config) Validate() error {
	if len(c.Servers) == 0 {
		return errors.New("no server conf found")
	}

	for _, serverConfig := range c.Servers {
		if err := serverConfig.Frontend.Validate(); err != nil {
			return err
		}
		if err := serverConfig.Backend.Validate(); err != nil {
			return err
		}
	}

	return nil
}

// ReadConfig reads and returns the configuration for the load balancer.
func ReadConfig(r io.Reader) (*Config, error) {
	data := new(bytes.Buffer)
	if _, err := data.ReadFrom(r); err != nil {
		return nil, fmt.Errorf("failed to read buffer: %w", err)
	}

	var config Config
	err := yaml.Unmarshal(data.Bytes(), &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parser conf: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate configuration: %w", err)
	}

	return &config, nil
}
