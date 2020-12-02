package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type getGit struct{}

func (g getGit) Get(m module, dst string) (string, error) {
	dir, err := ioutil.TempDir("/tmp", "terrafile")
	version := m.Version
	if m.Version == "" {
		version = "master"
	}
	if err != nil {
		return "", fmt.Errorf("Error creating temp folder at /tmp: %v", err)
	}
	defer os.RemoveAll(dir)
	log.Printf("Getting git module version: %s from %s\n", version, m.Source)
	cmd := exec.Command("git", "clone", "--single-branch", "--depth=1", "-b", version, m.Source, dir)
	err = cmd.Run()
	if err != nil {
		return "", fmt.Errorf("Error cloning module from git. repo: %s, version: %s: %v", m.Source, version, err)
	}
	cmd = exec.Command("git", "rev-parse", "HEAD")
	cmd.Dir = dir
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error reading local git ref at %s: %v", dir, err)
	}
	sfPath := filepath.Join(dir, m.Subfolder)
	os.RemoveAll(dst)
	if err = os.Rename(sfPath, dst); err != nil {
		return "", fmt.Errorf("Error copying repo subfolder: %v", err)
	}
	return out.String(), nil
}

func (g getGit) GetState(m module) (string, error) {
	version := m.Version
	if m.Version == "" {
		version = "master"
	}
	cmd := exec.Command("git", "ls-remote", "--quiet", m.Source, version, version+"^{}")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("Error reading git remote %s: %v", m.Source, err)
	}
	if out.String() == "" {
		return "", fmt.Errorf("Remote ref not found: %s, ref %v", m.Source, version)
	}
	lines := strings.Split(out.String(), "\n")
	line := lines[0]
	if len(lines) > 1 {
		for _, i := range lines {
			if strings.HasSuffix(i, "^{}") {
				line = i
			}
		}
	}
	commit := strings.Fields(line)[0]
	return commit, nil
}
