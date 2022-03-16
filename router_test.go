package lungo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouterInvalidRouteMethod(t *testing.T) {
	assertPanic(t, ErrInvalidRouteMethod, func() {
		router := NewRouter()

		router.Handle(Route{
			Path: "/",
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})
	})

	assertPanic(t, ErrInvalidRouteMethod, func() {
		router := NewRouter()

		router.Handle(Route{
			Method: "",
			Path:   "/",
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})
	})

	assertPanic(t, ErrInvalidRouteMethod, func() {
		router := NewRouter()

		router.Handle(Route{
			Method: "bad",
			Path:   "/",
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})
	})
}

func TestRouterHandleEmptyPath(t *testing.T) {
	assertPanic(t, ErrEmptyRoutePath, func() {
		router := NewRouter()

		router.Handle(Route{
			Method: http.MethodGet,
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})
	})

	assertPanic(t, ErrEmptyRoutePath, func() {
		router := NewRouter()

		router.Handle(Route{
			Method: http.MethodGet,
			Path:   "",
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})
	})
}

func TestRouterHandleNilHandler(t *testing.T) {
	assertPanic(t, ErrNilRouteHandler, func() {
		router := NewRouter()

		router.Handle(Route{
			Method: http.MethodGet,
			Path:   "/",
		})
	})

	assertPanic(t, ErrNilRouteHandler, func() {
		router := NewRouter()

		router.Handle(Route{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: nil,
		})
	})
}

func TestRouterHandleDuplicateHandler(t *testing.T) {
	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: HandlerFunc(func(c *Context) error {
			return nil
		}),
	})

	assertPanic(t, fmt.Sprintf(ErrDuplicateHandler, "/"), func() {
		router.Handle(Route{
			Method: http.MethodGet,
			Path:   "/",
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})
	})
}

func TestRouterMatch(t *testing.T) {
	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/a",
		Handler: HandlerFunc(func(c *Context) error {
			return nil
		}),
	})

	r := router.match(http.MethodGet, "/a")
	assertEqual(t, "/a", r.Path)
	assertEqual(t, http.MethodGet, r.Method)

	r = router.match(http.MethodGet, "/a/b")
	assertEqual(t, "/a", r.Path)
	assertEqual(t, http.MethodGet, r.Method)

	r = router.match(http.MethodGet, "/b")
	assertNil(t, r)

	r = router.match(http.MethodPost, "/a")
	assertNil(t, r)
}

func TestRouterShouldRedirectEmptyPath(t *testing.T) {
	req, err := http.NewRequest("GET", "", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/a",
		Handler: HandlerFunc(func(c *Context) error {
			return nil
		}),
	})

	_, redirect := router.shouldRedirect(http.MethodGet, req.URL.Path, req.URL)

	assertEqual(t, false, redirect)
}
func TestRouterShouldRedirectExistingPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/a",
		Handler: HandlerFunc(func(c *Context) error {
			return nil
		}),
	})

	_, redirect := router.shouldRedirect(http.MethodGet, req.URL.Path, req.URL)

	assertEqual(t, false, redirect)
}

func TestRouterShouldRedirectNonExistingPath(t *testing.T) {
	req, err := http.NewRequest("GET", "/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/b",
		Handler: HandlerFunc(func(c *Context) error {
			return nil
		}),
	})

	_, redirect := router.shouldRedirect(http.MethodGet, req.URL.Path, req.URL)

	assertEqual(t, false, redirect)
}

func TestRouterShouldRedirectWithTralingSlash(t *testing.T) {
	req, err := http.NewRequest("GET", "/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/a/",
		Handler: HandlerFunc(func(c *Context) error {
			return nil
		}),
	})

	_, redirect := router.shouldRedirect(http.MethodGet, req.URL.Path, req.URL)

	assertEqual(t, true, redirect)
}

func TestRouterShouldRedirectWithoutTralingSlash(t *testing.T) {
	req, err := http.NewRequest("GET", "/a/", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/a",
		Handler: HandlerFunc(func(c *Context) error {
			return nil
		}),
	})

	_, redirect := router.shouldRedirect(http.MethodGet, req.URL.Path, req.URL)

	assertEqual(t, true, redirect)
}

func TestRouterHandlerRedirect(t *testing.T) {
	{
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("CONNECT", "/a", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := NewRouter()

		router.Handle(Route{
			Method: http.MethodConnect,
			Path:   "/a/",
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})

		handler := router.Handler(req)

		handler.ServeHTTP(&Context{
			Request:  req,
			Response: rr,
		})

		header := rr.Header()

		assertEqual(t, http.StatusMovedPermanently, rr.Code)
		assertEqual(t, "/a/", header.Get("Location"))
	}

	{
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/a", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := NewRouter()

		router.Handle(Route{
			Method: http.MethodGet,
			Path:   "/a/",
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})

		handler := router.Handler(req)

		handler.ServeHTTP(&Context{
			Request:  req,
			Response: rr,
		})

		header := rr.Header()

		assertEqual(t, http.StatusMovedPermanently, rr.Code)
		assertEqual(t, "/a/", header.Get("Location"))
	}

	{
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", "/a/..", nil)
		if err != nil {
			t.Fatal(err)
		}

		router := NewRouter()

		router.Handle(Route{
			Method: http.MethodGet,
			Path:   "/",
			Handler: HandlerFunc(func(c *Context) error {
				return nil
			}),
		})

		handler := router.Handler(req)

		handler.ServeHTTP(&Context{
			Request:  req,
			Response: rr,
		})

		header := rr.Header()

		assertEqual(t, http.StatusMovedPermanently, rr.Code)
		assertEqual(t, "/", header.Get("Location"))
	}
}

func TestRouterHandlerNotFound(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/a", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/b",
		Handler: HandlerFunc(func(c *Context) error {
			return c.Text(http.StatusOK, "Hello, world!")
		}),
	})

	handler := router.Handler(req)

	err = handler.ServeHTTP(&Context{
		Request:  req,
		Response: rr,
	})

	re, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected Not Found RequestError.")
	}

	assertEqual(t, http.StatusNotFound, re.Code)
	assertEqual(t, "Not Found", re.Message)
}

func TestRouter(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: HandlerFunc(func(c *Context) error {
			return c.Text(http.StatusOK, "Hello, world!")
		}),
	})

	router.ServeHTTP(&Context{
		Request:  req,
		Response: rr,
	})

	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, "Hello, world!", rr.Body.String())
}

func TestRouterInvalidRequest(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	req.RequestURI = "*"

	err = router.ServeHTTP(&Context{
		Request:  req,
		Response: rr,
	})

	re, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError.")
	}

	assertEqual(t, http.StatusBadRequest, re.Code)
	assertEqual(t, http.StatusText(http.StatusBadRequest), re.Message)

	header := rr.Header()

	assertEqual(t, "close", header.Get("Connection"))
}

func TestRouterMiddleware(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	router := NewRouter()

	router.Handle(Route{
		Method: http.MethodGet,
		Path:   "/",
		Handler: HandlerFunc(func(c *Context) error {
			return c.Text(http.StatusOK, "Hello, world!")
		}),
	})

	router.Use(func(next Handler) Handler {
		return HandlerFunc(func(c *Context) error {
			c.Response.Header().Add("X-XSS-Protection", "1; mode=blockFilter")

			return next.ServeHTTP(c)
		})
	})

	router.ServeHTTP(&Context{
		Request:  req,
		Response: rr,
	})

	header := rr.Header()

	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, "Hello, world!", rr.Body.String())
	assertEqual(t, "1; mode=blockFilter", header.Get("X-XSS-Protection"))
}
