# Go Linter Tool

A simple configurable Go code linter that validates code naming conventions for folders, files, functions, variables, constants, structs, and interfaces. It also supports pre-commit hook integration to prevent bad naming conventions from being committed.

---

## 📁 Folder Structure

```
/go-linter-tool
├── main.go # CLI entry point (init/lint)
├── internal/
│   ├── check.go               # AST naming rule checks (functions, vars, consts, types)
│   ├── config.go              # Default configuration and JSON loading
│   ├── init.go                # Pre-commit hook initializer
│   └── linter.go              # Core linter logic and file walking
├── .go.linter.json            # Configuration file (auto-generated)
├── pre-commit-hook.sh         # Git hook script
└── README.md                 # Project documentation
```

---

## 🔧 Features

- Validates naming rules:
  - Folder and file names
  - Function and handler names
  - Variable and constant names
  - Struct and interface names
- Uses JSON config: `.go.linter.json`
- Skips ignored folders/files
- Git pre-commit hook integration

---

## 🚀 Usage
### 1. Install
```bash
go install github.com/yaza-putu/go-linter@latest
```

if it's not accessible, try moving it to the global bin
```bash
 mv ~/go/bin/go-linter /usr/local/bin
```

### 2. Generate Default Config + Hook

```sh
$ golinter init
```

This will:

- Create `.go.linter.json` if missing
- Copy `pre-commit-hook.sh` into `.git/hooks/pre-commit`

### 3. Lint Your Project

```sh
git commit -m "message"
```
or
```sh
golinter lint .
```
---

## 📂 File Overview

### `main.go`

- Main CLI entry point
- Handles commands:
  - `init`: generates config + sets hook
  - `lint <path>`: runs linting logic

### `internal/config.go`

- Defines `Config` struct
- Loads config from `.go.linter.json`
- Supplies defaults if config not found

### `internal/init.go`

- Handles setup of Git pre-commit hook
- Makes the script executable on UNIX-based systems

### `internal/linter.go`

- Walks through project folder
- Applies folder & file naming checks
- Parses `.go` files into ASTs

### `internal/check.go`

- Applies AST-level rules
  - Naming for: functions, variables, constants, interfaces, structs

### `.go.linter.json`

- JSON config file for all lint rules
- Customizable naming regex, descriptions, and exceptions

### `pre-commit-hook.sh`

- Executed automatically by Git before commit
- Runs `go run main.go lint .`
- Blocks commit if lint errors exist

---

## ✅ Example Output

```
[ERROR] my_handler.go:5:1 - Handler function 'my_handler' doesn't match pattern: Handler functions should be PascalCase and end with 'Handler'
  Suggestion: Use PascalCase and end with 'Handler'
```

---

## 🧪 Custom Rules

Edit `.go.linter.json` to change:

- Regex patterns
- Description messages
- Naming exceptions

---

## 🛡️ Pre-commit Protection

Once `init` is run, every `git commit` will trigger the linter. If violations exist, the commit will be blocked.

---

## 📜 License

MIT

