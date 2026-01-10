package update

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// Service coordinates update checking, downloading, and application.
type Service struct {
	config     UpdateConfig
	checker    *Checker
	downloader *Downloader
	applier    *Applier
	logger     *slog.Logger

	mu           sync.RWMutex
	status       UpdateStatus
	lastCheck    time.Time
	updateInfo   *UpdateInfo
	downloadPath string

	// For automatic checking
	stopChan chan struct{}
	running  bool
}

// NewService creates a new update service.
func NewService(config UpdateConfig, logger *slog.Logger) (*Service, error) {
	// Create download directory
	downloadDir := filepath.Join(os.TempDir(), "seed-updates")

	downloader, err := NewDownloader(downloadDir)
	if err != nil {
		return nil, fmt.Errorf("create downloader: %w", err)
	}

	applier, err := NewApplier()
	if err != nil {
		return nil, fmt.Errorf("create applier: %w", err)
	}

	return &Service{
		config:     config,
		checker:    NewChecker(config),
		downloader: downloader,
		applier:    applier,
		logger:     logger,
		status: UpdateStatus{
			State: StateIdle,
		},
		stopChan: make(chan struct{}),
	}, nil
}

// Start begins automatic update checking if enabled.
func (s *Service) Start(ctx context.Context) error {
	if !s.config.Enabled {
		return nil
	}

	s.mu.Lock()
	if s.running {
		s.mu.Unlock()
		return nil
	}
	s.running = true
	s.mu.Unlock()

	// Initial check
	go func() {
		// Delay initial check by 1 minute to let the app fully start
		select {
		case <-time.After(time.Minute):
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		}

		s.checkForUpdate(ctx)
	}()

	// Periodic checks
	go s.periodicCheck(ctx)

	return nil
}

// Stop stops automatic update checking.
func (s *Service) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return
	}

	close(s.stopChan)
	s.running = false
}

// periodicCheck runs update checks at the configured interval.
func (s *Service) periodicCheck(ctx context.Context) {
	ticker := time.NewTicker(s.config.CheckInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.checkForUpdate(ctx)
		case <-ctx.Done():
			return
		case <-s.stopChan:
			return
		}
	}
}

// checkForUpdate performs an update check and handles auto-download if enabled.
func (s *Service) checkForUpdate(ctx context.Context) {
	info, err := s.CheckForUpdate(ctx)
	if err != nil {
		if s.logger != nil {
			s.logger.WarnContext(ctx, "Update check failed", "error", err)
		}
		return
	}

	if info.Available && s.config.AutoDownload {
		_, downloadErr := s.DownloadUpdate(ctx, nil)
		if downloadErr != nil {
			if s.logger != nil {
				s.logger.WarnContext(ctx, "Auto-download failed", "error", downloadErr)
			}
		}
	}
}

// CheckForUpdate checks for available updates.
func (s *Service) CheckForUpdate(ctx context.Context) (*UpdateInfo, error) {
	s.mu.Lock()
	s.status = UpdateStatus{
		State:     StateChecking,
		StartedAt: time.Now(),
		Message:   "Checking for updates",
	}
	s.mu.Unlock()

	info, err := s.checker.CheckForUpdate(ctx)
	if err != nil {
		s.mu.Lock()
		s.status.State = StateFailed
		s.status.Error = err.Error()
		s.mu.Unlock()
		return nil, err
	}

	s.mu.Lock()
	s.updateInfo = info
	s.lastCheck = time.Now()
	s.status.State = StateIdle
	s.status.Message = ""
	s.mu.Unlock()

	if s.logger != nil {
		if info.Available {
			s.logger.InfoContext(ctx, "Update available",
				"current", info.CurrentVersion,
				"latest", info.LatestVersion)
		} else {
			s.logger.DebugContext(ctx, "No update available",
				"current", info.CurrentVersion)
		}
	}

	return info, nil
}

// GetUpdateInfo returns the latest update information.
func (s *Service) GetUpdateInfo() *UpdateInfo {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.updateInfo
}

// DownloadUpdate downloads the available update.
func (s *Service) DownloadUpdate(ctx context.Context, progressCb func(DownloadProgress)) (string, error) {
	s.mu.RLock()
	info := s.updateInfo
	s.mu.RUnlock()

	if info == nil || !info.Available {
		return "", errors.New("no update available")
	}

	if info.DownloadURL == "" {
		return "", errors.New("no download URL for this platform")
	}

	// Download the binary
	path, err := s.downloader.Download(ctx, info.DownloadURL, progressCb)
	if err != nil {
		s.mu.Lock()
		s.status.State = StateFailed
		s.status.Error = err.Error()
		s.mu.Unlock()
		return "", fmt.Errorf("download: %w", err)
	}

	// Verify checksum if available
	if info.ChecksumURL != "" {
		if verifyErr := s.downloader.VerifyChecksum(path, info.ChecksumURL); verifyErr != nil {
			// Remove failed download
			_ = os.Remove(path)
			s.mu.Lock()
			s.status.State = StateFailed
			s.status.Error = verifyErr.Error()
			s.mu.Unlock()
			return "", fmt.Errorf("verify checksum: %w", verifyErr)
		}
	}

	s.mu.Lock()
	s.downloadPath = path
	s.status.State = StateIdle
	s.status.Message = "Update downloaded and verified"
	s.mu.Unlock()

	if s.logger != nil {
		s.logger.InfoContext(ctx, "Update downloaded", "path", path)
	}

	return path, nil
}

// ApplyUpdate applies the downloaded update.
func (s *Service) ApplyUpdate(ctx context.Context) error {
	s.mu.RLock()
	downloadPath := s.downloadPath
	s.mu.RUnlock()

	if downloadPath == "" {
		return errors.New("no update downloaded")
	}

	if err := s.applier.Apply(ctx, downloadPath); err != nil {
		return fmt.Errorf("apply update: %w", err)
	}

	if s.logger != nil {
		s.logger.InfoContext(ctx, "Update applied successfully")
	}

	return nil
}

// RestartService restarts the service to use the new binary.
func (s *Service) RestartService(ctx context.Context) error {
	return s.applier.RestartService(ctx)
}

// Rollback reverts to the previous version.
func (s *Service) Rollback() error {
	if err := s.applier.Rollback(); err != nil {
		return fmt.Errorf("rollback: %w", err)
	}

	if s.logger != nil {
		s.logger.Info("Rolled back to previous version")
	}

	return nil
}

// GetStatus returns the current update service status.
func (s *Service) GetStatus() UpdateStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	// Merge status from sub-components
	if s.status.State == StateDownloading {
		return s.downloader.GetStatus()
	}
	if s.status.State == StateApplying || s.status.State == StateRestarting {
		return s.applier.GetStatus()
	}

	return s.status
}

// GetLastCheckTime returns when the last update check occurred.
func (s *Service) GetLastCheckTime() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.lastCheck
}

// IsUpdateDownloaded returns true if an update has been downloaded.
func (s *Service) IsUpdateDownloaded() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.downloadPath != ""
}

// CleanupDownloads removes downloaded update files.
func (s *Service) CleanupDownloads() error {
	s.mu.Lock()
	s.downloadPath = ""
	s.mu.Unlock()

	return s.downloader.CleanupDownloads()
}

// GetConfig returns the current update configuration.
func (s *Service) GetConfig() UpdateConfig {
	return s.config
}

// SetConfig updates the configuration.
func (s *Service) SetConfig(config UpdateConfig) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.config = config
	s.checker = NewChecker(config)
}
