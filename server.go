package go_http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type HTTPOption func(h *HTTPServer)
type HandleFunc func(ctx *Context)

// 用server封装net/http
// Handler为一个处理的函数，对原始的Handle进一步封装
// server应该有哪些功能?   开关，注册router

type server interface {
	// Handler 本质上也是一个接口，接口内置
	http.Handler // 底层是ServeHTTP()，实质上是对ServeHTTP的封装
	Start(addr string) error
	End() error
	// 非常核心的api不能暴露给外人。
	addRouter(method, pattern string, handleFunc HandleFunc)
}

type HTTPServer struct {
	srv  *http.Server
	stop func() error
	// 先维护一个map，使得特定的方法+访问路径有特定的func
	//routers map[string]HandleFunc
	*router
}

// HTTPServer的构造函数

func NewHTTPServer(opts ...HTTPOption) *HTTPServer {
	h := &HTTPServer{
		router: newRouter(),
	}
	for _, opt := range opts {
		opt(h)
	}
	return h
}

// 设置HTTPServer的关闭，这是opts中的一个值，所以在初始化时直接赋值
// 直接返回opt即可(优雅关闭)

func WithHTTPServerStop(fn func() error) HTTPOption {
	return func(h *HTTPServer) {
		if fn == nil {
			fn = func() error {
				fmt.Println("123123123")
				quit := make(chan os.Signal)
				signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
				<-quit
				log.Println("Shutdown Server...")

				ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()
				if err := h.srv.Shutdown(ctx); err != nil {
					log.Fatal("Server Shutdown...", err)
				}
				select {
				case <-ctx.Done():
					log.Println("timeout of 5s")
				}
				return nil
			}
		}
		h.stop = fn
	}
}

type HTTPS struct {
}

func (h *HTTPServer) Start(addr string) error {
	h.srv = &http.Server{
		Addr:    addr,
		Handler: h,
	}
	return h.srv.ListenAndServe()
	// 此处的h是http.handler(), 而http.handler是serveHTTP()
	// 本质上都是对serveHTTP()的调用，所以直接调用自己就行了
	// http.Handler是一个接口，而HTTP实现了这个接口
	// return http.ListenAndServe(addr, h)
}

func (h *HTTPServer) Stop() error {
	return h.stop()
}

// 如果定义了ServeHTTP，那么HTTPServer就是一个handler，就可以直接使用HTTPServer
// 进行多路复用，不用自己写
// 如果没有定义ServeHTTP，还想要用HTTPServer，就得写个ServerMux，不用HTTPServer的话，
// 得使用net自带的ListenAndServe，其内部进行了ServerMux的组合，最后由ListenAndServe启动

func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// 每一个请求都有一个上下文，只有请求来了才初始化。
	node, params, err := h.getRouter(r.Method, r.URL.Path)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("404 NOT FOUND"))
		return
	}
	// 构造上下文
	c := NewContext(w, r)
	c.Params = params
	node.handleFunc(c)
}

func (h *HTTPServer) GET(pattern string, handle HandleFunc) {
	h.addRouter(http.MethodGet, pattern, handle)
}

func (h *HTTPServer) POST(pattern string, handle HandleFunc) {
	h.addRouter(http.MethodPost, pattern, handle)
}

func (h *HTTPServer) DELETE(pattern string, handle HandleFunc) {
	h.addRouter(http.MethodDelete, pattern, handle)
}

func (h *HTTPServer) PUT(pattern string, handle HandleFunc) {
	h.addRouter(http.MethodPut, pattern, handle)
}

func (h *HTTPS) Start(addr string) error {
	return nil
}

func (h *HTTPS) End() error {
	return nil
}

func main() {
	//http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {})
	// 两者启动方式不一样，为了简化上层使用，需要进一步抽象
	// http
	//http.ListenAndServe(":8080", nil)
	// https
	//http.ListenAndServeTLS(":8080", "certFile", "keyFile", nil)
}
