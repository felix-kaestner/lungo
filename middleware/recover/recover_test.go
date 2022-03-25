package recover

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

func TestRecover(t *testing.T) {
	var tests = []struct {
		middleware lungo.Middleware
		eval       func(e error)
	}{
		{
			middleware: New(),
			eval: func(e error) {
				assertEqual(t, "Error", e.Error())
			},
		},
		{
			middleware: New(func(c *Config) {
				c.HandleStackTrace = func(e any) {
					assertEqual(t, "Error", e)
				}
			}),
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

		h := testcase.middleware(lungo.HandlerFunc(func(c *lungo.Context) (err error) {
			panic("Error")
		}))
		h.ServeHTTP(c)
	}
}
