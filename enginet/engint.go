package enginet

import (
	"log"
	"net/http"
)

type HandlerFunc func(ctx *Context)

type Engine struct {
	router *router
}

func (e *Engine) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c := newContext(w, r)
	e.router.handle(c)
}

func New() *Engine {
	return &Engine{router: newRouter()}
}

func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.router.addRoute("GET", pattern, handler)
}

func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.router.addRoute("POST", pattern, handler)
}

func (e *Engine) Run(addr string) error {
	log.Printf("[EngineT] Running at http://%s\n", addr)
	return http.ListenAndServe(addr, e)
}
