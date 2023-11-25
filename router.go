package go_http

import (
	"errors"
	"fmt"
	"strings"
)

type router struct {
	trees map[string]*TrieRoot
}

func newRouter() *router {
	return &router{make(map[string]*TrieRoot)}
}

func (r *router) addRouter(method, pattern string, handleFunc HandleFunc) {
	fmt.Printf("add router %s - %s\n", method, pattern)
	if r.trees == nil {
		r.trees = make(map[string]*TrieRoot)
	}
	trieRoot := r.trees[method]
	if trieRoot == nil {
		trieRoot = &TrieRoot{
			root: make(map[string]*Node),
		}
		r.trees[method] = trieRoot
	}
	trieRoot.AddTrieRoot(pattern, handleFunc)
}

func (r *router) getRouter(method, pattern string) (*Node, map[string]string, error) {
	params := make(map[string]string)
	fmt.Printf("request router %s - %s\n", method, pattern)
	if r.trees == nil {
		return nil, nil, errors.New("尚未添加路由")
	}
	trieRoot := r.trees[method]
	node, err := trieRoot.GetTrieRoot(pattern)
	if err != nil {
		return nil, nil, err
	}
	if strings.HasPrefix(node.part, ":") {
		params[node.part[1:]] = pattern
		return node, params, nil
	}
	return node, nil, nil
}
