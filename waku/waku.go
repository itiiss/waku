package waku

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"path"
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
	// for html render
	htmlTemplate *template.Template
	funcMap      template.FuncMap
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

// SetFuncMap  自定义的模版渲染函数
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.funcMap = funcMap
}

// LoadHTMLGlob 组合 模版 + 自定义的模版渲染函数
func (engine *Engine) LoadHTMLGlob(pattern string) {
	engine.htmlTemplate = template.Must(template.New("").Funcs(engine.funcMap).ParseGlob(pattern))
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

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) handleFunc {
	// 对于静态文件的path，拼上group的prefix得到绝对路径
	absolutePath := path.Join(group.prefix, relativePath)
	// 将文件绝对路径和根路径拼成最终路径
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		_, err := fs.Open(file)

		if err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (group *RouterGroup) Static(relativePath string, root string) {
	// 通过最终路径，得到一个返回该路径文件的 handler
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	// 得到该文件的相对路径作为 key
	urlPattern := path.Join(relativePath, "/*filepath")
	// 注册进路由表中
	group.Get(urlPattern, handler)
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
	c.engine = engine
	engine.router.handle(c)
}
