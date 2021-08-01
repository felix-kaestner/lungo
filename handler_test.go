package lungo

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWithContext(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	ctx := &Context{
		Request:  req,
		Response: rr,
	}

	handler := WithContext(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		c := r.Context().Value("context")
		assertEqual(t, ctx, c)

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "Hello, world!")
	}))

	handler.ServeHTTP(ctx)

	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, "Hello, world!", rr.Body.String())
}
