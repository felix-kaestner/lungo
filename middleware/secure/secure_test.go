package secure

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

func TestSecure(t *testing.T) {
	var tests = []struct {
		middleware lungo.Middleware
		eval       func(h http.Header)
	}{
		{
			middleware: New(),
			eval: func(h http.Header) {
				assertEqual(t, "1; mode=block", h.Get(lungo.HeaderXXSSProtection))
				assertEqual(t, "SAMEORIGIN", h.Get(lungo.HeaderXFrameOptions))
				assertEqual(t, "", h.Get(lungo.HeaderContentSecurityPolicy))
				assertEqual(t, "nosniff", h.Get(lungo.HeaderXContentTypeOptions))
				assertEqual(t, "", h.Get(lungo.HeaderReferrerPolicy))
			},
		},
		{
			middleware: New(func(c *Config) {
				c.XSSProtection = "1"
				c.XFrameOptions = "DENY"
				c.ContentSecurityPolicy = "default-src https"
				c.ContentTypeNosniff = "nosniff"
				c.ReferrerPolicy = "no-referrer"
			}),
			eval: func(h http.Header) {
				assertEqual(t, "1", h.Get(lungo.HeaderXXSSProtection))
				assertEqual(t, "DENY", h.Get(lungo.HeaderXFrameOptions))
				assertEqual(t, "default-src https", h.Get(lungo.HeaderContentSecurityPolicy))
				assertEqual(t, "nosniff", h.Get(lungo.HeaderXContentTypeOptions))
				assertEqual(t, "no-referrer", h.Get(lungo.HeaderReferrerPolicy))
			},
		},
	}

	for _, testcase := range tests {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		app := lungo.New()
		c := app.NewContext(rr, req)

		h := testcase.middleware(lungo.NotFoundHandler())
		h.ServeHTTP(c)

		testcase.eval(rr.Header())
	}
}
