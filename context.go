package go_http

import "net/http"

type Context struct {
	Response http.ResponseWriter
	Request  *http.Request
	Method   string
	Pattern  string
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Response: w,
		Request:  r,
		Method:   r.Method,
		Pattern:  r.URL.Path,
	}
}
