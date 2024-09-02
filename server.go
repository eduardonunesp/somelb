package main

import (
	"fmt"
	"log/slog"
	"net/http"
)

type Frontend interface {
	Serve() error
}

type Backend interface {
	Request(request *http.Request)
}

type Server struct {
	err      chan error
	stop     chan struct{}
	recv     *http.Request
	send     *Response
	frontend Frontend
	backend  Backend
	conf     ServerConfig
	logger   *slog.Logger
}

func NewServer(conf ServerConfig) (*Server, error) {
	recv := make(chan *http.Request)
	send := make(chan *Response)

	logger := slog.Default().WithGroup("server")

	newFrontend, err := NewFrontend(
		conf.Frontend,
		logger.WithGroup("frontend"),
		send,
		recv,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new frontned: %w", err)
	}

	newBackend, err := NewBackend(
		conf.Backend,
		logger.WithGroup("backend"),
		send,
		recv,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create new backend: %w", err)
	}

	return &Server{
		frontend: newFrontend,
		backend:  newBackend,
		err:      make(chan error),
		stop:     make(chan struct{}),
		conf:     conf,
		logger:   logger,
	}, nil
}

func (s *Server) Run() error {
	s.logger.Info("starting server")

	go func() {
		if err := s.frontend.Serve(); err != nil {
			errMsg := "failed to start frontend"
			s.logger.Error(errMsg, slog.Any("error", err))
			s.err <- fmt.Errorf("%s: %w", errMsg, err)
		}
	}()

	select {
	case err := <-s.err:
		s.logger.Error("server stopped with error", slog.Any("error", err))
		return err
	case <-s.stop:
		s.logger.Info("server stopped gracefully")
		return nil
	}
}
