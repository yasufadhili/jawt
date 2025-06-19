#!/bin/bash
set -e

# Directory where this script is located
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$DIR/../../../" && pwd )"

# Set up paths
GRAMMAR_DIR="$DIR/../grammar"
GRAMMAR_FILE="$GRAMMAR_DIR/JMLComponent.g4"
OUTPUT_DIR="$DIR/generated"
ANTLR_JAR="$PROJECT_ROOT/tools/antlr-4.13.2-complete.jar"

# Check if directories exist, create if not
mkdir -p "$GRAMMAR_DIR"
mkdir -p "$OUTPUT_DIR"

# Print colored messages
print_step() {
  echo -e "\033[1;34m== $1 ==\033[0m"
}

print_error() {
  echo -e "\033[1;31m✗ $1\033[0m" >&2
}

# Check if ANTLR JAR exists
if [ ! -f "$ANTLR_JAR" ]; then
    print_error "ANTLR JAR file not found at $ANTLR_JAR"
    exit 1
fi

# Check if the grammar file exists
if [ ! -f "$GRAMMAR_FILE" ]; then
    print_error "Grammar file not found at $GRAMMAR_FILE"
    exit 1
fi

# Run ANTLR4 to generate Go code
print_step "Generating Component Compiler parser from grammar: $GRAMMAR_FILE"
java -jar "$ANTLR_JAR" -Dlanguage=Go -package parser -visitor -o "$OUTPUT_DIR" "$GRAMMAR_FILE"

echo "Component Compiler parser generated successfully in $OUTPUT_DIR"