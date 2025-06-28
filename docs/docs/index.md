# Just Another Web Tool

**JAWT is a web development toolchain that enables building minimal, performant web applications using a declarative approach.**

Write your applications using **JML**, a domain-specific language that compiles to optimised web standards—HTML for pages, JavaScript for interactive components, and WebAssembly for performance-critical modules.

## What is JAWT?

JAWT provides a unified development experience through its intelligent compilation system. Instead of juggling multiple technologies and build configurations, you write everything in JML and let JAWT handle the complexity of modern web development.

### Key Features

- **Single Language, Multiple Targets**: Write JML code that compiles to HTML, JavaScript, and WebAssembly as needed
- **Component-Driven Architecture**: Build reusable components with clear interfaces and composition patterns
- **Zero Configuration**: Works out of the box with sensible defaults whilst remaining customisable
- **Hot Module Replacement**: See changes instantly during development with intelligent reloading
- **Optimised Builds**: Automatic code splitting, tree shaking, and performance optimisation
- **Type Safety**: Leverage TypeScript-like type checking for robust applications

## Understanding JML

JML combines the structural clarity of markup languages with TypeScript-like scripting capabilities. Every JML file begins with a document type declaration that determines its compilation target:

```jml
_doctype page home          // Compiles to HTML document
_doctype component Button   // Compiles to JavaScript web component  
_doctype module calculator  // Compiles to WebAssembly module
```

### Document Types Explained

**Pages** define complete web pages with metadata and structure. They compile to HTML documents and serve as your application's entry points:

```jml
_doctype page dashboard

import component Layout from "components/layout"

Page {
    title: "Dashboard - My App"
    description: "Application dashboard with analytics"
    
    Layout {
        section: "dashboard"
        showSidebar: true
    }
}
```

**Components** encapsulate reusable UI elements with properties, state, and event handling. They compile to modern JavaScript web components:

```jml
_doctype component UserCard

Container {
    style: "bg-white shadow-md rounded-lg p-6"
    
    Text {
        content: props.userName
        style: "text-xl font-semibold text-gray-800"
    }
    
    Button {
        text: "View Profile"
        onClick: () => navigateToProfile(props.userId)
        style: "mt-4 bg-blue-500 text-white px-4 py-2 rounded"
    }
}
```

**Modules** handle computational logic and performance-critical operations. They compile to WebAssembly for near-native performance:

```jml
_doctype module imageProcessor

export function processImage(data: ImageData): ImageData {
    // Heavy image processing logic
    return optimiseImage(data)
}

function optimiseImage(data: ImageData): ImageData {
    // WebAssembly-optimised processing
}
```

## Development Workflow

JAWT's CLI provides everything you need to build applications efficiently:

1. **Create** a new project: `jawt init my-app`
2. **Develop** with hot reload: `jawt run`
3. **Build** for production: `jawt build`
4. **Debug** when needed: `jawt debug`

The unified compiler handles all document types intelligently, resolving dependencies across your entire application and generating optimised output for each target.

## Architecture Philosophy

JAWT follows a clear separation of concerns:

- **Pages** handle structure and routing
- **Components** manage user interaction and state
- **Modules** provide computational performance

This architecture enables optimal loading strategies—pages load instantly, components activate when needed, and modules execute computations at near-native speed.

## Quick Start Example

Here's a complete JAWT application structure:

```
my-app/
├── app/
│   ├── index.jml              # Home page
│   └── about/index.jml        # About page  
├── components/
│   ├── layout.jml             # Shared layout
│   └── user-card.jml          # Reusable component
└── modules/
    └── analytics.jml          # Performance module
```

**Page** (`app/index.jml`):
```jml
_doctype page home

import component Layout from "components/layout"

Page {
    title: "Welcome to My App"
    
    Layout {
        showWelcome: true
    }
}
```

**Component** (`components/layout.jml`):
```jml
_doctype component Layout

import module analytics from "modules/analytics"

Container {
    style: "min-h-screen bg-gray-50"
    
    Header {
        style: "bg-white shadow-sm p-4"
        
        if (props.showWelcome) {
            Text {
                content: "Welcome!"
                style: "text-2xl font-bold"
            }
        }
    }
    
    onClick: () => analytics.trackPageView()
}
```

## Browser Support

JAWT generates modern web standards that work across all current browsers:

- **HTML5**: Semantic, accessible markup
- **ES2020+**: Modern JavaScript with automatic polyfills
- **WebAssembly**: Supported in all major browsers since 2017
- **CSS Grid/Flexbox**: Modern layout with Tailwind CSS integration

## Next Steps

Ready to start building with JAWT? Here's where to go next:

### Getting Started
- **[Installation & Setup](getting-started/installation.md)** - Install JAWT and create your first project
- **[Tutorial](tutorial/first-page.md)** - Build a simple application step by step
- **[Project Structure](getting-started/project-structure.md)** - Understand how JAWT projects are organised

### Language Reference
- **[JML Syntax](jml/syntax.md)** - Complete JML language specification
- **[Pages](jml/pages.md)** - Creating pages and handling routing
- **[Components](jml/components.md)** - Building interactive components
- **[Modules](jml/modules.md)** - Using modules for computational logic
- **[Scripts](jml/scripts.md)** - Writing scripts for interacting with components

### Advanced Topics
- **[Architecture](architecture/index.md)** - Understanding JAWT's compilation system
- **[CLI Reference](architecture/cli.md)** - Complete command-line interface guide
- **[Configuration](architecture/configuration.md)** - Customising your build process
- **[Deployment](deployment/index.md)** - Publishing your applications (Soon)

### Resources
- **[Examples](examples/index.md)** - Sample applications and patterns
- **[Migration Guide](resources/migration.md)** - Moving from other tools
- **[FAQ](resources/faq.md)** - Common questions and solutions
- **[Contributing](contributing/index.md)** - Help improve JAWT