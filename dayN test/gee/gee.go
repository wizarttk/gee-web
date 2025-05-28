package gee

import (
	"log"
	"net/http"
	"path"
	"strings"
)

type HandlerFunc func(*Context)

type (
	RouterGroup struct {
		prefix      string
		middlewares []HandlerFunc
		parent      *RouterGroup
		engine      *Engine
	}
	Engine struct {
		*RouterGroup
		router *router
		groups []*RouterGroup
	}
)

func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup}
	return engine
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}

	engine.groups = append(engine.groups, newGroup)
	return newGroup
}

func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares = append(group.middlewares, middlewares...)
}

func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := group.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.engine.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.engine.addRoute("POST", pattern, handler)
}

// NEW:
// createStaticHandler 创建一个处理静态文件请求的处理器，并确保在处理请求之前验证文件的存在性和可访问性。
// "relativePath" 客户端访问静态文件时的相对路径（表示相对于路由组前缀(group.prefix)的路径。这个路径将用于构建处理静态文件的绝对路径。）
// "fs"           提供静态文件的文件系统（通常使用 http.Dir,可以将文件夹作为文件系统暴露）
// 返回值：HandlerFunc类型的处理函数
//
/*
* 例子：
*  r := gee.New()
*  adminGroup := r.Group("/api")
*  adminGroup.Static("/assets","./static")
* 可能的请求路径：http://localhost:8080/api/assets/image.jpg
*   1. group.prefix 路由组的前缀：
*   2. relativePath 路由中配置的相对路径：/assets
*   3. absolutePath 绝对路径前缀： /api/assets
*   4. filepath     请求中的匹配的文件路径： /image.jpg
*  然后在配置的根目录("./static")下寻找对应文件
*  最终访问的是 ./static/image.jpg
 */
func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	// 将group.prefix 和 relativePath拼接成绝对路径（absolutePath），这将作为静态文件的绝对路径，用于http.StripPrefix的参数，目的是从请求的URL移除这部分路径
	absolutePath := path.Join(group.prefix, relativePath)
	// http.StripPrefix 创建了一个新的 http.Handler，它负责检查请求路径是否以指定前缀开头，
	// 如果是，则移除该前缀并将剩余部分交给另一个 http.Handler 处理；
	// 否则，返回 404 错误。它执行的是精确的字符串匹配。
	// 创建文件服务器，并去除URL中的路由前缀部分，指向静态文件所在的根目录
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath") // 从上下文c中获取请求中的文件路径参数（filepath）
		// 尝试打开请求的文件，如果文件不存在或者没有访问权限，返回HTTP 404状态码
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req) // 如果文件存在且有访问权限，使用创建的文件服务器 fileServer 处理请求，将响应写回客户端
	}
}

// NEW:
// Static 将磁盘上的某个文件夹 root 映射到路由 relative。
// 例如，将本地的 ./static 映射到 /assets，
// 用户访问 /assets/js/geektutu.js 就相当于访问本地的 ./static/js/geektutu.js。
// "relativePath" 客户端访问静态文件时的相对路径（URL路径前缀）
// "root"         静态文件在服务器上的根目录
func (group *RouterGroup) Static(relativePath string, root string) {
	// 调用 createStaticHandler 方法，创建静态文件处理器
	// http.Dir(root) 将静态文件在服务器上的根目录封装成一个 http.FileSystem，用于文件服务器。
	handler := group.createStaticHandler(relativePath, http.Dir(root))
	// 拼接路由模式，用于匹配客户端请求静态文件的路径
	// 形成 /relativePath/*filepath 的格式
	urlPattern := path.Join(relativePath, "/*filepath")
	// 注册 GET 请求的而路由处理器
	group.GET(urlPattern, handler)
}

func (engine *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, engine)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	for _, group := range engine.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}

	c := newContext(w, req)
	c.handlers = middlewares
	engine.router.handle(c)
}
