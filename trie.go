package go_http

import (
	"errors"
	"fmt"
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
		root, ok = root.addNode(part)
		if !ok {
			panic(fmt.Sprintf("add router %s 时冲突", pattern))
		}
	}
	if root.handleFunc != nil {
		panic(fmt.Sprintf("add router %s 时冲突", pattern))
	}
	root.handleFunc = handleFunc
	return nil
}

func (r *TrieRoot) GetTrieRoot(pattern string) (*Node, map[string]string, error) {
	patternCheck(pattern)
	params := make(map[string]string)
	flag := 1
	root, ok := r.root["/"]
	if !ok {
		return nil, nil, errors.New("root不存在")
	}
	parts := strings.Split(strings.Trim(pattern, "/"), "/")
	part := ""
	for _, part = range parts {
		if part == "" {
			return nil, nil, errors.New("pattern有空格")
		}
		root = root.getNode(part)
		if root == nil {
			return nil, nil, errors.New("不存在")
		}
		// add router /test
		// add router /test/*test
		// request router /test/test/test
		// expected: /test/test
		// actual: /test/test/test
		// 最多只有一个'*'
		// request router /test
		// expected: /test
		// actual: /test
		if strings.HasPrefix(root.part, "*") {
			params[root.part[1:]] = pattern[flag:]
			return root, params, nil
		}
		flag = flag + len(part) + 1
	}
	if root.handleFunc == nil {
		return nil, nil, errors.New("前缀正确，但是不存在handleFunc")
	}
	if strings.HasPrefix(root.part, ":") {
		params[root.part[1:]] = part
		return root, params, nil
	}
	return root, nil, nil
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
	starChild *Node
}

func (n *Node) addNode(part string) (*Node, bool) {
	if strings.HasPrefix(part, "*") {
		// 如果有paramChild，那么不管有什么路由都可以匹配上
		// 不能再注册starChild
		if n.paramChild != nil {
			return nil, false
		}
		// 原来starChild存在则直接替换
		child := &Node{
			part: part,
		}
		n.starChild = child
		return child, true
	}
	if strings.HasPrefix(part, ":") {
		if n.starChild != nil {
			return nil, false
		}
		if n.paramChild != nil {
			if n.paramChild.part != part {
				return nil, false
			}
		}
		n.paramChild = &Node{part: part}
		return n.paramChild, true
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
	return child, true
}

func (n *Node) getNode(part string) *Node {
	if n.children != nil {
		if n.children[part] != nil {
			return n.children[part]
		}
	}
	if n.paramChild != nil {
		return n.paramChild
	}
	if n.starChild != nil {
		return n.starChild
	}
	return nil
}

func patternCheck(pattern string) {
	if !strings.HasPrefix(pattern, "/") {
		panic("pattern不是以/开头")
	}
	if strings.HasSuffix(pattern, "/") {
		panic("pattern最后有/")
	}
}
