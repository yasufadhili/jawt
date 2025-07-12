# JAWT Roadmap

This document outlines the complete development journey for JAWT, from its foundational concepts to a robust, production-capable toolchain. This roadmap is designed for a personal project, focusing on a clear, sequential build-out of features.

## Core Philosophy

JAWT is being built with the following principles at its heart:
*   **Developer-centric**: To create a highly efficient and enjoyable development experience for myself and any collaborators.
*   **Declarative-first**: To define *what* the application should be, allowing JAWT to handle the underlying *how*.
*   **Zero-configuration**: To eliminate the need for complex setup and boilerplate.
*   **Opinionated design**: To provide clear, effective, and consistent ways of doing things.
*   **Self-contained portability**: To ensure applications are easy to manage, share, and deploy without external dependencies.

## Phased Development: From Concept to Production

The roadmap is structured into distinct phases, each building upon the last, detailing the progression from initial development to a fully-fledged toolchain.

### Phase 1: Foundational Development & Core Tooling (Initial Build)

This phase focuses on establishing the absolute core of JAWT: the language, the compiler, and the basic development environment.

#### 1.1 JML Language Definition & Parser
*   **JML Syntax Specification**: Formal definition of the JML language syntax.
*   **ANTLR Grammar Development**: Creation of the ANTLR grammar for JML.
*   **JML Parser Implementation**: Building the parser to transform JML source into an Abstract Syntax Tree (AST).

#### 1.2 Core Compiler & Code Generation
*   **Go-based Compiler Core**: Initial development of the Go compiler to process the JML AST.
*   **Basic Web Component Generation (Lit)**: Ability to compile simple JML components into Lit-based Web Components.
*   **TailwindCSS Integration**: Initial setup for processing TailwindCSS classes from JML `style` attributes.
*   **TypeScript Transpilation Integration**: Setting up the pipeline for TypeScript logic within JML to be transpiled.

#### 1.3 Command Line Interface (CLI) Basics
*   **`jawt init`**: Command to scaffold a new, empty JAWT project with basic directory structure and configuration files (`app.json`, `jawt.config.json`).
*   **`jawt dev`**: Command to start a basic development server with file watching and rudimentary browser refresh.
*   **`jawt build`**: Command to perform a basic compilation of JML into static HTML/JS assets.

#### 1.4 Initial Runtime Components
*   **Basic Built-in Components**: Implementation of fundamental Lit-based components like `Container`, `Text`, `Button`.
*   **Basic Routing**: Initial implementation of file-based routing for JML pages.

### Phase 2: Enhanced Developer Experience & Core Functionality

Building upon the foundation, this phase focuses on improving the development workflow and implementing essential features for practical application development.

#### 2.1 Advanced JML & Compiler Features
*   **JML Language Expansion**: Implementation of `for-loops`, conditional expressions, and type-checked `props` within JML.
*   **Improved Compiler Diagnostics**: More precise and actionable error messages for JML and TypeScript compilation issues.
*   **Golden File Testing**: Implementing comprehensive golden file tests for the compiler to ensure consistent and correct output.

#### 2.2 Robust Development Server
*   **Hot Module Replacement (HMR)**: Full implementation of HMR for instant feedback on code changes without full page reloads.
*   **Browser Error Overlay**: Developing a clear, informative error overlay for the development server.
*   **CLI Command Completion**: Implementing auto-completion for `jawt` commands in common shells.

#### 2.3 Runtime & API Expansion
*   **Comprehensive Built-in Components**: Expanding the default Lit-based component library with more common UI elements (e.g., advanced form controls, data display components).
*   **Core Runtime APIs**:
    *   **`network` API**: Full implementation of `fetch()`, `getJSON()`, `post()` with robust error handling and common HTTP patterns.
    *   **`store` API**: Comprehensive `get()`, `set(), `clear()`, `observe()` for unified client-side persistence (using NanoStores).
    *   **`events` API**: Stable `emit()`, `on()`, `off()` for inter-component communication.
    *   **`clipboard` API**: Reliable `copy()`, `paste()` functionality.
    *   **`date` API**: Utility functions for `now()`, `format()`, `parse()`.
    *   **`env` API**: Secure and reliable `isProd()`, `getEnv()` for environment-specific logic.

#### 2.4 Project Types & Reusability
*   **Application Project Type**: Full support for building self-contained SPAs.
*   **Library Project Type**: Ability to define and export JML components or TypeScript scripts as reusable libraries.
*   **`jawt add`**: Command to integrate JML component libraries from local paths or repositories.
*   **`jawt install`**: Command to integrate npm logic packages for use within Jawt projects.

### Phase 3: Production Readiness & Advanced Capabilities

This phase focuses on features critical for deploying production-grade applications, emphasising performance, security, and advanced development patterns.

#### 3.1 Performance & Optimisation
*   **Advanced Bundling**:
    *   **Automatic Code Splitting**: Implementing intelligent code splitting per route/component for faster initial page loads.
    *   **Comprehensive Tree Shaking**: Ensuring only truly used code is included in the final bundles.
*   **Image Optimisation Pipeline**: Built-in capabilities for image compression, responsive image generation (e.g., WebP conversion, multiple sizes), and lazy loading.
*   **Font Optimisation**: Font subsetting and preloading strategies.
*   **Critical CSS Extraction**: Automatically extracting and inlining critical CSS for above-the-fold content.

#### 3.2 Security & Reliability
*   **Content Security Policy (CSP) Generation**: Automated or easily configurable CSP headers for enhanced security.
*   **XSS/CSRF Guidance & Primitives**: Providing clear patterns and potentially built-in primitives to mitigate common web vulnerabilities.
*   **Secrets Management Integration**: Secure handling of environment variables and sensitive data during build and runtime.
*   **Dependency Security Scanning**: Integrating with tools to scan third-party dependencies for known vulnerabilities.
*   **Comprehensive Error Handling & Reporting**: Beyond dev server overlays, a production-ready mechanism to catch and report unhandled client-side exceptions.
*   **Monitoring & Observability Hooks**: Providing integration points for APM tools and structured logging.

#### 3.3 Advanced Runtime Features
*   **Advanced State Management Patterns**: Guidance and potentially new primitives for managing global state, derived state, and state synchronisation for complex applications.
*   **Form Handling & Validation**: Robust form primitives with integrated validation rules (potentially via Zod or similar), error display, and submission handling.
*   **Internationalisation (i18n) & Localisation (l10n)**: A comprehensive system for managing translations, date/number formatting, and locale-specific content.
*   **PWA Capabilities**: Built-in Service Worker generation for caching assets, offline access, and manifest file generation.

#### 3.4 Tooling & Ecosystem Maturity
*   **Full LSP (Language Server Protocol) Implementation**: Delivering a complete LSP for JML, offering advanced autocompletion, type checking, refactoring, and contextual help within IDEs.
*   **`jawt debug`**: Implementing a comprehensive debugger for inspecting runtime behaviour, component state, and performance metrics.
*   **Automated Testing Framework**: Providing a clear, integrated path for unit, integration, and end-to-end testing of JAWT applications.
*   **Visual Regression Testing Integration**: Providing a framework for visual regression tests to prevent unintended UI changes.

### Phase 4: Ecosystem & Long-Term Vision

This phase looks towards expanding JAWT's capabilities for broader use cases and long-term maintainability.

#### 4.1 Deployment & CI/CD
*   **Automated CI/CD Integration**: Providing templates and guidance for integrating JAWT builds and tests into CI/CD pipelines.
*   **Deployment Recipes**: Clear, documented pathways for deploying JAWT applications to various hosting environments.
*   **Rollback Strategy**: Built-in support for quick and safe rollbacks of deployments.

#### 4.2 Advanced Rendering & Data
*   **Server-Side Rendering (SSR) / Static Site Generation (SSG)**: Exploring and implementing mechanisms for pre-rendering JML pages on the server or at build time for improved SEO and initial load performance (likely a Go-based rendering engine).
*   **Efficient Data Fetching Strategies**: Advanced caching mechanisms (e.g., stale-while-revalidate) and data pre-fetching for upcoming routes or components.

#### 4.3 Refinement & Future-Proofing
*   **Official Plugin System**: A formalised plugin system for extending JAWT's compiler, runtime, or CLI with custom functionalities.
*   **Advanced Route Metadata & Layouts**: Full support for route metadata (`meta` block) and advanced layout components.
*   **Page Transitions**: Built-in support for smooth page transitions via router hooks or dedicated components.
*   **DevTools Overlay**: A more comprehensive in-browser DevTools overlay for development.

## Non-Goals

To maintain focus and uphold JAWT's core philosophy, the following are explicitly considered non-goals:
*   **General-purpose build tool**: JAWT will not become a generic bundler like Webpack or Vite.
*   **Framework Agnosticism**: JAWT is built *on* Web Components via Lit; it will not support other frontend frameworks (e.g., React, Vue) directly.
*   **Unopinionated Configuration**: JAWT will remain highly opinionated, prioritising convention over configuration.