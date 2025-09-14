#!/bin/bash

# SetupSuite Installation Script
# This script downloads and installs the SetupSuite CLI tool

set -e

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="setupsuite"
GITHUB_REPO="FlubioStudios/setupsuite"  # Update with actual repo when published
VERSION="latest"

echo "Installing SetupSuite..."

# Check if running as root
if [[ $EUID -eq 0 ]]; then
   echo "This script should not be run as root for the download, but SetupSuite itself requires root to run." 
   echo "Please run this script as a regular user."
   exit 1
fi

# Check if Go is installed (for building from source)
if ! command -v go &> /dev/null; then
    echo "Go is not installed. Installing Go..."
    
    # Detect architecture
    ARCH=$(uname -m)
    case $ARCH in
        x86_64) GOARCH="amd64" ;;
        aarch64) GOARCH="arm64" ;;
        armv6l) GOARCH="armv6l" ;;
        armv7l) GOARCH="armv6l" ;;
        *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
    esac
    
    # Download and install Go
    GO_VERSION="1.21.3"
    wget -O /tmp/go.tar.gz "https://golang.org/dl/go${GO_VERSION}.linux-${GOARCH}.tar.gz"
    sudo rm -rf /usr/local/go
    sudo tar -C /usr/local -xzf /tmp/go.tar.gz
    export PATH=$PATH:/usr/local/go/bin
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
    rm /tmp/go.tar.gz
fi

# Create temporary directory
TEMP_DIR=$(mktemp -d)
cd "$TEMP_DIR"

# Clone or download the repository
echo "Downloading SetupSuite source code..."
if command -v git &> /dev/null; then
    git clone "https://github.com/${GITHUB_REPO}.git" .
else
    echo "Git not found. Please install git first:"
    echo "sudo apt-get update && sudo apt-get install -y git"
    exit 1
fi

# Build the binary
echo "Building SetupSuite..."
go build -o "$BINARY_NAME" ./suite

# Install the binary
echo "Installing SetupSuite to $INSTALL_DIR..."
sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"

# Clean up
cd /
rm -rf "$TEMP_DIR"

echo "SetupSuite installed successfully!"
echo ""
echo "Usage:"
echo "  sudo setupsuite -help                     # Show help"
echo "  sudo setupsuite -generate -type web       # Generate web server config"
echo "  sudo setupsuite                           # Run setup with default config"
echo ""
echo "Example workflow:"
echo "  1. sudo setupsuite -generate -type web -config /etc/setupsuite/web.sscfg"
echo "  2. sudo nano /etc/setupsuite/web.sscfg    # Edit the config"
echo "  3. sudo setupsuite -config /etc/setupsuite/web.sscfg"
echo ""
echo "⚠️  WARNING: SetupSuite makes significant changes to your system."
echo "   Always test on a VM or non-production server first!"
