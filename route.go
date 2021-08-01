package lungo

// Route stores information to match and respond to requests.
type Route struct {
	Method  string  `json:"method"`
	Path    string  `json:"path"`
	Handler Handler `json:"-"`
}

// `ServeHTTP` implements the Handler interface
func (route *Route) ServeHTTP(c *Context) error {
	handler := route.Handler
	return handler.ServeHTTP(c)
}
