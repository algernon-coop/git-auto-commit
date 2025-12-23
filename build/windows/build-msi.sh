#!/bin/bash
# Build MSI installer for Windows using WiX Toolset
# This script is designed to run on Ubuntu with WiX Toolset installed via mono

set -e

ARCH=$1
VERSION=$2
BINARY_PATH=$3

if [ -z "$ARCH" ] || [ -z "$VERSION" ] || [ -z "$BINARY_PATH" ]; then
    echo "Usage: $0 <arch> <version> <binary_path>"
    echo "Example: $0 amd64 1.0.0 dist/git-auto-commit-windows-amd64.exe"
    exit 1
fi

# Set platform-specific variables
if [ "$ARCH" = "amd64" ]; then
    PLATFORM="x64"
    PROGRAM_FILES_FOLDER="ProgramFiles64Folder"
elif [ "$ARCH" = "arm64" ]; then
    PLATFORM="arm64"
    PROGRAM_FILES_FOLDER="ProgramFiles64Folder"
else
    echo "Unknown architecture: $ARCH"
    exit 1
fi

# Clean version (remove 'v' prefix if present)
CLEAN_VERSION="${VERSION#v}"

OUTPUT_DIR="dist"
WXS_FILE="build/windows/git-auto-commit.wxs"
WXS_OBJ="$OUTPUT_DIR/git-auto-commit-$ARCH.wixobj"
MSI_FILE="$OUTPUT_DIR/git-auto-commit-windows-$ARCH.msi"

echo "Building MSI for $PLATFORM architecture..."
echo "Version: $CLEAN_VERSION"
echo "Binary: $BINARY_PATH"

# Compile WiX source
wix build \
    -arch "$PLATFORM" \
    -d "Version=$CLEAN_VERSION" \
    -d "Platform=$PLATFORM" \
    -d "ProgramFilesFolder=$PROGRAM_FILES_FOLDER" \
    -d "BinaryPath=$BINARY_PATH" \
    -out "$MSI_FILE" \
    "$WXS_FILE"

echo "MSI created: $MSI_FILE"
