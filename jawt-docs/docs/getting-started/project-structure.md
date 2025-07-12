# How a JAWT Project is Laid Out

Knowing your way around a JAWT project is key. I've set up a simple, conventional structure so you don't have to waste time making decisions about where to put files. Everything has a place, and once you know the layout, you can get to the fun part: building.

## The Basic Blueprint

When you run `jawt init my-project`, you get this starting structure:

```
my-project/
├── app/                    # Your pages and routes
│   └── index.jml          # The main page (the "/" route)
├── components/            # Reusable UI bits and pieces
├── assets/                # Images, fonts, and other static stuff
├── app.json              # Info about your project
├── jawt.config.json      # Configuration for the build process
└── dist/                 # Where the final, compiled app goes
```

Let's break down what each of these does.

## The `app/` Directory - Your Application's Pages

The `app/` directory is where the pages of your site live. The structure of this directory defines your app's routes.

### File-Based Routing

JAWT uses the file system to define routes. It's simple and intuitive:

```
app/
├── index.jml              # → / (the root page)
├── about/
│   └── index.jml         # → /about
├── blog/
│   ├── index.jml         # → /blog
│   └── [slug].jml        # → /blog/a-cool-post (a dynamic route)
└── user/
    ├── [id].jml          # → /user/123 (another dynamic route)
    └── [id]/
        └── settings.jml  # → /user/123/settings
```

### Page File Rules

Every page file needs to:

1.  Be named `index.jml` for a static route, or use `[brackets]` for a dynamic one.
2.  Start with `_doctype page <name>`.
3.  Have a single `Page` component as its root element.

### Example Page

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

Use square brackets for parts of the URL that can change.

```jml
# app/blog/[slug].jml
_doctype page blogPost

import component BlogLayout from "components/blog-layout"

Page {
    title: "Blog Post"
    description: "Read our latest blog content"
    
    BlogLayout {
        slug: params.slug  // You can access the dynamic part like this
    }
}
```

## The `components/` Directory - Reusable Building Blocks

Think of components as your own custom HTML tags. You build them once and can reuse them anywhere. The `components/` directory is where they all live.

### Organising Components

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
└── forms/
    ├── contact-form.jml
    └── login-form.jml
```

### Component File Rules

Each component file must:

1.  Start with `_doctype component <ComponentName>`.
2.  Export a single component.
3.  Use `PascalCase` for the component name.

```jml
# components/ui/button.jml
_doctype component Button

Button {
    text: props.text || "Click me"
    style: `${props.style || ""} px-4 py-2 rounded bg-blue-500 text-white hover:bg-blue-600`
    onClick: props.onClick
}
```

### Importing Components

You can import components in a few different ways:

```jml
# Import from the global components directory
import component Header from "components/layout/header"

# Import from the same directory (for page-specific components)
import component Card from "card"

# Import with an alias
import component MainButton from "components/ui/button"
```

## The `assets/` Directory - Static Files

This is where you put all your non-code files: images, fonts, CSS, etc.

### Organizing Assets

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

You can reference assets from anywhere in your project.

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

### `app.json` - Project Info

This file holds basic info about your project.

```json
{
  "name": "my-jawt-app",
  "version": "1.0.0",
  "description": "My awesome JAWT application",
  "author": "Your Name"
}
```

### `jawt.config.json` - Build Settings

This file controls how JAWT builds and serves your project.

```json
{
  "build": {
    "outDir": "dist",
    "assetsDir": "assets",
    "minify": true
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
    "assets": "assets"
  }
}
```

## The `dist/` Directory - The Final Product

After you run `jawt build`, this directory will contain your compiled app, ready to be deployed.

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
└── manifest.json          # A map of the build
```

## Best Practices

### 1. Group by Feature

Organise your components by what they do, not what they are.

```
components/
├── auth/           # Auth-related components
│   ├── login-form.jml
│   └── signup-form.jml
├── dashboard/      # Dashboard-specific components
│   ├── stats-card.jml
│   └── chart.jml
└── common/         # Shared components used everywhere
    ├── button.jml
    └── modal.jml
```

### 2. Name Things Clearly

Use descriptive names for your files.

```
# Good
components/user/profile-card.jml
components/navigation/main-menu.jml

# Not so good
components/card.jml
components/menu.jml
```

### 3. Don't Go Too Deep

Try to keep your directory structure relatively flat (3-4 levels max).

```
# Good
components/forms/contact/contact-form.jml

# A bit much
components/ui/forms/contact/complex/contact-form.jml
```

## The Workflow

### Starting Out

```bash
# Go to your project directory
cd my-project

# Start the dev server
jawt run

# In another terminal, you can run the debugger (optional)
jawt debug
```

### Building for Production

```bash
# Create an optimized build
jawt build

# Serve the production build locally to test it
jawt serve
```

## Common Questions

### Import Not Working?

**Problem**: Component not found.

**Solution**: Double-check your import paths. Make sure they're relative to the correct directory.

```jml
# Correct
import component Button from "components/ui/button"

# Incorrect (usually)
import component Button from "ui/button"
```

### Page Not Showing Up?

**Problem**: A page isn't accessible at its URL.

**Solution**: Check the file name and directory structure in `app/`.

```
# For the /about route, this is correct:
app/about/index.jml

# This is wrong:
app/about.jml
```

## What's Next?

Now that you know how a JAWT project is structured, you can:

1.  **Organise your own projects** with this structure in mind.
2.  **Plan out new projects** before you start coding.
3.  **Create your own project templates** for different types of apps.

That's it! You're ready to start building well-organised, scalable apps with JAWT.