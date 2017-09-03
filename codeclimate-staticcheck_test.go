package main

import (
	"os"
	"reflect"
	"sort"
	"testing"
)

func TestGetPkgs(t *testing.T) {
	for _, test := range []struct {
		root  string
		paths []string
		pkgs  []string
	}{
		{
			"/go/src",
			[]string{
				"/go/src/app/foo.go",
				"/go/src/app/bar.go",
				"/go/src/app/bar/baz.go",
			},
			[]string{
				"app",
				"app/bar",
			},
		},
		{
			"/go/src",
			[]string{
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/.codeclimate.yml",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/.gitignore",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/codeclimate-staticcheck.go",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/codeclimate-staticcheck_test.go",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/Dockerfile",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/engine.json",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/LICENSE",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/Makefile",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/README.md",
				"/go/src/github.com/alexkappa/codeclimate-staticcheck/wercker.yml",
			},
			[]string{
				"github.com/alexkappa/codeclimate-staticcheck",
			},
		},
	} {
		goPath := os.Getenv("GOPATH")
		os.Setenv("GOPATH", "/go")
		defer os.Setenv("GOPATH", goPath)

		pkgs := getPkgs(test.paths)

		sort.Strings(pkgs)
		sort.Strings(test.pkgs)

		if !reflect.DeepEqual(test.pkgs, pkgs) {
			t.Errorf("Unexpected packages found.\n\tHave: %v\n\tWant: %v", test.pkgs, pkgs)
		}
	}
}
