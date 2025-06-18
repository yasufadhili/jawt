# JAWT - Just A Weird Web Tool

JAWT enables building minimal web applications using a declarative approach. 
Write your application structure and behaviour using Jawt Markup Language (JML), a domain-specific language for JAWT.

## Quick Start

Create a new project:
```bash
jawt init my-app
cd my-app
```

Run your application:
```bash
jawt run
```

Build for production:
```bash
jawt build
```

## JML Syntax

### Pages

Pages define the entry points of your application. Each page specifies metadata and content structure:

```jml
_doctype page index

import Layout from "components/layout"

Page {
  title: "Welcome to My App"
  description: "A simple web application built with JAWT"
  
  Layout {
    content: "Hello, world!"
  }
}
```

### Components

Components are reusable building blocks that encapsulate layout and styling:

```jml
_doctype component Layout

Container {
  style: "max-width-4xl mx-auto p-6"
  
  Header {
    style: "mb-8"
    text: "My Application"
  }
  
  Main {
    style: "flex-1"
    children: props.content
  }
}
```

### Styling

JAWT uses utility-based styling similar to Tailwind CSS:

```jml
Card {
  style: "bg-white rounded-lg shadow-md p-4 mb-4"
  
  Title {
    style: "text-xl font-bold text-gray-800"
    text: "Card Title"
  }
  
  Content {
    style: "text-gray-600 mt-2"
    text: "Card content goes here"
  }
}
```

## Project Structure

```
my-app/
├── components/
│   ├── layout.jml
│   └── card.jml
├── pages/
│   ├── index.jml
│   └── about.jml
├── assets/
│   └── styles.css
└── jawt.config.json
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `jawt init <name>` | Create a new JAWT project |
| `jawt run` | Start development server with hot reload |
| `jawt build` | Build optimised production bundle |
| `jawt serve` | Serve production build locally |

## Configuration

Customise your build process with `jawt.config.json`:

```json
{
  "port": 6500,
  "output": "dist",
  "publicPath": "/",
  "minify": true
}
```

## Features

- **Declarative Syntax**: Express your UI structure clearly and concisely
- **Component System**: Build reusable components with props and composition
- **Hot Reload**: See changes instantly during development
- **Minimal Bundle**: Generates lightweight, optimised output
- **Zero Configuration**: Works out of the box with sensible defaults