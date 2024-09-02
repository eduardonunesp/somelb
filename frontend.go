package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var requestID atomic.Uint64

type frontend struct {
	port   uint16
	alive  bool
	logger *slog.Logger
	send   chan *Response
	recv   chan *http.Request
}

func NewFrontend(
	conf FrontendConfig,
	logger *slog.Logger,
	recv chan *Response,
	send chan *http.Request,
) (*frontend, error) {
	newFrontend := &frontend{
		port:   conf.Port,
		logger: logger,
		send:   recv, recv: send,
	}
	return newFrontend, nil
}

func (f *frontend) request(w http.ResponseWriter, req *http.Request) {
	requestID.Add(1)
	requestIDStr := strconv.FormatUint(requestID.Load(), 10)
	req.Header.Set("__requestID", requestIDStr)
	f.recv <- req

	f.logger.Info("received request", slog.Any("req", req))

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()

	for {
		select {
		case <-ctx.Done():
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		case backendResponse := <-f.send:
			if requestIDStr == backendResponse.RequestID {
				for k, v := range backendResponse.HeaderMap {
					for _, vv := range v {
						w.Header().Set(k, vv)
					}
				}
				w.WriteHeader(backendResponse.Code)
				w.Write(backendResponse.Body.Bytes())
				return
			}
		}
	}
}

func (f *frontend) Serve() error {
	http.HandleFunc("/", f.request)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", f.port), nil); err != nil {
		return fmt.Errorf("failed to start frontend %w", err)
	}
	return nil
}
