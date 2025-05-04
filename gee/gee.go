package gee

import (
	"fmt"
	"net/http"
)

// HandlerFunc defines the request handler used by gee. 未来可能会带上context
type HandlerFunc func(http.ResponseWriter, *http.Request)

// Engine implements the interface of http.Handler
type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{
		router: make(map[string]HandlerFunc),
	}
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := e.router[key]; ok {
		go handler(w, req)
	} else {
		http.Error(w, fmt.Sprintf("404 NOT FOUND: %s", req.URL.Path), http.StatusNotFound)
	}
}

func (e *Engine) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, e)
	return err
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern // key中遇到的第一个"-"前面是方法，后面是路径
	e.router[key] = handler
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}
