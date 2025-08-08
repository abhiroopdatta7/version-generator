package gitType

import (
	"version-generator/versionSchemes"
)

// BaseGitHandler provides common functionality for git handlers
type BaseGitHandler struct {
	versionGenerator *versionSchemes.VersionGenerator
}

// NewBaseGitHandler creates a new base git handler
func NewBaseGitHandler() *BaseGitHandler {
	return &BaseGitHandler{
		versionGenerator: versionSchemes.NewVersionGenerator(),
	}
}

// GenerateVersionInfoFromComponents creates VersionInfo from git components
func (b *BaseGitHandler) GenerateVersionInfoFromComponents(branchName, shortHash, lastTag string, commitsSince int, dockerFormat bool) *VersionInfo {
	// Generate version string using legacy format for backward compatibility
	version := b.versionGenerator.GenerateLegacy(lastTag, commitsSince, shortHash, branchName, dockerFormat)

	return &VersionInfo{
		Branch:       branchName,
		LastTag:      lastTag,
		CommitsSince: commitsSince,
		ShortHash:    shortHash,
		Version:      version,
	}
}

// GenerateVersionInfoFromComponentsWithOptions creates VersionInfo with custom options
func (b *BaseGitHandler) GenerateVersionInfoFromComponentsWithOptions(branchName, shortHash, lastTag string, commitsSince int, options versionSchemes.VersioningOptions) *VersionInfo {
	// Generate version string using new options
	version := b.versionGenerator.GenerateVersion(lastTag, commitsSince, shortHash, branchName, options)

	return &VersionInfo{
		Branch:       branchName,
		LastTag:      lastTag,
		CommitsSince: commitsSince,
		ShortHash:    shortHash,
		Version:      version,
	}
}
