# JAWT - Just Another Web Tool

So, what's JAWT? Honestly, it's a personal scratch to an itch. I want to build web apps without getting bogged down in a swamp of configuration files, boilerplate, and a dozen different libraries just to get a "Hello, World" on screen.

JAWT is my take on a simpler, more direct way to build things for the web. It's a toolchain built around a declarative language called **JML**. The whole idea is to let you describe *what* you want your app to look like and do, and let the tool handle the messy parts.

> **Heads up!** This is still very much a work-in-progress. Things might change, break, or not work as expected.

## What's the big idea?

Jawt is built on a few core philosophies:

*   **Developer-focused**: Designed for a single developer with strong preferences.
*   **Declarative-first**: You write *what* you want, Jawt figures out *how* to do it.
*   **Zero-config**: Everything just works out of the box.
*   **Highly opinionated**: One (or two) clear ways to do things.
*   **Self-contained**: Apps are portable with no external dependencies cluttering your project.

It's a **toolchain** for building modern Single Page Applications (SPAs), a **compiler** for JML into Web Components (via Lit), a **runtime** with built-in routing and state, and a **CLI** that handles everything from scaffolding to packaging.

Jawt is **not** a general-purpose build tool, nor is it a framework like React or Vue; it builds *on* Web Components. It's tailored for a specific, efficient workflow.

## Where to go from here?

**Useful links to get started:**

- [Getting Started](https://yasufadhili.github.io/jawt/) — A quick tutorial to get your feet wet.
- [Documentation](https://yasufadhili.github.io/jawt/) — The full guide to building and using Jawt apps.
- [Architecture](https://yasufadhili.github.io/jawt/architecture/) — A deep dive into Jawt's design and internal workings.
- [JML](https://yasufadhili.github.io/jawt/jml/) — Learn the language for building Jawt apps.
- [Jawt CLI](https://yasufadhili.github.io/jawt/references/cli) — All the commands you'll need.
- [Building from source](BUILDING.MD)

## License

All parts of the Jawt Toolchain are licensed under the [MIT Licence](LICENSE).
