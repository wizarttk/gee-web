package gee

import (
	"net/http"
	"path"
)

// 使用相对路径可以将请求路径和真实文件路径解耦，便于映射，例如：
// 动态路由与通配符：
// 路由规则 /assets/*filepath：使用*通配符，可以匹配/assets/开头的任意路径
// 匹配的部分存储在参数 filepath 中，如 /assets/js/geektutu.js 对应 filepath = js/geektutu.js
//
// 请求路径：/assets/js/geektutu.js
// 相对路径：js/geektutu.js
// 真实路径：/usr/web/js/geektutu.js
// 通过这种设计，服务器可以灵活调整文件的存储位置，二无需更改URL设计
func (group *RouterGroup) createStaticHandle(relativePath string, fs http.FileSystem) HandlerFunc {
	// 将路由组的前缀与传入的静态文件相对路径拼接成绝对路径
	absolutePath := path.Join(group.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))

	return func(c *Context) {
		file := c.Param("filepath")
		if _, err := fs.Open(file); err != nil {
			c.Status(http.StatusNotFound)
			return
		}
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}
