# Emitter (`internal/emitter`)

The `internal/emitter` package is responsible for the code generation phase of the JAWT compilation process. Its primary responsibility is to take the Abstract Syntax Tree (AST) produced by the `compiler` package and transform it into executable web code.

## Purpose

The `emitter` package generates different output based on the JML document type:

*   **Pages (`_doctype page`)**: Emitted as standalone HTML files, including necessary `<script>` and `<style>` tags to link to generated JavaScript and CSS.
*   **Components (`_doctype component`)**: Emitted as Lit Components (a lightweight library for building web components). The generated Lit component code is primarily TypeScript, which is then processed by the TypeScript compiler.

This package bridges the gap between the JML AST and the final browser-executable code, leveraging existing web standards and libraries (like Lit) where appropriate.

## How it Works

1.  **AST Traversal**: The emitter traverses the JML AST, visiting each node.
2.  **Target-Specific Code Generation**: Based on the node type and the overall document type (page or component), the emitter generates corresponding HTML, JavaScript (specifically Lit component code for JML components), or CSS.
3.  **Lit Component Generation**: For JML components, the emitter generates TypeScript code that defines a custom HTML element using Lit's `LitElement` base class. This includes mapping JML properties to Lit properties, handling state, and rendering the component's template.
4.  **TypeScript Compiler Integration**: The generated TypeScript code (for Lit components and any inline scripts) is then passed to the TypeScript compiler (managed by the `process` package) for transpilation into browser-compatible JavaScript.
5.  **Output Placement**: The final HTML, JavaScript, and CSS files are saved to the appropriate output directories as configured in `jawt.project.json` (e.g., `build/`, `dist/`).

## Key Responsibilities

*   Transforming JML syntax into standard web technologies.
*   Generating efficient and performant code.
*   Ensuring proper integration with the TypeScript compilation pipeline.
*   Handling asset references and paths within the generated output.

## Future Implementation Details

When fully implemented, the `emitter` package will contain:

*   An `Emitter` struct with methods like `EmitPage(doc *ast.Document)` and `EmitComponent(doc *ast.Document)`.
*   Specialized visitor implementations (e.g., `HtmlEmitterVisitor`, `LitComponentEmitterVisitor`) to walk the AST and generate code for each node type based on the target format.
*   Logic to manage imports and dependencies within the generated code, ensuring all required modules are correctly bundled or referenced.
*   Configuration options to control output format, minification, and other build optimizations during the emission process.

**Example (Conceptual - Lit Component Emission):**

```go
// package emitter

// import (
// 	"github.com/yasufadhili/jawt/internal/ast"
// 	"github.com/yasufadhili/jawt/internal/core"
// 	"io"
// )

// type Emitter struct {
// 	ctx    *core.JawtContext
// 	writer io.Writer
// }

// func NewEmitter(ctx *core.JawtContext, writer io.Writer) *Emitter {
// 	return &Emitter{
// 		ctx:    ctx,
// 		writer: writer,
// 	}
// }

// func (e *Emitter) EmitComponent(doc *ast.Document) error {
// 	// Logic to traverse the AST and generate Lit component TypeScript code
// 	// Example: Generate a class extending LitElement, define properties, and render method
// 	// This generated TS code will then be compiled by the TypeScript compiler
// 	generatedTS := `
// 		import { LitElement, html, css } from 'lit';
// 		import { property } from 'lit/decorators.js';

// 		export class MyComponent extends LitElement {
// 			@property({ type: String }) name = 'World';

// 			static styles = css`
// 				:host { display: block; padding: 16px; }
// 				.greeting { color: blue; }
// 			`;

// 			render() {
// 				return html`<div class="greeting">Hello, ${this.name}!</div>`;
// 			}
// 		}
// 		customElements.define('my-component', MyComponent);
// 	`;
// 	
// 	// Write generatedTS to a .ts file in the build directory
// 	// The build system will then pick this up for TypeScript compilation
// 	return e.ctx.Paths.WriteBuildFile(doc.Name.Name + '.ts', generatedTS)
// }

// func (e *Emitter) EmitPage(doc *ast.Document) error {
// 	// Logic to traverse the AST and generate HTML file
// 	// This will include references to compiled Lit components and other assets
// 	generatedHTML := `
// 		<!DOCTYPE html>
// 		<html>
// 		<head>
// 			<title>${doc.Name.Name}</title>
// 			<script type="module" src="/build/components/${doc.Name.Name}.js"></script>
// 		</head>
// 		<body>
// 			<my-component name="JAWT"></my-component>
// 		</body>
// 		</html>
// 	`;
// 	
// 	// Write generatedHTML to a .html file in the build directory
// 	return e.ctx.Paths.WriteBuildFile(doc.Name.Name + '.html', generatedHTML)
// }
```