package main

import (
	"flag"
	"fmt"
	"go/parser"
	"os"
	"runtime"

	"github.com/codeclimate/cc-engine-go/engine"
	"golang.org/x/tools/go/loader"
	"honnef.co/go/tools/lint/lintutil"
	"honnef.co/go/tools/staticcheck"
)

var (
	Version = "latest"

	flagVersion bool
)

func init() {
	flag.BoolVar(&flagVersion, "v", false, "print version and exit")
	flag.Parse()
}

func main() {

	if flagVersion {
		fmt.Printf("codeclimate-staticcheck %s (%s_%s)\n", Version, runtime.GOOS, runtime.GOARCH)
		os.Exit(0)
	}

	config, err := engine.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s\n", err)
		os.Exit(1)
	}

	// path, _ := os.Getwd()
	path := "/code/"

	filenames, err := engine.GoFileWalk(path, engine.IncludePaths(path, config))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading files: %s\n", err)
		os.Exit(1)
	}

	loader := &loader.Config{ParserMode: parser.ImportsOnly, Cwd: path}
	loader.CreateFromFilenames(path, filenames...)
	program, err := loader.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading program: %s\n", err)
		os.Exit(1)
	}

	var packages []string
	for _, packageInfo := range program.InitialPackages() {
		packages = append(packages, packageInfo.Pkg.Name())
	}

	checker := staticcheck.NewChecker()
	problems, _, err := lintutil.Lint(checker, packages, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error linting program: %s\n", err)
		os.Exit(1)
	}

	for _, problem := range problems {
		position := program.Fset.Position(problem.Position)
		issue := &engine.Issue{
			Type:              "issue",
			Check:             "Staticcheck" + problem.Text,
			Description:       problem.Text,
			RemediationPoints: 5000,
			Location: &engine.Location{
				Lines: &engine.LinesOnlyPosition{
					Begin: position.Line,
					End:   position.Line,
				},
			},
			Categories: []string{"Style"},
		}
		engine.PrintIssue(issue)
	}
}
