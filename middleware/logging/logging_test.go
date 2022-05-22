package logging

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/felix-kaestner/lungo"
)

const LogRegex = `^\d{4}/\d{2}/\d{2} \d{2}:\d{2}:\d{2} GET / ((\d*[.])?\d+)µs$`

func assertEqual(t *testing.T, expected, actual any) {
	if reflect.DeepEqual(expected, actual) {
		return
	}

	t.Errorf("Test %s: Expected `%v` (type %v), Received `%v` (type %v)", t.Name(), expected, reflect.TypeOf(expected), actual, reflect.TypeOf(actual))
}

func TestLogging(t *testing.T) {
	var tests = []struct {
		middleware lungo.Middleware
		eval       func(b *bytes.Buffer)
	}{
		{
			middleware: New(func(c *Config) {
				c.Logger = nil
			}),
			eval: func(b *bytes.Buffer) {
				assertEqual(t, "", b.String())
			},
		},
		{
			middleware: New(func(c *Config) {
				// Invalid Template, fail on parse
				c.Template = "{{.Request.Method"
			}),
			eval: func(b *bytes.Buffer) {
				assertEqual(t, "", b.String())
			},
		},
		{
			middleware: New(func(c *Config) {
				// Invalid Template, fail on parse
				c.Template = `{{template "name"}}`
			}),
			eval: func(b *bytes.Buffer) {
				assertEqual(t, "", b.String())
			},
		},
		{
			middleware: New(func(c *Config) {
				c.Logger = log.Default()
			}),
			eval: func(b *bytes.Buffer) {
				match, err := regexp.MatchString(LogRegex, strings.TrimSpace(b.String()))
				if err != nil {
					t.Fatal(err)
				}
				assertEqual(t, true, match)
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

		b := &bytes.Buffer{}
		log.SetOutput(b)

		h := testcase.middleware(lungo.NotFoundHandler())
		h.ServeHTTP(c)

		testcase.eval(b)
	}
}

func TestTmp(t *testing.T) {
	tests := []string{
		"2021/05/09 13:10:59 GET / 73.9µs",
		"2021/05/09 13:13:09 GET / 37µs",
		"2021/05/09 13:14:32 GET / 39.6µs",
	}

	for _, s := range tests {
		match, err := regexp.MatchString(LogRegex, s)
		if err != nil {
			t.Fatal(err)
		}
		assertEqual(t, true, match)
	}
}
