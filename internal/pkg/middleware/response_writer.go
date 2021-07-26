package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
)

type MyResponseWriter struct {
	http.ResponseWriter
	status int
}

func (w *MyResponseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.status = code
}

func (w *MyResponseWriter) Status() int {
	return w.status
}
func (w *MyResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}

// StatusResponseWriter wrap go response writer with a record of the http return status
func StatusResponseWriter(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mw := &MyResponseWriter{
			ResponseWriter: w,
			status:         http.StatusOK,
		}
		h.ServeHTTP(mw, r)
	})
}
