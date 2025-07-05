package internal

import (
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := Default()
	assert.NotNil(t, cfg, "Default config should not be nil")

	// Test key existence
	expectedRules := []string{
		"folder_naming",
		"file_naming",
		"handler_naming",
		"variable_naming",
		"function_naming",
		"constant_naming",
		"struct_naming",
		"interface_naming",
	}

	for _, rule := range expectedRules {
		_, ok := cfg.Rules[rule]
		assert.True(t, ok, "Rule '%s' should exist in default config", rule)
	}

	// Sample test for a specific rule (e.g. variable_naming)
	varRule := cfg.Rules["variable_naming"]
	assert.Equal(t, "^[a-z][a-zA-Z0-9]*$", varRule.Pattern)
	assert.Contains(t, varRule.Exceptions, "id")
	assert.Contains(t, varRule.Exceptions, "_")
	assert.Equal(t, "Variables should use camelCase", varRule.Description)

	// Check interface naming suffix
	ifaceRule := cfg.Rules["interface_naming"]
	assert.Equal(t, "er", ifaceRule.Suffix)
	assert.Regexp(t, ifaceRule.Pattern, "Reader")

	// Check exclusions
	assert.Contains(t, cfg.Exclusions.Folders, "vendor")
	assert.Contains(t, cfg.Exclusions.Files, "*.pb.go")
}

func TestDefaultPatternsAreValidAndMatchCorrectly(t *testing.T) {
	cfg := Default()
	testCases := []struct {
		ruleName     string
		validNames   []string
		invalidNames []string
	}{
		{
			ruleName:     "folder_naming",
			validNames:   []string{"my_folder", "test123", "internal"},
			invalidNames: []string{"MyFolder", "123folder", "my-folder"},
		},
		{
			ruleName:     "file_naming",
			validNames:   []string{"main.go", "file_name.go"},
			invalidNames: []string{"FileName.go", "fileName.go", "main.GO"},
		},
		{
			ruleName:     "handler_naming",
			validNames:   []string{"UserHandler", "OrderHandler"},
			invalidNames: []string{"userhandler", "Userhandler", "HandlerUser"},
		},
		{
			ruleName:     "variable_naming",
			validNames:   []string{"userName", "id123", "ok"},
			invalidNames: []string{"UserName", "123id", "_name"},
		},
		{
			ruleName:     "function_naming",
			validNames:   []string{"InitApp", "doSomething", "main", "init"},
			invalidNames: []string{"123func", "_helper"},
		},
		{
			ruleName:     "constant_naming",
			validNames:   []string{"MAX_LIMIT", "API_KEY"},
			invalidNames: []string{"MaxLimit", "api_key"},
		},
		{
			ruleName:     "struct_naming",
			validNames:   []string{"User", "Order123", "privateStruct"},
			invalidNames: []string{"_User", "private_struct"},
		},
		{
			ruleName:     "interface_naming",
			validNames:   []string{"Reader", "Writer", "Closer"},
			invalidNames: []string{"reader", "myInterface", "user_handler"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.ruleName, func(t *testing.T) {
			rule := cfg.Rules[tc.ruleName]
			re, err := regexp.Compile(rule.Pattern)
			assert.NoError(t, err, "pattern '%s' must be a valid regex", rule.Pattern)

			for _, name := range tc.validNames {
				assert.True(t, re.MatchString(name), "Expected name '%s' to match pattern for rule %s", name, tc.ruleName)
			}

			for _, name := range tc.invalidNames {
				assert.False(t, re.MatchString(name), "Expected name '%s' to NOT match pattern for rule %s", name, tc.ruleName)
			}
		})
	}
}
