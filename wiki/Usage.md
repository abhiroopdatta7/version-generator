# Usage Examples

Comprehensive examples of using version-generator in real-world scenarios.

## 🚀 Basic Usage

### Simple Version Generation
```bash
# Get current version
./version-generator
# Output: v1.2.3+5

# Check help
./version-generator --help

# Check tool version
./version-generator --version
```

### Different Versioning Schemes
```bash
# Semantic versioning
./version-generator --semver
# Output: v1.2.3-dev.5

# Calendar versioning  
./version-generator --cal-ver
# Output: 2025.08.5

# Simple format
./version-generator --simple
# Output: v1.2.3

# Include git hash
./version-generator --hash
# Output: v1.2.3+5+abc1234
```

## 📄 File Generation

### Go Source Files
```bash
# Generate version.go
./version-generator --go --semver

# Custom path
./version-generator --go --go-path=pkg/version/version.go

# Directory with auto-filename
./version-generator --go --go-path=internal/
```

**Generated `version.go`:**
```go
package main

const Version = "v1.2.3-dev.5"
```

### C++ Header Files
```bash
# Generate version.h
./version-generator --cpp --semver

# Custom path
./version-generator --cpp --cpp-path=include/myapp/version.h
```

**Generated `version.h`:**
```cpp
#ifndef VERSION_H
#define VERSION_H

#define VERSION "v1.2.3-dev.5"

#endif // VERSION_H
```

### YAML Configuration
```bash
# Generate version.yaml
./version-generator --yaml --cal-ver

# Custom path
./version-generator --yaml --yaml-path=config/version.yaml
```

**Generated `version.yaml`:**
```yaml
version: 2025.08.5
```

### Plain Text Files
```bash
# Generate .VERSION file
./version-generator --file --simple

# Custom filename
./version-generator --file --file-path=VERSION.txt
```

## 🔧 Build Integration

### Makefile Integration
```makefile
# Get version for builds
VERSION := $(shell ./version-generator --semver)

build:
	@echo "Building version: $(VERSION)"
	go build -ldflags="-X main.Version=$(VERSION)" -o myapp

docker:
	docker build -t myapp:$(VERSION) .

release: build
	tar -czf myapp-$(VERSION).tar.gz myapp
	
.PHONY: build docker release
```

### Shell Scripts
```bash
#!/bin/bash

# Build script with version embedding
set -e

VERSION=$(./version-generator --semver)
echo "Building version: $VERSION"

# Build Go binary
go build -ldflags="-X main.Version=$VERSION" -o myapp

# Create archives
tar -czf myapp-$VERSION-linux.tar.gz myapp
zip myapp-$VERSION-linux.zip myapp

echo "Build complete: $VERSION"
```

### npm/package.json Updates
```bash
#!/bin/bash

# Update Node.js package version
VERSION=$(./version-generator --semver | sed 's/v//')
npm version $VERSION --no-git-tag-version
echo "Updated package.json to $VERSION"
```

## 🐳 Docker Integration

### Dockerfile with Version
```dockerfile
# Build stage
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

# Get version during build
RUN ./version-generator --semver > /tmp/version
RUN VERSION=$(cat /tmp/version) && \
    go build -ldflags="-X main.Version=$VERSION" -o myapp

# Runtime stage
FROM alpine:latest
COPY --from=builder /app/myapp /usr/local/bin/
CMD ["myapp"]
```

### Docker Build with Tags
```bash
# Build with version tag
VERSION=$(./version-generator --semver)
docker build -t myapp:$VERSION .
docker build -t myapp:latest .

# Push with version
docker push myapp:$VERSION
docker push myapp:latest
```

### Docker Compose
```yaml
version: '3.8'
services:
  app:
    build: .
    image: myapp:${VERSION:-latest}
    environment:
      - APP_VERSION=${VERSION}
```

```bash
# Run with version
VERSION=$(./version-generator --semver) docker-compose up
```

## 🔄 CI/CD Integration

### GitHub Actions
```yaml
name: Build and Release

on:
  push:
    branches: [ main ]
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v3
      with:
        fetch-depth: 0  # Important for version generation
        
    - name: Setup Go
      uses: actions/setup-go@v3
      with:
        go-version: 1.21
        
    - name: Get version
      id: version
      run: |
        VERSION=$(./version-generator --semver)
        echo "version=$VERSION" >> $GITHUB_OUTPUT
        echo "Version: $VERSION"
        
    - name: Build
      run: |
        go build -ldflags="-X main.Version=${{ steps.version.outputs.version }}" -o myapp
        
    - name: Create Release
      if: github.ref == 'refs/heads/main'
      uses: actions/create-release@v1
      with:
        tag_name: ${{ steps.version.outputs.version }}
        release_name: Release ${{ steps.version.outputs.version }}
```

### GitLab CI
```yaml
stages:
  - build
  - release

variables:
  GO_VERSION: "1.21"

before_script:
  - export VERSION=$(./version-generator --semver)
  - echo "Building version $VERSION"

build:
  stage: build
  image: golang:$GO_VERSION
  script:
    - go build -ldflags="-X main.Version=$VERSION" -o myapp
  artifacts:
    paths:
      - myapp
    expire_in: 1 hour

release:
  stage: release
  only:
    - main
  script:
    - echo "Releasing version $VERSION"
    - # Release commands here
```

### Jenkins Pipeline
```groovy
pipeline {
    agent any
    
    environment {
        VERSION = sh(script: './version-generator --semver', returnStdout: true).trim()
    }
    
    stages {
        stage('Build') {
            steps {
                echo "Building version ${VERSION}"
                sh "go build -ldflags='-X main.Version=${VERSION}' -o myapp"
            }
        }
        
        stage('Test') {
            steps {
                sh './myapp --version'
            }
        }
        
        stage('Release') {
            when { branch 'main' }
            steps {
                sh "docker build -t myapp:${VERSION} ."
                sh "docker push myapp:${VERSION}"
            }
        }
    }
}
```

## 🌐 Language-Specific Examples

### Go Projects
```bash
# Generate Go version file
./version-generator --go --semver

# Build with embedded version
VERSION=$(./version-generator --semver)
go build -ldflags="-X main.Version=$VERSION -X main.BuildTime=$(date -u +%Y-%m-%dT%H:%M:%SZ)" -o myapp

# Module versioning
git tag $(./version-generator --semver)
```

### Node.js Projects
```bash
# Update package.json
VERSION=$(./version-generator --semver | sed 's/v//')
npm version $VERSION --no-git-tag-version

# Build with version
VERSION=$(./version-generator --semver) npm run build

# Environment variable
export REACT_APP_VERSION=$(./version-generator --semver)
npm start
```

### Python Projects
```bash
# Update setup.py
VERSION=$(./version-generator --semver | sed 's/v//')
sed -i "s/version='.*'/version='$VERSION'/" setup.py

# Update __init__.py
echo "__version__ = '$(./version-generator --semver | sed 's/v//')'" > mypackage/__init__.py

# Build package
python setup.py sdist bdist_wheel
```

### Rust Projects
```bash
# Update Cargo.toml
VERSION=$(./version-generator --semver | sed 's/v//')
sed -i "s/^version = .*/version = \"$VERSION\"/" Cargo.toml

# Build with version
cargo build --release
```

### C++ Projects
```bash
# Generate header
./version-generator --cpp --semver

# CMake integration
VERSION=$(./version-generator --semver | sed 's/v//')
cmake -DVERSION=$VERSION ..
make
```

## 📊 Advanced Workflows

### Multi-Environment Deployment
```bash
#!/bin/bash

BRANCH=$(git branch --show-current)
BASE_VERSION=$(./version-generator --semver)

case $BRANCH in
  main)
    VERSION=$BASE_VERSION
    ENV="production"
    ;;
  develop)
    VERSION="$BASE_VERSION-dev"
    ENV="staging"
    ;;
  *)
    VERSION="$BASE_VERSION-$BRANCH"
    ENV="development"
    ;;
esac

echo "Deploying $VERSION to $ENV"
docker build -t myapp:$VERSION .
docker tag myapp:$VERSION myapp:$ENV
```

### Release Automation
```bash
#!/bin/bash
# complete-release.sh

set -e

echo "Starting release process..."

# Generate version
VERSION=$(./version-generator --semver)
echo "Release version: $VERSION"

# Update version files
./version-generator --go --cpp --yaml --file --semver

# Commit version files
git add version.go version.h version.yaml .VERSION
git commit -m "chore: update version files for $VERSION"

# Create and push tag
git tag -a $VERSION -m "Release $VERSION"
git push origin $VERSION

# Build release artifacts
go build -ldflags="-X main.Version=$VERSION" -o myapp-linux
GOOS=windows go build -ldflags="-X main.Version=$VERSION" -o myapp-windows.exe
GOOS=darwin go build -ldflags="-X main.Version=$VERSION" -o myapp-macos

# Create GitHub release
gh release create $VERSION \
  --title "Release $VERSION" \
  --generate-notes \
  myapp-linux myapp-windows.exe myapp-macos

echo "Release $VERSION complete!"
```

### Branch-specific Versioning
```bash
#!/bin/bash

BRANCH=$(git branch --show-current)

case $BRANCH in
  main|master)
    # Production: clean semver
    VERSION=$(./version-generator --semver)
    ;;
  develop|staging)
    # Staging: calendar version
    VERSION=$(./version-generator --cal-ver)
    ;;
  *)
    # Feature: hash for uniqueness
    VERSION=$(./version-generator --semver --hash)
    ;;
esac

echo "Branch: $BRANCH, Version: $VERSION"
```

## 🔍 Debugging and Troubleshooting

### Verbose Version Information
```bash
# Get detailed version info
echo "Tool version: $(./version-generator --version)"
echo "Generated version: $(./version-generator --semver)"
echo "Current branch: $(git branch --show-current)"
echo "Last tag: $(git describe --tags --abbrev=0 2>/dev/null || echo 'none')"
echo "Commits since tag: $(git rev-list --count HEAD ^$(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~999) 2>/dev/null || echo 0)"
```

### Test Different Backends
```bash
# Compare system git vs built-in git
echo "System git: $(./version-generator --semver)"
echo "Built-in git: $(./version-generator --semver --in-built-git)"
```

### Validation Scripts
```bash
#!/bin/bash
# validate-version.sh

VERSION=$(./version-generator --semver)

# Validate semver format
if [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+.*$ ]]; then
    echo "✅ Valid semver: $VERSION"
else
    echo "❌ Invalid semver: $VERSION"
    exit 1
fi

# Test all formats
echo "Default: $(./version-generator)"
echo "SemVer: $(./version-generator --semver)"
echo "CalVer: $(./version-generator --cal-ver)"
echo "Simple: $(./version-generator --simple)"
```
