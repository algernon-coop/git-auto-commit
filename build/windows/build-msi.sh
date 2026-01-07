#!/bin/bash
# Build MSI installer for Windows using go-msi
# This script is designed to run on any platform

set -e

ARCH=$1
VERSION=$2
BINARY_PATH=$3

if [ -z "$ARCH" ] || [ -z "$VERSION" ] || [ -z "$BINARY_PATH" ]; then
    echo "Usage: $0 <arch> <version> <binary_path>"
    echo "Example: $0 amd64 1.0.0 dist/git-auto-commit-windows-amd64.exe"
    exit 1
fi

# Clean version (remove 'v' prefix if present)
CLEAN_VERSION="${VERSION#v}"

OUTPUT_DIR="dist"
MSI_FILE="$OUTPUT_DIR/git-auto-commit-windows-$ARCH.msi"

echo "Building MSI for $ARCH architecture..."
echo "Version: $CLEAN_VERSION"
echo "Binary: $BINARY_PATH"

# Check if binary exists
if [ ! -f "$BINARY_PATH" ]; then
    echo "Error: Binary not found at $BINARY_PATH"
    exit 1
fi

# Build MSI using go-msi
go-msi make \
    --msi "$MSI_FILE" \
    --version "$CLEAN_VERSION" \
    --arch "$ARCH" \
    --path "wix.json" \
    --src "build/windows/templates"

echo "MSI created: $MSI_FILE"
