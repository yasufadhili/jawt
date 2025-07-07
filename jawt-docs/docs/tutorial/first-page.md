# Your First JAWT Page

This is your "Hello, World!" for JAWT. We'll create a simple web page using JML and some of the built-in components. No complex stuff, just the basics to get you started.

## What You'll Learn

-   How to create a new JAWT project.
-   The basic structure of a JML page.
-   How to use built-in components like `Container`, `Text`, and `Button`.
-   How to apply styles with Tailwind CSS.

## Prerequisites

-   You'll need JAWT installed.
-   A basic understanding of HTML and CSS will be helpful.

## Setting Up the Project

First, let's create a new project.

```bash
jawt init my-first-page
cd my-first-page
```

This will give you a basic project structure:
```
my-first-page/
├── app/
│   └── index.jml          # This is where we'll work
├── components/            # For later
├── assets/               # For images, etc.
├── app.json             # Project config
└── jawt.config.json     # Build config
```

## The Anatomy of a JML Page

Every JML page has a simple, predictable structure.

```jml
_doctype page pageName

Page {
    title: "Page Title"
    description: "A description of the page"
    
    // You can only have one component directly inside a Page
    Container {
        // The rest of your page content goes here
    }
}
```

### The Key Parts

1.  **`_doctype page pageName`**: This tells JAWT to compile this file into an HTML page.
2.  **`Page`**: A special component that holds metadata and the content of your page.
3.  **A Single Root Component**: You can only have one component directly inside `Page`. Usually, this is a `Container`.

## Let's Build Something

Open up `app/index.jml` and replace its content with this:

```jml
_doctype page welcome

Page {
    title: "Welcome to JAWT"
    description: "My first JAWT page using JML"
    
    Container {
        style: "min-h-screen bg-gray-50 flex flex-col items-center justify-center p-8"
        
        Text {
            text: "Hello, JAWT!"
            style: "text-4xl font-bold text-blue-600 mb-4"
        }
        
        Text {
            text: "This is my first page built with JML."
            style: "text-lg text-gray-700 mb-8 text-center"
        }
        
        Button {
            text: "Welcome Button"
            style: "bg-blue-500 hover:bg-blue-600 text-white px-6 py-2 rounded-lg shadow-md transition-colors"
        }
    }
}
```

## The Built-in Components

JAWT gives you a few basic components to start with.

### `Container`

This is your basic `<div>`. It's for grouping other elements.

```jml
Container {
    style: "flex flex-col space-y-4 p-6"
    
    // Other components go here
}
```

### `Text`

This is for displaying text. Think of it as a `<p>` or `<h1>` tag.

```jml
Text {
    text: "Some text here."
    style: "text-xl font-semibold text-gray-800"
}
```

### `Button`

This creates a clickable button.

```jml
Button {
    text: "Click me!"
    style: "bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
}
```

## Styling with Tailwind CSS

JAWT is set up to use Tailwind CSS for styling. You just add a `style` property with Tailwind's utility classes.

```jml
Container {
    style: "bg-white shadow-lg rounded-lg p-6 max-w-md mx-auto"
}
```

## A More Complete Example

Let's make a slightly more interesting page.

```jml
_doctype page portfolio

Page {
    title: "My Portfolio"
    description: "A simple portfolio page built with JAWT"
    favicon: "/favicon.ico"
    
    Container {
        style: "min-h-screen bg-gradient-to-br from-blue-50 to-indigo-100"
        
        Container {
            style: "max-w-4xl mx-auto p-8"
            
            Text {
                text: "John Doe"
                style: "text-5xl font-bold text-gray-800 text-center mb-2"
            }
            
            Text {
                text: "Web Developer & Designer"
                style: "text-xl text-gray-600 text-center mb-8"
            }
            
            Container {
                style: "bg-white rounded-xl shadow-lg p-8 mb-8"
                
                Text {
                    text: "About Me"
                    style: "text-2xl font-semibold text-gray-800 mb-4"
                }
                
                Text {
                    text: "I'm passionate about creating beautiful, functional web experiences using modern tools and technologies. JAWT allows me to build fast, efficient websites with clean, declarative code."
                    style: "text-gray-700 leading-relaxed mb-6"
                }
                
                Button {
                    text: "View My Work"
                    style: "bg-blue-500 hover:bg-blue-600 text-white px-6 py-3 rounded-lg font-medium transition-colors shadow-md"
                }
            }
        }
    }
}
```

## See It in Action

To see your page, run the dev server:

```bash
# Start the dev server
jawt run
```

Your page will be running at `http://localhost:6500`. The server has hot reloading, so any changes you make to the JML file will show up in the browser instantly.

## Page Properties

The `Page` component has a few properties you can set:

```jml
Page {
    title: "Page Title"                    // The <title> tag
    description: "Page description"         // The meta description
    favicon: "/favicon.ico"                 // The favicon
    name: "internal-page-name"             // An internal name for the page
    keywords: "keyword1, keyword2"         // SEO keywords
    author: "Your Name"                    // The author of the page
    viewport: "width=device-width, initial-scale=1.0"  // Viewport settings
    
    // Your page content
    Container {
        // ...
    }
}
```

## Building for Production

When you're ready to deploy your page, just run:

```bash
# Create an optimized build
jawt build
```

This will create a `dist/` directory with your compiled page, ready to be uploaded to any web server.

## Key Takeaways

1.  **JML is declarative**: You describe the "what," not the "how."
2.  **Pages have a single root component**: Usually a `Container`.
3.  **Built-in components are your friends**: `Container`, `Text`, and `Button` are great starting points.
4.  **Styling is done with Tailwind**: Use the `style` property.
5.  **Hot reload is awesome**: See your changes instantly.

## What's Next?

Now that you've built a basic page, you can explore:

-   Creating your own reusable components.
-   Adding interactivity with TypeScript.
-   Building multi-page apps with routing.

Congrats! You've built your first page with JAWT. You now know the fundamentals for building much more complex applications.