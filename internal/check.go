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

	if strings.HasSuffix(name, l.config.HandlerNaming.Suffix) {
		matched, _ := regexp.MatchString(l.config.HandlerNaming.Pattern, name)
		if !matched {
			l.addError(LintError{
				File:       path,
				Line:       pos.Line,
				Column:     pos.Column,
				Type:       "handler_naming",
				Message:    fmt.Sprintf("Handler '%s' should match: %s", name, l.config.HandlerNaming.Description),
				Suggestion: "Use PascalCase ending with 'Handler'",
				Severity:   "error",
			})
		}
	} else if ast.IsExported(name) {
		matched, _ := regexp.MatchString(l.config.FunctionNaming.Pattern, name)
		if !matched && !contains(l.config.FunctionNaming.Exceptions, name) {
			l.addError(LintError{
				File:       path,
				Line:       pos.Line,
				Column:     pos.Column,
				Type:       "function_naming",
				Message:    fmt.Sprintf("Function '%s' should match: %s", name, l.config.FunctionNaming.Description),
				Suggestion: "Use PascalCase for exported functions",
				Severity:   "error",
			})
		}
	}
}

func (l *Linter) checkType(spec *ast.TypeSpec, path string) {
	name := spec.Name.Name
	pos := l.fileSet.Position(spec.Pos())

	switch spec.Type.(type) {
	case *ast.InterfaceType:
		matched, _ := regexp.MatchString(l.config.InterfaceNaming.Pattern, name)
		if !matched {
			l.addError(LintError{
				File:       path,
				Line:       pos.Line,
				Column:     pos.Column,
				Type:       "interface_naming",
				Message:    fmt.Sprintf("Interface '%s' should match: %s", name, l.config.InterfaceNaming.Description),
				Suggestion: "End with 'er'",
				Severity:   "error",
			})
		}
	case *ast.StructType:
		matched, _ := regexp.MatchString(l.config.StructNaming.Pattern, name)
		if !matched {
			l.addError(LintError{
				File:       path,
				Line:       pos.Line,
				Column:     pos.Column,
				Type:       "struct_naming",
				Message:    fmt.Sprintf("Struct '%s' should match: %s", name, l.config.StructNaming.Description),
				Suggestion: "Use PascalCase",
				Severity:   "error",
			})
		}
	}
}

func (l *Linter) checkVariable(name string, pos token.Pos, path string) {
	if contains(l.config.VariableNaming.Exceptions, name) {
		return
	}
	if !ast.IsExported(name) && (len(name) == 1 || name == "id" || name == "db" || name == "ok" || name == "err") {
		return
	}
	matched, _ := regexp.MatchString(l.config.VariableNaming.Pattern, name)
	if !matched {
		p := l.fileSet.Position(pos)
		l.addError(LintError{
			File:       path,
			Line:       p.Line,
			Column:     p.Column,
			Type:       "variable_naming",
			Message:    fmt.Sprintf("Variable '%s' should match: %s", name, l.config.VariableNaming.Description),
			Suggestion: "Use camelCase",
			Severity:   "warning",
		})
	}
}

func (l *Linter) checkConstant(name string, pos token.Pos, path string) {
	if ast.IsExported(name) {
		matched, _ := regexp.MatchString(l.config.ConstantNaming.Pattern, name)
		if !matched {
			p := l.fileSet.Position(pos)
			l.addError(LintError{
				File:       path,
				Line:       p.Line,
				Column:     p.Column,
				Type:       "constant_naming",
				Message:    fmt.Sprintf("Constant '%s' should match: %s", name, l.config.ConstantNaming.Description),
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
