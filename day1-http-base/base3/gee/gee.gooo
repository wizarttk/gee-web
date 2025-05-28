package gee

import (
	"fmt"
	"log"
	"net/http"
)

// HandlerFunc 定义了请求处理函数的类型，接收 ResponseWriter 和 Request 作为参数
// 这样可以方便用户传入自己的处理逻辑
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine 是 Gee 框架的核心结构体，它实现了 http.Handler 接口
type Engine struct {
	// router 用于存储请求路径和对应的处理函数，key 的格式为"请求方法-路由"（例："Get-/hello"）
	router map[string]HandlerFunc
}

// New 创建并返回一个新的 Engine 实例,初始化其 router 字段。
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

// addRoute 注册路由的通用方法 ，参数 method 为HTTP请求方法，pattern 为请求路径，handler 为对应的处理函数
func (engine *Engine) addRoute(method, pattern string, handler HandlerFunc) {
	// 生成路由映射的 key，例如"GET-/hello"
	key := method + "-" + pattern
	// 打印日志，方便调试和查看路由注册情况
	log.Printf("Route %4s - %s", method, pattern)
	// 将 key 和 handler 存入 router 中
	engine.router[key] = handler
}

// GET、POST包装了addRoute，用来注册路由
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

// RUN 启动http serve服务器，并监听指定地址
func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

// ServerHTTP方法 实现了http.Handler接口，当有 HTTP 请求到来时，会调用此方法
// 它根据请求的 Method 和路径 URL.Path 查找对应的处理函数,如果找到则执行,否则返回 404 错误。
func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		fmt.Fprintf(w, "404 NOT FOUND:%s\n", req.URL)
	}
}
