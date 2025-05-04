package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type H map[string]interface{}

type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request
	// request info
	Method string
	Path   string
	// response info
	StatusCode int
}

func newContext(w http.ResponseWriter, r *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    r,
		Method: r.Method,
		Path:   r.URL.Path,
	}
}

/* ---------------------------------- request封装：解析参数 ---------------------------------- */

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
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
	c.SetStatus(code)
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
