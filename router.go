package geew

import (
	"fmt"
	"net/http"
	"strings"
)

type router struct {
	// handlers map[string]HandlerFunc
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

// roots key like : roots['GET'] roots['POST']
// handlers key like : handlers['GET-/p/:lang/doc'], handlers['POST-/p/book']
func newRouter() *router {
	//return &router{handlers: make(map[string]HandlerFunc)}
	return &router{roots: make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	v := strings.Split(pattern, "/")
	parts := make([]string, 0)

	for _, item := range v {
		if item != "" {
			parts = append(parts, item)
			// when the pats has '*' return
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}

func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	//log.Printf("Route %4s - %s", method, pattern)
	//key := method + "-" + pattern
	//r.handlers[key] = handler
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	sp := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]

	if !ok {
		return nil, nil
	}

	n := root.search(sp, 0)

	if n != nil {
		parts := parsePattern(n.pattern)
		for i, p := range parts {
			if p[0] == ':' {
				params[p[1:]] = sp[i]
			}
			if p[0] == '*' && len(p) > 1 {
				params[p[1:]] = strings.Join(sp[i:], "/")
				break
			}
		}
		return n, params
	}
	return nil, nil
}

func (r *router) handle(c *Context) {
	//key := c.Method + "-" + c.Path
	//if handler, ok := r.handlers[key]; ok {
	//	handler(c)
	//} else {
	//	c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
	//}
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		c.Params = params
		key := c.Method + "-" + n.pattern
		//r.handlers[key](c)
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		//c.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		c.handlers = append(c.handlers, func(c *Context) {
			c.JSON(http.StatusNotFound, H{
				"code": http.StatusNotFound,
				"msg": fmt.Sprintf("404 NOT FOUND: %s\n", c.Path),
			})
		})
	}
	c.Next()
}
