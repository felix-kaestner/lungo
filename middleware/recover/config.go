package recover

type Config struct {
	// Handle defines a callback function to handle
	// the stack trace of the panic.
	//
	// Optional. Default: nil
	HandleStackTrace func(e interface{})
}

var DefaultConfig = &Config{
	HandleStackTrace: nil,
}
