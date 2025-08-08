package versionSchemes

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

// VersioningOptions defines different versioning scheme options
type VersioningOptions struct {
	Semver bool // Use Semantic Versioning: v1.2.3-alpha.4 or v1.2.3-beta.4+branch
	CalVer bool // Use Calendar Versioning: 2024.08.4 or 2024.08.4-branch
	Simple bool // Use simple format: v1.2.3 (no branch/commit info)
	Hash   bool // Include short hash in version
}

// VersionGenerator provides methods to generate version strings using different schemes
type VersionGenerator struct{}

// NewVersionGenerator creates a new version generator
func NewVersionGenerator() *VersionGenerator {
	return &VersionGenerator{}
}

// GenerateVersion generates version string based on the provided options
func (vg *VersionGenerator) GenerateVersion(lastTag string, commitsSince int, shortHash, branchName string, options VersioningOptions) string {
	if commitsSince == 0 && !options.Hash {
		// We're exactly on a tag and no hash requested
		if options.Simple {
			return lastTag
		}
		if options.CalVer {
			return vg.GenerateCalVer(lastTag, 0, branchName, false, shortHash)
		}
		return lastTag
	}

	// Handle different versioning schemes
	switch {
	case options.CalVer:
		return vg.GenerateCalVer(lastTag, commitsSince, branchName, options.Hash, shortHash)
	case options.Semver:
		return vg.GenerateSemVer(lastTag, commitsSince, branchName, options.Hash, shortHash)
	case options.Simple:
		return vg.GenerateSimple(lastTag, shortHash, options.Hash)
	default:
		return vg.GenerateDefault(lastTag, commitsSince, shortHash, branchName, options.Hash)
	}
}

// GenerateLegacy generates version string in legacy format (for backward compatibility)
func (vg *VersionGenerator) GenerateLegacy(lastTag string, commitsSince int, shortHash, branchName string, dockerFormat bool) string {
	if commitsSince == 0 {
		// We're exactly on a tag
		return lastTag
	}

	// For main/master branches, don't include branch name in version
	if vg.isMainBranch(branchName) {
		if dockerFormat {
			// Docker format: <tag>-<count>
			return fmt.Sprintf("%s-%d", lastTag, commitsSince)
		} else {
			// Default format: <tag>+<count>
			return fmt.Sprintf("%s+%d", lastTag, commitsSince)
		}
	}

	// For other branches, clean branch name and include it in version string
	cleanBranch := vg.cleanBranchName(branchName)

	// Choose format based on dockerFormat flag
	if dockerFormat {
		// Docker format: <tag>-<branch>-<count of commit>
		return fmt.Sprintf("%s-%s-%d", lastTag, cleanBranch, commitsSince)
	} else {
		// Default format: <tag>-<branch>+<count of commit>
		return fmt.Sprintf("%s-%s+%d", lastTag, cleanBranch, commitsSince)
	}
}

// GenerateCalVer generates Calendar Versioning format
func (vg *VersionGenerator) GenerateCalVer(lastTag string, commitsSince int, branchName string, includeHash bool, shortHash string) string {
	now := time.Now()
	calVer := fmt.Sprintf("%d.%02d", now.Year(), now.Month())

	if commitsSince > 0 {
		calVer = fmt.Sprintf("%s.%d", calVer, commitsSince)
	}

	if !vg.isMainBranch(branchName) {
		cleanBranch := vg.cleanBranchName(branchName)
		calVer = fmt.Sprintf("%s-%s", calVer, cleanBranch)
	}

	if includeHash && shortHash != "" {
		calVer = fmt.Sprintf("%s+%s", calVer, shortHash)
	}

	return calVer
}

// GenerateSemVer generates Semantic Versioning format
func (vg *VersionGenerator) GenerateSemVer(lastTag string, commitsSince int, branchName string, includeHash bool, shortHash string) string {
	if commitsSince == 0 && !includeHash {
		return lastTag
	}

	// Parse the tag to extract semver parts
	version := lastTag
	if hasVersionPrefix(version) {
		version = version[1:]
	}

	if vg.isMainBranch(branchName) {
		if commitsSince > 0 {
			version = fmt.Sprintf("%s-dev.%d", version, commitsSince)
		}
	} else {
		cleanBranch := vg.cleanBranchName(branchName)
		if commitsSince > 0 {
			version = fmt.Sprintf("%s-%s.%d", version, cleanBranch, commitsSince)
		} else {
			version = fmt.Sprintf("%s-%s", version, cleanBranch)
		}
	}

	if includeHash && shortHash != "" {
		version = fmt.Sprintf("%s+%s", version, shortHash)
	}

	return ensureVersionPrefix(version)
}

// GenerateSimple generates simple version format
func (vg *VersionGenerator) GenerateSimple(lastTag string, shortHash string, includeHash bool) string {
	if includeHash {
		return fmt.Sprintf("%s+%s", lastTag, shortHash)
	}
	return lastTag
}

// GenerateDefault generates default format
func (vg *VersionGenerator) GenerateDefault(lastTag string, commitsSince int, shortHash, branchName string, includeHash bool) string {
	if commitsSince == 0 && !includeHash {
		return lastTag
	}

	version := lastTag
	if vg.isMainBranch(branchName) {
		if commitsSince > 0 {
			version = fmt.Sprintf("%s+%d", lastTag, commitsSince)
		}
	} else {
		cleanBranch := vg.cleanBranchName(branchName)
		if commitsSince > 0 {
			version = fmt.Sprintf("%s-%s+%d", lastTag, cleanBranch, commitsSince)
		} else {
			version = fmt.Sprintf("%s-%s", lastTag, cleanBranch)
		}
	}

	if includeHash {
		version = fmt.Sprintf("%s+%s", version, shortHash)
	}

	return version
}

// Helper functions

func (vg *VersionGenerator) isMainBranch(branchName string) bool {
	return branchName == "main" || branchName == "master" || branchName == "detached"
}

func (vg *VersionGenerator) cleanBranchName(branchName string) string {
	return regexp.MustCompile(`[^a-zA-Z0-9\-]`).ReplaceAllString(branchName, "-")
}

func hasVersionPrefix(version string) bool {
	return strings.HasPrefix(version, "v")
}

func ensureVersionPrefix(version string) string {
	if !hasVersionPrefix(version) {
		return "v" + version
	}
	return version
}
