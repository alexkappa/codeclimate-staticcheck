package main

import (
	"fmt"
	"go/token"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"honnef.co/go/tools/lint/lintutil"
	"honnef.co/go/tools/staticcheck"

	"github.com/codeclimate/cc-engine-go/engine"
	"github.com/kisielk/gotool"
)

func main() {
	config, err := engine.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
		os.Exit(1)
	}

	rootPath := "/code"

	pkgPath := filepath.Join(os.Getenv("GOPATH"), "src", "app")
	if engineConfig, ok := config["config"].(map[string]interface{}); ok {
		if packageDir, ok := engineConfig["package"].(string); ok {
			pkgPath = filepath.Join(os.Getenv("GOPATH"), "src", packageDir)
		}
	}

	err = copyDir(pkgPath, rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed copying dir to %q. %s\n", pkgPath, err)
		os.Exit(1)
	}

	os.Chdir(pkgPath)

	pkgs := gotool.ImportPaths([]string{"./..."})

	checker := staticcheck.NewChecker()
	problems, program, err := lintutil.Lint(checker, pkgs, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error linting program: %s\n", err)
		os.Exit(1)
	}

	includePaths := engine.IncludePaths(pkgPath, config)

	for _, problem := range problems {

		position := program.Fset.Position(problem.Position)

		if inSlice(position.Filename, includePaths) {
			issue := &engine.Issue{
				Type:              "issue",
				Check:             check(problem.Text),
				Description:       description(problem.Text),
				RemediationPoints: 5000,
				Location:          location(position, pkgPath),
				Categories:        []string{"Style"},
			}

			engine.PrintIssue(issue)
		}
	}
}

func check(s string) string {
	return "Go-StaticCheck/" + strings.TrimRight(strings.Split(s, "(")[1], ")")
}

func description(s string) string {
	return strings.Split(s, "(")[0]
}

func location(pos token.Position, path string) *engine.Location {
	return &engine.Location{
		Path: strings.SplitAfter(pos.Filename, path+"/")[1],
		Lines: &engine.LinesOnlyPosition{
			Begin: pos.Line,
			End:   pos.Line,
		},
	}
}

func copyDir(dst, src string) error {
	walkFn := func(srcPath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		relPath, err := filepath.Rel(src, srcPath)
		if err != nil {
			return fmt.Errorf("failed getting relative path. %s", err)
		}
		dstPath := filepath.Join(dst, relPath)
		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode())
		}
		return copyFile(dstPath, srcPath)
	}
	return filepath.Walk(src, walkFn)
}

func copyFile(dst, src string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()
	tmp, err := ioutil.TempFile(filepath.Dir(dst), filepath.Base(dst)+".")
	if err != nil {
		return err
	}
	_, err = io.Copy(tmp, in)
	if err != nil {
		tmp.Close()
		os.Remove(tmp.Name())
		return err
	}
	if err := tmp.Close(); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	const perm = 0644
	if err := os.Chmod(tmp.Name(), perm); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	if err := os.Rename(tmp.Name(), dst); err != nil {
		os.Remove(tmp.Name())
		return err
	}
	return nil
}

func inSlice(s string, slice []string) bool {
	for _, elem := range slice {
		if elem == s {
			return true
		}
	}
	return false
}
