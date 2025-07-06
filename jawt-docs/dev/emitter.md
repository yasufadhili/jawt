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

## Styling in JML Components

JML components support two primary ways of applying styles:

*   **Shadow DOM Styles**: Style blocks defined directly within a JML component's definition are pure CSS and are encapsulated within the component's Shadow DOM. These styles are isolated and do not leak outside the component, ensuring predictable styling.
*   **Light DOM Styles (via `style` prop)**: Styles passed to a component via its `style` property (typically from TypeScript or another JML component) are applied to the component's Light DOM. These are standard HTML `style` attributes or CSS classes that affect the component's outer element and its unencapsulated content.

JAWT integrates with the Tailwind CSS CLI to process and optimize styles, ensuring that only necessary CSS is included in the final build.

## Key Responsibilities

*   Transforming JML syntax into standard web technologies.
*   Generating efficient and performant code.
*   Ensuring proper integration with the TypeScript compilation pipeline.
*   Handling asset references and paths within the generated output.
*   Managing Shadow DOM and Light DOM styling for components.

## Implementation Details

The `emitter` package will be structured into specialized files for clarity and maintainability:

*   `emitter/page.go`: Contains the logic for `EmitPage`, responsible for generating HTML files for JML pages.
*   `emitter/component.go`: Contains the logic for `EmitComponent`, responsible for generating TypeScript code for Lit Components from JML components.

### Inbuilt Components

JAWT includes a set of inbuilt components written directly in TypeScript. These components are bundled with the JAWT executable. During the initial build run, these inbuilt components are also passed to the TypeScript compiler to ensure they are available and optimized alongside user-defined components.

### Helper Methods

The emitter will utilize several helper methods to streamline code generation and ensure correctness, including (but not limited to):

*   `write(content string)`: A utility for writing generated code to the output stream.
*   `sanitiseComponentName(name string)`: Ensures component names adhere to valid naming conventions for web components.
*   `addInbuiltComponents(ctx *core.JawtContext)`: Manages the inclusion and compilation of JAWT's predefined TypeScript components.
*   Other utilities for handling imports, property serialization, and template generation.

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
