package main

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"rlaas/internal/server"
)

func TestRunBuildsServer(t *testing.T) {
	path := filepath.Join("..", "..", "examples", "policies.json")
	_ = os.Setenv("RLAAS_POLICY_FILE", path)
	defer os.Unsetenv("RLAAS_POLICY_FILE")

	called := false
	err := run(func(s *server.HTTPServer) error {
		called = true
		if s == nil || s.Mux == nil {
			t.Fatalf("expected server")
		}
		return nil
	})
	if err != nil || !called {
		t.Fatalf("expected run success: %v", err)
	}
}

func TestRunMissingPolicyFile(t *testing.T) {
	_ = os.Setenv("RLAAS_POLICY_FILE", filepath.Join(t.TempDir(), "missing.json"))
	defer os.Unsetenv("RLAAS_POLICY_FILE")
	if err := run(func(s *server.HTTPServer) error { return nil }); err == nil || !strings.Contains(err.Error(), "policy file not found") {
		t.Fatalf("expected missing file error")
	}
}

func TestDefaultListenInvalid(t *testing.T) {
	if err := defaultListen(&server.HTTPServer{Addr: ":-1", Mux: nil}); err == nil {
		t.Fatalf("expected default listen error")
	}
}

func TestMainServerReturnsOnStartupError(t *testing.T) {
	_ = os.Setenv("RLAAS_POLICY_FILE", filepath.Join(t.TempDir(), "missing.json"))
	defer os.Unsetenv("RLAAS_POLICY_FILE")
	main()
}

func TestRunUsesDefaultPolicyFile(t *testing.T) {
	os.Unsetenv("RLAAS_POLICY_FILE")
	old, _ := os.Getwd()
	defer os.Chdir(old)
	_ = os.Chdir(filepath.Join("..", ".."))
	called := false
	err := run(func(s *server.HTTPServer) error {
		called = true
		return nil
	})
	if err != nil || !called {
		t.Fatalf("expected run success with default policy file")
	}
}

func TestRunListenError(t *testing.T) {
	path := filepath.Join("..", "..", "examples", "policies.json")
	_ = os.Setenv("RLAAS_POLICY_FILE", path)
	defer os.Unsetenv("RLAAS_POLICY_FILE")
	err := run(func(s *server.HTTPServer) error { return errors.New("listen failed") })
	if err == nil || !strings.Contains(err.Error(), "listen failed") {
		t.Fatalf("expected listen error")
	}
}
