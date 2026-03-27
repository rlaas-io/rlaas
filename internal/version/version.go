// Package version holds build-time metadata injected via -ldflags.
// At runtime, consuming code can call Info() or inspect the exported
// variables directly (e.g. to populate a /version endpoint or log line).
package version

import "fmt"

// Set at build time with:
//
//	-X github.com/rlaas-io/rlaas/internal/version.Version=$(git describe --tags --always --dirty)
//	-X github.com/rlaas-io/rlaas/internal/version.Commit=$(git rev-parse --short HEAD)
//	-X github.com/rlaas-io/rlaas/internal/version.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)
var (
	Version   = "dev"
	Commit    = "unknown"
	BuildTime = "unknown"
)

// Info returns a single-line human-readable version string suitable for log
// output and the /version endpoint.
func Info() string {
	return fmt.Sprintf("%s (commit=%s built=%s)", Version, Commit, BuildTime)
}
