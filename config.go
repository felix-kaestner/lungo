package lungo

type Config struct {
	// Maximum number of  bytes to read from the request body.
	// A request body larger than that will result in returning
	// a "http: request body too large" error.
	//
	// Set this value to -1 to allow all arbitrary large request bodies.
	//
	// Default: 1 * 1024 * 1024 = 1048576 Bytes = 1MiB
	MaxBodySize int `json:"max_body_size"`
}

const (
	DefaultMaxBodySize = 1048576 // 1 * 1024 * 1024 = 1048576 Bytes = 1MiB
)
