package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
)

type getFile struct{}

func (g getFile) Get(m module, dst string) (string, error) {
	// log.Printf("Getting local module from %s to %s\n", src, dst)
	if err := copy.Copy(m.Source, dst); err != nil {
		return "", fmt.Errorf("Error copying module from local path %s: %v", m.Source, err)
	}
	h, err := g.GetState(m)
	if err != nil {
		return "", fmt.Errorf("Error calculating checksum for module %s: %v", m.Source, err)
	}
	return h, nil

}

func (g getFile) GetState(m module) (string, error) {
	r, err := hashDir(m.Source)
	if err != nil {
		return "", fmt.Errorf("Error")
	}
	return r, nil
}

func hashDir(path string) (string, error) {
	hash := md5.New()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		io.WriteString(hash, path)
		return nil
	})
	if err != nil {
		return "", nil
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
