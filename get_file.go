package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/otiai10/copy"
	log "github.com/sirupsen/logrus"
)

type getFile struct{}

func (g getFile) Get(m module, dst string) (string, error) {
	log.Infof("Getting local module from %s to %s", m.Source, dst)
	os.RemoveAll(dst)
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
	h := md5.New()
	err := filepath.Walk(path, func(file string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		f, err := os.Open(file)
		if err != nil {
			return err
		}
		hf := md5.New()
		_, err = io.Copy(hf, f)
		if err != nil {
			return err
		}
		// extract relative path only
		rp := file
		if path != "." {
			rp = file[len(path)+1:]
		}
		fmt.Fprintf(h, "%x  %s\n", hf.Sum(nil), rp)
		return nil
	})
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
