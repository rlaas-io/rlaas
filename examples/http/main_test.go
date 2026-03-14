package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunExample(t *testing.T) {
	buf := &bytes.Buffer{}
	if err := run(buf); err != nil {
		t.Fatalf("run failed: %v", err)
	}
	out := buf.String()
	if !strings.Contains(out, "request 1") {
		t.Fatalf("unexpected output: %s", out)
	}
}

func TestMainExample(t *testing.T) {
	main()
}

func TestMainExamplePanicPath(t *testing.T) {
	old, _ := os.Getwd()
	defer os.Chdir(old)
	_ = os.Chdir(t.TempDir())
	defer func() { _ = recover() }()
	_ = filepath.Separator
	main()
}

func TestRunExampleErrorPath(t *testing.T) {
	old, _ := os.Getwd()
	defer os.Chdir(old)
	tmp := t.TempDir()
	_ = os.Chdir(tmp)
	_ = os.MkdirAll("examples", 0755)
	_ = os.WriteFile(filepath.Join("examples", "policies.json"), []byte("not-json"), 0644)
	if err := run(&bytes.Buffer{}); err != nil {
		t.Fatalf("expected fail-open run behavior, got error: %v", err)
	}
}
