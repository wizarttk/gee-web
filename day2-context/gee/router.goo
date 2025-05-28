package gee

import (
	"net/http"
)

// router 结构体用于保存所有路由映射关系。
// handlers 是一个 map，key 为 "METHOD-PATH" 格式，value 为对应的处理函数。
type router struct {
	handlers map[string]HandlerFunc
}

// newRouter 创建并返回一个新的 router 实例。
// 初始化时创建了一个空的 handlers 映射。
func newRouter() *router {
	return &router{handlers: make(map[string]HandlerFunc)}
}

// addRoute 用于注册一个新的路由规则。
// 参数 method 是 HTTP 方法，如 "GET" 或 "POST"。
// 参数 pattern 是路由路径。
// 参数 handler 是请求到达该路由时的处理函数。
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	// 构造 key，格式为 "METHOD-PATTERN"，确保不同请求可以唯一匹配到处理函数。
	key := method + "-" + pattern
	r.handlers[key] = handler
}

// handle 根据请求的 Context 查找对应的处理函数并执行。
// 如果找不到匹配的路由，则返回 404 错误。
func (r *router) handle(c *Context) {
	// 构造 key，用于从 handlers 中查找对应的处理函数。
	key := c.Method + "-" + c.Path
	// 如果找到了对应的处理函数，则调用；否则，返回 404 状态和提示信息。
	if handler, ok := r.handlers[key]; ok {
		handler(c)
	} else {
		c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	}
}
