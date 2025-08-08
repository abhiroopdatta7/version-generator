# Output Formats

Understanding the different output formats and their use cases.

## 📄 Format Types

version-generator supports multiple output formats to suit different integration needs:
- **Console Output**: Direct terminal output (default)
- **File Generation**: Create files with version information
- **Multiple Formats**: Generate several files at once

## 🖥️ Console Output Formats

### Default Format
```bash
./version-generator
# Output: v1.2.3+5
```
- **Pattern**: `v{major}.{minor}.{patch}+{commits}`
- **Use Case**: General purpose, build scripts
- **Character Set**: ASCII alphanumeric with `v`, `.`, `+`

### Semantic Versioning (SemVer)
```bash
./version-generator --semver
# Output: v1.2.3-dev.5
```
- **Pattern**: `v{major}.{minor}.{patch}-{prerelease}.{commits}`
- **Use Case**: Package managers, dependency management
- **Compliance**: [Semantic Versioning 2.0.0](https://semver.org/)

### Calendar Versioning (CalVer)
```bash
./version-generator --cal-ver
# Output: 2025.08.5
```
- **Pattern**: `{YYYY}.{MM}.{commits}`
- **Use Case**: Time-based releases, continuous deployment
- **Components**: Year (4-digit), Month (2-digit), Commit count

### Simple Format
```bash
./version-generator --simple
# Output: v1.2.3
```
- **Pattern**: `v{major}.{minor}.{patch}`
- **Use Case**: Clean tags, release names
- **Note**: No commit count information

### Hash-Enhanced Format
```bash
./version-generator --hash
# Output: v1.2.3+5+abc1234
```
- **Pattern**: `{base_version}+{git_hash}`
- **Use Case**: Unique identification, debugging
- **Hash**: First 7 characters of git commit SHA

## 📂 File Generation Formats

### Go Source Files

#### Basic Go File
```bash
./version-generator --go --semver
# Creates: version.go
```

**Generated Content:**
```go
package main

const Version = "v1.2.3-dev.5"
```

#### Custom Package and Path
```bash
./version-generator --go --go-path=pkg/version/version.go --semver
# Creates: pkg/version/version.go
```

**Generated Content:**
```go
package version

const Version = "v1.2.3-dev.5"
```

#### Directory Auto-naming
```bash
./version-generator --go --go-path=internal/ --semver
# Creates: internal/version.go
```

### C++ Header Files

#### Basic Header
```bash
./version-generator --cpp --semver
# Creates: version.h
```

**Generated Content:**
```cpp
#ifndef VERSION_H
#define VERSION_H

#define VERSION "v1.2.3-dev.5"

#endif // VERSION_H
```

#### Custom Header Path
```bash
./version-generator --cpp --cpp-path=include/myapp/version.h --semver
# Creates: include/myapp/version.h
```

**Generated Content:**
```cpp
#ifndef MYAPP_VERSION_H
#define MYAPP_VERSION_H

#define VERSION "v1.2.3-dev.5"

#endif // MYAPP_VERSION_H
```

### YAML Configuration Files

#### Basic YAML
```bash
./version-generator --yaml --cal-ver
# Creates: version.yaml
```

**Generated Content:**
```yaml
version: 2025.08.5
```

#### Custom YAML Path
```bash
./version-generator --yaml --yaml-path=config/app.yaml --semver
# Creates: config/app.yaml
```

**Generated Content:**
```yaml
version: v1.2.3-dev.5
```

### Plain Text Files

#### Basic Text File
```bash
./version-generator --file --simple
# Creates: .VERSION
```

**Generated Content:**
```
v1.2.3
```

#### Custom Text File
```bash
./version-generator --file --file-path=VERSION.txt --semver
# Creates: VERSION.txt
```

**Generated Content:**
```
v1.2.3-dev.5
```

## 🔧 Path Resolution Rules

### Automatic File Naming

When a directory is specified as a path, version-generator automatically determines the filename:

```bash
# Directory paths get default filenames
./version-generator --go --go-path=src/          # Creates: src/version.go
./version-generator --cpp --cpp-path=include/    # Creates: include/version.h
./version-generator --yaml --yaml-path=config/   # Creates: config/version.yaml
./version-generator --file --file-path=build/    # Creates: build/.VERSION
```

### Explicit File Naming

When a file path is specified, the exact path is used:

```bash
# File paths are used exactly as specified
./version-generator --go --go-path=src/app_version.go      # Creates: src/app_version.go
./version-generator --cpp --cpp-path=inc/app.h            # Creates: inc/app.h
./version-generator --yaml --yaml-path=cfg/release.yaml   # Creates: cfg/release.yaml
./version-generator --file --file-path=build/BUILD_ID     # Creates: build/BUILD_ID
```

### Directory Creation

version-generator automatically creates directories if they don't exist:

```bash
# Creates deep directory structure
./version-generator --go --go-path=deep/nested/path/version.go
# Result: Creates 'deep/nested/path/' directories and version.go file
```

## 📋 Format Comparison

| Format | Example | Length | Use Case | Special Characters |
|--------|---------|--------|----------|-------------------|
| Default | `v1.2.3+5` | Short | Build scripts | `v`, `.`, `+` |
| SemVer | `v1.2.3-dev.5` | Medium | Package managers | `v`, `.`, `-` |
| CalVer | `2025.08.5` | Short | Time-based releases | `.` |
| Simple | `v1.2.3` | Shortest | Clean tags | `v`, `.` |
| Hash | `v1.2.3+5+abc1234` | Longest | Debugging | `v`, `.`, `+` |

## 🎯 Use Case Matrix

| Scenario | Recommended Format | Reason |
|----------|-------------------|---------|
| Docker tags | CalVer or Simple | Clean, sortable |
| npm packages | SemVer | Package manager compliance |
| Build artifacts | Default or Hash | Unique identification |
| Git tags | Simple or SemVer | Clean repository history |
| CI/CD pipelines | CalVer or Default | Time-based tracking |
| Debug builds | Hash | Exact commit tracking |
| Release notes | Simple | Clean presentation |
| Internal builds | Default | Balance of info and brevity |

## 🔄 Multi-Format Generation

### Generate Multiple Files
```bash
# Create all file types with different versions
./version-generator --go --cpp --yaml --file --semver

# Results in:
# - version.go (Go source)
# - version.h (C++ header)  
# - version.yaml (YAML config)
# - .VERSION (text file)
```

### Custom Paths for Each
```bash
./version-generator \
  --go --go-path=src/version.go \
  --cpp --cpp-path=include/version.h \
  --yaml --yaml-path=config/version.yaml \
  --file --file-path=build/VERSION.txt \
  --semver
```

## 📊 Advanced Examples

### Environment-Specific Formats
```bash
#!/bin/bash

ENVIRONMENT=${1:-development}

case $ENVIRONMENT in
  production)
    # Clean version for production
    ./version-generator --simple > production.version
    ;;
  staging)
    # Calendar version for staging
    ./version-generator --cal-ver > staging.version
    ;;
  development)
    # Full version with hash for development
    ./version-generator --hash > development.version
    ;;
esac
```

### Build System Integration
```makefile
VERSION_FORMATS := go cpp yaml file
VERSION_TYPE := semver

.PHONY: version-files
version-files:
	./version-generator --$(VERSION_TYPE) $(addprefix --,$(VERSION_FORMATS))

.PHONY: clean-version
clean-version:
	rm -f version.go version.h version.yaml .VERSION
```

### Cross-Platform File Generation
```bash
#!/bin/bash
# generate-all-versions.sh

set -e

echo "Generating version files..."

# Console versions
echo "Console formats:"
echo "  Default: $(./version-generator)"
echo "  SemVer:  $(./version-generator --semver)"
echo "  CalVer:  $(./version-generator --cal-ver)"
echo "  Simple:  $(./version-generator --simple)"
echo "  Hash:    $(./version-generator --hash)"

# File generation
echo "Generating files..."
./version-generator --go --cpp --yaml --file --semver

echo "Generated files:"
ls -la version.go version.h version.yaml .VERSION

echo "File contents:"
echo "=== version.go ==="
cat version.go
echo "=== version.h ==="
cat version.h
echo "=== version.yaml ==="
cat version.yaml
echo "=== .VERSION ==="
cat .VERSION
```

## 🔍 Format Validation

### Semantic Version Validation
```bash
#!/bin/bash

VERSION=$(./version-generator --semver)

# Check SemVer compliance
if [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-[a-zA-Z0-9.-]+)?(\+[a-zA-Z0-9.-]+)?$ ]]; then
    echo "✅ Valid SemVer: $VERSION"
else
    echo "❌ Invalid SemVer: $VERSION"
    exit 1
fi
```

### Calendar Version Validation
```bash
#!/bin/bash

VERSION=$(./version-generator --cal-ver)

# Check CalVer format (YYYY.MM.patch)
if [[ $VERSION =~ ^[0-9]{4}\.[0-9]{2}\.[0-9]+$ ]]; then
    echo "✅ Valid CalVer: $VERSION"
else
    echo "❌ Invalid CalVer: $VERSION"
    exit 1
fi
```

### File Content Validation
```bash
#!/bin/bash

# Generate and validate Go file
./version-generator --go --semver

if go run -c version.go 2>/dev/null; then
    echo "✅ Generated Go file is valid"
else
    echo "❌ Generated Go file has syntax errors"
fi

# Validate C++ header
if echo '#include "version.h"' | gcc -x c++ -c - 2>/dev/null; then
    echo "✅ Generated C++ header is valid"
else
    echo "❌ Generated C++ header has syntax errors"
fi

# Validate YAML
if python -c "import yaml; yaml.safe_load(open('version.yaml'))" 2>/dev/null; then
    echo "✅ Generated YAML is valid"
else
    echo "❌ Generated YAML has syntax errors"
fi
```

## 📋 Output Format Summary

version-generator provides flexible output options to integrate seamlessly into any development workflow:

- **Multiple versioning schemes** for different use cases
- **Language-specific file generation** for Go, C++, YAML, and text
- **Automatic path resolution** with directory creation
- **Multi-format generation** in single command
- **Standards compliance** with SemVer and CalVer specifications

Choose the format that best fits your project's needs and build system requirements.
