package cors

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/felix-kaestner/lungo"
)

// New creates a new CORS middleware instance
func New(configure ...func(*Config)) lungo.Middleware {
	config := new(Config)
	*config = *DefaultConfig

	for _, c := range configure {
		c(config)
	}

	allowHeaders := strings.Join(config.AllowHeaders, ",")
	allowMethods := strings.Join(config.AllowMethods, ",")
	exposeHeaders := strings.Join(config.ExposeHeaders, ",")
	maxAge := strconv.Itoa(config.MaxAge)

	return func(next lungo.Handler) lungo.Handler {
		return lungo.HandlerFunc(func(c *lungo.Context) error {
			origin := c.Header(lungo.HeaderOrigin)

			if origin == "" {
				if c.Method() != http.MethodOptions {
					return next.ServeHTTP(c)
				}
				return c.NoContent()
			}

			allowOrigin := ""
			for _, o := range config.AllowOrigins {
				if o == "*" && config.AllowCredentials {
					allowOrigin = origin
					break
				}
				if o == "*" || o == origin {
					allowOrigin = o
					break
				}
			}

			if allowOrigin == "" {
				if c.Method() != http.MethodOptions {
					return next.ServeHTTP(c)
				}
				return c.NoContent()
			}

			// Set Vary: Origin
			c.SetHeader(lungo.HeaderVary, lungo.HeaderOrigin)

			// Set Allow-Origin
			c.SetHeader(lungo.HeaderAccessControlAllowOrigin, allowOrigin)

			// Set Allow-Credentials if set to true
			if config.AllowCredentials {
				c.SetHeader(lungo.HeaderAccessControlAllowCredentials, "true")
			}

			// If not preflight, dispatch
			if c.Method() != http.MethodOptions {
				if exposeHeaders != "" {
					c.SetHeader(lungo.HeaderAccessControlExposeHeaders, exposeHeaders)
				}

				return next.ServeHTTP(c)
			}

			// Preflight request
			c.AddHeader(lungo.HeaderVary, lungo.HeaderAccessControlRequestMethod)
			c.AddHeader(lungo.HeaderVary, lungo.HeaderAccessControlRequestHeaders)

			// Set Allow-Methods
			c.SetHeader(lungo.HeaderAccessControlAllowMethods, allowMethods)

			// Set Allow-Headers if not empty
			if allowHeaders != "" {
				c.SetHeader(lungo.HeaderAccessControlAllowHeaders, allowHeaders)
			}

			// Set MaxAge is set
			if config.MaxAge > 0 {
				c.SetHeader(lungo.HeaderAccessControlMaxAge, maxAge)
			}

			return c.NoContent()
		})
	}
}
