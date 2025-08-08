# CLI Reference

Complete command-line reference for version-generator.

## 📋 Basic Syntax

```bash
version-generator [flags]
```

## 🚩 Global Flags

### Version Information
```bash
-v, --version           Show version information
-h, --help              Show context-sensitive help
```

### Versioning Schemes
```bash
--semver                Use Semantic Versioning format
--cal-ver               Use Calendar Versioning format  
--simple                Use simple version format (no branch info)
--hash                  Include short hash in version
```

### Git Backend
```bash
-i, --in-built-git      Use built-in go-git library instead of system git
```

### File Generation
```bash
-g, --go                Generate Go format version file
--go-path=PATH          Path for Go file (default: version.go)

-c, --cpp               Generate C++ format version file  
--cpp-path=PATH         Path for C++ file (default: version.h)

-y, --yaml              Generate YAML format version file
--yaml-path=PATH        Path for YAML file (default: version.yaml)

-f, --file              Write version to file
--file-path=PATH        Path for file (default: .VERSION)
```

## 📖 Detailed Flag Documentation

### `--semver` - Semantic Versioning
Generates version strings following [Semantic Versioning](https://semver.org/) specification.

**Examples:**
```bash
./version-generator --semver
# On tag: v1.2.3
# After tag: v1.2.3-dev.5
# Feature branch: v1.2.3-feature-name.3
```

**Format:**
- `v{MAJOR}.{MINOR}.{PATCH}` (on exact tag)
- `v{MAJOR}.{MINOR}.{PATCH}-dev.{COUNT}` (main branch after tag)
- `v{MAJOR}.{MINOR}.{PATCH}-{BRANCH}.{COUNT}` (feature branch)

### `--cal-ver` - Calendar Versioning
Generates version strings based on current date following [Calendar Versioning](https://calver.org/).

**Examples:**
```bash
./version-generator --cal-ver
# Output: 2025.08
# With commits: 2025.08.5
# Feature branch: 2025.08.5-feature-name
```

**Format:**
- `{YEAR}.{MONTH}` (no commits since midnight)
- `{YEAR}.{MONTH}.{COUNT}` (with commits)
- `{YEAR}.{MONTH}.{COUNT}-{BRANCH}` (feature branch)

### `--simple` - Simple Format
Generates basic version strings without branch information.

**Examples:**
```bash
./version-generator --simple
# Output: v1.2.3 (always just the tag)
```

### `--hash` - Include Git Hash
Adds the short Git commit hash to the version string.

**Examples:**
```bash
./version-generator --hash
# Output: v1.2.3+abc1234

./version-generator --semver --hash
# Output: v1.2.3-dev.5+abc1234
```

### `-i, --in-built-git` - Built-in Git
Uses the built-in go-git library instead of system git executable.

**When to use:**
- System git not available
- Docker containers without git
- Environments with restricted system access
- Consistent behavior across platforms

**Examples:**
```bash
./version-generator --in-built-git
./version-generator -i --semver
```

## 📄 File Generation Flags

### `--go` - Generate Go File
Creates a Go source file with the version as a constant.

**Default filename:** `version.go`
**Custom path:** `--go-path=custom/path/version.go`

**Generated content:**
```go
package main

const Version = "v1.2.3"
```

**Examples:**
```bash
# Default file
./version-generator --go --semver

# Custom path  
./version-generator --go --go-path=pkg/version.go

# Custom directory (filename auto-added)
./version-generator --go --go-path=src/
```

### `--cpp` - Generate C++ File
Creates a C++ header file with version definitions.

**Default filename:** `version.h`
**Custom path:** `--cpp-path=include/version.h`

**Generated content:**
```cpp
#ifndef VERSION_H
#define VERSION_H

#define VERSION "v1.2.3"

#endif // VERSION_H
```

### `--yaml` - Generate YAML File
Creates a YAML file with version information.

**Default filename:** `version.yaml`
**Custom path:** `--yaml-path=config/version.yaml`

**Generated content:**
```yaml
version: v1.2.3
```

### `--file` - Generate Plain Text File
Creates a plain text file with just the version string.

**Default filename:** `.VERSION`
**Custom path:** `--file-path=VERSION.txt`

**Generated content:**
```
v1.2.3
```

## 🔄 Flag Combinations

### Multiple Schemes (Invalid)
```bash
# ❌ Cannot use multiple versioning schemes
./version-generator --semver --cal-ver
```

### Scheme + Hash
```bash
# ✅ Combine any scheme with hash
./version-generator --semver --hash
./version-generator --cal-ver --hash
./version-generator --simple --hash
```

### Multiple File Types
```bash
# ✅ Generate multiple files
./version-generator --go --cpp --yaml --file
```

### Scheme + Files
```bash
# ✅ Use specific scheme for file generation
./version-generator --semver --go --cpp
```

## 📁 Path Handling

### Relative Paths
```bash
./version-generator --go-path=src/version.go
./version-generator --cpp-path=../include/version.h
```

### Absolute Paths
```bash
./version-generator --go-path=/opt/app/version.go
```

### Directory Paths (ending with /)
```bash
# Automatically appends default filename
./version-generator --go-path=src/          # Creates: src/version.go
./version-generator --cpp-path=include/     # Creates: include/version.h
```

## 🚀 Advanced Usage

### Environment Integration
```bash
# Export version for use in scripts
export VERSION=$(./version-generator --semver)
echo $VERSION

# Use in build commands
go build -ldflags="-X main.Version=$(./version-generator --semver)"
```

### CI/CD Pipelines
```bash
# GitHub Actions
VERSION=$(./version-generator --semver)
echo "version=$VERSION" >> $GITHUB_OUTPUT

# GitLab CI
VERSION=$(./version-generator --semver)
echo "VERSION=$VERSION" >> build.env
```

### Docker Integration
```bash
# Build with version tag
docker build -t myapp:$(./version-generator --semver) .

# Multi-stage build
VERSION=$(./version-generator --semver) docker build --build-arg VERSION=$VERSION .
```

## 🔍 Exit Codes

| Code | Meaning |
|------|---------|
| 0 | Success |
| 1 | General error |
| 2 | Git repository not found |
| 3 | Git command failed |
| 4 | File write error |

## 📚 Examples by Use Case

### Development
```bash
# Quick version check
./version-generator

# Detailed version with hash
./version-generator --hash
```

### Build Scripts
```bash
# Makefile compatible
VERSION := $(shell ./version-generator --semver)
```

### Release Automation
```bash
# Create all version files
./version-generator --semver --go --cpp --yaml --file
```

### Container Builds
```bash
# Use built-in git for consistency  
./version-generator --in-built-git --semver
```
