// day4 直接调用匹配到的处理函数处理请求。
// day5 利用中间件链机制，将匹配的处理函数（或 404 函数）追加到 Context.handlers 中，并调用 c.Next() 触发整个链的执行。
// 中间件链的执行顺序是：先执行所有在分组或全局注册的中间件（按它们被添加的顺序），最后执行最终的路由处理函数（或者在没有匹配路由时执行 404 处理函数）。
package gee

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

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

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern

	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}

	n := root.search(searchParts, 0)
	if n != nil {
		parts := parsePattern(n.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
			}
		}
		return n, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

// PERF:
// r.handle 不再直接调用匹配的路由处理函数，
// 而是将最终处理函数（或404函数）追加到 Context 的中间件链中，
// 然后通过 c.Next() 依次执行整个链，从而实现中间件机制的扩展。
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)

	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		// c.handlers：当前请求上下文中的处理函数链，包含所有中间件和最终的路由处理函数，按顺序依次执行。
		// r.handlers：路由器中存储路由与其最终处理函数之间映射的集合，通过路由键（如 "GET-/path"）定位到对应的处理函数。
		c.handlers = append(c.handlers, r.handlers[key]) // NEW:在请求到来时，根据请求的方法和路径匹配相应的处理函数，并将其添加到 c.handlers 中。而不是立即执行。这为中间件链的处理提供了可能性
	} else {
		c.handlers = append(c.handlers, func(c *Context) { // NEW:  将404处理函数作为一个中间件添加到 c.handlers 列表中
			c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// NEW:在 r.handle 方法中，最后调用 c.Next() 是为了触发当前请求的处理流程，将处理函数链条（包括中间件和最终的路由处理函数）按顺序执行。（因为c.index初始化为-1）
	c.Next()
}
