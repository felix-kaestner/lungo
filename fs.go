package lungo

import (
	"net/http"
	"os"
	"path"
)

type FileServer struct {
	root http.Dir
}

// ServeHTTP implements the Handler interface for the file server.
func (f *FileServer) ServeHTTP(c *Context) error {
	fp := path.Join(string(f.root), c.Request.URL.Path)
	if _, err := os.Stat(fp); err != nil {
		return c.NotFound()
	}

	return c.File(fp)
}

// FileHandler creates a new Handler that serves static files in a directory.
func FileHandler(root string) Handler {
	return &FileServer{root: http.Dir(root)}
}
