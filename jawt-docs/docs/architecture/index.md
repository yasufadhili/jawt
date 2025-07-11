# Jawt – Just Another Web Tool

**Jawt is a playful, opinionated, zero-config toolchain for building modern web applications using a custom declarative language called JML.** It’s built to scratch the itch of building web UIs in a clear, composable, and developer-controlled way — with **no setup files, no bundler chaos, and no boilerplate jungles**.

> It’s my own curated frontend universe, compiled entirely through Go, designed for real-world usage and personal scalability.

---

## ✨ Philosophy

* **Developer-focused**: Jawt is built for a single developer (you) with strong preferences.
* **Declarative-first**: You write *what* you want, Jawt figures out *how* to do it.
* **Zero-config**: Everything just works. You never touch Webpack, TSConfig, or PostCSS again.
* **Highly opinionated**: One way (or two at most) to do things.
* **Self-contained**: Apps are portable. No `.env`, `node_modules`, or configs littering the userland.

---

## What Jawt Is

* A **toolchain** for building modern Single Page Applications (SPAs).
* A **compiler** that transforms `.jml` declarative files into Web Components written in Lit.
* A **runtime** environment with built-in router, state, and system APIs.
* A **library system** for reusing logic and UI components across projects.
* A **self-hosted dev server** with live reload and full integration.
* A **CLI** that handles everything from scaffolding to packaging.

---

## What Jawt Is Not

* Not a general-purpose build tool like Webpack or Vite.
* Not a framework like React/Vue — it builds **on Web Components** via **Lit**.
* Not aimed at supporting every use case — it supports *my* chosen path.
* Not exposing internal tooling — Node, Tailwind, etc. are hidden behind the toolchain.

---

## Core Toolchain & Technologies

| Layer | Tool |
|---|---|
| **Parser** | ANTLR |
| **Compiler & CLI** | Go |
| **Styling** | TailwindCSS |
| **Logic & Transpilation** | TypeScript |
| **Component Runtime** | Lit |
| **Routing** | Vaadin Router |
| **State Management** | NanoStores |
| **Bundling** | esbuild (optional) |
| **Dev Server + Watcher** | Custom Go-based live server with HMR |

---

## Notable Design Constraints

* No mixing of frameworks — no React/Vue support.
* All routing and UI flows are based on Web Components.
* You own the UX philosophy: consistency is preferred over flexibility.
* Everything builds down to HTML + JS — even routing is declarative.

---

## Final Thought

**Jawt is my own full-stack SPA engine**, shaped exactly how you want it. It's fast, expressive, local-first, and declarative at the core. It respects the web, uses its standards (TypeScript, Web Components), but gives you a language and system that’s tailored to my creativity — not the trends.
