package waku

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type handleFunc func(c *Context)

type RouterGroup struct {
	prefix     string       // 按path前缀分组
	middleware []handleFunc // 该分组需要的中间件
	parent     *RouterGroup // 直接父分组
	engine     *Engine      // 所有Group的根分组
}

// 实现serveHttp方法后成为一个通用handler
type Engine struct {
	router       *Router
	*RouterGroup // embedding RouterGroup  into Engine
	groups       []*RouterGroup
}

// NewEngine 构造函数
func NewEngine() *Engine {
	// init engine instance
	engine := &Engine{router: NewRouter()}
	// init RouterGroup instance and assign to engine.routerGroup
	engine.RouterGroup = &RouterGroup{engine: engine}
	// init groups array with [engine.routerGroup]
	// 将根group 初始化进groups slice中
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

// AddGroup 添加新的Group到engine中
func (group *RouterGroup) NewGroup(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine}
	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

// 基于Group添加路由
func (group *RouterGroup) AddRoute(method, comp string, handleFunc handleFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.AddRoute(method, pattern, handleFunc)
}

func (group *RouterGroup) Get(pattern string, handleFunc handleFunc) {
	group.AddRoute("GET", pattern, handleFunc)
}

func (group *RouterGroup) Post(pattern string, handleFunc handleFunc) {
	group.AddRoute("POST", pattern, handleFunc)
}

// 基于engine添加路由
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

// Use 注册middleware
func (group *RouterGroup) Use(middleware ...handleFunc) {
	group.middleware = append(group.middleware, middleware...)
}

// 通过req和res构建 context，handle函数再解析context获取信息，找到对应的handleFunc处理
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []handleFunc
	// 取得注册在group上的middlewares函数
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middleware...)
		}
	}
	c := NewContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
