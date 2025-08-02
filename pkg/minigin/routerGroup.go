package minigin

import (
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

// Use ミドルウェアを追加する
func (g *RouterGroup) Use(middlewares ...HandlerFunc) {
	g.middlewares = append(g.middlewares, middlewares...)
}

// Group 新しいルーターグループを作成する
func (g *RouterGroup) Group(prefix string) *RouterGroup {
	newGroup := &RouterGroup{
		prefix: g.prefix + prefix,
		parent: g,
		engine: g.engine,
	}
	return newGroup
}

// handle ルートとハンドラチェーンの構築
func (g *RouterGroup) handle(method, relativePath string, handler HandlerFunc) {
	absPath := path.Join(g.prefix, relativePath)

	// 1. 現在のグループから親へ遡り、ミドルウェアを逆順で収集
	var collectedMiddlewares []HandlerFunc
	group := g
	for group != nil {
		collectedMiddlewares = append(collectedMiddlewares, group.middlewares...)
		group = group.parent
	}

	// 2. 収集した逆順のミドルウェアを反転させ、正しい実行順序（親→子）にする
	for i, j := 0, len(collectedMiddlewares)-1; i < j; i, j = i+1, j-1 {
		collectedMiddlewares[i], collectedMiddlewares[j] = collectedMiddlewares[j], collectedMiddlewares[i]
	}

	// 3. ハンドラチェーンを作成
	finalHandlers := make([]HandlerFunc, 0, len(collectedMiddlewares)+1)
	finalHandlers = append(finalHandlers, collectedMiddlewares...)
	finalHandlers = append(finalHandlers, handler)

	g.engine.addRoute(method, absPath, finalHandlers)
}

// GET メソッドを追加
func (g *RouterGroup) GET(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodGet, relativePath, handler)
}

// POST メソッドを追加
func (g *RouterGroup) POST(relativePath string, handler HandlerFunc) {
	g.handle(http.MethodPost, relativePath, handler)
}
