package gitType

import (
	"fmt"
	"sort"
	"version-generator/versionSchemes"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// GoGitHandler implements GitHandler using go-git library
type GoGitHandler struct {
	repo *git.Repository
	*BaseGitHandler
}

// NewGoGitHandler creates a new go-git handler
func NewGoGitHandler(repoPath string) (*GoGitHandler, error) {
	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	return &GoGitHandler{
		repo:           repo,
		BaseGitHandler: NewBaseGitHandler(),
	}, nil
}

// GenerateVersionInfo generates version information using go-git
func (g *GoGitHandler) GenerateVersionInfo(dockerFormat bool) (*VersionInfo, error) {
	// Get current branch
	branchName, err := g.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	// Get short hash
	shortHash, err := g.GetShortHash()
	if err != nil {
		return nil, err
	}

	// Find the last tag
	lastTag, err := g.GetLastTag(branchName)
	if err != nil {
		return nil, err
	}

	// Count commits since last tag
	commitsSince, err := g.GetCommitsSinceTag(lastTag)
	if err != nil {
		return nil, err
	}

	// Use base handler to generate version info
	return g.GenerateVersionInfoFromComponents(branchName, shortHash, lastTag, commitsSince, dockerFormat), nil
}

// GenerateVersionInfoWithOptions generates version information using go-git with custom options
func (g *GoGitHandler) GenerateVersionInfoWithOptions(options versionSchemes.VersioningOptions) (*VersionInfo, error) {
	// Get current branch
	branchName, err := g.GetCurrentBranch()
	if err != nil {
		return nil, err
	}

	// Get short hash
	shortHash, err := g.GetShortHash()
	if err != nil {
		return nil, err
	}

	// Find the last tag
	lastTag, err := g.GetLastTag(branchName)
	if err != nil {
		return nil, err
	}

	// Count commits since last tag
	commitsSince, err := g.GetCommitsSinceTag(lastTag)
	if err != nil {
		return nil, err
	}

	// Use base handler to generate version info with options
	return g.GenerateVersionInfoFromComponentsWithOptions(branchName, shortHash, lastTag, commitsSince, options), nil
}

// GetCurrentBranch returns the current branch name
func (g *GoGitHandler) GetCurrentBranch() (string, error) {
	head, err := g.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	if head.Name().IsBranch() {
		return head.Name().Short(), nil
	}

	// If it's a detached HEAD, try to find which branch contains this commit
	currentHash := head.Hash()

	// Get all branch references
	refs, err := g.repo.References()
	if err != nil {
		return "detached", nil
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		if ref.Name().IsBranch() {
			// Check if this branch contains the current commit
			branchCommit, err := g.repo.CommitObject(ref.Hash())
			if err != nil {
				return nil // Continue to next branch
			}

			// Walk through the branch history to see if it contains our commit
			iter := object.NewCommitPreorderIter(branchCommit, nil, nil)
			defer iter.Close()

			found := false
			iter.ForEach(func(c *object.Commit) error {
				if c.Hash == currentHash {
					found = true
					return fmt.Errorf("found") // Break the loop
				}
				return nil
			})

			if found {
				return fmt.Errorf("branch:%s", ref.Name().Short()) // Break and return this branch
			}
		}
		return nil
	})

	if err != nil && err.Error() != "" {
		if errMsg := err.Error(); len(errMsg) > 7 && errMsg[:7] == "branch:" {
			return errMsg[7:], nil
		}
	}

	// If no branch found, return "detached"
	return "detached", nil
}

// GetShortHash returns the short hash of current commit
func (g *GoGitHandler) GetShortHash() (string, error) {
	head, err := g.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}
	return head.Hash().String()[:7], nil
}

// GetLastTag finds the last reachable tag
func (g *GoGitHandler) GetLastTag(branchName string) (string, error) {
	head, err := g.repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	// For non-main/master branches, find tags from the rebase point
	if branchName != "main" && branchName != "master" {
		return g.findTagFromRebasePoint(head.Hash(), branchName)
	}

	// For main/master branches, use the original logic
	return g.findTagFromCurrentBranch(head.Hash())
}

// GetCommitsSinceTag counts commits since the specified tag
func (g *GoGitHandler) GetCommitsSinceTag(tagName string) (int, error) {
	head, err := g.repo.Head()
	if err != nil {
		return 0, fmt.Errorf("failed to get HEAD: %w", err)
	}

	if tagName == "v0.0.0" {
		// Count all commits if no tag exists
		return g.countAllCommits(head.Hash())
	}

	// Find the tag commit hash
	tagRef, err := g.repo.Tag(tagName)
	if err != nil {
		return 0, fmt.Errorf("failed to get tag reference: %w", err)
	}

	var tagCommitHash plumbing.Hash
	obj, err := g.repo.Object(plumbing.AnyObject, tagRef.Hash())
	if err != nil {
		return 0, err
	}

	switch o := obj.(type) {
	case *object.Tag:
		tagCommitHash = o.Target
	case *object.Commit:
		tagCommitHash = tagRef.Hash()
	default:
		return 0, fmt.Errorf("tag points to unsupported object type")
	}

	if head.Hash() == tagCommitHash {
		return 0, nil
	}

	// Count commits between current and tag
	commit, err := g.repo.CommitObject(head.Hash())
	if err != nil {
		return 0, err
	}

	count := 0
	iter := object.NewCommitPreorderIter(commit, nil, nil)
	defer iter.Close()

	err = iter.ForEach(func(c *object.Commit) error {
		if c.Hash == tagCommitHash {
			return fmt.Errorf("found tag") // Break the loop
		}
		count++
		return nil
	})

	if err != nil && err.Error() != "found tag" {
		return 0, err
	}

	return count, nil
}

// findTagFromRebasePoint finds tags from the rebase point for feature branches
func (g *GoGitHandler) findTagFromRebasePoint(commitHash plumbing.Hash, branchName string) (string, error) {
	// Try to find main or master branch
	var mainBranch *plumbing.Reference

	// Try to get main branch first
	mainRef, err := g.repo.Reference(plumbing.NewBranchReferenceName("main"), true)
	if err == nil {
		mainBranch = mainRef
	} else {
		// Try master branch
		masterRef, err := g.repo.Reference(plumbing.NewBranchReferenceName("master"), true)
		if err == nil {
			mainBranch = masterRef
		}
	}

	if mainBranch == nil {
		// If no main/master branch found, fall back to current branch logic
		return g.findTagFromCurrentBranch(commitHash)
	}

	// Find common ancestor between current branch and main/master
	commonAncestor, err := g.findCommonAncestor(commitHash, mainBranch.Hash())
	if err != nil {
		// If can't find common ancestor, fall back to current branch logic
		return g.findTagFromCurrentBranch(commitHash)
	}

	// Find tags reachable from the common ancestor
	return g.findTagFromCurrentBranch(commonAncestor)
}

// findCommonAncestor finds the common ancestor between two commits
func (g *GoGitHandler) findCommonAncestor(commit1, commit2 plumbing.Hash) (plumbing.Hash, error) {
	if commit1 == commit2 {
		return commit1, nil
	}

	// Get commits for both hashes
	c1, err := g.repo.CommitObject(commit1)
	if err != nil {
		return plumbing.ZeroHash, err
	}

	c2, err := g.repo.CommitObject(commit2)
	if err != nil {
		return plumbing.ZeroHash, err
	}

	// Create a map of all ancestors of commit2
	ancestorsMap := make(map[plumbing.Hash]bool)
	iter := object.NewCommitPreorderIter(c2, nil, nil)
	defer iter.Close()

	err = iter.ForEach(func(c *object.Commit) error {
		ancestorsMap[c.Hash] = true
		return nil
	})
	if err != nil {
		return plumbing.ZeroHash, err
	}

	// Walk through ancestors of commit1 to find first common ancestor
	iter1 := object.NewCommitPreorderIter(c1, nil, nil)
	defer iter1.Close()

	var commonAncestor plumbing.Hash
	err = iter1.ForEach(func(c *object.Commit) error {
		if ancestorsMap[c.Hash] {
			commonAncestor = c.Hash
			return fmt.Errorf("found common ancestor") // Break the loop
		}
		return nil
	})

	if err != nil && err.Error() == "found common ancestor" {
		return commonAncestor, nil
	}

	return plumbing.ZeroHash, fmt.Errorf("no common ancestor found")
}

// findTagFromCurrentBranch finds tags reachable from current branch
func (g *GoGitHandler) findTagFromCurrentBranch(commitHash plumbing.Hash) (string, error) {
	// Get all tags
	tagRefs, err := g.repo.Tags()
	if err != nil {
		return "", fmt.Errorf("failed to get tags: %w", err)
	}

	var tags []struct {
		name string
		hash plumbing.Hash
		time int64
	}

	err = tagRefs.ForEach(func(ref *plumbing.Reference) error {
		tagName := ref.Name().Short()

		// Get the commit that the tag points to
		var tagCommitHash plumbing.Hash

		// Check if it's an annotated tag
		obj, err := g.repo.Object(plumbing.AnyObject, ref.Hash())
		if err != nil {
			return err
		}

		switch o := obj.(type) {
		case *object.Tag:
			// Annotated tag
			tagCommitHash = o.Target
		case *object.Commit:
			// Lightweight tag
			tagCommitHash = ref.Hash()
		default:
			return nil // Skip non-commit/tag objects
		}

		// Check if this tag is reachable from the current commit
		isReachable, err := g.isCommitReachable(commitHash, tagCommitHash)
		if err != nil {
			return err
		}

		if isReachable {
			commit, err := g.repo.CommitObject(tagCommitHash)
			if err != nil {
				return err
			}

			tags = append(tags, struct {
				name string
				hash plumbing.Hash
				time int64
			}{
				name: tagName,
				hash: tagCommitHash,
				time: commit.Committer.When.Unix(),
			})
		}

		return nil
	})

	if err != nil {
		return "", err
	}

	if len(tags) == 0 {
		return "v0.0.0", nil // No tags found
	}

	// Sort tags by commit time (newest first)
	sort.Slice(tags, func(i, j int) bool {
		return tags[i].time > tags[j].time
	})

	return tags[0].name, nil
}

// isCommitReachable checks if a commit is reachable from another commit
func (g *GoGitHandler) isCommitReachable(from, to plumbing.Hash) (bool, error) {
	if from == to {
		return true, nil
	}

	// Get commit iterator from 'from' commit
	commit, err := g.repo.CommitObject(from)
	if err != nil {
		return false, err
	}

	iter := object.NewCommitPreorderIter(commit, nil, nil)
	defer iter.Close()

	err = iter.ForEach(func(c *object.Commit) error {
		if c.Hash == to {
			return fmt.Errorf("found") // Use error to break the loop
		}
		return nil
	})

	return err != nil && err.Error() == "found", nil
}

// countAllCommits counts all commits from a given commit
func (g *GoGitHandler) countAllCommits(commitHash plumbing.Hash) (int, error) {
	commit, err := g.repo.CommitObject(commitHash)
	if err != nil {
		return 0, err
	}

	count := 0
	iter := object.NewCommitPreorderIter(commit, nil, nil)
	defer iter.Close()

	err = iter.ForEach(func(c *object.Commit) error {
		count++
		return nil
	})

	return count, err
}
