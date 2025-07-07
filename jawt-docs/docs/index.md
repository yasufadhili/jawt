# Welcome to JAWT

So you're here to check out JAWT. Awesome.

In a nutshell, **JAWT is a toolchain for building web apps in a more straightforward, declarative way.** I got tired of the endless setup and complexity that comes with modern web development, so I built JAWT around a simple idea: you should be able to describe your application's structure and logic, and the tool should handle the rest.

You'll be writing your apps using **JML**, a language I designed specifically for this, and **TypeScript** for all the dynamic bits.

## What's Jawt All About?

Instead of wrestling with bundlers, transpilers, and a million config files, Jawt gives you a single, unified workflow. You write JML for structure, TypeScript for logic, and Jawt's compiler figures out how to turn it all into an optimised, fast web app.

### Here's the Gist:

-   **One Language for Structure**: Use JML to lay out your pages and components. It's clean and easy to read.
-   **Component-Based Everything**: Build your UI out of reusable components. It's a tried-and-true way to keep things organised.
-   **Zero-Config, but Hackable**: It works right away with smart defaults, but you can still tweak things if you need to.
-   **Instant Feedback**: The dev server comes with hot module replacement, so you see your changes live.
-   **Optimised by Default**: It automatically handles things like code splitting and tree shaking to make sure your app is fast.
-   **Full-Fat TypeScript**: No compromises. Use all the TypeScript features you know and love.

## How JML Works

JML is the core of JAWT. Every JML file has a `_doctype` that tells the compiler what it is. There are two main types:

```jml
_doctype page home          // This is a whole web page.
_doctype component Button   // This is a reusable building block.
```

### Pages

Pages are the main entry points of your app. They define a full web page with its metadata and structure.

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

### Components

Components are the reusable UI bits. You build them once and use them everywhere.

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

## Bringing in the Logic with TypeScript

Jawt wouldn't be complete without a way to handle dynamic logic. That's where TypeScript comes in. You can write your functions and classes in `.ts` files and import them right into your JML.

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

## The Workflow

The `jawt` CLI is your best friend here. It's got everything you need:

1.  **Start a new project:** `jawt init my-app`
2.  **Run the dev server:** `jawt run`
3.  **Build for production:** `jawt build`
4.  **Debug your app:** `jawt debug`

The compiler is smart enough to understand all the different document types and how they depend on each other, spitting out a nice, optimised app at the end.

## The Philosophy

The idea behind Jawt's architecture is to keep things separate and clean:

-   **Pages** for structure and routing.
-   **Components** for UI and user interaction.
-   **Scripts** for business logic and dynamic stuff.

This makes your code easier to reason about and helps Jawt optimise how it loads everything.

## Quick Start Example

Here's what a simple Jawt app looks like:

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

Jawt spits out modern, standards-compliant code that works in all the usual suspects (Chrome, Firefox, Safari, Edge). It'll even add polyfills where needed.

## Where to Next?

Ready to dive in?

### Getting Started
- **[Installation & Setup](getting-started/installation.md)** - Get Jawt installed and make your first project.
- **[Tutorial](tutorial/first-page.md)** - Build a simple app from scratch.
- **[Project Structure](getting-started/project-structure.md)** - Learn how to organise your Jawt projects.

### Language Reference
- **[JML Syntax](jml/syntax.md)** - The full JML language spec.
- **[Pages](jml/pages.md)** - All about creating pages and routes.
- **[Components](jml/components.md)** - Master building interactive components.
- **[Scripts](jml/scripts.md)** - Using TypeScript for dynamic functionality.

### Advanced Stuff
- **[Architecture](architecture/index.md)** - A peek under the hood.
- **[CLI Reference](references/cli.md)** - A guide to all the CLI commands.
- **[Configuration](architecture/configuration.md)** - How to customise your build.
- **[Deployment](deployment/index.md)** - Getting your app out into the world.

### Resources
- **[Examples](examples/index.md)** - Sample apps and code patterns.
- **[Migration Guide](resources/migration.md)** - Coming from another tool?
- **[FAQ](resources/faq.md)** - Answers to common questions.
- **[Contributing](contributing/index.md)** - Want to help out?
