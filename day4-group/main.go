package main

import (
	"gee"
	"net/http"
)

func main() {
	r := gee.New()                         // 创建Engine 实例
	r.GET("/index", func(c *gee.Context) { // 为 /index 添加一个GET请求的处理函数
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})

	v1 := r.Group("/v1") // 创建第一个路由组 v1，前缀为 /v1 （Engine作为顶层分组）
	{                    // 用大括号包裹的代码是为了更清晰第表示路由组内的路由
		v1.GET("/", func(c *gee.Context) { // 为路径 /v1/添加一个GET请求的处理函数
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})
		v1.GET("/hello", func(c *gee.Context) { // 为路径 /v1/hello ...
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Query("name"), c.Path)
		})
	}

	v2 := r.Group("/v2") // 创建第二个路由组 v2，前缀为 /v2 （Engine作为顶层分组）
	{
		v2.GET("/hello/:name", func(c *gee.Context) { // 为路径 /v2/hello/:name 添加...
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
		v2.POST("/login", func(c *gee.Context) { // 为路径 /v2/login 添加一个POST请求的处理函数
			c.JSON(http.StatusOK, gee.H{ // 处理函数从POST表单中读取username和passwod并返回一个JSON响应
				"username": c.PostForm("username"),
				"password": c.PostForm("password"),
			})
		})
	}

	r.Run(":9999")
}
