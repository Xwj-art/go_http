package go_http

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
)

type H map[string]any

type Context struct {
	MyResponse *MyResponse
	Request    *http.Request
	Method     string
	Pattern    string
	Params     map[string]string
	cacheQuery url.Values
	cacheBody  io.ReadCloser
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		MyResponse: &MyResponse{
			ResponseWriter: w,
		},
		Request: r,
		Method:  r.Method,
		Pattern: r.URL.Path,
	}
}

func (c *Context) GetParams(key string) (string, error) {
	value, ok := c.Params[key]
	if !ok {
		return "", errors.New("获取param失败")
	}
	return value, nil
}

func (c *Context) Query(key string) (string, error) {
	if c.cacheQuery == nil {
		c.cacheQuery = c.Request.URL.Query()
	}
	value, ok := c.cacheQuery[key]
	if !ok {
		return "", errors.New("查询的参数不存在")
	}
	return value[0], nil
}

func (c *Context) Form(key string) (string, error) {
	if c.cacheBody == nil {
		c.cacheBody = c.Request.Body
	}
	err := c.Request.ParseForm()
	if err != nil {
		return "", errors.New("表单解析失败")
	}
	value := c.Request.FormValue(key)
	return value, nil
}

func (c *Context) BindJSON(dest any) error {
	if c.cacheBody == nil {
		c.cacheBody = c.Request.Body
	}
	decoder := json.NewDecoder(c.cacheBody)
	decoder.DisallowUnknownFields()
	return decoder.Decode(dest)
}

func (c *Context) JSON(code int, data any) {
	c.MyResponse.SetStatusCode(code)
	c.MyResponse.SetHeader("Context-Type", "application/json")
	tmp, err := json.Marshal(data)
	if err != nil {
		c.MyResponse.SetStatusCode(http.StatusInternalServerError)
		c.MyResponse.DelHeader("Context-Type")
		panic("json.Marshal调用失败")
	}
	c.MyResponse.SetData(tmp)
}
func (c *Context) HTML(code int, html string) {
	c.MyResponse.SetStatusCode(code)
	c.MyResponse.SetHeader("Context-Type", "text/html")
	c.MyResponse.SetData([]byte(html))
}
func (c *Context) TEXT(code int, text string) {
	c.MyResponse.SetStatusCode(code)
	c.MyResponse.SetHeader("Context-Type", "text/plain")
	c.MyResponse.SetData([]byte(text))
}

func (c *Context) flashDataToResponse() {
	// 顺序特定，不可更改
	c.MyResponse.WriteHeader(c.MyResponse.statusCode)
	for key, value := range c.MyResponse.header {
		c.MyResponse.Header().Set(key, value)
	}
	c.MyResponse.Write(c.MyResponse.Data)
}

// MyResponse 将MyResponse在底层封装，并暴露几个改变属性的api给Context用

type MyResponse struct {
	http.ResponseWriter
	header     map[string]string
	statusCode int
	Data       []byte
}

func (r *MyResponse) SetStatusCode(code int) { r.statusCode = code }
func (r *MyResponse) SetHeader(key, value string) {
	if r.header == nil {
		r.header = make(map[string]string)
	}
	r.header[key] = value
}
func (r *MyResponse) DelHeader(key string) { delete(r.header, key) }
func (r *MyResponse) SetData(data []byte)  { r.Data = data }
