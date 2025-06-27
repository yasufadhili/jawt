# JAWT—Just Another Web Tool

**JAWT is a toolchain that enables building minimal web applications using a declarative approach.** 

Write your application structure and behaviour using **JML**, a domain-specific language for JAWT.

>**NOTE:** Very early development

## Features

- **Declarative Syntax**: Express your UI structure clearly and concisely
- **Component System**: Build reusable components with props and composition
- **Hot Reload**: See changes instantly during development
- **Minimal Bundle**: Generates lightweight, optimized output
- **Zero Configuration**: Works out of the box with sensible defaults

## Resources

**Useful links to understand the Toolchain**

- [Getting Started](https://yasufadhili.github.io/jawt/) — Quick tutorial on Jawt
- [Documentation](https://yasufadhili.github.io/jawt/) — Build and use Jawt apps
- [Architecture](https://yasufadhili.github.io/jawt/) — How Jawt handles pages, components, modules and more
- [JML](https://yasufadhili.github.io/jawt/) — Jawt's DSL for app development
- [Jawt CLI](CLI.MD) — Useful commands when working with Jawt
- [Building](BUILDING.MD) this source code

## Components

**This project repo consists of several components that make up JAWT**

- [Compiler](internal/compiler) : Compiles Pages, Modules and Components
- [Build System](internal/build) : Orchestrates the entire build process
- [Development Server](internal/server) : Serves Jawt projects locally during development

## License

All parts of the Jawt Toolchain are licensed under the [MIT Licence](LICENSE).
