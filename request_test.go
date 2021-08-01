package lungo

import (
	"net/http"
	"testing"
)

func TestCanonical(t *testing.T) {
	assertEqual(t, "/", Canonical(""))
	assertEqual(t, "/", Canonical("."))
	assertEqual(t, "/", Canonical(".."))

	assertEqual(t, "/", Canonical("/."))
	assertEqual(t, "/", Canonical("/.."))
	assertEqual(t, "/", Canonical("/../"))
	assertEqual(t, "/", Canonical("./.."))
	assertEqual(t, "/", Canonical("/../."))
	assertEqual(t, "/", Canonical("/../.."))
	assertEqual(t, "/", Canonical("/test/.."))

	assertEqual(t, "/test", Canonical("test"))
	assertEqual(t, "/test", Canonical("/test"))

	assertEqual(t, "/test/", Canonical("test/"))
	assertEqual(t, "/test/", Canonical("/test/"))
	assertEqual(t, "/test/", Canonical("/test/./"))
}

func TestIsValidMethod(t *testing.T) {
	assertEqual(t, false, IsValidMethod(""))
	assertEqual(t, false, IsValidMethod(" "))
	assertEqual(t, false, IsValidMethod("BAD"))

	assertEqual(t, true, IsValidMethod(http.MethodGet))
	assertEqual(t, true, IsValidMethod(http.MethodHead))
	assertEqual(t, true, IsValidMethod(http.MethodPost))
	assertEqual(t, true, IsValidMethod(http.MethodPut))
	assertEqual(t, true, IsValidMethod(http.MethodPatch))
	assertEqual(t, true, IsValidMethod(http.MethodDelete))
	assertEqual(t, true, IsValidMethod(http.MethodConnect))
	assertEqual(t, true, IsValidMethod(http.MethodOptions))
	assertEqual(t, true, IsValidMethod(http.MethodTrace))

}
