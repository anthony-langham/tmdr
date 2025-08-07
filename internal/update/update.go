package update

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/anthonylangham/tmdr/internal/version"
)

const (
	githubAPIURL = "https://api.github.com/repos/anthony-langham/tmdr/releases/latest"
	timeout      = 3 * time.Second
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	HTMLURL string `json:"html_url"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// CheckForUpdate checks if a new version is available on GitHub
func CheckForUpdate() (bool, string, string, error) {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(githubAPIURL)
	if err != nil {
		return false, "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false, "", "", fmt.Errorf("failed to check for updates: %s", resp.Status)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return false, "", "", err
	}

	// Remove 'v' prefix for comparison
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := version.Version

	// Simple version comparison (works for semantic versioning)
	if compareVersions(latestVersion, currentVersion) > 0 {
		return true, latestVersion, release.HTMLURL, nil
	}

	return false, "", "", nil
}

// compareVersions compares two semantic version strings
// Returns: 1 if v1 > v2, -1 if v1 < v2, 0 if equal
func compareVersions(v1, v2 string) int {
	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	for i := 0; i < len(parts1) && i < len(parts2); i++ {
		var n1, n2 int
		fmt.Sscanf(parts1[i], "%d", &n1)
		fmt.Sscanf(parts2[i], "%d", &n2)

		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}

	return len(parts1) - len(parts2)
}

// UpdateInfo contains information about an available update
type UpdateInfo struct {
	Available   bool
	Version     string
	URL         string
	DownloadURL string
	AssetName   string
}

// CheckForUpdateAsync checks for updates in the background
func CheckForUpdateAsync() <-chan UpdateInfo {
	ch := make(chan UpdateInfo, 1)

	go func() {
		info := CheckForUpdateWithAssets()
		ch <- info
	}()

	return ch
}

// CheckForUpdateWithAssets checks for updates and finds the right asset to download
func CheckForUpdateWithAssets() UpdateInfo {
	client := &http.Client{
		Timeout: timeout,
	}

	resp, err := client.Get(githubAPIURL)
	if err != nil {
		return UpdateInfo{Available: false}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return UpdateInfo{Available: false}
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return UpdateInfo{Available: false}
	}

	// Remove 'v' prefix for comparison
	latestVersion := strings.TrimPrefix(release.TagName, "v")
	currentVersion := version.Version

	// Check if update is available
	if compareVersions(latestVersion, currentVersion) <= 0 {
		return UpdateInfo{Available: false}
	}

	// Find the right asset for this platform
	assetName := getAssetName()
	var downloadURL string
	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, assetName) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}

	return UpdateInfo{
		Available:   true,
		Version:     latestVersion,
		URL:         release.HTMLURL,
		DownloadURL: downloadURL,
		AssetName:   assetName,
	}
}

// getAssetName returns the asset name for the current platform
func getAssetName() string {
	os := runtime.GOOS
	arch := runtime.GOARCH
	
	// Map to our release naming convention
	if os == "darwin" {
		if arch == "arm64" {
			return "darwin-arm64"
		}
		return "darwin-amd64"
	} else if os == "linux" {
		if arch == "arm64" {
			return "linux-arm64"
		}
		return "linux-amd64"
	} else if os == "windows" {
		return "windows-amd64"
	}
	return ""
}

// DownloadUpdate downloads the update to a temporary file
func DownloadUpdate(downloadURL string, onProgress func(downloaded, total int64)) (string, error) {
	if downloadURL == "" {
		return "", fmt.Errorf("no download URL available")
	}

	// Download file
	resp, err := http.Get(downloadURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("download failed: %s", resp.Status)
	}

	// Get the size
	size := resp.ContentLength

	// Create a progress reader
	pr := &progressReader{
		Reader:     resp.Body,
		Total:      size,
		OnProgress: onProgress,
	}

	// Handle different file types
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, "tmdr-update")
	
	if strings.HasSuffix(downloadURL, ".tar.gz") {
		// Extract from tar.gz
		extractedPath, err := extractFromTarGz(pr)
		if err != nil {
			return "", err
		}
		return extractedPath, nil
	} else if strings.HasSuffix(downloadURL, ".zip") {
		// Extract from zip (Windows)
		extractedPath, err := extractFromZip(pr)
		if err != nil {
			return "", err
		}
		return extractedPath, nil
	} else {
		// Direct binary download
		if runtime.GOOS == "windows" {
			tempFile += ".exe"
		}

		out, err := os.Create(tempFile)
		if err != nil {
			return "", err
		}
		defer out.Close()

		_, err = io.Copy(out, pr)
		if err != nil {
			return "", err
		}

		// Make executable on Unix systems
		if runtime.GOOS != "windows" {
			err = os.Chmod(tempFile, 0755)
			if err != nil {
				return "", err
			}
		}

		return tempFile, nil
	}
}

type progressReader struct {
	io.Reader
	Total      int64
	Downloaded int64
	OnProgress func(downloaded, total int64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.Downloaded += int64(n)
	if pr.OnProgress != nil {
		pr.OnProgress(pr.Downloaded, pr.Total)
	}
	return n, err
}

// extractFromTarGz extracts the binary from a tar.gz archive
func extractFromTarGz(r io.Reader) (string, error) {
	// Read all data first
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	// Create gzip reader
	gr, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	defer gr.Close()

	// Create tar reader
	tr := tar.NewReader(gr)

	// Find and extract the binary
	tempDir := os.TempDir()
	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}

		// Look for the tmdr binary
		if strings.Contains(header.Name, "tmdr") && !strings.Contains(header.Name, "/") {
			tempFile := filepath.Join(tempDir, "tmdr-update")
			if runtime.GOOS == "windows" {
				tempFile += ".exe"
			}

			out, err := os.Create(tempFile)
			if err != nil {
				return "", err
			}
			defer out.Close()

			_, err = io.Copy(out, tr)
			if err != nil {
				return "", err
			}

			// Make executable
			if runtime.GOOS != "windows" {
				err = os.Chmod(tempFile, 0755)
				if err != nil {
					return "", err
				}
			}

			return tempFile, nil
		}
	}

	return "", fmt.Errorf("tmdr binary not found in archive")
}

// extractFromZip extracts the binary from a zip archive
func extractFromZip(r io.Reader) (string, error) {
	// Read all data first
	data, err := io.ReadAll(r)
	if err != nil {
		return "", err
	}

	// Create zip reader
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return "", err
	}

	// Find and extract the binary
	tempDir := os.TempDir()
	for _, file := range zr.File {
		if strings.Contains(file.Name, "tmdr") && strings.HasSuffix(file.Name, ".exe") {
			rc, err := file.Open()
			if err != nil {
				return "", err
			}
			defer rc.Close()

			tempFile := filepath.Join(tempDir, "tmdr-update.exe")
			out, err := os.Create(tempFile)
			if err != nil {
				return "", err
			}
			defer out.Close()

			_, err = io.Copy(out, rc)
			if err != nil {
				return "", err
			}

			return tempFile, nil
		}
	}

	return "", fmt.Errorf("tmdr.exe not found in archive")
}

// InstallUpdate replaces the current binary with the downloaded update
func InstallUpdate(updatePath string) error {
	// Get current executable path
	exePath, err := os.Executable()
	if err != nil {
		return err
	}

	// Resolve any symlinks
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return err
	}

	// On Windows, we need to rename the old binary first
	if runtime.GOOS == "windows" {
		oldPath := exePath + ".old"
		_ = os.Remove(oldPath) // Remove any existing .old file
		if err := os.Rename(exePath, oldPath); err != nil {
			return err
		}
		defer os.Remove(oldPath) // Clean up after
	}

	// Move the new binary into place
	if err := os.Rename(updatePath, exePath); err != nil {
		// Try copying if rename fails (might be across filesystems)
		input, err := os.Open(updatePath)
		if err != nil {
			return err
		}
		defer input.Close()

		output, err := os.Create(exePath)
		if err != nil {
			return err
		}
		defer output.Close()

		_, err = io.Copy(output, input)
		if err != nil {
			return err
		}

		// Ensure it's executable
		if runtime.GOOS != "windows" {
			err = os.Chmod(exePath, 0755)
			if err != nil {
				return err
			}
		}

		// Remove the temp file
		os.Remove(updatePath)
	}

	return nil
}