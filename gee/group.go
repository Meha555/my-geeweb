package gee

import (
	"log"
	"net/http"
	"path"
)

type routerGroup struct {
	prefix      string        // 为了记录路由Path
	middlewares []HandlerFunc // 为了获得添加中间件的能力
	engine      *Engine       // 为了获得访问router的能力
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

/* ---------------------------------- 模板渲染 ---------------------------------- */

func (g *routerGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := path.Join(g.prefix, relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		file := c.Param("filepath")
		// Check if file exists and/or if we have permission to access it
		if _, err := fs.Open(file); err != nil {
			c.SetStatus(http.StatusNotFound)
			return
		}
		// 我们做的只是接收到请求，把请求的路径地址映射到静态资源所在的真是地址，剩下的就交给静态资源服务器去做就好了。
		fileServer.ServeHTTP(c.Writer, c.Req)
	}
}

// Static serve static files
// relativePath is the relative path for client to indicate static files on route path
// root is the absolute path of static files on file-server
func (g *routerGroup) Static(relativePath string, root string) {
	handler := g.createStaticHandler(relativePath, http.Dir(root))
	urlPattern := path.Join(relativePath, "/*filepath")
	// Register GET handlers
	g.GET(urlPattern, handler)
}
