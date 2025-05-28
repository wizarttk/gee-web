package gee

import (
	"log"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		startTime := time.Now()

		c.Next()

		duration := time.Since(startTime)
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, duration)
	}
}
