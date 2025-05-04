package gee

import (
	"net/http"
)

// HandlerFunc defines the request handler used by gee. 未来可能会带上context
type HandlerFunc func(*Context)

// Engine implements the interface of http.Handler
type Engine struct {
	router *router
}

func New() *Engine {
	return &Engine{
		router: newRouter(),
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}

func (e *Engine) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, e)
	return err
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute("POST", pattern, handler)
}
