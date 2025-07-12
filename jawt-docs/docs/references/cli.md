# The JAWT CLI

The `jawt` command-line interface is your main tool for creating, running, and building JAWT applications.

## Getting Started

```bash
# Create a new project
jawt init my-app

# cd into the new directory
cd my-app

# Start the dev server
jawt dev

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

### `dev`

Starts the local development server with hot module replacement (HMR). It watches for changes to your JML files and automatically refreshes your browser.

#### Usage

```bash
jawt dev [options]
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
jawt dev

# Start on a different port
jawt dev -p 3000

# Start with a clean cache
jawt dev -c
```

---

### `build`

Compiles your app or library into a production-ready build. It optimises everything and spits out standard HTML, CSS, and JavaScript.

#### Usage

```bash
jawt build [options]
```

#### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-o <directory>` | Specify a custom output directory. | `dist` |
| `--as-lib` | Compile the project as a library for reuse. | `false` |

#### Examples

```bash
# Build to the default `dist` directory
jawt build

# Build to a `public` directory instead
jawt build -o public

# Build as a library
jawt build --as-lib
```

#### What it Generates

-   **HTML**: Your compiled JML pages.
-   **CSS**: Optimised and minified stylesheets.
-   **JavaScript**: Minified component bundles.
-   **WASM**: Compiled modules (if you have any).
-   **Assets**: Your static files, processed and optimized.

---

### `create page`

Scaffolds a new JML page with a basic structure.

#### Usage

```bash
jawt create page <name>
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `<name>` | The name of the new page. | Yes |

#### Example

```bash
jawt create page about
```

---

### `create component`

Scaffolds a new JML component with a basic structure.

#### Usage

```bash
jawt create component <name>
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `<name>` | The name of the new component. | Yes |

#### Example

```bash
jawt create component MyButton
```

---

### `add`

Adds a JML component library from a local path or remote repository.

#### Usage

```bash
jawt add <path/repo>
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `<path/repo>` | The path to the local library or its repository URL. | Yes |

#### Example

```bash
jawt add ./my-local-lib
jawt add https://github.com/user/some-jawt-lib
```

---

### `install`

Installs npm logic packages for use within Jawt projects.

#### Usage

```bash
jawt install <pkg>
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `<pkg>` | The name of the npm package to install. | Yes |

#### Example

```bash
jawt install lodash
```

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
jawt dev --help
jawt build --help
```
