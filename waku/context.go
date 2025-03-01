package waku

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	// origin context
	Writer  http.ResponseWriter
	Request *http.Request
	// request info
	Path   string
	Method string
	Params map[string]string
	// response info
	StatusCode int
	// middleware
	handlers []handleFunc
	index    int
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func NewContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer:  w,
		Request: r,
		Path:    r.URL.Path,
		Method:  r.Method,
		index:   -1}
}

type H map[string]interface{}

type T struct {
	username string
	password string
}

func (c *Context) Next() {
	c.index++
	s := len(c.handlers)
	// 调用当前idx到len的所有handlers
	for ; c.index < s; c.index++ {
		c.handlers[c.index](c)
	}
}

// PostForm
// 获取post请求中postForm中key字段的value
func (c *Context) PostForm(key string) string {
	return c.Request.FormValue(key)
}

// PostJSON
// 解析post 请求中的 json数据
func (c *Context) PostJSON(dest T) error {
	// 检查请求的 Content-Type 是否为 application/json
	contentType := c.Request.Header.Get("Content-Type")
	if contentType != "application/json" {
		return fmt.Errorf("expected Content-Type 'application/json', got %s", contentType)
	}
	// 解析 JSON 请求体
	err := json.NewDecoder(c.Request.Body).Decode(dest)
	if err != nil {
		return err
	}
	return nil
}

// Query
// 获得get请求query中key字段的value
func (c *Context) Query(key string) string {
	return c.Request.URL.Query().Get(key)
}

// Status
// 设置返回的 status code
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

// return string as result
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}

// JSON return JSON as result
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	err := encoder.Encode(obj)
	if err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

// HTML return html string template as result
func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	c.Writer.Write([]byte(html))
}

// Data return data blob as result
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func GenerateKey(method, path string) string {
	return method + "-" + path
}
