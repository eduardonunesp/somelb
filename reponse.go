package main

import (
	"bytes"
	"net/http"
)

type Response struct {
	RequestID string
	Code      int
	HeaderMap http.Header
	Body      *bytes.Buffer
}

func NewResponse(requestID string) *Response {
	return &Response{
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
		Code:      200,
		RequestID: requestID,
	}
}

func NewResponseWithStatus(requestID string, statusCode int) *Response {
	return &Response{
		RequestID: requestID,
		Code:      statusCode,
		HeaderMap: make(http.Header),
		Body:      new(bytes.Buffer),
	}
}

func (rw *Response) Header() http.Header {
	m := rw.HeaderMap
	if m == nil {
		m = make(http.Header)
		rw.HeaderMap = m
	}
	return m
}

func (rw *Response) writeHeader(b []byte, str string) {
	rw.WriteHeader(200)
}

func (rw *Response) Write(buf []byte) (int, error) {
	rw.writeHeader(buf, "")
	if rw.Body != nil {
		rw.Body.Write(buf)
	}
	return len(buf), nil
}

func (rw *Response) WriteHeader(code int) {
	rw.Code = code
	if rw.HeaderMap == nil {
		rw.HeaderMap = make(http.Header)
	}
}
