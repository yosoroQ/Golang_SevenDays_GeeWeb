package gee

import (
	"net/http"
)

type HandlerFunc func(http.ResponseWriter, *http.Request)

type Engine struct {
	router map[string]HandlerFunc
}

func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}

func (enginge *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "-" + pattern
	enginge.router[key] = handler
}

// Get
func (enginge *Engine) GET(pattern string, handler HandlerFunc) {
	enginge.addRoute("GET", pattern, handler)
}

// Post
func (enginge *Engine) Post(pattern string, handler HandlerFunc) {
	enginge.addRoute("POST", pattern, handler)
}

// Run
func (enginge *Engine) Run(addr string) (err error) {
	return http.ListenAndServe(addr, enginge)
}

func (engine *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key := req.Method + "-" + req.URL.Path
	if handler, ok := engine.router[key]; ok {
		handler(w, req)
	} else {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("404 NOT FOUND"))
	}
}
