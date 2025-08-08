# CI/CD Integration

Comprehensive guide for integrating version-generator into Continuous Integration and Continuous Deployment pipelines.

## 🚀 Overview

version-generator is designed for seamless CI/CD integration, providing consistent versioning across build, test, and deployment stages. This guide covers major CI/CD platforms and deployment strategies.

## 🔧 Prerequisites

### Repository Setup
```yaml
# Ensure proper git history for version generation
git:
  fetch-depth: 0  # Important: fetch full history for accurate version calculation
```

### Tool Availability
- version-generator binary in PATH or workspace
- Git repository with proper tag history
- Appropriate permissions for tag creation and release management

## 🎯 GitHub Actions

### Basic Integration
```yaml
name: CI/CD with Versioning

on:
  push:
    branches: [ main, develop ]
  pull_request:
    branches: [ main ]

jobs:
  version-and-build:
    runs-on: ubuntu-latest
    
    steps:
    - name: Checkout
      uses: actions/checkout@v4
      with:
        fetch-depth: 0  # Essential for version generation
        
    - name: Generate Version
      id: version
      run: |
        VERSION=$(./version-generator --semver)
        echo "version=$VERSION" >> $GITHUB_OUTPUT
        echo "Generated version: $VERSION"
        
    - name: Build Application
      run: |
        go build -ldflags="-X main.Version=${{ steps.version.outputs.version }}" -o myapp
        
    - name: Build Docker Image
      run: |
        docker build -t myapp:${{ steps.version.outputs.version }} .
        docker tag myapp:${{ steps.version.outputs.version }} myapp:latest
```

### Multi-Environment Deployment
```yaml
name: Multi-Environment Deploy

on:
  push:
    branches: [ main, develop, 'feature/**' ]

jobs:
  deploy:
    runs-on: ubuntu-latest
    
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
        
    - name: Determine Environment and Version
      id: env
      run: |
        BRANCH=${GITHUB_REF#refs/heads/}
        case $BRANCH in
          main)
            VERSION=$(./version-generator --simple)
            ENVIRONMENT="production"
            ;;
          develop)
            VERSION=$(./version-generator --cal-ver)
            ENVIRONMENT="staging"
            ;;
          feature/*)
            VERSION=$(./version-generator --hash)
            ENVIRONMENT="development"
            ;;
        esac
        echo "version=$VERSION" >> $GITHUB_OUTPUT
        echo "environment=$ENVIRONMENT" >> $GITHUB_OUTPUT
        
    - name: Deploy to Environment
      run: |
        echo "Deploying ${{ steps.env.outputs.version }} to ${{ steps.env.outputs.environment }}"
        # Your deployment commands here
```

### Release Automation
```yaml
name: Automated Release

on:
  push:
    branches: [ main ]

jobs:
  release:
    runs-on: ubuntu-latest
    if: "!contains(github.event.head_commit.message, '[skip-release]')"
    
    steps:
    - uses: actions/checkout@v4
      with:
        fetch-depth: 0
        token: ${{ secrets.GITHUB_TOKEN }}
        
    - name: Generate Release Version
      id: version
      run: |
        VERSION=$(./version-generator --simple)
        echo "version=$VERSION" >> $GITHUB_OUTPUT
        
    - name: Check if Release Needed
      id: check
      run: |
        if git tag -l | grep -q "^${{ steps.version.outputs.version }}$"; then
          echo "needs_release=false" >> $GITHUB_OUTPUT
          echo "Version ${{ steps.version.outputs.version }} already exists"
        else
          echo "needs_release=true" >> $GITHUB_OUTPUT
          echo "Creating new release ${{ steps.version.outputs.version }}"
        fi
        
    - name: Build Release Assets
      if: steps.check.outputs.needs_release == 'true'
      run: |
        VERSION=${{ steps.version.outputs.version }}
        
        # Build for multiple platforms
        GOOS=linux GOARCH=amd64 go build -ldflags="-X main.Version=$VERSION" -o myapp-linux-amd64
        GOOS=windows GOARCH=amd64 go build -ldflags="-X main.Version=$VERSION" -o myapp-windows-amd64.exe
        GOOS=darwin GOARCH=amd64 go build -ldflags="-X main.Version=$VERSION" -o myapp-darwin-amd64
        
        # Create archives
        tar -czf myapp-linux-amd64.tar.gz myapp-linux-amd64
        zip myapp-windows-amd64.zip myapp-windows-amd64.exe
        tar -czf myapp-darwin-amd64.tar.gz myapp-darwin-amd64
        
    - name: Create GitHub Release
      if: steps.check.outputs.needs_release == 'true'
      uses: actions/create-release@v1
      env:
        GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
      with:
        tag_name: ${{ steps.version.outputs.version }}
        release_name: Release ${{ steps.version.outputs.version }}
        draft: false
        prerelease: false
        
    - name: Upload Release Assets
      if: steps.check.outputs.needs_release == 'true'
      # Upload logic here
```

## 🦊 GitLab CI

### Basic Pipeline
```yaml
# .gitlab-ci.yml
stages:
  - version
  - build
  - test
  - deploy

variables:
  GIT_STRATEGY: clone
  GIT_DEPTH: 0  # Full history for version generation

generate-version:
  stage: version
  image: alpine:latest
  before_script:
    - apk add --no-cache git
  script:
    - VERSION=$(./version-generator --semver)
    - echo "VERSION=$VERSION" > version.env
    - echo "Generated version: $VERSION"
  artifacts:
    reports:
      dotenv: version.env
    expire_in: 1 hour

build:
  stage: build
  needs: ["generate-version"]
  script:
    - echo "Building version: $VERSION"
    - go build -ldflags="-X main.Version=$VERSION" -o myapp
  artifacts:
    paths:
      - myapp
    expire_in: 1 hour

deploy-staging:
  stage: deploy
  needs: ["build"]
  only:
    - develop
  script:
    - echo "Deploying $VERSION to staging"
    - docker build -t registry.gitlab.com/myproject/myapp:$VERSION .
    - docker push registry.gitlab.com/myproject/myapp:$VERSION

deploy-production:
  stage: deploy
  needs: ["build"]
  only:
    - main
  script:
    - echo "Deploying $VERSION to production"
    - docker build -t registry.gitlab.com/myproject/myapp:$VERSION .
    - docker tag registry.gitlab.com/myproject/myapp:$VERSION registry.gitlab.com/myproject/myapp:latest
    - docker push registry.gitlab.com/myproject/myapp:$VERSION
    - docker push registry.gitlab.com/myproject/myapp:latest
```

### Advanced GitLab Pipeline
```yaml
# .gitlab-ci.yml
include:
  - template: Security/SAST.gitlab-ci.yml
  - template: Security/Container-Scanning.gitlab-ci.yml

stages:
  - version
  - build
  - test
  - security
  - deploy
  - release

.version-template: &version-template
  image: alpine:latest
  before_script:
    - apk add --no-cache git
  script:
    - |
      case $CI_COMMIT_REF_NAME in
        main)
          VERSION=$(./version-generator --simple)
          ENVIRONMENT="production"
          ;;
        develop)
          VERSION=$(./version-generator --cal-ver)
          ENVIRONMENT="staging"
          ;;
        *)
          VERSION=$(./version-generator --hash)
          ENVIRONMENT="development"
          ;;
      esac
    - echo "VERSION=$VERSION" > version.env
    - echo "ENVIRONMENT=$ENVIRONMENT" >> version.env
  artifacts:
    reports:
      dotenv: version.env

version:
  stage: version
  <<: *version-template

build-docker:
  stage: build
  needs: ["version"]
  services:
    - docker:dind
  script:
    - docker build --build-arg VERSION=$VERSION -t $CI_REGISTRY_IMAGE:$VERSION .
    - docker push $CI_REGISTRY_IMAGE:$VERSION

test-integration:
  stage: test
  needs: ["build-docker"]
  script:
    - docker run --rm $CI_REGISTRY_IMAGE:$VERSION ./run-tests.sh

deploy-k8s:
  stage: deploy
  needs: ["test-integration"]
  script:
    - |
      envsubst < k8s-deployment.yaml | kubectl apply -f -
  environment:
    name: $ENVIRONMENT
    url: https://$ENVIRONMENT.myapp.com

create-release:
  stage: release
  only:
    - main
  needs: ["deploy-k8s"]
  script:
    - |
      curl --request POST \
           --header "PRIVATE-TOKEN: $CI_JOB_TOKEN" \
           --data "tag_name=$VERSION&ref=main&name=Release $VERSION" \
           "$CI_API_V4_URL/projects/$CI_PROJECT_ID/releases"
```

## 🔵 Azure DevOps

### Azure Pipelines YAML
```yaml
# azure-pipelines.yml
trigger:
  branches:
    include:
    - main
    - develop
    - feature/*

pool:
  vmImage: 'ubuntu-latest'

variables:
  GO_VERSION: '1.21'

stages:
- stage: Version
  displayName: 'Generate Version'
  jobs:
  - job: GenerateVersion
    displayName: 'Generate Version Number'
    steps:
    - checkout: self
      fetchDepth: 0  # Full history
      
    - script: |
        VERSION=$(./version-generator --semver)
        echo "##vso[task.setvariable variable=appVersion;isOutput=true]$VERSION"
        echo "Generated version: $VERSION"
      name: version
      displayName: 'Generate Application Version'

- stage: Build
  displayName: 'Build Application'
  dependsOn: Version
  variables:
    appVersion: $[ stageDependencies.Version.GenerateVersion.outputs['version.appVersion'] ]
  jobs:
  - job: Build
    displayName: 'Build and Test'
    steps:
    - task: GoTool@0
      inputs:
        version: $(GO_VERSION)
        
    - script: |
        go build -ldflags="-X main.Version=$(appVersion)" -o myapp
      displayName: 'Build Application'
      
    - task: Docker@2
      inputs:
        command: 'buildAndPush'
        repository: 'myapp'
        tags: |
          $(appVersion)
          latest

- stage: Deploy
  displayName: 'Deploy Application'
  dependsOn: Build
  condition: and(succeeded(), eq(variables['Build.SourceBranch'], 'refs/heads/main'))
  variables:
    appVersion: $[ stageDependencies.Version.GenerateVersion.outputs['version.appVersion'] ]
  jobs:
  - deployment: DeployProduction
    displayName: 'Deploy to Production'
    environment: 'production'
    strategy:
      runOnce:
        deploy:
          steps:
          - script: |
              echo "Deploying version $(appVersion) to production"
              # Deployment commands here
```

### Azure DevOps with Environments
```yaml
# azure-pipelines-environments.yml
resources:
  repositories:
  - repository: self
    fetchDepth: 0

variables:
- group: production-vars
- name: appVersion
  value: ''

stages:
- stage: GenerateVersion
  jobs:
  - job: Version
    steps:
    - script: |
        BRANCH_NAME=$(echo "$(Build.SourceBranch)" | sed 's/refs\/heads\///')
        case $BRANCH_NAME in
          main)
            VERSION=$(./version-generator --simple)
            ENVIRONMENT="production"
            ;;
          develop)
            VERSION=$(./version-generator --cal-ver)
            ENVIRONMENT="staging"
            ;;
          *)
            VERSION=$(./version-generator --hash)
            ENVIRONMENT="development"
            ;;
        esac
        echo "##vso[task.setvariable variable=appVersion;isOutput=true]$VERSION"
        echo "##vso[task.setvariable variable=targetEnvironment;isOutput=true]$ENVIRONMENT"
      name: setVersion
      
- stage: DeployToEnvironment
  dependsOn: GenerateVersion
  variables:
    appVersion: $[ stageDependencies.GenerateVersion.Version.outputs['setVersion.appVersion'] ]
    targetEnvironment: $[ stageDependencies.GenerateVersion.Version.outputs['setVersion.targetEnvironment'] ]
  jobs:
  - deployment: Deploy
    environment: $(targetEnvironment)
    strategy:
      runOnce:
        deploy:
          steps:
          - script: |
              echo "Deploying $(appVersion) to $(targetEnvironment)"
```

## ⚡ Jenkins

### Declarative Pipeline
```groovy
// Jenkinsfile
pipeline {
    agent any
    
    environment {
        GO_VERSION = '1.21'
        DOCKER_REGISTRY = 'myregistry.com'
    }
    
    stages {
        stage('Checkout') {
            steps {
                checkout([
                    $class: 'GitSCM',
                    branches: [[name: '*/main']],
                    extensions: [[$class: 'CloneOption', depth: 0, noTags: false, reference: '', shallow: false]],
                    userRemoteConfigs: [[url: 'https://github.com/myorg/myrepo.git']]
                ])
            }
        }
        
        stage('Generate Version') {
            steps {
                script {
                    def branchName = env.GIT_BRANCH.replaceAll('origin/', '')
                    def versionCommand
                    def environment
                    
                    switch(branchName) {
                        case 'main':
                            versionCommand = './version-generator --simple'
                            environment = 'production'
                            break
                        case 'develop':
                            versionCommand = './version-generator --cal-ver'
                            environment = 'staging'
                            break
                        default:
                            versionCommand = './version-generator --hash'
                            environment = 'development'
                    }
                    
                    env.APP_VERSION = sh(script: versionCommand, returnStdout: true).trim()
                    env.TARGET_ENVIRONMENT = environment
                    
                    echo "Generated version: ${env.APP_VERSION}"
                    echo "Target environment: ${env.TARGET_ENVIRONMENT}"
                }
            }
        }
        
        stage('Build') {
            parallel {
                stage('Build Binary') {
                    steps {
                        sh """
                            go build -ldflags="-X main.Version=${env.APP_VERSION}" -o myapp
                        """
                    }
                }
                
                stage('Build Docker') {
                    steps {
                        script {
                            def image = docker.build("${env.DOCKER_REGISTRY}/myapp:${env.APP_VERSION}")
                            docker.withRegistry("https://${env.DOCKER_REGISTRY}") {
                                image.push()
                                image.push('latest')
                            }
                        }
                    }
                }
            }
        }
        
        stage('Test') {
            steps {
                sh './run-tests.sh'
            }
        }
        
        stage('Deploy') {
            when {
                anyOf {
                    branch 'main'
                    branch 'develop'
                }
            }
            steps {
                script {
                    if (env.TARGET_ENVIRONMENT == 'production') {
                        input message: "Deploy ${env.APP_VERSION} to production?", ok: 'Deploy'
                    }
                    
                    sh """
                        echo "Deploying ${env.APP_VERSION} to ${env.TARGET_ENVIRONMENT}"
                        helm upgrade --install myapp ./helm-chart \\
                            --set image.tag=${env.APP_VERSION} \\
                            --set environment=${env.TARGET_ENVIRONMENT}
                    """
                }
            }
        }
        
        stage('Create Release') {
            when {
                branch 'main'
            }
            steps {
                script {
                    // Check if release already exists
                    def releaseExists = sh(
                        script: "git tag -l | grep -q '^${env.APP_VERSION}\$'",
                        returnStatus: true
                    ) == 0
                    
                    if (!releaseExists) {
                        sh """
                            git tag ${env.APP_VERSION}
                            git push origin ${env.APP_VERSION}
                        """
                        
                        // Create GitHub release if using GitHub
                        withCredentials([string(credentialsId: 'github-token', variable: 'GITHUB_TOKEN')]) {
                            sh """
                                gh release create ${env.APP_VERSION} \\
                                    --title "Release ${env.APP_VERSION}" \\
                                    --generate-notes \\
                                    myapp
                            """
                        }
                    }
                }
            }
        }
    }
    
    post {
        always {
            cleanWs()
        }
        success {
            echo "Pipeline completed successfully for version ${env.APP_VERSION}"
        }
        failure {
            echo "Pipeline failed for version ${env.APP_VERSION}"
        }
    }
}
```

### Scripted Pipeline
```groovy
// Jenkinsfile (Scripted)
node {
    try {
        stage('Checkout') {
            checkout([
                $class: 'GitSCM',
                branches: [[name: '*/main']],
                extensions: [[$class: 'CloneOption', depth: 0, noTags: false]],
                userRemoteConfigs: [[url: 'https://github.com/myorg/myrepo.git']]
            ])
        }
        
        stage('Version') {
            script {
                def version = sh(script: './version-generator --semver', returnStdout: true).trim()
                env.APP_VERSION = version
                currentBuild.displayName = "#${env.BUILD_NUMBER} - ${version}"
            }
        }
        
        stage('Build & Deploy') {
            parallel(
                'Build': {
                    sh "go build -ldflags='-X main.Version=${env.APP_VERSION}' -o myapp"
                },
                'Docker': {
                    def image = docker.build("myapp:${env.APP_VERSION}")
                    image.push()
                }
            )
        }
        
    } catch (Exception e) {
        currentBuild.result = 'FAILURE'
        throw e
    } finally {
        cleanWs()
    }
}
```

## 🏗️ Buildkite

### Basic Pipeline
```yaml
# .buildkite/pipeline.yml
steps:
  - label: ":git: Generate Version"
    key: "version"
    command: |
      VERSION=$$(./version-generator --semver)
      echo "Generated version: $$VERSION"
      buildkite-agent meta-data set "app-version" "$$VERSION"
    plugins:
      - docker#v3.8.0:
          image: alpine/git

  - label: ":hammer: Build"
    depends_on: "version"
    command: |
      VERSION=$$(buildkite-agent meta-data get "app-version")
      go build -ldflags="-X main.Version=$$VERSION" -o myapp
    plugins:
      - docker#v3.8.0:
          image: golang:1.21

  - label: ":docker: Docker Build"
    depends_on: "version"
    command: |
      VERSION=$$(buildkite-agent meta-data get "app-version")
      docker build -t myapp:$$VERSION .
      docker push myapp:$$VERSION

  - label: ":rocket: Deploy"
    depends_on: ["build", "docker-build"]
    branches: "main"
    command: |
      VERSION=$$(buildkite-agent meta-data get "app-version")
      echo "Deploying $$VERSION to production"
      # Deployment commands here
```

## 🌊 CircleCI

### Configuration
```yaml
# .circleci/config.yml
version: 2.1

orbs:
  go: circleci/go@1.7.1
  docker: circleci/docker@2.1.4

workflows:
  version: 2
  build-test-deploy:
    jobs:
      - generate-version
      - build:
          requires:
            - generate-version
      - test:
          requires:
            - build
      - deploy:
          requires:
            - test
          filters:
            branches:
              only: main

jobs:
  generate-version:
    docker:
      - image: alpine/git
    steps:
      - checkout
      - run:
          name: Generate Version
          command: |
            VERSION=$(./version-generator --semver)
            echo "export APP_VERSION=$VERSION" >> $BASH_ENV
            echo "Generated version: $VERSION"
      - persist_to_workspace:
          root: .
          paths:
            - .

  build:
    docker:
      - image: golang:1.21
    steps:
      - attach_workspace:
          at: .
      - run:
          name: Build Application
          command: |
            source $BASH_ENV
            go build -ldflags="-X main.Version=$APP_VERSION" -o myapp
      - store_artifacts:
          path: myapp

  test:
    docker:
      - image: golang:1.21
    steps:
      - attach_workspace:
          at: .
      - run:
          name: Run Tests
          command: go test ./...

  deploy:
    docker:
      - image: cimg/base:stable
    steps:
      - attach_workspace:
          at: .
      - setup_remote_docker
      - run:
          name: Deploy Application
          command: |
            source $BASH_ENV
            docker build -t myapp:$APP_VERSION .
            echo "Deploying $APP_VERSION to production"
```

## 🐳 Docker Integration

### Multi-stage Dockerfile with Versioning
```dockerfile
# Multi-stage build with version
FROM golang:1.21 AS builder

WORKDIR /app
COPY . .

# Generate version during build
RUN ./version-generator --semver > /tmp/version

# Build with embedded version
RUN VERSION=$(cat /tmp/version) && \
    go build -ldflags="-X main.Version=$VERSION" -o myapp

# Runtime stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/myapp .
COPY --from=builder /tmp/version ./VERSION

CMD ["./myapp"]
```

### Docker Compose with CI
```yaml
# docker-compose.ci.yml
version: '3.8'

services:
  app:
    build:
      context: .
      args:
        VERSION: ${CI_VERSION}
    image: myapp:${CI_VERSION}
    environment:
      - APP_VERSION=${CI_VERSION}
      - ENVIRONMENT=${CI_ENVIRONMENT}
```

```bash
#!/bin/bash
# ci-docker-build.sh

VERSION=$(./version-generator --semver)
ENVIRONMENT=${CI_ENVIRONMENT:-development}

export CI_VERSION=$VERSION
export CI_ENVIRONMENT=$ENVIRONMENT

docker-compose -f docker-compose.ci.yml build
docker-compose -f docker-compose.ci.yml push
```

## 🎯 Best Practices

### Version Consistency
```bash
#!/bin/bash
# ensure-version-consistency.sh

# Generate version once and reuse
VERSION=$(./version-generator --semver)
echo "BUILD_VERSION=$VERSION" > .env

# Use in all subsequent steps
source .env
echo "Building version: $BUILD_VERSION"
```

### Caching Strategy
```yaml
# GitHub Actions caching example
- name: Cache version-generator
  uses: actions/cache@v3
  with:
    path: ./version-generator
    key: version-generator-${{ runner.os }}-${{ hashFiles('go.sum') }}
```

### Error Handling
```bash
#!/bin/bash
# robust-version-generation.sh

set -euo pipefail

# Validate git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo "Error: Not a git repository"
    exit 1
fi

# Generate version with fallback
VERSION=$(./version-generator --semver 2>/dev/null) || {
    echo "Warning: version-generator failed, using fallback"
    VERSION="v0.0.0-unknown"
}

echo "Using version: $VERSION"
```

### Security Considerations
```yaml
# Secure token handling
- name: Create Release
  env:
    GITHUB_TOKEN: ${{ secrets.GITHUB_TOKEN }}
  run: |
    VERSION=$(./version-generator --simple)
    gh release create $VERSION --generate-notes
```

## 📊 Monitoring and Observability

### Version Tracking
```yaml
# Add version labels to deployments
metadata:
  labels:
    app.version: "${APP_VERSION}"
    deployment.timestamp: "${BUILD_TIMESTAMP}"
```

### Build Metrics
```bash
#!/bin/bash
# collect-build-metrics.sh

VERSION=$(./version-generator --semver)
BUILD_TIME=$(date -u +%Y-%m-%dT%H:%M:%SZ)

# Send metrics to monitoring system
curl -X POST "https://metrics.company.com/builds" \
  -H "Content-Type: application/json" \
  -d "{
    \"version\": \"$VERSION\",
    \"build_time\": \"$BUILD_TIME\",
    \"branch\": \"$CI_BRANCH\",
    \"environment\": \"$CI_ENVIRONMENT\"
  }"
```

CI/CD integration with version-generator enables consistent, automated versioning across your entire deployment pipeline, from development through production.
