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
- [Architecture](https://yasufadhili.github.io/jawt/architecture/) — How Jawt handles pages, components, modules and more
- [JML](https://yasufadhili.github.io/jawt/jml/) — Jawt's DSL for app development
- [Jawt CLI](https://yasufadhili.github.io/jawt/references/cli) — Useful commands when working with Jawt
- [Building](BUILDING.MD) this source code

## Components

**This project repo consists of several components that make up JAWT**

- [Compiler](internal/compiler) : Compiles Pages, Modules and Components
- [Build System](internal/build) : Orchestrates the entire build process
- [Development Server](internal/devserver) : Serves Jawt projects locally during development

## Licensing

This project is primarily licensed under the [MIT Licence](LICENSE). 

It also includes code from the new [TypeScript compiler](https://github.com/microsoft/typesctipt-go),
which is licensed under the [Apache License 2.0](./LICENSE-APACHE). See the `internal/tsc/` directory for details.

