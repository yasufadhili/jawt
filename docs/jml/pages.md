# JML Pages Overview

**Pages in JML represent complete web pages that compile to standalone HTML documents.** They serve as the entry points for an application's routes and define the fundamental structure and metadata for each page in a site.

## Page Structure

Every JML page follows a strict structural pattern with a single root component that contains all page content. This design encourages clean architecture and consistent layouts across an application.

### Basic Page Anatomy

```jml
_doctype page home

import component Layout from "components/layout"

Page {
    title: "Welcome to My App"
    description: "A minimal web application built with JAWT"
    
    Layout {
        content: "Hello, World!"
    }
}
```

## Single Child Component Pattern

**Pages can only have one direct child component.** This architectural constraint promotes the use of layout components that manage the overall page structure, leading to more maintainable and reusable code.

### Recommended Pattern: Layout Components

The most effective approach is to create dedicated layout components that handle the page's visual structure:

```jml
_doctype page about

import component MainLayout from "components/main-layout"

Page {
    title: "About Us"
    description: "Learn more about our company and mission"
    
    MainLayout {
        section: "about"
        showSidebar: true
    }
}
```

This pattern allows you to:
- Maintain consistent layouts across multiple pages
- Centralise navigation and common UI elements
- Easily modify page structure without touching individual pages
- Support different layout variants for different page types

## Page Properties

Pages support several built-in properties for controlling HTML document metadata and behaviour:

### Essential Properties

```jml
Page {
    title: "Page Title"           // Sets <title> tag
    description: "Page summary"   // Sets meta description
    favicon: "/favicon.ico"       // Sets favicon link
    name: "home"                  // Internal page identifier
    
    // Your single child component
    Layout {}
}
```

### Additional Metadata Properties

```jml
Page {
    title: "E-commerce Store"
    description: "Shop the latest products online"
    favicon: "/assets/favicon.png"
    name: "shop"
    
    // SEO and social media properties
    keywords: "shopping, online, products"
    author: "Your Company Name"
    viewport: "width=device-width, initial-scale=1.0"
    
    // Open Graph properties for social sharing
    ogTitle: "Amazing Products Online"
    ogDescription: "Discover our collection"
    ogImage: "/assets/og-image.jpg"
    
    Layout {
        pageType: "store"
    }
}
```

## Import System for Pages

Pages can only import components, not modules. This restriction maintains the clear separation between presentation logic (components) and computational logic (modules).

### Component Imports

```jml
_doctype page dashboard

// Import from global components directory
import component AdminLayout from "components/admin-layout"
import component Header from "components/header"

// Import from same directory (page-specific components)
import component DashboardStats from "stats"

Page {
    title: "Admin Dashboard"
    
    AdminLayout {
        header: Header
        stats: DashboardStats
    }
}
```

### Import Rules

- **Global Components**: Import from `components/` directory using full path
- **Local Components**: Import from same directory using filename only
- **No Module Imports**: Pages cannot directly import modules; this must be done through components
- **Component Resolution**: Build system resolves all component dependencies before page compilation

## Dynamic Routing and Parameters

While the Build System handles route detection and parameter extraction, pages can receive dynamic parameters that are passed to their child components.

### Route Parameter Handling

For a dynamic route like `app/user/[id].jml`:

```jml
_doctype page userProfile

import component UserLayout from "components/user-layout"

Page {
    title: "User Profile"
    description: "View user profile information"
    
    UserLayout {
        // Parameters from route are automatically available
        // and can be passed to child components
        userId: params.id
        section: "profile"
    }
}
```

### Multiple Parameters

For routes like `app/blog/[category]/[slug].jml`:

```jml
_doctype page blogPost

import component BlogLayout from "components/blog-layout"
import component PostContent from "post-content"

Page {
    title: "Blog Post"
    description: "Read our latest blog content"
    
    BlogLayout {
        category: params.category
        slug: params.slug
        
        PostContent {
            categoryId: params.category
            postSlug: params.slug
        }
    }
}
```

## File Organisation

Pages must be organised in the `app/` directory following specific naming conventions:

### Directory Structure
```
app/
├── index.jml              # Root route (/)
├── about/
│   └── index.jml         # /about route
├── blog/
│   ├── index.jml         # /blog route
│   └── [slug].jml        # /blog/:slug dynamic route
└── user/
    ├── index.jml         # /user route
    ├── [id].jml          # /user/:id dynamic route
    └── [id]/
        └── settings.jml  # /user/:id/settings nested route
```

### Naming Rules

- **Static Routes**: Use `index.jml` within directories
- **Dynamic Routes**: Use `[parameter].jml` notation
- **Nested Routes**: Create subdirectories for route nesting
- **Root Page**: `app/index.jml` serves as the application root

## Page Lifecycle

Pages go through a defined compilation process:

1. **Route Detection**: Build system identifies page files and extracts route patterns
2. **Dependency Resolution**: All imported components are resolved and compiled first
3. **Parameter Extraction**: Dynamic route parameters are identified and made available
4. **HTML Generation**: Page compiles to complete HTML document with metadata
5. **Asset Integration**: Stylesheets, scripts, and components are integrated into final output

## Best Practices

### Layout Component Strategy

Create layout components that handle common page patterns:

```jml
// components/standard-layout.jml
_doctype component StandardLayout

Container {
    style: "min-h-screen bg-gray-50"
    
    Header {
        title: props.pageTitle
        navigation: props.showNav
    }
    
    Main {
        style: "container mx-auto px-4 py-8"
        content: props.children
    }
    
    Footer {
        style: "bg-gray-800 text-white p-4"
    }
}
```

### Consistent Metadata

Establish consistent patterns for page metadata across your application:

```jml
Page {
    title: `${props.pageTitle} | Your App Name`
    description: props.description || "Default app description"
    favicon: "/assets/favicon.svg"
    
    Layout {
        // Pass through relevant properties
    }
}
```

### Parameter Validation

Consider parameter validation within your layout components:

```jml
// In your layout component
Container {
    // Validate parameters before use
    content: params.id ? UserProfile : ErrorPage
}
```

## Limitations and Considerations

- **Single Child Constraint**: Pages can only have one direct child component
- **No Direct Module Access**: Pages cannot import modules directly
- **Route Handling**: Dynamic routing logic is handled by the Build System, not in JML
- **Static Generation**: Pages compile to static HTML; dynamic behaviour comes from components
- **Metadata Scope**: Page properties only affect the HTML document head, not component behaviour

## Integration with Build System

The Build System handles:
- **Route Registration**: Automatically registers routes based on file structure
- **Parameter Extraction**: Extracts route parameters and makes them available to pages
- **Dependency Management**: Ensures all imported components are available
- **HTML Generation**: Combines page metadata with component output
- **Asset Bundling**: Integrates stylesheets and JavaScript as needed

Pages focus purely on structure and metadata, while the Build System manages the complex routing and compilation orchestration.

## Related Documentation

- **[Components](components.md)** - Understanding component development for page content
- **[Build System](../architecture/build-system.md)** - How routing and compilation work together
- **[JML Overview](index.md)** - General JML language concepts
- **[Project Structure](../jawt/project-structure.md)** - Organising your JAWT application