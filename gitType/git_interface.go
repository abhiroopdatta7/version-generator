package gitType

// VersionInfo contains git version information
type VersionInfo struct {
	Branch       string
	LastTag      string
	CommitsSince int
	ShortHash    string
	Version      string
}

// VersioningOptions defines different versioning scheme options
type VersioningOptions struct {
	Semver bool // Use Semantic Versioning: v1.2.3-alpha.4 or v1.2.3-beta.4+branch
	CalVer bool // Use Calendar Versioning: 2024.08.4 or 2024.08.4-branch
	Simple bool // Use simple format: v1.2.3 (no branch/commit info)
	Hash   bool // Include short hash in version
}

// GitHandler interface defines methods for git operations
type GitHandler interface {
	// GenerateVersionInfo generates version information from git repository
	GenerateVersionInfo(dockerFormat bool) (*VersionInfo, error)

	// GenerateVersionInfoWithOptions generates version with custom options
	GenerateVersionInfoWithOptions(options VersioningOptions) (*VersionInfo, error)

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
