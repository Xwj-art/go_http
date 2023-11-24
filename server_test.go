package go_http

import (
	"testing"
)

func Login(ctx *Context) {
	ctx.Response.Write([]byte("This is Loign"))
}

func Register(ctx *Context) {
	ctx.Response.Write([]byte("This is Register"))
}

func TestHTTPServer_Start(t *testing.T) {
	h := NewHTTPServer(WithHTTPServerStop(nil))
	h.GET("/login", Login)
	h.POST("/register", Register)
	err := h.Start(":8080")
	if err != nil {
		panic(err)
	}
}
