package lungo

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"errors"
	"math/big"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestAppContext(t *testing.T) {
	app := New()
	c := app.AcquireContext()
	assertNotNil(t, c.App)
	app.ReleaseContext(c)
}

func TestAppConfig(t *testing.T) {
	app := New()

	assertNotNil(t, app.Config())
	assertEqual(t, DefaultMaxBodySize, app.Config().MaxBodySize)

	app = New(func(c *Config) { c.MaxBodySize = 1024 })
	assertEqual(t, 1024, app.Config().MaxBodySize)
}

func TestAppServer(t *testing.T) {
	app := New()

	assertNil(t, app.Server())
}

func TestAppServe(t *testing.T) {
	app := New()

	ln, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		t.Fatal(err)
	}

	cerr := make(chan error)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cerr <- app.Serve(ln)
	}()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

wait:
	for {
		select {
		case <-time.After(time.Second):
			break wait

		case <-ticker.C:
			_, err := http.Get("http://0.0.0.0:8000")
			if err == nil {
				break wait
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}

	err = <-cerr
	if err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}
}

func TestAppServeTLS(t *testing.T) {
	app := New()

	ln, err := net.Listen("tcp", "0.0.0.0:8000")
	if err != nil {
		t.Fatal(err)
	}

	GenerateCertificate(t)

	cerr := make(chan error)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cerr <- app.ServeTLS(ln, "cert.pem", "key.pem")
	}()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

wait:
	for {
		select {
		case <-time.After(time.Second):
			break wait

		case <-ticker.C:
			_, err := http.Get("http://0.0.0.0:8000")
			if err == nil {
				break wait
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}

	err = <-cerr
	if err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}
}

func TestAppListen(t *testing.T) {
	app := New()

	cerr := make(chan error)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cerr <- app.Listen("0.0.0.0:8000")
	}()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

wait:
	for {
		select {
		case <-time.After(time.Second):
			break wait

		case <-ticker.C:
			_, err := http.Get("http://0.0.0.0:8000")
			if err == nil {
				break wait
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}

	err := <-cerr
	if err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}
}

func TestAppListenTLS(t *testing.T) {
	app := New()

	GenerateCertificate(t)

	cerr := make(chan error)
	go func() {
		time.Sleep(10 * time.Millisecond)
		cerr <- app.ListenTLS("0.0.0.0:8000", "cert.pem", "key.pem")
	}()

	ticker := time.NewTicker(5 * time.Millisecond)
	defer ticker.Stop()

wait:
	for {
		select {
		case <-time.After(time.Second):
			break wait

		case <-ticker.C:
			_, err := http.Get("http://0.0.0.0:8000")
			if err == nil {
				break wait
			}
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		t.Fatal(err)
	}

	err := <-cerr
	if err != nil && err != http.ErrServerClosed {
		t.Fatal(err)
	}
}

func TestAppShutdown(t *testing.T) {
	app := New()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := app.Shutdown(ctx)

	assertNotNil(t, err)
	assertEqual(t, "Shutdown: Server is not running.", err.Error())
}

func TestApp(t *testing.T) {
	var tests = []struct {
		path    string
		method  string
		status  int
		handler HandlerFunc
	}{
		{
			path:   "/get",
			method: http.MethodGet,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/head",
			method: http.MethodHead,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/post",
			method: http.MethodPost,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/put",
			method: http.MethodPut,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/patch",
			method: http.MethodPatch,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/delete",
			method: http.MethodDelete,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/connect",
			method: http.MethodConnect,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/options",
			method: http.MethodOptions,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/trace",
			method: http.MethodTrace,
			status: http.StatusOK,
			handler: func(c *Context) error {
				return c.Text(http.StatusOK, "Hello, world!")
			},
		},
		{
			path:   "/error/a",
			method: http.MethodGet,
			status: http.StatusBadRequest,
			handler: func(c *Context) error {
				return c.Error(http.StatusBadRequest)
			},
		},
		{
			path:   "/error/b",
			method: http.MethodGet,
			status: http.StatusInternalServerError,
			handler: func(c *Context) error {
				return errors.New("Foo")
			},
		},
	}
	for _, testcase := range tests {
		app := New()

		switch testcase.method {
		case http.MethodGet:
			app.Get(testcase.path, testcase.handler)
		case http.MethodHead:
			app.Head(testcase.path, testcase.handler)
		case http.MethodPost:
			app.Post(testcase.path, testcase.handler)
		case http.MethodPut:
			app.Put(testcase.path, testcase.handler)
		case http.MethodPatch:
			app.Patch(testcase.path, testcase.handler)
		case http.MethodDelete:
			app.Delete(testcase.path, testcase.handler)
		case http.MethodConnect:
			app.Connect(testcase.path, testcase.handler)
		case http.MethodOptions:
			app.Options(testcase.path, testcase.handler)
		case http.MethodTrace:
			app.Trace(testcase.path, testcase.handler)
		}

		rr := httptest.NewRecorder()
		req, err := http.NewRequest(testcase.method, testcase.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		app.ServeHTTP(rr, req)

		assertEqual(t, testcase.status, rr.Code)
	}
}

func TestAppAll(t *testing.T) {
	app := New()

	app.All("/", func(c *Context) error {
		return c.Text(http.StatusOK, "Hello, world!")
	})

	for _, method := range methods {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest(method, "/", nil)
		if err != nil {
			t.Fatal(err)
		}

		app.ServeHTTP(rr, req)

		assertEqual(t, http.StatusOK, rr.Code)
	}
}

func TestAppStatic(t *testing.T) {
	var tests = []struct {
		dir    string
		path   string
		status int
	}{
		{
			dir:    ".github/",
			path:   "/labeler.yml",
			status: http.StatusOK,
		},
		{
			dir:    "./",
			path:   "/labeler.yml",
			status: http.StatusNotFound,
		},
	}

	for _, testcase := range tests {
		app := New()
		app.Static("/", testcase.dir)

		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", testcase.path, nil)
		if err != nil {
			t.Fatal(err)
		}

		app.ServeHTTP(rr, req)

		assertEqual(t, testcase.status, rr.Code)
	}
}

func TestAppMiddleware(t *testing.T) {
	rr := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	app := New()

	app.Get("/", func(c *Context) error {
		return c.Text(http.StatusOK, "Hello, world!")
	})

	app.Use(func(next Handler) Handler {
		return HandlerFunc(func(c *Context) error {
			c.Response.Header().Add("X-XSS-Protection", "1; mode=blockFilter")

			return next.ServeHTTP(c)
		})
	})

	app.ServeHTTP(rr, req)

	header := rr.Header()

	assertEqual(t, http.StatusOK, rr.Code)
	assertEqual(t, "Hello, world!", rr.Body.String())
	assertEqual(t, "1; mode=blockFilter", header.Get("X-XSS-Protection"))
}

func TestAppMount(t *testing.T) {
	var tests = []struct {
		mountPath   string
		handlerPath string
		requestPath string
	}{
		{
			mountPath:   "/v1",
			handlerPath: "/",
			requestPath: "/v1",
		},
		{
			mountPath:   "/v1/",
			handlerPath: "/",
			requestPath: "/v1/",
		},
		{
			mountPath:   "/v1",
			handlerPath: "/a",
			requestPath: "/v1/a",
		},
	}

	for _, testcase := range tests {
		rr := httptest.NewRecorder()
		req, err := http.NewRequest("GET", testcase.requestPath, nil)
		if err != nil {
			t.Fatal(err)
		}

		group := New()
		group.Get(testcase.handlerPath, func(c *Context) error {
			return c.Text(http.StatusOK, "Hello, world!")
		})

		app := New()
		app.Mount(testcase.mountPath, group)
		app.ServeHTTP(rr, req)

		assertEqual(t, http.StatusOK, rr.Code)
		assertEqual(t, "Hello, world!", rr.Body.String())
	}
}

func GenerateCertificate(t *testing.T) {
	ca := &x509.Certificate{
		SerialNumber: big.NewInt(2020),
		Subject: pkix.Name{
			Organization: []string{"Lungo"},
			Country:      []string{"DE"},
			Province:     []string{"Saxony"},
			Locality:     []string{"Dresden"},
			PostalCode:   []string{"01069"},
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(1, 0, 0),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	key, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		t.Fatal(err)
	}

	b, err := x509.CreateCertificate(rand.Reader, ca, ca, &key.PublicKey, key)
	if err != nil {
		t.Fatal(err)
	}

	certOut, err := os.Create("cert.pem")
	if err != nil {
		t.Fatalf("Failed to open cert.pem for writing: %v", err)
	}
	if err := pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: b}); err != nil {
		t.Fatalf("Failed to write data to cert.pem: %v", err)
	}
	if err := certOut.Close(); err != nil {
		t.Fatalf("Error closing cert.pem: %v", err)
	}

	keyOut, err := os.OpenFile("key.pem", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		t.Fatalf("Failed to open key.pem for writing: %v", err)
	}
	if err := pem.Encode(keyOut, &pem.Block{Type: "PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}); err != nil {
		t.Fatalf("Failed to write data to key.pem: %v", err)
	}
	if err := keyOut.Close(); err != nil {
		t.Fatalf("Error closing key.pem: %v", err)
	}
}
