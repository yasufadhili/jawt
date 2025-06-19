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
   chmod +x internal/pc/parser/generate.sh
   chmod +x internal/cc/parser/generate.sh
```
### Architecture Overview

JAWT uses two separate compilers:
1. **Page Compiler (PC)**: Processes JML page files that define application entry points
2. **Component Compiler (CC)**: Processes JML component files that define reusable UI elements

Each compiler has its own ANTLR4 grammar and parser:
- Page Grammar: `internal/pc/grammar/JMLPage.g4`
- Component Grammar: `internal/cc/grammar/JMLComponent.g4`

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
# Clean, generate parsers, and build
make

# Generate parsers only
make generate

# Build only (assumes parsers are already generated)
make build

# Clean build artifacts
make clean
```
#### Manual Build Process

If you prefer to build manually:
1. Generate the parsers for both compilers:
``` bash
   # Generate Page Compiler parser
   cd internal/pc/parser
   ./generate.sh
   
   # Generate Component Compiler parser
   cd ../../cc/parser
   ./generate.sh
```
1. Build the main application:
``` bash
   mkdir -p build/bin
   go build -o build/bin/jawt main.go
```
### Grammar Development
#### Page Grammar (PC)
The Page grammar defines the syntax for application pages:
1. View the grammar:
``` bash
   cat internal/pc/grammar/JMLPage.g4
```
1. After modifying, regenerate the parser:
``` bash
   cd internal/pc/parser
   ./generate.sh
```
#### Component Grammar (CC)
The Component grammar defines the syntax for reusable UI components:
1. View the grammar:
``` bash
   cat internal/cc/grammar/JMLComponent.g4
```
1. After modifying, regenerate the parser:
``` bash
   cd internal/cc/parser
   ./generate.sh
```
### Parser Generation Details
The parser generation process uses ANTLR4 to transform grammar files into Go code:
1. For the Page Compiler:
``` bash
   java -jar tools/antlr-4.13.2-complete.jar -Dlanguage=Go -package parser -visitor -o internal/pc/parser internal/pc/grammar/JMLPage.g4
```
1. For the Component Compiler:
``` bash
   java -jar tools/antlr-4.13.2-complete.jar -Dlanguage=Go -package parser -visitor -o internal/cc/parser internal/cc/grammar/JMLComponent.g4
```

### Creating New Grammar Rules
When extending the JML language:
1. Add new rules to the appropriate grammar file
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
1. Check which specific grammar is failing (Page or Component)
2. Verify grammar syntax:
``` bash
   # For Page grammar
   java -jar tools/antlr-4.13.2-complete.jar -o /tmp internal/pc/grammar/JMLPage.g4
   
   # For Component grammar
   java -jar tools/antlr-4.13.2-complete.jar -o /tmp internal/cc/grammar/JMLComponent.g4
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


