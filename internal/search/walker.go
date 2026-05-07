// Package search
package search

import (
	"bytes"
	"io/fs"
	"os"
	"path/filepath"
	"slices"

	"github.com/charlievieth/fastwalk"
)

var skipDirs = []string{".git", "vendor", "node_modules"}

func StartWalking(root string) {
	go func() {
		defer close(pathCh)
		_ = fastwalk.Walk(nil, root, func(path string, d fs.DirEntry, err error) error {
			if d.IsDir() && slices.Contains(skipDirs, d.Name()) {
				return filepath.SkipDir
			}

			if err != nil || d.IsDir() || filepath.Base(path) == ".ff" || isBinaryFile(path) {
				return nil
			}

			pathCh <- path
			return nil
		})
	}()
}

// We don't want to process binary files. Read the first 1024 bytes of the file.
// If any null bytes are present, it's probably binary. Any errors mean we can't
// read the file, return true since there's no benefit in further processing.
func isBinaryFile(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return true
	}
	defer f.Close()

	buf := make([]byte, 1024)
	n, err := f.Read(buf)
	if err != nil && n == 0 {
		return true
	}

	return bytes.IndexByte(buf[:n], 0) != -1
}

var pathCh = make(chan string, 8)
