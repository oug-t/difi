package tree

import (
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/charmbracelet/bubbles/list"
)

// TreeItem represents a file or folder in the UI list.
type TreeItem struct {
	Path     string
	FullPath string
	IsDir    bool
	Depth    int
}

func (i TreeItem) FilterValue() string { return i.FullPath }
func (i TreeItem) Description() string { return "" }
func (i TreeItem) Title() string {
	indent := strings.Repeat("  ", i.Depth)
	icon := getIcon(i.Path, i.IsDir)
	return fmt.Sprintf("%s%s %s", indent, icon, i.Path)
}

// Build converts a list of file paths into a sorted tree list.
// Compaction is disabled to ensure tree stability.
func Build(paths []string) []list.Item {
	// Initialize root
	root := &node{
		children: make(map[string]*node),
		isDir:    true,
	}

	// 1. Build the raw tree structure
	for _, path := range paths {
		parts := strings.Split(path, "/")
		current := root
		for i, part := range parts {
			if _, exists := current.children[part]; !exists {
				isDir := i < len(parts)-1
				fullPath := strings.Join(parts[:i+1], "/")
				current.children[part] = &node{
					name:     part,
					fullPath: fullPath,
					children: make(map[string]*node),
					isDir:    isDir,
				}
			}
			current = current.children[part]
		}
	}

	// 2. Flatten to list items (Sorting happens here)
	var items []list.Item
	flatten(root, 0, &items)
	return items
}

// -- Helpers --

type node struct {
	name     string
	fullPath string
	children map[string]*node
	isDir    bool
}

// flatten recursively converts the tree into a linear list, sorting children by type and name.
func flatten(n *node, depth int, items *[]list.Item) {
	keys := make([]string, 0, len(n.children))
	for k := range n.children {
		keys = append(keys, k)
	}

	// Sort: Directories first, then alphabetical
	sort.Slice(keys, func(i, j int) bool {
		a, b := n.children[keys[i]], n.children[keys[j]]
		// Folders first
		if a.isDir && !b.isDir {
			return true
		}
		if !a.isDir && b.isDir {
			return false
		}
		return a.name < b.name
	})

	for _, k := range keys {
		child := n.children[k]
		// Add current node
		*items = append(*items, TreeItem{
			Path:     child.name,
			FullPath: child.fullPath,
			IsDir:    child.isDir,
			Depth:    depth,
		})

		// Recurse if directory
		if child.isDir {
			flatten(child, depth+1, items)
		}
	}
}

func getIcon(name string, isDir bool) string {
	if isDir {
		return " "
	}
	ext := filepath.Ext(name)
	switch strings.ToLower(ext) {
	case ".go":
		return " "
	case ".js", ".ts", ".tsx":
		return " "
	case ".svelte":
		return " "
	case ".md":
		return " "
	case ".json":
		return " "
	case ".yml", ".yaml":
		return " "
	case ".html":
		return " "
	case ".css":
		return " "
	case ".git":
		return " "
	case ".dockerfile":
		return " "
	default:
		return " "
	}
}
