#!/bin/bash

# Build script for version-generator with proper version embedding

set -e

# Get version information
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_DATE=$(date -u +"%Y-%m-%dT%H:%M:%SZ")

# If we're on a tag, use just the tag name
if git describe --exact-match --tags HEAD >/dev/null 2>&1; then
    VERSION=$(git describe --exact-match --tags HEAD)
fi

echo "Building version-generator..."
echo "Version: $VERSION"
echo "Git Commit: $GIT_COMMIT"
echo "Build Date: $BUILD_DATE"

# Build flags
LDFLAGS="-s -w"
LDFLAGS="$LDFLAGS -X main.Version=$VERSION"
LDFLAGS="$LDFLAGS -X main.GitCommit=$GIT_COMMIT"
LDFLAGS="$LDFLAGS -X main.BuildDate=$BUILD_DATE"

# Build the binary
go build -ldflags="$LDFLAGS" -o version-generator .

echo "Build complete: ./version-generator"
echo "Testing version output:"
./version-generator --version
