
#!/bin/bash
set -e


# Project root directory (where this script is located)
PROJECT_ROOT="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"



TOOLS_DIR="$PROJECT_ROOT/tools"

INTERNAL_DIR="$PROJECT_ROOT/internal"

COMPILER_DIR="$INTERNAL_DIR/compiler"

BUILD_DIR="$PROJECT_ROOT/build"

BIN_DIR="$BUILD_DIR/bin"



ANTLR_JAR="$TOOLS_DIR/antlr-4.13.2-complete.jar"


# Print coloured messages
print_step() {
  echo -e "\033[1;34m=== $1 ===\033[0m"
}

print_success() {
  echo -e "\033[1;32m✓ $1\033[0m"
}

print_error() {
  echo -e "\033[1;31m✗ $1\033[0m" >&2
}



if [ ! -f "$ANTLR_JAR" ]; then
    print_error "ANTLR JAR file not found at $ANTLR_JAR"
    exit 1
fi



mkdir -p "$BIN_DIR"



print_step "Generating Compiler Parser"
if [ -f "$COMPILER_DIR/parser/generate.sh" ]; then
    (cd "$COMPILER_DIR/parser" && ./generate.sh)
    if [ $? -ne 0 ]; then
        print_error "Failed to generate Component Compiler parser"
        exit 1
    fi
    print_success "Component Compiler parser generated successfully"
else
    print_error "Component Compiler generate script not found at $COMPILER_DIR/parser/generate.sh"
    exit 1
fi



print_step "Building JAWT Binary"
go build -o "$BIN_DIR/jawt" "$PROJECT_ROOT/main.go"
if [ $? -ne 0 ]; then
    print_error "Build failed"
    exit 1
fi




print_success "Build completed successfully!"
print_success "Binary available at: $BIN_DIR/jawt"

