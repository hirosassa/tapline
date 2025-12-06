#!/bin/bash
set -e

echo "Installing Tapline..."

# Build the binary
echo "Building binary..."
go build -o tapline ./cmd/tapline

# Check if GOPATH/bin is in PATH
if [ -d "$GOPATH/bin" ]; then
    INSTALL_DIR="$GOPATH/bin"
elif [ -d "$HOME/go/bin" ]; then
    INSTALL_DIR="$HOME/go/bin"
else
    INSTALL_DIR="/usr/local/bin"
fi

# Copy binary to install directory
echo "Installing to $INSTALL_DIR..."
if [ "$INSTALL_DIR" = "/usr/local/bin" ]; then
    sudo cp tapline "$INSTALL_DIR/"
else
    cp tapline "$INSTALL_DIR/"
fi

echo ""
echo "Tapline installed successfully!"
echo ""
echo "To use Tapline with Claude Code:"
echo "1. Ensure the .claude/hooks.json file is in your project directory"
echo "2. Start using Claude Code - logs will be automatically collected"
echo ""
echo "Manual usage:"
echo "  tapline conversation_start"
echo "  tapline user_prompt \"Your message\""
echo "  tapline assistant_response \"Assistant response\""
echo "  tapline conversation_end"
