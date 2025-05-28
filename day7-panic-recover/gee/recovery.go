package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
)

// 用于生成包含调用栈信息的调试信息
func trace(message string) string { // message 是传入的错误信息
	// uintptr 是一种无符号整数类型，其大小足以容纳任何指针的位模式。
	var pcs [32]uintptr // 用于存储调用栈的程序计数器（Program Counter）
	// 获取调用栈信息
	n := runtime.Callers(3, pcs[:]) // 跳过当前函数的前3个调用（包括自身和调用它的函数）；n 是实际写入 pcs 数组的调用栈帧数量

	var str strings.Builder                   // 创建一个 strings.Builder对象，用于高校拼接字符串
	str.WriteString(message + "\nTraceback:") // 将错误信息和 "Traceback:" 作为前缀写入
	// 遍历捕获的调用栈地址
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)                           // 根据程序计数器获取对应的函数信息
		file, line := fn.FileLine(pc)                         // 获取调用栈中具体的文件名和代码行号
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line)) // 拼接调用栈信息到 strings.Builder
	}
	return str.String() // 返回完整的调用栈字符串
}

// 数组越界错误发生时，向用户返回 Internal Server Error，并且在日志中打印必要的错误信息，方便进行错误定位
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)                               // 将捕获到的 panic 信息格式化为字符串，便于后续日志记录。
				log.Printf("%s\n\n", trace(message))                            // 将 panic 的堆栈信息打印到日志中。
				c.Fail(http.StatusInternalServerError, "Internal Server Error") // 向客户端返回一个通用的错误消息，避免暴露服务端的详细实现信息，增强安全性
			}
		}()

		c.Next()
	}
}
