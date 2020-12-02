package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetGitGetState(t *testing.T) {
	mm := map[string]module{
		"22f20976951cccf8ea7189817dd5686198572548": {
			Source:  "git@github.com:terraform-aws-modules/terraform-aws-vpc",
			Version: "v2.61.0",
		},
		"b0a5f0aadadf9f027cac4e07e15df6108ed9820a": {
			Source:  "https://github.com/kaminskip88/test_repo.git",
			Version: "immutable",
		},
		// TODO: add more cases
	}

	g := new(getGit)
	for k, m := range mm {
		r, err := g.GetState(m)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, k, r, "GetState() result not matched expected", m.Source, m.Version)
	}
}

func TestGetGitGet(t *testing.T) {
	type testCase struct {
		filePath    string
		fileContent string
		module      module
	}
	ts := []testCase{
		{
			filePath:    "dir1/abc",
			fileContent: "data for abc\n",
			module: module{
				Name:    "test_full",
				Source:  "https://github.com/kaminskip88/test_repo.git",
				Version: "1.0.0",
			},
		},
		{
			filePath:    "abc",
			fileContent: "data for abc\n",
			module: module{
				Name:      "test_dir1",
				Source:    "https://github.com/kaminskip88/test_repo.git",
				Version:   "1.0.0",
				Subfolder: "dir1",
			},
		},
		{
			filePath:    "21abc",
			fileContent: "data for 21abc\n",
			module: module{
				Name:      "test_dir2",
				Source:    "https://github.com/kaminskip88/test_repo.git",
				Version:   "1.0.0",
				Subfolder: "dir2/dir1",
			},
		},
	}

	dir, err := ioutil.TempDir("/tmp", "go-test-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)
	g := new(getGit)
	for _, m := range ts {
		_, err := g.Get(m.module, filepath.Join(dir, m.module.Name))
		if err != nil {
			t.Fatal(err)
		}
		f := readFile(t, filepath.Join(dir, m.module.Name, m.filePath))
		assert.Equal(t, m.fileContent, f, "GetState() file content not matched expected")
	}
}
