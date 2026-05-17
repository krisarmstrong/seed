//go:build dev

// Dev embed for the React frontend. Reads from disk so frontend changes
// are picked up without rebuilding the Go binary. Build with `-tags dev`
// (or via `make dev`) to activate this file.

package api

import (
	"io/fs"
	"os"
)

// getUIFS returns the on-disk frontend filesystem. Path is relative to
// the process working directory, which is the repo root for `make dev`.
func getUIFS() (fs.FS, error) {
	return os.DirFS("internal/api/ui"), nil
}

// isUIEmbedded reports whether the UI is compiled into the binary.
// Always false in dev builds — assets are served from disk.
func isUIEmbedded() bool {
	return false
}
