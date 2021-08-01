package logging

import "testing"

func TestColor(t *testing.T) {
	assertEqual(t, "\033[1;32mFoo\033[0m", Success.Sprintf("Foo"))
}
