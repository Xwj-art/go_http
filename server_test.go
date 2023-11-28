package go_http

import (
	"fmt"
	"net/http"
	"testing"
)

func Login(ctx *Context) {
	//ctx.Response.Write([]byte("This is Loign"))
	ctx.MyResponse.TEXT(http.StatusOK, fmt.Sprintf("text -- login"))
}

func Register(ctx *Context) {
	ctx.MyResponse.TEXT(http.StatusOK, fmt.Sprintf("text -- register"))
}

func ParamIndex(ctx *Context) {
	course, err := ctx.GetParams("param")
	if err != nil {
		ctx.MyResponse.TEXT(http.StatusNotFound, "param错误")
	}
	ctx.MyResponse.TEXT(http.StatusOK, fmt.Sprintf("%s -- %s", ctx.Pattern, course))
}

func All(ctx *Context) {
	course, err := ctx.GetParams("all")
	if err != nil {
		ctx.MyResponse.TEXT(http.StatusNotFound, "param错误")
	}
	ctx.MyResponse.TEXT(http.StatusOK, fmt.Sprintf("%s -- %s", ctx.Pattern, course))
}

func TestHTTPServer_Start(t *testing.T) {
	h := NewHTTPServer(WithHTTPServerStop(nil))
	h.GET("/login", Login)
	h.POST("/register", Register)
	h.POST("/study/:param", ParamIndex)
	h.POST("/all/*all", All)
	err := h.Start(":8080")
	if err != nil {
		panic(err)
	}
}
