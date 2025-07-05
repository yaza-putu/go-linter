package internal

import (
	"encoding/json"
	"os"
)

// RuleConfig defines a single rule's pattern, description, and optional fields
type RuleConfig struct {
	Pattern     string   `json:"pattern"`
	Description string   `json:"description"`
	Exceptions  []string `json:"exceptions,omitempty"`
	Suffix      string   `json:"suffix,omitempty"`
}

// Config represents the linter configuration
type Config struct {
	Rules      map[string]RuleConfig `json:"rules"`
	Exclusions struct {
		Folders []string `json:"folders"`
		Files   []string `json:"files"`
	} `json:"exclusions"`
}

// LoadConfig loads configuration from a JSON file
func LoadConfig(configPath string) (*Config, error) {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return Default(), nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}

// GenerateDefaultConfig generates a default configuration file
func GenerateDefaultConfig(path string) error {
	config := Default()
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		Rules: map[string]RuleConfig{
			"folder_naming": {
				Pattern:     "^[a-z][a-z0-9_]*$",
				Description: "Folder names should be lowercase with underscores",
				Exceptions:  []string{".git", ".github", "node_modules", "internal", "."},
			},
			"file_naming": {
				Pattern:     "^[a-z][a-z0-9_]*\\.go$",
				Description: "File names should be lowercase with underscores",
				Exceptions:  []string{"main.go", "go.mod", "go.sum"},
			},
			"handler_naming": {
				Pattern:     "^[A-Za-z][a-zA-Z0-9]*Handler$",
				Description: "Handlers should be PascalCase (exported) or camelCase (unexported) and end with 'Handler'",
				Suffix:      "Handler",
			},
			"variable_naming": {
				Pattern:     "^[A-Za-z][a-zA-Z0-9]*$",
				Description: "Variables should be PascalCase (exported) or camelCase (unexported)",
				Exceptions:  []string{"i", "j", "k", "id", "db", "ok", "err", "_"},
			},
			"function_naming": {
				// Allow PascalCase (exported) and camelCase (unexported)
				Pattern:     "^[A-Za-z][a-zA-Z0-9]*$",
				Description: "Functions should be PascalCase (exported) or camelCase (unexported)",
				Exceptions:  []string{"main", "init"},
			},
			"constant_naming": {
				Pattern:     "^[A-Z][A-Z0-9_]*$",
				Description: "Constants should be UPPER_CASE with underscores",
			},
			"struct_naming": {
				Pattern:     "^[A-Za-z][a-zA-Z0-9]*$",
				Description: "Struct name should be PascalCase (exported) or camelCase (unexported)",
			},
			"interface_naming": {
				Pattern:     "^[A-Z][a-zA-Z0-9]*$",
				Description: "Interfaces should be PascalCase",
				Suffix:      "er",
			},
		},
		Exclusions: struct {
			Folders []string `json:"folders"`
			Files   []string `json:"files"`
		}{
			Folders: []string{
				"vendor", ".git", "node_modules", "dist", "build", "third_party", "docs",
			},
			Files: []string{
				"*.pb.go", "*.gen.go", "*_test.go", "*_mock.go", "*_fixture.go",
			},
		},
	}
}
