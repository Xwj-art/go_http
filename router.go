package go_http

import (
	"errors"
	"fmt"
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

func (r *router) getRouter(method, pattern string) (*Node, error) {
	fmt.Printf("request router %s - %s\n", method, pattern)
	if r.trees == nil {
		return nil, errors.New("尚未添加路由")
	}
	trieRoot := r.trees[method]
	node, err := trieRoot.GetTrieRoot(pattern)
	if err != nil {
		return nil, errors.New("获取路由失败")
	}
	return node, nil
}
