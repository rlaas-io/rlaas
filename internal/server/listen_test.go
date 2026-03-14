package server

import "testing"

func TestListenAndServeInvalidAddress(t *testing.T) {
	s := &HTTPServer{Addr: ":-1", Mux: nil}
	if err := s.ListenAndServe(); err == nil {
		t.Fatalf("expected listen error for invalid address")
	}
}
