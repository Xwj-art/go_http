package go_http

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestRouterAdd(t *testing.T) {
	testCases := []struct {
		name, method, pattern, wantErr string
	}{
		{
			name:    "test1_add",
			method:  "GET",
			pattern: "/study/java",
		},
		{
			name:    "test2_add",
			method:  "POST",
			pattern: "java",
			wantErr: "pattern不是以/开头",
		},
	}
	r := newRouter()
	tmpFunc := func(ctx *Context) {}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r.addRouter(tc.method, tc.pattern, tmpFunc)
			//assert.PanicsWithError(t, tc.wantErr, func() {})
		})
	}
}

func TestRouterGet(t *testing.T) {
	testCases := []struct {
		name, method, pattern string
		wantErr               error
	}{
		{
			name:    "test1_get",
			method:  "POST",
			pattern: "/study/java",
		},
		{
			name:    "test2_get",
			method:  "GET",
			pattern: "/study//java",
			wantErr: errors.New("pattern有空格"),
		},
		{
			name:    "test3_get",
			method:  "POST",
			pattern: "/study",
			wantErr: errors.New("前缀正确，但是不存在handleFunc"),
		},
	}
	r := newRouter()
	tmpFunc := func(ctx *Context) {}
	r.addRouter("POST", "/study/java", tmpFunc)
	r.addRouter("GET", "/study/java", tmpFunc)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, _, err := r.getRouter(tc.method, tc.pattern)
			assert.Equal(t, tc.wantErr, err)
		})
	}
}

func TestParamRouterAdd(t *testing.T) {
	testCases := []struct {
		name, method, pattern, wantErr string
	}{
		{
			name:    "test1_param_add",
			method:  "GET",
			pattern: "/study/:course",
		},
	}
	r := newRouter()
	tmpFunc := func(ctx *Context) {}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			r.addRouter(tc.method, tc.pattern, tmpFunc)
		})
	}
}

func TestParamRouterGet(t *testing.T) {
	testCases := []struct {
		name, method, pattern, wantAnswer string
	}{
		{
			name:       "test1_param_get",
			method:     "GET",
			pattern:    "/study/java",
			wantAnswer: ":course",
		},
		{
			name:       "test2_param_get",
			method:     "GET",
			pattern:    "/study/java/answer",
			wantAnswer: "answer",
		},
		{
			name:       "test3_param_get",
			method:     "GET",
			pattern:    "/study/golang",
			wantAnswer: ":course",
		},
	}
	r := newRouter()
	tmpFunc := func(ctx *Context) {}
	r.addRouter("GET", "/study/:course", tmpFunc)
	r.addRouter("GET", "/study/:course/answer", tmpFunc)
	//r.addRouter("GET", "/study/golang", tmpFunc)
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			n, _, _ := r.getRouter(tc.method, tc.pattern)
			assert.NotNil(t, n)
			assert.Equal(t, tc.wantAnswer, n.part)
			//assert.True(t, strings.HasPrefix(n.part, ":"))
		})
	}
}
