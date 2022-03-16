package lungo

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
)

// Router is an HTTP request multiplexer.
type Router struct {
	mux         *http.ServeMux
	mutex       sync.RWMutex
	routes      map[string]map[string]Route
	middlewares []Middleware
}

const (
	// ErrInvalidRouteMethod is returned when the route method is invalid.
	ErrInvalidRouteMethod = "Invalid method. The method of a handler must be a valid http method."
	// ErrEmptyRoutePath is returned when the route path is empty.
	ErrEmptyRoutePath = "Invalid pattern. The path of a handler can't be empty."
	// ErrNilRouteHandler is returned when the route handler is nil.
	ErrNilRouteHandler = "Nil handler. The handler can't be nil."
	// ErrDuplicateHandler is returned when a route handler is already registered for a given path.
	ErrDuplicateHandler = "Duplicate path. Path `%s` already contains a http handler."
)

// NewRouter allocates and returns a new router instance.
func NewRouter() *Router {
	return &Router{routes: make(map[string]map[string]Route), middlewares: make([]Middleware, 0)}
}

// Handle registers a Handler for the given path.
// This will panic, if the path already has a Handler registered.
func (router *Router) Handle(route Route) {
	router.mutex.Lock()
	defer router.mutex.Unlock()

	if !IsValidMethod(route.Method) {
		panic(ErrInvalidRouteMethod)
	}
	if route.Path == "" {
		panic(ErrEmptyRoutePath)
	}
	if route.Handler == nil {
		panic(ErrNilRouteHandler)
	}

	if router.routes[route.Method] == nil {
		router.routes[route.Method] = make(map[string]Route)
	} else if _, exist := router.routes[route.Method][route.Path]; exist {
		panic(fmt.Sprintf(ErrDuplicateHandler, route.Path))
	}

	router.routes[route.Method][route.Path] = route
}

// Use adds a Middleware to the router.
// Middleware can be used to intercept or otherwise modify requests.
// The are executed in the order that they are applied to the Router (FIFO).
func (router *Router) Use(middlewares ...Middleware) {
	router.middlewares = append(router.middlewares, middlewares...)
}

// Find the Route instace for a given path and http method.
// The Route contains the Handler for this request
// This function returns nil if no route matches the request path.
func (router *Router) match(method, path string) *Route {
	router.mutex.RLock()
	defer router.mutex.RUnlock()

	// Check for exact match.
	if route, ok := router.routes[method][path]; ok {
		return &route
	}

	// Check for valid match by comparing the prefix of the route
	// TODO: alternative method for finding closest route via trie tree
	for p, route := range router.routes[method] {
		if strings.HasPrefix(path, p) {
			return &route
		}
	}

	return nil
}

// shouldRedirect determines if the request should be redirected to another path.
//
// This is the case if a handler for path + "/" was registered, but not for the
// path itself, or if a handler was registered for path without any trailing slash.
//
// If the path needs to be changed,it creates a new URL and will return true to
// indicate that a redirect should be made.
func (router *Router) shouldRedirect(method, path string, u *url.URL) (*url.URL, bool) {
	router.mutex.RLock()
	defer router.mutex.RUnlock()

	// In case the path is empty, do not redirect
	n := len(path)
	if n == 0 {
		return u, false
	}

	// If a route exist for the given path, do not redirect
	if _, exist := router.routes[method][path]; exist {
		return u, false
	}

	if path[n-1] != '/' {
		// Redirect, if a handler exists for path + "/"
		if _, exist := router.routes[method][path+"/"]; exist {
			return &url.URL{Path: path + "/", RawQuery: u.RawQuery}, true
		}
	} else {
		// Redirect, if a handler exists for path without trailing slash
		if _, exist := router.routes[method][path[:n-1]]; exist {
			return &url.URL{Path: path[:n-1], RawQuery: u.RawQuery}, true
		}
	}

	return u, false
}

// Handler returns the Handler to use for the given http request.
//
// If the path of the request is not in canonical form, then the
// returned handler will be a redirect to the canonical path.
//
// If no handler matches the given request, then the return
// will be a `Error 404: page not found` Handler.
func (router *Router) Handler(r *http.Request) Handler {
	if r.URL.Path == "" {
		r.URL.Path = "/"
	}

	method := r.Method
	path := r.URL.Path

	// Use canonicalized path except for CONNECT requests.
	if r.Method == http.MethodConnect {
		if u, redirect := router.shouldRedirect(method, path, r.URL); redirect {
			return RedirectHandler(u.String(), http.StatusMovedPermanently)
		}
	} else {
		path = Canonical(r.URL.Path)

		if u, redirect := router.shouldRedirect(method, path, r.URL); redirect {
			return RedirectHandler(u.String(), http.StatusMovedPermanently)
		}

		if path != r.URL.Path {
			u := *r.URL
			u.Path = path
			return RedirectHandler(u.String(), http.StatusMovedPermanently)
		}
	}

	route := router.match(method, path)
	if route == nil {
		return NotFoundHandler()
	}

	return route
}

// `ServeHTTP` implements the Handler interface. It handles the http request
// and dispatches it to the request handler whose path matches the request URL.
func (router *Router) ServeHTTP(c *Context) error {
	if c.Request.RequestURI == "*" {
		// Close all connections using HTTP/1.1 or HTTP/2
		if c.Request.ProtoAtLeast(1, 1) {
			c.Response.Header().Set("Connection", "close")
		}
		return &RequestError{
			Code:    http.StatusBadRequest,
			Message: http.StatusText(http.StatusBadRequest),
		}
	}

	// Set the initial handler function to serve the request
	handler := router.Handler(c.Request)

	// apply middlewares
	for _, middleware := range router.middlewares {
		// middleware is a function in the form (Handler) -> Handler
		// thus returning an instance of the Handler interface.
		handler = middleware(handler)
	}

	// write the response calling ServeHTTP on the Handler interface
	return handler.ServeHTTP(c)
}
