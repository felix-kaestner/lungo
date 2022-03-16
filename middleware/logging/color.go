package logging

import "fmt"

// ColorCode defines identifiers to be used
// as ANSI color escape sequences
type ColorCode int

const (
	// Error is the ANSI color code for red
	Error ColorCode = 31
	// Success is the ANSI color code for green
	Success ColorCode = 32
	// Warning is the ANSI color code for orange
	Warning ColorCode = 33
	// Info is the ANSI color code for blue
	Info ColorCode = 34
	// Debug is the ANSI color code for magenta
	Debug ColorCode = 35
	// Notice is the ANSI color code for cyan
	Notice ColorCode = 36
	// Reset is the ANSI color code to reset the color
	Reset ColorCode = 0
)

// Sprintf colorizes the provided string by wrapping
// it inside an ANSI color escape sequence
func (c ColorCode) Sprintf(a interface{}) string {
	return fmt.Sprintf("\033[1;%dm%s\033[0m", c, a)
}
