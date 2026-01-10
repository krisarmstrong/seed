// Package update provides in-app update functionality for Seed.
//
// Features:
//   - Check GitHub releases API for new versions
//   - Download and verify binaries with SHA256 checksums
//   - Apply updates with automatic rollback on failure
//   - Configurable update check intervals
//
// Security:
//   - HTTPS only for all update operations
//   - SHA256 checksum verification
//   - Binary backup before update for rollback
//   - No automatic updates without user consent
package update
