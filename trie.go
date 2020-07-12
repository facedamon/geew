package geew

import "strings"

type node struct {
	// 待匹配路由 like: /p/:lang
	pattern string
	// 路由关键词 like: p, :lang
	part string
	// 子节点 like [doc, info]
	children []*node
	// 是否精确匹配, part含有: 或 * 是为true
	isWild bool
}

// 第一次匹配part成功的节点
func (n *node) matchChild(part string) *node {
	for _, c := range n.children {
		if c.part == part || c.isWild {
			return c
		}
	}
	return nil
}

// 匹配所有part成功的节点
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, c := range n.children {
		if c.part == part || c.isWild {
			nodes = append(nodes, c)
		}
	}
	return nodes
}

// 插入路由part
// @param psttern like : /p/
// @param parts like : p, :lang
// @param height : the node level of tree
func (n *node) insert(pattern string, parts []string, height int) {
	if len(parts) == height {
		n.pattern = pattern
		return
	}

	part := parts[height]
	child := n.matchChild(part)
	if nil == child {
		child = &node{part: part, isWild: part[0] == ':' || part[0] == '*'}
		n.children = append(n.children, child)
	}
	child.insert(pattern, parts, height+1)
}

// 路由匹配
// @param parts like : p, :lang
// @param height : the node level of tree
func (n *node) search(parts []string, height int) *node {
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		// the root
		return n
	}
	part := parts[height]
	children := n.matchChildren(part)

	for _, c := range children {
		r := c.search(parts, height+1)
		if r != nil {
			return r
		}
	}
	return nil
}
