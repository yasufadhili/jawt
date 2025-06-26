# Jawt Page Compilation

The Page Compilation processes JML files located in the `app/` directory that define pages—such as `app/index.jml`, `app/about/index.jml`, or `app/user/[id].jml` for dynamic routes—and generates corresponding HTML files. It ensures that page-specific JML syntax is correctly translated into web-ready HTML, incorporating styling and dependencies as needed.

## How It Works

The Compiler receives a file path and an output path from the build system, which identifies page files and manages the overall compilation process. The compiler then transforms the provided JML file into HTML through a series of well-defined stages. It does not handle routing, nested route logic, or dependency management—these are delegated to the build system.

### Page File Requirements

For a JML file to be recognised and compiled as a page, it must meet the following criteria:
- **Location**: Reside in the `app/` directory, either directly (e.g., `app/index.jml`) or within a subfolder (e.g., `app/about/index.jml`).
- **Naming**: Named `index.jml` when inside a subfolder, reflecting the route structure (e.g., `app/contact/index.jml` for the `/contact` route).
- **Declaration**: Begin with `_doctype page <page-name>` as the first non-empty line, where `<page-name>` is a descriptive name for the page.

Examples of valid page files:
- `app/index.jml` → compiles to the root page (`/`)
- `app/about/index.jml` → compiles to `/about`
- `app/user/[id].jml` → compiles to a dynamic route (e.g., `/user/123`)

### Compilation Process

The Page Compiler follows these stages to transform a JML page file into HTML:

1. **Lexing and Parsing**  
   - The JML source code is processed using ANTLR4, to convert it into an abstract syntax tree (AST).

2. **Symbol Collection**  
   - Symbols such as variables, imported components, and modules are collected into a symbol table. This prepares the compiler for dependency checks and semantic validation.

3. **Semantic Analysis**  
   - The compiler verifies the correctness of the JML code by:
     - Ensuring the `_doctype page <page-name>` declaration is present and properly formatted.
     - Checking that imported components and modules exist and are resolved in the dependency graph (managed by the build system).
     - Validating the JML syntax and logic for errors.
   - This step guarantees that the page can be rendered correctly with all dependencies in place.

4. **Emit**  
   - The final HTML is generated from the AST, translating JML constructs into HTML elements and attributes. Tailwind utility classes specified in a property style are incorporated into the output, ensuring the page is styled as intended.

### Integration with the Build System

The Compiler operates as a focused tool within the JAWT ecosystem:
- **Input**: Receives the JML file path and output path from the build system.
- **Dependencies**: Relies on the prior compilation of components and modules (handled by the Component Compiler and Module Compiler) to ensure all referenced assets are available.
- **Routing**: The build system manages route detection (including nested and dynamic routes like `app/user/[id].jml`) and saves them in a routes table, passing only the compilation task to the Page Compiler.

### Styling

JML uses a property style for styling, accepting a string of Tailwind utility classes. For example:
```
Component {
    style: "flex justify-center items-center"
}
  
```
During the emit stage, the Compiler embeds these classes into the HTML, ensuring the output reflects the specified design.

## Key Features

- **Simplicity**: Compiles a single JML file into HTML without additional responsibilities like routing or bundling.
- **Dependency Validation**: Ensures all imported components and modules are resolved during semantic analysis.
- **Tailwind Integration**: Seamlessly processes Tailwind utility classes for styling.
- **Dynamic Route Support**: Compiles pages for dynamic routes (e.g., `app/user/[id].jml`), though routing logic is handled by the build system.

## Limitations and Assumptions

- The Page Compilation assumes that components and modules are compiled beforehand, as they are dependencies for pages.
- It does not perform optimisations like HTML minification; such tasks are left to the build system.
- Dynamic route parameters (e.g., `[id]`) are not processed by the Page Compiler; it generates static HTML, with dynamic behavior managed elsewhere in the toolchain.

## Example Workflow

1. **Project Structure**:
   ```
   app/
   ├── index.jml           # Root page
   ├── about/
   │   └── index.jml       # /about page
   └── user/
       └── [id].jml        # Dynamic route /user/:id
   ```

2. **Sample JML Page** (`app/index.jml`):
   ```
   _doctype page home

   import component MyComponent from "components/my"
   
   Page {
    title: "Home"

    Container {
        Text {
            content: "Welcome to JAWT"
        }
    }

    MyComponent {}

   }

   ```

3. **Output HTML**:
   ```html
   <!DOCTYPE html>
   <html>
   <head>
     <title>Home</title>
   </head>
   <body>
     <div class="container mx-auto p-4">
       <h1 class="text-2xl font-bold">Welcome to JAWT</h1>
       <my-component></my-component>
     </div>
   </body>
   </html>
   ```

## Potential Improvements

Soon enhancing the Page Compilation’s functionality within the JAWT toolchain through:
- **Error Reporting**: Provide detailed diagnostics during semantic analysis for better developer feedback.
- **Caching**: Implement a caching mechanism for parsed ASTs to speed up recompilation during development.
- **Validation Hooks**: Allow the build system to inject custom validation rules before the emit stage.

## Related Components

- [JML Documentation](../jml/readme.md) – JML syntax conventions.
- [Architecture Guide](../architecture/overview.md) – Ecosystem and toolchain design.

