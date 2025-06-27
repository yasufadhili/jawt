# Your First JAWT Page

In this tutorial, you'll learn how to create your first web page using JML. Think of this as your "Hello, World!" moment with JAWT—we'll focus purely on page creation using built-in components, without any external dependencies or complex features.

## What You'll Learn

- How to set up a basic JAWT project
- Understanding JML page syntax and structure
- Using built-in components like `Container`, `Text`, and `Button`
- Applying styling with Tailwind CSS classes
- Creating a complete, functional web page

## Prerequisites

- JAWT installed on your system
- Basic understanding of HTML and CSS concepts
- Familiarity with Tailwind CSS is helpful but not required

## Setting Up Your Project

Let's start by creating a new JAWT project:

```bash
jawt init my-first-page
cd my-first-page
```

This creates a project structure like this:
```
my-first-page/
├── app/
│   └── index.jml          # Your main page file
├── components/            # For future components
├── assets/               # Static assets
├── app.json             # Project configuration
└── jawt.config.json     # Build configuration
```

## Understanding JML Page Structure

Every JML page follows a specific pattern. Let's examine the basic structure:

```jml
_doctype page pageName

Page {
    title: "Page Title"
    description: "Page description"
    
    // Single root component goes here
    Container {
        // Page content
    }
}
```

### Key Components of a JML Page

1. **Document Type Declaration**: `_doctype page pageName` - This tells JAWT that this file should compile to an HTML page
2. **Page Component**: The `Page` wrapper that contains metadata and your page content
3. **Single Root Component**: Pages can only have one direct child component (like React's JSX)

## Creating Your First Page

Let's replace the default content in `app/index.jml` with a simple welcome page:

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
            text: "This is my first page built with JML"
            style: "text-lg text-gray-700 mb-8 text-center"
        }
        
        Button {
            text: "Welcome Button"
            style: "bg-blue-500 hover:bg-blue-600 text-white px-6 py-2 rounded-lg shadow-md transition-colors"
        }
    }
}
```

## Understanding Built-in Components

JAWT provides several built-in components that you can use immediately:

### Container
The `Container` component is like a `<div>` in HTML—it groups other elements together.

```jml
Container {
    style: "flex flex-col space-y-4 p-6"
    
    // Child components go here
}
```

### Text
The `Text` component displays text content, similar to `<p>`, `<h1>`, etc. in HTML.

```jml
Text {
    text: "Your text here"
    style: "text-xl font-semibold text-gray-800"
}
```

### Button
The `Button` component creates clickable buttons.

```jml
Button {
    text: "Click me!"
    style: "bg-green-500 text-white px-4 py-2 rounded hover:bg-green-600"
}
```

## Styling with Tailwind CSS

JAWT uses Tailwind CSS for styling. You apply styles using the `style` property with Tailwind utility classes:

```jml
Container {
    style: "bg-white shadow-lg rounded-lg p-6 max-w-md mx-auto"
}
```

### Common Tailwind Patterns

Here are some useful Tailwind class combinations:

**Centering content:**
```jml
style: "flex items-center justify-center"
```

**Card-like appearance:**
```jml
style: "bg-white shadow-md rounded-lg p-6"
```

**Responsive spacing:**
```jml
style: "p-4 md:p-8 lg:p-12"
```

**Colour schemes:**
```jml
style: "bg-blue-500 text-white hover:bg-blue-600"
```

## Building a Complete Example

Let's create a more comprehensive page that showcases different components and styling:

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
            
            Container {
                style: "bg-white rounded-xl shadow-lg p-8"
                
                Text {
                    text: "Skills"
                    style: "text-2xl font-semibold text-gray-800 mb-6"
                }
                
                Container {
                    style: "grid grid-cols-2 md:grid-cols-3 gap-4"
                    
                    Container {
                        style: "bg-blue-100 text-blue-800 px-4 py-2 rounded-lg text-center"
                        
                        Text {
                            text: "JAWT"
                            style: "font-medium"
                        }
                    }
                    
                    Container {
                        style: "bg-green-100 text-green-800 px-4 py-2 rounded-lg text-center"
                        
                        Text {
                            text: "JavaScript"
                            style: "font-medium"
                        }
                    }
                    
                    Container {
                        style: "bg-purple-100 text-purple-800 px-4 py-2 rounded-lg text-center"
                        
                        Text {
                            text: "CSS"
                            style: "font-medium"
                        }
                    }
                    
                    Container {
                        style: "bg-red-100 text-red-800 px-4 py-2 rounded-lg text-center"
                        
                        Text {
                            text: "HTML"
                            style: "font-medium"
                        }
                    }
                    
                    Container {
                        style: "bg-yellow-100 text-yellow-800 px-4 py-2 rounded-lg text-center"
                        
                        Text {
                            text: "UI Design"
                            style: "font-medium"
                        }
                    }
                    
                    Container {
                        style: "bg-indigo-100 text-indigo-800 px-4 py-2 rounded-lg text-center"
                        
                        Text {
                            text: "Tailwind"
                            style: "font-medium"
                        }
                    }
                }
            }
        }
    }
}
```

## Running Your Page

To see your page in action:

```bash
# Start the development server
jawt run

# Your page will be available at http://localhost:6500
```

The development server includes hot reload, so any changes you make to your JML file will automatically update in the browser.

## Page Properties Reference

The `Page` component supports several properties for controlling the HTML document:

```jml
Page {
    title: "Page Title"                    // Sets <title> tag
    description: "Page description"         // Sets meta description
    favicon: "/favicon.ico"                 // Sets favicon
    name: "internal-page-name"             // Internal identifier
    keywords: "keyword1, keyword2"         // SEO keywords
    author: "Your Name"                    // Page author
    viewport: "width=device-width, initial-scale=1.0"  // Viewport settings
    
    // Your page content
    Container {
        // ...
    }
}
```

## Building for Production

When you're ready to deploy your page:

```bash
# Build optimised version
jawt build

# Your HTML file will be in the dist/ directory
```

This creates a `dist/` directory with your compiled page file, optimised and ready for deployment to any web server.

## Key Takeaways

1. **JML is Declarative**: You describe what you want, not how to build it
2. **Single Root Component**: Pages can only have one direct child (usually a `Container` or `Main`)
3. **Built-in Components**: `Container`, `Text`, and `Button` cover most basic needs
4. **Tailwind Styling**: Use the `style` property with Tailwind utility classes
5. **Hot Reload**: Changes appear instantly during development

## Next Steps

Now that you've created your first JAWT page, you might want to explore:

- Creating reusable components for more complex applications
- Adding interactive behaviour with JML scripting
- Building multi-page applications with routing
- Integrating WebAssembly modules for performance-critical features

## Troubleshooting

**Page not loading?**
- Check that your `_doctype page` declaration is correct
- Ensure your JML syntax is valid (proper braces, quotes)
- Look for error messages in the terminal running `jawt run`

**Styling not appearing?**
- Verify Tailwind class names are correct
- Check that `style` properties are in quotes
- Remember that Tailwind uses specific class names (e.g., `text-blue-600`, not `text-blue`)

**Changes not reflecting?**
- Save your file and wait a moment for hot reload
- Check the browser console for any errors
- Try stopping and restarting `jawt run`

Congratulations! You've just created your first JAWT page using JML. You've learned the fundamental concepts that will serve as the foundation for building more complex web applications with JAWT.