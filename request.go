package lungo

import (
	"net/http"
	"path"
	"strings"
)

// Canonical returns a canonical path for p, eliminating . and .. elements.
func Canonical(p string) string {
	if p == "" {
		return "/"
	}

	if p[0] != '/' {
		p = "/" + p
	}

	cp := path.Clean(p)

	// path.Clean removes trailing slash except for root, thus
	// we need to put the trailing slash back if necessary.
	n := len(p) - 1
	if p[n] == '/' && cp != "/" {
		// Fast path for common case of p being the string we want
		if n == len(cp) && strings.HasPrefix(p, cp) {
			cp = p
		} else {
			cp += "/"
		}
	}

	return cp
}

// IsValidMethod checks if the provided http method is valid
func IsValidMethod(method string) bool {
	if method == "" {
		return false
	}

	/*
	   Method   = "OPTIONS"                ; RFC 7231 Section 9.2
	            | "GET"                    ; RFC 7231 Section 9.3
	            | "HEAD"                   ; RFC 7231 Section 9.4
	            | "POST"                   ; RFC 7231 Section 9.5
	            | "PUT"                    ; RFC 7231 Section 9.6
	            | "DELETE"                 ; RFC 7231 Section 9.7
	            | "TRACE"                  ; RFC 7231 Section 9.8
	            | "CONNECT"                ; RFC 7231 Section 9.9
	            | "PATCH"                  ; RFC 5789
	*/

	switch method {
	case http.MethodGet,
		http.MethodHead,
		http.MethodPost,
		http.MethodPut,
		http.MethodPatch,
		http.MethodDelete,
		http.MethodConnect,
		http.MethodOptions,
		http.MethodTrace:
		return true
	}

	return false
}
