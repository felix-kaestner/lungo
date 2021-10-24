package lungo

import (
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContextFlush(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	c.Flush()

	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, true, rr.Flushed)
}

type HijackableResponseRecorder struct {
	delegate *httptest.ResponseRecorder
	conn     net.Conn
}

func NewHijackableRecorder(conn net.Conn) http.ResponseWriter {
	return &HijackableResponseRecorder{
		delegate: httptest.NewRecorder(),
		conn:     conn,
	}
}

// Implement `Header` method of http.ResponseWriter interface
func (h *HijackableResponseRecorder) Header() http.Header {
	return h.delegate.Header()
}

// Implement `Write` method of http.ResponseWriter interface
func (h *HijackableResponseRecorder) Write(buf []byte) (int, error) {
	return h.delegate.Write(buf)
}

// Implement `WriteHeader` method of http.ResponseWriter interface
func (h *HijackableResponseRecorder) WriteHeader(code int) {
	h.delegate.WriteHeader(code)
}

// Implement `Hijack` method of http.Hijacker interface
func (h *HijackableResponseRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return h.conn, bufio.NewReadWriter(bufio.NewReader(h.conn), bufio.NewWriter(h.conn)), nil
}

func TestContextHijack(t *testing.T) {
	app := New()

	ln, err := net.Listen("tcp", ":8000")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		t.Fatal(err)
	}

	rr := NewHijackableRecorder(conn)

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	conn, bufwr, err := c.Hijack()
	if err != nil {
		t.Fatal(err)
	}

	bufwr.WriteString("HTTP/1.1 200 Awesome\n\n")
	bufwr.Flush()
	conn.Close()
}

func TestContextParam(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/?id=1", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	assertEqual(t, "1", c.Param("id"))

	c.SetParam("id", "2")
	assertEqual(t, "2", c.Param("id"))

	c.AddParam("name", "foo")
	assertEqual(t, "foo", c.Param("name"))
	assertEqual(t, "foo", c.ParamOrDefault("name", "bar"))

	c.DeleteParam("name")
	assertEqual(t, "", c.Param("name"))
	assertEqual(t, "bar", c.ParamOrDefault("name", "bar"))
}

func TestContextHeader(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	assertEqual(t, "", c.Header(HeaderOrigin))

	c.SetHeader(HeaderVary, HeaderOrigin)
	assertEqual(t, HeaderOrigin, c.Response.Header().Get(HeaderVary))

	c.AddHeader(HeaderVary, HeaderAccessControlRequestMethod)
	assertEqual(t, HeaderOrigin, c.Response.Header().Values(HeaderVary)[0])
	assertEqual(t, HeaderAccessControlRequestMethod, c.Response.Header().Values(HeaderVary)[1])
}

func TestContextCookie(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	header := rr.Header()
	c := app.NewContext(rr, req)

	assertEqual(t, 0, len(c.Cookies()))

	_, err = c.Cookie("User")
	assertEqual(t, http.ErrNoCookie, err)

	c.SetCookie(&http.Cookie{Name: "User", Value: "John", Path: "/"})
	assertEqual(t, "User=John; Path=/", header.Get("Set-Cookie"))
}

func TestContextParseMediaType(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)

	}

	req.Header.Set(HeaderContentType, MIMEApplicationJSON)

	c := app.NewContext(rr, req)

	mt, _, err := c.ParseMediaType()
	assertNil(t, err)
	assertEqual(t, mt, MIMEApplicationJSON)
}

func TestContextError(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	err = c.Error(http.StatusBadRequest)
	re, ok := err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError.")
	}

	assertEqual(t, http.StatusBadRequest, re.Code)
	assertEqual(t, http.StatusText(http.StatusBadRequest), re.Message)

	err = c.Errorf(http.StatusBadRequest, "Foo")
	re, ok = err.(*RequestError)
	if !ok {
		t.Errorf("Expected RequestError.")
	}

	assertEqual(t, http.StatusBadRequest, re.Code)
	assertEqual(t, "Foo", re.Message)
}

func TestContextFile(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	c.File("lungo.go")

	header := rr.Header()

	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, "text/x-go; charset=utf-8", header.Get(HeaderContentType))
}

func TestContextNoContent(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	c.NoContent()
	assertEqual(t, http.StatusNoContent, rr.Code)
	assertEqual(t, "", rr.Body.String())
}

func TestContextText(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	c.Text(http.StatusOK, "Hello, world!")
	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, "Hello, world!", rr.Body.String())
}

func TestContextJson(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	c.Json(http.StatusOK, Map{"message": "Hello, world!"})
	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, "{\"message\":\"Hello, world!\"}", strings.TrimSpace(rr.Body.String()))
}

func TestContextDecodeJSONBody(t *testing.T) {
	var tests = []struct {
		body      io.Reader
		mime      string
		configure func(*Config)
		eval      func(err error)
	}{
		{
			body: bytes.NewReader([]byte("{\"msg\":\"Hello, world!\"}")),
			mime: MIMEApplicationXML,
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusUnsupportedMediaType, re.Code)
			},
		},
		{
			body: nil,
			mime: MIMEApplicationJSON,
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusBadRequest, re.Code)
			},
		},
		{
			body: bytes.NewReader([]byte("{\"msg\": Hello}")),
			mime: MIMEApplicationJSON,
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusBadRequest, re.Code)
			},
		},
		{
			body: bytes.NewReader([]byte("{\"msg\":\"}")),
			mime: MIMEApplicationJSON,
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusBadRequest, re.Code)
			},
		},
		{
			body: bytes.NewReader([]byte("{\"msg\": []}")),
			mime: MIMEApplicationJSON,
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusBadRequest, re.Code)
			},
		},
		{
			body: bytes.NewReader([]byte("{\"bad\":\"Hello, world!\"}")),
			mime: MIMEApplicationJSON,
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusBadRequest, re.Code)
			},
		},
		{
			body: bytes.NewReader([]byte("")),
			mime: MIMEApplicationJSON,
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusBadRequest, re.Code)
			},
		},
		{
			body: bytes.NewReader([]byte("{\"msg\":\"Hello, world!\"}")),
			mime: MIMEApplicationJSON,
			configure: func(c *Config) {
				c.MaxBodySize = 0
			},
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusRequestEntityTooLarge, re.Code)
			},
		},
		{
			body: bytes.NewReader([]byte("{\"msg\":\"Hello, world!\"}{\"msg\":\"Hello, world!\"}")),
			mime: MIMEApplicationJSON,
			eval: func(err error) {
				re, ok := err.(*RequestError)
				if !ok {
					t.Errorf("Expected RequestError.")
				}
				assertEqual(t, http.StatusBadRequest, re.Code)
			},
		},
		{
			body: bytes.NewReader([]byte("{\"msg\":\"Hello, world!\"}")),
			mime: MIMEApplicationJSON,
			eval: func(err error) {
				assertNil(t, err)
			},
		},
	}

	for _, testcase := range tests {
		configure := testcase.configure
		if configure == nil {
			configure = func(c *Config) {}
		}
		app := New(configure)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("POST", "/", testcase.body)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set(HeaderContentType, testcase.mime)

		c := app.NewContext(rr, req)

		err = c.DecodeJSONBody(&struct {
			Msg string `json:"msg"`
		}{})

		testcase.eval(err)
	}
}

func TestContext(t *testing.T) {
	app := New()

	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	c := app.NewContext(rr, req)

	assertEqual(t, http.MethodGet, c.Method())
	assertEqual(t, "/", c.Path())
}
