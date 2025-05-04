package gee

import (
	"log"
	"path"
)

type routerGroup struct {
	prefix      string
	middlewares []HandlerFunc
	engine      *Engine // 为了获得访问router的能力
}

// Group is defined to create a new RouterGroup
func (g *routerGroup) Group(prefix string) *routerGroup {
	newGroup := &routerGroup{
		prefix: path.Join(g.prefix, prefix),
		engine: g.engine, // all groups share the same Engine instance
	}
	g.engine.groups = append(g.engine.groups, newGroup)
	return newGroup
}

func (g *routerGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("Route %4s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (g *routerGroup) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (g *routerGroup) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}

// Use is defined to add middleware to the group
func (g *routerGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}
