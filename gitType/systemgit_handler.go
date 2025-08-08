package gitType

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// SystemGitHandler implements GitHandler using system git executable
type SystemGitHandler struct {
	repoPath string
}

// NewSystemGitHandler creates a new system git handler
func NewSystemGitHandler(repoPath string) (*SystemGitHandler, error) {
	// Check if git is available
	_, err := exec.LookPath("git")
	if err != nil {
		return nil, fmt.Errorf("git executable not found: %w", err)
	}

	return &SystemGitHandler{repoPath: repoPath}, nil
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

	// Generate version string
	version := s.generateVersionString(lastTag, commitsSince, shortHash, branchName, dockerFormat)

	return &VersionInfo{
		Branch:       branchName,
		LastTag:      lastTag,
		CommitsSince: commitsSince,
		ShortHash:    shortHash,
		Version:      version,
	}, nil
}

// GenerateVersionInfoWithOptions generates version information using system git with custom options
func (s *SystemGitHandler) GenerateVersionInfoWithOptions(options VersioningOptions) (*VersionInfo, error) {
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

	// Generate version string using new options
	version := s.generateVersionStringWithOptions(lastTag, commitsSince, shortHash, branchName, options)

	return &VersionInfo{
		Branch:       branchName,
		LastTag:      lastTag,
		CommitsSince: commitsSince,
		ShortHash:    shortHash,
		Version:      version,
	}, nil
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

// generateVersionString generates the version string
func (s *SystemGitHandler) generateVersionString(lastTag string, commitsSince int, shortHash, branchName string, dockerFormat bool) string {
	if commitsSince == 0 {
		// We're exactly on a tag
		return lastTag
	}

	// For main/master branches, don't include branch name in version
	if branchName == "main" || branchName == "master" {
		if dockerFormat {
			// Docker format: <tag>-<count>
			return fmt.Sprintf("%s-%d", lastTag, commitsSince)
		} else {
			// Default format: <tag>+<count>
			return fmt.Sprintf("%s+%d", lastTag, commitsSince)
		}
	}

	// For other branches, clean branch name and include it in version string
	cleanBranch := regexp.MustCompile(`[^a-zA-Z0-9\-]`).ReplaceAllString(branchName, "-")

	// Choose format based on dockerFormat flag
	if dockerFormat {
		// Docker format: <tag>-<branch>-<count of commit>
		return fmt.Sprintf("%s-%s-%d", lastTag, cleanBranch, commitsSince)
	} else {
		// Default format: <tag>-<branch>+<count of commit>
		return fmt.Sprintf("%s-%s+%d", lastTag, cleanBranch, commitsSince)
	}
}

// generateVersionStringWithOptions generates version string with custom options
func (s *SystemGitHandler) generateVersionStringWithOptions(lastTag string, commitsSince int, shortHash, branchName string, options VersioningOptions) string {
	if commitsSince == 0 && !options.Hash {
		// We're exactly on a tag and no hash requested
		if options.Simple {
			return lastTag
		}
		if options.CalVer {
			return s.convertToCalVer(lastTag, 0, branchName, false)
		}
		return lastTag
	}

	// Handle different versioning schemes
	switch {
	case options.CalVer:
		return s.generateCalVerString(lastTag, commitsSince, shortHash, branchName, options.Hash)
	case options.Semver:
		return s.generateSemVerString(lastTag, commitsSince, shortHash, branchName, options.Hash)
	case options.Simple:
		return s.generateSimpleString(lastTag, commitsSince, shortHash, options.Hash)
	default:
		return s.generateDefaultString(lastTag, commitsSince, shortHash, branchName, options.Hash)
	}
}

// generateCalVerString generates Calendar Versioning format
func (s *SystemGitHandler) generateCalVerString(lastTag string, commitsSince int, shortHash, branchName string, includeHash bool) string {
	now := time.Now()
	calVer := fmt.Sprintf("%d.%02d", now.Year(), now.Month())
	
	if commitsSince > 0 {
		calVer = fmt.Sprintf("%s.%d", calVer, commitsSince)
	}
	
	if branchName != "main" && branchName != "master" && branchName != "detached" {
		cleanBranch := regexp.MustCompile(`[^a-zA-Z0-9\-]`).ReplaceAllString(branchName, "-")
		calVer = fmt.Sprintf("%s-%s", calVer, cleanBranch)
	}
	
	if includeHash {
		calVer = fmt.Sprintf("%s+%s", calVer, shortHash)
	}
	
	return calVer
}

// generateSemVerString generates Semantic Versioning format
func (s *SystemGitHandler) generateSemVerString(lastTag string, commitsSince int, shortHash, branchName string, includeHash bool) string {
	if commitsSince == 0 && !includeHash {
		return lastTag
	}

	// Parse the tag to extract semver parts
	version := lastTag
	if strings.HasPrefix(version, "v") {
		version = version[1:]
	}

	if branchName == "main" || branchName == "master" || branchName == "detached" {
		if commitsSince > 0 {
			version = fmt.Sprintf("%s-dev.%d", version, commitsSince)
		}
	} else {
		cleanBranch := regexp.MustCompile(`[^a-zA-Z0-9\-]`).ReplaceAllString(branchName, "-")
		if commitsSince > 0 {
			version = fmt.Sprintf("%s-%s.%d", version, cleanBranch, commitsSince)
		} else {
			version = fmt.Sprintf("%s-%s", version, cleanBranch)
		}
	}

	if includeHash {
		version = fmt.Sprintf("%s+%s", version, shortHash)
	}

	return "v" + version
}

// generateSimpleString generates simple version format
func (s *SystemGitHandler) generateSimpleString(lastTag string, commitsSince int, shortHash string, includeHash bool) string {
	if includeHash {
		return fmt.Sprintf("%s+%s", lastTag, shortHash)
	}
	return lastTag
}

// generateDefaultString generates default format
func (s *SystemGitHandler) generateDefaultString(lastTag string, commitsSince int, shortHash, branchName string, includeHash bool) string {
	if commitsSince == 0 && !includeHash {
		return lastTag
	}

	version := lastTag
	if branchName == "main" || branchName == "master" || branchName == "detached" {
		if commitsSince > 0 {
			version = fmt.Sprintf("%s+%d", lastTag, commitsSince)
		}
	} else {
		cleanBranch := regexp.MustCompile(`[^a-zA-Z0-9\-]`).ReplaceAllString(branchName, "-")
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

// convertToCalVer is a helper function for CalVer conversion
func (s *SystemGitHandler) convertToCalVer(lastTag string, commitsSince int, branchName string, includeHash bool) string {
	// This is a simplified implementation
	return s.generateCalVerString(lastTag, commitsSince, "", branchName, includeHash)
}
