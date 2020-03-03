package commands

import (
	"errors"
	"fmt"
	"strings"
)

// CTree structure represents a tree of kubectl commands
type CTree struct {
	Part     string
	Parent   *CTree
	Children []*CTree
}

// NewCTree creates a kubectl command tree
func NewCTree(commands []string) (*CTree, error) {
	root := CTree{
		Part: "kubectl",
	}

	for _, command := range commands {
		if err := root.MergeCommand(command); err != nil {
			return nil, err
		}
	}
	return &root, nil
}

// MergeCommand merges a kubectl command into the tree
func (tree *CTree) MergeCommand(command string) error {
	parts := split(command)

	if len(parts) == 0 || parts[0] != "kubectl" {
		return errors.New("This is not a kubectl command")
	}

	tree.addChildren(parts[1:])
	return nil
}

func split(command string) []string {
	parts := strings.Fields(command)
	var joinedParts []string

	// Potential ways to specify flags
	// --flag
	// -f
	// --flag value
	// --flag=value
	// -f value
	// -f=value
	for i, max := 0, len(parts); i < max; i++ {
		if i+1 < max && strings.HasPrefix(parts[i], "-") && !strings.HasPrefix(parts[i+1], "-") {
			joinedParts = append(joinedParts, strings.Join([]string{parts[i], parts[i+1]}, " "))
			i++
		} else {
			joinedParts = append(joinedParts, parts[i])
		}
	}

	return joinedParts
}

func (tree *CTree) addChildren(parts []string) {
	if len(parts) == 0 {
		return
	}

	for _, child := range tree.Children {
		if child.Part == parts[0] {
			if len(parts) > 0 {
				child.addChildren(parts[1:])
				return
			}
		}
	}

	child := CTree{
		Part: parts[0],
	}
	tree.addChild(&child)
	child.addChildren(parts[1:])
}

func (tree *CTree) addChild(child *CTree) *CTree {
	tree.Children = append(tree.Children, child)
	child.Parent = tree

	return tree
}

// GetCmd returns the executable kubectl command at a certain position in the tree
// (depth-first search)
//      1
//    /   \
//   2     5
//  / \   / \
// 3   4 6   8
//       |
//       7
func (tree *CTree) GetCmd(position int) *Cmd {
	current := tree.getTree(&position)
	if current != nil {
		return current.toCmd()
	}

	return nil
}

func (tree *CTree) getTree(position *int) *CTree {
	if *position <= 0 {
		return nil
	}

	if *position == 1 {
		return tree
	}

	if len(tree.Children) > 0 {
		for _, child := range tree.Children {
			*position = *position - 1
			found := child.getTree(position)
			if found != nil {
				return found
			}
		}
	}

	return nil
}

func (tree *CTree) toCmd() *Cmd {
	var args []string
	var current *CTree
	for current = tree; current.Parent != nil; current = current.Parent {
		parts := strings.Split(current.Part, " ")
		args = append(parts, args...)
	}
	return NewCmd(current.Part, args...)
}

// GetPosition returns the position of a kubectl command in the tree (depth-first search)
//      1
//    /   \
//   2     5
//  / \   / \
// 3   4 6   8
//       |
//       7
func (tree *CTree) GetPosition(command string) *int {
	parts := split(command)
	position := 1
	return tree.getPosition(parts, &position)
}

func (tree *CTree) getPosition(parts []string, position *int) *int {
	if len(parts) == 1 && parts[0] == tree.Part {
		return position
	}

	for _, child := range tree.Children {
		*position++
		if len(parts) > 0 && parts[0] == tree.Part {
			if p := child.getPosition(parts[1:], position); p != nil {
				return p
			}
		} else {
			child.getPosition(nil, position)
		}
	}

	return nil
}

// GetNextParts returns a list of potential next parts for a given command
func (tree *CTree) GetNextParts(command string) []string {
	parts := split(command)
	return tree.getNextParts(parts)
}

func (tree *CTree) getNextParts(parts []string) []string {
	if len(parts) == 0 {
		return nil
	}

	if tree.Part == parts[0] {
		var nextParts []string
		for _, child := range tree.Children {
			nextPart := child.getNextParts(parts[1:])
			nextParts = append(nextParts, nextPart...)
		}
		return nextParts
	}

	if !strings.Contains(tree.Part, parts[0]) {
		return nil
	}

	return []string{tree.Part}
}

// RemoveCommand removes a kubectl command and all its children commands from the tree
func (tree *CTree) RemoveCommand(position int) error {
	found := tree.getTree(&position)
	if found != nil {
		found.remove()
		return nil
	}
	return errors.New("Kubectl not found")
}

func (tree *CTree) remove() {
	if tree.Parent != nil {
		for index, sibling := range tree.Parent.Children {
			if sibling == tree {
				tree.Parent.Children = append(
					tree.Parent.Children[:index],
					tree.Parent.Children[index+1:]...,
				)
			}
		}
	}
}

// ToStrings returns a list of indented command parts.
// Example output:
//   "kubectl"
//   "  -n kubeflow"
//   "    get"
//   "      pod"
//   "      cronjob"
//   "  -n pipelines"
//   "    get"
//   "      pod"
func (tree *CTree) ToStrings(tabSize int) []string {
	var all []string
	return tree.toStrings(tabSize, 0, &all)
}

func (tree *CTree) toStrings(tabSize, depth int, all *[]string) []string {
	prefix := strings.Repeat(" ", depth*tabSize)
	current := fmt.Sprintf("%s%s", prefix, tree.Part)
	*all = append(*all, current)

	if len(tree.Children) > 0 {
		for _, child := range tree.Children {
			child.toStrings(tabSize, depth+1, all)
		}
	}

	return *all
}

// Serialize returns the complete list of commands required to rebuild the tree from scratch
// Example output:
//   "kubectl -n kubeflow get pod"
//   "kubectl -n kubeflow get cronjob"
//   "kubectl -n pipelines get pod"
func (tree *CTree) Serialize() []string {
	var all []string
	tree.serialize(&all)
	return all
}

func (tree *CTree) serialize(all *[]string) {
	if len(tree.Children) > 0 {
		for _, child := range tree.Children {
			child.serialize(all)
		}
	} else {
		command := tree.toCmd()
		*all = append(*all, command.ToString())
	}
}
