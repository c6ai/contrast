package main

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"time"
)

func main() {
	downloader := []downloader{
		// &curler{}, // No OOM
		// &wgetter{}, // No OOM
		&native{}, // OOM
	}
	resources := []string{
		// "https://github.com/edgelesssys/uplosi/releases/download/v0.1.3/uplosi_0.1.3_linux_amd64.tar.gz", // OOM
		"http://wikipedia.org", // OOM
		// "http://foo.bar", // No OOM
	}
	for {
		for _, res := range resources {
			for _, d := range downloader {
				ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
				defer cancel()
				hash, err := d.Download(ctx, res)
				if err != nil {
					fmt.Printf("Failed to download %s using %T: %v\n", res, d, err)
					continue
				}
				fmt.Printf("Downloaded %s using %T: %s\n", res, d, hash)
			}
		}
		time.Sleep(5 * time.Second)
	}
}

type downloader interface {
	// Downloader downloads the content from the given URL and returns its hash.
	Download(ctx context.Context, url string) (string, error)
}

type curler struct{}

func (c *curler) Download(ctx context.Context, url string) (string, error) {
	args := []string{
		"-f",      // Fail silently
		"-s",      // Silent mode
		"-S",      // Show error
		"-L",      // Follow redirects
		"-o", "-", // Write output to stdout
		"--max-time", "10", // Timeout
		url,
	}
	cmd := exec.CommandContext(ctx, "curl", args...)
	out, err := cmd.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return "", fmt.Errorf("curl failed with exit code %d: %s", exitErr.ExitCode(), exitErr.Stderr)
	} else if err != nil {
		return "", fmt.Errorf("curl failed: %w", err)
	}
	hash := sha256.Sum256(out)
	return hex.EncodeToString(hash[:]), nil
}

type wgetter struct{}

func (w *wgetter) Download(ctx context.Context, url string) (string, error) {
	args := []string{
		"-O-", // Write output to stdout
		url,
	}
	cmd := exec.CommandContext(ctx, "wget", args...)
	out, err := cmd.Output()
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		return "", fmt.Errorf("wget failed with exit code %d: %s", exitErr.ExitCode(), exitErr.Stderr)
	} else if err != nil {
		return "", fmt.Errorf("wget failed: %w", err)
	}
	hash := sha256.Sum256(out)
	return hex.EncodeToString(hash[:]), nil
}

type native struct{}

func (n *native) Download(ctx context.Context, url string) (string, error) {
	req, err := http.NewRequestWithContext(context.TODO(), http.MethodGet, url, http.NoBody)
	if err != nil {
		return "", fmt.Errorf("native http request failed to create: %w", err)
	}
	resp, err := http.DefaultClient.Do(req) // OOMs
	if err != nil {
		return "", fmt.Errorf("native http request failed: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("native http request failed with status code %d", resp.StatusCode)
	}
	hash := sha256.New()
	if _, err := io.Copy(hash, resp.Body); err != nil {
		return "", fmt.Errorf("native http request failed to read response body: %w", err)
	}
	return hex.EncodeToString(hash.Sum(nil)[:]), nil
}
