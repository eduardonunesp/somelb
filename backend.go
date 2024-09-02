// Package main implements the backend functionality for the server.
package main

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
)

// Error definitions
var (
	ErrFailedToLoadHost       = errors.New("failed to load backend")
	ErrNoBackendHostAvailable = errors.New("no backend host available")
)

// hostContainer holds a Host and its request counter.
type hostContainer struct {
	RequestCounter int
	Host           *Host
}

// NewHostContainer creates a new hostContainer with the given Host.
func NewHostContainer(host *Host) *hostContainer {
	return &hostContainer{
		RequestCounter: 0,
		Host:           host,
	}
}

// IsAlive checks if the Host in the hostContainer is alive.
func (h *hostContainer) IsAlive() bool {
	return h.Host.IsAlive()
}

// backend represents the backend server with multiple hosts.
type backend struct {
	hosts  []*hostContainer
	stop   chan struct{}
	recv   chan *Response
	send   chan *http.Request
	logger *slog.Logger
	mu     sync.Mutex
}

// NewBackend creates a new backend with the given configuration, logger, and channels.
func NewBackend(
	conf BackendConfig,
	logger *slog.Logger,
	recv chan *Response,
	send chan *http.Request,
) (*backend, error) {
	if err := conf.Validate(); err != nil {
		return nil, err
	}
	var newBackend = &backend{
		hosts:  make([]*hostContainer, 0),
		stop:   make(chan struct{}),
		recv:   recv,
		send:   send,
		logger: logger,
	}
	for hostIdx, hostConf := range conf.Hosts {
		newHost, err := NewHostFromString(
			hostConf.Url,
			logger.WithGroup(fmt.Sprintf("host_%d", hostIdx)),
		)
		if err != nil {
			return nil, errors.Join(ErrFailedToLoadHost, err)
		}
		newBackend.hosts = append(newBackend.hosts, NewHostContainer(newHost))
	}
	go newBackend.internalRoutines()
	return newBackend, nil
}

// internalRoutines handles incoming requests and stops the backend when needed.
func (b *backend) internalRoutines() {
	for {
		select {
		case <-b.stop:
			return
		case newRequest := <-b.send:
			b.logger.Info("new request to backend")
			b.Request(newRequest)
		}
	}
}

// nextHost selects the next available host with the least connections.
func (b *backend) nextHost() (*Host, error) {
	b.mu.Lock()
	defer b.mu.Unlock()

	if len(b.hosts) == 0 {
		return nil, ErrNoBackendHostAvailable
	}

	leastConnCounter := -1
	leastConnHost := b.hosts[0]

	for _, host := range b.hosts {
		if (host.RequestCounter < leastConnCounter || leastConnCounter == -1) && host.IsAlive() {
			leastConnCounter = host.RequestCounter
			leastConnHost = host
		}
	}

	if !leastConnHost.IsAlive() {
		return nil, ErrNoBackendHostAvailable
	}

	leastConnHost.RequestCounter += 1
	b.logger.Info("least connections host")

	return leastConnHost.Host, nil
}

// Request sends the given request to the next available host.
func (b *backend) Request(request *http.Request) {
	host, err := b.nextHost()
	if err != nil {
		errType := http.StatusInternalServerError
		if errors.Is(err, ErrNoBackendHostAvailable) {
			errType = http.StatusBadGateway
		}
		b.recv <- NewResponseWithStatus(
			request.Header.Get("__requestID"),
			errType,
		)
		b.logger.Error("error on get next Host", slog.Any("error", err))
		return
	}

	resp, err := host.SendRequest(request)
	if err != nil {
		b.logger.Error("error on send request", slog.Any("error", err))
	}
	b.recv <- resp
}
