# JAWT Development Roadmap

This roadmap outlines the development phases and features for JAWT (Just A Weird Web Tool), a toolchain for building minimal web applications using declarative JML syntax. 
The focus is on simplicity, clarity, and functionality for creating lightweight web apps.

## Foundation Phase ✅

- [x] **Core infrastructure established** to support the JAWT ecosystem.
- [x] **CLI commands** (`init`, `run`, `build`) implemented for project setup, development, and production builds.

---

## Phase 1: Core Compilation System [ ]

- [ ] **Page Compiler** (`internal/pc`):
    - Compiles JML page definitions into standard HTML documents.
    - Handles page metadata (title, description, etc.).
    - Manages page routing and structure.
    - Integrates with component system for page composition.
- [ ] **Component Compiler** (`internal/cc`):
    - Transforms JML component definitions into JavaScript components.
    - Supports component props, composition, and lifecycle management.
    - Uses Tailwind CSS utility classes for styling.
    - Optimises component rendering for performance.
- [ ] **Module Compiler** (`internal/mc`):
    - Compiles JML modules into WebAssembly (WASM) modules.
    - Handles module imports, exports, and dependencies.
    - Provides seamless JavaScript-WASM interoperability.
    - Optimises runtime performance for WASM modules.
- [ ] **Styling System**:
    - Generates Tailwind CSS utility classes for pages and components.
    - Processes custom styling declarations in JML.
    - Optimises CSS output by removing unused styles.
    - Supports responsive design patterns.

---

## Phase 2: Build Infrastructure [ ]

- [ ] **Build System** (`internal/build`):
    - Orchestrates compilation of HTML, CSS, JavaScript, and WASM.
    - Manages asset processing (images, fonts, etc.).
    - Generates production-ready output with Tailwind CSS purging.
- [ ] **Development Server** (`internal/server`):
    - Serves JAWT projects locally during development.
    - Provides hot reload for JML, CSS, and JavaScript changes.
    - Supports file watching and change detection.

---

## Phase 3: Developer Experience [ ]

- [ ] **Browser Debugger**:
    - Visual debugging interface in the browser.
    - Component hierarchy inspection.
    - JML syntax error highlighting.
    - WASM module debugging and profiling.
    - Tailwind CSS class usage analysis.
- [ ] **Language Server Protocol (LSP)**:
    - JML syntax highlighting and validation.
    - IntelliSense and auto-completion for JML.
    - Error diagnostics and quick fixes.
    - Support for go-to-definition and refactoring.
- [ ] **Editor Support**:
    - Visual Studio Code extension for JML.
    - Basic support for Vim/Neovim.
    - Optional plugins for other editors (e.g., Sublime Text).

---

## Phase 4: Advanced Features [ ]

- [ ] **Enhanced JML Features**:
    - Conditional rendering syntax for dynamic UI.
    - Loop and iteration constructs for lists.
    - Event handling declarations in JML.
    - Animation and transition support.
    - Form validation syntax.
- [ ] **Performance Optimisations**:
    - Tree-shaking for unused code.
    - Code splitting for lazy loading.
    - Critical CSS extraction for faster rendering.
    - Bundle size analysis tools.
- [ ] **Testing Framework**:
    - Unit testing for components.
    - Integration testing for pages.
    - Basic accessibility testing.

---

## Phase 5: Ecosystem Expansion [ ]

- [ ] **Plugin System**:
    - Simple architecture for custom plugins.
    - Hooks for extending the build process.
    - Support for community contributions.
- [ ] **Component Library**:
    - Basic reusable component library.
    - Documentation for common components.
    - Playground for testing components.
- [ ] **Deployment Tools**:
    - Simplified static site generation.
    - Basic Progressive Web App (PWA) features.

---


## Milestones

- **v0.1.0 - Foundation Release** ✅:
    - Core CLI functionality (`init`, `run`, `build`).
    - Full HTML, CSS, JavaScript, and WASM compilation pipeline.
    - Development server with hot reload.
- **v0.2.0 - Developer Experience** [ ]:
    - Browser debugger and basic LSP support.
    - Visual Studio Code extension for JML.
- **v1.0.0 - Production Ready** [ ]:
    - Stable JML syntax and API.
    - Comprehensive editor support.
    - Basic testing and performance tools.

---
