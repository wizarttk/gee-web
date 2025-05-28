/*
*  路由组结构示例
*        ⬇️
*  Engine (前缀: "")    -- r := gee.New()
*  ├── /v1              -- v1 :=  r.Group("v1")
*  │   ├── /v1/users    -- v1.GET("/users", callback)
*  │   └── /v1/orders
*  ├── /v2
*  │   ├── /v2/api
*  │   └── /v2/admin
*  └── /static
**/
package gee

import (
	"log"
	"net/http"
)

type HandlerFunc func(*Context)

// NEW:
// RouterGroup 用于对具有相同前缀的路由进行分组控制，并支持中间件和分组潜逃
type (
	RouterGroup struct { // 表示一组具有共同前缀和中间件的路由
		prefix      string        // 分组前缀，例如：/ 或者 /api
		middlewares []HandlerFunc // 分组中间件，是一些处理请求的函数，支持统一处理（如日志、鉴权）
		parent      *RouterGroup  // 父分组，要支持分组嵌套，需要知道当前分组的父亲(parent)是谁
		engine      *Engine       // 所有分组共享同一个 Engine 实例。engine字段指向Engine实例，整个框架的所有资源都是由Engine统一协调的，那么就可以通过Engine间接地访问各种接口
	}

	// PERF:
	// Engine 结构体是 Gee Web 框架的核心，它嵌入了 RouterGroup 结构体，使得 Engine 具备了 RouterGroup 的功能，同时又增加了路由器和路由组的管理功能。通过这种设计，Engine 能够：
	//  1. 作为一个最顶层的路由组来管理路由,拥有 RouterGroup 的的所有能力(匿名嵌入*RouterGroup)
	//  2. 使用路由组的功能来组织和管理路由。
	//  3. 管理所有的路由组和路由规则，实现复杂的路由分发和中间件处理逻辑。
	Engine struct {
		*RouterGroup                // （结构体匿名字段）Engine 结构体中嵌入RouterGourp，这样做的目的是使 Engine 本身也能作为一个路由组使用，具备 RouterGroup 的所有功能，如添加前缀、使用中间件等（顶级分组前缀为空）。
		router       *router        // router 字段指向一个 router 实例，负责管理所有的路由规则和请求的分发。
		groups       []*RouterGroup // 保存所有分组，便于全局管理
	}
)

// PERF:
// New 初始化 Engine 时，同时创建顶级分组，并将 Engine 自身作为顶级分组（路由系统中处于最上层的那个分组），
// 从而使得所有基于分组的路由注册和中间件添加都通过这个顶级分组来实现。
func New() *Engine {
	engine := &Engine{router: newRouter()}             // 创建一个新的Engine实例，并初始化路由器（router）
	engine.RouterGroup = &RouterGroup{engine: engine}  // 将 Engine 自身作为顶级路由分组，顶级分组的前缀为空，所有路由都从此开始
	engine.groups = []*RouterGroup{engine.RouterGroup} // 将顶级分组加入到 groups 列表中，便于后续管理
	return engine
}

// NEW:
// Group方法 用于创建一个子路由组
// 新分组的前缀为父分组前缀加上传入的prefix
// 这样可以实现分组嵌套，并且新分组共享同一个 Engine 实例。
func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine    // 获取当前的而 Engine 实例
	newGroup := &RouterGroup{ // 创建新的分组
		prefix: group.prefix + prefix, // 新的前缀为父分组的前缀 + 自己的prefix
		parent: group,                 // 制定父分组
		engine: engine,                // 使用同一个 Engine 实例，所有路由组可以访问和操作同一个路由表
	}

	engine.groups = append(engine.groups, newGroup) // 添加新分组到 groups 列表（groups列表用于存储所有定义的路由组，便于管理和遍历）
	return newGroup                                 // 返回新创建的 newGroup 实例，使得调用者能够进一步对其进行配置和使用
}

// PERF:
// 将 addRoute 方法从 Engine 移到 RouterGroup 后，可以更好地支持路由分组和层次化管理，提升代码的模块化和灵活性。
// 同时，通过路径前缀的自动拼接，简化了路由定义过程，增强了代码的可维护性和可读性。
// comp：路由的路径，例如"/hello"或"/hello/"
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	// 将 RouterGroup 的前缀(group.prefix)与传入的路由路径(comp)进行拼接，得到完整的路由模式(pattern)
	// 例如：如果prefix是"/v1"而comp是"/users"，那么最终的pattern是"/v1/users"
	// 这种设计允许路由组（RouterGroup）拥有一个公共的路由前缀，使得路由定义更加简洁和模块化
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

// PERF:
//
//	func (engine *Engine) GET(pattern string, handler HandlerFunc) {
//		engine.addRoute("GET", pattern, handler)
//	}
func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

// PERF:
//
//	func (engine *Engine) POST(pattern string, handler HandlerFunc) {
//		engine.addRoute("POST", pattern, handler)
//	}
func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}

func (engine *Engine) Run(address string) (err error) {
	return http.ListenAndServe(address, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	engine.router.handle(c)
}
