package minigin

import (
	"net/http"
)

type Engine struct {
	*RouterGroup

	// map[string]1個目: "/user/create"のようなパス
	// map[string]2個目: "POST"のようなHTTPメソッド
	// []HandlerFunc: そのルートに適用されるミドルウェアと最終的なハンドラ関数のスライス
	trees map[string]map[string][]HandlerFunc
}

// New はEngineインスタンスの作成
func New() *Engine {
	engine := &Engine{
		trees: make(map[string]map[string][]HandlerFunc),
	}
	engine.RouterGroup = &RouterGroup{engine: engine}
	return engine
}

// addRoute はエンジンにルートを追加する
func (e *Engine) addRoute(method, relativePath string, handlers []HandlerFunc) {
	if e.trees[relativePath] == nil {
		e.trees[relativePath] = make(map[string][]HandlerFunc)
	}
	e.trees[relativePath][method] = handlers
}

// ServeHTTP http.Handlerインターフェースの実装
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

	rw := &ResponseWriter{
		Writer: w,
		status: 0,
		size:   0,
	}

	c := &Context{
		Writer:   rw,
		Request:  r,
		handlers: methodHandlers,
		index:    -1,
	}

	c.Next()
}

// Run HTTPサーバを指定のアドレスで起動する
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}
