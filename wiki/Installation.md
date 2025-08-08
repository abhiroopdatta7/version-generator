# Installation Guide

This guide covers all the ways to install and set up version-generator.

## 📦 Method 1: Download Pre-built Binary (Recommended)

### Latest Release
```bash
# Download latest release for Linux
wget https://github.com/abhiroopdatta7/version-generator/releases/latest/download/version-generator

# Make executable
chmod +x version-generator

# Move to PATH (optional)
sudo mv version-generator /usr/local/bin/
```

### Specific Version
```bash
# Download v1.0.1
wget https://github.com/abhiroopdatta7/version-generator/releases/download/v1.0.1/version-generator
chmod +x version-generator
```

## 🛠️ Method 2: Build from Source

### Prerequisites
- Go 1.21 or later
- Git

### Clone and Build
```bash
# Clone repository
git clone https://github.com/abhiroopdatta7/version-generator.git
cd version-generator

# Build with embedded version
./build.sh

# Or build manually
go build -o version-generator .
```

### Install to PATH
```bash
# Copy to local bin
cp version-generator ~/.local/bin/

# Or system-wide
sudo cp version-generator /usr/local/bin/
```

## 🐳 Method 3: Docker (Coming Soon)

```bash
# Run in Docker container
docker run --rm -v $(pwd):/repo ghcr.io/abhiroopdatta7/version-generator
```

## 🔧 Method 4: Go Install

```bash
# Install directly with Go
go install github.com/abhiroopdatta7/version-generator@latest
```

## ✅ Verify Installation

```bash
# Check version
version-generator --version

# Basic usage test
version-generator --help
```

## 🖥️ Platform-Specific Instructions

### Windows
```powershell
# Download using PowerShell
Invoke-WebRequest -Uri "https://github.com/abhiroopdatta7/version-generator/releases/latest/download/version-generator" -OutFile "version-generator.exe"

# Run
.\version-generator.exe --version
```

### macOS
```bash
# Download and install
curl -L https://github.com/abhiroopdatta7/version-generator/releases/latest/download/version-generator -o version-generator
chmod +x version-generator
sudo mv version-generator /usr/local/bin/
```

### Linux (Package Managers)

#### Homebrew
```bash
# Coming soon
brew install abhiroopdatta7/tap/version-generator
```

#### Snap
```bash
# Coming soon
sudo snap install version-generator
```

## 🔧 Configuration

No configuration files are needed! version-generator works out of the box in any Git repository.

### Optional: Global Configuration
You can create shell aliases for convenience:
```bash
# Add to ~/.bashrc or ~/.zshrc
alias vgen='version-generator'
alias vgen-semver='version-generator --semver'
alias vgen-calver='version-generator --cal-ver'
```

## 🚀 Next Steps

- [Quick Start Guide](Quick-Start) - Get started in 5 minutes
- [Usage Examples](Usage) - Common use cases
- [CLI Reference](CLI-Reference) - All command options

## 🐛 Installation Issues

### Common Problems

**"Command not found"**
- Ensure the binary is executable: `chmod +x version-generator`
- Check if it's in your PATH: `which version-generator`

**"Permission denied"**
- Make the file executable: `chmod +x version-generator`
- Or run with explicit path: `./version-generator`

**"No such file or directory"**
- Verify download completed: `ls -la version-generator`
- Check file permissions: `ls -la version-generator`

For more help, see [Troubleshooting](Troubleshooting).
