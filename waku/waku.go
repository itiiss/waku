package waku

import (
	"fmt"
	"net/http"
)

type handleFunc func(c *Context)

// 实现serveHttp方法后成为一个通用handler
type Engine struct {
	router *Router
}

// NewEngine 构造函数
func NewEngine() *Engine {
	return &Engine{router: NewRouter()}
}

func (engine *Engine) Get(routePattern string, handleFunc handleFunc) {
	engine.router.AddRoute(http.MethodGet, routePattern, handleFunc)
}

func (engine *Engine) Post(routePattern string, handleFunc handleFunc) {
	engine.router.AddRoute(http.MethodPost, routePattern, handleFunc)
}

// Run 启动serve监听addr端口
func (engine *Engine) Run(addr string) error {
	fmt.Printf("waku server listening at %s\n", addr)
	return http.ListenAndServe(addr, engine)
}

// 通过req中的 method和path，找到对应的handler来处理请求，找不到则返回404
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := NewContext(w, req)
	engine.router.handle(c)
}
