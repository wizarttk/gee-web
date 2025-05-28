package gee

import (
	"log"
	"time"
)

// Logger 是一个返回 HandlerFunc 的函数，
// 它定义了一个日志中间件，用于记录每个 HTTP 请求的处理时间和响应状态。
func Logger() HandlerFunc {
	return func(c *Context) {
		// 记录当前时间，作为请求处理开始的时间点。
		startTime := time.Now()

		// 调用 c.Next()，以执行后续的中间件和最终的处理程序。
		c.Next()

		// 计算从开始到现在所经过的时间，即请求的处理耗时。
		duration := time.Since(startTime)

		// 使用 log.Printf 打印日志，包含响应状态码、请求的 URI 和处理时间。
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, duration)
	}
}

// 中间件处理逻辑：
//
// 在处理程序内部，首先记录当前时间 startTime，作为请求处理的起始时间。
//
// 调用 c.Next()，这会将控制权交给下一个中间件或最终的处理程序，直到所有处理完成后返回。
//
// 在 c.Next() 返回后，计算从 startTime 到当前时间的时间间隔 duration，表示请求的处理耗时。
//
// 最后，使用 log.Printf 输出包含响应状态码（c.StatusCode）、请求的 URI（c.Req.RequestURI）和处理时间（duration）的日志信息。
