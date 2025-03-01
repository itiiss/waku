package waku

import (
	"net/http"
	"strings"
)

type Router struct {
	roots    map[string]*node // key分别是GET，POST，DELETE，PUT这四个方法，value为Trie树的root节点
	handlers map[string]handleFunc
}

func NewRouter() *Router {
	return &Router{
		handlers: make(map[string]handleFunc),
		roots:    make(map[string]*node),
	}
}

// only can handle one *
func parsePattern(pattern string) []string {
	vs := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range vs {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

// AddRoute 添加key = 方法-路径和 value = handler 到routes map中
func (r *Router) AddRoute(method, routePattern string, handleFunc handleFunc) {
	parts := parsePattern(routePattern)
	key := GenerateKey(method, routePattern)

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}

	r.roots[method].insert(routePattern, parts, 0)
	r.handlers[key] = handleFunc
}

func (r *Router) getRoute(method, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)

	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	// 找到method和path都匹配的节点
	node := root.search(searchParts, 0)
	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			// 提取路由参数到params中
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}
	return nil, nil
}

func (r *Router) handle(c *Context) {

	node, params := r.getRoute(c.Method, c.Path)

	if node != nil {
		c.Params = params
		// node.pattern将动态路由path的:id 换成具体 123
		key := GenerateKey(c.Method, node.pattern)
		handleFunc := r.handlers[key]
		// 将业务handler放在middlewares的最后
		c.handlers = append(c.handlers, handleFunc)
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.String(http.StatusNotFound, "404 page not found", c.Path)
		})
	}
	// 从头开始执行middlewares中的func
	c.Next()
}
