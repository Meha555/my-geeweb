package gee

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Method string
	Path   string
	Params map[string]string
	// response info
	StatusCode int
	// middleware info
	handlers []HandlerFunc // 整个路由上绑定的中间件，最后一个应当是用户的业务Handler
	index    int           // 当前执行到第几个中间件
	// abort    bool          // 是否中断洋葱模型
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		Path:   r.URL.Path,
		index:  -1,
	}
}

/* ---------------------------------- request封装：解析参数 ---------------------------------- */

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

/* ---------------------------------- response封装：设置参数 ---------------------------------- */

func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}

func (c *Context) SetStatus(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}

/* ---------------------------------- response封装：返回响应 ---------------------------------- */

func (c *Context) Text(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.SetStatus(code)
	c.Writer.Write(fmt.Appendf(nil, format, values...))
}

func (c *Context) Data(code int, data []byte) {
	c.SetHeader("Content-Type", "application/octet-stream")
	c.SetStatus(code)
	c.Writer.Write(data)
}

func (c *Context) Error(code int, format string, values ...interface{}) {
	c.StatusCode = code
	http.Error(c.Writer, fmt.Sprintf(format, values...), code)
}

func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.SetStatus(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
	}
}

func (c *Context) HTML(code int, html string) {
	c.SetHeader("Content-Type", "text/html")
	c.SetStatus(code)
	c.Writer.Write([]byte(html))
}

/* ----------------------------------- 中间件 ---------------------------------- */

// Next 实现了DFS
func (c *Context) Next() {
	c.index++
	// if !c.abort {
	// 	c.handlers[c.index](c)
	// }
	l := len(c.handlers)
	// 这里实现了BFS，主要是避免用户实现的中间件忘记调用Next()
	for ; c.index < l; /*&& !c.abort*/ c.index++ {
		c.handlers[c.index](c)
	}
}

// Abort 终止后续中间件的执行。Abort+return终止此后所有代码的执行
func (c *Context) Abort() {
	log.Printf("Abort handling at handlers[%d]", c.index)
	// c.abort = true
	c.index = len(c.handlers)
}
