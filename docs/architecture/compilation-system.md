# Compilation System

The JAWT compilation system uses a compiler that processes all JML source files and transforms them into web assets. 
The compiler intelligently handles different document types (pages, components, and modules) through specialised compilation paths within a cohesive architecture.

##  Compiler Architecture

### Compiler (`internal/compiler`)

The compiler is responsible for processing all JML files regardless of their document type, providing a consistent compilation experience across the entire codebase.

**Key Responsibilities:**
- Parse all JML document types (`_doctype page`, `_doctype component`, `_doctype module`)
- Generate appropriate output based on document type:
   - Pages → HTML files
   - Components → JavaScript modules
   - Modules → WebAssembly binaries
- Handle cross-document dependencies and imports
- Manage symbol resolution across the entire project
- Provide  error reporting and diagnostics

**Input:** All `.jml` files throughout the project  
**Output:** HTML files, JavaScript modules, and WASM binaries as appropriate

## Document Type Compilation Paths

The compiler uses the `_doctype` declaration to determine the appropriate compilation path.

### Page Compilation Path
```
JML Page → Parser → Semantic Analysis → HTML Generator → .html
```

- Processes `_doctype page` files
- Generates complete HTML documents with proper DOCTYPE and structure
- Handles page metadata (title, description, etc.)
- Manages routing and navigation structure
- Integrates imported components into the page structure

### Component Compilation Path
```
JML Component → Parser → Semantic Analysis → JS Generator → .js
```

- Processes `_doctype component` files
- Transforms components to JavaScript web components
- Handles component props, state, and lifecycle
- Supports component composition and reusability
- Optimises rendering performance

### Module Compilation Path
```
JML Module → Parser → Semantic Analysis → WASM Generator → .wasm + .js
```

- Processes `_doctype module` files
- Compiles performance-critical code to WebAssembly
- Generates JavaScript glue code for seamless integration
- Manages memory between JavaScript and WASM
- Optimises module size and execution performance

##  Compilation Pipeline

The compiler follows this comprehensive pipeline:

1. **Project Discovery**
   - Scan project directories for all `.jml` files
   - Build a complete dependency graph across all document types
   - Determine optimal compilation order

2. ** Parsing**
   - Parse all JML files using the same lexer and parser
   - Build abstract syntax trees (ASTs) for all documents
   - Share common parsing logic across document types

3. **Cross-Document Symbol Resolution**
   - Collect symbols from all files into a unified symbol table
   - Resolve imports and dependencies across document types
   - Validate that all references can be satisfied

4. **Semantic Analysis**
   - Perform type checking and validation
   - Ensure component interfaces match usage
   - Validate module exports and imports
   - Check for circular dependencies

5. **Specialised Code Generation**
   - Route each document to its appropriate code generator
   - Pages generate HTML with embedded component references
   - Components generate JavaScript with lifecycle management
   - Modules generate WebAssembly with JavaScript bindings

6. **Cross-Reference Resolution**
   - Link component references in pages to compiled JavaScript
   - Connect module calls in components to WASM binaries
   - Generate import/export manifests for the build system

## Benefits of this Architecture

### Consistency
- Single parsing and validation logic for all JML code
- Consistent error messages and diagnostics
- Uniform handling of imports and dependencies

### Performance
- Compile all files in a single pass
- Share parsed ASTs between compilation phases
- Optimise cross-document references during compilation

### Maintainability
- Single codebase for all compilation logic
- Easier to add new document types or features
- Centralised error handling and reporting

### Developer Experience
- Consistent compilation behaviour regardless of document type
-  error reporting across pages, components, and modules
- Single point of configuration for compilation options

## Configuration

The compiler is configured through `jawt.config.json`:

```json
{
  "compiler": {
    "outputDir": "dist",
    "minify": true,
    "sourceMap": true,
    "target": "es2020"
  },
  "paths": {
    "pages": "app",
    "components": "components",
    "modules": "modules"
  },
  "optimization": {
    "treeShaking": true,
    "deadCodeElimination": true,
    "wasmOptLevel": "O2"
  }
}
```

## Error Handling and Diagnostics

The compiler provides comprehensive error reporting:

### Syntax Errors
- Consistent error messages across all document types
- Precise line and column information
- Helpful suggestions for common mistakes

### Semantic Errors
- Cross-document reference validation
- Type mismatches in component props
- Missing imports or circular dependencies
- Invalid module exports

### Build Errors
- Failed code generation
- Output file conflicts
- Resource allocation issues

### Example Error Output
```
Error: Component 'UserCard' not found
  --> app/dashboard/index.jml:15:5
   |
15 |     UserCard {
   |     ^^^^^^^^ Component not found in any imported module
   |
   Help: Did you forget to import from "components/user-card"?
```

## Integration with Build System

The  compiler integrates seamlessly with the build system:

1. **Dependency Management**: Build system provides resolved dependency graph
2. **Incremental Compilation**: Only recompile changed files and their dependents
3. **Asset Pipeline**: Compiler output feeds into optimisation and bundling
4. **Development Server**: Hot module replacement based on document type changes

## Performance Optimisations

### Compilation Performance
- Parallel processing of independent files
- Cached ASTs for unchanged files
- Incremental symbol table updates
- Memory-efficient compilation passes

### Output Optimisation
- Dead code elimination across document boundaries
- Tree shaking for unused components and modules
- WASM size optimisation
- JavaScript minification and compression

## Future Enhancements

The architecture enables several potential improvements:

- **Incremental Type Checking**: Faster recompilation during development
- **Cross-Document Optimisation**: Inline small components, optimise module calls
- **Advanced Error Recovery**: Continue compilation after non-fatal errors
- **Plugin Architecture**: Allow custom compilation steps and transformations