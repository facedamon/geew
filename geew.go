package geew

import (
	"net/http"
)

// HandlerFunc defines the request handler used by geew
//type HandlerFunc func(http.ResponseWriter, *http.Request)
type HandlerFunc func(*Context)

// Engine implement the interface of ServeHTTP
type Engine struct {
	//router map[string]HandlerFunc
	router *router
	*group
	// store all groups
	groups []*group
}

// New is the constructor of geew.Engine
func New() *Engine {
	//return &Engine{router: newRouter()}
	e := &Engine{router: newRouter()}
	e.group = &group{engine: e}
	e.groups = []*group{e.group}
	return e
}

func (e *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	e.router.addRoute(method, pattern, handler)
}

// GET defines the method to add GET request
func (e *Engine) GET(pattern string, handler HandlerFunc) {
	e.addRoute("GET", pattern, handler)
}

// POST defines the method to add POST request
func (e *Engine) POST(pattern string, handler HandlerFunc) {
	e.addRoute("POST", pattern, handler)
}

// Run defines the method to start up a http server
func (e *Engine) Run(addr string) error {
	return http.ListenAndServe(addr, e)
}

// Implement the HTTP ServeHTTP
func (e *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := newContext(w, req)
	e.router.handle(c)
}
