package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateSetGet(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "terrafile-test-state-set-get")
	assert.Nil(t, err, "error creating temp folder")
	defer os.RemoveAll(dir)
	sf := filepath.Join(dir, "terrafile.state")
	s := newState(sf)
	tv := "111222333444555666777"
	s.set("mod1", tv)
	assert.Equal(t, tv, s.state["mod1"], "state not matched expected")
	assert.Equal(t, s.get("mod1"), tv, "get() result not matched expected")
}

func TestStateSave(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "terrafile-test-save")
	assert.Nil(t, err, "error creating temp folder")
	defer os.RemoveAll(dir)
	sf := filepath.Join(dir, "terrafile.state")
	s := newState(sf)
	s.set("mod1", "111222333444555666777")
	s.set("mod2", "777666555444333222111")
	expected := `{"mod1":"111222333444555666777","mod2":"777666555444333222111"}`
	s.save()
	assert.Equal(t, expected, readFile(t, sf), "file content don't match expected")
}

func TestStateLoad(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "terrafile-test-load")
	assert.Nil(t, err, "error creating temp folder")
	defer os.RemoveAll(dir)
	sf := filepath.Join(dir, "terrafile.state")
	tf := `{"mod1":"111222333444555666777","mod2":"777666555444333222111"}`
	writeFile(t, sf, tf)
	s := newState(sf)
	s.load()
	r1 := s.get("mod1")
	r2 := s.get("mod2")
	assert.Equal(t, "111222333444555666777", r1)
	assert.Equal(t, "777666555444333222111", r2)
}
