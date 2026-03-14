package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	httpadapter "rlaas/internal/adapter/http"
	"rlaas/internal/server"
	"rlaas/internal/store/counter/memory"
	filestore "rlaas/internal/store/policy/file"
	"rlaas/pkg/rlaas"
)

// main starts a local HTTP server with file policies and memory counters.
func main() {
	if err := run(defaultListen); err != nil {
		log.Printf("rlaas-server startup failed: %v", err)
		return
	}
}

func run(listenFn func(*server.HTTPServer) error) error {
	policyFile := os.Getenv("RLAAS_POLICY_FILE")
	if policyFile == "" {
		policyFile = "examples/policies.json"
	}
	if _, err := os.Stat(policyFile); err != nil {
		return fmt.Errorf("policy file not found: %s", policyFile)
	}
	client := rlaas.New(rlaas.Options{
		PolicyStore:  filestore.New(policyFile),
		CounterStore: memory.New(),
		KeyPrefix:    "rlaas",
	})
	checkHandler := httpadapter.CheckHandler(client)
	httpServer := server.NewHTTPServer(":8080", checkHandler)
	mux := httpServer.Mux
	mw := httpadapter.NewMiddleware(client)
	// Demo endpoint applies middleware so you can observe enforcement quickly.
	mux.Handle("/demo", mw.Handler(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write([]byte("ok"))
	})))
	return listenFn(httpServer)
}

func defaultListen(s *server.HTTPServer) error {
	return s.ListenAndServe()
}
