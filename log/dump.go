package log

import (
	"bufio"
	"errors"
	"io"
	"net"
	"net/http"
)

type bodyReader struct {
	io.ReadCloser
	maxSize int
	bytes   int
}

func newBodyReader(reader io.ReadCloser, maxSize int) *bodyReader {
	return &bodyReader{
		ReadCloser: reader,
		maxSize:    maxSize,
		bytes:      0,
	}
}

type bodyWriter struct {
	http.ResponseWriter
	maxSize int
	bytes   int
	status  int
}

func newBodyWriter(writer http.ResponseWriter, maxSize int) *bodyWriter {
	return &bodyWriter{
		ResponseWriter: writer,
		maxSize:        maxSize,
		bytes:          0,
		status:         http.StatusOK,
	}
}

func (w *bodyWriter) Status() int {
	return w.status
}

// implements http.ResponseWriter
func (w *bodyWriter) Write(b []byte) (int, error) {
	w.bytes += len(b) //nolint:staticcheck
	return w.ResponseWriter.Write(b)
}

// implements http.ResponseWriter
func (r *bodyWriter) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

// implements http.Flusher
func (w *bodyWriter) Flush() {
	if w.ResponseWriter.(http.Flusher) != nil {
		w.ResponseWriter.(http.Flusher).Flush()
	}
}

// implements http.Hijacker
func (w *bodyWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if w.ResponseWriter.(http.Hijacker) != nil {
		return w.ResponseWriter.(http.Hijacker).Hijack()
	}

	return nil, nil, errors.New("Hijack not supported")
}
