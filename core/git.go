package core

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

// executes git commands
type GitExecutor struct {
	WorkDir string // Optional working directory
}

func NewGitExecutor(workDir string) *GitExecutor {
	return &GitExecutor{
		WorkDir: workDir,
	}
}

// run git command and return output
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
		return "", fmt.Errorf("git error: %s\nCommand: git %s\nOutput: %s",
			err, strings.Join(args, " "), stderr.String())
	}

	return stdout.String(), nil
}

// returns the name of the current branch
func (g *GitExecutor) GetCurrentBranch() (string, error) {
	output, err := g.Execute("rev-parse", "--abbrev-ref", "HEAD")
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(output), nil
}

// checks if the current branch has an upstream configured
func (g *GitExecutor) HasUpstream() (bool, error) {
	_, err := g.Execute("rev-parse", "--abbrev-ref", "--symbolic-full-name", "@{u}")
	if err != nil {
		// If there's an error, it means there's no upstream
		return false, nil
	}
	return true, nil
}

// CheckoutBranch checks out a branch
func (g *GitExecutor) CheckoutBranch(branch string) error {
	_, err := g.Execute("checkout", branch)
	return err
}

// creates a new branch from a parent branch
func (g *GitExecutor) CreateBranch(newBranch, parentBranch string) error {
	_, err := g.Execute("checkout", "-b", newBranch, parentBranch)
	return err
}

// rebases the current branch onto another branch
func (g *GitExecutor) RebaseBranch(targetBranch string) error {
	_, err := g.Execute("rebase", targetBranch)
	return err
}

// pushes the current branch with force-with-lease
func (g *GitExecutor) PushBranch() error {
	_, err := g.Execute("push", "--force-with-lease")
	return err
}

// sets the upstream for the current branch
func (g *GitExecutor) SetUpstreamBranch(branch string) error {
	_, err := g.Execute("push", "--set-upstream", "origin", branch, "--force-with-lease")
	return err
}

// fetches from the remote
func (g *GitExecutor) FetchRemote() error {
	_, err := g.Execute("fetch")
	return err
}

// shows the commit graph
func (g *GitExecutor) ShowGraph() (string, error) {
	return g.Execute("log", "--graph", "--oneline", "--decorate", "--all")
}
