//go:build linux

package truststore

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// pathExists reports whether a filesystem path exists. Used by
// detectLinuxStore to pick the trust-store layout the host uses.
func pathExists(p string) bool {
	_, err := os.Stat(p)
	return err == nil
}

// linuxStore describes how a particular Linux distribution stores its
// system-wide trust anchors and how to refresh the bundle afterwards.
//
// Refresh is a closure rather than a string slice so the gosec G204
// "exec from variable" rule sees literal argument lists at the call
// site. Each detectLinuxStore branch wires in its distribution-specific
// command directly.
type linuxStore struct {
	// AnchorDir is the directory CA-trust source PEMs go into.
	AnchorDir string
	// Suffix is the file extension expected by the distribution
	// (".crt" for Debian, ".pem" for RHEL/SUSE).
	Suffix string
	// Label is a human-readable name used in Result.Stores.
	Label string
	// RefreshLabel is the command string shown in error messages
	// (e.g. "update-ca-trust extract").
	RefreshLabel string
	// Refresh runs the post-write bundle refresh and returns the
	// combined stdout/stderr for diagnostics.
	Refresh func(context.Context) ([]byte, error)
}

// detectLinuxStore inspects the filesystem and returns the trust-store
// layout the current host uses, or false if none is recognized.
func detectLinuxStore() (linuxStore, bool) {
	candidates := []linuxStore{
		{
			AnchorDir:    "/usr/local/share/ca-certificates",
			Suffix:       ".crt",
			Label:        "Debian/Ubuntu system CA bundle",
			RefreshLabel: "update-ca-certificates",
			Refresh: func(ctx context.Context) ([]byte, error) {
				return exec.CommandContext(ctx, "update-ca-certificates").CombinedOutput()
			},
		},
		{
			AnchorDir:    "/etc/pki/ca-trust/source/anchors",
			Suffix:       ".pem",
			Label:        "RHEL/Fedora system CA bundle",
			RefreshLabel: "update-ca-trust extract",
			Refresh: func(ctx context.Context) ([]byte, error) {
				return exec.CommandContext(ctx, "update-ca-trust", "extract").CombinedOutput()
			},
		},
		{
			AnchorDir:    "/etc/ca-certificates/trust-source/anchors",
			Suffix:       ".crt",
			Label:        "Arch system CA bundle",
			RefreshLabel: "trust extract-compat",
			Refresh: func(ctx context.Context) ([]byte, error) {
				return exec.CommandContext(ctx, "trust", "extract-compat").CombinedOutput()
			},
		},
		{
			AnchorDir:    "/usr/share/pki/trust/anchors",
			Suffix:       ".pem",
			Label:        "openSUSE system CA bundle",
			RefreshLabel: "update-ca-certificates",
			Refresh: func(ctx context.Context) ([]byte, error) {
				return exec.CommandContext(ctx, "update-ca-certificates").CombinedOutput()
			},
		},
	}
	for _, c := range candidates {
		if pathExists(c.AnchorDir) {
			return c, true
		}
	}
	return linuxStore{}, false
}

// anchorPath is the destination path for the seed CA inside the host's
// anchor directory.
func (s linuxStore) anchorPath() string {
	return filepath.Join(s.AnchorDir, "seed-root"+s.Suffix)
}

// trustAnchorMode is the on-disk permission set the system bundle
// refresh tooling expects: world-readable, owner-writable. The file
// contents are a public certificate and contain no secret material.
const trustAnchorMode os.FileMode = 0o644

// writeAnchorFile writes pem to dst with trustAnchorMode. The write
// happens in two steps: create at 0o600 (so the file does not exist
// world-readable while the contents are being written), then chmod to
// trustAnchorMode so update-ca-certificates / update-ca-trust can
// find and read it.
//
// dst is bounded under root via [filepath.Rel] before any I/O. Both
// arguments originate from detectLinuxStore's hardcoded candidate
// list, so the bound is a defense-in-depth check rather than a
// security-critical gate, but it also forces [filepath.Clean]
// normalization right before [os.WriteFile] which is what gosec G703
// requires to clear path-traversal taint analysis.
func writeAnchorFile(root, dst string, pem []byte) (string, error) {
	cleanRoot := filepath.Clean(root)
	cleanDst := filepath.Clean(dst)
	rel, relErr := filepath.Rel(cleanRoot, cleanDst)
	if relErr != nil {
		return "", fmt.Errorf("resolve %s against %s: %w", dst, root, relErr)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("anchor path %q escapes %q", dst, root)
	}
	if err := os.WriteFile(cleanDst, pem, 0o600); err != nil {
		return "", fmt.Errorf("write %s: %w", cleanDst, err)
	}
	if err := os.Chmod(cleanDst, trustAnchorMode); err != nil {
		return "", fmt.Errorf("chmod %s: %w", cleanDst, err)
	}
	return cleanDst, nil
}

func installPlatform(ctx context.Context, certPath string) (Result, error) {
	store, ok := detectLinuxStore()
	if !ok {
		return Result{}, errors.New(
			"no supported system CA directory found " +
				"(tried /usr/local/share/ca-certificates, /etc/pki/ca-trust/source/anchors, " +
				"/etc/ca-certificates/trust-source/anchors, /usr/share/pki/trust/anchors)")
	}

	// #nosec G304 -- certPath is operator-supplied (validated by
	// ValidateCertFile in the caller) and the destination is a fixed
	// system directory.
	pemBytes, readErr := os.ReadFile(certPath)
	if readErr != nil {
		return Result{}, fmt.Errorf("read certificate: %w", readErr)
	}

	dst, writeErr := writeAnchorFile(store.AnchorDir, store.anchorPath(), pemBytes)
	if writeErr != nil {
		return Result{}, writeErr
	}

	out, refreshErr := store.Refresh(ctx)
	if refreshErr != nil {
		return Result{}, fmt.Errorf(
			"%s: %s: %w", store.RefreshLabel, trimError(out), refreshErr)
	}

	return Result{
		Stores: []string{store.Label + " (" + dst + ")"},
	}, nil
}

func uninstallPlatform(ctx context.Context, _ string) (Result, error) {
	store, ok := detectLinuxStore()
	if !ok {
		return Result{}, errors.New("no supported system CA directory found")
	}
	dst := store.anchorPath()
	res := Result{}
	if pathExists(dst) {
		if err := os.Remove(dst); err != nil {
			return Result{}, fmt.Errorf("remove %s: %w", dst, err)
		}
		res.Stores = append(res.Stores, store.Label+" ("+dst+")")
	} else {
		res.Skipped = append(res.Skipped, store.Label+": "+dst+" not present")
	}

	out, refreshErr := store.Refresh(ctx)
	if refreshErr != nil {
		return res, fmt.Errorf(
			"%s: %s: %w", store.RefreshLabel, trimError(out), refreshErr)
	}
	return res, nil
}
