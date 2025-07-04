# Just Another Web Tool (JAWT)

**Jawt is a development toolchain that enables building web applications using a declarative approach.**

Write your applications using **JML**, a domain-specific language for Jawt app structure, combined with TypeScript for dynamic functionality.

## What is Jawt?

Jawt provides a unified development experience through its intelligent compilation system. Instead of juggling multiple technologies and build configurations, you write everything in Jml and TypeScript, letting Jawt handle the complexity of modern web development.

### Key Features

- **Unified Language Approach**: Write JML for structure and TypeScript for logic
- **Component-Driven Architecture**: Build reusable components with clear interfaces and composition patterns
- **Zero Configuration**: Works out of the box with sensible defaults whilst remaining customisable
- **Hot Module Replacement**: See changes instantly during development with intelligent reloading
- **Optimised Builds**: Automatic code splitting, tree shaking, and performance optimisation
- **Type Safety**: Full TypeScript support for robust applications

## Understanding JML

JML combines the structural clarity of markup languages with component-driven development. Every JML file begins with a document type declaration that determines its purpose:

```jml
_doctype page home          // Defines a complete web page
_doctype component Button   // Defines a reusable component
```

### Document Types Explained

**Pages** define complete web pages with metadata and structure. They serve as your application's entry points:

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

**Components** encapsulate reusable UI elements with properties, state, and event handling:

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

## TypeScript Integration

JAWT seamlessly integrates TypeScript for dynamic functionality. Write your scripts in TypeScript and import them directly into your JML components:

```typescript
// scripts/analytics.ts
export function trackPageView(page: string): void {
    // Analytics logic
}

export function trackUserAction(action: string, data?: any): void {
    // User interaction tracking
}
```

```jml
_doctype component Analytics

import script analytics from "scripts/analytics"

Container {
    onClick: () => analytics.trackUserAction("button_click", { id: props.buttonId })
}
```

## Development Workflow

JAWT's CLI provides everything you need to build applications efficiently:

1. **Create** a new project: `jawt init my-app`
2. **Develop** with hot reload: `jawt run`
3. **Build** for production: `jawt build`
4. **Debug** when needed: `jawt debug`

The unified compiler handles all document types intelligently, resolving dependencies across your entire application and generating optimised output.

## Architecture Philosophy

JAWT follows a clear separation of concerns:

- **Pages** handle structure and routing
- **Components** manage user interaction and state
- **Scripts** provide dynamic functionality and business logic

This architecture enables optimal loading strategies and maintainable code organisation.

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
└── scripts/
    └── analytics.ts           # TypeScript functionality
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

import script analytics from "scripts/analytics"

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
    
    onClick: () => analytics.trackPageView("home")
}
```

**Script** (`scripts/analytics.ts`):
```typescript
export function trackPageView(page: string): void {
    console.log(`Page view: ${page}`)
    // Analytics implementation
}
```

## Browser Support

JAWT generates modern web standards that work across all current browsers with automatic polyfills and optimisation.

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
- **[Scripts](jml/scripts.md)** - Writing TypeScript scripts for dynamic functionality

### Advanced Topics
- **[Architecture](architecture/index.md)** - Understanding JAWT's compilation system
- **[CLI Reference](references/cli.md)** - Complete command-line interface guide
- **[Configuration](architecture/configuration.md)** - Customising your build process
- **[Deployment](deployment/index.md)** - Publishing your applications

### Resources
- **[Examples](examples/index.md)** - Sample applications and patterns
- **[Migration Guide](resources/migration.md)** - Moving from other tools
- **[FAQ](resources/faq.md)** - Common questions and solutions
- **[Contributing](contributing/index.md)** - Help improve JAWT