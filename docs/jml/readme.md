# JML Language Overview

**JML is a domain-specific language for Jawt.** It combines the structural clarity of markup languages with the power of TypeScript-like scripting, creating a unified development experience for web applications.

## Design Philosophy

JML embraces a declarative paradigm, where you describe *what* your application should look like and behave like, rather than *how* to build it step by step. This approach leads to more maintainable and intuitive code that closely mirrors the final user interface structure.

## Language Structure

JML syntax is fundamentally based on a subset of TypeScript, providing familiar programming constructs whilst maintaining the declarative nature essential for UI description. The language seamlessly integrates three core concepts:

- **Declarative UI Structure**: Components are described using a hierarchical, tree-like syntax
- **Styling Integration**: Built-in support for Tailwind CSS utility classes through the `style` property
- **Scripting Capabilities**: TypeScript-like functions and logic for interactive behaviour

## Core Components

### Document Types

Every JML file begins with a document type declaration that determines its compilation target:

```jml
_doctype page index        // Compiles to HTML page
_doctype component Button  // Compiles to JavaScript web component
_doctype module math       // Compiles to WebAssembly module
```

### Component Syntax

Components follow a declarative structure where properties and children are defined within curly braces:

```jml
Container {
    style: "flex flex-col items-center p-4"
    
    Text {
        content: "Welcome to JAWT"
        style: "text-2xl font-bold text-blue-600"
    }
    
    Button {
        text: "Click Me"
        onClick: handleClick
        style: "bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded"
    }
}
```

### Styling System

JML uses Tailwind CSS utility classes for styling, applied through the `style` property. This provides a consistent, utility-first approach to styling that compiles to optimised CSS:

```jml
Card {
    style: "bg-white shadow-lg rounded-lg p-6 hover:shadow-xl transition-shadow"
    
    // Component content
}
```

## Compilation Targets

JML's unique architecture allows different parts of your application to compile to different targets based on their purpose:

### Pages → HTML
Pages represent complete web pages and compile to standard HTML documents. They can import and use components, creating the overall structure of your application.

### Components → JavaScript
Components compile to modern JavaScript web components and ES6 modules. They handle DOM interaction, user events, and can call WebAssembly modules for performance-critical operations.

### Modules → WebAssembly
Modules compile to WebAssembly for performance-critical calculations and heavy computations. They cannot directly interact with the DOM but can be called from components.

## Import System

JML provides a straightforward import system for managing dependencies between different parts of your application:

```jml
// Import from project components directory
import component Layout from "components/layout"

// Import from same directory (page-specific)
import component Card from "card"

// Import browser APIs (components only)
import browser
```

## Scripting Integration

JML seamlessly integrates TypeScript-like scripting for interactive behaviour:

```jml
_doctype component InteractiveButton

Button {
    onClick: () => showMessage()
    text: "Greet"
    style: "bg-green-500 text-white px-4 py-2 rounded"
}

function showMessage(): void {
    browser.Alert("Hello from JAWT!")
}
```

## Development Experience

JML is designed with developer experience in mind:

- **Hot Reload**: See changes instantly during development
- **Component Reusability**: Build once, use anywhere in your application  
- **Type Safety**: Leverage TypeScript's type system for robust applications

## Architecture Benefits

The JML architecture provides several key advantages:

- **Clear Separation**: Pages handle structure, components manage interaction, modules provide performance
- **Optimised Output**: Each compilation target is optimised for its specific purpose
- **Seamless Integration**: JavaScript and WebAssembly work together transparently
- **Minimal Bundle Size**: Only include what you actually use

## Getting Started

To begin working with JML, you'll typically start with a page definition and build up your application using components:

```jml
_doctype page home

import component Header from "components/header"
import component MainContent from "components/main"

Page {
    title: "My JAWT Application"
    
    Header {
        title: "Welcome"
    }
    
    MainContent {
        style: "container mx-auto px-4"
    }
}
```

## Further Reading

For detailed information about specific aspects of JML:

- **[Pages](pages.md)** - Page structure and routing
- **[Components](components.md)** - Component development
- **[Modules](modules.md)** - WebAssembly module creation
- **[Scripting](scripting.md)** - JML's scripting capabilities

