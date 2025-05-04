package gee

import (
	"net/http"
	"strings"
)

type router struct {
	// 存储每种请求Method的Trie 树根节点 eg, roots['GET'] roots['POST']
	roots map[string]*node
	// 存储每个请求Path的 HandlerFunc eg, handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

// 整个Path中，只允许出现一次"*"
func parsePattern(pattern string) []string {
	items := strings.Split(pattern, "/")

	parts := make([]string, 0)
	for _, item := range items {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)

	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0) // 切片的下标从0开始，所以这里树高也规定为从0开始
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	node := root.search(searchParts, 0)
	// 如果存在这个节点
	if node != nil {
		parts := parsePattern(node.pattern)
		for index, part := range parts {
			if part[0] == ':' {
				params[part[1:]] = searchParts[index]
			}
			if part[0] == '*' && len(part) > 1 {
				params[part[1:]] = strings.Join(searchParts[index:], "/")
				break
			}
		}
		return node, params
	}

	return nil, nil
}

func (r *router) getRoutes(method string) []*node {
	root, ok := r.roots[method]
	if !ok {
		return nil
	}
	nodes := make([]*node, 0)
	root.travel(&nodes)
	return nodes
}

func (r *router) handle(c *Context) {
	node, params := r.getRoute(c.Method, c.Path)
	if node != nil {
		c.Params = params
		key := c.Method + "-" + node.pattern // 注意这里需要使用解析后的动态参数路由，所以不能用c.Path
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(c *Context) {
			c.Error(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	// 开始执行洋葱模型
	c.Next()
}
