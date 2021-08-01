package lungo

import (
	"net/http"
	"net/http/httptest"
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
			contentType: MIMETextPlainCharsetUTF8,
		},
		{
			dir:         ".github/",
			path:        "labeler.yml",
			contentType: MIMETextPlainCharsetUTF8,
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
		assertEqual(t, testcase.contentType, header.Get(HeaderContentType))
	}
}
