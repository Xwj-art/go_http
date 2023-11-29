package go_http

import (
	"fmt"
	"net/http"
	"testing"
)

func Login(ctx *Context) {
	//ctx.Response.Write([]byte("This is Loign"))
	ctx.JSON(http.StatusOK, H{
		"name": "小明",
		"id":   114514,
	})
}

func Register(ctx *Context) {
	ctx.HTML(http.StatusOK, `<h1 style="color:red;"> 芝士一个html测试的标题 </h1>`)
}

func ParamIndex(ctx *Context) {
	course, err := ctx.GetParams("param")
	if err != nil {
		ctx.TEXT(http.StatusNotFound, "param错误")
	}
	ctx.TEXT(http.StatusOK, fmt.Sprintf("%s -- %s", ctx.Pattern, course))
}

func All(ctx *Context) {
	course, err := ctx.GetParams("all")
	if err != nil {
		ctx.TEXT(http.StatusNotFound, "param错误")
	}
	ctx.TEXT(http.StatusOK, fmt.Sprintf("%s -- %s", ctx.Pattern, course))
}

// http://localhost:8080/query?username=xxx&password=aaa
func Query(c *Context) {
	username, err := c.Query("username")
	if err != nil {
		panic("query错了")
	}
	password, err := c.Query("password")
	if err != nil {
		panic("query错了")
	}
	c.JSON(http.StatusOK, H{
		"username": username,
		"password": password,
		"id":       1,
	})
}

func Form(c *Context) {
	username, err := c.Form("username")
	if err != nil {
		panic("form错了")
	}
	password, err := c.Form("password")
	if err != nil {
		panic("form错了")
	}
	c.JSON(http.StatusOK, H{
		"username": username,
		"password": password,
		"id":       1,
	})
}

func BindJson(c *Context) {
	type User struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	var user User
	// 绑定处理传来的json
	if err := c.BindJSON(&user); err != nil {
		panic("json绑定错误")
	}
	c.JSON(http.StatusOK, H{
		"username": user.Username,
		"password": user.Password,
	})
}

func TestHTTPServer_Start(t *testing.T) {
	router := NewHTTPServer(WithHTTPServerStop(nil))
	router.GET("/login", Login)
	router.GET("/register", Register)
	router.GET("/study/:param", ParamIndex)
	router.GET("/all/*all", All)
	router.GET("/query", Query)
	router.POST("/form", Form)
	router.POST("/bind_json", BindJson)
	err := router.Start(":8080")
	if err != nil {
		panic(err)
	}
}
