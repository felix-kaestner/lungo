package lungo

import (
	"context"
	"net/http"
)

type (
	// A Handler responds to an HTTP request.
	//
	// ServeHTTP should write reply headers and data to the ResponseWriter
	// and then return. Returning signals that the request is finished; it
	// is not valid to use the ResponseWriter or read from the
	// Request.Body after or concurrently with the completion of the
	// ServeHTTP call.
	//
	// Depending on the HTTP client software, HTTP protocol version, and
	// any intermediaries between the client and the Go server, it may not
	// be possible to read from the Request.Body after writing to the
	// ResponseWriter. Cautious handlers should read the Request.Body
	// first, and then reply.
	//
	// Except for reading the body, handlers should not modify the
	// provided Request.
	//
	// If ServeHTTP panics, the server (the caller of ServeHTTP) assumes
	// that the effect of the panic was isolated to the active request.
	// It recovers the panic, logs a stack trace to the server error log,
	// and either closes the network connection or sends an HTTP/2
	// RST_STREAM, depending on the HTTP protocol. To abort a handler so
	// the client sees an interrupted response but the server doesn't log
	// an error, panic with the value ErrAbortHandler.
	Handler interface {
		ServeHTTP(c *Context) error
	}

	// `Middleware` is a function which receives an Handler and returns another Handler.
	Middleware func(Handler) Handler

	// The HandlerFunc type is an adapter to allow the use of
	// ordinary functions as HTTP handlers. If f is a function
	// with the appropriate signature, HandlerFunc(f) is a
	// Handler that calls f.
	HandlerFunc func(c *Context) error
)

// `ServeHTTP` implements the Handler interface for
// HandlerFunc by simply returning the function call
// with the provided context
func (h HandlerFunc) ServeHTTP(c *Context) error { return h(c) }

// `WithContext` is an adapter to allow the usage of http.Handler
// with the context based API provided by Lungo
func WithContext(handler http.Handler) HandlerFunc {
	return HandlerFunc(func(c *Context) error {
		ctx := context.WithValue(c.Request.Context(), "context", c)
		handler.ServeHTTP(c.Response, c.Request.WithContext(ctx))
		return nil
	})
}
