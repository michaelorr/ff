// Package search
package search

import (
	"io/fs"
	"path/filepath"

	"github.com/charlievieth/fastwalk"
)

func StartWalking(root string) {
	go func() {
		defer close(pathCh)
		_ = fastwalk.Walk(nil, root, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}
			if filepath.Base(path) == ".ff" {
				return nil
			}
			itemCh <- path
			return nil
		})
	}()
}

var pathCh = make(chan string, 8)
