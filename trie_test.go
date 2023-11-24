package go_http

// trie测试，handleFunc为stirng
//
//import (
//	"errors"
//	"github.com/stretchr/testify/assert"
//	"testing"
//)
//
//func TestTrieRoot_AddTrieRoot(t *testing.T) {
//	testCases := []struct {
//		name, pattern, handleFunc string
//		wantTrieRoot        *TrieRoot
//	}{
//		{
//			name:    "xxx",
//			pattern: "/user/login",
//			handleFunc:    "hello",
//			wantTrieRoot: &TrieRoot{
//				map[string]*Node{
//					"/": {
//						part: "/",
//						children: map[string]*Node{
//							"user": {
//								part: "user",
//								children: map[string]*Node{
//									"login": {
//										part: "login",
//										handleFunc: "hello",
//									},
//								},
//							},
//						},
//					},
//				},
//			},
//		},
//	}
//	router := &TrieRoot{
//		map[string]*Node{
//			"/": {part: "/"},
//		},
//	}
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			router.AddTrieRoot(tc.pattern, tc.handleFunc)
//			assert.Equal(t, tc.wantTrieRoot, router)
//		})
//	}
//}
//
//func TestTrieRoot_GetTrieRoot(t *testing.T) {
//	testCases := []struct {
//		name, findPattern, wantData string
//		wantErr                     error
//	}{
//		{
//			name:        "success",
//			findPattern: "/user/login",
//			wantData:    "hello",
//		},
//		{
//			name:        "fail1",
//			findPattern: "/user//register",
//			wantData:    "world",
//			wantErr:     errors.New("pattern有空格"),
//		},
//		{
//			name:        "fail2",
//			findPattern: "/selsadsjior/register",
//			wantErr:     errors.New("不存在"),
//		},
//	}
//	router := &TrieRoot{
//		map[string]*Node{
//			"/": {part: "/"},
//		},
//	}
//	router.AddTrieRoot("/user/login", "hello")
//	router.AddTrieRoot("/user/register", "world")
//	router.AddTrieRoot("/study/golang", "good")
//	router.AddTrieRoot("/study/python", "aaa")
//	for _, tc := range testCases {
//		t.Run(tc.name, func(t *testing.T) {
//			n, err := router.GetTrieRoot(tc.findPattern)
//			assert.Equal(t, tc.wantErr, err)
//			if err != nil {
//				return
//			}
//			assert.Equal(t, tc.wantData, n.handleFunc)
//		})
//	}
//}
