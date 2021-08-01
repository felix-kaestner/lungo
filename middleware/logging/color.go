package logging

import "fmt"

type ColorCode int

const (
	Error   ColorCode = 31
	Success ColorCode = 32
	Warning ColorCode = 33
	Info    ColorCode = 34
	Debug   ColorCode = 35
	Notice  ColorCode = 36
	Reset   ColorCode = 0
)

func (c ColorCode) Sprintf(a interface{}) string {
	return fmt.Sprintf("\033[1;%dm%s\033[0m", c, a)
}
