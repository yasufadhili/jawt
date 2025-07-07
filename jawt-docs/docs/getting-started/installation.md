# Getting Set Up

Alright, let's get JAWT installed. It's pretty painless.

## The Quick Way

### Linux (Debian/Ubuntu)

I made a script that handles everything. Just pop this in your terminal:

```bash
curl -fsSL https://raw.githubusercontent.com/yasufadhili/jawt/main/install.sh | sudo bash
```

This little command will:

-   Figure out your system's architecture.
-   Grab the latest JAWT release.
-   Stick it in the right place.
-   Set up your PATH.
-   Handle all the permissions.

### Windows (Coming Soon)

I'll have an MSI installer for Windows soon. It'll be a simple download-and-click affair.

## Did it Work?

To make sure everything is installed correctly, just ask JAWT for its version:

```bash
jawt --version
```

You should see something like this:

```
JAWT v0.1.0
Platform: linux/amd64
```

## What You'll Need

JAWT is pretty low-maintenance.

### Supported Systems
-   **Linux**: x64, ARM64 (Ubuntu 18.04+, CentOS 7+, Debian 9+)
-   **macOS**: x64, Apple Silicon (macOS 10.15+)
-   **Windows**: x64, ARM64 (Windows 10+, Windows Server 2019+)

### Dependencies
-   **None.** JAWT is a single binary. No need to install anything else.
-   A modern browser (like Firefox, Edge, Chrome, or Safari).
-   Your favourite text editor.

## If Something Went Wrong

### Permission Problems (Linux/macOS)

If you get permission errors, you probably forgot to run the install script with `sudo`.

```bash
curl -fsSL https://raw.githubusercontent.com/yasufadhili/jawt/main/install.sh | sudo bash
```

### "Command Not Found"

If your terminal says `jawt: command not found`:

**Linux/macOS**: Try restarting your terminal or running:
```bash
source ~/.bashrc  # or ~/.zshrc if you're a zsh person
```

**Windows**: Restart your Command Prompt or PowerShell.

### Firewall/Antivirus Warnings

Some antivirus programs can be a bit jumpy. JAWT is safe, I promise. You can add an exception for it or, if you're feeling adventurous, build it from the source yourself.

## The Manual Way

If you'd rather do it yourself:

1.  Go to the [JAWT releases page](https://github.com/yasufadhili/jawt/releases).
2.  Download the right archive for your system.
3.  Extract the `jawt` binary and put it somewhere in your PATH.
4.  Make it executable (on Linux/macOS): `chmod +x jawt`

## Editor Setup

Right now, there's no fancy syntax highlighting or IntelliSense for JML in any major editor. It's on the to-do list! For now, you can just write it as plain text.

## Configuration

JAWT is designed to work without any configuration. But if you want to customise things:

### Global Config
JAWT looks for a config file here:

-   `~/.jawt/config.json` (Linux/macOS)
-   `%APPDATA%\jawt\config.json` (Windows)

### Project Config
Each project can have its own `jawt.config.json` for project-specific settings.

## Updating JAWT

### Automatic Updates

I've added a command to let JAWT update itself:
```bash
jawt update
```

### Manual Updates

Just run the installation script again to get the latest version.

**Linux/macOS**:
```bash
curl -fsSL https://raw.githubusercontent.com/yasufadhili/jawt/main/install.sh | sudo bash
```

**Windows**: Download and run the new MSI installer when it's ready.

## Getting Rid of JAWT

### Linux/macOS
```bash
sudo rm /usr/local/bin/jawt
sudo rm -rf ~/.jawt
```

### Windows

Use "Add or Remove Programs" or run:
```powershell
jawt uninstall
```

## What's Next?

Now that you're all set up, you're ready to build something.

-   **[Create Your First Project](../tutorial/first-page.md)** - Let's build a simple JAWT app.
-   **[Project Structure](../getting-started/project-structure.md)** - A look at how JAWT projects are organized.
-   **[JML Quick Start](../jml/index.md)** - The basics of the JML language.
-   **[CLI Reference](../references/cli.md)** - All the commands you can use.
