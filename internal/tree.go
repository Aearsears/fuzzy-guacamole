package internal

import (
	"fmt"
	"strings"
)

// used for folder structure in s3
type Tree struct {
	Root *TreeNode
}

type TreeNode struct {
	Value    string
	Children []*TreeNode
	childMap map[string]*TreeNode
	Parent   *TreeNode
	Level    int
}

func (t *Tree) Display() string {
	s := fmt.Sprintf("%s\n", t.Root.Value)
	for _, child := range t.Root.Children {
		s += t.displayNode(child, 1)
	}
	return s
}
func (t *Tree) displayNode(node *TreeNode, level int) string {
	if node == nil {
		return ""
	}
	s := fmt.Sprintf("%s%s\n", strings.Repeat(" ", level*2), node.Value)
	for _, child := range node.Children {
		s += t.displayNode(child, level+1)
	}
	return s
}

func (n *TreeNode) AddNode(path string, level int) {
	before, after, found := strings.Cut(path, "/")
	if n.childMap == nil {
		n.childMap = make(map[string]*TreeNode)
	}

	if !found {
		//create leaf node
		if _, ok := n.childMap[before]; ok {
			return
		}
		node := &TreeNode{Value: before, Level: level + 1, Parent: n}
		n.childMap[before] = node
		n.Children = append(n.Children, node)
	} else {
		//create dir node
		if _, ok := n.childMap[before]; ok {
			n.childMap[before].AddNode(after, level+1)
			return
		}
		node := &TreeNode{Value: before, Level: level + 1, Parent: n}
		n.childMap[before] = node
		n.Children = append(n.Children, node)
		node.AddNode(after, level+1)
	}
}

func CreateTree(objs []string) *Tree {
	t := &Tree{}
	t.Root = &TreeNode{
		Value:    "/",
		Level:    0,
		Children: []*TreeNode{},
		childMap: make(map[string]*TreeNode),
		Parent:   nil,
	}
	for _, obj := range objs {
		t.Root.AddNode(obj, 0)
	}
	return t
}
