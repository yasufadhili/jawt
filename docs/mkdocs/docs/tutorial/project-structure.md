# JAWT Project Structure

Understanding your JAWT project structure is key to learning how to organise a well-designed project—everything has its place, and knowing where to find (and put) things makes development much more efficient. This tutorial will guide you through the anatomy of a JAWT project and explain how each directory and file contributes to your application.

## Overview

A JAWT project follows a convention-over-configuration approach, meaning that by organising your files in the expected structure, everything works seamlessly without additional setup. Think of it as a well-organised workshop where every tool has its designated spot.

## Basic Project Structure

When you run `jawt init my-project`, you get this foundation:

```
my-project/
├── app/                    # Pages and routing
│   └── index.jml          # Root page (/)
├── components/            # Reusable UI components
├── modules/               # Modules for complex computations
├── assets/                # Static files (images, fonts, etc.)
├── app.json              # Project metadata
├── jawt.config.json      # Build configuration
└── dist/                 # Build output (created after jawt build)
```

Let's explore each part in detail.

## The `app/` Directory - Your Application Pages

The `app/` directory is where your application's pages live. It defines how one navigates through your site.

### Routing Convention

JAWT uses file-based routing, where the directory structure directly maps to URL routes:

```
app/
├── index.jml              # → / (root page)
├── about/
│   └── index.jml         # → /about
├── blog/
│   ├── index.jml         # → /blog
│   └── [slug].jml        # → /blog/hello-world (dynamic route)
├── contact/
│   └── index.jml         # → /contact
└── user/
    ├── index.jml         # → /user
    ├── [id].jml          # → /user/123 (dynamic route)
    └── [id]/
        └── settings.jml  # → /user/123/settings
```

### Page File Requirements

Each page file must:

1. Be named `index.jml` (for static routes) or use brackets for dynamic routes like `[id].jml`
2. Start with `_doctype page <name>`
3. Contain a single `Page` component as the root

### Example Page Structure

```jml
# app/about/index.jml
_doctype page about

import component Layout from "components/layout"

Page {
    title: "About Us"
    description: "Learn more about our company"
    
    Layout {
        section: "about"
    }
}
```

### Dynamic Routes

Dynamic routes use square brackets to indicate parameters:

```jml
# app/blog/[slug].jml
_doctype page blogPost

import component BlogLayout from "components/blog-layout"

Page {
    title: "Blog Post"
    description: "Read our latest blog content"
    
    BlogLayout {
        slug: params.slug  # Access the dynamic parameter
    }
}
```

## The `components/` Directory - Reusable UI Building Blocks

Think of components as prefabricated modules in construction—they're built once and can be used throughout your project. The `components/` directory houses all your reusable UI elements.

### Component Organisation

```
components/
├── layout/
│   ├── header.jml
│   ├── footer.jml
│   └── main-layout.jml
├── ui/
│   ├── button.jml
│   ├── card.jml
│   └── modal.jml
├── forms/
│   ├── contact-form.jml
│   └── login-form.jml
└── navigation/
    ├── navbar.jml
    └── sidebar.jml
```

### Component File Structure

Each component file must:

1. Start with `_doctype component <ComponentName>`
2. Export a single component
3. Use PascalCase for component names

```jml
# components/ui/button.jml
_doctype component Button

Button {
    text: props.text || "Click me"
    style: `${props.style || ""} px-4 py-2 rounded bg-blue-500 text-white hover:bg-blue-600`
    onClick: props.onClick
}
```

### Component Import Patterns

Components can be imported in different ways:

```jml
# Import from components directory (global)
import component Header from "components/layout/header"

# Import from same directory (local)
import component Card from "card"

# Import with alias
import component MainButton from "components/ui/button"
```

## The `modules/` Directory - Performance-Critical Code

Modules handle computationally intensive tasks that need to run at near-native speed. These compile to WebAssembly (WASM).

### Module Organisation

```
modules/
├── math/
│   ├── calculations.jml
│   └── geometry.jml
├── image/
│   └── processing.jml
└── data/
    ├── sorting.jml
    └── filtering.jml
```

### Module File Structure

```jml
# modules/math/calculations.jml
_doctype module calculations

export function fibonacci(n: number): number {
    if (n <= 1) return n;
    return fibonacci(n - 1) + fibonacci(n - 2);
}

export function isPrime(n: number): boolean {
    if (n < 2) return false;
    for (let i = 2; i <= Math.sqrt(n); i++) {
        if (n % i === 0) return false;
    }
    return true;
}
```

### Using Modules in Components

```jml
# components/calculator.jml
_doctype component Calculator

import browser // for Alert

import module math from "modules/math/calculations"

Container {
    style: "p-4"
    
    Button {
        text: "Calculate Fibonacci"
        onClick: () => {
            const result = math.fibonacci(10);
            browser.Alert(`Fibonacci(10) = ${result}`);
        }
    }
}
```

## The `assets/` Directory - Static Resources

The assets directory is like a storage room for all your non-code files—images, fonts, stylesheets, and other static resources.

### Asset Organisation

```
assets/
├── images/
│   ├── logo.svg
│   ├── hero-bg.jpg
│   └── icons/
│       ├── home.svg
│       └── user.svg
├── fonts/
│   ├── custom-font.woff2
│   └── icons.ttf
├── styles/
│   └── custom.css
└── data/
    ├── config.json
    └── content.json
```

### Using Assets

Assets can be referenced from anywhere in your project:

```jml
# In a page or component
Container {
    style: "bg-cover bg-center min-h-screen"
    backgroundImage: "url('/assets/images/hero-bg.jpg')"
    
    Image {
        src: "/assets/images/logo.svg"
        alt: "Company Logo"
        style: "w-32 h-32"
    }
}
```

## Configuration Files

### `app.json` - Project Metadata

This file contains basic information about your project:

```json
{
  "name": "my-jawt-app",
  "version": "1.0.0",
  "description": "My awesome JAWT application",
  "author": "Your Name"
}
```

### `jawt.config.json` - Build Configuration

This file controls how your project is built and served:

```json
{
  "build": {
    "outDir": "dist",
    "assetsDir": "assets",
    "minify": true,
  },
  "server": {
    "port": 6500,
    "host": "localhost",
    "https": false,
    "open": true
  },
  "paths": {
    "pages": "app",
    "components": "components",
    "modules": "modules",
    "assets": "assets"
  }
}
```

## The `dist/` Directory - Build Output

After running `jawt build`, the `dist/` directory contains your compiled application:

```
dist/
├── index.html              # Compiled root page
├── about/
│   └── index.html         # Compiled about page
├── assets/
│   ├── js/
│   │   ├── components.js  # Compiled components
│   │   └── modules.wasm   # Compiled modules
│   ├── css/
│   │   └── styles.css     # Compiled styles
│   └── images/
│       └── ...            # Optimised images
└── manifest.json          # Build manifest
```

## Project Structure Best Practices

### 1. Logical Grouping

Organise components by functionality rather than type:

```
components/
├── auth/           # Authentication-related components
│   ├── login-form.jml
│   └── signup-form.jml
├── dashboard/      # Dashboard-specific components
│   ├── stats-card.jml
│   └── chart.jml
└── common/         # Shared components
    ├── button.jml
    └── modal.jml
```

### 2. Consistent Naming

Use clear, descriptive names:

```
# Good
components/user/profile-card.jml
components/navigation/main-menu.jml

# Avoid
components/card.jml
components/menu.jml
```

### 3. Depth Considerations

Keep directory nesting reasonable (3-4 levels maximum):

```
# Good
components/forms/contact/contact-form.jml

# Too deep
components/ui/forms/contact/complex/contact-form.jml
```

## Complete Example Project Structure

Here's what a real-world JAWT project might look like:

```
my-blog/
├── app/
│   ├── index.jml                    # Home page
│   ├── about/
│   │   └── index.jml               # About page
│   ├── blog/
│   │   ├── index.jml               # Blog listing
│   │   └── [slug].jml              # Individual post
│   └── contact/
│       └── index.jml               # Contact page
├── components/
│   ├── layout/
│   │   ├── header.jml              # Site header
│   │   ├── footer.jml              # Site footer
│   │   └── main-layout.jml         # Main layout wrapper
│   ├── blog/
│   │   ├── post-card.jml           # Blog post preview
│   │   ├── post-content.jml        # Full post display
│   │   └── post-list.jml           # List of posts
│   ├── ui/
│   │   ├── button.jml              # Reusable button
│   │   ├── card.jml                # Card component
│   │   └── modal.jml               # Modal dialog
│   └── forms/
│       └── contact-form.jml        # Contact form
├── modules/
│   ├── content/
│   │   └── markdown-parser.jml     # Markdown processing
│   └── utils/
│       └── date-formatter.jml      # Date utilities
├── assets/
│   ├── images/
│   │   ├── logo.svg
│   │   └── blog/
│   │       ├── post1-hero.jpg
│   │       └── post2-hero.jpg
│   ├── fonts/
│   │   └── custom-font.woff2
│   └── data/
│       └── blog-posts.json
├── app.json
├── jawt.config.json
└── dist/                           # Generated after build
```

## Common Patterns and Conventions

### Page-Specific Components

Sometimes you need components that are only used by a single page:

```
app/
├── dashboard/
│   ├── index.jml                   # Dashboard page
│   ├── stats-widget.jml            # Page-specific component
│   └── chart-container.jml         # Page-specific component
```

Import page-specific components using relative paths:

```jml
# In app/dashboard/index.jml
import component StatsWidget from "stats-widget"
import component ChartContainer from "chart-container"
```

### Shared Layouts

Create layout components for consistent page structure:

```jml
# components/layout/standard-layout.jml
_doctype component StandardLayout

Container {
    style: "min-h-screen bg-gray-50"
    
    Header {
        title: props.title
    }
    
    Main {
        style: "container mx-auto px-4 py-8"
        content: props.children
    }
    
    Footer {}
}
```

### Configuration-Based Routing

Use configuration files to manage complex routing:

```json
# assets/data/routes.json
{
  "routes": [
    {
      "path": "/",
      "component": "home",
      "title": "Home"
    },
    {
      "path": "/blog/:slug",
      "component": "blog-post",
      "title": "Blog Post"
    }
  ]
}
```

## Development Workflow

### Starting Development

```bash
# Navigate to project directory
cd my-project

# Start development server
jawt run

# In another terminal, start debugger (optional)
jawt debug
```

### Building for Production

```bash
# Build optimised version
jawt build

# Serve production build locally for testing
jawt serve
```

### Project Maintenance

```bash
# Check project structure
ls -la app/ components/ modules/

# Validate configuration
cat jawt.config.json

# Clean build artifacts
rm -rf dist/
```

## Troubleshooting Common Issues

### Import Resolution Problems

**Problem**: Component not found error

**Solution**: Check import paths and ensure components are in the correct directory

```jml
# Correct
import component Button from "components/ui/button"

# Incorrect
import component Button from "ui/button"
```

### Routing Issues

**Problem**: Page not accessible

**Solution**: Verify file naming and directory structure

```
# Correct for /about route
app/about/index.jml

# Incorrect
app/about.jml
```

### Build Configuration Problems

**Problem**: Assets not found after build
**Solution**: Check `jawt.config.json` paths configuration

```json
{
  "paths": {
    "assets": "assets"  // Ensure this matches your directory
  }
}
```

## Next Steps

Now that you understand JAWT project structure, you can:

1. **Organise Existing Projects**: Restructure your current projects following these conventions
2. **Plan New Projects**: Design your directory structure before you start coding
3. **Explore Advanced Features**: Look into custom build configurations and optimisations
4. **Create Templates**: Build project templates for common patterns you use

## Key Takeaways

1. **Convention Over Configuration**: Following JAWT's expected structure makes everything work seamlessly
2. **Logical Organisation**: Group related files together for easier maintenance
3. **Clear Separation**: Pages, components, and modules each have their designated purpose and location
4. **Scalable Structure**: The pattern works for both small projects and large applications
5. **Import Flexibility**: Use both absolute and relative imports as appropriate

You can now create well-organised, maintainable applications that scale as your projects grow.