package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// GitError represents a specific Git command error with details
type GitError struct {
	Command string
	Args    []string
	Stderr  string
	Err     error
}

// Error implements the error interface for GitError
func (e *GitError) Error() string {
	return fmt.Sprintf("git error: %s\nCommand: git %s\nOutput: %s",
		e.Err, strings.Join(e.Args, " "), e.Stderr)
}

// BranchNotFoundError represents an error when a branch doesn't exist
type BranchNotFoundError struct {
	BranchName string
}

func (e *BranchNotFoundError) Error() string {
	return fmt.Sprintf("branch not found: %s", e.BranchName)
}

// RemoteError represents an error with the remote repository
type RemoteError struct {
	Remote string
	Err    error
}

func (e *RemoteError) Error() string {
	return fmt.Sprintf("remote error with %s: %s", e.Remote, e.Err)
}

// MergeConflictError represents a git merge conflict
type MergeConflictError struct {
	Branch string
	Target string
}

func (e *MergeConflictError) Error() string {
	return fmt.Sprintf("merge conflict when rebasing %s onto %s", e.Branch, e.Target)
}

// GitExecutor handles running Git commands
type GitExecutor struct {
	WorkDir string // Optional working directory
}

// NewGitExecutor creates a new GitExecutor
func NewGitExecutor(workDir string) *GitExecutor {
	return &GitExecutor{
		WorkDir: workDir,
	}
}

// Execute runs a git command and returns its output
func (g *GitExecutor) Execute(args ...string) (string, error) {
	cmd := exec.Command("git", args...)

	if g.WorkDir != "" {
		cmd.Dir = g.WorkDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		// Perform error detection based on stderr output
		stderrStr := stderr.String()
		
		// Check for common error patterns and return specific error types
		if strings.Contains(stderrStr, "not a git repository") {
			return "", fmt.Errorf("not in a git repository")
		}
		
		if strings.Contains(stderrStr, "did not match any") && 
		   (strings.Contains(stderrStr, "pathspec") || strings.Contains(stderrStr, "reference")) {
			// Try to extract the branch name from args
			branchName := ""
			for _, arg := range args {
				if !strings.HasPrefix(arg, "-") && arg != "checkout" && arg != "branch" {
					branchName = arg
					break
				}
			}
			return "", &BranchNotFoundError{BranchName: branchName}
		}
		
		if strings.Contains(stderrStr, "conflict") && 
		   (strings.Contains(stderrStr, "merge") || strings.Contains(stderrStr, "rebase")) {
			currentBranch, _ := g.GetCurrentBranch()
			targetBranch := ""
			for i, arg := range args {
				if arg == "rebase" && i+1 < len(args) {
					targetBranch = args[i+1]
					break
				}
			}
			return "", &MergeConflictError{Branch: currentBranch, Target: targetBranch}
		}
		
		if strings.Contains(stderrStr, "could not read from remote repository") {
			remote := "origin"
			for i, arg := range args {
				if arg == "push" && i+1 < len(args) {
					remote = args[i+1]
					break
				}
			}
			return "", &RemoteError{Remote: remote, Err: err}
		}
		
		// Return generic GitError for unrecognized errors
		return "", &GitError{
			Command: "git",
			Args:    args,
			Stderr:  stderrStr,
			Err:     err,
		}
	}

	return stdout.String(), nil
}

// GetCurrentBranch returns the name of the current branch
func (g *GitExecutor) GetCurrentBranch() (string, error) {
	output, err := g.Execute("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// HasUpstream checks if the current branch has an upstream configured
func (g *GitExecutor) HasUpstream() (bool, error) {
	_, err := g.Execute("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil {
		// Check if this is a "no upstream" error vs another type of error
		if gitErr, ok := err.(*GitError); ok {
			if strings.Contains(gitErr.Stderr, "upstream") && strings.Contains(gitErr.Stderr, "not set") {
				return false, nil
			}
		}
		// Otherwise it's a different error that should be returned
		return false, err
	}
	return true, nil
}

// CheckoutBranch checks out a branch
func (g *GitExecutor) CheckoutBranch(branch string) error {
	_, err := g.Execute("checkout", branch)
	return err
}

// CreateBranch creates a new branch from a parent branch
func (g *GitExecutor) CreateBranch(newBranch, parentBranch string) error {
	_, err := g.Execute("checkout", "-b", newBranch, parentBranch)
	return err
}

// RebaseBranch rebases the current branch onto another branch
func (g *GitExecutor) RebaseBranch(targetBranch string) error {
	_, err := g.Execute("rebase", targetBranch)
	return err
}

// PushBranch pushes the current branch with force-with-lease
func (g *GitExecutor) PushBranch() error {
	_, err := g.Execute("push", "--force-with-lease")
	return err
}

// SetUpstreamBranch sets the upstream for the current branch
func (g *GitExecutor) SetUpstreamBranch(branch string) error {
	_, err := g.Execute("push", "--set-upstream", "origin", branch, "--force-with-lease")
	return err
}

// FetchRemote fetches from the remote
func (g *GitExecutor) FetchRemote() error {
	_, err := g.Execute("fetch")
	return err
}

// ShowGraph shows the commit graph
func (g *GitExecutor) ShowGraph() (string, error) {
	return g.Execute("log", "--graph", "--oneline", "--decorate", "--all")
}

// ListBranches returns a list of local branches
func (g *GitExecutor) ListBranches() ([]string, error) {
	output, err := g.Execute("branch")
	if err != nil {
		return nil, err
	}
	
	branches := []string{}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if strings.HasPrefix(branch, "*") {
			branch = strings.TrimPrefix(branch, "* ")
		}
		branches = append(branches, branch)
	}
	
	return branches, nil
}

// ListRemoteBranches returns a list of remote branches
func (g *GitExecutor) ListRemoteBranches() ([]string, error) {
	output, err := g.Execute("branch", "-r")
	if err != nil {
		return nil, err
	}
	
	branches := []string{}
	lines := strings.Split(strings.TrimSpace(output), "\n")
	for _, line := range lines {
		branch := strings.TrimSpace(line)
		if !strings.Contains(branch, "HEAD") {
			branches = append(branches, branch)
		}
	}
	
	return branches, nil
}