# JML Language Overview

**JML is the heart of JAWT development.** It's a declarative language that combines structural clarity with the programming power of TypeScript, designed specifically for building modern applications.

## What Makes JML Special

Think of JML as a blueprint language for your applications. JML lets you describe what your application should be rather than how to build it piece by piece.

### The Three Pillars of JML

JML is built around three core document types, each serving a distinct purpose:

1. **Pages** — The foundations of your application
2. **Components** — The reusable building blocks
3. **Modules** — The performance engines

This separation is like having specialised tools in a workshop: you wouldn't use a hammer for precision work, nor tweezers for heavy lifting. Each document type excels at its intended purpose.

## Language Philosophy

### Declarative by Design

JML embraces a declarative approach where you describe the desired outcome rather than the steps to achieve it:

```jml
// You describe WHAT you want
Card {
    title: "Welcome"
    style: "bg-white shadow-lg rounded-lg p-6"
    
    Text {
        content: "Hello, world!"
        style: "text-gray-700"
    }
}
```

Compare this to imperative approaches where you'd manually create elements, set attributes, append children, and manage the DOM step by step.

### TypeScript-Inspired Syntax

JML borrows heavily from TypeScript's syntax, making it familiar to developers whilst maintaining its declarative nature:

```jml
// Variables and functions work like TypeScript
const greeting = "Hello, JAWT!"

function handleClick(): void {
    console.log("Button clicked!")
}

// But UI structure is declarative
Button {
    text: greeting
    onClick: handleClick
    style: "bg-blue-500 text-white px-4 py-2 rounded"
}
```

## Document Types Explained

### Pages: Your Application's Entry Points

Pages represent complete web pages and serve as the entry points to your application.

```jml
_doctype page home

import component Layout from "components/layout"

Page {
    title: "Welcome to My App"
    description: "A minimal web application built with JAWT"
    
    Layout {
        content: "Hello, world!"
    }
}
```

**Key characteristics:**

- Must begin with `_doctype page`

- Can only have one direct child component
- Compile to complete HTML documents
- Handle routing and page metadata

### Components: Reusable Building Blocks

Components are like Lego bricks—individual pieces that can be combined in countless ways to build complex structures. They encapsulate both appearance and behaviour.

```jml
_doctype component UserCard

Card {
    style: "border rounded-lg p-4 shadow"
    
    Avatar {
        src: props.avatarUrl
        alt: props.name
        style: "w-12 h-12 rounded-full"
    }
    
    Text {
        content: props.name
        style: "font-semibold text-lg"
    }
    
    Text {
        content: props.email
        style: "text-gray-600"
    }
}
```

**Key characteristics:**

- Begin with `_doctype component`
- Accept props for customisation
- Compile to JavaScript web components
- Can import other components and modules

### Modules: Performance Powerhouses

Modules are like specialised engines—they handle computationally intensive tasks with maximum efficiency. Think of them as the high-performance components.

```jml
_doctype module calculations

export function fibonacci(n: number): number {
    if (n <= 1) return n
    return fibonacci(n - 1) + fibonacci(n - 2)
}

export function isPrime(num: number): boolean {
    if (num <= 1) return false
    for (let i = 2; i <= Math.sqrt(num); i++) {
        if (num % i === 0) return false
    }
    return true
}
```

**Key characteristics:**

- Begin with `_doctype module`
- Compile to WebAssembly for maximum performance
- Cannot directly interact with the DOM
- Export functions for use by components

## The Import System

JML's import system is like a well-organised library—you can easily find and use exactly what you need:

```jml
// Import components from the global components directory
import component Button from "components/button"
import component Card from "components/card"

// Import from the same directory (page-specific components)
import component LocalWidget from "widget"

// Import modules for performance-critical operations
import module MathUtils from "modules/math"

// Import browser APIs (components only)
import browser
```

## Styling with Tailwind CSS

JML integrates seamlessly with Tailwind CSS through the `style` property — you can achieve any design without mixing from scratch:

```jml
Container {
    style: "max-w-4xl mx-auto px-4 py-8"
    
    Header {
        style: "text-3xl font-bold text-gray-900 mb-6"
        content: "My Application"
    }
    
    Grid {
        style: "grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6"
        
        // Grid items...
    }
}
```

## Interactive Behaviour

JML makes adding interactivity as natural as describing static content:

```jml
_doctype component Counter

Container {
    style: "flex items-center space-x-4"
    
    Button {
        text: "-"
        onClick: () => setCount(count - 1)
        style: "bg-red-500 text-white px-3 py-1 rounded"
    }
    
    Text {
        content: `Count: ${count}`
        style: "font-mono text-lg"
    }
    
    Button {
        text: "+"
        onClick: () => setCount(count + 1)
        style: "bg-green-500 text-white px-3 py-1 rounded"
    }
}

let count = 0

function setCount(newCount: number): void {
    count = newCount
    // Component automatically re-renders
}
```

## Development Experience

### Hot Module Replacement

Changes to your JML files are reflected instantly in the browser, like having a live preview that updates as you type. This creates a fluid development experience where you can see your changes immediately.

### Type Safety

JML leverages TypeScript's type system to catch errors early:

```jml
// TypeScript-style type annotations
function processUser(user: { name: string; age: number }): string {
    return `${user.name} is ${user.age} years old`
}

// Props are type-checked automatically
UserCard {
    name: "John Doe"        // ✅ Correct type
    age: "thirty"           // ❌ Type error: expected number
}
```

### Component Composition

JML encourages building complex UIs through composition—combining simple components into more sophisticated ones:

```jml
_doctype component BlogPost

Article {
    style: "max-w-2xl mx-auto"
    
    PostHeader {
        title: props.title
        author: props.author
        publishDate: props.date
    }
    
    PostContent {
        content: props.content
        style: "prose prose-lg"
    }
    
    PostFooter {
        tags: props.tags
        shareUrl: props.url
    }
}
```

## Project Organisation

A typical JAWT project follows a clear structure that mirrors the language's architecture:

```
my-jawt-app/
├── app/                    # Pages (routes)
│   ├── index.jml          # Home page
│   ├── about/
│   │   └── index.jml      # About page
│   └── blog/
│       ├── index.jml      # Blog listing
│       └── [slug].jml     # Individual blog posts
├── components/             # Reusable components
│   ├── layout.jml
│   ├── button.jml
│   └── card.jml
├── modules/               # WebAssembly modules
│   ├── math.jml
│   └── image-processing.jml
└── jawt.config.json       # Configuration
```

## Key Benefits

### Unified Development Experience

Instead of juggling HTML, CSS, JavaScript, and build configurations, JML provides a single, cohesive language for describing your entire application.

### Optimised Output

Each document type compiles to its optimal target:

- Pages become lean HTML with proper metadata
- Components become efficient Web Components
- Modules become high-performance WebAssembly

### Maintainable Code

The declarative nature and clear separation of concerns makes JML applications easier to understand, modify, and extend over time. Code reads like a description of what the application does rather than a complex set of instructions.

## Getting Started

The best way to understand JML is to start with a simple page and gradually add components:

```jml
_doctype page welcome

Page {
    title: "Welcome to JAWT"
    
    Container {
        style: "min-h-screen flex items-center justify-center bg-gray-100"
        
        Card {
            style: "bg-white p-8 rounded-lg shadow-lg text-center"
            
            Text {
                content: "Hello, JAWT!"
                style: "text-2xl font-bold text-gray-900 mb-4"
            }
            
            Text {
                content: "Your first JML application is running."
                style: "text-gray-600"
            }
        }
    }
}
```

## What's Next?

This overview provides the foundation for understanding JML. To dive deeper into specific aspects:

- **[Pages](./pages.md)** — Learn about routing, metadata, and page structure
- **[Components](./components.md)** — Master props, state, and component lifecycle
- **[Modules](./modules.md)** — Harness WebAssembly for performance-critical code
- **[Scripting](./scripts.md)** — Explore JML's programming capabilities in detail

---

**Version**: Early Development  
**Last Updated**: Development Roadmap Phase 1