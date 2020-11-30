package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	log "github.com/sirupsen/logrus"
	yaml "gopkg.in/yaml.v2"
)

const (
	stateFile = "terrafile.state"
	usage     = `usage: %s
Get modules declared in Terrafile

Options:
`
)

type get interface {
	Get(module, string) (string, error)
	GetState(module) (string, error)
}

type module struct {
	Name      string
	Source    string `yaml:"source"`
	Version   string `yaml:"version"`
	Subfolder string `yaml:"subfolder"`
}

func main() {
	// parse flags
	Terrafile := flag.String("f", "./Terrafile", "Path to Terrafile")
	ModuleDir := flag.String("m", "./.modules", "Module folder")
	Verbose := flag.Bool("v", false, "Verbose output")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, usage, os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	// set log level
	log.SetLevel(log.InfoLevel)
	if *Verbose {
		log.SetLevel(log.DebugLevel)
	}

	// read Terrafile
	modules, err := readTerrafile(*Terrafile)
	if err != nil {
		log.Fatalln(err)
	}

	// load state file
	statePath := filepath.Join(*ModuleDir, stateFile)
	s := newState(statePath)
	err = s.load()
	if err != nil {
		log.Fatalln(err)
	}

	// get modules
	var wg sync.WaitGroup
	for _, m := range modules {
		wg.Add(1)
		go func(m module) {
			defer wg.Done()
			g, err := detect(m.Source)
			if err != nil {
				log.Fatalln(err)
			}
			modDir := filepath.Join(*ModuleDir, m.Name)
			sn := s.get(m.Name)
			if sn == "" {
				log.Infof("Getting module %s", m.Name)
				ref, err := g.Get(m, modDir)
				if err != nil {
					log.Fatalln(err)
				}
				s.set(m.Name, ref)
			} else {
				// log.Debugf("Found state for module %s", m.Name)
				ref, err := g.GetState(m)
				if err != nil {
					log.Fatalln(err)
				}
				if ref != sn {
					// log.Debugf("Update required for for module %s", m.Name)
					_, err = g.Get(m, modDir)
					if err != nil {
						log.Fatalln(err)
					}
					s.set(m.Name, ref)
				} else {
					log.Infof("Module %s up-to-date", m.Name)
				}
			}
		}(m)
	}
	wg.Wait()

	// save state
	err = s.save()
	if err != nil {
		log.Fatalln(err)
	}
}

func readTerrafile(path string) ([]module, error) {
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var mod map[string]module
	if err := yaml.Unmarshal(yamlFile, &mod); err != nil {
		return nil, err
	}
	var ms []module
	for n, m := range mod {
		m.Name = n
		ms = append(ms, m)
	}
	return ms, nil
}

func detect(source string) (get, error) {
	if filepath.IsAbs(source) {
		return new(getFile), nil
	}
	return new(getGit), nil
}
