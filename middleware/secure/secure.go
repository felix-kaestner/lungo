package secure

import (
	"github.com/felix-kaestner/lungo"
)

func New(configure ...func(*Config)) lungo.Middleware {
	config := new(Config)
	*config = *DefaultConfig

	for _, c := range configure {
		c(config)
	}

	return func(next lungo.Handler) lungo.Handler {
		return lungo.HandlerFunc(func(c *lungo.Context) error {
			if config.XSSProtection != "" {
				c.AddHeader(lungo.HeaderXXSSProtection, config.XSSProtection)
			}
			if config.XFrameOptions != "" {
				c.AddHeader(lungo.HeaderXFrameOptions, config.XFrameOptions)
			}
			if config.ContentSecurityPolicy != "" {
				c.AddHeader(lungo.HeaderContentSecurityPolicy, config.ContentSecurityPolicy)
			}
			if config.ContentTypeNosniff != "" {
				c.AddHeader(lungo.HeaderXContentTypeOptions, config.ContentTypeNosniff)
			}
			if config.ReferrerPolicy != "" {
				c.AddHeader(lungo.HeaderReferrerPolicy, config.ReferrerPolicy)
			}
			return next.ServeHTTP(c)
		})
	}
}
