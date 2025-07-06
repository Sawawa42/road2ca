package minigin

import (
	"encoding/json"
	"log"
	"net/http"
	"path"
)

type HandlerFunc func(*Context)

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc
	parent      *RouterGroup
	engine      *Engine
}

type Context struct {
	Writer   http.ResponseWriter
	Request  *http.Request
	handlers []HandlerFunc
	index    int
}

type Engine struct {
	*RouterGroup

	// map[string]1個目: "/user/create"のようなパス
	// map[string]2個目: "POST"のようなHTTPメソッド
	// []HandlerFunc: そのルートに適用されるミドルウェアと最終的なハンドラ関数のスライス
	trees map[string]map[string][]HandlerFunc
}

func (c *Context) Next() {
	c.index++
	if c.index < len(c.handlers) {
		c.handlers[c.index](c)
	}
}

type H map[string]any

func (c *Context) JSON(code int, obj any) {
	json, err := json.Marshal(obj)
	if err != nil {
		log.Printf("Failed to json.Marshal: %v", err)
		c.JSON(http.StatusInternalServerError, H{"error": "Internal server error"})
		return
	}
	c.Writer.Header().Set("Content-Type", "application/json")
	c.Writer.WriteHeader(code)
	c.Writer.Write(json)
}

func New() *Engine {
	engine := &Engine{
		trees: make(map[string]map[string][]HandlerFunc),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

func (e *Engine) addRoute(method, relativePath string, handlers []HandlerFunc) {
	if e.trees[relativePath] == nil {
		e.trees[relativePath] = make(map[string][]HandlerFunc)
	}
	e.trees[relativePath][method] = handlers
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path

	pathHandlers, ok := e.trees[path]
	if !ok {
		http.NotFound(w, r)
		return
	}

	methodHandlers, ok := pathHandlers[r.Method]
	// ここでoptionsを考慮しないとCORS対応する前に404が返されてしまう
	if !ok && r.Method == http.MethodOptions {
		methodHandlers = e.RouterGroup.middlewares
	} else if !ok {
		http.NotFound(w, r)
		return
	}

	c := &Context{
		Writer:   w,
		Request:  r,
		handlers: methodHandlers,
		index:    -1,
	}

	c.Next()
}

func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

func (g *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: g.engine,
	}
	return newGroup
}

func (g *RouterGroup) handle(method, relativePath string, handler HandlerFunc) {
	absPath := path.Join(g.prefix, relativePath)
	var handlers []HandlerFunc
	group := g
	for group != nil {
		handlers = append(group.middlewares, handlers...)
		group = group.parent
	}
	handlers = append(handlers, handler)
	g.engine.addRoute(method, absPath, handlers)
}

func (g *RouterGroup) GET(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodGet, relativePath, handler)
}

func (g *RouterGroup) POST(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodPost, relativePath, handler)
}

func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
