package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

// onlyForV2 返回一个 gee.HandlerFunc 类型的中间件函数
func onlyForV2() gee.HandlerFunc {
	return func(c *gee.Context) {
		// 开始计时
		t := time.Now()
		// 模拟一个服务器错误
		c.Fail(500, "internal Serve Error")
		// 计算处理时间并记录
		log.Printf("[%d] %s in %v for group v2", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

func main() {
	r := gee.New()                    // 创建一个新的gee实例
	r.Use(gee.Logger())               // 全局中间件 gee.Logger
	r.GET("/", func(c *gee.Context) { // 给 'GET - /'添加处理函数（/是根路径）
		c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
	})

	v2 := r.Group("/v2") // 创建一个以 /v2 开头的路由组
	v2.Use(onlyForV2())  // 给 v2 路由组添加 onlyForV2 中间件
	{
		v2.GET("/hello/:name", func(c *gee.Context) { // 给 GET - /v2/hello/:name 添加处理函数
			// 预期路径如：/hello/geektutu
			// 返回字符串 "hello {name}, you're at {path}"
			c.String(http.StatusOK, "hello %s, you're at %s\n", c.Param("name"), c.Path)
		})
	}

	// 启动HTTP服务器，监听在端口9999
	r.Run(":9999")
}
