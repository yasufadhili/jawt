# The Emitter (`internal/emitter`)

The `emitter` is where the magic happens. It takes the AST that the compiler and checker have worked on and turns it into actual, runnable web code. This is the code generation phase.

## The Goal

The emitter has two main jobs, depending on the type of JML document:

-   **Pages (`_doctype page`)**: These get turned into standalone HTML files. The emitter will also inject the necessary `<script>` and `<style>` tags to link to the compiled JavaScript and CSS.
-   **Components (`_doctype component`)**: These get turned into Lit Components. Lit is a great little library from Google for building web components. The emitter generates TypeScript code that defines a Lit component, and then the TypeScript compiler takes over from there.

## How It Works

1.  **AST Traversal**: The emitter walks the AST, node by node.
2.  **Code Generation**: For each node, it generates the corresponding HTML or TypeScript code.
3.  **Lit Component Generation**: For JML components, it generates a TypeScript class that extends `LitElement`. It maps JML properties to Lit properties, handles state, and creates the `render` method.
4.  **Output**: The final HTML, JavaScript, and CSS files are saved to the build directory.

## Styling

JML components have two ways to handle styles:

-   **Shadow DOM Styles**: If you define a `<style>` block inside a JML component, those styles will be encapsulated in the component's Shadow DOM. This is great for creating truly reusable components that don't leak styles.
-   **Light DOM Styles**: You can also pass styles to a component through its `style` prop. These are just standard HTML style attributes that get applied to the component's outer element.

JAWT also integrates with Tailwind CSS to process and optimize all the styles.

## Inbuilt Components

JAWT has a few built-in components that are written in TypeScript. These get compiled along with the user-defined components to make sure they're available and optimized.