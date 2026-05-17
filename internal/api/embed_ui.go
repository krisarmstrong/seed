//go:build !dev

// Production embed for the React frontend. The build tag `!dev` means
// this file compiles into release binaries. Dev-mode builds use the
// disk fallback in embed_ui_dev.go.
//
// Canonical embed path shared with niac and stem: internal/api/ui/.
// Vite (ui/vite.config.ts) writes its build directly here, and this
// file embeds it. Single source of truth, no copying.

package api

import (
	"embed"
	"io/fs"
)

//go:embed all:ui
var uiFS embed.FS

// getUIFS returns the embedded frontend filesystem rooted at the ui/
// subdirectory.
func getUIFS() (fs.FS, error) {
	return fs.Sub(uiFS, "ui")
}

// isUIEmbedded reports whether the UI is compiled into the binary.
// Always true in production builds.
func isUIEmbedded() bool {
	return true
}
