package geew

import "log"

// group defines the router group
type group struct {
	prefix string
	// support middleware
	middlewares []HandlerFunc
	// support nesting
	parent *group
	// all groups share a Engine instance
	engine *Engine
}

// Group is defines to create a new group
// rember all groups share the same Engine instance
func (g *group) Group(prefix string) *group {
	e := g.engine
	newGroup := &group{
		prefix: g.prefix + prefix,
		parent: g,
		engine: e,
	}
	e.groups = append(e.groups, newGroup)
	return newGroup
}

func (g *group) addRoute(method string, comp string, handler HandlerFunc) {
	pattern := g.prefix + comp
	log.Printf("route %4s - %s", method, pattern)
	g.engine.router.addRoute(method, pattern, handler)
}

// Use is defines to add middleware to the group
func (g *group) Use(ms ...HandlerFunc) {
	g.middlewares = append(g.middlewares, ms...)
}

func (g *group) GET(pattern string, handler HandlerFunc) {
	g.addRoute("GET", pattern, handler)
}

func (g *group) POST(pattern string, handler HandlerFunc) {
	g.addRoute("POST", pattern, handler)
}
