package internal

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddNodeAndDisplayChildren(t *testing.T) {
	root := &TreeNode{
		Value:    "/",
		Level:    0,
		Children: []*TreeNode{},
		childMap: make(map[string]*TreeNode),
		Parent:   nil,
	}

	root.AddNode("folder1/file1.txt", 0)
	root.AddNode("folder1/file2.txt", 0)
	root.AddNode("folder2/file3.txt", 0)
	root.AddNode("file4.txt", 0)

	assert.Equal(t, 3, len(root.Children))
	assert.Contains(t, root.DisplayChildren(), "folder1")
	assert.Contains(t, root.DisplayChildren(), "folder2")
	assert.Contains(t, root.DisplayChildren(), "file4.txt")

	// Check children of folder1
	folder1 := root.childMap["folder1"]
	assert.NotNil(t, folder1)
	assert.Equal(t, 2, len(folder1.Children))
	assert.Contains(t, folder1.DisplayChildren(), "file1.txt")
	assert.Contains(t, folder1.DisplayChildren(), "file2.txt")
}

func TestCreateTreeAndDisplay(t *testing.T) {
	objects := []string{
		"dir1/file1.txt",
		"dir1/file2.txt",
		"dir2/file3.txt",
		"file4.txt",
	}
	tree := CreateTree(objects)
	assert.NotNil(t, tree)
	assert.NotNil(t, tree.Root)
	assert.Equal(t, "/", tree.Root.Value)
	assert.Equal(t, 3, len(tree.Root.Children))

	display := tree.Display()
	assert.True(t, strings.Contains(display, "dir1"))
	assert.True(t, strings.Contains(display, "dir2"))
	assert.True(t, strings.Contains(display, "file4.txt"))
	assert.True(t, strings.Contains(display, "file1.txt"))
	assert.True(t, strings.Contains(display, "file2.txt"))
	assert.True(t, strings.Contains(display, "file3.txt"))
}

func TestAddNodeDuplicate(t *testing.T) {
	root := &TreeNode{
		Value:    "/",
		Level:    0,
		Children: []*TreeNode{},
		childMap: make(map[string]*TreeNode),
		Parent:   nil,
	}

	root.AddNode("file.txt", 0)
	root.AddNode("file.txt", 0) // duplicate
	assert.Equal(t, 1, len(root.Children))
}

func TestDisplayNodeNil(t *testing.T) {
	tree := &Tree{}
	result := tree.displayNode(nil, 1)
	assert.Equal(t, "", result)
}
