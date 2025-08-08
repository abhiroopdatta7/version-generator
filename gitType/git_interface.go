package gitType

import "version-generator/versionSchemes"

// VersionInfo contains git version information
type VersionInfo struct {
	Branch       string
	LastTag      string
	CommitsSince int
	ShortHash    string
	Version      string
}

// VersioningOptions defines different versioning scheme options
// Deprecated: Use versionSchemes.VersioningOptions instead
type VersioningOptions = versionSchemes.VersioningOptions

// GitHandler interface defines methods for git operations
type GitHandler interface {
	// GenerateVersionInfo generates version information from git repository
	GenerateVersionInfo(dockerFormat bool) (*VersionInfo, error)

	// GenerateVersionInfoWithOptions generates version with custom options
	GenerateVersionInfoWithOptions(options versionSchemes.VersioningOptions) (*VersionInfo, error)

	// GetCurrentBranch returns the current branch name
	GetCurrentBranch() (string, error)

	// GetLastTag finds the last reachable tag
	GetLastTag(branchName string) (string, error)

	// GetCommitsSinceTag counts commits since the specified tag
	GetCommitsSinceTag(tagName string) (int, error)

	// GetShortHash returns the short hash of current commit
	GetShortHash() (string, error)
}

// GetGitHandler returns appropriate git handler based on inBuiltGit flag
func GetGitHandler(inBuiltGit bool, repoPath string) (GitHandler, error) {
	if inBuiltGit {
		return NewGoGitHandler(repoPath)
	}
	return NewSystemGitHandler(repoPath)
}
