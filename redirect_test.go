package lungo

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRedirect(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()

	handler := RedirectHandler("/hello", http.StatusMovedPermanently)

	handler.ServeHTTP(&Context{
		Request:  req,
		Response: rr,
	})

	header := rr.Header()

	assertEqual(t, http.StatusMovedPermanently, rr.Code)
	assertEqual(t, "text/html; charset=utf-8", header.Get("Content-Type"))
	assertEqual(t, "/hello", header.Get("Location"))
	assertEqual(t, "<a href=\"/hello\">"+http.StatusText(http.StatusMovedPermanently)+"</a>.", strings.Replace(rr.Body.String(), "\n", "", -1))
}
