# JAWT CLI Commands

The `jawt` command-line interface provides essential tools for creating, developing, and building JAWT applications. All commands should be run from your project directory unless otherwise specified.

## Commands Overview

| Command | Description | Status      |
|---------|-------------|-------------|
| [`init`](#init) | Create a new JAWT project | ✅ Available |
| [`run`](#run) | Start development server with hot reload | ✅ Available |
| [`build`](#build) | Build optimised production bundle | ✅ Available |
| [`serve`](#serve) | Serve production build locally | Unavailable |
| [`debug`](#debug) | Start JAWT debugger | Unavailable |

---

## init

Creates a new JAWT project with the default structure and configuration files.

### Usage

```bash
jawt init <project-name>
```

### Arguments

- `<project-name>` - Name of the new project (required)

### Example

```bash
jawt init my-awesome-app
cd my-awesome-app
```

This creates a new directory with the following structure:
```
my-awesome-app/
├── app/
│   └── index.jml
├── components/
├── assets/
├── app.json
└── jawt.config.json
```

---

## run

Starts the development server with hot reload functionality. Monitors your JML files for changes and automatically reloads the browser.

### Usage

```bash
jawt run [options]
```

### Options

- `-p <port>` - Specify custom port (default: 6500)
- `-c` - Run with cleared cache

### Prerequisites

- `app.json` must exist in the current directory
- `jawt.config.json` must exist in the current directory

### Examples

```bash
# Start on default port (6500)
jawt run

# Start on custom port
jawt run -p 3000

# Start with cleared cache
jawt run -c

# Combine options
jawt run -p 8080 -c
```

The development server will be available at `http://localhost:<port>` with hot reload enabled.

---

## build

Compiles your JAWT application into production-ready web standard files. Generates optimised HTML, CSS, and JavaScript output.

### Usage

```bash
jawt build [options]
```

### Options

- `-o <directory>` - Specify custom output directory (default: `dist`)

### Examples

```bash
# Build to default dist directory
jawt build

# Build to custom output directory
jawt build -o public

# Build to nested directory
jawt build -o build/production
```

### Output

The build process creates web standard files including:
- Compiled HTML pages
- Optimised CSS stylesheets
- Minified JavaScript bundles
- Static assets (images, fonts, etc.)

---

## serve

Serves the production build locally for previewing how your application will behave in production. This command requires you to run `jawt build` first.

### Usage

```bash
jawt serve
```

### Prerequisites

- Must have a built application (run `jawt build` first)
- Built files must exist in the output directory

### Example Workflow

```bash
# Build your application
jawt build

# Serve the production build
jawt serve
```

The production preview will be available at a local URL (check console output for the exact address).

---

## debug

Starts the JAWT debugger, which provides debugging tools and insights into your application's compilation and runtime behaviour.

### Usage

```bash
jawt debug [options]
```

### Options

- `-p <port>` - Specify custom port (default: 6501)

### Examples

```bash
# Start debugger on default port (6501)
jawt debug

# Start debugger on custom port
jawt debug -p 9000
```

The debugger interface will be available in your browser at `http://localhost:<port>`.

---

## Global Options

These options work with most JAWT commands:

- `--help` - Display help information for the command
- `--version` - Display JAWT version information

### Examples

```bash
jawt --version
jawt run --help
jawt build --help
```

---

## Common Workflows

### Starting a New Project

```bash
# Create new project
jawt init my-project
cd my-project

# Start development
jawt run
```

### Development Workflow

```bash
# Start development server
jawt run

# In another terminal, start debugger if needed
jawt debug
```

### Production Deployment

```bash
# Build for production
jawt build

# Preview production build locally
jawt serve

# Deploy contents of dist/ directory to your hosting provider
```

---

## Troubleshooting

### Common Issues

**Command not found: jawt**
- Ensure JAWT is properly installed and in your PATH

**Missing app.json or jawt.config.json**
- These files are required for `jawt run`
- Create them using `jawt init` or manually

**Port already in use**
- Use `-p` option to specify a different port
- Check what's running on the default ports (6500, 6501)

**Build fails**
- Check your JML syntax for errors
- Ensure all imported components exist
- Run `jawt debug` for detailed error information