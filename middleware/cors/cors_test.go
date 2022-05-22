package cors

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/felix-kaestner/lungo"
)

func assertEqual(t *testing.T, expected, actual any) {
	if reflect.DeepEqual(expected, actual) {
		return
	}

	t.Errorf("Test %s: Expected `%v` (type %v), Received `%v` (type %v)", t.Name(), expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
}

func TestCors(t *testing.T) {
	var tests = []struct {
		origin     string
		method     string
		status     int
		middleware lungo.Middleware
		eval       func(h http.Header)
	}{
		{
			origin:     "",
			method:     http.MethodGet,
			middleware: New(),
			eval: func(h http.Header) {
				assertEqual(t, "", h.Get(lungo.HeaderAccessControlAllowOrigin))
			},
		},
		{
			origin:     "",
			method:     http.MethodOptions,
			status:     http.StatusNoContent,
			middleware: New(),
			eval: func(h http.Header) {
				assertEqual(t, "", h.Get(lungo.HeaderAccessControlAllowOrigin))
			},
		},
		{
			origin: "http://example.de",
			method: http.MethodGet,
			middleware: New(func(c *Config) {
				c.AllowOrigins = []string{"http://example.com"}
			}),
			eval: func(h http.Header) {
				assertEqual(t, "", h.Get(lungo.HeaderAccessControlAllowOrigin))
			},
		},
		{
			origin: "http://example.de",
			method: http.MethodOptions,
			status: http.StatusNoContent,
			middleware: New(func(c *Config) {
				c.AllowOrigins = []string{"http://example.com"}
			}),
			eval: func(h http.Header) {
				assertEqual(t, "", h.Get(lungo.HeaderAccessControlAllowOrigin))
			},
		},
		{
			origin:     "localhost",
			method:     http.MethodGet,
			middleware: New(),
			eval: func(h http.Header) {
				assertEqual(t, lungo.HeaderOrigin, h.Get(lungo.HeaderVary))
				assertEqual(t, "*", h.Get(lungo.HeaderAccessControlAllowOrigin))
			},
		},
		{
			origin: "localhost",
			method: http.MethodGet,
			middleware: New(func(c *Config) {
				c.AllowCredentials = true
			}),
			eval: func(h http.Header) {
				assertEqual(t, lungo.HeaderOrigin, h.Get(lungo.HeaderVary))
				assertEqual(t, "true", h.Get(lungo.HeaderAccessControlAllowCredentials))
				assertEqual(t, "localhost", h.Get(lungo.HeaderAccessControlAllowOrigin))
			},
		},
		{
			origin: "localhost",
			method: http.MethodGet,
			middleware: New(func(c *Config) {
				c.AllowOrigins = []string{"localhost"}
				c.AllowCredentials = true
			}),
			eval: func(h http.Header) {
				assertEqual(t, lungo.HeaderOrigin, h.Get(lungo.HeaderVary))
				assertEqual(t, "true", h.Get(lungo.HeaderAccessControlAllowCredentials))
				assertEqual(t, "localhost", h.Get(lungo.HeaderAccessControlAllowOrigin))
			},
		},
		{
			origin: "localhost",
			method: http.MethodGet,
			middleware: New(func(c *Config) {
				c.AllowOrigins = []string{"localhost"}
				c.AllowCredentials = true
				c.ExposeHeaders = []string{lungo.HeaderAuthorization}
			}),
			eval: func(h http.Header) {
				assertEqual(t, lungo.HeaderOrigin, h.Get(lungo.HeaderVary))
				assertEqual(t, "true", h.Get(lungo.HeaderAccessControlAllowCredentials))
				assertEqual(t, "localhost", h.Get(lungo.HeaderAccessControlAllowOrigin))
				assertEqual(t, lungo.HeaderAuthorization, h.Get(lungo.HeaderAccessControlExposeHeaders))
			},
		},
		{
			origin: "localhost",
			method: http.MethodGet,
			middleware: New(func(c *Config) {
				c.AllowOrigins = []string{"localhost"}
				c.AllowCredentials = true
				c.ExposeHeaders = []string{lungo.HeaderAuthorization, lungo.HeaderContentType}
			}),
			eval: func(h http.Header) {
				assertEqual(t, "true", h.Get(lungo.HeaderAccessControlAllowCredentials))
				assertEqual(t, "localhost", h.Get(lungo.HeaderAccessControlAllowOrigin))
				assertEqual(t, lungo.HeaderAuthorization+","+lungo.HeaderContentType, h.Get(lungo.HeaderAccessControlExposeHeaders))
			},
		},
		{
			origin:     "localhost",
			method:     http.MethodOptions,
			middleware: New(),
			eval: func(h http.Header) {
				assertEqual(t, lungo.HeaderOrigin, h.Values(lungo.HeaderVary)[0])
				assertEqual(t, lungo.HeaderAccessControlRequestMethod, h.Values(lungo.HeaderVary)[1])
				assertEqual(t, lungo.HeaderAccessControlRequestHeaders, h.Values(lungo.HeaderVary)[2])
				assertEqual(t, "*", h.Get(lungo.HeaderAccessControlAllowOrigin))
				assertEqual(t, "GET,HEAD,PUT,POST,PATCH,DELETE", h.Get(lungo.HeaderAccessControlAllowMethods))
			},
		},
		{
			origin: "localhost",
			method: http.MethodOptions,
			middleware: New(func(c *Config) {
				c.AllowMethods = []string{http.MethodGet, http.MethodPost}
				c.AllowHeaders = []string{lungo.HeaderAuthorization, lungo.HeaderContentType}
				c.MaxAge = 3600
			}),
			eval: func(h http.Header) {
				assertEqual(t, lungo.HeaderOrigin, h.Values(lungo.HeaderVary)[0])
				assertEqual(t, lungo.HeaderAccessControlRequestMethod, h.Values(lungo.HeaderVary)[1])
				assertEqual(t, lungo.HeaderAccessControlRequestHeaders, h.Values(lungo.HeaderVary)[2])
				assertEqual(t, "*", h.Get(lungo.HeaderAccessControlAllowOrigin))
				assertEqual(t, "GET,POST", h.Get(lungo.HeaderAccessControlAllowMethods))
				assertEqual(t, lungo.HeaderAuthorization+","+lungo.HeaderContentType, h.Get(lungo.HeaderAccessControlAllowHeaders))
				assertEqual(t, "3600", h.Get(lungo.HeaderAccessControlMaxAge))
			},
		},
	}

	for _, testcase := range tests {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(testcase.method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		if testcase.origin != "" {
			req.Header.Set(lungo.HeaderOrigin, testcase.origin)
		}

		app := lungo.New()
		c := app.NewContext(rr, req)

		h := testcase.middleware(lungo.NotFoundHandler())
		h.ServeHTTP(c)

		testcase.eval(rr.Header())

		if testcase.status > 0 {
			assertEqual(t, testcase.status, rr.Code)
		}
	}
}
