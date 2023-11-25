package go_http

import (
	"errors"
	"strings"
)

type TrieRoot struct {
	root map[string]*Node
}

func (r *TrieRoot) AddTrieRoot(pattern string, handleFunc HandleFunc) *TrieRoot {
	patternCheck(pattern)
	if r.root == nil {
		r.root = make(map[string]*Node)
	}
	root, ok := r.root["/"]
	if !ok {
		root = &Node{
			part: "/",
		}
		r.root["/"] = root
	}
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, part := range parts {
		if part == "" {
			panic("空字符串不符合")
		}
		root = root.addNode(part)
	}
	root.handleFunc = handleFunc
	return nil
}

func (r *TrieRoot) GetTrieRoot(pattern string) (*Node, error) {
	patternCheck(pattern)
	root, ok := r.root["/"]
	if !ok {
		return nil, errors.New("root不存在")
	}
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	for _, part := range parts {
		if part == "" {
			return nil, errors.New("pattern有空格")
		}
		root = root.getNode(part)
		if root == nil {
			return nil, errors.New("不存在")
		}
	}
	if root.handleFunc == nil {
		return nil, errors.New("前缀正确，但是不存在handleFunc")
	}
	return root, nil
}

type Node struct {
	// 唯一标识
	part string
	// 静态路由后代
	children map[string]*Node
	// 数据
	handleFunc HandleFunc
	// 参数路由后代
	paramChild *Node
	// /study/golang | /study/:course | /study/:course/test
	// 如果查找/study/java，先去静态路由golang那里查，找不到就去:course
	// 静态路由一对一，动态路由一对n，只需要一个节点
}

func (n *Node) addNode(part string) *Node {
	if strings.HasPrefix(part, ":") {
		if n.paramChild == nil {
			child := &Node{
				part: part,
			}
			n.paramChild = child
			return child
		}
		return n.paramChild
	}
	if n.children == nil {
		n.children = make(map[string]*Node)
	}
	child, ok := n.children[part]
	if !ok {
		child = &Node{
			part:     part,
			children: nil,
		}
		n.children[part] = child
	}
	return child
}

func (n *Node) getNode(part string) *Node {
	if n.children == nil && n.paramChild == nil {
		return nil
	}
	if n.children != nil {
		node, ok := n.children[part]
		if !ok {
			return nil
		}
		return node
	}
	return n.paramChild
}

func patternCheck(pattern string) {
	if !strings.HasPrefix(pattern, "/") {
		panic("pattern不是以/开头")
	}
	if strings.HasSuffix(pattern, "/") {
		panic("pattern最后有/")
	}
}
