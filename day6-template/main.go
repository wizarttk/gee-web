package main

/*
(1) render array
$ curl http://localhost:9999/date
<html>
<body>
    <p>hello, gee</p>
    <p>Date: 2019-08-17</p>
</body>
</html>
*/

/*
(2) custom render function
$ curl http://localhost:9999/students
<html>
<body>
    <p>hello, gee</p>
    <p>0: Geektutu is 20 years old</p>
    <p>1: Jack is 22 years old</p>
</body>
</html>
*/

/*
(3) serve static files
$ curl http://localhost:9999/assets/css/geektutu.css
p {
    color: orange;
    font-weight: 700;
    font-size: 20px;
}
*/

import (
	"fmt"
	"gee"
	"html/template"
	"net/http"
	"time"
)

// 定义student结构体,包含Name和Age两个字段
type student struct {
	Name string
	Age  int8 // 学生年龄
}

// 自定义模板函数,用于格式化时间
func FormatAsDate(t time.Time) string {
	year, month, day := t.Date()                         // 从时间中提取年月日
	return fmt.Sprintf("%d-%02d-%02d", year, month, day) // 格式化为yyyy-mm-dd格式
}

func main() {
	r := gee.New() // 创建gee实例

	r.Use(gee.Logger()) // 注册Logger中间件

	// 设置模板自定义函数映射
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate": FormatAsDate,
	})

	r.LoadHTMLGlob("templates/*")   // 加载模板文件
	r.Static("/assets", "./static") // 设置静态文件路由

	// 创建两个student实例
	stu1 := &student{Name: "Geektutu", Age: 20}
	stu2 := &student{Name: "Jack", Age: 22}

	// 根路由处理函数
	r.GET("/", func(c *gee.Context) {
		c.HTML(http.StatusOK, "css.tmpl", nil)
	})

	// /students路由处理函数
	r.GET("/students", func(c *gee.Context) {
		c.HTML(http.StatusOK, "arr.tmpl", gee.H{
			"title":  "gee",
			"stuArr": [2]*student{stu1, stu2}, // 传入学生数组
		})
	})

	// /date路由处理函数
	r.GET("/date", func(c *gee.Context) {
		c.HTML(http.StatusOK, "custom_func.tmpl", gee.H{
			"title": "gee",
			"now":   time.Date(2019, 8, 17, 0, 0, 0, 0, time.UTC), // 创建特定时间
		})
	})

	// 启动服务器,监听9999端口
	r.Run(":9999")
}
