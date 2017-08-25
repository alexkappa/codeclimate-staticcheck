package main

import (
	"fmt"
	"go/token"
	"os"
	"path/filepath"
	"strings"

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

	rootPath := filepath.Join(os.Getenv("GOPATH"), "src", "app")
	if engineConfig, ok := config["config"].(map[string]interface{}); ok {
		if packageDir, ok := engineConfig["package_dir"].(string); ok {
			rootPath = filepath.Join(os.Getenv("GOPATH"), "src", packageDir)
		}
	}

	if _, err := os.Stat(rootPath); err != nil {
		rootPathDir := filepath.Dir(rootPath)
		err := os.MkdirAll(rootPathDir, 0744)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed creating package dir. %s\n", err)
			os.Exit(1)
		}
		err = os.Symlink("/code", rootPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed creating symlink from package dir to /code. %s\n", err)
			os.Exit(1)
		}
		os.Chdir(rootPath)
	}

	targetPath, err := filepath.EvalSymlinks(rootPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed evaluating $GOPATH symlink. %s\n", err)
		os.Exit(1)
	}

	filenames, err := engine.GoFileWalk(targetPath, engine.IncludePaths(targetPath, config))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading files: %s\n", err)
		os.Exit(1)
	}

	loader := &loader.Config{
		ParserMode:  0, // parser.ParseComments | DeclarationErrors,
		Cwd:         targetPath,
		AllowErrors: true,
	}
	loader.CreateFromFilenames(targetPath, filenames...)
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
