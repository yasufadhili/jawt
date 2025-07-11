# JML: The Language of JAWT

**JML is the language you'll be speaking when you build with JAWT.** I designed it to be a declarative language that gets out of your way, mixing the structural clarity of markup with the power of TypeScript.

## What's the Big Deal with JML?

Think of JML as a way to write blueprints for your app. Instead of telling the browser *how* to draw everything step-by-step, you just describe *what* you want the final thing to look like.

### JML Supports:

* Declarative UI trees
* Embedded TypeScript logic
* Page and component document types
* Imports for components, scripts, runtime APIs, and stores
* For-loops and expressions
* Type-checked props

### The Two Flavours of JML

JML has two main "document types," each with a specific job:

1.  **Pages** — The main canvases for your application.
2.  **Components** — The reusable Lego bricks you build with.

This separation keeps things clean. Pages handle the overall structure and routing, while components handle the interactive bits and pieces.

## The Philosophy: Declarative First

JML is all about being declarative. You describe the end result, not the process.

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

This is a lot cleaner than manually creating DOM elements, setting attributes, and appending them to each other.

### Familiar Syntax

If you've used TypeScript, JML should feel pretty familiar. I borrowed a lot of its syntax for variables, functions, and types.

```jml
// Variables and functions feel like TypeScript
const greeting = "Hello, JAWT!"

function handleClick(): void {
    console.log("Button clicked!")
}

// But the UI part is declarative
Button {
    text: greeting
    onClick: handleClick
    style: "bg-blue-500 text-white px-4 py-2 rounded"
}
```

## The Document Types in Detail

### Pages: Your App's Entry Points

A page represents a whole web page. It's where you define things like the page title and pull in the components that make up the UI.

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

**Key things to remember:**

-   Always starts with `_doctype page`.
-   Can only have one component directly inside the `Page` block.
-   Handles routing and page-level metadata.

### Components: Your Reusable Building Blocks

Components are the best part. They're like custom HTML tags you can reuse everywhere. They can have their own logic and can be composed to create really complex UIs.

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

**Key things to remember:**

-   Starts with `_doctype component`.
-   Uses `props` to get data from its parent.
-   Can import other components and scripts.

## The Import System

JML lets you pull in components and TypeScript code easily.

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

This is where JML really shines. You can write complex logic in TypeScript and use it right inside your components.

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

JML is set up to work with Tailwind CSS out of the box. Just use the `style` property.

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

## Making Things Interactive

Adding interactivity feels natural in JML.

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

## The Developer Experience

### Hot Module Replacement

When you save a JML file, the changes show up in your browser instantly. It makes for a really smooth workflow.

### Type Safety

JML uses TypeScript's type system, so you get type checking for free.

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

### Composition is Key

JML is all about building big things from small, simple pieces.

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

## How to Organise Your Project

A typical JAWT project has a structure that makes sense with the language's design:

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

## Why Bother?

### A Unified Experience

With JML, you don't have to constantly switch between HTML, CSS, and JavaScript files and mindsets. It's one cohesive way to build your app.

### Optimised for You

Jawt's compiler is smart about how it builds your code. Pages become lean HTML, components become efficient Web Components, and your TypeScript just works.

### Code That's Easy to Read

Because it's declarative, JML code often reads like a description of the UI, which makes it easier to come back to later.

## Getting Your Hands Dirty

The best way to learn is by doing. Start with a simple page and build from there.

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

Now that you've got the basics of JML, you can dive deeper:

-   **[Pages](./pages.md)** — Learn about routing and page structure.
-   **[Components](./components.md)** — Master props, state, and all things component.
-   **[Scripts](./scripts.md)** — Get the full scoop on TypeScript integration.
