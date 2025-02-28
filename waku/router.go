package waku

import "net/http"

type Router struct {
	handlers map[string]handleFunc
}

func NewRouter() *Router {
	return &Router{handlers: make(map[string]handleFunc)}
}

// 添加key = 方法-路径和 value = handler 到routes map中
func (engine *Engine) AddRoute(method, routePattern string, handleFunc handleFunc) {
	key := GenerateKey(method, routePattern)
	engine.router.handlers[key] = handleFunc
}

func (r *Router) handle(c *Context) {
	key := GenerateKey(c.Method, c.Path)
	handler, ok := r.handlers[key]
	if ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 page not found", c.Path)
	}

}
