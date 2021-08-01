package lungo

import (
	"net/http"
	"testing"
)

func TestError(t *testing.T) {
	err := &RequestError{
		Code:    http.StatusBadRequest,
		Message: "Bad Request",
	}

	assertEqual(t, "Bad Request", err.Error())
}
