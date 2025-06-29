#!/bin/bash

# Jawt Installation Script
# Supports full installation with Node.js, TypeScript, and TailwindCSS or executable-only installation
# Automatically detects system architecture, manages dependencies, and configures environment
# Includes version checking to avoid unnecessary re-downloads

set -euo pipefail

REPO="yasufadhili/jawt"
API_URL="https://api.github.com/repos/${REPO}/releases/latest"
NODE_API_URL="https://nodejs.org/dist/index.json"

JAWT_ROOT="/usr/local/jawt"
JAWT_BIN_DIR="${JAWT_ROOT}/bin"
JAWT_NODE_DIR="${JAWT_ROOT}/node"
JAWT_TSC_DIR="${JAWT_ROOT}/tsc"
JAWT_TAILWIND_DIR="${JAWT_ROOT}/tailwind"
SYSTEM_BIN_DIR="/usr/local/bin"

BINARY_NAME="jawt"
TMP_DIR=$(mktemp -d)

# Ensure temporary directory is cleaned up on exit
trap 'rm -rf "$TMP_DIR"' EXIT

# Colours for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Colour

print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
}

error_exit() {
    print_error "$1"
    exit 1
}

download_with_progress() {
    local url="$1"
    local output="$2"
    local description="$3"

    print_info "Downloading ${description}..."
    if command -v wget >/dev/null 2>&1; then
        wget --progress=bar:force:noscroll -O "$output" "$url" 2>&1 | \
        while IFS= read -r line; do
            if [[ $line =~ [0-9]+% ]]; then
                echo "..."
            fi
        done
        echo
    elif command -v curl >/dev/null 2>&1; then
        curl -# -fL -o "$output" "$url"
    else
        error_exit "Neither curl nor wget is available for downloading"
    fi
}

extract_with_progress() {
    local archive="$1"
    local destination="$2"
    local description="$3"

    print_info "Extracting ${description}..."
    if [[ "$archive" == *.tar.xz ]]; then
        tar -xf "$archive" -C "$destination" --strip-components=1
    elif [[ "$archive" == *.tar.gz ]]; then
        tar -xzf "$archive" -C "$destination" --strip-components=1
    elif [[ "$archive" == *.zip ]]; then
        unzip -q "$archive" -d "$TMP_DIR/extract"
        # Move contents from extracted folder to destination
        mv "$TMP_DIR/extract"/*/* "$destination"
    else
        error_exit "Unsupported archive format: $archive"
    fi
}

detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$OS" in
        linux) OS_NAME="linux" ;;
        darwin) OS_NAME="macos" ;;
        *) error_exit "Unsupported operating system: $OS" ;;
    esac

    case "$ARCH" in
        x86_64) ARCH_NAME="amd64" ;;
        aarch64|arm64) ARCH_NAME="arm64" ;;
        *) error_exit "Unsupported architecture: $ARCH" ;;
    esac

    print_info "Detected platform: ${OS_NAME}/${ARCH_NAME}"
}

check_dependencies() {
    print_info "Checking dependencies..."

    local missing_deps=()

    if ! command -v curl >/dev/null 2>&1 && ! command -v wget >/dev/null 2>&1; then
        missing_deps+=("curl or wget")
    fi

    if ! command -v jq >/dev/null 2>&1; then
        missing_deps+=("jq")
    fi

    if ! command -v tar >/dev/null 2>&1; then
        missing_deps+=("tar")
    fi

    if [ ${#missing_deps[@]} -ne 0 ]; then
        error_exit "Missing required dependencies: ${missing_deps[*]}"
    fi

    print_success "All dependencies are available"
}

prompt_installation_type() {
    if [ -d "$JAWT_ROOT" ]; then
        print_warning "Jawt installation directory already exists at $JAWT_ROOT"
        echo
        echo "Installation options:"
        echo "1) Full installation (Jawt + Node.js + TypeScript + TailwindCSS)"
        echo "2) Jawt executable only"
        echo "3) Cancel installation"
        echo

        while true; do
            read -p "Choose installation type (1-3): " choice
            case $choice in
                1) INSTALL_TYPE="full"; break ;;
                2) INSTALL_TYPE="executable"; break ;;
                3) print_info "Installation cancelled"; exit 0 ;;
                *) print_error "Invalid choice. Please select 1, 2, or 3." ;;
            esac
        done
    else
        echo
        echo "Installation options:"
        echo "1) Full installation (Jawt + Node.js + TypeScript + TailwindCSS) [Recommended]"
        echo "2) Jawt executable only"
        echo

        while true; do
            read -p "Choose installation type (1-2): " choice
            case $choice in
                1) INSTALL_TYPE="full"; break ;;
                2) INSTALL_TYPE="executable"; break ;;
                *) print_error "Invalid choice. Please select 1 or 2." ;;
            esac
        done
    fi

    print_info "Selected installation type: $INSTALL_TYPE"
}

create_directories() {
    print_info "Creating directory structure..."

    sudo mkdir -p "$JAWT_BIN_DIR"

    if [ "$INSTALL_TYPE" = "full" ]; then
        sudo mkdir -p "$JAWT_NODE_DIR"
        sudo mkdir -p "$JAWT_TSC_DIR"
        sudo mkdir -p "$JAWT_TAILWIND_DIR"
    fi

    print_success "Directory structure created"
}

get_node_lts_version() {
    print_info "Fetching Node.js LTS version information..."

    NODE_VERSION=$(curl -fsSL "$NODE_API_URL" | jq -r '[.[] | select(.lts != false)] | .[0].version')

    if [ -z "$NODE_VERSION" ] || [ "$NODE_VERSION" = "null" ]; then
        error_exit "Failed to fetch Node.js LTS version"
    fi

    print_info "Latest Node.js LTS version: $NODE_VERSION"
}

check_existing_node_version() {
    if [ -x "$JAWT_NODE_DIR/bin/node" ]; then
        CURRENT_NODE_VERSION=$("$JAWT_NODE_DIR/bin/node" --version 2>/dev/null || echo "")
        if [ "$CURRENT_NODE_VERSION" = "$NODE_VERSION" ]; then
            print_success "Node.js $NODE_VERSION is already installed"
            return 0
        else
            print_info "Current Node.js version ($CURRENT_NODE_VERSION) differs from LTS ($NODE_VERSION)"
            return 1
        fi
    else
        print_info "No existing Node.js installation found"
        return 1
    fi
}

install_nodejs() {
    if [ "$INSTALL_TYPE" != "full" ]; then
        return 0
    fi

    get_node_lts_version

    # Check if we already have the correct version
    if check_existing_node_version; then
        print_info "Skipping Node.js download - already at LTS version"
        return 0
    fi

    # Determine Node.js architecture naming
    case "$ARCH_NAME" in
        amd64) NODE_ARCH="x64" ;;
        arm64) NODE_ARCH="arm64" ;;
    esac

    # Determine Node.js platform naming and archive format
    case "$OS_NAME" in
        linux)
            NODE_PLATFORM="linux"
            ARCHIVE_EXT="tar.xz"
            ;;
        macos)
            NODE_PLATFORM="darwin"
            ARCHIVE_EXT="tar.gz"
            ;;
    esac

    NODE_ARCHIVE="node-${NODE_VERSION}-${NODE_PLATFORM}-${NODE_ARCH}.${ARCHIVE_EXT}"
    NODE_URL="https://nodejs.org/dist/${NODE_VERSION}/${NODE_ARCHIVE}"

    print_info "Downloading Node.js ${NODE_VERSION} for ${NODE_PLATFORM}/${NODE_ARCH}..."
    download_with_progress "$NODE_URL" "$TMP_DIR/$NODE_ARCHIVE" "Node.js $NODE_VERSION"

    print_info "Installing Node.js to $JAWT_NODE_DIR..."
    # Remove existing installation if upgrading
    if [ -d "$JAWT_NODE_DIR" ]; then
        sudo rm -rf "$JAWT_NODE_DIR"
        sudo mkdir -p "$JAWT_NODE_DIR"
    fi

    extract_with_progress "$TMP_DIR/$NODE_ARCHIVE" "$JAWT_NODE_DIR" "Node.js"

    # Set proper ownership and permissions
    sudo chown -R root:root "$JAWT_NODE_DIR"
    sudo chmod -R 755 "$JAWT_NODE_DIR"

    print_success "Node.js installed successfully"
}

check_package_installed() {
    local package_name="$1"
    local install_dir="$2"

    if [ -d "$install_dir/node_modules/$package_name" ]; then
        return 0
    else
        return 1
    fi
}

install_typescript() {
    if [ "$INSTALL_TYPE" != "full" ]; then
        return 0
    fi

    print_info "Installing TypeScript..."

    # Use the installed Node.js to install TypeScript globally in our jawt location
    export PATH="${JAWT_NODE_DIR}/bin:$PATH"

    # Check if TypeScript is already installed
    if check_package_installed "typescript" "$JAWT_TSC_DIR"; then
        print_success "TypeScript is already installed"
    else
        # Install TypeScript to our jawt directory
        "$JAWT_NODE_DIR/bin/npm" install -g --prefix "$JAWT_TSC_DIR" typescript
        print_success "TypeScript installed successfully"
    fi
}

install_tailwindcss() {
    if [ "$INSTALL_TYPE" != "full" ]; then
        return 0
    fi

    print_info "Installing TailwindCSS..."

    # Use the installed Node.js to install TailwindCSS globally in our jawt location
    export PATH="${JAWT_NODE_DIR}/bin:$PATH"

    # Check if TailwindCSS is already installed
    if check_package_installed "tailwindcss" "$JAWT_TAILWIND_DIR"; then
        print_success "TailwindCSS is already installed"
    else
        # Install TailwindCSS and its dependencies to our jawt directory
        "$JAWT_NODE_DIR/bin/npm" install -g --prefix "$JAWT_TAILWIND_DIR" tailwindcss @tailwindcss/cli autoprefixer postcss
        print_success "TailwindCSS installed successfully"
    fi
}

install_jawt_executable() {
    print_info "Fetching latest Jawt release..."

    # Determine binary name based on OS and architecture
    BINARY_PATTERN="jawt-${OS_NAME}-${ARCH_NAME}"

    # Fetch the latest release and find the appropriate asset
    DL_URL=$(curl -fsSL "$API_URL" | jq -r ".assets[] | select(.name | test(\"${BINARY_PATTERN}\")) | .browser_download_url")

    if [ -z "$DL_URL" ]; then
        error_exit "No matching asset found for ${OS_NAME}/${ARCH_NAME}"
    fi

    # Download the binary
    download_with_progress "$DL_URL" "$TMP_DIR/$BINARY_NAME" "Jawt"

    # Install the binary
    print_info "Installing Jawt..."
    sudo cp "$TMP_DIR/$BINARY_NAME" "$JAWT_BIN_DIR/$BINARY_NAME"
    sudo chmod +x "$JAWT_BIN_DIR/$BINARY_NAME"

    print_success "Jawt installed to $JAWT_BIN_DIR"
}

create_system_executable() {
    print_info "Creating system-wide executable..."

    # Create a wrapper script that sets up the environment
    cat << 'EOF' | sudo tee "$SYSTEM_BIN_DIR/$BINARY_NAME" > /dev/null
#!/bin/bash

# Jawt Wrapper Script
# Sets up the internal environment and executes Jawt

JAWT_ROOT="/usr/local/jawt"
JAWT_BIN_DIR="${JAWT_ROOT}/bin"
JAWT_NODE_DIR="${JAWT_ROOT}/node"
JAWT_TSC_DIR="${JAWT_ROOT}/tsc"
JAWT_TAILWIND_DIR="${JAWT_ROOT}/tailwind"

# Set internal environment variables for Jawt's use only
export JAWT_ROOT
export JAWT_NODE_PATH="${JAWT_NODE_DIR}/bin/node"
export JAWT_NPM_PATH="${JAWT_NODE_DIR}/bin/npm"
export JAWT_NPX_PATH="${JAWT_NODE_DIR}/bin/npm"
export JAWT_TSC_PATH="${JAWT_TSC_DIR}/bin/tsc"
export JAWT_TAILWIND_PATH="${JAWT_TAILWIND_DIR}/bin/tailwindcss"

# Execute Jawt with all provided arguments
# Jawt can use the above environment variables to access its internal tools
exec "${JAWT_BIN_DIR}/jawt" "$@"
EOF

    sudo chmod +x "$SYSTEM_BIN_DIR/$BINARY_NAME"

    print_success "System-wide executable created at $SYSTEM_BIN_DIR/$BINARY_NAME"
}

update_shell_config() {
    print_info "Updating shell configuration..."

    SHELL_CONFIG_FILES=("$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile" "$HOME/.zsh_profile")
    CONFIG_UPDATED=false

    # Environment variables to add
    ENV_VARS=(
        "export JAWT_ROOT=\"$JAWT_ROOT\""
    )

    for CONFIG_FILE in "${SHELL_CONFIG_FILES[@]}"; do
        if [ -f "$CONFIG_FILE" ]; then
            # Add Jawt environment configuration
            if ! grep -q "JAWT_ROOT" "$CONFIG_FILE"; then
                echo "" >> "$CONFIG_FILE"
                echo "# Jawt Environment Configuration" >> "$CONFIG_FILE"
                for env_var in "${ENV_VARS[@]}"; do
                    echo "$env_var" >> "$CONFIG_FILE"
                done
                print_info "Updated $CONFIG_FILE with Jawt environment"
            fi
            CONFIG_UPDATED=true
            break
        fi
    done

    if [ "$CONFIG_UPDATED" = false ]; then
        print_warning "No supported shell configuration file found."
        print_info "Please add the following to your shell configuration:"
        for env_var in "${ENV_VARS[@]}"; do
            echo "  $env_var"
        done
    fi
}

verify_installation() {
    print_info "Verifying installation..."

    # Test Jawt
    if "$SYSTEM_BIN_DIR/$BINARY_NAME" --version >/dev/null 2>&1; then
        print_success "Jawt is working"
    else
        print_warning "Jawt test failed - you may need to restart your shell"
    fi

    if [ "$INSTALL_TYPE" = "full" ]; then
        # Test Node.js
        if "$JAWT_NODE_DIR/bin/node" --version >/dev/null 2>&1; then
            NODE_VER=$("$JAWT_NODE_DIR/bin/node" --version)
            print_success "Node.js is working (version: $NODE_VER)"
        else
            print_warning "Node.js test failed"
        fi

        # Test TypeScript
        if "$JAWT_TSC_DIR/bin/tsc" --version >/dev/null 2>&1; then
            TSC_VER=$("$JAWT_TSC_DIR/bin/tsc" --version)
            print_success "TypeScript is working ($TSC_VER)"
        else
            print_warning "TypeScript test failed"
        fi

        # Test TailwindCSS
        if "$JAWT_TAILWIND_DIR/bin/tailwindcss" --version >/dev/null 2>&1; then
            TAILWIND_VER=$("$JAWT_TAILWIND_DIR/bin/tailwindcss" --version)
            print_success "TailwindCSS is working ($TAILWIND_VER)"
        else
            print_warning "TailwindCSS test failed"
        fi
    fi
}

print_summary() {
    echo
    print_success "=== Jawt Installation Complete ==="
    echo
    echo "Installation type: $INSTALL_TYPE"
    echo "Installation directory: $JAWT_ROOT"
    echo "System executable: $SYSTEM_BIN_DIR/$BINARY_NAME"
    echo

    if [ "$INSTALL_TYPE" = "full" ]; then
        echo "Installed components:"
        echo "  • Jawt: $JAWT_BIN_DIR/$BINARY_NAME"
        echo "  • Jawt Node.js: $JAWT_NODE_DIR"
        echo "  • Jawt TypeScript: $JAWT_TSC_DIR"
        echo "  • Jawt TailwindCSS: $JAWT_TAILWIND_DIR"
        echo
        echo "Available internal tools (accessible via environment variables):"
        echo "  • JAWT_NODE_PATH   - Node.js runtime"
        echo "  • JAWT_NPM_PATH    - NPM package manager"
        echo "  • JAWT_TSC_PATH    - TypeScript compiler"
        echo "  • JAWT_TAILWIND_PATH - TailwindCSS CLI"
        echo
        echo "Available commands:"
        echo "  • jawt         - Just Another Web Tool"
    else
        echo "Installed components:"
        echo "  • Jawt: $JAWT_BIN_DIR/$BINARY_NAME"
        echo
        echo "Available commands:"
        echo "  • jawt         - Jawt command-line tool"
    fi

    echo
    print_info "You can now run 'jawt' from anywhere in your system!"
    print_info "Restart your shell or run 'source ~/.bashrc' (or equivalent) to update your environment."
}

main() {
    echo
    print_info "=== Jawt Installation Script ==="
    echo

    # Check for sudo access
    if ! sudo -n true 2>/dev/null; then
        print_info "This installation requires administrative privileges."
        sudo -v || error_exit "Unable to obtain administrative privileges"
    fi

    detect_platform
    check_dependencies
    prompt_installation_type
    create_directories
    install_nodejs
    install_typescript
    install_tailwindcss
    install_jawt_executable
    create_system_executable
    update_shell_config
    verify_installation
    print_summary
}

# Run main function
main "$@"