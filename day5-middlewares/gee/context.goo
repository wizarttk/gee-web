// day4的 context.go 仅封装了基本的请求和相应操作，
// day5在此基础上增加了中间件链的支持
package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer     http.ResponseWriter
	Req        *http.Request
	Params     map[string]string
	Method     string
	Path       string
	StatusCode int
	// NEW:
	// 请求处理时通过调用 c.Next() 按顺序执行。
	handlers []HandlerFunc // 存储中间件链和最终的处理函数
	index    int           // 当前执行到的中间件索引
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Method: req.Method,
		Path:   req.URL.Path,
		index:  -1, // NEW: 初始情况下index被设置为-1，因此第一次调用Next方法时，index被设置为0，指向第一个中间件处理函数
	}
}

// NEW:
// Next方法 依次执行中间件链的下一个处理函数
func (c *Context) Next() {
	c.index++                      // 这行代码将当前的中间件的索引 index 加1。这样可以指向下一个中间件处理函数
	s := len(c.handlers)           // 将中间件处理函数链的长度存储在变量s中。这个值表示当前请求需要执行的中间件函数的总数
	for ; c.index < s; c.index++ { // 逐个执行中间件直到链结束
		c.handlers[c.index](c)
	}
}

// NEW:
// Fail方法 用于终止中间件的执行并立即返回错误响应。
// 可以在遇到错误时立即停止处理后续的中间件，并返回指定的错误信息和状态码
func (c *Context) Fail(code int, err string) {
	// 将 index 设置为 handlers 的长度，这意味着不能继续循环遍历中间件了
	// 因为 Next 方法在执行中间件时是通过for循环条件c.idnex < len(c.handlers)来决定是否继续执行下一个中间件的
	// 将 index 设置为 len(c.handlers) 可以确保 Next 方法不会执行任何后续的中间件
	c.index = len(c.handlers)
	// 发送一个带有错误信息的JSON响应，便于客户端解析和显示错误
	// H{"message": err} 是一个键值对,键为"message"，用来表示错误的具体内容
	c.JSON(code, H{"message": err})
}

func (c *Context) Param(key string) string {
	value, _ := c.Params[key]
	return value
}

func (c *Context) PostForm(key string) string {
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

func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
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
	c.Status(code)
	c.Writer.Write([]byte(html))
}
