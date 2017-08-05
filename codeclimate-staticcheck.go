package main

import (
	"fmt"
	"go/parser"
	"os"

	"github.com/codeclimate/cc-engine-go/engine"
	"golang.org/x/tools/go/loader"
	"honnef.co/go/tools/lint"
	"honnef.co/go/tools/staticcheck"
)

func main() {
	config, err := engine.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading config: %s", err)
		os.Exit(1)
	}

	path := "/code/"

	loader := &loader.Config{ParserMode: parser.ImportsOnly, Cwd: path}
	loader.CreateFromFilenames(path, engine.IncludePaths(path, config)...)
	program, err := loader.Load()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	checker := staticcheck.NewChecker()
	linter := lint.Linter{Checker: checker}
	problems := linter.Lint(program)

	for _, problem := range problems {

		position := program.Fset.Position(problem.Position)

		engine.PrintIssue(&engine.Issue{
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
		})
	}
}
