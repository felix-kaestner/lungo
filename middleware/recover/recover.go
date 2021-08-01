package recover

import (
	"fmt"

	"github.com/felix-kaestner/lungo"
)

func New(configure ...func(*Config)) lungo.Middleware {
	config := new(Config)
	*config = *DefaultConfig

	for _, c := range configure {
		c(config)
	}

	return func(next lungo.Handler) lungo.Handler {
		return lungo.HandlerFunc(func(c *lungo.Context) (err error) {
			defer func() {
				if r := recover(); r != nil {
					if config.HandleStackTrace != nil {
						config.HandleStackTrace(r)
					}

					var ok bool
					if err, ok = r.(error); !ok {
						// Set error that will call the global error handler
						err = fmt.Errorf("%v", r)
					}
				}
			}()

			return next.ServeHTTP(c)
		})
	}
}
