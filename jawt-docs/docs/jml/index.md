# JML Overview

**JML is the heart of JAWT development.** It's a declarative language that combines structural clarity with the programming power of TypeScript, designed specifically for building modern web applications.

## What Makes JML Special

Think of JML as a blueprint language for your applications. JML lets you describe what your application should be rather than how to build it piece by piece.

### The Two Pillars of JML

JML is built around two core document types, each serving a distinct purpose:

1. **Pages** — The foundations of your application
2. **Components** — The reusable building blocks

This separation is like having specialised tools in a workshop: pages handle structure and routing, whilst components provide reusable functionality and interactivity.

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
- Handle routing and page metadata
- Serve as application entry points

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
- Can import other components and scripts
- Encapsulate reusable functionality

## The Import System

JML's import system allows you to easily combine components and integrate TypeScript functionality:

```jml
// Import components from the global components directory
import component Button from "components/button"
import component Card from "components/card"

// Import from the same directory (page-specific components)
import component LocalWidget from "widget"

// Import TypeScript scripts for dynamic functionality
import script analytics from "scripts/analytics"
import script utils from "scripts/utils"

// Import browser APIs (components only)
import browser
```

## TypeScript Integration

JML seamlessly integrates with TypeScript through the script import system. Write your dynamic functionality in TypeScript and import it directly into your JML components:

```typescript
// scripts/counter.ts
export class Counter {
    private count: number = 0
    
    increment(): number {
        return ++this.count
    }
    
    decrement(): number {
        return --this.count
    }
    
    getValue(): number {
        return this.count
    }
}

export function formatCount(count: number): string {
    return `Count: ${count}`
}
```

```jml
_doctype component CounterWidget

import script counter from "scripts/counter"

Container {
    style: "flex items-center space-x-4"
    
    Button {
        text: "-"
        onClick: () => handleDecrement()
        style: "bg-red-500 text-white px-3 py-1 rounded"
    }
    
    Text {
        content: counter.formatCount(currentCount)
        style: "font-mono text-lg"
    }
    
    Button {
        text: "+"
        onClick: () => handleIncrement()
        style: "bg-green-500 text-white px-3 py-1 rounded"
    }
}

const counterInstance = new counter.Counter()
let currentCount = 0

function handleIncrement(): void {
    currentCount = counterInstance.increment()
}

function handleDecrement(): void {
    currentCount = counterInstance.decrement()
}
```

## Styling with Tailwind CSS

JML integrates seamlessly with Tailwind CSS through the `style` property:

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
_doctype component TodoList

import script todoManager from "scripts/todo-manager"

Container {
    style: "max-w-md mx-auto"
    
    Input {
        placeholder: "Add a new task..."
        onEnter: (value) => addTodo(value)
        style: "w-full p-2 border rounded mb-4"
    }
    
    List {
        style: "space-y-2"
        
        for (todo in todos) {
            TodoItem {
                text: todo.text
                completed: todo.completed
                onToggle: () => toggleTodo(todo.id)
                onDelete: () => deleteTodo(todo.id)
            }
        }
    }
}

let todos = []

function addTodo(text: string): void {
    todos = todoManager.addTodo(todos, text)
}

function toggleTodo(id: string): void {
    todos = todoManager.toggleTodo(todos, id)
}

function deleteTodo(id: string): void {
    todos = todoManager.deleteTodo(todos, id)
}
```

## Development Experience

### Hot Module Replacement

Changes to your JML files are reflected instantly in the browser, creating a fluid development experience where you can see your changes immediately.

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
├── scripts/               # TypeScript functionality
│   ├── analytics.ts
│   ├── utils.ts
│   └── api.ts
└── jawt.config.json       # Configuration
```

## Key Benefits

### Unified Development Experience

Instead of juggling HTML, CSS, JavaScript, and build configurations, JML provides a single, cohesive language for describing your application structure, with TypeScript handling the dynamic functionality.

### Optimised Output

Each document type compiles to its optimal target:

- Pages become lean HTML with proper metadata
- Components become efficient Web Components
- TypeScript scripts provide full JavaScript functionality

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
- **[Scripts](./scripts.md)** — Explore TypeScript integration and dynamic functionality

---

**Version**: Early Development  
**Last Updated**: Development Roadmap Phase 1