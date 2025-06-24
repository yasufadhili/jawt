# Compilation System

The JAWT compilation system is responsible for transforming JML source files into web assets. It consists of several specialised compilers that work together to process different aspects of the application.

## Core Compilers

### 1. Page Compiler (`internal/page_compiler`)

Responsible for processing page definitions and generating HTML output.

**Key Responsibilities:**
- Parse JML page definitions
- Generate HTML structure
- Handle page metadata (title, description, etc.)
- Manage routing and navigation
- Integrate with the component system

**Input:** `.jml` files in the `app/` directory  
**Output:** HTML files in the build directory

### 2. Component Compiler (`internal/cc`)

Processes reusable UI components defined in JML.

**Key Responsibilities:**
- Transform JML components to JavaScript
- Handle component props and state
- Manage component lifecycle
- Support component composition
- Optimise rendering performance

**Input:** `.jml` files in the `components/` directory  
**Output:** JavaScript modules in the build directory

### 3. Module Compiler (`internal/mc`)

Compiles performance-critical code to WebAssembly.

**Key Responsibilities:**
- Handle module imports/exports
- Manage memory between JavaScript and WASM
- Optimise WASM module size and performance

**Input:** Source files in supported languages  
**Output:** `.wasm` and JavaScript glue code

## Compilation Pipeline

1. **Initialisation**
   - Load project configuration
   - Resolve dependencies
   - Set up compilation context

2. **Dependency Resolution**
   - Build dependency graph
   - Determine compilation order
   - Cache previous compilation results if available

3. **Compilation**
   - Parse JML syntax
   - Validate component structure
   - Generate intermediate representation
   - Optimise output

4. **Code Generation**
   - Generate target code (HTML/JS/WASM)
   - Apply optimisations

## Configuration

The compilation process can be configured using `jawt.config.json`:

```json
{
  "compilerOptions": {
    "outputDir": "dist",
    "minify": true,
  },
  "pages": {
    "dir": "app",
    "extensions": [".jml"]
  },
  "components": {
    "dir": "components",
    "extensions": [".jml"]
  }
}
```

## Error Handling

The compilation system provides detailed error messages including:
- Syntax errors in JML
- Missing components or dependencies
- Type mismatches in component props
- Circular dependencies
- Invalid WASM module exports
