package lungo

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
)

// App is the top-level application instance
type App struct {
	mutex  sync.RWMutex
	pool   sync.Pool
	config *Config
	router *Router
	server *http.Server
}

// New creates an instance of App.
func New(configure ...func(*Config)) (app *App) {
	app = &App{
		router: NewRouter(),
		pool: sync.Pool{
			New: func() interface{} {
				return &Context{App: app}
			},
		},
		config: &Config{
			MaxBodySize: DefaultMaxBodySize,
		},
	}
	for _, c := range configure {
		c(app.config)
	}
	return
}

// NewContext returns a new Context instance.
//
// It serves as an adapter for http.Handlerfunc and converts the
// request to the context based API provided by Lungo.
func (app *App) NewContext(w http.ResponseWriter, r *http.Request) *Context {
	params, _ := url.ParseQuery(r.URL.RawQuery)
	return &Context{
		App:      app,
		Request:  r,
		Response: w,
		Params:   params,
	}
}

// Server returns the http.Server instance of the application
func (app *App) Server() *http.Server {
	return app.server
}

// Config returns the Config instance of the application
func (app *App) Config() *Config {
	return app.config
}

// AcquireContext acquires a empty context instance from the pool.
// This context instance must be released by calling `ReleaseContext()`.
func (app *App) AcquireContext() *Context {
	return app.pool.Get().(*Context)
}

// ReleaseContext releases the context instance back to the pool.
// The context instace must first be acquired by calling `AcquireContext()`.
func (app *App) ReleaseContext(c *Context) {
	app.pool.Put(c)
}

// Get adds a new Route with http method "GET" to the Router of the application.
func (app *App) Get(path string, handler HandlerFunc) {
	app.Handle(http.MethodGet, path, handler)
}

// Head adds a new Route with http method "HEAD" to the Router of the application.
func (app *App) Head(path string, handler HandlerFunc) {
	app.Handle(http.MethodHead, path, handler)
}

// Post adds a new Route with http method "POST" to the Router of the application.
func (app *App) Post(path string, handler HandlerFunc) {
	app.Handle(http.MethodPost, path, handler)
}

// Put adds a new Route with http method "PUT" to the Router of the application.
func (app *App) Put(path string, handler HandlerFunc) {
	app.Handle(http.MethodPut, path, handler)
}

// Patch adds a new Route with http method "PATCH" to the Router of the application.
func (app *App) Patch(path string, handler HandlerFunc) {
	app.Handle(http.MethodPatch, path, handler)
}

// Delete adds a new Route with http method "DELETE" to the Router of the application.
func (app *App) Delete(path string, handler HandlerFunc) {
	app.Handle(http.MethodDelete, path, handler)
}

// Connect adds a new Route with http method "CONNECT" to the Router of the application.
func (app *App) Connect(path string, handler HandlerFunc) {
	app.Handle(http.MethodConnect, path, handler)
}

// Options adds a new Route with http method "OPTIONS" to the Router of the application.
func (app *App) Options(path string, handler HandlerFunc) {
	app.Handle(http.MethodOptions, path, handler)
}

// Trace adds a new Route with http method "TRACE" to the Router of the application.
func (app *App) Trace(path string, handler HandlerFunc) {
	app.Handle(http.MethodTrace, path, handler)
}

// Handle adds a new Route with the specified http method to the Router of the application.
func (app *App) Handle(method, path string, handler HandlerFunc) {
	app.router.Handle(Route{Method: method, Path: path, Handler: handler})
}

// All adds a new Route on all HTTP methods to the Router of the application.
func (app *App) All(path string, handler HandlerFunc) {
	for _, method := range methods {
		app.Handle(method, path, handler)
	}
}

// Static adds a new Route to the Router of the application, which serves static files.
func (app *App) Static(path, root string) {
	app.router.Handle(Route{Method: http.MethodGet, Path: path, Handler: FileHandler(root)})
}

// Use adds a Middleware to the router.
// Middleware can be used to intercept or otherwise modify requests.
// The are executed in the order that they are applied to the Router (FIFO).
func (app *App) Use(middlewares ...Middleware) {
	app.router.Use(middlewares...)
}

// Mount adds a new app which handles requests on the specified pattern.
func (app *App) Mount(pattern string, group *App) {
	app.All(pattern, WithContext(http.StripPrefix(pattern, group)))
}

// HandleError is a centralized error handler function which resolves the
// provided error and replies to the request with the specified error message
// and HTTP code. The error message is written as plain text.
// JSON response with status code and message.
func (app *App) HandleError(c *Context, e error) {
	re, ok := e.(*RequestError)
	if !ok {
		re = &RequestError{
			Code:    http.StatusInternalServerError,
			Message: http.StatusText(http.StatusInternalServerError),
		}
	}

	http.Error(c.Response, re.Message, re.Code)
}

// ServeHTTP implements the http.Handler interface which
// is used by the http.Server to dispatch requests.
//
// It serves as an adapter for http.Handler and converts the
// request to the context based API provided by Lungo.
func (app *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, ok := r.Context().Value("context").(*Context)
	if !ok {
		c = app.AcquireContext()
	}
	c.Reset(w, r)
	if err := app.router.ServeHTTP(c); err != nil {
		app.HandleError(c, err)
	}
	if !ok {
		app.ReleaseContext(c)
	}
}

// Serve accepts incoming HTTP connections on the listener l,
// creating a new service goroutine for each. The service goroutines
// read requests and then call handler to reply to them.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// HTTP/2 support is only enabled if the Listener returns *tls.Conn
// connections and they were configured with "h2" in the TLS
// Config.NextProtos.
//
// Serve always returns a non-nil error.
func (app *App) Serve(l net.Listener) error {
	app.mutex.Lock()
	app.server = &http.Server{Handler: app}
	app.mutex.Unlock()
	return app.server.Serve(l)
}

// ServeTLS accepts incoming HTTPS connections on the listener l,
// creating a new service goroutine for each. The service goroutines
// read requests and then call handler to reply to them.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// Additionally, files containing a certificate and matching private key
// for the server must be provided. If the certificate is signed by a
// certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
//
// ServeTLS always returns a non-nil error.
func (app *App) ServeTLS(l net.Listener, certFile, keyFile string) error {
	app.mutex.Lock()
	app.server = &http.Server{Handler: app}
	app.mutex.Unlock()
	return app.server.ServeTLS(l, certFile, keyFile)
}

// Listen listens on the TCP network address addr and then calls
// Serve with handler to handle requests on incoming connections.
// Accepted connections are configured to enable TCP keep-alives.
//
// The handler is typically nil, in which case the DefaultServeMux is used.
//
// ListenAndServe always returns a non-nil error.
func (app *App) Listen(addr string) error {
	app.mutex.Lock()
	app.server = &http.Server{Addr: addr, Handler: app}
	app.mutex.Unlock()
	return app.server.ListenAndServe()
}

// ListenTLS acts identically to Listen, except that it
// expects HTTPS connections. Additionally, files containing a certificate and
// matching private key for the server must be provided. If the certificate
// is signed by a certificate authority, the certFile should be the concatenation
// of the server's certificate, any intermediates, and the CA's certificate.
func (app *App) ListenTLS(addr, certFile, keyFile string) error {
	app.mutex.Lock()
	app.server = &http.Server{Addr: addr, Handler: app}
	app.mutex.Unlock()
	return app.server.ListenAndServeTLS(certFile, keyFile)
}

// Shutdown gracefully shuts down the server without interrupting any
// active connections. Shutdown works by first closing all open
// listeners, then closing all idle connections, and then waiting
// indefinitely for connections to return to idle and then shut down.
// If the provided context expires before the shutdown is complete,
// Shutdown returns the context's error, otherwise it returns any
// error returned from closing the Server's underlying Listener(s).
//
// When Shutdown is called, Serve, ListenAndServe, and
// ListenAndServeTLS immediately return ErrServerClosed. Make sure the
// program doesn't exit and waits instead for Shutdown to return.
//
// Shutdown does not attempt to close nor wait for hijacked
// connections such as WebSockets. The caller of Shutdown should
// separately notify such long-lived connections of shutdown and wait
// for them to close, if desired. See RegisterOnShutdown for a way to
// register shutdown notification functions.
//
// Once Shutdown has been called on a server, it may not be reused;
// future calls to methods such as Serve will return ErrServerClosed.
func (app *App) Shutdown(ctx context.Context) error {
	app.mutex.RLock()
	defer app.mutex.RUnlock()
	if app.server == nil {
		return fmt.Errorf("Shutdown: Server is not running.")
	}
	return app.server.Shutdown(ctx)
}
