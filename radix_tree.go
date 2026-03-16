package gocket

import (
	"iter"
)

func newRadixTree() radixTree {
	return radixTree{
		root: radixTreeNode{
			path:          pathNode{},
			pathChildren:  []*radixTreeNode{},
			wildcardChild: optional[*radixTreeNode]{},
			end:           []Route{},
		},
	}
}

type radixTree struct {
	root radixTreeNode
}

func (tree radixTree) Routes() iter.Seq[*Route] {
	return tree.root.Routes
}

type optional[T any] struct {
	val    T
	isSome bool
}
type radixTreeNode struct {
	path          pathNode
	pathChildren  []*radixTreeNode
	wildcardChild optional[*radixTreeNode]
	// if empty this is not a valid endpoint
	// There can be multiple routes for a single path and they should be tried in the order in this array
	// This is because of queries and MiddleWare guards
	end []Route
}
type pathNode struct {
	// if this is a part of radixTreeNode.wildcardChildren then it doesn't matter what this is set to
	path string
}

func (tree *radixTree) addRoute(route Route) {
	path := route.Pattern.Path.Parts
	tree.root.addRoute(path, route)
}

func (node *radixTreeNode) addRoute(path []PathElement, route Route) {
	if len(path) == 0 {
		node.end = append(node.end, route)
		return
	}
	part := path[0]
	// try to find it
	if part.Wildcard {
		if !node.wildcardChild.isSome {
			newNode := radixTreeNode{
				pathNode{},
				make([]*radixTreeNode, 0),
				optional[*radixTreeNode]{&radixTreeNode{}, false},
				make([]Route, 0),
			}
			node.wildcardChild = optional[*radixTreeNode]{
				&newNode,
				true,
			}
		}
		node.wildcardChild.val.addRoute(path[1:], route)
		return
	}
	for _, child := range node.pathChildren {
		if part.Name == child.path.path {
			child.addRoute(path[1:], route)
			return
		}
	}
	newNode := radixTreeNode{
		pathNode{part.Name},
		make([]*radixTreeNode, 0),
		optional[*radixTreeNode]{&radixTreeNode{}, false},
		make([]Route, 0),
	}
	newNode.addRoute(path[1:], route)
	node.pathChildren = append(node.pathChildren, &newNode)
}

func (tree *radixTree) matchPath(path []string) []*Route {
	return tree.root.matchPath(path)
}

func (node *radixTreeNode) matchPath(path []string) []*Route {
	if len(path) == 0 {
		matchedPaths := make([]*Route, len(node.end))
		for i := range node.end {
			matchedPaths[i] = &node.end[i]
		}
		return matchedPaths
	}
	matchedPaths := make([]*Route, 0)
	part := path[0]
	for _, child := range node.pathChildren {
		if child.path.path == part {
			matchedRoutes := child.matchPath(path[1:])
			matchedPaths = append(matchedPaths, matchedRoutes...)
		}
	}
	if node.wildcardChild.isSome {
		matchedRoutes := node.wildcardChild.val.matchPath(path[1:])
		matchedPaths = append(matchedPaths, matchedRoutes...)
	}
	return matchedPaths
}

func (node *radixTreeNode) Routes(yield func(*Route) bool) {
	for _, route := range node.end {
		if !yield(&route) {
			return
		}
	}
	for _, child := range node.pathChildren {
		for route := range child.Routes {
			yield(route)
		}
	}
	if !node.wildcardChild.isSome {
		return
	}
	for route := range node.wildcardChild.val.Routes {
		yield(route)
	}
}
