package main

import (
	"fmt"
	"os"

	"github.com/yaza-putu/golang-linter/internal"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-linter [command] [args]")
		fmt.Println("Commands:")
		fmt.Println("  lint <project_path>     - Lint the project")
		fmt.Println("  init                    - Generate default config file")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "init":
		err := internal.GenerateDefaultConfig(".go.linter.json")
		if err != nil {
			fmt.Printf("Error generating config: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("Generated .go.linter.json configuration file")

		if err := internal.InitHook("."); err != nil {
			fmt.Fprintf(os.Stderr, "❌ %v\n", err)
			os.Exit(1)
		}

	case "lint":
		if len(os.Args) < 3 {
			fmt.Println("Usage: go-linter lint <project_path>")
			os.Exit(1)
		}

		projectPath := os.Args[2]
		l, err := internal.NewLinter(".go.linter.json")
		if err != nil {
			fmt.Printf("Error creating linter: %v\n", err)
			os.Exit(1)
		}

		err = l.LintProject(projectPath)
		if err != nil {
			fmt.Printf("Error linting project: %v\n", err)
			os.Exit(1)
		}

		if l.HasErrors() {
			l.PrintErrors()
			os.Exit(1)
		} else {
			fmt.Println("✅ All checks passed!")
		}
	default:
		fmt.Printf("Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
