package main

import (
	"errors"
	"log/slog"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"
)

const (
	maxFailedAttempts  = 3
	thresholdHealth    = 3
	timeoutHealthCheck = time.Second * 10
)

var (
	ErrInvalidURL = errors.New("invalid url string given")
)

type HealthCheck struct {
	failedAttempts  int
	successAttempts int
	alive           bool
	logger          *slog.Logger
}

type Host struct {
	stop        chan struct{}
	proxy       *httputil.ReverseProxy
	url         *url.URL
	logger      *slog.Logger
	healthCheck HealthCheck
}

// NewHostFromString creates a new urlStr representation from string
func NewHostFromString(urlStr string, logger *slog.Logger) (*Host, error) {
	nUrl, err := url.Parse(urlStr)
	if err != nil {
		return nil, errors.Join(ErrInvalidURL, err)
	}
	newHost := &Host{
		stop:   make(chan struct{}),
		proxy:  httputil.NewSingleHostReverseProxy(nUrl),
		url:    nUrl,
		logger: logger,
		healthCheck: HealthCheck{
			logger: logger,
		},
	}
	go newHost.internalRoutines()
	return newHost, nil
}

func (h *Host) IsAlive() bool {
	return h.healthCheck.alive
}

func (h *Host) internalRoutines() {
	for {
		select {
		case <-h.stop:
			return
		case <-time.After(timeoutHealthCheck):
			h.healthCheck.healthChecking(h.url)
		}
	}
}

func (h *HealthCheck) healthChecking(targetUrl *url.URL) {
	var requestOk bool
	defer func() {
		if !requestOk {
			if h.alive {
				h.failedAttempts++
				if h.failedAttempts >= maxFailedAttempts {
					h.alive = false
					h.failedAttempts = maxFailedAttempts
					h.logger.Info("Host is not alive", slog.String("url", targetUrl.String()))
				}
			}
			return
		}

		if !h.alive {
			h.successAttempts++
			if h.successAttempts >= thresholdHealth {
				h.alive = true
				h.failedAttempts = 0
				h.successAttempts = 0
				h.logger.Info("Host is alive", slog.String("url", targetUrl.String()))
			}
			h.logger.Info("Host given signs of alive", slog.String("url", targetUrl.String()))
		}
	}()

	response, err := http.Get(targetUrl.String())

	if err != nil {
		h.logger.Error("failed to check alive", slog.String("url", targetUrl.String()))
		return
	}

	if response.StatusCode != 200 {
		h.logger.Info("cannot health check Host", slog.String("url", targetUrl.String()))
		return
	}

	requestOk = true
}

func (h *Host) SendRequest(req *http.Request) (*Response, error) {
	response := NewResponse(req.Header.Get("__requestID"))
	h.proxy.ServeHTTP(response, req)
	return response, nil
}
