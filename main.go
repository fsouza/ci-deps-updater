// Copyright 2019 Francisco Souza. All rights reserved.
// Use of this source code is governed by an ISC-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
)

// TODO(fsouza): automatically detect most recent releases of Go.
var goVersions = []string{
	"1.11.x",
	"1.12.x",
	"1.13rc",
}

func main() {
	repo := flag.String("repo", "", "name of the repo to process")
	flag.Parse()
	dir, err := cloneRepo(*repo)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(dir)
}

func cloneRepo(repoURL string) (dir string, err error) {
	if repoURL == "" {
		return "", errors.New("repo is required")
	}
	tmpdir, err := ioutil.TempDir("", "")
	if err != nil {
		return "", err
	}
	dir = filepath.Join(tmpdir, "code")
	var buf bytes.Buffer
	cmd := exec.Command("git", "clone", repoURL, dir)
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	err = cmd.Run()
	if err != nil {
		return "", cmdError("git clone", err, &buf)
	}
	return dir, nil
}

func cmdError(cmd string, err error, output *bytes.Buffer) error {
	return fmt.Errorf("cannot run %q: %v\nOUTPUT:\n%s", cmd, err, output)
}
