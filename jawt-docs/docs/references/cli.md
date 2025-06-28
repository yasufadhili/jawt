# JAWT CLI Reference

The JAWT CLI provides essential commands for creating, developing, and building JAWT applications.

## Installation

Ensure JAWT is [installed](../getting-started/installation.md) and available in your system PATH before using these commands.

## Quick Start

```bash
# Create a new project
jawt init my-app

# Navigate to project directory
cd my-app

# Start development server
jawt run

# Build for production
jawt build
```

## Commands

### init

Creates a new JAWT project with the default structure and configuration files.

#### Synopsis

```bash
jawt init <project-name>
```

#### Arguments

| Argument | Description | Required |
|----------|-------------|----------|
| `<project-name>` | Name of the new project directory | Yes |

#### Examples

```bash
# Create a new project called "my-app"
jawt init my-app

# Create a project with a more descriptive name
jawt init portfolio-website
```

#### Generated Structure

```
project-name/
├── app/
│   └── index.jml          # Main application entry point
├── components/            # Reusable components directory
├── assets/               # Static assets (images, fonts, etc.)
├── app.json             # Application configuration
└── jawt.config.json     # JAWT toolchain configuration
```

---

### run

Starts the development server with hot reload functionality. Monitors JML files for changes and automatically reloads the browser.

#### Synopsis

```bash
jawt run [options]
```

#### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-p <port>` | Specify custom port | 6500 |
| `-c` | Run with cleared cache | - |

#### Prerequisites

- `app.json` must exist in the current directory
- `jawt.config.json` must exist in the current directory

#### Examples

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

#### Output

The development server will be available at `http://localhost:<port>` with:

- Hot reload enabled for JML files
- Automatic browser refresh on changes
- Real-time error reporting

---

### build

Compiles your JAWT application into production-ready web standard files. Generates optimised HTML, CSS, and JavaScript output.

#### Synopsis

```bash
jawt build [options]
```

#### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-o <directory>` | Specify custom output directory | `dist` |

#### Examples

```bash
# Build to default dist directory
jawt build

# Build to custom output directory
jawt build -o public

# Build to nested directory
jawt build -o build/production
```

#### Output Files

The build process generates:
- **HTML**: Compiled pages from JML definitions
- **CSS**: Optimised stylesheets with unused styles removed
- **JavaScript**: Minified component bundles
- **WASM**: Compiled modules (when applicable)
- **Assets**: Processed static files (images, fonts, etc.)

---

### serve

Serves the production build locally for previewing how your application will behave in production environment.

#### Synopsis

```bash
jawt serve
```

#### Prerequisites

- Must have a built application (run `jawt build` first)
- Built files must exist in the output directory

#### Example Workflow

```bash
# Build your application
jawt build

# Serve the production build
jawt serve
```

#### Status

⚠️ **Currently Unavailable** - This command is planned for future releases.

---

### debug

Starts the JAWT debugger, providing debugging tools and insights into your application's compilation and runtime behaviour.

#### Synopsis

```bash
jawt debug [options]
```

#### Options

| Option | Description | Default |
|--------|-------------|---------|
| `-p <port>` | Specify custom port | 6501 |

#### Examples

```bash
# Start debugger on default port (6501)
jawt debug

# Start debugger on custom port
jawt debug -p 9000
```

#### Features

The debugger interface provides:
- Component hierarchy inspection
- JML syntax error highlighting
- Build process insights
- Performance metrics

#### Status

⚠️ **Currently Unavailable** - This command is planned for future releases.

---

## Global Options

These options work with most JAWT commands:

| Option | Description |
|--------|-------------|
| `--help` | Display help information for the command |
| `--version` | Display JAWT version information |

#### Examples

```bash
jawt --version
jawt run --help
jawt build --help
```

## Configuration Files

### app.json

Application-specific configuration including metadata, routing, and build settings.

### jawt.config.json

JAWT toolchain configuration for compilation, development server, and build optimisations.

## Project Structure

A typical JAWT project follows this structure:

```
my-project/
├── app/                   # Application pages
│   ├── index.jml         # Main entry point
│   └── about.jml         # Additional pages
├── components/           # Reusable components
│   ├── Header.jml
│   └── Footer.jml
├── assets/              # Static assets
│   ├── images/
│   ├── fonts/
│   └── styles/
├── dist/                # Build output (generated)
├── app.json            # App configuration
└── jawt.config.json    # JAWT configuration
```

## Common Workflows

### Starting a New Project

```bash
# Create and set up new project
jawt init my-project
cd my-project

# Start development
jawt run
```

### Development Workflow

```bash
# Start development server
jawt run

# In another terminal, start debugger (when available)
jawt debug
```

### Production Deployment

```bash
# Build for production
jawt build

# Preview production build locally (when available)
jawt serve

# Deploy contents of dist/ directory to your hosting provider
```

### Working with Custom Ports

```bash
# If default ports are in use
jawt run -p 3000
jawt debug -p 3001
```

## Troubleshooting

### Common Issues

#### Command not found: jawt

**Problem**: JAWT CLI is not recognised by your shell.

**Solutions**:

- Ensure JAWT is properly [installed](../getting-started/installation.md)
- Verify JAWT is in your system PATH
- Restart your terminal after installation

#### Missing app.json or jawt.config.json

**Problem**: Required configuration files are missing.

**Solutions**:

- Create a new project using `jawt init`
- Manually create the required configuration files
- Check that you're running commands from the correct directory

#### Port already in use

**Problem**: Default ports (6500, 6501) are occupied by other processes.

**Solutions**:
```bash
# Use custom ports
jawt run -p 8080
jawt debug -p 8081

# Find and stop processes using default ports
lsof -ti:6500 | xargs kill -9  # macOS/Linux
netstat -ano | findstr :6500   # Windows
```

#### Build fails

**Problem**: Compilation errors during the build process.

**Solutions**:

- Check JML syntax for errors
- Ensure all imported components exist
- Verify file paths and dependencies
- Run `jawt debug` for detailed error information (when available)

#### Hot reload not working

**Problem**: Changes aren't reflected in the browser automatically.

**Solutions**:

- Ensure you're editing files within the project directory
- Check the browser console for connection errors
- Try running with cleared cache: `jawt run -c`
- Manually refresh the browser

### Getting Help

For additional support:

- Check the [JAWT documentation](../index.md)
- Review [JML](../jml/index.md) syntax guides
- Report issues on the [project repository](https://github.com/yasufadhili/jawt)

---

**Version**: Early Development  
**Last Updated**: Development Roadmap Phase 1