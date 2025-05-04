package gee

import (
	"html/template"
	"net/http"
	"strings"
)

// HandlerFunc defines the request handler used by gee. 未来可能会带上context
type HandlerFunc func(*Context)

// Engine implements the interface of http.Handler
type Engine struct {
	// Engine将会作为最顶层的分组，因此Engine也要作为一个RouterGroup,具有RouterGroup的所有能力
	// 这里让Engine嵌套routerGroup的原因是，Go语言的嵌套在其他语言中类似于继承，子类必然是比父类有更多的成员变量和方法。RouterGroup 仅仅是负责分组路由，Engine 除了分组路由外，还有很多其他的功能。RouterGroup 继承 Engine 的 Run()，ServeHTTP 等方法是没有意义的。
	*routerGroup
	groups        []*routerGroup
	router        *router
	htmlTemplates *template.Template // 将所有的模板加载进内存
	funcMap       template.FuncMap   // 所有的自定义模板渲染函数
}

func New() *Engine {
	e := &Engine{
		router: newRouter(),
	}
	e.routerGroup = &routerGroup{
		engine: e,
	}
	e.groups = []*routerGroup{e.routerGroup}
	return e
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	var middlewares []HandlerFunc
	// 接收到一个具体请求时，要判断该请求适用于哪些中间件，在这里我们简单通过 URL 的前缀来判断。
	for _, group := range e.groups {
		if strings.HasPrefix(req.URL.Path, group.prefix) {
			middlewares = append(middlewares, group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.engine = e
	c.handlers = middlewares
	e.router.handle(c)
}

func (e *Engine) Run(addr string) (err error) {
	err = http.ListenAndServe(addr, e)
	return err
}


// Default use Logger() & Recovery middlewares
func Default() *Engine {
	engine := New()
	engine.Use(Logger(), Recovery())
	return engine
}

func (e *Engine) SetFuncMap(funcMap template.FuncMap) {
	e.funcMap = funcMap
}

func (e *Engine) LoadHTMLGlob(pattern string) {
	e.htmlTemplates = template.Must(template.New("").Funcs(e.funcMap).ParseGlob(pattern))
}
