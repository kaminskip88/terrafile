package main

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFileGetState(t *testing.T) {
	dir, err := ioutil.TempDir("/tmp", "go-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	if err := ioutil.WriteFile(filepath.Join(dir, "abc"), []byte("data for abc"), 0644); err != nil {
		t.Fatal(err)
	}
	if err := os.Mkdir(filepath.Join(dir, "def"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(dir, "def", "xyz"), []byte("data for xyz"), 0644); err != nil {
		t.Fatal(err)
	}
	h := md5.New()
	s := "b9cf3d37dc9d46eb380cfade176b151e  abc\ndafe17e57a5c37a89e1381f325714c78  def/xyz\n"
	fmt.Fprintf(h, s)
	exp := fmt.Sprintf("%x", h.Sum(nil))
	m := module{
		Source: dir,
	}
	g := new(getFile)
	r1, err := g.GetState(m)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, exp, r1, "GetState() result not matching test datsa")
	if err := ioutil.WriteFile(filepath.Join(dir, "def", "xyz"), []byte("NEW data for xyz"), 0644); err != nil {
		t.Fatal(err)
	}
	r2, err := g.GetState(m)
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, r1, r2, "GetState() returns same changes after file changed")
}

func TestGetFileGet(t *testing.T) {
	src, err := ioutil.TempDir("/tmp", "go-test-")
	dst, err := ioutil.TempDir("/tmp", "go-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(src)
	defer os.RemoveAll(dst)
	if err := os.Mkdir(filepath.Join(src, "abc"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := ioutil.WriteFile(filepath.Join(src, "abc", "xyz"), []byte("data for xyz"), 0644); err != nil {
		t.Fatal(err)
	}
	m := module{
		Source: src,
	}
	g := new(getFile)
	s, err := g.Get(m, dst)
	e, err := g.GetState(m)
	if err != nil {
		t.Fatal(err)
	}
	assert.Equal(t, e, s, "Get() returned state not matched expected")
	assert.FileExists(t, filepath.Join(dst, "abc", "xyz"))
	f := readFile(t, filepath.Join(dst, "abc", "xyz"))
	assert.Equal(t, "data for xyz", f, "Get() dst file not matching expected")
	// TODO: test relative paths (./ ../ ../../)
}
