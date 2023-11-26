package go_http

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouter(t *testing.T) {
	testCases1 := []struct {
		name, pattern, method, wantAns string
		wantErr                        error
	}{
		{
			name:    "test_1_1",
			pattern: "/study/:course",
			method:  "GET",
		},
		{
			name:    "test_1_2",
			pattern: "/study/course/golang",
			method:  "GET",
		},
		{
			name:    "test_1_3",
			pattern: "/study/golang",
			method:  "GET",
		},
		{
			name:    "test_1_4",
			pattern: "/study/:course/golang",
			method:  "GET",
		},
		{
			name:    "test_1_5",
			pattern: "/study/:course/java",
			method:  "GET",
		},
		{
			name:    "test_1_6",
			pattern: "/test/test/test/*course",
			method:  "GET",
		},
		{
			name:    "test_1_7",
			pattern: "/test/*course",
			method:  "GET",
		},
	}
	testCases2 := []struct {
		name, pattern, method, wantAns string
		wantErr                        error
	}{
		{
			name:    "test_2_1",
			pattern: "/study/course",
			method:  "GET",
			wantErr: errors.New("前缀正确，但是不存在handleFunc"),
		},
		{
			name:    "test_2_2",
			pattern: "/study/course/golang",
			method:  "GET",
			wantAns: "golang",
		},
		{
			name:    "test_2_3",
			pattern: "/study/golang",
			method:  "GET",
			wantAns: "golang",
		},
		{
			name:    "test_2_4",
			pattern: "/study/param/java",
			method:  "GET",
			wantAns: "java",
		},
		{
			name:    "test_2_5",
			pattern: "/test/course",
			method:  "GET",
			wantAns: "course",
		},
		{
			name:    "test_2_6",
			pattern: "/test/test/test/test/test/test",
			method:  "GET",
			wantAns: "test/test/test",
		},
		{
			name:    "test_2_7",
			pattern: "/test/course/course",
			method:  "GET",
			wantAns: "course/course",
		},
	}
	r := newRouter()
	handleFunc := func(ctx *Context) {}
	for _, tc := range testCases1 {
		t.Run(tc.name, func(t *testing.T) {
			r.addRouter(tc.method, tc.pattern, handleFunc)
		})
	}
	for _, tc := range testCases2 {
		t.Run(tc.name, func(t *testing.T) {
			node, params, err := r.getRouter(tc.method, tc.pattern)
			if err != nil {
				assert.Equal(t, tc.wantErr, err)
				return
			}
			if params != nil {
				assert.Equal(t, tc.wantAns, params[node.part[1:]])
				return
			}
			assert.Equal(t, tc.wantAns, node.part)
		})
	}
}
