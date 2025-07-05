package internal

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type LintError struct {
	File       string `json:"file"`
	Line       int    `json:"line"`
	Column     int    `json:"column"`
	Type       string `json:"type"`
	Message    string `json:"message"`
	Suggestion string `json:"suggestion"`
	Severity   string `json:"severity"`
}

type Linter struct {
	config  *Config
	fileSet *token.FileSet
	errors  []LintError
}

func NewLinter(configPath string) (*Linter, error) {
	cfg, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("load config: %w", err)
	}

	return &Linter{
		config:  cfg,
		fileSet: token.NewFileSet(),
		errors:  []LintError{},
	}, nil
}

func (l *Linter) LintProject(root string) error {
	return filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			if l.isExcludedFolder(info.Name()) {
				return filepath.SkipDir
			}
			return l.checkFolder(path, info.Name())
		}
		if l.isExcludedFile(info.Name()) {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".go") {
			if err := l.checkFile(path, info.Name()); err != nil {
				return err
			}
			return l.parseAST(path)
		}
		return nil
	})
}

func (l *Linter) isExcludedFolder(name string) bool {
	for _, f := range l.config.Exclusions.Folders {
		if matched, _ := filepath.Match(f, name); matched {
			return true
		}
	}
	return false
}

func (l *Linter) isExcludedFile(name string) bool {
	for _, f := range l.config.Exclusions.Files {
		if matched, _ := filepath.Match(f, name); matched {
			return true
		}
	}
	return false
}

func (l *Linter) checkFolder(path, name string) error {
	for _, exc := range l.config.FolderNaming.Exceptions {
		if name == exc {
			return nil
		}
	}
	matched, _ := regexp.MatchString(l.config.FolderNaming.Pattern, name)
	if !matched {
		l.addError(LintError{
			File:       path,
			Line:       1,
			Column:     1,
			Type:       "folder_naming",
			Message:    fmt.Sprintf("Folder '%s' invalid: %s", name, l.config.FolderNaming.Description),
			Suggestion: "Use lowercase_with_underscores",
			Severity:   "warning",
		})
	}
	return nil
}

func (l *Linter) checkFile(path, name string) error {
	for _, exc := range l.config.FileNaming.Exceptions {
		if name == exc {
			return nil
		}
	}
	matched, _ := regexp.MatchString(l.config.FileNaming.Pattern, name)
	if !matched {
		l.addError(LintError{
			File:       path,
			Line:       1,
			Column:     1,
			Type:       "file_naming",
			Message:    fmt.Sprintf("File '%s' invalid: %s", name, l.config.FileNaming.Description),
			Suggestion: "Use lowercase_with_underscores.go",
			Severity:   "warning",
		})
	}
	return nil
}

func (l *Linter) parseAST(path string) error {
	src, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	file, err := parser.ParseFile(l.fileSet, path, src, parser.ParseComments)
	if err != nil {
		return err
	}

	ast.Inspect(file, func(n ast.Node) bool {
		if n == nil {
			return true
		}
		return true
	})
	return nil
}

func (l *Linter) addError(e LintError) {
	l.errors = append(l.errors, e)
}

func (l *Linter) GetErrors() []LintError {
	return l.errors
}

func (l *Linter) HasErrors() bool {
	return len(l.errors) > 0
}

func (l *Linter) PrintErrors() {
	for _, e := range l.errors {
		fmt.Printf("[%s] %s:%d:%d - %s\n", strings.ToUpper(e.Severity), e.File, e.Line, e.Column, e.Message)
		if e.Suggestion != "" {
			fmt.Printf("  Suggestion: %s\n", e.Suggestion)
		}
	}
}
