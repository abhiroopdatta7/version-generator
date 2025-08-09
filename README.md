# Version Generator

A Go application that generates semantic version numbers based on Git repository state. Supports both built-in go-git library and system git executable.

## Features

- **Dual Git Backend**: Choose between built-in go-git library or system git executable
- **Branch Detection**: Identifies current Git branch
- **Tag Analysis**: Finds the most recent reachable tag
- **Rebase-Aware Tag Discovery**: For feature branches, finds tags from the rebase point (common ancestor with main/master)
- **Commit Counting**: Counts commits since the last tag
- **Version Generation**: Creates semantic version strings with pre-release information
- **Multiple Output Formats**: Support for Go, C++, YAML, and plain text files
- **Version Formats**: Support for default and Docker-compatible version formats
- **File Output**: Write generated versions to files in various formats for CI/CD integration
- **Modern CLI**: Uses Kong for clean command line argument parsing
- **Path Support**: Flexible file path and directory support for output files

## Installation

```bash
go mod tidy
go build -o version-generator .
```

## Usage

Run the program from within a Git repository:

```bash
./version-generator [flags]
```

### Command Line Options

```
Flags:
  -h, --help              Show context-sensitive help.
    --semver                Use Semantic Versioning format
    --calver                Use Calendar Versioning format
    --simple                Use simple version format (no branch info)
    --hash                  Include short hash in version
  -i, --in-built-git      Use built-in go-git library instead of system git
  -g, --go                Generate Go format version file
      --go-path=PATH      Path for Go file (default: version.go)
  -c, --cpp               Generate C++ format version file
      --cpp-path=PATH     Path for C++ file (default: version.h)
  -y, --yaml              Generate YAML format version file
      --yaml-path=PATH    Path for YAML file (default: version.yaml)
  -f, --file              Write version to file
      --file-path=PATH    Path for file (default: .VERSION)
```

### Git Backend Options

- **System Git (default)**: Uses system git executable via command line
- **Built-in Go-Git**: Uses pure Go implementation with `-i/--in-built-git` flag

### Examples

```bash
# Print version to console (default, uses system git)
./version-generator

# Print version using built-in go-git library
./version-generator -i

# Generate Go source file with version constant
./version-generator -g

# Generate Go file with custom path
./version-generator -g --go-path=src/version.go

# Generate C++ header file
./version-generator -c --cpp-path=include/version.h

# Generate YAML configuration file
./version-generator -y --yaml-path=config/version.yaml

# Write plain text version to .VERSION file
./version-generator -f

# Write to custom file with specific path
./version-generator -f --file-path=build/VERSION.txt



# Generate in subdirectory (will create directories as needed)
./version-generator -g --go-path=build/generated/version.go
```

### Example Output

**Console Output (default):**
```
v1.2.3-feature-new-api+5
```

**Go File Output (`-g`):**
```go
package main

const Version = "v1.2.3-feature-new-api+5"
```

**C++ File Output (`-c`):**
```cpp
#define VERSION "v1.2.3-feature-new-api+5"
```

**YAML File Output (`-y`):**
```yaml
version: v1.2.3-feature-new-api+5
```

**Plain Text File Output (`-f`):**
```
v1.2.3-feature-new-api+5
```

## Version Generation Logic

### On a Tag
If you're exactly on a tag commit:
```
Generated Version: v1.2.3
```

### Development on Main/Master Branch
For commits after a tag on main/master:
```
Semver/Default: v1.2.3+5
```

### Feature Branch
For commits on other branches:
```
Semver/Default: v1.2.3-feature-branch+5
```

### Feature Branch with Rebase Logic
When on a feature branch, the tool finds the common ancestor with main/master and uses tags from that point:
```
# If rebased from v1.2.3 and has 3 commits since rebase
Generated Version: v1.2.3-feature-branch+3
```

### No Tags Found
If no tags exist in the repository:
```
Semver/Default: v0.0.0+10
```

## Version Format

The generated version follows this pattern:


### Version Format
```
<tag>-<branch>+<count>
```

Where:
- `tag`: The most recent reachable Git tag (or from rebase point for feature branches)
- `branch`: Current branch name (cleaned for version compatibility)
    - For main/master branches, branch name is omitted from version
    - Special characters are replaced with hyphens
- `count`: Number of commits since the last tag

## Git Backend Architecture

The application uses a modular git interface system with two implementations:

### System Git Handler (Default)
- Uses system git executable via command line
- Requires git to be installed and available in PATH
- Fast and reliable for most use cases
- Leverages native git performance

### Built-in Go-Git Handler
- Uses pure Go implementation (go-git library)
- No external dependencies on system git
- Cross-platform compatibility
- Useful in containerized environments without git installed

Both implementations provide identical functionality and produce the same results.

## How It Works

1. **Git Backend Selection**: Chooses between system git or built-in go-git based on `-i` flag
2. **Repository Analysis**: Opens the current directory as a Git repository
3. **Branch Detection**: Identifies the current branch or detached HEAD state
4. **Rebase-Aware Tag Discovery**: 
   - For main/master: Finds all tags reachable from current commit
   - For feature branches: Finds common ancestor with main/master, then finds tags from that point
5. **Tag Selection**: Selects the most recent tag based on commit timestamp
6. **Commit Counting**: Counts commits between current HEAD and the selected tag
7. **Version Assembly**: Constructs a semantic version string based on the collected information
8. **Output Generation**: Formats and writes version to console or files in specified format

## Output File Types

The application supports multiple output formats through a modular file type system:

### Go Source Files (`-g`)
Generates Go source files with version constant:
```go
package main

const Version = "v1.2.3+5"
```

### C++ Header Files (`-c`)
Generates C++ header files with version define:
```cpp
#define VERSION "v1.2.3+5"
```

### YAML Configuration Files (`-y`)
Generates YAML files with version key:
```yaml
version: v1.2.3+5
```

### Plain Text Files (`-f`)
Generates simple text files with version string:
```
v1.2.3+5
```

### Path and Directory Support
- All file types support custom paths with `--{type}-path` flags
- Directories are created automatically if they don't exist
- Files are overwritten if they already exist
- Supports both relative and absolute paths

## Use Cases

- **CI/CD Pipelines**: Generate build versions automatically and write to files in various formats
- **Docker Images**: Use Docker format for container tagging compatibility
- **Release Management**: Create consistent version numbers across environments
- **Development Builds**: Track development progress with meaningful version strings
- **Multi-branch Development**: Distinguish between different feature branches with proper rebase handling
- **Automated Deployments**: Write versions to files for deployment scripts
- **Source Code Integration**: Generate version constants directly in Go or C++ source files
- **Configuration Management**: Export version information to YAML configuration files
- **Cross-Platform Builds**: Use built-in git for environments without system git installed

## File Output Integration

The tool can write version strings to files in multiple formats for integration with build systems:

### Build System Examples

```bash
# Generate Go version file for embedding in builds
./version-generator -g --go-path=internal/version/version.go

# Generate C++ header for native applications
./version-generator -c --cpp-path=src/version.h

# Generate YAML for Kubernetes deployments
./version-generator -y --yaml-path=k8s/version.yaml



# Multi-format generation for complex projects
./version-generator -g --go-path=pkg/version.go
./version-generator -c --cpp-path=include/version.hpp
./version-generator -y --yaml-path=config/app-version.yaml
```

### CI/CD Integration Examples

```bash
# GitHub Actions / GitLab CI
- name: Generate Version
  run: |
    ./version-generator -g --go-path=internal/version.go
    ./version-generator -f --file-path=VERSION.txt
    echo "APP_VERSION=$(cat VERSION.txt)" >> $GITHUB_ENV

# Docker multi-stage builds
COPY version-generator /usr/local/bin/
RUN version-generator -f --file-path=/tmp/VERSION
RUN docker build -t myapp:$(cat /tmp/VERSION) .

# Makefile integration
version:
	./version-generator -g --go-path=pkg/version/version.go
	./version-generator -f --file-path=VERSION

build: version
	go build -ldflags="-X main.Version=$(cat VERSION)" ./cmd/myapp
```

## Dependencies

### Core Dependencies
- **[go-git/go-git](https://github.com/go-git/go-git)**: Pure Go implementation of Git (for built-in git backend)
- **[alecthomas/kong](https://github.com/alecthomas/kong)**: Modern command line parser with struct tags
- **[gopkg.in/yaml.v3](https://gopkg.in/yaml.v3)**: YAML processing for YAML file output

### Standard Library
- Go standard library: `fmt`, `log`, `os`, `path/filepath`, `regexp`, `sort`, `strings`, `strconv`, `os/exec`

### System Requirements
- **For system git backend**: Git executable installed and available in PATH
- **For built-in git backend**: No external dependencies required

## Command Line Interface

The tool uses the Kong library for modern command line parsing, providing:
- Struct-tag based flag definitions
- Automatic help generation
- POSIX-style flag parsing
- Boolean and string flag support
- Path placeholders and default value documentation
- Subcommand support (extensible for future features)

## Architecture

### Modular Design
- **Git Interface**: Pluggable git backend system (`gitType` package)
- **File Type System**: Extensible file format handlers (`fileType` package)
- **Kong CLI**: Modern argument parsing with comprehensive help
- **Clean Separation**: Business logic separated from implementation details

### Package Structure
```
version-generator/
├── main.go                 # Main application and CLI handling
├── gitType/               # Git backend implementations
│   ├── git_interface.go   # Git handler interface
│   ├── gogit_handler.go   # Built-in go-git implementation
│   └── systemgit_handler.go # System git implementation
└── fileType/              # File format implementations
    ├── filetype.go        # File type interface
    ├── basic.go           # Plain text files
    ├── golang.go          # Go source files
    ├── cpp.go             # C++ header files
    └── yaml.go            # YAML configuration files
```

## Error Handling

The application will exit with an error if:
- Not run from within a Git repository
- Git repository is corrupted or inaccessible
- Unable to read Git objects or references
- Unable to write to the specified output file or create required directories
- Git executable not found (when using system git backend)
- Invalid command line arguments or flag combinations

## Performance

### System Git Backend
- Fast execution using native git commands
- Minimal memory usage
- Leverages git's optimized algorithms
- Requires git installation

### Built-in Go-Git Backend  
- Pure Go implementation
- Higher memory usage for large repositories
- No external dependencies
- Slightly slower for very large repositories
- Better for containerized or restricted environments

## Extending the Application

### Adding New File Types
Create a new file in `fileType/` implementing the `FileType` interface:
```go
type FileType interface {
    WriteVersion(filePath string, version string) error
}
```

### Adding New Git Backends
Implement the `GitHandler` interface in `gitType/`:
```go
type GitHandler interface {
    GenerateVersionInfo(dockerFormat bool) (*VersionInfo, error)
    GetCurrentBranch() (string, error)
    GetLastTag(branchName string) (string, error)
    GetCommitsSinceTag(tagName string) (int, error)
    GetShortHash() (string, error)
}
```

## License

This project is provided as-is for educational and development purposes.
