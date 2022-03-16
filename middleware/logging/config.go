package logging

import (
	"log"
)

// Config defines the configuration options for the logging middleware
type Config struct {
	// Template defines the logging format as a text/template string.
	//
	// see: https://golang.org/pkg/text/template/
	//
	// Available Tags:
	// - "Request": *http.Request
	// - "Duration": *time.Duration
	//
	// Optional. Default:
	Template string

	// Logger defines the log.Logger instance to use for logging.
	// By default the standard logger is used, which makes it
	// equivalent to calling log.Printf().
	//
	// Optional. Default: log.Default()
	Logger *log.Logger

	// CallDepth defines the call stack depth to use for logging.
	//
	// Optional. Default: 2
	CallDepth int
}

// DefaultConfig contains the default value for the
// logging middleware configuration
var DefaultConfig = &Config{
	Template:  "{{.Request.Method}} {{.Request.URL.Path}} {{.Duration.String}}",
	Logger:    log.Default(),
	CallDepth: 2,
}
