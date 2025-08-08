package gitType

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"version-generator/versionSchemes"
)

// SystemGitHandler implements GitHandler using system git executable
type SystemGitHandler struct {
	repoPath string
	*BaseGitHandler
}

// NewSystemGitHandler creates a new system git handler
func NewSystemGitHandler(repoPath string) (*SystemGitHandler, error) {
	// Check if git is available
	_, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("git executable not found: %w", err)
	}

	return &SystemGitHandler{
		repoPath:       repoPath,
		BaseGitHandler: NewBaseGitHandler(),
	}, nil
}

// runGitCommand executes a git command and returns the output
func (s *SystemGitHandler) runGitCommand(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = s.repoPath

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("git command failed: %w", err)
	}

	return strings.TrimSpace(string(output)), nil
}

// GenerateVersionInfo generates version information using system git
func (s *SystemGitHandler) GenerateVersionInfo(dockerFormat bool) (*VersionInfo, error) {
	// Get current branch
	branchName, err := s.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	// Get short hash
	shortHash, err := s.GetShortHash()
	if err != nil {
		return nil, err
	}

	// Find the last tag
	lastTag, err := s.GetLastTag(branchName)
	if err != nil {
		return nil, err
	}

	// Count commits since last tag
	commitsSince, err := s.GetCommitsSinceTag(lastTag)
	if err != nil {
		return nil, err
	}

	// Use base handler to generate version info
	return s.GenerateVersionInfoFromComponents(branchName, shortHash, lastTag, commitsSince, dockerFormat), nil
}

// GenerateVersionInfoWithOptions generates version information using system git with custom options
func (s *SystemGitHandler) GenerateVersionInfoWithOptions(options versionSchemes.VersioningOptions) (*VersionInfo, error) {
	// Get current branch
	branchName, err := s.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	// Get short hash
	shortHash, err := s.GetShortHash()
	if err != nil {
		return nil, err
	}

	// Find the last tag
	lastTag, err := s.GetLastTag(branchName)
	if err != nil {
		return nil, err
	}

	// Count commits since last tag
	commitsSince, err := s.GetCommitsSinceTag(lastTag)
	if err != nil {
		return nil, err
	}

	// Use base handler to generate version info with options
	return s.GenerateVersionInfoFromComponentsWithOptions(branchName, shortHash, lastTag, commitsSince, options), nil
}

// GetCurrentBranch returns the current branch name
func (s *SystemGitHandler) GetCurrentBranch() (string, error) {
	output, err := s.runGitCommand("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get current branch: %w", err)
	}

	// If in detached HEAD state, try to find which branch contains this commit
	if output == "HEAD" {
		// Try to find a branch that contains the current commit
		branchOutput, err := s.runGitCommand("branch", "--contains", "HEAD")
		if err == nil && branchOutput != "" {
			// Parse the output to get the first branch name
			lines := strings.Split(strings.TrimSpace(branchOutput), "\n")
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" && !strings.HasPrefix(line, "*") {
					// Remove any leading characters and return the branch name
					if strings.HasPrefix(line, "* ") {
						line = line[2:]
					}
					// Skip detached HEAD indicators
					if !strings.Contains(line, "detached") && !strings.Contains(line, "HEAD") {
						return line, nil
					}
				}
			}
		}
		// If no branch found or all were detached indicators, return "detached"
		return "detached", nil
	}

	return output, nil
}

// GetShortHash returns the short hash of current commit
func (s *SystemGitHandler) GetShortHash() (string, error) {
	output, err := s.runGitCommand("rev-parse", "--short", "HEAD")
	if err != nil {
		return "", fmt.Errorf("failed to get short hash: %w", err)
	}
	return output, nil
}

// GetLastTag finds the last reachable tag
func (s *SystemGitHandler) GetLastTag(branchName string) (string, error) {
	// For non-main/master branches, find tags from the merge-base with main/master
	if branchName != "main" && branchName != "master" {
		return s.findTagFromRebasePoint(branchName)
	}

	// For main/master branches, find the most recent tag
	output, err := s.runGitCommand("describe", "--tags", "--abbrev=0")
	if err != nil {
		// No tags found
		return "v0.0.0", nil
	}

	return output, nil
}

// findTagFromRebasePoint finds tags from the rebase point for feature branches
func (s *SystemGitHandler) findTagFromRebasePoint(branchName string) (string, error) {
	// Try to find main or master branch and get merge-base
	var mergeBase string
	var err error

	// Try main first
	mergeBase, err = s.runGitCommand("merge-base", "HEAD", "main")
	if err != nil {
		// Try master
		mergeBase, err = s.runGitCommand("merge-base", "HEAD", "master")
		if err != nil {
			// If no main/master branch found, fall back to current branch logic
			return s.GetLastTag("main") // This will use the regular logic
		}
	}

	// Find the most recent tag reachable from the merge-base
	output, err := s.runGitCommand("describe", "--tags", "--abbrev=0", mergeBase)
	if err != nil {
		// No tags found
		return "v0.0.0", nil
	}

	return output, nil
}

// GetCommitsSinceTag counts commits since the specified tag
func (s *SystemGitHandler) GetCommitsSinceTag(tagName string) (int, error) {
	if tagName == "v0.0.0" {
		// Count all commits if no tag exists
		output, err := s.runGitCommand("rev-list", "--count", "HEAD")
		if err != nil {
			return 0, fmt.Errorf("failed to count all commits: %w", err)
		}

		count, err := strconv.Atoi(output)
		if err != nil {
			return 0, fmt.Errorf("failed to parse commit count: %w", err)
		}

		return count, nil
	}

	// Check if we're exactly on the tag
	currentHash, err := s.runGitCommand("rev-parse", "HEAD")
	if err != nil {
		return 0, fmt.Errorf("failed to get current commit hash: %w", err)
	}

	tagHash, err := s.runGitCommand("rev-parse", tagName+"^{commit}")
	if err != nil {
		return 0, fmt.Errorf("failed to get tag commit hash: %w", err)
	}

	if currentHash == tagHash {
		return 0, nil
	}

	// Count commits since tag
	output, err := s.runGitCommand("rev-list", "--count", "HEAD", "^"+tagName)
	if err != nil {
		return 0, fmt.Errorf("failed to count commits since tag: %w", err)
	}

	count, err := strconv.Atoi(output)
	if err != nil {
		return 0, fmt.Errorf("failed to parse commit count: %w", err)
	}

	return count, nil
}
