package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	Method string
	Path   string
	// NEW: Params 字段
	// Context 结构体新增了 Params 字段，用于存储动态路由匹配的参数
	// 在路径中可能会包含一些动态的部分，比如 /users/:id/profile，其中 :id 就是动态参数。当用户请求 /users/123/profile 时，:id 就会被解析为 id: "123" 并存储在 Params 中。
	// 这个 Params 映射会由 router 中的 getRoute 方法负责填充（在解析路由时），使得动态参数可以在后续处理函数中方便地访问。
	Params     map[string]string
	StatusCode int
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Method: req.Method,
		Path:   req.URL.Path,
	}
}

// NEW: Param 方法用于获取动态路由中的参数值
// 通过该方法可以方便地访问URL中的动态部分
func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) PostFrom(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) String(code int, format string, value ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, value...)))
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}

func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.Writer.Write([]byte(html))
}
