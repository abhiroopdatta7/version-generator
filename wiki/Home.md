# Version Generator Wiki

Welcome to the **Version Generator** documentation! This tool automatically generates semantic version numbers based on your Git repository state.

## 🚀 Quick Start

```bash
# Download the latest release
wget https://github.com/abhiroopdatta7/version-generator/releases/latest/download/version-generator

# Make it executable
chmod +x version-generator

# Generate version
./version-generator
```

## 📚 Documentation

### Getting Started
- **[Installation Guide](Installation)** - How to install and set up version-generator
- **[Quick Start Guide](Quick-Start)** - Get up and running in 5 minutes
- **[Basic Usage](Usage)** - Common use cases and examples

### Features
- **[Versioning Schemes](Versioning-Schemes)** - Semantic, Calendar, and Simple versioning
- **[Output Formats](Output-Formats)** - Go, C++, YAML, and plain text output
- **[Git Backends](Git-Backends)** - Built-in go-git vs system git
- **[Branch Handling](Branch-Handling)** - How different branches affect versioning

### Advanced
- **[CLI Reference](CLI-Reference)** - Complete command-line options
- **[Build from Source](Build-from-Source)** - Compilation and development
- **[CI/CD Integration](CI-CD-Integration)** - Using in automated pipelines
- **[Troubleshooting](Troubleshooting)** - Common issues and solutions

### Development
- **[Architecture](Architecture)** - Code structure and design
- **[Contributing](Contributing)** - How to contribute to the project
- **[Release Process](Release-Process)** - How releases are created

## 🎯 Key Features

| Feature | Description |
|---------|-------------|
| **Multiple Schemes** | Semantic, Calendar, and Simple versioning formats |
| **Dual Git Support** | Built-in go-git library or system git executable |
| **Branch Aware** | Intelligent handling of feature branches and main branches |
| **Multiple Outputs** | Generate Go, C++, YAML, or plain text files |
| **Zero Config** | Works out of the box in any Git repository |
| **Cross Platform** | Works on Windows, macOS, and Linux |

## 📖 Quick Examples

### Generate Semantic Version
```bash
./version-generator --semver
# Output: v1.2.3-feature-branch.5
```

### Generate Calendar Version
```bash
./version-generator --cal-ver
# Output: 2025.08.5-feature-branch
```

### Create Go Source File
```bash
./version-generator --go --semver
# Creates: version.go with embedded version
```

### Include Git Hash
```bash
./version-generator --hash
# Output: v1.2.3+abc1234
```

## 🔗 Links

- **[GitHub Repository](https://github.com/abhiroopdatta7/version-generator)**
- **[Latest Release](https://github.com/abhiroopdatta7/version-generator/releases/latest)**
- **[Issue Tracker](https://github.com/abhiroopdatta7/version-generator/issues)**

## 📝 Recent Updates

- **v1.0.1** - Fixed `--version` flag in distributed binaries
- **v1.0.0** - Major refactoring with improved architecture
- Added comprehensive .gitignore and build scripts
- Centralized version generation logic
