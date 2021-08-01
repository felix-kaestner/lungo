package lungo

// RequestError stores information about errors during the handling a request.
// The provided code must be a valid HTTP 1xx-5xx status code.
type RequestError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// `Error` implements the error interface
func (e *RequestError) Error() string {
	return e.Message
}

// NotFoundHandler returns a simple request handler
// that replies to each request with a `Not Found` reply.
func NotFoundHandler() Handler {
	return HandlerFunc(func(c *Context) error {
		return c.NotFound()
	})
}
