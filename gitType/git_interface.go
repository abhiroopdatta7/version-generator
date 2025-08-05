package gitType

// VersionInfo contains git version information
type VersionInfo struct {
	Branch       string
	LastTag      string
	CommitsSince int
	ShortHash    string
	Version      string
}

// GitHandler interface defines methods for git operations
type GitHandler interface {
	// GenerateVersionInfo generates version information from git repository
	GenerateVersionInfo(dockerFormat bool) (*VersionInfo, error)

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
