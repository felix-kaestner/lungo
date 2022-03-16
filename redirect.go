package lungo

import "net/http"

// Redirect to a fixed URL
type Redirect struct {
	url  string
	code int
}

// `ServeHTTP` implements the Handler interface
func (redirect *Redirect) ServeHTTP(c *Context) error {
	http.Redirect(c.Response, c.Request, redirect.url, redirect.code)
	return nil
}

// RedirectHandler creates a new Handler that redirects requests it receives
// to a given URL using the provided status code.
//
// The provided status code should be in the 3xx range and is usually
// http.StatusMovedPermanently, http.StatusFound or http.StatusSeeOther.
func RedirectHandler(url string, code int) Handler {
	return &Redirect{url, code}
}
