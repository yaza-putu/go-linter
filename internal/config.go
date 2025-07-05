package internal

import (
	"encoding/json"
	"os"
)

// Config represents the linter configuration
type Config struct {
	FolderNaming struct {
		Pattern     string   `json:"pattern"`
		Description string   `json:"description"`
		Exceptions  []string `json:"exceptions"`
	} `json:"folder_naming"`

	FileNaming struct {
		Pattern     string   `json:"pattern"`
		Description string   `json:"description"`
		Exceptions  []string `json:"exceptions"`
	} `json:"file_naming"`

	HandlerNaming struct {
		Pattern     string `json:"pattern"`
		Description string `json:"description"`
		Suffix      string `json:"suffix"`
	} `json:"handler_naming"`

	VariableNaming struct {
		Pattern     string   `json:"pattern"`
		Description string   `json:"description"`
		Exceptions  []string `json:"exceptions"`
	} `json:"variable_naming"`

	FunctionNaming struct {
		Pattern     string   `json:"pattern"`
		Description string   `json:"description"`
		Exceptions  []string `json:"exceptions"`
	} `json:"function_naming"`

	ConstantNaming struct {
		Pattern     string `json:"pattern"`
		Description string `json:"description"`
	} `json:"constant_naming"`

	StructNaming struct {
		Pattern     string `json:"pattern"`
		Description string `json:"description"`
	} `json:"struct_naming"`

	InterfaceNaming struct {
		Pattern     string `json:"pattern"`
		Description string `json:"description"`
		Suffix      string `json:"suffix"`
	} `json:"interface_naming"`

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
		FolderNaming: struct {
			Pattern     string   `json:"pattern"`
			Description string   `json:"description"`
			Exceptions  []string `json:"exceptions"`
		}{
			Pattern:     "^[a-z][a-z0-9_]*$",
			Description: "Folder names should be lowercase with underscores",
			Exceptions:  []string{".git", ".github", "node_modules"},
		},
		FileNaming: struct {
			Pattern     string   `json:"pattern"`
			Description string   `json:"description"`
			Exceptions  []string `json:"exceptions"`
		}{
			Pattern:     "^[a-z][a-z0-9_]*\\.go$",
			Description: "Go files should be lowercase with underscores",
			Exceptions:  []string{"main.go", "go.mod", "go.sum"},
		},
		HandlerNaming: struct {
			Pattern     string `json:"pattern"`
			Description string `json:"description"`
			Suffix      string `json:"suffix"`
		}{
			Pattern:     "^[A-Z][a-zA-Z0-9]*Handler$",
			Description: "Handler functions should be PascalCase and end with 'Handler'",
			Suffix:      "Handler",
		},
		VariableNaming: struct {
			Pattern     string   `json:"pattern"`
			Description string   `json:"description"`
			Exceptions  []string `json:"exceptions"`
		}{
			Pattern:     "^[a-z][a-zA-Z0-9]*$",
			Description: "Variables should be camelCase",
			Exceptions:  []string{"i", "j", "k", "id", "db", "ok", "err"},
		},
		FunctionNaming: struct {
			Pattern     string   `json:"pattern"`
			Description string   `json:"description"`
			Exceptions  []string `json:"exceptions"`
		}{
			Pattern:     "^[A-Z][a-zA-Z0-9]*$",
			Description: "Public functions should be PascalCase",
			Exceptions:  []string{"main", "init"},
		},
		ConstantNaming: struct {
			Pattern     string `json:"pattern"`
			Description string `json:"description"`
		}{
			Pattern:     "^[A-Z][A-Z0-9_]*$",
			Description: "Constants should be UPPERCASE with underscores",
		},
		StructNaming: struct {
			Pattern     string `json:"pattern"`
			Description string `json:"description"`
		}{
			Pattern:     "^[A-Z][a-zA-Z0-9]*$",
			Description: "Structs should be PascalCase",
		},
		InterfaceNaming: struct {
			Pattern     string `json:"pattern"`
			Description string `json:"description"`
			Suffix      string `json:"suffix"`
		}{
			Pattern:     "^[A-Z][a-zA-Z0-9]*er$",
			Description: "Interfaces should be PascalCase and end with 'er'",
			Suffix:      "er",
		},
		Exclusions: struct {
			Folders []string `json:"folders"`
			Files   []string `json:"files"`
		}{
			Folders: []string{"vendor", ".git", "node_modules", "dist", "build"},
			Files:   []string{"*.pb.go", "*.gen.go", "*_test.go"},
		},
	}
}
