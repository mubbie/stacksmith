package core

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// StackConfig represents the stored branch relationships
type StackConfig struct {
	Relationships map[string]string `yaml:"relationships"`
	Metadata      struct {
		MainBranch  string    `yaml:"main_branch"`
		LastUpdated time.Time `yaml:"last_updated"`
	} `yaml:"metadata"`
}

// BranchNode represent a git branch in the stack
type BranchNode struct {
	Name      string
	CommitSHA string
	IsHead    bool
	IsOrphan  bool
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
	Orphans    []*BranchNode
}

// BranchInfo hold details about the branch (intermediary date store)
type BranchInfo struct {
	Name      string
	CommitSHA string
	ParentSHA string
}

// SaveStackConfig: Save the stack config to a file
func (g *GitExecutor) SaveStackConfig(config *StackConfig) error {
	// Get project root directory
	rootDir, err := g.Execute("rev-parse", "--show-toplevel")
	if err != nil {
		return err
	}
	rootDir = strings.TrimSpace(rootDir)

	// Create .git/stacksmith directory if it doesn't exist
	stacksmithDir := filepath.Join(rootDir, ".git", "stacksmith")
	if err := os.MkdirAll(stacksmithDir, 0755); err != nil {
		return err
	}

	// Set last updated time
	config.Metadata.LastUpdated = time.Now()

	// Add header comment and marshal to YAML
	yamlData, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	// Add header comment
	header := "# Stacksmith branch relationships\n" +
		fmt.Sprintf("# Last updated: %s\n\n", config.Metadata.LastUpdated.Format("2006-01-02 15:04:05"))

	// Write to file
	filePath := filepath.Join(stacksmithDir, "stack.yml")
	return os.WriteFile(filePath, []byte(header+string(yamlData)), 0644)
}

// LoadStackConfig loads branch relationships from YAML
func (g *GitExecutor) LoadStackConfig() (*StackConfig, error) {
	// Get project root directory
	rootDir, err := g.Execute("rev-parse", "--show-toplevel")
	if err != nil {
		return nil, err
	}
	rootDir = strings.TrimSpace(rootDir)

	// Construct file path
	filePath := filepath.Join(rootDir, ".git", "stacksmith", "stack.yml")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Return empty config if file doesn't exist
		config := &StackConfig{
			Relationships: make(map[string]string),
		}

		// Determine main branch
		for _, branch := range []string{"main", "master"} {
			if _, err := g.Execute("rev-parse", "--verify", branch); err == nil {
				config.Metadata.MainBranch = branch
				break
			}
		}

		return config, nil
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML
	var config StackConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	// Initialize if nil
	if config.Relationships == nil {
		config.Relationships = make(map[string]string)
	}

	return &config, nil
}

// RecordBranchRelationship: Record a branch relationship in the stack config
func (g *GitExecutor) RecordBranchRelationship(childBranch, parentBranch string) error {
	config, err := g.LoadStackConfig()
	if err != nil {
		return err
	}

	// Add or update relationship
	config.Relationships[childBranch] = parentBranch

	return g.SaveStackConfig(config)
}

// getBranchesWithCommits: Return all local branches with their HEAD commit SHAs
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
		if len(parts) >= 2 { // At least commit and one parent
			parentSHA := parts[1] // Use first parent if multiple (merge commit)

			result[branchName] = &BranchInfo{
				Name:      branchName,
				CommitSHA: commitSHA,
				ParentSHA: parentSHA,
			}
		}
	}

	return result, nil
}

// FindMostLikelyParent: finds the most likely parent branch for a branch
func (g *GitExecutor) FindMostLikelyParent(branch string, branchInfos map[string]*BranchInfo) string {
	// Can't find parent if branch info is missing
	branchInfo, exists := branchInfos[branch]
	if !exists {
		return ""
	}

	// Check direct commit parent relationship
	for otherName, otherInfo := range branchInfos {
		if otherName == branch {
			continue // Skip self
		}

		if otherInfo.CommitSHA == branchInfo.ParentSHA {
			return otherName // Found direct parent
		}
	}

	// If no direct parent, find branch containing this branch's parent commit
	for otherName, _ := range branchInfos {
		if otherName == branch {
			continue // Skip self
		}

		// Check if parent commit is contained in the other branch
		output, err := g.Execute("branch", "--contains", branchInfo.ParentSHA)
		if err != nil {
			continue
		}

		// Parse branches containing the commit
		containingBranches := strings.Split(strings.TrimSpace(output), "\n")
		for _, line := range containingBranches {
			name := strings.TrimSpace(strings.TrimPrefix(line, "*"))
			name = strings.TrimSpace(name)

			if name == otherName {
				return otherName
			}
		}
	}

	return "" // No parent found
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

	// load saved stack config
	config, err := g.LoadStackConfig()
	if err != nil {
		return nil, err
	}

	// build the tree structure
	stack, err := g.buildBranchTree(branchesWithCommits, branchesWithParents, config)
	if err != nil {
		return nil, err
	}

	// save updated relationships to config
	if err = g.SaveStackConfig(config); err != nil {
		// Non-fatal error, just continue for now
		fmt.Fprintf(os.Stderr, "Warning: Failed to save stack config: %v\n", err)
	}

	return stack, nil
}

// buildBranchTree: Construct the branch hierarchy with self healing
func (g *GitExecutor) buildBranchTree(
	branches map[string]string,
	branchInfos map[string]*BranchInfo,
	config *StackConfig) (*BranchStack, error) {

	// Create nodes for each branch
	nodes := make(map[string]*BranchNode)
	for name, commit := range branches {
		nodes[name] = &BranchNode{
			Name:      name,
			CommitSHA: commit,
			Children:  []*BranchNode{},
		}
	}

	// Track which branches are already processed
	processedBranches := make(map[string]bool)
	branchHasParent := make(map[string]bool)

	// Determine main branch if not set
	mainBranch := config.Metadata.MainBranch
	if mainBranch == "" {
		// Try to detect main branch
		for _, name := range []string{"main", "master"} {
			if _, exists := nodes[name]; exists {
				mainBranch = name
				config.Metadata.MainBranch = name
				break
			}
		}
	}

	// Clean up deleted branches from config
	for child := range config.Relationships {
		if nodes[child] == nil {
			delete(config.Relationships, child)
		}
	}

	// Clear duplicate entries in relationships
	// Find branches with multiple parents
	childToParent := make(map[string]string)
	for child, parent := range config.Relationships {
		// if existing, exists := childToParent[child]; exists {
		//     // For simplicity, for branch with multiple parents we'll just keep the last one we process
		// 	// TODO: Implement logging for this
		// }
		childToParent[child] = parent
	}

	// Rebuild clean relationships
	cleanRelationships := make(map[string]string)
	for child, parent := range childToParent {
		cleanRelationships[child] = parent
	}
	config.Relationships = cleanRelationships

	// Connect parent-child relationships
	for child, parent := range config.Relationships {
		// Skip if either branch doesn't exist
		if nodes[child] == nil || nodes[parent] == nil {
			continue
		}

		// Mark this branch as having a parent
		branchHasParent[child] = true

		// Only add to children if not already processed
		if !processedBranches[child] {
			nodes[parent].Children = append(nodes[parent].Children, nodes[child])
			processedBranches[child] = true
		}
	}

	// Process branches without explicit relationships
	for name, node := range nodes {
		if processedBranches[name] || name == mainBranch || branchHasParent[name] {
			continue // Skip already processed or main branch
		}

		// Try to find most likely parent
		parentName := g.FindMostLikelyParent(name, branchInfos)
		if parentName != "" && nodes[parentName] != nil {
			// Found parent, update relationships
			nodes[parentName].Children = append(nodes[parentName].Children, node)
			config.Relationships[name] = parentName
			processedBranches[name] = true
			branchHasParent[name] = true
		}
	}

	// Collect root nodes and orphans
	var rootNodes []*BranchNode
	var orphanNodes []*BranchNode

	// Add main branch as first root if it exists
	if mainBranch != "" && nodes[mainBranch] != nil {
		rootNodes = append(rootNodes, nodes[mainBranch])
		processedBranches[mainBranch] = true
	}

	// Process remaining branches that haven't been added to the tree yet
	for name, node := range nodes {
		if processedBranches[name] {
			continue // Skip branches already in the tree
		}

		if !branchHasParent[name] {
			// This branch doesn't have a parent

			// Check if it has children - if so, it's a root
			hasChildren := len(node.Children) > 0

			if hasChildren {
				rootNodes = append(rootNodes, node)
				processedBranches[name] = true
			} else {
				// No parent and no children = truly orphaned
				node.IsOrphan = true
				orphanNodes = append(orphanNodes, node)
				processedBranches[name] = true
			}
		}
	}

	// Mark current branch as HEAD
	currentBranch, err := g.GetCurrentBranch()
	if err == nil && nodes[currentBranch] != nil {
		nodes[currentBranch].IsHead = true
	}

	// Add health information (ahead/behind counts)
	for childName, parentName := range config.Relationships {
		// Skip if either branch is missing
		if nodes[childName] == nil || nodes[parentName] == nil {
			continue
		}

		// Calculate ahead/behind
		ahead, behind, err := g.GetAheadBehind(childName, parentName)
		if err == nil {
			nodes[childName].Ahead = ahead
			nodes[childName].Behind = behind
		}

		// Check if merged
		merged, err := g.IsBranchMerged(childName, parentName)
		if err == nil {
			nodes[childName].IsMerged = merged
		}
	}

	return &BranchStack{
		Roots:      rootNodes,
		AllNodes:   nodes,
		MainBranch: mainBranch,
		Orphans:    orphanNodes,
	}, nil
}
