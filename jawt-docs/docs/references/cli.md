# The JAWT CLI

The `jawt` command-line interface is your main tool for creating, running, and building JAWT applications.

## Getting Started

```bash
# Create a new project
jawt init my-app

# cd into the new directory
cd my-app

# Start the dev server
jawt run

# Build the app for production
jawt build
```

## Commands

### `init`

Kicks off a new JAWT project with a default file structure and config files.

#### Usage

```bash
jawt init <project-name>
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `<project-name>` | The name of the new project directory. | Yes |

#### Example

```bash
jawt init my-awesome-app
```

#### What it Creates

```
my-awesome-app/
├── app/
│   └── index.jml          # Your main page
├── components/            # For reusable components
├── assets/               # For static files
├── app.json             # App config
└── jawt.config.json     # JAWT config
```

---

### `run`

Fires up the development server, which comes with hot reloading. It watches for changes to your JML files and automatically refreshes your browser.

#### Usage

```bash
jawt run [options]
```

#### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-p <port>` | Use a custom port. | 6500 |
| `-c` | Start with a clean cache. | - |

#### Prerequisites

-   You need an `app.json` and `jawt.config.json` in your project's root directory.

#### Examples

```bash
# Start the server on the default port (6500)
jawt run

# Start on a different port
jawt run -p 3000

# Start with a clean cache
jawt run -c
```

---

### `build`

Compiles your app into a production-ready build. It optimises everything and spits out standard HTML, CSS, and JavaScript.

#### Usage

```bash
jawt build [options]
```

#### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-o <directory>` | Specify a custom output directory. | `dist` |

#### Examples

```bash
# Build to the default `dist` directory
jawt build

# Build to a `public` directory instead
jawt build -o public
```

#### What it Generates

-   **HTML**: Your compiled JML pages.
-   **CSS**: Optimised and minified stylesheets.
-   **JavaScript**: Minified component bundles.
-   **WASM**: Compiled modules (if you have any).
-   **Assets**: Your static files, processed and optimized.

---

### `serve`

Serves your production build locally. This is useful for checking how your app will behave in a production environment.

#### Usage

```bash
jawt serve
```

#### Prerequisites

-   You need to have a production build. Run `jawt build` first.

#### Example

```bash
# Build your app
jawt build

# Serve the production build
jawt serve
```

#### Status

⚠️ **Coming soon!** This command isn't implemented yet.

---

### `debug`

Starts the JAWT debugger, which gives you tools to inspect your app's compilation and runtime behavior.

#### Usage

```bash
jawt debug [options]
```

#### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-p <port>` | Use a custom port. | 6501 |

#### Examples

```bash
# Start the debugger on the default port (6501)
jawt debug

# Start on a different port
jawt debug -p 9000
```

#### Features

-   Inspect the component tree.
-   Get detailed error messages.
-   See build process insights.
-   Check performance metrics.

#### Status

⚠️ **Coming soon!** This command isn't implemented yet.

---

## Global Options

These options work with most commands.

| Option | Description |
|--------|-------------|
| `--help` | Show help for a command. |
| `--version` | Show the JAWT version. |

#### Examples

```bash
jawt --version
jawt run --help
jawt build --help
```