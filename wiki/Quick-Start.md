# Quick Start Guide

Get up and running with version-generator in 5 minutes!

## 🚀 1. Install

```bash
# Download latest release
wget https://github.com/abhiroopdatta7/version-generator/releases/latest/download/version-generator
chmod +x version-generator
```

## 📂 2. Navigate to Your Git Repository

```bash
cd /path/to/your/git/repository
```

## ⚡ 3. Generate Your First Version

```bash
# Basic version generation
./version-generator
# Output: v1.2.3 (or similar based on your tags)
```

## 🎯 4. Try Different Formats

```bash
# Semantic versioning
./version-generator --semver
# Output: v1.2.3-feature-branch.5

# Calendar versioning
./version-generator --cal-ver  
# Output: 2025.08.5-feature-branch

# Simple format
./version-generator --simple
# Output: v1.2.3

# Include git hash
./version-generator --hash
# Output: v1.2.3+abc1234
```

## 📄 5. Generate Files

```bash
# Create Go source file
./version-generator --go --semver
# Creates: version.go

# Create C++ header
./version-generator --cpp
# Creates: version.h

# Create YAML file
./version-generator --yaml
# Creates: version.yaml

# Create plain text file
./version-generator --file
# Creates: .VERSION
```

## 🔧 6. Common Use Cases

### CI/CD Pipeline
```bash
# Get version for Docker tag
VERSION=$(./version-generator --semver)
docker build -t myapp:$VERSION .
```

### Release Script
```bash
#!/bin/bash
VERSION=$(./version-generator --semver)
echo "Building release $VERSION"
go build -ldflags="-X main.Version=$VERSION" -o myapp
```

### Git Tag Creation
```bash
# Generate version and create tag
VERSION=$(./version-generator --semver)
git tag -a $VERSION -m "Release $VERSION"
```

## 🌟 What's Next?

### Learn More
- **[Usage Guide](Usage)** - Detailed examples and use cases
- **[Versioning Schemes](Versioning-Schemes)** - Understanding different formats
- **[CLI Reference](CLI-Reference)** - Complete command options

### Advanced Features
- **[Branch Handling](Branch-Handling)** - How branches affect versioning
- **[Output Formats](Output-Formats)** - File generation options
- **[CI/CD Integration](CI-CD-Integration)** - Automation examples

## 💡 Pro Tips

### 1. Use Aliases
```bash
# Add to ~/.bashrc
alias vgen='./version-generator'
alias vsem='./version-generator --semver'
alias vcal='./version-generator --cal-ver'
```

### 2. Environment Variables
```bash
# In CI/CD
export VERSION=$(./version-generator --semver)
echo "Building version: $VERSION"
```

### 3. Makefile Integration
```makefile
VERSION := $(shell ./version-generator --semver)

build:
	go build -ldflags="-X main.Version=$(VERSION)" -o myapp

docker:
	docker build -t myapp:$(VERSION) .
```

### 4. Package.json Update
```bash
# Update Node.js package.json version
VERSION=$(./version-generator --semver | sed 's/v//')
npm version $VERSION --no-git-tag-version
```

## 🏷️ Understanding Your Output

### Tag-based Versioning
- **On a tag**: `v1.2.3` (exact tag)
- **After tag**: `v1.2.3+5` (5 commits since tag)
- **Feature branch**: `v1.2.3-feature-name+3` (3 commits on feature branch)

### Branch Behavior
- **main/master**: Clean versions without branch names
- **Feature branches**: Include sanitized branch name
- **Detached HEAD**: Uses "detached" as branch name

## ❓ Need Help?

- **Common issues**: [Troubleshooting](Troubleshooting)
- **File a bug**: [GitHub Issues](https://github.com/abhiroopdatta7/version-generator/issues)
- **Request feature**: [GitHub Issues](https://github.com/abhiroopdatta7/version-generator/issues)

## 📚 Examples by Project Type

### Go Project
```bash
# Generate Go version file
./version-generator --go --semver

# Use in build
go build -ldflags="-X main.Version=$(./version-generator --semver)"
```

### Node.js Project  
```bash
# Update package.json
VERSION=$(./version-generator --semver | sed 's/v//')
npm version $VERSION --no-git-tag-version
```

### Docker Project
```bash
# Build with version tag
docker build -t myapp:$(./version-generator --semver) .
```

### Python Project
```bash
# Update __version__ in setup.py
VERSION=$(./version-generator --semver | sed 's/v//')
sed -i "s/__version__ = .*/__version__ = '$VERSION'/" setup.py
```
