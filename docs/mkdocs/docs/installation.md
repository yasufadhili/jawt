# Installation & Setup

Getting JAWT up and running is straightforward. Choose your platform below and follow the simple installation steps.

## Quick Installation

### Linux & macOS (Soon)

Install JAWT with a single command that automatically downloads and configures everything:

```bash
curl -fsSL https://raw.githubusercontent.com/yasufadhili/jawt/main/install.sh | sudo bash
```


This script will:

- Detect your system architecture automatically
- Download the latest JAWT release for your platform
- Install JAWT to the appropriate system directory
- Configure your PATH environment variable
- Set up all necessary permissions

### Windows (Soon)

Download and run the MSI installer:

**[Download JAWT for Windows](#)**

The installer will:

- Install JAWT to `Program Files`
- Add JAWT to your system PATH
- Create Start Menu shortcuts
- Set up file associations for `.jml` files

## Verify Installation

Once installed, verify JAWT is working correctly by checking the version:

```bash
jawt --version
```

You should see output similar to:
```
JAWT v0.1.0
Platform: linux/amd64
```

## System Requirements

JAWT has minimal system requirements and works on:

### Supported Platforms
- **Linux**: x64, ARM64 (Ubuntu 18.04+, CentOS 7+, Debian 9+)
- **macOS**: x64, Apple Silicon (macOS 10.15+)
- **Windows**: x64, ARM64 (Windows 10+, Windows Server 2019+)

### Dependencies
- **No external dependencies required** - JAWT is a single binary
- **Modern browser** for development server (Firefox, Edge, Chrome, Safari)
- **Text editor** of your choice (VS Code, Vim, Emacs, etc.)

## Troubleshooting

### Permission Issues (Linux/macOS)

If you encounter permission errors, ensure you're running the install script with `sudo`:

```bash
curl -fsSL https://raw.githubusercontent.com/yasufadhili/jawt/main/install.sh | sudo bash
```

### Command Not Found

If `jawt --version` returns "command not found":

**Linux/macOS**: Restart your terminal or run:
```bash
source ~/.bashrc  # or ~/.zshrc if using zsh
```

**Windows**: Restart Command Prompt or PowerShell

### Firewall/Antivirus Warnings

Some antivirus software may flag the installer. JAWT is safe to install - you can:
- Add an exception for the JAWT installer
- Download directly from GitHub releases if needed
- Build from source if you prefer (see Contributing guide)

## Manual Installation

If you prefer manual installation or need a specific version:

1. Visit the [JAWT releases page](https://github.com/yasufadhili/jawt/releases)
2. Download the appropriate archive for your platform
3. Extract the `jawt` binary to a directory in your PATH
4. Make the binary executable (Linux/macOS): `chmod +x jawt`

## Development Environment Setup

JAWT currently has no support in any major text editor or IDE. You can still just write JML code in any text editor
## Configuration

JAWT works with zero configuration, but you can customise behaviour:

### Global Configuration
JAWT looks for configuration in:

- `~/.jawt/config.json` (Linux/macOS)

- `%APPDATA%\jawt\config.json` (Windows)

### Project Configuration
Each project can have its own `jawt.config.json` file for project-specific settings.

## Updating JAWT

### Automatic Updates
JAWT can update itself:
```bash
jawt update
```

### Manual Updates
Re-run the installation script to get the latest version:

**Linux/macOS**:
```bash
curl -fsSL https://raw.githubusercontent.com/yasufadhili/jawt/main/install.sh | sudo bash
```

**Windows**: Download and run the latest MSI installer

## Uninstalling JAWT

### Linux/macOS
```bash
sudo rm /usr/local/bin/jawt
sudo rm -rf ~/.jawt
```

### Windows
Use "Add or Remove Programs" in Windows Settings, or run:
```powershell
jawt uninstall
```

## Next Steps

Now that JAWT is installed, you're ready to start building:

- **[Create Your First Project](../tutorial/first-project.md)** - Build a simple JAWT application from scratch
- **[Project Structure](project-structure.md)** - Understand how JAWT projects are organised
- **[JML Quick Start](../jml/quick-start.md)** - Learn the basics of JML syntax
- **[CLI Reference](../cli/index.md)** - Explore all available JAWT commands