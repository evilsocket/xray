package main

import (
	"net/http"
	"strings"

	assetfs "github.com/elazarl/go-bindata-assetfs"
)

type BFS struct {
	fs http.FileSystem
}

func (b *BFS) Open(name string) (http.File, error) {
	return b.fs.Open(name)
}

func (b *BFS) Exists(prefix string, filepath string) bool {
	if p := strings.TrimPrefix(filepath, prefix); len(p) < len(filepath) {
		if _, err := b.fs.Open(p); err != nil {
			return false
		}
		return true
	}
	return false
}

func NewBFS(root string) *BFS {
	fs := &assetfs.AssetFS{Asset, AssetDir, AssetInfo, root, ""}
	return &BFS{
		fs,
	}
}
