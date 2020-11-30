package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type state struct {
	filePath string
	state    map[string]string
}

func newState(path string) state {
	return state{path, make(map[string]string)}
}

func (s state) set(k, v string) {
	s.state[k] = strings.TrimSpace(v)
}

func (s state) get(k string) string {
	return s.state[k]
}

func (s state) load() error {
	if _, err := os.Stat(s.filePath); err != nil {
		s.state = make(map[string]string)
		return nil
	}
	jsonFile, err := os.Open(s.filePath)
	if err != nil {
		return fmt.Errorf("Error opening state file: %v", err)
	}
	defer jsonFile.Close()

	jsonByte, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return fmt.Errorf("Error reading state file: %v", err)
	}
	json.Unmarshal(jsonByte, &s.state)
	return nil
}

func (s state) save() error {
	jsonState, err := json.Marshal(s.state)
	if err != nil {
		return fmt.Errorf("Error writing state: %v", err)
	}

	err = ioutil.WriteFile(s.filePath, jsonState, 0644)
	if err != nil {
		return fmt.Errorf("Error writing state file: %v", err)
	}
	return nil
}
