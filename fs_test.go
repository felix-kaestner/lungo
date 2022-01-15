package lungo

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestFileServer(t *testing.T) {
	var tests = []struct {
		dir         string
		path        string
		contentType string
	}{
		{
			dir:         ".github/",
			path:        "/labeler.yml",
			contentType: "application/x-yaml",
		},
		{
			dir:         ".github/",
			path:        "labeler.yml",
			contentType: "application/x-yaml",
		},
		{
			dir:         "./",
			path:        "labeler.yml",
			contentType: "",
		},
	}

	for _, testcase := range tests {
		req, err := http.NewRequest("GET", testcase.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		rr := httptest.NewRecorder()

		h := FileHandler(testcase.dir)

		h.ServeHTTP(&Context{
			Request:  req,
			Response: rr,
		})

		header := rr.Header()

		ct := header.Get(HeaderContentType)
		if ct != testcase.contentType && ct != MIMETextPlainCharsetUTF8 {
			t.Errorf("Test %s: Expected HTTP Content-Type to be either `%v` or `%v` (type string), Received `%v` (type %v)", t.Name(), testcase.contentType, "text/plain; charset=utf-8", ct, reflect.TypeOf(ct))
		}
	}
}
