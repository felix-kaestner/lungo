package template

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

func TestTemplate(t *testing.T) {
	var tests = []struct {
		middleware lungo.Middleware
		eval       func(h http.Header)
	}{
		{
			middleware: New(),
			eval:       func(h http.Header) {},
		},
		{
			middleware: New(func(c *Config) {

			}),
			eval: func(h http.Header) {},
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
