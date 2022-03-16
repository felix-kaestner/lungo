package logging

import (
	"bytes"
	"text/template"
	"time"

	"github.com/felix-kaestner/lungo"
)

// New creates a new logging middleware instance
func New(configure ...func(*Config)) lungo.Middleware {
	config := new(Config)
	*config = *DefaultConfig

	for _, c := range configure {
		c(config)
	}

	return func(next lungo.Handler) lungo.Handler {
		return lungo.HandlerFunc(func(c *lungo.Context) (err error) {
			start := time.Now()
			defer func() {
				if config.Logger == nil {
					return
				}

				var t *template.Template
				if t, err = template.New("log").Parse(config.Template); err != nil {
					return
				}

				data := lungo.Map{
					"Request":  c.Request,
					"Duration": time.Since(start),
				}

				var b bytes.Buffer
				if err = t.Execute(&b, data); err != nil {
					return
				}

				err = config.Logger.Output(config.CallDepth, b.String())
			}()

			return next.ServeHTTP(c)
		})
	}
}
