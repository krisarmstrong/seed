package update

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	// downloadBufferSize is the buffer size for download operations.
	downloadBufferSize = 32 * 1024 // 32KB
	// maxDownloadSize is the maximum allowed download size (500MB).
	maxDownloadSize = 500 * 1024 * 1024
	// downloadTimeout is the timeout for large file downloads.
	downloadTimeout = 30 * time.Minute
	// checksumTimeout is the timeout for checksum file downloads.
	checksumTimeout = 30 * time.Second
	// maxBodyReadSize is the maximum size for error body reads.
	maxBodyReadSize = 1024
	// percentMultiplier is the multiplier to convert ratio to percentage.
	percentMultiplier = 100
)

// DownloadProgress reports download progress.
type DownloadProgress struct {
	Downloaded int64
	Total      int64
	Percent    float64
}

// Downloader handles downloading and verifying update binaries.
type Downloader struct {
	httpClient  *http.Client
	downloadDir string

	mu       sync.RWMutex
	progress DownloadProgress
	status   UpdateStatus
}

// NewDownloader creates a new downloader with the given download directory.
func NewDownloader(downloadDir string) (*Downloader, error) {
	// Ensure download directory exists
	if err := os.MkdirAll(downloadDir, 0o750); err != nil {
		return nil, fmt.Errorf("create download dir: %w", err)
	}

	return &Downloader{
		httpClient: &http.Client{
			Timeout: downloadTimeout,
		},
		downloadDir: downloadDir,
		status: UpdateStatus{
			State: StateIdle,
		},
	}, nil
}

// Download downloads a file from the given URL and saves it to the download directory.
// Returns the path to the downloaded file.
//
//nolint:gocognit,funlen // Download operations inherently require multiple steps and error handling.
func (d *Downloader) Download(ctx context.Context, url string, progressCb func(DownloadProgress)) (string, error) {
	d.mu.Lock()
	d.status = UpdateStatus{
		State:     StateDownloading,
		StartedAt: time.Now(),
	}
	d.mu.Unlock()

	// Create request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		d.setError(fmt.Errorf("create request: %w", err))
		return "", err
	}

	// Execute request
	resp, err := d.httpClient.Do(req)
	if err != nil {
		d.setError(fmt.Errorf("download request: %w", err))
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		statusErr := fmt.Errorf("download failed: HTTP %d", resp.StatusCode)
		d.setError(statusErr)
		return "", statusErr
	}

	// Check content length
	totalSize := resp.ContentLength
	if totalSize > maxDownloadSize {
		sizeErr := fmt.Errorf("file too large: %d bytes (max %d)", totalSize, maxDownloadSize)
		d.setError(sizeErr)
		return "", sizeErr
	}

	// Update status with total size
	d.mu.Lock()
	d.status.TotalBytes = totalSize
	d.progress.Total = totalSize
	d.mu.Unlock()

	// Create destination file
	filename := filepath.Base(url)
	destPath := filepath.Join(d.downloadDir, filename)

	file, err := os.Create(destPath)
	if err != nil {
		d.setError(fmt.Errorf("create file: %w", err))
		return "", err
	}
	defer file.Close()

	// Download with progress tracking
	buf := make([]byte, downloadBufferSize)
	var downloaded int64

	for {
		select {
		case <-ctx.Done():
			// Clean up partial download
			_ = file.Close()
			_ = os.Remove(destPath)
			ctxErr := ctx.Err()
			d.setError(ctxErr)
			return "", ctxErr
		default:
		}

		n, readErr := resp.Body.Read(buf)
		if n > 0 {
			_, writeErr := file.Write(buf[:n])
			if writeErr != nil {
				d.setError(fmt.Errorf("write file: %w", writeErr))
				return "", writeErr
			}

			downloaded += int64(n)

			// Update progress
			d.mu.Lock()
			d.status.DownloadedBytes = downloaded
			d.progress.Downloaded = downloaded
			if totalSize > 0 {
				d.progress.Percent = float64(downloaded) / float64(totalSize) * percentMultiplier
				d.status.Progress = d.progress.Percent
			}
			d.mu.Unlock()

			// Notify callback
			if progressCb != nil {
				progressCb(d.progress)
			}
		}

		if readErr == io.EOF {
			break
		}
		if readErr != nil {
			d.setError(fmt.Errorf("read response: %w", readErr))
			return "", readErr
		}
	}

	d.mu.Lock()
	d.status.State = StateVerifying
	d.status.Progress = progressComplete
	d.mu.Unlock()

	return destPath, nil
}

// VerifyChecksum verifies the SHA256 checksum of a file.
func (d *Downloader) VerifyChecksum(filePath, checksumURL string) error {
	// Download checksum file
	ctx, cancel := context.WithTimeout(context.Background(), checksumTimeout)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, checksumURL, nil)
	if err != nil {
		return fmt.Errorf("create checksum request: %w", err)
	}

	resp, err := d.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("download checksum: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("checksum download failed: HTTP %d", resp.StatusCode)
	}

	// Read checksum content
	checksumContent, err := io.ReadAll(io.LimitReader(resp.Body, maxBodyReadSize))
	if err != nil {
		return fmt.Errorf("read checksum: %w", err)
	}

	// Parse expected checksum (format: "checksum  filename" or just "checksum")
	expectedChecksum := strings.Fields(string(checksumContent))[0]
	expectedChecksum = strings.TrimSpace(expectedChecksum)

	// Calculate actual checksum
	actualChecksum, err := d.calculateSHA256(filePath)
	if err != nil {
		return fmt.Errorf("calculate checksum: %w", err)
	}

	// Compare
	if !strings.EqualFold(actualChecksum, expectedChecksum) {
		return fmt.Errorf("checksum mismatch: expected %s, got %s", expectedChecksum, actualChecksum)
	}

	return nil
}

// calculateSHA256 calculates the SHA256 checksum of a file.
func (d *Downloader) calculateSHA256(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, copyErr := io.Copy(hash, file); copyErr != nil {
		return "", copyErr
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

// GetStatus returns the current download status.
func (d *Downloader) GetStatus() UpdateStatus {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.status
}

// GetProgress returns the current download progress.
func (d *Downloader) GetProgress() DownloadProgress {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.progress
}

// setError sets an error state.
func (d *Downloader) setError(err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.status.State = StateFailed
	d.status.Error = err.Error()
}

// CleanupDownloads removes old downloaded files.
func (d *Downloader) CleanupDownloads() error {
	entries, err := os.ReadDir(d.downloadDir)
	if err != nil {
		return fmt.Errorf("read download dir: %w", err)
	}

	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		entryPath := filepath.Join(d.downloadDir, entry.Name())
		// Ignore errors on individual file removal
		_ = os.Remove(entryPath)
	}

	return nil
}
