package lungo

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoute(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	route := &Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: HandlerFunc(func(c *Context) error {
			return c.Text(http.StatusOK, "Hello, world!")
		}),
	}

	route.ServeHTTP(&Context{
		Request:  req,
		Response: rr,
	})

	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, "Hello, world!", rr.Body.String())
}
