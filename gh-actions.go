// Copyright 2019 Francisco Souza. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type workflow struct {
	Name string
	On   yaml.MapSlice
	Jobs map[string]job
}

type job struct {
	Strategy *strategy `yaml:",omitempty"`
	Name     string
	RunsOn   string `yaml:"runs-on"`
	Steps    []step
}

type strategy struct {
	Matrix map[string]interface{}
}

func (s *strategy) hasGo() bool {
	// "go_version" is a convention. could parameterize.
	_, ok := s.Matrix["go_version"]
	return ok
}

type step struct {
	Uses string            `yaml:",omitempty"`
	ID   string            `yaml:",omitempty"`
	With map[string]string `yaml:",omitempty"`
	Run  string            `yaml:",omitempty"`
	Env  yaml.MapSlice     `yaml:",omitempty"`
}

func (s *step) isGo() bool {
	return strings.HasPrefix(s.Uses, "actions/setup-go") || strings.HasPrefix(s.Uses, "docker://golang")
}

func handleGHActions(dir string) (bool, error) {
	workflows, err := loadWorkflows(dir)
	if err != nil {
		return false, err
	}
	if len(workflows) < 1 {
		return false, errors.New("couldn't load any workflows")
	}
	for filename, wflow := range workflows {
	}
	return false, nil
}

func loadWorkflows(repoDir string) (map[string]workflow, error) {
	fullDirPath := filepath.Join(repoDir, ".github", "workflows")
	wdir, err := os.Open(fullDirPath)
	if err != nil {
		return nil, fmt.Errorf("cannot open dir %q: %w", fullDirPath, err)
	}
	defer wdir.Close()
	fis, err := wdir.Readdir(0)
	if err != nil {
		return nil, fmt.Errorf("cannot read dir %q: %w", fullDirPath, err)
	}
	workflows := map[string]workflow{}
	for _, fi := range fis {
		if ext := filepath.Ext(fi.Name()); ext != "yaml" && ext != "yml" {
			continue
		}
		filePath := filepath.Join(fullDirPath, fi.Name())
		w, err := loadWorkflow(filePath)
		if err != nil {
			log.Printf("[WARNING] invalid workflow in file %q: %v", filePath, err)
			continue
		}
		workflows[filePath] = w
	}
	return workflows, nil
}

func loadWorkflow(file string) (workflow, error) {
	var w workflow
	f, err := os.Open(file)
	if err != nil {
		return w, err
	}
	defer f.Close()
	decoder := yaml.NewDecoder(f)
	decoder.SetStrict(true)
	err = decoder.Decode(&w)
	return w, err
}
