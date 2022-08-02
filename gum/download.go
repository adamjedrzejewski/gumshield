package gum

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

func downloadSources(sources []string, dir string) error {
	for _, source := range sources {
		_, outPath := filepath.Split(source)
		outPath = filepath.Join(dir, outPath)
		if _, err := os.Stat(outPath); os.IsNotExist(err) {
			if err := downloadFile(source, outPath); err != nil {
				return err
			}
		}
	}

	return nil
}

func downloadFile(url, outPath string) error {
	out, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer out.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
