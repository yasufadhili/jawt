# JAWT—Just Another Web Tool

**JAWT is a toolchain that enables building minimal web applications using a declarative approach.** 

Write your application structure and behaviour using **JML**, a domain-specific language for JAWT.

## Features

- **Declarative Syntax**: Express your UI structure clearly and concisely
- **Component System**: Build reusable components with props and composition
- **Hot Reload**: See changes instantly during development
- **Minimal Bundle**: Generates lightweight, optimized output
- **Zero Configuration**: Works out of the box with sensible defaults

## Resources

**Useful links to understand the Toolchain**

- [Getting Started](docs/tutorial) — Quick tutorial on Jawt
- [Documentation](docs/jawt) — Build and use Jawt apps
- [Architecture](docs/architecture) — How Jawt handles pages, components, modules and more
- [JML](docs/jml) — Jawt's DSL for app development
- [Jawt CLI](CLI.MD) — Useful commands when working with Jawt
- [Building](BUILDING.MD) this source code

## Components

**This project repo consists of several components that make up JAWT**

- [Page Compiler](internal/pc) : Compiles Pages
- [Component Compiler](internal/cc) : Compiles Components
- [Module Compiler](internal/mc) : Compiles Modules
- [Build System](internal/build) : Orchestrates the entire build process
- [Development Server](internal/server) : Serves Jawt projects locally during development

## License

All parts of the Jawt Toolchain are licensed under the [MIT Licence](LICENSE).
