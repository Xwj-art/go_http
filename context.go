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

func (r *MyResponse) JSON(code int, data any) {
	r.SetStatusCode(code)
	r.SetHeader("Context-Type", "application/json")
	tmp, err := json.Marshal(data)
	if err != nil {
		r.SetStatusCode(http.StatusInternalServerError)
		r.DelHeader("Context-Type")
		panic("json.Marshal调用失败")
	}
	r.SetData(tmp)
}
func (r *MyResponse) HTML(code int, html string) {
	r.SetStatusCode(code)
	r.SetHeader("Context-Type", "text/html")
	r.SetData([]byte(html))
}
func (r *MyResponse) TEXT(code int, text string) {
	r.SetStatusCode(code)
	r.SetHeader("Context-Type", "text/plain")
	r.SetData([]byte(text))
}

func (r *MyResponse) flashDataToResponse() {
	// 顺序特定，不可更改
	r.WriteHeader(r.statusCode)
	for key, value := range r.header {
		r.Header().Set(key, value)
	}
	r.Write(r.Data)
}
