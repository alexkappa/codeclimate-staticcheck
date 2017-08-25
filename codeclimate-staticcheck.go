package main

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"

	"github.com/alexkappa/codeclimate-staticcheck/fileutil"
	"github.com/codeclimate/cc-engine-go/engine"
	"golang.org/x/tools/go/loader"
	"honnef.co/go/tools/lint"
	"honnef.co/go/tools/staticcheck"
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
		if packageDir, ok := engineConfig["package_dir"].(string); ok {
			pkgPath = filepath.Join(os.Getenv("GOPATH"), "src", packageDir)
		}
	}

	err = fileutil.CopyDir(pkgPath, rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed copying dir to %q. %s\n", pkgPath, err)
		os.Exit(1)
	}

	filenames, err := fileutil.GoFileWalk(pkgPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading files: %s\n", err)
		os.Exit(1)
	}

	loader := &loader.Config{
		ParserMode:  0,
		Cwd:         pkgPath,
		AllowErrors: true,
	}
	loader.CreateFromFilenames(pkgPath, filenames...)
	program, err := loader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading program: %s\n", err)
		os.Exit(1)
	}

	checker := staticcheck.NewChecker()
	linter := lint.Linter{Checker: checker}
	problems := linter.Lint(program)

	for _, problem := range problems {

		issue := &engine.Issue{
			Type:              "issue",
			Check:             check(problem.Text),
			Description:       description(problem.Text),
			RemediationPoints: 5000,
			Location:          location(program, problem.Position, rootPath),
			Categories:        []string{"Style"},
		}

		engine.PrintIssue(issue)
	}
}

func check(s string) string {
	return "Staticcheck/" + strings.TrimRight(strings.Split(s, "(")[1], ")")
}

func description(s string) string {
	return strings.Split(s, "(")[0]
}

func location(p *loader.Program, pos token.Pos, path string) *engine.Location {
	position := p.Fset.Position(pos)
	return &engine.Location{
		Path: strings.SplitAfter(position.Filename, path)[1],
		Lines: &engine.LinesOnlyPosition{
			Begin: position.Line,
			End:   position.Line,
		},
	}
}
