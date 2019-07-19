// Note: This is a copy of the Mat's router "way" at https://github.com/matryer/way
// I used as a way to learn how the routing can be made without using gorilla/mux
// The implementation is the same, but the comments are mine for my own understanding of the subject

package src

import (
	"context"
	"net/http"
	"strings"
)

type muxContextKey string

type route struct {
	method string
	segments []string
	handler http.Handler
	prefix bool
}

type Router struct {
	routes []*route
	NotFound http.Handler
}

func NewRouter() *Router {
	return &Router{
		NotFound: http.NotFoundHandler(),
	}
}

// If we have the next URL: /blog/21/name/jorge, we gonna get
// [blog, 21, name, jorge]
func (r *Router) pathSegments(pattern string) []string {
	return strings.Split(strings.Trim(pattern, "/"), "/")
}

func (r *Router) Handle(method, pattern string, handler http.Handler) {
	route := &route{
		method: strings.ToLower(method),
		segments: r.pathSegments(pattern),
		handler: handler,
		prefix: strings.HasPrefix(pattern, "/") || strings.HasSuffix(pattern, "..."),
	}
	r.routes = append(r.routes, route)
}

func (r *Router) HandleFunc(method, pattern string, fn http.HandlerFunc)  {
	r.Handle(method, pattern, fn)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := strings.ToLower(req.Method)
	segments := r.pathSegments(req.URL.Path)
	for _, route := range r.routes {
		if route.method != method && route.method != "*" {
			continue
		}

		if ctx, ok := route.match(req.Context(), r, segments); ok {
			route.handler.ServeHTTP(w, req.WithContext(ctx))
			return
		}
	}
	r.NotFound.ServeHTTP(w, req)
}

func Param(ctx context.Context, param string) string {
	value, ok :=  ctx.Value(muxContextKey(param)).(string)
	if !ok {
		return ""
	}
	return value
}

func (r *route) match(ctx context.Context, router *Router, segments []string) (context.Context, bool) {
	if len(segments) > len(r.segments) && !r.prefix {
		return nil, false
	}

	for i, seg := range r.segments {
		if i > len(segments) -1{
			return nil, false
		}
		isParam := false
		if strings.HasPrefix(seg, ":") {
			isParam = true
			seg = strings.TrimPrefix(seg, ":")
		}
		if !isParam {
			if strings.HasSuffix(seg, "...") {
				if strings.HasPrefix(segments[i], seg[:len(seg)-3]) {
					return ctx, true
				}
			}
			if seg != segments[i] {
				return nil, false
			}
		}

		if isParam {
			ctx = context.WithValue(ctx, muxContextKey(seg), segments[i])
		}
	}
	return ctx, true
}