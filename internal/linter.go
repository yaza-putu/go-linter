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
			return l.checkRule("folder_naming", path, info.Name(), 1, 1)
		}
		if l.isExcludedFile(info.Name()) {
			return nil
		}
		if strings.HasSuffix(info.Name(), ".go") {
			if err := l.checkRule("file_naming", path, info.Name(), 1, 1); err != nil {
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

func (l *Linter) checkRule(ruleName, path, name string, line, col int) error {
	rule, ok := l.config.Rules[ruleName]
	if !ok {
		return nil
	}
	if isException(name, rule.Exceptions) {
		return nil
	}
	if !matchRegex(rule.Pattern, name) {
		l.addError(LintError{
			File:       path,
			Line:       line,
			Column:     col,
			Type:       ruleName,
			Message:    fmt.Sprintf("'%s' invalid: %s", name, rule.Description),
			Suggestion: defaultSuggestion(ruleName, rule.Suffix),
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
		switch node := n.(type) {
		case *ast.GenDecl:
			for _, spec := range node.Specs {
				if vs, ok := spec.(*ast.ValueSpec); ok {
					for _, name := range vs.Names {
						pos := l.fileSet.Position(name.Pos())
						ruleName := "variable_naming"
						if node.Tok == token.CONST {
							ruleName = "constant_naming"
						}
						l.checkDynamicRule(ruleName, name.Name, pos)
					}
				}
			}
		case *ast.FuncDecl:
			pos := l.fileSet.Position(node.Name.Pos())
			name := node.Name.Name
			if rule, ok := l.config.Rules["handler_naming"]; ok && rule.Suffix != "" && strings.HasSuffix(name, rule.Suffix) {
				l.checkDynamicRule("handler_naming", name, pos)
			} else {
				l.checkDynamicRule("function_naming", name, pos)
			}
		case *ast.TypeSpec:
			pos := l.fileSet.Position(node.Name.Pos())
			switch node.Type.(type) {
			case *ast.StructType:
				l.checkDynamicRule("struct_naming", node.Name.Name, pos)
			case *ast.InterfaceType:
				l.checkDynamicRule("interface_naming", node.Name.Name, pos)
			}
		case *ast.AssignStmt:
			assign := node
			if assign.Tok != token.DEFINE {
				break
			}
			for _, expr := range assign.Lhs {
				if ident, ok := expr.(*ast.Ident); ok {
					pos := l.fileSet.Position(ident.Pos())
					l.checkDynamicRule("variable_naming", ident.Name, pos)
				}
			}
		}
		return true
	})
	return nil
}

func (l *Linter) checkDynamicRule(ruleName, name string, pos token.Position) {
	rule, ok := l.config.Rules[ruleName]
	if !ok {
		return
	}
	if isException(name, rule.Exceptions) {
		return
	}
	if !matchRegex(rule.Pattern, name) {
		suggestion := defaultSuggestion(ruleName, rule.Suffix)
		l.addError(LintError{
			File:       pos.Filename,
			Line:       pos.Line,
			Column:     pos.Column,
			Type:       ruleName,
			Message:    fmt.Sprintf("'%s' doesn't match: %s", name, rule.Description),
			Suggestion: suggestion,
			Severity:   "warning",
		})
	}
}

func isException(name string, exceptions []string) bool {
	for _, exc := range exceptions {
		if name == exc {
			return true
		}
	}
	return false
}

func matchRegex(pattern, name string) bool {
	re := regexp.MustCompile(pattern)
	return re.MatchString(name)
}

func defaultSuggestion(ruleName, suffix string) string {
	switch ruleName {
	case "folder_naming":
		return "Use lowercase_with_underscores"
	case "file_naming":
		return "Use lowercase_with_underscores.go"
	case "function_naming":
		return "Use PascalCase for exported functions"
	case "variable_naming":
		return "Use camelCase"
	case "constant_naming":
		return "Use UPPERCASE_WITH_UNDERSCORES"
	case "struct_naming":
		return "Use PascalCase"
	case "interface_naming":
		if suffix != "" {
			return "End with '" + suffix + "'"
		}
		return "Use PascalCase"
	case "handler_naming":
		if suffix != "" {
			return "End with '" + suffix + "'"
		}
		return "Use PascalCase + Handler"
	default:
		return "Follow naming convention"
	}
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
