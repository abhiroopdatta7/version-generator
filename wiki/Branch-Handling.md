# Branch Handling

Understanding how version-generator works with different Git branching strategies and workflows.

## 🌳 Branch Strategy Overview

version-generator adapts to different branching strategies by analyzing Git repository state, tags, and commit history. The version generation behavior changes based on the current branch and its relationship to tags and other branches.

## 🎯 Core Branch Concepts

### Tag-Based Versioning
version-generator uses Git tags as the foundation for version generation:
- **Latest Tag**: The most recent Git tag determines the base version
- **Commit Distance**: Number of commits since the latest tag
- **Branch Context**: Current branch affects version format and metadata

### Version Components
Different branch types influence how version components are calculated:
- **Major.Minor.Patch**: Derived from the latest Git tag
- **Pre-release**: Added based on branch type (develop, feature, etc.)
- **Build Metadata**: Includes commit count and optionally commit hash

## 🚀 Main/Master Branch

The main branch (master/main) represents the stable, production-ready code.

### Version Format
```bash
# On main branch at tag
git checkout main
git tag v1.2.3
./version-generator --semver
# Output: v1.2.3

# On main branch with commits after tag
git commit -m "feat: new feature"
./version-generator --semver  
# Output: v1.2.4-alpha.1
```

### Characteristics
- **Clean Versions**: Tags on main produce clean version numbers
- **Alpha Pre-release**: Commits after tags get `-alpha.{count}` suffix
- **Release Ready**: Suitable for production releases

### Use Cases
```bash
# Release tagging
git checkout main
VERSION=$(./version-generator --simple)
git tag $VERSION
git push origin $VERSION

# Production builds
docker build -t myapp:$(./version-generator --simple) .
```

## 🔄 Development Branch

Development branches (develop, dev, next) contain the latest development changes.

### Version Format
```bash
# On develop branch
git checkout develop
./version-generator --semver
# Output: v1.2.4-dev.5

# Calendar versioning on develop
./version-generator --cal-ver
# Output: 2025.08.15
```

### Characteristics
- **Dev Pre-release**: Uses `-dev.{count}` suffix for SemVer
- **Calendar Versioning**: Often suitable for time-based development builds
- **Continuous Integration**: Ideal for automated builds and testing

### CI/CD Integration
```yaml
# GitHub Actions example
- name: Build Development Version
  if: github.ref == 'refs/heads/develop'
  run: |
    VERSION=$(./version-generator --cal-ver)
    docker build -t myapp:dev-$VERSION .
```

## 🌿 Feature Branches

Feature branches contain work on specific features or enhancements.

### Version Format
```bash
# On feature branch
git checkout feature/new-authentication
./version-generator --semver
# Output: v1.2.4-feature.3

# With hash for uniqueness
./version-generator --hash
# Output: v1.2.4-feature.3+abc1234
```

### Characteristics
- **Feature Pre-release**: Uses `-feature.{count}` or branch name in version
- **Hash Identification**: Often includes git hash for uniqueness
- **Temporary Versions**: Not intended for production use

### Development Workflow
```bash
# Feature development
git checkout -b feature/user-dashboard
./version-generator --semver --hash > VERSION.txt

# Build feature version
VERSION=$(./version-generator --semver)
docker build -t myapp:$VERSION .
```

## 🐛 Bugfix/Hotfix Branches

Branches for fixing critical issues or bugs.

### Version Format
```bash
# On hotfix branch
git checkout hotfix/security-fix
./version-generator --semver
# Output: v1.2.4-hotfix.2

# For immediate patch release
./version-generator --simple
# Output: v1.2.4 (if this would be the next patch)
```

### Characteristics
- **Hotfix Pre-release**: Uses `-hotfix.{count}` suffix
- **Patch Candidates**: Often represent the next patch version
- **Urgent Builds**: Suitable for emergency deployments

### Emergency Release Process
```bash
# Create hotfix
git checkout -b hotfix/critical-bug main

# Generate patch version
PATCH_VERSION=$(./version-generator --simple)
echo "Preparing hotfix: $PATCH_VERSION"

# After fix, tag and release
git tag $PATCH_VERSION
gh release create $PATCH_VERSION --title "Hotfix $PATCH_VERSION"
```

## 🔀 Release Branches

Branches preparing for the next release.

### Version Format
```bash
# On release branch
git checkout release/v1.3.0
./version-generator --semver
# Output: v1.3.0-rc.4

# Pre-release candidates
./version-generator --simple
# Output: v1.3.0
```

### Characteristics
- **Release Candidates**: Uses `-rc.{count}` suffix
- **Stabilization**: Versions represent release candidates
- **Final Prep**: Last step before production release

### Release Workflow
```bash
# Create release branch
git checkout -b release/v1.3.0 develop

# Generate release candidate versions
RC_VERSION=$(./version-generator --semver)
echo "Release candidate: $RC_VERSION"

# Test and stabilize
./run-tests.sh
./build-release.sh $RC_VERSION

# Final release
git checkout main
git merge release/v1.3.0
git tag v1.3.0
```

## 🔧 Branch Detection Logic

### Automatic Branch Recognition
version-generator automatically detects branch types based on naming patterns:

```bash
# Branch patterns and their version formats
main|master     → v1.2.3-alpha.N
develop|dev     → v1.2.3-dev.N
feature/*       → v1.2.3-feature.N
hotfix/*        → v1.2.3-hotfix.N
release/*       → v1.2.3-rc.N
bugfix/*        → v1.2.3-bugfix.N
```

### Custom Branch Handling
```bash
#!/bin/bash
# custom-branch-versioning.sh

BRANCH=$(git branch --show-current)
BASE_VERSION=$(./version-generator --simple)

case $BRANCH in
  main|master)
    VERSION=$BASE_VERSION
    TAG="stable"
    ;;
  develop|dev)
    VERSION=$(./version-generator --cal-ver)
    TAG="dev"
    ;;
  feature/*)
    VERSION=$(./version-generator --semver --hash)
    TAG="feature"
    ;;
  hotfix/*)
    VERSION=$(./version-generator --semver)
    TAG="hotfix"
    ;;
  *)
    VERSION=$(./version-generator --hash)
    TAG="experimental"
    ;;
esac

echo "Branch: $BRANCH"
echo "Version: $VERSION" 
echo "Tag: $TAG"
```

## 🔄 GitFlow Integration

### GitFlow Branch Mapping
```bash
#!/bin/bash
# gitflow-versions.sh

BRANCH=$(git branch --show-current)

case $BRANCH in
  master)
    # Production releases
    ./version-generator --simple
    ;;
  develop)
    # Development builds
    ./version-generator --cal-ver
    ;;
  feature/*)
    # Feature development
    ./version-generator --semver --hash
    ;;
  release/*)
    # Release candidates
    ./version-generator --semver
    ;;
  hotfix/*)
    # Hotfix builds
    ./version-generator --semver
    ;;
esac
```

### GitFlow CI/CD
```yaml
name: GitFlow Versioning

on:
  push:
    branches: [ master, develop, 'feature/**', 'release/**', 'hotfix/**' ]

jobs:
  version:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
      with:
        fetch-depth: 0
        
    - name: Generate Version
      id: version
      run: |
        BRANCH=${GITHUB_REF#refs/heads/}
        case $BRANCH in
          master)
            VERSION=$(./version-generator --simple)
            ENVIRONMENT="production"
            ;;
          develop)
            VERSION=$(./version-generator --cal-ver)
            ENVIRONMENT="staging"
            ;;
          feature/*)
            VERSION=$(./version-generator --semver --hash)
            ENVIRONMENT="development"
            ;;
          release/*)
            VERSION=$(./version-generator --semver)
            ENVIRONMENT="preproduction"
            ;;
          hotfix/*)
            VERSION=$(./version-generator --semver)
            ENVIRONMENT="hotfix"
            ;;
        esac
        echo "version=$VERSION" >> $GITHUB_OUTPUT
        echo "environment=$ENVIRONMENT" >> $GITHUB_OUTPUT
```

## 🎋 GitHub Flow Integration

### Simplified Branching
```bash
#!/bin/bash
# github-flow-versions.sh

BRANCH=$(git branch --show-current)

if [ "$BRANCH" = "main" ]; then
    # Main branch: clean versions
    ./version-generator --simple
else
    # Feature branches: development versions
    ./version-generator --semver --hash
fi
```

### Pull Request Builds
```yaml
name: PR Builds

on:
  pull_request:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v3
      with:
        fetch-depth: 0
        
    - name: Generate PR Version
      run: |
        # Use hash for PR uniqueness
        VERSION=$(./version-generator --hash)
        echo "PR Version: $VERSION"
        docker build -t myapp:pr-$VERSION .
```

## 🌐 Multi-Environment Deployment

### Environment-Based Versioning
```bash
#!/bin/bash
# multi-env-deploy.sh

BRANCH=$(git branch --show-current)
VERSION=$(./version-generator --semver)

case $BRANCH in
  main)
    ENVIRONMENT="production"
    REGISTRY="prod-registry"
    ;;
  develop)
    ENVIRONMENT="staging"
    REGISTRY="staging-registry"
    ;;
  *)
    ENVIRONMENT="development"
    REGISTRY="dev-registry"
    ;;
esac

echo "Deploying $VERSION to $ENVIRONMENT"
docker tag myapp:$VERSION $REGISTRY/myapp:$VERSION
docker tag myapp:$VERSION $REGISTRY/myapp:$ENVIRONMENT
docker push $REGISTRY/myapp:$VERSION
docker push $REGISTRY/myapp:$ENVIRONMENT
```

### Kubernetes Deployment
```yaml
# k8s-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: myapp
  labels:
    app: myapp
    version: ${VERSION}
    branch: ${BRANCH}
spec:
  template:
    spec:
      containers:
      - name: myapp
        image: myapp:${VERSION}
        env:
        - name: APP_VERSION
          value: ${VERSION}
        - name: GIT_BRANCH
          value: ${BRANCH}
```

## 🔍 Branch Analysis Tools

### Branch Information Script
```bash
#!/bin/bash
# branch-info.sh

echo "=== Branch Analysis ==="
echo "Current Branch: $(git branch --show-current)"
echo "Latest Tag: $(git describe --tags --abbrev=0 2>/dev/null || echo 'none')"
echo "Commits Since Tag: $(git rev-list --count HEAD ^$(git describe --tags --abbrev=0 2>/dev/null || echo HEAD~999) 2>/dev/null || echo 0)"
echo "Last Commit: $(git log -1 --format='%h %s')"
echo ""
echo "=== Version Generation ==="
echo "Default: $(./version-generator)"
echo "SemVer: $(./version-generator --semver)"
echo "CalVer: $(./version-generator --cal-ver)"
echo "Simple: $(./version-generator --simple)"
echo "Hash: $(./version-generator --hash)"
```

### Version Validation
```bash
#!/bin/bash
# validate-branch-version.sh

BRANCH=$(git branch --show-current)
VERSION=$(./version-generator --semver)

# Validate version format based on branch
case $BRANCH in
  main|master)
    if [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+(-alpha\.[0-9]+)?$ ]]; then
      echo "✅ Valid main branch version: $VERSION"
    else
      echo "❌ Invalid main branch version: $VERSION"
      exit 1
    fi
    ;;
  develop)
    if [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+-dev\.[0-9]+$ ]]; then
      echo "✅ Valid develop branch version: $VERSION"
    else
      echo "❌ Invalid develop branch version: $VERSION"
      exit 1
    fi
    ;;
  feature/*)
    if [[ $VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+-feature\.[0-9]+$ ]]; then
      echo "✅ Valid feature branch version: $VERSION"
    else
      echo "❌ Invalid feature branch version: $VERSION"
      exit 1
    fi
    ;;
esac
```

## 📋 Branch Strategy Recommendations

### Development Team Guidelines

1. **Main Branch**
   - Use simple versioning for releases: `./version-generator --simple`
   - Tag releases immediately after merging
   - Keep main branch stable and deployable

2. **Development Branch**
   - Use calendar versioning for continuous integration: `./version-generator --cal-ver`
   - Deploy automatically to staging environments
   - Run comprehensive test suites

3. **Feature Branches**
   - Use hash-enhanced versions: `./version-generator --hash`
   - Build and test in isolation
   - Clean up after merge

4. **Release Branches**
   - Use semantic versioning: `./version-generator --semver`
   - Generate release candidates
   - Document version changes

### CI/CD Best Practices

1. **Version Consistency**
   - Use the same version generation across all pipeline stages
   - Store versions as build artifacts
   - Tag Docker images with generated versions

2. **Environment Promotion**
   - Use the same version identifier across environments
   - Track version deployments
   - Enable easy rollbacks

3. **Branch Protection**
   - Validate version formats in CI
   - Require version generation tests
   - Enforce tagging policies

version-generator's branch handling provides flexible versioning that adapts to your team's branching strategy while maintaining consistency and traceability across all environments.
