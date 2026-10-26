package downloader

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

type ProgressReporter interface {
	Add(bytes int64)
	SetTotal(total int64)
}

type progressWriter struct {
	w        io.Writer
	reporter ProgressReporter
}

func (pr *progressWriter) Write(p []byte) (int, error) {
	n, err := pr.w.Write(p)
	if n > 0 {
		pr.reporter.Add(int64(n))
	}
	return n, err
}

func Download(ctx context.Context, client *http.Client, rawURL, outputPath string, reporter ProgressReporter) error {
	if _, err := os.Stat(outputPath); err == nil {
		return fmt.Errorf("output file already exists: %s", outputPath)
	} else if !errors.Is(err, os.ErrNotExist) {
		return fmt.Errorf("check output file: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected HTTP status: %s", resp.Status)
	}

	reporter.SetTotal(resp.ContentLength)

	partFile, err := os.CreateTemp(filepath.Dir(outputPath), ".godl-*.part")
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}

	partPath := partFile.Name()

	pr := progressWriter{w: partFile, reporter: reporter}

	_, copyErr := io.Copy(&pr, resp.Body)
	closeErr := partFile.Close()

	if copyErr != nil {
		_ = os.Remove(partPath)
		return fmt.Errorf("stream file: %w", copyErr)
	}

	if closeErr != nil {
		_ = os.Remove(partPath)
		return fmt.Errorf("close file: %w", closeErr)
	}

	if err := os.Link(partPath, outputPath); err != nil {
		_ = os.Remove(partPath)
		return fmt.Errorf("publish file: %w", err)
	}

	if err := os.Remove(partPath); err != nil {
		return fmt.Errorf("file saved, but temporary file cleanup failed: %w", err)
	}

	return nil

}
