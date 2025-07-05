package internal

import (
	"fmt"
	"go/ast"
	"go/token"
	"regexp"
	"strings"
)

func (l *Linter) checkFunction(fn *ast.FuncDecl, path string) {
	if fn.Name == nil {
		return
	}
	name := fn.Name.Name
	pos := l.fileSet.Position(fn.Pos())

	// Handler check
	if rule, ok := l.config.Rules["handler_naming"]; ok {
		if strings.HasSuffix(name, rule.Suffix) {
			matched, _ := regexp.MatchString(rule.Pattern, name)
			if !matched {
				l.addError(LintError{
					File:       path,
					Line:       pos.Line,
					Column:     pos.Column,
					Type:       "handler_naming",
					Message:    fmt.Sprintf("Handler '%s' should match: %s", name, rule.Description),
					Suggestion: "Use PascalCase ending with '" + rule.Suffix + "'",
					Severity:   "error",
				})
			}
			return
		}
	}

	// Function check
	if rule, ok := l.config.Rules["function_naming"]; ok {
		if ast.IsExported(name) && !contains(rule.Exceptions, name) {
			matched, _ := regexp.MatchString(rule.Pattern, name)
			if !matched {
				l.addError(LintError{
					File:       path,
					Line:       pos.Line,
					Column:     pos.Column,
					Type:       "function_naming",
					Message:    fmt.Sprintf("Function '%s' should match: %s", name, rule.Description),
					Suggestion: "Use PascalCase for exported functions",
					Severity:   "error",
				})
			}
		}
	}
}

func (l *Linter) checkType(spec *ast.TypeSpec, path string) {
	name := spec.Name.Name
	pos := l.fileSet.Position(spec.Pos())

	switch spec.Type.(type) {
	case *ast.InterfaceType:
		if rule, ok := l.config.Rules["interface_naming"]; ok {
			matched, _ := regexp.MatchString(rule.Pattern, name)
			if !matched {
				l.addError(LintError{
					File:       path,
					Line:       pos.Line,
					Column:     pos.Column,
					Type:       "interface_naming",
					Message:    fmt.Sprintf("Interface '%s' should match: %s", name, rule.Description),
					Suggestion: "End with '" + rule.Suffix + "'",
					Severity:   "error",
				})
			}
		}
	case *ast.StructType:
		if rule, ok := l.config.Rules["struct_naming"]; ok {
			matched, _ := regexp.MatchString(rule.Pattern, name)
			if !matched {
				l.addError(LintError{
					File:       path,
					Line:       pos.Line,
					Column:     pos.Column,
					Type:       "struct_naming",
					Message:    fmt.Sprintf("Struct '%s' should match: %s", name, rule.Description),
					Suggestion: "Use PascalCase",
					Severity:   "error",
				})
			}
		}
	}
}

func (l *Linter) checkVariable(name string, pos token.Pos, path string) {
	rule, ok := l.config.Rules["variable_naming"]
	if !ok {
		return
	}

	if contains(rule.Exceptions, name) {
		return
	}

	if !ast.IsExported(name) && (len(name) == 1 || name == "id" || name == "db" || name == "ok" || name == "err") {
		return
	}

	matched, _ := regexp.MatchString(rule.Pattern, name)
	if !matched {
		p := l.fileSet.Position(pos)
		l.addError(LintError{
			File:       path,
			Line:       p.Line,
			Column:     p.Column,
			Type:       "variable_naming",
			Message:    fmt.Sprintf("Variable '%s' should match: %s", name, rule.Description),
			Suggestion: "Use camelCase",
			Severity:   "warning",
		})
	}
}

func (l *Linter) checkConstant(name string, pos token.Pos, path string) {
	rule, ok := l.config.Rules["constant_naming"]
	if !ok {
		return
	}

	if ast.IsExported(name) {
		matched, _ := regexp.MatchString(rule.Pattern, name)
		if !matched {
			p := l.fileSet.Position(pos)
			l.addError(LintError{
				File:       path,
				Line:       p.Line,
				Column:     p.Column,
				Type:       "constant_naming",
				Message:    fmt.Sprintf("Constant '%s' should match: %s", name, rule.Description),
				Suggestion: "Use UPPERCASE_WITH_UNDERSCORES",
				Severity:   "error",
			})
		}
	}
}

func contains(list []string, target string) bool {
	for _, v := range list {
		if v == target {
			return true
		}
	}
	return false
}
