#!/bin/bash

# Go Code Style Pre-commit Hook
# This script runs the Go linter before each commit

set -e

echo "🔍 Running Go code style checks..."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if go-linter binary exists
if ! command -v golinter &> /dev/null; then
    echo -e "${YELLOW}⚠️  golinter not found...${NC}"
    echo "please run command : go install github.com/yaza-putu/golinter@latest"
fi

# Check if config file exists
if [ ! -f ".go.linter.json" ]; then
    echo -e "${YELLOW}⚠️  .go.linter.json not found. Generating default config...${NC}"
    ./go-linter init
fi

# Get staged Go files
STAGED_GO_FILES=$(git diff --cached --name-only --diff-filter=ACM | grep '\.go$' || true)

if [ -z "$STAGED_GO_FILES" ]; then
    echo -e "${GREEN}✅ No Go files to check${NC}"
    exit 0
fi

echo -e "${YELLOW}📂 Staged Go files:${NC}"
echo "$STAGED_GO_FILES" | while read file; do
    echo "  - $file"
done
echo

# Run the linter on the project
echo -e "${YELLOW}🔍 Running linter...${NC}"
if ./go-linter lint .; then
    echo -e "${GREEN}✅ All code style checks passed!${NC}"
    exit 0
else
    echo -e "${RED}❌ Code style checks failed!${NC}"
    echo -e "${YELLOW}💡 Please fix the issues above before committing.${NC}"
    echo -e "${YELLOW}💡 You can customize the rules in .go.linter.json${NC}"
    exit 1
fi