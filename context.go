package lungo

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net"
	"net/http"
	"net/url"
	"strings"
)

// Context represents the context of a request, including request,
// response, path, parameters, headers and handler for common operations.
type Context struct {
	App      *App                `json:"-"`
	Request  *http.Request       `json:"-"`
	Response http.ResponseWriter `json:"-"`
	Params   url.Values          `json:"params"`
}

// Reset applies the given request to the Context instance.
//
// The `App` property is assumed to be static, thus it will not be changed.
func (c *Context) Reset(w http.ResponseWriter, r *http.Request) {
	params, _ := url.ParseQuery(r.URL.RawQuery)

	c.Request = r
	c.Response = w
	c.Params = params
}

// Flush implements the http.Flusher interface to allow an HTTP handler to flush
// buffered data to the client.
// See [http.Flusher](https://golang.org/pkg/net/http/#Flusher)
func (c *Context) Flush() {
	c.Response.(http.Flusher).Flush()
}

// Hijack implements the http.Hijacker interface to allow an HTTP handler to
// take over the connection.
// See [http.Hijacker](https://golang.org/pkg/net/http/#Hijacker)
func (c *Context) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return c.Response.(http.Hijacker).Hijack()
}

// Method specifies the HTTP method (GET, POST, PUT, etc.).
// For client requests, an empty string means GET.
//
// Go's HTTP client does not support sending a request with
// the CONNECT method. See the documentation on Transport for
// details.
func (c *Context) Method() string {
	return c.Request.Method
}

// Path specifies either the URI being requested (for server
// requests) or the URL to access (for client requests).
//
// For server requests, the URL is parsed from the URI
// supplied on the Request-Line as stored in RequestURI.  For
// most requests, fields other than Path and RawQuery will be
// empty. (See RFC 7230, Section 5.3)
//
// For client requests, the URL's Host specifies the server to
// connect to, while the Request's Host field optionally
// specifies the Host header value to send in the HTTP
// request.
func (c *Context) Path() string {
	return c.Request.URL.Path
}

// Param gets the first value associated with the given parameter.
// If there are no values associated with the key, Get returns
// the empty string. To access multiple values, use the map
// directly.
func (c *Context) Param(key string) string {
	return c.Params.Get(key)
}

// ParamOrDefault gets the first value associated with the given parameter
// or the provided default value if there are no values associated
// with the key.
func (c *Context) ParamOrDefault(key, val string) string {
	p := c.Param(key)
	if p == "" {
		return val
	}
	return p
}

// SetParam sets the parameter value. It replaces any existing
// values.
func (c *Context) SetParam(key, value string) {
	c.Params.Set(key, value)
}

// AddParam adds the value to the parameter's values. It appends
// to any existing values associated with key.
func (c *Context) AddParam(key, value string) {
	c.Params.Add(key, value)
}

// DeleteParam deletes the values associated with the parameter.
func (c *Context) DeleteParam(key string) {
	c.Params.Del(key)
}

// Header gets the first value associated with the given key. If
// there are no values associated with the key, it returns "".
// It is case insensitive; textproto.CanonicalMIMEHeaderKey is
// used to canonicalize the provided key. To use non-canonical keys,
// access the map directly.
func (c *Context) Header(key string) string {
	return c.Request.Header.Get(key)
}

// SetHeader sets the header entries associated with key to the
// single element value. It replaces any existing values
// associated with key. The key is case insensitive; it is
// canonicalized by textproto.CanonicalMIMEHeaderKey.
// To use non-canonical keys, assign to the map directly.
func (c *Context) SetHeader(key, value string) {
	c.Response.Header().Set(key, value)
}

// AddHeader adds the key, value pair to the header.
// It appends to any existing values associated with key.
// The key is case insensitive; it is canonicalized by
// CanonicalHeaderKey.
func (c *Context) AddHeader(key, value string) {
	c.Response.Header().Add(key, value)
}

// WriteHeader sends an HTTP response header with the provided
// status code.
//
// If WriteHeader is not called explicitly, the first call to Write
// will trigger an implicit WriteHeader(http.StatusOK).
// Thus explicit calls to WriteHeader are mainly used to
// send error codes.
//
// The provided code must be a valid HTTP 1xx-5xx status code.
// Only one header may be written. Go does not currently
// support sending user-defined 1xx informational headers,
// with the exception of 100-continue response header that the
// Server sends automatically when the Request.Body is read.
func (c *Context) WriteHeader(code int) {
	c.Response.WriteHeader(code)
}

// Cookie returns the named cookie provided in the request or
// ErrNoCookie if not found.
// If multiple cookies match the given name, only one cookie will
// be returned.
func (c *Context) Cookie(name string) (*http.Cookie, error) {
	return c.Request.Cookie(name)
}

// SetCookie adds a Set-Cookie header to the provided ResponseWriter's headers.
// The provided cookie must have a valid Name. Invalid cookies may be
// silently dropped.
func (c *Context) SetCookie(cookie *http.Cookie) {
	http.SetCookie(c.Response, cookie)
}

// Cookies parses and returns the HTTP cookies sent with the request.
func (c *Context) Cookies() []*http.Cookie {
	return c.Request.Cookies()
}

// ParseMediaType parses a media type value and any optional
// parameters, per RFC 1521.  Media types are the values in
// Content-Type and Content-Disposition headers (RFC 2183).
// On success, ParseMediaType returns the media type converted
// to lowercase and trimmed of white space and a non-nil map.
// If there is an error parsing the optional parameter,
// the media type will be returned along with the error
// ErrInvalidMediaParameter.
// The returned map, params, maps from the lowercase
// attribute to the attribute value with its case preserved.
func (c *Context) ParseMediaType() (mediatype string, params map[string]string, err error) {
	ct := c.Header(HeaderContentType)
	return mime.ParseMediaType(ct)
}

// Error dispatches a error response.
// Use the method parameter `code` to set status code of the response.
// The message will automatically be set to the corresponding http.StatusText.
func (c *Context) Error(code int) (err error) {
	err = &RequestError{Code: code, Message: http.StatusText(code)}
	return
}

// Errorf dispatches a error response with custom message.
// Use the method parameter `code` to set status code of the response.
// Use the method parameter `message` to set the message of the response.
func (c *Context) Errorf(code int, message interface{}) (err error) {
	msg := fmt.Sprintf("%v", message)
	err = &RequestError{Code: code, Message: msg}
	return
}

// File dispatches a response serving a file.
// Use the method parameter `name` to supply the filename.
func (c *Context) File(name string) (err error) {
	http.ServeFile(c.Response, c.Request, name)
	return
}

// NoContent dispatchss an empty response with an HTTP 204 status code.
func (c *Context) NoContent() (err error) {
	c.WriteHeader(http.StatusNoContent)
	return
}

// NotFound replies to the request with an HTTP 404 not found error.
func (c *Context) NotFound() error {
	return c.Error(http.StatusNotFound)
}

// Text dispatches a text response.
// Use the method parameter `code` to set the header status code.
// Use the method parameter `value` to supply the object to be serialized to text.
func (c *Context) Text(code int, value interface{}) (err error) {
	c.SetHeader(HeaderContentType, MIMETextPlain)
	c.WriteHeader(code)
	_, err = fmt.Fprint(c.Response, value)
	return
}

// Json dispatches a JSON response.
// Use the method parameter `code` to set the header status code.
// Use the method parameter `value` to supply the object to be serialized to JSON.
func (c *Context) Json(code int, value interface{}) (err error) {
	c.SetHeader(HeaderContentType, MIMEApplicationJSON)
	c.WriteHeader(code)
	err = json.NewEncoder(c.Response).Encode(value)
	return
}

// DecodeJSONBody decodes an object with a given interface from a JSON request body.
func (c *Context) DecodeJSONBody(dst interface{}) error {

	// check that the the Content-Type header has the value application/json.
	if mt, _, err := c.ParseMediaType(); err != nil || !strings.HasPrefix(mt, MIMEApplicationJSON) {
		msg := fmt.Sprintf("Request header for `%s` is not set to `%s`", HeaderContentType, MIMEApplicationJSON)
		return &RequestError{Code: http.StatusUnsupportedMediaType, Message: msg}
	}

	// check that the request contains a body
	if c.Request.Body == nil {
		return &RequestError{Code: http.StatusBadRequest, Message: "Request body is unset"}
	}

	// Use http.MaxBytesReader to enforce a maximum read from the
	// response body according to application configuration.
	// A request body larger than that will now result in
	// DecodeJSONBody() returning a "http: request body too large" error.
	if c.App.config != nil {
		if c.App.config.MaxBodySize > -1 {
			c.Request.Body = http.MaxBytesReader(c.Response, c.Request.Body, int64(c.App.config.MaxBodySize))
		}
	}

	// Setup the decoder and call the DisallowUnknownFields() method on it.
	// This will cause DecodeJSONBody() to return a "json: unknown field ..." error
	// if it encounters any extra unexpected fields in the JSON. Strictly
	// speaking, it returns an error for "keys which do not match any
	// non-ignored, exported fields in the destination".
	dec := json.NewDecoder(c.Request.Body)
	dec.DisallowUnknownFields()

	err := dec.Decode(dst)
	if err != nil {
		var syntaxError *json.SyntaxError
		var unmarshalTypeError *json.UnmarshalTypeError

		switch {
		// Catch any syntax errors in the JSON and send an error message
		// which interpolates the location of the problem to make it
		// easier for the client to fix.
		case errors.As(err, &syntaxError):
			msg := fmt.Sprintf("Request body contains badly-formed JSON (at position %d)", syntaxError.Offset)
			return &RequestError{Code: http.StatusBadRequest, Message: msg}

		// In some circumstances DecodeJSONBody() may also return an
		// io.ErrUnexpectedEOF error for syntax errors in the JSON. There
		// is an open issue regarding this at
		// https://github.com/golang/go/issues/25956.
		case errors.Is(err, io.ErrUnexpectedEOF):
			msg := "Request body contains badly-formed JSON"
			return &RequestError{Code: http.StatusBadRequest, Message: msg}

		// Catch any type errors, like trying to assign a string in the
		// JSON request body to a int field. We can
		// interpolate the relevant field name and position into the error
		// message to make it easier for the client to fix.
		case errors.As(err, &unmarshalTypeError):
			msg := fmt.Sprintf("Request body contains an invalid value for the %q field (at position %d)", unmarshalTypeError.Field, unmarshalTypeError.Offset)
			return &RequestError{Code: http.StatusBadRequest, Message: msg}

		// Catch the error caused by extra unexpected fields in the request
		// body. We extract the field name from the error message and
		// interpolate it in our custom error message. There is an open
		// issue at https://github.com/golang/go/issues/29035 regarding
		// turning this into a sentinel error.
		case strings.HasPrefix(err.Error(), "json: unknown field "):
			fieldName := strings.TrimPrefix(err.Error(), "json: unknown field ")
			msg := fmt.Sprintf("Request body contains unknown field %s", fieldName)
			return &RequestError{Code: http.StatusBadRequest, Message: msg}

		// An io.EOF error is returned by DecodeJSONBody() if the request body
		// is empty.
		case errors.Is(err, io.EOF):
			msg := "Request body must not be empty"
			return &RequestError{Code: http.StatusBadRequest, Message: msg}

		// Catch the error caused by the request body being too large. Again
		// there is an open issue regarding turning this into a sentinel
		// error at https://github.com/golang/go/issues/30715.
		case err.Error() == "http: request body too large":
			msg := fmt.Sprintf("Request body must not be larger than %s", byteToBinaryIEC(int64(c.App.config.MaxBodySize)))
			return &RequestError{Code: http.StatusRequestEntityTooLarge, Message: msg}

		default:
			return err
		}
	}

	err = dec.Decode(&struct{}{})
	if err != io.EOF {
		msg := "Request body must only contain a single JSON object"
		return &RequestError{Code: http.StatusBadRequest, Message: msg}
	}

	return nil
}
