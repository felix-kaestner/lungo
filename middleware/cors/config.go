package cors

import "net/http"

type Config struct {
	// AllowOrigins defines a list of origins that may access the resource.
	//
	// Optional. Default value []string{"*"}.
	AllowOrigins []string

	// AllowMethods defines a list of http methods allowed when accessing the resource.
	// This is used in response to a preflight request.
	//
	// Optional. Default value "GET,HEAD,PUT,POST,PATCH,DELETE"
	AllowMethods []string

	// AllowHeaders defines a list of request headers that can be used when
	// making the actual request. This is in response to a preflight request.
	//
	// Optional. Default value []string{}.
	AllowHeaders []string

	// AllowCredentials defines whether or not the response to the request
	// can be exposed when the credentials flag is true. When used as part of
	// a response to a preflight request, this indicates whether or not the
	// actual request can be made using credentials.
	//
	// Optional. Default value false.
	AllowCredentials bool

	// ExposeHeaders defines a whitelist headers that clients are allowed to
	// access.
	//
	// Optional. Default value []string{}.
	ExposeHeaders []string

	// MaxAge indicates how long (in seconds) the results of a preflight request
	// can be cached.
	//
	// Optional. Default value 0.
	MaxAge int
}

var DefaultConfig = &Config{
	AllowOrigins:     []string{"*"},
	AllowMethods:     []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPost, http.MethodPatch, http.MethodDelete},
	AllowHeaders:     []string{},
	AllowCredentials: false,
	ExposeHeaders:    []string{},
	MaxAge:           0,
}
