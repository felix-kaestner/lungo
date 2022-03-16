package recover

// Config defines the configuration options for the recover middleware
type Config struct {
	// Handle defines a callback function to handle
	// the stack trace of the panic.
	//
	// Optional. Default: nil
	HandleStackTrace func(e interface{})
}

// DefaultConfig contains the default value for the
// recover middleware configuration
var DefaultConfig = &Config{
	HandleStackTrace: nil,
}
