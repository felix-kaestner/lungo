package template

import (
	"github.com/felix-kaestner/lungo"
)

// New creates a new template middleware instance
func New(configure ...func(*Config)) lungo.Middleware {
	config := new(Config)
	*config = *DefaultConfig

	for _, c := range configure {
		c(config)
	}

	return func(next lungo.Handler) lungo.Handler {
		return lungo.HandlerFunc(func(c *lungo.Context) error {
			return next.ServeHTTP(c)
		})
	}
}
