package waku

import (
	"fmt"
	"net/http"
)

type handleFunc func(w http.ResponseWriter, r *http.Request)

// 实现serveHttp方法后成为一个通用handler
type Engine struct {
	routes map[string]handleFunc
}

// NewEngine 构造函数
func NewEngine() *Engine {
	return &Engine{routes: make(map[string]handleFunc)}
}

// 添加key = 方法-路径和 value = handler 到routes map中
func (engine *Engine) addRoute(method, routePattern string, handleFunc handleFunc) {
	key := generateKey(method, routePattern)
	engine.routes[key] = handleFunc
}

func (engine *Engine) Get(routePattern string, handleFunc handleFunc) {
	engine.addRoute(http.MethodGet, routePattern, handleFunc)
}

func (engine *Engine) Post(routePattern string, handleFunc handleFunc) {
	engine.addRoute(http.MethodPost, routePattern, handleFunc)
}

// Run 启动serve监听addr端口
func (engine *Engine) Run(addr string) error {
	fmt.Printf("waku server listening at %s\n", addr)
	return http.ListenAndServe(addr, engine)
}

// 通过req中的 method和path，找到对应的handler来处理请求，找不到则返回404
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := generateKey(req.Method, req.URL.Path)
	fmt.Printf("waku server handle request %s\n", key)
	if handler, ok := engine.routes[key]; ok {
		handler(w, req)
	} else {
		http.NotFound(w, req)
	}
}

func generateKey(method, path string) string {
	return method + "-" + path
}
