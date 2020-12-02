package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func readFile(t *testing.T, path string) string {
	if _, err := os.Stat(path); err != nil {
		assert.FailNow(t, "File not exists: %s", path)
	}
	f, err := ioutil.ReadFile(path)
	if err != nil {
		assert.FailNow(t, "can't read file: %s", path)
	}
	return string(f)
}

func writeFile(t *testing.T, path, content string) {
	f, err := os.Create(path)
	if err != nil {
		assert.FailNow(t, "can't open file: %s", path)
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		assert.FailNow(t, "can't write to file: %s", path)
	}
}

func TestReadTerrafile(t *testing.T) {
	modules, err := readTerrafile("examples/Terrafile")
	if err != nil {
		assert.FailNow(t, "Unable to read Terrafile: %v", err)
	}
	mod1 := module{
		Name:      "aws-vpc",
		Source:    "git@github.com:terraform-aws-modules/terraform-aws-vpc",
		Version:   "v1.46.0",
		Subfolder: "",
	}
	mod2 := module{
		Name:      "simple-vpc",
		Source:    "git@github.com:terraform-aws-modules/terraform-aws-vpc",
		Version:   "v2.52.0",
		Subfolder: "examples/simple-vpc",
	}
	mod3 := module{
		Name:      "aws-vpc-master",
		Source:    "git@github.com:terraform-aws-modules/terraform-aws-vpc",
		Version:   "",
		Subfolder: "",
	}
	expected := []module{mod1, mod2, mod3}
	assert.ElementsMatch(t, expected, modules, "read modules don't match expected")
}

func TestDetect(t *testing.T) {
	type detectTest struct {
		source   string
		expected get
	}
	detectTests := []detectTest{
		{"/tmp/whatever", new(getFile)},
		{"whatever/folder", new(getFile)},
		{"./whatever/folder", new(getFile)},
		{"../whatever/folder", new(getFile)},
		{"git@github.com:kaminskip88/terrafile.git", new(getGit)},
		{"https://github.com/kaminskip88/terrafile.git", new(getGit)},
		{"git@bitbucket.org:org/terraform-modules.git", new(getGit)},
		{"https://user@bitbucket.org/org/terraform-modules.git", new(getGit)},
	}
	for _, i := range detectTests {
		g := detect(i.source)
		assert.Equal(t, i.expected, g, "detect() results don't match expected for source: %s", i.source)
	}
}
