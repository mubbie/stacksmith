package core

import (
	"strings"
)

// BranchNode represent a git branch in the stack
type BranchNode struct {
	Name      string
	CommitSHA string
	IsHead    bool
	Children  []*BranchNode

	Ahead    int
	Behind   int
	IsMerged bool
}

// BranchStack represents the stack of branches
type BranchStack struct {
	Roots      []*BranchNode
	AllNodes   map[string]*BranchNode
	MainBranch string // main or master
}

// BranchInfo hold details about the branch (intermediary date store)
type BranchInfo struct {
	Name      string
	CommitSHA string
	ParentSHA string
}

// BuildBranchStack: Analyze git history and build the branch stack
func (g *GitExecutor) BuildBranchStack() (*BranchStack, error) {
	// get all local branches
	branchesWithCommits, err := g.getBranchesWithCommits()
	if err != nil {
		return nil, err
	}

	// find parent commit SHA for each branch's HEAD commit
	branchesWithParents, err := g.findParentCommits(branchesWithCommits)
	if err != nil {
		return nil, err
	}

	// build the tree structure
	stack, err := g.buildBranchTree(branchesWithParents)
	if err != nil {
		return nil, err
	}

	return stack, nil
}

// getBranchesWithCommits: Returns all local branches with the HEAD commits
func (g *GitExecutor) getBranchesWithCommits() (map[string]string, error) {
	output, err := g.Execute("for-each-ref", "--format=%(refname:short) %(objectname)", "refs/heads/")
	if err != nil {
		return nil, err
	}

	branches := make(map[string]string)
	lines := strings.Split(strings.TrimSpace(output), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}
		parts := strings.SplitN(strings.TrimSpace(line), " ", 2)
		if len(parts) == 2 {
			branchName := parts[0]
			commitSHA := parts[1]
			branches[branchName] = commitSHA
		}
	}

	return branches, nil
}

// findParentCommits: gets the parents commit for each branch's HEAD
func (g *GitExecutor) findParentCommits(branches map[string]string) (map[string]*BranchInfo, error) {
	result := make(map[string]*BranchInfo)

	for branchName, commitSHA := range branches {
		// Get the parent commit SHA
		output, err := g.Execute("rev-list", "--parents", "-n", "1", commitSHA)
		if err != nil {
			continue // Skip if error
		}

		// Parse the output: <commit-sha> <parent1-sha> [<parent2-sha> ...]
		parts := strings.Fields(strings.TrimSpace(output))
		if len(parts) >= 2 {
			parentSHA := parts[1]

			result[branchName] = &BranchInfo{
				Name:      branchName,
				CommitSHA: commitSHA,
				ParentSHA: parentSHA,
			}
		}
	}

	return result, nil
}

// buildBranchTree: Construct the branch hierarchy
func (g *GitExecutor) buildBranchTree(branches map[string]*BranchInfo) (*BranchStack, error) {
	// create nodes for each branch
	nodes := make(map[string]*BranchNode)
	for name, info := range branches {
		nodes[name] = &BranchNode{
			Name:      name,
			CommitSHA: info.CommitSHA,
			Children:  []*BranchNode{},
		}
	}

	// build commit SHA to branch mapping
	shaToBranch := make(map[string]string)
	for name, info := range branches {
		shaToBranch[info.CommitSHA] = name
	}

	// build map of branches to their parent branches
	branchToParent := make(map[string]string)
	for name, node := range nodes {
		for _, child := range node.Children {
			branchToParent[child.Name] = name
		}
	}

	// connect parent-child relationships
	var rootNodes []*BranchNode

	for name, info := range branches {
		// Find the branch that has this branch's parent commit as its HEAD
		parentName := ""
		for bn, bi := range branches {
			if bi.CommitSHA == info.ParentSHA {
				parentName = bn
				break
			}
		}

		if parentName != "" {
			// Found parent branch, add as child
			nodes[parentName].Children = append(nodes[parentName].Children, nodes[name])
		} else {
			// No parent found, add to root nodes
			rootNodes = append(rootNodes, nodes[name])
		}
	}

	// determine main branch (master or main - won't work if other branch name is used)
	mainBranch := "main"
	_, err := g.Execute("rev-parse", "--verify", "main")
	if err != nil {
		// Try master if main doesn't exist
		_, err = g.Execute("rev-parse", "--verify", "master")
		if err == nil {
			mainBranch = "master"
		}
	}

	// if main/master is found, ensure it's first in the roots list
	if nodes[mainBranch] != nil {
		// Check if it's already in the roots
		foundInRoots := false
		for i, node := range rootNodes {
			if node.Name == mainBranch {
				// Move to front
				foundInRoots = true
				if i > 0 {
					rootNodes = append([]*BranchNode{node}, append(rootNodes[:i], rootNodes[i+1:]...)...)
				}
				break
			}
		}

		// If not in roots, it might be a child - force it to be a root
		if !foundInRoots {
			rootNodes = append([]*BranchNode{nodes[mainBranch]}, rootNodes...)
		}
	}

	// mark head branch
	currentBranch, err := g.GetCurrentBranch()
	if err == nil && nodes[currentBranch] != nil {
		nodes[currentBranch].IsHead = true
	}

	// Check if branches are merged to parent branch
	for name, node := range nodes {
		// Skip branches without a parent (root branches)
		parentName, hasParent := branchToParent[name]
		if !hasParent {
			continue
		}

		// Check if this branch is merged into its parent
		merged, err := g.IsBranchMerged(name, parentName)
		if err == nil {
			node.IsMerged = merged
		}
	}

	// add health information (ahead/behind count)
	for name, node := range nodes {
		// Find direct parent branch
		for parentName, parentNode := range nodes {
			for _, child := range parentNode.Children {
				if child.Name == name {
					// Found the parent, calculate ahead/behind
					ahead, behind, err := g.GetAheadBehind(name, parentName)
					if err == nil {
						node.Ahead = ahead
						node.Behind = behind
					}
					break
				}
			}
		}
	}

	return &BranchStack{
		Roots:      rootNodes,
		AllNodes:   nodes,
		MainBranch: mainBranch,
	}, nil
}
