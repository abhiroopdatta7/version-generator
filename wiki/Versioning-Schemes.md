# Versioning Schemes

Version-generator supports multiple versioning schemes to fit different project needs and conventions.

## 🎯 Overview

| Scheme | Format | Use Case | Example |
|--------|--------|----------|---------|
| **Default** | Tag-based with commits | General purpose | `v1.2.3+5` |
| **Semantic** | [SemVer](https://semver.org/) compatible | Libraries, APIs | `v1.2.3-dev.5` |
| **Calendar** | [CalVer](https://calver.org/) based | Date-driven releases | `2025.08.5` |
| **Simple** | Clean tag only | Simple projects | `v1.2.3` |

## 📋 Default Versioning

The default scheme provides a balance between informativeness and simplicity.

### Format Rules
- **On exact tag**: `{TAG}` (e.g., `v1.2.3`)
- **Main branch**: `{TAG}+{COUNT}` (e.g., `v1.2.3+5`)
- **Feature branch**: `{TAG}-{BRANCH}+{COUNT}` (e.g., `v1.2.3-feature+3`)
- **Docker format**: Uses `-` instead of `+` for compatibility

### Examples
```bash
# Basic usage
./version-generator
```

**Sample outputs:**
```
v1.2.3          # On tag v1.2.3
v1.2.3+5        # 5 commits after v1.2.3 on main
v1.2.3-feat+2   # 2 commits on feature branch
```

## 🔄 Semantic Versioning (--semver)

Follows [Semantic Versioning 2.0.0](https://semver.org/) specification for maximum compatibility.

### Format Rules
- **On exact tag**: `v{MAJOR}.{MINOR}.{PATCH}`
- **Main branch**: `v{MAJOR}.{MINOR}.{PATCH}-dev.{COUNT}`
- **Feature branch**: `v{MAJOR}.{MINOR}.{PATCH}-{BRANCH}.{COUNT}`
- **With hash**: Appends `+{HASH}`

### Examples
```bash
./version-generator --semver
```

**Sample outputs:**
```
v1.2.3                    # On tag v1.2.3
v1.2.3-dev.5             # 5 commits after tag on main
v1.2.3-feature-auth.3    # 3 commits on feature/auth branch
v1.2.3-dev.5+abc1234     # With git hash
```

### Pre-release Identifiers
- `dev` - Development builds on main branch
- `{branch-name}` - Feature branch builds
- Numbers indicate commit count

### Compatibility
✅ **Compatible with:**
- npm/yarn package managers
- Go module versioning
- Docker image tags
- Most package managers

## 📅 Calendar Versioning (--cal-ver)

Uses current date for version numbers, following [Calendar Versioning](https://calver.org/) patterns.

### Format Rules
- **Base format**: `{YEAR}.{MONTH}`
- **With commits**: `{YEAR}.{MONTH}.{COUNT}`
- **Feature branch**: `{YEAR}.{MONTH}.{COUNT}-{BRANCH}`
- **With hash**: Appends `+{HASH}`

### Examples
```bash
./version-generator --cal-ver
```

**Sample outputs:**
```
2025.08              # Current month, no commits
2025.08.5            # 5 commits this month
2025.08.5-feature    # On feature branch
2025.08.5+abc1234    # With git hash
```

### Use Cases
- **Regular releases**: Monthly or time-based releases
- **Date tracking**: When release date matters
- **Marketing versions**: Easy to communicate versions
- **Long-term projects**: Clear temporal progression

### Benefits
- **Intuitive**: Easy to understand when version was created
- **No conflicts**: Versions naturally increment with time
- **Marketing friendly**: "2025.08 release" is clear

## 🎯 Simple Versioning (--simple)

Minimal versioning that only shows the tag without additional information.

### Format Rules
- **Always**: `{TAG}` (e.g., `v1.2.3`)
- **With hash**: `{TAG}+{HASH}` (e.g., `v1.2.3+abc1234`)

### Examples
```bash
./version-generator --simple
```

**Sample outputs:**
```
v1.2.3           # Clean tag version
v1.2.3+abc1234   # With hash flag
```

### Use Cases
- **Release builds**: When you only want clean version tags
- **Documentation**: Simple version display
- **Marketing**: Clean version numbers for user-facing displays

## 🏷️ Hash Integration (--hash)

Any scheme can include the git commit hash for precise identification.

### Examples
```bash
./version-generator --hash
./version-generator --semver --hash
./version-generator --cal-ver --hash
./version-generator --simple --hash
```

**Sample outputs:**
```
v1.2.3+abc1234              # Default with hash
v1.2.3-dev.5+abc1234        # SemVer with hash
2025.08.5+abc1234           # CalVer with hash
v1.2.3+abc1234              # Simple with hash
```

### Use Cases
- **Debug builds**: Precise commit identification
- **Troubleshooting**: Link version to exact commit
- **Development**: Track specific builds

## 🌿 Branch Handling

All schemes handle branches intelligently based on branch name.

### Main Branches
Branches `main`, `master`, and `detached` are treated as primary branches:
- **Clean versions**: No branch name in version
- **Dev indicators**: Use `-dev` suffix in SemVer

### Feature Branches
All other branches include sanitized branch name:
- **Sanitization**: Replace invalid characters with `-`
- **Examples**: `feature/auth` → `feature-auth`

### Branch Sanitization Rules
```
feature/auth-system  →  feature-auth-system
bug/fix-#123        →  bug-fix-123
release/v1.2        →  release-v1-2
user@domain         →  user-domain
```

## 🔄 Scheme Comparison

### When to Use Each Scheme

| Scenario | Recommended Scheme | Reason |
|----------|-------------------|---------|
| **Library/API** | Semantic | SemVer compatibility |
| **Web Application** | Default | Balanced information |
| **Mobile App** | Calendar | Marketing-friendly |
| **CLI Tool** | Semantic | Package manager compatibility |
| **Internal Tool** | Simple | Clean, minimal |
| **Debug Build** | Any + Hash | Precise tracking |

### Format Comparison

| Input State | Default | Semantic | Calendar | Simple |
|-------------|---------|----------|-----------|---------|
| On tag v1.2.3 | `v1.2.3` | `v1.2.3` | `2025.08` | `v1.2.3` |
| 5 commits after | `v1.2.3+5` | `v1.2.3-dev.5` | `2025.08.5` | `v1.2.3` |
| Feature branch | `v1.2.3-feat+3` | `v1.2.3-feat.3` | `2025.08.3-feat` | `v1.2.3` |
| With hash | `v1.2.3+5+abc123` | `v1.2.3-dev.5+abc123` | `2025.08.5+abc123` | `v1.2.3+abc123` |

## 🚀 Advanced Usage

### Conditional Schemes
```bash
# Use different schemes based on branch
if [[ $(git branch --show-current) == "main" ]]; then
    VERSION=$(./version-generator --semver)
else
    VERSION=$(./version-generator --cal-ver)
fi
```

### Environment-based Schemes
```bash
# Production vs development
if [[ "$ENV" == "prod" ]]; then
    VERSION=$(./version-generator --simple)
else
    VERSION=$(./version-generator --semver --hash)
fi
```

### Multi-format Generation
```bash
# Generate multiple formats
SEMVER=$(./version-generator --semver)
CALVER=$(./version-generator --cal-ver) 
SIMPLE=$(./version-generator --simple)

echo "SemVer: $SEMVER"
echo "CalVer: $CALVER"
echo "Simple: $SIMPLE"
```

## 🔧 Integration Examples

### Package.json (Node.js)
```bash
# Use SemVer without 'v' prefix
VERSION=$(./version-generator --semver | sed 's/v//')
npm version $VERSION --no-git-tag-version
```

### Docker Tags
```bash
# Use CalVer for date-based releases
docker build -t myapp:$(./version-generator --cal-ver) .
```

### Go Build
```bash
# Embed SemVer in binary
go build -ldflags="-X main.Version=$(./version-generator --semver)"
```

### CMake (C++)
```bash
# Generate C++ header
./version-generator --cpp --semver
# Use in CMakeLists.txt
```
