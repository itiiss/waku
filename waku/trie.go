package waku

import (
	"log"
	"strings"
)

type node struct {
	pattern  string  // 待匹配的路由整体
	part     string  // 当前路由节点
	children []*node // 下一级路由的所有节点
	isWild   bool    // 是否精准匹配
}

// 第一个匹配成功的节点，用于定位插入
func (n *node) matchChild(part string) *node {
	for _, child := range n.children {
		if child.part == part || child.isWild {
			return child
		}
	}
	return nil
}

// 所有匹配节点，用于查找结果
func (n *node) matchChildren(part string) []*node {
	nodes := make([]*node, 0)
	for _, child := range n.children {
		if child.part == part || child.isWild {
			nodes = append(nodes, child)
		}
	}
	return nodes
}

func (n *node) insert(pattern string, parts []string, height int) {
	// base情况，到达叶子节点

	if len(parts) == height {
		n.pattern = pattern
		return
	}

	// 取出当前层的part，查找children中匹配的节点用来插入
	part := parts[height]
	child := n.matchChild(part)

	log.Printf("insert pattern: %s,part: %s, height: %v", pattern, part, height)

	//未找到匹配节点，在children中添加一个新节点
	if child == nil {
		child = &node{
			part:   part,
			isWild: part[0] == ':' || part[0] == '*',
		}
		n.children = append(n.children, child)
	}

	// 找到匹配节点，以该节点为根节点继续查找
	child.insert(pattern, parts, height+1)
}

func (n *node) search(parts []string, height int) *node {

	// base情况，当前节点是叶子节点或者*通配节点
	if len(parts) == height || strings.HasPrefix(n.part, "*") {
		if n.pattern == "" {
			return nil
		}
		return n
	}

	// 取出当前层的part，查找children中所有匹配的节点
	part := parts[height]
	children := n.matchChildren(part)

	// 以每个匹配子节点作为根节点，递归的进行search，直到到达base情况
	for _, child := range children {
		result := child.search(parts, height+1)
		if result != nil {
			return result
		}
	}
	return nil
}
