# Building JAWT

### Prerequisites
- Go 1.24 or higher
- Java Runtime Environment (JRE) 8 or higher (for ANTLR4)

### Setting Up Development Environment

1. Install Go dependencies:
``` bash
   go mod download
```

2. Make sure the build scripts are executable:
``` bash
   chmod +x build.sh
   chmod +x internal/compiler/parser/generate.sh
```
### Architecture Overview

JAWT uses a JML compiler (internal/compiler):
- Processes JML page, files that define application entry points
- Processes JML component files that define reusable UI elements

The compiler has its ANTLR4 grammar and parser:
- ANTLR4 Grammar: `grammar/JML.g4`

### Building JAWT

#### Using the Build Script
The simplest way to build the project is using the provided build script:
``` bash
./build.sh
```
This will:
- Generate both parsers from their respective grammar files
- Build the project and place the binary in `build/bin/jawt`

#### Using Make
Alternatively, you can use the Makefile:
``` bash
# Clean, generate parser, and build
make

# Generate parser only
make generate

# Build only (assumes parser are already generated)
make build

# Clean build artifacts
make clean
```
#### Manual Build Process

If you prefer to build manually:
1. Generate the parsers for both compilers:
``` bash
   # Generate Compiler parser
   cd internal/compiler/parser
   ./generate.sh
```
1. Build the main application:
``` bash
   mkdir -p build/bin
   go build -o build/bin/jawt main.go
```
### Grammar Development

#### JML Grammar
The JML grammar defines the syntax for JAWT pages, components:
1. View the grammar:
``` bash
   cat grammar/JML.g4
```
1. After modifying, regenerate the parser:
``` bash
   cd internal/compiler/parser
   ./generate.sh
```
### Parser Generation Details
The parser generation process uses ANTLR4 to transform grammar files into Go code:

``` bash
   java -jar tools/antlr-4.13.2-complete.jar -Dlanguage=Go -package parser -visitor -o internal/compiler/parser grammar/JML.g4
```

### Creating New Grammar Rules
When extending the JML language:
1. Add new rules to the grammar file
2. Regenerate the parser
3. Implement visitor methods to handle the new rules
4. Add tests for the new functionality

### Running Tests
Run the test suite:
``` bash
go test ./...
```
Run tests with a coverage report:
``` bash
go test -cover ./...
```
### Creating a Release
To create a release build:
``` bash
./build.sh
```
The compiled binary will be available at `build/bin/jawt`.

## Troubleshooting

### Parser Generation Issues

If you encounter errors during parser generation:
1. Check that the grammar is being generated
2. Verify grammar syntax:
``` bash
   java -jar tools/antlr-4.13.2-complete.jar -o /tmp grammar/JML.g4
```
1. Ensure ANTLR4 JAR is accessible:
``` bash
   ls -la tools/antlr-4.13.2-complete.jar
```
### Build Errors
For Go build errors:
1. Update dependencies:
``` bash
   go mod tidy
```
1. Clean and rebuild:
``` bash
   make clean
   make
```


