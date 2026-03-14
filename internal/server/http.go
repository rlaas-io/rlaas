package server

import (
	"log"
	"net/http"
)

// HTTPServer contains address and mux for the HTTP transport layer.
type HTTPServer struct {
	Addr string
	Mux  *http.ServeMux
}

// NewHTTPServer registers service endpoints.
// /v1/check returns decision responses.
// /healthz returns plain ok when server is healthy.
func NewHTTPServer(addr string, checkHandler http.Handler) *HTTPServer {
	mux := http.NewServeMux()
	mux.Handle("/v1/check", checkHandler)
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte("ok")) })
	return &HTTPServer{Addr: addr, Mux: mux}
}

// ListenAndServe starts the HTTP server.
func (s *HTTPServer) ListenAndServe() error {
	log.Printf("rlaas http server listening on %s", s.Addr)
	return http.ListenAndServe(s.Addr, s.Mux)
}
