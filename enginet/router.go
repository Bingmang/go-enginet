package enginet

import (
	"net/http"
	"strings"
)

type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}

func parsePattern(pattern string) []string {
	patternList := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, pattern := range patternList {
		if pattern == "" {
			continue
		}
		parts = append(parts, pattern)
		if pattern[0] == '*' {
			break
		}
	}
	return parts
}

func (r *router) addRoute(method, pattern string, handler HandlerFunc) {
	parts := parsePattern(pattern)
	key := method + "-" + pattern
	if _, ok := r.roots[method]; !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}

func (r *router) getRoute(method, path string) (*node, map[string]string) {
	// 把Path拆分为parts（每段路径为一个part）
	searchParts := parsePattern(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok {
		return nil, nil
	}
	// 从Trie树中搜索target route
	n := root.search(searchParts, 0)
	if n == nil {
		return nil, nil
	}
	// 把target route的pattern分解成parts，并搜索需要解析成参数的part
	parts := parsePattern(n.pattern)
	for i, part := range parts {
		// 如果是":"，则把当前part作为参数
		if part[0] == ':' {
			params[part[1:]] = searchParts[i]
		}
		// 如果是通配符，则停止搜索，把后面的全部当作参数
		if part[0] == '*' && len(part) > 1 {
			params[part[1:]] = strings.Join(searchParts[i:], "/")
			break
		}
	}
	return n, params
}

func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n != nil {
		key := c.Method + "-" + n.pattern
		c.Params = params
		c.handlers = append(c.handlers, r.handlers[key])
	} else {
		c.handlers = append(c.handlers, func(ctx *Context) {
			ctx.String(http.StatusNotFound, "404 NOT FOUND: %s\n", c.Path)
		})
	}
	c.Next()
}
