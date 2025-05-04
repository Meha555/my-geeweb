package gee

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func Logger() HandlerFunc {
	return func(c *Context) {
		// Start timer
		t := time.Now()
		// Process request
		c.Next()
		// Calculate resolution time
		log.Printf("[%d] %s in %v", c.StatusCode, c.Req.RequestURI, time.Since(t))
	}
}

// Recovery 要放在第一个
func Recovery() HandlerFunc {
	return func(c *Context) {
		defer func() {
			if err := recover(); err != nil {
				message := fmt.Sprintf("%s", err)
				log.Printf("%s\n\n", trace(message))
				c.Error(http.StatusInternalServerError, "Internal Server Error")
			}
			// c.Abort()
		}()
		// 必须要有Next()。因为defer recover机制只能针对于当前函数以及直接调用的函数的panic，所以需要让整条中间件调用链必须是DFS的而不能是BFS的，BFS的部分不是Recovery直接调用的函数，无法recover，panic会被net/http自带的recover机制捕获
		c.Next()
	}
}

// trace acquire stack trace for debug
func trace(message string) string {
	var pcs [32]uintptr
	n := runtime.Callers(4, pcs[:]) // skip first 3 caller

	var builder strings.Builder
	builder.WriteString(message + "\nTraceback:")
	// for _, pc := range pcs[:n] {
	// 	fn := runtime.FuncForPC(pc)
	// 	file, line := fn.FileLine(pc)
	// 	builder.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	// }
	frames := runtime.CallersFrames(pcs[:n])
	for {
		frame, more := frames.Next()
		builder.WriteString(fmt.Sprintf("\n\t%s:%d", frame.File, frame.Line))
		if !more {
			break
		}
	}
	return builder.String()
}

func CORSMiddleware() HandlerFunc {
	return func(c *Context) {
		fmt.Print("1")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Req.Method == "OPTIONS" {
			c.Error(http.StatusNoContent, "cors error")
			return
		}

		c.Next()
	}
}
